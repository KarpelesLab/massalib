package massalib

import (
	"fmt"
	"io"
)

// EncodeProtobufVarint encodes the given integer using protobuf varint conventions.
func EncodeProtobufVarint(x uint64) []byte {
	var buf [10]byte // varint for 64-bit fits in at most 10 bytes
	n := 0
	for x >= 0x80 {
		buf[n] = byte(x&0x7F | 0x80)
		x >>= 7
		n++
	}
	buf[n] = byte(x)
	return buf[:n+1]
}

// DecodeProtobufVarint decodes a varint encoded following protobuf conventions
// from the given byte slice. It returns the decoded value, the number of bytes
// consumed, and any error encountered.
func DecodeProtobufVarint(buf []byte) (result uint64, buflen int, err error) {
	var shift uint
	for i, b := range buf {
		// If this is the 10th byte (i == 9), it can only hold a single bit (if any). We use b because this can't continue into next bit either.
		if i == 9 {
			if b > 1 {
				return 0, 0, fmt.Errorf("protobuf varint overflow: too many bits")
			}
			result |= uint64(b) << 63
			return result, i + 1, nil
		}
		// Extract the 7 bits of data
		val := uint64(b & 0x7F)

		result |= val << shift

		if (b & 0x80) == 0 {
			// Last byte
			return result, i + 1, nil
		}

		shift += 7
	}
	// If we get here, we ran out of buffer before finding a termination byte
	return 0, 0, fmt.Errorf("buffer ended prematurely while decoding varint")
}

// ReadProtobufVarint reads a protobuf varint from the given reader. It returns
// the decoded value, the number of bytes read, and any error encountered.
func ReadProtobufVarint(r bytereader) (uint64, int64, error) {
	var result uint64
	var shift uint
	var ln int64

	for i := 0; ; i++ { // Varint won't exceed 10 bytes
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return 0, ln, err
		}
		ln += 1
		// If this is the 10th byte read, check for overflow & return
		if i == 9 {
			if b > 1 {
				return 0, ln, fmt.Errorf("protobuf varint overflow: too many bits")
			}
			result |= uint64(b) << 63
			return result, ln, nil
		}
		val := uint64(b & 0x7F)

		result |= val << shift

		if (b & 0x80) == 0 {
			// Last byte read
			return result, ln, nil
		}

		shift += 7
	}
}
