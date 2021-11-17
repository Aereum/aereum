package state

import (
	"crypto"

	"github.com/Aereum/aereum/core/instruction"
	"github.com/Aereum/aereum/core/store"
)

type Validator interface {
	Validate(i *instruction.Instruction) bool
}

type State struct {
	Epoch           uint64
	Subscribers     store.HashVault
	Captions        store.HashVault
	Wallets         store.Wallet
	Audiences       store.Audience
	SponsorOffers   store.Sponsor
	SponsorGranted  store.Sponsor
	PowerOfAttorney store.HashVault
	SponsorExpire   map[uint64]crypto.Hash
}

func (s *State) IncorporateMutations(epoch uint64, mutations *StateMutation) {
	for hash, delta := range mutations.DeltaWallets {
		if delta > 0 {
			s.Wallets.Credit(hash, uint64(delta))
		} else if delta < 0 {
			s.Wallets.Debit(hash, uint64(-delta))
		}
	}
	for hash := range mutations.GrantPower {
		s.PowerOfAttorney.Insert(hash)
	}
	for hash := range mutations.RevokePower {
		s.PowerOfAttorney.Remove(hash)
	}
	for hash := range mutations.NewSubscriber {
		s.Subscribers.Insert(hash)
	}
	for hash := range mutations.NewCaption {
		s.Captions.Insert(hash)
	}
	for hash, keys := range mutations.NewAudiences {
		s.Audiences.SetKeys(hash, keys)
	}
	// sponsorship
}
