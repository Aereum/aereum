// Copyright 2021 The aereum Authors
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
package message

import (
	"reflect"
	"testing"
)

func TestSubscribe(t *testing.T) {
	s := &Subscribe{
		Token:   []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Caption: "Caption",
		Details: "Details",
	}
	bytes := s.Serialize()
	r := ParseSubscribe(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for subscribe: %v, %v\n", *s, *r)
	}
}

func TestAbout(t *testing.T) {
	s := &About{
		Details: "Details",
	}
	bytes := s.Serialize()
	r := ParseAbout(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for About: %v, %v\n", *s, *r)
	}
}

func TestCreateAudience(t *testing.T) {
	s := &CreateAudience{
		Token:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Description: "Description",
	}
	bytes := s.Serialize()
	r := ParseCreateAudience(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for CreateAudience: %v, %v\n", *s, *r)
	}
}

func TestJoinAudience(t *testing.T) {
	s := &JoinAudience{
		Audience: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Expire:   15,
	}
	bytes := s.Serialize()
	r := ParseJoinAudience(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for JoinAudience: %v, %v\n", *s, *r)
	}
}

func TestChangeAudience(t *testing.T) {
	f1 := Follower{
		Token:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Secret: []byte{1, 2, 3, 4, 5, 6, 7, 8, 13},
	}
	f2 := Follower{
		Token:  []byte{1, 2, 3, 4, 5, 6, 7, 9, 10},
		Secret: []byte{1, 2, 3, 4, 5, 6, 7, 13},
	}

	s := &ChangeAudience{
		Audience: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Add:      []*Follower{&f1, &f2, &f1},
		Remove:   []*Follower{&f2, &f1},
		Details:  "Details",
	}
	bytes := s.Serialize()
	r := ParseChangeAudience(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for ChangeAudience: %v, %v\n", *s, *r)
	}
}

func TestAdvertisingOffer(t *testing.T) {
	s := &AdvertisingOffer{
		Token:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Audience:       []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		ContentType:    "string",
		ContentData:    []byte{13, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		AdvertisingFee: 82982989282,
		Expire:         15,
	}
	bytes := s.Serialize()
	r := ParseAdvertisingOffer(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for AdvertisingOffer: %v, %v\n", *s, *r)
	}
}

func TestContent(t *testing.T) {
	s := &Content{
		Audience:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		ContentSecret:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 1},
		ContentType:      "string",
		ContentData:      []byte{13, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		AdvertisingToken: []byte{1, 2, 5, 6, 7, 8},
	}
	bytes := s.Serialize()
	r := ParseContent(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for Content: %v, %v\n", *s, *r)
	}
}

func TestGrantPowerOfAttorney(t *testing.T) {
	s := &GrantPowerOfAttorney{
		Token:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		Expire: 167,
	}
	bytes := s.Serialize()
	r := ParseGrantPowerOfAttorney(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for GrantPowerOfAttorney: %v, %v\n", *s, *r)
	}
}

func TestRevokePowerOfAttorney(t *testing.T) {
	s := &RevokePowerOfAttorney{
		Token: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
	}
	bytes := s.Serialize()
	r := ParseRevokePowerOfAttorney(bytes)
	if !reflect.DeepEqual(*s, *r) {
		t.Errorf("Serialize and Parsing incompatible for RevokePowerOfAttorney: %v, %v\n", *s, *r)
	}
}
