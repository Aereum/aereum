package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
	"github.com/Aereum/aereum/core/util"
)

type CreateStage struct {
	authored    *authoredInstruction
	audience    crypto.Token
	submission  crypto.Token
	moderation  crypto.Token
	flag        byte
	description string
}

func (a *CreateStage) Epoch() uint64 {
	return a.authored.epoch
}

func (stage *CreateStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(stage.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.HashToken(stage.audience)
	if stage := v.GetAudienceKeys(audienceHash); stage != nil {
		return false
	}
	v.AddFeeCollected(stage.authored.fee)
	stageKeys := store.StageKeys{
		Moderate: stage.moderation,
		Submit:   stage.submission,
		Stage:    stage.audience,
		Flag:     stage.flag,
	}
	return v.SetNewAudience(audienceHash, stageKeys)
}

func (stage *CreateStage) Payments() *Payment {
	return stage.authored.payments()
}

func (stage *CreateStage) Kind() byte {
	return iCreateAudience
}

func (stage *CreateStage) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(stage.audience, &bytes)
	util.PutToken(stage.submission, &bytes)
	util.PutToken(stage.moderation, &bytes)
	bytes = append(bytes, stage.flag)
	util.PutString(stage.description, &bytes)
	return bytes
}

func (stage *CreateStage) Serialize() []byte {
	return stage.authored.serialize(iCreateAudience, stage.serializeBulk())
}

func ParseCreateStage(data []byte) *CreateStage {
	if data[0] != 0 || data[1] != iCreateAudience {
		return nil
	}
	stage := CreateStage{
		authored: &authoredInstruction{},
	}
	position := stage.authored.parseHead(data)
	stage.audience, position = util.ParseToken(data, position)
	stage.submission, position = util.ParseToken(data, position)
	stage.moderation, position = util.ParseToken(data, position)
	stage.flag, position = util.ParseByte(data, position)
	stage.description, position = util.ParseString(data, position)
	if stage.authored.parseTail(data, position) {
		return &stage
	}
	return nil
}

type JoinStage struct {
	authored     *authoredInstruction
	audience     crypto.Token
	diffHellKey  crypto.Token
	presentation string
}

func (a *JoinStage) Epoch() uint64 {
	return a.authored.epoch
}

func (join *JoinStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(join.authored.authorHash()) {
		return false
	}
	if keys := v.GetAudienceKeys(crypto.HashToken(join.audience)); keys == nil {
		return false
	}
	v.AddFeeCollected(join.authored.fee)
	return true
}

func (join *JoinStage) Payments() *Payment {
	return join.authored.payments()
}

func (join *JoinStage) Kind() byte {
	return iJoinAudience
}

func (join *JoinStage) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(join.audience, &bytes)
	util.PutToken(join.diffHellKey, &bytes)
	util.PutString(join.presentation, &bytes)
	return bytes
}

func (stage *JoinStage) Serialize() []byte {
	return stage.authored.serialize(iJoinAudience, stage.serializeBulk())
}

func ParseJoinStage(data []byte) *JoinStage {
	if data[0] != 0 || data[1] != iJoinAudience {
		return nil
	}
	stage := JoinStage{
		authored: &authoredInstruction{},
	}
	position := stage.authored.parseHead(data)
	stage.audience, position = util.ParseToken(data, position)
	stage.diffHellKey, position = util.ParseToken(data, position)
	stage.presentation, position = util.ParseString(data, position)
	if stage.authored.parseTail(data, position) {
		return &stage
	}
	return nil
}

type AcceptJoinStage struct {
	authored     *authoredInstruction
	stage        crypto.Token
	member       crypto.Token
	diffHellKey  crypto.Token
	read         []byte
	submit       []byte
	moderate     []byte
	modSignature crypto.Signature
}

func (a *AcceptJoinStage) Epoch() uint64 {
	return a.authored.epoch
}

func (accept *AcceptJoinStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(accept.authored.authorHash()) {
		return false
	}
	keys := v.GetAudienceKeys(crypto.HashToken(accept.stage))
	if keys == nil || keys.Moderate == crypto.ZeroToken {
		return false
	}
	if !keys.Moderate.Verify(accept.serializeModBulk(), accept.modSignature) {
		return false
	}
	//hashed := crypto.Hasher(accept.Serialize())
	//if bytes.Equal(keys[0:crypto.Size], hashed[:]) {
	v.AddFeeCollected(accept.authored.fee)
	return true
	//}
	//return false
}

func (accept *AcceptJoinStage) Payments() *Payment {
	return accept.authored.payments()
}

func (accept *AcceptJoinStage) Kind() byte {
	return iJoinAudience
}

func (accept *AcceptJoinStage) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(accept.stage, &bytes)
	util.PutToken(accept.member, &bytes)
	util.PutToken(accept.diffHellKey, &bytes)
	util.PutByteArray(accept.read, &bytes)
	util.PutByteArray(accept.submit, &bytes)
	util.PutByteArray(accept.moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinStage) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutSignature(accept.modSignature, &bytes)
	return bytes
}

func (accept *AcceptJoinStage) Serialize() []byte {
	return accept.authored.serialize(iAcceptJoinRequest, accept.serializeBulk())
}

func ParseAcceptJoinStage(data []byte) *AcceptJoinStage {
	if data[0] != 0 || data[1] != iAcceptJoinRequest {
		return nil
	}
	accept := AcceptJoinStage{
		authored: &authoredInstruction{},
	}
	position := accept.authored.parseHead(data)
	accept.stage, position = util.ParseToken(data, position)
	accept.member, position = util.ParseToken(data, position)
	accept.diffHellKey, position = util.ParseToken(data, position)
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

type UpdateStage struct {
	authored     *authoredInstruction
	stage        crypto.Token // existing audience public token - it doesn't change
	submission   crypto.Token // new submission public token
	moderation   crypto.Token // new moderation public token
	diffHellKey  crypto.Token
	flag         byte
	description  string
	readMembers  TokenCiphers
	subMembers   TokenCiphers
	modMembers   TokenCiphers
	audSignature crypto.Signature
}

func (a *UpdateStage) Epoch() uint64 {
	return a.authored.epoch
}

func (update *UpdateStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.authored.authorHash()) {
		return false
	}
	hashed := crypto.HashToken(update.stage)
	stageKeys := store.StageKeys{
		Moderate: update.moderation,
		Submit:   update.submission,
		Stage:    update.stage,
		Flag:     update.flag,
	}
	if v.UpdateAudience(hashed, stageKeys) {
		v.AddFeeCollected(update.authored.fee)
		return true
	}
	return false
}

func (update *UpdateStage) Payments() *Payment {
	return update.authored.payments()
}

func (update *UpdateStage) Kind() byte {
	return iUpdateAudience
}

func (update *UpdateStage) serializeAudBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(update.stage, &bytes)
	util.PutToken(update.submission, &bytes)
	util.PutToken(update.moderation, &bytes)
	util.PutToken(update.diffHellKey, &bytes)
	util.PutByte(update.flag, &bytes)
	util.PutString(update.description, &bytes)
	putTokenCiphers(update.readMembers, &bytes)
	putTokenCiphers(update.subMembers, &bytes)
	putTokenCiphers(update.modMembers, &bytes)
	return bytes
}

func (update *UpdateStage) serializeBulk() []byte {
	bytes := update.serializeAudBulk()
	util.PutSignature(update.audSignature, &bytes)
	return bytes
}

func (update *UpdateStage) Serialize() []byte {
	return update.authored.serialize(iUpdateAudience, update.serializeBulk())
}

func ParseUpdateStage(data []byte) *UpdateStage {
	if data[0] != 0 || data[1] != iUpdateAudience {
		return nil
	}
	update := UpdateStage{
		authored: &authoredInstruction{},
	}
	position := update.authored.parseHead(data)
	update.stage, position = util.ParseToken(data, position)
	update.submission, position = util.ParseToken(data, position)
	update.moderation, position = util.ParseToken(data, position)
	update.diffHellKey, position = util.ParseToken(data, position)
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
