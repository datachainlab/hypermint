package app

import (
	"bytes"
	"encoding/json"

	amino "github.com/tendermint/go-amino"
)

var cdc *amino.Codec

func GetCodec() *amino.Codec {
	return cdc
}

// attempt to make some pretty json
func MarshalJSONIndent(cdc *amino.Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func init() {
	cdc = amino.NewCodec()
}
