package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

// Content creation instruction
type Content struct {
	epoch           uint64
	published       uint64
	author          []byte
	audience        []byte
	contentType     string
	content         []byte
	hash            []byte
	sponsored       bool
	encrypted       bool
	subSignature    []byte
	moderator       []byte
	modSignature    []byte
	attorney        []byte
	signature       []byte
	wallet          []byte
	fee             uint64
	walletSignature []byte
}

func (a *Content) Epoch() uint64 {
	return a.epoch
}

func (content *Content) Validate(v InstructionValidator) bool {
	if content.epoch > v.Epoch() {
		return false
	}
	if !v.HasMember(crypto.Hasher(content.author)) {
		return false
	}
	audienceHash := crypto.Hasher(content.audience)
	keys := v.GetAudienceKeys(audienceHash)
	if content.sponsored {
		if content.encrypted {
			return false
		}
		if len(content.subSignature) != 0 || len(content.modSignature) != 0 {
			return false
		}
		hash := crypto.Hasher(append(content.author, content.audience...))
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
	subKey, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	hash := crypto.Hasher(content.serializeSubBulk()[10:])
	if !subKey.Verify(hash[:], content.subSignature) {
		return false
	}
	if len(content.moderator) != 0 {
		modKey, err := crypto.PublicKeyFromBytes(keys[crypto.PublicKeySize:])
		if err != nil {
			return false
		}
		hash := crypto.Hasher(content.serializeModBulk())
		if !modKey.Verify(hash[:], content.modSignature) {
			return false
		}
	}
	v.AddFeeCollected(content.fee)
	return true
}

func (a *Content) Payments() *Payment {
	if len(a.wallet) > 0 {
		return NewPayment(crypto.Hasher(a.wallet), a.fee)
	}
	if len(a.attorney) > 0 {
		return NewPayment(crypto.Hasher(a.attorney), a.fee)
	}
	return NewPayment(crypto.Hasher(a.author), a.fee)
}

func (content *Content) Kind() byte {
	return iContent
}

func (content *Content) serializeSubBulk() []byte {
	bytes := []byte{0, iContent}
	util.PutUint64(content.epoch, &bytes)
	util.PutUint64(content.published, &bytes)
	util.PutByteArray(content.author, &bytes)
	util.PutByteArray(content.audience, &bytes)
	util.PutString(content.contentType, &bytes)
	util.PutByteArray(content.content, &bytes)
	util.PutByteArray(content.hash, &bytes)
	util.PutBool(content.sponsored, &bytes)
	util.PutBool(content.encrypted, &bytes)
	return bytes
}

func (content *Content) serializeModBulk() []byte {
	bytes := content.serializeSubBulk()
	util.PutByteArray(content.subSignature, &bytes)
	util.PutByteArray(content.moderator, &bytes)
	return bytes
}

func (content *Content) serializeSignBulk() []byte {
	bytes := content.serializeModBulk()
	util.PutByteArray(content.modSignature, &bytes)
	util.PutByteArray(content.attorney, &bytes)
	return bytes
}

func (content *Content) serializeWalletBulk() []byte {
	bytes := content.serializeSignBulk()
	util.PutByteArray(content.signature, &bytes)
	util.PutByteArray(content.wallet, &bytes)
	util.PutUint64(content.fee, &bytes)
	return bytes
}

func (content *Content) Serialize() []byte {
	bytes := content.serializeWalletBulk()
	util.PutByteArray(content.walletSignature, &bytes)
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
	content.author, position = util.ParseByteArray(data, position)
	content.audience, position = util.ParseByteArray(data, position)
	content.contentType, position = util.ParseString(data, position)
	content.content, position = util.ParseByteArray(data, position)
	content.hash, position = util.ParseByteArray(data, position)
	content.sponsored, position = util.ParseBool(data, position)
	content.encrypted, position = util.ParseBool(data, position)
	content.subSignature, position = util.ParseByteArray(data, position)
	content.moderator, position = util.ParseByteArray(data, position)
	content.modSignature, position = util.ParseByteArray(data, position)
	if len(content.moderator) == 0 && (content.epoch != content.published) {
		return nil
	}
	content.attorney, position = util.ParseByteArray(data, position)
	hash := crypto.Hasher(data[0:position])
	var pubKey crypto.PublicKey
	var err error
	if len(content.attorney) > 0 {
		pubKey, err = crypto.PublicKeyFromBytes(content.attorney)
	} else if len(content.moderator) > 0 {
		pubKey, err = crypto.PublicKeyFromBytes(content.moderator)
	} else {
		pubKey, err = crypto.PublicKeyFromBytes(content.author)
	}
	if err != nil {
		return nil
	}
	content.signature, position = util.ParseByteArray(data, position)
	if !pubKey.Verify(hash[:], content.signature) {
		return nil
	}
	content.wallet, position = util.ParseByteArray(data, position)
	content.fee, position = util.ParseUint64(data, position)
	hash = crypto.Hasher(data[0:position])
	content.walletSignature, _ = util.ParseByteArray(data, position)
	if len(content.wallet) > 0 {
		pubKey, err = crypto.PublicKeyFromBytes(content.wallet)
		if err != nil {
			return nil
		}
	}
	if !pubKey.Verify(hash[:], content.walletSignature) {
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
