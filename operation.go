package massalib

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"slices"

	"lukechampine.com/blake3"
)

// OperationType identifies the kind of operation in the Massa network.
type OperationType uint32

// Massa operation types.
const (
	OpTransaction OperationType = iota // Transfer MAS coins
	OpRollBuy                          // Buy rolls for staking
	OpRollSell                         // Sell staking rolls
	OpExecuteSC                        // Execute smart contract bytecode
	OpCallSC                           // Call a smart contract function
)

var _ OperationBody = (*BodyTransaction)(nil)

// OperationBody is the interface implemented by all operation body types.
type OperationBody interface {
	Bytes() []byte
	UnmarshalBinary(data []byte) error
	Type() OperationType
	ReadFrom(r io.Reader) (int64, error)
}

// Operation represents a Massa operation consisting of a fee, expiration period, and a typed body.
type Operation struct {
	Fee    Amount
	Expire uint64 // expire_period, typically current period + 10
	Body   OperationBody
}

// Bytes returns the binary representation of the operation (fee, expire, type, body).
func (o *Operation) Bytes() []byte {
	// Fee, Expire, Type & body
	return slices.Concat(o.Fee.Bytes(), EncodeProtobufVarint(o.Expire), EncodeProtobufVarint(uint64(o.Body.Type())), o.Body.Bytes())
}

// Sign signs the given operation and returns a serialized signed operation suitable
// for submission to the network. The key must be an Ed25519 signer.
//
// See: https://docs.massa.net/docs/learn/operation-format-execution
func (o *Operation) Sign(chainId ChainId, key crypto.Signer) ([]byte, error) {
	pub, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("invalid key passed to Sign")
	}
	h := o.Hash(chainId, pub)
	sig, err := key.Sign(rand.Reader, h, crypto.Hash(0))
	if err != nil {
		return nil, err
	}
	return slices.Concat(EncodeProtobufVarint(0), sig, EncodeProtobufVarint(0), pub, o.Bytes()), nil
}

// Hash returns the blake3 hash of the operation contents prefixed with the chain ID
// and public key, as required by the Massa signing scheme.
func (o *Operation) Hash(chainId ChainId, pubKey ed25519.PublicKey) []byte {
	// chainid + pubkey vers[0] + pubkey + Bytes
	buf := slices.Concat(chainId.Bytes(), EncodeProtobufVarint(0), pubKey, o.Bytes())
	h := blake3.Sum256(buf)
	return h[:]
}

// UnmarshalBinary decodes an operation from its binary representation.
func (o *Operation) UnmarshalBinary(b []byte) error {
	_, err := o.ReadFrom(bytes.NewReader(b))
	return err
}

// ReadFrom reads an operation from the given reader.
func (o *Operation) ReadFrom(r io.Reader) (int64, error) {
	b := asBytereader(r)
	// read fee, expire, type
	fee, n1, err := ReadProtobufVarint(b)
	if err != nil {
		return n1, err
	}
	expire, n2, err := ReadProtobufVarint(b)
	if err != nil {
		return n1 + n2, err
	}
	typ, n3, err := ReadProtobufVarint(b)
	if err != nil {
		return n1 + n2 + n3, err
	}

	switch OperationType(typ) {
	case OpTransaction:
		o.Body = &BodyTransaction{}
	default:
		return n1 + n2 + n3, fmt.Errorf("unsupported tx type %d", typ)
	}

	n4, err := o.Body.ReadFrom(b)
	if err != nil {
		return n1 + n2 + n3 + n4, err
	}

	o.Fee = Amount(fee)
	o.Expire = expire
	return n1 + n2 + n3 + n4, nil
}

// BodyTransaction represents a MAS coin transfer operation.
type BodyTransaction struct {
	Destination *Address
	Amount      Amount
}

// Bytes returns the binary encoding of the transaction body.
func (tx *BodyTransaction) Bytes() []byte {
	return slices.Concat(tx.Destination.Bytes(), tx.Amount.Bytes())
}

// Type returns OpTransaction.
func (tx *BodyTransaction) Type() OperationType {
	return OpTransaction
}

// UnmarshalBinary decodes a transaction body from its binary representation.
func (tx *BodyTransaction) UnmarshalBinary(buf []byte) error {
	_, err := tx.ReadFrom(bytes.NewReader(buf))
	return err
}

// ReadFrom reads a transaction body from the given reader.
func (tx *BodyTransaction) ReadFrom(r io.Reader) (int64, error) {
	b := asBytereader(r)

	addr := &Address{}
	n1, err := addr.ReadFrom(b)
	if err != nil {
		return n1, err
	}
	var amt Amount
	n2, err := amt.ReadFrom(b)
	if err != nil {
		return n1 + n2, err
	}

	tx.Destination = addr
	tx.Amount = amt
	return n1 + n2, nil
}

// BodyRollBuy represents a roll purchase operation for staking.
type BodyRollBuy struct {
	Rolls uint64
}

// BodyRollSell represents a roll sale operation.
type BodyRollSell struct {
	Rolls uint64
}

// DataStoreItem represents a key-value pair in a smart contract datastore.
type DataStoreItem struct {
	Key   []byte
	Value []byte
}

// BodyExecuteSC represents a smart contract execution operation.
type BodyExecuteSC struct {
	MaxGas    uint64
	MaxCoins  Amount
	Bytecode  []byte           // Raw bytes of bytecode to execute (up to 10MB)
	Datastore []*DataStoreItem // Concatenated datastore items
}

// BodyCallSC represents a smart contract function call operation.
type BodyCallSC struct {
	MaxGas   uint64
	MaxCoins Amount
	Target   *Address
	Function string // Name of the function to call encoded as UTF-8 string
	Param    []byte
}
