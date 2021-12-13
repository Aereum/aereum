package crypto

import (
	"bytes"
	"testing"
)

func TestCipher(t *testing.T) {
	key := NewCipherKey()
	cipher := CipherFromKey(key)
	data := []byte{1, 4, 5, 6, 8, 10, 23, 45, 89, 113}
	secret := cipher.Seal(data)
	secretdata, err := cipher.Open(secret)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, secretdata) {
		t.Error("cipher not working")
	}
}

func TestCipherNonce(t *testing.T) {
	key := NewCipherKey()
	cipher := CipherNonceFromKey(key)
	data := []byte{1, 4, 5, 6, 8, 10, 23, 45, 89, 113}
	secret := cipher.Seal(data)
	secretdata, err := cipher.Open(secret)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, secretdata) {
		t.Error("cipher not working")
	}
}

func TestCipherNewNonce(t *testing.T) {
	key := NewCipherKey()
	cipher := CipherNonceFromKey(key)
	cipher2 := CipherNonceFromKey(key)
	data := []byte{1, 4, 5, 6, 8, 10, 23, 45, 89, 113}
	secret, nonce := cipher.SealWithNewNonce(data)
	secretdata, err := cipher2.OpenNewNonce(secret, nonce)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, secretdata) {
		t.Error("cipher not working")
	}
}

func TestPublicKey(t *testing.T) {
	pub, prv := RandomAsymetricKey()
	data := []byte{1, 4, 5, 6, 8, 10, 23, 45, 89, 113}
	sign := prv.Sign(data)
	if !pub.Verify(data, sign) {
		t.Errorf("signature not working")
	}
}
