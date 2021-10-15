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

