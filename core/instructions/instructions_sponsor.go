package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

// 	Content creation instruction
type SponsorshipOffer struct {
	Authored    *AuthoredInstruction
	Stage       crypto.Token
	ContentType string
	Content     []byte
	Expiry      uint64
	Revenue     uint64
}

func (a *SponsorshipOffer) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *SponsorshipOffer) Epoch() uint64 {
	return a.Authored.epoch
}

func (sponsored *SponsorshipOffer) Validate(v InstructionValidator) bool {
	if !v.HasMember(sponsored.Authored.authorHash()) {
		return false
	}
	stageHash := crypto.HashToken(sponsored.Stage)
	stageKeys := v.GetAudienceKeys(stageHash)
	if stageKeys == nil {
		return false
	}
	if sponsored.Expiry <= v.Epoch() {
		return false
	}
	var balance uint64
	if sponsored.Authored.Wallet != crypto.ZeroToken {
		balance = v.Balance(crypto.HashToken(sponsored.Authored.Wallet))
	} else if sponsored.Authored.Attorney != crypto.ZeroToken {
		balance = v.Balance(crypto.HashToken(sponsored.Authored.Attorney))
	} else {
		balance = v.Balance(crypto.HashToken(sponsored.Authored.Author))
	}
	if sponsored.Revenue+sponsored.Authored.Fee > balance {
		return false
	}
	hash := crypto.Hasher(sponsored.Serialize())
	if v.SetNewSpnOffer(hash, sponsored.Expiry) {
		v.AddFeeCollected(sponsored.Authored.Fee)
		return true
	}
	return false
}

func (sponsored *SponsorshipOffer) Payments() *Payment {
	return sponsored.Authored.payments()
}

func (sponsored *SponsorshipOffer) Kind() byte {
	return ISponsorshipOffer
}

func (sponsored *SponsorshipOffer) serializeBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(sponsored.Stage, &bytes)
	util.PutString(sponsored.ContentType, &bytes)
	util.PutByteArray(sponsored.Content, &bytes)
	util.PutUint64(sponsored.Expiry, &bytes)
	util.PutUint64(sponsored.Revenue, &bytes)
	return bytes
}

func (sponsored *SponsorshipOffer) Serialize() []byte {
	return sponsored.Authored.serialize(ISponsorshipOffer, sponsored.serializeBulk())
}

func ParseSponsorshipOffer(data []byte) *SponsorshipOffer {
	if data[0] != 0 || data[1] != ISponsorshipOffer {
		return nil
	}
	sponsored := SponsorshipOffer{
		Authored: &AuthoredInstruction{},
	}
	position := sponsored.Authored.parseHead(data)
	sponsored.Stage, position = util.ParseToken(data, position)
	sponsored.ContentType, position = util.ParseString(data, position)
	sponsored.Content, position = util.ParseByteArray(data, position)
	sponsored.Expiry, position = util.ParseUint64(data, position)
	sponsored.Revenue, position = util.ParseUint64(data, position)
	if sponsored.Authored.parseTail(data, position) {
		return &sponsored
	}
	return nil
}

// Reaction instruction
type SponsorshipAcceptance struct {
	Authored     *AuthoredInstruction
	Stage        crypto.Token
	Offer        *SponsorshipOffer
	modSignature crypto.Signature
}

func (a *SponsorshipAcceptance) Authority() crypto.Token {
	return a.Authored.Author
}

func (a *SponsorshipAcceptance) Epoch() uint64 {
	return a.Authored.epoch
}

func (accept *SponsorshipAcceptance) Validate(v InstructionValidator) bool {

	// NAO ENCONTREI A FUNCAO E NAO CONSEGUI MONTAR A FUNCAO
	// block.validator.setNewHash(accept.offer.content)
	if !v.HasMember(accept.Authored.authorHash()) {
		return false
	}
	stageHash := crypto.HashToken(accept.Stage)
	stageKeys := v.GetAudienceKeys(stageHash)
	if stageKeys == nil {
		return false
	}
	if accept.Offer.Expiry < v.Epoch() {
		return false
	}
	offerHash := crypto.Hasher(accept.Offer.Serialize())
	if v.SponsorshipOffer(offerHash) == 0 {
		return false
	}
	//hash := crypto.Hasher(accept.serializeModBulk())
	if !stageKeys.Moderate.Verify(accept.serializeModBulk(), accept.modSignature) {
		return false
	}
	if v.SetNewUseSpnOffer(offerHash) {
		v.AddFeeCollected(accept.Authored.Fee)
		return true
	}
	return false
}

func (accept *SponsorshipAcceptance) Payments() *Payment {
	payments := accept.Authored.payments()
	payments.NewCredit(crypto.HashToken(accept.Stage), accept.Offer.Revenue)
	if accept.Offer.Authored.Wallet != crypto.ZeroToken {
		payments.NewDebit(crypto.HashToken(accept.Offer.Authored.Wallet), accept.Offer.Revenue)
	} else if accept.Offer.Authored.Attorney != crypto.ZeroToken {
		payments.NewDebit(crypto.HashToken(accept.Offer.Authored.Attorney), accept.Offer.Revenue)
	} else {
		payments.NewDebit(crypto.HashToken(accept.Offer.Authored.Author), accept.Offer.Revenue)
	}
	return payments
}

func (accept *SponsorshipAcceptance) Kind() byte {
	return ISponsorshipAcceptance
}

func (accept *SponsorshipAcceptance) serializeModBulk() []byte {
	bytes := make([]byte, 0)
	util.PutToken(accept.Stage, &bytes)
	util.PutByteArray(accept.Offer.Serialize(), &bytes)
	return bytes
}

func (accept *SponsorshipAcceptance) serializeBulk() []byte {
	bytes := accept.serializeModBulk()
	util.PutSignature(accept.modSignature, &bytes)
	return bytes
}

func (accept *SponsorshipAcceptance) Serialize() []byte {
	return accept.Authored.serialize(ISponsorshipAcceptance, accept.serializeBulk())
}

func ParseSponsorshipAcceptance(data []byte) *SponsorshipAcceptance {
	if data[0] != 0 || data[1] != ISponsorshipAcceptance {
		return nil
	}
	accept := SponsorshipAcceptance{
		Authored: &AuthoredInstruction{},
	}
	position := accept.Authored.parseHead(data)
	accept.Stage, position = util.ParseToken(data, position)
	var offerBytes []byte
	offerBytes, position = util.ParseByteArray(data, position)
	accept.Offer = ParseSponsorshipOffer(offerBytes)
	if accept.Offer == nil {
		return nil
	}
	accept.modSignature, position = util.ParseSignature(data, position)
	if accept.Authored.parseTail(data, position) {
		return &accept
	}
	return nil
}
