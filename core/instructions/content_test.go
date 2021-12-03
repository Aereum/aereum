package instructions

import (
	"fmt"
	"reflect"
	"testing"
)

func TestContent(t *testing.T) {

	content := author.NewContent(audienceTest, "test", []byte("content of the content instruction"), true, true, 10, 2000)
	fmt.Printf("%+v\n", *content)
	content2 := ParseContent(content.Serialize())
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
	reaction := author.NewReact([]byte{1, 2, 3}, 2, 10, 2000)
	reaction2 := ParseReact(reaction.Serialize())
	if reaction2 == nil {
		t.Error("could not parse Reaction")
		return
	}
	if !reflect.DeepEqual(reaction, reaction2) {
		t.Error("Parse and Serialize not working for Reaction")
	}
}
