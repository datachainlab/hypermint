package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/config"
	"github.com/bluele/hypermint/pkg/node"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/bluele/hypermint/pkg/validator"
)

//Parameter names, for init gen-tx command
var (
	FlagName       = "name"
	FlagClientHome = "home-client"
	FlagOWK        = "owk"
)

//parameter names, init command
var (
	FlagOverwrite = "overwrite"
	FlagWithTxs   = "with-txs"
	FlagIP        = "ip"
	FlagChainID   = "chain-id"
)

// genesis piece structure for creating combined genesis
type GenesisTx struct {
	NodeID    string                   `json:"node_id"`
	IP        string                   `json:"ip"`
	Validator tmtypes.GenesisValidator `json:"validator"`
	AppGenTx  json.RawMessage          `json:"app_gen_tx"`
}

// Storage for init command input parameters
type InitConfig struct {
	ChainID     string
	GenTxs      bool
	GenTxsDir   string
	Overwrite   bool
	GenesisTime time.Time
}

// get cmd to initialize all files for tendermint and application
func GenTxCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen-tx",
		Short: "Create genesis transaction file (under [--home]/config/gentx/gentx-[nodeID].json)",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {

			c := ctx.Config
			c.SetRoot(viper.GetString(tmcli.HomeFlag))

			ip := viper.GetString(FlagIP)
			if len(ip) == 0 {
				eip, err := externalIP()
				if err != nil {
					return err
				}
				ip = eip
			}

			genTxConfig := config.GenTx{
				viper.GetString(FlagName),
				viper.GetString(FlagClientHome),
				viper.GetBool(FlagOWK),
				ip,
			}
			cliPrint, genTxFile, err := gentxWithConfig(cdc, appInit, c, genTxConfig)
			if err != nil {
				return err
			}
			toPrint := struct {
				AppMessage json.RawMessage `json:"app_message"`
				GenTxFile  json.RawMessage `json:"gen_tx_file"`
			}{
				cliPrint,
				genTxFile,
			}
			out, err := app.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().String(FlagIP, "", "external facing IP to use if left blank IP will be retrieved from this machine")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx)
	return cmd
}

// get cmd to initialize all files for tendermint and application
func initCmd(ctx *app.Context, cdc *amino.Codec, appInit app.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(tmcli.HomeFlag))
			initConfig := InitConfig{
				viper.GetString(FlagChainID),
				viper.GetBool(FlagWithTxs),
				filepath.Join(config.RootDir, "config", "gentx"),
				viper.GetBool(FlagOverwrite),
				tmtime.Now(),
			}

			chainID, nodeID, appMessage, err := InitWithConfig(cdc, appInit, config, initConfig)
			if err != nil {
				return err
			}
			// print out some key information
			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}
			out, err := app.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}
			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().String(FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool(FlagWithTxs, false, "apply existing genesis transactions from [--home]/config/gentx/")
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenState)
	cmd.Flags().AddFlagSet(appInit.FlagsAppGenTx) // need to add this flagset for when no GenTx's provided
	cmd.AddCommand(GenTxCmd(ctx, cdc, appInit))
	return cmd
}

// append a genesis-piece
func processGenTxs(genTxsDir string, cdc *amino.Codec) (
	validators []tmtypes.GenesisValidator, appGenTxs []json.RawMessage, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return
	}

	genTxs := make(map[string]GenesisTx)
	var nodeIDs []string
	for _, fo := range fos {
		filename := path.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (path.Ext(filename) != ".json") {
			continue
		}

		// get the genTx
		var bz []byte
		bz, err = ioutil.ReadFile(filename)
		if err != nil {
			return
		}
		var genTx GenesisTx
		err = cdc.UnmarshalJSON(bz, &genTx)
		if err != nil {
			return
		}

		genTxs[genTx.NodeID] = genTx
		nodeIDs = append(nodeIDs, genTx.NodeID)
	}

	sort.Strings(nodeIDs)

	for _, nodeID := range nodeIDs {
		genTx := genTxs[nodeID]

		// combine some stuff
		validators = append(validators, genTx.Validator)
		appGenTxs = append(appGenTxs, genTx.AppGenTx)

		// Add a persistent peer
		comma := ","
		if len(persistentPeers) == 0 {
			comma = ""
		}
		persistentPeers += fmt.Sprintf("%s%s@%s:26656", comma, genTx.NodeID, genTx.IP)
	}

	return
}

func InitWithConfig(cdc *amino.Codec, appInit app.AppInit, c *cfg.Config, initConfig InitConfig) (
	chainID string, nodeID string, appMessage json.RawMessage, err error) {

	nodeKey, err := node.LoadNodeKey(c.NodeKeyFile())
	if err != nil {
		return
	}
	nodeID = string(nodeKey.ID())

	if initConfig.ChainID == "" {
		initConfig.ChainID = fmt.Sprintf("test-chain-%v", cmn.RandStr(6))
	}
	chainID = initConfig.ChainID

	genFile := c.GenesisFile()
	if !initConfig.Overwrite && cmn.FileExists(genFile) {
		err = fmt.Errorf("genesis.json file already exists: %v", genFile)
		return
	}

	// process genesis transactions, or otherwise create one for defaults
	var appGenTxs []json.RawMessage
	var validators []types.GenesisValidator
	var persistentPeers string

	if initConfig.GenTxs {
		validators, appGenTxs, persistentPeers, err = processGenTxs(initConfig.GenTxsDir, cdc)
		if err != nil {
			return
		}
		c.P2P.PersistentPeers = persistentPeers
		config.SaveConfig(c)
	} else {
		genTxConfig := config.GenTx{
			viper.GetString(FlagName),
			viper.GetString(FlagClientHome),
			viper.GetBool(FlagOWK),
			"127.0.0.1",
		}

		// Write updated config with moniker
		// c.Moniker = genTxConfig.Name
		// config.SaveConfig(c)
		appGenTx, am, validator, err := appInit.AppGenTx(cdc, nodeKey.PubKey(), genTxConfig)
		appMessage = am
		if err != nil {
			return "", "", nil, err
		}
		validators = []types.GenesisValidator{validator}
		appGenTxs = []json.RawMessage{appGenTx}
	}

	appState, err := appInit.AppGenState(cdc, appGenTxs)
	if err != nil {
		return
	}

	err = writeGenesisFile(cdc, genFile, initConfig.ChainID, validators, appState, initConfig.GenesisTime)
	if err != nil {
		return
	}

	return
}

//________________________________________________________________________________________

// writeGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
// nolint: unparam
func writeGenesisFile(cdc *amino.Codec, genesisFile, chainID string, validators []tmtypes.GenesisValidator, appState json.RawMessage, genesisTime time.Time) error {
	genDoc := tmtypes.GenesisDoc{
		GenesisTime: genesisTime,
		ChainID:     chainID,
		Validators:  validators,
		AppState:    appState,
		ConsensusParams: &types.ConsensusParams{
			Block:     types.DefaultBlockParams(),
			Evidence:  types.DefaultEvidenceParams(),
			Validator: types.ValidatorParams{PubKeyTypes: []string{types.ABCIPubKeyTypeSecp256k1}},
		},
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genesisFile)
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// TODO there must be a better way to get external IP
func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}
	return false
}

func addrToIP(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

func gentxWithConfig(cdc *amino.Codec, appInit app.AppInit, config *cfg.Config, genTxConfig config.GenTx) (
	cliPrint json.RawMessage, genTxFile json.RawMessage, err error) {

	pv := validator.GenFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
		secp256k1.GenPrivKey(),
	)
	nodeKey, err := node.GenNodeKeyByPrivKey(config.NodeKeyFile(), pv.Key.PrivKey)
	if err != nil {
		return
	}
	nodeID := string(nodeKey.ID())

	appGenTx, cliPrint, validator, err := appInit.AppGenTx(cdc, pv.GetPubKey(), genTxConfig)
	if err != nil {
		return
	}

	tx := app.GenesisTx{
		NodeID:    nodeID,
		IP:        genTxConfig.IP,
		Validator: validator,
		AppGenTx:  appGenTx,
	}

	bz, err := app.MarshalJSONIndent(cdc, tx)
	if err != nil {
		return
	}
	genTxFile = json.RawMessage(bz)
	name := fmt.Sprintf("gentx-%v.json", nodeID)
	writePath := filepath.Join(config.RootDir, "config", "gentx")
	file := filepath.Join(writePath, name)
	err = cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return
	}
	err = cmn.WriteFile(file, bz, 0644)
	if err != nil {
		return
	}

	// Write updated config with moniker
	config.Moniker = genTxConfig.Name
	configFilePath := filepath.Join(config.RootDir, "config", "config.toml")
	cfg.WriteConfigFile(configFilePath, config)

	return
}

func getIP(i int) (ip string, err error) {
	ip = viper.GetString(startingIPAddress)
	if len(ip) == 0 {
		ip, err = util.ExternalIP()
		if err != nil {
			return "", err
		}
	} else {
		ip, err = calculateIP(ip, i)
		if err != nil {
			return "", err
		}
	}
	return ip, nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)
	err := cmn.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}
	err = cmn.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}
	return nil
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}
	return ipv4.String(), nil
}
