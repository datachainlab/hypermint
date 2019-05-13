package types

import "github.com/bluele/hypermint/pkg/abci/codec"

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
}
