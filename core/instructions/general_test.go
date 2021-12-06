package instructions

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

var (
	jsonString1     = `{"teste":1}`
	jsonString1_new = `{"update":1}`
)

func TestGeneral(t *testing.T) {
	state, token := NewGenesisState()
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance == 0 {
		t.Error("wrong genesis")
	}
	validator := &Validator{State: state}
	block := NewBlock(crypto.Hasher([]byte{}), 0, 1, token.PublicKey().ToBytes(), validator)
	creator := &Author{token: &token}
	pubKey, prvKey := crypto.RandomAsymetricKey()
	firstAuthor := &Author{token: &prvKey, wallet: &token}
	join := creator.NewJoinNetworkThirdParty(pubKey.ToBytes(), "First Member", jsonString1, 1, 1)
	if block.Incorporate(join) != true {
		t.Error("could not add new member")
	}
	state.IncorporateBlock(block)
	if !state.Members.Exists(crypto.Hasher(pubKey.ToBytes())) {
		t.Error("state did not add new member")
	}
	if !state.Captions.Exists(crypto.Hasher([]byte("First Member"))) {
		t.Error("state did not add new caption")
	}
	block = NewBlock(crypto.Hasher([]byte{}), 0, 2, token.PublicKey().ToBytes(), validator)
	update := firstAuthor.NewUpdateInfo(jsonString1_new, 12, 10)
	block.Incorporate(update)
	state.IncorporateBlock(block)
}
