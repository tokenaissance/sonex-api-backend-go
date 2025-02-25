// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package common

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// XtokenMetaData contains all meta data concerning the Xtoken contract.
var XtokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"AddressAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedAddresses\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"allowedAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"setAllowedAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// XtokenABI is the input ABI used to generate the binding from.
// Deprecated: Use XtokenMetaData.ABI instead.
var XtokenABI = XtokenMetaData.ABI

// Xtoken is an auto generated Go binding around an Ethereum contract.
type Xtoken struct {
	XtokenCaller     // Read-only binding to the contract
	XtokenTransactor // Write-only binding to the contract
	XtokenFilterer   // Log filterer for contract events
}

// XtokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type XtokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XtokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XtokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XtokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XtokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XtokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XtokenSession struct {
	Contract     *Xtoken           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XtokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XtokenCallerSession struct {
	Contract *XtokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// XtokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XtokenTransactorSession struct {
	Contract     *XtokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XtokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type XtokenRaw struct {
	Contract *Xtoken // Generic contract binding to access the raw methods on
}

// XtokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XtokenCallerRaw struct {
	Contract *XtokenCaller // Generic read-only contract binding to access the raw methods on
}

// XtokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XtokenTransactorRaw struct {
	Contract *XtokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXtoken creates a new instance of Xtoken, bound to a specific deployed contract.
func NewXtoken(address common.Address, backend bind.ContractBackend) (*Xtoken, error) {
	contract, err := bindXtoken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Xtoken{XtokenCaller: XtokenCaller{contract: contract}, XtokenTransactor: XtokenTransactor{contract: contract}, XtokenFilterer: XtokenFilterer{contract: contract}}, nil
}

// NewXtokenCaller creates a new read-only instance of Xtoken, bound to a specific deployed contract.
func NewXtokenCaller(address common.Address, caller bind.ContractCaller) (*XtokenCaller, error) {
	contract, err := bindXtoken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XtokenCaller{contract: contract}, nil
}

// NewXtokenTransactor creates a new write-only instance of Xtoken, bound to a specific deployed contract.
func NewXtokenTransactor(address common.Address, transactor bind.ContractTransactor) (*XtokenTransactor, error) {
	contract, err := bindXtoken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XtokenTransactor{contract: contract}, nil
}

// NewXtokenFilterer creates a new log filterer instance of Xtoken, bound to a specific deployed contract.
func NewXtokenFilterer(address common.Address, filterer bind.ContractFilterer) (*XtokenFilterer, error) {
	contract, err := bindXtoken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XtokenFilterer{contract: contract}, nil
}

// bindXtoken binds a generic wrapper to an already deployed contract.
func bindXtoken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := XtokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Xtoken *XtokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Xtoken.Contract.XtokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Xtoken *XtokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Xtoken.Contract.XtokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Xtoken *XtokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Xtoken.Contract.XtokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Xtoken *XtokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Xtoken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Xtoken *XtokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Xtoken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Xtoken *XtokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Xtoken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Xtoken *XtokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Xtoken *XtokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Xtoken.Contract.Allowance(&_Xtoken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Xtoken *XtokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Xtoken.Contract.Allowance(&_Xtoken.CallOpts, owner, spender)
}

// AllowedAddresses is a free data retrieval call binding the contract method 0x4120657a.
//
// Solidity: function allowedAddresses(address ) view returns(bool)
func (_Xtoken *XtokenCaller) AllowedAddresses(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "allowedAddresses", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedAddresses is a free data retrieval call binding the contract method 0x4120657a.
//
// Solidity: function allowedAddresses(address ) view returns(bool)
func (_Xtoken *XtokenSession) AllowedAddresses(arg0 common.Address) (bool, error) {
	return _Xtoken.Contract.AllowedAddresses(&_Xtoken.CallOpts, arg0)
}

// AllowedAddresses is a free data retrieval call binding the contract method 0x4120657a.
//
// Solidity: function allowedAddresses(address ) view returns(bool)
func (_Xtoken *XtokenCallerSession) AllowedAddresses(arg0 common.Address) (bool, error) {
	return _Xtoken.Contract.AllowedAddresses(&_Xtoken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Xtoken *XtokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Xtoken *XtokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Xtoken.Contract.BalanceOf(&_Xtoken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Xtoken *XtokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Xtoken.Contract.BalanceOf(&_Xtoken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Xtoken *XtokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Xtoken *XtokenSession) Decimals() (uint8, error) {
	return _Xtoken.Contract.Decimals(&_Xtoken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_Xtoken *XtokenCallerSession) Decimals() (uint8, error) {
	return _Xtoken.Contract.Decimals(&_Xtoken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Xtoken *XtokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Xtoken *XtokenSession) Name() (string, error) {
	return _Xtoken.Contract.Name(&_Xtoken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Xtoken *XtokenCallerSession) Name() (string, error) {
	return _Xtoken.Contract.Name(&_Xtoken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Xtoken *XtokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Xtoken *XtokenSession) Owner() (common.Address, error) {
	return _Xtoken.Contract.Owner(&_Xtoken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Xtoken *XtokenCallerSession) Owner() (common.Address, error) {
	return _Xtoken.Contract.Owner(&_Xtoken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Xtoken *XtokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Xtoken *XtokenSession) Symbol() (string, error) {
	return _Xtoken.Contract.Symbol(&_Xtoken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Xtoken *XtokenCallerSession) Symbol() (string, error) {
	return _Xtoken.Contract.Symbol(&_Xtoken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Xtoken *XtokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Xtoken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Xtoken *XtokenSession) TotalSupply() (*big.Int, error) {
	return _Xtoken.Contract.TotalSupply(&_Xtoken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Xtoken *XtokenCallerSession) TotalSupply() (*big.Int, error) {
	return _Xtoken.Contract.TotalSupply(&_Xtoken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Xtoken *XtokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Xtoken *XtokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Approve(&_Xtoken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_Xtoken *XtokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Approve(&_Xtoken.TransactOpts, spender, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_Xtoken *XtokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "mint", to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_Xtoken *XtokenSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Mint(&_Xtoken.TransactOpts, to, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (_Xtoken *XtokenTransactorSession) Mint(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Mint(&_Xtoken.TransactOpts, to, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Xtoken *XtokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Xtoken *XtokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _Xtoken.Contract.RenounceOwnership(&_Xtoken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Xtoken *XtokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Xtoken.Contract.RenounceOwnership(&_Xtoken.TransactOpts)
}

// SetAllowedAddress is a paid mutator transaction binding the contract method 0xfae36986.
//
// Solidity: function setAllowedAddress(address allowedAddress, bool authorized) returns()
func (_Xtoken *XtokenTransactor) SetAllowedAddress(opts *bind.TransactOpts, allowedAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "setAllowedAddress", allowedAddress, authorized)
}

// SetAllowedAddress is a paid mutator transaction binding the contract method 0xfae36986.
//
// Solidity: function setAllowedAddress(address allowedAddress, bool authorized) returns()
func (_Xtoken *XtokenSession) SetAllowedAddress(allowedAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _Xtoken.Contract.SetAllowedAddress(&_Xtoken.TransactOpts, allowedAddress, authorized)
}

// SetAllowedAddress is a paid mutator transaction binding the contract method 0xfae36986.
//
// Solidity: function setAllowedAddress(address allowedAddress, bool authorized) returns()
func (_Xtoken *XtokenTransactorSession) SetAllowedAddress(allowedAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _Xtoken.Contract.SetAllowedAddress(&_Xtoken.TransactOpts, allowedAddress, authorized)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Transfer(&_Xtoken.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.Transfer(&_Xtoken.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.TransferFrom(&_Xtoken.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_Xtoken *XtokenTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Xtoken.Contract.TransferFrom(&_Xtoken.TransactOpts, from, to, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Xtoken *XtokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Xtoken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Xtoken *XtokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Xtoken.Contract.TransferOwnership(&_Xtoken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Xtoken *XtokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Xtoken.Contract.TransferOwnership(&_Xtoken.TransactOpts, newOwner)
}

// XtokenAddressAuthorizedIterator is returned from FilterAddressAuthorized and is used to iterate over the raw logs and unpacked data for AddressAuthorized events raised by the Xtoken contract.
type XtokenAddressAuthorizedIterator struct {
	Event *XtokenAddressAuthorized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XtokenAddressAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XtokenAddressAuthorized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XtokenAddressAuthorized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XtokenAddressAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XtokenAddressAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XtokenAddressAuthorized represents a AddressAuthorized event raised by the Xtoken contract.
type XtokenAddressAuthorized struct {
	ContractAddress common.Address
	Authorized      bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAddressAuthorized is a free log retrieval operation binding the contract event 0x97a3679c37536161ce94576f0df5c286fe717e7a0eec7d467758d7abb6337001.
//
// Solidity: event AddressAuthorized(address indexed contractAddress, bool authorized)
func (_Xtoken *XtokenFilterer) FilterAddressAuthorized(opts *bind.FilterOpts, contractAddress []common.Address) (*XtokenAddressAuthorizedIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Xtoken.contract.FilterLogs(opts, "AddressAuthorized", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &XtokenAddressAuthorizedIterator{contract: _Xtoken.contract, event: "AddressAuthorized", logs: logs, sub: sub}, nil
}

// WatchAddressAuthorized is a free log subscription operation binding the contract event 0x97a3679c37536161ce94576f0df5c286fe717e7a0eec7d467758d7abb6337001.
//
// Solidity: event AddressAuthorized(address indexed contractAddress, bool authorized)
func (_Xtoken *XtokenFilterer) WatchAddressAuthorized(opts *bind.WatchOpts, sink chan<- *XtokenAddressAuthorized, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Xtoken.contract.WatchLogs(opts, "AddressAuthorized", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XtokenAddressAuthorized)
				if err := _Xtoken.contract.UnpackLog(event, "AddressAuthorized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddressAuthorized is a log parse operation binding the contract event 0x97a3679c37536161ce94576f0df5c286fe717e7a0eec7d467758d7abb6337001.
//
// Solidity: event AddressAuthorized(address indexed contractAddress, bool authorized)
func (_Xtoken *XtokenFilterer) ParseAddressAuthorized(log types.Log) (*XtokenAddressAuthorized, error) {
	event := new(XtokenAddressAuthorized)
	if err := _Xtoken.contract.UnpackLog(event, "AddressAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XtokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Xtoken contract.
type XtokenApprovalIterator struct {
	Event *XtokenApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XtokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XtokenApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XtokenApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XtokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XtokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XtokenApproval represents a Approval event raised by the Xtoken contract.
type XtokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Xtoken *XtokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*XtokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Xtoken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &XtokenApprovalIterator{contract: _Xtoken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Xtoken *XtokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *XtokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Xtoken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XtokenApproval)
				if err := _Xtoken.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Xtoken *XtokenFilterer) ParseApproval(log types.Log) (*XtokenApproval, error) {
	event := new(XtokenApproval)
	if err := _Xtoken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XtokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Xtoken contract.
type XtokenOwnershipTransferredIterator struct {
	Event *XtokenOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XtokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XtokenOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XtokenOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XtokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XtokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XtokenOwnershipTransferred represents a OwnershipTransferred event raised by the Xtoken contract.
type XtokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Xtoken *XtokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*XtokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Xtoken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &XtokenOwnershipTransferredIterator{contract: _Xtoken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Xtoken *XtokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *XtokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Xtoken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XtokenOwnershipTransferred)
				if err := _Xtoken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Xtoken *XtokenFilterer) ParseOwnershipTransferred(log types.Log) (*XtokenOwnershipTransferred, error) {
	event := new(XtokenOwnershipTransferred)
	if err := _Xtoken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XtokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Xtoken contract.
type XtokenTransferIterator struct {
	Event *XtokenTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XtokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XtokenTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XtokenTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XtokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XtokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XtokenTransfer represents a Transfer event raised by the Xtoken contract.
type XtokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Xtoken *XtokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*XtokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Xtoken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &XtokenTransferIterator{contract: _Xtoken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Xtoken *XtokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *XtokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Xtoken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XtokenTransfer)
				if err := _Xtoken.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Xtoken *XtokenFilterer) ParseTransfer(log types.Log) (*XtokenTransfer, error) {
	event := new(XtokenTransfer)
	if err := _Xtoken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
