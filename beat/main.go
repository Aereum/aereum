package main

import (
	"fmt"
	"time"

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
	conns := make([]*network.SecureConnection, 10)

	for n := 0; n < 10; n++ {
		var err error
		time.Sleep(time.Millisecond)
		conns[n], err = network.NewInstructionClient(":7802", token, token.PublicKey())
		if err != nil {
			panic(err)
		}
	}

	createMsg := make([][]byte, 50000)
	for n := 0; n < len(createMsg); n++ {
		_, authors := crypto.RandomAsymetricKey()
		inst := instructions.NewSingleReciepientTransfer(token, authors.PublicKey().ToBytes(), "whatever", 10, 1, 10)
		createMsg[n] = inst.Serialize()
	}

	//conn, err := network.NewInstructionClient(":7802", token, token.PublicKey())
	//if err != nil {
	//	panic(err)
	//}

	ww := time.NewTicker(4 * time.Second)
	for n := 0; n < len(createMsg); n++ {
		time.Sleep(time.Microsecond)
		if err := conns[n%10].WriteMessage(createMsg[n]); err != nil {
			fmt.Println(err)
			break
		}
	}
	<-ww.C
}
