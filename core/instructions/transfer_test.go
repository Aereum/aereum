package instructions

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestTransfer(t *testing.T) {

	token, tokenPrivate := crypto.RandomAsymetricKey()

	token_r0, _ := crypto.RandomAsymetricKey()
	token_r1, _ := crypto.RandomAsymetricKey()
	recipients := make([]Recipient, 2)
	recipients[0] = Recipient{Token: token_r0.ToBytes(), Value: 10}
	recipients[1] = Recipient{Token: token_r1.ToBytes(), Value: 100}

	message := &Transfer{
		Version:         0,
		InstructionType: 0,
		epoch:           10928298,
		From:            token.ToBytes(),
		To:              recipients,
		Reason:          "Sponsored content recieved",
		Fee:             10,
	}
	hash := crypto.Hasher(message.serializeWithoutSignature())
	message.Signature, _ = tokenPrivate.Sign(hash[:])
	copy := ParseTransfer(message.Serialize())
	if copy == nil {
		t.Error("Could not Transfer.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for Transfer messages.")
	}
}

func TestDeposit(t *testing.T) {

	token, _ := crypto.RandomAsymetricKey()
	message := &Deposit{
		Version:         0,
		InstructionType: 0,
		epoch:           10928298,
		Token:           token.ToBytes(),
		Value:           10,
		Reason:          "Sponsored content recieved",
		Fee:             10,
		Signature:       []byte{1, 2, 3},
	}

	deposit := message.Serialize()
	deposit2 := ParseDeposit(deposit)
	if deposit2 == nil {
		t.Error("Could not Deposit.")
		return
	}
	if ok := reflect.DeepEqual(*message, *deposit2); !ok {
		t.Error("Parse and Serialization not working for Deposit messages.")
	}

}
