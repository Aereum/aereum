package instructionsnew

import "github.com/Aereum/aereum/core/crypto"

type Validator struct {
	State     *State
	Mutations *Mutation
}

func (c *Validator) Balance(hash crypto.Hash) uint64 {
	_, balance := c.State.Wallets.Balance(hash)
	delta := c.Mutations.DeltaBalance(hash)
	if delta < 0 {
		balance = balance - uint64(-delta)
	} else {
		balance = balance + uint64(delta)
	}
	return balance
}

func (c *Validator) PowerOfAttorney(hash crypto.Hash) bool {
	if c.Mutations.HasRevokePower(hash) {
		return false
	}
	if c.Mutations.HasGrantPower(hash) {
		return true
	}
	return c.State.PowerOfAttorney.Exists(hash)
}

func (c *Validator) SponsorshipOffer(hash crypto.Hash) uint64 {
	if c.Mutations.HasUsedSponsorOffer(hash) {
		return 0
	}
	if offer := c.Mutations.GetSponsorOffer(hash); !offer {
		return 0
	}
	expire := c.State.SponsorOffers.Exists(hash)
	return expire
}

func (c *Validator) HasMember(hash crypto.Hash) bool {
	if c.Mutations.HasMember(hash) {
		return true
	}
	return c.State.Members.Exists(hash)
}

func (c *Validator) HasCaption(hash crypto.Hash) bool {
	if c.Mutations.HasCaption(hash) {
		return true
	}
	return c.State.Captions.Exists(hash)
}

func (c *Validator) GetAudienceKeys(hash crypto.Hash) []byte {
	if audience := c.Mutations.GetAudience(hash); audience != nil {
		return audience
	}
	ok, keys := c.State.Audiences.GetKeys(hash)
	if !ok {
		return nil
	}
	return keys
}

func (c *Validator) GetEphemeralExpire(hash crypto.Hash) (bool, uint64) {
	if ok, expire := c.Mutations.HasEphemeral(hash); ok {
		return true, expire
	}
	expire := c.State.EphemeralTokens.Exists(hash)
	return expire > 0, expire
}
