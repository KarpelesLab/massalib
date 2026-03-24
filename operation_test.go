package massalib

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

func TestOperationBytesRoundtrip(t *testing.T) {
	// Decode a known address for the destination
	addr, err := DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	op := &Operation{
		Fee:    1000,
		Expire: 100,
		Body: &BodyTransaction{
			Destination: addr,
			Amount:      5_000_000_000, // 5 MAS
		},
	}

	data := op.Bytes()
	if len(data) == 0 {
		t.Fatal("Bytes() returned empty")
	}

	// Roundtrip
	var op2 Operation
	err = op2.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	if op2.Fee != op.Fee {
		t.Errorf("Fee: got %d, want %d", op2.Fee, op.Fee)
	}
	if op2.Expire != op.Expire {
		t.Errorf("Expire: got %d, want %d", op2.Expire, op.Expire)
	}

	tx := op2.Body.(*BodyTransaction)
	if tx.Amount != 5_000_000_000 {
		t.Errorf("Amount: got %d, want 5000000000", tx.Amount)
	}
	if tx.Destination.String() != addr.String() {
		t.Errorf("Destination: got %s, want %s", tx.Destination.String(), addr.String())
	}
}

func TestOperationReadFromReader(t *testing.T) {
	addr, err := DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	op := &Operation{
		Fee:    0,
		Expire: 10,
		Body: &BodyTransaction{
			Destination: addr,
			Amount:      1,
		},
	}

	data := op.Bytes()

	var op2 Operation
	n, err := op2.ReadFrom(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ReadFrom: %v", err)
	}
	if n != int64(len(data)) {
		t.Errorf("ReadFrom read %d bytes, want %d", n, len(data))
	}
}

func TestOperationUnsupportedType(t *testing.T) {
	// Encode an operation with type 99 (unsupported)
	data := append(EncodeProtobufVarint(0), EncodeProtobufVarint(10)...)
	data = append(data, EncodeProtobufVarint(99)...)

	var op Operation
	err := op.UnmarshalBinary(data)
	if err == nil {
		t.Error("expected error for unsupported operation type")
	}
}

func TestOperationSignAndHash(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	addr, err := DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	op := &Operation{
		Fee:    1000,
		Expire: 100,
		Body: &BodyTransaction{
			Destination: addr,
			Amount:      1_000_000_000,
		},
	}

	// Hash should be deterministic
	h1 := op.Hash(MainNet, pub)
	h2 := op.Hash(MainNet, pub)
	if !bytes.Equal(h1, h2) {
		t.Error("Hash not deterministic")
	}
	if len(h1) != 32 {
		t.Errorf("Hash length: got %d, want 32", len(h1))
	}

	// Different chain IDs should produce different hashes
	h3 := op.Hash(BuildNet, pub)
	if bytes.Equal(h1, h3) {
		t.Error("different chain IDs should produce different hashes")
	}

	// Sign should produce valid output
	signed, err := op.Sign(MainNet, priv)
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if len(signed) == 0 {
		t.Error("Sign returned empty output")
	}
}

func TestBodyTransactionRoundtrip(t *testing.T) {
	addr, err := DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
	if err != nil {
		t.Fatalf("DecodeAddress: %v", err)
	}

	tx := &BodyTransaction{
		Destination: addr,
		Amount:      42,
	}

	data := tx.Bytes()

	var tx2 BodyTransaction
	err = tx2.UnmarshalBinary(data)
	if err != nil {
		t.Fatalf("UnmarshalBinary: %v", err)
	}

	if tx2.Amount != tx.Amount {
		t.Errorf("Amount: got %d, want %d", tx2.Amount, tx.Amount)
	}
	if tx2.Type() != OpTransaction {
		t.Errorf("Type: got %d, want %d", tx2.Type(), OpTransaction)
	}
}
