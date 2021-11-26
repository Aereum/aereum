package instructionsnew

import (
	"github.com/Aereum/aereum/core/crypto"
)

// Content creation instruction
type Content struct {
	authored        *authoredInstruction
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
	attorney        []byte
	signature       []byte
	moderator       []byte
	modSignature    []byte
	wallet          []byte
	fee             uint64
	walletSignature []byte
}

func (content *Content) Validate(block *Block) bool {
	if content.epoch > block.Epoch {
		return false
	}
	if !block.validator.HasMember(crypto.Hasher(content.author)) {
		return false
	}
	audienceHash := crypto.Hasher(content.audience)
	keys := block.validator.GetAudienceKeys(audienceHash)
	if content.sponsored {
		if content.encrypted {
			return false
		}
		if len(content.subSignature) != 0 || len(content.modSignature) != 0 {
			return false
		}
		hash := crypto.Hasher(append(content.authored.author, content.audience...))
		ok, contentHash := block.validator.HasGrantedSponser(hash)
		if !ok {
			return false
		}
		if !crypto.Hasher(content.content).Equal(contentHash) {
			return false
		}
		return block.SetPublishSponsor(hash)
	}
	subKey, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	hash := crypto.Hasher(content.serializeSubBulk())
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
	return true
}

func (content *Content) Payments() *Payment {
	return content.authored.payments()
}

func (content *Content) Kind() byte {
	return iContent
}

func (content *Content) serializeSubBulk() []byte {
	bytes := []byte{0, iContent}
	PutUint64(content.epoch, &bytes)
	PutUint64(content.published, &bytes)
	PutByteArray(content.author, &bytes)
	PutByteArray(content.audience, &bytes)
	PutString(content.contentType, &bytes)
	PutByteArray(content.content, &bytes)
	PutByteArray(content.hash, &bytes)
	PutBool(content.sponsored, &bytes)
	PutBool(content.encrypted, &bytes)
	return bytes
}

func (content *Content) serializeModBulk() []byte {
	bytes := content.serializeSubBulk()
	PutByteArray(content.attorney, &bytes)
	PutByteArray(content.signature, &bytes)
	PutByteArray(content.moderator, &bytes)
	return bytes
}

func (content *Content) serializeWalletBulk() []byte {
	bytes := content.serializeModBulk()
	PutByteArray(content.modSignature, &bytes)
	PutByteArray(content.wallet, &bytes)
	PutUint64(content.fee, &bytes)
	return bytes
}

func (content *Content) Serialize() []byte {
	bytes := content.serializeWalletBulk()
	PutByteArray(content.walletSignature, &bytes)
	return bytes
}

func ParseContent(data []byte) *Content {
	if data[0] != 0 || data[1] != iContent {
		return nil
	}
	var content Content
	position := content.authored.parseHead(data)
	content.epoch, position = ParseUint64(data, position)
	content.published, position = ParseUint64(data, position)
	content.author, position = ParseByteArray(data, position)
	content.audience, position = ParseByteArray(data, position)
	content.contentType, position = ParseString(data, position)
	content.content, position = ParseByteArray(data, position)
	content.hash, position = ParseByteArray(data, position)
	content.sponsored, position = ParseBool(data, position)
	content.encrypted, position = ParseBool(data, position)
	content.subSignature, position = ParseByteArray(data, position)
	content.attorney, position = ParseByteArray(data, position)
	hash := crypto.Hasher(append(data[0:2], data[16:position]...))
	var pubKey crypto.PublicKey
	var err error
	if len(content.attorney) > 0 {
		pubKey, err = crypto.PublicKeyFromBytes(content.attorney)
	} else {
		pubKey, err = crypto.PublicKeyFromBytes(content.author)
	}
	if err != nil {
		return nil
	}
	content.signature, position = ParseByteArray(data, position)
	if !pubKey.Verify(hash[:], content.signature) {
		return nil
	}
	content.moderator, position = ParseByteArray(data, position)
	content.modSignature, position = ParseByteArray(data, position)
	if len(content.moderator) == 0 && (content.epoch != content.published) {
		return nil
	}
	content.wallet, position = ParseByteArray(data, position)
	content.fee, position = ParseUint64(data, position)
	hash = crypto.Hasher(data[0:position])
	content.walletSignature, _ = ParseByteArray(data, position)
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

func (react *React) Validate(block *Block) bool {
	return block.validator.HasMember(react.authored.authorHash())
}

func (react *React) Payments() *Payment {
	return react.authored.payments()
}

func (react *React) Kind() byte {
	return iContent
}

func (react *React) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(react.hash, &bytes)
	PutByte(react.reaction, &bytes)
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
	react.hash, position = ParseByteArray(data, position)
	react.reaction, position = ParseByte(data, position)
	if react.authored.parseTail(data, position) {
		return &react
	}
	return nil
}

// CREATING TEST ENTRY
// type ContentBase struct {
// 	audience     crypto.PrivateKey
// 	author       crypto.PrivateKey
// 	contentType  string
// 	content      []byte
// 	hash         []byte
// 	sponsored    bool
// 	encrypted    bool
// 	subSignature []byte
// 	modSignature []byte
// }
