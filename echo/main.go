package main

import (
	"fmt"
	"log"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/network"
)

var beatKey = []byte{
	48, 72, 2, 65, 0, 179, 59, 167, 214, 225, 1, 89, 231, 2, 23, 190,
	49, 1, 161, 85, 93, 45, 21, 46, 182, 248, 160, 129, 237, 169, 176,
	15, 99, 212, 242, 150, 204, 103, 186, 151, 44, 53, 48, 209, 179,
	93, 157, 143, 231, 71, 251, 185, 128, 178, 222, 249, 133, 230, 211,
	76, 3, 216, 218, 124, 239, 243, 80, 133, 245, 2, 3, 1, 0, 1}

func main() {
	_, token := crypto.RandomAsymetricKey()
	beatPubKey, err := crypto.PublicKeyFromBytes(beatKey)
	if err != nil {
		log.Fatal("invalid remote key")
	}
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
		block := instructions.ParseBlock(data)
		if block != nil {
			fmt.Println(block.JSONSimple())
		} else {
			fmt.Println(".....")
		}
	}
}
