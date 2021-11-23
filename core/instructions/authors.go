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
	"encoding/json"

	"github.com/Aereum/aereum/core/crypto"
)

// Join network
type JoinNetwork struct {
	Caption string
	Details string
}

func (s *JoinNetwork) Validate(validator Validator) bool {
	if validator.HasCaption(crypto.Hasher([]byte(s.Caption))) {
		return false
	}
	return true
}

func (s *JoinNetwork) Kind() byte {
	return IJoinNetwork
}

func (s *JoinNetwork) Serialize() []byte {
	bytes := make([]byte, 0)
	PutString(s.Caption, &bytes)
	PutString(s.Details, &bytes)
	return bytes
}

func ParseJoinNetwork(data []byte) *JoinNetwork {
	p := JoinNetwork{}
	position := 0
	p.Caption, position = ParseString(data, position)
	p.Details, position = ParseString(data, position)
	if !json.Valid([]byte(p.Details)) {
		return nil
	}
	if position == len(data) {
		return &p
	}
	return nil
}

// Update member information
type UpdateInfo struct {
	Details string
}

func (s *UpdateInfo) Validate(validator Validator) bool {
	return true
}

func (s *UpdateInfo) Kind() byte {
	return IUpdateInfo
}

func (s *UpdateInfo) Serialize() []byte {
	bytes := make([]byte, 0)
	PutString(s.Details, &bytes)
	return bytes
}

func ParseUpdateInfo(data []byte) *UpdateInfo {
	p := UpdateInfo{}
	position := 0
	p.Details, position = ParseString(data, position)
	if position == len(data) {
		return &p
	}
	if !json.Valid([]byte(p.Details)) {
		return nil
	}
	return nil
}

// Grant power of attorney to a network member
type GrantPowerOfAttorney struct {
	Attorney []byte
}

func (s *GrantPowerOfAttorney) Validate(validator Validator) bool {
	return validator.HasMember(crypto.Hasher(s.Attorney))
}

func (s *GrantPowerOfAttorney) Kind() byte {
	return IGrantPowerOfAttorney
}

func (s *GrantPowerOfAttorney) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Attorney, &bytes)
	return bytes
}

func ParseGrantPowerOfAttorney(data []byte) *GrantPowerOfAttorney {
	p := GrantPowerOfAttorney{}
	position := 0
	p.Attorney, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Attorney); err != nil {
		return nil
	}
	if position == len(data) {
		return &p
	}
	return nil
}

// Revoke power of attorney previously granted
type RevokePowerOfAttorney struct {
	Attorney []byte
}

func (s *RevokePowerOfAttorney) Validate(validator Validator) bool {
	return true
}

func (s *RevokePowerOfAttorney) Kind() byte {
	return IRevokePowerOfAttorney
}

func (s *RevokePowerOfAttorney) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.Attorney, &bytes)
	return bytes
}

func ParseRevokePowerOfAttorney(data []byte) *RevokePowerOfAttorney {
	p := RevokePowerOfAttorney{}
	position := 0
	p.Attorney, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.Attorney); err != nil {
		return nil
	}
	if position == len(data) {
		return &p
	}
	return nil
}

// Create ephemeral token for anonymous messages
type CreateEphemeral struct {
	EphemeralToken []byte
	Expiry         uint64
}

func (s *CreateEphemeral) Validate(validator Validator) bool {
	// TODO
	return true
}

func (s *CreateEphemeral) Kind() byte {
	return ICreateEphemeral
}

func (s *CreateEphemeral) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.EphemeralToken, &bytes)
	PutUint64(s.Expiry, &bytes)
	return bytes
}

func ParseCreateEphemeral(data []byte) *CreateEphemeral {
	p := CreateEphemeral{}
	position := 0
	p.EphemeralToken, position = ParseByteArray(data, position)
	if _, err := crypto.PublicKeyFromBytes(p.EphemeralToken); err != nil {
		return nil
	}
	p.Expiry, position = ParseUint64(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}

// Anonymously establish message exchange with another network member
type SecureChannel struct {
	TokenRange     []byte
	Nonce          uint64
	EncryptedNonce []byte
	Content        []byte
}

func (s *SecureChannel) Validate(validator Validator) bool {
	// TODO
	return true
}

func (s *SecureChannel) Kind() byte {
	return ISecureChannel
}

func (s *SecureChannel) Serialize() []byte {
	bytes := make([]byte, 0)
	PutByteArray(s.TokenRange, &bytes)
	PutUint64(s.Nonce, &bytes)
	PutByteArray(s.EncryptedNonce, &bytes)
	PutByteArray(s.Content, &bytes)
	return bytes
}

func ParseSecureChannel(data []byte) *SecureChannel {
	p := SecureChannel{}
	position := 0
	p.TokenRange, position = ParseByteArray(data, position)
	p.Nonce, position = ParseUint64(data, position)
	p.EncryptedNonce, position = ParseByteArray(data, position)
	p.Content, position = ParseByteArray(data, position)
	if position == len(data) {
		return &p
	}
	return nil
}
