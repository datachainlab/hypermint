package contract

import (
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
)

type ContractManager struct {
	mapper ContractMapper
}

func NewContractManager(mp ContractMapper) *ContractManager {
	return &ContractManager{
		mapper: mp,
	}
}

func (cm *ContractManager) GetContract(ctx types.Context, addr common.Address) (*Contract, error) {
	return cm.mapper.Get(ctx, addr)
}

func (cm *ContractManager) SaveContract(ctx types.Context, addr common.Address, c *Contract) error {
	cm.mapper.Put(ctx, addr, c)
	return nil
}

func (cm *ContractManager) DeployContract(ctx types.Context, tx *transaction.ContractDeployTx) (common.Address, error) {
	c := TxToContract(tx)
	addr := c.Address()
	_, err := cm.GetContract(ctx, addr)
	if err == nil {
		return addr, fmt.Errorf("already exists: %v", addr.Hex())
	} else if err != ErrContractNotFound {
		return addr, err
	}
	if err := cm.SaveContract(ctx, addr, c); err != nil {
		return addr, err
	}
	return addr, nil
}
