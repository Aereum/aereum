package instructionsnew

import (
	"reflect"
	"testing"
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
		t.Error("Parse and Serialize not working for AcceptJoinAudience")
	}
}
