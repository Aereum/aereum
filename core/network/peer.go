package network

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
)

const maxEpochReceiveMessage = 100

type Peers struct {
	messageBroadcast map[crypto.Hash]*SecureConnection
}

func (p *Peers) BroadcastMessage(msg []byte) {
	for _, conn := range p.messageBroadcast {
		conn.WriteMessage(msg)
	}
}

type HashedMessage struct {
	msg   []byte
	hash  crypto.Hash
	epoch int
}

// Receive Queue spins a goroutine that receives all messages from all peers and
// strips out duplicated messages.
func ReceiveQueue(prvKey crypto.PrivateKey, validator chan *message.Message, blockformation chan struct{}, epoch int) {
	msgChan := make(chan *HashedMessage)
	currentEpoch := epoch
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
						parsed, err := message.ParseMessage(hashMsg.msg)
						if err != nil {
							validator <- parsed
						}
					}
				}
			case <-blockformation:
				currentEpoch += 1
				// ignore all the messages marked to die at current epoch
				newMap := make(map[crypto.Hash]struct{})
				recentHashes = append([]map[crypto.Hash]struct{}{newMap}, recentHashes[1:]...)
			}
		}
	}()
	go ListenTCP(messageReceiveConnectionPort, NewMessageReceiver(msgChan), prvKey)
}

func NewMessageReceiver(msgChan chan *HashedMessage) handlePort {
	return func(conn *SecureConnection) {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			return
		}
		hashed := HashedMessage{msg: data}
		hashed.hash, hashed.epoch = message.GetHashAndEpochFromMessage(data)
		msgChan <- &hashed
	}
}
