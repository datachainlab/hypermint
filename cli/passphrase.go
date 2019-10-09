package cli

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func ReadPassphrase() (string, error) {
	_, _ = fmt.Fprint(os.Stderr, "Passphrase: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	_, _ = fmt.Fprintf(os.Stderr, "\n")
	if err != nil {
		return "", err
	} else {
		return string(password), nil
	}
}
