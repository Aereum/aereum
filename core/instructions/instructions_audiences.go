package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
	"github.com/Aereum/aereum/core/util"
)

type CreateStage struct {
	Authored    *AuthoredInstruction
	Audience    crypto.Token
	Submission  crypto.Token
	Moderation  crypto.Token
	Flag        byte
	Description string
}

func (a *CreateStage) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *CreateStage) Epoch() uint64 {
	return a.Authored.epoch
}

func (stage *CreateStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(stage.Authored.authorHash()) {
		return false
	}
	audienceHash := crypto.HashToken(stage.Audience)
	if stage := v.GetAudienceKeys(audienceHash); stage != nil {
		return false
	}
	v.AddFeeCollected(stage.Authored.Fee)
	stageKeys := store.StageKeys{
		Moderate: stage.Moderation,
		Submit:   stage.Submission,
		Stage:    stage.Audience,
		Flag:     stage.Flag,
	}
	return v.SetNewAudience(audienceHash, stageKeys)
}

func (stage *CreateStage) Payments() *Payment {
	return stage.Authored.payments()
}

func (stage *CreateStage) Kind() byte {
	return ICreateAudience
}

func (stage *CreateStage) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(stage.Audience, &bytes)
	util.PutToken(stage.Submission, &bytes)
	util.PutToken(stage.Moderation, &bytes)
	bytes = append(bytes, stage.Flag)
	util.PutString(stage.Description, &bytes)
	return bytes
}

func (stage *CreateStage) Serialize() []byte {
	return stage.Authored.serialize(ICreateAudience, stage.serializeBulk())
}

func ParseCreateStage(data []byte) *CreateStage {
	if data[0] != 0 || data[1] != ICreateAudience {
		return nil
	}
	stage := CreateStage{
		Authored: &AuthoredInstruction{},
	}
	position := stage.Authored.parseHead(data)
	stage.Audience, position = util.ParseToken(data, position)
	stage.Submission, position = util.ParseToken(data, position)
	stage.Moderation, position = util.ParseToken(data, position)
	stage.Flag, position = util.ParseByte(data, position)
	stage.Description, position = util.ParseString(data, position)
	if stage.Authored.parseTail(data, position) {
		return &stage
	}
	return nil
}

type JoinStage struct {
	Authored     *AuthoredInstruction
	Audience     crypto.Token
	DiffHellKey  crypto.Token
	Presentation string
}

func (a *JoinStage) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *JoinStage) Epoch() uint64 {
	return a.Authored.epoch
}

func (join *JoinStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(join.Authored.authorHash()) {
		return false
	}
	if keys := v.GetAudienceKeys(crypto.HashToken(join.Audience)); keys == nil {
		return false
	}
	v.AddFeeCollected(join.Authored.Fee)
	return true
}

func (join *JoinStage) Payments() *Payment {
	return join.Authored.payments()
}

func (join *JoinStage) Kind() byte {
	return IJoinAudience
}

func (join *JoinStage) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(join.Audience, &bytes)
	util.PutToken(join.DiffHellKey, &bytes)
	util.PutString(join.Presentation, &bytes)
	return bytes
}

func (stage *JoinStage) Serialize() []byte {
	return stage.Authored.serialize(IJoinAudience, stage.serializeBulk())
}

func ParseJoinStage(data []byte) *JoinStage {
	if data[0] != 0 || data[1] != IJoinAudience {
		return nil
	}
	stage := JoinStage{
		Authored: &AuthoredInstruction{},
	}
	position := stage.Authored.parseHead(data)
	stage.Audience, position = util.ParseToken(data, position)
	stage.DiffHellKey, position = util.ParseToken(data, position)
	stage.Presentation, position = util.ParseString(data, position)
	if stage.Authored.parseTail(data, position) {
		return &stage
	}
	return nil
}

type AcceptJoinStage struct {
	Authored     *AuthoredInstruction
	Stage        crypto.Token
	Member       crypto.Token
	DiffHellKey  crypto.Token
	Read         []byte
	Submit       []byte
	Moderate     []byte
	modSignature crypto.Signature
}

func (a *AcceptJoinStage) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *AcceptJoinStage) Epoch() uint64 {
	return a.Authored.epoch
}

func (accept *AcceptJoinStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(accept.Authored.authorHash()) {
		return false
	}
	keys := v.GetAudienceKeys(crypto.HashToken(accept.Stage))
	if keys == nil || keys.Moderate == crypto.ZeroToken {
		return false
	}
	if !keys.Moderate.Verify(accept.serializeModBulk(), accept.modSignature) {
		return false
	}
	//hashed := crypto.Hasher(accept.Serialize())
	//if bytes.Equal(keys[0:crypto.Size], hashed[:]) {
	v.AddFeeCollected(accept.Authored.Fee)
	return true
	//}
	//return false
}

func (accept *AcceptJoinStage) Payments() *Payment {
	return accept.Authored.payments()
}

func (accept *AcceptJoinStage) Kind() byte {
	return IJoinAudience
}

func (accept *AcceptJoinStage) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(accept.Stage, &bytes)
	util.PutToken(accept.Member, &bytes)
	util.PutToken(accept.DiffHellKey, &bytes)
	util.PutByteArray(accept.Read, &bytes)
	util.PutByteArray(accept.Submit, &bytes)
	util.PutByteArray(accept.Moderate, &bytes)
	return bytes
}

func (accept *AcceptJoinStage) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutSignature(accept.modSignature, &bytes)
	return bytes
}

func (accept *AcceptJoinStage) Serialize() []byte {
	return accept.Authored.serialize(IAcceptJoinRequest, accept.serializeBulk())
}

func ParseAcceptJoinStage(data []byte) *AcceptJoinStage {
	if data[0] != 0 || data[1] != IAcceptJoinRequest {
		return nil
	}
	accept := AcceptJoinStage{
		Authored: &AuthoredInstruction{},
	}
	position := accept.Authored.parseHead(data)
	accept.Stage, position = util.ParseToken(data, position)
	accept.Member, position = util.ParseToken(data, position)
	accept.DiffHellKey, position = util.ParseToken(data, position)
	accept.Read, position = util.ParseByteArray(data, position)
	accept.Submit, position = util.ParseByteArray(data, position)
	accept.Moderate, position = util.ParseByteArray(data, position)
	accept.modSignature, position = util.ParseSignature(data, position)
	if accept.Authored.parseTail(data, position) {
		return &accept
	}
	return nil
}

type TokenCipher struct {
	Token  crypto.Token
	Cipher []byte
}

type TokenCiphers []TokenCipher

func putTokenCipher(tc TokenCipher, data *[]byte) {
	util.PutToken(tc.Token, data)
	util.PutByteArray(tc.Cipher, data)
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
	tc.Token, position = util.ParseToken(data, position)
	tc.Cipher, position = util.ParseByteArray(data, position)
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
	Authored     *AuthoredInstruction
	Stage        crypto.Token // existing audience public token - it doesn't change
	Submission   crypto.Token // new submission public token
	Moderation   crypto.Token // new moderation public token
	DiffHellKey  crypto.Token
	Flag         byte
	Description  string
	ReadMembers  TokenCiphers
	SubMembers   TokenCiphers
	ModMembers   TokenCiphers
	audSignature crypto.Signature
}

func (a *UpdateStage) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *UpdateStage) Epoch() uint64 {
	return a.Authored.epoch
}

func (update *UpdateStage) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.Authored.authorHash()) {
		return false
	}
	hashed := crypto.HashToken(update.Stage)
	stageKeys := store.StageKeys{
		Moderate: update.Moderation,
		Submit:   update.Submission,
		Stage:    update.Stage,
		Flag:     update.Flag,
	}
	if v.UpdateAudience(hashed, stageKeys) {
		v.AddFeeCollected(update.Authored.Fee)
		return true
	}
	return false
}

func (update *UpdateStage) Payments() *Payment {
	return update.Authored.payments()
}

func (update *UpdateStage) Kind() byte {
	return IUpdateAudience
}

func (update *UpdateStage) serializeAudBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(update.Stage, &bytes)
	util.PutToken(update.Submission, &bytes)
	util.PutToken(update.Moderation, &bytes)
	util.PutToken(update.DiffHellKey, &bytes)
	util.PutByte(update.Flag, &bytes)
	util.PutString(update.Description, &bytes)
	putTokenCiphers(update.ReadMembers, &bytes)
	putTokenCiphers(update.SubMembers, &bytes)
	putTokenCiphers(update.ModMembers, &bytes)
	return bytes
}

func (update *UpdateStage) serializeBulk() []byte {
	bytes := update.serializeAudBulk()
	util.PutSignature(update.audSignature, &bytes)
	return bytes
}

func (update *UpdateStage) Serialize() []byte {
	return update.Authored.serialize(IUpdateAudience, update.serializeBulk())
}

func ParseUpdateStage(data []byte) *UpdateStage {
	if data[0] != 0 || data[1] != IUpdateAudience {
		return nil
	}
	update := UpdateStage{
		Authored: &AuthoredInstruction{},
	}
	position := update.Authored.parseHead(data)
	update.Stage, position = util.ParseToken(data, position)
	update.Submission, position = util.ParseToken(data, position)
	update.Moderation, position = util.ParseToken(data, position)
	update.DiffHellKey, position = util.ParseToken(data, position)
	update.Flag, position = util.ParseByte(data, position)
	update.Description, position = util.ParseString(data, position)
	update.ReadMembers, position = parseTokenCiphers(data, position)
	update.SubMembers, position = parseTokenCiphers(data, position)
	update.ModMembers, position = parseTokenCiphers(data, position)
	update.audSignature, position = util.ParseSignature(data, position)
	if update.Authored.parseTail(data, position) {
		return &update
	}
	return nil
}
