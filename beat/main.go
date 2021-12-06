package main

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/consensus/authority"
	"github.com/Aereum/aereum/core/network"
)

func main() {
	chain, token := consensus.NewGenesisBlockChain()
	consensus := authority.NewProofOfAtuhority(chain, token)
	network.NewNode(token, make(map[crypto.PublicKey]string), consensus, 0)
	//conns := make([]*network.SecureConnection, 10)
	//var err error
	//for n := 0; n < 10; n++ {
	//	conns[n], err = network.NewInstructionClient(":7802", token, token.PublicKey())
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	conn, err := network.NewInstructionClient(":7802", token, token.PublicKey())
	if err != nil {
		panic(err)
	}

	authors := make([]crypto.PrivateKey, 50000)
	for n := 0; n < 50000; n++ {
		_, authors[n] = crypto.RandomAsymetricKey()
		inst := instructions.NewSingleReciepientTransfer(token, authors[n].PublicKey().ToBytes(), "whatever", 10, 1, 10)
		if err := conn.WriteMessage(inst.Serialize()); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("***********************************************************")
}
