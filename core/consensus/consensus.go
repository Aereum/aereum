package consensus

import (
	"bytes"
	"sort"
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

/*
	Interface:

		Receives New Instructions
		Receives New Blocks
		Receives Block signatures
		Peer validation request
*/

type PeerRequest struct {
	Token    crypto.Hash
	Response chan bool
}

type Consensus struct {
	PeerRequest    chan *PeerRequest
	NewBlock       chan *instructions.Block
	BlockSignature chan *Signature
	Checkpoint     chan *Checkpoint
	blockChain     *BlockChain
	pool           *InstructionPool
}

func NewConsensus(blockchain *BlockChain) *Consensus {
	consensus := Consensus{
		PeerRequest:    make(chan *PeerRequest),
		NewBlock:       make(chan *instructions.Block),
		BlockSignature: make(chan *Signature),
		Checkpoint:     make(chan *Checkpoint),
		blockChain:     blockchain,
		pool:           NewInstructionPool(),
	}
	go func() {
		for {
			select {
			case peer := <-consensus.PeerRequest:
				peer.Response <- consensus.blockChain.Engine.RegisterPeer(peer.Token)
			case block := <-consensus.NewBlock:
				consensus.PushNewBlock(block)
			case signature := <-consensus.BlockSignature:
				consensus.blockChain.AppendSignature(*signature)
			}
		}
	}()
	return &consensus
}

func (c *Consensus) PushNewBlock(block *instructions.Block) {
	epoch := block.Epoch
	newSignedBlock := SignedBlock{
		Block:      block,
		Signatures: []Signature{},
	}
	if blocks, ok := c.blockChain.CandidateBlocks[epoch]; ok {
		blocks = append(blocks, &newSignedBlock)
	} else {
		c.blockChain.CandidateBlocks[epoch] = SignedBlocks{&newSignedBlock}
	}
	consensus := c.blockChain.Engine.GetConsensus(c.blockChain.CandidateBlocks[epoch])
	if consensus != nil {
		c.blockChain.RecentBlocks = append(c.blockChain.RecentBlocks, &newSignedBlock)
		delete(c.blockChain.CandidateBlocks, epoch)
		sort.Sort(c.blockChain.RecentBlocks)
		c.blockChain.IncorporateBlocks()
	}
}

type ConsensusEngine interface {
	RegisterPeer(crypto.Hash) bool
	DropPeer(crypto.Hash)
	BlockFormationSignal() chan uint64
	IsBlockLeader(uint64) (bool, time.Time)
	IsBlockValidator(uint64) bool
	GetConsensus([]*SignedBlock) *SignedBlock
}

type BlockChain struct {
	GenesisTime     time.Time
	TotalStake      uint64
	Epoch           uint64
	CurrentState    *instructions.State
	RecentBlocks    SignedBlocks
	CandidateBlocks map[uint64]SignedBlocks
	Engine          ConsensusEngine
	AcceptPeers     bool
	MinimumStake    uint64
}

func (b *BlockChain) IncoporateBlocks() {

}

func IntervalToNewEpoch(epoch uint64, genesis time.Time) time.Duration {
	return time.Until(genesis.Add(time.Duration(int64(epoch) * 1000000000)))
}

func (b *BlockChain) AppendSignature(signature Signature) {
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

func LauchNewGenesisConsensus(egine ConsensusEngine) {
	//pool := NewInstructionPool()

	//processInstruction := make(chan instructions.Instruction)
	go func() {
		for {
			select {
			//case newInstruction := <-processInstruction:
			//pool.Queue(newInstruction)
			}
		}
	}()
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
