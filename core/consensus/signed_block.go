package consensus

import (
	"time"

	"github.com/Aereum/aereum/core/chain"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type Signature struct {
	Hash      crypto.Hash
	Token     []byte
	Signature []byte
}

type SignedBlock struct {
	Block      *chain.Block
	Signatures []Signature
}

type SignedBlocks []*SignedBlock

func (blocks SignedBlocks) Less(i, j int) bool {
	return blocks[i].Block.Epoch() < blocks[j].Block.Epoch()
}

func (blocks SignedBlocks) Len() int {
	return len(blocks)
}

func (blocks SignedBlocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

type Blocks []*chain.Block

func (blocks Blocks) Less(i, j int) bool {
	return blocks[i].Epoch() < blocks[j].Epoch()
}

func (blocks Blocks) Len() int {
	return len(blocks)
}

func (blocks Blocks) Swap(i, j int) {
	blocks[i], blocks[j] = blocks[j], blocks[i]
}

type Checkpoint struct {
	Validator       *chain.MutatingState
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

func BlockBuilder(checkpoint *Checkpoint, epoch uint64, token crypto.PrivateKey, finish time.Time, pool *InstructionPool) chan *chain.Block {
	block := chain.NewBlock(checkpoint.CheckpointHash, checkpoint.CheckpointEpoch, epoch, token.PublicKey().ToBytes(), checkpoint.Validator)
	stop := time.NewTicker(time.Until(finish))
	communication := make(chan processInstruction)
	finished := make(chan *chain.Block)
	running := true
	cache := make([]instructionCache, 0)
	go func() {
		for {
			select {
			case <-stop.C:
				finished <- block
				running = false
				for _, cached := range cache {
					pool.Queue(cached.instruction, cached.hash)
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
			newInstruction, newHash := pool.Unqueue()
			if newInstruction != nil {
				cache = append(cache, instructionCache{newInstruction, newHash})
				communication <- processInstruction{
					instruction: newInstruction,
					valid:       valid,
				}
				if <-valid {
					cache = cache[0 : len(cache)-1]
				}
			}
		}
	}()
	return finished
}
