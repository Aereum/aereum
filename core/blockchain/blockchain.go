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
	transfers    []*message.Transfer
}

func (s *StateMutations) Validate(msg []byte) bool {
	if message.IsTransfer(msg) {
		transfer, _ := message.ParseTranfer(msg)
		if transfer != nil {
			return false
		}
		if s.Debit(crypto.Hasher(transfer.From), transfer.Value + transfer.Fee) {
			s.Credit(transfer.To, int(transfer.Value))
			s.transfers = append(s.transfers, transfer)
			return true
		}
	} 
	if !message.IsMessage(msg) {
		return false
	}
	message, err := message.ParseMessage(msg)
	if message == nil || err != nil {
		return false
	}
	payments := message.Payments()
	if !s.CanPay(payments) {
		return false
	}
	switch message.MessageType(msg) {
	case SubscribeMsg:
		subscribe := message.AsSubscribe()
		if subscribe == nil {
			return false
		} 
		//
	case AboutMsg:
		about := message.AsAbout()
		if about == nil {
			return false
		}
		//
	case CreateAudienceMsg:
		createAudience := message.AsCreateAudiece()
		if createAudience == nil {
			return false
		}

	case JoinAudienceMsg:
		joinAudience := message.AsJoinAudience()
		if joinAudience == nil {
			return false
		}

	case AudienceChangeMsg:
		audienceChange := message.AsChangeAudience()
		if audienceChange == nil {
			return false
		}
	
	case AdvertisingOfferMsg:
		advertisingOffer := message.AsAdvertisingOffer()
		if advertisingOffer == nil {
			return false
		}

	case ContentMsg:
		content := message.AsContent()
		if content == nil {
			return false
		}

	case GrantPowerOfAttorneyMsg:
		grantPower := message.AsGrantPowerOfAttorney()
		if grantPower == nil {
			return false
		}

	case RevokePowerOfAttorneyMsg:
		revokePower := message.AsRevokePowerOfAttorney()
		if revokePower == nil {
			return false
		}	
	}
	for n, acc := range payments.DebitAcc {
		s.Debit(acc, int(payments.DebitValue[n]))
	}
	for n, acc := range payments.CreditAcc {
		s.Credit(acc, int(payments.CreditValue[n]))
	}
	return true
}

func (s *State) FromNewBlock() {
	block := StateMutations{
		State:        s,
		DeltaWallets: make(map[crypto.Hash]int),
		messages:     make([]*message.Message, 0),
	}
	validator := make(chan []byte)
	sealBlock := make(chan struct{}) 
	go func () {
		for {
			select{
			case msg := <-validator:
				//
			case <-sealBlock:
				//
			}
		}
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

func (s *StateMutations) Debit(acc crypto.Hash, value int) bool {
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

func (s *StateMutations) CanPay(payment message.Payment) bool {
	for n, acc := range payment.DebitAcc {
		_, funds := s.State.Wallets.Balance(acc)
		delta := s.DeltaWallets[payment.DebitValue[n]]
		if int(funds)+delta < value {
			return false
		}
	}
	return true
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
