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

package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
)

// Post sponsored content offer to an audience
type SponsorshipOffer struct {
	Audience    []byte
	ContentType string
	Content     []byte // NAO SEI SE PRECISARIA INCLUIR UM CAMPO COM HASH DO CONTENT
	Expiry      uint64
	Revenue     uint64
}

func (s *SponsorshipOffer) Kind() byte {
	return ISponsorshipOffer
}

func (s *SponsorshipOffer) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutString(s.ContentType, &bytes)
	PutByteArray(s.Content, &bytes)
	PutUint64(s.Expiry, &bytes)
	PutUint64(s.Revenue, &bytes)
	return bytes
}

func ParseSponsorshipOffer(data []byte) *SponsorshipOffer {
	p := SponsorshipOffer{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.ContentType, position = ParseString(data, position)
	p.Content, position = ParseByteArray(data, position)
	p.Expiry, position = ParseUint64(data, position)
	p.Revenue, position = ParseUint64(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}

// Accept the sponsored content offer
type SponsorshipAcceptance struct {
	Audience     []byte
	Hash         []byte
	ModSignature []byte
}

func (s *SponsorshipAcceptance) Kind() byte {
	return ISponsorshipAcceptance
}

func (s *SponsorshipAcceptance) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutByteArray(s.Hash, &bytes)
	PutByteArray(s.ModSignature, &bytes)
	return bytes
}

func ParseSponsorshipAcceptance(data []byte) *SponsorshipAcceptance {
	p := SponsorshipAcceptance{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.Hash, position = ParseByteArray(data, position)
	p.ModSignature, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.ModSignature); err != nil {
		return nil
	}
	if position == len(data) {
		return &p
	}
	return nil
}
