package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

// Content creation instruction
type Content struct {
	epoch           uint64
	Published       uint64
	Author          crypto.Token
	Audience        crypto.Token
	ContentType     string
	Content         []byte
	Hash            []byte
	Sponsored       bool
	Encrypted       bool
	SubSignature    crypto.Signature
	Moderator       crypto.Token
	ModSignature    crypto.Signature
	Attorney        crypto.Token
	Signature       crypto.Signature
	Wallet          crypto.Token
	Fee             uint64
	WalletSignature crypto.Signature
}

func (a *Content) Authority() crypto.Token {
	return a.Author
}

func (a *Content) Epoch() uint64 {
	return a.epoch
}

func (content *Content) Validate(v InstructionValidator) bool {
	if content.epoch > v.Epoch() {
		return false
	}
	if !v.HasMember(crypto.HashToken(content.Author)) {
		return false
	}
	audienceHash := crypto.HashToken(content.Audience)
	stageKeys := v.GetAudienceKeys(audienceHash)
	if stageKeys == nil {
		return false
	}
	if content.Sponsored {
		if content.Encrypted {
			return false
		}
		if len(content.SubSignature) != 0 || len(content.ModSignature) != 0 {
			return false
		}
		hash := crypto.Hasher(append(content.Author[:], content.Audience[:]...))
		ok, contentHash := v.HasGrantedSponser(hash)
		if !ok {
			return false
		}
		if !crypto.Hasher(content.Content).Equal(contentHash) {
			return false
		}
		v.AddFeeCollected(content.Fee)
		return v.SetPublishSponsor(hash)
	}
	if !stageKeys.Submit.Verify(content.serializeSubBulk()[10:], content.SubSignature) {
		return false
	}
	if content.Moderator != crypto.ZeroToken {
		if !stageKeys.Moderate.Verify(content.serializeModBulk(), content.ModSignature) {
			return false
		}
	}
	v.AddFeeCollected(content.Fee)
	return true
}

func (a *Content) Payments() *Payment {
	if len(a.Wallet) > 0 {
		return NewPayment(crypto.HashToken(a.Wallet), a.Fee)
	}
	if len(a.Attorney) > 0 {
		return NewPayment(crypto.HashToken(a.Attorney), a.Fee)
	}
	return NewPayment(crypto.HashToken(a.Author), a.Fee)
}

func (content *Content) Kind() byte {
	return IContent
}

func (content *Content) serializeSubBulk() []byte {
	bytes := []byte{0, IContent}
	util.PutUint64(content.epoch, &bytes)
	util.PutUint64(content.Published, &bytes)
	util.PutToken(content.Author, &bytes)
	util.PutToken(content.Audience, &bytes)
	util.PutString(content.ContentType, &bytes)
	util.PutByteArray(content.Content, &bytes)
	util.PutByteArray(content.Hash, &bytes)
	util.PutBool(content.Sponsored, &bytes)
	util.PutBool(content.Encrypted, &bytes)
	return bytes
}

func (content *Content) serializeModBulk() []byte {
	bytes := content.serializeSubBulk()
	util.PutSignature(content.SubSignature, &bytes)
	util.PutToken(content.Moderator, &bytes)
	return bytes
}

func (content *Content) serializeSignBulk() []byte {
	bytes := content.serializeModBulk()
	util.PutSignature(content.ModSignature, &bytes)
	util.PutToken(content.Attorney, &bytes)
	return bytes
}

func (content *Content) serializeWalletBulk() []byte {
	bytes := content.serializeSignBulk()
	util.PutSignature(content.Signature, &bytes)
	util.PutToken(content.Wallet, &bytes)
	util.PutUint64(content.Fee, &bytes)
	return bytes
}

func (content *Content) Serialize() []byte {
	bytes := content.serializeWalletBulk()
	util.PutSignature(content.WalletSignature, &bytes)
	return bytes
}

func ParseContent(data []byte) *Content {
	if data[0] != 0 || data[1] != IContent {
		return nil
	}
	var content Content
	position := 2
	content.epoch, position = util.ParseUint64(data, position)
	content.Published, position = util.ParseUint64(data, position)
	content.Author, position = util.ParseToken(data, position)
	content.Audience, position = util.ParseToken(data, position)
	content.ContentType, position = util.ParseString(data, position)
	content.Content, position = util.ParseByteArray(data, position)
	content.Hash, position = util.ParseByteArray(data, position)
	content.Sponsored, position = util.ParseBool(data, position)
	content.Encrypted, position = util.ParseBool(data, position)
	content.SubSignature, position = util.ParseSignature(data, position)
	content.Moderator, position = util.ParseToken(data, position)
	content.ModSignature, position = util.ParseSignature(data, position)
	if len(content.Moderator) == 0 && (content.epoch != content.Published) {
		return nil
	}
	content.Attorney, position = util.ParseToken(data, position)
	msg := data[0:position]
	token := content.Author
	if len(content.Attorney) > 0 {
		token = content.Attorney
	} else if len(content.Moderator) > 0 {
		token = content.Moderator
	}
	content.Signature, position = util.ParseSignature(data, position)
	if !token.Verify(msg, content.Signature) {
		return nil
	}
	content.Wallet, position = util.ParseToken(data, position)
	content.Fee, position = util.ParseUint64(data, position)
	msg = data[0:position]
	content.WalletSignature, _ = util.ParseSignature(data, position)
	if content.Wallet != crypto.ZeroToken {
		token = content.Wallet
	}
	if !token.Verify(msg, content.WalletSignature) {
		return nil
	}
	return &content
}

// Reaction instruction
type React struct {
	Authored *AuthoredInstruction
	Hash     []byte
	Reaction byte
}

func (a *React) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *React) Epoch() uint64 {
	return a.Authored.epoch
}

func (react *React) Validate(v InstructionValidator) bool {
	if v.HasMember(react.Authored.authorHash()) {
		v.AddFeeCollected(react.Authored.Fee)
		return true
	}
	return false
}

func (react *React) Payments() *Payment {
	return react.Authored.payments()
}

func (react *React) Kind() byte {
	return IContent
}

func (react *React) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(react.Hash, &bytes)
	util.PutByte(react.Reaction, &bytes)
	return bytes
}

func (react *React) Serialize() []byte {
	return react.Authored.serialize(IReact, react.serializeBulk())
}

func ParseReact(data []byte) *React {
	if data[0] != 0 || data[1] != IReact {
		return nil
	}
	react := React{
		Authored: &AuthoredInstruction{},
	}
	position := react.Authored.parseHead(data)
	react.Hash, position = util.ParseByteArray(data, position)
	react.Reaction, position = util.ParseByte(data, position)
	if react.Authored.parseTail(data, position) {
		return &react
	}
	return nil
}
