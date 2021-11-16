package consensus

import (
	"github.com/Aereum/aereum/core/crypto"

	"github.com/Aereum/aereum/core/message"
)

type BlockSignature struct {
	Epoch     uint64
	Token     crypto.PublicKey
	Signature []byte
}

type Consensus interface {
	AppendSignature(signature BlockSignature) bool
	NewInstruction(m *message.Message)
}
