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

	// Starting block
	block := NewBlock(crypto.Hasher([]byte{}), 0, 1, blockFormationToken.PublicKey().ToBytes(), validator)
	eve := &Author{token: &token}
	eveBalance := 1e6
	joinFee := 10
	count := 0

	// First member crypto data
	pubKey1, prvKey1 := crypto.RandomAsymetricKey()
	_, prvWal1 := crypto.RandomAsymetricKey()
	firstAuthor := &Author{token: &prvKey1, wallet: &prvWal1}
	firstBalance := 0

	// Join Network sent by eve (pq nao posso usar join network normal?)
	join := eve.NewJoinNetworkThirdParty(pubKey1.ToBytes(), "member1", jsonString1, 1, uint64(joinFee))
	if block.Incorporate(join) != true {
		t.Error("could not add new member")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// Block 1 incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Members.Exists(crypto.Hasher(pubKey1.ToBytes())) {
		t.Error("state did not add new member")
	}
	if !state.Captions.Exists(crypto.Hasher([]byte("member1"))) {
		t.Error("state did not add new member caption")
	}
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance != uint64(eveBalance) {
		fmt.Print(balance)
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor := state.Wallets.Balance(crypto.Hasher(firstAuthor.wallet.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(0) {
		t.Error("first author wallet must start with zero aero")
	}
	_, balanceBlockFormator := state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee")
	}

	// BLOCK 2
	block = NewBlock(crypto.Hasher([]byte{}), 1, 2, blockFormationToken.PublicKey().ToBytes(), validator)

	// Transfer from eve to first member
	transfer := NewSingleReciepientTransfer(*eve.token, firstAuthor.wallet.PublicKey().ToBytes(), "first transfer", 100, 2, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add transfer")
	}
	firstBalance = firstBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Join Network - Second member, sent by eve
	pubKey2, prvKey2 := crypto.RandomAsymetricKey()
	_, prvWal2 := crypto.RandomAsymetricKey()
	secondAuthor := &Author{token: &prvKey2, wallet: &prvWal2}
	secondBalance := 0
	join2 := eve.NewJoinNetworkThirdParty(pubKey2.ToBytes(), "member2", jsonString1, 2, uint64(joinFee))
	if block.Incorporate(join2) != true {
		t.Error("could not add member2")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// Join Network - Third member, sent by member1
	pubKey3, prvKey3 := crypto.RandomAsymetricKey()
	_, prvWal3 := crypto.RandomAsymetricKey()
	thirdAuthor := &Author{token: &prvKey3, wallet: &prvWal3}
	thirdBalance := 0
	join3 := eve.NewJoinNetworkThirdParty(pubKey3.ToBytes(), "member3", jsonString1, 2, uint64(joinFee))
	if block.Incorporate(join3) != true {
		t.Error("could not add member3")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// First author update info
	update := eve.NewUpdateInfo(jsonString1_new, 2, uint64(joinFee))
	if !block.Incorporate(update) {
		t.Error("could not add update")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// Block 2 incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Members.Exists(crypto.Hasher(pubKey2.ToBytes())) {
		t.Error("state did not add second member")
	}
	if !state.Captions.Exists(crypto.Hasher([]byte("member2"))) {
		t.Error("state did not add second member caption")
	}
	if !state.Members.Exists(crypto.Hasher(pubKey3.ToBytes())) {
		t.Error("state did not add third member")
	}
	if !state.Captions.Exists(crypto.Hasher([]byte("member3"))) {
		t.Error("state did not add third member caption")
	}
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance != uint64(eveBalance) {
		fmt.Print(balance)
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor = state.Wallets.Balance(crypto.Hasher(firstAuthor.wallet.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not receive eve transfer")
	}
	_, balanceBlockFormator = state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 3
	block = NewBlock(crypto.Hasher([]byte{}), 2, 3, blockFormationToken.PublicKey().ToBytes(), validator)

	// Transfer from eve to second member
	transfer = NewSingleReciepientTransfer(*eve.token, secondAuthor.wallet.PublicKey().ToBytes(), "second transfer", 100, 3, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add second transfer")
	}
	secondBalance = secondBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Transfer from eve to third member
	transfer = NewSingleReciepientTransfer(*eve.token, thirdAuthor.wallet.PublicKey().ToBytes(), "third transfer", 100, 3, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add third transfer")
	}
	thirdBalance = thirdBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Create audience
	audienceTest := NewAudience()
	createAudience := firstAuthor.NewCreateAudience(audienceTest, 1, "first audience", 3, uint64(joinFee))
	if !block.Incorporate(createAudience) {
		t.Error("could not add create audience")
	}
	count = count + 1
	firstBalance = firstBalance - joinFee

	// Power of attorney sent by first author with third author as attorney
	poa := firstAuthor.NewGrantPowerOfAttorney(thirdAuthor.token.PublicKey().ToBytes(), 3, uint64(joinFee))
	if !block.Incorporate(poa) {
		t.Error("could not add poa")
	}
	firstAuthor.attorney = thirdAuthor.token
	count = count + 1
	firstBalance = firstBalance - joinFee

	// Block 3 incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Audiences.Exists(crypto.Hasher(audienceTest.token.PublicKey().ToBytes())) {
		t.Error("state did not create audience")
	}
	hashAttorney := crypto.Hasher(append(firstAuthor.token.PublicKey().ToBytes(), thirdAuthor.token.PublicKey().ToBytes()...))
	if !state.PowerOfAttorney.Exists(hashAttorney) {
		t.Error("power of attorney was not granted")
	}
	if _, balance := state.Wallets.Balance(crypto.Hasher(token.PublicKey().ToBytes())); balance != uint64(eveBalance) {
		fmt.Print(balance)
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor = state.Wallets.Balance(crypto.Hasher(firstAuthor.wallet.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spent on instructions")
	}
	_, balanceSecondAuthor := state.Wallets.Balance(crypto.Hasher(secondAuthor.wallet.PublicKey().ToBytes()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not receive transfer")
	}
	_, balanceThirdAuthor := state.Wallets.Balance(crypto.Hasher(thirdAuthor.wallet.PublicKey().ToBytes()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 4
	block = NewBlock(crypto.Hasher([]byte{}), 3, 4, blockFormationToken.PublicKey().ToBytes(), validator)

	// Join audience sent by second member
	joinAudience := secondAuthor.NewJoinAudience(audienceTest.token.ToBytes(), "first audience member", 4, uint64(joinFee))
	if !block.Incorporate(joinAudience) {
		t.Error("could not send join audience instruction")
	}
	secondBalance = secondBalance - joinFee
	count = count + 1

	// Content
	content := firstAuthor.NewContent(audienceTest, "text", []byte("first content"), true, true, 4, uint64(joinFee))
	if !block.Incorporate(content) {
		t.Error("could not publish content to audience")
	}
	firstBalance = firstBalance - joinFee
	count = count + 1

	// Sponsorship Offer
	sponsorOffer := thirdAuthor.NewSponsorshipOffer(audienceTest, "txt", []byte("sponsor"), 20, 20, 4, uint64(joinFee))
	if !block.Incorporate(sponsorOffer) {
		t.Error("could not publish sponsor offer to audience")
	}
	thirdBalance = thirdBalance - joinFee
	count = count + 1

	// Block 4 incorporation and balance checks
	state.IncorporateBlock(block)

	// COMO CHECAR SPONSOR OFFER
	// if !state.SponsorOffers.Exists() {
	// 	t.Error("sponsor offer was not incorporated")
	// }
	_, balanceFirstAuthor = state.Wallets.Balance(crypto.Hasher(firstAuthor.wallet.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spent on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.Balance(crypto.Hasher(secondAuthor.wallet.PublicKey().ToBytes()))
	fmt.Print(balanceSecondAuthor)
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not receive transfer")
	}
	_, balanceThirdAuthor = state.Wallets.Balance(crypto.Hasher(thirdAuthor.wallet.PublicKey().ToBytes()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 5
	block = NewBlock(crypto.Hasher([]byte{}), 4, 5, blockFormationToken.PublicKey().ToBytes(), validator)

	// Accept join audience
	acceptJoin := thirdAuthor.NewAcceptJoinAudience(audienceTest, secondAuthor.token.PublicKey(), 2, 2, uint64(joinFee))
	if !block.Incorporate(acceptJoin) {
		t.Error("could not accept join request to audience")
	}
	// firstBalance = firstBalance - joinFee // attorney esta enviando em nome de member1, porem quem paga eh member1
	thirdBalance = thirdBalance - joinFee
	count = count + 1

	// Accept sponsor offer
	sponsordAccept := thirdAuthor.NewSponsorshipAcceptance(audienceTest, sponsorOffer, 5, uint64(joinFee))
	if !block.Incorporate(sponsordAccept) {
		t.Error("could not accept sponsorship acceptance")
	}
	// firstBalance = firstBalance - joinFee
	thirdBalance = thirdBalance - joinFee
	count = count + 1

	state.IncorporateBlock(block)
	_, balanceFirstAuthor = state.Wallets.Balance(crypto.Hasher(firstAuthor.token.PublicKey().ToBytes()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spent on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.Balance(crypto.Hasher(secondAuthor.token.PublicKey().ToBytes()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not receive transfer")
	}
	_, balanceThirdAuthor = state.Wallets.Balance(crypto.Hasher(thirdAuthor.token.PublicKey().ToBytes()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.Balance(crypto.Hasher(blockFormationToken.PublicKey().ToBytes()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

}
