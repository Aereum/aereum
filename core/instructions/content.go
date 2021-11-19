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
package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
)

// Post content to an existing audience
type Content struct {
	Audience     []byte
	ContentType  string
	Content      []byte
	Hash         []byte
	Sponsored    bool
	Encrypted    bool
	SubSignature []byte
	ModSignature []byte
}

func (s *Content) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutString(s.ContentType, &bytes)
	PutByteArray(s.Content, &bytes)
	PutByteArray(s.Hash, &bytes) // NAO SEI SE ESTA CERTA ESSA SERIALIZACAO DE HASH
	PutBool(s.Sponsored, &bytes)
	PutBool(s.Encrypted, &bytes)
	PutByteArray(s.SubSignature, &bytes)
	PutByteArray(s.ModSignature, &bytes)
	return bytes
}

// PRECISA AJUSTAR O PARSE PARA OS CAMPOS OPCIONAIS
func ParseContent(data []byte) *Content {
	p := Content{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.ContentType, position = ParseString(data, position)
	p.Content, position = ParseByteArray(data, position)
	p.Hash, position = ParseByteArray(data, position)
	p.Sponsored, position = ParseBool(data, position)
	p.Encrypted, position = ParseBool(data, position)
	p.SubSignature, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.SubSignature); err != nil {
		return nil
	}
	p.ModSignature, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.ModSignature); err != nil {
		return nil
	}
	if position == len(data) {
		return &p
	}
	return nil
}

// Reaction to a content message
type React struct {
	Hash     []byte
	Reaction byte
}

func (s *React) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Hash, &bytes)
	bytes = append(bytes, s.Reaction)
	return bytes
}

func ParseReact(data []byte) *React {
	p := React{}
	position := 0
	p.Hash, position = ParseByteArray(data, position)
	p.Reaction, position = ParseByte(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}
