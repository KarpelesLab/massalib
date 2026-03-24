package massalib

import (
	"encoding/binary"
	"testing"
)

func TestChainIdBytes(t *testing.T) {
	tests := []struct {
		name string
		id   ChainId
	}{
		{"MainNet", MainNet},
		{"BuildNet", BuildNet},
		{"SecureNet", SecureNet},
		{"LabNet", LabNet},
		{"Sandbox", Sandbox},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.id.Bytes()
			if len(b) != 8 {
				t.Fatalf("Bytes() returned %d bytes, want 8", len(b))
			}
			got := binary.BigEndian.Uint64(b)
			if got != uint64(tt.id) {
				t.Errorf("Bytes() roundtrip: got %d, want %d", got, uint64(tt.id))
			}
		})
	}
}

func TestChainIdAppendBinary(t *testing.T) {
	prefix := []byte{0xAA, 0xBB}
	result, err := MainNet.AppendBinary(prefix)
	if err != nil {
		t.Fatalf("AppendBinary error: %v", err)
	}
	if len(result) != 10 {
		t.Fatalf("AppendBinary returned %d bytes, want 10", len(result))
	}
	if result[0] != 0xAA || result[1] != 0xBB {
		t.Error("AppendBinary did not preserve prefix")
	}
	got := binary.BigEndian.Uint64(result[2:])
	if got != uint64(MainNet) {
		t.Errorf("AppendBinary value: got %d, want %d", got, uint64(MainNet))
	}
}
