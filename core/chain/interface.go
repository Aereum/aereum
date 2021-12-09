package chain

import "github.com/Aereum/aereum/core/crypto"

type InstructionValidator interface {
	SetNewGrantPower(hash crypto.Hash) bool
	SetNewRevokePower(hash crypto.Hash) bool
	SetNewUseSpnOffer(hash crypto.Hash) bool
	SetNewSpnOffer(hash crypto.Hash, expire uint64) bool
	SetPublishSponsor(hash crypto.Hash) bool
	SetNewEphemeralToken(hash crypto.Hash, expire uint64) bool
	SetNewMember(tokenHash crypto.Hash, captionHashe crypto.Hash) bool
	SetNewAudience(hash crypto.Hash, keys []byte) bool
	UpdateAudience(hash crypto.Hash, keys []byte) bool
	GetBalance(hash crypto.Hash) uint64
	PowerOfAttorney(hash crypto.Hash) bool
	SponsorshipOffer(hash crypto.Hash) uint64
	HasMember(hash crypto.Hash) bool
	HasCaption(hash crypto.Hash) bool
	HasGrantedSponser(hash crypto.Hash) (bool, crypto.Hash)
	GetAudienceKeys(hash crypto.Hash) []byte
	GetEphemeralExpire(hash crypto.Hash) (bool, uint64)
	Epoch() uint64
}
