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
// connection but also the connection public key on aereum network id.
//
// after establishing connection, the caller send the called the following
// message: its public key naked and a AES-256 symetric key encrypted using the
// caller aereum network public key.
// The message is signed with the caller network private key.
//
// the called then sends the caller the following message: the decrypted
// received AES-256 symetric key and another AES-256 symetric key, both
// encrypted by the caller ephemeral public key.
// The capacity to decrypt the original SHA-256 attests that the called is
// in possession of the secret key associated without the need of a signature.
//
// the symetric keys are used to encryption and authentication for sucessive
// messages using random nonce generated at the header of each message. The
// caller uses the key provided itself provided the called its own.

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

//func prependLength(msg []byte) []byte {
//	return append([]byte{byte(len(msg))}, msg...)
//}

/*func prependRead(msg []byte) ([]byte, []byte) {
	if len(msg) < 1 {
		return nil, nil
	}
	length := msg[0]
	if len(msg) < int(length)+1 {
		return nil, nil
	}
	return msg[1 : length+1], msg[length+1:]
}*/

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
