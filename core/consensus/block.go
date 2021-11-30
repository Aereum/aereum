package consensus

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type Signature struct {
	Hash      crypto.Hash
	Token     []byte
	Signature []byte
}

type SignedBlock struct {
	Block      *instructions.Block
	Signatures []Signature
}

type SignedBlocks []*SignedBlock

func (blocks SignedBlocks) Less(i, j int) bool {
	return blocks[i].Block.Epoch < blocks[j].Block.Epoch
}

func (blocks SignedBlocks) Len() int {
	return len(blocks)
}

func (blocks SignedBlocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

type Blocks []*instructions.Block

func (blocks Blocks) Less(i, j int) bool {
	return blocks[i].Epoch < blocks[j].Epoch
}

func (blocks Blocks) Len() int {
	return len(blocks)
}

func (blocks Blocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

type Checkpoint struct {
	Validator       *instructions.Validator
	CheckpointEpoch uint64
	CheckpointHash  crypto.Hash
}

type processInstruction struct {
	instruction instructions.Instruction
	valid       chan bool
}

type instructionCache struct {
	instruction instructions.Instruction
	hash        crypto.Hash
}

func BlockBuilder(checkpoint *Checkpoint, epoch uint64, token crypto.PrivateKey, finish time.Time, consensus *Consensus) {
	block := instructions.NewBlock(checkpoint.CheckpointHash, checkpoint.CheckpointEpoch, epoch, token.PublicKey().ToBytes(), checkpoint.Validator)
	stop := time.NewTicker(time.Until(finish))
	communication := make(chan processInstruction)
	running := true
	cache := make([]instructionCache, 0)
	go func() {
		for {
			select {
			case <-stop.C:
				consensus.PushNewBlock(block)
				running = false
				for _, cached := range cache {
					consensus.pool.Queue(cached.instruction, cached.hash)
				}
				return
			case process := <-communication:
				process.valid <- block.Incorporate(process.instruction)
			}
		}
	}()

	go func() {
		valid := make(chan (bool))
		for {
			if !running {
				break
			}
			newInstruction, newHash := consensus.pool.Unqueue()
			cache = append(cache, instructionCache{newInstruction, newHash})
			if newInstruction != nil {
				communication <- processInstruction{
					instruction: newInstruction,
					valid:       valid,
				}
				if <-valid {
					consensus.pool.Delete(newHash)
				}
			}
		}
	}()

}
