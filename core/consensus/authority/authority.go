package authority

import (
	"time"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
)

var blockInterval = time.Second

func NewProofOfAtuhority(chain *consensus.BlockChain, token crypto.PrivateKey) *consensus.Communication {
	comm := consensus.NewCommunication()
	pool := consensus.NewInstructionPool()
	go func() {
		for {
			select {
			case peer := <-comm.PeerRequest:
				peer.Response <- false
			case <-comm.BlockSignature:
				// do nothing
			case <-comm.Checksum:
				// do nothing
			case sync := <-comm.Synchronization:
				sync.Ok <- false
			case hashedInst := <-comm.Instructions:
				pool.Queue(hashedInst.Instruction, hashedInst.Hash)
			}
		}
	}()

	go func() {
		epoch := chain.Epoch + 1
		for {
			nextBlock := time.Now().Add(consensus.IntervalToNewEpoch(epoch, chain.GenesisTime))
			newBlock := <-consensus.BlockBuilder(chain.GetLastCheckpoint(), epoch, token, nextBlock, pool)
			chain.CurrentState.IncorporateBlock(newBlock)
		}
	}()

	return comm
}
