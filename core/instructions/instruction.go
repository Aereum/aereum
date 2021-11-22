package instructions

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)

const (
	ITransfer byte = iota
	IDeposit
	IWithdraw
	IJoinNetwork
	IUpdateInfo
	ICreateAudience
	IJoinAudience
	IAcceptJoinRequest
	IContent
	IUpdateAudience
	IGrantPowerOfAttorney
	IRevokePowerOfAttorney
	ISponsorshipOffer
	ISponsorshipAcceptance
	ICreateEphemeral
	ISecureChannel
	IReact
	IUnkown
)

var (
	ErrCouldNotParseMessage   = errors.New("could not parse message")
	ErrCouldNotParseSignature = errors.New("could not parse signature")
	ErrInvalidSignature       = errors.New("invalid signature")
)

type Payment struct {
	DebitAcc    []crypto.Hash
	DebitValue  []uint64
	CreditAcc   []crypto.Hash
	CreditValue []uint64
}

type Instruction interface {
	Kind() byte
	Epoch() uint64
	Serialize() []byte
}

type KindSerializer interface {
	Kind() byte
	Serialize() []byte
}

type Serializer interface {
	Serialize() []byte
}

func GetEpochFromByteArray(msg []byte) uint64 {
	epoch, _ := ParseUint64(msg, 1)
	return epoch
}

func CollectFees(instruction Instruction, token []byte) *Payment {
	switch v := instruction.(type) {
	case *AuthoredInstruction:
		pay := Payment{
			DebitValue:  []uint64{v.Fee},
			CreditValue: []uint64{v.Fee},
			CreditAcc:   []crypto.Hash{crypto.Hasher(token)},
		}
		if v.Wallet != nil {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Wallet)}
		} else if v.Attorney != nil {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Attorney)}
		} else {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Author)}
		}
		return &pay
	default:
		return nil
	}
}

func GetPayments(instruction Instruction) *Payment {
	switch v := instruction.(type) {
	case *AuthoredInstruction:
		pay := Payment{DebitValue: []uint64{v.Fee}}
		if v.Wallet != nil {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Wallet)}
		} else if v.Attorney != nil {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Attorney)}
		} else {
			pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.Author)}
		}
		if v.Kind() == ISponsorshipAcceptance {
			acceptance := v.AsSponsorshipAcceptance()
			sponsor := acceptance.Offer.AsSponsorshipOffer()
			pay.DebitAcc = append(pay.DebitAcc, crypto.Hasher(v.Wallet))
			pay.DebitValue = append(pay.DebitValue, sponsor.Revenue)
			pay.CreditAcc = append(pay.CreditAcc, crypto.Hasher(acceptance.Audience))
			pay.CreditValue = append(pay.CreditValue, sponsor.Revenue)
		}
		return &pay
	case *Transfer:
		total := uint64(0)
		pay := Payment{
			CreditAcc:   make([]crypto.Hash, len(v.To)),
			CreditValue: make([]uint64, len(v.To)),
		}
		for n, receipient := range v.To {
			pay.CreditAcc[n] = crypto.Hasher(v.To[n].Token)
			pay.CreditValue[n] = v.To[n].Value
			total += receipient.Value
		}
		pay.DebitAcc = []crypto.Hash{crypto.Hasher(v.From)}
		pay.DebitValue = []uint64{total}
	case *Deposit:
		pay := Payment{
			DebitAcc:    []crypto.Hash{crypto.Hasher(v.Token)},
			DebitValue:  []uint64{v.Value},
			CreditAcc:   []crypto.Hash{},
			CreditValue: []uint64{},
		}
		return &pay
	case *Withdraw:
		pay := Payment{
			CreditAcc:   []crypto.Hash{crypto.Hasher(v.Token)},
			CreditValue: []uint64{v.Value},
			DebitAcc:    []crypto.Hash{},
			DebitValue:  []uint64{},
		}
		return &pay
	}
	return nil
}

func IsAuthoredInstruction(instruction Instruction) bool {
	return instruction.Kind() >= IJoinNetwork && instruction.Kind() < IUnkown
}
