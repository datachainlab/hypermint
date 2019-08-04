package transaction

import (
	"crypto/ecdsa"
	"fmt"
	"testing"
	"time"

	clihelper "github.com/bluele/hypermint/pkg/client/helper"
	"github.com/bluele/hypermint/pkg/transaction"
	icommon "github.com/bluele/hypermint/tests/integration/common"
	"github.com/bluele/hypermint/tests/integration/helper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	mnemonic = "token dash time stand brisk fatal health honey frozen brown flight kitchen"
	password = "password"
)

type TransferTestSuite struct {
	icommon.NodeTestSuite
	owner *ecdsa.PrivateKey
	alice *ecdsa.PrivateKey
	bob   *ecdsa.PrivateKey
}

func (ts *TransferTestSuite) SetupSuite() {
	ts.owner = helper.GetPrivKey(nil, mnemonic, "m/44'/60'/0'/0/0")
	ts.NodeTestSuite.SetupSuite(crypto.PubkeyToAddress(ts.owner.PublicKey))
	_, err := ts.KS.ImportECDSA(ts.owner, password)
	ts.NoError(err)

	ts.alice = helper.GetPrivKey(ts.KS, mnemonic, "m/44'/60'/0'/0/1")
	ts.bob = helper.GetPrivKey(ts.KS, mnemonic, "m/44'/60'/0'/0/2")
	viper.Set(clihelper.FlagPassword, password)
}

func (ts *TransferTestSuite) TestTransfer() {
	ownerAddr := crypto.PubkeyToAddress(ts.owner.PublicKey)
	aliceAddr := crypto.PubkeyToAddress(ts.alice.PublicKey)

	// check if initial state is valid
	ob, err := ts.GetNodeClientContext(ts.CliDir, ownerAddr).GetBalanceByAddress(ownerAddr)
	ts.NoError(err)
	ts.EqualValues(100, ob)

	var steps = []struct {
		sender          common.Address
		receiver        common.Address
		amount          uint64
		senderBalance   uint64
		receiverBalance uint64
		hasError        bool
	}{
		{aliceAddr, ownerAddr, 10, 0, 0, true},
		{ownerAddr, aliceAddr, 10, 90, 10, false},
		{ownerAddr, aliceAddr, 20, 70, 30, false},
		{aliceAddr, ownerAddr, 10, 20, 80, false},
		{ownerAddr, aliceAddr, 81, 0, 0, true},
		{aliceAddr, ownerAddr, 20, 0, 100, false},
		{aliceAddr, ownerAddr, 10, 0, 0, true},
	}

	for i, s := range steps {
		ts.T().Run(fmt.Sprint(i), func(t *testing.T) {
			ctx := ts.GetNodeClientContext(ts.CliDir, s.sender)
			tx := &transaction.TransferTx{
				Common: transaction.CommonTx{
					Code:  transaction.TRANSFER,
					From:  s.sender,
					Gas:   1,
					Nonce: uint64(time.Now().UnixNano()),
				},
				To:     s.receiver,
				Amount: s.amount,
			}

			if err := ctx.SignAndBroadcastTx(tx, s.sender); s.hasError {
				ts.Error(err)
				return
			} else {
				ts.NoError(err)
			}

			time.Sleep(2 * ts.Config.Consensus.TimeoutCommit)

			sb, err := ctx.GetBalanceByAddress(s.sender)
			ts.NoError(err)
			ts.EqualValues(s.senderBalance, sb)

			rb, err := ctx.GetBalanceByAddress(s.receiver)
			ts.NoError(err)
			ts.EqualValues(s.receiverBalance, rb)
		})
	}
}

func (ts *TransferTestSuite) TearDownSuite() {
	ts.NodeTestSuite.TearDownSuite()
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
