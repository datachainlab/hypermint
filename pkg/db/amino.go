package db

import (
	"github.com/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterAmino(cdc)
}

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(RWSets{},
		"hypermint/RWSets", nil)
	cdc.RegisterConcrete(RWSet{},
		"hypermint/RWSet", nil)
}
