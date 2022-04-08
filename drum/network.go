package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

var echoToken = crypto.Token{112, 13, 121, 57, 244, 25, 96, 109, 189, 81, 99, 217, 170, 140, 153, 62,
	143, 118, 71, 18, 223, 22, 218, 85, 228, 96, 63, 81, 110, 82, 79, 65}

type InstructionBroker struct {
	Received chan instructions.Instruction
	Send     chan instructions.Instruction
	Epoch    uint64
}

func EchoConnection(port int, token crypto.PrivateKey) *InstructionBroker {
	dial, err := net.Dial("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		panic(err)
	}
	secureConnection, err := network.PerformClientHandShake(dial, token, echoToken)
	if err != nil {
		log.Fatal(err)

	}
	broker := InstructionBroker{
		Received: make(chan instructions.Instruction),
		Send:     make(chan instructions.Instruction),
	}
	go func() {
		for {
			msg, err := secureConnection.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}
			if len(msg) == 8 {
				broker.Epoch = binary.LittleEndian.Uint64(msg)
			} else {
				instruction := instructions.ParseInstruction(msg)
				if instruction != nil {
					broker.Received <- instruction
				}
			}
		}
	}()

	go func() {
		for {
			instruction := <-broker.Send
			secureConnection.WriteMessage(instruction.Serialize())
		}
	}()
	return &broker
}
