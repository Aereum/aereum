package instructionsnew

import (
	"github.com/Aereum/aereum/core/crypto"
)

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

func (a *authoredInstruction) authorHash() crypto.Hash {
	return crypto.Hasher(a.author)
}

func (a *authoredInstruction) payments() *Payment {
	if len(a.wallet) < 0 {
		return &Payment{
			Credit: []Wallet{},
			Debit:  []Wallet{Wallet{Account: crypto.Hasher(a.wallet), FungibleTokens: a.fee}},
		}
	}
	if len(a.attorney) < 0 {
		return &Payment{
			Credit: []Wallet{},
			Debit:  []Wallet{Wallet{Account: crypto.Hasher(a.attorney), FungibleTokens: a.fee}},
		}
	}
	return &Payment{
		Credit: []Wallet{},
		Debit:  []Wallet{Wallet{Account: crypto.Hasher(a.author), FungibleTokens: a.fee}},
	}

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
		author, err = crypto.PublicKeyFromBytes(a.author)
	}
	if err != nil {
		return false
	}
	a.signature, position = ParseByteArray(data, position)
	if !author.Verify(hash[:], a.signature) {
		return false
	}
	if len(a.wallet) > 0 {
		wallet, err = crypto.PublicKeyFromBytes(a.wallet)
		if err != nil {
			return false
		}
		hash = crypto.Hasher(data[0:position])
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
		authored.wallet = a.wallet.PublicKey().ToBytes()
	}
	if a.attorney != nil {
		authored.attorney = a.attorney.PublicKey().ToBytes()
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

func (a *Author) NewUpdateInfo(details string, epoch, fee uint64) *UpdateInfo {
	update := UpdateInfo{
		authored: a.NewAuthored(epoch, fee),
		details:  details,
	}
	bulk := update.serializeBulk()
	if a.sign(update.authored, bulk, iUpdateInfo) {
		return &update
	}
	return nil
}

func (a *Author) NewGrantPowerOfAttorney(attorney []byte, epoch, fee uint64) *GrantPowerOfAttorney {
	grant := GrantPowerOfAttorney{
		authored: a.NewAuthored(epoch, fee),
		attorney: attorney,
	}
	bulk := grant.serializeBulk()
	if a.sign(grant.authored, bulk, iGrantPowerOfAttorney) {
		return &grant
	}
	return nil
}

func (a *Author) NewRevokePowerOfAttorney(attorney []byte, epoch, fee uint64) *RevokePowerOfAttorney {
	revoke := RevokePowerOfAttorney{
		authored: a.NewAuthored(epoch, fee),
		attorney: attorney,
	}
	bulk := revoke.serializeBulk()
	if a.sign(revoke.authored, bulk, iRevokePowerOfAttorney) {
		return &revoke
	}
	return nil
}

func (a *Author) NewCreateEphemeral(token []byte, expiry, epoch, fee uint64) *CreateEphemeral {
	ephemeral := CreateEphemeral{
		authored:       a.NewAuthored(epoch, fee),
		ephemeralToken: token,
		expiry:         expiry,
	}
	bulk := ephemeral.serializeBulk()
	if a.sign(ephemeral.authored, bulk, iCreateEphemeral) {
		return &ephemeral
	}
	return nil
}

func (a *Author) NewSecureChannel(tokenRange []byte, nonce uint64, encryptedNonce, content []byte, epoch, fee uint64) *SecureChannel {
	secure := SecureChannel{
		authored:       a.NewAuthored(epoch, fee),
		tokenRange:     tokenRange,
		nonce:          nonce,
		encryptedNonce: encryptedNonce,
		content:        content,
	}
	bulk := secure.serializeBulk()
	if a.sign(secure.authored, bulk, iSecureChannel) {
		return &secure
	}
	return nil
}

func (a *Author) NewCreateAudience(audience *Audience, flag byte, description string, epoch, fee uint64) *CreateAudience {
	newAudience := CreateAudience{
		authored:      a.NewAuthored(epoch, fee),
		audience:      audience.token.PublicKey().ToBytes(),
		submission:    audience.submission.PublicKey().ToBytes(),
		moderation:    audience.moderation.PublicKey().ToBytes(),
		audienceKey:   audience.SealedToken(),
		submissionKey: audience.SealedSubmission(),
		moderationKey: audience.SealedModeration(),
		flag:          flag,
		description:   description,
	}
	bulk := newAudience.serializeBulk()
	if a.sign(newAudience.authored, bulk, iCreateAudience) {
		return &newAudience
	}
	return nil
}

func (a *Author) NewJoinAudience(audience []byte, presentation string, epoch, fee uint64) *JoinAudience {
	join := JoinAudience{
		authored:     a.NewAuthored(epoch, fee),
		audience:     audience,
		presentation: presentation,
	}
	bulk := join.serializeBulk()
	if a.sign(join.authored, bulk, iJoinAudience) {
		return &join
	}
	return nil
}

func (a *Author) NewAcceptJoinAudience(audience *Audience, member crypto.PublicKey, level byte, epoch, fee uint64) *AcceptJoinAudience {
	accept := AcceptJoinAudience{
		authored: a.NewAuthored(epoch, fee),
		member:   member.ToBytes(),
		submit:   []byte{},
		moderate: []byte{},
		audience: []byte{},
	}
	accept.read, _ = member.Encrypt(audience.readCipher)
	if level > 0 {
		accept.submit, _ = member.Encrypt(audience.submitKeyCipher)
	}
	if level > 1 {
		accept.moderate, _ = member.Encrypt(audience.moderateKeyCipher)

	}
	if level > 2 {
		accept.audience, _ = member.Encrypt(audience.audienceKeyCipher)
	}
	bulk := accept.serializeBulk()
	if a.sign(accept.authored, bulk, iAcceptJoinRequest) {
		return &accept
	}
	return nil
}

func (a *Author) NewUpdateAudience(audience *Audience, readers, submiters, moderators []crypto.PublicKey, flag byte, description string, epoch, fee uint64) *UpdateAudience {
	update := UpdateAudience{
		authored:      a.NewAuthored(epoch, fee),
		audience:      audience.token.PublicKey().ToBytes(),
		submission:    audience.submission.PublicKey().ToBytes(),
		moderation:    audience.moderation.PublicKey().ToBytes(),
		audienceKey:   audience.SealedToken(),
		submissionKey: audience.SealedSubmission(),
		moderationKey: audience.SealedModeration(),
		flag:          flag,
		description:   description,
		readMembers:   audience.ReadTokenCiphers(readers),
		subMembers:    audience.SubmitTokenCiphers(submiters),
		modMembers:    audience.ModerateTokenCiphers(moderators),
	}
	bulk := update.serializeBulk()
	if a.sign(update.authored, bulk, iUpdateAudience) {
		return &update
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
		PutByteArray(authored.signature, &bytes)
		hash = crypto.Hasher(bytes)
		authored.walletSignature, err = a.wallet.Sign(hash[:])
	}
	return err == nil
}
