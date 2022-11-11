package main

import (
	"log"

	"github.com/Aereum/aereum/core/chain"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/network"
)

func BlockListener(token crypto.PrivateKey, db *DB) {
	conn, err := network.NewAttendeeClient(":7801", token, beatPubKey)
	if err != nil {
		log.Fatal(err)
	}
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("connection error")
		}
		block := chain.ParseBlock(data)
		go func(instructions [][]byte) {
			if block != nil {
				for _, msg := range block.Instructions {
					db.Incorporate(msg)
				}
			}
		}(block.Instructions)
	}
}

func InstructionBroker(token crypto.PrivateKey) chan []byte {
	conn, err := network.NewInstructionClient(":7802", token, beatPubKey)
	if err != nil {
		log.Fatal(err)
	}
	broker := make(chan []byte)
	go func() {
		for {
			msg := <-broker
			if msg[0] == 255 {

			} else if err := conn.WriteMessage(msg); err != nil {
				panic(err)
			}
		}
	}()
	return broker
}
