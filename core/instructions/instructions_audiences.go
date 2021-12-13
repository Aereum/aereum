package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type CreateAudience struct {
	authored    *authoredInstruction
	audience    crypto.Token
	submission  crypto.Token
	moderation  crypto.Token
	flag        byte
	description string
}

func (a *CreateAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (audience *CreateAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(audience.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.HashToken(audience.audience)
	if ok, _, _, _ := v.GetAudienceKeys(audienceHash); !ok {
		return false
	}
	v.AddFeeCollected(audience.authored.fee)
	return v.SetNewAudience(audienceHash, audience.moderation, audience.submission, audience.flag)
}

func (audience *CreateAudience) Payments() *Payment {
	return audience.authored.payments()
}

func (audience *CreateAudience) Kind() byte {
	return iCreateAudience
}

func (audience *CreateAudience) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(audience.audience, &bytes)
	util.PutToken(audience.submission, &bytes)
	util.PutToken(audience.moderation, &bytes)
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
	audience.audience, position = util.ParseToken(data, position)
	audience.submission, position = util.ParseToken(data, position)
	audience.moderation, position = util.ParseToken(data, position)
	audience.flag, position = util.ParseByte(data, position)
	audience.description, position = util.ParseString(data, position)
	if audience.authored.parseTail(data, position) {
		return &audience
	}
	return nil
}

type JoinAudience struct {
	authored     *authoredInstruction
	audience     crypto.Token
	key          crypto.Token
	presentation string
}

func (a *JoinAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (join *JoinAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(join.authored.authorHash()) {
		return false
	}
	if ok, _, _, _ := v.GetAudienceKeys(crypto.HashToken(join.audience)); !ok {
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
	util.PutToken(audience.audience, &bytes)
	util.PutToken(audience.key, &bytes)
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
	audience.audience, position = util.ParseToken(data, position)
	audience.key, position = util.ParseToken(data, position)
	audience.presentation, position = util.ParseString(data, position)
	if audience.authored.parseTail(data, position) {
		return &audience
	}
	return nil
}

type AcceptJoinAudience struct {
	authored     *authoredInstruction
	audience     crypto.Token
	member       crypto.Token
	key          crypto.Token
	read         []byte
	submit       []byte
	moderate     []byte
	modSignature crypto.Signature
}

func (a *AcceptJoinAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (accept *AcceptJoinAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(accept.authored.authorHash()) {
		return false
	}
	ok, moderate, _, _ := v.GetAudienceKeys(crypto.HashToken(accept.audience))
	if !ok || moderate == crypto.ZeroToken {
		return false
	}
	if !moderate.Verify(accept.serializeModBulk(), accept.modSignature) {
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
	util.PutToken(accept.audience, &bytes)
	util.PutToken(accept.member, &bytes)
	util.PutToken(accept.key, &bytes)
	util.PutByteArray(accept.read, &bytes)
	util.PutByteArray(accept.submit, &bytes)
	util.PutByteArray(accept.moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinAudience) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutSignature(accept.modSignature, &bytes)
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
	accept.audience, position = util.ParseToken(data, position)
	accept.member, position = util.ParseToken(data, position)
	accept.key, position = util.ParseToken(data, position)
	accept.read, position = util.ParseByteArray(data, position)
	accept.submit, position = util.ParseByteArray(data, position)
	accept.moderate, position = util.ParseByteArray(data, position)
	accept.modSignature, position = util.ParseSignature(data, position)
	if accept.authored.parseTail(data, position) {
		return &accept
	}
	return nil
}

type TokenCipher struct {
	token  crypto.Token
	cipher []byte
}

type TokenCiphers []TokenCipher

func putTokenCipher(tc TokenCipher, data *[]byte) {
	util.PutToken(tc.token, data)
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
	tc.token, position = util.ParseToken(data, position)
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
	authored     *authoredInstruction
	audience     crypto.Token // existing audience public token - it doesn't change
	submission   crypto.Token // new submission public token
	moderation   crypto.Token // new moderation public token
	key          crypto.Token
	flag         byte
	description  string
	readMembers  TokenCiphers
	subMembers   TokenCiphers
	modMembers   TokenCiphers
	audSignature crypto.Signature
}

func (a *UpdateAudience) Epoch() uint64 {
	return a.authored.epoch
}

func (update *UpdateAudience) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.authored.authorHash()) {
		return false
	}
	hashed := crypto.HashToken(update.audience)
	if v.UpdateAudience(hashed, update.moderation, update.submission, update.flag) {
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
	util.PutToken(update.audience, &bytes)
	util.PutToken(update.submission, &bytes)
	util.PutToken(update.moderation, &bytes)
	util.PutToken(update.key, &bytes)
	util.PutByte(update.flag, &bytes)
	util.PutString(update.description, &bytes)
	putTokenCiphers(update.readMembers, &bytes)
	putTokenCiphers(update.subMembers, &bytes)
	putTokenCiphers(update.modMembers, &bytes)
	return bytes
}

func (update *UpdateAudience) serializeBulk() []byte {
	bytes := update.serializeAudBulk()
	util.PutSignature(update.audSignature, &bytes)
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
	update.audience, position = util.ParseToken(data, position)
	update.submission, position = util.ParseToken(data, position)
	update.moderation, position = util.ParseToken(data, position)
	update.key, position = util.ParseToken(data, position)
	update.flag, position = util.ParseByte(data, position)
	update.description, position = util.ParseString(data, position)
	update.readMembers, position = parseTokenCiphers(data, position)
	update.subMembers, position = parseTokenCiphers(data, position)
	update.modMembers, position = parseTokenCiphers(data, position)
	update.audSignature, position = util.ParseSignature(data, position)
	if update.authored.parseTail(data, position) {
		return &update
	}
	return nil
}
