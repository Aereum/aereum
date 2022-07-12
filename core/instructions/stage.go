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
	Cipher      crypto.Cipher
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
	diffHellPrv, diffHellPub := dh.NewEphemeralKey()
	update := UpdateStage{
		Stage:       a.PrivateKey.PublicKey(),
		Submission:  a.Submission.PublicKey(),
		Moderation:  a.Moderation.PublicKey(),
		Flag:        a.Flag,
		DiffHellKey: diffHellPub,
		Description: a.Description,
		ReadMembers: make(TokenCiphers, 0),
		SubMembers:  make(TokenCiphers, 0),
		ModMembers:  make(TokenCiphers, 0),
	}
	for token, key := range a.Readers {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{Token: token, Cipher: cipher.Seal(a.CipherKey)}
		update.ReadMembers = append(update.ReadMembers, tc)
	}
	for token, key := range a.Moderators {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{Token: token, Cipher: cipher.Seal(a.Moderation[:32])}
		update.ModMembers = append(update.ModMembers, tc)
	}
	for token, key := range a.Submittors {
		cipher := dh.ConsensusCipher(diffHellPrv, key)
		tc := TokenCipher{Token: token, Cipher: cipher.Seal(a.Submission[:32])}
		update.SubMembers = append(update.SubMembers, tc)
	}
	return &update
}

func (a *Stage) JoinRequestAccepted(accept *AcceptJoinStage) error {
	cipher := dh.ConsensusCipher(a.PrivateKey, accept.DiffHellKey)
	if accept.Read != nil {
		var err error
		if a.CipherKey, err = cipher.Open(accept.Read); err != nil {
			a.Cipher = crypto.CipherFromKey(a.CipherKey)
		} else {
			return err
		}
	}
	if accept.Submit != nil {
		if key, err := cipher.Open(accept.Submit); err != nil {
			copy(a.Submission[:], key)
		} else {
			return err
		}
	}
	if accept.Moderate != nil {
		if key, err := cipher.Open(accept.Moderate); err != nil {
			copy(a.Moderation[:], key)
		} else {
			return err
		}
	}
	a.PrivateKey = crypto.ZeroPrivateKey
	return nil
}

func (a *Stage) AcceptJoinRequest(req *JoinStage, level byte, author *Author, epoch, fee uint64) *AcceptJoinStage {
	accept := AcceptJoinStage{
		Authored: author.NewAuthored(epoch, fee),
		Stage:    a.PrivateKey.PublicKey(),
		Member:   req.Authored.Author,
		Read:     []byte{},
		Submit:   []byte{},
		Moderate: []byte{},
	}
	prv, pub := dh.NewEphemeralKey()
	cipher := dh.ConsensusCipher(prv, req.DiffHellKey)
	accept.Read = cipher.Seal(a.CipherKey)
	if level > 0 {
		accept.Submit = cipher.Seal(a.Submission[:32])
	}
	if level > 1 {
		accept.Moderate = cipher.Seal(a.Moderation[:32])
	}
	accept.DiffHellKey = pub
	modbulk := accept.serializeModBulk()
	accept.modSignature = a.Moderation.Sign(modbulk)
	bulk := accept.serializeBulk()
	if author.sign(accept.Authored, bulk, IAcceptJoinRequest) {
		return &accept
	}
	return nil
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

func NewJoinStage() {}
