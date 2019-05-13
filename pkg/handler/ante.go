package handler

import (
	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/account"
)

func NewAnteHandler(am account.AccountMapper) types.AnteHandler {
	return func(
		ctx types.Context, tt types.Tx, simulate bool,
	) (_ types.Context, _ types.Result, abort bool) {
		return ctx, types.Result{}, false
	}
}
