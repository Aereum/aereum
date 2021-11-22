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

package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
)

type Recipient struct {
	Token []byte
	Value uint64
}

// Transfer aero from a wallet to a series of other wallets
type Transfer struct {
	Version         byte
	InstructionType byte
	epoch           uint64
	From            []byte
	To              []Recipient
	Reason          string
	Fee             uint64
	Signature       []byte
}

func (a *Transfer) Kind() byte {
	return a.InstructionType
}

func (a *Transfer) Epoch() uint64 {
	return a.epoch
}

func (s *Transfer) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByte(s.Version, &bytes)
	PutByte(s.InstructionType, &bytes)
	PutUint64(s.epoch, &bytes)
	PutByteArray(s.From, &bytes)
	bytes = append(bytes, byte(len(s.To)))
	for _, v := range s.To {
		PutByteArray(v.Token, &bytes)
		PutUint64(v.Value, &bytes)
	}
	PutString(s.Reason, &bytes)
	PutUint64(s.Fee, &bytes)
	PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseTransfer(data []byte) (*Transfer, error) {
	p := Transfer{}
	position := 0
	counter := uint64(0)
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.From, position = ParseByteArray(data, position)

	counter, position = ParseUint64(data, position) // Precisa corrigir essa funcao de parse pq eh um int nao um uint
	recipients := make([]Recipient, 0)
	for i := 0; i < int(counter); i++ {
		t, position := ParseByteArray(data, position)
		v, position := ParseUint64(data, position)
		n := Recipient{Token: t, Value: v}
		recipients = append(recipients, n)
	}
	p.To = recipients
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)

	msgToVerify := data[0:position]
	p.Signature, position = ParseByteArray(data, position)
	token := p.From
	if publicKey, err := crypto.PublicKeyFromBytes(token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, p.Signature) {
			return nil, ErrInvalidSignature
		}
	}
	return &p, nil
}

// Deposit aero in a wallet
type Deposit struct {
	Version         byte
	InstructionType byte
	epoch           uint64
	Token           []byte
	Value           uint64
	Reason          string
	Fee             uint64
	Signature       []byte
}

func (a *Deposit) Kind() byte {
	return a.InstructionType
}

func (a *Deposit) Epoch() uint64 {
	return a.epoch
}

func (s *Deposit) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByte(s.Version, &bytes)
	PutByte(s.InstructionType, &bytes)
	PutUint64(s.epoch, &bytes)
	PutByteArray(s.Token, &bytes)
	PutUint64(s.Value, &bytes)
	PutString(s.Reason, &bytes)
	PutUint64(s.Fee, &bytes)
	PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseDeposit(data []byte) (*Deposit, error) {
	p := Deposit{}
	position := 0
	counter := uint64(0)
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.Token, position = ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
	p.Value = ParseUint64(data, position)
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)
	msgToVerify := data[0:position]
	p.Signature, position = ParseByteArray(data, position)
	token := p.From
	if publicKey, err := crypto.PublicKeyFromBytes(token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, p.Signature) {
			return nil, ErrInvalidSignature
		}
	}
	return &p, nil
}


// Withdraw aero from a wallet
type Withdraw struct {
	Version         byte
	InstructionType byte
	epoch           uint64
	Token           []byte
	Value           uint64
	Reason          string
	Fee             uint64
	Signature       []byte
}

func (a *Withdraw) Kind() byte {
	return a.InstructionType
}

func (a *Withdraw) Epoch() uint64 {
	return a.epoch
}

func (s *Withdraw) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByte(s.Version, &bytes)
	PutByte(s.InstructionType, &bytes)
	PutUint64(s.epoch, &bytes)
	PutByteArray(s.Token, &bytes)
	PutUint64(s.Value, &bytes)
	PutString(s.Reason, &bytes)
	PutUint64(s.Fee, &bytes)
	PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseWithdraw(data []byte) (*Withdraw, error) {
	p := Withdraw{}
	position := 0
	counter := uint64(0)
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.Token, position = ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
	p.Value = ParseUint64(data, position)
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)
	msgToVerify := data[0:position]
	p.Signature, position = ParseByteArray(data, position)
	token := p.From
	if publicKey, err := crypto.PublicKeyFromBytes(token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, p.Signature) {
			return nil, ErrInvalidSignature
		}
	}
	return &p, nil
}
