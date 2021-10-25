package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/blockchain"
	"github.com/Aereum/aereum/core/crypto"
)

type AttendeeNetwork struct {
	attendees map[crypto.Hash]*SecureConnection
	comm      chan *blockchain.Block
}

func NewAttendeeNetwork(port int,
	prvKey crypto.PrivateKey, comm chan *blockchain.Block, validator chan ValidatedConnection) AttendeeNetwork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := AttendeeNetwork{
		attendees: make(map[crypto.Hash]*SecureConnection),
		comm:      comm,
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
					network.handleAttendeeConnection(secureConnection, comm)
				}
			}
		}
	}()
	return network
}

func (m AttendeeNetwork) handleAttendeeConnection(conn *SecureConnection, comm chan *blockchain.Block) {
	m.attendees[conn.hash] = conn
	go func() {
		_, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			delete(m.attendees, conn.hash)
			return
		}
	}()
	for {
		block := <-comm
		for hash, conn := range m.attendees {
			if err := conn.WriteMessage(block.Serialize()); err != nil {
				conn.conn.Close()
				delete(m.attendees, hash)
			}
		}
	}
}
