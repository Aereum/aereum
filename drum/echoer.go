package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

type TextOpenContent struct {
	ActorToken  crypto.Token
	ActorName   string
	ContentType string
	Text        string
}

type AudienceDescription struct {
	Token       crypto.Token
	Description string
	Readers     map[crypto.Token]string
	Submittors  map[crypto.Token]string
	Moderators  map[crypto.Token]string
	Flag        byte
}

type Echoer interface {
	Send(instructions.Instruction)
	Recieve() *TextOpenContent
	Subscribe(crypto.Token)
	Member(crypto.Token) string
	Stage(crypto.Token) *AudienceDescription
}

type Theatre struct {
	channel      chan instructions.Instruction
	subscription map[crypto.Token]struct{}
	members      map[crypto.Token]string
	stages       map[crypto.Token]*AudienceDescription
}

func (t *Theatre) Member(token crypto.Token) string {
	return t.members[token]
}

func (t *Theatre) Send(i instructions.Instruction) {
	if i.Kind() == instructions.IContent {
		t.channel <- i
	}
}

func (t *Theatre) Stage(token crypto.Token) *AudienceDescription {
	return t.stages[token]
}

func (t *Theatre) Recieve() *TextOpenContent {
	for {
		instr := <-t.channel
		if _, ok := t.subscription[instr.Authority()]; ok {
			if instr.Kind() == instructions.IContent {
				if content, ok := instr.(*instructions.Content); ok {
					name := t.members[instr.Authority()]
					return &TextOpenContent{
						ActorToken:  instr.Authority(),
						ActorName:   name,
						ContentType: content.ContentType,
						Text:        string(content.Content),
					}
				}
			}
		}
	}
}

func (t *Theatre) Subscribe(token crypto.Token) {
	t.subscription[token] = struct{}{}
}

func readPlays() Echoer {
	theatre := Theatre{
		channel:      make(chan instructions.Instruction),
		subscription: make(map[crypto.Token]struct{}),
		members:      make(map[crypto.Token]string),
		stages:       make(map[crypto.Token]*AudienceDescription),
	}
	file, err := os.Open("plays.csv")
	if err != nil {
		log.Fatalf("Could not open palys.csv")
		return nil
	}
	playsCSV := csv.NewReader(file)
	playsCSV.Comma = '\t'
	playsCSV.FieldsPerRecord = 3
	plays := make(map[string][][2]string)
	stages := make(map[string]*instructions.Stage)
	characters := make(map[string]*instructions.Author)

	for {
		if line, _ := playsCSV.Read(); line != nil {
			if _, ok := characters[line[1]]; !ok {
				_, secret := crypto.RandomAsymetricKey()
				characters[line[1]] = &instructions.Author{
					PrivateKey: secret,
					Wallet:     secret,
					Attorney:   crypto.ZeroPrivateKey,
				}
			}
			speach := [2]string{line[1], line[2]}
			if play, ok := plays[line[0]]; ok {
				plays[line[0]] = append(play, speach)
			} else {
				plays[line[0]] = [][2]string{speach}
				stages[line[0]] = instructions.NewStage(0, line[0])
			}
		} else {
			break
		}
	}
	go func() {
		countRow := make(map[string]int)
		for playName, _ := range plays {
			countRow[playName] = 0
		}
		for {
			for playName, text := range plays {
				row := countRow[playName]
				if row >= len(text) {
					row = 0
				}
				countRow[playName] = row + 1
				speach := text[row]
				author, authorExists := characters[speach[0]]
				if !authorExists {
					log.Fatalf("could not find character %v", speach[0])
				}
				stage, stageExists := stages[playName]
				if !stageExists {
					log.Fatal("could not find stage")
				}
				if content := author.NewContent(stage, "text", []byte(speach[1]), false, false, 0, 0); content != nil {
					theatre.channel <- content
				}
			}
		}
	}()
	return &theatre
}
