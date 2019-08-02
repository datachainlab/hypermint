package helper

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/app/cmd"
	clictx "github.com/bluele/hypermint/pkg/client/context"
	"github.com/bluele/hypermint/pkg/logger"
	hnode "github.com/bluele/hypermint/pkg/node"
	"github.com/bluele/hypermint/pkg/validator"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tm-db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/rpc/client"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type NodeTestSuite struct {
	suite.Suite
	NodeDir string
	CliDir  string
	node    *node.Node
	Config  *config.Config
	KS      *keystore.KeyStore
}

func (ts *NodeTestSuite) SetupSuite(genesisOwner common.Address) {
	baseDir := path.Join(os.TempDir(), "hypermint-test", cmn.RandStr(8))
	nodeDir := path.Join(baseDir, "node")
	if _, err := os.Stat(nodeDir); os.IsExist(err) {
		os.RemoveAll(nodeDir)
	}
	viper.Set(tmcli.HomeFlag, nodeDir)
	ts.NodeDir = nodeDir

	ctx := &app.Context{Logger: logger.GetDefaultLogger("*:debug")}
	ts.NoError(app.SetupContext(ctx))
	ctx.Config.Consensus.TimeoutCommit = time.Second
	ts.NoError(initValidator(ctx))
	ts.NoError(initChain(ctx.Config, genesisOwner))
	ts.Config = ctx.Config

	nd, err := cmd.StartInProcess(
		ctx,
		app.ConstructAppCreator(newApp, "testapp"),
	)
	ts.NoError(err)
	ts.node = nd

	cliDir := path.Join(baseDir, "cli")
	ts.KS = keystore.NewKeyStore(cliDir, keystore.StandardScryptN, keystore.StandardScryptP)
	ts.CliDir = cliDir
}

func (ts *NodeTestSuite) TearDownSuite() {
	ts.NoError(ts.node.Stop())
}

func (ts *NodeTestSuite) GetNodeClientContext(homeDir string, sender common.Address) *clictx.Context {
	return &clictx.Context{
		HomeDir:        homeDir,
		NodeURI:        ts.Config.RPC.ListenAddress,
		InputAddresses: []common.Address{sender},
		Client:         client.NewHTTP(ts.Config.RPC.ListenAddress, "/websocket"),
		Verbose:        true,
	}
}

func newApp(lg log.Logger, db db.DB, traceStore io.Writer) abci.Application {
	logger.SetLogger(lg)
	return app.NewChain(lg, db, traceStore)
}

func initValidator(ctx *app.Context) error {
	config := ctx.Config
	prv := secp256k1.GenPrivKey()
	pv := validator.GenFilePV(config.PrivValidatorKeyFile(), config.PrivValidatorStateFile(), prv)
	if _, err := hnode.GenNodeKeyByPrivKey(config.NodeKeyFile(), pv.Key.PrivKey); err != nil {
		return err
	}
	return nil
}

func initChain(cfg *config.Config, genesisOwner common.Address) error {
	viper.Set("address", genesisOwner.Hex())
	initConfig := cmd.InitConfig{
		"chainid",
		false,
		filepath.Join(cfg.RootDir, "config", "gentx"),
		false,
		tmtime.Now(),
	}
	_, _, _, err := cmd.InitWithConfig(
		app.GetCodec(),
		app.NewAppInit(),
		cfg,
		initConfig,
	)
	return err
}
