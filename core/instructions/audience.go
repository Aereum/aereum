package instructions

import (
	"github.com/Aereum/aereum/core/crypto"
)

type Audience struct {
	PrivateKey crypto.PrivateKey
	Submission crypto.PrivateKey
	Moderation crypto.PrivateKey
	CipherKey  []byte
	Members    map[crypto.Token]crypto.Token
	Submittors []crypto.Token
	Moderators []crypto.Token
	Readers    []crypto.Token
}

func (a *Audience) SealedToken(key []byte) []byte {
	cipher := crypto.CipherFromKey(key)
	return cipher.Seal(a.PrivateKey[0:32])
}

func (a *Audience) SealedSubmission(key []byte) []byte {
	cipher := crypto.CipherFromKey(key)
	return cipher.Seal(a.Submission[0:32])
}

func (a *Audience) SealedModeration(key []byte) []byte {
	cipher := crypto.CipherFromKey(key)
	return cipher.Seal(a.Moderation[0:32])
}

func NewAudience() *Audience {
	audience := Audience{}
	_, audience.PrivateKey = crypto.RandomAsymetricKey()
	_, audience.Submission = crypto.RandomAsymetricKey()
	_, audience.Moderation = crypto.RandomAsymetricKey()
	audience.CipherKey = crypto.NewCipherKey()
	return &audience
}
