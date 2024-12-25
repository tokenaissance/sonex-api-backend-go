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

// IUniswapV3StakerIncentiveKey is an auto generated low-level Go binding around an user-defined struct.
type IUniswapV3StakerIncentiveKey struct {
	RewardToken common.Address
	Pool        common.Address
	StartTime   *big.Int
	EndTime     *big.Int
	Refundee    common.Address
}

// UniswapV3StakerMetaData contains all meta data concerning the UniswapV3Staker contract.
var UniswapV3StakerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIUniswapV3Factory\",\"name\":\"_factory\",\"type\":\"address\"},{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"_nonfungiblePositionManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_maxIncentiveStartLeadTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_maxIncentiveDuration\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"DepositTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"incentiveId\",\"type\":\"bytes32\"}],\"name\":\"IncentiveCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"incentiveId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"name\":\"IncentiveEnded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"RewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"incentiveId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"}],\"name\":\"TokenStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"incentiveId\",\"type\":\"bytes32\"}],\"name\":\"TokenUnstaked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountRequested\",\"type\":\"uint256\"}],\"name\":\"claimReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"createIncentive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"numberOfStakes\",\"type\":\"uint48\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"}],\"name\":\"endIncentive\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"contractIUniswapV3Factory\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getRewardInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"secondsInsideX128\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"}],\"name\":\"getTotalTokensInPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalToken0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalToken1\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"incentives\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"totalRewardUnclaimed\",\"type\":\"uint256\"},{\"internalType\":\"uint160\",\"name\":\"totalSecondsClaimedX128\",\"type\":\"uint160\"},{\"internalType\":\"uint96\",\"name\":\"numberOfStakes\",\"type\":\"uint96\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxIncentiveDuration\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxIncentiveStartLeadTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nonfungiblePositionManager\",\"outputs\":[{\"internalType\":\"contractINonfungiblePositionManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"rewards\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"stakeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"stakedTokenIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"incentiveId\",\"type\":\"bytes32\"}],\"name\":\"stakes\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"secondsPerLiquidityInsideInitialX128\",\"type\":\"uint160\"},{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"unstakeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20Minimal\",\"name\":\"rewardToken\",\"type\":\"address\"},{\"internalType\":\"contractIUniswapV3Pool\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"refundee\",\"type\":\"address\"}],\"internalType\":\"structIUniswapV3Staker.IncentiveKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address owner, uint48 numberOfStakes, int24 tickLower, int24 tickUpper)
func (_UniswapV3Staker *UniswapV3StakerCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner          common.Address
	NumberOfStakes *big.Int
	TickLower      *big.Int
	TickUpper      *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "deposits", arg0)

	outstruct := new(struct {
		Owner          common.Address
		NumberOfStakes *big.Int
		TickLower      *big.Int
		TickUpper      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.NumberOfStakes = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.TickLower = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.TickUpper = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address owner, uint48 numberOfStakes, int24 tickLower, int24 tickUpper)
func (_UniswapV3Staker *UniswapV3StakerSession) Deposits(arg0 *big.Int) (struct {
	Owner          common.Address
	NumberOfStakes *big.Int
	TickLower      *big.Int
	TickUpper      *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Deposits(&_UniswapV3Staker.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(address owner, uint48 numberOfStakes, int24 tickLower, int24 tickUpper)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Deposits(arg0 *big.Int) (struct {
	Owner          common.Address
	NumberOfStakes *big.Int
	TickLower      *big.Int
	TickUpper      *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Deposits(&_UniswapV3Staker.CallOpts, arg0)
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

// GetRewardInfo is a free data retrieval call binding the contract method 0xd953186e.
//
// Solidity: function getRewardInfo((address,address,uint256,uint256,address) key, uint256 tokenId) view returns(uint256 reward, uint160 secondsInsideX128)
func (_UniswapV3Staker *UniswapV3StakerCaller) GetRewardInfo(opts *bind.CallOpts, key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (struct {
	Reward            *big.Int
	SecondsInsideX128 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "getRewardInfo", key, tokenId)

	outstruct := new(struct {
		Reward            *big.Int
		SecondsInsideX128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reward = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.SecondsInsideX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRewardInfo is a free data retrieval call binding the contract method 0xd953186e.
//
// Solidity: function getRewardInfo((address,address,uint256,uint256,address) key, uint256 tokenId) view returns(uint256 reward, uint160 secondsInsideX128)
func (_UniswapV3Staker *UniswapV3StakerSession) GetRewardInfo(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (struct {
	Reward            *big.Int
	SecondsInsideX128 *big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetRewardInfo(&_UniswapV3Staker.CallOpts, key, tokenId)
}

// GetRewardInfo is a free data retrieval call binding the contract method 0xd953186e.
//
// Solidity: function getRewardInfo((address,address,uint256,uint256,address) key, uint256 tokenId) view returns(uint256 reward, uint160 secondsInsideX128)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) GetRewardInfo(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (struct {
	Reward            *big.Int
	SecondsInsideX128 *big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetRewardInfo(&_UniswapV3Staker.CallOpts, key, tokenId)
}

// GetTotalTokensInPool is a free data retrieval call binding the contract method 0xa9989c18.
//
// Solidity: function getTotalTokensInPool((address,address,uint256,uint256,address) key) view returns(uint256 totalToken0, uint256 totalToken1)
func (_UniswapV3Staker *UniswapV3StakerCaller) GetTotalTokensInPool(opts *bind.CallOpts, key IUniswapV3StakerIncentiveKey) (struct {
	TotalToken0 *big.Int
	TotalToken1 *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "getTotalTokensInPool", key)

	outstruct := new(struct {
		TotalToken0 *big.Int
		TotalToken1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalToken0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalToken1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTotalTokensInPool is a free data retrieval call binding the contract method 0xa9989c18.
//
// Solidity: function getTotalTokensInPool((address,address,uint256,uint256,address) key) view returns(uint256 totalToken0, uint256 totalToken1)
func (_UniswapV3Staker *UniswapV3StakerSession) GetTotalTokensInPool(key IUniswapV3StakerIncentiveKey) (struct {
	TotalToken0 *big.Int
	TotalToken1 *big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetTotalTokensInPool(&_UniswapV3Staker.CallOpts, key)
}

// GetTotalTokensInPool is a free data retrieval call binding the contract method 0xa9989c18.
//
// Solidity: function getTotalTokensInPool((address,address,uint256,uint256,address) key) view returns(uint256 totalToken0, uint256 totalToken1)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) GetTotalTokensInPool(key IUniswapV3StakerIncentiveKey) (struct {
	TotalToken0 *big.Int
	TotalToken1 *big.Int
}, error) {
	return _UniswapV3Staker.Contract.GetTotalTokensInPool(&_UniswapV3Staker.CallOpts, key)
}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256 totalRewardUnclaimed, uint160 totalSecondsClaimedX128, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerCaller) Incentives(opts *bind.CallOpts, arg0 [32]byte) (struct {
	TotalRewardUnclaimed    *big.Int
	TotalSecondsClaimedX128 *big.Int
	NumberOfStakes          *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "incentives", arg0)

	outstruct := new(struct {
		TotalRewardUnclaimed    *big.Int
		TotalSecondsClaimedX128 *big.Int
		NumberOfStakes          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalRewardUnclaimed = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalSecondsClaimedX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NumberOfStakes = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256 totalRewardUnclaimed, uint160 totalSecondsClaimedX128, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerSession) Incentives(arg0 [32]byte) (struct {
	TotalRewardUnclaimed    *big.Int
	TotalSecondsClaimedX128 *big.Int
	NumberOfStakes          *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Incentives(&_UniswapV3Staker.CallOpts, arg0)
}

// Incentives is a free data retrieval call binding the contract method 0x60777795.
//
// Solidity: function incentives(bytes32 ) view returns(uint256 totalRewardUnclaimed, uint160 totalSecondsClaimedX128, uint96 numberOfStakes)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Incentives(arg0 [32]byte) (struct {
	TotalRewardUnclaimed    *big.Int
	TotalSecondsClaimedX128 *big.Int
	NumberOfStakes          *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Incentives(&_UniswapV3Staker.CallOpts, arg0)
}

// MaxIncentiveDuration is a free data retrieval call binding the contract method 0x3dc0714b.
//
// Solidity: function maxIncentiveDuration() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCaller) MaxIncentiveDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "maxIncentiveDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxIncentiveDuration is a free data retrieval call binding the contract method 0x3dc0714b.
//
// Solidity: function maxIncentiveDuration() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerSession) MaxIncentiveDuration() (*big.Int, error) {
	return _UniswapV3Staker.Contract.MaxIncentiveDuration(&_UniswapV3Staker.CallOpts)
}

// MaxIncentiveDuration is a free data retrieval call binding the contract method 0x3dc0714b.
//
// Solidity: function maxIncentiveDuration() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) MaxIncentiveDuration() (*big.Int, error) {
	return _UniswapV3Staker.Contract.MaxIncentiveDuration(&_UniswapV3Staker.CallOpts)
}

// MaxIncentiveStartLeadTime is a free data retrieval call binding the contract method 0x01b75440.
//
// Solidity: function maxIncentiveStartLeadTime() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCaller) MaxIncentiveStartLeadTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "maxIncentiveStartLeadTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxIncentiveStartLeadTime is a free data retrieval call binding the contract method 0x01b75440.
//
// Solidity: function maxIncentiveStartLeadTime() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerSession) MaxIncentiveStartLeadTime() (*big.Int, error) {
	return _UniswapV3Staker.Contract.MaxIncentiveStartLeadTime(&_UniswapV3Staker.CallOpts)
}

// MaxIncentiveStartLeadTime is a free data retrieval call binding the contract method 0x01b75440.
//
// Solidity: function maxIncentiveStartLeadTime() view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) MaxIncentiveStartLeadTime() (*big.Int, error) {
	return _UniswapV3Staker.Contract.MaxIncentiveStartLeadTime(&_UniswapV3Staker.CallOpts)
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

// Rewards is a free data retrieval call binding the contract method 0xe70b9e27.
//
// Solidity: function rewards(address , address ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCaller) Rewards(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "rewards", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Rewards is a free data retrieval call binding the contract method 0xe70b9e27.
//
// Solidity: function rewards(address , address ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerSession) Rewards(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV3Staker.Contract.Rewards(&_UniswapV3Staker.CallOpts, arg0, arg1)
}

// Rewards is a free data retrieval call binding the contract method 0xe70b9e27.
//
// Solidity: function rewards(address , address ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Rewards(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV3Staker.Contract.Rewards(&_UniswapV3Staker.CallOpts, arg0, arg1)
}

// StakedTokenIds is a free data retrieval call binding the contract method 0xa6aed5c4.
//
// Solidity: function stakedTokenIds(bytes32 , uint256 ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCaller) StakedTokenIds(opts *bind.CallOpts, arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "stakedTokenIds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakedTokenIds is a free data retrieval call binding the contract method 0xa6aed5c4.
//
// Solidity: function stakedTokenIds(bytes32 , uint256 ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerSession) StakedTokenIds(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _UniswapV3Staker.Contract.StakedTokenIds(&_UniswapV3Staker.CallOpts, arg0, arg1)
}

// StakedTokenIds is a free data retrieval call binding the contract method 0xa6aed5c4.
//
// Solidity: function stakedTokenIds(bytes32 , uint256 ) view returns(uint256)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) StakedTokenIds(arg0 [32]byte, arg1 *big.Int) (*big.Int, error) {
	return _UniswapV3Staker.Contract.StakedTokenIds(&_UniswapV3Staker.CallOpts, arg0, arg1)
}

// Stakes is a free data retrieval call binding the contract method 0xc36c1ea5.
//
// Solidity: function stakes(uint256 tokenId, bytes32 incentiveId) view returns(uint160 secondsPerLiquidityInsideInitialX128, uint128 liquidity)
func (_UniswapV3Staker *UniswapV3StakerCaller) Stakes(opts *bind.CallOpts, tokenId *big.Int, incentiveId [32]byte) (struct {
	SecondsPerLiquidityInsideInitialX128 *big.Int
	Liquidity                            *big.Int
}, error) {
	var out []interface{}
	err := _UniswapV3Staker.contract.Call(opts, &out, "stakes", tokenId, incentiveId)

	outstruct := new(struct {
		SecondsPerLiquidityInsideInitialX128 *big.Int
		Liquidity                            *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SecondsPerLiquidityInsideInitialX128 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Liquidity = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Stakes is a free data retrieval call binding the contract method 0xc36c1ea5.
//
// Solidity: function stakes(uint256 tokenId, bytes32 incentiveId) view returns(uint160 secondsPerLiquidityInsideInitialX128, uint128 liquidity)
func (_UniswapV3Staker *UniswapV3StakerSession) Stakes(tokenId *big.Int, incentiveId [32]byte) (struct {
	SecondsPerLiquidityInsideInitialX128 *big.Int
	Liquidity                            *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Stakes(&_UniswapV3Staker.CallOpts, tokenId, incentiveId)
}

// Stakes is a free data retrieval call binding the contract method 0xc36c1ea5.
//
// Solidity: function stakes(uint256 tokenId, bytes32 incentiveId) view returns(uint160 secondsPerLiquidityInsideInitialX128, uint128 liquidity)
func (_UniswapV3Staker *UniswapV3StakerCallerSession) Stakes(tokenId *big.Int, incentiveId [32]byte) (struct {
	SecondsPerLiquidityInsideInitialX128 *big.Int
	Liquidity                            *big.Int
}, error) {
	return _UniswapV3Staker.Contract.Stakes(&_UniswapV3Staker.CallOpts, tokenId, incentiveId)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x2f2d783d.
//
// Solidity: function claimReward(address rewardToken, address to, uint256 amountRequested) returns(uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerTransactor) ClaimReward(opts *bind.TransactOpts, rewardToken common.Address, to common.Address, amountRequested *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "claimReward", rewardToken, to, amountRequested)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x2f2d783d.
//
// Solidity: function claimReward(address rewardToken, address to, uint256 amountRequested) returns(uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerSession) ClaimReward(rewardToken common.Address, to common.Address, amountRequested *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.ClaimReward(&_UniswapV3Staker.TransactOpts, rewardToken, to, amountRequested)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x2f2d783d.
//
// Solidity: function claimReward(address rewardToken, address to, uint256 amountRequested) returns(uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) ClaimReward(rewardToken common.Address, to common.Address, amountRequested *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.ClaimReward(&_UniswapV3Staker.TransactOpts, rewardToken, to, amountRequested)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x5cc5e3d9.
//
// Solidity: function createIncentive((address,address,uint256,uint256,address) key, uint256 reward) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) CreateIncentive(opts *bind.TransactOpts, key IUniswapV3StakerIncentiveKey, reward *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "createIncentive", key, reward)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x5cc5e3d9.
//
// Solidity: function createIncentive((address,address,uint256,uint256,address) key, uint256 reward) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) CreateIncentive(key IUniswapV3StakerIncentiveKey, reward *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.CreateIncentive(&_UniswapV3Staker.TransactOpts, key, reward)
}

// CreateIncentive is a paid mutator transaction binding the contract method 0x5cc5e3d9.
//
// Solidity: function createIncentive((address,address,uint256,uint256,address) key, uint256 reward) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) CreateIncentive(key IUniswapV3StakerIncentiveKey, reward *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.CreateIncentive(&_UniswapV3Staker.TransactOpts, key, reward)
}

// EndIncentive is a paid mutator transaction binding the contract method 0xb5ada6e4.
//
// Solidity: function endIncentive((address,address,uint256,uint256,address) key) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerTransactor) EndIncentive(opts *bind.TransactOpts, key IUniswapV3StakerIncentiveKey) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "endIncentive", key)
}

// EndIncentive is a paid mutator transaction binding the contract method 0xb5ada6e4.
//
// Solidity: function endIncentive((address,address,uint256,uint256,address) key) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerSession) EndIncentive(key IUniswapV3StakerIncentiveKey) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.EndIncentive(&_UniswapV3Staker.TransactOpts, key)
}

// EndIncentive is a paid mutator transaction binding the contract method 0xb5ada6e4.
//
// Solidity: function endIncentive((address,address,uint256,uint256,address) key) returns(uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) EndIncentive(key IUniswapV3StakerIncentiveKey) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.EndIncentive(&_UniswapV3Staker.TransactOpts, key)
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
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes data) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "onERC721Received", arg0, from, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes data) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerSession) OnERC721Received(arg0 common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.OnERC721Received(&_UniswapV3Staker.TransactOpts, arg0, from, tokenId, data)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address from, uint256 tokenId, bytes data) returns(bytes4)
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) OnERC721Received(arg0 common.Address, from common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.OnERC721Received(&_UniswapV3Staker.TransactOpts, arg0, from, tokenId, data)
}

// StakeToken is a paid mutator transaction binding the contract method 0xf2d2909b.
//
// Solidity: function stakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) StakeToken(opts *bind.TransactOpts, key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "stakeToken", key, tokenId)
}

// StakeToken is a paid mutator transaction binding the contract method 0xf2d2909b.
//
// Solidity: function stakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) StakeToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.StakeToken(&_UniswapV3Staker.TransactOpts, key, tokenId)
}

// StakeToken is a paid mutator transaction binding the contract method 0xf2d2909b.
//
// Solidity: function stakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) StakeToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.StakeToken(&_UniswapV3Staker.TransactOpts, key, tokenId)
}

// TransferDeposit is a paid mutator transaction binding the contract method 0x26bfee04.
//
// Solidity: function transferDeposit(uint256 tokenId, address to) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) TransferDeposit(opts *bind.TransactOpts, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "transferDeposit", tokenId, to)
}

// TransferDeposit is a paid mutator transaction binding the contract method 0x26bfee04.
//
// Solidity: function transferDeposit(uint256 tokenId, address to) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) TransferDeposit(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.TransferDeposit(&_UniswapV3Staker.TransactOpts, tokenId, to)
}

// TransferDeposit is a paid mutator transaction binding the contract method 0x26bfee04.
//
// Solidity: function transferDeposit(uint256 tokenId, address to) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) TransferDeposit(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.TransferDeposit(&_UniswapV3Staker.TransactOpts, tokenId, to)
}

// UnstakeToken is a paid mutator transaction binding the contract method 0xf549ab42.
//
// Solidity: function unstakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) UnstakeToken(opts *bind.TransactOpts, key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "unstakeToken", key, tokenId)
}

// UnstakeToken is a paid mutator transaction binding the contract method 0xf549ab42.
//
// Solidity: function unstakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) UnstakeToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.UnstakeToken(&_UniswapV3Staker.TransactOpts, key, tokenId)
}

// UnstakeToken is a paid mutator transaction binding the contract method 0xf549ab42.
//
// Solidity: function unstakeToken((address,address,uint256,uint256,address) key, uint256 tokenId) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) UnstakeToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.UnstakeToken(&_UniswapV3Staker.TransactOpts, key, tokenId)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0xb62f046c.
//
// Solidity: function withdrawToken((address,address,uint256,uint256,address) key, uint256 tokenId, address to, bytes data) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactor) WithdrawToken(opts *bind.TransactOpts, key IUniswapV3StakerIncentiveKey, tokenId *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.contract.Transact(opts, "withdrawToken", key, tokenId, to, data)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0xb62f046c.
//
// Solidity: function withdrawToken((address,address,uint256,uint256,address) key, uint256 tokenId, address to, bytes data) returns()
func (_UniswapV3Staker *UniswapV3StakerSession) WithdrawToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.WithdrawToken(&_UniswapV3Staker.TransactOpts, key, tokenId, to, data)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0xb62f046c.
//
// Solidity: function withdrawToken((address,address,uint256,uint256,address) key, uint256 tokenId, address to, bytes data) returns()
func (_UniswapV3Staker *UniswapV3StakerTransactorSession) WithdrawToken(key IUniswapV3StakerIncentiveKey, tokenId *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV3Staker.Contract.WithdrawToken(&_UniswapV3Staker.TransactOpts, key, tokenId, to, data)
}

// UniswapV3StakerDepositTransferredIterator is returned from FilterDepositTransferred and is used to iterate over the raw logs and unpacked data for DepositTransferred events raised by the UniswapV3Staker contract.
type UniswapV3StakerDepositTransferredIterator struct {
	Event *UniswapV3StakerDepositTransferred // Event containing the contract specifics and raw log

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
func (it *UniswapV3StakerDepositTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV3StakerDepositTransferred)
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
		it.Event = new(UniswapV3StakerDepositTransferred)
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
func (it *UniswapV3StakerDepositTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV3StakerDepositTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV3StakerDepositTransferred represents a DepositTransferred event raised by the UniswapV3Staker contract.
type UniswapV3StakerDepositTransferred struct {
	TokenId  *big.Int
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositTransferred is a free log retrieval operation binding the contract event 0xcdfc765b85e1048bee3c6a0f9d1c91fc7c4631f5fe5745a55fc6843db5c3260f.
//
// Solidity: event DepositTransferred(uint256 indexed tokenId, address indexed oldOwner, address indexed newOwner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterDepositTransferred(opts *bind.FilterOpts, tokenId []*big.Int, oldOwner []common.Address, newOwner []common.Address) (*UniswapV3StakerDepositTransferredIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "DepositTransferred", tokenIdRule, oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerDepositTransferredIterator{contract: _UniswapV3Staker.contract, event: "DepositTransferred", logs: logs, sub: sub}, nil
}

// WatchDepositTransferred is a free log subscription operation binding the contract event 0xcdfc765b85e1048bee3c6a0f9d1c91fc7c4631f5fe5745a55fc6843db5c3260f.
//
// Solidity: event DepositTransferred(uint256 indexed tokenId, address indexed oldOwner, address indexed newOwner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchDepositTransferred(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerDepositTransferred, tokenId []*big.Int, oldOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var oldOwnerRule []interface{}
	for _, oldOwnerItem := range oldOwner {
		oldOwnerRule = append(oldOwnerRule, oldOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "DepositTransferred", tokenIdRule, oldOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV3StakerDepositTransferred)
				if err := _UniswapV3Staker.contract.UnpackLog(event, "DepositTransferred", log); err != nil {
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

// ParseDepositTransferred is a log parse operation binding the contract event 0xcdfc765b85e1048bee3c6a0f9d1c91fc7c4631f5fe5745a55fc6843db5c3260f.
//
// Solidity: event DepositTransferred(uint256 indexed tokenId, address indexed oldOwner, address indexed newOwner)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseDepositTransferred(log types.Log) (*UniswapV3StakerDepositTransferred, error) {
	event := new(UniswapV3StakerDepositTransferred)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "DepositTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	EndTime     *big.Int
	Refundee    common.Address
	Reward      *big.Int
	IncentiveId [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIncentiveCreated is a free log retrieval operation binding the contract event 0xb18ff4445e125a58578c06f9ed112ec820db80c93bb32bafa6e23fc75ff53881.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 endTime, address refundee, uint256 reward, bytes32 incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterIncentiveCreated(opts *bind.FilterOpts, rewardToken []common.Address, pool []common.Address) (*UniswapV3StakerIncentiveCreatedIterator, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "IncentiveCreated", rewardTokenRule, poolRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerIncentiveCreatedIterator{contract: _UniswapV3Staker.contract, event: "IncentiveCreated", logs: logs, sub: sub}, nil
}

// WatchIncentiveCreated is a free log subscription operation binding the contract event 0xb18ff4445e125a58578c06f9ed112ec820db80c93bb32bafa6e23fc75ff53881.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 endTime, address refundee, uint256 reward, bytes32 incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchIncentiveCreated(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerIncentiveCreated, rewardToken []common.Address, pool []common.Address) (event.Subscription, error) {

	var rewardTokenRule []interface{}
	for _, rewardTokenItem := range rewardToken {
		rewardTokenRule = append(rewardTokenRule, rewardTokenItem)
	}
	var poolRule []interface{}
	for _, poolItem := range pool {
		poolRule = append(poolRule, poolItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "IncentiveCreated", rewardTokenRule, poolRule)
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

// ParseIncentiveCreated is a log parse operation binding the contract event 0xb18ff4445e125a58578c06f9ed112ec820db80c93bb32bafa6e23fc75ff53881.
//
// Solidity: event IncentiveCreated(address indexed rewardToken, address indexed pool, uint256 startTime, uint256 endTime, address refundee, uint256 reward, bytes32 incentiveId)
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
	IncentiveId [32]byte
	Refund      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterIncentiveEnded is a free log retrieval operation binding the contract event 0x65124e6175aa9904f40735e87e2a37c76e87a609b855287bb4d1aba8257d9763.
//
// Solidity: event IncentiveEnded(bytes32 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterIncentiveEnded(opts *bind.FilterOpts, incentiveId [][32]byte) (*UniswapV3StakerIncentiveEndedIterator, error) {

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

// WatchIncentiveEnded is a free log subscription operation binding the contract event 0x65124e6175aa9904f40735e87e2a37c76e87a609b855287bb4d1aba8257d9763.
//
// Solidity: event IncentiveEnded(bytes32 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchIncentiveEnded(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerIncentiveEnded, incentiveId [][32]byte) (event.Subscription, error) {

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

// ParseIncentiveEnded is a log parse operation binding the contract event 0x65124e6175aa9904f40735e87e2a37c76e87a609b855287bb4d1aba8257d9763.
//
// Solidity: event IncentiveEnded(bytes32 indexed incentiveId, uint256 refund)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseIncentiveEnded(log types.Log) (*UniswapV3StakerIncentiveEnded, error) {
	event := new(UniswapV3StakerIncentiveEnded)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "IncentiveEnded", log); err != nil {
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
	To     common.Address
	Reward *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed to, uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterRewardClaimed(opts *bind.FilterOpts, to []common.Address) (*UniswapV3StakerRewardClaimedIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "RewardClaimed", toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerRewardClaimedIterator{contract: _UniswapV3Staker.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed to, uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerRewardClaimed, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "RewardClaimed", toRule)
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

// ParseRewardClaimed is a log parse operation binding the contract event 0x106f923f993c2149d49b4255ff723acafa1f2d94393f561d3eda32ae348f7241.
//
// Solidity: event RewardClaimed(address indexed to, uint256 reward)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseRewardClaimed(log types.Log) (*UniswapV3StakerRewardClaimed, error) {
	event := new(UniswapV3StakerRewardClaimed)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
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
	TokenId     *big.Int
	IncentiveId [32]byte
	Liquidity   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenStaked is a free log retrieval operation binding the contract event 0x3fe90ccd0a34e28f2b4b7a1e8323415ed9dd595f4eec5dfd461d18c2df336dbd.
//
// Solidity: event TokenStaked(uint256 indexed tokenId, bytes32 indexed incentiveId, uint128 liquidity)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterTokenStaked(opts *bind.FilterOpts, tokenId []*big.Int, incentiveId [][32]byte) (*UniswapV3StakerTokenStakedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "TokenStaked", tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTokenStakedIterator{contract: _UniswapV3Staker.contract, event: "TokenStaked", logs: logs, sub: sub}, nil
}

// WatchTokenStaked is a free log subscription operation binding the contract event 0x3fe90ccd0a34e28f2b4b7a1e8323415ed9dd595f4eec5dfd461d18c2df336dbd.
//
// Solidity: event TokenStaked(uint256 indexed tokenId, bytes32 indexed incentiveId, uint128 liquidity)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchTokenStaked(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerTokenStaked, tokenId []*big.Int, incentiveId [][32]byte) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "TokenStaked", tokenIdRule, incentiveIdRule)
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

// ParseTokenStaked is a log parse operation binding the contract event 0x3fe90ccd0a34e28f2b4b7a1e8323415ed9dd595f4eec5dfd461d18c2df336dbd.
//
// Solidity: event TokenStaked(uint256 indexed tokenId, bytes32 indexed incentiveId, uint128 liquidity)
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
	TokenId     *big.Int
	IncentiveId [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTokenUnstaked is a free log retrieval operation binding the contract event 0xe1ba67e807ae0efa0a9549f9520ddc15c27f0a4dae2bc045e800ca66a940778f.
//
// Solidity: event TokenUnstaked(uint256 indexed tokenId, bytes32 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) FilterTokenUnstaked(opts *bind.FilterOpts, tokenId []*big.Int, incentiveId [][32]byte) (*UniswapV3StakerTokenUnstakedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.FilterLogs(opts, "TokenUnstaked", tokenIdRule, incentiveIdRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV3StakerTokenUnstakedIterator{contract: _UniswapV3Staker.contract, event: "TokenUnstaked", logs: logs, sub: sub}, nil
}

// WatchTokenUnstaked is a free log subscription operation binding the contract event 0xe1ba67e807ae0efa0a9549f9520ddc15c27f0a4dae2bc045e800ca66a940778f.
//
// Solidity: event TokenUnstaked(uint256 indexed tokenId, bytes32 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) WatchTokenUnstaked(opts *bind.WatchOpts, sink chan<- *UniswapV3StakerTokenUnstaked, tokenId []*big.Int, incentiveId [][32]byte) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var incentiveIdRule []interface{}
	for _, incentiveIdItem := range incentiveId {
		incentiveIdRule = append(incentiveIdRule, incentiveIdItem)
	}

	logs, sub, err := _UniswapV3Staker.contract.WatchLogs(opts, "TokenUnstaked", tokenIdRule, incentiveIdRule)
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

// ParseTokenUnstaked is a log parse operation binding the contract event 0xe1ba67e807ae0efa0a9549f9520ddc15c27f0a4dae2bc045e800ca66a940778f.
//
// Solidity: event TokenUnstaked(uint256 indexed tokenId, bytes32 indexed incentiveId)
func (_UniswapV3Staker *UniswapV3StakerFilterer) ParseTokenUnstaked(log types.Log) (*UniswapV3StakerTokenUnstaked, error) {
	event := new(UniswapV3StakerTokenUnstaked)
	if err := _UniswapV3Staker.contract.UnpackLog(event, "TokenUnstaked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
