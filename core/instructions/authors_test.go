package instructions

import (
	"reflect"
	"testing"
)

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
		t.Error("Parse and Serialization not working for message with power of attorney.")
	}
}
