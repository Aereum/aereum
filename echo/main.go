package main

import (
	"fmt"
	"log"

	"github.com/Aereum/aereum/core/chain"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/network"
)

var beatPubKey = crypto.Token{
	209, 223, 10, 121, 58, 83, 59, 194, 78, 158, 215, 85, 205, 174, 40, 196,
	47, 41, 218, 173, 89, 50, 139, 155, 130, 24, 102, 241, 51, 69, 156, 236,
}

func main() {
	_, token := crypto.RandomAsymetricKey()
	conn, err := network.NewAttendeeClient(":7801", token, beatPubKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ok")
	for {
		data, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("connection error")
		}
		block := chain.ParseBlock(data)
		if block != nil {
			fmt.Println(block.JSONSimple())
		} else {
			fmt.Println(".....")
		}
	}
}
