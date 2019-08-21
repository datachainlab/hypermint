package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bluele/hypermint/pkg/client/helper"
	ecommon "github.com/bluele/hypermint/tests/e2e/common"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
	cmn "github.com/tendermint/tendermint/libs/common"
	"golang.org/x/xerrors"
)

const (
	testContractPath         = "../build/contract_test.wasm"
	testExternalContractPath = "../build/external_contract_test.wasm"
)

type E2ETestSuite struct {
	ecommon.NodeTestSuite
}

func (ts *E2ETestSuite) SetupTest() {
	pjRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		ts.FailNow("failed to call Abs()", err.Error())
	}
	// TODO these values should be configurable?
	ts.Setup(
		filepath.Join(pjRoot, "build"),
		filepath.Join(filepath.Join(pjRoot, ".hm"), cmn.RandStr(8)),
	)
}

func (ts *E2ETestSuite) TearDownTest() {
	ts.TearDown()
}

func (ts *E2ETestSuite) TestBalance() {
	ctx := context.Background()
	{
		balance, err := ts.GetBalance(ctx, ts.Account(1))
		ts.NoError(err)
		ts.Equal(100, balance)
	}
	{
		_, err := ts.GetBalance(ctx, ts.Account(2))
		ts.Error(err)
	}
}

func (ts *E2ETestSuite) TestTransfer() {
	ctx := context.Background()

	{
		balance, err := ts.GetBalance(ctx, ts.Account(1))
		ts.NoError(err)
		ts.Equal(100, balance)
	}

	ts.NoError(ts.Transfer(ctx, ts.Account(1), ts.Account(2), 10))

	{
		balance, err := ts.GetBalance(ctx, ts.Account(1))
		ts.NoError(err)
		ts.Equal(90, balance)
	}
	{
		balance, err := ts.GetBalance(ctx, ts.Account(2))
		ts.NoError(err)
		ts.Equal(10, balance)
	}

	ts.NoError(ts.Transfer(ctx, ts.Account(2), ts.Account(1), 10))
	{
		balance, err := ts.GetBalance(ctx, ts.Account(1))
		ts.NoError(err)
		ts.Equal(100, balance)
	}
	{
		balance, err := ts.GetBalance(ctx, ts.Account(2))
		ts.NoError(err)
		ts.Equal(0, balance)
	}
}

func (ts *E2ETestSuite) TestContract() {
	ctx := context.Background()
	contract, err := ts.DeployContract(ctx, ts.Account(1), testContractPath)
	if !ts.NoError(err) {
		return
	}
	ts.T().Logf("contract address is %v", contract.Hex())

	const key = "key"
	const value = "value"

	ts.Run("check if update state successfully", func() {
		_, err := ts.CallContract(ctx, ts.Account(1), contract, "test_write_state", []string{key, value}, false)
		ts.NoError(err)

		out, err := ts.CallContract(ctx, ts.Account(1), contract, "test_read_state", []string{key}, true)
		ts.NoError(err)
		ts.Equal(value, string(out))

		ts.Run("ensure that expected event is happened", func() {
			_, err := ts.CallContract(ctx, ts.Account(1), contract, "test_emit_event", []string{"first", "second"}, false)
			ts.NoError(err)
			count, err := ts.SearchEvent(ctx, contract, "test-event-name-0")
			ts.NoError(err)
			ts.Equal(1, count)
			count, err = ts.SearchEvent(ctx, contract, "test-event-name-1")
			ts.NoError(err)
			ts.Equal(1, count)
		})
	})

	ts.Run("get a proof of updated state, and check if its proof is valid", func() {
		cli := ts.RPCClient()
		kvp, err := helper.GetKVProofInfo(cli, contract, 0, []byte(key), []byte(value))
		if ts.NoError(err) {
			_, err := kvp.Marshal()
			ts.NoError(err)
			c, err := cli.Commit(&kvp.Height)
			ts.NoError(err)
			err = kvp.VerifyWithHeader(c.SignedHeader.Header)
			ts.NoError(err)
		}
	})
}

func (ts *E2ETestSuite) TestCallExternalContract() {
	ctx := context.Background()
	contractAddress, err := ts.DeployContract(ctx, ts.Account(1), testContractPath)
	if !ts.NoError(err) {
		return
	}
	ts.T().Logf("contract address is %v", contractAddress.Hex())

	exContractAddress, err := ts.DeployContract(ctx, ts.Account(1), testExternalContractPath)
	if !ts.NoError(err) {
		return
	}
	ts.T().Logf("external contract address is %v", exContractAddress.Hex())

	ts.Run("call contract simply", func() {
		out, err := ts.CallContract(ctx, ts.Account(1), exContractAddress, "test_plus", []string{"1", "2"}, true)
		ts.NoError(err)
		ts.Equal("3", string(out))
	})

	ts.Run("call contract via contract", func() {
		out, err := ts.CallContract(ctx, ts.Account(1), contractAddress, "test_call_external_contract", []string{exContractAddress.Hex(), "1", "2"}, true)
		ts.NoError(err)
		ts.Equal("3", string(out))
	})

	ts.Run("check if caller address of external contract is an address of original contract", func() {
		out, err := ts.CallContract(ctx, ts.Account(1), contractAddress, "test_call_who_am_i_on_external_contract", []string{exContractAddress.Hex()}, true)
		ts.NoError(err)
		ts.Equal(contractAddress.Bytes(), out)
	})
}

func (ts *E2ETestSuite) GetBalance(ctx context.Context, addr common.Address) (int, error) {
	cmd := fmt.Sprintf("balance --address=%v", addr.Hex())
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return 0, xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	} else {
		return strconv.Atoi(string(out))
	}
}

func (ts *E2ETestSuite) Transfer(ctx context.Context, from, to common.Address, amount int) error {
	cmd := fmt.Sprintf("transfer --address=%v --amount=10 --to=%v --gas=1 --password=password", from.Hex(), to.Hex())
	_, err := ts.sendTxCMD(ctx, cmd)
	return err
}

func (ts *E2ETestSuite) DeployContract(ctx context.Context, from common.Address, path string) (common.Address, error) {
	cmd := fmt.Sprintf("contract deploy --address=%v --path=%v --gas=1 --password=password", from.Hex(), path)
	address, err := ts.sendTxCMD(ctx, cmd)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(string(address)), nil
}

func (ts *E2ETestSuite) CallContract(ctx context.Context, from, contract common.Address, fn string, args []string, isSimulate bool) ([]byte, error) {
	cmd := fmt.Sprintf(
		`contract call --address=%v --contract=%v --func="%v" --args=%#v --password=password --gas=1`,
		from.Hex(),
		contract.Hex(),
		fn,
		strings.Join(args, ","),
	)
	if isSimulate {
		cmd += " --simulate --silent"
	}
	return ts.sendTxCMD(ctx, cmd)
}

func (ts *E2ETestSuite) SearchEvent(ctx context.Context, contract common.Address, event string) (int, error) {
	cmd := fmt.Sprintf(
		`contract event search --address=%v --event=%v --count`,
		contract.Hex(),
		event,
	)
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return 0, xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	} else {
		return strconv.Atoi(string(out))
	}
}

func (ts *E2ETestSuite) sendTxCMD(ctx context.Context, cmd string) ([]byte, error) {
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return nil, xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	} else {
		time.Sleep(2 * ts.TimeoutCommit)
		return out, nil
	}
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
