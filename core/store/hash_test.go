package store

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestHashVault(t *testing.T) {
	vault := NewHashVault("teste", 0, 8)
	hash := crypto.Hasher([]byte{1, 2, 3, 4})
	vault.Insert(hash)
	if !vault.Exists(hash) {
		t.Errorf("vault insert/exists not working")
	}
	if !vault.Remove(hash) || vault.Exists(hash) {
		t.Errorf("vault remove not working")
	}

}
