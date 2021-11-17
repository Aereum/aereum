package state

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type CheckPoint struct {
	State     *State
	Mutations *StateMutation
}

func (s *CheckPoint) Validate(i *instruction.Instruction) bool {
	return true
}

func (s *CheckPoint) AuthorExists(m *instruction.Message) bool {
	hash := crypto.Hasher(m.Author)
	if s.State.Subscribers.Exists(hash) {
		return true
	}
	if _, ok := s.Mutations.NewSubscriber[hash]; ok {
		return true
	}
	return false
}

func (s *CheckPoint) CaptionExists(caption string) bool {
	hash := crypto.Hasher([]byte(caption))
	if s.State.Captions.Exists(hash) {
		return true
	}
	if _, ok := s.Mutations.NewCaption[hash]; ok {
		return true
	}
	return false
}

func (s *CheckPoint) Balance(hash crypto.Hash) int {
	_, balance := s.State.Wallets.Balance(hash)
	if delta, ok := s.Mutations.DeltaWallets[hash]; ok {
		return int(balance) + delta
	}
	return int(balance)
}

func (s *CheckPoint) GetAudince(hash crypto.Hash) (bool, []byte) {
	if audience, ok := s.Mutations.NewAudiences[hash]; ok {
		return true, audience
	}
	return s.State.Audiences.GetKeys(hash)
}

func (s *CheckPoint) HasPowerOfAttorney(hash crypto.Hash) bool {
	has := s.State.PowerOfAttorney.Exists(hash)
	if !has {
		_, has = s.Mutations.GrantPower[hash]
	}
	if !has {
		return false
	}
	_, revoked := s.Mutations.RevokePower[hash]
	return !revoked
}

func (s *CheckPoint) HasAdvOffer(hash crypto.Hash) *instruction.AdvertisingOffer {
	return nil
}
