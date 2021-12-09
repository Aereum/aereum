package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

// 	Content creation instruction
type SponsorshipOffer struct {
	authored    *authoredInstruction
	audience    []byte
	contentType string
	content     []byte
	expiry      uint64
	revenue     uint64
}

func (a *SponsorshipOffer) Epoch() uint64 {
	return a.authored.epoch
}

func (sponsored *SponsorshipOffer) Validate(v InstructionValidator) bool {
	if !v.HasMember(sponsored.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(sponsored.audience)
	keys := v.GetAudienceKeys(audienceHash)
	if keys == nil {
		return false
	}
	if sponsored.expiry <= v.Epoch() {
		return false
	}
	var balance uint64
	if sponsored.authored.wallet != nil {
		balance = v.Balance(crypto.Hasher(sponsored.authored.wallet))
	} else if sponsored.authored.attorney != nil {
		balance = v.Balance(crypto.Hasher(sponsored.authored.attorney))
	} else {
		balance = v.Balance(crypto.Hasher(sponsored.authored.author))
	}
	if sponsored.revenue+sponsored.authored.fee > balance {
		return false
	}
	hash := crypto.Hasher(sponsored.Serialize())
	if v.SetNewSpnOffer(hash, sponsored.expiry) {
		v.AddFeeCollected(sponsored.authored.fee)
		return true
	}
	return false
}

func (sponsored *SponsorshipOffer) Payments() *Payment {
	return sponsored.authored.payments()
}

func (sponsored *SponsorshipOffer) Kind() byte {
	return iSponsorshipOffer
}

func (sponsored *SponsorshipOffer) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(sponsored.audience, &bytes)
	util.PutString(sponsored.contentType, &bytes)
	util.PutByteArray(sponsored.content, &bytes)
	util.PutUint64(sponsored.expiry, &bytes)
	util.PutUint64(sponsored.revenue, &bytes)
	return bytes
}

func (sponsored *SponsorshipOffer) Serialize() []byte {
	return sponsored.authored.serialize(iSponsorshipOffer, sponsored.serializeBulk())
}

func ParseSponsorshipOffer(data []byte) *SponsorshipOffer {
	if data[0] != 0 || data[1] != iSponsorshipOffer {
		return nil
	}
	sponsored := SponsorshipOffer{
		authored: &authoredInstruction{},
	}
	position := sponsored.authored.parseHead(data)
	sponsored.audience, position = util.ParseByteArray(data, position)
	sponsored.contentType, position = util.ParseString(data, position)
	sponsored.content, position = util.ParseByteArray(data, position)
	sponsored.expiry, position = util.ParseUint64(data, position)
	sponsored.revenue, position = util.ParseUint64(data, position)
	if sponsored.authored.parseTail(data, position) {
		return &sponsored
	}
	return nil
}

// Reaction instruction
type SponsorshipAcceptance struct {
	authored     *authoredInstruction
	audience     []byte
	offer        *SponsorshipOffer
	modSignature []byte
}

func (a *SponsorshipAcceptance) Epoch() uint64 {
	return a.authored.epoch
}

func (accept *SponsorshipAcceptance) Validate(v InstructionValidator) bool {

	// NAO ENCONTREI A FUNCAO E NAO CONSEGUI MONTAR A FUNCAO
	// block.validator.setNewHash(accept.offer.content)
	if !v.HasMember(accept.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(accept.audience)
	keys := v.GetAudienceKeys(audienceHash)
	if keys == nil {
		return false
	}
	if accept.offer.expiry < v.Epoch() {
		return false
	}
	offerHash := crypto.Hasher(accept.offer.Serialize())
	if v.SponsorshipOffer(offerHash) == 0 {
		return false
	}
	//hash := crypto.Hasher(accept.serializeModBulk())
	modKey, err := crypto.PublicKeyFromBytes(keys[crypto.PublicKeySize : 2*crypto.PublicKeySize])
	if err != nil {
		return false
	}
	if !modKey.Verify(accept.serializeModBulk(), accept.modSignature) {
		return false
	}
	if v.SetNewUseSpnOffer(offerHash) {
		v.AddFeeCollected(accept.authored.fee)
		return true
	}
	return false
}

func (accept *SponsorshipAcceptance) Payments() *Payment {
	payments := accept.authored.payments()
	payments.NewCredit(crypto.Hasher(accept.audience), accept.offer.revenue)
	if accept.offer.authored.wallet != nil {
		payments.NewDebit(crypto.Hasher(accept.offer.authored.wallet), accept.offer.revenue)
	} else if accept.offer.authored.attorney != nil {
		payments.NewDebit(crypto.Hasher(accept.offer.authored.attorney), accept.offer.revenue)
	} else {
		payments.NewDebit(crypto.Hasher(accept.offer.authored.author), accept.offer.revenue)
	}
	return payments
}

func (accept *SponsorshipAcceptance) Kind() byte {
	return iSponsorshipAcceptance
}

func (accept *SponsorshipAcceptance) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	util.PutByteArray(accept.audience, &bytes)
	util.PutByteArray(accept.offer.Serialize(), &bytes)
	return bytes
}

func (accept *SponsorshipAcceptance) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutByteArray(accept.modSignature, &bytes)
	return bytes
}

func (accept *SponsorshipAcceptance) Serialize() []byte {
	return accept.authored.serialize(iSponsorshipAcceptance, accept.serializeBulk())
}

func ParseSponsorshipAcceptance(data []byte) *SponsorshipAcceptance {
	if data[0] != 0 || data[1] != iSponsorshipAcceptance {
		return nil
	}
	accept := SponsorshipAcceptance{
		authored: &authoredInstruction{},
	}
	position := accept.authored.parseHead(data)
	accept.audience, position = util.ParseByteArray(data, position)
	var offerBytes []byte
	offerBytes, position = util.ParseByteArray(data, position)
	accept.offer = ParseSponsorshipOffer(offerBytes)
	if accept.offer == nil {
		return nil
	}
	accept.modSignature, position = util.ParseByteArray(data, position)
	if accept.authored.parseTail(data, position) {
		return &accept
	}
	return nil
}
