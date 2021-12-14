package network

import (
	"fmt"
	"net"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

// pool of connections that are ready to receive primitive instructions
// it receives instructions, calcultaes the hash and sends to the instruction
// queue that will check if it is well formed, brodcast to peer network and
// send for the consensus engine for appropriate action. There is no response
// for any instruction.

type InstructionNetWork map[crypto.Hash]*SecureConnection

func NewInstructionClient(address string, prv crypto.PrivateKey, rmt crypto.Token) (*SecureConnection, error) {
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

func NewInstructionNetwork(port int, prvKey crypto.PrivateKey, broker InstructionBroker, validator chan consensus.ValidatedConnection) InstructionNetWork {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	network := make(InstructionNetWork)
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
					//network.handleMessengerConnection(secureConnection, broker)
					go InstructionConnectionHandler(secureConnection, broker)
				}
			}
		}
	}()
	return network
}

func InstructionConnectionHandler(conn *SecureConnection, broker InstructionBroker) {
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			return
		}
		hashed := HashedInstructionBytes{msg: data}
		hashed.hash = crypto.Hasher(data)
		hashed.epoch = int(instructions.GetEpochFromByteArray(data))

		broker <- &hashed
	}
}
