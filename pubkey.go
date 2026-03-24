package massalib

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"slices"

	"github.com/KarpelesLab/base58"
	"lukechampine.com/blake3"
)

// PublicKey represents a Massa Ed25519 public key with a version prefix.
type PublicKey struct {
	Version uint64
	PubKey  ed25519.PublicKey
}

// AsAddress derives the user account Address from the public key by hashing
// the serialized key bytes with blake3.
func (pk *PublicKey) AsAddress() *Address {
	h := blake3.Sum256(pk.Bytes())

	res := &Address{
		Category: 0,
		Version:  0,
		Hash:     h[:],
	}

	return res
}

// Bytes returns the binary encoding of the public key (version varint followed by
// the raw Ed25519 key bytes).
func (pk *PublicKey) Bytes() []byte {
	return slices.Concat(EncodeProtobufVarint(pk.Version), pk.PubKey)
}

// MarshalBinary encodes the public key in massa byte format (for compatibility).
func (pk *PublicKey) MarshalBinary() ([]byte, error) {
	return pk.Bytes(), nil
}

// String returns the Massa text representation of the public key (P prefix
// followed by base58check-encoded version and key bytes).
func (pk *PublicKey) String() string {
	buf := pk.Bytes()
	cksum := multiHash(buf, sha256.New, sha256.New)
	return "P" + base58.Bitcoin.Encode(slices.Concat(buf, cksum[:4]))
}

// MarshalText implements encoding.TextMarshaler.
func (pk *PublicKey) MarshalText() (text []byte, err error) {
	return []byte(pk.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler by parsing a Massa public key
// string (e.g. "P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd").
func (pk *PublicKey) UnmarshalText(text []byte) error {
	if len(text) < 18 {
		return errors.New("invalid public key text: length must be higher than 18")
	}
	if text[0] != 'P' {
		return errors.New("invalid public key text: must start with P")
	}

	text = text[1:]
	buf, err := base58.Bitcoin.Decode(string(text))
	if err != nil {
		return fmt.Errorf("invalid public key text: failed to decode: %w", err)
	}

	// check checksum
	cksum := buf[len(buf)-4:]
	buf = buf[:len(buf)-4]
	h := multiHash(buf, sha256.New, sha256.New)
	if subtle.ConstantTimeCompare(cksum, h[:4]) == 0 {
		return errors.New("invalid public key text: bad checksum")
	}

	vers, l, err := DecodeProtobufVarint(buf)
	if err != nil {
		return fmt.Errorf("invalid public key text: failed to read version: %w", err)
	}
	buf = buf[l:]

	pk.Version = vers
	pk.PubKey = ed25519.PublicKey(buf)
	return nil
}
