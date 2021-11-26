package instructionsnew

import (
	"github.com/Aereum/aereum/core/crypto"
)

// 	Content creation instruction
type SponsorshipOffer struct {
	authored    *authoredInstruction
	audience    []byte
	contentType string
	expiry      uint64
	revenue     uint64
}

func (sponsored *SponsorshipOffer) Validate(block *Block) bool {
	if !block.validator.HasMember(sponsored.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(sponsored.audience)
	keys := block.validator.GetAudienceKeys(audienceHash)
	if keys == nil {
		return false
	}
	if sponsored.expiry <= block.Epoch {
		return false
	}
	var balance uint64
	if sponsored.authored.wallet != nil {
		balance = block.validator.Balance(crypto.Hasher(sponsored.authored.wallet))
	} else if sponsored.authored.attorney != nil {
		balance = block.validator.Balance(crypto.Hasher(sponsored.authored.attorney))
	} else {
		balance = block.validator.Balance(crypto.Hasher(sponsored.authored.author))
	}
	if sponsored.revenue+sponsored.authored.fee > balance {
		return false
	}
	hash := crypto.Hasher(sponsored.Serialize())
	return block.SetNewSpnOffer(hash, sponsored.expiry)
}

func (sponsored *SponsorshipOffer) Payments() *Payment {
	return sponsored.authored.payments()
}

func (sponsored *SponsorshipOffer) Kind() byte {
	return iSponsorshipOffer
}

func (sponsored *SponsorshipOffer) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(sponsored.audience, &bytes)
	PutString(sponsored.contentType, &bytes)
	PutUint64(sponsored.expiry, &bytes)
	PutUint64(sponsored.revenue, &bytes)
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
	sponsored.audience, position = ParseByteArray(data, position)
	sponsored.contentType, position = ParseString(data, position)
	sponsored.expiry, position = ParseUint64(data, position)
	sponsored.revenue, position = ParseUint64(data, position)
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

func (accept *SponsorshipAcceptance) Validate(block *Block) bool {
	if !block.validator.HasMember(accept.authored.authorHash()) {
		return false
	}
	audienceHash := crypto.Hasher(accept.audience)
	keys := block.validator.GetAudienceKeys(audienceHash)
	if keys == nil {
		return false
	}
	if accept.offer.expiry >= block.Epoch {
		return false
	}
	offerHash := crypto.Hasher(accept.offer.Serialize())
	if block.validator.SponsorshipOffer(offerHash) != 0 {
		return false
	}
	hash := crypto.Hasher(accept.serializeModBulk())
	modKey, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	if !modKey.Verify(hash[:], accept.modSignature) {
		return false
	}
	return block.SetNewUseSonOffer(offerHash)

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
	PutByteArray(accept.audience, &bytes)
	PutByteArray(accept.offer.Serialize(), &bytes)
	return bytes
}

func (accept *SponsorshipAcceptance) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	PutByteArray(accept.modSignature, &bytes)
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
	accept.audience, position = ParseByteArray(data, position)
	var offerBytes []byte
	offerBytes, position = ParseByteArray(data, position)
	accept.offer = ParseSponsorshipOffer(offerBytes)
	if accept.offer == nil {
		return nil
	}
	accept.modSignature, position = ParseByteArray(data, position)
	if accept.authored.parseTail(data, position) {
		return &accept
	}
	return nil
}
