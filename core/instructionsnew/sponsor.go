package instructionsnew

// 	Content creation instruction
type SponsorshipOffer struct {
	authored    *authoredInstruction
	audience    []byte
	contentType string
	expiry      uint64
	revenue     uint64
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
	offer        *authoredInstruction
	modSignature []byte
}

func (sponsoracceptance *SponsorshipAcceptance) Payments() *Payment {
	return sponsoracceptance.authored.payments()
}

func (sponsoracceptance *SponsorshipAcceptance) Kind() byte {
	return iSponsorshipAcceptance
}

func (sponsoracceptance *SponsorshipAcceptance) serializeBulk() []byte {
	bytes := make([]byte, 0)
	PutByteArray(sponsoracceptance.audience, &bytes)
	// aqui nao sei como resolver
	PutByteArray(sponsoracceptance.modSignature[], &bytes)
	return bytes
}

func (sponsoracceptance *SponsorshipAcceptance) Serialize() []byte {
	return sponsoracceptance.authored.serialize(iSponsorshipAcceptance, react.serializeBulk())
}

func ParseSponsorshipAcceptance(data []byte) *SponsorshipAcceptance {
	if data[0] != 0 || data[1] != iSponsorshipAcceptance {
		return nil
	}
	sponsoracceptance := SponsorshipAcceptance{
		authored: &authoredInstruction{},
	}
	position := sponsoracceptance.authored.parseHead(data)
	sponsoracceptance.audience, position = ParseByteArray(data, position)
	// precisa ler segundo authored
	sponsoracceptance.modSignature, position = ParseByteArray(data, position)
	if sponsoracceptance.authored.parseTail(data, position) {
		return &sponsoracceptance
	}
	return nil
}
