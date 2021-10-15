package blockchain

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
)

type StateMutations struct {
	State        *State
	DeltaWallets map[crypto.Hash]int
	messages     []*message.Message
	transfers    []*message.Transfer
}
