package instructionsnew

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

func (content *Content) Payments() *Payment {
	return content.authored.payments()
}

func (content *Content) Kind() byte {
	return iContent
}

func (content *Content) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(content.audience, &bytes)
	PutByteArray(content.author, &bytes)
	PutString(content.contentType, &bytes)
	PutByteArray(content.content, &bytes)
	PutByteArray(content.hash, &bytes)
	PutBool(content.sponsored, &bytes)
	PutBool(content.encrypted, &bytes)
	PutByteArray(content.subSignature, &bytes)
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
