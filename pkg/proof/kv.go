package proof

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/store"
	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/db"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/types"
)

var prt = store.DefaultProofRuntime()

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
	kp = kp.AppendKey([]byte(app.ContractStoreKey.Name()), merkle.KeyEncodingURL)
	kp = kp.AppendKey(key, merkle.KeyEncodingHex)

	if err := prt.VerifyValue(
		p.Proof,
		h.AppHash,
		kp.String(),
		db.ValueObject{Value: p.Value, Version: ver}.Marshal(),
	); err != nil {
		return fmt.Errorf("failed to verify: %v", err)
	}

	return nil
}
