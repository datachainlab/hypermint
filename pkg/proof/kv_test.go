package proof

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

func TestKVProofOp(t *testing.T) {
	assert := assert.New(t)
	f := func(h Header) bool {
		header := h.Value()
		op, err := MakeKVProofOp(header)
		assert.NoError(err)
		po, err := prt.Decode(op)
		assert.NoError(err)
		root, err := po.Run([][]byte{header.AppHash})
		assert.NoError(err)
		if assert.Equal(1, len(root)) {
			assert.EqualValues(header.Hash(), root[0])
		}
		return true
	}
	c := &quick.Config{
		MaxCountScale: 1000,
	}
	if err := quick.Check(f, c); err != nil {
		t.Error(err)
	}
}

type (
	Hash        [32]byte
	PartsHeader struct {
		Total int
		Hash  Hash
	}
	BlockID struct {
		Hash        Hash
		PartsHeader PartsHeader
	}
	WrappedTime int64
)

func randInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func (t WrappedTime) Generate(rand *rand.Rand, size int) reflect.Value {
	const secRange = 60 * 60 * 24 * 365 * 10
	v := WrappedTime(time.Now().Unix() + (randInt(-1*secRange, secRange)))
	return reflect.ValueOf(v)
}

// a copy from github.com/tendermint/tendermint/types/block.gotypes.Header
// type of each field is different from the original(for testing)
type Header struct {
	// basic block info
	Version  version.Consensus `json:"version"`
	ChainID  string            `json:"chain_id"`
	Height   int64             `json:"height"`
	Time     WrappedTime       `json:"time"`
	NumTxs   int64             `json:"num_txs"`
	TotalTxs int64             `json:"total_txs"`

	// prev block info
	LastBlockID BlockID `json:"last_block_id"`

	// hashes of block data
	LastCommitHash Hash `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       Hash `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     Hash `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash Hash `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      Hash `json:"consensus_hash"`       // consensus params for current block
	AppHash            Hash `json:"app_hash"`             // state after txs from the previous block
	LastResultsHash    Hash `json:"last_results_hash"`    // root hash of all results from the txs from the previous block

	// consensus info
	EvidenceHash    Hash `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress Hash `json:"proposer_address"` // original proposer of the block
}

func (h Header) Value() *types.Header {
	return &types.Header{
		Version:  h.Version,
		ChainID:  h.ChainID,
		Height:   h.Height,
		Time:     time.Unix(int64(h.Time), 0),
		NumTxs:   h.NumTxs,
		TotalTxs: h.TotalTxs,
		LastBlockID: types.BlockID{
			Hash: h.LastBlockID.Hash[:],
			PartsHeader: types.PartSetHeader{
				Total: h.LastBlockID.PartsHeader.Total,
				Hash:  h.LastBlockID.PartsHeader.Hash[:],
			},
		},
		LastCommitHash:     h.LastCommitHash[:],
		DataHash:           h.DataHash[:],
		ValidatorsHash:     h.ValidatorsHash[:],
		NextValidatorsHash: h.NextValidatorsHash[:],
		ConsensusHash:      h.ConsensusHash[:],
		AppHash:            h.AppHash[:],
		LastResultsHash:    h.LastCommitHash[:],
		EvidenceHash:       h.EvidenceHash[:],
		ProposerAddress:    h.ProposerAddress[:],
	}
}
