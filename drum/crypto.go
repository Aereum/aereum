package main

import (
	"io"
	"log"
	"os"

	"golang.org/x/crypto/scrypt"

	"github.com/Aereum/aereum/core/crypto"
)

const (
	TypePrivateKey byte = iota
	TypeStageSecrets
)

type StageSecrets struct {
	Ownership  crypto.PrivateKey
	Moderation crypto.PrivateKey
	Submission crypto.PrivateKey
	CipherKey  []byte
}

type SecureVault struct {
	SecretKey crypto.PrivateKey
	Stages    map[crypto.Token]*StageSecrets
	Secrets   map[crypto.Token]crypto.PrivateKey
	Ephemeral []crypto.PrivateKey
	file      io.WriteCloser
	cipher    crypto.Cipher
}

func (vault *SecureVault) AppendData(SecretType byte, data []byte) {
	sealed := append([]byte{SecretType}, vault.cipher.Seal(data)...)
	if n, err := vault.file.Write(sealed); n != len(sealed) || err != nil {
		log.Fatal("Could not write to secure vault.")
	}
}

func (vault *SecureVault) NewStageKeys(encrypted, public bool) *StageSecrets {
	token, own := crypto.RandomAsymetricKey()
	_, moderate := crypto.RandomAsymetricKey()
	var submit crypto.PrivateKey
	if !public {
		_, submit = crypto.RandomAsymetricKey()
	}
	cipherKey := make([]byte, crypto.CipherKeySize)
	if encrypted {
		cipherKey = crypto.NewCipherKey()
	}
	data := append(token[:], own[:]...)
	data = append(data, moderate[:]...)
	data = append(data, submit[:]...)
	data = append(data, cipherKey...)
	vault.AppendData(TypeStageSecrets, data)
	secrets := &StageSecrets{
		Ownership:  own,
		Moderation: moderate,
		Submission: submit,
		CipherKey:  cipherKey,
	}
	vault.Stages[token] = secrets
	return secrets
}

func (vault *SecureVault) Moderation(stage crypto.Token) crypto.PrivateKey {
	var secret crypto.PrivateKey
	if secrets, ok := vault.Stages[stage]; ok {
		secret = secrets.Moderation
	}
	return secret
}

func (vault *SecureVault) Owenership(stage crypto.Token) crypto.PrivateKey {
	var secret crypto.PrivateKey
	if secrets, ok := vault.Stages[stage]; ok {
		secret = secrets.Ownership
	}
	return secret
}

func (vault *SecureVault) Submission(stage crypto.Token) crypto.PrivateKey {
	var secret crypto.PrivateKey
	if secrets, ok := vault.Stages[stage]; ok {
		secret = secrets.Submission
	}
	return secret
}

func (vault *SecureVault) CipherKey(stage crypto.Token) []byte {
	if secrets, ok := vault.Stages[stage]; ok {
		return secrets.CipherKey
	}
	return nil
}

func (vault *SecureVault) Secret(token crypto.Token) crypto.PrivateKey {
	var zeroToken crypto.PrivateKey
	if secret, ok := vault.Secrets[token]; ok {
		return secret
	}
	return zeroToken
}

func (vault *SecureVault) NewKey() crypto.PrivateKey {
	_, secret := crypto.RandomAsymetricKey()
	vault.AppendData(TypePrivateKey, secret[:])
	return secret
}

func NewSecureVault(password string, file io.Writer) {
	public, secret := crypto.RandomAsymetricKey()
	if n, err := file.Write(public[:]); n != len(public) || err != nil {
		log.Fatal("Could not open secret vault file.")
	}
	key, err := scrypt.Key([]byte(password), public[:], 32768, 8, 1, 32)
	if err != nil {
		log.Fatal("Could not derive cipher from password.")
	}
	cipher := crypto.CipherFromKey(key)
	sealed := cipher.Seal(secret[:])
	if n, err := file.Write(sealed); n != len(sealed) || err != nil {
		log.Fatal("Could not open secret vault file.")
	}
}

func ReadKey(data []byte, position int, cipher crypto.Cipher) (crypto.PrivateKey, int) {
	key, err := cipher.Open(data[position : position+crypto.Size+crypto.NonceSize])
	if err != nil {
		log.Fatal("Could not parse Secure Vault")
	}
	var pk crypto.PrivateKey
	copy(pk[0:crypto.PrivateKeySize], key)
	return pk, position + crypto.Size + crypto.NonceSize
}

func OpenVaultFromPassword(password []byte, fileName string) *SecureVault {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Could not open Secret Vault")
	}
	vault, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Could not read Secret Vault")
	}
	if err := file.Close(); err != nil {
		log.Fatal("Could not close Secret Vault")
	}
	var token crypto.Token
	copy(token[:], vault[0:crypto.Size])
	key, err := scrypt.Key(password, token[:], 32768, 8, 1, 32)
	cipher := crypto.CipherFromKey(key)
	secretKey, _ := ReadKey(vault, 0, cipher)
	if !secretKey.PublicKey().Equal(token) {
		log.Fatal("Wrong password.")
	}
	return nil
}
