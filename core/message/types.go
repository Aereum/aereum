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

import "github.com/Aereum/aereum/core/crypto"

type Serializer interface {
	Serialize() []byte
	Kind() byte
}

type Subscribe struct {
	Token   []byte
	Caption string
	Details string
}

func (s *Subscribe) Kind() byte {
	return SubscribeMsg
}

func (s *Subscribe) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Token, &bytes)
	PutString(s.Caption, &bytes)
	PutString(s.Details, &bytes)
	return bytes
}

func ParseSubscribe(data []byte) *Subscribe {
	s := Subscribe{}
	position := 0
	s.Token, position = ParseByteArray(data, position)
	s.Caption, position = ParseString(data, position)
	s.Details, position = ParseString(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type About struct {
	Details string
}

func (s *About) Kind() byte {
	return AboutMsg
}

func (s *About) Serialize() []byte {
	bytes := make([]byte, 0)
	PutString(s.Details, &bytes)
	return bytes
}

func ParseAbout(data []byte) *About {
	s := About{}
	position := 0
	s.Details, position = ParseString(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type CreateAudience struct {
	Token []byte // Private Key allows change in this structure
	//Moderate []byte // Private Key allows to validate join requests
	//Write    []byte // Private Key allows to submit messages to audience
	//Reply    []byte // Private Key allows to interact with messages
	//Read     []byte // Private Key allows to decrypt messages to audience
	Description string
}

func (s *CreateAudience) Kind() byte {
	return CreateAudienceMsg
}

func (s *CreateAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Token, &bytes)
	PutString(s.Description, &bytes)
	return bytes
}

func ParseCreateAudience(data []byte) *CreateAudience {
	s := CreateAudience{}
	position := 0
	s.Token, position = ParseByteArray(data, position)
	s.Description, position = ParseString(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type JoinAudience struct {
	Audience []byte
	Expire   uint64
}

func (s *JoinAudience) Kind() byte {
	return JoinAudienceMsg
}

func (s *JoinAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutUint64(s.Expire, &bytes)
	return bytes
}

func ParseJoinAudience(data []byte) *JoinAudience {
	s := JoinAudience{}
	position := 0
	s.Audience, position = ParseByteArray(data, position)
	s.Expire, position = ParseUint64(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type Follower struct {
	Token  []byte
	Secret []byte
}

type ChangeAudience struct {
	Audience []byte
	Add      []*Follower
	Remove   []*Follower
	Details  string
}

func (s *ChangeAudience) Kind() byte {
	return AudienceChangeMsg
}

func (s *ChangeAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	l := len(s.Add)
	if l > 65535 {
		l = 65535
	}
	bytes = append(bytes, byte(l), byte(l>>8))
	for n := 0; n < l; n++ {
		PutByteArray(s.Add[n].Token, &bytes)
		PutByteArray(s.Add[n].Secret, &bytes)
	}
	l = len(s.Remove)
	if l > 65535 {
		l = 65535
	}
	bytes = append(bytes, byte(l), byte(l>>8))
	for n := 0; n < l; n++ {
		PutByteArray(s.Remove[n].Token, &bytes)
		PutByteArray(s.Remove[n].Secret, &bytes)
	}
	PutString(s.Details, &bytes)
	return bytes
}

func ParseChangeAudience(data []byte) *ChangeAudience {
	length := len(data)
	s := ChangeAudience{}
	position := 0
	s.Audience, position = ParseByteArray(data, position)
	if position+1 >= length {
		return nil
	}
	addLength := int(data[position]) | int(data[position+1])<<8
	position += 2
	// add
	add := make([]*Follower, addLength)
	for n := 0; n < addLength; n++ {
		f := Follower{}
		f.Token, position = ParseByteArray(data, position)
		f.Secret, position = ParseByteArray(data, position)
		add[n] = &f
	}
	s.Add = add
	if position+1 >= length {
		return nil
	}
	removeLength := int(data[position]) | int(data[position+1])<<8
	position += 2
	remove := make([]*Follower, removeLength)
	for n := 0; n < removeLength; n++ {
		f := Follower{}
		f.Token, position = ParseByteArray(data, position)
		f.Secret, position = ParseByteArray(data, position)
		remove[n] = &f
	}
	s.Remove = remove
	s.Details, position = ParseString(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type AdvertisingOffer struct {
	Token          []byte
	Audience       []byte
	ContentType    string
	ContentData    []byte
	AdvertisingFee uint64
	Expire         uint64
}

func (s *AdvertisingOffer) Kind() byte {
	return AdvertisingOfferMsg
}

func (s *AdvertisingOffer) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Token, &bytes)
	PutByteArray(s.Audience, &bytes)
	PutString(s.ContentType, &bytes)
	PutByteArray(s.ContentData, &bytes)
	PutUint64(s.AdvertisingFee, &bytes)
	PutUint64(s.Expire, &bytes)
	return bytes
}

func ParseAdvertisingOffer(data []byte) *AdvertisingOffer {
	s := AdvertisingOffer{}
	position := 0
	s.Token, position = ParseByteArray(data, position)
	s.Audience, position = ParseByteArray(data, position)
	s.ContentType, position = ParseString(data, position)
	s.ContentData, position = ParseByteArray(data, position)
	s.AdvertisingFee, position = ParseUint64(data, position)
	s.Expire, position = ParseUint64(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type Content struct {
	Audience         []byte
	ContentSecret    []byte
	ContentType      string
	ContentData      []byte
	AdvertisingToken []byte
	HashContent      []byte
	SubmitSignature  []byte
	PublishSignature []byte
	SubmitHash       crypto.Hash
	PublishHash      crypto.Hash
}

func (s *Content) Kind() byte {
	return ContentMsg
}

func (s *Content) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutByteArray(s.ContentSecret, &bytes)
	PutString(s.ContentType, &bytes)
	PutByteArray(s.ContentData, &bytes)
	PutByteArray(s.AdvertisingToken, &bytes)
	PutByteArray(s.HashContent, &bytes)
	PutByteArray(s.SubmitSignature, &bytes)
	PutByteArray(s.PublishSignature, &bytes)
	return bytes
}

func ParseContent(data []byte) *Content {
	s := Content{}
	position := 0
	s.Audience, position = ParseByteArray(data, position)
	s.ContentSecret, position = ParseByteArray(data, position)
	s.ContentType, position = ParseString(data, position)
	s.ContentData, position = ParseByteArray(data, position)
	s.AdvertisingToken, position = ParseByteArray(data, position)
	s.HashContent, position = ParseByteArray(data, position)
	s.SubmitHash = crypto.Hasher(data[0:position])
	s.SubmitSignature, position = ParseByteArray(data, position)
	s.PublishHash = crypto.Hasher(data[0:position])
	s.PublishSignature, position = ParseByteArray(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type GrantPowerOfAttorney struct {
	Token  []byte
	Expire uint64
}

func (s *GrantPowerOfAttorney) Kind() byte {
	return GrantPowerOfAttorneyMsg
}

func (s *GrantPowerOfAttorney) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Token, &bytes)
	PutUint64(s.Expire, &bytes)
	return bytes
}

func ParseGrantPowerOfAttorney(data []byte) *GrantPowerOfAttorney {
	s := GrantPowerOfAttorney{}
	position := 0
	s.Token, position = ParseByteArray(data, position)
	s.Expire, position = ParseUint64(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}

type RevokePowerOfAttorney struct {
	Token []byte
}

func (s *RevokePowerOfAttorney) Kind() byte {
	return RevokePowerOfAttorneyMsg
}

func (s *RevokePowerOfAttorney) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Token, &bytes)
	return bytes
}

func ParseRevokePowerOfAttorney(data []byte) *RevokePowerOfAttorney {
	s := RevokePowerOfAttorney{}
	position := 0
	s.Token, position = ParseByteArray(data, position)
	if position == len(data) {
		return &s
	}
	return nil
}
