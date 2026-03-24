package massalib_test

import (
	"testing"

	"github.com/KarpelesLab/massalib"
)

func TestAddress(t *testing.T) {
	a, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("failed to parse addr: %s", err)
	}

	if a.String() != "AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A" {
		t.Errorf("failed to reformat address: %s != AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A", a.String())
	}

	a, err = massalib.DecodeAddress("P12TZEhzNX7rNGZ271jqCLThfwFhvQQuDpUM6n8uuNiupqsiUaCs")
	if err != nil {
		t.Fatalf("failed to parse addr from pubkey: %s", err)
	}

	if a.String() != "AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH" {
		t.Errorf("failed to reformat address: %s != %s", a.String(), "AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH")
	}
}

func TestDecodeAddressErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"too short", "A"},
		{"bad prefix", "XX128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A"},
		{"bad base58", "AU!!!invalid!!!"},
		{"bad checksum", "AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := massalib.DecodeAddress(tt.input)
			if err == nil {
				t.Errorf("expected error for input %q", tt.input)
			}
		})
	}
}

func TestAddressThread(t *testing.T) {
	a, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	thread := a.Thread()
	if thread > 31 {
		t.Errorf("Thread() = %d, should be <= 31", thread)
	}
}

func TestAddressMarshalBinary(t *testing.T) {
	a, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	data, err := a.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary: %v", err)
	}
	if len(data) == 0 {
		t.Error("MarshalBinary returned empty")
	}
}

func TestAddressSmartContract(t *testing.T) {
	// Decode a user address, modify to smart contract, and check prefix
	a, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	a.Category = 1
	s := a.String()
	if s[:2] != "AS" {
		t.Errorf("smart contract address should start with AS, got %s", s[:2])
	}
}
