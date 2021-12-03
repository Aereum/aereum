package store

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestHashExpire(t *testing.T) {
	vault := NewExpireHashVault("teste", 0, 8)
	hash := crypto.Hasher([]byte{1, 2, 3, 4})
	vault.Insert(hash, 12)
	if vault.Exists(hash) != 12 {
		t.Errorf("vault insert/exists not working")
	}
	if !vault.Remove(hash) || vault.Exists(hash) != 0 {
		t.Errorf("vault remove not working")
	}

}
