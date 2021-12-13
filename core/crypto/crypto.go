package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

// defines temporary crypto primitives

const (
	NonceSize      = 12
	Size           = sha256.Size
	CipherKeySize  = 32
	CipherSize     = NonceSize + CipherKeySize
	PublicKeySize  = 32
	TokenSize      = 32
	PrivateKeySize = 64
	SignatureSize  = 64
)

type Cipher struct {
	cipher cipher.AEAD
}

type CipherNonce struct {
	cipher cipher.AEAD
	nonce  []byte
}

func NewCipherKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

func CipherNonceFromKey(key []byte) CipherNonce {
	if len(key) != 32 {
		panic("wrong cipher key size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if n, err := rand.Read(nonce); n != gcm.NonceSize() {
		panic(err)
	}
	return CipherNonce{cipher: gcm, nonce: nonce}
}

func CipherFromKey(key []byte) Cipher {
	if len(key) != 32 {
		panic("wrong cipher key size")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return Cipher{cipher: gcm}
}

func (c Cipher) Seal(msg []byte) []byte {
	nonce := make([]byte, NonceSize)
	return c.cipher.Seal(nil, nonce, msg, nil)
}

func (c CipherNonce) Seal(msg []byte) []byte {
	return c.cipher.Seal(nil, c.nonce, msg, nil)
}

func (c CipherNonce) SetNonce(nonce []byte) {
	c.nonce = nonce
}

func (c CipherNonce) SealWithNewNonce(msg []byte) ([]byte, []byte) {
	if n, err := rand.Read(c.nonce); n != c.cipher.NonceSize() {
		panic(err)
	}
	sealed := c.cipher.Seal(nil, c.nonce, msg, nil)
	return sealed, c.nonce
}

func (c Cipher) Open(msg []byte) ([]byte, error) {
	nonce := make([]byte, NonceSize)
	return c.cipher.Open(nil, nonce, msg, nil)
}

func (c CipherNonce) Open(msg []byte) ([]byte, error) {
	return c.cipher.Open(nil, c.nonce, msg, nil)
}

func (c CipherNonce) OpenNewNonce(msg []byte, nonce []byte) ([]byte, error) {
	c.nonce = nonce
	return c.cipher.Open(nil, c.nonce, msg, nil)
}

func Nonce() []byte {
	nonce := make([]byte, NonceSize)
	rand.Read(nonce)
	return nonce
}
