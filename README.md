# massalib

[![Tests](https://github.com/KarpelesLab/massalib/actions/workflows/test.yml/badge.svg)](https://github.com/KarpelesLab/massalib/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/KarpelesLab/massalib/badge.svg?branch=master)](https://coveralls.io/github/KarpelesLab/massalib?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/KarpelesLab/massalib.svg)](https://pkg.go.dev/github.com/KarpelesLab/massalib)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Go library for interacting with the [Massa](https://massa.net/) blockchain. Provides address encoding/decoding, public key handling, operation construction and signing, and a gRPC client for communicating with Massa nodes.

## Install

```bash
go get github.com/KarpelesLab/massalib
```

## Features

- **Addresses** — Decode, encode, and derive Massa addresses (user `AU` and smart contract `AS` prefixes) with base58check validation
- **Public keys** — Parse, serialize, and convert Ed25519 public keys to addresses
- **Operations** — Build, sign, hash, and serialize Massa operations (transfers, roll buy/sell, smart contract execute/call)
- **gRPC client** — Connect to a Massa node and query status, stream slot transfers, and submit operations
- **Chain IDs** — Constants for MainNet, BuildNet, SecureNet, LabNet, and Sandbox

## Usage

### Decode an address

```go
addr, err := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")
if err != nil {
    log.Fatal(err)
}
fmt.Println(addr.Thread()) // blockclique thread (0–31)
```

### Parse a public key and derive its address

```go
pk := &massalib.PublicKey{}
err := pk.UnmarshalText([]byte("P1t4JZwHhWNLt4xYabCbukyVNxSbhYPdF6wCYuRmDuHD784juxd"))
if err != nil {
    log.Fatal(err)
}
addr := pk.AsAddress()
fmt.Println(addr) // AU1zyQ2XEA6ZCCtR3CDQVRC7Q1rBpNZDmQe1EN5FXrQCq9435ziH
```

### Build and sign a transaction

```go
dest, _ := massalib.DecodeAddress("AU128qq86hv2NzXqhowRzaeMruThQyxQQC3PgW3cgHg2ttgXMTa1A")

op := &massalib.Operation{
    Fee:    1000,
    Expire: currentPeriod + 10,
    Body: &massalib.BodyTransaction{
        Destination: dest,
        Amount:      1_000_000_000, // 1 MAS
    },
}

signed, err := op.Sign(massalib.MainNet, privateKey)
if err != nil {
    log.Fatal(err)
}
```

### Connect to a node via gRPC

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

rpc, err := massalib.New("localhost:33037",
    grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer rpc.Close()

status, err := rpc.GetStatus(context.Background())
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Node %s running %s\n", status.NodeId, status.Version)
```

### Submit operations

```go
ids, err := rpc.SendOperations(ctx, signed)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Operation IDs:", ids)
```

## Types

| Type | Description |
|------|-------------|
| `Address` | Massa address (user or smart contract) with base58check encoding |
| `PublicKey` | Ed25519 public key with version prefix |
| `Amount` | Coin amount as uint64 with 9 decimal places (max ~18.4B MAS) |
| `ChainId` | Network identifier (MainNet, BuildNet, etc.) |
| `Operation` | Fee + expiry + typed body |
| `BodyTransaction` | MAS transfer |
| `BodyRollBuy` / `BodyRollSell` | Staking roll operations |
| `BodyExecuteSC` | Smart contract bytecode execution |
| `BodyCallSC` | Smart contract function call |
| `RPC` | gRPC client wrapper |

## License

MIT
