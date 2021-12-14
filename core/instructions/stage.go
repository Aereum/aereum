package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
)

type Stage struct {
	PrivateKey  crypto.PrivateKey
	Submission  crypto.PrivateKey
	Moderation  crypto.PrivateKey
	CipherKey   []byte
	Submittors  map[crypto.Token]crypto.Token // token -> DiffieHellmanKey
	Moderators  map[crypto.Token]crypto.Token // token -> DiffieHellmanKey
	Readers     map[crypto.Token]crypto.Token // token -> DiffieHellmanKey
	Flag        byte
	Description string
}

func (a *Stage) ResetKeys() *UpdateStage {
	_, a.Submission = crypto.RandomAsymetricKey()
	_, a.Moderation = crypto.RandomAsymetricKey()
	a.CipherKey = crypto.NewCipherKey()
	diffHellPub, diffHellPrv := dh.NewEphemeralKey()
	update := UpdateStage{
		stage:       a.PrivateKey.PublicKey(),
		submission:  a.Submission.PublicKey(),
		moderation:  a.Moderation.PublicKey(),
		flag:        a.Flag,
		diffHellKey: diffHellPub,
		description: a.Description,
		readMembers: make(TokenCiphers, 0),
		subMembers:  make(TokenCiphers, 0),
		modMembers:  make(TokenCiphers, 0),
	}
	for token, key := range a.Readers {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{token: token, cipher: cipher.Seal(a.CipherKey)}
		update.readMembers = append(update.readMembers, tc)
	}
	for token, key := range a.Moderators {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{token: token, cipher: cipher.Seal(a.Moderation[:32])}
		update.modMembers = append(update.modMembers, tc)
	}
	for token, key := range a.Submittors {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{token: token, cipher: cipher.Seal(a.Submission[:32])}
		update.subMembers = append(update.subMembers, tc)
	}
	return &update
}

func NewStage(flag byte, description string) *Stage {
	stage := Stage{Flag: flag, Description: description}
	_, stage.PrivateKey = crypto.RandomAsymetricKey()
	_, stage.Submission = crypto.RandomAsymetricKey()
	_, stage.Moderation = crypto.RandomAsymetricKey()
	stage.CipherKey = crypto.NewCipherKey()
	stage.Submittors = make(map[crypto.Token]crypto.Token)
	stage.Moderators = make(map[crypto.Token]crypto.Token)
	stage.Readers = make(map[crypto.Token]crypto.Token)
	return &stage
}
