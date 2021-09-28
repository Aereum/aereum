package message

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

const (
	genesis byte = 0
	// version 0
	transferMsg
	subscribeMsg
	aboutMsg
	createAudienceMsg
	joinAudienceMsg
	audienceChangeMsg
	advertisingOfferMsg
	contentMsg
	grantPowerOfAttorneyMsg
	revokePowerOfAttorneyMsg
	unkownMessageType
	// to be used in other versions
)

type Transfer struct {
	From      []byte
	To        []byte
	Value     uint64
	Epoch     uint64
	Signature []byte
}

func (t *Transfer) Serialize() []byte {
	bytes := &[]byte{transferMsg}
	PutByteArray(t.From, bytes)
	PutByteArray(t.To, bytes)
	PutUint64(t.Value, bytes)
	PutUint64(t.Epoch, bytes)
	PutByteArray(t.Signature, bytes)
	return *bytes
}

func ParseTranfer(data []byte) *Transfer {
	if len(data) == 0 || data[0] != transferMsg {
		return nil
	}
	length := len(data)
	var msg Transfer
	position := 1
	msg.From, position = ParseByteArray(data, position)
	msg.To, position = ParseByteArray(data, position)
	msg.Value, position = ParseUint64(data, position)
	msg.Epoch, position = ParseUint64(data, position)
	if position >= length {
		return nil
	}
	hashed := sha256.Sum256(data[0 : position-1])
	msg.Signature, position = ParseByteArray(data, position)
	if position-1 > length || len(msg.Signature) == 0 {
		return nil
	}
	// check signature
	if publicKey, err := x509.ParsePKCS1PublicKey(msg.From); err != nil {
		return nil
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.Signature) != nil {
			return nil
		}
	}
	return &msg
}

type Message struct {
	MessageType     byte
	Author          []byte
	Message         []byte
	FeeWallet       []byte
	FeeValue        uint64
	Epoch           uint64
	PowerOfAttorney []byte
	Signature       []byte
	WalletSignature []byte
}

func (m *Message) Serialize() []byte {
	bytes := &[]byte{m.MessageType}
	PutByteArray(m.Author, bytes)
	PutByteArray(m.Message, bytes)
	PutByteArray(m.FeeWallet, bytes)
	PutUint64(m.FeeValue, bytes)
	PutByteArray(m.PowerOfAttorney, bytes)
	PutByteArray(m.Signature, bytes)
	PutByteArray(m.WalletSignature, bytes)
	return *bytes
}

func ParseMessage(data []byte) *Message {
	if data[0] >= unkownMessageType || data[0] <= transferMsg {
		return nil
	}
	length := len(data)
	var msg Message
	msg.MessageType = data[0]
	position := 1
	msg.Author, position = ParseByteArray(data, position)
	msg.Message, position = ParseByteArray(data, position)
	msg.FeeWallet, position = ParseByteArray(data, position)
	msg.FeeValue, position = ParseUint64(data, position)
	msg.Epoch, position = ParseUint64(data, position)
	msg.PowerOfAttorney, position = ParseByteArray(data, position)
	if position-1 > length {
		return nil
	}
	hashed := sha256.Sum256(data[0 : position-1])
	msg.Signature, position = ParseByteArray(data, position)
	if position-1 > length {
		return nil
	}
	// check author or power of attorney signature
	token := msg.Author
	if len(msg.PowerOfAttorney) > 0 {
		token = msg.PowerOfAttorney
	}
	if publicKey, err := x509.ParsePKCS1PublicKey(token); err != nil {
		return nil
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.Signature) != nil {
			return nil
		}
	}
	// check wallet signature
	msg.WalletSignature, position = ParseByteArray(data, position)
	if position != length {
		return nil // this must be the last byte of the array
	}
	if publicKey, err := x509.ParsePKCS1PublicKey(msg.FeeWallet); err != nil {
		return nil
	} else {
		if rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], msg.WalletSignature) != nil {
			return nil
		}
	}
	return &msg
}
