package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
)

type State struct {
	Epoch           uint64
	Members         store.HashVault
	Captions        store.HashVault
	Wallets         store.Wallet
	Audiences       store.Audience
	SponsorOffers   store.Sponsor
	SponsorGranted  store.Sponsor
	PowerOfAttorney store.HashVault
	SponsorExpire   map[uint64]crypto.Hash
}

func GroupMutations(mutations []*StateMutation) *StateMutation {
	grouped := NewStateMutation(mutations[0].Epoch, mutations[len(mutations)-1].State)
	for _, mutation := range mutations {
		for acc, balance := range mutation.DeltaWallets {
			if oldBalance, ok := grouped.DeltaWallets[acc]; ok {
				grouped.DeltaWallets[acc] = oldBalance + balance
			} else {
				grouped.DeltaWallets[acc] = balance
			}
		}
		for hash := range mutation.GrantPower {
			grouped.GrantPower[hash] = struct{}{}
		}
		for hash := range mutation.RevokePower {
			grouped.RevokePower[hash] = struct{}{}
			delete(grouped.GrantPower, hash)
		}
		for hash := range mutation.UseSpnOffer {
			grouped.UseSpnOffer[hash] = struct{}{}
			delete(grouped.NewSpnOffer, hash)
		}
		for hash, offer := range mutation.NewSpnOffer {
			grouped.NewSpnOffer[hash] = offer
		}
		for hash := range mutation.NewMembers {
			grouped.NewMembers[hash] = struct{}{}
		}
		for hash := range mutation.NewCaption {
			grouped.NewCaption[hash] = struct{}{}
		}
		for hash, keys := range mutation.NewAudiences {
			grouped.NewAudiences[hash] = keys
		}

	}
	return grouped
}

func NewStateMutation(epoch uint64, state *State) *StateMutation {
	return &StateMutation{
		Epoch:        epoch,
		State:        state,
		DeltaWallets: make(map[crypto.Hash]int),
		Hashes:       make(map[crypto.Hash]struct{}),
		GrantPower:   make(map[crypto.Hash]struct{}),
		RevokePower:  make(map[crypto.Hash]struct{}),
		UseSpnOffer:  make(map[crypto.Hash]struct{}),
		NewSpnOffer:  make(map[crypto.Hash]*SponsorshipOffer),
		NewMembers:   make(map[crypto.Hash]struct{}),
		NewCaption:   make(map[crypto.Hash]struct{}),
		NewAudiences: make(map[crypto.Hash][]byte),
	}
}

type StateMutation struct {
	Epoch        uint64
	State        *State
	DeltaWallets map[crypto.Hash]int
	Hashes       map[crypto.Hash]struct{}
	GrantPower   map[crypto.Hash]struct{}
	RevokePower  map[crypto.Hash]struct{}
	UseSpnOffer  map[crypto.Hash]struct{}
	NewSpnOffer  map[crypto.Hash]*SponsorshipOffer
	NewMembers   map[crypto.Hash]struct{}
	NewCaption   map[crypto.Hash]struct{}
	NewAudiences map[crypto.Hash][]byte
}

func (s *StateMutation) SetNewHash(hash crypto.Hash) bool {
	if _, ok := s.Hashes[hash]; ok {
		return false
	}
	s.Hashes[hash] = struct{}{}
	return true
}

func (s *StateMutation) SetNewGrantPower(hash crypto.Hash) bool {
	if _, ok := s.GrantPower[hash]; ok {
		return false
	}
	s.GrantPower[hash] = struct{}{}
	return true
}

func (s *StateMutation) SetNewRevokePower(hash crypto.Hash) bool {
	if _, ok := s.RevokePower[hash]; ok {
		return false
	}
	s.RevokePower[hash] = struct{}{}
	return true
}

func (s *StateMutation) SetNewUseAdvOffer(hash crypto.Hash, expire uint64) bool {
	if _, ok := s.UseSpnOffer[hash]; ok {
		return false
	}
	s.UseSpnOffer[hash] = struct{}{}
	return true
}

/*func (s *StateMutation) SetNewAdvOffer(hash crypto.Hash, expire uint64) bool {
	if _, ok := s.UseAdvOffer[hash]; ok {
		return false
	}
	s.NewAdvOffer[hash] = expire
	return true
}*/

func (s *StateMutation) SetNewSubscriber(tokenHash crypto.Hash, captionHash crypto.Hash) bool {
	if _, ok := s.NewMembers[tokenHash]; ok {
		return false
	}
	if _, ok := s.NewMembers[captionHash]; ok {
		return false
	}
	s.NewMembers[tokenHash] = struct{}{}
	s.NewCaption[captionHash] = struct{}{}
	return true
}

func (s *StateMutation) SetNewAudience(hash crypto.Hash, keys []byte) bool {
	if _, ok := s.NewAudiences[hash]; ok {
		return false
	}
	s.NewAudiences[hash] = keys
	return true
}

func (m *StateMutation) CanPay(payments Payment) bool {
	for n, debitAcc := range payments.DebitAcc {
		ok, stateBalance := m.State.Wallets.Balance(debitAcc)
		if !ok {
			return false
		}
		if delta, ok := m.DeltaWallets[debitAcc]; ok {
			if int(stateBalance)+delta < int(payments.DebitValue[n]) {
				return false
			}
		} else {
			if stateBalance < payments.DebitValue[n] {
				return false
			}
		}
	}
	return true
}

func (m *StateMutation) TransferPayments(payments Payment) {
	for n, debitAcc := range payments.DebitAcc {
		if delta, ok := m.DeltaWallets[debitAcc]; ok {
			m.DeltaWallets[debitAcc] = delta - int(payments.DebitValue[n])
		} else {
			m.DeltaWallets[debitAcc] = -int(payments.DebitValue[n])
		}
	}
	for n, creditAcc := range payments.CreditAcc {
		if delta, ok := m.DeltaWallets[creditAcc]; ok {
			m.DeltaWallets[creditAcc] = delta + int(payments.CreditValue[n])
		} else {
			m.DeltaWallets[creditAcc] = int(payments.CreditValue[n])
		}
	}
}
