// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// ConsensusValidatorEntrypointMetaData contains all meta data concerning the ConsensusValidatorEntrypoint contract.
var ConsensusValidatorEntrypointMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"initialCollateralOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferCollateralOwnership\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"prevOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjail\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateExtraVotingPower\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"extraVotingPower\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateralOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maturesAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgDepositCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"collateralOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgRegisterValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"initialCollateralOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"initialCollateralAmountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgTransferCollateralOwnership\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"prevOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgUnjail\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgUpdateExtraVotingPower\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"extraVotingPowerWei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgWithdrawCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"collateralOwner\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amountGwei\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maturesAt\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PermittedCallerSet\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidParameter\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"}]}]",
}

// ConsensusValidatorEntrypointABI is the input ABI used to generate the binding from.
// Deprecated: Use ConsensusValidatorEntrypointMetaData.ABI instead.
var ConsensusValidatorEntrypointABI = ConsensusValidatorEntrypointMetaData.ABI

// ConsensusValidatorEntrypoint is an auto generated Go binding around an Ethereum contract.
type ConsensusValidatorEntrypoint struct {
	ConsensusValidatorEntrypointCaller     // Read-only binding to the contract
	ConsensusValidatorEntrypointTransactor // Write-only binding to the contract
	ConsensusValidatorEntrypointFilterer   // Log filterer for contract events
}

// ConsensusValidatorEntrypointCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusValidatorEntrypointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusValidatorEntrypointSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusValidatorEntrypointSession struct {
	Contract     *ConsensusValidatorEntrypoint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ConsensusValidatorEntrypointCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusValidatorEntrypointCallerSession struct {
	Contract *ConsensusValidatorEntrypointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// ConsensusValidatorEntrypointTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusValidatorEntrypointTransactorSession struct {
	Contract     *ConsensusValidatorEntrypointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// ConsensusValidatorEntrypointRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointRaw struct {
	Contract *ConsensusValidatorEntrypoint // Generic contract binding to access the raw methods on
}

// ConsensusValidatorEntrypointCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointCallerRaw struct {
	Contract *ConsensusValidatorEntrypointCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusValidatorEntrypointTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusValidatorEntrypointTransactorRaw struct {
	Contract *ConsensusValidatorEntrypointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensusValidatorEntrypoint creates a new instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypoint(address common.Address, backend bind.ContractBackend) (*ConsensusValidatorEntrypoint, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypoint{ConsensusValidatorEntrypointCaller: ConsensusValidatorEntrypointCaller{contract: contract}, ConsensusValidatorEntrypointTransactor: ConsensusValidatorEntrypointTransactor{contract: contract}, ConsensusValidatorEntrypointFilterer: ConsensusValidatorEntrypointFilterer{contract: contract}}, nil
}

// NewConsensusValidatorEntrypointCaller creates a new read-only instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointCaller(address common.Address, caller bind.ContractCaller) (*ConsensusValidatorEntrypointCaller, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointCaller{contract: contract}, nil
}

// NewConsensusValidatorEntrypointTransactor creates a new write-only instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusValidatorEntrypointTransactor, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointTransactor{contract: contract}, nil
}

// NewConsensusValidatorEntrypointFilterer creates a new log filterer instance of ConsensusValidatorEntrypoint, bound to a specific deployed contract.
func NewConsensusValidatorEntrypointFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusValidatorEntrypointFilterer, error) {
	contract, err := bindConsensusValidatorEntrypoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointFilterer{contract: contract}, nil
}

// bindConsensusValidatorEntrypoint binds a generic wrapper to an already deployed contract.
func bindConsensusValidatorEntrypoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConsensusValidatorEntrypointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.ConsensusValidatorEntrypointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusValidatorEntrypoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusValidatorEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusValidatorEntrypoint.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusValidatorEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusValidatorEntrypoint.CallOpts)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) IsPermittedCaller(opts *bind.CallOpts, caller common.Address) (bool, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "isPermittedCaller", caller)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusValidatorEntrypoint.Contract.IsPermittedCaller(&_ConsensusValidatorEntrypoint.CallOpts, caller)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusValidatorEntrypoint.Contract.IsPermittedCaller(&_ConsensusValidatorEntrypoint.CallOpts, caller)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Owner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.Owner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) Owner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.Owner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) PendingOwner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.PendingOwner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) PendingOwner() (common.Address, error) {
	return _ConsensusValidatorEntrypoint.Contract.PendingOwner(&_ConsensusValidatorEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ConsensusValidatorEntrypoint.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusValidatorEntrypoint.Contract.ProxiableUUID(&_ConsensusValidatorEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusValidatorEntrypoint.Contract.ProxiableUUID(&_ConsensusValidatorEntrypoint.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.AcceptOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.AcceptOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x9a3d08fb.
//
// Solidity: function depositCollateral(address valAddr, address collateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) DepositCollateral(opts *bind.TransactOpts, valAddr common.Address, collateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "depositCollateral", valAddr, collateralOwner)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x9a3d08fb.
//
// Solidity: function depositCollateral(address valAddr, address collateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) DepositCollateral(valAddr common.Address, collateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.DepositCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, collateralOwner)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x9a3d08fb.
//
// Solidity: function depositCollateral(address valAddr, address collateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) DepositCollateral(valAddr common.Address, collateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.DepositCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, collateralOwner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Initialize(opts *bind.TransactOpts, owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "initialize", owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Initialize(&_ConsensusValidatorEntrypoint.TransactOpts, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Initialize(&_ConsensusValidatorEntrypoint.TransactOpts, owner_)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbc990314.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey, address initialCollateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) RegisterValidator(opts *bind.TransactOpts, valAddr common.Address, pubKey []byte, initialCollateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "registerValidator", valAddr, pubKey, initialCollateralOwner)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbc990314.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey, address initialCollateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) RegisterValidator(valAddr common.Address, pubKey []byte, initialCollateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RegisterValidator(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, pubKey, initialCollateralOwner)
}

// RegisterValidator is a paid mutator transaction binding the contract method 0xbc990314.
//
// Solidity: function registerValidator(address valAddr, bytes pubKey, address initialCollateralOwner) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) RegisterValidator(valAddr common.Address, pubKey []byte, initialCollateralOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RegisterValidator(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, pubKey, initialCollateralOwner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RenounceOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.RenounceOwnership(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) SetPermittedCaller(opts *bind.TransactOpts, caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "setPermittedCaller", caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.SetPermittedCaller(&_ConsensusValidatorEntrypoint.TransactOpts, caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.SetPermittedCaller(&_ConsensusValidatorEntrypoint.TransactOpts, caller, isPermitted)
}

// TransferCollateralOwnership is a paid mutator transaction binding the contract method 0x4fc85a0c.
//
// Solidity: function transferCollateralOwnership(address valAddr, address prevOwner, address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) TransferCollateralOwnership(opts *bind.TransactOpts, valAddr common.Address, prevOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "transferCollateralOwnership", valAddr, prevOwner, newOwner)
}

// TransferCollateralOwnership is a paid mutator transaction binding the contract method 0x4fc85a0c.
//
// Solidity: function transferCollateralOwnership(address valAddr, address prevOwner, address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) TransferCollateralOwnership(valAddr common.Address, prevOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferCollateralOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, prevOwner, newOwner)
}

// TransferCollateralOwnership is a paid mutator transaction binding the contract method 0x4fc85a0c.
//
// Solidity: function transferCollateralOwnership(address valAddr, address prevOwner, address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) TransferCollateralOwnership(valAddr common.Address, prevOwner common.Address, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferCollateralOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, prevOwner, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.TransferOwnership(&_ConsensusValidatorEntrypoint.TransactOpts, newOwner)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Unjail(opts *bind.TransactOpts, valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "unjail", valAddr)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Unjail(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Unjail(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// Unjail is a paid mutator transaction binding the contract method 0x449ecfe6.
//
// Solidity: function unjail(address valAddr) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Unjail(valAddr common.Address) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Unjail(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) UpdateExtraVotingPower(opts *bind.TransactOpts, valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "updateExtraVotingPower", valAddr, extraVotingPower)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UpdateExtraVotingPower(valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpdateExtraVotingPower(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, extraVotingPower)
}

// UpdateExtraVotingPower is a paid mutator transaction binding the contract method 0xf09d3149.
//
// Solidity: function updateExtraVotingPower(address valAddr, uint256 extraVotingPower) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) UpdateExtraVotingPower(valAddr common.Address, extraVotingPower *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpdateExtraVotingPower(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, extraVotingPower)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpgradeToAndCall(&_ConsensusValidatorEntrypoint.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.UpgradeToAndCall(&_ConsensusValidatorEntrypoint.TransactOpts, newImplementation, data)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59fbcd70.
//
// Solidity: function withdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amount, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) WithdrawCollateral(opts *bind.TransactOpts, valAddr common.Address, collateralOwner common.Address, receiver common.Address, amount *big.Int, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.Transact(opts, "withdrawCollateral", valAddr, collateralOwner, receiver, amount, maturesAt)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59fbcd70.
//
// Solidity: function withdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amount, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) WithdrawCollateral(valAddr common.Address, collateralOwner common.Address, receiver common.Address, amount *big.Int, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.WithdrawCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, collateralOwner, receiver, amount, maturesAt)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x59fbcd70.
//
// Solidity: function withdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amount, uint48 maturesAt) returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) WithdrawCollateral(valAddr common.Address, collateralOwner common.Address, receiver common.Address, amount *big.Int, maturesAt *big.Int) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.WithdrawCollateral(&_ConsensusValidatorEntrypoint.TransactOpts, valAddr, collateralOwner, receiver, amount, maturesAt)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Fallback(&_ConsensusValidatorEntrypoint.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Fallback(&_ConsensusValidatorEntrypoint.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointSession) Receive() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Receive(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointTransactorSession) Receive() (*types.Transaction, error) {
	return _ConsensusValidatorEntrypoint.Contract.Receive(&_ConsensusValidatorEntrypoint.TransactOpts)
}

// ConsensusValidatorEntrypointInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointInitializedIterator struct {
	Event *ConsensusValidatorEntrypointInitialized // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointInitialized)
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
		it.Event = new(ConsensusValidatorEntrypointInitialized)
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
func (it *ConsensusValidatorEntrypointInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointInitialized represents a Initialized event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterInitialized(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointInitializedIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointInitializedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointInitialized) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointInitialized)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseInitialized(log types.Log) (*ConsensusValidatorEntrypointInitialized, error) {
	event := new(ConsensusValidatorEntrypointInitialized)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgDepositCollateralIterator is returned from FilterMsgDepositCollateral and is used to iterate over the raw logs and unpacked data for MsgDepositCollateral events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgDepositCollateralIterator struct {
	Event *ConsensusValidatorEntrypointMsgDepositCollateral // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgDepositCollateral)
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
		it.Event = new(ConsensusValidatorEntrypointMsgDepositCollateral)
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
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgDepositCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgDepositCollateral represents a MsgDepositCollateral event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgDepositCollateral struct {
	ValAddr         common.Address
	CollateralOwner common.Address
	AmountGwei      *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMsgDepositCollateral is a free log retrieval operation binding the contract event 0xa7ec5f2adeb921d4cf9f5086c7b8d1d983b5a72ab97d00cd0cf60006449e4cf2.
//
// Solidity: event MsgDepositCollateral(address valAddr, address collateralOwner, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgDepositCollateral(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgDepositCollateralIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgDepositCollateral")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgDepositCollateralIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgDepositCollateral", logs: logs, sub: sub}, nil
}

// WatchMsgDepositCollateral is a free log subscription operation binding the contract event 0xa7ec5f2adeb921d4cf9f5086c7b8d1d983b5a72ab97d00cd0cf60006449e4cf2.
//
// Solidity: event MsgDepositCollateral(address valAddr, address collateralOwner, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgDepositCollateral(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgDepositCollateral) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgDepositCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgDepositCollateral)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgDepositCollateral", log); err != nil {
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

// ParseMsgDepositCollateral is a log parse operation binding the contract event 0xa7ec5f2adeb921d4cf9f5086c7b8d1d983b5a72ab97d00cd0cf60006449e4cf2.
//
// Solidity: event MsgDepositCollateral(address valAddr, address collateralOwner, uint256 amountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgDepositCollateral(log types.Log) (*ConsensusValidatorEntrypointMsgDepositCollateral, error) {
	event := new(ConsensusValidatorEntrypointMsgDepositCollateral)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgDepositCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgRegisterValidatorIterator is returned from FilterMsgRegisterValidator and is used to iterate over the raw logs and unpacked data for MsgRegisterValidator events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgRegisterValidatorIterator struct {
	Event *ConsensusValidatorEntrypointMsgRegisterValidator // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgRegisterValidator)
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
		it.Event = new(ConsensusValidatorEntrypointMsgRegisterValidator)
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
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgRegisterValidatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgRegisterValidator represents a MsgRegisterValidator event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgRegisterValidator struct {
	ValAddr                     common.Address
	PubKey                      []byte
	InitialCollateralOwner      common.Address
	InitialCollateralAmountGwei *big.Int
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterMsgRegisterValidator is a free log retrieval operation binding the contract event 0x4eb098ea3e3d659c98c11bffb8d2a6d6d31607a54e122b807b42605247cdd5c6.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, address initialCollateralOwner, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgRegisterValidator(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgRegisterValidatorIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgRegisterValidator")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgRegisterValidatorIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgRegisterValidator", logs: logs, sub: sub}, nil
}

// WatchMsgRegisterValidator is a free log subscription operation binding the contract event 0x4eb098ea3e3d659c98c11bffb8d2a6d6d31607a54e122b807b42605247cdd5c6.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, address initialCollateralOwner, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgRegisterValidator(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgRegisterValidator) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgRegisterValidator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgRegisterValidator)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgRegisterValidator", log); err != nil {
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

// ParseMsgRegisterValidator is a log parse operation binding the contract event 0x4eb098ea3e3d659c98c11bffb8d2a6d6d31607a54e122b807b42605247cdd5c6.
//
// Solidity: event MsgRegisterValidator(address valAddr, bytes pubKey, address initialCollateralOwner, uint256 initialCollateralAmountGwei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgRegisterValidator(log types.Log) (*ConsensusValidatorEntrypointMsgRegisterValidator, error) {
	event := new(ConsensusValidatorEntrypointMsgRegisterValidator)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgRegisterValidator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator is returned from FilterMsgTransferCollateralOwnership and is used to iterate over the raw logs and unpacked data for MsgTransferCollateralOwnership events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator struct {
	Event *ConsensusValidatorEntrypointMsgTransferCollateralOwnership // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgTransferCollateralOwnership)
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
		it.Event = new(ConsensusValidatorEntrypointMsgTransferCollateralOwnership)
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
func (it *ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgTransferCollateralOwnership represents a MsgTransferCollateralOwnership event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgTransferCollateralOwnership struct {
	ValAddr   common.Address
	PrevOwner common.Address
	NewOwner  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMsgTransferCollateralOwnership is a free log retrieval operation binding the contract event 0x9d6b2fb9b2bb5189c38dbcd4827c1d1af06009156feba0a43821e8a616d5bee3.
//
// Solidity: event MsgTransferCollateralOwnership(address valAddr, address prevOwner, address newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgTransferCollateralOwnership(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgTransferCollateralOwnership")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgTransferCollateralOwnershipIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgTransferCollateralOwnership", logs: logs, sub: sub}, nil
}

// WatchMsgTransferCollateralOwnership is a free log subscription operation binding the contract event 0x9d6b2fb9b2bb5189c38dbcd4827c1d1af06009156feba0a43821e8a616d5bee3.
//
// Solidity: event MsgTransferCollateralOwnership(address valAddr, address prevOwner, address newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgTransferCollateralOwnership(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgTransferCollateralOwnership) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgTransferCollateralOwnership")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgTransferCollateralOwnership)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgTransferCollateralOwnership", log); err != nil {
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

// ParseMsgTransferCollateralOwnership is a log parse operation binding the contract event 0x9d6b2fb9b2bb5189c38dbcd4827c1d1af06009156feba0a43821e8a616d5bee3.
//
// Solidity: event MsgTransferCollateralOwnership(address valAddr, address prevOwner, address newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgTransferCollateralOwnership(log types.Log) (*ConsensusValidatorEntrypointMsgTransferCollateralOwnership, error) {
	event := new(ConsensusValidatorEntrypointMsgTransferCollateralOwnership)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgTransferCollateralOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgUnjailIterator is returned from FilterMsgUnjail and is used to iterate over the raw logs and unpacked data for MsgUnjail events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUnjailIterator struct {
	Event *ConsensusValidatorEntrypointMsgUnjail // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgUnjail)
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
		it.Event = new(ConsensusValidatorEntrypointMsgUnjail)
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
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgUnjailIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgUnjail represents a MsgUnjail event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUnjail struct {
	ValAddr common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMsgUnjail is a free log retrieval operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgUnjail(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgUnjailIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgUnjail")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgUnjailIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgUnjail", logs: logs, sub: sub}, nil
}

// WatchMsgUnjail is a free log subscription operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgUnjail(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgUnjail) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgUnjail")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgUnjail)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUnjail", log); err != nil {
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

// ParseMsgUnjail is a log parse operation binding the contract event 0xc2c03e4fbe86816915120cf64410100921065083a63a94c0a510d190bb79a893.
//
// Solidity: event MsgUnjail(address valAddr)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgUnjail(log types.Log) (*ConsensusValidatorEntrypointMsgUnjail, error) {
	event := new(ConsensusValidatorEntrypointMsgUnjail)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUnjail", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator is returned from FilterMsgUpdateExtraVotingPower and is used to iterate over the raw logs and unpacked data for MsgUpdateExtraVotingPower events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator struct {
	Event *ConsensusValidatorEntrypointMsgUpdateExtraVotingPower // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
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
		it.Event = new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
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
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgUpdateExtraVotingPower represents a MsgUpdateExtraVotingPower event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgUpdateExtraVotingPower struct {
	ValAddr             common.Address
	ExtraVotingPowerWei *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterMsgUpdateExtraVotingPower is a free log retrieval operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerWei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgUpdateExtraVotingPower(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgUpdateExtraVotingPower")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgUpdateExtraVotingPowerIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgUpdateExtraVotingPower", logs: logs, sub: sub}, nil
}

// WatchMsgUpdateExtraVotingPower is a free log subscription operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerWei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgUpdateExtraVotingPower(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgUpdateExtraVotingPower) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgUpdateExtraVotingPower")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUpdateExtraVotingPower", log); err != nil {
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

// ParseMsgUpdateExtraVotingPower is a log parse operation binding the contract event 0x38a463da3fce9952cb48e4bbe1b35ba1c3cfed5aec3043a4a0461a3168191d9c.
//
// Solidity: event MsgUpdateExtraVotingPower(address valAddr, uint256 extraVotingPowerWei)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgUpdateExtraVotingPower(log types.Log) (*ConsensusValidatorEntrypointMsgUpdateExtraVotingPower, error) {
	event := new(ConsensusValidatorEntrypointMsgUpdateExtraVotingPower)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgUpdateExtraVotingPower", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointMsgWithdrawCollateralIterator is returned from FilterMsgWithdrawCollateral and is used to iterate over the raw logs and unpacked data for MsgWithdrawCollateral events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgWithdrawCollateralIterator struct {
	Event *ConsensusValidatorEntrypointMsgWithdrawCollateral // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
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
		it.Event = new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
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
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointMsgWithdrawCollateralIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointMsgWithdrawCollateral represents a MsgWithdrawCollateral event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointMsgWithdrawCollateral struct {
	ValAddr         common.Address
	CollateralOwner common.Address
	Receiver        common.Address
	AmountGwei      *big.Int
	MaturesAt       *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMsgWithdrawCollateral is a free log retrieval operation binding the contract event 0xbdab795133b31a218b5ed9311c4ba8d56ca979b26a525191ab0456ba45f94c70.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amountGwei, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterMsgWithdrawCollateral(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointMsgWithdrawCollateralIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "MsgWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointMsgWithdrawCollateralIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "MsgWithdrawCollateral", logs: logs, sub: sub}, nil
}

// WatchMsgWithdrawCollateral is a free log subscription operation binding the contract event 0xbdab795133b31a218b5ed9311c4ba8d56ca979b26a525191ab0456ba45f94c70.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amountGwei, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchMsgWithdrawCollateral(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointMsgWithdrawCollateral) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "MsgWithdrawCollateral")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgWithdrawCollateral", log); err != nil {
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

// ParseMsgWithdrawCollateral is a log parse operation binding the contract event 0xbdab795133b31a218b5ed9311c4ba8d56ca979b26a525191ab0456ba45f94c70.
//
// Solidity: event MsgWithdrawCollateral(address valAddr, address collateralOwner, address receiver, uint256 amountGwei, uint48 maturesAt)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseMsgWithdrawCollateral(log types.Log) (*ConsensusValidatorEntrypointMsgWithdrawCollateral, error) {
	event := new(ConsensusValidatorEntrypointMsgWithdrawCollateral)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "MsgWithdrawCollateral", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferStartedIterator struct {
	Event *ConsensusValidatorEntrypointOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointOwnershipTransferStarted)
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
		it.Event = new(ConsensusValidatorEntrypointOwnershipTransferStarted)
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
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusValidatorEntrypointOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointOwnershipTransferStartedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointOwnershipTransferStarted)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseOwnershipTransferStarted(log types.Log) (*ConsensusValidatorEntrypointOwnershipTransferStarted, error) {
	event := new(ConsensusValidatorEntrypointOwnershipTransferStarted)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferredIterator struct {
	Event *ConsensusValidatorEntrypointOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointOwnershipTransferred)
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
		it.Event = new(ConsensusValidatorEntrypointOwnershipTransferred)
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
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointOwnershipTransferred represents a OwnershipTransferred event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusValidatorEntrypointOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointOwnershipTransferredIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointOwnershipTransferred)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseOwnershipTransferred(log types.Log) (*ConsensusValidatorEntrypointOwnershipTransferred, error) {
	event := new(ConsensusValidatorEntrypointOwnershipTransferred)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointPermittedCallerSetIterator is returned from FilterPermittedCallerSet and is used to iterate over the raw logs and unpacked data for PermittedCallerSet events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointPermittedCallerSetIterator struct {
	Event *ConsensusValidatorEntrypointPermittedCallerSet // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointPermittedCallerSet)
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
		it.Event = new(ConsensusValidatorEntrypointPermittedCallerSet)
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
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointPermittedCallerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointPermittedCallerSet represents a PermittedCallerSet event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointPermittedCallerSet struct {
	Caller      common.Address
	IsPermitted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPermittedCallerSet is a free log retrieval operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterPermittedCallerSet(opts *bind.FilterOpts) (*ConsensusValidatorEntrypointPermittedCallerSetIterator, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointPermittedCallerSetIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "PermittedCallerSet", logs: logs, sub: sub}, nil
}

// WatchPermittedCallerSet is a free log subscription operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchPermittedCallerSet(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointPermittedCallerSet) (event.Subscription, error) {

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointPermittedCallerSet)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
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

// ParsePermittedCallerSet is a log parse operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParsePermittedCallerSet(log types.Log) (*ConsensusValidatorEntrypointPermittedCallerSet, error) {
	event := new(ConsensusValidatorEntrypointPermittedCallerSet)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusValidatorEntrypointUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointUpgradedIterator struct {
	Event *ConsensusValidatorEntrypointUpgraded // Event containing the contract specifics and raw log

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
func (it *ConsensusValidatorEntrypointUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusValidatorEntrypointUpgraded)
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
		it.Event = new(ConsensusValidatorEntrypointUpgraded)
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
func (it *ConsensusValidatorEntrypointUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusValidatorEntrypointUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusValidatorEntrypointUpgraded represents a Upgraded event raised by the ConsensusValidatorEntrypoint contract.
type ConsensusValidatorEntrypointUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ConsensusValidatorEntrypointUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusValidatorEntrypointUpgradedIterator{contract: _ConsensusValidatorEntrypoint.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ConsensusValidatorEntrypointUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusValidatorEntrypoint.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusValidatorEntrypointUpgraded)
				if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusValidatorEntrypoint *ConsensusValidatorEntrypointFilterer) ParseUpgraded(log types.Log) (*ConsensusValidatorEntrypointUpgraded, error) {
	event := new(ConsensusValidatorEntrypointUpgraded)
	if err := _ConsensusValidatorEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
