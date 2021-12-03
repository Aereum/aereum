package main

import (
	"github.com/Aereum/aereum/core/crypto"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/consensus/authority"
	"github.com/Aereum/aereum/core/network"
)

func main() {
	chain, token := consensus.NewGenesisBlockChain()
	consensus := authority.NewProofOfAtuhority(chain, token)
	network.NewNode(token, make(map[crypto.PublicKey]string), consensus, 0)

}
