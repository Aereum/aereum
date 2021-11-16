// Copyright 2021 The Aereum Authors
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
package instruction

import (
	"errors"
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
)

var ErrCouldNotParseMessage = errors.New("could not parse message")
var ErrCouldNotParseSignature = errors.New("could not parse signature")
var ErrInvalidSignature = errors.New("invalid signature")

const (
	GenesisMsg byte = iota
	// version 0
	TransferMsg
	SubscribeMsg
	AboutMsg
	CreateAudienceMsg
	JoinAudienceMsg
	AcceptJoinAudienceMsg
	AudienceChangeMsg
	AdvertisingOfferMsg
	ContentMsg
	GrantPowerOfAttorneyMsg
	RevokePowerOfAttorneyMsg
	UnkownMessageType // to be used in other versions
)

func NewMessage(AuthorKey, WalletKey crypto.PrivateKey, msg Serializer,
	FeeValue, Epoch uint64, PowerOfAttorney crypto.PrivateKey) *Message {
	message := &Message{
		MessageType: msg.Kind(),
		Author:      AuthorKey.PublicKey().ToBytes(),
		Message:     msg.Serialize(),
		FeeWallet:   WalletKey.PublicKey().ToBytes(),
		Epoch:       Epoch,
	}
	if PowerOfAttorney.IsValid() {
		message.PowerOfAttorney = PowerOfAttorney.PublicKey().ToBytes()
		message.Sign(PowerOfAttorney, WalletKey)
	} else {
		message.PowerOfAttorney = []byte{}
		message.Sign(AuthorKey, WalletKey)
	}
	return message
}

func MessageType(msg []byte) byte {
	if len(msg) == 0 {
		return UnkownMessageType
	}
	if msg[0] >= UnkownMessageType {
		return UnkownMessageType
	}
	return msg[0]
}

func IsMessage(msg []byte) bool {
	msgType := MessageType(msg)
	return msgType > TransferMsg && msgType < UnkownMessageType
}

type Instruction []byte

type Message struct {
	MessageType     byte
	Epoch           uint64
	Author          []byte // must be a subscriber public key, can be anonoymous on transfer
	Message         []byte
	FeeWallet       []byte // can be any wallet
	FeeValue        uint64
	PowerOfAttorney []byte // must be authorized by the subscriber
	Signature       []byte // either author or power of attorney
	WalletSignature []byte
}

type Payment struct {
	DebitAcc    []crypto.Hash
	DebitValue  []uint64
	CreditAcc   []crypto.Hash
	CreditValue []uint64
}

func (m *Message) Payments() Payment {
	return Payment{
		DebitAcc:   []crypto.Hash{crypto.Hasher(m.FeeWallet)},
		DebitValue: []uint64{m.FeeValue},
	}
}

func GetHashAndEpochFromMessage(msg []byte) (crypto.Hash, int) {
	epoch, _ := ParseUint64(msg, 1)
	return crypto.Hasher(msg), int(epoch)
}

func (m *Message) serializeWithoutSignatures() []byte {
	bytes := []byte{m.MessageType}
	PutUint64(m.Epoch, &bytes)
	PutByteArray(m.Author, &bytes)
	PutByteArray(m.Message, &bytes)
	PutByteArray(m.FeeWallet, &bytes)
	PutUint64(m.FeeValue, &bytes)
	PutByteArray(m.PowerOfAttorney, &bytes)
	return bytes
}

func (m *Message) Sign(author, wallet crypto.PrivateKey) bool {
	bytes := m.serializeWithoutSignatures()
	signAuthor, err := author.Sign(bytes)
	if err != nil {
		return false
	}
	PutByteArray(signAuthor, &bytes)
	signWallet, errWallet := wallet.Sign(bytes)
	if errWallet != nil {
		return false
	}
	m.Signature = signAuthor
	m.WalletSignature = signWallet
	return true
}

func (m *Message) Serialize() []byte {
	bytes := m.serializeWithoutSignatures()
	PutByteArray(m.Signature, &bytes)
	PutByteArray(m.WalletSignature, &bytes)
	return bytes
}

func ParseGenesis(data []byte) *Genesis {
	if data[0] != GenesisMsg {
		return nil
	}
	return &Genesis{}
}

func ParseMessage(data []byte) (*Message, error) {
	if data[0] >= UnkownMessageType || data[0] <= TransferMsg {
		return nil, fmt.Errorf("wrong message type")
	}
	length := len(data)
	var msg Message
	msg.MessageType = data[0]
	position := 1
	msg.Epoch, position = ParseUint64(data, position)
	msg.Author, position = ParseByteArray(data, position)
	msg.Message, position = ParseByteArray(data, position)
	msg.FeeWallet, position = ParseByteArray(data, position)
	msg.FeeValue, position = ParseUint64(data, position)
	msg.PowerOfAttorney, position = ParseByteArray(data, position)
	// check author or power of attorney signature
	if position-1 > length {
		return nil, ErrCouldNotParseMessage
	}
	msgToVerify := data[0:position]
	msg.Signature, position = ParseByteArray(data, position)
	token := msg.Author
	if len(msg.PowerOfAttorney) > 0 {
		token = msg.PowerOfAttorney
	}
	if publicKey, err := crypto.PublicKeyFromBytes(token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, msg.Signature) {
			return nil, ErrInvalidSignature
		}
	}

	// check wallet signature
	if position-1 > length {
		return nil, ErrCouldNotParseMessage
	}
	msgToVerify = data[0:position]
	msg.WalletSignature, position = ParseByteArray(data, position)
	if position != length {
		return nil, ErrCouldNotParseMessage
	}
	if publicKey, err := crypto.PublicKeyFromBytes(msg.FeeWallet); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, msg.WalletSignature) {
			return nil, ErrInvalidSignature
		}
	}
	return &msg, nil
}

func (m *Message) AsSubscribe() *Subscribe {
	return ParseSubscribe(m.Message)
}

func (m *Message) AsAbout() *About {
	return ParseAbout(m.Message)
}

func (m *Message) AsCreateAudiece() *CreateAudience {
	return ParseCreateAudience(m.Message)
}

func (m *Message) AsJoinAudience() *JoinAudience {
	return ParseJoinAudience(m.Message)
}

func (m *Message) AsAcceptJoinAudience() *AcceptJoinAudience {
	return ParseAcceptJoinAudience(m.Message)
}

func (m *Message) AsChangeAudience() *ChangeAudience {
	return ParseChangeAudience(m.Message)
}

func (m *Message) AsAdvertisingOffer() *AdvertisingOffer {
	return ParseAdvertisingOffer(m.Message)
}

func (m *Message) AsContent() *Content {
	return ParseContent(m.Message)
}

func (m *Message) AsGrantPowerOfAttorney() *GrantPowerOfAttorney {
	return ParseGrantPowerOfAttorney(m.Message)
}

func (m *Message) AsRevokePowerOfAttorney() *RevokePowerOfAttorney {
	return ParseRevokePowerOfAttorney(m.Message)
}
