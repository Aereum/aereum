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
package instructionsold

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestTransfer(t *testing.T) {
	FromPublic, FromPrivate := crypto.RandomAsymetricKey()
	transfer := &Transfer{
		From:  FromPublic.ToBytes(),
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
	AuthorPublic, AuthorPrivate := crypto.RandomAsymetricKey()
	WalletPublic, WalletPrivate := crypto.RandomAsymetricKey()
	message := &Message{
		MessageType:     CreateAudienceMsg,
		Author:          AuthorPublic.ToBytes(),
		Message:         []byte{1, 2, 3, 4, 5},
		FeeWallet:       WalletPublic.ToBytes(),
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
	_, AuthorPrivate := crypto.RandomAsymetricKey()
	_, WalletPrivate := crypto.RandomAsymetricKey()
	_, AttorneyPrivate := crypto.RandomAsymetricKey()

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
