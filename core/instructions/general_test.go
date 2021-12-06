package instructions

import (
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

var (
	_, token1    = crypto.RandomAsymetricKey()
	_, attorney1 = crypto.RandomAsymetricKey()
	_, wallet1   = crypto.RandomAsymetricKey()
	author1      = &Author{
		token:    &token,
		attorney: &attorney,
		wallet:   &wallet,
	}
	jsonString1     = `{"teste":1}`
	jsonString1_new = `{"update":1}`
)

var (
	_, token2    = crypto.RandomAsymetricKey()
	_, attorney2 = crypto.RandomAsymetricKey()
	_, wallet2   = crypto.RandomAsymetricKey()
	author2      = &Author{
		token:    &token,
		attorney: &attorney,
		wallet:   &wallet,
	}
	jsonString2 = `{"teste":1}`
)

var (
	state1 = State()

	validator1 = Validator()
)

func TestGeneral(t *testing.T) {

	hash := crypto.Hash()
	checkpoint := 10
	publisher_token, _ := crypto.RandomAsymetricKey()

	block := NewBlock(hash, checkpoint, 10, publisher_token.ToBytes(), validator)

	author1.NewJoinNetwork("member1", jsonString1, 10, 20)
	author2.NewJoinNetwork("member2", jsonString2, 11, 20)

	// fmt.Print(author1.token.ToBytes())
	author1.NewUpdateInfo(jsonString1_new, 12, 10)

	recipients := make([]Recipient, 1)
	recipients[0] = Recipient{Token: author2.token.ToBytes(), Value: 10}

	transfer := &Transfer{
		Version:         0,
		InstructionType: 0,
		epoch:           20,
		From:            author1.token.ToBytes(),
		To:              recipients,
		Reason:          "Sponsored content recieved",
		Fee:             10,
	}

}
