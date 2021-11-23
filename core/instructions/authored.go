package instructions

import (
	"fmt"

	"github.com/Aereum/aereum/core/crypto"
)

type AuthoredInstruction struct {
	Version         byte
	InstructionType byte
	epoch           uint64
	Author          []byte
	Message         KindSerializer
	Wallet          []byte
	Fee             uint64
	Attorney        []byte
	Signature       []byte
	WalletSignature []byte
}

func (a *AuthoredInstruction) Validate(validator Validator) bool {
	return a.Message.Validate(validator)
}

func NewAuthoredInstruction(author crypto.PrivateKey, instruction KindSerializer,
	epoch uint64, fee uint64, attorney, wallet *crypto.PrivateKey) Instruction {

	newInstruction := AuthoredInstruction{
		Version:         0,
		InstructionType: instruction.Kind(),
		epoch:           epoch,
		Author:          author.PublicKey().ToBytes(),
		Message:         instruction,
		Fee:             fee,
	}

	if (*attorney).IsValid() {
		newInstruction.Attorney = (*attorney).PublicKey().ToBytes()
		if wallet != nil {
			newInstruction.sign(*attorney, *wallet)
		} else {
			newInstruction.sign(*attorney, author)
		}
	} else {
		newInstruction.Attorney = []byte{}
		if wallet != nil {
			newInstruction.sign(author, *wallet)
		} else {
			bytes := newInstruction.serializeWithoutSignatures()
			newInstruction.Signature, _ = author.Sign(bytes)
		}
	}
	return &newInstruction
}

func (a *AuthoredInstruction) serializeWithoutSignatures() []byte {
	bytes := []byte{0, a.InstructionType}
	PutUint64(a.epoch, &bytes)
	PutByteArray(a.Author, &bytes)
	PutByteArray(a.Message.Serialize(), &bytes)
	PutByteArray(a.Wallet, &bytes)
	PutUint64(a.Fee, &bytes)
	PutByteArray(a.Attorney, &bytes)
	return bytes
}

func (a *AuthoredInstruction) Serialize() []byte {
	bytes := a.serializeWithoutSignatures()
	PutByteArray(a.Signature, &bytes)
	PutByteArray(a.WalletSignature, &bytes)
	return bytes
}

func (a *AuthoredInstruction) Kind() byte {
	return a.InstructionType
}

func (a *AuthoredInstruction) Epoch() uint64 {
	return a.epoch
}

func (a *AuthoredInstruction) sign(author, wallet crypto.PrivateKey) bool {
	bytes := a.serializeWithoutSignatures()
	signAuthor, err := author.Sign(bytes)
	if err != nil {
		return false
	}
	PutByteArray(signAuthor, &bytes)
	signWallet, errWallet := wallet.Sign(bytes)
	if errWallet != nil {
		return false
	}
	a.Signature = signAuthor
	a.WalletSignature = signWallet
	return true
}

func ParseAuthoredInstruction(data []byte) (*AuthoredInstruction, error) {
	if data[0] != 0 {
		return nil, fmt.Errorf("wrong instruction version")
	}
	if data[1] >= IUnkown || data[1] <= IJoinNetwork {
		return nil, fmt.Errorf("wrong instruction type")
	}
	length := len(data)
	var msg AuthoredInstruction
	msg.InstructionType = data[1]
	position := 2
	msg.epoch, position = ParseUint64(data, position)
	msg.Author, position = ParseByteArray(data, position)
	var bytes []byte
	bytes, position = ParseByteArray(data, position)
	switch msg.InstructionType {
	case IJoinNetwork:
		msg.Message = ParseJoinNetwork(bytes)
	case IUpdateInfo:
		msg.Message = ParseJoinNetwork(bytes)
	case ICreateAudience:
		msg.Message = ParseCreateAudience(bytes)
	case IJoinAudience:
		msg.Message = ParseJoinAudience(bytes)
	case IAcceptJoinRequest:
		msg.Message = ParseAcceptJoinAudience(bytes)
	case IContent:
		msg.Message = ParseContent(bytes)
	case IUpdateAudience:
		msg.Message = ParseUpdateAudience(bytes)
	case IGrantPowerOfAttorney:
		msg.Message = ParseGrantPowerOfAttorney(bytes)
	case IRevokePowerOfAttorney:
		msg.Message = ParseRevokePowerOfAttorney(bytes)
	case ISponsorshipOffer:
		msg.Message = ParseSponsorshipOffer(bytes)
	case ISponsorshipAcceptance:
		msg.Message = ParseSponsorshipAcceptance(bytes)
	case ICreateEphemeral:
		msg.Message = ParseCreateEphemeral(bytes)
	case ISecureChannel:
		msg.Message = ParseSecureChannel(bytes)
	case IReact:
		msg.Message = ParseReact(bytes)
	}
	if msg.Message == nil {
		return nil, fmt.Errorf("Could not parse instruction.")
	}
	msg.Wallet, position = ParseByteArray(data, position)
	msg.Fee, position = ParseUint64(data, position)
	msg.Attorney, position = ParseByteArray(data, position)
	// check author or power of attorney signature
	if position-1 > length {
		return nil, ErrCouldNotParseMessage
	}
	msgToVerify := data[0:position]
	msg.Signature, position = ParseByteArray(data, position)
	token := msg.Author
	if len(msg.Attorney) > 0 {
		token = msg.Attorney
	}
	if publicKey, err := crypto.PublicKeyFromBytes(token); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, msg.Signature) {
			return nil, ErrInvalidSignature
		}
	}

	// check wallet signature
	if position-1 > length {
		return nil, ErrCouldNotParseMessage
	}
	msgToVerify = data[0:position]
	msg.WalletSignature, position = ParseByteArray(data, position)
	if position != length {
		return nil, ErrCouldNotParseMessage
	}
	if publicKey, err := crypto.PublicKeyFromBytes(msg.Wallet); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(msgToVerify, msg.WalletSignature) {
			return nil, ErrInvalidSignature
		}
	}
	return &msg, nil
}
