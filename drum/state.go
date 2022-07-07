package main

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
	"github.com/Aereum/aereum/core/instructions"
)

type WalletBalance struct {
	Token   crypto.PrivateKey
	Balance uint64
}

type StageInfo struct {
	Stage   *instructions.Stage
	Content []StageContentInfo
}

type MyState struct {
	Vault          *SecureVault
	MyToken        crypto.Token
	Myself         *instructions.Author
	MySecret       crypto.PrivateKey
	MyAttorneys    map[crypto.Token]struct{}
	Stages         map[crypto.Token]*StageInfo
	Wallets        map[crypto.Token]WalletBalance
	Epoch          uint64
	MyInstructions []instructions.Instruction
	MyHashes       map[crypto.Hash]int
	Validated      []uint64
	instructionsIO *PersistentByteArray
	hashesIO       *PersistentByteArray
}

func (s *MyState) IncorporateGrantAttorney(poa instructions.GrantPowerOfAttorney) {
	s.MyAttorneys[poa.Attorney] = struct{}{}
}

func (s *MyState) IncorporateRevokeAttorney(poa instructions.GrantPowerOfAttorney) {
	delete(s.MyAttorneys, poa.Attorney)
}

func (s *MyState) IncorporateCreateStage(stage instructions.CreateStage) {
	newStage := instructions.Stage{
		Flag:        stage.Flag,
		Description: stage.Description,
	}
	s.Stages[stage.Audience] = &StageInfo{
		Stage: &newStage,
		Content: make([]StageContentInfo,0)
	}
}

func (s *MyState) IncorporateUpdateStage(stage instructions.UpdateStage) {
	// TODO
}

func (s *MyState) IncoporateAcceptMyJoinStageRequest(accept instructions.AcceptJoinStage) {
	var key crypto.PrivateKey
	if accept.DiffHellKey.Equal(s.MyToken) {
		token = s.MySecret
	} else {
		var ok bool
		key, ok = s.Vault.Secrets[accept.DiffHellKey]
		if !ok {
			return
		}
	}
	cipher := dh.ConsensusCipher(key, )
	tc := TokenCipher{Token: token, Cipher: cipher.Seal(a.CipherKey)}
	update.ReadMembers = append(update.ReadMembers, tc)
}

func (s *MyState) IncorporateContent(c instructions.Content) {
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

func (s *MyState) Incorporate(i instructions.Instruction) {
	data := i.Serialize()
	if author := i.Authority(); author.Equal(s.MyToken) {
		hash := crypto.Hasher(data)
		s.hashesIO.Append(hash[:])
	} else {
		s.instructionsIO.Append(data)
		// TODO own instruction from another device
	}
}

func NewMyState(token crypto.PrivateKey, broker *InstructionBroker) *MyState {
	state := &MyState{
		Myself:  &instructions.Author{PrivateKey: token, Wallet: token, Attorney: crypto.ZeroPrivateKey},
		Stages:  make(map[crypto.Token]*StageInfo),
		Wallets: make(map[crypto.Token]WalletBalance),
	}

	go func() {
		for {
			instruction := <-broker.Received
			if instruction.Kind() == instructions.ICreateAudience {
				create, _ := instruction.(*instructions.CreateStage)
				if stage, ok := state.Stages[create.Audience]; ok {

				}
			}
		}
	}()

	return state
}

func (s *MyState) Post(content []byte, contentType string, stage crypto.Token, encrypted, hashed bool) *instructions.Content {
	if stage, ok := s.Stages[stage]; !ok {
		return nil
	} else {
		return s.Myself.NewContent(stage.Stage, contentType, content, hashed, encrypted, s.Epoch, 0)
	}
}

type WalletKeys struct {
	Wallet crypto.Token
	Secret crypto.PrivateKey
}

func (s *MyState) CreateStage(description string, flag byte) *instructions.CreateStage {
	stage := instructions.NewStage(flag, description)
	return s.Myself.NewCreateAudience(stage, flag, description, s.Epoch, 0)
}

func NewState(token crypto.PrivateKey, broker *InstructionBroker) *State {
	state := &MyState{
		Author:    instructions.Author{PrivateKey: token, Wallet: token},
		Wallets:   []crypto.PrivateKey{token},
		Attorneys: []crypto.Token{},
		Stages:    make(map[crypto.Token]*Stage),
		Posts:     make([]*instructions.Content, 0),
		Broker:    broker,
	}
	go func() {
		for {
			instruction := <-broker.Received
			if instruction.Kind() == instructions.ICreateAudience {
				create, _ := instruction.(*instructions.CreateStage)
				if stage, ok := state.Stages[create.Audience]; ok {
					stage.live = true
				}
			}
		}
	}()
	return state
}
