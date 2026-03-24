package massalib

import (
	"io"
)

// Amount represents a Massa coin amount as a uint64 with 9 decimal places of precision.
// The maximum representable amount is 18,446,744,073.709551615 MAS.
//
// See: https://docs.massa.net/docs/learn/operation-format-execution#coin-amounts
type Amount uint64

// Bytes returns the protobuf varint encoding of the amount.
func (a Amount) Bytes() []byte {
	return EncodeProtobufVarint(uint64(a))
}

// ReadFrom reads a varint-encoded amount from the given reader.
func (a *Amount) ReadFrom(r io.Reader) (int64, error) {
	b := asBytereader(r)
	v, ln, err := ReadProtobufVarint(b)
	if err != nil {
		return ln, err
	}
	*a = Amount(v)
	return ln, nil
}
