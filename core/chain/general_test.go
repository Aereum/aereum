package chain

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/instructions"
)

var (
	jsonString1     = `{"teste":1}`
	jsonString1_new = `{"update":1}`
)

func TestGeneral(t *testing.T) {

	state, token := NewGenesisState()
	_, balance := state.Wallets.BalanceHash(crypto.HashToken(token.PublicKey()))
	if balance == 0 {
		t.Error("wrong genesis")
	}

	validator := &MutatingState{State: state}
	_, blockFormationToken := crypto.RandomAsymetricKey()

	// BLOCK 1

	// Starting block
	block := NewBlock(crypto.Hasher([]byte{}), 0, 1, blockFormationToken.PublicKey(), validator)
	eve := &instructions.Author{PrivateKey: token}
	eveBalance := 1e6
	joinFee := 10
	count := 0

	// First member crypto data
	pubKey1, prvKey1 := crypto.RandomAsymetricKey()
	_, prvWal1 := crypto.RandomAsymetricKey()
	firstAuthor := &instructions.Author{PrivateKey: prvKey1, Wallet: prvWal1}
	firstBalance := 0

	// Join Network sent by eve (pq nao posso usar join network normal?)
	join := eve.NewJoinNetworkThirdParty(pubKey1, "member1", jsonString1, 1, uint64(joinFee))
	if block.Incorporate(join) != true {
		t.Error("could not add new member")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// Block 1 incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Members.ExistsHash(crypto.HashToken(pubKey1)) {
		t.Error("state did not add new member")
	}
	if !state.Captions.ExistsHash(crypto.Hasher([]byte("member1"))) {
		t.Error("state did not add new member caption")
	}
	if _, balance := state.Wallets.BalanceHash(crypto.HashToken(token.PublicKey())); balance != uint64(eveBalance) {
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor := state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(0) {
		t.Error("first author wallet must start with zero aero")
	}
	_, balanceBlockFormator := state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee")
	}

	// BLOCK 2
	block = NewBlock(crypto.Hasher([]byte{}), 1, 2, blockFormationToken.PublicKey(), validator)

	// Transfer from eve to first member
	transfer := instructions.NewSingleReciepientTransfer(eve.PrivateKey, firstAuthor.Wallet.PublicKey(), "first transfer", 100, 2, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add transfer")
	}
	firstBalance = firstBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Join Network - Second member, sent by eve
	pubKey2, prvKey2 := crypto.RandomAsymetricKey()
	_, prvWal2 := crypto.RandomAsymetricKey()
	secondAuthor := &instructions.Author{PrivateKey: prvKey2, Wallet: prvWal2}
	secondBalance := 0
	join2 := eve.NewJoinNetworkThirdParty(pubKey2, "member2", jsonString1, 2, uint64(joinFee))
	if block.Incorporate(join2) != true {
		t.Error("could not add member2")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	// Join Network - Third member, sent by member1
	pubKey3, prvKey3 := crypto.RandomAsymetricKey()
	_, prvWal3 := crypto.RandomAsymetricKey()
	thirdAuthor := &instructions.Author{PrivateKey: prvKey3, Wallet: prvWal3}
	thirdBalance := 0
	join3 := eve.NewJoinNetworkThirdParty(pubKey3, "member3", jsonString1, 2, uint64(joinFee))
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
	if !state.Members.ExistsHash(crypto.HashToken(pubKey2)) {
		t.Error("state did not add second member")
	}
	if !state.Captions.ExistsHash(crypto.Hasher([]byte("member2"))) {
		t.Error("state did not add second member caption")
	}
	if !state.Members.ExistsHash(crypto.HashToken(pubKey3)) {
		t.Error("state did not add third member")
	}
	if !state.Captions.ExistsHash(crypto.Hasher([]byte("member3"))) {
		t.Error("state did not add third member caption")
	}
	if _, balance := state.Wallets.BalanceHash(crypto.HashToken(token.PublicKey())); balance != uint64(eveBalance) {
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not receive eve transfer")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 3
	block = NewBlock(crypto.Hasher([]byte{}), 2, 3, blockFormationToken.PublicKey(), validator)

	// Transfer from eve to second member
	transfer = instructions.NewSingleReciepientTransfer(eve.PrivateKey, secondAuthor.Wallet.PublicKey(), "second transfer", 100, 3, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add second transfer")
	}
	secondBalance = secondBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Transfer from eve to third member
	transfer = instructions.NewSingleReciepientTransfer(eve.PrivateKey, thirdAuthor.Wallet.PublicKey(), "third transfer", 100, 3, uint64(joinFee))
	if !block.Incorporate(transfer) {
		t.Error("could not add third transfer")
	}
	thirdBalance = thirdBalance + 100
	eveBalance = eveBalance - 100 - float64(joinFee)
	count = count + 1

	// Create audience
	audienceTest := instructions.NewStage(0, "teste")
	createAudience := firstAuthor.NewCreateAudience(audienceTest, 1, "first audience", 3, uint64(joinFee))
	if !block.Incorporate(createAudience) {
		t.Error("could not add create audience")
	}
	count = count + 1
	firstBalance = firstBalance - joinFee

	// Power of attorney sent by first author with third author as attorney
	poa := firstAuthor.NewGrantPowerOfAttorney(thirdAuthor.PrivateKey.PublicKey(), 3, uint64(joinFee))
	if !block.Incorporate(poa) {
		t.Error("could not add poa")
	}
	firstAuthor.Attorney = thirdAuthor.PrivateKey
	count = count + 1
	firstBalance = firstBalance - joinFee

	// Block 3 incorporation and balance checks
	state.IncorporateBlock(block)
	if !state.Stages.Exists(crypto.HashToken(audienceTest.PrivateKey.PublicKey())) {
		t.Error("state did not create audience")
	}
	firstAuthorKey := firstAuthor.PrivateKey.PublicKey()
	thirdAuthorKey := thirdAuthor.PrivateKey.PublicKey()
	hashAttorney := crypto.Hasher(append(firstAuthorKey[:], thirdAuthorKey[:]...))
	if !state.PowerOfAttorney.ExistsHash(hashAttorney) {
		t.Error("power of attorney was not granted")
	}
	if _, balance := state.Wallets.BalanceHash(crypto.HashToken(token.PublicKey())); balance != uint64(eveBalance) {
		t.Error("state did not debit wallet", balance)
	}
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spent on instructions")
	}
	_, balanceSecondAuthor := state.Wallets.BalanceHash(crypto.HashToken(secondAuthor.Wallet.PublicKey()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not receive transfer")
	}
	_, balanceThirdAuthor := state.Wallets.BalanceHash(crypto.HashToken(thirdAuthor.Wallet.PublicKey()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 4
	block = NewBlock(crypto.Hasher([]byte{}), 3, 4, blockFormationToken.PublicKey(), validator)

	// Join audience sent by second member
	joinAudience := secondAuthor.NewJoinAudience(audienceTest.PrivateKey.PublicKey(), "first audience member", 4, uint64(joinFee))
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

	// COMO CHECAR SPONSOR OFFER NO ESTADO (???)
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spent on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.BalanceHash(crypto.HashToken(secondAuthor.Wallet.PublicKey()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not receive transfer")
	}
	_, balanceThirdAuthor = state.Wallets.BalanceHash(crypto.HashToken(thirdAuthor.Wallet.PublicKey()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not receive transfer")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 5
	block = NewBlock(crypto.Hasher([]byte{}), 4, 5, blockFormationToken.PublicKey(), validator)

	// Accept join audience
	acceptJoin := firstAuthor.NewAcceptJoinAudience(audienceTest, secondAuthor.PrivateKey.PublicKey(), secondAuthor.PrivateKey.PublicKey(), 2, 5, uint64(joinFee))
	if !block.Incorporate(acceptJoin) {
		t.Error("could not accept join request to audience")
	}
	firstBalance = firstBalance - joinFee
	count = count + 1

	// Accept sponsor offer
	sponsordAccept := firstAuthor.NewSponsorshipAcceptance(audienceTest, sponsorOffer, 5, uint64(joinFee))
	if !block.Incorporate(sponsordAccept) {
		t.Error("could not accept sponsorship acceptance")
	}
	firstBalance = firstBalance - joinFee
	thirdBalance = thirdBalance - 20
	count = count + 1

	state.IncorporateBlock(block)
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spend on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.BalanceHash(crypto.HashToken(secondAuthor.Wallet.PublicKey()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not spend on instructions")
	}
	_, balanceThirdAuthor = state.Wallets.BalanceHash(crypto.HashToken(thirdAuthor.Wallet.PublicKey()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not spend on instructions")
	}
	_, balanceAudience := state.Wallets.BalanceHash(crypto.HashToken(audienceTest.PrivateKey.PublicKey()))
	if balanceAudience != 20 {
		t.Error("audience did not receive revenue from sponsor")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}

	// BLOCK 6
	block = NewBlock(crypto.Hasher([]byte{}), 5, 6, blockFormationToken.PublicKey(), validator)

	// React to content sent by member2
	react := secondAuthor.NewReact([]byte("teste"), 1, 6, uint64(joinFee))
	if !block.Incorporate(react) {
		t.Error("could not accept react to content")
	}
	secondBalance = secondBalance - joinFee
	count = count + 1

	// Update audience keys by member1
	readers := make(map[crypto.Token]crypto.Token, 3)
	for n := 0; n < 3; n++ {
		token, _ := crypto.RandomAsymetricKey()
		readers[token] = token
	}
	updateAudience := firstAuthor.NewUpdateAudience(audienceTest, readers, readers, readers, 1, "removing member2 from audience", 6, uint64(joinFee))
	if !block.Incorporate(updateAudience) {
		t.Error("could not accept update audience instruction")
	}
	firstBalance = firstBalance - joinFee
	count = count + 1

	// Create Ephemeral token by member 2
	pubEph, prvEph := crypto.RandomAsymetricKey()
	ephemeralAuthor := &instructions.Author{PrivateKey: prvEph, Wallet: token} // ephemeral token using eve wallet
	ephemeral := secondAuthor.NewCreateEphemeral(pubEph, 20, 6, uint64(joinFee))
	if !block.Incorporate(ephemeral) {
		t.Error("could not accept create ephemeral token instruction")
	}
	secondBalance = secondBalance - joinFee
	count = count + 1

	state.IncorporateBlock(block)
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spend on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.BalanceHash(crypto.HashToken(secondAuthor.Wallet.PublicKey()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not spend on instructions")
	}
	_, balanceThirdAuthor = state.Wallets.BalanceHash(crypto.HashToken(thirdAuthor.Wallet.PublicKey()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not spend on instructions")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}
	if epoch := state.EphemeralTokens.Exists(crypto.HashToken(pubEph)); epoch != 20 {
		t.Error("ephemeral token not incorporated")
	}

	// BLOCK 7
	block = NewBlock(crypto.Hasher([]byte{}), 6, 7, blockFormationToken.PublicKey(), validator)

	// Secure Channel by member 2
	secure := ephemeralAuthor.NewSecureChannel([]byte("teste"), uint64(1), []byte("encryptedNonce"), []byte("content"), 7, uint64(joinFee))
	if !block.Incorporate(secure) {
		t.Error("could not accept secure channel instruction")
	}
	eveBalance = eveBalance - float64(joinFee)
	count = count + 1

	state.IncorporateBlock(block)
	_, balanceFirstAuthor = state.Wallets.BalanceHash(crypto.HashToken(firstAuthor.Wallet.PublicKey()))
	if balanceFirstAuthor != uint64(firstBalance) {
		t.Error("first author did not spend on instructions")
	}
	_, balanceSecondAuthor = state.Wallets.BalanceHash(crypto.HashToken(secondAuthor.Wallet.PublicKey()))
	if balanceSecondAuthor != uint64(secondBalance) {
		t.Error("second author did not spend on instructions")
	}
	_, balanceThirdAuthor = state.Wallets.BalanceHash(crypto.HashToken(thirdAuthor.Wallet.PublicKey()))
	if balanceThirdAuthor != uint64(thirdBalance) {
		t.Error("third author did not spend on instructions")
	}
	_, balanceBlockFormator = state.Wallets.BalanceHash(crypto.HashToken(blockFormationToken.PublicKey()))
	if balanceBlockFormator != uint64(count*joinFee) {
		t.Error("block formator has not received fee for processed instructions")
	}
}
