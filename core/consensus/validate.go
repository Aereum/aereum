package consensus

import (
	"github.com/Aereum/aereum/core/chain"
	"github.com/Aereum/aereum/core/instructions"
)

func ValidateBlock(data []byte, validator chain.MutatingState) *chain.Block {
	block := chain.ParseBlock(data)
	block.SetValidator(&validator)
	for _, instructionBytes := range block.Instructions {
		instruction := instructions.ParseInstruction(instructionBytes)
		if instruction == nil {
			return nil
		}
		if !block.Incorporate(instruction) {
			return nil
		}
	}
	return block
}
