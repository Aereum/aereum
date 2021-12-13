package instructions

import (
	"encoding/json"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

type JoinNetwork struct {
	authored *authoredInstruction
	caption  string
	details  string
}

func (a *JoinNetwork) Epoch() uint64 {
	return a.authored.epoch
}

func (join *JoinNetwork) Validate(v InstructionValidator) bool {
	captionHash := crypto.Hasher([]byte(join.caption))
	if v.HasCaption(captionHash) {
		return false
	}
	authorHash := crypto.Hasher(join.authored.author[:])
	if v.HasMember(authorHash) {
		return false
	}
	if !json.Valid([]byte(join.details)) {
		return false
	}
	if v.SetNewMember(authorHash, captionHash) {
		v.AddFeeCollected(join.authored.fee)
		return true
	}
	return false
}

func (join *JoinNetwork) Payments() *Payment {
	return join.authored.payments()
}

func (join *JoinNetwork) Kind() byte {
	return iJoinNetwork
}

func (join *JoinNetwork) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutString(join.caption, &bytes)
	util.PutString(join.details, &bytes)
	return bytes
}

func (join *JoinNetwork) Serialize() []byte {
	return join.authored.serialize(iJoinNetwork, join.serializeBulk())
}

func ParseJoinNetwork(data []byte) *JoinNetwork {
	if data[0] != 0 || data[1] != iJoinNetwork {
		return nil
	}
	join := JoinNetwork{
		authored: &authoredInstruction{},
	}
	position := join.authored.parseHead(data)
	join.caption, position = util.ParseString(data, position)
	join.details, position = util.ParseString(data, position)
	if !json.Valid([]byte(join.details)) {
		return nil
	}
	if join.authored.parseTail(data, position) {
		return &join
	}
	return nil
}

type UpdateInfo struct {
	authored *authoredInstruction
	details  string
}

func (a *UpdateInfo) Epoch() uint64 {
	return a.authored.epoch
}

func (update *UpdateInfo) Validate(v InstructionValidator) bool {
	if !v.HasMember(update.authored.authorHash()) {
		return false
	}
	if json.Valid([]byte(update.details)) {
		v.AddFeeCollected(update.authored.fee)
		return true
	}
	return false
}

func (update *UpdateInfo) Payments() *Payment {
	return update.authored.payments()
}

func (update *UpdateInfo) Kind() byte {
	return iUpdateInfo
}

func (update *UpdateInfo) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutString(update.details, &bytes)
	return bytes
}

func (update *UpdateInfo) Serialize() []byte {
	return update.authored.serialize(iUpdateInfo, update.serializeBulk())
}

func ParseUpdateInfo(data []byte) *UpdateInfo {
	if data[0] != 0 || data[1] != iUpdateInfo {
		return nil
	}
	update := UpdateInfo{
		authored: &authoredInstruction{},
	}
	position := update.authored.parseHead(data)
	update.details, position = util.ParseString(data, position)
	if !json.Valid([]byte(update.details)) {
		return nil
	}
	if update.authored.parseTail(data, position) {
		return &update
	}
	return nil
}

type GrantPowerOfAttorney struct {
	authored *authoredInstruction
	attorney []byte
}

func (a *GrantPowerOfAttorney) Epoch() uint64 {
	return a.authored.epoch
}

func (grant *GrantPowerOfAttorney) Validate(v InstructionValidator) bool {
	if !v.HasMember(grant.authored.authorHash()) {
		return false
	}
	if !v.HasMember(crypto.Hasher(grant.attorney)) {
		return false
	}
	hash := crypto.Hasher(append(grant.authored.author[:], grant.attorney...))
	if v.PowerOfAttorney(hash) {
		return false
	}
	if v.SetNewGrantPower(hash) {
		v.AddFeeCollected(grant.authored.fee)
		return true
	}
	return false
}

func (grant *GrantPowerOfAttorney) Payments() *Payment {
	return grant.authored.payments()
}

func (grant *GrantPowerOfAttorney) Kind() byte {
	return iGrantPowerOfAttorney
}

func (grant *GrantPowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(grant.attorney, &bytes)
	return bytes
}

func (grant *GrantPowerOfAttorney) Serialize() []byte {
	return grant.authored.serialize(iGrantPowerOfAttorney, grant.serializeBulk())
}

func ParseGrantPowerOfAttorney(data []byte) *GrantPowerOfAttorney {
	if data[0] != 0 || data[1] != iGrantPowerOfAttorney {
		return nil
	}
	grant := GrantPowerOfAttorney{
		authored: &authoredInstruction{},
	}
	position := grant.authored.parseHead(data)
	grant.attorney, position = util.ParseByteArray(data, position)
	if grant.authored.parseTail(data, position) {
		return &grant
	}
	return nil
}

type RevokePowerOfAttorney struct {
	authored *authoredInstruction
	attorney []byte
}

func (a *RevokePowerOfAttorney) Epoch() uint64 {
	return a.authored.epoch
}

func (revoke *RevokePowerOfAttorney) Validate(v InstructionValidator) bool {
	if !v.HasMember(revoke.authored.authorHash()) {
		return false
	}
	if !v.HasMember(crypto.Hasher(revoke.attorney)) {
		return false
	}
	hash := crypto.Hasher(append(revoke.authored.author[:], revoke.attorney...))
	if !v.PowerOfAttorney(hash) {
		return false
	}
	if v.SetNewRevokePower(hash) {
		v.AddFeeCollected(revoke.authored.fee)
		return true
	}
	return false
}

func (revoke *RevokePowerOfAttorney) Payments() *Payment {
	return revoke.authored.payments()
}

func (revoke *RevokePowerOfAttorney) Kind() byte {
	return iRevokePowerOfAttorney
}

func (revoke *RevokePowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(revoke.attorney, &bytes)
	return bytes
}

func (revoke *RevokePowerOfAttorney) Serialize() []byte {
	return revoke.authored.serialize(iRevokePowerOfAttorney, revoke.serializeBulk())
}

func ParseRevokePowerOfAttorney(data []byte) *RevokePowerOfAttorney {
	if data[0] != 0 || data[1] != iRevokePowerOfAttorney {
		return nil
	}
	revoke := RevokePowerOfAttorney{
		authored: &authoredInstruction{},
	}
	position := revoke.authored.parseHead(data)
	revoke.attorney, position = util.ParseByteArray(data, position)
	if revoke.authored.parseTail(data, position) {
		return &revoke
	}
	return nil
}

type CreateEphemeral struct {
	authored       *authoredInstruction
	ephemeralToken []byte
	expiry         uint64
}

func (a *CreateEphemeral) Epoch() uint64 {
	return a.authored.epoch
}

func (ephemeral *CreateEphemeral) Validate(v InstructionValidator) bool {
	if !v.HasMember(ephemeral.authored.authorHash()) {
		return false
	}
	if ephemeral.expiry <= v.Epoch() {
		return false
	}
	hash := crypto.Hasher(ephemeral.ephemeralToken)
	if ok, expire := v.GetEphemeralExpire(hash); ok && expire > v.Epoch() {
		return false
	}
	if v.SetNewEphemeralToken(hash, ephemeral.expiry) {
		v.AddFeeCollected(ephemeral.authored.fee)
		return true
	}
	return false
}

func (ephemeral *CreateEphemeral) Payments() *Payment {
	return ephemeral.authored.payments()
}

func (ephemeral *CreateEphemeral) Kind() byte {
	return iCreateEphemeral
}

func (ephemeral *CreateEphemeral) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(ephemeral.ephemeralToken, &bytes)
	util.PutUint64(ephemeral.expiry, &bytes)
	return bytes
}

func (ephemeral *CreateEphemeral) Serialize() []byte {
	return ephemeral.authored.serialize(iCreateEphemeral, ephemeral.serializeBulk())
}

func ParseCreateEphemeral(data []byte) *CreateEphemeral {
	if data[0] != 0 || data[1] != iCreateEphemeral {
		return nil
	}
	ephemeral := CreateEphemeral{
		authored: &authoredInstruction{},
	}
	position := ephemeral.authored.parseHead(data)
	ephemeral.ephemeralToken, position = util.ParseByteArray(data, position)
	ephemeral.expiry, position = util.ParseUint64(data, position)
	if ephemeral.authored.parseTail(data, position) {
		return &ephemeral
	}
	return nil
}

type SecureChannel struct {
	authored       *authoredInstruction
	tokenRange     []byte
	nonce          uint64
	encryptedNonce []byte
	content        []byte
}

func (a *SecureChannel) Epoch() uint64 {
	return a.authored.epoch
}

func (secure *SecureChannel) Validate(v InstructionValidator) bool {
	authorHash := crypto.Hasher(secure.authored.author[:])
	if _, expire := v.GetEphemeralExpire(authorHash); expire <= v.Epoch() {
		return false
	}
	if len(secure.tokenRange) >= crypto.Size {
		return false
	}
	v.AddFeeCollected(secure.authored.fee)
	return true
}

func (secure *SecureChannel) Payments() *Payment {
	return secure.authored.payments()
}

func (secure *SecureChannel) Kind() byte {
	return iSecureChannel
}

func (secure *SecureChannel) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(secure.tokenRange, &bytes)
	util.PutUint64(secure.nonce, &bytes)
	util.PutByteArray(secure.encryptedNonce, &bytes)
	util.PutByteArray(secure.content, &bytes)
	return bytes
}

func (secure *SecureChannel) Serialize() []byte {
	return secure.authored.serialize(iSecureChannel, secure.serializeBulk())
}

func ParseSecureChannel(data []byte) *SecureChannel {
	if data[0] != 0 || data[1] != iSecureChannel {
		return nil
	}
	secure := SecureChannel{
		authored: &authoredInstruction{},
	}
	position := secure.authored.parseHead(data)
	secure.tokenRange, position = util.ParseByteArray(data, position)
	secure.nonce, position = util.ParseUint64(data, position)
	secure.encryptedNonce, position = util.ParseByteArray(data, position)
	secure.content, position = util.ParseByteArray(data, position)

	if secure.authored.parseTail(data, position) {
		return &secure
	}
	return nil
}
