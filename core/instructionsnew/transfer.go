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

package instructionsnew

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

func (t *Transfer) Payments() *Payment {
	total := uint64(0)
	payment := &Payment{
		Credit: make([]Wallet, 0),
		Debit:  make([]Wallet, 0),
	}
	for _, credit := range t.To {
		payment.NewCredit(crypto.Hasher(credit.Token), credit.Value)
		total += credit.Value
	}
	payment.NewDebit(crypto.Hasher(t.From), total+t.Fee)
	return payment
}

func (t *Transfer) Validate(block *Block) bool {
	return true
}

func (a *Transfer) Kind() byte {
	return a.InstructionType
}

func (a *Transfer) Epoch() uint64 {
	return a.epoch
}

func (s *Transfer) serializeWithoutSignature() []byte {
	bytes := make([]byte, 0)
	PutByte(s.Version, &bytes)
	PutByte(s.InstructionType, &bytes)
	PutUint64(s.epoch, &bytes)
	PutByteArray(s.From, &bytes)
	PutUint16(uint16(len(s.To)), &bytes)
	count := len(s.To)
	if len(s.To) > 1<<16-1 {
		count = 1 << 16
	}
	for n := 0; n < count; n++ {
		PutByteArray(s.To[n].Token, &bytes)
		PutUint64(s.To[n].Value, &bytes)
	}
	PutString(s.Reason, &bytes)
	PutUint64(s.Fee, &bytes)
	return bytes
}

func (s *Transfer) Serialize() []byte {
	bytes := s.serializeWithoutSignature()
	PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseTransfer(data []byte) *Transfer {
	p := Transfer{}
	position := 0
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.From, position = ParseByteArray(data, position)
	var count uint16
	count, position = ParseUint16(data, position)
	p.To = make([]Recipient, int(count))
	for i := 0; i < int(count); i++ {
		p.To[i].Token, position = ParseByteArray(data, position)
		p.To[i].Value, position = ParseUint64(data, position)
	}
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)
	hash := crypto.Hasher(data[0:position])
	p.Signature, _ = ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.From); err != nil {
		return nil
	} else {
		if !publicKey.Verify(hash[:], p.Signature) {
			return nil
		}
	}
	return &p
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

func (d *Deposit) Payments() *Payment {
	return NewPayment(crypto.Hasher(d.Token), d.Value)
}

func (t *Deposit) Validate(block *Block) bool {
	return true
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

func ParseDeposit(data []byte) *Deposit {
	p := Deposit{}
	position := 0
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.Token, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	}
	p.Value, position = ParseUint64(data, position)
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)
	msgToVerify := data[0:position]
	p.Signature, _ = ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	} else {
		if !publicKey.Verify(msgToVerify, p.Signature) {
			return nil
		}
	}
	return &p
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

func (w *Withdraw) Payments() *Payment {
	return &Payment{
		Credit: []Wallet{{crypto.Hasher(w.Token), w.Value}},
		Debit:  []Wallet{},
	}
}

func (t *Withdraw) Validate(block *Block) bool {
	return true
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

func ParseWithdraw(data []byte) *Withdraw {
	p := Withdraw{}
	position := 0
	p.Version, position = ParseByte(data, position)
	p.InstructionType, position = ParseByte(data, position)
	p.epoch, position = ParseUint64(data, position)
	p.Token, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	}
	p.Value, position = ParseUint64(data, position)
	p.Reason, position = ParseString(data, position)
	p.Fee, position = ParseUint64(data, position)
	msgToVerify := data[0:position]
	p.Signature, _ = ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	} else {
		if !publicKey.Verify(msgToVerify, p.Signature) {
			return nil
		}
	}
	return &p
}
