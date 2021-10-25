package blockchain

import (
	"bytes"

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
	AdvertisingOffers map[uint64]wallet.HashVault
	PowerOfAttorney   wallet.HashVault // power of attonery token hash
	Frozen            *StateMutations
	IsMutating        bool
	Mutations         *StateMutations
}

func (s *State) AuthorExists(m *message.Message) bool {
	return s.Subscribers.Exists(crypto.Hasher(m.Author))
}

func (s *State) ValidadeSubscribe(msg *message.Message) bool {
	subscribe := msg.AsSubscribe()
	if subscribe == nil {
		return false
	}
	// token must be new... caption must be new.
	isNew := s.Subscribers.Exists(crypto.Hasher(msg.Author)) ||
		s.Captions.Exists(crypto.Hasher([]byte(subscribe.Caption)))
	if isNew {
		return false
	}
	if s.IsMutating {
		if !s.Mutations.SetNewSubscriber(crypto.Hasher(subscribe.Token), crypto.Hasher([]byte(subscribe.Caption))) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidateAbout(msg *message.Message) bool {
	about := msg.AsAbout()
	if about == nil {
		return false
	}
	// no further tests are necessary
	if s.IsMutating {
		hash := crypto.Hasher(msg.Author)
		if !s.Mutations.SetNewHash(hash) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadeCreateAudience(msg *message.Message) bool {
	createAudience := msg.AsCreateAudiece()
	if createAudience == nil {
		return false
	}
	// must be a new audience token
	hash := crypto.Hasher(createAudience.Token)
	if s.Audiences.Exists(crypto.Hasher(createAudience.Token)) {
		return false
	}
	if s.IsMutating {
		if !s.Mutations.SetNewAudience(hash, append(createAudience.Moderate, createAudience.Submit...)) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadeJoinAudience(msg *message.Message) bool {
	joinAudience := msg.AsJoinAudience()
	if joinAudience == nil {
		return false
	}
	hash := crypto.Hasher(joinAudience.Audience)
	if !s.Audiences.Exists(hash) {
		return false
	}
	if s.IsMutating {
		if !s.Mutations.SetNewHash(hash) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadeAcceptJoinAudience(msg *message.Message) bool {
	acceptJoinAudience := msg.AsAcceptJoinAudience()
	if acceptJoinAudience == nil {
		return false
	}
	// check if moderator signature is valid
	request := acceptJoinAudience.Request.AsJoinAudience()
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
	if !moderator.Verify(request.Serialize(), acceptJoinAudience.ModeratorSignature) {
		return false
	}
	if s.IsMutating {
		hash := crypto.Hasher(append(request.Audience, acceptJoinAudience.Request.Author...))
		if !s.Mutations.SetNewHash(hash) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadeAudienceChange(msg *message.Message) bool {
	audienceChange := msg.AsChangeAudience()
	if audienceChange == nil {
		return false
	}
	if s.IsMutating {
		if !s.Mutations.SetNewAudience(crypto.Hasher(audienceChange.Audience), append(audienceChange.Moderate, audienceChange.Submit...)) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadateAdvertisingOffer(msg *message.Message) bool {
	advertisingOffer := msg.AsAdvertisingOffer()
	if advertisingOffer == nil {
		return false
	}
	if advertisingOffer.Expire <= s.Epoch {
		return false
	}
	if !s.Audiences.Exists(crypto.Hasher(advertisingOffer.Audience)) {
		return false
	}
	if s.IsMutating {
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidateContent(msg *message.Message) bool {
	m := msg.AsContent()
	if m == nil {
		return false
	}
	// check signatures
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
		offer := m.AdvertisingOffer
		hashed := crypto.Hasher(m.AdvertisingOffer.Serialize())
		if advHashStore, ok := s.AdvertisingOffers[m.AdvertisingOffer.Expire]; ok {
			if !(advHashStore.Exists(hashed)) {
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
		if s.IsMutating {
			if !s.Mutations.SetNewUseAdvOffer(hashed) {
				return false
			}
			return s.Mutations.IncorporateMessage(msg)
		}
		return true
	}
	return false
}

func (s *State) ValidateGrantPowerOfAttorney(msg *message.Message) bool {
	grantPower := msg.AsGrantPowerOfAttorney()
	if grantPower == nil {
		return false
	}
	if !s.AuthorExists(msg) {
		return false
	}
	// TOCHECK is it possible to recycle a Grant PoA after revoking?
	if s.IsMutating {
		hash := crypto.Hasher(append(msg.Author, grantPower.Token...))
		if !s.Mutations.SetNewGrantPower(hash) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) ValidadeRevokePowerOfAttorney(msg *message.Message) bool {
	revokePower := msg.AsRevokePowerOfAttorney()
	if revokePower == nil {
		return false
	}
	hash := crypto.Hasher(append(msg.Author, revokePower.Token...))
	if !s.PowerOfAttorney.Exists(hash) {
		return false
	}
	if s.IsMutating {
		if !s.Mutations.SetNewRevokePower(hash) {
			return false
		}
		return s.Mutations.IncorporateMessage(msg)
	}
	return true
}

func (s *State) Validate(info []byte) bool {
	if message.IsTransfer(info) {
		return s.ValidateTransfer(info)
	}
	if message.IsMessage(info) {
		return s.ValidateMessage(info)
	}
	return false
}

func (s *State) ValidateTransfer(trf []byte) bool {
	transfer, err := message.ParseTranfer(trf)
	if err != nil {
		return false
	}
	if s.IsMutating {
		payments := transfer.Payments()
		if !s.Mutations.CanPay(payments) {
			return false
		}
		s.Mutations.TransferPayments(payments)
		s.Mutations.transfers = append(s.Mutations.transfers, transfer)
	}
	return true
}

// State incorporates only the necessary information on the blockchain to
// validate new messages. It should be used on validation nodes.
func (s *State) ValidateMessage(msg []byte) bool {
	if !message.IsMessage(msg) {
		return false
	}
	parsed, err := message.ParseMessage(msg)
	if parsed == nil || err != nil {
		return false
	}
	if !s.Subscribers.Exists(crypto.Hasher(parsed.Author)) &&
		message.MessageType(msg) != message.SubscribeMsg {
		return false
	}
	switch message.MessageType(msg) {
	case message.SubscribeMsg:
		return s.ValidadeSubscribe(parsed)
	case message.AboutMsg:
		return s.ValidateAbout(parsed)
	case message.CreateAudienceMsg:
		return s.ValidadeCreateAudience(parsed)
	case message.JoinAudienceMsg:
		return s.ValidadeJoinAudience(parsed)
	case message.AcceptJoinAudienceMsg:
		return s.ValidadeAcceptJoinAudience(parsed)
	case message.AudienceChangeMsg:
		return s.ValidadeAudienceChange(parsed)
	case message.AdvertisingOfferMsg:
		return s.ValidadateAdvertisingOffer(parsed)
	case message.ContentMsg:
		return s.ValidateContent(parsed)
	case message.GrantPowerOfAttorneyMsg:
		return s.ValidateGrantPowerOfAttorney(parsed)
	case message.RevokePowerOfAttorneyMsg:
		return s.ValidateGrantPowerOfAttorney(parsed)
	}
	return true
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
