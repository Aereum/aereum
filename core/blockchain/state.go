package blockchain

import (
	"bytes"
	"sync"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
	"github.com/Aereum/aereum/core/wallet"
)

const maxAdvertisingOfferDelay = 1000

type State struct {
	Epoch             uint64
	Subscribers       wallet.HashVault // subscriber token hash
	Captions          wallet.HashVault // caption string hash
	Wallets           wallet.Wallet    // wallet token hash
	Audiences         wallet.Audience  // audience + Follower hash
	AdvertisingOffers map[uint64]wallet.HashStore
	PowerOfAttorney   wallet.HashStore // power of attonery token hash
	*sync.Mutex
}

// State incorporates only the necessary information on the blockchain to
// validate new messages. It should be used on validation nodes.
func (s *State) Validate(msg []byte) bool {
	if !message.IsMessage(msg) {
		return false
	}
	message, err := message.ParseMessage(msg)
	if message == nil || err != nil {
		return false
	}
	switch message.MessageType(msg) {
	case SubscribeMsg:
		subscribe := message.AsSubscribe()
		if subscribe == nil {
			return false
		}
		return !s.AuthorExists(msg)
		//
	case AboutMsg:
		about := message.AsAbout()
		if about == nil {
			return false
		}
		return s.AuthorExists(msg)
		//
	case CreateAudienceMsg:
		createAudience := message.AsCreateAudiece()
		if createAudience == nil {
			return false
		}
		return !s.Audiences.Exists(crypto.Hasher(createAudince.Token))
		//
	case JoinAudienceMsg:
		joinAudience := message.AsJoinAudience()
		if joinAudience == nil {
			return false
		}
		return s.Audiences.Exists(crypto.Hasher(joinAudience.Audience))
		//
	case AcceptJoinAudienceMsg:
		acceptJoinAudience := message.AsAcceptJoinAudience()
		if acceptJoinAudience == nil {
			return false
		}
		return CanAcceptJoinAudience(acceptJoinAudience)
		//
	case AudienceChangeMsg:
		audienceChange := message.AsChangeAudience()
		if audienceChange == nil {
			return false
		}
		// TODO
	case AdvertisingOfferMsg:
		advertisingOffer := message.AsAdvertisingOffer()
		if advertisingOffer == nil {
			return false
		}
		// TODO
	case ContentMsg:
		content := message.AsContent()
		if content == nil {
			return false
		}
		return s.CanPublish(content)
		//
	case GrantPowerOfAttorneyMsg:
		grantPower := message.AsGrantPowerOfAttorney()
		if grantPower == nil {
			return false
		}
		if !s.AuthorExists(crypto.Hasher(grantPower.Token)) {
			return false
		}
		join := append(message.Author, grantPower.Token)
		return !s.PowerOfAttorney.Exists(crypto.Hasher(join))
		//
	case RevokePowerOfAttorneyMsg:
		revokePower := message.AsRevokePowerOfAttorney()
		if revokePower == nil {
			return false
		}
		join := append(message.Author, grantPower.Token)
		return s.PowerOfAttorney.Exists(crypto.Hasher(join))
		//
	}
	return true
}

func (s *State) AuthorExists(m *message.Message) bool {
	return s.Subscribers.Exists(crypto.Hasher(Author))
}

func (s *State) CanPublish(m *message.Content) bool {
	ok, keys := s.Audiences.GetKeys(crypto.Hasher(m.Audience))
	if !ok {
		return false
	}
	submissionPub, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	if !submissionPub.VerifyHash(m.SubmitHash, m.SubmitSignature) {
		return false
	}
	if len(m.PublishSignature) > 0 {
		pulishPub, err := crypto.PublicKeyFromBytes(keys[crypto.PublicKeySize : 2*crypto.PublicKeySize])
		if err != nil {
			return false
		}
		if !pulishPub.VerifyHash(m.PublishHash, m.PublishSignature) {
			return false
		}
	}
	// does not check if the advertisement offer has resources in the walltet to
	// pay, only if the offer exists and the content matches
	if m.AdvertisingOffer != nil {
		// check if living adv offer is within the block state
		hashed := crypto.Hasher(m.AdvertisingOffer.Serialize())
		if advHashStore, ok := s.AdvertisingOffers[m.AdvertisingOffer.Expire]; ok {
			if !advHashStore.Exists(hashed) {
				return false
			}
		} else {
			return false
		}
		// check content
		if !bytes.Equal(offer.Audience, m.Audience) {
			return false
		}
		if offer.ContentType != m.ContentType {
			return false
		}
		if !bytes.Equal(offer.ContentData, m.ContentData) {
			return false
		}
		return true
	} else {
		return false
	}
	return true
}

func (s *State) CanAcceptJoinAudience(m *message.AcceptJoinAudience) bool {
	request := m.Request.AsJoinAudience()
	if request == nil {
		return false
	}
	ok, keys := s.Audiences.GetKeys(crypto.Hasher(request.Audience))
	if !ok {
		return false
	}
	moderator, err := crypto.PublicKeyFromBytes(keys[0:crypto.PublicKeySize])
	if err != nil {
		return false
	}
	return moderator.Verify(request.Serialize(), m.ModeratorSignature)
}

/*if message.IsTransfer(msg) {
	transfer, _ := message.ParseTranfer(msg)
	if transfer != nil {
		return false
	}
	if s.Debit(crypto.Hasher(transfer.From), transfer.Value+transfer.Fee) {
		s.Credit(transfer.To, int(transfer.Value))
		s.transfers = append(s.transfers, transfer)
		return true
	}
}*/

//for n, acc := range payments.DebitAcc {
//	s.Debit(acc, int(payments.DebitValue[n]))
//}
//for n, acc := range payments.CreditAcc {
//	s.Credit(acc, int(payments.CreditValue[n]))
//}
//payments := message.Payments()
//if !s.CanPay(payments) {
//	return false
//}
