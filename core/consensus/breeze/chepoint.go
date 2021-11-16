package breeze

import (
	"bytes"

	"github.com/Aereum/aereum/core/blockchain"
	"github.com/Aereum/aereum/core/crypto"
)

type CheckPoint struct {
	Block            blockchain.Block
	Hash             crypto.Hash
	Signatures       [][]byte
	DropInstructions []crypto.Hash
}

func (c *CheckPoint) appendSignature(token crypto.PublicKey, signature []byte) (bool, int) {
	for _, s := range c.Signatures {
		if bytes.Equal(signature, s) {
			return false, len(c.Signatures)
		}
	}
	if token.VerifyHash(c.Hash, signature) {
		c.Signatures = append(c.Signatures, signature)
		return true, len(c.Signatures)
	}
	return false, len(c.Signatures)

}
