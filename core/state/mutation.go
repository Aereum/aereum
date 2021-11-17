package state

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type StateMutation struct {
	Epoch         uint64
	State         *State
	DeltaWallets  map[crypto.Hash]int
	Hashes        map[crypto.Hash]struct{}
	GrantPower    map[crypto.Hash]struct{}
	RevokePower   map[crypto.Hash]struct{}
	UseAdvOffer   map[crypto.Hash]uint64
	NewAdvOffer   map[crypto.Hash]*instruction.AdvertisingOffer
	NewSubscriber map[crypto.Hash]struct{}
	NewCaption    map[crypto.Hash]struct{}
	NewAudiences  map[crypto.Hash][]byte
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
	if _, ok := s.UseAdvOffer[hash]; ok {
		return false
	}
	s.UseAdvOffer[hash] = expire
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
	if _, ok := s.NewSubscriber[tokenHash]; ok {
		return false
	}
	if _, ok := s.NewCaption[captionHash]; ok {
		return false
	}
	s.NewSubscriber[tokenHash] = struct{}{}
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

func (m *StateMutation) CanPay(payments instruction.Payment) bool {
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

func (m *StateMutation) TransferPayments(payments instruction.Payment) {
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
