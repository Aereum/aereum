// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ed25519 implements the Ed25519 signature algorithm. See
// https://ed25519.cr.yp.to/.
//
// These functions are also compatible with the “Ed25519” function defined in
// RFC 8032. However, unlike RFC 8032's formulation, this package's private key
// representation includes a public key suffix to make multiple signing
// operations with the same key more efficient. This package refers to the RFC
// 8032 private key as the “seed”.
//
// This is a interface adaptation of the original file.
package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"

	"github.com/Aereum/aereum/core/crypto/edwards25519"
)

type Signature [SignatureSize]byte

var ZeroToken Token
var ZeroPrivateKey PrivateKey

func RandomAsymetricKey() (Token, PrivateKey) {

	var public [32]byte
	var private PrivateKey

	seed := make([]byte, PublicKeySize)
	rand.Read(seed)
	digest := sha512.Sum512(seed)
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest[:])
	edwards25519.GeScalarMultBase(&A, &hBytes)
	A.ToBytes(&public)
	copy(private[0:32], seed)
	copy(private[32:], public[:])
	return Token(public), private
}

func PrivateKeyFromSeed(seed [32]byte) PrivateKey {
	var public [32]byte
	var private PrivateKey

	digest := sha512.Sum512(seed[:])
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest[:])
	edwards25519.GeScalarMultBase(&A, &hBytes)
	A.ToBytes(&public)
	copy(private[0:32], seed[:])
	copy(private[32:], public[:])
	return private
}

type PrivateKey [PrivateKeySize]byte

func (p PrivateKey) PublicKey() Token {
	var token Token
	copy(token[:], p[32:])
	return token
}

func (p PrivateKey) Sign(msg []byte) Signature {

	var signature Signature

	h := sha512.New()
	h.Write(p[:32])
	var digest1, messageDigest, hramDigest [64]byte
	var expandedSecretKey [32]byte
	h.Sum(digest1[:0])
	copy(expandedSecretKey[:], digest1[:])
	expandedSecretKey[0] &= 248
	expandedSecretKey[31] &= 63
	expandedSecretKey[31] |= 64

	h.Reset()
	h.Write(digest1[32:])
	h.Write(msg)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	edwards25519.ScReduce(&messageDigestReduced, &messageDigest)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	h.Reset()
	h.Write(encodedR[:])
	h.Write(p[32:])
	h.Write(msg)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, &expandedSecretKey, &messageDigestReduced)

	copy(signature[:32], encodedR[:])
	copy(signature[32:], s[:])
	return signature
}

type Token [TokenSize]byte

func (t Token) Equal(another Token) bool {
	return t == another
}

func (t Token) Verify(msg []byte, signature Signature) bool {
	if signature[63]&224 != 0 {
		return false
	}

	var A edwards25519.ExtendedGroupElement
	publicKey := [32]byte(t)
	if !A.FromBytes(&publicKey) {
		return false
	}
	edwards25519.FeNeg(&A.X, &A.X)
	edwards25519.FeNeg(&A.T, &A.T)

	h := sha512.New()
	h.Write(signature[:32])
	h.Write(publicKey[:])
	h.Write(msg)
	var digest [64]byte
	h.Sum(digest[:0])

	var hReduced [32]byte
	edwards25519.ScReduce(&hReduced, &digest)

	var R edwards25519.ProjectiveGroupElement
	var s [32]byte
	copy(s[:], signature[32:])

	// https://tools.ietf.org/html/rfc8032#section-5.1.7 requires that s be in
	// the range [0, order) in order to prevent signature malleability.
	if !edwards25519.ScMinimal(&s) {
		return false
	}

	edwards25519.GeDoubleScalarMultVartime(&R, &hReduced, &A, &s)

	var checkR [32]byte
	R.ToBytes(&checkR)
	return bytes.Equal(signature[:32], checkR[:])

}
