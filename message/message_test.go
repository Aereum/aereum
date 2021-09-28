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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"reflect"
	"testing"
)

func TestTransfer(t *testing.T) {
	FromPrivate, _ := rsa.GenerateKey(rand.Reader, 512)
	FromPublic := x509.MarshalPKCS1PublicKey(&FromPrivate.PublicKey)
	transfer := &Transfer{
		From:  FromPublic,
		To:    []byte{1, 2, 3, 4},
		Value: 12,
		Epoch: 1265,
	}
	if ok := transfer.Sign(FromPrivate); !ok {
		t.Error("could not sign transfer")
	}
	bytes := transfer.Serialize()
	copy, err := ParseTranfer(bytes)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if ok := reflect.DeepEqual(*transfer, *copy); !ok {
		t.Error("Parse and Serialization not working for Transfer")
	}
}

func TestMessage(t *testing.T) {
	AuthorPrivate, _ := rsa.GenerateKey(rand.Reader, 512)
	AuthorPublic := x509.MarshalPKCS1PublicKey(&AuthorPrivate.PublicKey)
	WalletPrivate, _ := rsa.GenerateKey(rand.Reader, 512)
	WalletPublic := x509.MarshalPKCS1PublicKey(&WalletPrivate.PublicKey)
	message := &Message{
		MessageType:     createAudienceMsg,
		Author:          AuthorPublic,
		Message:         []byte{1, 2, 3, 4, 5},
		FeeWallet:       WalletPublic,
		FeeValue:        124,
		Epoch:           1245252,
		PowerOfAttorney: []byte{},
	}
	if ok := message.Sign(AuthorPrivate, WalletPrivate); !ok {
		t.Error("could not sign transfer")
	}
	bytes := message.Serialize()
	copy, err := ParseMessage(bytes)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for Message")
	}
}

func TestMessagePowerOfAttorney(t *testing.T) {
	AuthorPrivate, _ := rsa.GenerateKey(rand.Reader, 512)
	WalletPrivate, _ := rsa.GenerateKey(rand.Reader, 512)
	AttorneyPrivate, _ := rsa.GenerateKey(rand.Reader, 512)

	message := NewMessage(AuthorPrivate, WalletPrivate,
		&About{Details: "Details"}, 10, 100, AttorneyPrivate)
	bytes := message.Serialize()
	copy, err := ParseMessage(bytes)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for message with power of attorney.")
	}
}
