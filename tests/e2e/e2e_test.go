package e2e

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/bluele/hypermint/tests/e2e/helper"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"
)

const (
	testContractPath = "../build/contract_test.wasm"
)

type E2ETestSuite struct {
	helper.NodeTestSuite
}

func (ts *E2ETestSuite) SetupTest() {
	// TODO these values should be configurable?
	ts.Setup("../../build", filepath.Join("../../.hm", cmn.RandStr(8)))
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

func (ts *E2ETestSuite) TestDeploy() {
	ctx := context.Background()
	ts.NoError(ts.DeployContract(ctx, ts.Account(1), testContractPath))
}

func (ts *E2ETestSuite) GetBalance(ctx context.Context, addr common.Address) (int, error) {
	cmd := fmt.Sprintf("balance --address=%v", addr.Hex())
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return 0, xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	} else {
		return strconv.Atoi(strings.TrimRight(string(out), "\n"))
	}
}

func (ts *E2ETestSuite) Transfer(ctx context.Context, from, to common.Address, amount int) error {
	cmd := fmt.Sprintf("transfer --address=%v --amount=10 --to=%v --gas=1 --password=password", from.Hex(), to.Hex())
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	}
	time.Sleep(2 * ts.TimeoutCommit)
	return nil
}

func (ts *E2ETestSuite) DeployContract(ctx context.Context, from common.Address, path string) error {
	cmd := fmt.Sprintf("contract deploy --address=%v --path=%v --gas=1 --password=password", from.Hex(), path)
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return xerrors.Errorf("%v:%v:%v", string(out), string(e), err)
	}
	time.Sleep(2 * ts.TimeoutCommit)
	return nil
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
