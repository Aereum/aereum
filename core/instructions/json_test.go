package instructions

import (
	"encoding/json"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestJSON(t *testing.T) {
	_, prvKey1 := crypto.RandomAsymetricKey()
	firstAuthor := &Author{token: &prvKey1, wallet: &token}
	join := firstAuthor.NewJoinNetwork("First Author", `{"update":1}`, 10, 10)
	if !json.Valid([]byte(join.JSON())) {
		t.Error("invalid join network json")
	}
}
