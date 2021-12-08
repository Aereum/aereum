package network

import (
	"fmt"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

const maxEpochReceiveMessage = 100

type HashedInstructionBytes struct {
	nonpeer bool // true if received from a peer instruction broadcast
	msg     []byte
	hash    crypto.Hash
	epoch   int
}

type InstructionBroker chan *HashedInstructionBytes

func NewInstructionBroker(
	token crypto.PrivateKey,
	peers *ValidatorNetwork,
	comm *consensus.Communication,
	newBlockSignal chan uint64,
	epoch uint64,
) InstructionBroker {
	broker := make(InstructionBroker)
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	currentEpoch := int(epoch)
	go func() {
		for {
			select {
			case hashInst := <-broker:
				if deltaEpoch := currentEpoch - int(hashInst.epoch); deltaEpoch < 100 && deltaEpoch >= 0 {
					if _, exists := recentHashes[deltaEpoch][hashInst.hash]; !exists {
						recentHashes[deltaEpoch][hashInst.hash] = struct{}{}
						if instruction := instructions.ParseInstruction(hashInst.msg); instruction != nil {
							comm.Instructions <- &instructions.HashInstruction{
								Instruction: instruction,
								Hash:        hashInst.hash,
							}
							if hashInst.nonpeer {
								message := NewNetworkMessage(BroadcastInstruction(hashInst.msg), token, false)
								peers.Broadcast(message)
							}
						}
						//broker <- hashInst
						// if instruction was not received from peer it should be broadcasted
					}
				}
			case newEpoch := <-newBlockSignal:
				deltaEpoch := int(newEpoch) - currentEpoch
				if deltaEpoch != 1 {
					panic(fmt.Sprintf("TODO: decide what to do... %v, %v", newEpoch, currentEpoch))
				}
				recentHashes = append(recentHashes[1:], make(map[crypto.Hash]struct{}))
				currentEpoch = int(newEpoch)
				fmt.Printf("current epoch: %v\n", currentEpoch)
			}
		}
	}()
	return broker
}
