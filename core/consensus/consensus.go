package consensus

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type BlockSignature struct {
	Epoch     uint64
	Token     crypto.PublicKey
	Signature []byte
}

type HashInstruction struct {
	Instruction *instruction.Instruction
	Hash        crypto.Hash
}

type Consensus interface {
	AppendSignature(signature BlockSignature) bool
	NewInstruction(m *instruction.Instruction)
}
