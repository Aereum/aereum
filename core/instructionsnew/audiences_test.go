package instructionsnew

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

var (
	audienceTest *Audience = NewAudience()
)

func TestCreateteAudience(t *testing.T) {
	audience := author.NewCreateAudience(audienceTest, 1, "teste", 10, 2000)
	audience2 := ParseCreateAudience(audience.Serialize())
	if audience2 == nil {
		t.Error("could not parse CreateAudience")
		return
	}
	if !reflect.DeepEqual(audience, audience2) {
		t.Error("Parse and Serialize not working for CreateAudience")
	}
}

func TestJoinAudience(t *testing.T) {
	join := author.NewJoinAudience(audienceTest.token.PublicKey().ToBytes(), "teste", 10, 2000)
	join2 := ParseJoinAudience(join.Serialize())
	if join2 == nil {
		t.Error("could not parse JoinAudience")
		return
	}
	if !reflect.DeepEqual(join, join2) {
		t.Error("Parse and Serialize not working for JoinAudience")
	}
}

func TestAcceptJoinAudience(t *testing.T) {
	accept := author.NewAcceptJoinAudience(audienceTest, author.token.PublicKey(), 3, 10, 2000)
	accept2 := ParseAcceptJoinAudience(accept.Serialize())
	if accept2 == nil {
		t.Error("could not parse AcceptJoinAudience")
		return
	}
	if !reflect.DeepEqual(accept, accept2) {
		fmt.Println(*accept)
		fmt.Println(*accept2)
		t.Error("Parse and Serialize not working for AcceptJoinAudience")
	}
}

func TestUpdateAudience(t *testing.T) {
	readers := make([]crypto.PublicKey, 3)
	for n := 0; n < 3; n++ {
		readers[n], _ = crypto.RandomAsymetricKey()
	}
	update := author.NewUpdateAudience(audienceTest, readers, readers, readers, 2, "teste", 10, 2000)
	update2 := ParseUpdateAudience(update.Serialize())
	if update2 == nil {
		t.Error("could not parse UpdateAudience")
		return
	}
	if !reflect.DeepEqual(update, update2) {
		t.Error("Parse and Serialize not working for UpdateAudience")
	}
}
