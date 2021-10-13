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
	"github.com/Aereum/aereum/core/hashdb"
	"github.com/Aereum/aereum/core/message"
	"github.com/Aereum/aereum/core/network"
	"github.com/Aereum/aereum/core/wallet"
)

type Blockchain struct {
	Messages []message.Message
}

type AudienceState struct {
	Token     []byte
	Followers []*message.Follower
}

type SyncJob struct {
	EpochStart uint64
	socket     *network.SyncSocket
}

// State incorporates only the necessary information on the blockchain to
// validate new messages. It should be used on validation nodes.
type State struct {
	Epoch             uint64
	Subscribers       hashdb.HashStore                    // subscriber token hash
	Captions          hashdb.HashStore                    // caption string hash
	Wallets           wallet.Wallet                       // wallet token hash
	Audiences         hashdb.HashStore                    // audience + Follower hash
	AudienceRequests  map[crypto.Hash]*[]*message.Message // audience hash
	PowerOfAttorney   hashdb.HashStore                    // power of attonery token hash
	AdvertisingOffers map[hashdb.Hash]*message.Message    // advertising offer hash
	SyncJobs          []SyncJob
	*sync.Mutex
}

func (s *State) IncorporateMutations(mut NewBlockMuttations) {
	s.Epoch += 1
	for subscriber, _ := range mut.NewSubsribers {
		s.Subscribers.InsertIfNotExists(subscriber)
		// TODO: Treat Error
	}
	for caption, _ := range mut.NewCaptions {
		s.Captions.InsertIfNotExists(caption)
		// TODO: Treat Error
	}
	for wallet, delta := range mut.DeltaWallets {
		if delta < 0 {
			s.Wallets.Debit(wallet, -uint64(delta))
		} else {
			s.Wallets.Credit(wallet, uint64(delta))
		}
	}

}

type NewBlockMuttations struct {
	State                     *State
	NewSubsribers             map[crypto.Hash]struct{}
	NewCaptions               map[crypto.Hash]struct{}
	DeltaWallets              map[crypto.Hash]int
	NewAudiences              map[crypto.Hash]*[]*message.Follower
	ChangeAudicences          map[crypto.Hash]*[]*message.Follower
	NewAudinceRequests        map[crypto.Hash]*message.Message
	GrantPowerOfAttorney      map[crypto.Hash]crypto.Hash
	RevokePowerOfAttorney     map[crypto.Hash]struct{}
	NewAdvertisingOffers      map[crypto.Hash]*message.Message
	AcceptedAdvertisingOffers map[crypto.Hash]*message.Message
	Messages                  []*message.Message
	Transfers                 []*message.Transfer
}

func (s *NewBlockMuttations) Withdraw(acc Hash, value int) bool {
	funds := s.State.Wallets[acc]
	delta := s.DeltaWallets[acc]
	if funds+delta > value {
		s.DeltaWallets[acc] = delta - value
		return true
	}
	return false
}

func (s *NewBlockMuttations) Credit(acc Hash, value int) {
	delta := s.DeltaWallets[acc]
	s.DeltaWallets[acc] = delta + value
}

func (s *NewBlockMuttations) Transfer(t *message.Transfer) bool {
	hashFrom := Hash256(t.From)
	funds := s.State.Wallets[hashFrom]
	delta := s.DeltaWallets[hashFrom]
	value := int(t.Value)
	if funds+delta < value {
		return false
	}
	hashTo := Hash256(t.To)
	deltaTo := s.DeltaWallets[hashTo]
	s.DeltaWallets[hashFrom] = delta - value
	s.DeltaWallets[hashTo] = deltaTo + value
	s.Transfers = append(s.Transfers, t)
	return true
}

func (s *NewBlockMuttations) RedistributeAdvertisemenetFee(value int, author crypto.Hash, audience []*message.Follower) {
	// 100% author provisory
	s.Credit(author, value)
}

func (s *NewBlockMuttations) CanContent(m *message.Content, author, wallet crypto.Hash, fee int) bool {
	if len(m.AdvertisingToken) > 0 {
		hash := crypto.Hasher(m.AdvertisingToken)
		if offerMsg, ok := s.State.AdvertisingOffers[hash]; ok {
			if _, ok := s.AcceptedAdvertisingOffers[hash]; ok {
				return false // message already reclaimed in the new block
			}
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
			// check if advertiser wallet have funds to pay
			value := int(offer.AdvertisingFee)
			if !s.Withdraw(Hash256(offerMsg.FeeWallet), value) {
				return false
			}
			// use protocol redistribution rule
			s.RedistributeAdvertisemenetFee(value, author, nil)
			// mark offer as accepted
			s.AcceptedAdvertisingOffers[hash] = offerMsg
			return true
		} else {
			return false
		}
	}
	return true
}

func (s *NewBlockMuttations) CanSubscribe(m *message.Subscribe, author Hash) bool {
	caption := Hash256([]byte(m.Caption))
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

func (s *NewBlockMuttations) CanCreateAudience(m *message.CreateAudience) bool {
	audience := crypto.Hasher(m.Token)

	if _, ok := s.State.Audiences[audience]; ok {
		return false
	}
	if _, ok := s.NewAudieces[audience]; ok {
		return false
	}

}
