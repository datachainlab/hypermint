package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	gcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func PrvKeyToCryptoKey(prv *ecdsa.PrivateKey) crypto.PrivKey {
	var p secp256k1.PrivKeySecp256k1
	copy(p[:], prv.D.Bytes())
	return p
}

func validateEcrecoverParams(hash, v, r, s []byte) error {
	if len(v) != 1 {
		return fmt.Errorf("length of v should be 1, got %v", len(v))
	} else if len(r) != 32 {
		return fmt.Errorf("length of r should be 32, got %v", len(r))
	} else if len(s) != 32 {
		return fmt.Errorf("length of s should be 32, got %v", len(s))
	} else if len(hash) != 32 {
		return fmt.Errorf("length of hash should be 32, got %v", len(hash))
	}
	return nil
}

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, v, r, s []byte) ([]byte, error) {
	if err := validateEcrecoverParams(hash, v, r, s); err != nil {
		return nil, err
	}
	sig := make([]byte, 65)
	copy(sig[0:32], r)
	copy(sig[32:64], s)

	vv := v[0]
	if vv >= 27 {
		vv -= 27
	}
	sig[64] = vv
	return gcrypto.Ecrecover(hash, sig)
}

func EcrecoverAddress(hash, v, r, s []byte) (common.Address, error) {
	s, err := Ecrecover(hash, v, r, s)
	if err != nil {
		return common.Address{}, err
	}

	x, y := elliptic.Unmarshal(gcrypto.S256(), s)
	pub := &ecdsa.PublicKey{Curve: gcrypto.S256(), X: x, Y: y}
	return gcrypto.PubkeyToAddress(*pub), nil
}
