package main

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

func getTokens(msg []byte) []crypto.Token {
	switch instructions.InstructionKind(msg) {
	case instructions.IContent:
		if content := instructions.ParseContent(msg); content != nil {
			return []crypto.Token{content.Author, content.Attorney, content.Wallet, content.Audience}
		}
	case instructions.ITransfer:
		if transfer := instructions.ParseTransfer(msg); transfer != nil {
			tokens := []crypto.Token{transfer.From}
			for _, reciepient := range transfer.To {
				tokens = append(tokens, reciepient.Token)
			}
			return tokens
		}
	case instructions.IDeposit:
		if deposit := instructions.ParseDeposit(msg); deposit != nil {
			return []crypto.Token{deposit.Token}
		}
	case instructions.IWithdraw:
		if withdraw := instructions.ParseWithdraw(msg); withdraw != nil {
			return []crypto.Token{withdraw.Token}
		}
	case instructions.IJoinNetwork:
		if join := instructions.ParseJoinNetwork(msg); join != nil {
			return authoredTokens(join.Authored)
		}
	case instructions.IUpdateInfo:
		if update := instructions.ParseUpdateInfo(msg); update != nil {
			return authoredTokens(update.Authored)
		}
	case instructions.ICreateAudience:
		if join := instructions.ParseCreateStage(msg); join != nil {
			return append(authoredTokens(join.Authored), join.Audience)
		}
	case instructions.IJoinAudience:
		if join := instructions.ParseJoinStage(msg); join != nil {
			return append(authoredTokens(join.Authored), join.Audience)
		}
	case instructions.IAcceptJoinRequest:
		if join := instructions.ParseAcceptJoinStage(msg); join != nil {
			return append(authoredTokens(join.Authored), join.Stage)
		}
	case instructions.IUpdateAudience:
		// TODO
	case instructions.IGrantPowerOfAttorney:
		if grant := instructions.ParseGrantPowerOfAttorney(msg); grant != nil {
			return append(authoredTokens(grant.Authored), grant.Attorney)
		}
	case instructions.IRevokePowerOfAttorney:
		if revoke := instructions.ParseRevokePowerOfAttorney(msg); revoke != nil {
			return append(authoredTokens(revoke.Authored), revoke.Attorney)
		}
	case instructions.ISponsorshipOffer:
		if offer := instructions.ParseSponsorshipOffer(msg); offer != nil {
			return append(authoredTokens(offer.Authored), offer.Stage)
		}
	case instructions.ISponsorshipAcceptance:
		if accept := instructions.ParseSponsorshipAcceptance(msg); accept != nil {
			return append(authoredTokens(accept.Authored), accept.Stage, accept.Offer.Authored.Author)
		}
	case instructions.ICreateEphemeral:
		if ephemeral := instructions.ParseCreateEphemeral(msg); ephemeral != nil {
			return append(authoredTokens(ephemeral.Authored), ephemeral.EphemeralToken)
		}
	case instructions.ISecureChannel:
		if secure := instructions.ParseSecureChannel(msg); secure != nil {
			// TODO: token range
			return authoredTokens(secure.Authored)
		}
	case instructions.IReact:
		if react := instructions.ParseReact(msg); react != nil {
			// TODO: token range
			return authoredTokens(react.Authored)
		}
	}
	return nil
}
