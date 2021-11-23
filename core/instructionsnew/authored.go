package instructionsnew

type BulkSerializer interface {
	serializeBulk() []byte
	InstructionType() byte
}

type authoredInstruction struct {
	epoch           uint64
	author          []byte
	wallet          []byte
	fee             uint64
	attorney        []byte
	signature       []byte
	walletSignature []byte
}

func (a *authoredInstruction) serializeWithoutSignature(instType byte, bulk []byte) []byte {
	bytes := []byte{0, instType}
	PutUint64(a.epoch, &bytes)
	PutByteArray(a.author, &bytes)
	bytes = append(bytes, bulk...)
	PutByteArray(a.wallet, &bytes)
	PutUint64(a.fee, &bytes)
	PutByteArray(a.attorney, &bytes)
	return bytes
}

func (a *authoredInstruction) serialize(instType byte, bulk []byte) []byte {
	bytes := a.serializeWithoutSignature(instType, bulk)
	PutByteArray(a.signature, &bytes)
	return bytes
}

func (a *authoredInstruction) parseHead(data []byte) int {
	position := 2
	a.epoch, position = ParseUint64(data, position)
	a.author, position = ParseByteArray(data, position)
	return position
}

func (a *authoredInstruction) parseTail(data []byte, position int) bool {
	a.wallet, position = ParseByteArray(data, position)
	a.fee, position = ParseUint64(data, position)
	a.attorney, position = ParseByteArray(data, position)
	a.signature, position = ParseByteArray(data, position)
	a.walletSignature, position = ParseByteArray(data, position)
	if position != len(data) {
		return false
	}
	return true
}
