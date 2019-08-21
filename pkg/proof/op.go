package proof

import (
	bytes "bytes"
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/store"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

const (
	HeaderOp           = "header"
	ProofOpHeaderField = "header:f"
)

var prt = store.DefaultProofRuntime()

func init() {
	prt.RegisterOpDecoder(ProofOpHeaderField, HeaderFieldOpDecoder)
}

type HeaderFieldOp struct {
	// Encoded in ProofOp.Key.
	key []byte

	// To encode in ProofOp.Data
	Proof *merkle.SimpleProof `json:"simple_proof"`
}

var _ merkle.ProofOperator = HeaderFieldOp{}

func NewHeaderFieldOp(key []byte, proof *merkle.SimpleProof) HeaderFieldOp {
	return HeaderFieldOp{
		key:   key,
		Proof: proof,
	}
}

func HeaderFieldOpDecoder(pop merkle.ProofOp) (merkle.ProofOperator, error) {
	if pop.Type != ProofOpHeaderField {
		return nil, cmn.NewError("unexpected ProofOp.Type; got %v, want %v", pop.Type, ProofOpHeaderField)
	}
	var op HeaderFieldOp // a bit strange as we'll discard this, but it works.
	err := cdc.UnmarshalBinaryLengthPrefixed(pop.Data, &op)
	if err != nil {
		return nil, cmn.ErrorWrap(err, "decoding ProofOp.Data into HeaderFieldOp")
	}
	return NewHeaderFieldOp(pop.Key, op.Proof), nil
}

func (op HeaderFieldOp) ProofOp() merkle.ProofOp {
	bz := cdc.MustMarshalBinaryLengthPrefixed(op)
	return merkle.ProofOp{
		Type: ProofOpHeaderField,
		Key:  op.key,
		Data: bz,
	}
}

func (op HeaderFieldOp) String() string {
	return fmt.Sprintf("HeaderFieldOp{%v}", op.GetKey())
}

func (op HeaderFieldOp) Run(args [][]byte) ([][]byte, error) {
	if len(args) != 1 {
		return nil, cmn.NewError("expected 1 arg, got %v", len(args))
	}
	value := cdcEncode(args[0])
	vhash := leafHash(value)

	if !bytes.Equal(vhash, op.Proof.LeafHash) {
		return nil, cmn.NewError("leaf hash mismatch: want %X got %X", op.Proof.LeafHash, vhash)
	}

	return [][]byte{
		op.Proof.ComputeRootHash(),
	}, nil
}

func (op HeaderFieldOp) GetKey() []byte {
	return op.key
}

// hash functions on tendermint
// (copy from https://github.com/bluele/tendermint/blob/ec53ce359bb8f011e4dbb715da098bea08c32ded/crypto/merkle/hash.go)

// TODO: make these have a large predefined capacity
var (
	leafPrefix  = []byte{0}
	innerPrefix = []byte{1}
)

// returns tmhash(0x00 || leaf)
func leafHash(leaf []byte) []byte {
	return tmhash.Sum(append(leafPrefix, leaf...))
}

// returns tmhash(0x01 || left || right)
func innerHash(left []byte, right []byte) []byte {
	return tmhash.Sum(append(innerPrefix, append(left, right...)...))
}

// cdcEncode returns nil if the input is nil, otherwise returns
// cdc.MustMarshalBinaryBare(item)
func cdcEncode(item interface{}) []byte {
	if item != nil && !cmn.IsTypedNil(item) && !cmn.IsEmpty(item) {
		return types.GetCodec().MustMarshalBinaryBare(item)
	}
	return nil
}
