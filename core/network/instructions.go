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
					network.handleMessengerConnection(secureConnection, broker)
				}
			}
		}
	}()
	return network
}

func (net InstructionNetWork) handleMessengerConnection(conn *SecureConnection, broker InstructionBroker) {
	net[conn.hash] = conn
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			conn.conn.Close()
			delete(net, conn.hash)
			return
		}
		hashed := HashedInstructionBytes{msg: data}
		hashed.hash = crypto.Hasher(data)
		hashed.epoch = int(instructions.GetEpochFromByteArray(data))
		broker <- &hashed
	}
}
