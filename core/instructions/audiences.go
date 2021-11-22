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

// Create a new audience
type CreateAudience struct {
	Audience      []byte
	Submission    []byte
	Moderation    []byte
	AudienceKey   []byte
	SubmissionKey []byte
	ModerationKey []byte
	Flag          byte
	Description   string
}

func (s *CreateAudience) Validate(validator Validator) bool {
	return validator.GetAudienceKeys(crypto.Hasher(s.Audience)) == nil
}

func (s *CreateAudience) Kind() byte {
	return ICreateAudience
}

func (s *CreateAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutByteArray(s.Submission, &bytes)
	PutByteArray(s.Moderation, &bytes)
	PutByteArray(s.AudienceKey, &bytes)
	PutByteArray(s.SubmissionKey, &bytes)
	PutByteArray(s.ModerationKey, &bytes)
	bytes = append(bytes, s.Flag)
	PutString(s.Description, &bytes)
	return bytes
}

func ParseCreateAudience(data []byte) *CreateAudience {
	p := CreateAudience{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.Submission, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Submission); err != nil {
		return nil
	}
	p.Moderation, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Moderation); err != nil {
		return nil
	}
	p.AudienceKey, position = ParseByteArray(data, position)
	p.SubmissionKey, position = ParseByteArray(data, position)
	p.ModerationKey, position = ParseByteArray(data, position)
	if position >= len(data) {
		return nil
	}
	p.Flag = data[position]
	position += 1
	p.Description, position = ParseString(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}

// Join request for an existing audience
type JoinAudience struct {
	Audience     []byte
	Presentation string
}

func (s *JoinAudience) Validate(validator Validator) bool {
	return validator.GetAudienceKeys(crypto.Hasher(s.Audience)) != nil
}

func (s *JoinAudience) Kind() byte {
	return IJoinAudience
}

func (s *JoinAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutString(s.Presentation, &bytes)
	return bytes
}

func ParseJoinAudience(data []byte) *JoinAudience {
	p := JoinAudience{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.Presentation, position = ParseString(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}

// Accept an audience's join request
type AcceptJoinAudience struct {
	Audience []byte
	Member   []byte
	Read     []byte
	Submit   []byte
	Moderate []byte
}

func (s *AcceptJoinAudience) Validate(validator Validator) bool {
	check := validator.GetAudienceKeys(crypto.Hasher(s.Audience)) != nil
	check = check && validator.HasMember(crypto.Hasher(s.Member))
	return check
}

func (s *AcceptJoinAudience) Kind() byte {
	return IAcceptJoinRequest
}

func (s *AcceptJoinAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutByteArray(s.Member, &bytes)
	PutByteArray(s.Read, &bytes)
	PutByteArray(s.Submit, &bytes)
	PutByteArray(s.Moderate, &bytes)
	return bytes
}

// PRECISA AJUSTAR O PARSE PARA OS CAMPOS OPCIONAIS
func ParseAcceptJoinAudience(data []byte) *AcceptJoinAudience {
	p := AcceptJoinAudience{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.Member, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Member); err != nil {
		return nil
	}
	p.Read, position = ParseByteArray(data, position)
	p.Submit, position = ParseByteArray(data, position)
	p.Moderate, position = ParseByteArray(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}

// Update audience access keys
type UpdateAudience struct {
	Audience      []byte
	Submission    []byte
	Moderation    []byte
	AudienceKey   []byte
	SubmissionKey []byte
	ModerationKey []byte
	ReadMembers   []byte
	SubMembers    []byte
	ModMembers    []byte
}

func (s *UpdateAudience) Validate(validator Validator) bool {
	// TODO: Incorporate logic
	return true
}

func (s *UpdateAudience) Kind() byte {
	return IUpdateAudience
}

func (s *UpdateAudience) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Audience, &bytes)
	PutByteArray(s.Submission, &bytes)
	PutByteArray(s.Moderation, &bytes)
	PutByteArray(s.AudienceKey, &bytes)
	PutByteArray(s.SubmissionKey, &bytes)
	PutByteArray(s.ModerationKey, &bytes)
	PutByteArray(s.ReadMembers, &bytes)
	PutByteArray(s.SubMembers, &bytes)
	PutByteArray(s.ModMembers, &bytes)
	return bytes
}

func ParseUpdateAudience(data []byte) *UpdateAudience {
	p := UpdateAudience{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil {
		return nil
	}
	p.Submission, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Submission); err != nil {
		return nil
	}
	p.Moderation, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Moderation); err != nil {
		return nil
	}
	p.AudienceKey, position = ParseByteArray(data, position)
	p.SubmissionKey, position = ParseByteArray(data, position)
	p.ModerationKey, position = ParseByteArray(data, position)
	p.ReadMembers, position = ParseByteArray(data, position)
	p.SubMembers, position = ParseByteArray(data, position)
	p.ModMembers, position = ParseByteArray(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}
