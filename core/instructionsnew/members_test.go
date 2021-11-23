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
		return
	}
	if !reflect.DeepEqual(network, network2) {
		t.Error("Parse and Serialize not working for JoinNetwork")
	}
}

func TestUpdateInfo(t *testing.T) {
	update := author.NewUpdateInfo(jsonString, 10, 2000)
	update2 := ParseUpdateInfo(update.Serialize())
	if update2 == nil {
		t.Error("could not parse UpdateInfo")
		return
	}
	if !reflect.DeepEqual(update, update2) {
		t.Error("Parse and Serialize not working for UpdateInfo")
	}
}

func TestGrantPowerOfAttorney(t *testing.T) {
	grant := author.NewGrantPowerOfAttorney([]byte{1, 2, 3, 6, 7, 8}, 10, 2000)
	grant2 := ParseGrantPowerOfAttorney(grant.Serialize())
	if grant2 == nil {
		t.Error("could not parse GrantPowerOfAttorney")
		return
	}
	if !reflect.DeepEqual(grant, grant2) {
		t.Error("Parse and Serialize not working for GrantPowerOfAttorney")
	}
}

func TestRevokePowerOfAttorney(t *testing.T) {
	revoke := author.NewRevokePowerOfAttorney([]byte{1, 2, 3, 6, 7, 8}, 10, 2000)
	revoke2 := ParseRevokePowerOfAttorney(revoke.Serialize())
	if revoke2 == nil {
		t.Error("could not parse RevokePowerOfAttorney")
		return
	}
	if !reflect.DeepEqual(revoke, revoke2) {
		t.Error("Parse and Serialize not working for RevokePowerOfAttorney")
	}
}

func TestCreateEphemeral(t *testing.T) {
	ephemeral := author.NewCreateEphemeral([]byte{1, 2, 3, 6, 7, 8}, 20, 10, 2000)
	ephemeral2 := ParseCreateEphemeral(ephemeral.Serialize())
	if ephemeral2 == nil {
		t.Error("could not parse CreateEphemeral")
		return
	}
	if !reflect.DeepEqual(ephemeral, ephemeral2) {
		t.Error("Parse and Serialize not working for CreateEphemeral")
	}
}

func TestSecureChannel(t *testing.T) {
	secure := author.NewSecureChannel([]byte{1, 2, 3, 6, 7, 8}, 20, []byte{1, 2, 3, 6, 7, 28}, []byte{1, 2, 3, 6, 7, 9, 10}, 10, 2000)
	secure2 := ParseSecureChannel(secure.Serialize())
	if secure2 == nil {
		t.Error("could not parse SecureChannel")
		return
	}
	if !reflect.DeepEqual(secure, secure2) {
		t.Error("Parse and Serialize not working for SecureChannel")
	}
}
