package network

import (
	"time"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
)

// architecture consists of
//
// Peers ----|
//           |-> Instruction Queue (Parses Instruction) -> Consensus Engine
// Others ---|
//
// Conensus Engine -> Instruction Queue -> Peer Brodast if from "Others"

const maxEpochReceiveMessage = 100
const validatorBuffer = 1000

type HashedMessage struct {
	nonpeer bool // true if received from a peer instruction broadcast
	msg     []byte
	hash    crypto.Hash
	epoch   int
}

type ValidatedMessage struct {
	msg *HashedMessage
	ok  bool
}

type InstructionQueue chan *HashedMessage

func NewInstructionQueue(token crypto.PrivateKey, peers ValidatorNetwork, comm consensus.Communication) {
	//queue := make(chan *HashedMessage)
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	nextBlock := time.Now().Truncate(BlockWindow).Add(BlockWindow)
	epoch := int(time.Since(GenesisTime).Truncate(BlockWindow).Seconds())
	blockTick := time.NewTicker(time.Until(nextBlock))
	go func() {
		for {
			select {
			case hashMsg := <-queue:
				if deltaEpoch := int(hashMsg.epoch) - epoch; deltaEpoch < 100 && deltaEpoch > 0 {
					isNew := true
					for hash, _ := range recentHashes[deltaEpoch] {
						if hash.Equal(hashMsg.hash) {
							isNew = false
							break
						}
					}
					if isNew {
						recentHashes[deltaEpoch][hashMsg.hash] = struct{}{}
						queue <- hashMsg
						// if instruction was not received from peer it should be broadcasted
						if hashMsg.nonpeer {
							message := NewNetworkMessage(BroadcastInstruction(hashMsg.msg), token, false)
							peers.Broadcast(message)
						}
					}
				}
			case newBlockTime := <-blockTick.C:
				epoch = int(time.Since(GenesisTime).Truncate(BlockWindow).Seconds())
				recentHashes = append(recentHashes[1:], make(map[crypto.Hash]struct{}))
				blockTick.Reset(time.Until(newBlockTime.Truncate(BlockWindow).Add(BlockWindow)))
			}
		}
	}()
}

// ReceiveQueue spins a goroutine that receives messages strips out duplicated
// messages, send to validator.
/*
func ReceiveQueue(state blockchain.State, blockformation chan struct{}) chan *HashedMessage {
	// one channel to receive messages from peer conections
	queue := make(chan *HashedMessage)
	// one channel to send non-repeated message for validade against state
	toValidateChan := make(chan *HashedMessage, validatorBuffer) // buffered
	// one channel to receive validation of message from state
	validatedChan := make(chan ValidatedMessage)
	// stores all received hashes for each recent epoch

	currentEpoch := int(state.Epoch)
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	// receiver message go-routine
	go func() {
		for {
			select {
			case hashMsg := <-queue:
				if deltaEpoch := int(hashMsg.epoch) - currentEpoch; deltaEpoch < 100 && deltaEpoch > 0 {
					isNew := true
					for hash, _ := range recentHashes[deltaEpoch] {
						if hash.Equal(hashMsg.hash) {
							isNew = false
							break
						}
					}
					if isNew {
						recentHashes[deltaEpoch][hashMsg.hash] = struct{}{}
						toValidateChan <- hashMsg
					}
				}
			case validated := <-validatedChan:
				if validated.ok && validated.msg.nonpeer {
					// TODO broadcast
				}
			case <-blockformation:
				currentEpoch += 1
				// ignore all the messages marked to die at current epoch
				newMap := make(map[crypto.Hash]struct{})
				recentHashes = append([]map[crypto.Hash]struct{}{newMap}, recentHashes[1:]...)
			}
		}
	}()
	// validation message//connection go-routine
	go func() {
		for {
			select {
			case msg := <-toValidateChan:
				ok := state.Validate(msg.msg)
				validatedChan <- ValidatedMessage{msg: msg, ok: ok}
			case validConn := <-validConnChan:
				validConn.ok <- state.Subscribers.Exists(validConn.token)
			}
		}
	}()

	return queue
*/
