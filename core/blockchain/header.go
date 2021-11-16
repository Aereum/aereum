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
package blockchain

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type Header struct {
	Token        []byte
	Parent       crypto.Hash
	ProofOfChain []byte
	Mutations    StateMutation
}

type StateMutation struct {
	parentState               *State
	NewSubsribers             map[crypto.Hash]struct{} // hash token -> hash caption
	NewCaptions               map[crypto.Hash]struct{}
	DeltaWallets              map[crypto.Hash]int
	NewAudiences              map[crypto.Hash]struct{} // Author + Token hash
	GrantPowerOfAttorney      map[crypto.Hash]struct{}
	RevokePowerOfAttorney     map[crypto.Hash]struct{}
	NewAdvertisingOffers      map[crypto.Hash]*instruction.Message
	AcceptedAdvertisingOffers map[crypto.Hash]struct{}
	messages                  []*[]byte
}

func (s *StateMutation) Serialize() []byte {
	return nil
}
