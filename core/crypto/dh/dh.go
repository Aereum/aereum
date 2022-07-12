package dh

// package dh implements the Diffie-Hellman secret exchange defined on RFC 7748
// using the X25519 function
// (begin of quote)
// Alice generates 32 random bytes in a[0] to a[31] and transmits K_A =
// X25519(a, 9) to Bob, where 9 is the u-coordinate of the base point
// and is encoded as a byte with value 9, followed by 31 zero bytes.
//
// Bob similarly generates 32 random bytes in b[0] to b[31], computes
// K_B = X25519(b, 9), and transmits it to Alice.
//
// Using their generated values and the received input, Alice computes
// X25519(a, K_B) and Bob computes X25519(b, K_A).
//
// Both now share K = X25519(a, X25519(b, 9)) = X25519(b, X25519(a, 9))
// as a shared secret.  Both MAY check, without leaking extra
// information about the value of K, whether K is the all-zero value and
// abort if so (see below).  Alice and Bob can then use a key-derivation
// function that includes K, K_A, and K_B to derive a symmetric key.
//
// The check for the all-zero value results from the fact that the
// X25519 function produces that value if it operates on an input
// corresponding to a point with small order, where the order divides
// the cofactor of the curve (see Section 7).  The check may be
// performed by ORing all the bytes together and checking whether the
// result is zero, as this eliminates standard side-channels in software
// implementations.
// (end of quote)
// the SHA256 hash on the agreed key is used as a key for an AES 256 Cipher.

import (
	"fmt"
	"math/rand"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh/curve25519"
)

type Party struct {
	key       []byte
	keyX      []byte
	agreedKey []byte
}

func NewEphemeralKey() (crypto.PrivateKey, crypto.Token) {
	var pubToken crypto.Token
	var prvToken crypto.PrivateKey
	rand.Read(prvToken[0:32])
	pub, err := curve25519.X25519(prvToken[0:32], curve25519.Basepoint)
	if err == nil {
		copy(prvToken[32:], pub)
		copy(pubToken[:], pub)
	}
	return prvToken, pubToken
}

func ConsensusKey(local crypto.PrivateKey, remote crypto.Token) []byte {
	agreedKey, err := curve25519.X25519(local[0:32], remote[:])
	if err != nil {
		fmt.Println("------------------", err)
		return nil
	}
	hashed := crypto.Hasher(agreedKey)
	return hashed[:]
}

func ConsensusCipher(local crypto.PrivateKey, remote crypto.Token) crypto.Cipher {
	return crypto.CipherFromKey(ConsensusKey(local, remote))
}

func NewEphemeralRequest() *Party {
	rnd := make([]byte, 32)
	rand.Read(rnd)
	rndX, err := curve25519.X25519(rnd, curve25519.Basepoint)
	if err != nil {
		return nil
	}
	return &Party{key: rnd, keyX: rndX}
}

func NewEphemeralResponse(aliceKeyX []byte) *Party {
	rnd := make([]byte, 32)
	rand.Read(rnd)
	rndX, err := curve25519.X25519(rnd, curve25519.Basepoint)
	if err != nil {
		return nil
	}
	agreedKey, err := curve25519.X25519(rnd, aliceKeyX)
	if err != nil {
		return nil
	}
	return &Party{key: rnd, keyX: rndX, agreedKey: agreedKey}
}

func (p *Party) IncorporateResponse(otherKeyX []byte) bool {
	agreedKey, err := curve25519.X25519(p.key, otherKeyX)
	if err != nil {
		return false
	}
	p.agreedKey = agreedKey
	return true
}

func (p *Party) Cipher() crypto.Cipher {
	hashed := crypto.Hasher(p.agreedKey)
	return crypto.CipherFromKey(hashed[:])
}

func (p *Party) CipherNonce() crypto.CipherNonce {
	hashed := crypto.Hasher(p.agreedKey)
	return crypto.CipherNonceFromKey(hashed[:])
}
