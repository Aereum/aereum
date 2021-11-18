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
package instructions

import (
	"errors"

	"github.com/Aereum/aereum/core/crypto"
)


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
