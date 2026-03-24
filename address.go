package massalib

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/KarpelesLab/base58"
	"lukechampine.com/blake3"
)

// Address represents a Massa blockchain address, which can be either a user account (Category 0)
// or a smart contract (Category 1).
//
// See: https://docs.massa.net/docs/learn/operation-format-execution#address
type Address struct {
	Category uint64 // 0 for User, 1 for Smart contract
	Version  uint64
	Hash     []byte // Underlying blake3 hash
}

// Thread returns the blockclique sharding thread of the address, computed by taking
// the first 5 bits of the underlying hash.
func (a *Address) Thread() byte {
	return a.Hash[0] & 0x1f
}

// SetValue computes the hash of the provided buffer and sets it as address Hash.
func (a *Address) SetValue(in []byte) {
	h := blake3.Sum256(in)
	a.Hash = h[:]
}

// Bytes encodes the address in massa byte format.
func (a *Address) Bytes() []byte {
	return slices.Concat(EncodeProtobufVarint(a.Category), EncodeProtobufVarint(a.Version), a.Hash)
}

// MarshalBinary encodes the address in massa byte format (for compatibility).
func (a *Address) MarshalBinary() ([]byte, error) {
	return a.Bytes(), nil
}

// ReadFrom reads a binary address from a given source, which is made of two varint
// (category, version) followed by 32 bytes of hash.
func (a *Address) ReadFrom(r io.Reader) (int64, error) {
	b := asBytereader(r)

	v1, n1, err := ReadProtobufVarint(b)
	if err != nil {
		return n1, err
	}
	v2, n2, err := ReadProtobufVarint(b)
	if err != nil {
		return n1 + n2, err
	}
	buf := make([]byte, 32)
	n3, err := io.ReadFull(b, buf)
	if err != nil {
		return n1 + n2 + int64(n3), err
	}

	a.Category = v1
	a.Version = v2
	a.Hash = buf
	return n1 + n2 + int64(n3), nil
}

// String encodes the address as a string following Massa encoding conventions.
func (a *Address) String() string {
	addrBytes := slices.Concat(EncodeProtobufVarint(a.Version), a.Hash)
	// Note that the base58check encoding is the same one as the one used by bitcoin but without version number as it is already included in the underlying binary serialization of the operation ID.
	cksum := multiHash(addrBytes, sha256.New, sha256.New)
	suffix := base58.Bitcoin.Encode(slices.Concat(addrBytes, cksum[:4]))

	switch a.Category {
	case 0:
		return "AU" + suffix
	case 1:
		return "AS" + suffix
	default:
		return "A?" + suffix
	}
}

// DecodeAddress parses a Massa address string (AU for user, AS for smart contract)
// or a public key string (P prefix) into an Address. It validates the base58check
// checksum and returns an error if the address is malformed.
func DecodeAddress(addr string) (*Address, error) {
	if len(addr) < 2 {
		return nil, errors.New("invalid massa address")
	}

	if addr[0] == 'P' {
		// from a public key

		buf, err := base58.Bitcoin.Decode(addr[1:])
		if err != nil {
			return nil, err
		}
		if len(buf) < 4 {
			return nil, errors.New("invalid massa address")
		}

		// check checksum
		cksum := buf[len(buf)-4:]
		buf = buf[:len(buf)-4]
		computedCksum := multiHash(buf, sha256.New, sha256.New)
		if subtle.ConstantTimeCompare(cksum, computedCksum[:4]) != 1 {
			return nil, errors.New("invalid massa public key")
		}

		// hash pubkey
		h := blake3.Sum256(buf)

		res := &Address{
			Category: 0,
			Version:  0,
			Hash:     h[:],
		}

		return res, nil
	}

	var cat uint64
	switch addr[:2] {
	case "AU":
		cat = 0
	case "AS":
		cat = 1
	default:
		return nil, fmt.Errorf("invalid massa address prefix %s", addr[:2])
	}

	addr = addr[2:]
	buf, err := base58.Bitcoin.Decode(addr)
	if err != nil {
		return nil, err
	}
	if len(buf) < 4 {
		return nil, errors.New("invalid massa address")
	}

	// check checksum
	cksum := buf[len(buf)-4:]
	buf = buf[:len(buf)-4]
	computedCksum := multiHash(buf, sha256.New, sha256.New)
	if subtle.ConstantTimeCompare(cksum, computedCksum[:4]) != 1 {
		return nil, errors.New("invalid massa address")
	}

	// let's decode the version
	vers, buflen, err := DecodeProtobufVarint(buf)
	if err != nil {
		return nil, err
	}

	res := &Address{
		Category: cat,
		Version:  vers,
		Hash:     buf[buflen:],
	}

	return res, nil
}
