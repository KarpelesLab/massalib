package massalib

import "hash"

// multiHash performs sequential multi-level hashing: the output of each hash
// function is fed as input to the next. This is used for base58check checksums
// (double SHA-256).
func multiHash(b []byte, alg ...func() hash.Hash) []byte {
	var x []byte
	for _, a := range alg {
		h := a()
		h.Write(b)
		b = h.Sum(x)
		x = b[:0]
	}
	return b
}
