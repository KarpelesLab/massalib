package massalib

import (
	"bytes"
	"testing"
)

func TestAmountBytes(t *testing.T) {
	tests := []struct {
		name string
		amt  Amount
	}{
		{"zero", 0},
		{"one nano", 1},
		{"one MAS", 1_000_000_000},
		{"large amount", 18_446_744_073_709_551_615},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.amt.Bytes()
			if len(b) == 0 {
				t.Fatal("Bytes() returned empty slice")
			}

			// verify roundtrip via ReadFrom
			var got Amount
			_, err := got.ReadFrom(bytes.NewReader(b))
			if err != nil {
				t.Fatalf("ReadFrom error: %v", err)
			}
			if got != tt.amt {
				t.Errorf("roundtrip: got %d, want %d", got, tt.amt)
			}
		})
	}
}

func TestAmountReadFromError(t *testing.T) {
	var a Amount
	_, err := a.ReadFrom(bytes.NewReader(nil))
	if err == nil {
		t.Error("expected error reading from empty reader")
	}
}
