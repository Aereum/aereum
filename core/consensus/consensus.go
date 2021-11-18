package consensus

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
	"github.com/Aereum/aereum/core/network"
	"github.com/Aereum/aereum/core/state"
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

type ConsensusEngine func(initial state.State, peers network.ValidatorNetwork) *Consensus

type Consensus struct {
	Signature   chan *BlockSignature
	Instruction chan *instruction.Instruction
}

func (c *Consensus) NewInstruction(i *instruction.Instruction) {
	c.Instruction <- i
}

func (c *Consensus) BlockSignature(sign *BlockSignature) {
	c.Signature <- sign
}
