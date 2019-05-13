package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bgentry/speakeasy"
	"github.com/ethereum/go-ethereum/common"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
)

const (
	// Minimum acceptable password length
	MinPassLength = 4

	// Flags
	FlagHomeDir  = "home"
	FlagVerbose  = "verbose"
	FlagNode     = "node"
	FlagAddress  = "address"
	FlagPassword = "password"
)

// Allows for reading prompts for stdin
func BufferStdin() *bufio.Reader {
	return bufio.NewReader(os.Stdin)
}

// Prompts for a password one-time
// Enforces minimum password length
func GetPassword(prompt string, buf *bufio.Reader) (pass string, err error) {
	if inputIsTty() {
		pass, err = speakeasy.Ask(prompt)
	} else {
		pass, err = readLineFromBuf(buf)
	}
	if err != nil {
		return "", err
	}
	if len(pass) < MinPassLength {
		return "", fmt.Errorf("Password must be at least %d characters", MinPassLength)
	}
	return pass, nil
}

// Prompts for a password twice to verify they match
func GetCheckPassword(prompt, prompt2 string, buf *bufio.Reader) (string, error) {
	if !inputIsTty() {
		return GetPassword(prompt, buf)
	}

	pass, err := GetPassword(prompt, buf)
	if err != nil {
		return "", err
	}
	pass2, err := GetPassword(prompt2, buf)
	if err != nil {
		return "", err
	}
	if pass != pass2 {
		return "", errors.New("Passphrases did not match")
	}
	return pass, nil
}

// Convert string to Ethereum Address
func StrToAddress(addrStr string) (common.Address, error) {
	if !common.IsHexAddress(strings.TrimSpace(addrStr)) {
		return common.Address{}, errors.New("invalid address provided, please use hex format")
	}
	return common.HexToAddress(addrStr), nil
}

// Returns true iff we have an interactive prompt
func inputIsTty() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

// reads one line from stdin
func readLineFromBuf(buf *bufio.Reader) (string, error) {
	pass, err := buf.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(pass), nil
}

func ParseAddrs(addrsStr string) ([]common.Address, error) {
	var addrs []common.Address
	as := strings.Split(addrsStr, ",")
	for _, a := range as {
		a = strings.TrimSpace(a)
		if a == "" {
			break
		}
		addr, err := StrToAddress(a)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}
