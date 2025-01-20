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

// IStakingFactoryStakingPoolInfo is an auto generated low-level Go binding around an user-defined struct.
type IStakingFactoryStakingPoolInfo struct {
	StakingToken         common.Address
	StakingTokenName     string
	StakingTokenSymbol   string
	StakingTokenDecimals uint8
	RewardToken          common.Address
	RewardTokenName      string
	RewardTokenSymbol    string
	RewardTokenDecimals  uint8
	RewardPerBlock       *big.Int
}

// StakingFactoryMetaData contains all meta data concerning the StakingFactory contract.
var StakingFactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"OperatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_stakingToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_rewardToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_rewardPerBlock\",\"type\":\"uint256\"}],\"name\":\"PoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_ethBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_stakingBalance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_rewardBalance\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"PoolDestroied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_old\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_new\",\"type\":\"uint256\"}],\"name\":\"RewardPerBlockSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_rewardToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rewardPerBlock\",\"type\":\"uint256\"}],\"name\":\"createPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"destroyPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_ethBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakingBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rewardBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllPools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_pools\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"}],\"name\":\"getPoolInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"stakingToken\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"stakingTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"stakingTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"stakingTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"rewardTokenName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"rewardTokenSymbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"rewardTokenDecimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"rewardPerBlock\",\"type\":\"uint256\"}],\"internalType\":\"structIStakingFactory.StakingPoolInfo\",\"name\":\"_poolInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getPools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_pools\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_nativeCoinName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_nativeCoinSymbol\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"_nativeCoinDecimals\",\"type\":\"uint8\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"operators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"removeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_rewardPerBlock\",\"type\":\"uint256\"}],\"name\":\"setRewardPerBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakingFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingFactoryMetaData.ABI instead.
var StakingFactoryABI = StakingFactoryMetaData.ABI

// StakingFactory is an auto generated Go binding around an Ethereum contract.
type StakingFactory struct {
	StakingFactoryCaller     // Read-only binding to the contract
	StakingFactoryTransactor // Write-only binding to the contract
	StakingFactoryFilterer   // Log filterer for contract events
}

// StakingFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingFactorySession struct {
	Contract     *StakingFactory   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingFactoryCallerSession struct {
	Contract *StakingFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// StakingFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingFactoryTransactorSession struct {
	Contract     *StakingFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// StakingFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingFactoryRaw struct {
	Contract *StakingFactory // Generic contract binding to access the raw methods on
}

// StakingFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingFactoryCallerRaw struct {
	Contract *StakingFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// StakingFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingFactoryTransactorRaw struct {
	Contract *StakingFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingFactory creates a new instance of StakingFactory, bound to a specific deployed contract.
func NewStakingFactory(address common.Address, backend bind.ContractBackend) (*StakingFactory, error) {
	contract, err := bindStakingFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingFactory{StakingFactoryCaller: StakingFactoryCaller{contract: contract}, StakingFactoryTransactor: StakingFactoryTransactor{contract: contract}, StakingFactoryFilterer: StakingFactoryFilterer{contract: contract}}, nil
}

// NewStakingFactoryCaller creates a new read-only instance of StakingFactory, bound to a specific deployed contract.
func NewStakingFactoryCaller(address common.Address, caller bind.ContractCaller) (*StakingFactoryCaller, error) {
	contract, err := bindStakingFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryCaller{contract: contract}, nil
}

// NewStakingFactoryTransactor creates a new write-only instance of StakingFactory, bound to a specific deployed contract.
func NewStakingFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingFactoryTransactor, error) {
	contract, err := bindStakingFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryTransactor{contract: contract}, nil
}

// NewStakingFactoryFilterer creates a new log filterer instance of StakingFactory, bound to a specific deployed contract.
func NewStakingFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingFactoryFilterer, error) {
	contract, err := bindStakingFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryFilterer{contract: contract}, nil
}

// bindStakingFactory binds a generic wrapper to an already deployed contract.
func bindStakingFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StakingFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingFactory *StakingFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingFactory.Contract.StakingFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingFactory *StakingFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingFactory.Contract.StakingFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingFactory *StakingFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingFactory.Contract.StakingFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingFactory *StakingFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingFactory *StakingFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingFactory *StakingFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingFactory.Contract.contract.Transact(opts, method, params...)
}

// GetAllPools is a free data retrieval call binding the contract method 0xd88ff1f4.
//
// Solidity: function getAllPools() view returns(address[] _pools)
func (_StakingFactory *StakingFactoryCaller) GetAllPools(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _StakingFactory.contract.Call(opts, &out, "getAllPools")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAllPools is a free data retrieval call binding the contract method 0xd88ff1f4.
//
// Solidity: function getAllPools() view returns(address[] _pools)
func (_StakingFactory *StakingFactorySession) GetAllPools() ([]common.Address, error) {
	return _StakingFactory.Contract.GetAllPools(&_StakingFactory.CallOpts)
}

// GetAllPools is a free data retrieval call binding the contract method 0xd88ff1f4.
//
// Solidity: function getAllPools() view returns(address[] _pools)
func (_StakingFactory *StakingFactoryCallerSession) GetAllPools() ([]common.Address, error) {
	return _StakingFactory.Contract.GetAllPools(&_StakingFactory.CallOpts)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x06bfa938.
//
// Solidity: function getPoolInfo(address _pool) view returns((address,string,string,uint8,address,string,string,uint8,uint256) _poolInfo)
func (_StakingFactory *StakingFactoryCaller) GetPoolInfo(opts *bind.CallOpts, _pool common.Address) (IStakingFactoryStakingPoolInfo, error) {
	var out []interface{}
	err := _StakingFactory.contract.Call(opts, &out, "getPoolInfo", _pool)

	if err != nil {
		return *new(IStakingFactoryStakingPoolInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IStakingFactoryStakingPoolInfo)).(*IStakingFactoryStakingPoolInfo)

	return out0, err

}

// GetPoolInfo is a free data retrieval call binding the contract method 0x06bfa938.
//
// Solidity: function getPoolInfo(address _pool) view returns((address,string,string,uint8,address,string,string,uint8,uint256) _poolInfo)
func (_StakingFactory *StakingFactorySession) GetPoolInfo(_pool common.Address) (IStakingFactoryStakingPoolInfo, error) {
	return _StakingFactory.Contract.GetPoolInfo(&_StakingFactory.CallOpts, _pool)
}

// GetPoolInfo is a free data retrieval call binding the contract method 0x06bfa938.
//
// Solidity: function getPoolInfo(address _pool) view returns((address,string,string,uint8,address,string,string,uint8,uint256) _poolInfo)
func (_StakingFactory *StakingFactoryCallerSession) GetPoolInfo(_pool common.Address) (IStakingFactoryStakingPoolInfo, error) {
	return _StakingFactory.Contract.GetPoolInfo(&_StakingFactory.CallOpts, _pool)
}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address _user) view returns(address[] _pools)
func (_StakingFactory *StakingFactoryCaller) GetPools(opts *bind.CallOpts, _user common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _StakingFactory.contract.Call(opts, &out, "getPools", _user)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address _user) view returns(address[] _pools)
func (_StakingFactory *StakingFactorySession) GetPools(_user common.Address) ([]common.Address, error) {
	return _StakingFactory.Contract.GetPools(&_StakingFactory.CallOpts, _user)
}

// GetPools is a free data retrieval call binding the contract method 0x5c39f467.
//
// Solidity: function getPools(address _user) view returns(address[] _pools)
func (_StakingFactory *StakingFactoryCallerSession) GetPools(_user common.Address) ([]common.Address, error) {
	return _StakingFactory.Contract.GetPools(&_StakingFactory.CallOpts, _user)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_StakingFactory *StakingFactoryCaller) Operators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _StakingFactory.contract.Call(opts, &out, "operators", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_StakingFactory *StakingFactorySession) Operators(arg0 common.Address) (bool, error) {
	return _StakingFactory.Contract.Operators(&_StakingFactory.CallOpts, arg0)
}

// Operators is a free data retrieval call binding the contract method 0x13e7c9d8.
//
// Solidity: function operators(address ) view returns(bool)
func (_StakingFactory *StakingFactoryCallerSession) Operators(arg0 common.Address) (bool, error) {
	return _StakingFactory.Contract.Operators(&_StakingFactory.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingFactory *StakingFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingFactory *StakingFactorySession) Owner() (common.Address, error) {
	return _StakingFactory.Contract.Owner(&_StakingFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StakingFactory *StakingFactoryCallerSession) Owner() (common.Address, error) {
	return _StakingFactory.Contract.Owner(&_StakingFactory.CallOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _operator) returns()
func (_StakingFactory *StakingFactoryTransactor) AddOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "addOperator", _operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _operator) returns()
func (_StakingFactory *StakingFactorySession) AddOperator(_operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.AddOperator(&_StakingFactory.TransactOpts, _operator)
}

// AddOperator is a paid mutator transaction binding the contract method 0x9870d7fe.
//
// Solidity: function addOperator(address _operator) returns()
func (_StakingFactory *StakingFactoryTransactorSession) AddOperator(_operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.AddOperator(&_StakingFactory.TransactOpts, _operator)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address _stakingToken, address _rewardToken, uint256 _rewardPerBlock) payable returns(address _pool)
func (_StakingFactory *StakingFactoryTransactor) CreatePool(opts *bind.TransactOpts, _stakingToken common.Address, _rewardToken common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "createPool", _stakingToken, _rewardToken, _rewardPerBlock)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address _stakingToken, address _rewardToken, uint256 _rewardPerBlock) payable returns(address _pool)
func (_StakingFactory *StakingFactorySession) CreatePool(_stakingToken common.Address, _rewardToken common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.Contract.CreatePool(&_StakingFactory.TransactOpts, _stakingToken, _rewardToken, _rewardPerBlock)
}

// CreatePool is a paid mutator transaction binding the contract method 0x51810fb5.
//
// Solidity: function createPool(address _stakingToken, address _rewardToken, uint256 _rewardPerBlock) payable returns(address _pool)
func (_StakingFactory *StakingFactoryTransactorSession) CreatePool(_stakingToken common.Address, _rewardToken common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.Contract.CreatePool(&_StakingFactory.TransactOpts, _stakingToken, _rewardToken, _rewardPerBlock)
}

// DestroyPool is a paid mutator transaction binding the contract method 0x63dbccac.
//
// Solidity: function destroyPool(address _pool, address _recipient) returns(uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance)
func (_StakingFactory *StakingFactoryTransactor) DestroyPool(opts *bind.TransactOpts, _pool common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "destroyPool", _pool, _recipient)
}

// DestroyPool is a paid mutator transaction binding the contract method 0x63dbccac.
//
// Solidity: function destroyPool(address _pool, address _recipient) returns(uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance)
func (_StakingFactory *StakingFactorySession) DestroyPool(_pool common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.DestroyPool(&_StakingFactory.TransactOpts, _pool, _recipient)
}

// DestroyPool is a paid mutator transaction binding the contract method 0x63dbccac.
//
// Solidity: function destroyPool(address _pool, address _recipient) returns(uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance)
func (_StakingFactory *StakingFactoryTransactorSession) DestroyPool(_pool common.Address, _recipient common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.DestroyPool(&_StakingFactory.TransactOpts, _pool, _recipient)
}

// Initialize is a paid mutator transaction binding the contract method 0xf6d2ee86.
//
// Solidity: function initialize(address _owner, string _nativeCoinName, string _nativeCoinSymbol, uint8 _nativeCoinDecimals) returns()
func (_StakingFactory *StakingFactoryTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address, _nativeCoinName string, _nativeCoinSymbol string, _nativeCoinDecimals uint8) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "initialize", _owner, _nativeCoinName, _nativeCoinSymbol, _nativeCoinDecimals)
}

// Initialize is a paid mutator transaction binding the contract method 0xf6d2ee86.
//
// Solidity: function initialize(address _owner, string _nativeCoinName, string _nativeCoinSymbol, uint8 _nativeCoinDecimals) returns()
func (_StakingFactory *StakingFactorySession) Initialize(_owner common.Address, _nativeCoinName string, _nativeCoinSymbol string, _nativeCoinDecimals uint8) (*types.Transaction, error) {
	return _StakingFactory.Contract.Initialize(&_StakingFactory.TransactOpts, _owner, _nativeCoinName, _nativeCoinSymbol, _nativeCoinDecimals)
}

// Initialize is a paid mutator transaction binding the contract method 0xf6d2ee86.
//
// Solidity: function initialize(address _owner, string _nativeCoinName, string _nativeCoinSymbol, uint8 _nativeCoinDecimals) returns()
func (_StakingFactory *StakingFactoryTransactorSession) Initialize(_owner common.Address, _nativeCoinName string, _nativeCoinSymbol string, _nativeCoinDecimals uint8) (*types.Transaction, error) {
	return _StakingFactory.Contract.Initialize(&_StakingFactory.TransactOpts, _owner, _nativeCoinName, _nativeCoinSymbol, _nativeCoinDecimals)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingFactory *StakingFactoryTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingFactory *StakingFactorySession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _StakingFactory.Contract.Multicall(&_StakingFactory.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_StakingFactory *StakingFactoryTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _StakingFactory.Contract.Multicall(&_StakingFactory.TransactOpts, data)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _operator) returns()
func (_StakingFactory *StakingFactoryTransactor) RemoveOperator(opts *bind.TransactOpts, _operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "removeOperator", _operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _operator) returns()
func (_StakingFactory *StakingFactorySession) RemoveOperator(_operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.RemoveOperator(&_StakingFactory.TransactOpts, _operator)
}

// RemoveOperator is a paid mutator transaction binding the contract method 0xac8a584a.
//
// Solidity: function removeOperator(address _operator) returns()
func (_StakingFactory *StakingFactoryTransactorSession) RemoveOperator(_operator common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.RemoveOperator(&_StakingFactory.TransactOpts, _operator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingFactory *StakingFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingFactory *StakingFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingFactory.Contract.RenounceOwnership(&_StakingFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_StakingFactory *StakingFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _StakingFactory.Contract.RenounceOwnership(&_StakingFactory.TransactOpts)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0x5059e016.
//
// Solidity: function setRewardPerBlock(address _pool, uint256 _rewardPerBlock) returns()
func (_StakingFactory *StakingFactoryTransactor) SetRewardPerBlock(opts *bind.TransactOpts, _pool common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "setRewardPerBlock", _pool, _rewardPerBlock)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0x5059e016.
//
// Solidity: function setRewardPerBlock(address _pool, uint256 _rewardPerBlock) returns()
func (_StakingFactory *StakingFactorySession) SetRewardPerBlock(_pool common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.Contract.SetRewardPerBlock(&_StakingFactory.TransactOpts, _pool, _rewardPerBlock)
}

// SetRewardPerBlock is a paid mutator transaction binding the contract method 0x5059e016.
//
// Solidity: function setRewardPerBlock(address _pool, uint256 _rewardPerBlock) returns()
func (_StakingFactory *StakingFactoryTransactorSession) SetRewardPerBlock(_pool common.Address, _rewardPerBlock *big.Int) (*types.Transaction, error) {
	return _StakingFactory.Contract.SetRewardPerBlock(&_StakingFactory.TransactOpts, _pool, _rewardPerBlock)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingFactory *StakingFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _StakingFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingFactory *StakingFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.TransferOwnership(&_StakingFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_StakingFactory *StakingFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _StakingFactory.Contract.TransferOwnership(&_StakingFactory.TransactOpts, newOwner)
}

// StakingFactoryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the StakingFactory contract.
type StakingFactoryInitializedIterator struct {
	Event *StakingFactoryInitialized // Event containing the contract specifics and raw log

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
func (it *StakingFactoryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryInitialized)
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
		it.Event = new(StakingFactoryInitialized)
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
func (it *StakingFactoryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryInitialized represents a Initialized event raised by the StakingFactory contract.
type StakingFactoryInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakingFactory *StakingFactoryFilterer) FilterInitialized(opts *bind.FilterOpts) (*StakingFactoryInitializedIterator, error) {

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &StakingFactoryInitializedIterator{contract: _StakingFactory.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakingFactory *StakingFactoryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *StakingFactoryInitialized) (event.Subscription, error) {

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryInitialized)
				if err := _StakingFactory.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_StakingFactory *StakingFactoryFilterer) ParseInitialized(log types.Log) (*StakingFactoryInitialized, error) {
	event := new(StakingFactoryInitialized)
	if err := _StakingFactory.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the StakingFactory contract.
type StakingFactoryOperatorAddedIterator struct {
	Event *StakingFactoryOperatorAdded // Event containing the contract specifics and raw log

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
func (it *StakingFactoryOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryOperatorAdded)
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
		it.Event = new(StakingFactoryOperatorAdded)
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
func (it *StakingFactoryOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryOperatorAdded represents a OperatorAdded event raised by the StakingFactory contract.
type StakingFactoryOperatorAdded struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) FilterOperatorAdded(opts *bind.FilterOpts, _operator []common.Address) (*StakingFactoryOperatorAddedIterator, error) {

	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "OperatorAdded", _operatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryOperatorAddedIterator{contract: _StakingFactory.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *StakingFactoryOperatorAdded, _operator []common.Address) (event.Subscription, error) {

	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "OperatorAdded", _operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryOperatorAdded)
				if err := _StakingFactory.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ParseOperatorAdded is a log parse operation binding the contract event 0xac6fa858e9350a46cec16539926e0fde25b7629f84b5a72bffaae4df888ae86d.
//
// Solidity: event OperatorAdded(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) ParseOperatorAdded(log types.Log) (*StakingFactoryOperatorAdded, error) {
	event := new(StakingFactoryOperatorAdded)
	if err := _StakingFactory.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryOperatorRemovedIterator is returned from FilterOperatorRemoved and is used to iterate over the raw logs and unpacked data for OperatorRemoved events raised by the StakingFactory contract.
type StakingFactoryOperatorRemovedIterator struct {
	Event *StakingFactoryOperatorRemoved // Event containing the contract specifics and raw log

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
func (it *StakingFactoryOperatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryOperatorRemoved)
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
		it.Event = new(StakingFactoryOperatorRemoved)
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
func (it *StakingFactoryOperatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryOperatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryOperatorRemoved represents a OperatorRemoved event raised by the StakingFactory contract.
type StakingFactoryOperatorRemoved struct {
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorRemoved is a free log retrieval operation binding the contract event 0x80c0b871b97b595b16a7741c1b06fed0c6f6f558639f18ccbce50724325dc40d.
//
// Solidity: event OperatorRemoved(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) FilterOperatorRemoved(opts *bind.FilterOpts, _operator []common.Address) (*StakingFactoryOperatorRemovedIterator, error) {

	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "OperatorRemoved", _operatorRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryOperatorRemovedIterator{contract: _StakingFactory.contract, event: "OperatorRemoved", logs: logs, sub: sub}, nil
}

// WatchOperatorRemoved is a free log subscription operation binding the contract event 0x80c0b871b97b595b16a7741c1b06fed0c6f6f558639f18ccbce50724325dc40d.
//
// Solidity: event OperatorRemoved(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) WatchOperatorRemoved(opts *bind.WatchOpts, sink chan<- *StakingFactoryOperatorRemoved, _operator []common.Address) (event.Subscription, error) {

	var _operatorRule []interface{}
	for _, _operatorItem := range _operator {
		_operatorRule = append(_operatorRule, _operatorItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "OperatorRemoved", _operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryOperatorRemoved)
				if err := _StakingFactory.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
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

// ParseOperatorRemoved is a log parse operation binding the contract event 0x80c0b871b97b595b16a7741c1b06fed0c6f6f558639f18ccbce50724325dc40d.
//
// Solidity: event OperatorRemoved(address indexed _operator)
func (_StakingFactory *StakingFactoryFilterer) ParseOperatorRemoved(log types.Log) (*StakingFactoryOperatorRemoved, error) {
	event := new(StakingFactoryOperatorRemoved)
	if err := _StakingFactory.contract.UnpackLog(event, "OperatorRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the StakingFactory contract.
type StakingFactoryOwnershipTransferredIterator struct {
	Event *StakingFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *StakingFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryOwnershipTransferred)
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
		it.Event = new(StakingFactoryOwnershipTransferred)
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
func (it *StakingFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the StakingFactory contract.
type StakingFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingFactory *StakingFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakingFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryOwnershipTransferredIterator{contract: _StakingFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_StakingFactory *StakingFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakingFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryOwnershipTransferred)
				if err := _StakingFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_StakingFactory *StakingFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*StakingFactoryOwnershipTransferred, error) {
	event := new(StakingFactoryOwnershipTransferred)
	if err := _StakingFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryPoolCreatedIterator is returned from FilterPoolCreated and is used to iterate over the raw logs and unpacked data for PoolCreated events raised by the StakingFactory contract.
type StakingFactoryPoolCreatedIterator struct {
	Event *StakingFactoryPoolCreated // Event containing the contract specifics and raw log

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
func (it *StakingFactoryPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryPoolCreated)
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
		it.Event = new(StakingFactoryPoolCreated)
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
func (it *StakingFactoryPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryPoolCreated represents a PoolCreated event raised by the StakingFactory contract.
type StakingFactoryPoolCreated struct {
	Pool           common.Address
	StakingToken   common.Address
	RewardToken    common.Address
	RewardPerBlock *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPoolCreated is a free log retrieval operation binding the contract event 0xebbbe9dc3a19d2f959ac76ac0372b4983cdfb945f5d6aef4873c36fabb2ba8aa.
//
// Solidity: event PoolCreated(address indexed _pool, address indexed _stakingToken, address indexed _rewardToken, uint256 _rewardPerBlock)
func (_StakingFactory *StakingFactoryFilterer) FilterPoolCreated(opts *bind.FilterOpts, _pool []common.Address, _stakingToken []common.Address, _rewardToken []common.Address) (*StakingFactoryPoolCreatedIterator, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}
	var _stakingTokenRule []interface{}
	for _, _stakingTokenItem := range _stakingToken {
		_stakingTokenRule = append(_stakingTokenRule, _stakingTokenItem)
	}
	var _rewardTokenRule []interface{}
	for _, _rewardTokenItem := range _rewardToken {
		_rewardTokenRule = append(_rewardTokenRule, _rewardTokenItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "PoolCreated", _poolRule, _stakingTokenRule, _rewardTokenRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryPoolCreatedIterator{contract: _StakingFactory.contract, event: "PoolCreated", logs: logs, sub: sub}, nil
}

// WatchPoolCreated is a free log subscription operation binding the contract event 0xebbbe9dc3a19d2f959ac76ac0372b4983cdfb945f5d6aef4873c36fabb2ba8aa.
//
// Solidity: event PoolCreated(address indexed _pool, address indexed _stakingToken, address indexed _rewardToken, uint256 _rewardPerBlock)
func (_StakingFactory *StakingFactoryFilterer) WatchPoolCreated(opts *bind.WatchOpts, sink chan<- *StakingFactoryPoolCreated, _pool []common.Address, _stakingToken []common.Address, _rewardToken []common.Address) (event.Subscription, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}
	var _stakingTokenRule []interface{}
	for _, _stakingTokenItem := range _stakingToken {
		_stakingTokenRule = append(_stakingTokenRule, _stakingTokenItem)
	}
	var _rewardTokenRule []interface{}
	for _, _rewardTokenItem := range _rewardToken {
		_rewardTokenRule = append(_rewardTokenRule, _rewardTokenItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "PoolCreated", _poolRule, _stakingTokenRule, _rewardTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryPoolCreated)
				if err := _StakingFactory.contract.UnpackLog(event, "PoolCreated", log); err != nil {
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

// ParsePoolCreated is a log parse operation binding the contract event 0xebbbe9dc3a19d2f959ac76ac0372b4983cdfb945f5d6aef4873c36fabb2ba8aa.
//
// Solidity: event PoolCreated(address indexed _pool, address indexed _stakingToken, address indexed _rewardToken, uint256 _rewardPerBlock)
func (_StakingFactory *StakingFactoryFilterer) ParsePoolCreated(log types.Log) (*StakingFactoryPoolCreated, error) {
	event := new(StakingFactoryPoolCreated)
	if err := _StakingFactory.contract.UnpackLog(event, "PoolCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryPoolDestroiedIterator is returned from FilterPoolDestroied and is used to iterate over the raw logs and unpacked data for PoolDestroied events raised by the StakingFactory contract.
type StakingFactoryPoolDestroiedIterator struct {
	Event *StakingFactoryPoolDestroied // Event containing the contract specifics and raw log

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
func (it *StakingFactoryPoolDestroiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryPoolDestroied)
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
		it.Event = new(StakingFactoryPoolDestroied)
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
func (it *StakingFactoryPoolDestroiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryPoolDestroiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryPoolDestroied represents a PoolDestroied event raised by the StakingFactory contract.
type StakingFactoryPoolDestroied struct {
	Pool           common.Address
	EthBalance     *big.Int
	StakingBalance *big.Int
	RewardBalance  *big.Int
	Recipient      common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterPoolDestroied is a free log retrieval operation binding the contract event 0x8a1ea1f55ef0525ec0ea0b63984600112e22fa49fff322243c05a6325d1dae03.
//
// Solidity: event PoolDestroied(address indexed _pool, uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance, address indexed _recipient)
func (_StakingFactory *StakingFactoryFilterer) FilterPoolDestroied(opts *bind.FilterOpts, _pool []common.Address, _recipient []common.Address) (*StakingFactoryPoolDestroiedIterator, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}

	var _recipientRule []interface{}
	for _, _recipientItem := range _recipient {
		_recipientRule = append(_recipientRule, _recipientItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "PoolDestroied", _poolRule, _recipientRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryPoolDestroiedIterator{contract: _StakingFactory.contract, event: "PoolDestroied", logs: logs, sub: sub}, nil
}

// WatchPoolDestroied is a free log subscription operation binding the contract event 0x8a1ea1f55ef0525ec0ea0b63984600112e22fa49fff322243c05a6325d1dae03.
//
// Solidity: event PoolDestroied(address indexed _pool, uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance, address indexed _recipient)
func (_StakingFactory *StakingFactoryFilterer) WatchPoolDestroied(opts *bind.WatchOpts, sink chan<- *StakingFactoryPoolDestroied, _pool []common.Address, _recipient []common.Address) (event.Subscription, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}

	var _recipientRule []interface{}
	for _, _recipientItem := range _recipient {
		_recipientRule = append(_recipientRule, _recipientItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "PoolDestroied", _poolRule, _recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryPoolDestroied)
				if err := _StakingFactory.contract.UnpackLog(event, "PoolDestroied", log); err != nil {
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

// ParsePoolDestroied is a log parse operation binding the contract event 0x8a1ea1f55ef0525ec0ea0b63984600112e22fa49fff322243c05a6325d1dae03.
//
// Solidity: event PoolDestroied(address indexed _pool, uint256 _ethBalance, uint256 _stakingBalance, uint256 _rewardBalance, address indexed _recipient)
func (_StakingFactory *StakingFactoryFilterer) ParsePoolDestroied(log types.Log) (*StakingFactoryPoolDestroied, error) {
	event := new(StakingFactoryPoolDestroied)
	if err := _StakingFactory.contract.UnpackLog(event, "PoolDestroied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// StakingFactoryRewardPerBlockSetIterator is returned from FilterRewardPerBlockSet and is used to iterate over the raw logs and unpacked data for RewardPerBlockSet events raised by the StakingFactory contract.
type StakingFactoryRewardPerBlockSetIterator struct {
	Event *StakingFactoryRewardPerBlockSet // Event containing the contract specifics and raw log

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
func (it *StakingFactoryRewardPerBlockSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakingFactoryRewardPerBlockSet)
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
		it.Event = new(StakingFactoryRewardPerBlockSet)
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
func (it *StakingFactoryRewardPerBlockSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakingFactoryRewardPerBlockSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakingFactoryRewardPerBlockSet represents a RewardPerBlockSet event raised by the StakingFactory contract.
type StakingFactoryRewardPerBlockSet struct {
	Pool common.Address
	Old  *big.Int
	New  *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRewardPerBlockSet is a free log retrieval operation binding the contract event 0xfd3fcd47034370bad754e6710d41361299866b465f544bf8512273dce1fe0463.
//
// Solidity: event RewardPerBlockSet(address indexed _pool, uint256 _old, uint256 _new)
func (_StakingFactory *StakingFactoryFilterer) FilterRewardPerBlockSet(opts *bind.FilterOpts, _pool []common.Address) (*StakingFactoryRewardPerBlockSetIterator, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}

	logs, sub, err := _StakingFactory.contract.FilterLogs(opts, "RewardPerBlockSet", _poolRule)
	if err != nil {
		return nil, err
	}
	return &StakingFactoryRewardPerBlockSetIterator{contract: _StakingFactory.contract, event: "RewardPerBlockSet", logs: logs, sub: sub}, nil
}

// WatchRewardPerBlockSet is a free log subscription operation binding the contract event 0xfd3fcd47034370bad754e6710d41361299866b465f544bf8512273dce1fe0463.
//
// Solidity: event RewardPerBlockSet(address indexed _pool, uint256 _old, uint256 _new)
func (_StakingFactory *StakingFactoryFilterer) WatchRewardPerBlockSet(opts *bind.WatchOpts, sink chan<- *StakingFactoryRewardPerBlockSet, _pool []common.Address) (event.Subscription, error) {

	var _poolRule []interface{}
	for _, _poolItem := range _pool {
		_poolRule = append(_poolRule, _poolItem)
	}

	logs, sub, err := _StakingFactory.contract.WatchLogs(opts, "RewardPerBlockSet", _poolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakingFactoryRewardPerBlockSet)
				if err := _StakingFactory.contract.UnpackLog(event, "RewardPerBlockSet", log); err != nil {
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

// ParseRewardPerBlockSet is a log parse operation binding the contract event 0xfd3fcd47034370bad754e6710d41361299866b465f544bf8512273dce1fe0463.
//
// Solidity: event RewardPerBlockSet(address indexed _pool, uint256 _old, uint256 _new)
func (_StakingFactory *StakingFactoryFilterer) ParseRewardPerBlockSet(log types.Log) (*StakingFactoryRewardPerBlockSet, error) {
	event := new(StakingFactoryRewardPerBlockSet)
	if err := _StakingFactory.contract.UnpackLog(event, "RewardPerBlockSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
