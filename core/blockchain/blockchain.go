package blockchain

import (
	"bytes"
	"crypto/sha256"
	"sync"

	"github.com/Aereum/aereum/core/message"
)

type Blockchain struct {
	Messages []message.Message
}

type AudienceState struct {
	Token     []byte
	Followers []*message.Follower
}

type Hash [sha256.Size]byte

func Hash256(data []byte) Hash {
	return Hash(sha256.Sum256(data))
}

// State incorporates only the necessary information on the blockchain to 
// validate new messages. It should be used on validation nodes. 
type State struct {
	Epoch             uint64
	Subscribers       map[Hash]*Hash                // subscriber token hash
	Captions          map[Hash]struct{}             // caption string hash
	Wallets           map[Hash]int                  // wallet token hash
	Audiences         map[Hash]*[]*Hash             // audience hash
	AudienceRequests  map[Hash]*[]*message.Message  // audience hash
	PowerOfAttorney   map[Hash]Hash                 // power of attonery token hash
	AdvertisingOffers map[Hash]*message.Message     // advertising offer hash
	*sync.Mutex
}

func (s *State) IncorporateMutations(mut NewBlockMuttations) {
	s.Epoch += 1
	for subscriber, _ := range mut.NewSubsribers {
		s.Subscribers[subscriber] = struct{}{}
	}
	for caption, _ := range mut.NewCaptions {
		s.Captions[caption] = struct{}{}
	}
	for wallet, delta := range mut.DeltaWallets {
		s.Wallets[wallet] += delta
	}
	for  
}


type NewBlockMuttations struct {
	State                     *State
	NewSubsribers             map[Hash]struct{}
	NewCaptions               map[Hash]struct{}
	DeltaWallets              map[Hash]int
	NewAudiences              map[Hash]*[]*message.Follower
	ChangeAudicences          map[Hash]*[]*message.Follower
	NewAudinceRequests        map[Hash]*message.Message
	GrantPowerOfAttorney      map[Hash]Hash
	RevokePowerOfAttorney     map[Hash]struct{}
	NewAdvertisingOffers      map[Hash]*message.Message
	AcceptedAdvertisingOffers map[Hash]*message.Message
	Messages                  []*message.Message
	Transfers                 []*message.Transfer
}

func (s *NewBlockMuttations) Withdraw(acc Hash, value int) bool {
	funds := s.State.Wallets[acc]
	delta := s.DeltaWallets[acc]
	if funds+delta > value {
		s.DeltaWallets[acc] = delta - value
		return true
	}
	return false
}

func (s *NewBlockMuttations) Credit(acc Hash, value int) {
	delta := s.DeltaWallets[acc]
	s.DeltaWallets[acc] = delta + value
}

func (s *NewBlockMuttations) Transfer(t *message.Transfer) bool {
	hashFrom := Hash256(t.From)
	funds := s.State.Wallets[hashFrom]
	delta := s.DeltaWallets[hashFrom]
	value := int(t.Value)
	if funds+delta < value {
		return false
	}
	hashTo := Hash256(t.To)
	deltaTo := s.DeltaWallets[hashTo]
	s.DeltaWallets[hashFrom] = delta - value
	s.DeltaWallets[hashTo] = deltaTo + value
	s.Transfers = append(s.Transfers, t)
	return true
}

func (s *NewBlockMuttations) RedistributeAdvertisemenetFee(value int, author Hash, audience []*message.Follower) {
	// 100% author provisory
	s.Credit(author, value)
}

func (s *NewBlockMuttations) CanContent(m *message.Content, author, wallet Hash, fee int) bool {
	if len(m.AdvertisingToken) > 0 {
		hash := Hash256(m.AdvertisingToken)
		if offerMsg, ok := s.State.AdvertisingOffers[hash]; ok {
			if _, ok := s.AcceptedAdvertisingOffers[hash]; ok {
				return false // message already reclaimed in the new block
			}
			// check if advertising claim is valid
			offer := offerMsg.AsAdvertisingOffer()
			if !bytes.Equal(offer.Audience, m.Audience) {
				return false
			}
			if offer.ContentType != m.ContentType {
				return false
			}
			if !bytes.Equal(offer.ContentData, m.ContentData) {
				return false
			}
			// check if advertiser wallet have funds to pay
			value := int(offer.AdvertisingFee)
			if !s.Withdraw(Hash256(offerMsg.FeeWallet), value) {
				return false
			}
			// use protocol redistribution rule
			s.RedistributeAdvertisemenetFee(value, author, nil)
			// mark offer as accepted
			s.AcceptedAdvertisingOffers[hash] = offerMsg
			return true
		} else {
			return false
		}
	}
	return true
}

func (s *NewBlockMuttations) CanSubscribe(m *message.Subscribe, author Hash) bool {
	caption := Hash256([]byte(m.Caption))
	if _, ok := s.State.Subscribers[author]; ok {
		return false
	}
	if _, ok := s.NewSubsribers[author]; ok {
		return false
	}
	if _, ok := s.State.Captions[caption]; ok {
		return false
	}
	if _, ok := s.NewCaptions[caption]; ok {
		return false
	}
	s.NewSubsribers[author] = struct{}{}
	s.NewCaptions[caption] = struct{}{}
	return true
}

func (s *NewBlockMuttations) CanCreateAudience(m *message.CreateAudience) bool {
	audience := Hash256(m.Token)
	if _,ok := s.State.Audiences[audience]; ok {
		return false
	}
	if _, ok :+ s.NewAudieces[audience]; ok {
		return false
	}


}
