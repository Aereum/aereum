package network

import (
	"crypto/subtle"
	"errors"
	"net"

	"github.com/Aereum/aereum/core/crypto"
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
	signature, err := prv.Sign(msg)
	if err != nil {
		panic(err)
	}
	if err := writehs(conn, msg); err != nil {
		return err
	}
	if err := writehs(conn, signature); err != nil {
		return err
	}
	return nil
}

func prependLength(msg []byte) []byte {
	return append([]byte{byte(len(msg))}, msg...)
}

func prependRead(msg []byte) ([]byte, []byte) {
	if len(msg) < 1 {
		return nil, nil
	}
	length := msg[0]
	if len(msg) < int(length)+1 {
		return nil, nil
	}
	return msg[1 : length+1], msg[length+1:]
}

func PerformClientHandShake(conn net.Conn, prvKey crypto.PrivateKey, remotePub crypto.PublicKey) (*SecureConnection, error) {
	// send public key and ephemeral public key
	// 	subtle.ConstantTimeCompare()
	PubKeyBytes := prvKey.PublicKey().ToBytes()
	msg := prependLength(PubKeyBytes)
	key := crypto.NewCipherKey()
	keyEncrypted, err := remotePub.Encrypt(key)
	if err != nil {
		return nil, err
	}
	msg = append(msg, prependLength(keyEncrypted)...)
	writehsSigned(conn, msg, prvKey)

	// receive server public key
	resp, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	keyBackEncrypted, resp := prependRead(resp)
	keyBack, err := prvKey.Decrypt(keyBackEncrypted)
	if err != nil || subtle.ConstantTimeCompare(key, keyBack) != 1 {
		return nil, errCouldNotSecure
	}
	remoteKeyEncrypted, resp := prependRead(resp)
	if remoteKeyEncrypted == nil {
		return nil, errCouldNotSecure
	}
	remoteKey, err := prvKey.Decrypt(remoteKeyEncrypted)
	if err != nil {
		return nil, errCouldNotSecure
	}
	if len(resp) != 0 {
		return nil, errCouldNotSecure
	}
	return &SecureConnection{
		conn:         conn,
		cipher:       crypto.CipherFromKey(key),
		cipherRemote: crypto.CipherFromKey(remoteKey),
	}, nil
}

func PerformServerHandShake(conn net.Conn, prvKey crypto.PrivateKey) (*SecureConnection, error) {
	resp, err := readhs(conn)
	if err != nil {
		return nil, err
	}
	remoteKeyBytes, resp := prependRead(resp)
	if remoteKeyBytes == nil {
		return nil, errCouldNotSecure
	}
	remoteKey, err := crypto.PublicKeyFromBytes(remoteKeyBytes)
	if err != nil {
		return nil, err
	}
	remoteCipherKeyEncrypted, resp := prependRead(resp)
	if remoteCipherKeyEncrypted == nil {
		return nil, errCouldNotSecure
	}
	if len(resp) != 0 {
		return nil, errCouldNotSecure
	}
	remoteCipherKey, err := prvKey.Decrypt(remoteCipherKeyEncrypted)
	if err != nil {
		return nil, err
	}
	remoteCipherKeyBackEncrypted, err := remoteKey.Encrypt(remoteCipherKey)
	if err != nil {
		return nil, err
	}
	key := crypto.NewCipherKey()
	keyEncrypted, err := remoteKey.Encrypt(key)
	if err != nil {
		return nil, err
	}
	msgToSend := prependLength(remoteCipherKeyBackEncrypted)
	msgToSend = append(msgToSend, prependLength(keyEncrypted)...)
	writehs(conn, msgToSend)
	return &SecureConnection{
		conn:         conn,
		cipher:       crypto.CipherFromKey(key),
		cipherRemote: crypto.CipherFromKey(remoteCipherKey),
	}, nil
}
