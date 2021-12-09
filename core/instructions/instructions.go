package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/util"
)

const (
	iTransfer byte = iota
	iDeposit
	iWithdraw
	iJoinNetwork
	iUpdateInfo
	iCreateAudience
	iJoinAudience
	iAcceptJoinRequest
	iContent
	iUpdateAudience
	iGrantPowerOfAttorney
	iRevokePowerOfAttorney
	iSponsorshipOffer
	iSponsorshipAcceptance
	iCreateEphemeral
	iSecureChannel
	iReact
	iUnkown
)

type InstructionValidator interface {
	SetNewGrantPower(hash crypto.Hash) bool
	SetNewRevokePower(hash crypto.Hash) bool
	SetNewUseSpnOffer(hash crypto.Hash) bool
	SetNewSpnOffer(hash crypto.Hash, expire uint64) bool
	SetPublishSponsor(hash crypto.Hash) bool
	SetNewEphemeralToken(hash crypto.Hash, expire uint64) bool
	SetNewMember(tokenHash crypto.Hash, captionHashe crypto.Hash) bool
	SetNewAudience(hash crypto.Hash, keys []byte) bool
	UpdateAudience(hash crypto.Hash, keys []byte) bool
	Balance(hash crypto.Hash) uint64
	PowerOfAttorney(hash crypto.Hash) bool
	SponsorshipOffer(hash crypto.Hash) uint64
	HasMember(hash crypto.Hash) bool
	HasCaption(hash crypto.Hash) bool
	HasGrantedSponser(hash crypto.Hash) (bool, crypto.Hash)
	GetAudienceKeys(hash crypto.Hash) []byte
	GetEphemeralExpire(hash crypto.Hash) (bool, uint64)
	AddFeeCollected(uint64)
	Epoch() uint64
}

type HashInstruction struct {
	Instruction Instruction
	Hash        crypto.Hash
}

type Wallet struct {
	Account        crypto.Hash
	FungibleTokens uint64
}

type Payment struct {
	Debit  []Wallet
	Credit []Wallet
}

func GetEpochFromByteArray(inst []byte) uint64 {
	epoch, _ := util.ParseUint64(inst, 2)
	return epoch
}

func NewPayment(debitAcc crypto.Hash, value uint64) *Payment {
	return &Payment{
		Debit:  []Wallet{{debitAcc, value}},
		Credit: []Wallet{},
	}
}

func (p *Payment) NewCredit(account crypto.Hash, value uint64) {
	for _, credit := range p.Credit {
		if credit.Account.Equal(account) {
			credit.FungibleTokens += value
			return
		}
	}
	p.Credit = append(p.Credit, Wallet{Account: account, FungibleTokens: value})
}

func (p *Payment) NewDebit(account crypto.Hash, value uint64) {
	for _, debit := range p.Debit {
		if debit.Account.Equal(account) {
			debit.FungibleTokens += value
			return
		}
	}
	p.Debit = append(p.Debit, Wallet{Account: account, FungibleTokens: value})
}

type Instruction interface {
	Validate(InstructionValidator) bool
	Payments() *Payment
	Serialize() []byte
	Epoch() uint64
	Kind() byte
}

func ParseInstruction(data []byte) Instruction {
	if data[0] != 0 {
		return nil
	}
	switch data[1] {
	case iTransfer:
		return ParseTransfer(data)
	case iDeposit:
		return ParseDeposit(data)
	case iWithdraw:
		return ParseWithdraw(data)
	case iJoinNetwork:
		return ParseJoinNetwork(data)
	case iUpdateInfo:
		return ParseUpdateInfo(data)
	case iCreateAudience:
		return ParseCreateAudience(data)
	case iJoinAudience:
		return ParseJoinAudience(data)
	case iAcceptJoinRequest:
		return ParseAcceptJoinAudience(data)
	case iContent:
		return ParseContent(data)
	case iUpdateAudience:
		return ParseUpdateAudience(data)
	case iGrantPowerOfAttorney:
		return ParseGrantPowerOfAttorney(data)
	case iRevokePowerOfAttorney:
		return ParseRevokePowerOfAttorney(data)
	case iSponsorshipOffer:
		return ParseSponsorshipOffer(data)
	case iSponsorshipAcceptance:
		return ParseSponsorshipAcceptance(data)
	case iCreateEphemeral:
		return ParseCreateEphemeral(data)
	case iSecureChannel:
		return ParseSecureChannel(data)
	case iReact:
		return ParseReact(data)
	}
	return nil
}
