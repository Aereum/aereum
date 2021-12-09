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
	"github.com/Aereum/aereum/core/util"
)

type Recipient struct {
	Token []byte
	Value uint64
}

func NewSingleReciepientTransfer(from crypto.PrivateKey, to []byte, reason string, value, epoch, fee uint64) *Transfer {
	transfer := &Transfer{
		From:   from.PublicKey().ToBytes(),
		To:     []Recipient{{to, value}},
		Reason: reason,
		epoch:  epoch,
		Fee:    fee,
	}
	bytes := crypto.Hasher(transfer.serializeWithoutSignature())
	transfer.Signature, _ = from.Sign(bytes[:])
	return transfer
}

func NewDeposit(from crypto.PrivateKey, value, epoch, fee uint64) *Deposit {
	deposit := &Deposit{
		Token: from.PublicKey().ToBytes(),
		epoch: epoch,
		Fee:   fee,
	}
	bytes := crypto.Hasher(deposit.serializeWithoutSignature())
	deposit.Signature, _ = from.Sign(bytes[:])
	return deposit
}

func NewWithdraw(from crypto.PrivateKey, value, epoch, fee uint64) *Withdraw {
	deposit := &Withdraw{
		Token: from.PublicKey().ToBytes(),
		epoch: epoch,
		Fee:   fee,
	}
	bytes := crypto.Hasher(deposit.serializeWithoutSignature())
	deposit.Signature, _ = from.Sign(bytes[:])
	return deposit
}

// Transfer aero from a wallet to a series of other wallets
type Transfer struct {
	epoch     uint64
	From      []byte
	To        []Recipient
	Reason    string
	Fee       uint64
	Signature []byte
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

func (t *Transfer) Validate(v InstructionValidator) bool {
	v.AddFeeCollected(t.Fee)
	return true
}

func (a *Transfer) Kind() byte {
	return iTransfer
}

func (a *Transfer) Epoch() uint64 {
	return a.epoch
}

func (s *Transfer) serializeWithoutSignature() []byte {
	bytes := []byte{0, iTransfer}
	util.PutUint64(s.epoch, &bytes)
	util.PutByteArray(s.From, &bytes)
	util.PutUint16(uint16(len(s.To)), &bytes)
	count := len(s.To)
	if len(s.To) > 1<<16-1 {
		count = 1 << 16
	}
	for n := 0; n < count; n++ {
		util.PutByteArray(s.To[n].Token, &bytes)
		util.PutUint64(s.To[n].Value, &bytes)
	}
	util.PutString(s.Reason, &bytes)
	util.PutUint64(s.Fee, &bytes)
	return bytes
}

func (s *Transfer) Serialize() []byte {
	bytes := s.serializeWithoutSignature()
	util.PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseTransfer(data []byte) *Transfer {
	if len(data) < 2 || data[1] != iTransfer {
		return nil
	}
	p := Transfer{}
	position := 2
	p.epoch, position = util.ParseUint64(data, position)
	p.From, position = util.ParseByteArray(data, position)
	var count uint16
	count, position = util.ParseUint16(data, position)
	p.To = make([]Recipient, int(count))
	for i := 0; i < int(count); i++ {
		p.To[i].Token, position = util.ParseByteArray(data, position)
		p.To[i].Value, position = util.ParseUint64(data, position)
	}
	p.Reason, position = util.ParseString(data, position)
	p.Fee, position = util.ParseUint64(data, position)
	hash := crypto.Hasher(data[0:position])
	p.Signature, _ = util.ParseByteArray(data, position)
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
	epoch     uint64
	Token     []byte
	Value     uint64
	Fee       uint64
	Signature []byte
}

func (d *Deposit) Payments() *Payment {
	return NewPayment(crypto.Hasher(d.Token), d.Value)
}

func (t *Deposit) Validate(v InstructionValidator) bool {
	v.AddFeeCollected(t.Fee)
	return true
}

func (a *Deposit) Kind() byte {
	return iDeposit
}

func (a *Deposit) Epoch() uint64 {
	return a.epoch
}

func (s *Deposit) serializeWithoutSignature() []byte {
	bytes := []byte{0, iDeposit}
	util.PutUint64(s.epoch, &bytes)
	util.PutByteArray(s.Token, &bytes)
	util.PutUint64(s.Value, &bytes)
	util.PutUint64(s.Fee, &bytes)
	return bytes
}

func (s *Deposit) Serialize() []byte {
	bytes := s.serializeWithoutSignature()
	util.PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseDeposit(data []byte) *Deposit {
	if len(data) < 2 || data[1] != iDeposit {
		return nil
	}
	p := Deposit{}
	position := 2
	p.epoch, position = util.ParseUint64(data, position)
	p.Token, position = util.ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	}
	p.Value, position = util.ParseUint64(data, position)
	p.Fee, position = util.ParseUint64(data, position)
	msgToVerify := crypto.Hasher(data[0:position])
	p.Signature, _ = util.ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	} else {
		if !publicKey.Verify(msgToVerify[:], p.Signature) {
			return nil
		}
	}
	return &p
}

// Withdraw aero from a wallet
type Withdraw struct {
	epoch     uint64
	Token     []byte
	Value     uint64
	Fee       uint64
	Signature []byte
}

func (w *Withdraw) Payments() *Payment {
	return &Payment{
		Credit: []Wallet{{crypto.Hasher(w.Token), w.Value}},
		Debit:  []Wallet{},
	}
}

func (t *Withdraw) Validate(v InstructionValidator) bool {
	v.AddFeeCollected(t.Fee)
	return true
}

func (a *Withdraw) Kind() byte {
	return iWithdraw
}

func (a *Withdraw) Epoch() uint64 {
	return a.epoch
}

func (s *Withdraw) serializeWithoutSignature() []byte {
	bytes := []byte{0, iWithdraw}
	util.PutUint64(s.epoch, &bytes)
	util.PutByteArray(s.Token, &bytes)
	util.PutUint64(s.Value, &bytes)
	util.PutUint64(s.Fee, &bytes)
	return bytes
}

func (s *Withdraw) Serialize() []byte {
	bytes := s.serializeWithoutSignature()
	util.PutByteArray(s.Signature, &bytes)
	return bytes
}

func ParseWithdraw(data []byte) *Withdraw {
	if len(data) < 2 || data[1] != iWithdraw {
		return nil
	}
	p := Withdraw{}
	position := 2
	p.epoch, position = util.ParseUint64(data, position)
	p.Token, position = util.ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	}
	p.Value, position = util.ParseUint64(data, position)
	p.Fee, position = util.ParseUint64(data, position)
	msgToVerify := crypto.Hasher(data[0:position])
	p.Signature, _ = util.ParseByteArray(data, position)
	if publicKey, err := crypto.PublicKeyFromBytes(p.Token); err != nil {
		return nil
	} else {
		if !publicKey.Verify(msgToVerify[:], p.Signature) {
			return nil
		}
	}
	return &p
}
