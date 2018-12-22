package types

import (
	cmn "github.com/tendermint/tendermint/libs/common"
)

// Result is the union of ResponseDeliverTx and ResponseCheckTx.
type Result struct {

	// Code is the response code, is stored back on the chain.
	Code CodeType

	// Codespace is the string referring to the domain of an error
	Codespace CodespaceType

	// Data is any data returned from the app.
	Data []byte

	// Log is just debug information. NOTE: nondeterministic.
	Log string

	// GasWanted is the maximum units of work we allow this tx to perform.
	GasWanted uint64

	// GasUsed is the amount of gas actually consumed. NOTE: unimplemented
	GasUsed uint64

	// Tags are used for transaction indexing and pubsub.
	Tags cmn.KVPairs
}

// TODO: In the future, more codes may be OK.
func (res Result) IsOK() bool {
	return res.Code.IsOK()
}
