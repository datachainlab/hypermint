package validator

import (
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
}

// RegisterAmino registers all go-crypto related types in the given (amino) codec.
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(Validator{},
		"hypermint/Validator", nil)
}
