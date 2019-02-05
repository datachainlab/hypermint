package handler

import (
	"reflect"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/account"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/bluele/hypermint/pkg/validator"
	"github.com/ethereum/go-ethereum/common"
)

func NewHandler(am account.AccountMapper, cm *contract.ContractManager, envm *contract.EnvManager, valm validator.ValidatorMapper) types.Handler {
	return func(ctx types.Context, tx types.Tx) types.Result {
		switch tx := tx.(type) {
		case *transaction.TransferTx:
			return handleTransferTx(ctx, am, tx)
		case *transaction.ContractDeployTx:
			return handleContractDeployTx(ctx, cm, envm, tx)
		case *transaction.ContractCallTx:
			return handleContractCallTx(ctx, cm, envm, tx)
		case *transaction.ValidatorAddTx:
			return handleValidatorAddTx(ctx, tx, valm)
		case *transaction.ValidatorRemoveTx:
			return handleValidatorRemoveTx(ctx, tx, valm)
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

// TODO subtract bonded amount from the balance
func handleValidatorAddTx(ctx types.Context, tx *transaction.ValidatorAddTx, valm validator.ValidatorMapper) types.Result {
	val := validator.MakeValidatorFromTx(tx)
	if err := valm.Set(ctx, val); err != nil {
		return transaction.ErrInvalidValidatorAdd(transaction.DefaultCodespace, err.Error()).Result()
	}
	return types.Result{}
}

// TODO return bonded amount to the balance
func handleValidatorRemoveTx(ctx types.Context, tx *transaction.ValidatorRemoveTx, valm validator.ValidatorMapper) types.Result {
	valm.Remove(ctx, common.BytesToAddress(tx.PubKey.Bytes()))
	return types.Result{}
}
