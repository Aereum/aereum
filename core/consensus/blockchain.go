package consensus

import (
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type BlockChain struct {
	GenesisTime     time.Time
	TotalStake      uint64
	Epoch           uint64
	CurrentState    *instructions.State
	RecentBlocks    SignedBlocks
	CandidateBlocks map[uint64]SignedBlocks
}

func (b *BlockChain) GetLastCheckpoint() *Checkpoint {
	starting := b.CurrentState.Epoch
	if len(b.RecentBlocks) == 0 || b.RecentBlocks[0].Block.Epoch != starting+1 {
		return &Checkpoint{
			Validator: &instructions.Validator{
				State:     b.CurrentState,
				Mutations: instructions.NewMutation(),
			},
			CheckpointEpoch: b.CurrentState.Epoch,
		}
	}
	sequential := make([]*instructions.Block, 0)
	for _, block := range b.RecentBlocks {
		if block.Block.Epoch != starting+1 {
			break
		}
		starting += 1
		sequential = append(sequential, block.Block)
	}
	return &Checkpoint{
		Validator: &instructions.Validator{
			State:     b.CurrentState,
			Mutations: instructions.GroupBlockMutations(sequential),
		},
		CheckpointEpoch: sequential[len(sequential)-1].Epoch,
		CheckpointHash:  sequential[len(sequential)-1].Hash,
	}
}

func NewGenesisBlockChain(token crypto.PrivateKey) *BlockChain {
	state := instructions.NewGenesisStateWithToken(token)
	chain := BlockChain{
		GenesisTime:     time.Now(),
		TotalStake:      1000000,
		Epoch:           0,
		CurrentState:    state,
		RecentBlocks:    make(SignedBlocks, 0),
		CandidateBlocks: make(map[uint64]SignedBlocks),
	}
	return &chain
}
