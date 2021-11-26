package instructions

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
	Audiences       *store.Audience
	SponsorOffers   *store.HashExpireVault
	SponsorGranted  *store.Sponsor
	PowerOfAttorney *store.HashVault
	EphemeralTokens *store.HashExpireVault
	SponsorExpire   map[uint64]crypto.Hash
	EphemeralExpire map[uint64]crypto.Hash
}

func NewGenesisState() (*State, crypto.PrivateKey) {
	pubKey, prvKey := crypto.RandomAsymetricKey()
	hash := crypto.Hasher(pubKey.ToBytes())
	state := State{
		Epoch:           0,
		Members:         store.NewHashVault("members", 0, 8),
		Captions:        store.NewHashVault("captions", 0, 8),
		Wallets:         store.NewMemoryWalletStore(0, 8),
		Audiences:       store.NewMemoryAudienceStore(0, 8),
		SponsorOffers:   store.NewExpireHashVault("sponsoroffer", 0, 8),
		SponsorGranted:  store.NewSponsorShipOfferStore(0, 8),
		PowerOfAttorney: store.NewHashVault("poa", 0, 8),
		EphemeralTokens: store.NewExpireHashVault("ephemeral", 0, 8),
		SponsorExpire:   make(map[uint64]crypto.Hash),
		EphemeralExpire: make(map[uint64]crypto.Hash),
	}
	state.Members.Insert(hash)
	state.Captions.Insert(crypto.Hasher([]byte("Aereum Network Genesis")))
	state.Wallets.Credit(hash, 1e6)
	return &state, prvKey
}

func (s *State) IncorporateBlock(b *Block) {
	for hash := range b.mutations.NewCaption {
		s.Captions.Insert(hash)
	}
	for hash := range b.mutations.NewMembers {
		s.Members.Insert(hash)
	}
	for acc, delta := range b.mutations.DeltaWallets {
		if delta > 0 {
			s.Wallets.Credit(acc, uint64(delta))
		} else if delta < 0 {
			s.Wallets.Debit(acc, uint64(-delta))
		}
	}
	for hash := range b.mutations.GrantPower {
		s.PowerOfAttorney.Insert(hash)
	}
	for hash := range b.mutations.RevokePower {
		s.PowerOfAttorney.Remove(hash)
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
	for hash, keys := range b.mutations.NewAudiences {
		s.Audiences.SetKeys(hash, keys)
	}
	for hash, keys := range b.mutations.UpdAudiences {
		s.Audiences.SetKeys(hash, keys)
	}
}
