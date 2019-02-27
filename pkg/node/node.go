package node

import (
	"fmt"
	"io/ioutil"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/p2p"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

func GenNodeKeyByPrivKey(filePath string, privKey crypto.PrivKey) (*p2p.NodeKey, error) {
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}

	jsonBytes, err := cdc.MarshalJSON(nodeKey)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}

func GenNodeKey(filePath string) (*p2p.NodeKey, error) {
	return GenNodeKeyByPrivKey(filePath, secp256k1.GenPrivKey())
}

func LoadNodeKey(filePath string) (*p2p.NodeKey, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := new(p2p.NodeKey)
	err = cdc.UnmarshalJSON(jsonBytes, nodeKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading NodeKey from %v: %v", filePath, err)
	}
	return nodeKey, nil
}
