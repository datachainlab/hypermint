package app

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/account"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// State to Unmarshal
type GenesisState struct {
	Accounts []account.Account `json:"accounts"`
}

func GetInitChainer(am account.AccountMapper) func(types.Context, abci.RequestInitChain) abci.ResponseInitChain {
	return func(ctx types.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		stateJSON := req.AppStateBytes
		// TODO is this now the whole genesis file?

		var genesisState GenesisState
		err := json.Unmarshal(stateJSON, &genesisState)
		if err != nil {
			panic(err)
			// return sdk.ErrGenesisParse("").TraceCause(err, "")
		}

		for _, acc := range genesisState.Accounts {
			if _, err := am.AddBalance(ctx, acc.Address, acc.Amount); err != nil {
				panic(err)
			}
			fmt.Printf("addr=%v amount=%v\n",
				acc.Address.Hex(),
				acc.Amount,
			)

		}

		// load the initial stake information
		return abci.ResponseInitChain{}
	}
}

// Create the core parameters for genesis initialization
// note that the pubkey input is this machines pubkey
func AppGenState(cdc *amino.Codec, appGenTxs []json.RawMessage) (genesisState GenesisState, err error) {
	if len(appGenTxs) == 0 {
		err = errors.New("must provide at least genesis transaction")
		return
	}

	// get genesis flag account information
	accounts := make([]account.Account, 0, len(appGenTxs))
	for _, appGenTx := range appGenTxs {

		var genTx AppGenTx
		err = cdc.UnmarshalJSON(appGenTx, &genTx)
		if err != nil {
			return
		}

		if genTx.Address != "" {
			accounts = append(accounts, account.Account{
				Address: common.HexToAddress(genTx.Address),
				Amount:  100,
			})
		}
	}

	// create the final app state
	genesisState = GenesisState{
		Accounts: accounts,
	}
	return
}

// AppGenState but with JSON
func AppGenStateJSON(cdc *amino.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := AppGenState(cdc, appGenTxs)
	if err != nil {
		return nil, err
	}
	appState, err = json.Marshal(genesisState)
	return
}
