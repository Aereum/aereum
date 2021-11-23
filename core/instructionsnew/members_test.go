package instructionsnew

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

var (
	_, token    = crypto.RandomAsymetricKey()
	_, attorney = crypto.RandomAsymetricKey()
	_, wallet   = crypto.RandomAsymetricKey()
	author      = &Author{
		token:    &token,
		attorney: &attorney,
		wallet:   &wallet,
	}
	jsonString = `{"teste":1}`
)

func TestJoinNetwork(t *testing.T) {
	network := author.NewJoinNetwork("teste", jsonString, 10, 2000)
	network2 := ParseJoinNetwork(network.Serialize())
	if network2 == nil {
		t.Error("could not parse JoinNetwork")
	}
	if !reflect.DeepEqual(network, network2) {
		t.Error("Parse and Serialize not working for JoinNetwork")
	}
}
