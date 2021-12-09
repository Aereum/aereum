package instructions

import (
	"bytes"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type BulkSerializer interface {
	serializeBulk() []byte
	InstructionType() byte
}

type Author struct {
	Token    *crypto.PrivateKey
	Wallet   *crypto.PrivateKey
	Attorney *crypto.PrivateKey
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

func NewAuthor(token, wallet, attorney *crypto.PrivateKey) *Author {
	author := &Author{
		Token:    token,
		Wallet:   wallet,
		Attorney: attorney,
	}
	return author
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

func (a *Author) NewAuthored(epoch, fee uint64) *authoredInstruction {
	if a.Token == nil {
		return nil
	}
	authored := authoredInstruction{
		epoch:           epoch,
		author:          a.Token.PublicKey().ToBytes(),
		wallet:          []byte{},
		fee:             fee,
		attorney:        []byte{},
		signature:       []byte{},
		walletSignature: []byte{},
	}
	if a.Wallet != nil {
		authored.wallet = a.Wallet.PublicKey().ToBytes()
	}
	if a.Attorney != nil {
		authored.attorney = a.Attorney.PublicKey().ToBytes()
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

func (a *Author) NewJoinNetworkThirdParty(token []byte, caption string, details string, epoch, fee uint64) *JoinNetwork {
	authored := authoredInstruction{
		epoch:    epoch,
		author:   token,
		fee:      fee,
		attorney: []byte{},
		wallet:   []byte{},
	}
	if a.Attorney != nil {
		authored.attorney = a.Attorney.PublicKey().ToBytes()
	}
	if a.Wallet != nil {
		authored.wallet = a.Wallet.PublicKey().ToBytes()
	} else if a.Attorney != nil {
		authored.wallet = a.Attorney.PublicKey().ToBytes()
	} else {
		authored.wallet = a.Token.PublicKey().ToBytes()
	}
	join := JoinNetwork{
		authored: &authored,
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
		audience:      audience.Token.PublicKey().ToBytes(),
		submission:    audience.Submission.PublicKey().ToBytes(),
		moderation:    audience.Moderation.PublicKey().ToBytes(),
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
		audience: audience.Token.PublicKey().ToBytes(),
		member:   member.ToBytes(),
		read:     []byte{},
		submit:   []byte{},
		moderate: []byte{},
	}
	accept.read, _ = member.Encrypt(audience.AudienceKeyCipher)
	if level > 0 {
		accept.submit, _ = member.Encrypt(audience.SubmitKeyCipher)
	}
	if level > 1 {
		accept.moderate, _ = member.Encrypt(audience.ModerateKeyCipher)

	}
	modbulk := accept.serializeModBulk()
	var sign []byte
	var err error
	sign, err = audience.Moderation.Sign(modbulk)
	if err != nil {
		return nil
	}
	accept.modSignature = sign
	bulk := accept.serializeBulk()
	if a.sign(accept.authored, bulk, iAcceptJoinRequest) {
		return &accept
	}
	return nil
}

func (a *Author) NewUpdateAudience(audience *Audience, readers, submiters, moderators []crypto.PublicKey, flag byte, description string, epoch, fee uint64) *UpdateAudience {
	update := UpdateAudience{
		authored:      a.NewAuthored(epoch, fee),
		audience:      audience.Token.PublicKey().ToBytes(),
		submission:    audience.Submission.PublicKey().ToBytes(),
		moderation:    audience.Moderation.PublicKey().ToBytes(),
		submissionKey: audience.SealedSubmission(),
		moderationKey: audience.SealedModeration(),
		flag:          flag,
		description:   description,
		readMembers:   audience.ReadTokenCiphers(readers),
		subMembers:    audience.SubmitTokenCiphers(submiters),
		modMembers:    audience.ModerateTokenCiphers(moderators),
	}
	audBulk := update.serializeAudBulk()
	var sign []byte
	var err error
	sign, err = audience.Token.Sign(audBulk)
	if err != nil {
		return nil
	}
	update.audSignature = sign
	bulk := update.serializeBulk()
	if a.sign(update.authored, bulk, iUpdateAudience) {
		return &update
	}
	return nil
}

func (a *Author) ModerateContent(audience *Audience, content *Content, epoch, fee uint64) *Content {
	if audience == nil || audience.Moderation == nil {
		return nil
	}
	if !bytes.Equal(audience.Token.ToBytes(), content.audience) {
		return nil
	}
	newContent := &Content{
		epoch:        epoch,
		published:    content.epoch,
		author:       content.author,
		audience:     content.audience,
		contentType:  content.contentType,
		content:      content.content,
		sponsored:    content.sponsored,
		encrypted:    content.encrypted,
		subSignature: content.subSignature,
		moderator:    a.Token.ToBytes(),
		attorney:     []byte{},
		wallet:       []byte{},
		fee:          fee,
	}
	hash := crypto.Hasher(newContent.serializeModBulk())
	var err error
	newContent.modSignature, err = audience.Moderation.Sign(hash[:])
	if err != nil {
		return nil
	}
	if a.Attorney != nil {
		newContent.attorney = a.Attorney.ToBytes()
		hash = crypto.Hasher(newContent.serializeSignBulk())
		newContent.signature, err = a.Attorney.Sign(hash[:])
	} else {
		hash = crypto.Hasher(newContent.serializeSignBulk())
		newContent.signature, err = a.Token.Sign(hash[:])
	}
	if err != nil {
		return nil
	}
	if a.Wallet != nil {
		newContent.wallet = a.Wallet.ToBytes()
		hash = crypto.Hasher(newContent.serializeWalletBulk())
		newContent.walletSignature, err = a.Attorney.Sign(hash[:])
	} else {
		hash = crypto.Hasher(newContent.serializeWalletBulk())
		newContent.walletSignature, err = a.Token.Sign(hash[:])
	}
	if err != nil {
		return nil
	}
	return newContent
}

func (a *Author) NewContent(audience *Audience, contentType string, message []byte, hash, encrypted bool, epoch, fee uint64) *Content {
	if audience == nil {
		return nil
	}
	content := &Content{
		epoch:        epoch,
		published:    epoch,
		author:       a.Token.PublicKey().ToBytes(),
		audience:     audience.Token.PublicKey().ToBytes(),
		contentType:  contentType,
		hash:         []byte{},
		sponsored:    false,
		encrypted:    encrypted,
		attorney:     []byte{},
		moderator:    []byte{},
		modSignature: []byte{},
		wallet:       []byte{},
		fee:          fee,
	}
	if a.Attorney != nil {
		content.attorney = a.Attorney.PublicKey().ToBytes()
	}
	if a.Wallet != nil {
		content.wallet = a.Wallet.PublicKey().ToBytes()
	}
	if encrypted {
		cipher := crypto.CipherFromKey(audience.AudienceKeyCipher)
		content.content = cipher.Seal(message)
	} else {
		content.content = message
	}
	if hash {
		hashed := crypto.Hasher(message)
		content.hash = hashed[:]
	}
	subBulk := content.serializeSubBulk()
	hashed := crypto.Hasher(subBulk[10:])
	var sign []byte
	var err error
	sign, err = audience.Submission.Sign(hashed[:])
	if err != nil {
		return nil
	}
	content.subSignature = sign
	util.PutByteArray(content.subSignature, &subBulk)
	if audience.Moderation != nil {
		content.moderator = a.Token.PublicKey().ToBytes()
		hashed = crypto.Hasher(content.serializeModBulk())
		content.modSignature, err = audience.Moderation.Sign(hashed[:])
		if err != nil {
			return nil
		}
	}
	hashed = crypto.Hasher(content.serializeSignBulk())
	if a.Attorney != nil {
		content.signature, err = a.Attorney.Sign(hashed[:])
	} else {
		content.signature, err = a.Token.Sign(hashed[:])
	}
	if err != nil {
		return nil
	}
	hashed = crypto.Hasher(content.serializeWalletBulk())
	if a.Wallet != nil {
		sign, err = a.Wallet.Sign(hashed[:])
	} else {
		sign, err = a.Token.Sign(hashed[:])
	}
	if err != nil {
		return nil
	}
	content.walletSignature = sign
	return content
}

func (a *Author) NewReact(hash []byte, reaction byte, epoch, fee uint64) *React {
	react := React{
		authored: a.NewAuthored(epoch, fee),
		hash:     hash,
		reaction: reaction,
	}
	bulk := react.serializeBulk()
	if a.sign(react.authored, bulk, iReact) {
		return &react
	}
	return nil
}

func (a *Author) NewSponsorshipOffer(audience *Audience, contentType string, content []byte, expiry, revenue, epoch, fee uint64) *SponsorshipOffer {
	if audience == nil {
		return nil
	}
	sponsorOffer := SponsorshipOffer{
		authored:    a.NewAuthored(epoch, fee),
		audience:    audience.Token.PublicKey().ToBytes(),
		contentType: contentType,
		content:     content,
		expiry:      expiry,
		revenue:     revenue,
	}
	bulk := sponsorOffer.serializeBulk()
	if a.sign(sponsorOffer.authored, bulk, iSponsorshipOffer) {
		return &sponsorOffer
	}
	return nil
}

func (a *Author) NewSponsorshipAcceptance(audience *Audience, offer *SponsorshipOffer, epoch, fee uint64) *SponsorshipAcceptance {
	if audience == nil {
		return nil
	}
	sponsorAcceptance := SponsorshipAcceptance{
		authored:     a.NewAuthored(epoch, fee),
		offer:        offer,
		audience:     audience.Token.PublicKey().ToBytes(),
		modSignature: []byte{},
	}
	var err error
	sponsorAcceptance.modSignature, err = audience.Moderation.Sign(sponsorAcceptance.serializeModBulk())
	if err != nil {
		return nil
	}
	bulk := sponsorAcceptance.serializeBulk()
	if a.sign(sponsorAcceptance.authored, bulk, iSponsorshipAcceptance) {
		return &sponsorAcceptance
	}
	return nil
}

func (a *Author) sign(authored *authoredInstruction, bulk []byte, insType byte) bool {
	bytes := authored.serializeWithoutSignature(insType, bulk)
	hash := crypto.Hasher(bytes)
	var err error
	if a.Attorney != nil {
		authored.signature, err = a.Attorney.Sign(hash[:])
	} else {
		authored.signature, err = a.Token.Sign(hash[:])
	}
	if err != nil {
		return false
	}
	util.PutByteArray(authored.signature, &bytes)
	hash = crypto.Hasher(bytes)
	if a.Wallet != nil {
		authored.walletSignature, err = a.Wallet.Sign(hash[:])
	} else if a.Attorney != nil {
		authored.walletSignature, err = a.Attorney.Sign(hash[:])
	} else {
		authored.walletSignature, err = a.Token.Sign(hash[:])
	}
	return err == nil
}
