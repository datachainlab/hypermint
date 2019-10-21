package bind

import (
	"fmt"
	"github.com/bluele/hypermint/pkg/account/cli"
	"github.com/ethereum/go-ethereum/common"
	)

func GetFlagPassphrase(getPassphrase func() (string, error)) func() (string, error) {
	return func() (string, error) {
		passphrase, err := getPassphrase()
		if err != nil {
			if pp, err := cli.ReadPassphrase(); err != nil {
				return "", err
			} else {
				return pp, nil
			}
		} else {
			return passphrase, nil
		}
	}
}

func GetFlagAddress(getAddress func() (string, error), alias map[string]common.Address) func() (common.Address, error) {
	return func () (common.Address, error) {
		var address common.Address
		a, err := getAddress()
		if err != nil {
			return common.Address{}, err
		}
		if addr, ok := alias[a]; ok {
			return addr, nil
		} else {
			if !common.IsHexAddress(a) {
				return common.Address{}, fmt.Errorf("invalid address: %v", a)
			}
			address = common.HexToAddress(a)
		}
		return address, nil
	}
}
