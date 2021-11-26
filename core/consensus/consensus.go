package consensus

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

type BlockSignature struct {
	Hash      crypto.Hash
	Token     []byte
	Stake     uint64
	Signature []byte
}

type HashInstruction struct {
	Instruction *instructions.Instruction
	Hash        crypto.Hash
}

type ConsensusEngine func(initial instructions.State, peers network.ValidatorNetwork) *Consensus

type Consensus struct {
	Signature   chan *BlockSignature
	Instruction chan instructions.Instruction
}

func (c *Consensus) NewInstruction(i instructions.Instruction) {
	c.Instruction <- i
}

func (c *Consensus) BlockSignature(sign *BlockSignature) {
	c.Signature <- sign
}

type Consensual interface {
	IsLeader(time.Time) bool
	IsValidator(time.Time) bool
	ValidateBlock(*instructions.Block) bool
	IsConsensus(*instructions.Block) bool
}

func IntervalToNewEpoch(epoch uint64, genesis time.Time) time.Duration {
	return time.Until(genesis.Add(time.Duration(int64(epoch) * 1000000000)))
}

type BlockChain struct {
	GenesisTime         time.Time
	TotalStake          uint64
	Epoch               uint64
	CurrentState        *instructions.State
	RecentBlocks        Blocks
	CandidateBlocks     SignedBlocks
	CandidateSignatures map[uint64][]BlockSignature
}

func (b *BlockChain) GetLastValidator() *instructions.Validator {
	starting := b.CurrentState.Epoch
	if len(b.RecentBlocks) == 0 || b.RecentBlocks[0].Epoch != starting+1 {
		return &instructions.Validator{
			State:     b.CurrentState,
			Mutations: instructions.NewMutation(),
		}
	}
	sequential := make([]*instructions.Block, 0)
	for _, block := range b.RecentBlocks {
		if block.Epoch != starting+1 {
			break
		}
		starting += 1
		sequential = append(sequential, block)
	}
	return &instructions.Validator{
		State:     b.CurrentState,
		Mutations: instructions.GroupBlockMutations(sequential),
	}
}

func NewGenesisConsensus(engine Consensual) {
	pool := NewInstructionPool()
	state, token := instructions.NewGenesisState()
	timeOfGenesis := time.Now()
	epoch := int64(0)
	blockTick := time.NewTimer(time.Until(timeOfGenesis.Add(time.Duration((epoch + 1) * 1000000000))))
	go func() {
		for {
			starting := <-blockTick.C
			if engine.IsLeader(starting) {
				BlockBuilder()
			}
		}
	}()
}
