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

// ConsensusGovernanceEntrypointMetaData contains all meta data concerning the ConsensusGovernanceEntrypoint contract.
var ConsensusGovernanceEntrypointMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"execute\",\"inputs\":[{\"name\":\"messages\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPermittedCaller\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MsgExecute\",\"inputs\":[{\"name\":\"messages\",\"type\":\"string[]\",\"indexed\":false,\"internalType\":\"string[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PermittedCallerSet\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"isPermitted\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[]}]",
}

// ConsensusGovernanceEntrypointABI is the input ABI used to generate the binding from.
// Deprecated: Use ConsensusGovernanceEntrypointMetaData.ABI instead.
var ConsensusGovernanceEntrypointABI = ConsensusGovernanceEntrypointMetaData.ABI

// ConsensusGovernanceEntrypoint is an auto generated Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypoint struct {
	ConsensusGovernanceEntrypointCaller     // Read-only binding to the contract
	ConsensusGovernanceEntrypointTransactor // Write-only binding to the contract
	ConsensusGovernanceEntrypointFilterer   // Log filterer for contract events
}

// ConsensusGovernanceEntrypointCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusGovernanceEntrypointTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusGovernanceEntrypointFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusGovernanceEntrypointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusGovernanceEntrypointSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusGovernanceEntrypointSession struct {
	Contract     *ConsensusGovernanceEntrypoint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                  // Call options to use throughout this session
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ConsensusGovernanceEntrypointCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusGovernanceEntrypointCallerSession struct {
	Contract *ConsensusGovernanceEntrypointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                        // Call options to use throughout this session
}

// ConsensusGovernanceEntrypointTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusGovernanceEntrypointTransactorSession struct {
	Contract     *ConsensusGovernanceEntrypointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                        // Transaction auth options to use throughout this session
}

// ConsensusGovernanceEntrypointRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypointRaw struct {
	Contract *ConsensusGovernanceEntrypoint // Generic contract binding to access the raw methods on
}

// ConsensusGovernanceEntrypointCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypointCallerRaw struct {
	Contract *ConsensusGovernanceEntrypointCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusGovernanceEntrypointTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusGovernanceEntrypointTransactorRaw struct {
	Contract *ConsensusGovernanceEntrypointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensusGovernanceEntrypoint creates a new instance of ConsensusGovernanceEntrypoint, bound to a specific deployed contract.
func NewConsensusGovernanceEntrypoint(address common.Address, backend bind.ContractBackend) (*ConsensusGovernanceEntrypoint, error) {
	contract, err := bindConsensusGovernanceEntrypoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypoint{ConsensusGovernanceEntrypointCaller: ConsensusGovernanceEntrypointCaller{contract: contract}, ConsensusGovernanceEntrypointTransactor: ConsensusGovernanceEntrypointTransactor{contract: contract}, ConsensusGovernanceEntrypointFilterer: ConsensusGovernanceEntrypointFilterer{contract: contract}}, nil
}

// NewConsensusGovernanceEntrypointCaller creates a new read-only instance of ConsensusGovernanceEntrypoint, bound to a specific deployed contract.
func NewConsensusGovernanceEntrypointCaller(address common.Address, caller bind.ContractCaller) (*ConsensusGovernanceEntrypointCaller, error) {
	contract, err := bindConsensusGovernanceEntrypoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointCaller{contract: contract}, nil
}

// NewConsensusGovernanceEntrypointTransactor creates a new write-only instance of ConsensusGovernanceEntrypoint, bound to a specific deployed contract.
func NewConsensusGovernanceEntrypointTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusGovernanceEntrypointTransactor, error) {
	contract, err := bindConsensusGovernanceEntrypoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointTransactor{contract: contract}, nil
}

// NewConsensusGovernanceEntrypointFilterer creates a new log filterer instance of ConsensusGovernanceEntrypoint, bound to a specific deployed contract.
func NewConsensusGovernanceEntrypointFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusGovernanceEntrypointFilterer, error) {
	contract, err := bindConsensusGovernanceEntrypoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointFilterer{contract: contract}, nil
}

// bindConsensusGovernanceEntrypoint binds a generic wrapper to an already deployed contract.
func bindConsensusGovernanceEntrypoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConsensusGovernanceEntrypointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusGovernanceEntrypoint.Contract.ConsensusGovernanceEntrypointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.ConsensusGovernanceEntrypointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.ConsensusGovernanceEntrypointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusGovernanceEntrypoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConsensusGovernanceEntrypoint.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusGovernanceEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ConsensusGovernanceEntrypoint.Contract.UPGRADEINTERFACEVERSION(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCaller) IsPermittedCaller(opts *bind.CallOpts, caller common.Address) (bool, error) {
	var out []interface{}
	err := _ConsensusGovernanceEntrypoint.contract.Call(opts, &out, "isPermittedCaller", caller)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusGovernanceEntrypoint.Contract.IsPermittedCaller(&_ConsensusGovernanceEntrypoint.CallOpts, caller)
}

// IsPermittedCaller is a free data retrieval call binding the contract method 0x7fcc389f.
//
// Solidity: function isPermittedCaller(address caller) view returns(bool)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerSession) IsPermittedCaller(caller common.Address) (bool, error) {
	return _ConsensusGovernanceEntrypoint.Contract.IsPermittedCaller(&_ConsensusGovernanceEntrypoint.CallOpts, caller)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusGovernanceEntrypoint.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) Owner() (common.Address, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Owner(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerSession) Owner() (common.Address, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Owner(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConsensusGovernanceEntrypoint.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) PendingOwner() (common.Address, error) {
	return _ConsensusGovernanceEntrypoint.Contract.PendingOwner(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerSession) PendingOwner() (common.Address, error) {
	return _ConsensusGovernanceEntrypoint.Contract.PendingOwner(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ConsensusGovernanceEntrypoint.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusGovernanceEntrypoint.Contract.ProxiableUUID(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ConsensusGovernanceEntrypoint.Contract.ProxiableUUID(&_ConsensusGovernanceEntrypoint.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.AcceptOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.AcceptOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// Execute is a paid mutator transaction binding the contract method 0x1c9a67d5.
//
// Solidity: function execute(string[] messages) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) Execute(opts *bind.TransactOpts, messages []string) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "execute", messages)
}

// Execute is a paid mutator transaction binding the contract method 0x1c9a67d5.
//
// Solidity: function execute(string[] messages) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) Execute(messages []string) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Execute(&_ConsensusGovernanceEntrypoint.TransactOpts, messages)
}

// Execute is a paid mutator transaction binding the contract method 0x1c9a67d5.
//
// Solidity: function execute(string[] messages) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) Execute(messages []string) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Execute(&_ConsensusGovernanceEntrypoint.TransactOpts, messages)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) Initialize(opts *bind.TransactOpts, owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "initialize", owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Initialize(&_ConsensusGovernanceEntrypoint.TransactOpts, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Initialize(&_ConsensusGovernanceEntrypoint.TransactOpts, owner_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.RenounceOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.RenounceOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) SetPermittedCaller(opts *bind.TransactOpts, caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "setPermittedCaller", caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.SetPermittedCaller(&_ConsensusGovernanceEntrypoint.TransactOpts, caller, isPermitted)
}

// SetPermittedCaller is a paid mutator transaction binding the contract method 0x1727b6f3.
//
// Solidity: function setPermittedCaller(address caller, bool isPermitted) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) SetPermittedCaller(caller common.Address, isPermitted bool) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.SetPermittedCaller(&_ConsensusGovernanceEntrypoint.TransactOpts, caller, isPermitted)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.TransferOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.TransferOwnership(&_ConsensusGovernanceEntrypoint.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.UpgradeToAndCall(&_ConsensusGovernanceEntrypoint.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.UpgradeToAndCall(&_ConsensusGovernanceEntrypoint.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Fallback(&_ConsensusGovernanceEntrypoint.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Fallback(&_ConsensusGovernanceEntrypoint.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointSession) Receive() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Receive(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointTransactorSession) Receive() (*types.Transaction, error) {
	return _ConsensusGovernanceEntrypoint.Contract.Receive(&_ConsensusGovernanceEntrypoint.TransactOpts)
}

// ConsensusGovernanceEntrypointInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointInitializedIterator struct {
	Event *ConsensusGovernanceEntrypointInitialized // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointInitialized)
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
		it.Event = new(ConsensusGovernanceEntrypointInitialized)
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
func (it *ConsensusGovernanceEntrypointInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointInitialized represents a Initialized event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterInitialized(opts *bind.FilterOpts) (*ConsensusGovernanceEntrypointInitializedIterator, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointInitializedIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointInitialized) (event.Subscription, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointInitialized)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParseInitialized(log types.Log) (*ConsensusGovernanceEntrypointInitialized, error) {
	event := new(ConsensusGovernanceEntrypointInitialized)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusGovernanceEntrypointMsgExecuteIterator is returned from FilterMsgExecute and is used to iterate over the raw logs and unpacked data for MsgExecute events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointMsgExecuteIterator struct {
	Event *ConsensusGovernanceEntrypointMsgExecute // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointMsgExecuteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointMsgExecute)
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
		it.Event = new(ConsensusGovernanceEntrypointMsgExecute)
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
func (it *ConsensusGovernanceEntrypointMsgExecuteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointMsgExecuteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointMsgExecute represents a MsgExecute event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointMsgExecute struct {
	Messages []string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMsgExecute is a free log retrieval operation binding the contract event 0xc363b25c95578dcdd12780ca814bf3b5212f34826d54c2b380a442a4414369f0.
//
// Solidity: event MsgExecute(string[] messages)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterMsgExecute(opts *bind.FilterOpts) (*ConsensusGovernanceEntrypointMsgExecuteIterator, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "MsgExecute")
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointMsgExecuteIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "MsgExecute", logs: logs, sub: sub}, nil
}

// WatchMsgExecute is a free log subscription operation binding the contract event 0xc363b25c95578dcdd12780ca814bf3b5212f34826d54c2b380a442a4414369f0.
//
// Solidity: event MsgExecute(string[] messages)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchMsgExecute(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointMsgExecute) (event.Subscription, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "MsgExecute")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointMsgExecute)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "MsgExecute", log); err != nil {
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

// ParseMsgExecute is a log parse operation binding the contract event 0xc363b25c95578dcdd12780ca814bf3b5212f34826d54c2b380a442a4414369f0.
//
// Solidity: event MsgExecute(string[] messages)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParseMsgExecute(log types.Log) (*ConsensusGovernanceEntrypointMsgExecute, error) {
	event := new(ConsensusGovernanceEntrypointMsgExecute)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "MsgExecute", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusGovernanceEntrypointOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointOwnershipTransferStartedIterator struct {
	Event *ConsensusGovernanceEntrypointOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointOwnershipTransferStarted)
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
		it.Event = new(ConsensusGovernanceEntrypointOwnershipTransferStarted)
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
func (it *ConsensusGovernanceEntrypointOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusGovernanceEntrypointOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointOwnershipTransferStartedIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointOwnershipTransferStarted)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParseOwnershipTransferStarted(log types.Log) (*ConsensusGovernanceEntrypointOwnershipTransferStarted, error) {
	event := new(ConsensusGovernanceEntrypointOwnershipTransferStarted)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusGovernanceEntrypointOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointOwnershipTransferredIterator struct {
	Event *ConsensusGovernanceEntrypointOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointOwnershipTransferred)
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
		it.Event = new(ConsensusGovernanceEntrypointOwnershipTransferred)
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
func (it *ConsensusGovernanceEntrypointOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointOwnershipTransferred represents a OwnershipTransferred event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConsensusGovernanceEntrypointOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointOwnershipTransferredIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointOwnershipTransferred)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParseOwnershipTransferred(log types.Log) (*ConsensusGovernanceEntrypointOwnershipTransferred, error) {
	event := new(ConsensusGovernanceEntrypointOwnershipTransferred)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusGovernanceEntrypointPermittedCallerSetIterator is returned from FilterPermittedCallerSet and is used to iterate over the raw logs and unpacked data for PermittedCallerSet events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointPermittedCallerSetIterator struct {
	Event *ConsensusGovernanceEntrypointPermittedCallerSet // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointPermittedCallerSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointPermittedCallerSet)
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
		it.Event = new(ConsensusGovernanceEntrypointPermittedCallerSet)
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
func (it *ConsensusGovernanceEntrypointPermittedCallerSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointPermittedCallerSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointPermittedCallerSet represents a PermittedCallerSet event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointPermittedCallerSet struct {
	Caller      common.Address
	IsPermitted bool
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPermittedCallerSet is a free log retrieval operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterPermittedCallerSet(opts *bind.FilterOpts) (*ConsensusGovernanceEntrypointPermittedCallerSetIterator, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointPermittedCallerSetIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "PermittedCallerSet", logs: logs, sub: sub}, nil
}

// WatchPermittedCallerSet is a free log subscription operation binding the contract event 0x58b0246a79531a991271a8abe150f2c09805dec04338c87eca66ed423855d4c5.
//
// Solidity: event PermittedCallerSet(address caller, bool isPermitted)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchPermittedCallerSet(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointPermittedCallerSet) (event.Subscription, error) {

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "PermittedCallerSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointPermittedCallerSet)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
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
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParsePermittedCallerSet(log types.Log) (*ConsensusGovernanceEntrypointPermittedCallerSet, error) {
	event := new(ConsensusGovernanceEntrypointPermittedCallerSet)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "PermittedCallerSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConsensusGovernanceEntrypointUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointUpgradedIterator struct {
	Event *ConsensusGovernanceEntrypointUpgraded // Event containing the contract specifics and raw log

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
func (it *ConsensusGovernanceEntrypointUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConsensusGovernanceEntrypointUpgraded)
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
		it.Event = new(ConsensusGovernanceEntrypointUpgraded)
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
func (it *ConsensusGovernanceEntrypointUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConsensusGovernanceEntrypointUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConsensusGovernanceEntrypointUpgraded represents a Upgraded event raised by the ConsensusGovernanceEntrypoint contract.
type ConsensusGovernanceEntrypointUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ConsensusGovernanceEntrypointUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ConsensusGovernanceEntrypointUpgradedIterator{contract: _ConsensusGovernanceEntrypoint.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ConsensusGovernanceEntrypointUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ConsensusGovernanceEntrypoint.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConsensusGovernanceEntrypointUpgraded)
				if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ConsensusGovernanceEntrypoint *ConsensusGovernanceEntrypointFilterer) ParseUpgraded(log types.Log) (*ConsensusGovernanceEntrypointUpgraded, error) {
	event := new(ConsensusGovernanceEntrypointUpgraded)
	if err := _ConsensusGovernanceEntrypoint.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
