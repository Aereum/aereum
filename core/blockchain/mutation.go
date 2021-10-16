package blockchain

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
)

type StateMutations struct {
	State        *State
	DeltaWallets map[crypto.Hash]int
	Captions     map[crypto.Hash]struct{}
	Audiences    map[crypto.Hash]struct{}
	AdvOffers    map[crypto.Hash]struct{}
	Attorney     map[crypto.Hash]struct{}
	messages     []*message.Message
	transfers    []*message.Transfer
}

func (m *StateMutations) CanPay(payments message.Payment) bool {
	for n, debitAcc := range payments.DebitAcc {
		ok, stateBalance := m.State.Wallets.Balance(debitAcc)
		if !ok {
			return false
		}
		if delta, ok := m.DeltaWallets[debitAcc]; ok {
			if int(stateBalance)+delta < int(payments.DebitValue[n]) {
				return false
			}
		} else {
			if stateBalance < payments.DebitValue[n] {
				return false
			}
		}
	}
	return true
}

func (m *StateMutation) TransferPayments(payments message.Payment) {
	for n, debitAcc := range payments.DebitAcc {
		if delta, ok := m.DeltaWallets[debitAcc]; ok {
			m.DeltaWallets[debitAcc] = delta - int(payments.DebitValue[n])
		} else {
			m.DeltaWallets[debitAcc] = -int(payments.DebitValue[n])
		}
	}
	for n, creditAcc := range payments.creditAcc {
		if delta, ok := m.DeltaWallets[creditAcc]; ok {
			m.DeltaWallets[creditAcc] = delta + int(payments.CreditValue[n])
		} else {
			m.DeltaWallets[creditAcc] = int(payments.CreditValue[n])
		}
	}
}

func (m *StateMutation) IncorporateSubscriber(subscribe *message.Subscribe) bool {
	if subscribe == nil {
		return false
	}
	hashed := crypto.Hasher([]byte(subscribe.Caption))
	if _, ok := m.Captions[hashed]; ok {
		return false
	}
	m.Captions[hashed] = struct{}{}
	return true
}

func (m *StateMutation) IncorporateAudienceChange(chg message.ChangeAudience) bool {
	if chg == nil {
		return false
	}
	hashed := crypto.Hasher(chg.Audience)
	if _, ok := m.Audiences[hashed]; ok {
		return false
	}
	m.Audiences[hashed] = struct{}{}
	return true
}

func (m *StateMutation) IncorporateContent(content message.Content) bool {
	if content == nil {
		return false
	}
	if adv := content.AdvertisingOffer; adv != nil {
		hashed := crypto.Hasher(adv.Serialize())
		if _, ok := m.AdvOffers[hashed]; ok {
			return false
		}
		m.AdvOffers[hashed] = struct{}{}
	}
	return true
}

func (m *StateMutation) IncorporateMessage(msg message.Message) bool {
	payment = msg.Payments()
	if !m.CanPay(payment) {
		return false
	}
	if msg.MessageType == message.SubscribeMsg {
		if !m.IncorporateSubscriber(msg.AsSubscribe()) {
			return false
		}
	} else if msg.MessageType == message.AudienceChangeMsg {
		if !m.IncorporateAudienceChange(message.AsChangeAudience()) {
			return false
		}
	} else if msg.MessageType == message.ContentMsg {
		if !m.IncorporateContent(message.AsContent()) {
			return false
	} else if msg.MessageType == message.GrantPowerOfAttorneyMsg {
		
		hashed := crypto.Hasher(append(msg.Author))

	}
	m.TransferPayments(payments)
}
