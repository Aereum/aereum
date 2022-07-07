package main

import (
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
	"github.com/Aereum/aereum/core/instructions"
)

type Participant struct {
	Token         crypto.Token
	DiffieHellman crypto.Token
}

type StageContentInfo struct {
	Author      string
	Epoch       uint64
	ContentType string
	Content     []byte
	Moderated   bool
	Sponsored   bool
}

type Stage struct {
	Token                crypto.Token
	Secret               crypto.PrivateKey
	Moderate             crypto.PrivateKey
	Submit               crypto.PrivateKey
	DiffieHellman        crypto.Token
	DiffieHellmanSecret  crypto.Token
	CipherKey            []byte
	Cipher               crypto.Cipher
	Moderators           map[crypto.Token]struct{}
	Submittors           map[crypto.Token]struct{}
	Readers              map[crypto.Token]struct{}
	MembersDiffieHellman map[crypto.Token]crypto.Token
	Messages             []StageContentInfo
	Encrypted            bool
	Open                 bool
}

func (s *Stage) JoinRequest(req instructions.JoinStage, level byte, author *instructions.Author, epoch uint64, fee uint64) *instructions.AcceptJoinStage {
	if !req.Audience.Equal(s.Token) {
		return nil
	}

	s.MembersDiffieHellman[req.Authored.Author] = req.DiffHellKey
	accept := instructions.AcceptJoinStage{
		Authored: author.NewAuthored(epoch, fee),
		Stage:    s.Token,
		Member:   req.Authored.Author,
		Read:     []byte{},
		Submit:   []byte{},
		Moderate: []byte{},
	}

	cipher := dh.ConsensusCipher(s.DiffieHellman, req.DiffHellKey)
	accept.Read = cipher.Seal(s.CipherKey)
	s.Readers[req.Authored.Author] = struct{}{}
	if level > 0 {
		accept.Submit = cipher.Seal(s.Submit[:32])
		s.Submittors[req.Authored.Author] = struct{}{}
	}
	if level > 1 {
		s.Moderators[req.Authored.Author] = struct{}{}
		accept.Moderate = cipher.Seal(s.Moderate[:32])
	}
	modbulk := accept.serializeModBulk()
	accept.modSignature = audience.Moderation.Sign(modbulk)
	bulk := accept.serializeBulk()
	if a.sign(accept.Authored, bulk, IAcceptJoinRequest) {
		return &accept
	}
	return nil
}
