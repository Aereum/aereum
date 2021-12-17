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

import (
	"fmt"
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/store"
	"github.com/Aereum/aereum/core/util"
)

type Block struct {
	epoch         uint64
	Parent        crypto.Hash
	CheckPoint    uint64
	Publisher     crypto.Token
	PublishedAt   time.Time
	Instructions  [][]byte
	Hash          crypto.Hash
	FeesCollected uint64
	Signature     crypto.Signature
	validator     *MutatingState
	mutations     *mutation
}

func NewBlock(parent crypto.Hash, checkpoint, epoch uint64, publisher crypto.Token, validator *MutatingState) *Block {
	return &Block{
		Parent:       parent,
		epoch:        epoch,
		CheckPoint:   checkpoint,
		Publisher:    publisher,
		Instructions: make([][]byte, 0),
		validator:    validator,
		mutations:    NewMutation(),
	}
}

func (b *Block) Incorporate(instruction instructions.Instruction) bool {
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

func (b *Block) CanPay(payments *instructions.Payment) bool {
	for _, debit := range payments.Debit {
		existingBalance := b.validator.balance(debit.Account)
		delta := b.mutations.DeltaBalance(debit.Account)
		if int(existingBalance) < int(debit.FungibleTokens)+delta {
			return false
		}
	}
	return true
}

func (b *Block) TransferPayments(payments *instructions.Payment) {
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

func (b *Block) SetNewUseSpnOffer(hash crypto.Hash) bool {
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

func (b *Block) SetNewAudience(hash crypto.Hash, stage store.StageKeys) bool {
	if _, ok := b.mutations.NewStages[hash]; ok {
		return false
	}
	b.mutations.NewStages[hash] = stage
	return true
}

func (b *Block) UpdateAudience(hash crypto.Hash, stage store.StageKeys) bool {
	if _, ok := b.mutations.StageUpdate[hash]; ok {
		return false
	}
	b.mutations.StageUpdate[hash] = stage
	return true
}

func (b *Block) PowerOfAttorney(hash crypto.Hash) bool {
	return b.validator.powerOfAttorney(hash)
}

func (b *Block) SponsorshipOffer(hash crypto.Hash) uint64 {
	return b.validator.sponsorshipOffer(hash)
}

func (b *Block) HasMember(hash crypto.Hash) bool {
	return b.validator.hasMember(hash)
}

func (b *Block) HasCaption(hash crypto.Hash) bool {
	return b.validator.hasCaption(hash)
}

func (b *Block) HasGrantedSponser(hash crypto.Hash) (bool, crypto.Hash) {
	return b.validator.hasGrantedSponser(hash)
}

func (b *Block) GetAudienceKeys(hash crypto.Hash) *store.StageKeys {
	return b.validator.getAudienceKeys(hash)
}

func (b *Block) GetEphemeralExpire(hash crypto.Hash) (bool, uint64) {
	return b.validator.getEphemeralExpire(hash)
}

func (b *Block) Balance(hash crypto.Hash) uint64 {
	return b.validator.balance(hash)
}

func (b *Block) AddFeeCollected(value uint64) {
	b.FeesCollected += value
}

func (b *Block) Epoch() uint64 {
	return b.epoch
}

func (b *Block) Sign(token crypto.PrivateKey) {
	b.Signature = token.Sign(b.serializeWithoutSignature())
}

func (b *Block) Serialize() []byte {
	bytes := b.serializeWithoutSignature()
	util.PutSignature(b.Signature, &bytes)
	return bytes
}

func (b *Block) serializeWithoutSignature() []byte {
	bytes := make([]byte, 0)
	util.PutUint64(b.epoch, &bytes)
	util.PutByteArray(b.Parent[:], &bytes)
	util.PutUint64(b.CheckPoint, &bytes)
	util.PutByteArray(b.Publisher[:], &bytes)
	util.PutTime(b.PublishedAt, &bytes)
	util.PutUint16(uint16(len(b.Instructions)), &bytes)
	for _, instruction := range b.Instructions {
		util.PutByteArray(instruction, &bytes)
	}
	util.PutByteArray(b.Hash[:], &bytes)
	util.PutUint64(b.FeesCollected, &bytes)
	return bytes
}

func ParseBlock(data []byte) *Block {
	position := 0
	block := Block{}
	block.epoch, position = util.ParseUint64(data, position)
	block.Parent, position = util.ParseHash(data, position)
	block.CheckPoint, position = util.ParseUint64(data, position)
	block.Publisher, position = util.ParseToken(data, position)
	block.PublishedAt, position = util.ParseTime(data, position)
	block.Instructions, position = util.ParseByteArrayArray(data, position)
	block.Hash, position = util.ParseHash(data, position)
	block.FeesCollected, position = util.ParseUint64(data, position)
	msg := data[0:position]
	block.Signature, _ = util.ParseSignature(data, position)
	if !block.Publisher.Verify(msg, block.Signature) {
		fmt.Println("wrong signature")
		return nil
	}
	block.mutations = NewMutation()
	return &block
}

func (b *Block) SetValidator(validator *MutatingState) {
	b.validator = validator
}

func GetBlockEpoch(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	epoch, _ := util.ParseUint64(data, 0)
	return epoch
}

func (b *Block) JSONSimple() string {
	bulk := &util.JSONBuilder{}
	bulk.PutUint64("epoch", b.epoch)
	bulk.PutHex("parent", b.Parent[:])
	bulk.PutUint64("checkpoint", b.CheckPoint)
	bulk.PutHex("publisher", b.Publisher[:])
	bulk.PutTime("publishedAt", b.PublishedAt)
	bulk.PutUint64("instructionsCount", uint64(len(b.Instructions)))
	bulk.PutHex("hash", b.Parent[:])
	bulk.PutUint64("feesCollectes", b.FeesCollected)
	bulk.PutBase64("signature", b.Signature[:])
	return bulk.ToString()
}
