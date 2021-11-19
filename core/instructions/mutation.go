package instructions

import "github.com/Aereum/aereum/core/crypto"

type sponsorOfferState struct {
	contentHash crypto.Hash
	expire      uint64
}

type Mutation struct {
	DeltaWallets map[crypto.Hash]int
	GrantPower   map[crypto.Hash]struct{}
	RevokePower  map[crypto.Hash]struct{}
	UseSpnOffer  map[crypto.Hash]struct{}
	NewSpnOffer  map[crypto.Hash]*sponsorOfferState
	NewMembers   map[crypto.Hash]struct{}
	NewCaption   map[crypto.Hash]struct{}
	NewAudiences map[crypto.Hash][]byte
	NewEphemeral map[crypto.Hash]uint64
}

func NewMutation() *Mutation {
	return &Mutation{
		DeltaWallets: make(map[crypto.Hash]int),
		GrantPower:   make(map[crypto.Hash]struct{}),
		RevokePower:  make(map[crypto.Hash]struct{}),
		UseSpnOffer:  make(map[crypto.Hash]struct{}),
		NewSpnOffer:  make(map[crypto.Hash]*sponsorOfferState),
		NewMembers:   make(map[crypto.Hash]struct{}),
		NewCaption:   make(map[crypto.Hash]struct{}),
		NewAudiences: make(map[crypto.Hash][]byte),
		NewEphemeral: make(map[crypto.Hash]uint64),
	}
}

func (m *Mutation) DeltaBalance(hash crypto.Hash) int {
	balance := m.DeltaWallets[hash]
	return balance
}

func (m *Mutation) HasGrantPower(hash crypto.Hash) bool {
	_, ok := m.GrantPower[hash]
	return ok
}

func (m *Mutation) HasRevokePower(hash crypto.Hash) bool {
	_, ok := m.RevokePower[hash]
	return ok
}

func (m *Mutation) HasUsedSponsorOffer(hash crypto.Hash) bool {
	_, ok := m.UseSpnOffer[hash]
	return ok
}

func (m *Mutation) GetSponsorOffer(hash crypto.Hash) *sponsorOfferState {
	offer := m.NewSpnOffer[hash]
	return offer
}

func (m *Mutation) HasMember(hash crypto.Hash) bool {
	_, ok := m.NewMembers[hash]
	return ok
}

func (m *Mutation) HasCaption(hash crypto.Hash) bool {
	_, ok := m.NewCaption[hash]
	return ok
}

func (m *Mutation) GetAudience(hash crypto.Hash) []byte {
	audience := m.NewAudiences[hash]
	return audience
}

func (m *Mutation) HasEphemeral(hash crypto.Hash) (bool, uint64) {
	expire, ok := m.NewEphemeral[hash]
	return ok, expire
}

func GroupMutations(mutations []*Mutation) *Mutation {
	grouped := NewMutation()
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

/*func (m *StateMutation) CanPay(payments Payment) bool {
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
*/

/*
func setNewHash(hash crypto.Hash, store map[crypto.Hash]struct{}) bool {
	if _, ok := store[hash]; ok {
		return false
	}
	store[hash] = struct{}{}
	return true
}

func (s *Mutation) SetNewHash(hash crypto.Hash) bool {
	return setNewHash(hash, s.Hashes)
}

func (s *Mutation) SetNewGrantPower(hash crypto.Hash) bool {
	return setNewHash(hash, s.GrantPower)
}

func (s *Mutation) SetNewRevokePower(hash crypto.Hash) bool {
	return setNewHash(hash, s.RevokePower)
}

func (s *Mutation) SetNewUseSonOffer(hash crypto.Hash, expire uint64) bool {
	return setNewHash(hash, s.UseSpnOffer)
}

func (s *Mutation) SetNewAdvOffer(hash crypto.Hash, offer SponsorshipOffer) bool {
	if _, ok := s.NewSpnOffer[hash]; ok {
		return false
	}
	s.NewSpnOffer[hash] = &offer
	return true
}

func (s *Mutation) SetNewMember(tokenHash crypto.Hash, captionHash crypto.Hash) bool {
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

func (s *Mutation) SetNewAudience(hash crypto.Hash, keys []byte) bool {
	if _, ok := s.NewAudiences[hash]; ok {
		return false
	}
	s.NewAudiences[hash] = keys
	return true
}
*/
