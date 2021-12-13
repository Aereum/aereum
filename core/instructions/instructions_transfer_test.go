package instructions

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestTransfer(t *testing.T) {

	_, from := crypto.RandomAsymetricKey()
	to, _ := crypto.RandomAsymetricKey()

	message := NewSingleReciepientTransfer(from, to, "whatever", 10, 10, 10)
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

	_, from := crypto.RandomAsymetricKey()
	message := NewDeposit(from, 10, 10, 10)
	deposit2 := ParseDeposit(message.Serialize())
	if deposit2 == nil {
		t.Error("Could not Deposit.")
		return
	}
	if ok := reflect.DeepEqual(*message, *deposit2); !ok {
		t.Error("Parse and Serialization not working for Deposit messages.")
	}

}

func TestWithdraw(t *testing.T) {

	_, to := crypto.RandomAsymetricKey()
	withdraw := NewWithdraw(to, 10, 10, 10)
	withdraw2 := ParseWithdraw(withdraw.Serialize())
	if withdraw2 == nil {
		t.Error("Could not Withdraw.")
		return
	}
	if ok := reflect.DeepEqual(*withdraw, *withdraw2); !ok {
		t.Error("Parse and Serialization not working for Withdraw messages.")
	}

}
