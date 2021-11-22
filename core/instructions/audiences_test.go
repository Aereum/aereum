package instructions

import (
	"reflect"
	"testing"

	"github.com/Aereum/aereum/core/crypto"
)

func TestCreateAudience(t *testing.T) {
	
	// Creating 3 public-private keys for audience identification, submission and moderation
	public_audience, private_audience := crypto.RandomAsymetricKey()
	public_submission, private_submission := crypto.RandomAsymetricKey()
	public_moderation, private_moderation := crypto.RandomAsymetricKey()
	
	cipher_key_audience := crypto.NewCipherKey()
	cipher_key_submission := crypto.NewCipherKey()
	cipher_key_moderation := crypto.NewCipherKey()
	
	cipher_submission := crypto.CipherNonceFromKey(cipher_key_submission)
	cipher_moderation := crypto.CipherNonceFromKey(cipher_key_moderation)

	audience_cipher_submission.Seal(private_submission)

	message := &CreateAudience{
		Audience:    public_audience.ToBytes(),
		Submission:  public_submission.ToBytes(),
		Moderation:  public_moderation.ToBytes(),
		AudienceKey: 

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

func Test(t *testing.T) {
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
