package contract

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/bluele/hypermint/pkg/abci/store"
	sdk "github.com/bluele/hypermint/pkg/abci/types"
	"github.com/bluele/hypermint/pkg/contract"
	"github.com/bluele/hypermint/pkg/db"
	"github.com/bluele/hypermint/pkg/util"
	"github.com/bluele/hypermint/pkg/util/wallet"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tm-db"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	testMnemonic = "math razor capable expose worth grape metal sunset metal sudden usage scheme"
	hdwPath      = "m/44'/60'/0'/0/"

	testContractPath = "../../build/contract_test.wasm"
)

type ContractTestSuite struct {
	suite.Suite

	contract contract.Contract
	owner    *ecdsa.PrivateKey

	mainKey     *sdk.KVStoreKey
	cmsProvider func() store.CommitMultiStore
}

func (ts *ContractTestSuite) SetupTest() {
	assert := ts.Assert()
	b, err := ioutil.ReadFile(testContractPath)
	assert.NoError(err)

	ts.owner, err = ts.GetPrvkey(0)
	assert.NoError(err)
	ts.mainKey = sdk.NewKVStoreKey("main")
	ts.cmsProvider = func() store.CommitMultiStore {
		cms := store.NewCommitMultiStore(dbm.NewMemDB())
		cms.MountStoreWithDB(ts.mainKey, sdk.StoreTypeIAVL, nil)
		assert.NoError(cms.LoadLatestVersion())
		return cms
	}
	ts.contract = contract.Contract{
		Owner: crypto.PubkeyToAddress(ts.owner.PublicKey),
		Code:  b,
	}
}

func (ts *ContractTestSuite) GetPrvkey(index uint32) (*ecdsa.PrivateKey, error) {
	hp, err := wallet.ParseHDPathLevel("m/44'/60'/0'/0/" + fmt.Sprint(index))
	if err != nil {
		return nil, err
	}
	return wallet.GetPrvKeyFromHDWallet(bip39.NewSeed(testMnemonic, ""), hp)
}

func (ts *ContractTestSuite) TestKeccak256() {
	assert := ts.Assert()
	cms := ts.cmsProvider()

	msg := common.RandBytes(32)
	args := contract.NewArgs([][]byte{msg})

	env := &contract.Env{
		Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
		Contract: &ts.contract,
		DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
		Args:     args,
	}
	res, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_keccak256")
	assert.NoError(err)
	h, err := util.Keccak256(msg)
	assert.NoError(err)
	assert.Equal(h, res.Response)
}

func (ts *ContractTestSuite) TestSha256() {
	assert := ts.Assert()
	cms := ts.cmsProvider()

	msg := common.RandBytes(32)
	args := contract.NewArgs([][]byte{msg})

	env := &contract.Env{
		Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
		Contract: &ts.contract,
		DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
		Args:     args,
	}
	res, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_sha256")
	assert.NoError(err)
	h := util.Sha256(msg)
	assert.Equal(h, res.Response)
}

func (ts *ContractTestSuite) TestECRecover() {
	msg := common.RandBytes(32)
	var makeMsgHash = func(idx uint8) []byte {
		b := make([]byte, 32)
		copy(b[:], msg)
		b[len(b)-1] = idx
		return crypto.Keccak256(b)
	}

	var cases = []struct {
		signer      uint32
		sender      uint32
		signHashIdx uint8
		argHashIdx  uint8
		hasError    bool
	}{
		{0, 0, 0, 0, false},
		{1, 1, 0, 0, false},
		{1, 0, 0, 0, true},
		{0, 0, 1, 1, false},
		{0, 0, 0, 1, true},
		{0, 1, 0, 1, true},
	}

	for i, cs := range cases {
		ts.T().Run(fmt.Sprint(i), func(t *testing.T) {
			assert := assert.New(t)
			cms := ts.cmsProvider()
			var args contract.Args

			signer, err := ts.GetPrvkey(cs.signer)
			assert.NoError(err)
			sender, err := ts.GetPrvkey(cs.sender)
			assert.NoError(err)
			sh := makeMsgHash(cs.signHashIdx)
			ah := makeMsgHash(cs.argHashIdx)

			sig, err := crypto.Sign(sh, signer)
			assert.NoError(err)

			args.PushBytes(ah)
			args.PushBytes(sig)

			env := &contract.Env{
				Sender:   crypto.PubkeyToAddress(sender.PublicKey),
				Contract: &ts.contract,
				DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
				Args:     args,
			}
			res, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "check_signature")
			if cs.hasError {
				assert.Error(err)
				assert.EqualValues(-1, res.Code)
			} else {
				assert.NoError(err)
				assert.EqualValues(0, res.Code)
			}
		})
	}
}

func (ts *ContractTestSuite) TestCannotReadUncommittedState() {
	assert := ts.Assert()
	cms := ts.cmsProvider()

	env := &contract.Env{
		Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
		Contract: &ts.contract,
		DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
	}
	_, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_read_uncommitted_state")
	assert.Error(err)
}

func (ts *ContractTestSuite) TestReadWriteState() {
	assert := ts.Assert()
	cms := ts.cmsProvider()

	{ // Write a value to state
		env := &contract.Env{
			Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
			Contract: &ts.contract,
			DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
			Args:     contract.NewArgsFromStrings([]string{"key", "value"}),
		}
		_, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_write_state")
		assert.NoError(err)
	}
	cms.Commit()

	{ // Check the value on state
		env := &contract.Env{
			Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
			Contract: &ts.contract,
			DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{2, 1}),
			Args:     contract.NewArgsFromStrings([]string{"key"}),
		}
		_, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_read_state")
		assert.NoError(err)
	}
}

func (ts *ContractTestSuite) TestEmitEvent() {
	assert := ts.Assert()
	cms := ts.cmsProvider()
	msg0 := common.RandBytes(32)
	msg1 := common.RandBytes(32)
	args := contract.NewArgs([][]byte{msg0, msg1})

	env := &contract.Env{
		Sender:   crypto.PubkeyToAddress(ts.owner.PublicKey),
		Contract: &ts.contract,
		DB:       db.NewVersionedDB(cms.GetKVStore(ts.mainKey), db.Version{1, 1}),
		Args:     args,
	}
	res, err := env.Exec(sdk.NewContext(cms, abci.Header{}, false, nil), "test_emit_event")
	assert.NoError(err)
	assert.Equal(2, len(res.Events))

	assert.Equal([]byte("test-event-name-0"), res.Events[0].Name)
	assert.Equal(msg0, res.Events[0].Value)

	assert.Equal([]byte("test-event-name-1"), res.Events[1].Name)
	assert.Equal(msg1, res.Events[1].Value)
}

func TestContractTestSuite(t *testing.T) {
	suite.Run(t, new(ContractTestSuite))
}
