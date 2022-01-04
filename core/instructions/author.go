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

func (a *Author) NewAuthored(epoch, fee uint64) *AuthoredInstruction {
	if a.PrivateKey == crypto.ZeroPrivateKey {
		return nil
	}
	authored := AuthoredInstruction{
		epoch:  epoch,
		Author: a.PrivateKey.PublicKey(),
		Fee:    fee,
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		authored.Wallet = a.Wallet.PublicKey()
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		authored.Attorney = a.Attorney.PublicKey()
	}
	return &authored
}

func (a *Author) NewJoinNetwork(caption string, details string, epoch, fee uint64) *JoinNetwork {
	join := JoinNetwork{
		Authored: a.NewAuthored(epoch, fee),
		Caption:  caption,
		Details:  details,
	}
	bulk := join.serializeBulk()
	if a.sign(join.Authored, bulk, IJoinNetwork) {
		return &join
	}
	return nil
}

func (a *Author) NewJoinNetworkThirdParty(token crypto.Token, caption string, details string, epoch, fee uint64) *JoinNetwork {
	authored := AuthoredInstruction{
		epoch:  epoch,
		Author: token,
		Fee:    fee,
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		authored.Attorney = a.Attorney.PublicKey()
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		authored.Wallet = a.Wallet.PublicKey()
	} else if a.Attorney != crypto.ZeroPrivateKey {
		authored.Wallet = a.Attorney.PublicKey()
	} else {
		authored.Wallet = a.PrivateKey.PublicKey()
	}
	join := JoinNetwork{
		Authored: &authored,
		Caption:  caption,
		Details:  details,
	}
	bulk := join.serializeBulk()
	if a.sign(join.Authored, bulk, IJoinNetwork) {
		return &join
	}
	return nil
}

func (a *Author) NewUpdateInfo(details string, epoch, fee uint64) *UpdateInfo {
	update := UpdateInfo{
		Authored: a.NewAuthored(epoch, fee),
		Details:  details,
	}
	bulk := update.serializeBulk()
	if a.sign(update.Authored, bulk, IUpdateInfo) {
		return &update
	}
	return nil
}

func (a *Author) NewGrantPowerOfAttorney(attorney crypto.Token, epoch, fee uint64) *GrantPowerOfAttorney {
	grant := GrantPowerOfAttorney{
		Authored: a.NewAuthored(epoch, fee),
		Attorney: attorney,
	}
	bulk := grant.serializeBulk()
	if a.sign(grant.Authored, bulk, IGrantPowerOfAttorney) {
		return &grant
	}
	return nil
}

func (a *Author) NewRevokePowerOfAttorney(attorney crypto.Token, epoch, fee uint64) *RevokePowerOfAttorney {
	revoke := RevokePowerOfAttorney{
		Authored: a.NewAuthored(epoch, fee),
		Attorney: attorney,
	}
	bulk := revoke.serializeBulk()
	if a.sign(revoke.Authored, bulk, IRevokePowerOfAttorney) {
		return &revoke
	}
	return nil
}

func (a *Author) NewCreateEphemeral(token crypto.Token, expiry, epoch, fee uint64) *CreateEphemeral {
	ephemeral := CreateEphemeral{
		Authored:       a.NewAuthored(epoch, fee),
		EphemeralToken: token,
		Expiry:         expiry,
	}
	bulk := ephemeral.serializeBulk()
	if a.sign(ephemeral.Authored, bulk, ICreateEphemeral) {
		return &ephemeral
	}
	return nil
}

func (a *Author) NewSecureChannel(tokenRange []byte, nonce uint64, encryptedNonce, content []byte, epoch, fee uint64) *SecureChannel {
	secure := SecureChannel{
		Authored:       a.NewAuthored(epoch, fee),
		TokenRange:     tokenRange,
		Nonce:          nonce,
		EncryptedNonce: encryptedNonce,
		Content:        content,
	}
	bulk := secure.serializeBulk()
	if a.sign(secure.Authored, bulk, ISecureChannel) {
		return &secure
	}
	return nil
}

func (a *Author) NewCreateAudience(audience *Stage, flag byte, description string, epoch, fee uint64) *CreateStage {
	newAudience := CreateStage{
		Authored:    a.NewAuthored(epoch, fee),
		Audience:    audience.PrivateKey.PublicKey(),
		Submission:  audience.Submission.PublicKey(),
		Moderation:  audience.Moderation.PublicKey(),
		Flag:        flag,
		Description: description,
	}
	bulk := newAudience.serializeBulk()
	if a.sign(newAudience.Authored, bulk, ICreateAudience) {
		return &newAudience
	}
	return nil
}

func (a *Author) NewJoinAudience(audience crypto.Token, presentation string, epoch, fee uint64) *JoinStage {
	join := JoinStage{
		Authored:     a.NewAuthored(epoch, fee),
		Audience:     audience,
		Presentation: presentation,
	}
	bulk := join.serializeBulk()
	if a.sign(join.Authored, bulk, IJoinAudience) {
		return &join
	}
	return nil
}

func (a *Author) NewAcceptJoinAudience(audience *Stage, member, key crypto.Token, level byte, epoch, fee uint64) *AcceptJoinStage {
	accept := AcceptJoinStage{
		Authored: a.NewAuthored(epoch, fee),
		Stage:    audience.PrivateKey.PublicKey(),
		Member:   member,
		Read:     []byte{},
		Submit:   []byte{},
		Moderate: []byte{},
	}
	pub, prv := dh.NewEphemeralKey()
	cipher := dh.ConsensusCipher(prv, key)
	accept.Read = cipher.Seal(audience.CipherKey)
	if level > 0 {
		accept.Submit = cipher.Seal(audience.Submission[:32])
	}
	if level > 1 {
		accept.Moderate = cipher.Seal(audience.Moderation[:32])
	}
	accept.DiffHellKey = pub
	modbulk := accept.serializeModBulk()
	accept.modSignature = audience.Moderation.Sign(modbulk)
	bulk := accept.serializeBulk()
	if a.sign(accept.Authored, bulk, IAcceptJoinRequest) {
		return &accept
	}
	return nil
}

func (a *Author) NewUpdateAudience(audience *Stage, readers, submiters, moderators map[crypto.Token]crypto.Token, flag byte, description string, epoch, fee uint64) *UpdateStage {
	update := UpdateStage{
		Authored:   a.NewAuthored(epoch, fee),
		Stage:      audience.PrivateKey.PublicKey(),
		Submission: audience.Submission.PublicKey(),
		Moderation: audience.Moderation.PublicKey(),

		Flag:        flag,
		Description: description,
		ReadMembers: make(TokenCiphers, 0),
		SubMembers:  make(TokenCiphers, 0),
		ModMembers:  make(TokenCiphers, 0),
	}
	prv, pub := dh.NewEphemeralKey()
	update.DiffHellKey = pub
	for token, key := range readers {
		cipher := dh.ConsensusCipher(prv, key)
		update.ReadMembers = append(update.ReadMembers, TokenCipher{Token: token, Cipher: cipher.Seal(audience.CipherKey)})
	}
	for token, key := range submiters {
		cipher := dh.ConsensusCipher(prv, key)
		update.SubMembers = append(update.SubMembers, TokenCipher{Token: token, Cipher: cipher.Seal(audience.Submission[:32])})
	}
	for token, key := range moderators {
		cipher := dh.ConsensusCipher(prv, key)
		update.ModMembers = append(update.ModMembers, TokenCipher{Token: token, Cipher: cipher.Seal(audience.Moderation[:32])})
	}
	update.audSignature = audience.PrivateKey.Sign(update.serializeAudBulk())
	bulk := update.serializeBulk()
	if a.sign(update.Authored, bulk, IUpdateAudience) {
		return &update
	}
	return nil
}

func (a *Author) ModerateContent(audience *Stage, content *Content, epoch, fee uint64) *Content {
	if audience == nil || audience.Moderation == crypto.ZeroPrivateKey {
		return nil
	}
	if audience.PrivateKey.PublicKey() != content.Audience {
		return nil
	}
	newContent := &Content{
		epoch:        epoch,
		Published:    content.epoch,
		Author:       content.Author,
		Audience:     content.Audience,
		ContentType:  content.ContentType,
		Content:      content.Content,
		Sponsored:    content.Sponsored,
		Encrypted:    content.Encrypted,
		SubSignature: content.SubSignature,
		Moderator:    audience.PrivateKey.PublicKey(),
		Attorney:     a.Attorney.PublicKey(),
		Wallet:       a.Wallet.PublicKey(),
		Fee:          fee,
	}
	msg := newContent.serializeModBulk()
	newContent.ModSignature = audience.Moderation.Sign(msg)
	if a.Attorney != crypto.ZeroPrivateKey {
		newContent.Attorney = a.Attorney.PublicKey()
		newContent.Signature = a.Attorney.Sign(newContent.serializeSignBulk())
	} else {
		newContent.Signature = a.PrivateKey.Sign(newContent.serializeSignBulk())
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		newContent.Wallet = a.Wallet.PublicKey()
		newContent.WalletSignature = a.Attorney.Sign(newContent.serializeWalletBulk())
	} else {
		newContent.WalletSignature = a.PrivateKey.Sign(newContent.serializeWalletBulk())
	}
	return newContent
}

func (a *Author) NewContent(audience *Stage, contentType string, message []byte, hash, encrypted bool, epoch, fee uint64) *Content {
	if audience == nil {
		return nil
	}
	content := &Content{
		epoch:       epoch,
		Published:   epoch,
		Author:      a.PrivateKey.PublicKey(),
		Audience:    audience.PrivateKey.PublicKey(),
		ContentType: contentType,
		Hash:        []byte{},
		Sponsored:   false,
		Encrypted:   encrypted,
		Attorney:    a.Attorney.PublicKey(),
		Wallet:      a.Wallet.PublicKey(),
		Fee:         fee,
	}
	if encrypted {
		cipher := crypto.CipherFromKey(audience.CipherKey)
		content.Content = cipher.Seal(message)
	} else {
		content.Content = message
	}
	if hash {
		hashed := crypto.Hasher(message)
		content.Hash = hashed[:]
	}
	subBulk := content.serializeSubBulk()
	content.SubSignature = audience.Submission.Sign(subBulk[10:])
	util.PutSignature(content.SubSignature, &subBulk)
	if audience.Moderation != crypto.ZeroPrivateKey {
		content.Moderator = a.PrivateKey.PublicKey()
		content.ModSignature = audience.Moderation.Sign(content.serializeModBulk())
	}
	if a.Attorney != crypto.ZeroPrivateKey {
		content.Signature = a.Attorney.Sign(content.serializeSignBulk())
	} else {
		content.Signature = a.PrivateKey.Sign(content.serializeSignBulk())
	}
	if a.Wallet != crypto.ZeroPrivateKey {
		content.WalletSignature = a.Wallet.Sign(content.serializeWalletBulk())
	} else {
		content.WalletSignature = a.PrivateKey.Sign(content.serializeWalletBulk())
	}
	return content
}

func (a *Author) NewReact(hash []byte, reaction byte, epoch, fee uint64) *React {
	react := React{
		Authored: a.NewAuthored(epoch, fee),
		Hash:     hash,
		Reaction: reaction,
	}
	bulk := react.serializeBulk()
	if a.sign(react.Authored, bulk, IReact) {
		return &react
	}
	return nil
}

func (a *Author) NewSponsorshipOffer(audience *Stage, contentType string, content []byte, expiry, revenue, epoch, fee uint64) *SponsorshipOffer {
	if audience == nil {
		return nil
	}
	sponsorOffer := SponsorshipOffer{
		Authored:    a.NewAuthored(epoch, fee),
		Stage:       audience.PrivateKey.PublicKey(),
		ContentType: contentType,
		Content:     content,
		Expiry:      expiry,
		Revenue:     revenue,
	}
	bulk := sponsorOffer.serializeBulk()
	if a.sign(sponsorOffer.Authored, bulk, ISponsorshipOffer) {
		return &sponsorOffer
	}
	return nil
}

func (a *Author) NewSponsorshipAcceptance(audience *Stage, offer *SponsorshipOffer, epoch, fee uint64) *SponsorshipAcceptance {
	if audience == nil {
		return nil
	}
	sponsorAcceptance := SponsorshipAcceptance{
		Authored: a.NewAuthored(epoch, fee),
		Offer:    offer,
		Stage:    audience.PrivateKey.PublicKey(),
	}
	sponsorAcceptance.modSignature = audience.Moderation.Sign(sponsorAcceptance.serializeModBulk())
	bulk := sponsorAcceptance.serializeBulk()
	if a.sign(sponsorAcceptance.Authored, bulk, ISponsorshipAcceptance) {
		return &sponsorAcceptance
	}
	return nil
}

func (a *Author) sign(authored *AuthoredInstruction, bulk []byte, insType byte) bool {
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
