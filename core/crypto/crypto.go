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

const (
	NonceSize     = 12
	Size          = sha256.Size
	CipherKeySize = 32
	CipherSize    = NonceSize + CipherKeySize
	PublicKeySize = 42
)

type Hash [Size]byte

func (hash Hash) ToInt64() int64 {
	return int64(hash[0]) + (int64(hash[1]) << 8) + (int64(hash[2]) << 16) + (int64(hash[3]) << 24)
}

func BytesToHash(bytes []byte) Hash {
	if len(bytes) != Size {
		panic("invalid hash")
	}
	var h Hash
	for n := 0; n < Size; n++ {
		h[n] = bytes[n]
	}
	return h
}

func (h Hash) Equal(another Hash) bool {
	for n := 0; n < Size; n++ {
		if h[n] != another[n] {
			return false
		}
	}
	return true
}

func (h Hash) Equals(another []byte) bool {
	for n := 0; n < Size; n++ {
		if h[n] != another[n] {
			return false
		}
	}
	return true
}

func Hasher(data []byte) Hash {
	return Hash(sha256.Sum256(data))
}

func RandomAsymetricKey() (PublicKey, PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		panic(err)
	}
	publicKey := key.PublicKey
	return PublicKey{key: &publicKey}, PrivateKey{key: key}
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

func (p PrivateKey) Decrypt(msg []byte) ([]byte, error) {
	key := make([]byte, 32)
	err := rsa.DecryptPKCS1v15SessionKey(rand.Reader, p.key, msg, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (p PublicKey) Encrypt(msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, p.key, msg)
}

func (p PrivateKey) Sign(msg []byte) ([]byte, error) {
	hashed := sha256.Sum256(msg)
	return rsa.SignPKCS1v15(nil, p.key, crypto.SHA256, hashed[:])
}

func (p PublicKey) VerifyHash(hash Hash, signature []byte) bool {
	return rsa.VerifyPKCS1v15(p.key, crypto.SHA256, hash[:], signature) == nil
}

func (p PublicKey) Verify(msg []byte, signature []byte) bool {
	hashed := sha256.Sum256(msg)
	return rsa.VerifyPKCS1v15(p.key, crypto.SHA256, hashed[:], signature) == nil
}

func (p PublicKey) ToBytes() []byte {
	return x509.MarshalPKCS1PublicKey(p.key)
}

func PublicKeyFromBytes(bytes []byte) (PublicKey, error) {
	key, err := x509.ParsePKCS1PublicKey(bytes)
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey{key: key}, nil
}

func NewCipherKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

type Cipher struct {
	cipher cipher.AEAD
}

type CipherNonce struct {
	cipher cipher.AEAD
	nonce  []byte
}

func CipherNonceFromKey(key []byte) Cipher {
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
	return Cipher{cipher: gcm}
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
	var nonce []byte
	return c.cipher.Seal(nil, nonce, msg, nil)
}

func (c CipherNonce) Seal(msg []byte) []byte {
	return c.cipher.Seal(nil, c.nonce, msg, nil)
}

func (c CipherNonce) SealNewNonce(msg []byte) []byte {
	sealed := c.cipher.Seal(nil, c.nonce, msg, nil)
	if n, err := rand.Read(c.nonce); n != c.cipher.NonceSize() {
		panic(err)
	}
	return sealed
}

func (c Cipher) Open(msg []byte) ([]byte, error) {
	var nonce []byte
	return c.cipher.Open(nil, nonce, msg, nil)
}

func (c CipherNonce) Open(msg []byte) ([]byte, error) {
	return c.cipher.Open(nil, c.nonce, msg, nil)
}

func (c CipherNonce) OpenNewNonce(msg []byte, nonce []byte) ([]byte, error) {
	c.nonce = nonce
	return c.cipher.Open(nil, c.nonce, msg, nil)
}

func Equal() {

}

func Nonce() []byte {
	nonce := make([]byte, NonceSize)
	rand.Read(nonce)
	return nonce
}
