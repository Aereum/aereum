package crypto

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

// defines temporary crypto primitives

func RandomAsymetricKey() (PublicKey, PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return PublicKey{}, PrivateKey{}
	}
	return PublicKey{key: &key.PublicKey}, PrivateKey{key: key}
}

type PublicKey struct {
	key *rsa.PublicKey
}

type PrivateKey struct {
	key *rsa.PrivateKey
}

func (p PrivateKey) PublicKey() PublicKey {
	if p.key == nil {
		return PublicKey{}
	}
	return PublicKey{key: &p.key.PublicKey}
}

func (p PublicKey) IsValid() bool {
	return p.key != nil
}

func (p PrivateKey) IsValid() bool {
	return p.key != nil
}

func (p *PrivateKey) Decrypt(msg []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(nil, p.key, msg)
}

func (p *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(nil, p.key, msg)
}

func (p *PrivateKey) Sign(msg []byte) ([]byte, error) {
	hashed := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(nil, p.key, crypto.SHA256, hashed[:])
}

func (p *PublicKey) Verify(msg []byte, signature []byte) bool {
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(p.key, crypto.SHA256, hashed[:], signature) == nil
}

func (p *PublicKey) ToBytes() []byte {
	return x509.MarshalPKCS1PublicKey(p.key)
}

func PublicKeyFromBytes(bytes []byte) (PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(bytes)
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey{key: key}, nil
}

type Cipher struct {
	cipher cipher.AEAD
	nonce  []byte
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
	nonce := make([]byte, gcm.NonceSize())
	if n, err := rand.Read(nonce); n != gcm.NonceSize() {
		panic(err)
	}
	return Cipher{cipher: gcm, nonce: nonce}
}

func (c Cipher) Seal(msg []byte) []byte {
	return c.cipher.Seal(nil, c.nonce, msg, nil)
}

func (c Cipher) Open(msg []byte) ([]byte, error) {
	return c.cipher.Open(nil, c.nonce, msg, nil)
}
