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
package message

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

const (
	GenesisMsg byte = iota
	// version 0
	TransferMsg
	SubscribeMsg
	AboutMsg
	CreateAudienceMsg
	JoinAudienceMsg
	AudienceChangeMsg
	AdvertisingOfferMsg
	ContentMsg
	GrantPowerOfAttorneyMsg
	RevokePowerOfAttorneyMsg
	UnkownMessageType
	// to be used in other versions
)

type Genesis struct {
}

type Transfer struct {
	From      []byte
	To        []byte
	Value     uint64
	Epoch     uint64
	Signature []byte
}

func (t *Transfer) serializeWithouSignature() []byte {
	bytes := []byte{TransferMsg}
	PutByteArray(t.From, &bytes)
	PutByteArray(t.To, &bytes)
	PutUint64(t.Value, &bytes)
	PutUint64(t.Epoch, &bytes)
	return bytes
}

func (t *Transfer) Sign(privateKey *rsa.PrivateKey) bool {
	bytes := t.serializeWithouSignature()
	hashed := sha256.Sum256(bytes)
	sign, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return false
	}
	t.Signature = sign
	return true
}

func (t *Transfer) Serialize() []byte {
	bytes := t.serializeWithouSignature()
	PutByteArray(t.Signature, &bytes)
	return bytes
}

func ParseTranfer(data []byte) (*Transfer, error) {
	if len(data) == 0 || data[0] != TransferMsg {
		return nil, fmt.Errorf("wrong message type")
	}
	length := len(data)
	var msg Transfer
	position := 1
	msg.From, position = ParseByteArray(data, position)
	msg.To, position = ParseByteArray(data, position)
	msg.Value, position = ParseUint64(data, position)
	msg.Epoch, position = ParseUint64(data, position)
	if position >= length {
		return nil, fmt.Errorf("could not parse message")
	}
	hashed := sha256.Sum256(data[0:position])
	msg.Signature, position = ParseByteArray(data, position)
	if position-1 > length || len(msg.Signature) == 0 {
		return nil, fmt.Errorf("could not parse message")
	}
	// check signature
	if publicKey, err := x509.ParsePKCS1PublicKey(msg.From); err != nil {
		return nil, fmt.Errorf("could not parse signature")
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.Signature) != nil {
			return nil, fmt.Errorf("invalid signature")
		}
	}
	return &msg, nil
}

func NewMessage(AuthorKey, WalletKey *rsa.PrivateKey, msg Serializer,
	FeeValue, Epoch uint64, PowerOfAttorney *rsa.PrivateKey) *Message {
	message := &Message{
		MessageType: msg.Kind(),
		Author:      x509.MarshalPKCS1PublicKey(&AuthorKey.PublicKey),
		Message:     msg.Serialize(),
		FeeWallet:   x509.MarshalPKCS1PublicKey(&WalletKey.PublicKey),
		Epoch:       Epoch,
	}
	if PowerOfAttorney != nil {
		message.PowerOfAttorney = x509.MarshalPKCS1PublicKey(&PowerOfAttorney.PublicKey)
		message.Sign(PowerOfAttorney, WalletKey)
	} else {
		message.PowerOfAttorney = []byte{}
		message.Sign(AuthorKey, WalletKey)
	}
	return message
}

type Message struct {
	MessageType     byte
	Author          []byte
	Message         []byte
	FeeWallet       []byte
	FeeValue        uint64
	Epoch           uint64
	PowerOfAttorney []byte
	Signature       []byte
	WalletSignature []byte
}

func (m *Message) serializeWithoutSignatures() []byte {
	bytes := []byte{m.MessageType}
	PutByteArray(m.Author, &bytes)
	PutByteArray(m.Message, &bytes)
	PutByteArray(m.FeeWallet, &bytes)
	PutUint64(m.FeeValue, &bytes)
	PutUint64(m.Epoch, &bytes)
	PutByteArray(m.PowerOfAttorney, &bytes)
	return bytes
}

func (m *Message) Sign(author, wallet *rsa.PrivateKey) bool {
	bytes := m.serializeWithoutSignatures()
	hashed := sha256.Sum256(bytes)
	signAuthor, err := rsa.SignPKCS1v15(nil, author, crypto.SHA256, hashed[:])
	if err != nil {
		return false
	}
	PutByteArray(signAuthor, &bytes)
	hashed = sha256.Sum256(bytes)
	signWallet, errWallet := rsa.SignPKCS1v15(nil, wallet, crypto.SHA256, hashed[:])
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
	msg.Author, position = ParseByteArray(data, position)
	msg.Message, position = ParseByteArray(data, position)
	msg.FeeWallet, position = ParseByteArray(data, position)
	msg.FeeValue, position = ParseUint64(data, position)
	msg.Epoch, position = ParseUint64(data, position)
	msg.PowerOfAttorney, position = ParseByteArray(data, position)
	// check author or power of attorney signature
	if position-1 > length {
		return nil, fmt.Errorf("could not parse message")
	}
	hashed := sha256.Sum256(data[0:position])
	msg.Signature, position = ParseByteArray(data, position)
	token := msg.Author
	if len(msg.PowerOfAttorney) > 0 {
		token = msg.PowerOfAttorney
	}
	if publicKey, err := x509.ParsePKCS1PublicKey(token); err != nil {
		return nil, fmt.Errorf("could not parse author key")
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.Signature) != nil {
			return nil, fmt.Errorf("invalid author signature")
		}
	}

	// check wallet signature
	if position-1 > length {
		return nil, fmt.Errorf("could not parse message")
	}
	hashed = sha256.Sum256(data[0:position])
	msg.WalletSignature, position = ParseByteArray(data, position)
	if position != length {
		return nil, fmt.Errorf("could not parse message")
	}
	if publicKey, err := x509.ParsePKCS1PublicKey(msg.FeeWallet); err != nil {
		return nil, fmt.Errorf("could not parse wallet key")
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.WalletSignature) != nil {
			return nil, fmt.Errorf("invalid wallet signature")
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
