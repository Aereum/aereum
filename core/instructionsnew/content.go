package instructionsnew

import (
	"bytes"

	"github.com/Aereum/aereum/core/crypto"
)

// Content creation instruction
type Content struct {
	authored     *authoredInstruction
	audience     []byte
	author       []byte
	contentType  string
	content      []byte
	hash         []byte
	sponsored    bool
	encrypted    bool
	subSignature []byte
	modSignature []byte
}

func (content *Content) Validate(block *Block) bool {
	if !block.validator.HasMember(content.authored.authorHash()) {
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
	if len(content.modSignature) == 0 {
		if !bytes.Equal(content.author, content.authored.author) {
			return false
		}
	} else {
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
	bytes := make([]byte, 0)
	PutByteArray(content.audience, &bytes)
	PutByteArray(content.author, &bytes)
	PutString(content.contentType, &bytes)
	PutByteArray(content.content, &bytes)
	PutByteArray(content.hash, &bytes)
	PutBool(content.sponsored, &bytes)
	PutBool(content.encrypted, &bytes)
	return bytes
}

func (content *Content) serializeModBulk() []byte {
	bytes := content.serializeSubBulk()
	PutByteArray(content.subSignature, &bytes)
	return bytes
}

func (content *Content) serializeBulk() []byte {
	bytes := content.serializeModBulk()
	PutByteArray(content.modSignature, &bytes)
	return bytes
}

func (content *Content) Serialize() []byte {
	return content.authored.serialize(iContent, content.serializeBulk())
}

func ParseContent(data []byte) *Content {
	if data[0] != 0 || data[1] != iContent {
		return nil
	}
	content := Content{
		authored: &authoredInstruction{},
	}
	position := content.authored.parseHead(data)
	content.audience, position = ParseByteArray(data, position)
	content.author, position = ParseByteArray(data, position)
	content.contentType, position = ParseString(data, position)
	content.content, position = ParseByteArray(data, position)
	content.hash, position = ParseByteArray(data, position)
	content.sponsored, position = ParseBool(data, position)
	content.encrypted, position = ParseBool(data, position)
	content.subSignature, position = ParseByteArray(data, position)
	content.modSignature, position = ParseByteArray(data, position)
	if content.authored.parseTail(data, position) {
		return &content
	}
	return nil
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
