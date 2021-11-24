package instructionsnew

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
)

type Block struct {
	Parent       crypto.Hash
	Epoch        uint64
	CheckPoint   uint64
	Publisher    []byte
	PublishedAt  time.Time
	Instructions [][]byte
	Hash         crypto.Hash
	validator    *Validator
	mutations    *Mutation
}

func (b *Block) Incorporate(instruction Instruction) bool {
	if !instruction.Validate(*b.validator) {
		return false
	}
	payments := GetPayments(instruction)
	if !b.CanPay(payments) {
		return false
	}
	return true
}

func (b *Block) CanPay(payments *Payment) bool {
	for n, debitAcc := range payments.DebitAcc {
		existingBalance := b.validator.Balance(debitAcc)
		delta := b.mutations.DeltaBalance(debitAcc)
		if int(existingBalance) < int(payments.DebitValue[n])+delta {
			return false
		}
	}
	return true
}

func (b *Block) TransferPayments(payments *Payment) {
	for n, debitAcc := range payments.DebitAcc {
		if delta, ok := b.mutations.DeltaWallets[debitAcc]; ok {
			b.mutations.DeltaWallets[debitAcc] = delta - int(payments.DebitValue[n])
		} else {
			b.mutations.DeltaWallets[debitAcc] = -int(payments.DebitValue[n])
		}
	}
	for n, creditAcc := range payments.CreditAcc {
		if delta, ok := b.mutations.DeltaWallets[creditAcc]; ok {
			b.mutations.DeltaWallets[creditAcc] = delta + int(payments.CreditValue[n])
		} else {
			b.mutations.DeltaWallets[creditAcc] = int(payments.CreditValue[n])
		}
	}
}

func setNewHash(hash crypto.Hash, store map[crypto.Hash]struct{}) bool {
	if _, ok := store[hash]; ok {
		return false
	}
	store[hash] = struct{}{}
	return true
}

/*func (b *Block) SetNewHash(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.Hashes)
}
*/

func (b *Block) SetNewGrantPower(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.GrantPower)
}

func (b *Block) SetNewRevokePower(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.RevokePower)
}

func (b *Block) SetNewUseSonOffer(hash crypto.Hash, expire uint64) bool {
	return setNewHash(hash, b.mutations.UseSpnOffer)
}

func (b *Block) SetNewAdvOffer(hash crypto.Hash, offer SponsorshipOffer) bool {
	if _, ok := b.mutations.NewSpnOffer[hash]; ok {
		return false
	}
	b.mutations.NewSpnOffer[hash] = &sponsorOfferState{
		contentHash: crypto.Hasher(offer.Content),
		expire:      offer.Expiry,
	}
	return true
}

func (b *Block) SetNewMember(tokenHash crypto.Hash, captionHash crypto.Hash) bool {
	if _, ok := b.mutations.NewMembers[tokenHash]; ok {
		return false
	}
	if _, ok := b.mutations.NewMembers[captionHash]; ok {
		return false
	}
	b.mutations.NewMembers[tokenHash] = struct{}{}
	b.mutations.NewCaption[captionHash] = struct{}{}
	return true
}

func (b *Block) SetNewAudience(hash crypto.Hash, keys []byte) bool {
	if _, ok := b.mutations.NewAudiences[hash]; ok {
		return false
	}
	b.mutations.NewAudiences[hash] = keys
	return true
}
