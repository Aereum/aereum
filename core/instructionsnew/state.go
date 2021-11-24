package instructionsnew

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
)

type State struct {
	Epoch           uint64
	Members         store.HashVault
	Captions        store.HashVault
	Wallets         store.Wallet
	Audiences       store.Audience
	SponsorOffers   store.Sponsor
	SponsorGranted  store.Sponsor
	PowerOfAttorney store.HashVault
	EphemeralTokens store.HashExpireVault
	SponsorExpire   map[uint64]crypto.Hash
	EphemeralExpire map[uint64]crypto.Hash
}

/*
func (s *Block) ValidadeAcceptJoinAudience(msg *instruction.Message) bool {
	acceptJoinAudience := msg.AsAcceptJoinAudience()
	if acceptJoinAudience == nil {
		return false
	}
	// check if moderator signature is valid
	request := acceptJoinAudience.Request.AsJoinAudience()
	if request == nil {
		return false
	}
	ok, keys := s.state.GetAudince(crypto.Hasher(request.Audience))
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
	hash := crypto.Hasher(append(request.Audience, acceptJoinAudience.Request.Author...))
	if !s.mutations.SetNewHash(hash) {
		return false
	}
	return s.IncorporateMessage(msg)
}

func (s *Block) ValidadeAudienceChange(msg *instruction.Message) bool {
	audienceChange := msg.AsChangeAudience()
	if audienceChange == nil {
		return false
	}
	if !s.mutations.SetNewAudience(crypto.Hasher(audienceChange.Audience), append(audienceChange.Moderate, audienceChange.Submit...)) {
		return false
	}
	return s.IncorporateMessage(msg)
}

func (s *Block) ValidadeJoinAudience(msg *instruction.Message) bool {
	joinAudience := msg.AsJoinAudience()
	if joinAudience == nil {
		return false
	}
	hash := crypto.Hasher(joinAudience.Audience)
	if ok, _ := s.state.GetAudince(hash); !ok {
		return false
	}
	hash = crypto.Hasher(append(hash[:], msg.Author...))
	if !s.mutations.SetNewHash(hash) {
		return false
	}
	return s.IncorporateMessage(msg)
}


func (b *Block) SerializeWithoutHash() []byte {
	serialized := b.Parent[:]
	instructions.PutByteArray(b.Publisher, &serialized)
	instructions.PutUint64(b.Epoch, &serialized)
	instructions.PutUint64(uint64(len(b.Instructions)), &serialized)
	for _, msg := range b.Instructions {
		instructions.PutByteArray(msg, &serialized)
	}
	instructions.PutUint64(uint64(b.PublishedAt.UnixNano()), &serialized)
	return serialized
}

func (b *Block) Serialize() ([]byte, crypto.Hash) {
	serialized := b.SerializeWithoutHash()
	hash := crypto.Hasher(serialized)
	return append(serialized[0:crypto.Size], hash[:]...), hash
}

func ParseBlock(data []byte) *Block {
	block := &Block{}
	block.Parent = crypto.BytesToHash(data[0:crypto.Size])
	position := crypto.Size
	block.Publisher, position = instructions.ParseByteArray(data, position)
	block.PublishedAt, position = instructions.ParseTime(data, position)
	var count uint64
	count, position = instructions.ParseUint64(data, position)
	block.Instructions = make([][]byte, int(count))
	for n := 0; n < int(count); n++ {
		block.Instructions[n], position = instructions.ParseByteArray(data, position)
	}
	if len(data)-position != crypto.Size {
		return nil
	}
	block.Hash = crypto.BytesToHash(data[position:])
	return block
}

func (m *Block) CanPay(payments instruction.Payment) bool {
	for n, debitAcc := range payments.DebitAcc {
		stateBalance := m.state.Balance(debitAcc)
		if delta, ok := m.mutations.DeltaWallets[debitAcc]; ok {
			if int(stateBalance)+delta < int(payments.DebitValue[n]) {
				return false
			}
		} else {
			if stateBalance < int(payments.DebitValue[n]) {
				return false
			}
		}
	}
	return true
}

func (m *Block) TransferPayments(payments instruction.Payment) {
	for n, debitAcc := range payments.DebitAcc {
		if delta, ok := m.mutations.DeltaWallets[debitAcc]; ok {
			m.mutations.DeltaWallets[debitAcc] = delta - int(payments.DebitValue[n])
		} else {
			m.mutations.DeltaWallets[debitAcc] = -int(payments.DebitValue[n])
		}
	}
	for n, creditAcc := range payments.CreditAcc {
		if delta, ok := m.mutations.DeltaWallets[creditAcc]; ok {
			m.mutations.DeltaWallets[creditAcc] = delta + int(payments.CreditValue[n])
		} else {
			m.mutations.DeltaWallets[creditAcc] = int(payments.CreditValue[n])
		}
	}
}

func (m *Block) IncorporateMessage(msg *instruction.Message) bool {
	payment := msg.Payments()
	if !m.CanPay(payment) {
		return false
	}
	m.TransferPayments(payment)
	return true
}
*/
