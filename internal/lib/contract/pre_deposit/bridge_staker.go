// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package pre_deposit

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

// BridgeStakerMetaData contains all meta data concerning the BridgeStaker contract.
var BridgeStakerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_unlockTime\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"ContractAuthorized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"wastr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vastr\",\"type\":\"address\"}],\"name\":\"TokensSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newTime\",\"type\":\"uint256\"}],\"name\":\"UnlockTimeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"wastrAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"vastrAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"authorizedContracts\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"wastrAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"vastrAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalWastr\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVastr\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isUnlocked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"authorized\",\"type\":\"bool\"}],\"name\":\"setAuthorizedContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_wastr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vastr\",\"type\":\"address\"}],\"name\":\"setTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newTime\",\"type\":\"uint256\"}],\"name\":\"setUnlockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVastrStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalWastrStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vastr\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"vastrBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wastr\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"wastrBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BridgeStakerABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeStakerMetaData.ABI instead.
var BridgeStakerABI = BridgeStakerMetaData.ABI

// BridgeStaker is an auto generated Go binding around an Ethereum contract.
type BridgeStaker struct {
	BridgeStakerCaller     // Read-only binding to the contract
	BridgeStakerTransactor // Write-only binding to the contract
	BridgeStakerFilterer   // Log filterer for contract events
}

// BridgeStakerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeStakerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeStakerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeStakerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeStakerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeStakerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeStakerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeStakerSession struct {
	Contract     *BridgeStaker     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeStakerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeStakerCallerSession struct {
	Contract *BridgeStakerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BridgeStakerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeStakerTransactorSession struct {
	Contract     *BridgeStakerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BridgeStakerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeStakerRaw struct {
	Contract *BridgeStaker // Generic contract binding to access the raw methods on
}

// BridgeStakerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeStakerCallerRaw struct {
	Contract *BridgeStakerCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeStakerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeStakerTransactorRaw struct {
	Contract *BridgeStakerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridgeStaker creates a new instance of BridgeStaker, bound to a specific deployed contract.
func NewBridgeStaker(address common.Address, backend bind.ContractBackend) (*BridgeStaker, error) {
	contract, err := bindBridgeStaker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BridgeStaker{BridgeStakerCaller: BridgeStakerCaller{contract: contract}, BridgeStakerTransactor: BridgeStakerTransactor{contract: contract}, BridgeStakerFilterer: BridgeStakerFilterer{contract: contract}}, nil
}

// NewBridgeStakerCaller creates a new read-only instance of BridgeStaker, bound to a specific deployed contract.
func NewBridgeStakerCaller(address common.Address, caller bind.ContractCaller) (*BridgeStakerCaller, error) {
	contract, err := bindBridgeStaker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerCaller{contract: contract}, nil
}

// NewBridgeStakerTransactor creates a new write-only instance of BridgeStaker, bound to a specific deployed contract.
func NewBridgeStakerTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeStakerTransactor, error) {
	contract, err := bindBridgeStaker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerTransactor{contract: contract}, nil
}

// NewBridgeStakerFilterer creates a new log filterer instance of BridgeStaker, bound to a specific deployed contract.
func NewBridgeStakerFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeStakerFilterer, error) {
	contract, err := bindBridgeStaker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerFilterer{contract: contract}, nil
}

// bindBridgeStaker binds a generic wrapper to an already deployed contract.
func bindBridgeStaker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BridgeStakerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeStaker *BridgeStakerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeStaker.Contract.BridgeStakerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeStaker *BridgeStakerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeStaker.Contract.BridgeStakerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeStaker *BridgeStakerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeStaker.Contract.BridgeStakerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BridgeStaker *BridgeStakerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BridgeStaker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BridgeStaker *BridgeStakerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeStaker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BridgeStaker *BridgeStakerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BridgeStaker.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedContracts is a free data retrieval call binding the contract method 0xd5b9221b.
//
// Solidity: function authorizedContracts(address ) view returns(bool)
func (_BridgeStaker *BridgeStakerCaller) AuthorizedContracts(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "authorizedContracts", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizedContracts is a free data retrieval call binding the contract method 0xd5b9221b.
//
// Solidity: function authorizedContracts(address ) view returns(bool)
func (_BridgeStaker *BridgeStakerSession) AuthorizedContracts(arg0 common.Address) (bool, error) {
	return _BridgeStaker.Contract.AuthorizedContracts(&_BridgeStaker.CallOpts, arg0)
}

// AuthorizedContracts is a free data retrieval call binding the contract method 0xd5b9221b.
//
// Solidity: function authorizedContracts(address ) view returns(bool)
func (_BridgeStaker *BridgeStakerCallerSession) AuthorizedContracts(arg0 common.Address) (bool, error) {
	return _BridgeStaker.Contract.AuthorizedContracts(&_BridgeStaker.CallOpts, arg0)
}

// GetBalances is a free data retrieval call binding the contract method 0xc84aae17.
//
// Solidity: function getBalances(address user) view returns(uint256 wastrAmount, uint256 vastrAmount)
func (_BridgeStaker *BridgeStakerCaller) GetBalances(opts *bind.CallOpts, user common.Address) (struct {
	WastrAmount *big.Int
	VastrAmount *big.Int
}, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "getBalances", user)

	outstruct := new(struct {
		WastrAmount *big.Int
		VastrAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.WastrAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.VastrAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetBalances is a free data retrieval call binding the contract method 0xc84aae17.
//
// Solidity: function getBalances(address user) view returns(uint256 wastrAmount, uint256 vastrAmount)
func (_BridgeStaker *BridgeStakerSession) GetBalances(user common.Address) (struct {
	WastrAmount *big.Int
	VastrAmount *big.Int
}, error) {
	return _BridgeStaker.Contract.GetBalances(&_BridgeStaker.CallOpts, user)
}

// GetBalances is a free data retrieval call binding the contract method 0xc84aae17.
//
// Solidity: function getBalances(address user) view returns(uint256 wastrAmount, uint256 vastrAmount)
func (_BridgeStaker *BridgeStakerCallerSession) GetBalances(user common.Address) (struct {
	WastrAmount *big.Int
	VastrAmount *big.Int
}, error) {
	return _BridgeStaker.Contract.GetBalances(&_BridgeStaker.CallOpts, user)
}

// GetTotalStaked is a free data retrieval call binding the contract method 0x0917e776.
//
// Solidity: function getTotalStaked() view returns(uint256 totalWastr, uint256 totalVastr)
func (_BridgeStaker *BridgeStakerCaller) GetTotalStaked(opts *bind.CallOpts) (struct {
	TotalWastr *big.Int
	TotalVastr *big.Int
}, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "getTotalStaked")

	outstruct := new(struct {
		TotalWastr *big.Int
		TotalVastr *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalWastr = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalVastr = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTotalStaked is a free data retrieval call binding the contract method 0x0917e776.
//
// Solidity: function getTotalStaked() view returns(uint256 totalWastr, uint256 totalVastr)
func (_BridgeStaker *BridgeStakerSession) GetTotalStaked() (struct {
	TotalWastr *big.Int
	TotalVastr *big.Int
}, error) {
	return _BridgeStaker.Contract.GetTotalStaked(&_BridgeStaker.CallOpts)
}

// GetTotalStaked is a free data retrieval call binding the contract method 0x0917e776.
//
// Solidity: function getTotalStaked() view returns(uint256 totalWastr, uint256 totalVastr)
func (_BridgeStaker *BridgeStakerCallerSession) GetTotalStaked() (struct {
	TotalWastr *big.Int
	TotalVastr *big.Int
}, error) {
	return _BridgeStaker.Contract.GetTotalStaked(&_BridgeStaker.CallOpts)
}

// IsUnlocked is a free data retrieval call binding the contract method 0x8380edb7.
//
// Solidity: function isUnlocked() view returns(bool)
func (_BridgeStaker *BridgeStakerCaller) IsUnlocked(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "isUnlocked")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsUnlocked is a free data retrieval call binding the contract method 0x8380edb7.
//
// Solidity: function isUnlocked() view returns(bool)
func (_BridgeStaker *BridgeStakerSession) IsUnlocked() (bool, error) {
	return _BridgeStaker.Contract.IsUnlocked(&_BridgeStaker.CallOpts)
}

// IsUnlocked is a free data retrieval call binding the contract method 0x8380edb7.
//
// Solidity: function isUnlocked() view returns(bool)
func (_BridgeStaker *BridgeStakerCallerSession) IsUnlocked() (bool, error) {
	return _BridgeStaker.Contract.IsUnlocked(&_BridgeStaker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BridgeStaker *BridgeStakerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BridgeStaker *BridgeStakerSession) Owner() (common.Address, error) {
	return _BridgeStaker.Contract.Owner(&_BridgeStaker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BridgeStaker *BridgeStakerCallerSession) Owner() (common.Address, error) {
	return _BridgeStaker.Contract.Owner(&_BridgeStaker.CallOpts)
}

// TotalVastrStaked is a free data retrieval call binding the contract method 0x6ff08be5.
//
// Solidity: function totalVastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerCaller) TotalVastrStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "totalVastrStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVastrStaked is a free data retrieval call binding the contract method 0x6ff08be5.
//
// Solidity: function totalVastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerSession) TotalVastrStaked() (*big.Int, error) {
	return _BridgeStaker.Contract.TotalVastrStaked(&_BridgeStaker.CallOpts)
}

// TotalVastrStaked is a free data retrieval call binding the contract method 0x6ff08be5.
//
// Solidity: function totalVastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerCallerSession) TotalVastrStaked() (*big.Int, error) {
	return _BridgeStaker.Contract.TotalVastrStaked(&_BridgeStaker.CallOpts)
}

// TotalWastrStaked is a free data retrieval call binding the contract method 0x6d0f89a8.
//
// Solidity: function totalWastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerCaller) TotalWastrStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "totalWastrStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalWastrStaked is a free data retrieval call binding the contract method 0x6d0f89a8.
//
// Solidity: function totalWastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerSession) TotalWastrStaked() (*big.Int, error) {
	return _BridgeStaker.Contract.TotalWastrStaked(&_BridgeStaker.CallOpts)
}

// TotalWastrStaked is a free data retrieval call binding the contract method 0x6d0f89a8.
//
// Solidity: function totalWastrStaked() view returns(uint256)
func (_BridgeStaker *BridgeStakerCallerSession) TotalWastrStaked() (*big.Int, error) {
	return _BridgeStaker.Contract.TotalWastrStaked(&_BridgeStaker.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_BridgeStaker *BridgeStakerCaller) UnlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "unlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_BridgeStaker *BridgeStakerSession) UnlockTime() (*big.Int, error) {
	return _BridgeStaker.Contract.UnlockTime(&_BridgeStaker.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_BridgeStaker *BridgeStakerCallerSession) UnlockTime() (*big.Int, error) {
	return _BridgeStaker.Contract.UnlockTime(&_BridgeStaker.CallOpts)
}

// Vastr is a free data retrieval call binding the contract method 0x85720141.
//
// Solidity: function vastr() view returns(address)
func (_BridgeStaker *BridgeStakerCaller) Vastr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "vastr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vastr is a free data retrieval call binding the contract method 0x85720141.
//
// Solidity: function vastr() view returns(address)
func (_BridgeStaker *BridgeStakerSession) Vastr() (common.Address, error) {
	return _BridgeStaker.Contract.Vastr(&_BridgeStaker.CallOpts)
}

// Vastr is a free data retrieval call binding the contract method 0x85720141.
//
// Solidity: function vastr() view returns(address)
func (_BridgeStaker *BridgeStakerCallerSession) Vastr() (common.Address, error) {
	return _BridgeStaker.Contract.Vastr(&_BridgeStaker.CallOpts)
}

// VastrBalance is a free data retrieval call binding the contract method 0xb284106a.
//
// Solidity: function vastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerCaller) VastrBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "vastrBalance", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VastrBalance is a free data retrieval call binding the contract method 0xb284106a.
//
// Solidity: function vastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerSession) VastrBalance(arg0 common.Address) (*big.Int, error) {
	return _BridgeStaker.Contract.VastrBalance(&_BridgeStaker.CallOpts, arg0)
}

// VastrBalance is a free data retrieval call binding the contract method 0xb284106a.
//
// Solidity: function vastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerCallerSession) VastrBalance(arg0 common.Address) (*big.Int, error) {
	return _BridgeStaker.Contract.VastrBalance(&_BridgeStaker.CallOpts, arg0)
}

// Wastr is a free data retrieval call binding the contract method 0x0ab2a9a2.
//
// Solidity: function wastr() view returns(address)
func (_BridgeStaker *BridgeStakerCaller) Wastr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "wastr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Wastr is a free data retrieval call binding the contract method 0x0ab2a9a2.
//
// Solidity: function wastr() view returns(address)
func (_BridgeStaker *BridgeStakerSession) Wastr() (common.Address, error) {
	return _BridgeStaker.Contract.Wastr(&_BridgeStaker.CallOpts)
}

// Wastr is a free data retrieval call binding the contract method 0x0ab2a9a2.
//
// Solidity: function wastr() view returns(address)
func (_BridgeStaker *BridgeStakerCallerSession) Wastr() (common.Address, error) {
	return _BridgeStaker.Contract.Wastr(&_BridgeStaker.CallOpts)
}

// WastrBalance is a free data retrieval call binding the contract method 0x25f3b5bf.
//
// Solidity: function wastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerCaller) WastrBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BridgeStaker.contract.Call(opts, &out, "wastrBalance", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WastrBalance is a free data retrieval call binding the contract method 0x25f3b5bf.
//
// Solidity: function wastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerSession) WastrBalance(arg0 common.Address) (*big.Int, error) {
	return _BridgeStaker.Contract.WastrBalance(&_BridgeStaker.CallOpts, arg0)
}

// WastrBalance is a free data retrieval call binding the contract method 0x25f3b5bf.
//
// Solidity: function wastrBalance(address ) view returns(uint256)
func (_BridgeStaker *BridgeStakerCallerSession) WastrBalance(arg0 common.Address) (*big.Int, error) {
	return _BridgeStaker.Contract.WastrBalance(&_BridgeStaker.CallOpts, arg0)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeStaker *BridgeStakerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeStaker *BridgeStakerSession) RenounceOwnership() (*types.Transaction, error) {
	return _BridgeStaker.Contract.RenounceOwnership(&_BridgeStaker.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BridgeStaker *BridgeStakerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BridgeStaker.Contract.RenounceOwnership(&_BridgeStaker.TransactOpts)
}

// SetAuthorizedContract is a paid mutator transaction binding the contract method 0xf2fa7392.
//
// Solidity: function setAuthorizedContract(address contractAddress, bool authorized) returns()
func (_BridgeStaker *BridgeStakerTransactor) SetAuthorizedContract(opts *bind.TransactOpts, contractAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "setAuthorizedContract", contractAddress, authorized)
}

// SetAuthorizedContract is a paid mutator transaction binding the contract method 0xf2fa7392.
//
// Solidity: function setAuthorizedContract(address contractAddress, bool authorized) returns()
func (_BridgeStaker *BridgeStakerSession) SetAuthorizedContract(contractAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetAuthorizedContract(&_BridgeStaker.TransactOpts, contractAddress, authorized)
}

// SetAuthorizedContract is a paid mutator transaction binding the contract method 0xf2fa7392.
//
// Solidity: function setAuthorizedContract(address contractAddress, bool authorized) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) SetAuthorizedContract(contractAddress common.Address, authorized bool) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetAuthorizedContract(&_BridgeStaker.TransactOpts, contractAddress, authorized)
}

// SetTokens is a paid mutator transaction binding the contract method 0xcbc7854e.
//
// Solidity: function setTokens(address _wastr, address _vastr) returns()
func (_BridgeStaker *BridgeStakerTransactor) SetTokens(opts *bind.TransactOpts, _wastr common.Address, _vastr common.Address) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "setTokens", _wastr, _vastr)
}

// SetTokens is a paid mutator transaction binding the contract method 0xcbc7854e.
//
// Solidity: function setTokens(address _wastr, address _vastr) returns()
func (_BridgeStaker *BridgeStakerSession) SetTokens(_wastr common.Address, _vastr common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetTokens(&_BridgeStaker.TransactOpts, _wastr, _vastr)
}

// SetTokens is a paid mutator transaction binding the contract method 0xcbc7854e.
//
// Solidity: function setTokens(address _wastr, address _vastr) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) SetTokens(_wastr common.Address, _vastr common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetTokens(&_BridgeStaker.TransactOpts, _wastr, _vastr)
}

// SetUnlockTime is a paid mutator transaction binding the contract method 0xdace4557.
//
// Solidity: function setUnlockTime(uint256 newTime) returns()
func (_BridgeStaker *BridgeStakerTransactor) SetUnlockTime(opts *bind.TransactOpts, newTime *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "setUnlockTime", newTime)
}

// SetUnlockTime is a paid mutator transaction binding the contract method 0xdace4557.
//
// Solidity: function setUnlockTime(uint256 newTime) returns()
func (_BridgeStaker *BridgeStakerSession) SetUnlockTime(newTime *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetUnlockTime(&_BridgeStaker.TransactOpts, newTime)
}

// SetUnlockTime is a paid mutator transaction binding the contract method 0xdace4557.
//
// Solidity: function setUnlockTime(uint256 newTime) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) SetUnlockTime(newTime *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.Contract.SetUnlockTime(&_BridgeStaker.TransactOpts, newTime)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(address token, address sender, uint256 amount) returns()
func (_BridgeStaker *BridgeStakerTransactor) Stake(opts *bind.TransactOpts, token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "stake", token, sender, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(address token, address sender, uint256 amount) returns()
func (_BridgeStaker *BridgeStakerSession) Stake(token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.Contract.Stake(&_BridgeStaker.TransactOpts, token, sender, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(address token, address sender, uint256 amount) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) Stake(token common.Address, sender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BridgeStaker.Contract.Stake(&_BridgeStaker.TransactOpts, token, sender, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BridgeStaker *BridgeStakerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BridgeStaker *BridgeStakerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.TransferOwnership(&_BridgeStaker.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.TransferOwnership(&_BridgeStaker.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token) returns()
func (_BridgeStaker *BridgeStakerTransactor) Withdraw(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _BridgeStaker.contract.Transact(opts, "withdraw", token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token) returns()
func (_BridgeStaker *BridgeStakerSession) Withdraw(token common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.Withdraw(&_BridgeStaker.TransactOpts, token)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address token) returns()
func (_BridgeStaker *BridgeStakerTransactorSession) Withdraw(token common.Address) (*types.Transaction, error) {
	return _BridgeStaker.Contract.Withdraw(&_BridgeStaker.TransactOpts, token)
}

// BridgeStakerContractAuthorizedIterator is returned from FilterContractAuthorized and is used to iterate over the raw logs and unpacked data for ContractAuthorized events raised by the BridgeStaker contract.
type BridgeStakerContractAuthorizedIterator struct {
	Event *BridgeStakerContractAuthorized // Event containing the contract specifics and raw log

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
func (it *BridgeStakerContractAuthorizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerContractAuthorized)
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
		it.Event = new(BridgeStakerContractAuthorized)
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
func (it *BridgeStakerContractAuthorizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerContractAuthorizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerContractAuthorized represents a ContractAuthorized event raised by the BridgeStaker contract.
type BridgeStakerContractAuthorized struct {
	ContractAddress common.Address
	Authorized      bool
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterContractAuthorized is a free log retrieval operation binding the contract event 0x3213ab8ec2ed122402f038736b41878b69ae9f61f49a7edd0666bce891d29fc9.
//
// Solidity: event ContractAuthorized(address indexed contractAddress, bool authorized)
func (_BridgeStaker *BridgeStakerFilterer) FilterContractAuthorized(opts *bind.FilterOpts, contractAddress []common.Address) (*BridgeStakerContractAuthorizedIterator, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "ContractAuthorized", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerContractAuthorizedIterator{contract: _BridgeStaker.contract, event: "ContractAuthorized", logs: logs, sub: sub}, nil
}

// WatchContractAuthorized is a free log subscription operation binding the contract event 0x3213ab8ec2ed122402f038736b41878b69ae9f61f49a7edd0666bce891d29fc9.
//
// Solidity: event ContractAuthorized(address indexed contractAddress, bool authorized)
func (_BridgeStaker *BridgeStakerFilterer) WatchContractAuthorized(opts *bind.WatchOpts, sink chan<- *BridgeStakerContractAuthorized, contractAddress []common.Address) (event.Subscription, error) {

	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "ContractAuthorized", contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerContractAuthorized)
				if err := _BridgeStaker.contract.UnpackLog(event, "ContractAuthorized", log); err != nil {
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

// ParseContractAuthorized is a log parse operation binding the contract event 0x3213ab8ec2ed122402f038736b41878b69ae9f61f49a7edd0666bce891d29fc9.
//
// Solidity: event ContractAuthorized(address indexed contractAddress, bool authorized)
func (_BridgeStaker *BridgeStakerFilterer) ParseContractAuthorized(log types.Log) (*BridgeStakerContractAuthorized, error) {
	event := new(BridgeStakerContractAuthorized)
	if err := _BridgeStaker.contract.UnpackLog(event, "ContractAuthorized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeStakerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BridgeStaker contract.
type BridgeStakerOwnershipTransferredIterator struct {
	Event *BridgeStakerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BridgeStakerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerOwnershipTransferred)
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
		it.Event = new(BridgeStakerOwnershipTransferred)
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
func (it *BridgeStakerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerOwnershipTransferred represents a OwnershipTransferred event raised by the BridgeStaker contract.
type BridgeStakerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BridgeStaker *BridgeStakerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BridgeStakerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerOwnershipTransferredIterator{contract: _BridgeStaker.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BridgeStaker *BridgeStakerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BridgeStakerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerOwnershipTransferred)
				if err := _BridgeStaker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BridgeStaker *BridgeStakerFilterer) ParseOwnershipTransferred(log types.Log) (*BridgeStakerOwnershipTransferred, error) {
	event := new(BridgeStakerOwnershipTransferred)
	if err := _BridgeStaker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeStakerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the BridgeStaker contract.
type BridgeStakerStakedIterator struct {
	Event *BridgeStakerStaked // Event containing the contract specifics and raw log

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
func (it *BridgeStakerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerStaked)
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
		it.Event = new(BridgeStakerStaked)
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
func (it *BridgeStakerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerStaked represents a Staked event raised by the BridgeStaker contract.
type BridgeStakerStaked struct {
	User      common.Address
	Token     common.Address
	Amount    *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x6c86f3fd5118b3aa8bb4f389a617046de0a3d3d477de1a1673d227f802f616dc.
//
// Solidity: event Staked(address indexed user, address indexed token, uint256 amount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*BridgeStakerStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "Staked", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerStakedIterator{contract: _BridgeStaker.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x6c86f3fd5118b3aa8bb4f389a617046de0a3d3d477de1a1673d227f802f616dc.
//
// Solidity: event Staked(address indexed user, address indexed token, uint256 amount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *BridgeStakerStaked, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "Staked", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerStaked)
				if err := _BridgeStaker.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x6c86f3fd5118b3aa8bb4f389a617046de0a3d3d477de1a1673d227f802f616dc.
//
// Solidity: event Staked(address indexed user, address indexed token, uint256 amount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) ParseStaked(log types.Log) (*BridgeStakerStaked, error) {
	event := new(BridgeStakerStaked)
	if err := _BridgeStaker.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeStakerTokensSetIterator is returned from FilterTokensSet and is used to iterate over the raw logs and unpacked data for TokensSet events raised by the BridgeStaker contract.
type BridgeStakerTokensSetIterator struct {
	Event *BridgeStakerTokensSet // Event containing the contract specifics and raw log

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
func (it *BridgeStakerTokensSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerTokensSet)
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
		it.Event = new(BridgeStakerTokensSet)
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
func (it *BridgeStakerTokensSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerTokensSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerTokensSet represents a TokensSet event raised by the BridgeStaker contract.
type BridgeStakerTokensSet struct {
	Wastr common.Address
	Vastr common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTokensSet is a free log retrieval operation binding the contract event 0x51278cfe696b5a9595ada5bd0281911ab71445bc087ca9393087ea5c7b6edf30.
//
// Solidity: event TokensSet(address wastr, address vastr)
func (_BridgeStaker *BridgeStakerFilterer) FilterTokensSet(opts *bind.FilterOpts) (*BridgeStakerTokensSetIterator, error) {

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "TokensSet")
	if err != nil {
		return nil, err
	}
	return &BridgeStakerTokensSetIterator{contract: _BridgeStaker.contract, event: "TokensSet", logs: logs, sub: sub}, nil
}

// WatchTokensSet is a free log subscription operation binding the contract event 0x51278cfe696b5a9595ada5bd0281911ab71445bc087ca9393087ea5c7b6edf30.
//
// Solidity: event TokensSet(address wastr, address vastr)
func (_BridgeStaker *BridgeStakerFilterer) WatchTokensSet(opts *bind.WatchOpts, sink chan<- *BridgeStakerTokensSet) (event.Subscription, error) {

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "TokensSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerTokensSet)
				if err := _BridgeStaker.contract.UnpackLog(event, "TokensSet", log); err != nil {
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

// ParseTokensSet is a log parse operation binding the contract event 0x51278cfe696b5a9595ada5bd0281911ab71445bc087ca9393087ea5c7b6edf30.
//
// Solidity: event TokensSet(address wastr, address vastr)
func (_BridgeStaker *BridgeStakerFilterer) ParseTokensSet(log types.Log) (*BridgeStakerTokensSet, error) {
	event := new(BridgeStakerTokensSet)
	if err := _BridgeStaker.contract.UnpackLog(event, "TokensSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeStakerUnlockTimeUpdatedIterator is returned from FilterUnlockTimeUpdated and is used to iterate over the raw logs and unpacked data for UnlockTimeUpdated events raised by the BridgeStaker contract.
type BridgeStakerUnlockTimeUpdatedIterator struct {
	Event *BridgeStakerUnlockTimeUpdated // Event containing the contract specifics and raw log

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
func (it *BridgeStakerUnlockTimeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerUnlockTimeUpdated)
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
		it.Event = new(BridgeStakerUnlockTimeUpdated)
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
func (it *BridgeStakerUnlockTimeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerUnlockTimeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerUnlockTimeUpdated represents a UnlockTimeUpdated event raised by the BridgeStaker contract.
type BridgeStakerUnlockTimeUpdated struct {
	NewTime *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnlockTimeUpdated is a free log retrieval operation binding the contract event 0x6669a6c779a56629103733bcba993d66b30b061559b41f3c64bd2661cd9ea3c1.
//
// Solidity: event UnlockTimeUpdated(uint256 newTime)
func (_BridgeStaker *BridgeStakerFilterer) FilterUnlockTimeUpdated(opts *bind.FilterOpts) (*BridgeStakerUnlockTimeUpdatedIterator, error) {

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "UnlockTimeUpdated")
	if err != nil {
		return nil, err
	}
	return &BridgeStakerUnlockTimeUpdatedIterator{contract: _BridgeStaker.contract, event: "UnlockTimeUpdated", logs: logs, sub: sub}, nil
}

// WatchUnlockTimeUpdated is a free log subscription operation binding the contract event 0x6669a6c779a56629103733bcba993d66b30b061559b41f3c64bd2661cd9ea3c1.
//
// Solidity: event UnlockTimeUpdated(uint256 newTime)
func (_BridgeStaker *BridgeStakerFilterer) WatchUnlockTimeUpdated(opts *bind.WatchOpts, sink chan<- *BridgeStakerUnlockTimeUpdated) (event.Subscription, error) {

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "UnlockTimeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerUnlockTimeUpdated)
				if err := _BridgeStaker.contract.UnpackLog(event, "UnlockTimeUpdated", log); err != nil {
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

// ParseUnlockTimeUpdated is a log parse operation binding the contract event 0x6669a6c779a56629103733bcba993d66b30b061559b41f3c64bd2661cd9ea3c1.
//
// Solidity: event UnlockTimeUpdated(uint256 newTime)
func (_BridgeStaker *BridgeStakerFilterer) ParseUnlockTimeUpdated(log types.Log) (*BridgeStakerUnlockTimeUpdated, error) {
	event := new(BridgeStakerUnlockTimeUpdated)
	if err := _BridgeStaker.contract.UnpackLog(event, "UnlockTimeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeStakerWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the BridgeStaker contract.
type BridgeStakerWithdrawnIterator struct {
	Event *BridgeStakerWithdrawn // Event containing the contract specifics and raw log

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
func (it *BridgeStakerWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeStakerWithdrawn)
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
		it.Event = new(BridgeStakerWithdrawn)
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
func (it *BridgeStakerWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeStakerWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeStakerWithdrawn represents a Withdrawn event raised by the BridgeStaker contract.
type BridgeStakerWithdrawn struct {
	User        common.Address
	WastrAmount *big.Int
	VastrAmount *big.Int
	Timestamp   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x75e161b3e824b114fc1a33274bd7091918dd4e639cede50b78b15a4eea956a21.
//
// Solidity: event Withdrawn(address indexed user, uint256 wastrAmount, uint256 vastrAmount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) FilterWithdrawn(opts *bind.FilterOpts, user []common.Address) (*BridgeStakerWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BridgeStaker.contract.FilterLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &BridgeStakerWithdrawnIterator{contract: _BridgeStaker.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x75e161b3e824b114fc1a33274bd7091918dd4e639cede50b78b15a4eea956a21.
//
// Solidity: event Withdrawn(address indexed user, uint256 wastrAmount, uint256 vastrAmount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *BridgeStakerWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _BridgeStaker.contract.WatchLogs(opts, "Withdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeStakerWithdrawn)
				if err := _BridgeStaker.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x75e161b3e824b114fc1a33274bd7091918dd4e639cede50b78b15a4eea956a21.
//
// Solidity: event Withdrawn(address indexed user, uint256 wastrAmount, uint256 vastrAmount, uint256 timestamp)
func (_BridgeStaker *BridgeStakerFilterer) ParseWithdrawn(log types.Log) (*BridgeStakerWithdrawn, error) {
	event := new(BridgeStakerWithdrawn)
	if err := _BridgeStaker.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
