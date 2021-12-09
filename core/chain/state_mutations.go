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
// along with the aereum library. If not, see <http://www.gnu.org/licenses/>.
package chain

import "github.com/Aereum/aereum/core/crypto"

type mutation struct {
	DeltaWallets map[crypto.Hash]int
	GrantPower   map[crypto.Hash]struct{}
	RevokePower  map[crypto.Hash]struct{}
	UseSpnOffer  map[crypto.Hash]struct{}
	GrantSponsor map[crypto.Hash]crypto.Hash // hash of sponsor token + audience -> content hash
	PublishSpn   map[crypto.Hash]struct{}
	NewSpnOffer  map[crypto.Hash]uint64
	NewMembers   map[crypto.Hash]struct{}
	NewCaption   map[crypto.Hash]struct{}
	NewAudiences map[crypto.Hash][]byte
	UpdAudiences map[crypto.Hash][]byte
	NewEphemeral map[crypto.Hash]uint64
}

func NewMutation() *mutation {
	return &mutation{
		DeltaWallets: make(map[crypto.Hash]int),
		GrantPower:   make(map[crypto.Hash]struct{}),
		RevokePower:  make(map[crypto.Hash]struct{}),
		UseSpnOffer:  make(map[crypto.Hash]struct{}),
		GrantSponsor: make(map[crypto.Hash]crypto.Hash),
		PublishSpn:   make(map[crypto.Hash]struct{}),
		NewSpnOffer:  make(map[crypto.Hash]uint64),
		NewMembers:   make(map[crypto.Hash]struct{}),
		NewCaption:   make(map[crypto.Hash]struct{}),
		NewAudiences: make(map[crypto.Hash][]byte),
		UpdAudiences: make(map[crypto.Hash][]byte),
		NewEphemeral: make(map[crypto.Hash]uint64),
	}
}

func (m *mutation) DeltaBalance(hash crypto.Hash) int {
	balance := m.DeltaWallets[hash]
	return balance
}

func (m *mutation) HasGrantedSponsorship(hash crypto.Hash) (bool, crypto.Hash) {
	if _, ok := m.PublishSpn[hash]; ok {
		return false, crypto.Hasher([]byte{})
	}
	contentHash, ok := m.GrantSponsor[hash]
	return ok, contentHash
}

func (m *mutation) HasGrantPower(hash crypto.Hash) bool {
	_, ok := m.GrantPower[hash]
	return ok
}

func (m *mutation) HasRevokePower(hash crypto.Hash) bool {
	_, ok := m.RevokePower[hash]
	return ok
}

func (m *mutation) HasUsedSponsorOffer(hash crypto.Hash) bool {
	_, ok := m.UseSpnOffer[hash]
	return ok
}

func (m *mutation) GetSponsorOffer(hash crypto.Hash) bool {
	_, ok := m.NewSpnOffer[hash]
	return ok
}

func (m *mutation) HasMember(hash crypto.Hash) bool {
	_, ok := m.NewMembers[hash]
	return ok
}

func (m *mutation) HasCaption(hash crypto.Hash) bool {
	_, ok := m.NewCaption[hash]
	return ok
}

func (m *mutation) GetAudience(hash crypto.Hash) []byte {
	if audience, ok := m.UpdAudiences[hash]; ok {
		return audience
	}
	audience := m.NewAudiences[hash]
	return audience
}

func (m *mutation) HasEphemeral(hash crypto.Hash) (bool, uint64) {
	expire, ok := m.NewEphemeral[hash]
	return ok, expire
}

func GroupBlockMutations(blocks []*Block) *mutation {
	grouped := NewMutation()
	for _, block := range blocks {
		for acc, balance := range block.mutations.DeltaWallets {
			if oldBalance, ok := grouped.DeltaWallets[acc]; ok {
				grouped.DeltaWallets[acc] = oldBalance + balance
			} else {
				grouped.DeltaWallets[acc] = balance
			}
		}
		for hash := range block.mutations.GrantPower {
			grouped.GrantPower[hash] = struct{}{}
		}
		for hash := range block.mutations.RevokePower {
			grouped.RevokePower[hash] = struct{}{}
			delete(grouped.GrantPower, hash)
		}
		for hash := range block.mutations.UseSpnOffer {
			grouped.UseSpnOffer[hash] = struct{}{}
			delete(grouped.NewSpnOffer, hash)
		}
		for hash, offer := range block.mutations.NewSpnOffer {
			grouped.NewSpnOffer[hash] = offer
		}
		for hash := range block.mutations.NewMembers {
			grouped.NewMembers[hash] = struct{}{}
		}
		for hash := range block.mutations.NewCaption {
			grouped.NewCaption[hash] = struct{}{}
		}
		for hash, keys := range block.mutations.NewAudiences {
			grouped.NewAudiences[hash] = keys
		}
		// incorporate fees to block publisher
		if balance, ok := grouped.DeltaWallets[crypto.Hasher(block.Publisher)]; ok {
			grouped.DeltaWallets[crypto.Hasher(block.Publisher)] = balance + int(block.FeesCollected)
		} else {
			grouped.DeltaWallets[crypto.Hasher(block.Publisher)] = int(block.FeesCollected)
		}
	}
	return grouped
}
