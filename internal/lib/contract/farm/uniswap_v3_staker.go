// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package farm

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

// UniswapV3StakerMetaData contains all meta data concerning the UniswapV3Staker contract.
var UniswapV3StakerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIUniswapV3Factory\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"_nonfungiblePositionManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"IncentiveCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"name\":\"IncentiveEnded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"RewardRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"TokenDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"TokenStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"TokenUnstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"}],\"name\":\"createIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"endIncentive\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"contractIUniswapV3Factory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"getIncentiveInfo\",\"outputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalClaimed\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"numberOfStakes\",\"type\":\"uint96\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"getPendingRewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalReward\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"}],\"name\":\"getUserStakes\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"incentives\",\"outputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalClaimed\",\"type\":\"uint256\"},{\"internalType\":\"uint96\",\"name\":\"numberOfStakes\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonfungiblePositionManager\",\"outputs\":[{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"positions\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isStaked\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"}],\"name\":\"removeRewardRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"rewardPools\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastUpdateTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalWeight\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newRatePerSecond\",\"type\":\"uint256\"}],\"name\":\"setRewardRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"incentiveId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// UniswapV3StakerABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV3StakerMetaData.ABI instead.
var UniswapV3StakerABI = UniswapV3StakerMetaData.ABI

// UniswapV3Staker is an auto generated Go binding around an Ethereum contract.
type UniswapV3Staker struct {
	UniswapV3StakerCaller     // Read-only binding to the contract
	UniswapV3StakerTransactor // Write-only binding to the contract
	UniswapV3StakerFilterer   // Log filterer for contract events
}

// UniswapV3StakerCaller is an auto generated read-only Go binding around an Ethereum contract.
type UniswapV3StakerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StakerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UniswapV3StakerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StakerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UniswapV3StakerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV3StakerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UniswapV3StakerSession struct {
	Contract     *UniswapV3Staker  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV3StakerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UniswapV3StakerCallerSession struct {
	Contract *UniswapV3StakerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UniswapV3StakerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UniswapV3StakerTransactorSession struct {
	Contract     *UniswapV3StakerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UniswapV3StakerRaw is an auto generated low-level Go binding around an Ethereum contract.
type UniswapV3StakerRaw struct {
	Contract *UniswapV3Staker // Generic contract binding to access the raw methods on
}

// UniswapV3StakerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UniswapV3StakerCallerRaw struct {
	Contract *UniswapV3StakerCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV3StakerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UniswapV3StakerTransactorRaw struct {
	Contract *UniswapV3StakerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV3Staker creates a new instance of UniswapV3Staker, bound to a specific deployed contract.
func NewUniswapV3Staker(address common.Address, backend bind.ContractBackend) (*UniswapV3Staker, error) {
	contract, err := bindUniswapV3Staker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV3Staker{UniswapV3StakerCaller: UniswapV3StakerCaller{contract: contract}, UniswapV3StakerTransactor: UniswapV3StakerTransactor{contract: contract}, UniswapV3StakerFilterer: UniswapV3StakerFilterer{contract: contract}}, nil
}

// NewUniswapV3StakerCaller creates a new read-only instance of UniswapV3Staker, bound to a specific deployed contract.
func NewUniswapV3StakerCaller(address common.Address, caller bind.ContractCaller) (*UniswapV3StakerCaller, error) {
	contract, err := bindUniswapV3Staker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerCaller{contract: contract}, nil
}

// NewUniswapV3StakerTransactor creates a new write-only instance of UniswapV3Staker, bound to a specific deployed contract.
func NewUniswapV3StakerTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV3StakerTransactor, error) {
	contract, err := bindUniswapV3Staker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTransactor{contract: contract}, nil
}

// NewUniswapV3StakerFilterer creates a new log filterer instance of UniswapV3Staker, bound to a specific deployed contract.
func NewUniswapV3StakerFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV3StakerFilterer, error) {
	contract, err := bindUniswapV3Staker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerFilterer{contract: contract}, nil
}

// bindUniswapV3Staker binds a generic wrapper to an already deployed contract.
func bindUniswapV3Staker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV3StakerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Staker *UniswapV3StakerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Staker.Contract.UniswapV3StakerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Staker *UniswapV3StakerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.UniswapV3StakerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Staker *UniswapV3StakerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.UniswapV3StakerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV3Staker *UniswapV3StakerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV3Staker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV3Staker *UniswapV3StakerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV3Staker *UniswapV3StakerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.contract.Transact(opts, method, params...)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerSession) Factory() (common.Address, error) {
	return _UniswapV3Staker.Contract.Factory(&_UniswapV3Staker.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Factory() (common.Address, error) {
	return _UniswapV3Staker.Contract.Factory(&_UniswapV3Staker.CallOpts)
}

// GetIncentiveInfo is a free data retrieval call binding the contract method 0xd4e9f658.
//
// Solidity: function getIncentiveInfo(uint256 incentiveId) view returns(address rewardToken, address pool, uint256 startTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes, uint256[] tokenIds)
func (_UniswapV3Staker *UniswapV3StakerCaller) GetIncentiveInfo(opts *bind.CallOpts, incentiveId *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
	TokenIds       []*big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "getIncentiveInfo", incentiveId)

	outstruct := new(struct {
		RewardToken    common.Address
		Pool           common.Address
		StartTime      *big.Int
		Weight         *big.Int
		TotalClaimed   *big.Int
		NumberOfStakes *big.Int
		TokenIds       []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RewardToken = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Pool = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.StartTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Weight = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimed = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.NumberOfStakes = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TokenIds = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetIncentiveInfo is a free data retrieval call binding the contract method 0xd4e9f658.
//
// Solidity: function getIncentiveInfo(uint256 incentiveId) view returns(address rewardToken, address pool, uint256 startTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes, uint256[] tokenIds)
func (_UniswapV3Staker *UniswapV3StakerSession) GetIncentiveInfo(incentiveId *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
	TokenIds       []*big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetIncentiveInfo(&_UniswapV3Staker.CallOpts, incentiveId)
}

// GetIncentiveInfo is a free data retrieval call binding the contract method 0xd4e9f658.
//
// Solidity: function getIncentiveInfo(uint256 incentiveId) view returns(address rewardToken, address pool, uint256 startTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes, uint256[] tokenIds)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) GetIncentiveInfo(incentiveId *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
	TokenIds       []*big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetIncentiveInfo(&_UniswapV3Staker.CallOpts, incentiveId)
}

// GetPendingRewards is a free data retrieval call binding the contract method 0xe4e7c6b4.
//
// Solidity: function getPendingRewards(address user, uint256 incentiveId) view returns(uint256 totalReward)
func (_UniswapV3Staker *UniswapV3StakerCaller) GetPendingRewards(opts *bind.CallOpts, user common.Address, incentiveId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "getPendingRewards", user, incentiveId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPendingRewards is a free data retrieval call binding the contract method 0xe4e7c6b4.
//
// Solidity: function getPendingRewards(address user, uint256 incentiveId) view returns(uint256 totalReward)
func (_UniswapV3Staker *UniswapV3StakerSession) GetPendingRewards(user common.Address, incentiveId *big.Int) (*big.Int, error) {
	return _UniswapV3Staker.Contract.GetPendingRewards(&_UniswapV3Staker.CallOpts, user, incentiveId)
}

// GetPendingRewards is a free data retrieval call binding the contract method 0xe4e7c6b4.
//
// Solidity: function getPendingRewards(address user, uint256 incentiveId) view returns(uint256 totalReward)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) GetPendingRewards(user common.Address, incentiveId *big.Int) (*big.Int, error) {
	return _UniswapV3Staker.Contract.GetPendingRewards(&_UniswapV3Staker.CallOpts, user, incentiveId)
}

// GetUserStakes is a free data retrieval call binding the contract method 0x419bc8ad.
//
// Solidity: function getUserStakes(address user, uint256 incentiveId) view returns(uint256[])
func (_UniswapV3Staker *UniswapV3StakerCaller) GetUserStakes(opts *bind.CallOpts, user common.Address, incentiveId *big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "getUserStakes", user, incentiveId)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserStakes is a free data retrieval call binding the contract method 0x419bc8ad.
//
// Solidity: function getUserStakes(address user, uint256 incentiveId) view returns(uint256[])
func (_UniswapV3Staker *UniswapV3StakerSession) GetUserStakes(user common.Address, incentiveId *big.Int) ([]*big.Int, error) {
	return _UniswapV3Staker.Contract.GetUserStakes(&_UniswapV3Staker.CallOpts, user, incentiveId)
}

// GetUserStakes is a free data retrieval call binding the contract method 0x419bc8ad.
//
// Solidity: function getUserStakes(address user, uint256 incentiveId) view returns(uint256[])
func (_UniswapV3Staker *UniswapV3StakerCallerSession) GetUserStakes(user common.Address, incentiveId *big.Int) ([]*big.Int, error) {
	return _UniswapV3Staker.Contract.GetUserStakes(&_UniswapV3Staker.CallOpts, user, incentiveId)
}

// Incentives is a free data retrieval call binding the contract method 0xe121c102.
//
// Solidity: function incentives(uint256 ) view returns(address rewardToken, address pool, uint256 startTime, uint256 endTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerCaller) Incentives(opts *bind.CallOpts, arg0 *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	EndTime        *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "incentives", arg0)

	outstruct := new(struct {
		RewardToken    common.Address
		Pool           common.Address
		StartTime      *big.Int
		EndTime        *big.Int
		Weight         *big.Int
		TotalClaimed   *big.Int
		NumberOfStakes *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RewardToken = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Pool = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.StartTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.EndTime = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Weight = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TotalClaimed = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.NumberOfStakes = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Incentives is a free data retrieval call binding the contract method 0xe121c102.
//
// Solidity: function incentives(uint256 ) view returns(address rewardToken, address pool, uint256 startTime, uint256 endTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerSession) Incentives(arg0 *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	EndTime        *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Incentives(&_UniswapV3Staker.CallOpts, arg0)
}

// Incentives is a free data retrieval call binding the contract method 0xe121c102.
//
// Solidity: function incentives(uint256 ) view returns(address rewardToken, address pool, uint256 startTime, uint256 endTime, uint256 weight, uint256 totalClaimed, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Incentives(arg0 *big.Int) (struct {
	RewardToken    common.Address
	Pool           common.Address
	StartTime      *big.Int
	EndTime        *big.Int
	Weight         *big.Int
	TotalClaimed   *big.Int
	NumberOfStakes *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Incentives(&_UniswapV3Staker.CallOpts, arg0)
}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCaller) NonfungiblePositionManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "nonfungiblePositionManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerSession) NonfungiblePositionManager() (common.Address, error) {
	return _UniswapV3Staker.Contract.NonfungiblePositionManager(&_UniswapV3Staker.CallOpts)
}

// NonfungiblePositionManager is a free data retrieval call binding the contract method 0xb44a2722.
//
// Solidity: function nonfungiblePositionManager() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) NonfungiblePositionManager() (common.Address, error) {
	return _UniswapV3Staker.Contract.NonfungiblePositionManager(&_UniswapV3Staker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerSession) Owner() (common.Address, error) {
	return _UniswapV3Staker.Contract.Owner(&_UniswapV3Staker.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Owner() (common.Address, error) {
	return _UniswapV3Staker.Contract.Owner(&_UniswapV3Staker.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Staker *UniswapV3StakerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Staker *UniswapV3StakerSession) Paused() (bool, error) {
	return _UniswapV3Staker.Contract.Paused(&_UniswapV3Staker.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Paused() (bool, error) {
	return _UniswapV3Staker.Contract.Paused(&_UniswapV3Staker.CallOpts)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(address owner, bool isStaked)
func (_UniswapV3Staker *UniswapV3StakerCaller) Positions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner    common.Address
	IsStaked bool
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "positions", arg0)

	outstruct := new(struct {
		Owner    common.Address
		IsStaked bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.IsStaked = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(address owner, bool isStaked)
func (_UniswapV3Staker *UniswapV3StakerSession) Positions(arg0 *big.Int) (struct {
	Owner    common.Address
	IsStaked bool
}, error) {
	return _UniswapV3Staker.Contract.Positions(&_UniswapV3Staker.CallOpts, arg0)
}

// Positions is a free data retrieval call binding the contract method 0x99fbab88.
//
// Solidity: function positions(uint256 ) view returns(address owner, bool isStaked)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Positions(arg0 *big.Int) (struct {
	Owner    common.Address
	IsStaked bool
}, error) {
	return _UniswapV3Staker.Contract.Positions(&_UniswapV3Staker.CallOpts, arg0)
}

// RewardPools is a free data retrieval call binding the contract method 0x372d9448.
//
// Solidity: function rewardPools(address ) view returns(uint256 totalReward, uint256 lastUpdateTime, uint256 totalWeight)
func (_UniswapV3Staker *UniswapV3StakerCaller) RewardPools(opts *bind.CallOpts, arg0 common.Address) (struct {
	TotalReward    *big.Int
	LastUpdateTime *big.Int
	TotalWeight    *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "rewardPools", arg0)

	outstruct := new(struct {
		TotalReward    *big.Int
		LastUpdateTime *big.Int
		TotalWeight    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalReward = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LastUpdateTime = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TotalWeight = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RewardPools is a free data retrieval call binding the contract method 0x372d9448.
//
// Solidity: function rewardPools(address ) view returns(uint256 totalReward, uint256 lastUpdateTime, uint256 totalWeight)
func (_UniswapV3Staker *UniswapV3StakerSession) RewardPools(arg0 common.Address) (struct {
	TotalReward    *big.Int
	LastUpdateTime *big.Int
	TotalWeight    *big.Int
}, error) {
	return _UniswapV3Staker.Contract.RewardPools(&_UniswapV3Staker.CallOpts, arg0)
}

// RewardPools is a free data retrieval call binding the contract method 0x372d9448.
//
// Solidity: function rewardPools(address ) view returns(uint256 totalReward, uint256 lastUpdateTime, uint256 totalWeight)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) RewardPools(arg0 common.Address) (struct {
	TotalReward    *big.Int
	LastUpdateTime *big.Int
	TotalWeight    *big.Int
}, error) {
	return _UniswapV3Staker.Contract.RewardPools(&_UniswapV3Staker.CallOpts, arg0)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 incentiveId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) Claim(opts *bind.TransactOpts, incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "claim", incentiveId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 incentiveId) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) Claim(incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Claim(&_UniswapV3Staker.TransactOpts, incentiveId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 incentiveId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Claim(incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Claim(&_UniswapV3Staker.TransactOpts, incentiveId)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x23448106.
//
// Solidity: function createIncentive(address rewardToken, address pool, uint256 startTime, uint256 weight) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) CreateIncentive(opts *bind.TransactOpts, rewardToken common.Address, pool common.Address, startTime *big.Int, weight *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "createIncentive", rewardToken, pool, startTime, weight)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x23448106.
//
// Solidity: function createIncentive(address rewardToken, address pool, uint256 startTime, uint256 weight) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) CreateIncentive(rewardToken common.Address, pool common.Address, startTime *big.Int, weight *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.CreateIncentive(&_UniswapV3Staker.TransactOpts, rewardToken, pool, startTime, weight)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x23448106.
//
// Solidity: function createIncentive(address rewardToken, address pool, uint256 startTime, uint256 weight) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) CreateIncentive(rewardToken common.Address, pool common.Address, startTime *big.Int, weight *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.CreateIncentive(&_UniswapV3Staker.TransactOpts, rewardToken, pool, startTime, weight)
}

// EndIncentive is a paid mutator transaction binding the contract method 0x38ee2e30.
//
// Solidity: function endIncentive(uint256 incentiveId) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerTransactor) EndIncentive(opts *bind.TransactOpts, incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "endIncentive", incentiveId)
}

// EndIncentive is a paid mutator transaction binding the contract method 0x38ee2e30.
//
// Solidity: function endIncentive(uint256 incentiveId) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerSession) EndIncentive(incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.EndIncentive(&_UniswapV3Staker.TransactOpts, incentiveId)
}

// EndIncentive is a paid mutator transaction binding the contract method 0x38ee2e30.
//
// Solidity: function endIncentive(uint256 incentiveId) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) EndIncentive(incentiveId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.EndIncentive(&_UniswapV3Staker.TransactOpts, incentiveId)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_UniswapV3Staker *UniswapV3StakerTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_UniswapV3Staker *UniswapV3StakerSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Multicall(&_UniswapV3Staker.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) payable returns(bytes[] results)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Multicall(&_UniswapV3Staker.TransactOpts, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes ) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, from common.Address, tokenId *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "onERC721Received", arg0, from, tokenId, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes ) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerSession) OnERC721Received(arg0 common.Address, from common.Address, tokenId *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.OnERC721Received(&_UniswapV3Staker.TransactOpts, arg0, from, tokenId, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes ) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) OnERC721Received(arg0 common.Address, from common.Address, tokenId *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.OnERC721Received(&_UniswapV3Staker.TransactOpts, arg0, from, tokenId, arg3)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Staker *UniswapV3StakerSession) Pause() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Pause(&_UniswapV3Staker.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Pause() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Pause(&_UniswapV3Staker.TransactOpts)
}

// RemoveRewardRate is a paid mutator transaction binding the contract method 0xdbda470c.
//
// Solidity: function removeRewardRate(address rewardToken) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) RemoveRewardRate(opts *bind.TransactOpts, rewardToken common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "removeRewardRate", rewardToken)
}

// RemoveRewardRate is a paid mutator transaction binding the contract method 0xdbda470c.
//
// Solidity: function removeRewardRate(address rewardToken) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) RemoveRewardRate(rewardToken common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.RemoveRewardRate(&_UniswapV3Staker.TransactOpts, rewardToken)
}

// RemoveRewardRate is a paid mutator transaction binding the contract method 0xdbda470c.
//
// Solidity: function removeRewardRate(address rewardToken) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) RemoveRewardRate(rewardToken common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.RemoveRewardRate(&_UniswapV3Staker.TransactOpts, rewardToken)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV3Staker *UniswapV3StakerSession) RenounceOwnership() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.RenounceOwnership(&_UniswapV3Staker.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.RenounceOwnership(&_UniswapV3Staker.TransactOpts)
}

// SetRewardRate is a paid mutator transaction binding the contract method 0x812a71d5.
//
// Solidity: function setRewardRate(address rewardToken, uint256 newRatePerSecond) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) SetRewardRate(opts *bind.TransactOpts, rewardToken common.Address, newRatePerSecond *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "setRewardRate", rewardToken, newRatePerSecond)
}

// SetRewardRate is a paid mutator transaction binding the contract method 0x812a71d5.
//
// Solidity: function setRewardRate(address rewardToken, uint256 newRatePerSecond) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) SetRewardRate(rewardToken common.Address, newRatePerSecond *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.SetRewardRate(&_UniswapV3Staker.TransactOpts, rewardToken, newRatePerSecond)
}

// SetRewardRate is a paid mutator transaction binding the contract method 0x812a71d5.
//
// Solidity: function setRewardRate(address rewardToken, uint256 newRatePerSecond) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) SetRewardRate(rewardToken common.Address, newRatePerSecond *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.SetRewardRate(&_UniswapV3Staker.TransactOpts, rewardToken, newRatePerSecond)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) Stake(opts *bind.TransactOpts, incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "stake", incentiveId, tokenId)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) Stake(incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Stake(&_UniswapV3Staker.TransactOpts, incentiveId, tokenId)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Stake(incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Stake(&_UniswapV3Staker.TransactOpts, incentiveId, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.TransferOwnership(&_UniswapV3Staker.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.TransferOwnership(&_UniswapV3Staker.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Staker *UniswapV3StakerSession) Unpause() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Unpause(&_UniswapV3Staker.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Unpause() (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Unpause(&_UniswapV3Staker.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x9e2c8a5b.
//
// Solidity: function unstake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) Unstake(opts *bind.TransactOpts, incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "unstake", incentiveId, tokenId)
}

// Unstake is a paid mutator transaction binding the contract method 0x9e2c8a5b.
//
// Solidity: function unstake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) Unstake(incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Unstake(&_UniswapV3Staker.TransactOpts, incentiveId, tokenId)
}

// Unstake is a paid mutator transaction binding the contract method 0x9e2c8a5b.
//
// Solidity: function unstake(uint256 incentiveId, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) Unstake(incentiveId *big.Int, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.Unstake(&_UniswapV3Staker.TransactOpts, incentiveId, tokenId)
}

// UniswapV3StakerIncentiveCreatedIterator is returned from FilterIncentiveCreated and is used to iterate over the raw logs and unpacked data for IncentiveCreated events raised by the UniswapV3Staker contract.
type UniswapV3StakerIncentiveCreatedIterator struct {
	Event *UniswapV3StakerIncentiveCreated // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerIncentiveCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerIncentiveCreated)
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
		it.Event = new(UniswapV3StakerIncentiveCreated)
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
func (it *UniswapV3StakerIncentiveCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerIncentiveCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerIncentiveCreated represents a IncentiveCreated event raised by the UniswapV3Staker contract.
type UniswapV3StakerIncentiveCreated struct {
	RewardToken common.Address
	Pool        common.Address
	StartTime   *big.Int
	Weight      *big.Int
	IncentiveId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIncentiveCreated is a free log retrieval operation binding the contract event 0x03e6bb5c16a3101a95920236f2aa7ccf5f60826ad84f910037336d78745e6859.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 weight, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterIncentiveCreated(opts *bind.FilterOpts, rewardToken []common.Address, pool []common.Address, incentiveId []*big.Int) (*UniswapV3StakerIncentiveCreatedIterator, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "IncentiveCreated", rewardTokenRule, poolRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerIncentiveCreatedIterator{contract: _UniswapV3Staker.contract, event: "IncentiveCreated", logs: logs, sub: sub}, nil
}

// WatchIncentiveCreated is a free log subscription operation binding the contract event 0x03e6bb5c16a3101a95920236f2aa7ccf5f60826ad84f910037336d78745e6859.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 weight, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchIncentiveCreated(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerIncentiveCreated, rewardToken []common.Address, pool []common.Address, incentiveId []*big.Int) (event.Subscription, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "IncentiveCreated", rewardTokenRule, poolRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerIncentiveCreated)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "IncentiveCreated", log); err != nil {
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

// ParseIncentiveCreated is a log parse operation binding the contract event 0x03e6bb5c16a3101a95920236f2aa7ccf5f60826ad84f910037336d78745e6859.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 weight, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseIncentiveCreated(log types.Log) (*UniswapV3StakerIncentiveCreated, error) {
	event := new(UniswapV3StakerIncentiveCreated)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "IncentiveCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerIncentiveEndedIterator is returned from FilterIncentiveEnded and is used to iterate over the raw logs and unpacked data for IncentiveEnded events raised by the UniswapV3Staker contract.
type UniswapV3StakerIncentiveEndedIterator struct {
	Event *UniswapV3StakerIncentiveEnded // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerIncentiveEndedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerIncentiveEnded)
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
		it.Event = new(UniswapV3StakerIncentiveEnded)
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
func (it *UniswapV3StakerIncentiveEndedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerIncentiveEndedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerIncentiveEnded represents a IncentiveEnded event raised by the UniswapV3Staker contract.
type UniswapV3StakerIncentiveEnded struct {
	IncentiveId *big.Int
	Refund      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIncentiveEnded is a free log retrieval operation binding the contract event 0x13f1242f870865ae5fbac1ae3bf7fd3f447fd503ce7444dd46eb6daf0dfcb33e.
//
// Solidity: event IncentiveEnded(uint256 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterIncentiveEnded(opts *bind.FilterOpts, incentiveId []*big.Int) (*UniswapV3StakerIncentiveEndedIterator, error) {

	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "IncentiveEnded", incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerIncentiveEndedIterator{contract: _UniswapV3Staker.contract, event: "IncentiveEnded", logs: logs, sub: sub}, nil
}

// WatchIncentiveEnded is a free log subscription operation binding the contract event 0x13f1242f870865ae5fbac1ae3bf7fd3f447fd503ce7444dd46eb6daf0dfcb33e.
//
// Solidity: event IncentiveEnded(uint256 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchIncentiveEnded(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerIncentiveEnded, incentiveId []*big.Int) (event.Subscription, error) {

	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "IncentiveEnded", incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerIncentiveEnded)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "IncentiveEnded", log); err != nil {
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

// ParseIncentiveEnded is a log parse operation binding the contract event 0x13f1242f870865ae5fbac1ae3bf7fd3f447fd503ce7444dd46eb6daf0dfcb33e.
//
// Solidity: event IncentiveEnded(uint256 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseIncentiveEnded(log types.Log) (*UniswapV3StakerIncentiveEnded, error) {
	event := new(UniswapV3StakerIncentiveEnded)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "IncentiveEnded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the UniswapV3Staker contract.
type UniswapV3StakerOwnershipTransferredIterator struct {
	Event *UniswapV3StakerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerOwnershipTransferred)
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
		it.Event = new(UniswapV3StakerOwnershipTransferred)
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
func (it *UniswapV3StakerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerOwnershipTransferred represents a OwnershipTransferred event raised by the UniswapV3Staker contract.
type UniswapV3StakerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*UniswapV3StakerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerOwnershipTransferredIterator{contract: _UniswapV3Staker.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerOwnershipTransferred)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseOwnershipTransferred(log types.Log) (*UniswapV3StakerOwnershipTransferred, error) {
	event := new(UniswapV3StakerOwnershipTransferred)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the UniswapV3Staker contract.
type UniswapV3StakerPausedIterator struct {
	Event *UniswapV3StakerPaused // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerPaused)
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
		it.Event = new(UniswapV3StakerPaused)
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
func (it *UniswapV3StakerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerPaused represents a Paused event raised by the UniswapV3Staker contract.
type UniswapV3StakerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterPaused(opts *bind.FilterOpts, account []common.Address) (*UniswapV3StakerPausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerPausedIterator{contract: _UniswapV3Staker.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerPaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "Paused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerPaused)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParsePaused(log types.Log) (*UniswapV3StakerPaused, error) {
	event := new(UniswapV3StakerPaused)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the UniswapV3Staker contract.
type UniswapV3StakerRewardClaimedIterator struct {
	Event *UniswapV3StakerRewardClaimed // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerRewardClaimed)
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
		it.Event = new(UniswapV3StakerRewardClaimed)
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
func (it *UniswapV3StakerRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerRewardClaimed represents a RewardClaimed event raised by the UniswapV3Staker contract.
type UniswapV3StakerRewardClaimed struct {
	User        common.Address
	IncentiveId *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0xf01da32686223933d8a18a391060918c7f11a3648639edd87ae013e2e2731743.
//
// Solidity: event RewardClaimed(address indexed user, uint256 indexed incentiveId, uint256 amount)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterRewardClaimed(opts *bind.FilterOpts, user []common.Address, incentiveId []*big.Int) (*UniswapV3StakerRewardClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "RewardClaimed", userRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerRewardClaimedIterator{contract: _UniswapV3Staker.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0xf01da32686223933d8a18a391060918c7f11a3648639edd87ae013e2e2731743.
//
// Solidity: event RewardClaimed(address indexed user, uint256 indexed incentiveId, uint256 amount)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerRewardClaimed, user []common.Address, incentiveId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "RewardClaimed", userRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerRewardClaimed)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
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

// ParseRewardClaimed is a log parse operation binding the contract event 0xf01da32686223933d8a18a391060918c7f11a3648639edd87ae013e2e2731743.
//
// Solidity: event RewardClaimed(address indexed user, uint256 indexed incentiveId, uint256 amount)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseRewardClaimed(log types.Log) (*UniswapV3StakerRewardClaimed, error) {
	event := new(UniswapV3StakerRewardClaimed)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerRewardRateUpdatedIterator is returned from FilterRewardRateUpdated and is used to iterate over the raw logs and unpacked data for RewardRateUpdated events raised by the UniswapV3Staker contract.
type UniswapV3StakerRewardRateUpdatedIterator struct {
	Event *UniswapV3StakerRewardRateUpdated // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerRewardRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerRewardRateUpdated)
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
		it.Event = new(UniswapV3StakerRewardRateUpdated)
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
func (it *UniswapV3StakerRewardRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerRewardRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerRewardRateUpdated represents a RewardRateUpdated event raised by the UniswapV3Staker contract.
type UniswapV3StakerRewardRateUpdated struct {
	RewardToken common.Address
	OldRate     *big.Int
	NewRate     *big.Int
	Timestamp   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRewardRateUpdated is a free log retrieval operation binding the contract event 0xf986b07af98aeade9c7b51c1a6db30c137edfe3691250c21cc2417eea814dc70.
//
// Solidity: event RewardRateUpdated(address indexed rewardToken, uint256 oldRate, uint256 newRate, uint256 timestamp)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterRewardRateUpdated(opts *bind.FilterOpts, rewardToken []common.Address) (*UniswapV3StakerRewardRateUpdatedIterator, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "RewardRateUpdated", rewardTokenRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerRewardRateUpdatedIterator{contract: _UniswapV3Staker.contract, event: "RewardRateUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardRateUpdated is a free log subscription operation binding the contract event 0xf986b07af98aeade9c7b51c1a6db30c137edfe3691250c21cc2417eea814dc70.
//
// Solidity: event RewardRateUpdated(address indexed rewardToken, uint256 oldRate, uint256 newRate, uint256 timestamp)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchRewardRateUpdated(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerRewardRateUpdated, rewardToken []common.Address) (event.Subscription, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "RewardRateUpdated", rewardTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerRewardRateUpdated)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "RewardRateUpdated", log); err != nil {
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

// ParseRewardRateUpdated is a log parse operation binding the contract event 0xf986b07af98aeade9c7b51c1a6db30c137edfe3691250c21cc2417eea814dc70.
//
// Solidity: event RewardRateUpdated(address indexed rewardToken, uint256 oldRate, uint256 newRate, uint256 timestamp)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseRewardRateUpdated(log types.Log) (*UniswapV3StakerRewardRateUpdated, error) {
	event := new(UniswapV3StakerRewardRateUpdated)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "RewardRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerTokenDepositedIterator is returned from FilterTokenDeposited and is used to iterate over the raw logs and unpacked data for TokenDeposited events raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenDepositedIterator struct {
	Event *UniswapV3StakerTokenDeposited // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerTokenDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerTokenDeposited)
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
		it.Event = new(UniswapV3StakerTokenDeposited)
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
func (it *UniswapV3StakerTokenDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerTokenDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerTokenDeposited represents a TokenDeposited event raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenDeposited struct {
	TokenId *big.Int
	Owner   common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTokenDeposited is a free log retrieval operation binding the contract event 0xc7948d631280f3931f619912c8d54c139cc8d369aea73f6c92f7763ea0683dd0.
//
// Solidity: event TokenDeposited(uint256 indexed tokenId, address indexed owner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterTokenDeposited(opts *bind.FilterOpts, tokenId []*big.Int, owner []common.Address) (*UniswapV3StakerTokenDepositedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "TokenDeposited", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTokenDepositedIterator{contract: _UniswapV3Staker.contract, event: "TokenDeposited", logs: logs, sub: sub}, nil
}

// WatchTokenDeposited is a free log subscription operation binding the contract event 0xc7948d631280f3931f619912c8d54c139cc8d369aea73f6c92f7763ea0683dd0.
//
// Solidity: event TokenDeposited(uint256 indexed tokenId, address indexed owner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchTokenDeposited(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerTokenDeposited, tokenId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "TokenDeposited", tokenIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerTokenDeposited)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenDeposited", log); err != nil {
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

// ParseTokenDeposited is a log parse operation binding the contract event 0xc7948d631280f3931f619912c8d54c139cc8d369aea73f6c92f7763ea0683dd0.
//
// Solidity: event TokenDeposited(uint256 indexed tokenId, address indexed owner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseTokenDeposited(log types.Log) (*UniswapV3StakerTokenDeposited, error) {
	event := new(UniswapV3StakerTokenDeposited)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerTokenStakedIterator is returned from FilterTokenStaked and is used to iterate over the raw logs and unpacked data for TokenStaked events raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenStakedIterator struct {
	Event *UniswapV3StakerTokenStaked // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerTokenStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerTokenStaked)
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
		it.Event = new(UniswapV3StakerTokenStaked)
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
func (it *UniswapV3StakerTokenStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerTokenStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerTokenStaked represents a TokenStaked event raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenStaked struct {
	User        common.Address
	TokenId     *big.Int
	IncentiveId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenStaked is a free log retrieval operation binding the contract event 0x6173e4d2d9dd52aae0ed37afed3adcf924a490639b759ca93d32dc43366c17d2.
//
// Solidity: event TokenStaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterTokenStaked(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int, incentiveId []*big.Int) (*UniswapV3StakerTokenStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "TokenStaked", userRule, tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTokenStakedIterator{contract: _UniswapV3Staker.contract, event: "TokenStaked", logs: logs, sub: sub}, nil
}

// WatchTokenStaked is a free log subscription operation binding the contract event 0x6173e4d2d9dd52aae0ed37afed3adcf924a490639b759ca93d32dc43366c17d2.
//
// Solidity: event TokenStaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchTokenStaked(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerTokenStaked, user []common.Address, tokenId []*big.Int, incentiveId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "TokenStaked", userRule, tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerTokenStaked)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenStaked", log); err != nil {
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

// ParseTokenStaked is a log parse operation binding the contract event 0x6173e4d2d9dd52aae0ed37afed3adcf924a490639b759ca93d32dc43366c17d2.
//
// Solidity: event TokenStaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseTokenStaked(log types.Log) (*UniswapV3StakerTokenStaked, error) {
	event := new(UniswapV3StakerTokenStaked)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenStaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerTokenUnstakedIterator is returned from FilterTokenUnstaked and is used to iterate over the raw logs and unpacked data for TokenUnstaked events raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenUnstakedIterator struct {
	Event *UniswapV3StakerTokenUnstaked // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerTokenUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerTokenUnstaked)
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
		it.Event = new(UniswapV3StakerTokenUnstaked)
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
func (it *UniswapV3StakerTokenUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerTokenUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerTokenUnstaked represents a TokenUnstaked event raised by the UniswapV3Staker contract.
type UniswapV3StakerTokenUnstaked struct {
	User        common.Address
	TokenId     *big.Int
	IncentiveId *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenUnstaked is a free log retrieval operation binding the contract event 0x563fc70fca9d440d883299c76a721644d5aca7c45df98aeb5292c4cb823579bf.
//
// Solidity: event TokenUnstaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterTokenUnstaked(opts *bind.FilterOpts, user []common.Address, tokenId []*big.Int, incentiveId []*big.Int) (*UniswapV3StakerTokenUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "TokenUnstaked", userRule, tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTokenUnstakedIterator{contract: _UniswapV3Staker.contract, event: "TokenUnstaked", logs: logs, sub: sub}, nil
}

// WatchTokenUnstaked is a free log subscription operation binding the contract event 0x563fc70fca9d440d883299c76a721644d5aca7c45df98aeb5292c4cb823579bf.
//
// Solidity: event TokenUnstaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchTokenUnstaked(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerTokenUnstaked, user []common.Address, tokenId []*big.Int, incentiveId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "TokenUnstaked", userRule, tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerTokenUnstaked)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenUnstaked", log); err != nil {
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

// ParseTokenUnstaked is a log parse operation binding the contract event 0x563fc70fca9d440d883299c76a721644d5aca7c45df98aeb5292c4cb823579bf.
//
// Solidity: event TokenUnstaked(address indexed user, uint256 indexed tokenId, uint256 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseTokenUnstaked(log types.Log) (*UniswapV3StakerTokenUnstaked, error) {
	event := new(UniswapV3StakerTokenUnstaked)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenUnstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV3StakerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the UniswapV3Staker contract.
type UniswapV3StakerUnpausedIterator struct {
	Event *UniswapV3StakerUnpaused // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerUnpaused)
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
		it.Event = new(UniswapV3StakerUnpaused)
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
func (it *UniswapV3StakerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerUnpaused represents a Unpaused event raised by the UniswapV3Staker contract.
type UniswapV3StakerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterUnpaused(opts *bind.FilterOpts, account []common.Address) (*UniswapV3StakerUnpausedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerUnpausedIterator{contract: _UniswapV3Staker.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerUnpaused, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "Unpaused", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerUnpaused)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed account)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseUnpaused(log types.Log) (*UniswapV3StakerUnpaused, error) {
	event := new(UniswapV3StakerUnpaused)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
