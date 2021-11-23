package instructions

/*
func TestCreateAudience(t *testing.T) {

	// Creating 3 public-private keys for audience identification, submission and moderation
	public_audience, private_audience := crypto.RandomAsymetricKey()
	public_submission, private_submission := crypto.RandomAsymetricKey()
	public_moderation, private_moderation := crypto.RandomAsymetricKey()

	cipher_key_audience := crypto.NewCipherKey()
	cipher_key_submission := crypto.NewCipherKey()
	cipher_key_moderation := crypto.NewCipherKey()

	cipher_audience := crypto.CipherNonceFromKey(cipher_key_audience)
	cipher_submission := crypto.CipherNonceFromKey(cipher_key_submission)
	cipher_moderation := crypto.CipherNonceFromKey(cipher_key_moderation)

	ciphered_aud := cipher_audience.Seal(private_audience.ToBytes())
	ciphered_sub := cipher_submission.Seal(private_submission.ToBytes())
	ciphered_mod := cipher_moderation.Seal(private_moderation.ToBytes())

	message := &CreateAudience{
		Audience:      public_audience.ToBytes(),
		Submission:    public_submission.ToBytes(),
		Moderation:    public_moderation.ToBytes(),
		AudienceKey:   ciphered_aud,
		SubmissionKey: ciphered_sub,
		ModerationKey: ciphered_mod,
		Flag:          0,
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
	public_read, _ := crypto.RandomAsymetricKey()
	message := &AcceptJoinAudience{
		Audience: public_audience.ToBytes(),
		Member:   public_member.ToBytes(),
		Read:	
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
*/
