package main

import (
	"fmt"
	"log"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
	"github.com/Aereum/aereum/core/instructions"
)

type WalletBalance struct {
	Token   crypto.PrivateKey
	Balance uint64
}

type StageInfo struct {
	Stage      *instructions.Stage
	Content    []StageContentInfo
	Creator    crypto.Token
	Moderators []crypto.Token
	Submitters []crypto.Token
	Readers    []crypto.Token
	Live       bool
}

type MemberDetails struct {
	Token  crypto.Token
	Handle string
	Media  string
}

type StageContentInfo struct {
	Author      string
	Epoch       uint64
	ContentType string
	Content     []byte
	Moderated   bool
	Sponsored   bool
}

// Members is sent to front end to decide best presentation of other members

type MyState struct {
	Vault            *SecureVault
	MyToken          crypto.Token
	Myself           *instructions.Author
	MySecret         crypto.PrivateKey
	Members          map[crypto.Hash]*MemberDetails // hash of handle or token
	MyAttorneys      map[crypto.Token]struct{}
	Stages           map[crypto.Token]*StageInfo
	Wallets          map[crypto.Hash]*WalletBalance
	Stakes           map[crypto.Token]uint64 // wallet token to staked balance
	Epoch            uint64
	SentInstructions map[crypto.Hash]instructions.Instruction
	MyStakes         uint64
	Validated        []uint64
	instructionsIO   *PersistentByteArray
	hashesIO         *PersistentByteArray
	queue            *SocketQueue
	broker           *InstructionQueue
}

// NewMyState returns a pointer to a MyState instance linked to provided
// broker and queue channels to communicate with the aereum network and
// the front end respectively. The token is used to read the instruction
// cache files and associated with the password to read the secure vault
// where secret keys are securely stored.
func NewMyState(token crypto.Token, password string, broker *InstructionQueue, queue *SocketQueue) *MyState {
	// open files
	vaulFileName := fmt.Sprintf("%v_vault.dat", TokenToHex(token))
	vault := OpenVaultFromPassword([]byte(password), vaulFileName)
	if vault == nil {
		log.Fatalf("Could not open or create vault file.")
	}
	instructionsFileName := fmt.Sprintf("%v_instructions_cache.dat", TokenToHex(token))
	instructionsDB := OpenParsistentByteArray(instructionsFileName)
	if instructionsDB == nil {
		log.Fatalf("Could not open or create instructions cache file.")
	}
	hashesFileName := fmt.Sprintf("%v_hashes_cache.dat", TokenToHex(token))
	hashesDB := OpenParsistentByteArray(hashesFileName)
	if hashesDB == nil {
		log.Fatalf("Could not open or create hash confirmation file.")
	}
	if vault == nil || instructionsDB == nil || hashesDB == nil {
		log.Fatalf("Could not open or create necessary files.")
	}
	// set state to its most persistent version
	state := &MyState{
		Vault:            vault,
		MyToken:          token,
		Myself:           &instructions.Author{PrivateKey: vault.SecretKey, Wallet: vault.SecretKey, Attorney: crypto.ZeroPrivateKey},
		MyAttorneys:      make(map[crypto.Token]struct{}),
		Stages:           make(map[crypto.Token]*StageInfo),
		Wallets:          make(map[crypto.Hash]*WalletBalance),
		Stakes:           make(map[crypto.Token]uint64),
		Epoch:            0,
		SentInstructions: make(map[crypto.Hash]instructions.Instruction),
		instructionsIO:   instructionsDB,
		hashesIO:         hashesDB,
		queue:            queue,
		broker:           broker,
	}
	hashes := make(map[crypto.Hash]struct{})
	for {
		hash := hashesDB.Read()
		if hash == nil {
			break
		}
		var newHash crypto.Hash
		copy(newHash[:], hash)
		hashes[newHash] = struct{}{}
	}
	for {
		bytes := instructionsDB.Read()
		if bytes == nil {
			break
		}
		instruction := instructions.ParseInstruction(bytes)
		if instruction.Authority().Equal(token) {
			hash := crypto.Hasher(bytes)
			if _, ok := hashes[hash]; !ok {
				state.SentInstructions[hash] = instruction
			}
		}
		state.Incorporate(instruction)
	}
	// send stage, wallets and attorneys to front end
	for _, wallet := range state.Wallets {
		state.queue.Send(NewWalletJSON(wallet.Token.PublicKey(), wallet.Balance))

	}
	for _, stageInfo := range state.Stages {
		state.queue.Send(NewStageJSON(stageInfo))
	}
	for attorney := range state.MyAttorneys {
		state.queue.Send(NewAttorneyGrant(attorney))
	}
	return state
}

/* STAGE ACTIONS */

func (s *MyState) CreateStage(flag byte, description string) {
	newStage := instructions.NewStage(flag, description)
	stageInfo := &StageInfo{
		Stage:   newStage,
		Content: make([]StageContentInfo, 0),
		Live:    false,
		Creator: s.MyToken,
	}
	if instruction := s.Myself.NewCreateAudience(newStage, s.Epoch, 0); instruction != nil {
		s.Stages[newStage.PrivateKey.PublicKey()] = stageInfo
		s.broker.Send(instruction)
	}
	s.queue.Send(NewStageJSON(stageInfo))
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

func removeFromMap(original map[crypto.Token]crypto.Token, toRemove []crypto.Token) map[crypto.Token]crypto.Token {
	newMap := make(map[crypto.Token]crypto.Token)
	if len(toRemove) == 0 {
		for key, value := range original {
			newMap[key] = value
		}
	} else {
		for key, value := range original {
			mantain := true
			for _, exclude := range toRemove {
				if key.Equal(exclude) {
					mantain = false
					break
				}
			}
			if mantain {
				newMap[key] = value
			}
		}
	}
	return newMap
}

func (s *MyState) UpdateStage(token crypto.Token, flag byte, description string, exReading, exWritting, exModerating []crypto.Token) {
	stage, ok := s.Stages[token]
	if !ok {
		return
	}
	readers := removeFromMap(stage.Stage.Readers, exReading)
	submittors := removeFromMap(stage.Stage.Submittors, exWritting)
	moderators := removeFromMap(stage.Stage.Moderators, exModerating)
	s.Myself.NewUpdateAudience(stage.Stage, readers, submittors, moderators, flag, description, s.Epoch, 0)
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

func (s *MyState) AcceptJoinRequest(req *instructions.JoinStage, level byte) {
	stage, ok := s.Stages[req.Audience]
	if !ok {
		return
	}
	accept := stage.Stage.AcceptJoinRequest(req, level, s.Myself, s.Epoch, 0)
	if accept != nil {
		return
	}
	s.broker.Send(accept)
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

func (s *MyState) RequireJoinRequest(token crypto.Token, presentation string) {
	stage, ok := s.Stages[token]
	if !ok {
		return
	}
	prvKey, pubKey := dh.NewEphemeralKey()
	s.Vault.Secrets[pubKey] = prvKey
	stage.Stage.PrivateKey = prvKey
	request := s.Myself.NewJoinAudience(token, pubKey, presentation, s.Epoch, 0)
	if request == nil {
		return
	}
	s.broker.Send(request)
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
			s.queue.Send(NewWalletJSON(wallet.Token.PublicKey(), wallet.Balance))
		}
	}
}

func (s *MyState) IncorporateUpdateInfo(update *instructions.UpdateInfo) {
	if update != nil {
		return
	}
	// TODO LOGIC
}

func (s *MyState) CreateNewWallet() {
	walletSecret := s.Vault.NewKey()
	token := walletSecret.PublicKey()
	s.Wallets[crypto.HashToken(token)] = &WalletBalance{Token: walletSecret, Balance: 0}
	s.queue.Send(NewWalletJSON(token, 0))
}

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
			if stake, ok := s.Stakes[deposit.Token]; ok {
				s.Stakes[deposit.Token] = stake + deposit.Value
			} else {
				s.Stakes[deposit.Token] = deposit.Value
			}
		}
		return
	case instructions.IWithdraw:
		if withdraw, ok := i.(*instructions.Withdraw); ok {
			if stake, ok := s.Stakes[withdraw.Token]; ok {
				if withdraw.Value == stake {
					delete(s.Stakes, withdraw.Token)
				} else {
					s.Stakes[withdraw.Token] = stake - withdraw.Value
				}
			}
		}
		return
	case instructions.IJoinNetwork:
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

	case instructions.ISponsorshipAcceptance:

	case instructions.ICreateEphemeral:

	case instructions.ISecureChannel:

	case instructions.IReact:

	}

}

func (s *MyState) GrantAttorney(attorney crypto.Token) {
	grant := s.Myself.NewGrantPowerOfAttorney(attorney, s.Epoch, 0)
	if grant == nil {
		return
	}
	s.broker.Send(grant)
}

func (s *MyState) RevokeAttorney(attorney crypto.Token) {
	revoke := s.Myself.NewRevokePowerOfAttorney(attorney, s.Epoch, 0)
	if revoke == nil {
		return
	}
	s.broker.Send(revoke)
}

func (s *MyState) IncorporateGrantAttorney(poa *instructions.GrantPowerOfAttorney) {
	s.MyAttorneys[poa.Attorney] = struct{}{}
	s.queue.Send(NewAttorneysJSON(s.MyAttorneys))
}

func (s *MyState) IncorporateRevokeAttorney(poa *instructions.RevokePowerOfAttorney) {
	delete(s.MyAttorneys, poa.Attorney)
	s.queue.Send(NewAttorneysJSON(s.MyAttorneys))
}

func (s *MyState) Post(content []byte, contentType string, stage crypto.Token, encrypted, hashed bool) *instructions.Content {
	if stage, ok := s.Stages[stage]; !ok {
		return nil
	} else {
		return s.Myself.NewContent(stage.Stage, contentType, content, hashed, encrypted, s.Epoch, 0)
	}
}

/*

Regras:

	Create New Stage -> stage criado com live = false
		-> JSON confimando o novo stage non-live
	Incorporate New Stage -> live = true
		-> JSON confirmando o novo stage live
		-> JSON
	Accept Join Request
		-> Incorporate
	Audience Join Request
		-> incorporate stage without details in Stages
	Incorporate Join Request Accepted
		-> save deciphered keys on stage

*/
