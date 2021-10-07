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
	"github.com/Aereum/aereum/core/hashdb"
	"github.com/Aereum/aereum/core/message"
)

type Header struct {
	Token        []byte
	Parent       hashdb.Hash
	ProofOfChain []byte
	Mutations    StateMutation
}

type StateMutation struct {
	parentState               *State
	NewSubsribers             map[hashdb.Hash]struct{} // hash token -> hash caption
	NewCaptions               map[hashdb.Hash]struct{}
	DeltaWallets              map[hashdb.Hash]int
	NewAudiences              map[hashdb.Hash]struct{} // Author + Token hash
	GrantPowerOfAttorney      map[hashdb.Hash]struct{}
	RevokePowerOfAttorney     map[hashdb.Hash]struct{}
	NewAdvertisingOffers      map[hashdb.Hash]*message.Message
	AcceptedAdvertisingOffers map[hashdb.Hash]struct{}
	messages                  []*[]byte
}

func (s *StateMutation) Serialize() []byte {

}
