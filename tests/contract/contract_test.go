package contract

import (
	"io/ioutil"
	"testing"

	"github.com/bluele/hypermint/pkg/abci/store"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/util/wallet"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	testMnemonic = "math razor capable expose worth grape metal sunset metal sudden usage scheme"
)

func TestContract(t *testing.T) {
	assert := assert.New(t)
	contractPath := "./contract-test/target/wasm32-unknown-unknown/debug/contract_test.wasm"
	b, err := ioutil.ReadFile(contractPath)
	assert.NoError(err)

	hp, err := wallet.ParseHDPathLevel("m/44'/60'/0'/0/0")
	assert.NoError(err)
	prv, err := wallet.GetPrvKeyFromHDWallet(bip39.NewSeed(testMnemonic, ""), hp)
	assert.NoError(err)

	mdb := dbm.NewMemDB()
	defer mdb.Close()
	cms := store.NewCommitMultiStore(mdb)
	var key = sdk.NewKVStoreKey("main")
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	assert.NoError(cms.LoadLatestVersion())
	kvs := cms.GetKVStore(key)

	from := crypto.PubkeyToAddress(prv.PublicKey)
	var args contract.Args

	msgHash := crypto.Keccak256(common.RandBytes(32))
	sig, err := crypto.Sign(msgHash, prv)
	assert.NoError(err)

	args.PushBytes(msgHash)
	args.PushBytes(sig)

	env := &contract.Env{
		Sender: from,
		Contract: &contract.Contract{
			Owner: from,
			Code:  b,
		},
		DB:   db.NewVersionedDB(kvs, db.Version{1, 1}),
		Args: args,
	}
	c := sdk.NewContext(cms, abci.Header{}, false, nil)
	_, err = env.Exec(c, "check_signature")
	assert.NoError(err)
}
