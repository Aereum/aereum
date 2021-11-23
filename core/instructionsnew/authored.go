package instructionsnew

import "github.com/Aereum/aereum/core/crypto"

type BulkSerializer interface {
	serializeBulk() []byte
	InstructionType() byte
}

type Author struct {
	token    *crypto.PrivateKey
	wallet   *crypto.PrivateKey
	attorney *crypto.PrivateKey
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
	PutUint64(a.epoch, &bytes)
	PutByteArray(a.author, &bytes)
	bytes = append(bytes, bulk...)
	PutByteArray(a.wallet, &bytes)
	PutUint64(a.fee, &bytes)
	PutByteArray(a.attorney, &bytes)
	return bytes
}

func (a *authoredInstruction) serialize(instType byte, bulk []byte) []byte {
	bytes := a.serializeWithoutSignature(instType, bulk)
	PutByteArray(a.signature, &bytes)
	PutByteArray(a.walletSignature, &bytes)
	return bytes
}

func (a *authoredInstruction) parseHead(data []byte) int {
	position := 2
	a.epoch, position = ParseUint64(data, position)
	a.author, position = ParseByteArray(data, position)
	return position
}

func (a *authoredInstruction) parseTail(data []byte, position int) bool {
	a.wallet, position = ParseByteArray(data, position)
	a.fee, position = ParseUint64(data, position)
	a.attorney, position = ParseByteArray(data, position)
	hash := crypto.Hasher(data[0:position])
	var author, wallet crypto.PublicKey
	var err error
	if len(a.attorney) > 0 {
		author, err = crypto.PublicKeyFromBytes(a.attorney)
	} else {
		author, err = crypto.PublicKeyFromBytes(a.attorney)
	}
	if err != nil {
		return false
	}
	a.signature, position = ParseByteArray(data, position)
	if !author.Verify(hash[:], a.signature) {
		return false
	}
	hash = crypto.Hasher(data[0:position])
	if len(a.wallet) > 0 {
		wallet, err = crypto.PublicKeyFromBytes(a.wallet)
		if err != nil {
			return false
		}
		a.walletSignature, position = ParseByteArray(data, position)
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

func (a *Author) NewAuthored(epoch, fee uint64) *authoredInstruction {
	if a.token == nil {
		return nil
	}
	authored := authoredInstruction{
		epoch:           epoch,
		author:          a.token.PublicKey().ToBytes(),
		wallet:          []byte{},
		fee:             fee,
		attorney:        []byte{},
		signature:       []byte{},
		walletSignature: []byte{},
	}
	if a.wallet != nil {
		authored.wallet = a.wallet.ToBytes()
	}
	if a.attorney != nil {
		authored.attorney = a.attorney.ToBytes()
	}
	return &authored
}

func (a *Author) NewJoinNetwork(caption string, details string, epoch, fee uint64) *JoinNetwork {
	join := JoinNetwork{
		authored: a.NewAuthored(epoch, fee),
		caption:  caption,
		details:  details,
	}
	bulk := join.serializeBulk()
	if a.sign(join.authored, bulk, iJoinNetwork) {
		return &join
	}
	return nil
}

func (a *Author) sign(authored *authoredInstruction, bulk []byte, insType byte) bool {
	bytes := authored.serializeWithoutSignature(insType, bulk)
	hash := crypto.Hasher(bytes)
	var err error
	if a.attorney != nil {
		authored.signature, err = a.attorney.Sign(hash[:])
	} else {
		authored.signature, err = a.token.Sign(hash[:])
	}
	if a.wallet != nil {
		bytes = append(bytes, authored.signature...)
		hash = crypto.Hasher(bytes)
		authored.walletSignature, err = a.wallet.Sign(hash[:])
	}
	return err == nil
}
