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
	"time"

	"github.com/Aereum/aereum/core/crypto"
)

type Block struct {
	Parent       crypto.Hash
	Hash         crypto.Hash
	Publisher    []byte
	PublishedAt  time.Time
	Messages     [][]byte
	Transactions [][]byte
}

func (b *Block) Serialize() []byte {
	return nil
}

func (s *State) FreezeBlock() {

}

func (s *State) SealBlock() {
	for hash, delta := range s.Mutations.DeltaWallets {
		if delta > 0 {
			s.Wallets.Credit(hash, uint64(delta))
		} else if delta < 0 {
			s.Wallets.Debit(hash, uint64(-delta))
		}
	}
	for hash := range s.Mutations.GrantPower {
		s.PowerOfAttorney.Insert(hash)
	}
	for hash := range s.Mutations.RevokePower {
		s.PowerOfAttorney.Remove(hash)
	}
	for hash := range s.Mutations.NewSubscriber {
		s.Subscribers.Insert(hash)
	}
	for hash := range s.Mutations.NewCaption {
		s.Captions.Insert(hash)
	}
	for audience, keys := range s.Mutations.NewAudiences {
		s.Audiences.SetKeys(audience, keys)
	}
}

/*func (s *State) IncorporateMutations(mut StateMutations) {
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
*/
