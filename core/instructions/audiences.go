package instructions

import (
	"bytes"

	"github.com/Aereum/aereum/core/crypto"
)

type Audience struct {
	token             *crypto.PrivateKey
	submission        *crypto.PrivateKey
	moderation        *crypto.PrivateKey
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
		// tc.cipher, err = member.Encrypt(a.readCipher)
		tc.cipher, err = member.Encrypt(a.audienceKeyCipher)
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
	_, audtoken := crypto.RandomAsymetricKey()
	_, subtoken := crypto.RandomAsymetricKey()
	_, modtoken := crypto.RandomAsymetricKey()
	audience.token = &audtoken
	audience.submission = &subtoken
	audience.moderation = &modtoken
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

func (audience *CreateAudience) Validate(block *Block) bool {
	if !block.validator.HasMember(audience.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(audience.audience)
	if block.validator.GetAudienceKeys(audienceHash) != nil {
		return false
	}
	keys := append(audience.submission, audience.moderation...)
	block.FeesCollected += audience.authored.fee
	return block.SetNewAudience(audienceHash, keys)
}

func (audience *CreateAudience) Payments() *Payment {
	return audience.authored.payments()
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

func (join *JoinAudience) Validate(block *Block) bool {
	if !block.validator.HasMember(join.authored.authorHash()) {
		return false
	}
	if block.validator.GetAudienceKeys(crypto.Hasher(join.audience)) != nil {
		return false
	}
	block.FeesCollected += join.authored.fee
	return true
}

func (join *JoinAudience) Payments() *Payment {
	return join.authored.payments()
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
	authored     *authoredInstruction
	audience     []byte
	member       []byte
	read         []byte
	submit       []byte
	moderate     []byte
	modSignature []byte
}

func (accept *AcceptJoinAudience) Validate(block *Block) bool {
	if !block.validator.HasMember(accept.authored.authorHash()) {
		return false
	}
	keys := block.validator.GetAudienceKeys(crypto.Hasher(accept.audience))
	if keys == nil {
		return false
	}
	modPublic, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	hashed := crypto.Hasher(accept.serializeModBulk())
	if !modPublic.Verify(hashed[:], accept.modSignature) {
		return false
	}
	if bytes.Equal(keys[0:crypto.Size], hashed[:]) {
		block.FeesCollected += accept.authored.fee
		return true
	}
	return false
}

func (accept *AcceptJoinAudience) Payments() *Payment {
	return accept.authored.payments()
}

func (accept *AcceptJoinAudience) Kind() byte {
	return iJoinAudience
}

func (accept *AcceptJoinAudience) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(accept.audience, &bytes)
	PutByteArray(accept.member, &bytes)
	PutByteArray(accept.read, &bytes)
	PutByteArray(accept.submit, &bytes)
	PutByteArray(accept.moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinAudience) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	PutByteArray(accept.modSignature, &bytes)
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
	accept.modSignature, position = ParseByteArray(data, position)
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
	audience      []byte // existing audience public token - it doesn't change
	submission    []byte // new submission public token
	moderation    []byte // new moderation public token
	submissionKey []byte // ciphered private submission key
	moderationKey []byte // ciphered private moderation key
	flag          byte
	description   string
	readMembers   TokenCiphers
	subMembers    TokenCiphers
	modMembers    TokenCiphers
	audSignature  []byte
}

func (update *UpdateAudience) Validate(block *Block) bool {
	if !block.validator.HasMember(update.authored.authorHash()) {
		return false
	}
	hashed := crypto.Hasher(update.audience)
	newKeys := append(update.submission, update.moderation...)
	if block.UpdateAudience(hashed, newKeys) {
		block.FeesCollected += update.authored.fee
		return true
	}
	return false
}

func (update *UpdateAudience) Payments() *Payment {
	return update.authored.payments()
}

func (update *UpdateAudience) Kind() byte {
	return iUpdateAudience
}

func (update *UpdateAudience) serializeAudBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(update.audience, &bytes)
	PutByteArray(update.submission, &bytes)
	PutByteArray(update.moderation, &bytes)
	PutByteArray(update.submissionKey, &bytes)
	PutByteArray(update.moderationKey, &bytes)
	PutByte(update.flag, &bytes)
	PutString(update.description, &bytes)
	PutTokenCiphers(update.readMembers, &bytes)
	PutTokenCiphers(update.subMembers, &bytes)
	PutTokenCiphers(update.modMembers, &bytes)
	return bytes
}

func (update *UpdateAudience) serializeBulk() []byte {
	bytes := update.serializeAudBulk()
	PutByteArray(update.audSignature, &bytes)
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
	update.submissionKey, position = ParseByteArray(data, position)
	update.moderationKey, position = ParseByteArray(data, position)
	update.flag, position = ParseByte(data, position)
	update.description, position = ParseString(data, position)
	update.readMembers, position = ParseTokenCiphers(data, position)
	update.subMembers, position = ParseTokenCiphers(data, position)
	update.modMembers, position = ParseTokenCiphers(data, position)
	// hashed := crypto.Hasher(data[0:position])
	update.audSignature, position = ParseByteArray(data, position)
	// audPublic, err := crypto.PublicKeyFromBytes(update.authored.author)
	// if err != nil {
	// 	return nil
	// }
	// if !audPublic.Verify(hashed[:], update.audSignature) {
	// 	return nil
	// }
	if update.authored.parseTail(data, position) {
		return &update
	}
	return nil
}
