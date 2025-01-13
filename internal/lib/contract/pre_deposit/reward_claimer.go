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

// RewardClaimerMetaData contains all meta data concerning the RewardClaimer contract.
var RewardClaimerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_signerAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ECDSAInvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"ECDSAInvalidSignatureLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ECDSAInvalidSignatureS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidShortString\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"StringTooLong\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EIP712DomainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"RewardClaimed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"claimReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"claims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"eip712Domain\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"fields\",\"type\":\"bytes1\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"version\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifyingContract\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"extensions\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getUserTotalClaims\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newSigner\",\"type\":\"address\"}],\"name\":\"setSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"signerAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// RewardClaimerABI is the input ABI used to generate the binding from.
// Deprecated: Use RewardClaimerMetaData.ABI instead.
var RewardClaimerABI = RewardClaimerMetaData.ABI

// RewardClaimer is an auto generated Go binding around an Ethereum contract.
type RewardClaimer struct {
	RewardClaimerCaller     // Read-only binding to the contract
	RewardClaimerTransactor // Write-only binding to the contract
	RewardClaimerFilterer   // Log filterer for contract events
}

// RewardClaimerCaller is an auto generated read-only Go binding around an Ethereum contract.
type RewardClaimerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardClaimerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RewardClaimerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardClaimerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RewardClaimerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RewardClaimerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RewardClaimerSession struct {
	Contract     *RewardClaimer    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RewardClaimerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RewardClaimerCallerSession struct {
	Contract *RewardClaimerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// RewardClaimerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RewardClaimerTransactorSession struct {
	Contract     *RewardClaimerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// RewardClaimerRaw is an auto generated low-level Go binding around an Ethereum contract.
type RewardClaimerRaw struct {
	Contract *RewardClaimer // Generic contract binding to access the raw methods on
}

// RewardClaimerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RewardClaimerCallerRaw struct {
	Contract *RewardClaimerCaller // Generic read-only contract binding to access the raw methods on
}

// RewardClaimerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RewardClaimerTransactorRaw struct {
	Contract *RewardClaimerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRewardClaimer creates a new instance of RewardClaimer, bound to a specific deployed contract.
func NewRewardClaimer(address common.Address, backend bind.ContractBackend) (*RewardClaimer, error) {
	contract, err := bindRewardClaimer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RewardClaimer{RewardClaimerCaller: RewardClaimerCaller{contract: contract}, RewardClaimerTransactor: RewardClaimerTransactor{contract: contract}, RewardClaimerFilterer: RewardClaimerFilterer{contract: contract}}, nil
}

// NewRewardClaimerCaller creates a new read-only instance of RewardClaimer, bound to a specific deployed contract.
func NewRewardClaimerCaller(address common.Address, caller bind.ContractCaller) (*RewardClaimerCaller, error) {
	contract, err := bindRewardClaimer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RewardClaimerCaller{contract: contract}, nil
}

// NewRewardClaimerTransactor creates a new write-only instance of RewardClaimer, bound to a specific deployed contract.
func NewRewardClaimerTransactor(address common.Address, transactor bind.ContractTransactor) (*RewardClaimerTransactor, error) {
	contract, err := bindRewardClaimer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RewardClaimerTransactor{contract: contract}, nil
}

// NewRewardClaimerFilterer creates a new log filterer instance of RewardClaimer, bound to a specific deployed contract.
func NewRewardClaimerFilterer(address common.Address, filterer bind.ContractFilterer) (*RewardClaimerFilterer, error) {
	contract, err := bindRewardClaimer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RewardClaimerFilterer{contract: contract}, nil
}

// bindRewardClaimer binds a generic wrapper to an already deployed contract.
func bindRewardClaimer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RewardClaimerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardClaimer *RewardClaimerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardClaimer.Contract.RewardClaimerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardClaimer *RewardClaimerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardClaimer.Contract.RewardClaimerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardClaimer *RewardClaimerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardClaimer.Contract.RewardClaimerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RewardClaimer *RewardClaimerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RewardClaimer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RewardClaimer *RewardClaimerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardClaimer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RewardClaimer *RewardClaimerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RewardClaimer.Contract.contract.Transact(opts, method, params...)
}

// Claims is a free data retrieval call binding the contract method 0x3f733c6b.
//
// Solidity: function claims(address , address ) view returns(uint256)
func (_RewardClaimer *RewardClaimerCaller) Claims(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "claims", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Claims is a free data retrieval call binding the contract method 0x3f733c6b.
//
// Solidity: function claims(address , address ) view returns(uint256)
func (_RewardClaimer *RewardClaimerSession) Claims(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _RewardClaimer.Contract.Claims(&_RewardClaimer.CallOpts, arg0, arg1)
}

// Claims is a free data retrieval call binding the contract method 0x3f733c6b.
//
// Solidity: function claims(address , address ) view returns(uint256)
func (_RewardClaimer *RewardClaimerCallerSession) Claims(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _RewardClaimer.Contract.Claims(&_RewardClaimer.CallOpts, arg0, arg1)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_RewardClaimer *RewardClaimerCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_RewardClaimer *RewardClaimerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _RewardClaimer.Contract.Eip712Domain(&_RewardClaimer.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_RewardClaimer *RewardClaimerCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _RewardClaimer.Contract.Eip712Domain(&_RewardClaimer.CallOpts)
}

// GetUserTotalClaims is a free data retrieval call binding the contract method 0x371a5b34.
//
// Solidity: function getUserTotalClaims(address user, address token) view returns(uint256)
func (_RewardClaimer *RewardClaimerCaller) GetUserTotalClaims(opts *bind.CallOpts, user common.Address, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "getUserTotalClaims", user, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserTotalClaims is a free data retrieval call binding the contract method 0x371a5b34.
//
// Solidity: function getUserTotalClaims(address user, address token) view returns(uint256)
func (_RewardClaimer *RewardClaimerSession) GetUserTotalClaims(user common.Address, token common.Address) (*big.Int, error) {
	return _RewardClaimer.Contract.GetUserTotalClaims(&_RewardClaimer.CallOpts, user, token)
}

// GetUserTotalClaims is a free data retrieval call binding the contract method 0x371a5b34.
//
// Solidity: function getUserTotalClaims(address user, address token) view returns(uint256)
func (_RewardClaimer *RewardClaimerCallerSession) GetUserTotalClaims(user common.Address, token common.Address) (*big.Int, error) {
	return _RewardClaimer.Contract.GetUserTotalClaims(&_RewardClaimer.CallOpts, user, token)
}

// Nonces is a free data retrieval call binding the contract method 0x502e1a16.
//
// Solidity: function nonces(address , uint256 ) view returns(bool)
func (_RewardClaimer *RewardClaimerCaller) Nonces(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "nonces", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x502e1a16.
//
// Solidity: function nonces(address , uint256 ) view returns(bool)
func (_RewardClaimer *RewardClaimerSession) Nonces(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _RewardClaimer.Contract.Nonces(&_RewardClaimer.CallOpts, arg0, arg1)
}

// Nonces is a free data retrieval call binding the contract method 0x502e1a16.
//
// Solidity: function nonces(address , uint256 ) view returns(bool)
func (_RewardClaimer *RewardClaimerCallerSession) Nonces(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _RewardClaimer.Contract.Nonces(&_RewardClaimer.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardClaimer *RewardClaimerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardClaimer *RewardClaimerSession) Owner() (common.Address, error) {
	return _RewardClaimer.Contract.Owner(&_RewardClaimer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_RewardClaimer *RewardClaimerCallerSession) Owner() (common.Address, error) {
	return _RewardClaimer.Contract.Owner(&_RewardClaimer.CallOpts)
}

// SignerAddress is a free data retrieval call binding the contract method 0x5b7633d0.
//
// Solidity: function signerAddress() view returns(address)
func (_RewardClaimer *RewardClaimerCaller) SignerAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _RewardClaimer.contract.Call(opts, &out, "signerAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SignerAddress is a free data retrieval call binding the contract method 0x5b7633d0.
//
// Solidity: function signerAddress() view returns(address)
func (_RewardClaimer *RewardClaimerSession) SignerAddress() (common.Address, error) {
	return _RewardClaimer.Contract.SignerAddress(&_RewardClaimer.CallOpts)
}

// SignerAddress is a free data retrieval call binding the contract method 0x5b7633d0.
//
// Solidity: function signerAddress() view returns(address)
func (_RewardClaimer *RewardClaimerCallerSession) SignerAddress() (common.Address, error) {
	return _RewardClaimer.Contract.SignerAddress(&_RewardClaimer.CallOpts)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x61519f42.
//
// Solidity: function claimReward(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_RewardClaimer *RewardClaimerTransactor) ClaimReward(opts *bind.TransactOpts, token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _RewardClaimer.contract.Transact(opts, "claimReward", token, amount, nonce, signature)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x61519f42.
//
// Solidity: function claimReward(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_RewardClaimer *RewardClaimerSession) ClaimReward(token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _RewardClaimer.Contract.ClaimReward(&_RewardClaimer.TransactOpts, token, amount, nonce, signature)
}

// ClaimReward is a paid mutator transaction binding the contract method 0x61519f42.
//
// Solidity: function claimReward(address token, uint256 amount, uint256 nonce, bytes signature) returns()
func (_RewardClaimer *RewardClaimerTransactorSession) ClaimReward(token common.Address, amount *big.Int, nonce *big.Int, signature []byte) (*types.Transaction, error) {
	return _RewardClaimer.Contract.ClaimReward(&_RewardClaimer.TransactOpts, token, amount, nonce, signature)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardClaimer *RewardClaimerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RewardClaimer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardClaimer *RewardClaimerSession) RenounceOwnership() (*types.Transaction, error) {
	return _RewardClaimer.Contract.RenounceOwnership(&_RewardClaimer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_RewardClaimer *RewardClaimerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _RewardClaimer.Contract.RenounceOwnership(&_RewardClaimer.TransactOpts)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_RewardClaimer *RewardClaimerTransactor) SetSigner(opts *bind.TransactOpts, _newSigner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.contract.Transact(opts, "setSigner", _newSigner)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_RewardClaimer *RewardClaimerSession) SetSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.Contract.SetSigner(&_RewardClaimer.TransactOpts, _newSigner)
}

// SetSigner is a paid mutator transaction binding the contract method 0x6c19e783.
//
// Solidity: function setSigner(address _newSigner) returns()
func (_RewardClaimer *RewardClaimerTransactorSession) SetSigner(_newSigner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.Contract.SetSigner(&_RewardClaimer.TransactOpts, _newSigner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardClaimer *RewardClaimerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardClaimer *RewardClaimerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.Contract.TransferOwnership(&_RewardClaimer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_RewardClaimer *RewardClaimerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RewardClaimer.Contract.TransferOwnership(&_RewardClaimer.TransactOpts, newOwner)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_RewardClaimer *RewardClaimerTransactor) WithdrawToken(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RewardClaimer.contract.Transact(opts, "withdrawToken", token, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_RewardClaimer *RewardClaimerSession) WithdrawToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RewardClaimer.Contract.WithdrawToken(&_RewardClaimer.TransactOpts, token, amount)
}

// WithdrawToken is a paid mutator transaction binding the contract method 0x9e281a98.
//
// Solidity: function withdrawToken(address token, uint256 amount) returns()
func (_RewardClaimer *RewardClaimerTransactorSession) WithdrawToken(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _RewardClaimer.Contract.WithdrawToken(&_RewardClaimer.TransactOpts, token, amount)
}

// RewardClaimerEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the RewardClaimer contract.
type RewardClaimerEIP712DomainChangedIterator struct {
	Event *RewardClaimerEIP712DomainChanged // Event containing the contract specifics and raw log

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
func (it *RewardClaimerEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardClaimerEIP712DomainChanged)
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
		it.Event = new(RewardClaimerEIP712DomainChanged)
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
func (it *RewardClaimerEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardClaimerEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardClaimerEIP712DomainChanged represents a EIP712DomainChanged event raised by the RewardClaimer contract.
type RewardClaimerEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_RewardClaimer *RewardClaimerFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*RewardClaimerEIP712DomainChangedIterator, error) {

	logs, sub, err := _RewardClaimer.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &RewardClaimerEIP712DomainChangedIterator{contract: _RewardClaimer.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_RewardClaimer *RewardClaimerFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *RewardClaimerEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _RewardClaimer.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardClaimerEIP712DomainChanged)
				if err := _RewardClaimer.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
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

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_RewardClaimer *RewardClaimerFilterer) ParseEIP712DomainChanged(log types.Log) (*RewardClaimerEIP712DomainChanged, error) {
	event := new(RewardClaimerEIP712DomainChanged)
	if err := _RewardClaimer.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardClaimerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RewardClaimer contract.
type RewardClaimerOwnershipTransferredIterator struct {
	Event *RewardClaimerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RewardClaimerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardClaimerOwnershipTransferred)
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
		it.Event = new(RewardClaimerOwnershipTransferred)
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
func (it *RewardClaimerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardClaimerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardClaimerOwnershipTransferred represents a OwnershipTransferred event raised by the RewardClaimer contract.
type RewardClaimerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RewardClaimer *RewardClaimerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RewardClaimerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardClaimer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RewardClaimerOwnershipTransferredIterator{contract: _RewardClaimer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_RewardClaimer *RewardClaimerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RewardClaimerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RewardClaimer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardClaimerOwnershipTransferred)
				if err := _RewardClaimer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_RewardClaimer *RewardClaimerFilterer) ParseOwnershipTransferred(log types.Log) (*RewardClaimerOwnershipTransferred, error) {
	event := new(RewardClaimerOwnershipTransferred)
	if err := _RewardClaimer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RewardClaimerRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the RewardClaimer contract.
type RewardClaimerRewardClaimedIterator struct {
	Event *RewardClaimerRewardClaimed // Event containing the contract specifics and raw log

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
func (it *RewardClaimerRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RewardClaimerRewardClaimed)
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
		it.Event = new(RewardClaimerRewardClaimed)
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
func (it *RewardClaimerRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RewardClaimerRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RewardClaimerRewardClaimed represents a RewardClaimed event raised by the RewardClaimer contract.
type RewardClaimerRewardClaimed struct {
	User   common.Address
	Token  common.Address
	Amount *big.Int
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x7a84a08b02c91f3c62d572853f966fc799bbd121e8ad7833a4494ab8dcfcb404.
//
// Solidity: event RewardClaimed(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_RewardClaimer *RewardClaimerFilterer) FilterRewardClaimed(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*RewardClaimerRewardClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _RewardClaimer.contract.FilterLogs(opts, "RewardClaimed", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &RewardClaimerRewardClaimedIterator{contract: _RewardClaimer.contract, event: "RewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x7a84a08b02c91f3c62d572853f966fc799bbd121e8ad7833a4494ab8dcfcb404.
//
// Solidity: event RewardClaimed(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_RewardClaimer *RewardClaimerFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *RewardClaimerRewardClaimed, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _RewardClaimer.contract.WatchLogs(opts, "RewardClaimed", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RewardClaimerRewardClaimed)
				if err := _RewardClaimer.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
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

// ParseRewardClaimed is a log parse operation binding the contract event 0x7a84a08b02c91f3c62d572853f966fc799bbd121e8ad7833a4494ab8dcfcb404.
//
// Solidity: event RewardClaimed(address indexed user, address indexed token, uint256 amount, uint256 nonce)
func (_RewardClaimer *RewardClaimerFilterer) ParseRewardClaimed(log types.Log) (*RewardClaimerRewardClaimed, error) {
	event := new(RewardClaimerRewardClaimed)
	if err := _RewardClaimer.contract.UnpackLog(event, "RewardClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
