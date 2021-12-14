package store

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestHashVault(t *testing.T) {
	vault := NewHashVault("teste", 0, 8)
	hash := crypto.Hasher([]byte{1, 2, 3, 4})
	vault.InsertHash(hash)
	if !vault.ExistsHash(hash) {
		t.Errorf("vault insert/exists not working")
	}
	if !vault.RemoveHash(hash) || vault.ExistsHash(hash) {
		t.Errorf("vault remove not working")
	}

}
