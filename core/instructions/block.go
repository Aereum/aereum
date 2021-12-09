package instructions

import (
	"fmt"
	"time"

	"github.com/Aereum/aereum/core/crypto"
)

type Block struct {
	Epoch         uint64
	Parent        crypto.Hash
	CheckPoint    uint64
	Publisher     []byte
	PublishedAt   time.Time
	Instructions  [][]byte
	Hash          crypto.Hash
	FeesCollected uint64
	Signature     []byte
	validator     *Validator
	mutations     *Mutation
}

func NewBlock(parent crypto.Hash, checkpoint, epoch uint64, publisher []byte, validator *Validator) *Block {
	return &Block{
		Parent:       parent,
		Epoch:        epoch,
		CheckPoint:   checkpoint,
		Publisher:    publisher,
		Instructions: make([][]byte, 0),
		validator:    validator,
		mutations:    NewMutation(),
	}
}

func (b *Block) Incorporate(instruction Instruction) bool {
	payments := instruction.Payments()
	if !b.CanPay(payments) {
		if instruction.Kind() == iCreateAudience {
			fmt.Println("--------------------")
		}
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
			// fmt.Println(debit.Account, debit.FungibleTokens)
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

func (b *Block) Sign(token crypto.PrivateKey) {
	hashed := crypto.Hasher(b.serializeWithoutSignature())
	b.Signature, _ = token.Sign(hashed[:])
}

func (b *Block) Serialize() []byte {
	bytes := b.serializeWithoutSignature()
	PutByteArray(b.Signature, &bytes)
	return bytes
}

func (b *Block) serializeWithoutSignature() []byte {
	bytes := make([]byte, 0)
	PutUint64(b.Epoch, &bytes)
	PutByteArray(b.Parent[:], &bytes)
	PutByteArray(b.Publisher[:], &bytes)
	PutTime(b.PublishedAt, &bytes)
	PutUint16(uint16(len(b.Instructions)), &bytes)
	for _, instruction := range b.Instructions {
		PutByteArray(instruction, &bytes)
	}
	PutByteArray(b.Hash[:], &bytes)
	PutUint64(b.FeesCollected, &bytes)
	return bytes
}

func ParseBlock(data []byte) *Block {
	position := 0
	block := Block{}
	block.Epoch, position = ParseUint64(data, position)
	block.Parent, position = ParseHash(data, position)
	block.Publisher, position = ParseByteArray(data, position)
	block.PublishedAt, position = ParseTime(data, position)
	block.Instructions, position = ParseByteArrayArray(data, position)
	block.Hash, position = ParseHash(data, position)
	block.FeesCollected, position = ParseUint64(data, position)
	hashed := crypto.Hasher(data[0:position])
	block.Signature, _ = ParseByteArray(data, position)
	pubkey, err := crypto.PublicKeyFromBytes(block.Publisher)
	if err != nil {
		return nil
	}
	if !pubkey.Verify(hashed[:], block.Signature) {
		fmt.Println("wrong signature")
		return nil
	}
	block.mutations = NewMutation()
	return &block
}

func (b *Block) SetValidator(validator *Validator) {
	b.validator = validator
}

func GetBlockEpoch(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	epoch, _ := ParseUint64(data, 0)
	return epoch
}
