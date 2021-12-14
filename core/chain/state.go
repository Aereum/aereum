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

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
)

var (
	ErrNotSubsequentBlock = errors.New("cannot incorporate a non-subsequent block")
	ErrIncorporationError = errors.New("could not incorporate block")
)

type State struct {
	Epoch           uint64
	Members         *store.HashVault
	Captions        *store.HashVault
	Wallets         *store.Wallet
	Stages          *store.Stage
	SponsorOffers   *store.HashExpireVault
	SponsorGranted  *store.Sponsor
	PowerOfAttorney *store.HashVault
	EphemeralTokens *store.HashExpireVault
	SponsorExpire   map[uint64]crypto.Hash
	EphemeralExpire map[uint64]crypto.Hash
}

func NewGenesisState() (*State, crypto.PrivateKey) {
	pubKey, prvKey := crypto.RandomAsymetricKey()
	state := State{
		Epoch:           0,
		Members:         store.NewHashVault("members", 0, 8),
		Captions:        store.NewHashVault("captions", 0, 8),
		Wallets:         store.NewMemoryWalletStore(0, 8),
		Stages:          store.NewMemoryAudienceStore(0, 8),
		SponsorOffers:   store.NewExpireHashVault("sponsoroffer", 0, 8),
		SponsorGranted:  store.NewSponsorShipOfferStore(0, 8),
		PowerOfAttorney: store.NewHashVault("poa", 0, 8),
		EphemeralTokens: store.NewExpireHashVault("ephemeral", 0, 8),
		SponsorExpire:   make(map[uint64]crypto.Hash),
		EphemeralExpire: make(map[uint64]crypto.Hash),
	}
	state.Members.InsertToken(pubKey)
	state.Captions.InsertHash(crypto.Hasher([]byte("Aereum Network Genesis")))
	state.Wallets.Credit(pubKey, 1e6)
	return &state, prvKey
}

func NewGenesisStateWithToken(token crypto.PrivateKey) *State {
	hash := crypto.HashToken(token.PublicKey())
	state := State{
		Epoch:           0,
		Members:         store.NewHashVault("members", 0, 8),
		Captions:        store.NewHashVault("captions", 0, 8),
		Wallets:         store.NewMemoryWalletStore(0, 8),
		Stages:          store.NewMemoryAudienceStore(0, 8),
		SponsorOffers:   store.NewExpireHashVault("sponsoroffer", 0, 8),
		SponsorGranted:  store.NewSponsorShipOfferStore(0, 8),
		PowerOfAttorney: store.NewHashVault("poa", 0, 8),
		EphemeralTokens: store.NewExpireHashVault("ephemeral", 0, 8),
		SponsorExpire:   make(map[uint64]crypto.Hash),
		EphemeralExpire: make(map[uint64]crypto.Hash),
	}
	state.Members.InsertHash(hash)
	state.Captions.InsertHash(crypto.Hasher([]byte("Aereum Network Genesis")))
	state.Wallets.CreditHash(hash, 1e6)
	return &state
}

func (s *State) IncorporateBlock(b *Block) {
	for hash := range b.mutations.NewCaption {
		s.Captions.InsertHash(hash)
	}
	for hash := range b.mutations.NewMembers {
		s.Members.InsertHash(hash)
	}
	for acc, delta := range b.mutations.DeltaWallets {
		if delta > 0 {
			s.Wallets.CreditHash(acc, uint64(delta))
		} else if delta < 0 {
			s.Wallets.DebitHash(acc, uint64(-delta))
		}
	}
	for hash := range b.mutations.GrantPower {
		s.PowerOfAttorney.InsertHash(hash)
	}
	for hash := range b.mutations.RevokePower {
		s.PowerOfAttorney.RemoveHash(hash)
	}
	for hash := range b.mutations.PublishSpn {
		s.SponsorGranted.RemoveContentHash(hash)
	}
	for token, contentHash := range b.mutations.GrantSponsor {
		s.SponsorGranted.SetContentHash(token, contentHash[:])
	}
	for hash, expire := range b.mutations.NewSpnOffer {
		s.SponsorOffers.Insert(hash, expire)
	}
	for hash, keys := range b.mutations.NewStages {
		s.Stages.SetKeys(hash, &keys)
	}
	for hash, keys := range b.mutations.StageUpdate {
		s.Stages.SetKeys(hash, &keys)
	}
	s.Wallets.CreditHash(crypto.HashToken(b.Publisher), b.FeesCollected)
}
