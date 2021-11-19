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
// along with the aereum library. If not, see <http://www.gnu.org/licenses/>.

/*
	Package instructions provides data structures for the defined instructions
	of the aereum communication layer protocol.

	They are:
		Transfer
		Deposit
		Withdraw
		JoinNetwork
		UpdateInfo
		CreateAudience
		JoinAudience
		AcceptJoinRequest
		Content
		UpdateAudience
		GrantPowerOfAttorney
		RevokePowerOfAttorney
		SponsorshipOffer
		SponsorshipAcceptance
		CreateEphemeral
		SecureChannel
		React

	Each instruction has its canonical binary encoding rules implemented.

	It provides also data structures for state and state mutations, together
	with the prescribed validation rules against a given state.

*/
package instructions
