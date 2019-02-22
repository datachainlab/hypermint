package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bip39 "github.com/tyler-smith/go-bip39"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/util/wallet"
	"github.com/bluele/hypermint/pkg/validator"
)

const (
	flagMnemonic = "mnemonic"
	flagHDWPath  = "hdw_path"
)

func validatorCmd(ctx *app.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init-validator",
		Short: "initialize a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			viper.BindPFlags(cmd.Flags())

			cfg := ctx.Config
			mnemonic := viper.GetString(flagMnemonic)
			if !bip39.IsMnemonicValid(mnemonic) {
				return errors.New("invalid mnemonic")
			}
			hp, err := wallet.ParseHDPathLevel(viper.GetString(flagHDWPath))
			if err != nil {
				return err
			}
			seed := bip39.NewSeed(mnemonic, "")
			prv, err := wallet.GetPrvKeyFromHDWallet(seed, hp)
			if err != nil {
				return err
			}
			validator.GenFilePVWithECDSA(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile(), prv)
			return nil
		},
	}
	cmd.Flags().String(flagMnemonic, "", "mnemonic string")
	cmd.Flags().String(flagHDWPath, "", "HD Wallet path")
	return cmd
}
