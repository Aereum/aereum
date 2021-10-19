package blockchain

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/message"
)

type StateMutations struct {
	State        *State
	DeltaWallets map[crypto.Hash]int
	Hashes       map[crypto.Hash]struct{} // hashes of incorporation into state
	messages     []*message.Message
	transfers    []*message.Transfer
}

func (s *StateMutations) SetNewHash(hash crypto.Hash) bool {
	if _, ok := s.Hashes; ok {
		return false
	}
	s.Hashes[hash] = struct{}{}
	return true
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
		}
	} else if msg.MessageType == message.GrantPowerOfAttorneyMsg {

		hashed := crypto.Hasher(append(msg.Author))

	}
	m.messages = append(m.messages, msg)
	m.TransferPayments(payments)
	return true
}
