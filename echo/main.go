package main

import (
	"fmt"
	"log"

	"github.com/Aereum/aereum/core/chain"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/network"
)

var beatKey = []byte{
	74, 48, 72, 2, 65, 0, 164, 30, 27, 3, 105, 65, 113, 154, 203, 211, 234, 109, 171, 152, 46, 129,
	203, 144, 185, 161, 238, 20, 228, 245, 157, 246, 3, 135, 197, 187, 142, 105, 252, 15, 218,
	189, 71, 232, 55, 12, 179, 108, 217, 216, 202, 195, 114, 6, 177, 198, 223, 142, 200, 93, 165,
	134, 159, 215, 171, 117, 231, 52, 99, 3, 2, 3, 1, 0, 1, 0}

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
		block := chain.ParseBlock(data)
		if block != nil {
			fmt.Println(block.JSONSimple())
		} else {
			fmt.Println(".....")
		}
	}
}
