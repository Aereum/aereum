package content

import "github.com/Aereum/aereum/core/crypto"

type Audience struct {
	Token [crypto.PublicKeySize]byte
	//Flags           byte
	PublishingToken [crypto.PublicKeySize]byte
	SubmissionToken [crypto.PublicKeySize]byte
	Description     []byte // can submit? can forward?
}

type AudienceKeyRotation struct {
	Flags           byte
	PublishingToken [crypto.PublicKeySize]byte
	SubmissionToken [crypto.PublicKeySize]byte
	Description     []byte
	Readers         [][crypto.Size + crypto.CipherKeySize]byte
	Publishers      [][crypto.Size + crypto.CipherKeySize]byte
	Submitors       [][crypto.Size + crypto.CipherKeySize]byte
}

type Content struct {
	Audience            crypto.PublicKey
	ContentType         byte
	ContentNonce        []byte
	ContentData         []byte
	ContentHash         crypto.Hash
	SubmissionSignature []byte // author + epoch +
	PublishingSignature []byte
}
