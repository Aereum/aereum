package network

import (
	"github.com/Aereum/aereum/core/crypto"

	"github.com/Aereum/aereum/core/blockchain"
)

const (
	validationNodePort             = 7080
	blockBroadcastPort             = 7801
	messageReceiveConnectionPort   = 7802
	messageBroadcastConnectionPort = 7803
	syncPort                       = 7804
)

type MsgValidator struct {
	msg []byte
	ok  chan bool
}

type MsgValidatorChan chan *MsgValidator

type Node struct {
	State      blockchain.State
	Validators ValidatorNetwork
	Messengers MessengerNetwork
	Atendees   AttendeeNetwork
}

func NewNode(state blockchain.State,
	prvKey crypto.PrivateKey,
	trusted map[crypto.PublicKey]string) *Node {
	//
	hashedMsgChan, validateConnChan := ReceiveQueue(state, make(chan struct{}))
	node := &Node{}
	blockChan := make(chan *blockchain.Block)
	node.Validators = NewValidatorNetwork(
		validationNodePort,
		prvKey,
		hashedMsgChan,
		validateConnChan,
		trusted,
	)

	node.Messengers = *NewMessengerNetwork(
		messageReceiveConnectionPort,
		prvKey,
		hashedMsgChan,
		validateConnChan,
	)

	node.Atendees = NewAttendeeNetwork(
		blockBroadcastPort,
		prvKey,
		blockChan,
		validateConnChan,
	)
	return node
}
