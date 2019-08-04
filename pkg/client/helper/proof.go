package helper

import (
	"bytes"
	"fmt"

	"github.com/bluele/hypermint/pkg/app"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/proof"

	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpclient "github.com/tendermint/tendermint/rpc/client"
)

// GetKVProofInfo returns a proof of specified key-value pair existence
func GetKVProofInfo(cli rpclient.Client, contractAddr common.Address, height int64, key, value cmn.HexBytes) (*proof.KVProofInfo, error) {
	path := fmt.Sprintf("/store/%v/key", app.ContractStoreKey.Name())
	res, err := cli.ABCIQueryWithOptions(
		path,
		append(contractAddr.Bytes(), key.Bytes()...),
		rpclient.ABCIQueryOptions{
			Height: height,
			Prove:  true,
		},
	)
	if err != nil {
		return nil, err
	}
	vo, err := db.BytesToValueObject(res.Response.Value)
	if err != nil {
		return nil, err
	}
	if value != nil && !bytes.Equal(value, vo.Value) {
		return nil, fmt.Errorf("value is mismatch: %v(%v) != %v(%v)",
			string(value), value.Bytes(),
			string(vo.Value), vo.Value,
		)
	}

	h := res.Response.Height + 1
	c, err := cli.Commit(&h)
	if err != nil {
		return nil, err
	}
	header := c.SignedHeader.Header
	op, err := proof.MakeKVProofOp(header)
	if err != nil {
		return nil, err
	}
	p := res.Response.Proof
	p.Ops = append(p.Ops, op)

	kvp := proof.MakeKVProofInfo(
		header.Height,
		p,
		contractAddr,
		key,
		vo,
	)
	if err := kvp.VerifyWithHeader(header); err != nil {
		return nil, err
	}
	return kvp, nil
}
