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

// Validator consists of a state and permanent mutations not incorporated into
// the state. It provides an interface to check the state that would result
// in the incorporation of mutations to the original state without actually
// incorporating.
type MutatingState struct {
	State     *State
	Mutations *mutation
}

// Balance returns the balance of fungible tokens associated to the hash.
// It returns zero if the hash is not found.
func (c *MutatingState) balance(hash crypto.Hash) uint64 {
	_, balance := c.State.Wallets.Balance(hash)
	if c.Mutations == nil {
		return balance
	}
	delta := c.Mutations.DeltaBalance(hash)
	if delta < 0 {
		balance = balance - uint64(-delta)
	} else {
		balance = balance + uint64(delta)
	}
	return balance
}

// PowerOfAttorney checks if an attorney can sign on behalf of an author.
func (c *MutatingState) powerOfAttorney(hash crypto.Hash) bool {
	if c.Mutations != nil {
		if c.Mutations.HasRevokePower(hash) {
			return false
		}
		if c.Mutations.HasGrantPower(hash) {
			return true
		}
	}
	return c.State.PowerOfAttorney.Exists(hash)
}

// SponsorshipOffer returns the expire epoch of an SponsorshipOffer. It returns
// zero if no offer is found of the given hash
func (c *MutatingState) sponsorshipOffer(hash crypto.Hash) uint64 {
	if c.Mutations != nil {
		if c.Mutations.HasUsedSponsorOffer(hash) {
			return 0
		}
		if offer := c.Mutations.GetSponsorOffer(hash); !offer {
			return 0
		}
	}
	expire := c.State.SponsorOffers.Exists(hash)
	return expire
}

// HasMeber returns the existance of a member.
func (c *MutatingState) hasMember(hash crypto.Hash) bool {
	if c.Mutations != nil && c.Mutations.HasMember(hash) {
		return true
	}
	return c.State.Members.Exists(hash)
}

// HasGrantedSponsor returns the existence and the hash of the grantee +
// audience.
func (c *MutatingState) hasGrantedSponser(hash crypto.Hash) (bool, crypto.Hash) {
	if c.Mutations != nil {
		if ok, contentHash := c.Mutations.HasGrantedSponsorship(hash); ok {
			return true, contentHash
		}
	}
	ok, contentHash := c.State.SponsorGranted.GetContentHash(hash)
	return ok, crypto.Hasher(contentHash)
}

// HasCaption returns the existence of the caption
func (c *MutatingState) hasCaption(hash crypto.Hash) bool {
	if c.Mutations != nil && c.Mutations.HasCaption(hash) {
		return true
	}
	return c.State.Captions.Exists(hash)
}

// GetAudienceKeys returns the audience keys
func (c *MutatingState) getAudienceKeys(hash crypto.Hash) []byte {
	if c.Mutations != nil {
		if audience := c.Mutations.GetAudience(hash); audience != nil {
			return audience
		}
	}
	ok, keys := c.State.Audiences.GetKeys(hash)
	if !ok {
		return nil
	}
	return keys
}

// GetEphemeralExpire returns the expire epoch of the associated ephemeral token
// It returns zero if the token is not found.
func (c *MutatingState) getEphemeralExpire(hash crypto.Hash) (bool, uint64) {
	if c.Mutations != nil {
		if ok, expire := c.Mutations.HasEphemeral(hash); ok {
			return true, expire
		}
	}
	expire := c.State.EphemeralTokens.Exists(hash)
	return expire > 0, expire
}
