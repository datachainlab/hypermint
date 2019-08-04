package proof

import (
	"bytes"
	"fmt"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

func (p KVProofInfo) VerifyWithHeader(h *types.Header) error {
	if p.Height != h.Height {
		return fmt.Errorf("height is mismatch: %v != %v", p.Height, h.Height)
	}
	ver, err := db.MakeVersion(p.Version)
	if err != nil {
		return err
	}

	key := append(p.Contract, p.Key...)
	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte(HeaderOp), merkle.KeyEncodingURL)
	kp = kp.AppendKey([]byte(app.ContractStoreKey.Name()), merkle.KeyEncodingURL)
	kp = kp.AppendKey(key, merkle.KeyEncodingHex)

	if err := prt.VerifyValue(
		p.Proof,
		h.Hash(),
		kp.String(),
		db.ValueObject{Value: p.Value, Version: ver}.Marshal(),
	); err != nil {
		return fmt.Errorf("failed to verify: %v", err)
	}

	return nil
}

func MakeKVProofOp(h *types.Header) (merkle.ProofOp, error) {
	root, proofs := merkle.SimpleProofsFromByteSlices([][]byte{
		cdcEncode(h.Version),
		cdcEncode(h.ChainID),
		cdcEncode(h.Height),
		cdcEncode(h.Time),
		cdcEncode(h.NumTxs),
		cdcEncode(h.TotalTxs),
		cdcEncode(h.LastBlockID),
		cdcEncode(h.LastCommitHash),
		cdcEncode(h.DataHash),
		cdcEncode(h.ValidatorsHash),
		cdcEncode(h.NextValidatorsHash),
		cdcEncode(h.ConsensusHash),
		cdcEncode(h.AppHash),
		cdcEncode(h.LastResultsHash),
		cdcEncode(h.EvidenceHash),
		cdcEncode(h.ProposerAddress),
	})
	if !bytes.Equal(h.Hash(), root) {
		return merkle.ProofOp{}, fmt.Errorf("invalid block hash")
	}
	return NewHeaderFieldOp([]byte(HeaderOp), proofs[12]).ProofOp(), nil
}

func MakeKVProofInfo(height int64, proof *merkle.Proof, contract common.Address, key cmn.HexBytes, value *db.ValueObject) *KVProofInfo {
	kvp := &KVProofInfo{
		Height:   height,
		Proof:    proof,
		Contract: contract.Bytes(),
		Key:      key.Bytes(),
		Value:    value.Value,
		Version:  value.Version.Bytes(),
	}
	return kvp
}
