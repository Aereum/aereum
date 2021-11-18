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
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)

// Create a new audience
type CreateAudience struct {
	Audience		[]byte	
	Sumission		[]byte
	Moderation		[]byte
	AudienceKey		[]byte
	SumissionKey	[]byte
	ModerationKey	[]byte
	Flag			byte
	Description		string
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
	PutString(s.Description)
	return bytes
}

func ParseCreateAudience(data []byte) *CreateAudience {
    p := CreateAudience{}
    position := 0
    p.Audience, position = ParseByteArray(data, position)
    if _, err := crypto.PublicKeyFromBytes(p.Audience); err != nil  {
        return nil
    }
    p.Submission, position = ParseByteArray(data, position)
    if _, err := crypto.PublicKeyFromBytes(p.Submission); err != nil  {
        return nil
    }
    p.Moderation, position = ParseByteArray(data, position)
    if _, err := crypto.PublicKeyFromBytes(p.Moderation); err != nil  {
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
type JoinAudiece struct {
	Audience		[]byte
	Presentation	string
}

func (s *JoinAudiece) Serialize() []byte {
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
	Audience	[]byte
	Member		[]byte
	Read		[]byte
	Submit		[]byte
	Moderate	[]byte
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

func ParseAcceptJoinAudience(data []byte) *AcceptJoinAudience {
	p := AcceptJoinAudiece{}
	position := 0
	p.Audience, position = ParseByteArray(data, position)
	p.Audience, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(s.Audience); err != nil {
		return nil
	}
	p.Presentation, position = ParseString(data, position)
	if position == len(data) {
        return &p
    }
    return nil

}