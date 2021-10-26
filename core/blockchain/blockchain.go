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
	"github.com/Aereum/aereum/core/message"
)

type Block struct {
	Parent       crypto.Hash
	Publisher    []byte
	PublishedAt  time.Time
	Messages     [][]byte
	Transactions [][]byte
	Hash         crypto.Hash
}

func (b *Block) SerializeWithoutHash() []byte {
	serialized := b.Parent[:]
	message.PutByteArray(b.Publisher, &serialized)
	message.PutUint64(uint64(b.PublishedAt.UnixNano()), &serialized)
	message.PutUint64(uint64(len(b.Messages)), &serialized)
	for _, msg := range b.Messages {
		message.PutByteArray(msg, &serialized)
	}
	message.PutUint64(uint64(b.PublishedAt.UnixNano()), &serialized)
	message.PutUint64(uint64(len(b.Transactions)), &serialized)
	for _, transaction := range b.Transactions {
		message.PutByteArray(transaction, &serialized)
	}
	return serialized
}

func (b *Block) Serialize() ([]byte, crypto.Hash) {
	serialized := b.SerializeWithoutHash()
	hash := crypto.Hasher(serialized)
	return append(serialized[0:crypto.Size], hash[:]...), hash
}

func ParseBlock(data []byte) *Block {
	block := &Block{}
	block.Parent = crypto.BytesToHash(data[0:crypto.Size])
	position := crypto.Size
	block.Publisher, position = message.ParseByteArray(data, position)
	block.PublishedAt, position = message.ParseTime(data, position)
	var count uint64
	count, position = message.ParseUint64(data, position)
	block.Messages = make([][]byte, int(count))
	for n := 0; n < int(count); n++ {
		block.Messages[n], position = message.ParseByteArray(data, position)
	}
	count, position = message.ParseUint64(data, position)
	block.Transactions = make([][]byte, int(count))
	for n := 0; n < int(count); n++ {
		block.Transactions[n], position = message.ParseByteArray(data, position)
	}
	if len(data)-position != crypto.Size {
		return nil
	}
	block.Hash = crypto.BytesToHash(data[position:])
	return block
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
