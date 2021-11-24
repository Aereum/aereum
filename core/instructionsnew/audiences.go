package instructionsnew

import "github.com/Aereum/aereum/core/crypto"

type Audience struct {
	token             crypto.PrivateKey
	submission        crypto.PrivateKey
	moderation        crypto.PrivateKey
	readCipher        []byte
	audienceKeyCipher []byte
	submitKeyCipher   []byte
	moderateKeyCipher []byte
}

func (a *Audience) SealedToken() []byte {
	cipher := crypto.CipherFromKey(a.audienceKeyCipher)
	return cipher.Seal(a.token.ToBytes())
}

func (a *Audience) SealedSubmission() []byte {
	cipher := crypto.CipherFromKey(a.submitKeyCipher)
	return cipher.Seal(a.submission.ToBytes())
}

func (a *Audience) SealedModeration() []byte {
	cipher := crypto.CipherFromKey(a.moderateKeyCipher)
	return cipher.Seal(a.moderation.ToBytes())
}

func (a *Audience) ReadTokenCiphers(members []crypto.PublicKey) TokenCiphers {
	readTokens := make(TokenCiphers, 0)
	var err error
	for _, member := range members {
		tc := TokenCipher{
			token: member.ToBytes(),
		}
		tc.cipher, err = member.Encrypt(a.readCipher)
		if err == nil {
			readTokens = append(readTokens, tc)
		}
	}
	return readTokens
}

func (a *Audience) SubmitTokenCiphers(members []crypto.PublicKey) TokenCiphers {
	readTokens := make(TokenCiphers, 0)
	var err error
	for _, member := range members {
		tc := TokenCipher{
			token: member.ToBytes(),
		}
		tc.cipher, err = member.Encrypt(a.submitKeyCipher)
		if err == nil {
			readTokens = append(readTokens, tc)
		}
	}
	return readTokens
}

func (a *Audience) ModerateTokenCiphers(members []crypto.PublicKey) TokenCiphers {
	readTokens := make(TokenCiphers, 0)
	var err error
	for _, member := range members {
		tc := TokenCipher{
			token: member.ToBytes(),
		}
		tc.cipher, err = member.Encrypt(a.moderateKeyCipher)
		if err == nil {
			readTokens = append(readTokens, tc)
		}
	}
	return readTokens
}

func NewAudience() *Audience {
	audience := Audience{}
	_, audience.token = crypto.RandomAsymetricKey()
	_, audience.submission = crypto.RandomAsymetricKey()
	_, audience.moderation = crypto.RandomAsymetricKey()
	audience.readCipher = crypto.NewCipherKey()
	audience.audienceKeyCipher = crypto.NewCipherKey()
	audience.submitKeyCipher = crypto.NewCipherKey()
	audience.moderateKeyCipher = crypto.NewCipherKey()
	return &audience
}

type CreateAudience struct {
	authored      *authoredInstruction
	audience      []byte
	submission    []byte
	moderation    []byte
	audienceKey   []byte
	submissionKey []byte
	moderationKey []byte
	flag          byte
	description   string
}

func (audience *CreateAudience) Kind() byte {
	return iCreateAudience
}

func (audience *CreateAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(audience.audience, &bytes)
	PutByteArray(audience.submission, &bytes)
	PutByteArray(audience.moderation, &bytes)
	PutByteArray(audience.audienceKey, &bytes)
	PutByteArray(audience.submissionKey, &bytes)
	PutByteArray(audience.moderationKey, &bytes)
	bytes = append(bytes, audience.flag)
	PutString(audience.description, &bytes)
	return bytes
}

func (audience *CreateAudience) Serialize() []byte {
	return audience.authored.serialize(iCreateAudience, audience.serializeBulk())
}

func ParseCreateAudience(data []byte) *CreateAudience {
	if data[0] != 0 || data[1] != iCreateAudience {
		return nil
	}
	audience := CreateAudience{
		authored: &authoredInstruction{},
	}
	position := audience.authored.parseHead(data)
	audience.audience, position = ParseByteArray(data, position)
	audience.submission, position = ParseByteArray(data, position)
	audience.moderation, position = ParseByteArray(data, position)
	audience.audienceKey, position = ParseByteArray(data, position)
	audience.submissionKey, position = ParseByteArray(data, position)
	audience.moderationKey, position = ParseByteArray(data, position)
	audience.flag, position = ParseByte(data, position)
	audience.description, position = ParseString(data, position)
	if audience.authored.parseTail(data, position) {
		return &audience
	}
	return nil
}

type JoinAudience struct {
	authored     *authoredInstruction
	audience     []byte
	presentation string
}

func (audience *JoinAudience) Kind() byte {
	return iJoinAudience
}

func (audience *JoinAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(audience.audience, &bytes)
	PutString(audience.presentation, &bytes)
	return bytes
}

func (audience *JoinAudience) Serialize() []byte {
	return audience.authored.serialize(iJoinAudience, audience.serializeBulk())
}

func ParseJoinAudience(data []byte) *JoinAudience {
	if data[0] != 0 || data[1] != iJoinAudience {
		return nil
	}
	audience := JoinAudience{
		authored: &authoredInstruction{},
	}
	position := audience.authored.parseHead(data)
	audience.audience, position = ParseByteArray(data, position)
	audience.presentation, position = ParseString(data, position)
	if audience.authored.parseTail(data, position) {
		return &audience
	}
	return nil
}

type AcceptJoinAudience struct {
	authored *authoredInstruction
	audience []byte
	member   []byte
	read     []byte
	submit   []byte
	moderate []byte
}

func (accept *AcceptJoinAudience) Kind() byte {
	return iJoinAudience
}

func (accept *AcceptJoinAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(accept.audience, &bytes)
	PutByteArray(accept.member, &bytes)
	PutByteArray(accept.read, &bytes)
	PutByteArray(accept.submit, &bytes)
	PutByteArray(accept.moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinAudience) Serialize() []byte {
	return accept.authored.serialize(iAcceptJoinRequest, accept.serializeBulk())
}

func ParseAcceptJoinAudience(data []byte) *AcceptJoinAudience {
	if data[0] != 0 || data[1] != iAcceptJoinRequest {
		return nil
	}
	accept := AcceptJoinAudience{
		authored: &authoredInstruction{},
	}
	position := accept.authored.parseHead(data)
	accept.audience, position = ParseByteArray(data, position)
	accept.member, position = ParseByteArray(data, position)
	accept.read, position = ParseByteArray(data, position)
	accept.submit, position = ParseByteArray(data, position)
	accept.moderate, position = ParseByteArray(data, position)
	if accept.authored.parseTail(data, position) {
		return &accept
	}
	return nil
}

type TokenCipher struct {
	token  []byte
	cipher []byte
}

type TokenCiphers []TokenCipher

type UpdateAudience struct {
	authored      *authoredInstruction
	audience      []byte
	submission    []byte
	moderation    []byte
	audienceKey   []byte
	submissionKey []byte
	moderationKey []byte
	flag          byte
	description   string
	readMembers   TokenCiphers
	subMembers    TokenCiphers
	modMembers    TokenCiphers
}

func (update *UpdateAudience) Kind() byte {
	return iUpdateAudience
}

func (update *UpdateAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(update.audience, &bytes)
	PutByteArray(update.submission, &bytes)
	PutByteArray(update.moderation, &bytes)
	PutByteArray(update.audienceKey, &bytes)
	PutByteArray(update.submissionKey, &bytes)
	PutByteArray(update.moderationKey, &bytes)
	PutByte(update.flag, &bytes)
	PutString(update.description, &bytes)
	PutTokenCiphers(update.readMembers, &bytes)
	PutTokenCiphers(update.subMembers, &bytes)
	PutTokenCiphers(update.modMembers, &bytes)
	return bytes
}

func (update *UpdateAudience) Serialize() []byte {
	return update.authored.serialize(iUpdateAudience, update.serializeBulk())
}

func ParseUpdateAudience(data []byte) *UpdateAudience {
	if data[0] != 0 || data[1] != iUpdateAudience {
		return nil
	}
	update := UpdateAudience{
		authored: &authoredInstruction{},
	}
	position := update.authored.parseHead(data)
	update.audience, position = ParseByteArray(data, position)
	update.submission, position = ParseByteArray(data, position)
	update.moderation, position = ParseByteArray(data, position)
	update.audienceKey, position = ParseByteArray(data, position)
	update.submissionKey, position = ParseByteArray(data, position)
	update.moderationKey, position = ParseByteArray(data, position)
	update.flag, position = ParseByte(data, position)
	update.description, position = ParseString(data, position)
	update.readMembers, position = ParseTokenCiphers(data, position)
	update.subMembers, position = ParseTokenCiphers(data, position)
	update.modMembers, position = ParseTokenCiphers(data, position)
	if update.authored.parseTail(data, position) {
		return &update
	}
	return nil
}
