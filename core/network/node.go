package network

import (
	"crypto"

	"github.com/Aereum/aereum/core/blockchain"
)


type MsgValidator struct {
	msg []byte
	ok  chan bool
}

type MsgValidatorChan chan *MsgValidator 
}

type Node struct {
	State      blockchain.State
	validation chan *HashedMessage
	Validators ValidatorNetwork
}

type NodeState struct {
	epoch         int
	trustedPeers  map[crypto.Hash]string
	validadorPort int
	prvKey        crypto.PrivateKey
}

func NewNode(state NodeState) *Node {
	node := &Node{}
	blockchain.MsgValidator
	msgChan := ReceiveQueue(,) 
	node.Validators = NewValidatorNetwork(state.validadorPort, state.prvKey, ,state.trustedPeers)
}
