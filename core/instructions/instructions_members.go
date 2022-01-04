package instructions

import (
	"encoding/json"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type JoinNetwork struct {
	Authored *AuthoredInstruction
	Caption  string
	Details  string
}

func (a *JoinNetwork) Epoch() uint64 {
	return a.Authored.epoch
}

func (join *JoinNetwork) Validate(v InstructionValidator) bool {
	captionHash := crypto.Hasher([]byte(join.Caption))
	if v.HasCaption(captionHash) {
		return false
	}
	authorHash := crypto.Hasher(join.Authored.Author[:])
	if v.HasMember(authorHash) {
		return false
	}
	if !json.Valid([]byte(join.Details)) {
		return false
	}
	if v.SetNewMember(authorHash, captionHash) {
		v.AddFeeCollected(join.Authored.Fee)
		return true
	}
	return false
}

func (join *JoinNetwork) Payments() *Payment {
	return join.Authored.payments()
}

func (join *JoinNetwork) Kind() byte {
	return IJoinNetwork
}

func (join *JoinNetwork) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutString(join.Caption, &bytes)
	util.PutString(join.Details, &bytes)
	return bytes
}

func (join *JoinNetwork) Serialize() []byte {
	return join.Authored.serialize(IJoinNetwork, join.serializeBulk())
}

func ParseJoinNetwork(data []byte) *JoinNetwork {
	if data[0] != 0 || data[1] != IJoinNetwork {
		return nil
	}
	join := JoinNetwork{
		Authored: &AuthoredInstruction{},
	}
	position := join.Authored.parseHead(data)
	join.Caption, position = util.ParseString(data, position)
	join.Details, position = util.ParseString(data, position)
	if !json.Valid([]byte(join.Details)) {
		return nil
	}
	if join.Authored.parseTail(data, position) {
		return &join
	}
	return nil
}

type UpdateInfo struct {
	Authored *AuthoredInstruction
	Details  string
}

func (a *UpdateInfo) Epoch() uint64 {
	return a.Authored.epoch
}

func (update *UpdateInfo) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.Authored.authorHash()) {
		return false
	}
	if json.Valid([]byte(update.Details)) {
		v.AddFeeCollected(update.Authored.Fee)
		return true
	}
	return false
}

func (update *UpdateInfo) Payments() *Payment {
	return update.Authored.payments()
}

func (update *UpdateInfo) Kind() byte {
	return IUpdateInfo
}

func (update *UpdateInfo) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutString(update.Details, &bytes)
	return bytes
}

func (update *UpdateInfo) Serialize() []byte {
	return update.Authored.serialize(IUpdateInfo, update.serializeBulk())
}

func ParseUpdateInfo(data []byte) *UpdateInfo {
	if data[0] != 0 || data[1] != IUpdateInfo {
		return nil
	}
	update := UpdateInfo{
		Authored: &AuthoredInstruction{},
	}
	position := update.Authored.parseHead(data)
	update.Details, position = util.ParseString(data, position)
	if !json.Valid([]byte(update.Details)) {
		return nil
	}
	if update.Authored.parseTail(data, position) {
		return &update
	}
	return nil
}

type GrantPowerOfAttorney struct {
	Authored *AuthoredInstruction
	Attorney crypto.Token
}

func (a *GrantPowerOfAttorney) Epoch() uint64 {
	return a.Authored.epoch
}

func (grant *GrantPowerOfAttorney) Validate(v InstructionValidator) bool {
	if !v.HasMember(grant.Authored.authorHash()) {
		return false
	}
	if !v.HasMember(crypto.HashToken(grant.Attorney)) {
		return false
	}
	hash := crypto.Hasher(append(grant.Authored.Author[:], grant.Attorney[:]...))
	if v.PowerOfAttorney(hash) {
		return false
	}
	if v.SetNewGrantPower(hash) {
		v.AddFeeCollected(grant.Authored.Fee)
		return true
	}
	return false
}

func (grant *GrantPowerOfAttorney) Payments() *Payment {
	return grant.Authored.payments()
}

func (grant *GrantPowerOfAttorney) Kind() byte {
	return IGrantPowerOfAttorney
}

func (grant *GrantPowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(grant.Attorney, &bytes)
	return bytes
}

func (grant *GrantPowerOfAttorney) Serialize() []byte {
	return grant.Authored.serialize(IGrantPowerOfAttorney, grant.serializeBulk())
}

func ParseGrantPowerOfAttorney(data []byte) *GrantPowerOfAttorney {
	if data[0] != 0 || data[1] != IGrantPowerOfAttorney {
		return nil
	}
	grant := GrantPowerOfAttorney{
		Authored: &AuthoredInstruction{},
	}
	position := grant.Authored.parseHead(data)
	grant.Attorney, position = util.ParseToken(data, position)
	if grant.Authored.parseTail(data, position) {
		return &grant
	}
	return nil
}

type RevokePowerOfAttorney struct {
	Authored *AuthoredInstruction
	Attorney crypto.Token
}

func (a *RevokePowerOfAttorney) Epoch() uint64 {
	return a.Authored.epoch
}

func (revoke *RevokePowerOfAttorney) Validate(v InstructionValidator) bool {
	if !v.HasMember(revoke.Authored.authorHash()) {
		return false
	}
	if !v.HasMember(crypto.HashToken(revoke.Attorney)) {
		return false
	}
	hash := crypto.Hasher(append(revoke.Authored.Author[:], revoke.Attorney[:]...))
	if !v.PowerOfAttorney(hash) {
		return false
	}
	if v.SetNewRevokePower(hash) {
		v.AddFeeCollected(revoke.Authored.Fee)
		return true
	}
	return false
}

func (revoke *RevokePowerOfAttorney) Payments() *Payment {
	return revoke.Authored.payments()
}

func (revoke *RevokePowerOfAttorney) Kind() byte {
	return IRevokePowerOfAttorney
}

func (revoke *RevokePowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(revoke.Attorney, &bytes)
	return bytes
}

func (revoke *RevokePowerOfAttorney) Serialize() []byte {
	return revoke.Authored.serialize(IRevokePowerOfAttorney, revoke.serializeBulk())
}

func ParseRevokePowerOfAttorney(data []byte) *RevokePowerOfAttorney {
	if data[0] != 0 || data[1] != IRevokePowerOfAttorney {
		return nil
	}
	revoke := RevokePowerOfAttorney{
		Authored: &AuthoredInstruction{},
	}
	position := revoke.Authored.parseHead(data)
	revoke.Attorney, position = util.ParseToken(data, position)
	if revoke.Authored.parseTail(data, position) {
		return &revoke
	}
	return nil
}

type CreateEphemeral struct {
	Authored       *AuthoredInstruction
	EphemeralToken crypto.Token
	Expiry         uint64
}

func (a *CreateEphemeral) Epoch() uint64 {
	return a.Authored.epoch
}

func (ephemeral *CreateEphemeral) Validate(v InstructionValidator) bool {
	if !v.HasMember(ephemeral.Authored.authorHash()) {
		return false
	}
	if ephemeral.Expiry <= v.Epoch() {
		return false
	}
	hash := crypto.HashToken(ephemeral.EphemeralToken)
	if ok, expire := v.GetEphemeralExpire(hash); ok && expire > v.Epoch() {
		return false
	}
	if v.SetNewEphemeralToken(hash, ephemeral.Expiry) {
		v.AddFeeCollected(ephemeral.Authored.Fee)
		return true
	}
	return false
}

func (ephemeral *CreateEphemeral) Payments() *Payment {
	return ephemeral.Authored.payments()
}

func (ephemeral *CreateEphemeral) Kind() byte {
	return ICreateEphemeral
}

func (ephemeral *CreateEphemeral) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(ephemeral.EphemeralToken, &bytes)
	util.PutUint64(ephemeral.Expiry, &bytes)
	return bytes
}

func (ephemeral *CreateEphemeral) Serialize() []byte {
	return ephemeral.Authored.serialize(ICreateEphemeral, ephemeral.serializeBulk())
}

func ParseCreateEphemeral(data []byte) *CreateEphemeral {
	if data[0] != 0 || data[1] != ICreateEphemeral {
		return nil
	}
	ephemeral := CreateEphemeral{
		Authored: &AuthoredInstruction{},
	}
	position := ephemeral.Authored.parseHead(data)
	ephemeral.EphemeralToken, position = util.ParseToken(data, position)
	ephemeral.Expiry, position = util.ParseUint64(data, position)
	if ephemeral.Authored.parseTail(data, position) {
		return &ephemeral
	}
	return nil
}

type SecureChannel struct {
	Authored       *AuthoredInstruction
	TokenRange     []byte
	Nonce          uint64
	EncryptedNonce []byte
	Content        []byte
}

func (a *SecureChannel) Epoch() uint64 {
	return a.Authored.epoch
}

func (secure *SecureChannel) Validate(v InstructionValidator) bool {
	authorHash := crypto.Hasher(secure.Authored.Author[:])
	if _, expire := v.GetEphemeralExpire(authorHash); expire <= v.Epoch() {
		return false
	}
	if len(secure.TokenRange) >= crypto.Size {
		return false
	}
	v.AddFeeCollected(secure.Authored.Fee)
	return true
}

func (secure *SecureChannel) Payments() *Payment {
	return secure.Authored.payments()
}

func (secure *SecureChannel) Kind() byte {
	return ISecureChannel
}

func (secure *SecureChannel) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(secure.TokenRange, &bytes)
	util.PutUint64(secure.Nonce, &bytes)
	util.PutByteArray(secure.EncryptedNonce, &bytes)
	util.PutByteArray(secure.Content, &bytes)
	return bytes
}

func (secure *SecureChannel) Serialize() []byte {
	return secure.Authored.serialize(ISecureChannel, secure.serializeBulk())
}

func ParseSecureChannel(data []byte) *SecureChannel {
	if data[0] != 0 || data[1] != ISecureChannel {
		return nil
	}
	secure := SecureChannel{
		Authored: &AuthoredInstruction{},
	}
	position := secure.Authored.parseHead(data)
	secure.TokenRange, position = util.ParseByteArray(data, position)
	secure.Nonce, position = util.ParseUint64(data, position)
	secure.EncryptedNonce, position = util.ParseByteArray(data, position)
	secure.Content, position = util.ParseByteArray(data, position)

	if secure.Authored.parseTail(data, position) {
		return &secure
	}
	return nil
}
