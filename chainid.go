package massalib

import "encoding/binary"

// ChainId identifies a Massa blockchain network.
type ChainId uint64

// Well-known Massa chain IDs.
const (
	MainNet   ChainId = 77658377
	BuildNet  ChainId = 77658366
	SecureNet ChainId = 77658383
	LabNet    ChainId = 77658376
	Sandbox   ChainId = 77
)

// Bytes returns the big-endian 8-byte encoding of the chain ID.
func (c ChainId) Bytes() []byte {
	v, _ := c.AppendBinary(nil)
	return v
}

// AppendBinary appends the big-endian 8-byte encoding of the chain ID to b.
func (c ChainId) AppendBinary(b []byte) ([]byte, error) {
	return binary.BigEndian.AppendUint64(b, uint64(c)), nil
}
