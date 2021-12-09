package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type BulkSerializer interface {
	serializeBulk() []byte
	InstructionType() byte
}

type authoredInstruction struct {
	epoch           uint64
	author          []byte
	wallet          []byte
	fee             uint64
	attorney        []byte
	signature       []byte
	walletSignature []byte
}

func (a *authoredInstruction) authorHash() crypto.Hash {
	return crypto.Hasher(a.author)
}

func (a *authoredInstruction) payments() *Payment {
	if len(a.wallet) > 0 {
		return NewPayment(crypto.Hasher(a.wallet), a.fee)
	}
	if len(a.attorney) > 0 {
		return NewPayment(crypto.Hasher(a.attorney), a.fee)
	}
	return NewPayment(crypto.Hasher(a.author), a.fee)
}

func (a *authoredInstruction) Clone() *authoredInstruction {
	clone := &authoredInstruction{
		epoch: a.epoch,
		fee:   a.fee,
	}
	copy(clone.author, a.author)
	copy(clone.wallet, a.wallet)
	copy(clone.attorney, a.attorney)
	return clone
}

func (a *authoredInstruction) serializeWithoutSignature(instType byte, bulk []byte) []byte {
	bytes := []byte{0, instType}
	util.PutUint64(a.epoch, &bytes)
	util.PutByteArray(a.author, &bytes)
	bytes = append(bytes, bulk...)
	util.PutByteArray(a.wallet, &bytes)
	util.PutUint64(a.fee, &bytes)
	util.PutByteArray(a.attorney, &bytes)
	return bytes
}

func (a *authoredInstruction) serialize(instType byte, bulk []byte) []byte {
	bytes := a.serializeWithoutSignature(instType, bulk)
	util.PutByteArray(a.signature, &bytes)
	util.PutByteArray(a.walletSignature, &bytes)
	return bytes
}

func (a *authoredInstruction) parseHead(data []byte) int {
	position := 2
	a.epoch, position = util.ParseUint64(data, position)
	a.author, position = util.ParseByteArray(data, position)
	return position
}

func (a *authoredInstruction) parseTail(data []byte, position int) bool {
	a.wallet, position = util.ParseByteArray(data, position)
	a.fee, position = util.ParseUint64(data, position)
	a.attorney, position = util.ParseByteArray(data, position)
	hash := crypto.Hasher(data[0:position])
	var author, wallet crypto.PublicKey
	var err error
	if len(a.attorney) > 0 {
		author, err = crypto.PublicKeyFromBytes(a.attorney)
	} else {
		author, err = crypto.PublicKeyFromBytes(a.author)
	}
	if err != nil {
		return false
	}
	a.signature, position = util.ParseByteArray(data, position)
	if !author.Verify(hash[:], a.signature) {
		return false
	}
	if len(a.wallet) > 0 {
		wallet, err = crypto.PublicKeyFromBytes(a.wallet)
		if err != nil {
			return false
		}
		hash = crypto.Hasher(data[0:position])
		a.walletSignature, position = util.ParseByteArray(data, position)
		if position != len(data) {
			return false
		}
		return wallet.Verify(hash[:], a.walletSignature)
	} else {
		return position == len(data)
	}
}

func NewAuthored(epoch, fee uint64, author crypto.PrivateKey, wallet *crypto.PrivateKey, attorney *crypto.PrivateKey) *authoredInstruction {
	authored := &authoredInstruction{
		epoch:  epoch,
		author: author.PublicKey().ToBytes(),
		fee:    fee,
	}
	if wallet != nil {
		authored.wallet = (*wallet).PublicKey().ToBytes()
	} else {
		authored.wallet = []byte{}
	}
	if attorney != nil {
		authored.attorney = (*wallet).PublicKey().ToBytes()
	} else {
		authored.attorney = []byte{}
	}
	return authored
}
