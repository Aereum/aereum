package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
)

const (
	validatorNetworkMsg byte = iota
	broadcastValidatedMsg
	broadcastBlockMsg
	pingMsg
	pongMsg
)

type ValidatorNetwork map[crypto.Hash]*SecureConnection

func (v ValidatorNetwork) Broadcast(msg []byte, msgType byte, token crypto.Hash) {
	for _, peer := range v {
		/*msgToSend := []byte{msgType}
		msgToSend = append(msgToSend, token[:]...)
		msgToSend = append(msgToSend, msg...)
		peer.WriteMessage(msgToSend)*/
		peer.WriteMessage(msg)

	}
}

func NewValidatorNetwork(port int, prvKey crypto.PrivateKey, comm chan *HashedMessage,
	dial map[crypto.PublicKey]string) ValidatorNetwork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := make(ValidatorNetwork)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey)
				if err != nil {
					conn.Close()
				} else {
					network[secureConnection.hash] = secureConnection
					handleValidatorConnection(secureConnection, comm)
				}
			}
		}
	}()
	for publicKey, address := range dial {
		go func() {
			net, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn, err := PerformClientHandShake(net, prvKey, publicKey)
			if err != nil {
				return
			}
			network[crypto.Hasher(publicKey.ToBytes())] = conn
			handleValidatorConnection(conn, comm)
		}()
	}
	return network
}

func handleValidatorConnection(conn *SecureConnection, comm chan *HashedMessage) {
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			return
		}
		hashed := HashedMessage{msg: data}
		hashed.hash, hashed.epoch = message.GetHashAndEpochFromMessage(data)
		comm <- &hashed
	}
}
