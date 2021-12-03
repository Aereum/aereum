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
	jsonString1 = `{"teste":1}`
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

func generalTesting(t testing.T) {
	// member1 := author1.NewJoinNetwork("member1", jsonString1, 10, 20)
	// member2 := author2.NewJoinNetwork("member2", jsonString2, 11, 20)

}
