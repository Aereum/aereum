package main

import (
	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type StageSecret struct {
	stage  crypto.Hash
	secret [crypto.CipherSize]byte
}

var blocks = make([]consensus.SignedBlock, 0)

type instructionIndex struct {
	block    int
	position int
}

type Account struct {
	token         crypto.Hash
	caption       string
	wallet        uint64
	attoneys      []crypto.Hash
	stageOwner    []StageSecret
	stagePublish  []StageSecret
	stageModerate []StageSecret
	instructions  []instructionIndex
}

func (a *Account) filterInstructionByType(instructionType byte) []instructions.Instruction {
	output := make([]instructions.Instruction, 0)
	for _, index := range a.instructions {
		if newInst := getInstruction(index, instructionType); newInst != nil {
			output = append(output, newInst)
		}
	}
	return output
}

func getInstruction(index instructionIndex, instType byte) instructions.Instruction {
	if index.block > len(blocks) {
		return nil
	}
	blockInstructions := blocks[index.block].Block.Instructions
	if index.position > len(blockInstructions) {
		return nil
	}
	return instructions.ParseInstruction(blockInstructions[index.position])
}

type Audiences map[crypto.Hash]string
