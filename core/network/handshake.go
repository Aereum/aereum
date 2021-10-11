package network

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"net"
)

type HandshakeSecrets struct {
	prvKey          *rsa.PrivateKey
	cypherKey       [32]byte
	cypherRemoteKey [32]byte
	remotePubKey    *rsa.PublicKey
}

func send256Message(conn *net.TCPConn, msg []byte, prv *rsa.PrivateKey) error {
	if len(msg) > 256 {
		return errors.New("msg to large to send256Message")
	}
	msgToSend := append([]byte{byte(len(msg))}, msg...)
	nonce := make([]byte, 32)
	if n, err := rand.Read(nonce); n != 32 {
		return err
	}
	msgToSend = append(msgToSend, nonce...)
	hashed := sha256.Sum256(msgToSend)
	signature, err := rsa.SignPKCS1v15(nil, prv, crypto.SHA256, hashed[:])
	if err != nil {
		return err
	}
	signatureWithLength := append([]byte{byte(len(signature))}, signature...)
	msgToSend = append(msgToSend, signatureWithLength...)
	if n, err := conn.Write(msgToSend); n != len(msgToSend) {
		return err
	}
	return nil
}

func receive256Message(conn *net.TCPConn, pub *rsa.PublicKey) ([]byte, error) {
	length := make([]byte, 1)
	if n, err := conn.Read(length); n != 1 {
		return nil, err
	}
	msg := make([]byte, length[0])
	if n, err := conn.Read(msg); n != 1 {
		return nil, err
	}
	if pub == nil {
		publicKeyBytes := msg[0 : len(msg)-32]
		var err error
		pub, err = x509.ParsePKCS1PublicKey(publicKeyBytes)
		if err != nil {
			return nil, err
		}
	}
	signatureLength := make([]byte, 1)
	if n, err := conn.Read(signatureLength); n != 1 {
		return nil, err
	}
	signature := make([]byte, signatureLength[0])
	if n, err := conn.Read(signature); n != 1 {
		return nil, err
	}
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, msg, signature); err != nil {
		return nil, err
	}
	return msg[0 : len(msg)-32], nil
}

func PerformClientHandShake(conn *net.TCPConn, prvKey *rsa.PrivateKey) error {
	var hs HandshakeSecrets
	hs.prvKey = prvKey
	// receive server public key
	resp, err := receive256Message(conn, nil)
	if err != nil {
		return err
	}
	hs.remotePubKey, err = x509.ParsePKCS1PublicKey(resp)
	if err != nil {
		return err
	}

	// send public key
	pubKey := &prvKey.PublicKey
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	if err := send256Message(conn, pubKeyBytes, prvKey); err != nil {
		return err
	}

	// receice ephemeral key
	resp, err = receive256Message(conn, hs.remotePubKey)
	if err != nil {
		return err
	}
	if len(resp) != 32 {
		return errors.New("wrone cypher length")
	}
	for n := 0; n < 32; n++ {
		hs.cypherRemoteKey[n] = resp[n]
	}

}

func PerformServerHandShake(conn *net.TCPConn, prvKey *rsa.PrivateKey) error {
	var hs HandshakeSecrets
	hs.prvKey = prvKey
	// send public key
	pubKey := &prvKey.PublicKey
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	if err := send256Message(conn, pubKeyBytes, prvKey); err != nil {
		return err
	}
	// receive public key
	resp, err := receive256Message(conn, nil)
	if err != nil {
		return err
	}
	hs.remotePubKey, err = x509.ParsePKCS1PublicKey(resp)
	if err != nil {
		return err
	}
	// create cypher key
	cypherKey := make([]byte, 32)
	if n, err := rand.Read(cypherKey); n != 32 {
		return err
	}
	for n := 0; n < 32; n++ {
		hs.cypherKey[n] = cypherKey[n]
	}
	// encrypt it
	var cryptoCypher []byte
	cryptoCypher, err = rsa.EncryptPKCS1v15(nil, hs.remotePubKey, cypherKey)
	if err != nil {
		return err
	}
	send256Message(conn, cryptoCypher, prvKey)
}
