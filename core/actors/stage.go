package actors

import (
	"github.com/Aereum/aereum/core/crypto"
)

type Stage struct {
	PrivateKey  crypto.PrivateKey
	Submission  crypto.PrivateKey
	Moderation  crypto.PrivateKey
	CipherKey   []byte
	Submittors  []Author
	Moderators  []Author
	Readers     []Author
	Flag        byte
	Description string
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
