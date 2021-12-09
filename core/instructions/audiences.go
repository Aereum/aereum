package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type Audience struct {
	Token             *crypto.PrivateKey
	Submission        *crypto.PrivateKey
	Moderation        *crypto.PrivateKey
	AudienceKeyCipher []byte
	SubmitKeyCipher   []byte
	ModerateKeyCipher []byte
}

func (a *Audience) SealedToken() []byte {
	cipher := crypto.CipherFromKey(a.AudienceKeyCipher)
	return cipher.Seal(a.Token.ToBytes())
}

func (a *Audience) SealedSubmission() []byte {
	cipher := crypto.CipherFromKey(a.SubmitKeyCipher)
	return cipher.Seal(a.Submission.ToBytes())
}

func (a *Audience) SealedModeration() []byte {
	cipher := crypto.CipherFromKey(a.ModerateKeyCipher)
	return cipher.Seal(a.Moderation.ToBytes())
}

func (a *Audience) ReadTokenCiphers(members []crypto.PublicKey) TokenCiphers {
	readTokens := make(TokenCiphers, 0)
	var err error
	for _, member := range members {
		tc := TokenCipher{
			token: member.ToBytes(),
		}
		// tc.cipher, err = member.Encrypt(a.readCipher)
		tc.cipher, err = member.Encrypt(a.AudienceKeyCipher)
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
		tc.cipher, err = member.Encrypt(a.SubmitKeyCipher)
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
		tc.cipher, err = member.Encrypt(a.ModerateKeyCipher)
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
	audience.Token = &audtoken
	audience.Submission = &subtoken
	audience.Moderation = &modtoken
	audience.AudienceKeyCipher = crypto.NewCipherKey()
	audience.SubmitKeyCipher = crypto.NewCipherKey()
	audience.ModerateKeyCipher = crypto.NewCipherKey()
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

func (a *CreateAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (audience *CreateAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(audience.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(audience.audience)
	if v.GetAudienceKeys(audienceHash) != nil {
		return false
	}
	keys := append(audience.submission, audience.moderation...)
	v.AddFeeCollected(audience.authored.fee)
	return v.SetNewAudience(audienceHash, keys)
}

func (audience *CreateAudience) Payments() *Payment {
	return audience.authored.payments()
}

func (audience *CreateAudience) Kind() byte {
	return iCreateAudience
}

func (audience *CreateAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(audience.audience, &bytes)
	util.PutByteArray(audience.submission, &bytes)
	util.PutByteArray(audience.moderation, &bytes)
	util.PutByteArray(audience.audienceKey, &bytes)
	util.PutByteArray(audience.submissionKey, &bytes)
	util.PutByteArray(audience.moderationKey, &bytes)
	bytes = append(bytes, audience.flag)
	util.PutString(audience.description, &bytes)
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
	audience.audience, position = util.ParseByteArray(data, position)
	audience.submission, position = util.ParseByteArray(data, position)
	audience.moderation, position = util.ParseByteArray(data, position)
	audience.audienceKey, position = util.ParseByteArray(data, position)
	audience.submissionKey, position = util.ParseByteArray(data, position)
	audience.moderationKey, position = util.ParseByteArray(data, position)
	audience.flag, position = util.ParseByte(data, position)
	audience.description, position = util.ParseString(data, position)
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

func (a *JoinAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (join *JoinAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(join.authored.authorHash()) {
		return false
	}
	if v.GetAudienceKeys(crypto.Hasher(join.audience)) != nil {
		return false
	}
	v.AddFeeCollected(join.authored.fee)
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
	util.PutByteArray(audience.audience, &bytes)
	util.PutString(audience.presentation, &bytes)
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
	audience.audience, position = util.ParseByteArray(data, position)
	audience.presentation, position = util.ParseString(data, position)
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

func (a *AcceptJoinAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (accept *AcceptJoinAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(accept.authored.authorHash()) {
		return false
	}
	keys := v.GetAudienceKeys(crypto.Hasher(accept.audience))
	if keys == nil {
		return false
	}
	modPublic, err := crypto.PublicKeyFromBytes(keys[crypto.PublicKeySize : 2*crypto.PublicKeySize])
	if err != nil {
		return false
	}
	if !modPublic.Verify(accept.serializeModBulk(), accept.modSignature) {
		return false
	}
	//hashed := crypto.Hasher(accept.Serialize())
	//if bytes.Equal(keys[0:crypto.Size], hashed[:]) {
	v.AddFeeCollected(accept.authored.fee)
	return true
	//}
	//return false
}

func (accept *AcceptJoinAudience) Payments() *Payment {
	return accept.authored.payments()
}

func (accept *AcceptJoinAudience) Kind() byte {
	return iJoinAudience
}

func (accept *AcceptJoinAudience) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(accept.audience, &bytes)
	util.PutByteArray(accept.member, &bytes)
	util.PutByteArray(accept.read, &bytes)
	util.PutByteArray(accept.submit, &bytes)
	util.PutByteArray(accept.moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinAudience) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutByteArray(accept.modSignature, &bytes)
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
	accept.audience, position = util.ParseByteArray(data, position)
	accept.member, position = util.ParseByteArray(data, position)
	accept.read, position = util.ParseByteArray(data, position)
	accept.submit, position = util.ParseByteArray(data, position)
	accept.moderate, position = util.ParseByteArray(data, position)
	accept.modSignature, position = util.ParseByteArray(data, position)
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

func putTokenCipher(tc TokenCipher, data *[]byte) {
	util.PutByteArray(tc.token, data)
	util.PutByteArray(tc.cipher, data)
}

func putTokenCiphers(tcs TokenCiphers, data *[]byte) {
	if len(tcs) == 0 {
		*data = append(*data, 0, 0)
		return
	}
	maxLen := len(tcs)
	if len(tcs) > 1<<16-1 {
		maxLen = 1 << 16
	}
	*data = append(*data, byte(maxLen), byte(maxLen>>8))
	for n := 0; n < maxLen; n++ {
		putTokenCipher(tcs[n], data)
	}
}

func parseTokenCipher(data []byte, position int) (TokenCipher, int) {
	tc := TokenCipher{}
	if position+1 >= len(data) {
		return tc, position
	}
	tc.token, position = util.ParseByteArray(data, position)
	tc.cipher, position = util.ParseByteArray(data, position)
	return tc, position
}

func parseTokenCiphers(data []byte, position int) (TokenCiphers, int) {
	if position+1 >= len(data) {
		return TokenCiphers{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	position += 2
	if length == 0 {
		return TokenCiphers{}, position + 2
	}
	if position+length+2 > len(data) {
		return TokenCiphers{}, position + length + 2
	}
	tcs := make(TokenCiphers, length)
	for n := 0; n < length; n++ {
		tcs[n], position = parseTokenCipher(data, position)
	}
	return tcs, position
}

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

func (a *UpdateAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (update *UpdateAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.authored.authorHash()) {
		return false
	}
	hashed := crypto.Hasher(update.audience)
	newKeys := append(update.submission, update.moderation...)
	if v.UpdateAudience(hashed, newKeys) {
		v.AddFeeCollected(update.authored.fee)
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
	util.PutByteArray(update.audience, &bytes)
	util.PutByteArray(update.submission, &bytes)
	util.PutByteArray(update.moderation, &bytes)
	util.PutByteArray(update.submissionKey, &bytes)
	util.PutByteArray(update.moderationKey, &bytes)
	util.PutByte(update.flag, &bytes)
	util.PutString(update.description, &bytes)
	putTokenCiphers(update.readMembers, &bytes)
	putTokenCiphers(update.subMembers, &bytes)
	putTokenCiphers(update.modMembers, &bytes)
	return bytes
}

func (update *UpdateAudience) serializeBulk() []byte {
	bytes := update.serializeAudBulk()
	util.PutByteArray(update.audSignature, &bytes)
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
	update.audience, position = util.ParseByteArray(data, position)
	update.submission, position = util.ParseByteArray(data, position)
	update.moderation, position = util.ParseByteArray(data, position)
	update.submissionKey, position = util.ParseByteArray(data, position)
	update.moderationKey, position = util.ParseByteArray(data, position)
	update.flag, position = util.ParseByte(data, position)
	update.description, position = util.ParseString(data, position)
	update.readMembers, position = parseTokenCiphers(data, position)
	update.subMembers, position = parseTokenCiphers(data, position)
	update.modMembers, position = parseTokenCiphers(data, position)
	// hashed := crypto.Hasher(data[0:position])
	update.audSignature, position = util.ParseByteArray(data, position)
	// audPublic, err := crypto.PublicKeyFromBytes(update.authored.author)
	// if err != nil {
	// 	return nil
	// }
	// if !audPublic.Verify(hashed[:], update.audSignature) {
	// 	return nil
	// }
	// fmt.Printf(string(update.audSignature))
	if update.authored.parseTail(data, position) {
		return &update
	}
	return nil
}
