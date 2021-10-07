package blockchain

import (
	"github.com/Aereum/aereum/core/hashdb"
	"github.com/Aereum/aereum/core/message"
)

type Header struct {
	Token        []byte
	Parent       hashdb.Hash
	ProofOfChain []byte
	Mutations    StateMutation
}

type StateMutation struct {
	parentState               *State
	NewSubsribers             map[hashdb.Hash]struct{} // hash token -> hash caption
	NewCaptions               map[hashdb.Hash]struct{}
	DeltaWallets              map[hashdb.Hash]int
	NewAudiences              map[hashdb.Hash]struct{} // Author + Token hash
	GrantPowerOfAttorney      map[hashdb.Hash]struct{}
	RevokePowerOfAttorney     map[hashdb.Hash]struct{}
	NewAdvertisingOffers      map[hashdb.Hash]*message.Message
	AcceptedAdvertisingOffers map[hashdb.Hash]struct{}
	messages                  []*[]byte
}

func (s *StateMutation) Serialize() []byte {

}
