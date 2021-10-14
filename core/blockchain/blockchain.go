// Copyright 2021 The aereum Authors
// This file is part of the aereum library.
//
// The aereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The aereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package message contains data types related to aereum network.
package blockchain

import (
	"bytes"
	"sync"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
	"github.com/Aereum/aereum/core/wallet"
)

type Blockchain struct {
	Messages []message.Message
}

type AudienceState struct {
	Token     []byte
	Followers []*message.Follower
}

// State incorporates only the necessary information on the blockchain to
// validate new messages. It should be used on validation nodes.
type State struct {
	Epoch       uint64
	Subscribers wallet.HashVault // subscriber token hash
	Captions    wallet.HashVault // caption string hash
	Wallets     wallet.Wallet    // wallet token hash
	Audiences   wallet.Audience  // audience + Follower hash
	//AudienceRequests  map[crypto.Hash]*[]*message.Message // audience hash
	PowerOfAttorney   wallet.HashStore                 // power of attonery token hash
	AdvertisingOffers map[crypto.Hash]*message.Message // advertising offer hash
	//SyncJobs          []SyncJob
	*sync.Mutex
}

// validate message against current epoch state.
func (s *State) ValidateMessage(msg *message.Message) bool {
	pays := msg.Payments()
	if ok, balance := s.Wallets.Balance(pays.DebitAcc); !ok || pays.DebitValue > balance {
		return false
	}
	return true
}

type StateMutations struct {
	State        *State
	DeltaWallets map[crypto.Hash]int
	messages     []*message.Message
}

func (s *State) FromNewBlock() {
	block := StateMutations{
		State:        s,
		DeltaWallets: make(map[crypto.Hash]int),
		messages:     make([]*message.Message, 0),
	}
}

func (s *State) IncorporateMutations(mut StateMutations) {
	s.Epoch += 1
	for wallet, delta := range mut.DeltaWallets {
		if delta < 0 {
			s.Wallets.Debit(wallet, -uint64(delta))
		} else {
			s.Wallets.Credit(wallet, uint64(delta))
		}
	}

}

func (s *StateMutations) Withdraw(acc crypto.Hash, value int) bool {
	_, funds := s.State.Wallets.Balance(acc)
	delta := s.DeltaWallets[acc]
	if int(funds)+delta > value {
		s.DeltaWallets[acc] = delta - value
		return true
	}
	return false
}

func (s *StateMutations) Credit(acc crypto.Hash, value int) {
	delta := s.DeltaWallets[acc]
	s.DeltaWallets[acc] = delta + value
}

func (s *StateMutations) Transfer(t *message.Transfer) bool {
	hashFrom := crypto.Hasher(t.From)
	_, funds := s.State.Wallets.Balance(hashFrom)
	delta := s.DeltaWallets[hashFrom]
	value := int(t.Value)
	if int(funds)+delta < value {
		return false
	}
	hashTo := crypto.Hasher(t.To)
	deltaTo := s.DeltaWallets[hashTo]
	s.DeltaWallets[hashFrom] = delta - value
	s.DeltaWallets[hashTo] = deltaTo + value
	s.Transfers = append(s.Transfers, t)
	return true
}

func (s *StateMutations) RedistributeAdvertisemenetFee(value int, author crypto.Hash, audience []*message.Follower) {
	// 100% author provisory
	s.Credit(author, value)
}

func (s *State) CanSubscribe(m *message.Subscribe, author crypto.Hash) bool {
	caption := crypto.Hasher([]byte(m.Caption))
	if _, ok := s.State.Subscribers[author]; ok {
		return false
	}
	if _, ok := s.NewSubsribers[author]; ok {
		return false
	}
	if _, ok := s.State.Captions[caption]; ok {
		return false
	}
	if _, ok := s.NewCaptions[caption]; ok {
		return false
	}
	s.NewSubsribers[author] = struct{}{}
	s.NewCaptions[caption] = struct{}{}
	return true
}

func (s *State) CanPublish(m *message.Content) bool {
	ok, keys := s.Audiences.GetKeys(crypto.Hasher(m.Audience))
	if !ok {
		return false
	}
	submissionPub, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	if !submissionPub.VerifyHash(m.SubmitHash, m.SubmitSignature) {
		return false
	}
	if len(m.PublishSignature) > 0 {
		pulishPub, err := crypto.PublicKeyFromBytes(keys[crypto.PublicKeySize : 2*crypto.PublicKeySize])
		if err != nil {
			return false
		}
		if !pulishPub.VerifyHash(m.PublishHash, m.PublishSignature) {
			return false
		}
	}
	// does not check if the advertisement offer has resources in the walltet to
	// pay, only if the offer exists and the content matches
	if len(m.AdvertisingToken) > 0 {
		hash := crypto.Hasher(m.AdvertisingToken)
		if offerMsg, ok := s.AdvertisingOffers[hash]; ok {
			// check if advertising claim is valid
			offer := offerMsg.AsAdvertisingOffer()
			if !bytes.Equal(offer.Audience, m.Audience) {
				return false
			}
			if offer.ContentType != m.ContentType {
				return false
			}
			if !bytes.Equal(offer.ContentData, m.ContentData) {
				return false
			}
			return true
		} else {
			return false
		}
	}
	return true
}

func (s *State) CanCreateAudience(m *message.CreateAudience) bool {
	token := crypto.Hasher(m.Token)
	if !s.Audiences.Exists(token) {
		return false
	}
	return true
}
