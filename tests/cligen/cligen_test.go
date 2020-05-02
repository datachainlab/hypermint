package cligen

import (
	"bytes"
	"context"
	"fmt"
	ecommon "github.com/bluele/hypermint/tests/e2e/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	cmn "github.com/tendermint/tendermint/libs/common"
)

const (
	testContractPath         = "../build/contract_test.wasm"
)

func TestCligenTestSuite(t *testing.T) {
	suite.Run(t, new(CligenTestSuite))
}

type CligenTestSuite struct {
	ecommon.NodeTestSuite
}

func (ts *CligenTestSuite) SetupTest() {
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

func (ts *CligenTestSuite) TearDownTest() {
	ts.TearDown()
}

func (ts *CligenTestSuite) TestCLI() {
	ctx := context.Background()
	c, err := ts.DeployContract(ctx, ts.Account(1), testContractPath)
	if !ts.NoError(err) {
		return
	}

	caller := ts.Account(1)

	for _, tc := range []struct{
		suc bool
		args []string
		outregex string
	}{
		{ true, []string{"help"}, `.*$`},
		{ true, []string{"test-get-sender"}, caller.Hex()+`\n$`},
		{ true, []string{"test-get-contract-address"}, c.Hex()+`\n$`},
		{ true, []string{"test-write-state", "--key", "a2V5", "--value", "dmFsdWU="}, `$`},
		{ true, []string{"test-read-state", "--key", "a2V5"}, `dmFsdWU=\n$`},
	} {
		var rootCmd = &cobra.Command{
			Use:   "test",
			Short: "Test CLI",
		}

		args := append([]string{"contract-test", "--passphrase", "password", "--ksdir", ts.GetCLIHomeDir()}, tc.args...)
		buf := new(bytes.Buffer)
		rootCmd.SetOut(buf)
		rootCmd.SetArgs(args)
		rootCmd.AddCommand(ContractTestCmd(c.Hex(), caller.Hex()))
		if err := rootCmd.Execute(); tc.suc != (err == nil) {
			if err != nil {
				ts.T().Error(err)
			} else {
				ts.T().Fail()
			}
		}
		output := buf.String()
		r := regexp.MustCompile(tc.outregex)
		if !r.MatchString(output) {
			ts.T().Error(output, "does not match with the regexp", "'"+tc.outregex+"'", args)
		}
	}
}

func (ts *CligenTestSuite) DeployContract(ctx context.Context, from common.Address, path string) (common.Address, error) {
	cmd := fmt.Sprintf("contract deploy --address=%v --path=%v --gas=1 --password=password", from.Hex(), path)
	address, err := ts.sendTxCMD(ctx, cmd)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(string(address)), nil
}

func (ts *CligenTestSuite) sendTxCMD(ctx context.Context, cmd string) ([]byte, error) {
	if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
		return nil, fmt.Errorf("%v:%v:%v", string(out), string(e), err)
	} else {
		time.Sleep(2 * ts.TimeoutCommit)
		return out, nil
	}
}
