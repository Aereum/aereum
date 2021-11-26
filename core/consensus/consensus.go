package consensus

import (
	"bytes"
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

/*
	Interface:

		Receives New Instructions
		Receives New Blocks
		Receives Block signatures
		Peer validation request
*/

type BlockSignature struct {
	Hash      crypto.Hash
	Token     []byte
	Stake     uint64
	Signature []byte
	Weight    float64
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
	Weight(stake, totalStake uint64) float64
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
	Engine              Consensual
	AcceptPeers         bool
	MinimumStake        uint64
}

func (b *BlockChain) AppendSignature(signature BlockSignature) {
	if len(b.CandidateBlocks) != 0 {
		for _, block := range b.CandidateBlocks {
			if block.Block.Hash.Equals(signature.Hash[:]) {
				for _, signBlock := range block.Signatures {
					if bytes.Equal(signBlock.Token, signature.Token) {
						return
					}
				}
				newSignature := Signature{
					Hash:      signature.Hash,
					Token:     signature.Token,
					Weight:    b.Engine.Weight(signature.Stake, b.TotalStake),
					Signature: signature.Signature,
				}
				block.Signatures = append(block.Signatures, newSignature)
				//if b.Engine.IsConsensus()
				return
			}
		}
	}
	// TODO append to recent blocks also
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

/*
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
*/
