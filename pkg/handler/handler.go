package handler

import (
	"reflect"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/account"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/transaction"
)

func NewHandler(am account.AccountMapper, cm *contract.ContractManager, envm *contract.EnvManager) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result {
		switch tx := tx.(type) {
		case *transaction.TransferTx:
			return handleTransferTx(ctx, am, tx)
		case *transaction.ContractDeployTx:
			return handleContractDeployTx(ctx, cm, envm, tx)
		case *transaction.ContractCallTx:
			return handleContractCallTx(ctx, cm, envm, tx)
		default:
			errMsg := "Unrecognized Tx type: " + reflect.TypeOf(tx).Name()
			return types.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleTransferTx(ctx types.Context, am account.AccountMapper, tx *transaction.TransferTx) types.Result {
	if err := am.Transfer(ctx, tx.From, tx.Amount, tx.To); err != nil {
		return transaction.ErrFailTransfer(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}

func handleContractDeployTx(ctx types.Context, cm *contract.ContractManager, envm *contract.EnvManager, tx *transaction.ContractDeployTx) types.Result {
	addr, err := cm.DeployContract(ctx, tx)
	if err != nil {
		return transaction.ErrInvalidDeploy(transaction.DefaultCodespace, err.Error()).Result()
	}
	return handleContractCallTx(ctx, cm, envm, &transaction.ContractCallTx{
		Address:  addr,
		Func:     transaction.ContractInitFunc,
		CommonTx: tx.CommonTx,
	})
}

func handleContractCallTx(ctx types.Context, cm *contract.ContractManager, envm *contract.EnvManager, tx *transaction.ContractCallTx) types.Result {
	env, err := envm.Get(ctx, tx.Address, tx.Args)
	if err != nil {
		return transaction.ErrInvalidCall(transaction.DefaultCodespace, err.Error()).Result()
	}
	res, err := env.Exec(ctx, tx.Func)
	if err != nil {
		return transaction.ErrInvalidCall(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{
		Data: res,
	}
}
