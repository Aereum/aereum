package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

// Content creation instruction
type Content struct {
	epoch           uint64
	published       uint64
	author          crypto.Token
	audience        crypto.Token
	contentType     string
	content         []byte
	hash            []byte
	sponsored       bool
	encrypted       bool
	subSignature    crypto.Signature
	moderator       crypto.Token
	modSignature    crypto.Signature
	attorney        crypto.Token
	signature       crypto.Signature
	wallet          crypto.Token
	fee             uint64
	walletSignature crypto.Signature
}

func (a *Content) Epoch() uint64 {
	return a.epoch
}

func (content *Content) Validate(v InstructionValidator) bool {
	if content.epoch > v.Epoch() {
		return false
	}
	if !v.HasMember(crypto.HashToken(content.author)) {
		return false
	}
	audienceHash := crypto.HashToken(content.audience)
	stageKeys := v.GetAudienceKeys(audienceHash)
	if stageKeys == nil {
		return false
	}
	if content.sponsored {
		if content.encrypted {
			return false
		}
		if len(content.subSignature) != 0 || len(content.modSignature) != 0 {
			return false
		}
		hash := crypto.Hasher(append(content.author[:], content.audience[:]...))
		ok, contentHash := v.HasGrantedSponser(hash)
		if !ok {
			return false
		}
		if !crypto.Hasher(content.content).Equal(contentHash) {
			return false
		}
		v.AddFeeCollected(content.fee)
		return v.SetPublishSponsor(hash)
	}
	if !stageKeys.Submit.Verify(content.serializeSubBulk()[10:], content.subSignature) {
		return false
	}
	if content.moderator != crypto.ZeroToken {
		if !stageKeys.Moderate.Verify(content.serializeModBulk(), content.modSignature) {
			return false
		}
	}
	v.AddFeeCollected(content.fee)
	return true
}

func (a *Content) Payments() *Payment {
	if len(a.wallet) > 0 {
		return NewPayment(crypto.HashToken(a.wallet), a.fee)
	}
	if len(a.attorney) > 0 {
		return NewPayment(crypto.HashToken(a.attorney), a.fee)
	}
	return NewPayment(crypto.HashToken(a.author), a.fee)
}

func (content *Content) Kind() byte {
	return iContent
}

func (content *Content) serializeSubBulk() []byte {
	bytes := []byte{0, iContent}
	util.PutUint64(content.epoch, &bytes)
	util.PutUint64(content.published, &bytes)
	util.PutToken(content.author, &bytes)
	util.PutToken(content.audience, &bytes)
	util.PutString(content.contentType, &bytes)
	util.PutByteArray(content.content, &bytes)
	util.PutByteArray(content.hash, &bytes)
	util.PutBool(content.sponsored, &bytes)
	util.PutBool(content.encrypted, &bytes)
	return bytes
}

func (content *Content) serializeModBulk() []byte {
	bytes := content.serializeSubBulk()
	util.PutSignature(content.subSignature, &bytes)
	util.PutToken(content.moderator, &bytes)
	return bytes
}

func (content *Content) serializeSignBulk() []byte {
	bytes := content.serializeModBulk()
	util.PutSignature(content.modSignature, &bytes)
	util.PutToken(content.attorney, &bytes)
	return bytes
}

func (content *Content) serializeWalletBulk() []byte {
	bytes := content.serializeSignBulk()
	util.PutSignature(content.signature, &bytes)
	util.PutToken(content.wallet, &bytes)
	util.PutUint64(content.fee, &bytes)
	return bytes
}

func (content *Content) Serialize() []byte {
	bytes := content.serializeWalletBulk()
	util.PutSignature(content.walletSignature, &bytes)
	return bytes
}

func ParseContent(data []byte) *Content {
	if data[0] != 0 || data[1] != iContent {
		return nil
	}
	var content Content
	position := 2
	content.epoch, position = util.ParseUint64(data, position)
	content.published, position = util.ParseUint64(data, position)
	content.author, position = util.ParseToken(data, position)
	content.audience, position = util.ParseToken(data, position)
	content.contentType, position = util.ParseString(data, position)
	content.content, position = util.ParseByteArray(data, position)
	content.hash, position = util.ParseByteArray(data, position)
	content.sponsored, position = util.ParseBool(data, position)
	content.encrypted, position = util.ParseBool(data, position)
	content.subSignature, position = util.ParseSignature(data, position)
	content.moderator, position = util.ParseToken(data, position)
	content.modSignature, position = util.ParseSignature(data, position)
	if len(content.moderator) == 0 && (content.epoch != content.published) {
		return nil
	}
	content.attorney, position = util.ParseToken(data, position)
	msg := data[0:position]
	token := content.author
	if len(content.attorney) > 0 {
		token = content.attorney
	} else if len(content.moderator) > 0 {
		token = content.moderator
	}
	content.signature, position = util.ParseSignature(data, position)
	if !token.Verify(msg, content.signature) {
		return nil
	}
	content.wallet, position = util.ParseToken(data, position)
	content.fee, position = util.ParseUint64(data, position)
	msg = data[0:position]
	content.walletSignature, _ = util.ParseSignature(data, position)
	if content.wallet != crypto.ZeroToken {
		token = content.wallet
	}
	if !token.Verify(msg, content.walletSignature) {
		return nil
	}
	return &content
}

// Reaction instruction
type React struct {
	authored *authoredInstruction
	hash     []byte
	reaction byte
}

func (a *React) Epoch() uint64 {
	return a.authored.epoch
}

func (react *React) Validate(v InstructionValidator) bool {
	if v.HasMember(react.authored.authorHash()) {
		v.AddFeeCollected(react.authored.fee)
		return true
	}
	return false
}

func (react *React) Payments() *Payment {
	return react.authored.payments()
}

func (react *React) Kind() byte {
	return iContent
}

func (react *React) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(react.hash, &bytes)
	util.PutByte(react.reaction, &bytes)
	return bytes
}

func (react *React) Serialize() []byte {
	return react.authored.serialize(iReact, react.serializeBulk())
}

func ParseReact(data []byte) *React {
	if data[0] != 0 || data[1] != iReact {
		return nil
	}
	react := React{
		authored: &authoredInstruction{},
	}
	position := react.authored.parseHead(data)
	react.hash, position = util.ParseByteArray(data, position)
	react.reaction, position = util.ParseByte(data, position)
	if react.authored.parseTail(data, position) {
		return &react
	}
	return nil
}
