package instructions

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestCreateAudience(t *testing.T) {

	// Creating 3 public-private keys for audience identification, submission and moderation

	public_audience, _ := crypto.RandomAsymetricKey()
	public_submission, _ := crypto.RandomAsymetricKey()
	public_moderation, _ := crypto.RandomAsymetricKey()
	message := &CreateAudience{
		Audience:      public_audience.ToBytes(),
		Submission:    public_submission.ToBytes(),
		Moderation:    public_moderation.ToBytes(),
		AudienceKey:   []byte("teste"),
		SubmissionKey: []byte("teste"),
		ModerationKey: []byte("teste"),
		Flag:          byte(0),
		Description:   "Very first audience",
	}
	bytes := message.Serialize()
	copy := ParseCreateAudience(bytes)
	if copy == nil {
		t.Error("Could not ParseCreateAudience")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for CreateAudience messages")
	}
}

func TestJoinAudience(t *testing.T) {
	public_audience, _ := crypto.RandomAsymetricKey()
	message := &JoinAudience{
		Audience:     public_audience.ToBytes(),
		Presentation: "New member for existing audience",
	}
	bytes := message.Serialize()
	copy := ParseJoinAudience(bytes)
	if copy == nil {
		t.Error("Could not ParseJoinAudience.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for JoinAudience messages.")
	}
}

func TestAcceptJoinAudience(t *testing.T) {
	public_audience, _ := crypto.RandomAsymetricKey()
	public_member, _ := crypto.RandomAsymetricKey()
	message := &AcceptJoinAudience{
		Audience: public_audience.ToBytes(),
		Member:   public_member.ToBytes(),
		Read:     []byte("teste"),
		Submit:   []byte("teste"),
		Moderate: []byte("teste"),
	}
	bytes := message.Serialize()
	copy := ParseAcceptJoinAudience(bytes)
	if copy == nil {
		t.Error("Could not ParseAcceptJoinAudience.")
		return
	}
	if ok := reflect.DeepEqual(*message, *copy); !ok {
		t.Error("Parse and Serialization not working for AcceptJoinAudience messages.")
	}
}
