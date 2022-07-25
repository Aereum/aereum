package instructions

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
)

var (
	audienceTest *Stage = NewStage(0, "teste")
)

func TestCreateteAudience(t *testing.T) {
	audience := author.NewCreateAudience(audienceTest, 10, 2000)
	audience2 := ParseCreateStage(audience.Serialize())
	if audience2 == nil {
		t.Error("could not parse CreateAudience")
		return
	}
	if !reflect.DeepEqual(audience, audience2) {
		t.Error("Parse and Serialize not working for CreateAudience")
	}
}

func TestJoinAudience(t *testing.T) {
	join := author.NewJoinAudience(audienceTest.PrivateKey.PublicKey(), audienceTest.PrivateKey.PublicKey(), "teste", 10, 2000)
	join2 := ParseJoinStage(join.Serialize())
	if join2 == nil {
		t.Error("could not parse JoinAudience")
		return
	}
	if !reflect.DeepEqual(join, join2) {
		t.Error("Parse and Serialize not working for JoinAudience")
	}
}

func TestAcceptJoinAudience(t *testing.T) {
	_, key := dh.NewEphemeralKey()
	accept := author.NewAcceptJoinAudience(audienceTest, author.PrivateKey.PublicKey(), key, 2, 10, 2000)
	accept2 := ParseAcceptJoinStage(accept.Serialize())
	if accept2 == nil {
		t.Error("could not parse AcceptJoinAudience")
		return
	}
	if !reflect.DeepEqual(accept, accept2) {
		t.Error("Parse and Serialize not working for AcceptJoinAudience")
	}
}

func TestUpdateAudience(t *testing.T) {
	readers := make(map[crypto.Token]crypto.Token, 3)
	for n := 0; n < 3; n++ {
		token, _ := crypto.RandomAsymetricKey()
		_, key := dh.NewEphemeralKey()
		readers[token] = key
	}
	update := author.NewUpdateAudience(audienceTest, readers, readers, readers, 2, "teste", 10, 2000)
	fmt.Printf("%+v", *update)
	update2 := ParseUpdateStage(update.Serialize())
	// fmt.Printf(string(update2.audience))
	if update2 == nil {
		t.Error("could not parse UpdateAudience")
		return
	}
	if !reflect.DeepEqual(update, update2) {
		t.Error("Parse and Serialize not working for UpdateAudience")
	}
}
