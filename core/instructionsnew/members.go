package instructionsnew

import (
	"encoding/json"

	"github.com/Aereum/aereum/core/crypto"
)

type JoinNetwork struct {
	authored *authoredInstruction
	caption  string
	details  string
}

func (join *JoinNetwork) Validate(block *Block) bool {
	captionHash := crypto.Hasher([]byte(join.caption))
	if block.validator.HasCaption(captionHash) {
		return false
	}
	authorHash := crypto.Hasher([]byte(join.authored.author))
	if block.validator.HasMember(authorHash) {
		return false
	}
	if !json.Valid([]byte(join.details)) {
		return false
	}
	return block.SetNewMember(authorHash, captionHash)
}

func (join *JoinNetwork) Payments() *Payment {
	return join.authored.payments()
}

func (join *JoinNetwork) Kind() byte {
	return iJoinNetwork
}

func (join *JoinNetwork) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutString(join.caption, &bytes)
	PutString(join.details, &bytes)
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
	join.caption, position = ParseString(data, position)
	join.details, position = ParseString(data, position)
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

func (update *UpdateInfo) Validate(block *Block) bool {
	if !block.validator.HasMember(update.authored.authorHash()) {
		return false
	}
	return !json.Valid([]byte(update.details))
}

func (update *UpdateInfo) Payments() *Payment {
	return update.authored.payments()
}

func (update *UpdateInfo) Kind() byte {
	return iUpdateInfo
}

func (update *UpdateInfo) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutString(update.details, &bytes)
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
	update.details, position = ParseString(data, position)
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

func (grant *GrantPowerOfAttorney) Validate(block *Block) bool {
	if !block.validator.HasMember(grant.authored.authorHash()) {
		return false
	}
	if !block.validator.HasMember(crypto.Hasher(grant.attorney)) {
		return false
	}
	hash := crypto.Hasher(append(grant.authored.author, grant.attorney...))
	if block.validator.PowerOfAttorney(hash) {
		return false
	}
	return block.SetNewGrantPower(hash)
}

func (grant *GrantPowerOfAttorney) Payments() *Payment {
	return grant.authored.payments()
}

func (grant *GrantPowerOfAttorney) Kind() byte {
	return iGrantPowerOfAttorney
}

func (grant *GrantPowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(grant.attorney, &bytes)
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
	grant.attorney, position = ParseByteArray(data, position)
	if grant.authored.parseTail(data, position) {
		return &grant
	}
	return nil
}

type RevokePowerOfAttorney struct {
	authored *authoredInstruction
	attorney []byte
}

func (revoke *RevokePowerOfAttorney) Validate(block *Block) bool {
	if !block.validator.HasMember(revoke.authored.authorHash()) {
		return false
	}
	if !block.validator.HasMember(crypto.Hasher(revoke.attorney)) {
		return false
	}
	hash := crypto.Hasher(append(revoke.authored.author, revoke.attorney...))
	if !block.validator.PowerOfAttorney(hash) {
		return false
	}
	return block.SetNewRevokePower(hash)
}

func (revoke *RevokePowerOfAttorney) Payments() *Payment {
	return revoke.authored.payments()
}

func (revoke *RevokePowerOfAttorney) Kind() byte {
	return iRevokePowerOfAttorney
}

func (revoke *RevokePowerOfAttorney) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(revoke.attorney, &bytes)
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
	revoke.attorney, position = ParseByteArray(data, position)
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

func (ephemeral *CreateEphemeral) Payments() *Payment {
	return ephemeral.authored.payments()
}

func (ephemeral *CreateEphemeral) Kind() byte {
	return iCreateEphemeral
}

func (ephemeral *CreateEphemeral) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(ephemeral.ephemeralToken, &bytes)
	PutUint64(ephemeral.expiry, &bytes)
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
	ephemeral.ephemeralToken, position = ParseByteArray(data, position)
	ephemeral.expiry, position = ParseUint64(data, position)
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

func (secure *SecureChannel) Payments() *Payment {
	return secure.authored.payments()
}

func (secure *SecureChannel) Kind() byte {
	return iSecureChannel
}

func (secure *SecureChannel) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(secure.tokenRange, &bytes)
	PutUint64(secure.nonce, &bytes)
	PutByteArray(secure.encryptedNonce, &bytes)
	PutByteArray(secure.content, &bytes)
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
	secure.tokenRange, position = ParseByteArray(data, position)
	secure.nonce, position = ParseUint64(data, position)
	secure.encryptedNonce, position = ParseByteArray(data, position)
	secure.content, position = ParseByteArray(data, position)

	if secure.authored.parseTail(data, position) {
		return &secure
	}
	return nil
}
