package network

import (
	"github.com/Aereum/aereum/core/crypto"
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
	hash  crypto.Hash
	epoch int
	msg   []byte
}

func StartReceiveQueue(prvKey crypto.PrivateKey, validator chan HashedMessage, blockformation chan []HashedMessage, epoch int) {
	msgChan := make(chan HashedMessage)
	currentEpoch := epoch
	recentHashes := make([]map[crypto.Hash]struct{}, maxEpochReceiveMessage)
	for n := 0; n < maxEpochReceiveMessage; n++ {
		recentHashes[n] = make(map[crypto.Hash]struct{})
	}
	go func() {
		for {
			select {
			case hashMsg := <-msgChan:
				if deltaEpoch := hashMsg.epoch - currentEpoch; deltaEpoch < 100 && deltaEpoch > 0 {
					isNew := true
					for hash, _ := range recentHashes[deltaEpoch] {
						if hash.Equal(hashMsg.hash) {
							isNew = false
							break
						}
					}
					if isNew {
						validator <- hashMsg
					}
				}
			case blockMessages := <-blockformation:
				currentEpoch += 1
				// ignore all the messages marked to die at current epoch
				newMap := make(map[crypto.Hash]struct{})
				recentHashes = append([]map[crypto.Hash]struct{}{newMap}, recentHashes[1:]...)
				// clear the queue with messages formed at the block
				for _, msg := range blockMessages {
					deltaEpoch := msg.epoch - currentEpoch
					if deltaEpoch > 0 {
						delete(recentHashes[deltaEpoch], msg.hash)
					}
				}
			}
		}
	}()

	ListenTCP(messageReceiveConnectionPort, NewMessageReceiver, prvKey)
}

func NewMessageReceiver(m *MessageQueue, conn *SecureConnection) {
	m.messageSenders[conn.hash] = conn
	go func() {
		for {
			msg, err := conn.ReadMessage()
			if err != nil {
				conn.conn.Close()
				delete(m.messageSenders, conn.hash)
				return
			}
			m.msgChan <- hashedMessage{hash: crypto.Hasher(msg), msg: msg}
		}
	}()
}
