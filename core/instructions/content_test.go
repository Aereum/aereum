package instructions

import (
	"reflect"
	"testing"
)

func TestContent(t *testing.T) {
	byte_message := make([]byte, 0)
	PutString("content of the content instruction", &byte_message)
	content := author.NewContent(audienceTest, "test", byte_message, byte_message, true, true, 10, 2000)
	content2 := ParseCreateAudience(content.Serialize())
	if content2 == nil {
		t.Error("could not parse Content")
		return
	}
	if !reflect.DeepEqual(content, content2) {
		t.Error("Parse and Serialize not working for Content")
	}
}

func TestReact(t *testing.T) {
	byte_message := make([]byte, 0)
	PutString("this should be a hash", &byte_message)
	reaction := author.NewReact(audienceTest, byte_message, 0, 10, 2000)
	reaction2 := ParseCreateAudience(reaction.Serialize())
	if reaction2 == nil {
		t.Error("could not parse Reaction")
		return
	}
	if !reflect.DeepEqual(reaction, reaction2) {
		t.Error("Parse and Serialize not working for Reaction")
	}
}
