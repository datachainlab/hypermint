package proof

import (
	"bytes"
	"fmt"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/db"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/types"
)

func (p *KVProof) AppendHeaderProofOp(h *types.Header) error {
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
		return fmt.Errorf("invalid block hash")
	}
	p.Proof.Ops = append(p.Proof.Ops, NewHeaderFieldOp([]byte(HeaderOp), proofs[12]).ProofOp())
	return nil
}

func (p KVProof) VerifyWithHeader(h *types.Header) error {
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
