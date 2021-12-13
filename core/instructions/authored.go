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
	author          crypto.Token
	wallet          crypto.Token
	fee             uint64
	attorney        crypto.Token
	signature       crypto.Signature
	walletSignature crypto.Signature
}

func (a *authoredInstruction) authorHash() crypto.Hash {
	return crypto.Hasher(a.author[:])
}

func (a *authoredInstruction) payments() *Payment {
	if len(a.wallet) > 0 {
		return NewPayment(crypto.Hasher(a.wallet[:]), a.fee)
	}
	if len(a.attorney) > 0 {
		return NewPayment(crypto.Hasher(a.attorney[:]), a.fee)
	}
	return NewPayment(crypto.Hasher(a.author[:]), a.fee)
}

func (a *authoredInstruction) Clone() *authoredInstruction {
	clone := &authoredInstruction{
		epoch: a.epoch,
		fee:   a.fee,
	}
	clone.author = a.author
	clone.wallet = a.wallet
	clone.attorney = a.attorney
	return clone
}

func (a *authoredInstruction) serializeWithoutSignature(instType byte, bulk []byte) []byte {
	bytes := []byte{0, instType}
	util.PutUint64(a.epoch, &bytes)
	util.PutToken(a.author, &bytes)
	bytes = append(bytes, bulk...)
	util.PutToken(a.wallet, &bytes)
	util.PutUint64(a.fee, &bytes)
	util.PutToken(a.attorney, &bytes)
	return bytes
}

func (a *authoredInstruction) serialize(instType byte, bulk []byte) []byte {
	bytes := a.serializeWithoutSignature(instType, bulk)
	util.PutSignature(a.signature, &bytes)
	util.PutSignature(a.walletSignature, &bytes)
	return bytes
}

func (a *authoredInstruction) parseHead(data []byte) int {
	position := 2
	a.epoch, position = util.ParseUint64(data, position)
	a.author, position = util.ParseToken(data, position)
	return position
}

func (a *authoredInstruction) parseTail(data []byte, position int) bool {
	a.wallet, position = util.ParseToken(data, position)
	a.fee, position = util.ParseUint64(data, position)
	a.attorney, position = util.ParseToken(data, position)
	var token, wallet crypto.Token
	if a.attorney != crypto.ZeroToken {
		token = a.attorney
	} else {
		token = a.author
	}
	msg := data[0:position]
	a.signature, position = util.ParseSignature(data, position)
	if !token.Verify(msg, a.signature) {
		return false
	}
	if a.wallet != crypto.ZeroToken {
		wallet = a.wallet
		msg = data[0:position]
		a.walletSignature, position = util.ParseSignature(data, position)
		if position != len(data) {
			return false
		}
		return wallet.Verify(msg, a.walletSignature)
	} else {
		return position == len(data)
	}
}

func NewAuthored(epoch, fee uint64, author crypto.PrivateKey, wallet *crypto.PrivateKey, attorney *crypto.PrivateKey) *authoredInstruction {
	authored := &authoredInstruction{
		epoch:  epoch,
		author: author.PublicKey(),
		fee:    fee,
	}
	if wallet != nil {
		authored.wallet = (*wallet).PublicKey()
	} else {
		authored.wallet = crypto.ZeroToken
	}
	if attorney != nil {
		authored.attorney = (*wallet).PublicKey()
	} else {
		authored.attorney = crypto.ZeroToken
	}
	return authored
}
