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

	// BLOCK 1
	block := NewBlock(crypto.Hasher([]byte{}), 0, 1, blockFormationToken.PublicKey().ToBytes(), validator)
	eve := &Author{token: &token}
	pubKey1, prvKey1 := crypto.RandomAsymetricKey()
	firstAuthor := &Author{token: &prvKey1, wallet: &token}
	joinFee := 10

	// Join Network
	join := eve.NewJoinNetworkThirdParty(pubKey1.ToBytes(), "First Member", jsonString1, 1, uint64(joinFee))
	if block.Incorporate(join) != true {
		t.Error("could not add new member")
	}

	// Block incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Members.Exists(crypto.Hasher(pubKey1.ToBytes())) {
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

	// BLOCK 2
	block = NewBlock(crypto.Hasher([]byte{}), 1, 2, blockFormationToken.PublicKey().ToBytes(), validator)
	count := 1

	// Join Network - Second member, sent by first member
	//_, wallet2 := crypto.RandomAsymetricKey()
	pubKey2, prvKey2 := crypto.RandomAsymetricKey()
	secondAuthor := &Author{token: &prvKey2, wallet: &token} // estou colocando a wallet da eve pra todo mundo
	join2 := secondAuthor.NewJoinNetwork("Second Member", jsonString1, 2, uint64(joinFee))
	if block.Incorporate(join2) != true {
		t.Error("could not add second member")
	}
	count = count + 1

	// Join Network - Third member, sent by first member
	//_, wallet3 := crypto.RandomAsymetricKey()
	_, prvKey3 := crypto.RandomAsymetricKey()
	thirdAuthor := &Author{token: &prvKey3, wallet: &token} // estou colocando a wallet da eve pra todo mundo
	join3 := thirdAuthor.NewJoinNetwork("Third Member", jsonString1, 2, uint64(joinFee))
	if block.Incorporate(join3) != true {
		t.Error("could not add third member")
	}
	count = count + 1

	// First author update info
	update := firstAuthor.NewUpdateInfo(jsonString1_new, 2, uint64(joinFee))
	if !block.Incorporate(update) {
		t.Error("could not add update")
	}
	count = count + 1

	// Create audience
	audienceTest := NewAudience()
	createAudience := firstAuthor.NewCreateAudience(audienceTest, 1, "first audience", 2, uint64(joinFee))
	if !block.Incorporate(createAudience) {
		t.Error("could not add create audience")
	}
	count = count + 1

	// Transfer from eve to first member
	transfer := NewSingleReciepientTransfer(*eve.token, secondAuthor.token.PublicKey().ToBytes(), "first transfer", 100, 2, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add transfer")
	}
	count = count + 1

	// Power of attorney
	poa := firstAuthor.NewGrantPowerOfAttorney(token.PublicKey().ToBytes(), 2, uint64(joinFee))
	if !block.Incorporate(poa) {
		t.Error("could not add poa")
	}
	count = count + 1

	// Block incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Members.Exists(crypto.Hasher(pubKey2.ToBytes())) {
		t.Error("state did not add second member")
	}
	if !state.Captions.Exists(crypto.Hasher([]byte("Second Member"))) {
		t.Error("state did not add second member caption")
	}
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance != 1e6-uint64(count*joinFee-100) {
		fmt.Print(balance)
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor = state.Wallets.Balance(crypto.Hasher(secondAuthor.token.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(100) {
		t.Error("first author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 3

	// Join audience
	joinAudience := secondAuthor.NewJoinAudience(audienceTest.token.ToBytes(), "first audience member", 2, uint64(joinFee))
	block.Incorporate(joinAudience)
	count = count + 1

	// Accept join audience
	firstAuthor.NewAcceptJoinAudience(audienceTest, secondAuthor.token.PublicKey(), 2, 2, uint64(joinFee))

	// Content
	firstAuthor.NewContent(audienceTest, "text", []byte("first content"), true, true, 2, uint64(joinFee))
}
