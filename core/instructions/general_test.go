package instructions

import (
	"fmt"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

var (
	jsonString1     = `{"teste":1}`
	jsonString1_new = `{"update":1}`
)

func TestGeneral(t *testing.T) {

	state, token := NewGenesisState()
	_, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes()))
	if balance == 0 {
		t.Error("wrong genesis")
	}

	validator := &Validator{State: state}
	_, blockFormationToken := crypto.RandomAsymetricKey()

	// First Block
	block := NewBlock(crypto.Hasher([]byte{}), 0, 1, blockFormationToken.PublicKey().ToBytes(), validator)
	creator := &Author{token: &token}
	pubKey, prvKey := crypto.RandomAsymetricKey()
	firstAuthor := &Author{token: &prvKey, wallet: &token}
	joinFee := 10

	// Join Network
	join := creator.NewJoinNetworkThirdParty(pubKey.ToBytes(), "First Member", jsonString1, 1, uint64(joinFee))
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
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance != 1e6-uint64(joinFee) {
		fmt.Print(balance)
		t.Error("state did not debit wallet", balance)
	}

	_, balanceFirstAuthor := state.Wallets.Balance(crypto.Hasher(firstAuthor.token.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(0) {
		t.Error("first author wallet must start with zero aero")
	}

	_, balanceBlockFormator := state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(joinFee) {
		t.Error("block formator has not received fee")
	}

	// Second member Network instruction sent by first member
	_, prvKey2 := crypto.RandomAsymetricKey()
	secondAuthor := &Author{token: &prvKey2, wallet: &token}
	join2 := secondAuthor.NewJoinNetwork("Second Member", jsonString1, 1, uint64(joinFee))
	if block.Incorporate(join2) != true {
		t.Error("could not add new member")
	}

	// Second Block
	block = NewBlock(crypto.Hasher([]byte{}), 0, 2, token.PublicKey().ToBytes(), validator)

	update := firstAuthor.NewUpdateInfo(jsonString1_new, 12, 10)
	block.Incorporate(update)
	state.IncorporateBlock(block)

	// First author update info
	firstAuthor.NewUpdateInfo(jsonString1_new, 1, uint64(joinFee))

	// Create audience
	var (
		audienceTest *Audience = NewAudience()
	)
	firstAuthor.NewCreateAudience(audienceTest, 1, "first audience", 2, uint64(joinFee))

	// Join audience
	secondAuthor.NewJoinAudience(audienceTest.token.ToBytes(), "first audience member", 2, uint64(joinFee))

	// Accept join audience
	firstAuthor.NewAcceptJoinAudience(audienceTest, secondAuthor.token.PublicKey(), 2, 2, uint64(joinFee))

	// Content
	// firstAuthor.NewContent(audienceTest, "text", PutString("first post"), 1, 1, 2, uint64(joinFee))
}
