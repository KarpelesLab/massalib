package massalib

import (
	"bytes"
	"math"
	"testing"
)

func TestEncodeDecodeProtobufVarint(t *testing.T) {
	tests := []struct {
		name  string
		value uint64
	}{
		{"zero", 0},
		{"one", 1},
		{"single byte max", 127},
		{"two byte min", 128},
		{"300", 300},
		{"16384", 16384},
		{"max uint32", math.MaxUint32},
		{"max uint64", math.MaxUint64},
		{"large value", 1<<63 | 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := EncodeProtobufVarint(tt.value)
			decoded, n, err := DecodeProtobufVarint(encoded)
			if err != nil {
				t.Fatalf("DecodeProtobufVarint(%v) error: %v", encoded, err)
			}
			if n != len(encoded) {
				t.Errorf("DecodeProtobufVarint consumed %d bytes, want %d", n, len(encoded))
			}
			if decoded != tt.value {
				t.Errorf("roundtrip: got %d, want %d", decoded, tt.value)
			}
		})
	}
}

func TestDecodeProtobufVarintErrors(t *testing.T) {
	// empty buffer
	_, _, err := DecodeProtobufVarint(nil)
	if err == nil {
		t.Error("expected error for empty buffer")
	}

	// unterminated varint (all continuation bits set)
	_, _, err = DecodeProtobufVarint([]byte{0x80, 0x80})
	if err == nil {
		t.Error("expected error for unterminated varint")
	}

	// overflow: 10th byte > 1
	overflow := make([]byte, 10)
	for i := 0; i < 9; i++ {
		overflow[i] = 0x80
	}
	overflow[9] = 0x02
	_, _, err = DecodeProtobufVarint(overflow)
	if err == nil {
		t.Error("expected overflow error")
	}
}

func TestReadProtobufVarint(t *testing.T) {
	tests := []uint64{0, 1, 127, 128, 300, 16384, math.MaxUint32, math.MaxUint64}

	for _, val := range tests {
		encoded := EncodeProtobufVarint(val)
		r := bytes.NewReader(encoded)
		decoded, n, err := ReadProtobufVarint(r)
		if err != nil {
			t.Fatalf("ReadProtobufVarint(%d) error: %v", val, err)
		}
		if n != int64(len(encoded)) {
			t.Errorf("ReadProtobufVarint(%d) read %d bytes, want %d", val, n, len(encoded))
		}
		if decoded != val {
			t.Errorf("ReadProtobufVarint roundtrip: got %d, want %d", decoded, val)
		}
	}
}

func TestReadProtobufVarintEOF(t *testing.T) {
	r := bytes.NewReader(nil)
	_, _, err := ReadProtobufVarint(r)
	if err == nil {
		t.Error("expected error for empty reader")
	}
}

func TestEncodeProtobufVarintLength(t *testing.T) {
	// zero should be 1 byte
	if n := len(EncodeProtobufVarint(0)); n != 1 {
		t.Errorf("EncodeProtobufVarint(0) = %d bytes, want 1", n)
	}
	// max uint64 should be 10 bytes
	if n := len(EncodeProtobufVarint(math.MaxUint64)); n != 10 {
		t.Errorf("EncodeProtobufVarint(MaxUint64) = %d bytes, want 10", n)
	}
}
