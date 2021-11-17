package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
)

type MessengerNetwork struct {
	messengers map[crypto.Hash]*SecureConnection
	comm       chan *HashedMessage
}

func NewMessengerNetwork(port int,
	prvKey crypto.PrivateKey, comm chan *HashedMessage, validator chan ValidatedConnection) *MessengerNetwork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := &MessengerNetwork{
		messengers: make(map[crypto.Hash]*SecureConnection),
		comm:       comm,
	}

	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey, validator)
				if err != nil {
					conn.Close()
				} else {
					network.handleMessengerConnection(secureConnection, comm)
				}
			}
		}
	}()
	return network
}

func (m *MessengerNetwork) handleMessengerConnection(conn *SecureConnection, comm chan *HashedMessage) {
	m.messengers[conn.hash] = conn
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			delete(m.messengers, conn.hash)
			return
		}
		hashed := HashedMessage{msg: data}
		hashed.hash, hashed.epoch = instruction.GetHashAndEpochFromMessage(data)
		comm <- &hashed
	}
}
