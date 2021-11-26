package instructions

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
)

type Block struct {
	Parent        crypto.Hash
	Epoch         uint64
	CheckPoint    uint64
	Publisher     []byte
	PublishedAt   time.Time
	Instructions  [][]byte
	Hash          crypto.Hash
	Signature     []byte
	FeesCollected uint64
	validator     *Validator
	mutations     *Mutation
}

func (b *Block) Incorporate(instruction Instruction) bool {
	payments := instruction.Payments()
	if !b.CanPay(payments) {
		return false
	}
	if !instruction.Validate(b) {
		return false
	}
	b.TransferPayments(payments)
	b.Instructions = append(b.Instructions, instruction.Serialize())
	return true
}

func (b *Block) CanPay(payments *Payment) bool {
	for _, debit := range payments.Debit {
		existingBalance := b.validator.Balance(debit.Account)
		delta := b.mutations.DeltaBalance(debit.Account)
		if int(existingBalance) < int(debit.FungibleTokens)+delta {
			return false
		}
	}
	return true
}

func (b *Block) TransferPayments(payments *Payment) {
	for _, debit := range payments.Debit {
		if delta, ok := b.mutations.DeltaWallets[debit.Account]; ok {
			b.mutations.DeltaWallets[debit.Account] = delta - int(debit.FungibleTokens)
		} else {
			b.mutations.DeltaWallets[debit.Account] = -int(debit.FungibleTokens)
		}
	}
	for _, credit := range payments.Credit {
		if delta, ok := b.mutations.DeltaWallets[credit.Account]; ok {
			b.mutations.DeltaWallets[credit.Account] = delta + int(credit.FungibleTokens)
		} else {
			b.mutations.DeltaWallets[credit.Account] = int(credit.FungibleTokens)
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

func (b *Block) SetNewGrantPower(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.GrantPower)
}

func (b *Block) SetNewRevokePower(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.RevokePower)
}

func (b *Block) SetNewUseSonOffer(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.UseSpnOffer)
}

func (b *Block) SetNewSpnOffer(hash crypto.Hash, expire uint64) bool {
	if _, ok := b.mutations.NewSpnOffer[hash]; ok {
		return false
	}
	b.mutations.NewSpnOffer[hash] = expire
	return true
}

func (b *Block) SetPublishSponsor(hash crypto.Hash) bool {
	return setNewHash(hash, b.mutations.PublishSpn)
}

func (b *Block) SetNewEphemeralToken(hash crypto.Hash, expire uint64) bool {
	if _, ok := b.mutations.NewEphemeral[hash]; ok {
		return false
	}
	b.mutations.NewEphemeral[hash] = expire
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

func (b *Block) UpdateAudience(hash crypto.Hash, keys []byte) bool {
	if _, ok := b.mutations.UpdAudiences[hash]; ok {
		return false
	}
	b.mutations.UpdAudiences[hash] = keys
	return true
}
