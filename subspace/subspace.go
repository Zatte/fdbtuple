/*
 * subspace.go
 *
 * This source file is part of the FoundationDB open source project
 *
 * Copyright 2013-2018 Apple Inc. and the FoundationDB project authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// FoundationDB Go Subspace Layer

// Package subspace provides a convenient way to use FoundationDB tuples to
// define namespaces for different categories of data. The namespace is
// specified by a prefix tuple which is prepended to all tuples packed by the
// subspace. When unpacking a key with the subspace, the prefix tuple will be
// removed from the result.
//
// As a best practice, API clients should use at least one subspace for
// application data. For general guidance on subspace usage, see the Subspaces
// section of the Developer Guide
// (https://apple.github.io/foundationdb/developer-guide.html#subspaces).
package subspace

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/zatte/fdbtuple"
)

// Subspace represents a well-defined region of keyspace in a FoundationDB
// database.
type Subspace interface {
	// Sub returns a new Subspace whose prefix extends this Subspace with the
	// encoding of the provided element(s). If any of the elements are not a
	// valid fdbtuple.TupleElement, Sub will panic.
	Sub(el ...fdbtuple.TupleElement) Subspace

	// Bytes returns the literal bytes of the prefix of this Subspace.
	Bytes() []byte

	// Pack returns the key encoding the specified Tuple with the prefix of this
	// Subspace prepended.
	Pack(t fdbtuple.Tuple) fdbtuple.Key

	// PackWithVersionstamp returns the key encoding the specified tuple in
	// the subspace so that it may be used as the key in fdbtuple.Transaction's
	// SetVersionstampedKey() method. The passed tuple must contain exactly
	// one incomplete fdbtuple.Versionstamp instance or the method will return
	// with an error. The behavior here is the same as if one used the
	// fdbtuple.PackWithVersionstamp() method to appropriately pack together this
	// subspace and the passed fdbtuple.
	PackWithVersionstamp(t fdbtuple.Tuple) (fdbtuple.Key, error)

	// Unpack returns the Tuple encoded by the given key with the prefix of this
	// Subspace removed. Unpack will return an error if the key is not in this
	// Subspace or does not encode a well-formed fdbtuple.
	Unpack(k fdbtuple.KeyConvertible) (fdbtuple.Tuple, error)

	// Contains returns true if the provided key starts with the prefix of this
	// Subspace, indicating that the Subspace logically contains the key.
	Contains(k fdbtuple.KeyConvertible) bool

	// All Subspaces implement fdbtuple.KeyConvertible and may be used as
	// FoundationDB keys (corresponding to the prefix of this Subspace).
	fdbtuple.KeyConvertible

	// All Subspaces implement fdbtuple.ExactRange and fdbtuple.Range, and describe all
	// keys logically in this Subspace.
	fdbtuple.ExactRange
}

type subspace struct {
	rawPrefix []byte
}

// AllKeys returns the Subspace corresponding to all keys in a FoundationDB
// database.
func AllKeys() Subspace {
	return subspace{}
}

// Sub returns a new Subspace whose prefix is the encoding of the provided
// element(s). If any of the elements are not a valid fdbtuple.TupleElement, a
// runtime panic will occur.
func Sub(el ...fdbtuple.TupleElement) Subspace {
	return subspace{fdbtuple.Tuple(el).Pack()}
}

// FromBytes returns a new Subspace from the provided bytes.
func FromBytes(b []byte) Subspace {
	s := make([]byte, len(b))
	copy(s, b)
	return subspace{s}
}

// String implements the fmt.Stringer interface and return the subspace
// as a human readable byte string provided by fdbtuple.Printable.
func (s subspace) String() string {
	return fmt.Sprintf("Subspace(rawPrefix=%s)", fdbtuple.Printable(s.rawPrefix))
}

func (s subspace) Sub(el ...fdbtuple.TupleElement) Subspace {
	return subspace{concat(s.Bytes(), fdbtuple.Tuple(el).Pack()...)}
}

func (s subspace) Bytes() []byte {
	return s.rawPrefix
}

func (s subspace) Pack(t fdbtuple.Tuple) fdbtuple.Key {
	return fdbtuple.Key(concat(s.rawPrefix, t.Pack()...))
}

func (s subspace) PackWithVersionstamp(t fdbtuple.Tuple) (fdbtuple.Key, error) {
	return t.PackWithVersionstamp(s.rawPrefix)
}

func (s subspace) Unpack(k fdbtuple.KeyConvertible) (fdbtuple.Tuple, error) {
	key := k.FDBKey()
	if !bytes.HasPrefix(key, s.rawPrefix) {
		return nil, errors.New("key is not in subspace")
	}
	return fdbtuple.Unpack(key[len(s.rawPrefix):])
}

func (s subspace) Contains(k fdbtuple.KeyConvertible) bool {
	return bytes.HasPrefix(k.FDBKey(), s.rawPrefix)
}

func (s subspace) FDBKey() fdbtuple.Key {
	return fdbtuple.Key(s.rawPrefix)
}

func (s subspace) FDBRangeKeys() (fdbtuple.KeyConvertible, fdbtuple.KeyConvertible) {
	return fdbtuple.Key(concat(s.rawPrefix, 0x00)), fdbtuple.Key(concat(s.rawPrefix, 0xFF))
}

func (s subspace) FDBRangeKeySelectors() (fdbtuple.Selectable, fdbtuple.Selectable) {
	begin, end := s.FDBRangeKeys()
	return fdbtuple.FirstGreaterOrEqual(begin), fdbtuple.FirstGreaterOrEqual(end)
}

func concat(a []byte, b ...byte) []byte {
	r := make([]byte, len(a)+len(b))
	copy(r, a)
	copy(r[len(a):], b)
	return r
}
