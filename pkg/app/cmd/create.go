package cmd

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/bluele/hypermint/pkg/util/wallet"
	"github.com/bluele/hypermint/pkg/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	flagGenesisConfig = "genesis"
)

func createCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new node from genesis.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())
			config := ctx.Config
			genesisPath := viper.GetString(flagGenesisConfig)
			if _, err := os.Stat(genesisPath); err != nil && !os.IsExist(err) {
				return fmt.Errorf("genesis path not found path=%v err=%v", genesisPath, err)
			}

			var prv *ecdsa.PrivateKey
			var err error
			if mnemonic := viper.GetString(flagMnemonic); mnemonic != "" {
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
				hp, err := wallet.ParseHDPathLevel(viper.GetString(flagHDWPath))
				if err != nil {
					return err
				}
				seed := bip39.NewSeed(mnemonic, "")
				prv, err = wallet.GetPrvKeyFromHDWallet(seed, hp)
				if err != nil {
					return err
				}
			} else {
				prv, err = crypto.GenerateKey()
				if err != nil {
					return err
				}
			}
			validator.GenFilePVWithECDSA(config.PrivValidatorFile(), prv)
			return util.CopyFile(genesisPath, filepath.Join(config.RootDir, "config/genesis.json"))
		},
	}

	cmd.Flags().String(flagGenesisConfig, "", "path for genesis config")
	cmd.Flags().String(flagMnemonic, "", "mnemonic string")
	cmd.Flags().String(flagHDWPath, "", "HD Wallet path")
	util.CheckRequiredFlag(cmd, flagGenesisConfig)
	return cmd
}
