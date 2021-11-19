package blockchain

/*
import (
	"bytes"
	"time"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instruction"
	"github.com/Aereum/aereum/core/store"
)

const maxAdvertisingOfferDelay = 1000

type State struct {
	Hash              crypto.Hash
	Epoch             uint64
	Subscribers       store.HashVault // subscriber token hash
	Captions          store.HashVault // caption string hash
	Wallets           store.Wallet    // wallet token hash
	Audiences         store.Audience  // audience + Follower hash
	AdvertisingOffers map[uint64]store.HashVault
	PowerOfAttorney   store.HashVault // power of attonery token hash
	//Frozen            *StateMutations
	IsMutating bool
	Mutations  *StateMutations
	Key        crypto.PrivateKey
}

func (s *State) IncorporateMutations() []byte {
	for hash, delta := range s.Mutations.DeltaWallets {
		if delta > 0 {
			s.Wallets.Credit(hash, uint64(delta))
		} else if delta < 0 {
			s.Wallets.Debit(hash, uint64(-delta))
		}
	}
	for hash := range s.Mutations.GrantPower {
		s.PowerOfAttorney.Insert(hash)
	}
	for hash := range s.Mutations.RevokePower {
		s.PowerOfAttorney.Remove(hash)
	}
	for hash := range s.Mutations.NewSubscriber {
		s.Subscribers.Insert(hash)
	}
	for hash := range s.Mutations.NewCaption {
		s.Captions.Insert(hash)
	}
	for hash, keys := range s.Mutations.NewAudiences {
		s.Audiences.SetKeys(hash, keys)
	}
	for hash, epoch := range s.Mutations.UseAdvOffer {
		hashvault := s.AdvertisingOffers[epoch]
		hashvault.Remove(hash)
	}
	for hash, epoch := range s.Mutations.NewAdvOffer {
		if hashvault, ok := s.AdvertisingOffers[epoch]; ok {
			hashvault.Insert(hash)
		} else {
			//wallet.NewHashVault("advertising", epoch, )
		}
	}
	if hashvault, ok := s.AdvertisingOffers[s.Epoch+1]; ok {
		hashvault.Close()
		delete(s.AdvertisingOffers, s.Epoch+1)
	}
	s.Epoch += 1
	block := Block{
		Parent:       s.Hash,
		Publisher:    s.Key.PublicKey().ToBytes(),
		PublishedAt:  time.Now(),
		Messages:     make([][]byte, len(s.Mutations.messages)),
		Transactions: make([][]byte, len(s.Mutations.transfers)),
	}
	for n, msg := range s.Mutations.messages {
		block.Messages[n] = msg.Serialize()
	}
	for n, trf := range s.Mutations.transfers {
		block.Transactions[n] = trf.Serialize()
	}
	s.Mutations = NewStateMutation(s)
	data, hash := block.Serialize()
	s.Hash = hash
	return data
}

func (s *State) AuthorExists(m *instruction.Message) bool {
	return s.Subscribers.Exists(crypto.Hasher(m.Author))
}

func (s *State) ValidadeSubscribe(msg *instruction.Message) bool {
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

func (s *State) ValidateAbout(msg *instruction.Message) bool {
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

func (s *State) ValidadeCreateAudience(msg *instruction.Message) bool {
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

func (s *State) ValidadeJoinAudience(msg *instruction.Message) bool {
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

func (s *State) ValidadeAcceptJoinAudience(msg *instruction.Message) bool {
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

func (s *State) ValidadeAudienceChange(msg *instruction.Message) bool {
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

func (s *State) ValidadateAdvertisingOffer(msg *instruction.Message) bool {
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

func (s *State) ValidateContent(msg *instruction.Message) bool {
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
			if !s.Mutations.SetNewUseAdvOffer(hashed, offer.Expire) {
				return false
			}
			return s.Mutations.IncorporateMessage(msg)
		}
		return true
	}
	return false
}

func (s *State) ValidateGrantPowerOfAttorney(msg *instruction.Message) bool {
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

func (s *State) ValidadeRevokePowerOfAttorney(msg *instruction.Message) bool {
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
	if instruction.IsTransfer(info) {
		return s.ValidateTransfer(info)
	}
	if instruction.IsMessage(info) {
		return s.ValidateMessage(info)
	}
	return false
}

func (s *State) ValidateTransfer(trf []byte) bool {
	transfer, err := instruction.ParseTranfer(trf)
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
	if !instruction.IsMessage(msg) {
		return false
	}
	parsed, err := instruction.ParseMessage(msg)
	if parsed == nil || err != nil {
		return false
	}
	if !s.Subscribers.Exists(crypto.Hasher(parsed.Author)) &&
		instruction.MessageType(msg) != instruction.SubscribeMsg {
		return false
	}
	switch instruction.MessageType(msg) {
	case instruction.SubscribeMsg:
		return s.ValidadeSubscribe(parsed)
	case instruction.AboutMsg:
		return s.ValidateAbout(parsed)
	case instruction.CreateAudienceMsg:
		return s.ValidadeCreateAudience(parsed)
	case instruction.JoinAudienceMsg:
		return s.ValidadeJoinAudience(parsed)
	case instruction.AcceptJoinAudienceMsg:
		return s.ValidadeAcceptJoinAudience(parsed)
	case instruction.AudienceChangeMsg:
		return s.ValidadeAudienceChange(parsed)
	case instruction.AdvertisingOfferMsg:
		return s.ValidadateAdvertisingOffer(parsed)
	case instruction.ContentMsg:
		return s.ValidateContent(parsed)
	case instruction.GrantPowerOfAttorneyMsg:
		return s.ValidateGrantPowerOfAttorney(parsed)
	case instruction.RevokePowerOfAttorneyMsg:
		return s.ValidateGrantPowerOfAttorney(parsed)
	}
	return true
}
*/
