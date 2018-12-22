package wallet

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	bip39 "github.com/tyler-smith/go-bip39"
)

func TestWallet(t *testing.T) {
	var cases = []struct {
		mnemonic string
		path     string
		address  string
	}{
		{
			"math razor capable expose worth grape metal sunset metal sudden usage scheme",
			"m/44'/60'/0'/0/0",
			"0xa89F47C6b463f74d87572b058427dA0A13ec5425",
		},
		{
			"math razor capable expose worth grape metal sunset metal sudden usage scheme",
			"m/44'/60'/0'/0/1",
			"0xcBED645B1C1a6254f1149Df51d3591c6B3803007",
		},
		{
			"math razor capable expose worth grape metal sunset metal sudden usage scheme",
			"m/44'/60'/0'/0/2",
			"0x00731540cd6060991D6B9C57CE295998d9bC2faB",
		},
	}

	for _, cs := range cases {
		hp, err := ParseHDPathLevel(cs.path)
		if err != nil {
			t.Fatal(err)
		}
		if cs.path != hp.String() {
			t.Fatalf("%v != %v", cs.path, hp.String())
		}
		seed := bip39.NewSeed(cs.mnemonic, "")
		prv, err := GetPrvKeyFromHDWallet(seed, hp)
		if err != nil {
			t.Fatal(err)
		}
		addr := crypto.PubkeyToAddress(prv.PublicKey)
		if addr.Hex() != cs.address {
			t.Fatalf("%v != %v", addr.Hex(), cs.address)
		}
	}
}
