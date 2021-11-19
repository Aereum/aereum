package instructions

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestJoinNetwork(t *testing.T) {
	message := &JoinNetwork{
		Caption: "larissa",
		Details: `
{"Name": "Larissa", "Static Content": {"/img/": "https://www.*******.com/fotos/"}}`,
	}
	bytes := message.Serialize()
	copy := ParseJoinNetwork(bytes)
	if copy == nil {
		t.Error("Could not ParseJoinNetwork")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for JoinNetwork messages")
	}
}

func TestUpdateInfo(t *testing.T) {
	message := &UpdateInfo{
		Details: `
{"Name": "Larissa2", "Static Content": {"/img/": "https://www.*******.com/fotos2/"}}
`,
	}
	bytes := message.Serialize()
	copy := ParseUpdateInfo(bytes)
	if copy == nil {
		t.Error("Could not ParseUpdateInfo")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for UpdateInfo messages")
	}
}

func TestGrantPowerOfAttorney(t *testing.T) {
	token, _ := crypto.RandomAsymetricKey()
	message := &GrantPowerOfAttorney{
		Attorney: token.ToBytes(),
	}
	bytes := message.Serialize()
	copy := ParseGrantPowerOfAttorney(bytes)
	if copy == nil {
		t.Error("Could not ParseGrantPowerOfAttorney.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for GrantPowerOfAttorney messages.")
	}
}

func TestRevokePowerOfAttorney(t *testing.T) {
	token, _ := crypto.RandomAsymetricKey()
	message := &RevokePowerOfAttorney{
		Attorney: token.ToBytes(),
	}
	bytes := message.Serialize()
	copy := ParseRevokePowerOfAttorney(bytes)
	if copy == nil {
		t.Error("Could not RevokeGrantPowerOfAttorney.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for RevokePowerOfAttorney messages.")
	}
}

func TestCreateEphemeral(t *testing.T) {
	token, _ := crypto.RandomAsymetricKey()
	message := &CreateEphemeral{
		EphemeralToken: token.ToBytes(),
		Expiry:         25,
	}
	bytes := message.Serialize()
	copy := ParseCreateEphemeral(bytes)
	if copy == nil {
		t.Error("Could not CreateEphemeral.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for CreateEphemeral messages.")
	}
}

func TestSecureChannel(t *testing.T) {
	message := &SecureChannel{
		TokenRange:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Nonce:          121,
		EncryptedNonce: []byte{5, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Content:        []byte{5, 2, 3, 4, 5, 6, 7, 8, 9, 12},
	}
	bytes := message.Serialize()
	copy := ParseSecureChannel(bytes)
	if copy == nil {
		t.Error("Could not ParseSecureChannel.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for SecureChannel messages.")
	}
}
