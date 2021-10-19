package message

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)

func IsTransfer(msg []byte) bool {
	return len(msg) > 0 && msg[0] == TransferMsg
}

func NewTransfer(wallet crypto.PrivateKey, to crypto.Hash, value uint64, fee uint64, epoch uint64) *Transfer {
	t := &Transfer{
		MessageType: TransferMsg,
		Epoch:       epoch,
		From:        wallet.PublicKey().ToBytes(),
		To:          to[:],
		Value:       value,
		Fee:         fee,
	}
	hashed := crypto.Hasher(t.serializeWithouSignature())
	var err error
	t.Signature, err = wallet.Sign(hashed[:])
	if err != nil {
		return nil
	}
	return t
}

type Transfer struct {
	MessageType byte
	Epoch       uint64
	From        []byte
	To          []byte
	Value       uint64
	Fee         uint64
	Signature   []byte
}

func (t *Transfer) Payments() Payment {
	return Payment{
		DebitAcc:    []crypto.Hash{crypto.Hasher(t.From)},
		DebitValue:  []uint64{t.Value + t.Fee},
		CreditAcc:   []crypto.Hash{crypto.Hasher(t.To)},
		CreditValue: []uint64{t.Value},
	}
}

func (t *Transfer) serializeWithouSignature() []byte {
	bytes := []byte{TransferMsg}
	PutUint64(t.Epoch, &bytes)
	PutByteArray(t.From, &bytes)
	PutByteArray(t.To, &bytes)
	PutUint64(t.Value, &bytes)
	PutUint64(t.Fee, &bytes)
	return bytes
}

func (t *Transfer) Sign(privateKey crypto.PrivateKey) bool {
	bytes := t.serializeWithouSignature()
	sign, err := privateKey.Sign(bytes)
	if err != nil {
		return false
	}
	t.Signature = sign
	return true
}

func (t *Transfer) Serialize() []byte {
	bytes := t.serializeWithouSignature()
	PutByteArray(t.Signature, &bytes)
	return bytes
}

// tries to parse and verifies signature.
func ParseTranfer(data []byte) (*Transfer, error) {
	if len(data) == 0 || data[0] != TransferMsg {
		return nil, errors.New("wrong message type")
	}
	length := len(data)
	var msg Transfer
	position := 1
	msg.Epoch, position = ParseUint64(data, position)
	msg.From, position = ParseByteArray(data, position)
	msg.To, position = ParseByteArray(data, position)
	msg.Value, position = ParseUint64(data, position)
	msg.Fee, position = ParseUint64(data, position)
	if position >= length {
		return nil, ErrCouldNotParseMessage
	}
	hashed := crypto.Hasher(data[0:position])
	msg.Signature, position = ParseByteArray(data, position)
	if position-1 > length || len(msg.Signature) == 0 {
		return nil, ErrCouldNotParseMessage
	}
	// check signature
	if publicKey, err := crypto.PublicKeyFromBytes(msg.From); err != nil {
		return nil, ErrCouldNotParseSignature
	} else {
		if !publicKey.Verify(hashed[:], msg.Signature) {
			return nil, ErrInvalidSignature
		}
	}
	return &msg, nil
}
