package main

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
)

var beatPubKey = crypto.Token{
	209, 223, 10, 121, 58, 83, 59, 194, 78, 158, 215, 85, 205, 174, 40, 196,
	47, 41, 218, 173, 89, 50, 139, 155, 130, 24, 102, 241, 51, 69, 156, 236,
}

var echoSeed = crypto.Token{
	52, 223, 10, 121, 58, 83, 59, 194, 78, 158, 215, 85, 205, 174, 40, 196,
	47, 41, 218, 173, 89, 50, 139, 155, 130, 24, 102, 241, 51, 69, 156, 236,
}

func main() {
	token := crypto.PrivateKeyFromSeed(echoSeed)
	fmt.Println(token.PublicKey())
	broker := InstructionBroker(token)
	db := NewDB(broker)
	BlockListener(token, db)
	Serve(7900, token, db)
}
