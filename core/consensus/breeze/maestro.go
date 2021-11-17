package breeze

import (
	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/network"
)

const ChecksumWindows = 50000
const ValidatorsCount = 10

type Maestro struct {
	Peers          network.ValidatorNetwork
	Slots          []network.SecureConnection
	LiveCheckPoint uint64
	Proposed       map[uint64]*CheckPoint
	Instruction    chan *consensus.HashInstruction
	Signature      chan *consensus.BlockSignature
	IsLeader       bool
	Pool           *InstructionPool
}

func (maestro *Maestro) AppendSignature(signature *consensus.BlockSignature) {
	maestro.Signature <- signature
}

func (maestro *Maestro) appendSignature(signature *consensus.BlockSignature) {
	checkpoint, ok := maestro.Proposed[signature.Epoch]
	if !ok {
		return
	}
	ok, N := checkpoint.appendSignature(signature.Token, signature.Signature)
	if !ok {
		return
	}
	if N > ValidatorsCount/2 {
		if checkpoint.Block.Epoch == maestro.LiveCheckPoint+1 {
			maestro.LiveCheckPoint += 1
		}
	}
}

func (maestro *Maestro) NewInstruction(m *consensus.HashInstruction) {
	maestro.Instruction <- m
}

func (maestro *Maestro) newInstruction(m *consensus.HashInstruction) {
	if !maestro.IsLeader {
		maestro.Pool.Queue(m.Instruction, m.Hash)
	}
}

func NewMaestro() *Maestro {
	maestro := Maestro{
		Instruction: make(chan *consensus.HashInstruction),
		Signature:   make(chan *consensus.BlockSignature),
	}
	go func() {
		for {
			select {
			case signature := <-maestro.Signature:
				maestro.appendSignature(signature)
			case instruction := <-maestro.Instruction:
				maestro.newInstruction(instruction)
			}
		}
	}()
	return &maestro
}
