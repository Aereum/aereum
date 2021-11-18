package network

import (
	"time"

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

var BlockWindow, _ = time.ParseDuration("1s")
var GenesisTime = time.Date(2021, time.November, 18, 0, 0, 0, 0, time.UTC)

type MsgValidator struct {
	msg []byte
	ok  chan bool
}

type MsgValidatorChan chan *MsgValidator

type Node struct {
	State      blockchain.State
	Validators ValidatorNetwork
	Messengers InstructionNetwork
	Atendees   AttendeeNetwork
}

func NewNode(state blockchain.State,
	prvKey crypto.PrivateKey,
	trusted map[crypto.PublicKey]string) *Node {
	//

	trustedConnections := ConnectTCPPool(trusted, prvKey)
	instructionQueue := NewInstructionQueue(prvKey)

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

	node.Messengers = *InstructionNetwork(
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
