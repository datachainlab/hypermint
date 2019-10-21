// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/bluele/hypermint/pkg/account/abi/bind"
	"github.com/bluele/hypermint/pkg/transaction"
	"github.com/ethereum/go-ethereum/common"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bind.Bind
	_ = common.Big1
	_ = transaction.ContractInitFunc
	_ = bytes.NewBuffer
	_ = binary.Read
	_ = json.NewEncoder
)

var TokenABI string = "{\"functions\":[{\"name\":\"get_balance\",\"type\":\"function\",\"simulate\":true,\"inputs\":[],\"outputs\":[{\"type\":\"i64\",\"name\":\"_\"}]},{\"name\":\"transfer\",\"type\":\"function\",\"simulate\":false,\"inputs\":[{\"type\":\"address\",\"name\":\"to\"},{\"type\":\"i64\",\"name\":\"amount\"}],\"outputs\":[{\"type\":\"i64\",\"name\":\"_\"}]},{\"name\":\"init\",\"type\":\"function\",\"simulate\":false,\"inputs\":[],\"outputs\":[{\"type\":\"bytes\",\"name\":\"_\"}]}],\"events\":null,\"structs\":null}"

var TokenEventDecoder = bind.NewEventDecoder()

func init() {
}

type TokenContract interface {
	GetBalance(opts *bind.TransactOpts) (int64, error)

	Transfer(opts *bind.TransactOpts, to common.Address, amount int64) (*bind.SyncResult, error)
	TransferCommit(opts *bind.TransactOpts, to common.Address, amount int64) (*bind.CommitResult, error)

	Init(opts *bind.TransactOpts) (*bind.SyncResult, error)
	InitCommit(opts *bind.TransactOpts) (*bind.CommitResult, error)
}

type Token struct {
	TokenSimulator  // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
}

type TokenSimulator struct {
	contract *bind.BoundContract
}

type TokenTransactor struct {
	contract *bind.BoundContract
}

type TokenRaw struct {
	Contract *Token
}

type TokenSimulatorRaw struct {
	Contract *TokenSimulator
}

type TokenTransactorRaw struct {
	Contract *TokenTransactor
}

func NewToken(address common.Address, backend bind.ContractBackend) (*Token, error) {
	contract, err := bindToken(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Token{TokenSimulator: TokenSimulator{contract: contract}, TokenTransactor: TokenTransactor{contract: contract}}, nil
}

func NewTokenSimulator(address common.Address, simulator bind.ContractSimulator) (*TokenSimulator, error) {
	contract, err := bindToken(address, simulator, nil)
	if err != nil {
		return nil, err
	}
	return &TokenSimulator{contract: contract}, nil
}

func NewTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenTransactor, error) {
	contract, err := bindToken(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &TokenTransactor{contract: contract}, nil
}

func bindToken(address common.Address, simulator bind.ContractSimulator, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	return bind.NewBoundContract(address, simulator, transactor), nil
}

func (_Token *TokenRaw) Transact(opts *bind.TransactOpts, fn string, params ...[]byte) (*bind.SyncResult, error) {
	return _Token.Contract.TokenTransactor.contract.Transact(opts, fn, params...)
}

func (_Token *TokenSimulatorRaw) Simulate(opts *bind.TransactOpts, fn string, params ...[]byte) (*bind.SimulateResult, error) {
	return _Token.Contract.contract.Simulate(opts, fn, params...)
}

func (_Token *TokenSimulator) GetBalance(opts *bind.TransactOpts) (int64, error) {
	result, err := _Token.contract.Simulate(opts, "get_balance", bind.Args()...)
	if err != nil {
		return int64(0), err
	}

	buf := bytes.NewBuffer(result.Data)

	var v0 int64
	binary.Read(buf, binary.BigEndian, &v0)

	return v0, nil
}

func (_Token *TokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount int64) (*bind.SyncResult, error) {
	return _Token.contract.Transact(opts, "transfer", bind.Args(
		bind.Address(to),
		bind.I64(amount))...)
}

func (_Token *TokenTransactor) TransferCommit(opts *bind.TransactOpts, to common.Address, amount int64) (*bind.CommitResult, error) {
	return _Token.contract.TransactCommit(opts, "transfer", bind.Args(
		bind.Address(to),
		bind.I64(amount))...)
}

func (_Token *TokenTransactor) Init(opts *bind.TransactOpts) (*bind.SyncResult, error) {
	return _Token.contract.Transact(opts, "init", bind.Args()...)
}

func (_Token *TokenTransactor) InitCommit(opts *bind.TransactOpts) (*bind.CommitResult, error) {
	return _Token.contract.TransactCommit(opts, "init", bind.Args()...)
}
