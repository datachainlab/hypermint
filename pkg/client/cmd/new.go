package cmd

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	bip39 "github.com/tyler-smith/go-bip39"

	"github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/util/wallet"
)

const (
	flagPassword = "password"
	flagSilent   = "silent"
	flagMnemonic = "mnemonic"
	flagHDWPath  = "hdw_path"
)

func init() {
	rootCmd.AddCommand(newAccountCmd)
	newAccountCmd.Flags().String(flagPassword, "", "passphrase")
	newAccountCmd.Flags().Bool(flagSilent, false, "silent output")
	newAccountCmd.Flags().String(flagMnemonic, "", "mnemonic string")
	newAccountCmd.Flags().String(flagHDWPath, "", "HD Wallet path")
}

var newAccountCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new account",
	Long:  `Add an encrypted account to the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		dir := viper.GetString(helper.FlagHomeDir)
		ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
		mnemonic := viper.GetString(flagMnemonic)
		hdPath := viper.GetString(flagHDWPath)
		password, err := getPassword()
		if err != nil {
			return err
		}

		var ac *accounts.Account
		if mnemonic != "" && hdPath != "" {
			ac, err = importAccountFromHDW(ks, mnemonic, hdPath, password)
		} else {
			ac, err = createAccountWithPassword(ks, password)
		}
		fmt.Println(ac.Address.Hex())
		return nil
	},
}

func getPassword() (string, error) {
	var err error
	pass := viper.GetString(flagPassword)
	if pass == "" {
		buf := helper.BufferStdin()
		pass, err = helper.GetCheckPassword("Enter a passphrase for your key:", "Repeat the passphrase:", buf)
		if err != nil {
			return "", err
		}
	}
	return pass, nil
}

func createAccountWithPassword(ks *keystore.KeyStore, password string) (*accounts.Account, error) {
	acc, err := ks.NewAccount(password)
	if err != nil {
		return nil, err
	}

	if !viper.GetBool(flagSilent) {
		fmt.Println("\n**Important** do not lose your passphrase.")
		fmt.Println("It is the only way to recover your account")
		fmt.Println("You should export this account and store it in a secure location")
		fmt.Printf("Your new account address is: %s\n", acc.Address.Hex())
	}

	return &acc, nil
}

func importAccountFromHDW(ks *keystore.KeyStore, mnemonic, path, password string) (*accounts.Account, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	hp, err := wallet.ParseHDPathLevel(path)
	if err != nil {
		return nil, err
	}
	prv, err := wallet.GetPrvKeyFromHDWallet(bip39.NewSeed(mnemonic, ""), hp)
	if err != nil {
		return nil, err
	}
	ac, err := ks.ImportECDSA(prv, password)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}
