package main

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type WalletBalance struct {
	Token   crypto.Token
	Balance uint64
}

type State struct {
	Author    instructions.Author
	Wallets   []crypto.PrivateKey
	Attorneys []crypto.Token
	Stages    map[crypto.Token]*Stage
	Posts     []*instructions.Content
	Broker    *InstructionBroker
}

type Stage struct {
	stage *instructions.Stage
	live  bool
}

func NewState(token crypto.PrivateKey, broker *InstructionBroker) *State {
	state := &State{
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
			} else if instruction.Kind() ==
		}
	}()
	return state
}

func (s *State) CreateStage(description string, flag byte) {
	stage := instructions.NewStage(flag, description)
	instruction := s.Author.NewCreateAudience(stage, flag, description, s.Broker.Epoch, 1)
	s.Broker.Send <- instruction
}

func (s *State) Publish(stage crypto.Token, message string, contentType string, encrypted bool) {
	if stage, ok := s.Stages[stage]; ok {
		msg := []byte(message)
		if encrypted {
			// TODO
		}
		instruction := s.Author.NewContent(stage, contentType, msg, false, false, s.Broker.Epoch, 1)
		s.Broker.Send <- instruction
	}
}
