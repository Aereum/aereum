// Copyright 2021 The Aereum Authors
// This file is part of the aereum library.
//
// The aereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The aereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package message contains data types related to aereum network.
package instruction

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)

// Basic template used for all message types
type MessageTemplate struct {
	Version			byte
	Instruction		byte
	Epoch			uint64
	Author			[]byte	//public key token
	Message			[]byte
	Wallet			[]byte	//public key token
	Fee				uint64
	WalletSignature	[]byte
	Attorney		[]byte	//public key token
	Signature		[]byte 
}

func NewMessageTemplate(epoch uint64, author crypto.PrivateKey, message []byte, wallet crypto.PrivateKey,
	fee uint64, attorney crypto.PrivateKey) *MessageTemplate {
	m := &MessageTemplate {
		Version:		0,
		Instruction:	1,
		Epoch: 			epoch,
		Author:			author.PublicKey().ToBytes(), 
		Message:		message,
	}
	// If there's a wallet assigned, they are paying for instruction fee
	if wallet != nil{
		m.Wallet := wallet.PublicKey().ToBytes(), 
		m.Fee := fee,
		hashed := crypto.Hasher(t.serializeWithouSignature())
		var err error
		m.WalletSignature, err = wallet.Sign(hashed[:])
		if err != nil {
			return nil
		}
		return m
	// Else author is paying and there will be no wallet field
	} else {
		m.Fee := fee,
		hashed := crypto.Hasher(t.serializeWithouSignature())
		var err error
		m.WalletSignature, err = author.Sign(hashed[:])
		if err != nil {
			return nil
		}
		return m
	}
	// If there's an attorney assigned they must sign the instruction
	if attorney != nil {
		m.Attorney := attorney.PublicKey().ToBytes()
		hashed := crypto.Hasher(t.serializeWithouSignature())
		var err error
		m.Signature, err = attorney.Sign(hashed[:])
		if err != nil {
			return nil
		}
		return m
	// Else instruction is signed by the author and there will be no attorney field
	} else {
		hashed := crypto.Hasher(t.serializeWithouSignature())
		var err error
		m.Signature, err = author.Sign(hashed[:])
		if err != nil {
			return nil
		}
		return m
	}
}

// "Message" field in MessageTemplate struct can be of one of the following types
// JoinNetwork, UpdateInfo, CreateAudience..

type JoinNetwork struct {
	MessageType		byte
	Caption			string
	Details			map[string]
}

func NewJoinNetwork(caption string, details map[string]) *JoinNetwork {
	jn := &JoinNetwork {
		MessageType:	0,
		Caption:		string,
		Details: 		map[string], // nao sei exatamente como declarar que vai ser um json aqui
	}
	return jn
}

// // Vai precisar de uma funcao pra mandar pra byte array pra MessageTemplate receber do jeito certo
// func JoinNetworkToByteArray() *JoinNetwork {
// 	//???
// }

type UpdateInfo struct {
	MessageType	byte
	Details		map[string]
}

func NewUpdateInfo(details map[string]) *UpdateInfo {
	ui := &UpdateInfo {
		MessageType:	1,
		Details:		details,
	}
	return ui
}

// func UpdateInfoToByteArray()

type CreateAudience {
	MessageType byte
	Audience		[]byte
	Sumission		[]byte
	Moderation		[]byte
	AudienceKey		[]byte
	SumissionKey	[]byte
	ModerationKey	[]byte
	Flag			byte
	Description		string
}

func NewCreateAudience(audience crypto.PrivateKey, submission crypto.PrivateKey, 
	moderation crypto.PrivateKey, audienceKey 
)
