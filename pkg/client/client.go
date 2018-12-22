package client

import (
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"

	"github.com/bluele/hypermint/pkg/client/context"
	"github.com/bluele/hypermint/pkg/client/helper"
)

// Return a new context with parameters from the command line
func NewClientContextFromViper() (*context.Context, error) {
	nodeURI := viper.GetString(helper.FlagNode)
	var rpc client.Client
	if nodeURI != "" {
		rpc = client.NewHTTP(nodeURI, "/websocket")
	}
	addrs, err := helper.ParseAddrs(viper.GetString(helper.FlagAddress))
	if err != nil {
		return nil, err
	}
	return &context.Context{
		HomeDir:        viper.GetString(helper.FlagHomeDir),
		Verbose:        viper.GetBool(helper.FlagVerbose),
		InputAddresses: addrs,
		NodeURI:        nodeURI,
		Client:         rpc,
	}, nil
}
