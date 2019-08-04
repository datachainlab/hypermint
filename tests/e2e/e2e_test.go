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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cmn "github.com/tendermint/tendermint/libs/common"
	"golang.org/x/xerrors"
)

const (
	testContractPath = "../build/contract_test.wasm"
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

	ts.T().Run("check if update state successfully", func(t *testing.T) {
		_, err := ts.CallContract(ctx, ts.Account(1), contract, "test_write_state", []string{key, value}, false)
		assert.NoError(t, err)

		out, err := ts.CallContract(ctx, ts.Account(1), contract, "test_read_state", []string{key}, true)
		assert.NoError(t, err)
		assert.Equal(t, value, string(out))

		t.Run("ensure that expected event is happened", func(t *testing.T) {
			_, err := ts.CallContract(ctx, ts.Account(1), contract, "test_emit_event", []string{"first", "second"}, false)
			assert.NoError(t, err)
			count, err := ts.SearchEvent(ctx, contract, "test-event-name-0")
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
			count, err = ts.SearchEvent(ctx, contract, "test-event-name-1")
			assert.NoError(t, err)
			assert.Equal(t, 1, count)
		})
	})

	ts.T().Run("get a proof of updated state, and check if its proof is valid", func(t *testing.T) {
		cli := ts.RPCClient()
		kvp, err := helper.GetKVProofInfo(cli, contract, 0, []byte(key), []byte(value))
		if assert.NoError(t, err) {
			_, err := kvp.Marshal()
			assert.NoError(t, err)
			c, err := cli.Commit(&kvp.Height)
			assert.NoError(t, err)
			err = kvp.VerifyWithHeader(c.SignedHeader.Header)
			assert.NoError(t, err)
		}
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
