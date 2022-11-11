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
	Stake   uint64
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

// StageContentInfo should store relevant unencrypted information on content.
type StageContentInfo struct {
	Author      string
	Epoch       uint64
	ContentType string
	Content     []byte
	Moderated   bool
	Sponsored   bool
}

// Members is sent to front end to decide best presentation of other members

// MyState centralizes all the information relevant to provide an user interface
// to an specific aereum member.
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

func (s *MyState) CreateNewWallet() {
	walletSecret := s.Vault.NewKey()
	token := walletSecret.PublicKey()
	s.Wallets[crypto.HashToken(token)] = &WalletBalance{Token: walletSecret, Balance: 0}
	s.queue.Send(NewWalletJSON(token, 0))
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

func (s *MyState) Post(content []byte, contentType string, stage crypto.Token, encrypted, hashed bool) *instructions.Content {
	if stage, ok := s.Stages[stage]; !ok {
		return nil
	} else {
		return s.Myself.NewContent(stage.Stage, contentType, content, hashed, encrypted, s.Epoch, 0)
	}
}

/*


MyState
	--> Generate Instruction
	--> Incorporate Instruction




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
