package main

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
	"github.com/Aereum/aereum/core/instructions"
)

func (s *MyState) Incorporate(i instructions.Instruction) {
	// if instructions is from myself, mark it as incorporated to the blockchain
	data := i.Serialize()
	if author := i.Authority(); author.Equal(s.MyToken) {
		hash := crypto.Hasher(data)
		s.hashesIO.Append(hash[:])
		delete(s.SentInstructions, hash)
	}
	// store it on the instruction cache
	s.instructionsIO.Append(data)
	// update wallets with payments if relevant
	if payments := i.Payments(); payments != nil {
		s.IncoporatePayments(payments)
	}
	// perform appropriate actions accodingly to its type
	switch i.Kind() {

	case instructions.ITransfer:
		// already affected in Incorporte Payments
		return
	case instructions.IDeposit:
		if deposit, ok := i.(*instructions.Deposit); ok {
			s.IncorporateDeposit(deposit)
		}
		return
	case instructions.IWithdraw:
		if withdraw, ok := i.(*instructions.Withdraw); ok {
			s.IncorportateWithdraw(withdraw)
		}
		return
	case instructions.IJoinNetwork:
		if join, ok := i.(*instructions.JoinNetwork); ok {
			s.IncorporateJoinNetwork(join)
		}
		return
	case instructions.IUpdateInfo:
		if update, ok := i.(*instructions.UpdateInfo); ok {
			s.IncorporateUpdateInfo(update)
		}
		return
	case instructions.ICreateAudience:
		if createStage, ok := i.(*instructions.CreateStage); ok {
			s.IncorporateCreateStage(createStage)
		}
		return
	case instructions.IJoinAudience:
		// we only keep track of accepted requests
		return
	case instructions.IAcceptJoinRequest:
		if accept, ok := i.(*instructions.AcceptJoinStage); ok {
			if accept.Member.Equal(s.MyToken) {
				s.IncoporateAcceptMyJoinStageRequest(accept)
			} else {
				s.IncorporateAcceptOtherJoinRequest(accept)
			}
		}
		return
	case instructions.IContent:
		if content, ok := i.(*instructions.Content); ok {
			s.IncorporateContent(content)
		}
		return

	case instructions.IUpdateAudience:
		if update, ok := i.(*instructions.UpdateStage); ok {
			s.IncorporateUpdateStage(update)
		}
		return

	case instructions.IGrantPowerOfAttorney:
		if grant, ok := i.(*instructions.GrantPowerOfAttorney); ok {
			s.IncorporateGrantAttorney(grant)
		}
		return
	case instructions.IRevokePowerOfAttorney:
		if revoke, ok := i.(*instructions.RevokePowerOfAttorney); ok {
			s.IncorporateRevokeAttorney(revoke)
		}
		return

	case instructions.ISponsorshipOffer:
		// do it later
		return
	case instructions.ISponsorshipAcceptance:
		// do it later
		return

	case instructions.ICreateEphemeral:
		// do it later
		return

	case instructions.ISecureChannel:
		// do it later
		return

	case instructions.IReact:
		// do it later
		return

	}

}

func (s *MyState) IncorporateCreateStage(stage *instructions.CreateStage) {
	var newStageInfo *StageInfo
	if stage.Authored.Author.Equal(s.MyToken) {
		newStageInfo = s.Stages[stage.Audience]
		if newStageInfo == nil {
			return
		}
		newStageInfo.Live = true
	} else {
		newStageInfo = &StageInfo{
			Stage: &instructions.Stage{
				CipherKey:   make([]byte, 0),
				Readers:     make(map[crypto.Token]crypto.Token),
				Submittors:  make(map[crypto.Token]crypto.Token),
				Moderators:  make(map[crypto.Token]crypto.Token),
				Flag:        stage.Flag,
				Description: stage.Description,
			},
			Content: make([]StageContentInfo, 0),
			Creator: stage.Authored.Author,
			Live:    true,
		}
		s.Stages[stage.Audience] = newStageInfo
	}
	s.queue.Send(NewStageJSON(newStageInfo))
}

func (s *MyState) IncorporateUpdateStage(stage *instructions.UpdateStage) {
	mystage, ok := s.Stages[stage.Stage]
	if !ok {
		return
	}
	// invalidate existing moderation and submission keys if not compatible with
	// respective tokens anymore.
	if !mystage.Stage.Moderation.PublicKey().Equal(stage.Moderation) {
		mystage.Stage.Moderation = crypto.ZeroPrivateKey
	}
	if !mystage.Stage.Submission.PublicKey().Equal(stage.Submission) {
		mystage.Stage.Submission = crypto.ZeroPrivateKey
	}
	// search for a known token in secure vault for moderation
	for _, tokenCipher := range stage.ModMembers {
		if prv, ok := s.Vault.Secrets[tokenCipher.Token]; ok && prv != crypto.ZeroPrivateKey {
			cipher := dh.ConsensusCipher(prv, stage.DiffHellKey)
			var prvKey crypto.PrivateKey
			modKey, _ := cipher.Open(tokenCipher.Cipher)
			copy(prvKey[:], modKey)
			if prvKey.PublicKey().Equal(stage.Moderation) {
				mystage.Stage.Moderation = prvKey
			}
			break
		}
	}
	// search for a known token in secure vault for submission
	for _, tokenCipher := range stage.SubMembers {
		if prv, ok := s.Vault.Secrets[tokenCipher.Token]; ok && prv != crypto.ZeroPrivateKey {
			cipher := dh.ConsensusCipher(prv, stage.DiffHellKey)
			var prvKey crypto.PrivateKey
			subKey, _ := cipher.Open(tokenCipher.Cipher)
			copy(prvKey[:], subKey)
			if prvKey.PublicKey().Equal(stage.Submission) {
				mystage.Stage.Submission = prvKey
			}
			break
		}
	}
	// search for a known token in secure vault for reading cipher
	for _, tokenCipher := range stage.ReadMembers {
		if prv, ok := s.Vault.Secrets[tokenCipher.Token]; ok && prv != crypto.ZeroPrivateKey {
			cipher := dh.ConsensusCipher(prv, stage.DiffHellKey)
			cipherKey, _ := cipher.Open(tokenCipher.Cipher)
			mystage.Stage.CipherKey = cipherKey
			break
		}
	}
}

func (s *MyState) IncorporateAcceptJoinStageRequest(request *instructions.AcceptJoinStage) {
	// TODO: Decide is Mine or Other
}

func (s *MyState) IncoporateAcceptMyJoinStageRequest(accept *instructions.AcceptJoinStage) {
	stage, ok := s.Stages[accept.Stage]
	if !ok {
		return
	}
	stage.Stage.JoinRequestAccepted(accept)
	s.queue.Send(NewStageJSON(stage))
}

func (s *MyState) IncorporateAcceptOtherJoinRequest(accept *instructions.AcceptJoinStage) {
	stage, ok := s.Stages[accept.Stage]
	if !ok {
		return
	}
	if accept.Moderate != nil {
		stage.Moderators = append(stage.Moderators, accept.Member)
	}
	if accept.Submit != nil {
		stage.Submitters = append(stage.Submitters, accept.Member)
	}
	if accept.Read != nil {
		stage.Readers = append(stage.Readers, accept.Member)
	}
}

func (s *MyState) IncorporateContent(c *instructions.Content) {
	stage, ok := s.Stages[c.Audience]
	if !ok {
		return
	}
	content := StageContentInfo{
		Epoch:       c.Epoch(),
		ContentType: c.ContentType,
		Moderated:   !c.Moderator.Equal(crypto.ZeroToken),
		Sponsored:   c.Sponsored,
	}
	if c.Encrypted {
		cipher := crypto.CipherFromKey(stage.Stage.CipherKey)
		content.Content, _ = cipher.Open(c.Content)
	}
	stage.Content = append(stage.Content, content)
}

func (s *MyState) IncoporatePayments(pay *instructions.Payment) {
	for _, credit := range pay.Credit {
		if wallet, ok := s.Wallets[credit.Account]; ok {
			wallet.Balance = wallet.Balance + credit.FungibleTokens
			s.queue.Send(NewWalletJSON(wallet.Token.PublicKey(), wallet.Balance, 0))
		}
	}
	for _, debit := range pay.Debit {
		if wallet, ok := s.Wallets[debit.Account]; ok {
			wallet.Balance = wallet.Balance - debit.FungibleTokens
			s.queue.Send(NewWalletJSON(wallet.Token.PublicKey(), wallet.Balance, 0))
		}
	}
}

func (s *MyState) IncorporateGrantAttorney(poa *instructions.GrantPowerOfAttorney) {
	s.MyAttorneys[poa.Attorney] = struct{}{}
	s.queue.Send(NewAttorneysJSON(s.MyAttorneys))
}

func (s *MyState) IncorporateRevokeAttorney(poa *instructions.RevokePowerOfAttorney) {
	delete(s.MyAttorneys, poa.Attorney)
	s.queue.Send(NewAttorneysJSON(s.MyAttorneys))
}

func (m *MyState) IncorporateTransfer(transfer *instructions.Transfer) {
	total := uint64(0)
	for _, to := range transfer.To {
		if wallet, ok := m.Wallets[crypto.HashToken(to.Token)]; ok {
			wallet.Balance += to.Value
		}
		total += to.Value
	}
	if wallet, ok := m.Wallets[crypto.HashToken(transfer.From)]; ok {
		wallet.Balance -= total + transfer.Fee
	}
	// send wallets
}

func (m *MyState) IncorporateDeposit(deposit *instructions.Deposit) {
	if wallet, ok := m.Wallets[crypto.HashToken(deposit.Token)]; ok {
		wallet.Balance -= deposit.Value
		wallet.Stake += deposit.Value
	}
	// send wallets
}

func (m *MyState) IncorportateWithdraw(withdraw *instructions.Withdraw) {
	if wallet, ok := m.Wallets[crypto.HashToken(withdraw.Token)]; ok {
		wallet.Balance += withdraw.Value
		wallet.Stake -= withdraw.Value
	}
	// send wallets
}

func (m *MyState) IncorporateJoinNetwork(join *instructions.JoinNetwork) {
	detail := &MemberDetails{
		Token:  join.Authored.Author,
		Handle: join.Caption,
	}
	m.Members[crypto.HashToken(join.Authored.Author)] = detail
	m.Members[crypto.Hasher([]byte(join.Caption))] = detail
	// TODO
}

func (m *MyState) IncorporateUpdateInfo(update *instructions.UpdateInfo) {
	if update.Authored.Author.Equal(m.MyToken) {
		return
	}
	if member, ok := m.Members[crypto.HashToken(update.Authority())]; ok {
		oldHandle := member.Handle
		member.Handle = update.Details
		if member.Handle != update.Details {
			delete(m.Members, crypto.Hasher([]byte(oldHandle)))
			m.Members[crypto.Hasher([]byte(update.Details))] = member
		}
	}
	// TODO LOGIC
}

// COMPLETE THESE LATER

func (m *MyState) IncorporateSponsorshipOffer(offer *instructions.SponsorshipOffer) {
	// TODO
}

func (m *MyState) IncorporateSponsorshipAcceptance(accept *instructions.SponsorshipAcceptance) {
	// TODO
}

func (m *MyState) IncorporateCreateEphemeral(ephemeral *instructions.CreateEphemeral) {
	// TODO
}

func (m *MyState) IncorporateSecureChannel(secure *instructions.SecureChannel) {
	// TODO
}

func (m *MyState) IncorporateReact(react *instructions.React) {
	// TODO
}
