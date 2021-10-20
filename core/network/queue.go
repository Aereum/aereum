package network

import (
	"github.com/Aereum/aereum/core/blockchain"
	"github.com/Aereum/aereum/core/crypto"
)

const maxEpochReceiveMessage = 100
const validatorBuffer = 1000

type HashedMessage struct {
	nonpeer bool
	msg     []byte
	hash    crypto.Hash
	epoch   int
}

type ValidatedMessage struct {
	msg *HashedMessage
	ok  bool
}

// ReceiveQueue spins a goroutine that receives messages strips out duplicated
// messages, send to validator.
func ReceiveQueue(state blockchain.State, blockformation chan struct{}) chan *HashedMessage {
	// one channel to receive messages from peer conections
	msgReceiveChan := make(chan *HashedMessage)
	// one channel to send non-repeated message for validade against state
	toValidateChan := make(chan *HashedMessage, validatorBuffer) // buffered
	// one channel to receive validation of message from state
	validatedChan := make(chan ValidatedMessage)
	// stores all received hashes for each recent epoch
	currentEpoch := state.epoch
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	// receiver message go-routine
	go func() {
		for {
			select {
			case hashMsg := <-msgReceiveChan:
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
						toValidateChan <- ValidateMessage{msg: hashMsg, ok: validatedChan}
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
	// validation message go-routine
	go func() {
		for {
			msg := <-toValidateChan
			ok := state.Validate(msg.msg)
			validatedChan <- ValidatedMessage{msg: msg, ok: ok}
		}
	}()

	return msgChan
}
