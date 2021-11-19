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
//
package instructions

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)

// Transfer aero from a wallet to another wallet
type Transfer struct {
	To		[]byte
	From	[]byte
	Value	uint64
  }

func (s *Transfer) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.To, &bytes)
	PutByteArray(s.From, &bytes)
	PutUint64(s.Value, &bytes)
	return bytes
}

func ParseTransfer(data []byte) *Transfer {
	p := Transfer{}
	position := 0
	p.To, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.To); err != nil {
		return nil
	}
	p.From, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.From); err != nil {
		return nil
	}
	p.Value, position = ParseUint64(data, position)
	if position == len(data) {
        return &p
    }
    return nil
}

// Deposit aero in a wallet
type Deposit struct {
	To		[]byte
	Value	uint64
}

func (s *Deposit) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.To, &bytes)
	PutUint64(s.Value, &bytes)
	return bytes
}

func ParseDeposit(data []byte) *Deposit {
	p := Deposit{}
	position := 0
	p.To, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.To); err != nil {
		return nil
	}
	p.Value, position = ParseUint64(data, position)
	if position == len(data) {
        return &p
    }
    return nil
}

// Withdraw aero from a wallet
type Withdraw struct {
	From	[]byte
	Value	uint64
}

func (s *Withdraw) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.From, &bytes)
	PutUint64(s.Value, &bytes)
	return bytes
}

func ParseWithdraw(data []byte) *Withdraw {
	p := Withdraw{}
	position := 0
	p.From, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.From); err != nil {
		return nil
	}
	p.Value, position = ParseUint64(data, position)
	if position == len(data) {
        return &p
    }
    return nil
}