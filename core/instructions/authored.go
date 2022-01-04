package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type BulkSerializer interface {
	serializeBulk() []byte
	InstructionType() byte
}

type AuthoredInstruction struct {
	epoch           uint64
	Author          crypto.Token
	Wallet          crypto.Token
	Fee             uint64
	Attorney        crypto.Token
	signature       crypto.Signature
	walletSignature crypto.Signature
}

func (a *AuthoredInstruction) authorHash() crypto.Hash {
	return crypto.Hasher(a.Author[:])
}

func (a *AuthoredInstruction) payments() *Payment {
	if len(a.Wallet) > 0 {
		return NewPayment(crypto.Hasher(a.Wallet[:]), a.Fee)
	}
	if len(a.Attorney) > 0 {
		return NewPayment(crypto.Hasher(a.Attorney[:]), a.Fee)
	}
	return NewPayment(crypto.Hasher(a.Author[:]), a.Fee)
}

func (a *AuthoredInstruction) Clone() *AuthoredInstruction {
	clone := &AuthoredInstruction{
		epoch: a.epoch,
		Fee:   a.Fee,
	}
	clone.Author = a.Author
	clone.Wallet = a.Wallet
	clone.Attorney = a.Attorney
	return clone
}

func (a *AuthoredInstruction) serializeWithoutSignature(instType byte, bulk []byte) []byte {
	bytes := []byte{0, instType}
	util.PutUint64(a.epoch, &bytes)
	util.PutToken(a.Author, &bytes)
	bytes = append(bytes, bulk...)
	util.PutToken(a.Wallet, &bytes)
	util.PutUint64(a.Fee, &bytes)
	util.PutToken(a.Attorney, &bytes)
	return bytes
}

func (a *AuthoredInstruction) serialize(instType byte, bulk []byte) []byte {
	bytes := a.serializeWithoutSignature(instType, bulk)
	util.PutSignature(a.signature, &bytes)
	util.PutSignature(a.walletSignature, &bytes)
	return bytes
}

func (a *AuthoredInstruction) parseHead(data []byte) int {
	position := 2
	a.epoch, position = util.ParseUint64(data, position)
	a.Author, position = util.ParseToken(data, position)
	return position
}

func (a *AuthoredInstruction) parseTail(data []byte, position int) bool {
	a.Wallet, position = util.ParseToken(data, position)
	a.Fee, position = util.ParseUint64(data, position)
	a.Attorney, position = util.ParseToken(data, position)
	var token, wallet crypto.Token
	if a.Attorney != crypto.ZeroToken {
		token = a.Attorney
	} else {
		token = a.Author
	}
	msg := data[0:position]
	a.signature, position = util.ParseSignature(data, position)
	if !token.Verify(msg, a.signature) {
		return false
	}
	if a.Wallet != crypto.ZeroToken {
		wallet = a.Wallet
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

func NewAuthored(epoch, fee uint64, author crypto.PrivateKey, wallet *crypto.PrivateKey, attorney *crypto.PrivateKey) *AuthoredInstruction {
	authored := &AuthoredInstruction{
		epoch:  epoch,
		Author: author.PublicKey(),
		Fee:    fee,
	}
	if wallet != nil {
		authored.Wallet = (*wallet).PublicKey()
	} else {
		authored.Wallet = crypto.ZeroToken
	}
	if attorney != nil {
		authored.Attorney = (*wallet).PublicKey()
	} else {
		authored.Attorney = crypto.ZeroToken
	}
	return authored
}
