package common

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bluele/hypermint/pkg/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mattn/go-shellwords"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	cfg "github.com/tendermint/tendermint/config"
	rpclient "github.com/tendermint/tendermint/rpc/client"
)

const (
	testMnemonic = "math razor capable expose worth grape metal sunset metal sudden usage scheme"
	hdwPath      = "m/44'/60'/0'/0/"
	validatorIdx = 0
)

type NodeTestSuite struct {
	suite.Suite

	hmdProcess *os.Process
	accounts   []common.Address

	TimeoutCommit time.Duration

	hmdHomeDir   string
	hmcliHomeDir string
	hmdPath      string
	hmcliPath    string
}

func (ts *NodeTestSuite) Setup(binDir, homeDir string) {
	ts.TimeoutCommit = 100 * time.Millisecond

	ts.hmdHomeDir = filepath.Join(homeDir, "hmd")
	ts.hmcliHomeDir = filepath.Join(homeDir, "hmcli")
	ts.hmdPath = filepath.Join(binDir, "hmd")
	ts.hmcliPath = filepath.Join(binDir, "hmcli")

	ts.InitNode()
	ts.RunNode()
}

func (ts *NodeTestSuite) TearDown() {
	ts.StopNode()
}

func (ts *NodeTestSuite) InitNode() {
	ctx := context.Background()
	ts.accounts = ts.MakeAccounts(ctx, 0, 2)
	{
		cmd := fmt.Sprintf(`tendermint init-validator --mnemonic="%v" --hdw_path="%v%v"`, testMnemonic, hdwPath, validatorIdx)
		out, e, err := ts.ExecDCommand(ctx, cmd)
		if err != nil {
			ts.FailNowf("failed to call init-validator", "%v: %v: %v", string(out), err, string(e))
		}
	}
	{
		cmd := fmt.Sprintf(`init --address=%v`, ts.accounts[1].String())
		out, e, err := ts.ExecDCommand(ctx, cmd)
		if err != nil {
			ts.FailNowf("failed to call init", "%v: %v: %v", string(out), err, string(e))
		}
	}
}

func (ts *NodeTestSuite) MakeAccounts(ctx context.Context, start, end int) []common.Address {
	var addrs []common.Address
	for i := start; i <= end; i++ {
		cmd := fmt.Sprintf(`new --password=password --silent --mnemonic="%v" --hdw_path="%v%v"`, testMnemonic, hdwPath, i)
		if out, e, err := ts.ExecCLICommand(ctx, cmd); err != nil {
			ts.FailNowf("failed to make account", "%v: %v: %v", string(out), err, string(e))
		} else {
			s := strings.TrimRight(string(out), "\n")
			addrs = append(addrs, common.HexToAddress(s))
		}
	}
	return addrs
}

func (ts *NodeTestSuite) RunNode() {
	args, err := parseCmdf(`start --log_level="*:error" --home=%v`, ts.hmdHomeDir)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(ts.hmdPath, args...)
	params := fmt.Sprintf("consensus.timeout_commit=%v,rpc.max_body_bytes=10000000,mempool.max_tx_bytes=10000000", ts.TimeoutCommit.String())
	cmd.Env = []string{WithEnv("TM_PARAMS", params).String()}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		ts.FailNowf("failed to start node", "%v", err)
	}
	ts.hmdProcess = cmd.Process
	// wait for first block to be created
	time.Sleep(4 * ts.TimeoutCommit)
}

func (ts *NodeTestSuite) StopNode() {
	if ts.hmdProcess != nil {
		if err := ts.hmdProcess.Kill(); err != nil {
			panic(err)
		}
	}
}

func (ts *NodeTestSuite) ExecDCommand(ctx context.Context, cmdStr string, envs ...Env) ([]byte, []byte, error) {
	return ts.ExecCommand(ctx, ts.hmdPath, cmdStr+" --home="+ts.hmdHomeDir, envs...)
}

func (ts *NodeTestSuite) ExecCLICommand(ctx context.Context, cmdStr string, envs ...Env) ([]byte, []byte, error) {
	return ts.ExecCommand(ctx, ts.hmcliPath, cmdStr+" --home="+ts.hmcliHomeDir, envs...)
}

func (ts *NodeTestSuite) ExecCommand(ctx context.Context, binName, cmdStr string, envs ...Env) ([]byte, []byte, error) {
	ts.T().Logf("command: %v %v", binName, cmdStr)
	args, err := parseCmd(cmdStr)
	if err != nil {
		return nil, nil, err
	}
	cmd := exec.CommandContext(ctx, binName, args...)
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	for _, e := range envs {
		cmd.Env = append(cmd.Env, e.String())
	}
	err = cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

func (ts *NodeTestSuite) Account(idx int) common.Address {
	return ts.accounts[idx]
}

func (ts *NodeTestSuite) GetNodeConfig() *cfg.Config {
	viper.AddConfigPath(ts.hmdHomeDir + "/config")
	viper.ReadInConfig()
	c, err := config.GetConfig(ts.hmdHomeDir)
	if err != nil {
		ts.FailNow("config not found", err.Error())
	}
	return c
}

func (ts *NodeTestSuite) RPCClient() rpclient.Client {
	c := ts.GetNodeConfig()
	fmt.Println("address is ", c.RPC.ListenAddress)
	return rpclient.NewHTTP(c.RPC.ListenAddress, "/websocket")
}

func (ts *NodeTestSuite) GetCLIHomeDir() string {
	return ts.hmcliHomeDir
}

type Env interface {
	String() string
}

type env struct {
	key   string
	value string
}

func (e env) String() string {
	return fmt.Sprintf("%v=%v", e.key, e.value)
}

func WithEnv(key, value string) Env {
	return env{key, value}
}

func parseCmd(s string) ([]string, error) {
	return shellwords.Parse(s)
}

func parseCmdf(s string, args ...interface{}) ([]string, error) {
	cmd := fmt.Sprintf(s, args...)
	return parseCmd(cmd)
}
