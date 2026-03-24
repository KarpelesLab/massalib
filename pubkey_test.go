package massalib_test

import (
	"testing"

	"github.com/KarpelesLab/massalib"
)

func TestPubkeyParse(t *testing.T) {
	// from https://docs.massa.net/docs/learn/operation-format-execution#public-key
	pk := &massalib.PublicKey{}
	err := pk.UnmarshalText([]byte("P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd"))
	if err != nil {
		t.Fatalf("failed to parse key: %s", err)
	}
}

func TestPubkeyRoundtrip(t *testing.T) {
	input := "P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd"
	pk := &massalib.PublicKey{}
	if err := pk.UnmarshalText([]byte(input)); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}

	// String roundtrip
	if pk.String() != input {
		t.Errorf("String: got %s, want %s", pk.String(), input)
	}

	// MarshalText roundtrip
	text, err := pk.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	if string(text) != input {
		t.Errorf("MarshalText: got %s, want %s", text, input)
	}

	// MarshalBinary / Bytes
	data, err := pk.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if len(data) == 0 {
		t.Error("MarshalBinary returned empty")
	}
	if len(pk.Bytes()) != len(data) {
		t.Error("Bytes() and MarshalBinary() length mismatch")
	}
}

func TestPubkeyAsAddress(t *testing.T) {
	pk := &massalib.PublicKey{}
	if err := pk.UnmarshalText([]byte("P12TZEhzNX7rNGZ271jqCLThfwFhvQQuDpUM6n8uuNiupqsiUaCs")); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}

	addr := pk.AsAddress()
	if addr.String() != "AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH" {
		t.Errorf("AsAddress: got %s, want AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH", addr.String())
	}
}

func TestPubkeyUnmarshalTextErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"too short", "P123"},
		{"wrong prefix", "X1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd"},
		{"bad base58", "P!!!invalidbase58!!!!!!!!!!!!!!!!!!!!!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pk := &massalib.PublicKey{}
			err := pk.UnmarshalText([]byte(tt.input))
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}
