package instructions

import "github.com/Aereum/aereum/core/crypto"

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
	Validate(*Block) bool
	Payments() *Payment
	Serialize() []byte
	Epoch() uint64
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
