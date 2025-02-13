// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package xpool

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

// StakingPoolMetaData contains all meta data concerning the StakingPool contract.
var StakingPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rewardPerBlock\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_withdrawed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_claimed\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"name\":\"Withdrawed\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"claim\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"claimed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_claimed\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"destroy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"earned\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_earned\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isDestroied\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_rewardBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardPerBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rewardPerBlock\",\"type\":\"uint256\"}],\"name\":\"setRewardPerBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"staked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_staked\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakedNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_stakedNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakingToken\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalClaimed\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_withdrawed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_claimed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_height\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// StakingPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingPoolMetaData.ABI instead.
var StakingPoolABI = StakingPoolMetaData.ABI

// StakingPool is an auto generated Go binding around an Ethereum contract.
type StakingPool struct {
	StakingPoolCaller     // Read-only binding to the contract
	StakingPoolTransactor // Write-only binding to the contract
	StakingPoolFilterer   // Log filterer for contract events
}

// StakingPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingPoolSession struct {
	Contract     *StakingPool      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingPoolCallerSession struct {
	Contract *StakingPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakingPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingPoolTransactorSession struct {
	Contract     *StakingPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakingPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingPoolRaw struct {
	Contract *StakingPool // Generic contract binding to access the raw methods on
}

// StakingPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingPoolCallerRaw struct {
	Contract *StakingPoolCaller // Generic read-only contract binding to access the raw methods on
}

// StakingPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingPoolTransactorRaw struct {
	Contract *StakingPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingPool creates a new instance of StakingPool, bound to a specific deployed contract.
func NewStakingPool(address common.Address, backend bind.ContractBackend) (*StakingPool, error) {
	contract, err := bindStakingPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingPool{StakingPoolCaller: StakingPoolCaller{contract: contract}, StakingPoolTransactor: StakingPoolTransactor{contract: contract}, StakingPoolFilterer: StakingPoolFilterer{contract: contract}}, nil
}

// NewStakingPoolCaller creates a new read-only instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolCaller(address common.Address, caller bind.ContractCaller) (*StakingPoolCaller, error) {
	contract, err := bindStakingPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolCaller{contract: contract}, nil
}

// NewStakingPoolTransactor creates a new write-only instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingPoolTransactor, error) {
	contract, err := bindStakingPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolTransactor{contract: contract}, nil
}

// NewStakingPoolFilterer creates a new log filterer instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingPoolFilterer, error) {
	contract, err := bindStakingPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingPoolFilterer{contract: contract}, nil
}

// bindStakingPool binds a generic wrapper to an already deployed contract.
func bindStakingPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingPoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPool *StakingPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPool.Contract.StakingPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPool *StakingPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.Contract.StakingPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPool *StakingPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPool.Contract.StakingPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPool *StakingPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPool *StakingPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPool *StakingPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPool.Contract.contract.Transact(opts, method, params...)
}

// Claimed is a free data retrieval call binding the contract method 0xc884ef83.
//
// Solidity: function claimed(address _user) view returns(uint256 _claimed)
func (_StakingPool *StakingPoolCaller) Claimed(opts *bind.CallOpts, _user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "claimed", _user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Claimed is a free data retrieval call binding the contract method 0xc884ef83.
//
// Solidity: function claimed(address _user) view returns(uint256 _claimed)
func (_StakingPool *StakingPoolSession) Claimed(_user common.Address) (*big.Int, error) {
	return _StakingPool.Contract.Claimed(&_StakingPool.CallOpts, _user)
}

// Claimed is a free data retrieval call binding the contract method 0xc884ef83.
//
// Solidity: function claimed(address _user) view returns(uint256 _claimed)
func (_StakingPool *StakingPoolCallerSession) Claimed(_user common.Address) (*big.Int, error) {
	return _StakingPool.Contract.Claimed(&_StakingPool.CallOpts, _user)
}

// Earned is a free data retrieval call binding the contract method 0x008cc262.
//
// Solidity: function earned(address _user) view returns(uint256 _earned, uint256 _height)
func (_StakingPool *StakingPoolCaller) Earned(opts *bind.CallOpts, _user common.Address) (struct {
	Earned *big.Int
	Height *big.Int
}, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "earned", _user)

	outstruct := new(struct {
		Earned *big.Int
		Height *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Earned = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Height = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Earned is a free data retrieval call binding the contract method 0x008cc262.
//
// Solidity: function earned(address _user) view returns(uint256 _earned, uint256 _height)
func (_StakingPool *StakingPoolSession) Earned(_user common.Address) (struct {
	Earned *big.Int
	Height *big.Int
}, error) {
	return _StakingPool.Contract.Earned(&_StakingPool.CallOpts, _user)
}

// Earned is a free data retrieval call binding the contract method 0x008cc262.
//
// Solidity: function earned(address _user) view returns(uint256 _earned, uint256 _height)
func (_StakingPool *StakingPoolCallerSession) Earned(_user common.Address) (struct {
	Earned *big.Int
	Height *big.Int
}, error) {
	return _StakingPool.Contract.Earned(&_StakingPool.CallOpts, _user)
}

// IsDestroied is a free data retrieval call binding the contract method 0x7a884b3e.
//
// Solidity: function isDestroied() view returns(bool)
func (_StakingPool *StakingPoolCaller) IsDestroied(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "isDestroied")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDestroied is a free data retrieval call binding the contract method 0x7a884b3e.
//
// Solidity: function isDestroied() view returns(bool)
func (_StakingPool *StakingPoolSession) IsDestroied() (bool, error) {
	return _StakingPool.Contract.IsDestroied(&_StakingPool.CallOpts)
}

// IsDestroied is a free data retrieval call binding the contract method 0x7a884b3e.
//
// Solidity: function isDestroied() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) IsDestroied() (bool, error) {
	return _StakingPool.Contract.IsDestroied(&_StakingPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingPool *StakingPoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingPool *StakingPoolSession) Owner() (common.Address, error) {
	return _StakingPool.Contract.Owner(&_StakingPool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingPool *StakingPoolCallerSession) Owner() (common.Address, error) {
	return _StakingPool.Contract.Owner(&_StakingPool.CallOpts)
}

// RewardBalance is a free data retrieval call binding the contract method 0xaa5c3ab4.
//
// Solidity: function rewardBalance() view returns(uint256 _rewardBalance)
func (_StakingPool *StakingPoolCaller) RewardBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "rewardBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardBalance is a free data retrieval call binding the contract method 0xaa5c3ab4.
//
// Solidity: function rewardBalance() view returns(uint256 _rewardBalance)
func (_StakingPool *StakingPoolSession) RewardBalance() (*big.Int, error) {
	return _StakingPool.Contract.RewardBalance(&_StakingPool.CallOpts)
}

// RewardBalance is a free data retrieval call binding the contract method 0xaa5c3ab4.
//
// Solidity: function rewardBalance() view returns(uint256 _rewardBalance)
func (_StakingPool *StakingPoolCallerSession) RewardBalance() (*big.Int, error) {
	return _StakingPool.Contract.RewardBalance(&_StakingPool.CallOpts)
}

// RewardPerBlock is a free data retrieval call binding the contract method 0x8ae39cac.
//
// Solidity: function rewardPerBlock() view returns(uint256)
func (_StakingPool *StakingPoolCaller) RewardPerBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "rewardPerBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardPerBlock is a free data retrieval call binding the contract method 0x8ae39cac.
//
// Solidity: function rewardPerBlock() view returns(uint256)
func (_StakingPool *StakingPoolSession) RewardPerBlock() (*big.Int, error) {
	return _StakingPool.Contract.RewardPerBlock(&_StakingPool.CallOpts)
}

// RewardPerBlock is a free data retrieval call binding the contract method 0x8ae39cac.
//
// Solidity: function rewardPerBlock() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) RewardPerBlock() (*big.Int, error) {
	return _StakingPool.Contract.RewardPerBlock(&_StakingPool.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_StakingPool *StakingPoolCaller) RewardToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "rewardToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_StakingPool *StakingPoolSession) RewardToken() (common.Address, error) {
	return _StakingPool.Contract.RewardToken(&_StakingPool.CallOpts)
}

// RewardToken is a free data retrieval call binding the contract method 0xf7c618c1.
//
// Solidity: function rewardToken() view returns(address)
func (_StakingPool *StakingPoolCallerSession) RewardToken() (common.Address, error) {
	return _StakingPool.Contract.RewardToken(&_StakingPool.CallOpts)
}

// Staked is a free data retrieval call binding the contract method 0x98807d84.
//
// Solidity: function staked(address _user) view returns(uint256 _staked)
func (_StakingPool *StakingPoolCaller) Staked(opts *bind.CallOpts, _user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "staked", _user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Staked is a free data retrieval call binding the contract method 0x98807d84.
//
// Solidity: function staked(address _user) view returns(uint256 _staked)
func (_StakingPool *StakingPoolSession) Staked(_user common.Address) (*big.Int, error) {
	return _StakingPool.Contract.Staked(&_StakingPool.CallOpts, _user)
}

// Staked is a free data retrieval call binding the contract method 0x98807d84.
//
// Solidity: function staked(address _user) view returns(uint256 _staked)
func (_StakingPool *StakingPoolCallerSession) Staked(_user common.Address) (*big.Int, error) {
	return _StakingPool.Contract.Staked(&_StakingPool.CallOpts, _user)
}

// StakedNumber is a free data retrieval call binding the contract method 0xaab3a2de.
//
// Solidity: function stakedNumber() view returns(uint256 _stakedNumber)
func (_StakingPool *StakingPoolCaller) StakedNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "stakedNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakedNumber is a free data retrieval call binding the contract method 0xaab3a2de.
//
// Solidity: function stakedNumber() view returns(uint256 _stakedNumber)
func (_StakingPool *StakingPoolSession) StakedNumber() (*big.Int, error) {
	return _StakingPool.Contract.StakedNumber(&_StakingPool.CallOpts)
}

// StakedNumber is a free data retrieval call binding the contract method 0xaab3a2de.
//
// Solidity: function stakedNumber() view returns(uint256 _stakedNumber)
func (_StakingPool *StakingPoolCallerSession) StakedNumber() (*big.Int, error) {
	return _StakingPool.Contract.StakedNumber(&_StakingPool.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingPool *StakingPoolCaller) StakingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "stakingToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingPool *StakingPoolSession) StakingToken() (common.Address, error) {
	return _StakingPool.Contract.StakingToken(&_StakingPool.CallOpts)
}

// StakingToken is a free data retrieval call binding the contract method 0x72f702f3.
//
// Solidity: function stakingToken() view returns(address)
func (_StakingPool *StakingPoolCallerSession) StakingToken() (common.Address, error) {
	return _StakingPool.Contract.StakingToken(&_StakingPool.CallOpts)
}

// TotalClaimed is a free data retrieval call binding the contract method 0xd54ad2a1.
//
// Solidity: function totalClaimed() view returns(uint256)
func (_StakingPool *StakingPoolCaller) TotalClaimed(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "totalClaimed")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalClaimed is a free data retrieval call binding the contract method 0xd54ad2a1.
//
// Solidity: function totalClaimed() view returns(uint256)
func (_StakingPool *StakingPoolSession) TotalClaimed() (*big.Int, error) {
	return _StakingPool.Contract.TotalClaimed(&_StakingPool.CallOpts)
}

// TotalClaimed is a free data retrieval call binding the contract method 0xd54ad2a1.
//
// Solidity: function totalClaimed() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) TotalClaimed() (*big.Int, error) {
	return _StakingPool.Contract.TotalClaimed(&_StakingPool.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingPool *StakingPoolCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "totalStaked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingPool *StakingPoolSession) TotalStaked() (*big.Int, error) {
	return _StakingPool.Contract.TotalStaked(&_StakingPool.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) TotalStaked() (*big.Int, error) {
	return _StakingPool.Contract.TotalStaked(&_StakingPool.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolTransactor) Claim(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "claim")
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolSession) Claim() (*types.Transaction, error) {
	return _StakingPool.Contract.Claim(&_StakingPool.TransactOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x4e71d92d.
//
// Solidity: function claim() returns(uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolTransactorSession) Claim() (*types.Transaction, error) {
	return _StakingPool.Contract.Claim(&_StakingPool.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) payable returns()
func (_StakingPool *StakingPoolTransactor) Deposit(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "deposit", _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) payable returns()
func (_StakingPool *StakingPoolSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.Deposit(&_StakingPool.TransactOpts, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 _amount) payable returns()
func (_StakingPool *StakingPoolTransactorSession) Deposit(_amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.Deposit(&_StakingPool.TransactOpts, _amount)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address _recipient) returns()
func (_StakingPool *StakingPoolTransactor) Destroy(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "destroy", _recipient)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address _recipient) returns()
func (_StakingPool *StakingPoolSession) Destroy(_recipient common.Address) (*types.Transaction, error) {
	return _StakingPool.Contract.Destroy(&_StakingPool.TransactOpts, _recipient)
}

// Destroy is a paid mutator transaction binding the contract method 0x00f55d9d.
//
// Solidity: function destroy(address _recipient) returns()
func (_StakingPool *StakingPoolTransactorSession) Destroy(_recipient common.Address) (*types.Transaction, error) {
	return _StakingPool.Contract.Destroy(&_StakingPool.TransactOpts, _recipient)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingPool *StakingPoolTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingPool *StakingPoolSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.Multicall(&_StakingPool.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingPool *StakingPoolTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.Multicall(&_StakingPool.TransactOpts, data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPool *StakingPoolTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPool *StakingPoolSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingPool.Contract.RenounceOwnership(&_StakingPool.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingPool *StakingPoolTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingPool.Contract.RenounceOwnership(&_StakingPool.TransactOpts)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0xbb872b4a.
//
// Solidity: function setRewardPerBlock(uint256 _rewardPerBlock) returns()
func (_StakingPool *StakingPoolTransactor) SetRewardPerBlock(opts *bind.TransactOpts, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "setRewardPerBlock", _rewardPerBlock)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0xbb872b4a.
//
// Solidity: function setRewardPerBlock(uint256 _rewardPerBlock) returns()
func (_StakingPool *StakingPoolSession) SetRewardPerBlock(_rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.SetRewardPerBlock(&_StakingPool.TransactOpts, _rewardPerBlock)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0xbb872b4a.
//
// Solidity: function setRewardPerBlock(uint256 _rewardPerBlock) returns()
func (_StakingPool *StakingPoolTransactorSession) SetRewardPerBlock(_rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.SetRewardPerBlock(&_StakingPool.TransactOpts, _rewardPerBlock)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPool *StakingPoolTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPool *StakingPoolSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingPool.Contract.TransferOwnership(&_StakingPool.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingPool *StakingPoolTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingPool.Contract.TransferOwnership(&_StakingPool.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns(uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns(uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.Withdraw(&_StakingPool.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns(uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _StakingPool.Contract.Withdraw(&_StakingPool.TransactOpts, _amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakingPool *StakingPoolTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakingPool *StakingPoolSession) Receive() (*types.Transaction, error) {
	return _StakingPool.Contract.Receive(&_StakingPool.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_StakingPool *StakingPoolTransactorSession) Receive() (*types.Transaction, error) {
	return _StakingPool.Contract.Receive(&_StakingPool.TransactOpts)
}

// StakingPoolClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the StakingPool contract.
type StakingPoolClaimedIterator struct {
	Event *StakingPoolClaimed // Event containing the contract specifics and raw log

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
func (it *StakingPoolClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolClaimed)
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
		it.Event = new(StakingPoolClaimed)
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
func (it *StakingPoolClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolClaimed represents a Claimed event raised by the StakingPool contract.
type StakingPoolClaimed struct {
	User   common.Address
	Amount *big.Int
	Height *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x987d620f307ff6b94d58743cb7a7509f24071586a77759b77c2d4e29f75a2f9a.
//
// Solidity: event Claimed(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) FilterClaimed(opts *bind.FilterOpts, _user []common.Address) (*StakingPoolClaimedIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.FilterLogs(opts, "Claimed", _userRule)
	if err != nil {
		return nil, err
	}
	return &StakingPoolClaimedIterator{contract: _StakingPool.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x987d620f307ff6b94d58743cb7a7509f24071586a77759b77c2d4e29f75a2f9a.
//
// Solidity: event Claimed(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *StakingPoolClaimed, _user []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.WatchLogs(opts, "Claimed", _userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolClaimed)
				if err := _StakingPool.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0x987d620f307ff6b94d58743cb7a7509f24071586a77759b77c2d4e29f75a2f9a.
//
// Solidity: event Claimed(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) ParseClaimed(log types.Log) (*StakingPoolClaimed, error) {
	event := new(StakingPoolClaimed)
	if err := _StakingPool.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the StakingPool contract.
type StakingPoolDepositedIterator struct {
	Event *StakingPoolDeposited // Event containing the contract specifics and raw log

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
func (it *StakingPoolDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolDeposited)
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
		it.Event = new(StakingPoolDeposited)
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
func (it *StakingPoolDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolDeposited represents a Deposited event raised by the StakingPool contract.
type StakingPoolDeposited struct {
	User   common.Address
	Amount *big.Int
	Height *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) FilterDeposited(opts *bind.FilterOpts, _user []common.Address) (*StakingPoolDepositedIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.FilterLogs(opts, "Deposited", _userRule)
	if err != nil {
		return nil, err
	}
	return &StakingPoolDepositedIterator{contract: _StakingPool.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *StakingPoolDeposited, _user []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.WatchLogs(opts, "Deposited", _userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolDeposited)
				if err := _StakingPool.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x73a19dd210f1a7f902193214c0ee91dd35ee5b4d920cba8d519eca65a7b488ca.
//
// Solidity: event Deposited(address indexed _user, uint256 _amount, uint256 _height)
func (_StakingPool *StakingPoolFilterer) ParseDeposited(log types.Log) (*StakingPoolDeposited, error) {
	event := new(StakingPoolDeposited)
	if err := _StakingPool.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the StakingPool contract.
type StakingPoolOwnershipTransferredIterator struct {
	Event *StakingPoolOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakingPoolOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolOwnershipTransferred)
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
		it.Event = new(StakingPoolOwnershipTransferred)
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
func (it *StakingPoolOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolOwnershipTransferred represents a OwnershipTransferred event raised by the StakingPool contract.
type StakingPoolOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingPool *StakingPoolFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakingPoolOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingPool.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakingPoolOwnershipTransferredIterator{contract: _StakingPool.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingPool *StakingPoolFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingPoolOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingPool.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolOwnershipTransferred)
				if err := _StakingPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_StakingPool *StakingPoolFilterer) ParseOwnershipTransferred(log types.Log) (*StakingPoolOwnershipTransferred, error) {
	event := new(StakingPoolOwnershipTransferred)
	if err := _StakingPool.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingPoolWithdrawedIterator is returned from FilterWithdrawed and is used to iterate over the raw logs and unpacked data for Withdrawed events raised by the StakingPool contract.
type StakingPoolWithdrawedIterator struct {
	Event *StakingPoolWithdrawed // Event containing the contract specifics and raw log

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
func (it *StakingPoolWithdrawedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingPoolWithdrawed)
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
		it.Event = new(StakingPoolWithdrawed)
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
func (it *StakingPoolWithdrawedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingPoolWithdrawedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingPoolWithdrawed represents a Withdrawed event raised by the StakingPool contract.
type StakingPoolWithdrawed struct {
	User       common.Address
	Withdrawed *big.Int
	Claimed    *big.Int
	Height     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterWithdrawed is a free log retrieval operation binding the contract event 0x992aa94ba074f022eacd638616e322641ce0911ddc743e0dc48c338afafc1ff0.
//
// Solidity: event Withdrawed(address indexed _user, uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolFilterer) FilterWithdrawed(opts *bind.FilterOpts, _user []common.Address) (*StakingPoolWithdrawedIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.FilterLogs(opts, "Withdrawed", _userRule)
	if err != nil {
		return nil, err
	}
	return &StakingPoolWithdrawedIterator{contract: _StakingPool.contract, event: "Withdrawed", logs: logs, sub: sub}, nil
}

// WatchWithdrawed is a free log subscription operation binding the contract event 0x992aa94ba074f022eacd638616e322641ce0911ddc743e0dc48c338afafc1ff0.
//
// Solidity: event Withdrawed(address indexed _user, uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolFilterer) WatchWithdrawed(opts *bind.WatchOpts, sink chan<- *StakingPoolWithdrawed, _user []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}

	logs, sub, err := _StakingPool.contract.WatchLogs(opts, "Withdrawed", _userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingPoolWithdrawed)
				if err := _StakingPool.contract.UnpackLog(event, "Withdrawed", log); err != nil {
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

// ParseWithdrawed is a log parse operation binding the contract event 0x992aa94ba074f022eacd638616e322641ce0911ddc743e0dc48c338afafc1ff0.
//
// Solidity: event Withdrawed(address indexed _user, uint256 _withdrawed, uint256 _claimed, uint256 _height)
func (_StakingPool *StakingPoolFilterer) ParseWithdrawed(log types.Log) (*StakingPoolWithdrawed, error) {
	event := new(StakingPoolWithdrawed)
	if err := _StakingPool.contract.UnpackLog(event, "Withdrawed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
