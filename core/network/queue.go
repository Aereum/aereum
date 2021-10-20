package network

import (
	"github.com/Aereum/aereum/core/crypto"
)

const maxEpochReceiveMessage = 100

type HashedMessage struct {
	nonpeer bool
	msg     []byte
	hash    crypto.Hash
	epoch   int
}

type ValidateMessage struct {
	msg *HashedMessage
	ok  chan ValidatedMessage
}

type ValidatedMessage struct {
	msg *HashedMessage
	ok  bool
}

// ReceiveQueue spins a goroutine that receives messages strips out duplicated
// messages, send to validator.
func ReceiveQueue(validator chan ValidateMessage, blockformation chan struct{},
	epoch int) chan *HashedMessage {

	msgChan := make(chan *HashedMessage)
	validatedChan := make(chan ValidatedMessage)
	currentEpoch := epoch
	// stores all received hashes for each recent epoch
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	go func() {
		for {
			select {
			case hashMsg := <-msgChan:
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
						validator <- ValidateMessage{msg: hashMsg, ok: validatedChan}
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
	return msgChan
}
