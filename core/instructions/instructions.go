package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/store"
	"github.com/Aereum/aereum/core/util"
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
	SetNewAudience(hash crypto.Hash, stage store.StageKeys) bool
	UpdateAudience(hash crypto.Hash, stage store.StageKeys) bool
	Balance(hash crypto.Hash) uint64
	PowerOfAttorney(hash crypto.Hash) bool
	SponsorshipOffer(hash crypto.Hash) uint64
	HasMember(hash crypto.Hash) bool
	HasCaption(hash crypto.Hash) bool
	HasGrantedSponser(hash crypto.Hash) (bool, crypto.Hash)
	GetAudienceKeys(hash crypto.Hash) *store.StageKeys
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
	JSON() string
}

func ParseInstruction(data []byte) Instruction {
	if data[0] != 0 {
		return nil
	}
	switch data[1] {
	case ITransfer:
		return ParseTransfer(data)
	case IDeposit:
		return ParseDeposit(data)
	case IWithdraw:
		return ParseWithdraw(data)
	case IJoinNetwork:
		return ParseJoinNetwork(data)
	case IUpdateInfo:
		return ParseUpdateInfo(data)
	case ICreateAudience:
		return ParseCreateStage(data)
	case IJoinAudience:
		return ParseJoinStage(data)
	case IAcceptJoinRequest:
		return ParseAcceptJoinStage(data)
	case IContent:
		return ParseContent(data)
	case IUpdateAudience:
		return ParseUpdateStage(data)
	case IGrantPowerOfAttorney:
		return ParseGrantPowerOfAttorney(data)
	case IRevokePowerOfAttorney:
		return ParseRevokePowerOfAttorney(data)
	case ISponsorshipOffer:
		return ParseSponsorshipOffer(data)
	case ISponsorshipAcceptance:
		return ParseSponsorshipAcceptance(data)
	case ICreateEphemeral:
		return ParseCreateEphemeral(data)
	case ISecureChannel:
		return ParseSecureChannel(data)
	case IReact:
		return ParseReact(data)
	}
	return nil
}

func InstructionKind(msg []byte) byte {
	if len(msg) < 2 {
		return iUnkown
	}
	return msg[1]
}
