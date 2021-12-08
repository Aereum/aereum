package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
)

type AttendeeNetwork struct {
	attendees map[crypto.Hash]*SecureConnection
	comm      chan *consensus.SignedBlock
}

func NewAttendeeClient(address string, prv crypto.PrivateKey, rmt crypto.PublicKey) (*SecureConnection, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	secure, err := PerformClientHandShake(conn, prv, rmt)
	if err != nil {
		return nil, err
	}
	return secure, nil
}

func NewAttendeeNetwork(port int,
	prvKey crypto.PrivateKey, comm *consensus.Communication) AttendeeNetwork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := AttendeeNetwork{
		attendees: make(map[crypto.Hash]*SecureConnection),
		comm:      make(chan *consensus.SignedBlock),
	}

	if err != nil {
		panic(err)
	}
	// listener loop
	go func() {
		for {
			conn, err := listener.Accept()
			if err == nil {
				secureConnection, err := PerformServerHandShake(conn, prvKey, comm.ValidateConn)
				if err != nil {
					conn.Close()
				} else {
					network.attendees[secureConnection.hash] = secureConnection
				}
			}
		}
	}()
	// write loop
	go func() {
		for {
			block := <-network.comm
			blockBytes := block.Block.Serialize()
			for hash, conn := range network.attendees {
				if err := conn.WriteMessage(blockBytes); err != nil {
					conn.conn.Close()
					delete(network.attendees, hash)
				}
			}
		}
	}()

	return network
}
