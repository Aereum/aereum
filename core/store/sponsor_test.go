package store

import (
	"bytes"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestSponsor(t *testing.T) {
	sponsor := NewSponsorShipOfferStore(0, 8)
	hash := crypto.Hasher([]byte{1, 2, 3, 4})
	hash2 := crypto.Hasher([]byte{1, 4, 3, 4})
	sponsor.SetContentHash(hash, hash2[:])
	ok, hash3 := sponsor.GetContentHash(hash)
	if !ok || bytes.Equal(hash2[:], hash3) {
		t.Errorf("sponsor set/get not working")
	}
	sponsor.RemoveContentHash(hash)
	if ok, _ := sponsor.GetContentHash(hash); ok {
		t.Errorf("sponsor remove not working")
	}

}
