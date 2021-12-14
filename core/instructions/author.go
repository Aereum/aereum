package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
	"github.com/Aereum/aereum/core/util"
)

type Author struct {
	PrivateKey crypto.PrivateKey
	Wallet     crypto.PrivateKey
	Attorney   crypto.PrivateKey
}

func (a *Author) NewAuthored(epoch, fee uint64) *authoredInstruction {
	if a.PrivateKey == crypto.ZeroPrivateKey {
		return nil
	}
	authored := authoredInstruction{
		epoch:  epoch,
		author: a.PrivateKey.PublicKey(),
		fee:    fee,
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		authored.wallet = a.Wallet.PublicKey()
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		authored.attorney = a.Attorney.PublicKey()
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

func (a *Author) NewJoinNetworkThirdParty(token crypto.Token, caption string, details string, epoch, fee uint64) *JoinNetwork {
	authored := authoredInstruction{
		epoch:  epoch,
		author: token,
		fee:    fee,
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		authored.attorney = a.Attorney.PublicKey()
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		authored.wallet = a.Wallet.PublicKey()
	} else if a.Attorney != crypto.ZeroPrivateKey {
		authored.wallet = a.Attorney.PublicKey()
	} else {
		authored.wallet = a.PrivateKey.PublicKey()
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

func (a *Author) NewGrantPowerOfAttorney(attorney crypto.Token, epoch, fee uint64) *GrantPowerOfAttorney {
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

func (a *Author) NewRevokePowerOfAttorney(attorney crypto.Token, epoch, fee uint64) *RevokePowerOfAttorney {
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

func (a *Author) NewCreateEphemeral(token crypto.Token, expiry, epoch, fee uint64) *CreateEphemeral {
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

func (a *Author) NewCreateAudience(audience *Stage, flag byte, description string, epoch, fee uint64) *CreateStage {
	newAudience := CreateStage{
		authored:    a.NewAuthored(epoch, fee),
		audience:    audience.PrivateKey.PublicKey(),
		submission:  audience.Submission.PublicKey(),
		moderation:  audience.Moderation.PublicKey(),
		flag:        flag,
		description: description,
	}
	bulk := newAudience.serializeBulk()
	if a.sign(newAudience.authored, bulk, iCreateAudience) {
		return &newAudience
	}
	return nil
}

func (a *Author) NewJoinAudience(audience crypto.Token, presentation string, epoch, fee uint64) *JoinStage {
	join := JoinStage{
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

func (a *Author) NewAcceptJoinAudience(audience *Stage, member, key crypto.Token, level byte, epoch, fee uint64) *AcceptJoinStage {
	accept := AcceptJoinStage{
		authored: a.NewAuthored(epoch, fee),
		stage:    audience.PrivateKey.PublicKey(),
		member:   member,
		read:     []byte{},
		submit:   []byte{},
		moderate: []byte{},
	}
	pub, prv := dh.NewEphemeralKey()
	cipher := dh.ConsensusCipher(prv, key)
	accept.read = cipher.Seal(audience.CipherKey)
	if level > 0 {
		accept.submit = cipher.Seal(audience.Submission[:32])
	}
	if level > 1 {
		accept.moderate = cipher.Seal(audience.Moderation[:32])
	}
	accept.diffHellKey = pub
	modbulk := accept.serializeModBulk()
	accept.modSignature = audience.Moderation.Sign(modbulk)
	bulk := accept.serializeBulk()
	if a.sign(accept.authored, bulk, iAcceptJoinRequest) {
		return &accept
	}
	return nil
}

func (a *Author) NewUpdateAudience(audience *Stage, readers, submiters, moderators map[crypto.Token]crypto.Token, flag byte, description string, epoch, fee uint64) *UpdateStage {
	update := UpdateStage{
		authored:   a.NewAuthored(epoch, fee),
		stage:      audience.PrivateKey.PublicKey(),
		submission: audience.Submission.PublicKey(),
		moderation: audience.Moderation.PublicKey(),

		flag:        flag,
		description: description,
		readMembers: make(TokenCiphers, 0),
		subMembers:  make(TokenCiphers, 0),
		modMembers:  make(TokenCiphers, 0),
	}
	prv, pub := dh.NewEphemeralKey()
	update.diffHellKey = pub
	for token, key := range readers {
		cipher := dh.ConsensusCipher(prv, key)
		update.readMembers = append(update.readMembers, TokenCipher{token: token, cipher: cipher.Seal(audience.CipherKey)})
	}
	for token, key := range submiters {
		cipher := dh.ConsensusCipher(prv, key)
		update.subMembers = append(update.subMembers, TokenCipher{token: token, cipher: cipher.Seal(audience.Submission[:32])})
	}
	for token, key := range moderators {
		cipher := dh.ConsensusCipher(prv, key)
		update.modMembers = append(update.modMembers, TokenCipher{token: token, cipher: cipher.Seal(audience.Moderation[:32])})
	}
	update.audSignature = audience.PrivateKey.Sign(update.serializeAudBulk())
	bulk := update.serializeBulk()
	if a.sign(update.authored, bulk, iUpdateAudience) {
		return &update
	}
	return nil
}

func (a *Author) ModerateContent(audience *Stage, content *Content, epoch, fee uint64) *Content {
	if audience == nil || audience.Moderation == crypto.ZeroPrivateKey {
		return nil
	}
	if audience.PrivateKey.PublicKey() != content.audience {
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
		moderator:    audience.PrivateKey.PublicKey(),
		attorney:     a.Attorney.PublicKey(),
		wallet:       a.Wallet.PublicKey(),
		fee:          fee,
	}
	msg := newContent.serializeModBulk()
	newContent.modSignature = audience.Moderation.Sign(msg)
	if a.Attorney != crypto.ZeroPrivateKey {
		newContent.attorney = a.Attorney.PublicKey()
		newContent.signature = a.Attorney.Sign(newContent.serializeSignBulk())
	} else {
		newContent.signature = a.PrivateKey.Sign(newContent.serializeSignBulk())
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		newContent.wallet = a.Wallet.PublicKey()
		newContent.walletSignature = a.Attorney.Sign(newContent.serializeWalletBulk())
	} else {
		newContent.walletSignature = a.PrivateKey.Sign(newContent.serializeWalletBulk())
	}
	return newContent
}

func (a *Author) NewContent(audience *Stage, contentType string, message []byte, hash, encrypted bool, epoch, fee uint64) *Content {
	if audience == nil {
		return nil
	}
	content := &Content{
		epoch:       epoch,
		published:   epoch,
		author:      a.PrivateKey.PublicKey(),
		audience:    audience.PrivateKey.PublicKey(),
		contentType: contentType,
		hash:        []byte{},
		sponsored:   false,
		encrypted:   encrypted,
		attorney:    a.Attorney.PublicKey(),
		wallet:      a.Wallet.PublicKey(),
		fee:         fee,
	}
	if encrypted {
		cipher := crypto.CipherFromKey(audience.CipherKey)
		content.content = cipher.Seal(message)
	} else {
		content.content = message
	}
	if hash {
		hashed := crypto.Hasher(message)
		content.hash = hashed[:]
	}
	subBulk := content.serializeSubBulk()
	content.subSignature = audience.Submission.Sign(subBulk[10:])
	util.PutSignature(content.subSignature, &subBulk)
	if audience.Moderation != crypto.ZeroPrivateKey {
		content.moderator = a.PrivateKey.PublicKey()
		content.modSignature = audience.Moderation.Sign(content.serializeModBulk())
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		content.signature = a.Attorney.Sign(content.serializeSignBulk())
	} else {
		content.signature = a.PrivateKey.Sign(content.serializeSignBulk())
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		content.walletSignature = a.Wallet.Sign(content.serializeWalletBulk())
	} else {
		content.walletSignature = a.PrivateKey.Sign(content.serializeWalletBulk())
	}
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

func (a *Author) NewSponsorshipOffer(audience *Stage, contentType string, content []byte, expiry, revenue, epoch, fee uint64) *SponsorshipOffer {
	if audience == nil {
		return nil
	}
	sponsorOffer := SponsorshipOffer{
		authored:    a.NewAuthored(epoch, fee),
		stage:       audience.PrivateKey.PublicKey(),
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

func (a *Author) NewSponsorshipAcceptance(audience *Stage, offer *SponsorshipOffer, epoch, fee uint64) *SponsorshipAcceptance {
	if audience == nil {
		return nil
	}
	sponsorAcceptance := SponsorshipAcceptance{
		authored: a.NewAuthored(epoch, fee),
		offer:    offer,
		stage:    audience.PrivateKey.PublicKey(),
	}
	sponsorAcceptance.modSignature = audience.Moderation.Sign(sponsorAcceptance.serializeModBulk())
	bulk := sponsorAcceptance.serializeBulk()
	if a.sign(sponsorAcceptance.authored, bulk, iSponsorshipAcceptance) {
		return &sponsorAcceptance
	}
	return nil
}

func (a *Author) sign(authored *authoredInstruction, bulk []byte, insType byte) bool {
	bytes := authored.serializeWithoutSignature(insType, bulk)
	if a.Attorney != crypto.ZeroPrivateKey {
		authored.signature = a.Attorney.Sign(bytes)
	} else {
		authored.signature = a.PrivateKey.Sign(bytes)
	}
	util.PutSignature(authored.signature, &bytes)
	if a.Wallet != crypto.ZeroPrivateKey {
		authored.walletSignature = a.Wallet.Sign(bytes)
	} else if a.Attorney != crypto.ZeroPrivateKey {
		authored.walletSignature = a.Attorney.Sign(bytes)
	} else {
		authored.walletSignature = a.PrivateKey.Sign(bytes)
	}
	return true
}
