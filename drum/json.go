package main

import (
	"encoding/hex"
	"encoding/json"

	"github.com/Aereum/aereum/core/crypto"
)

type jsoner interface {
	JSON() []byte
}

type AttorneysJSON struct {
	Info      string   `JSON:"info"`
	Attorneys []string `JSON:"attorneys"`
}

func TokenToHex(token crypto.Token) string {
	return hex.EncodeToString(token[:])
}

func NewAttorneysJSON(attorneys map[crypto.Token]struct{}) AttorneysJSON {
	JSON := AttorneysJSON{
		Info:      "Attorneys",
		Attorneys: make([]string, 0),
	}
	for attorney := range attorneys {
		JSON.Attorneys = append(JSON.Attorneys, attorney.Hex())
	}
	return JSON
}

func (a AttorneysJSON) JSON() []byte {
	encoded, _ := json.Marshal(a)
	return encoded
}

type StageJSON struct {
	Info         string `JSON:"info"`
	Token        string
	Description  string
	Flag         byte
	Readers      []string
	Submitters   []string
	Moderators   []string
	Creator      string
	MessageCount int
	Live         bool
}

func TokenMapToList(tokenMap map[crypto.Token]crypto.Token) []string {
	tokens := make([]string, 0)
	for token := range tokenMap {
		tokens = append(tokens, TokenToHex(token))
	}
	return tokens
}

func NewStageJSON(s *StageInfo) StageJSON {
	return StageJSON{
		Info:         "Stage",
		Token:        TokenToHex(s.Stage.PrivateKey.PublicKey()),
		Description:  s.Stage.Description,
		Flag:         s.Stage.Flag,
		Readers:      TokenMapToList(s.Stage.Readers),
		Submitters:   TokenMapToList(s.Stage.Submittors),
		Moderators:   TokenMapToList(s.Stage.Moderators),
		Creator:      TokenToHex(s.Creator),
		MessageCount: len(s.Content),
		Live:         s.Live,
	}
}
func (s StageJSON) JSON() []byte {
	encoded, _ := json.Marshal(s)
	return encoded
}

type WalletJSON struct {
	Info    string `JSON:"info"`
	Token   string `JSON:"token"`
	Balance int    `JSON:"balance"`
}

func NewWalletJSON(wallet crypto.Token, balance uint64) WalletJSON {
	return WalletJSON{
		Info:    "wallet",
		Token:   wallet.Hex(),
		Balance: int(balance),
	}
}

func (a WalletJSON) JSON() []byte {
	encoded, _ := json.Marshal(a)
	return encoded
}

type AttorneyJSON struct {
	Info  string `JSON:"info"`
	Token string `JSON:"token"`
}

func NewAttorneyGrant(token crypto.Token) AttorneyJSON {
	return AttorneyJSON{
		Info:  "Grant Attorney",
		Token: token.Hex(),
	}
}

func NewAttorneyRevoke(token crypto.Token) AttorneyJSON {
	return AttorneyJSON{
		Info:  "Revoke Attorney",
		Token: token.Hex(),
	}
}

func (a AttorneyJSON) JSON() []byte {
	encoded, _ := json.Marshal(a)
	return encoded
}
