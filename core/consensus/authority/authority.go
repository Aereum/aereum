package authority

import (
	"fmt"
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
			case validate := <-comm.ValidateConn:
				validate.Ok <- true
			}
		}
	}()

	go func() {
		epoch := chain.Epoch + 1
		for {
			fmt.Println(epoch)
			nextBlock := time.Now().Add(consensus.IntervalToNewEpoch(epoch, chain.GenesisTime))
			//fmt.Println(nextBlock)
			newBlock := <-consensus.BlockBuilder(chain.GetLastCheckpoint(), epoch, token, nextBlock, pool)
			newBlock.Sign(token)
			chain.CurrentState.IncorporateBlock(newBlock)
			comm.Checkpoint <- &consensus.SignedBlock{Block: newBlock, Signatures: make([]consensus.Signature, 0)}
			epoch += 1
		}
	}()

	return comm
}
