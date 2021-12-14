package network

import (
	"crypto/subtle"
	"errors"
	"net"

	"github.com/Aereum/aereum/core/consensus"
	"github.com/Aereum/aereum/core/crypto"
	"github.com/Aereum/aereum/core/crypto/dh"
)

var errCouldNotSecure = errors.New("could not secure communication")

// Simple implementation of hasdshake for secure communication between nodes.
// The secure channel should not be used to transmit confidential information.
//
// Both parties in the communication channel must be registered users within
// the aereum network consensus blockchain.
//
// The caller must know from the onset not only the IP and port number for the
// connection but also the connection token on aereum network to check the
// identity of the server.
//
// After establishing connection, the caller send the called the following
// message: its token naked and an X25519 ephemeral public key for the diffie
// hellman consensus secret.
//
// The called then sends the caller the following message: a copy of the
// ephemeral public key received and another epheral public key, and the
// signature of this message. It uses the two ephemeral token to derive the
// diffie hellman secret key.
//
// The caller checks the validity of the information sent by the called and
// derives the diffie hellman secret key. It finally sends the ephemeral public
// key sent by the called signed.

// The called confirms the information sent by the caller and if checks the
// handshake is terminated and the secure connection estabilished.
//
// If any information is not valid, the connection is promptly terminated by
// any party.

// read the first byte (n) and read subsequent n-bytes from connection
func readhs(conn net.Conn) ([]byte, error) {
	length := make([]byte, 1)
	if n, err := conn.Read(length); n != 1 {
		return nil, err
	}
	msg := make([]byte, length[0])
	if n, err := conn.Read(msg); n != int(length[0]) {
		return nil, err
	}
	return msg, nil
}

func writehs(conn net.Conn, msg []byte) error {
	if len(msg) > 256 {
		return errors.New("msg too large to send")
	}
	msgToSend := append([]byte{byte(len(msg))}, msg...)
	if n, err := conn.Write(msgToSend); n != len(msgToSend) {
		return err
	}
	return nil
}

func writehsSigned(conn net.Conn, msg []byte, prv crypto.PrivateKey) error {
	if len(msg) > 256 {
		return errors.New("msg too large to send")
	}
	signature := prv.Sign(msg)
	if err := writehs(conn, msg); err != nil {
		return err
	}
	if err := writehs(conn, signature[:]); err != nil {
		return err
	}
	return nil
}

func PerformClientHandShake(conn net.Conn, prvKey crypto.PrivateKey, remotePub crypto.Token) (*SecureConnection, error) {
	// send public key and ephemeral public key for diffie hellman
	pubKey := prvKey.PublicKey()
	ephPrv, ephPub := dh.NewEphemeralKey()
	msg := append(pubKey[:], ephPub[:]...)
	writehs(conn, msg)

	// receive from server copy of the sent public key and another pub ephemeral key
	// signed
	resp, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	// test if the copy matches with subtle
	if subtle.ConstantTimeCompare(resp[0:crypto.TokenSize], ephPub[:]) != 1 {
		return nil, errors.New("client: copy of ephemeral key does not match")
	}
	// read and check signature
	respSign, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	if len(respSign) != crypto.SignatureSize {
		return nil, errors.New("client: signature size does not match")
	}
	var sign crypto.Signature
	copy(sign[:], respSign)
	if !remotePub.Verify(resp, sign) {
		return nil, errors.New("client: signature does not match")
	}
	// calculate diffie hellman shared secret
	var remoteEphToken crypto.Token
	copy(remoteEphToken[:], resp[crypto.TokenSize:])
	cipherKey := dh.ConsensusKey(ephPrv, remoteEphToken)
	// sent received ephemeral key signed to prove identity
	writehsSigned(conn, remoteEphToken[:], prvKey)
	return &SecureConnection{
		hash:         crypto.HashToken(remotePub),
		conn:         conn,
		cipher:       crypto.CipherNonceFromKey(cipherKey),
		cipherRemote: crypto.CipherNonceFromKey(cipherKey),
	}, nil
}

func PerformServerHandShake(conn net.Conn, prvKey crypto.PrivateKey, validator chan consensus.ValidatedConnection) (*SecureConnection, error) {
	resp, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	if len(resp) != 2*crypto.TokenSize {
		return nil, errors.New("server: public key + ephemeral key of wrong size")
	}
	// check if public key is a member: TODO check if is a validator
	ok := make(chan bool)
	var remoteToken crypto.Token
	copy(remoteToken[:], resp[:crypto.TokenSize])
	validator <- consensus.ValidatedConnection{Token: crypto.HashToken(remoteToken), Ok: ok}
	if !<-ok {
		conn.Close()
		return nil, errors.New("server: not a valid public key in the network")
	}
	var remoteEphToken crypto.Token
	copy(remoteEphToken[:], resp[crypto.TokenSize:])
	ephPrv, ephPub := dh.NewEphemeralKey()
	cipherKey := dh.ConsensusKey(ephPrv, remoteEphToken)
	msg := append(remoteEphToken[:], ephPub[:]...)
	if err := writehsSigned(conn, msg, prvKey); err != nil {
		return nil, err
	}
	// receive a copy of the sent key signed
	resp, err = readhs(conn)
	if err != nil {
		return nil, err
	}
	// test if the copy matches with subtle
	if len(resp) != crypto.TokenSize {
		return nil, errors.New("server: invalid token size")
	}
	if subtle.ConstantTimeCompare(resp, ephPub[:]) != 1 {
		return nil, errors.New("server: copy of ephemeral key does not match")
	}
	// read and check signature
	respSign, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	if len(respSign) != crypto.SignatureSize {
		return nil, errors.New("server: wrong signature size")
	}
	var sign crypto.Signature
	copy(sign[:], respSign)
	if !remoteToken.Verify(resp, sign) {
		return nil, errors.New("server: signature does not match")
	}
	return &SecureConnection{
		hash:         crypto.HashToken(remoteToken),
		conn:         conn,
		cipher:       crypto.CipherNonceFromKey(cipherKey),
		cipherRemote: crypto.CipherNonceFromKey(cipherKey),
	}, nil
}
