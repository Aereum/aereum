package breeze

import (
	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/message"
	"github.com/Aereum/aereum/core/network"
)

const ChecksumWindows = 50000
const ValidatorsCount = 10

type Maestro struct {
	Peers          network.ValidatorNetwork
	Slots          []network.SecureConnection
	LiveCheckPoint uint64
	Proposed       map[uint64]*CheckPoint
	Instruction    chan *message.Message
	Signature      chan *consensus.BlockSignature
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

func (maestro *Maestro) NewInstruction(m *message.Message) {
	maestro.Instruction <- m
}

func (maestro *Maestro) newInstruction(m *message.Message) {

}

func NewMaestro() *Maestro {

	maestro = Maestro{}

	go func() {

	}()

}
