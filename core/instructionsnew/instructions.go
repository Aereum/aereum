package instructionsnew

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

type Wallet struct {
	Account        crypto.Hash
	FungibleTokens uint64
}

type Payment struct {
	Debit  []Wallet
	Credit []Wallet
}
