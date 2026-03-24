package massalib

import (
	"crypto/sha256"
	"testing"
)

func TestMultiHash(t *testing.T) {
	data := []byte("hello world")

	// single hash should equal sha256
	single := multiHash(data, sha256.New)
	h := sha256.Sum256(data)
	if string(single) != string(h[:]) {
		t.Error("single hash mismatch")
	}

	// double hash (sha256(sha256(data)))
	double := multiHash(data, sha256.New, sha256.New)
	h2 := sha256.Sum256(h[:])
	if string(double) != string(h2[:]) {
		t.Error("double hash mismatch")
	}

	// empty alg list returns input unchanged
	result := multiHash(data)
	if string(result) != string(data) {
		t.Error("no-op hash should return input")
	}
}
