package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
	"github.com/Aereum/aereum/core/util"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/consensus/authority"
	"github.com/Aereum/aereum/core/network"
)

func main() {

	var token crypto.PrivateKey
	generate := false
	total := 10000
	var data []byte

	if generate {
		file, err := os.Create("teste.dat")
		if err != nil {
			log.Fatal(err)
		}
		_, token = crypto.RandomAsymetricKey()
		data := make([]byte, 0)
		util.PutByteArray(token.ToBytes(), &data)
		for n := 0; n < total; n++ {
			_, authors := crypto.RandomAsymetricKey()
			inst := instructions.NewSingleReciepientTransfer(token, authors.PublicKey().ToBytes(), "whatever", 10, 1, 10)
			util.PutByteArray(inst.Serialize(), &data)
		}
		if n, err := file.Write(data); n != len(data) || err != nil {
			file.Close()
			log.Fatal(err)
		}
	}
	file, err := os.Open("teste.dat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err = io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	position := 0
	bytes, position := util.ParseByteArray(data, position)
	token, err = crypto.PrivateKeyFromBytes(bytes)
	fmt.Println(token.PublicKey().ToBytes())
	if err != nil {
		log.Fatal(err)
	}
	createMsg := make([][]byte, total)
	for n := 0; n < len(createMsg); n++ {
		createMsg[n], position = util.ParseByteArray(data, position)
		if inst := instructions.ParseTransfer(createMsg[n]); inst == nil {
			fmt.Print(createMsg[n])
			return
		}
	}

	chain := consensus.NewGenesisBlockChain(token)
	consensus := authority.NewProofOfAtuhority(chain, token)
	network.NewNode(token, make(map[crypto.PublicKey]string), consensus, 1)

	conns := make([]*network.SecureConnection, 10)
	for n := 0; n < 10; n++ {
		var err error
		time.Sleep(time.Microsecond)
		conns[n], err = network.NewInstructionClient(":7802", token, token.PublicKey())
		if err != nil {
			panic(err)
		}
	}
	for n := 0; n < total; n++ {
		time.Sleep(10 * time.Millisecond)
		if err := conns[n%10].WriteMessage(createMsg[n]); err != nil {
			fmt.Println(err)
			break
		}
	}
}
