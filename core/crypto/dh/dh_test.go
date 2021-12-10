package dh

import (
	"bytes"
	"testing"
)

func TestDH(t *testing.T) {
	alice := NewAlice()
	bob := NewBob(alice.keyX)
	if !alice.Incorporate(bob.keyX) {
		t.Error("dh scheme not working")
	}
	if !bytes.Equal(alice.agreedKey, bob.agreedKey) {
		t.Error("dh scheme not working")
	}

	aliceCipher := alice.Cipher()
	bobCipher := bob.Cipher()
	secret := aliceCipher.Seal([]byte("Testando"))
	original, err := bobCipher.Open(secret)
	if err != nil || string(original) != "Testando" {
		t.Error("dh scheme cipher not working")
	}

}
