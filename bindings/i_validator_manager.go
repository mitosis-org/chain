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

// IValidatorManagerCreateValidatorRequest is an auto generated low-level Go binding around an user-defined struct.
type IValidatorManagerCreateValidatorRequest struct {
	Operator            common.Address
	RewardManager       common.Address
	WithdrawalRecipient common.Address
	CommissionRate      *big.Int
	Metadata            []byte
}

// IValidatorManagerGlobalValidatorConfigResponse is an auto generated low-level Go binding around an user-defined struct.
type IValidatorManagerGlobalValidatorConfigResponse struct {
	InitialValidatorDeposit          *big.Int
	CollateralWithdrawalDelaySeconds *big.Int
	MinimumCommissionRate            *big.Int
	CommissionRateUpdateDelayEpoch   *big.Int
}

// IValidatorManagerSetGlobalValidatorConfigRequest is an auto generated low-level Go binding around an user-defined struct.
type IValidatorManagerSetGlobalValidatorConfigRequest struct {
	InitialValidatorDeposit          *big.Int
	CollateralWithdrawalDelaySeconds *big.Int
	MinimumCommissionRate            *big.Int
	CommissionRateUpdateDelayEpoch   *big.Int
}

// IValidatorManagerUpdateRewardConfigRequest is an auto generated low-level Go binding around an user-defined struct.
type IValidatorManagerUpdateRewardConfigRequest struct {
	CommissionRate *big.Int
}

// IValidatorManagerValidatorInfoResponse is an auto generated low-level Go binding around an user-defined struct.
type IValidatorManagerValidatorInfoResponse struct {
	ValAddr                          common.Address
	PubKey                           []byte
	Operator                         common.Address
	RewardManager                    common.Address
	WithdrawalRecipient              common.Address
	CommissionRate                   *big.Int
	PendingCommissionRate            *big.Int
	PendingCommissionRateUpdateEpoch *big.Int
	Metadata                         []byte
}

// IValidatorManagerMetaData contains all meta data concerning the IValidatorManager contract.
var IValidatorManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"MAX_COMMISSION_RATE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createValidator\",\"inputs\":[{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"request\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.CreateValidatorRequest\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"depositCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"entrypoint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIConsensusValidatorEntrypoint\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"epochFeeder\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIEpochFeeder\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalValidatorConfig\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.GlobalValidatorConfigResponse\",\"components\":[{\"name\":\"initialValidatorDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"collateralWithdrawalDelaySeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimumCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commissionRateUpdateDelayEpoch\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setFee\",\"inputs\":[{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGlobalValidatorConfig\",\"inputs\":[{\"name\":\"request\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.SetGlobalValidatorConfigRequest\",\"components\":[{\"name\":\"initialValidatorDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"collateralWithdrawalDelaySeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimumCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commissionRateUpdateDelayEpoch\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unjailValidator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"updateMetadata\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateOperator\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRewardConfig\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"request\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.UpdateRewardConfigRequest\",\"components\":[{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRewardManager\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateWithdrawalRecipient\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"validatorAt\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorInfo\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.ValidatorInfoResponse\",\"components\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingCommissionRateUpdateEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorInfoAt\",\"inputs\":[{\"name\":\"epoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIValidatorManager.ValidatorInfoResponse\",\"components\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"pendingCommissionRateUpdateEpoch\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"validatorPubKeyToAddress\",\"inputs\":[{\"name\":\"pubKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"withdrawCollateral\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"CollateralDeposited\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CollateralWithdrawn\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EntrypointUpdated\",\"inputs\":[{\"name\":\"entrypoint\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIConsensusValidatorEntrypoint\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochFeederUpdated\",\"inputs\":[{\"name\":\"epochFeeder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"contractIEpochFeeder\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeePaid\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSet\",\"inputs\":[{\"name\":\"previousFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GlobalValidatorConfigUpdated\",\"inputs\":[{\"name\":\"request\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIValidatorManager.SetGlobalValidatorConfigRequest\",\"components\":[{\"name\":\"initialValidatorDeposit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"collateralWithdrawalDelaySeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minimumCommissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"commissionRateUpdateDelayEpoch\",\"type\":\"uint96\",\"internalType\":\"uint96\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OperatorUpdated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardConfigUpdated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"request\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIValidatorManager.UpdateRewardConfigRequest\",\"components\":[{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardManagerUpdated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorCreated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pubKey\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"},{\"name\":\"initialDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"request\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structIValidatorManager.CreateValidatorRequest\",\"components\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardManager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"withdrawalRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnjailed\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WithdrawalRecipientUpdated\",\"inputs\":[{\"name\":\"valAddr\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"IValidatorManager__InsufficientFee\",\"inputs\":[]}]",
}

// IValidatorManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use IValidatorManagerMetaData.ABI instead.
var IValidatorManagerABI = IValidatorManagerMetaData.ABI

// IValidatorManager is an auto generated Go binding around an Ethereum contract.
type IValidatorManager struct {
	IValidatorManagerCaller     // Read-only binding to the contract
	IValidatorManagerTransactor // Write-only binding to the contract
	IValidatorManagerFilterer   // Log filterer for contract events
}

// IValidatorManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type IValidatorManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IValidatorManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IValidatorManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IValidatorManagerSession struct {
	Contract     *IValidatorManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IValidatorManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IValidatorManagerCallerSession struct {
	Contract *IValidatorManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IValidatorManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IValidatorManagerTransactorSession struct {
	Contract     *IValidatorManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IValidatorManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type IValidatorManagerRaw struct {
	Contract *IValidatorManager // Generic contract binding to access the raw methods on
}

// IValidatorManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IValidatorManagerCallerRaw struct {
	Contract *IValidatorManagerCaller // Generic read-only contract binding to access the raw methods on
}

// IValidatorManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IValidatorManagerTransactorRaw struct {
	Contract *IValidatorManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIValidatorManager creates a new instance of IValidatorManager, bound to a specific deployed contract.
func NewIValidatorManager(address common.Address, backend bind.ContractBackend) (*IValidatorManager, error) {
	contract, err := bindIValidatorManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IValidatorManager{IValidatorManagerCaller: IValidatorManagerCaller{contract: contract}, IValidatorManagerTransactor: IValidatorManagerTransactor{contract: contract}, IValidatorManagerFilterer: IValidatorManagerFilterer{contract: contract}}, nil
}

// NewIValidatorManagerCaller creates a new read-only instance of IValidatorManager, bound to a specific deployed contract.
func NewIValidatorManagerCaller(address common.Address, caller bind.ContractCaller) (*IValidatorManagerCaller, error) {
	contract, err := bindIValidatorManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerCaller{contract: contract}, nil
}

// NewIValidatorManagerTransactor creates a new write-only instance of IValidatorManager, bound to a specific deployed contract.
func NewIValidatorManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*IValidatorManagerTransactor, error) {
	contract, err := bindIValidatorManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerTransactor{contract: contract}, nil
}

// NewIValidatorManagerFilterer creates a new log filterer instance of IValidatorManager, bound to a specific deployed contract.
func NewIValidatorManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*IValidatorManagerFilterer, error) {
	contract, err := bindIValidatorManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerFilterer{contract: contract}, nil
}

// bindIValidatorManager binds a generic wrapper to an already deployed contract.
func bindIValidatorManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IValidatorManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IValidatorManager *IValidatorManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IValidatorManager.Contract.IValidatorManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IValidatorManager *IValidatorManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IValidatorManager.Contract.IValidatorManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IValidatorManager *IValidatorManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IValidatorManager.Contract.IValidatorManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IValidatorManager *IValidatorManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IValidatorManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IValidatorManager *IValidatorManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IValidatorManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IValidatorManager *IValidatorManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IValidatorManager.Contract.contract.Transact(opts, method, params...)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCaller) MAXCOMMISSIONRATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "MAX_COMMISSION_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_IValidatorManager *IValidatorManagerSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _IValidatorManager.Contract.MAXCOMMISSIONRATE(&_IValidatorManager.CallOpts)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCallerSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _IValidatorManager.Contract.MAXCOMMISSIONRATE(&_IValidatorManager.CallOpts)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns(address)
func (_IValidatorManager *IValidatorManagerCaller) Entrypoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "entrypoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns(address)
func (_IValidatorManager *IValidatorManagerSession) Entrypoint() (common.Address, error) {
	return _IValidatorManager.Contract.Entrypoint(&_IValidatorManager.CallOpts)
}

// Entrypoint is a free data retrieval call binding the contract method 0xa65d69d4.
//
// Solidity: function entrypoint() view returns(address)
func (_IValidatorManager *IValidatorManagerCallerSession) Entrypoint() (common.Address, error) {
	return _IValidatorManager.Contract.Entrypoint(&_IValidatorManager.CallOpts)
}

// EpochFeeder is a free data retrieval call binding the contract method 0x8af4a41f.
//
// Solidity: function epochFeeder() view returns(address)
func (_IValidatorManager *IValidatorManagerCaller) EpochFeeder(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "epochFeeder")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EpochFeeder is a free data retrieval call binding the contract method 0x8af4a41f.
//
// Solidity: function epochFeeder() view returns(address)
func (_IValidatorManager *IValidatorManagerSession) EpochFeeder() (common.Address, error) {
	return _IValidatorManager.Contract.EpochFeeder(&_IValidatorManager.CallOpts)
}

// EpochFeeder is a free data retrieval call binding the contract method 0x8af4a41f.
//
// Solidity: function epochFeeder() view returns(address)
func (_IValidatorManager *IValidatorManagerCallerSession) EpochFeeder() (common.Address, error) {
	return _IValidatorManager.Contract.EpochFeeder(&_IValidatorManager.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IValidatorManager *IValidatorManagerSession) Fee() (*big.Int, error) {
	return _IValidatorManager.Contract.Fee(&_IValidatorManager.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCallerSession) Fee() (*big.Int, error) {
	return _IValidatorManager.Contract.Fee(&_IValidatorManager.CallOpts)
}

// GlobalValidatorConfig is a free data retrieval call binding the contract method 0x1169d671.
//
// Solidity: function globalValidatorConfig() view returns((uint256,uint256,uint256,uint96))
func (_IValidatorManager *IValidatorManagerCaller) GlobalValidatorConfig(opts *bind.CallOpts) (IValidatorManagerGlobalValidatorConfigResponse, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "globalValidatorConfig")

	if err != nil {
		return *new(IValidatorManagerGlobalValidatorConfigResponse), err
	}

	out0 := *abi.ConvertType(out[0], new(IValidatorManagerGlobalValidatorConfigResponse)).(*IValidatorManagerGlobalValidatorConfigResponse)

	return out0, err

}

// GlobalValidatorConfig is a free data retrieval call binding the contract method 0x1169d671.
//
// Solidity: function globalValidatorConfig() view returns((uint256,uint256,uint256,uint96))
func (_IValidatorManager *IValidatorManagerSession) GlobalValidatorConfig() (IValidatorManagerGlobalValidatorConfigResponse, error) {
	return _IValidatorManager.Contract.GlobalValidatorConfig(&_IValidatorManager.CallOpts)
}

// GlobalValidatorConfig is a free data retrieval call binding the contract method 0x1169d671.
//
// Solidity: function globalValidatorConfig() view returns((uint256,uint256,uint256,uint96))
func (_IValidatorManager *IValidatorManagerCallerSession) GlobalValidatorConfig() (IValidatorManagerGlobalValidatorConfigResponse, error) {
	return _IValidatorManager.Contract.GlobalValidatorConfig(&_IValidatorManager.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address valAddr) view returns(bool)
func (_IValidatorManager *IValidatorManagerCaller) IsValidator(opts *bind.CallOpts, valAddr common.Address) (bool, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "isValidator", valAddr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address valAddr) view returns(bool)
func (_IValidatorManager *IValidatorManagerSession) IsValidator(valAddr common.Address) (bool, error) {
	return _IValidatorManager.Contract.IsValidator(&_IValidatorManager.CallOpts, valAddr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address valAddr) view returns(bool)
func (_IValidatorManager *IValidatorManagerCallerSession) IsValidator(valAddr common.Address) (bool, error) {
	return _IValidatorManager.Contract.IsValidator(&_IValidatorManager.CallOpts, valAddr)
}

// ValidatorAt is a free data retrieval call binding the contract method 0x32e0aa1f.
//
// Solidity: function validatorAt(uint256 index) view returns(address)
func (_IValidatorManager *IValidatorManagerCaller) ValidatorAt(opts *bind.CallOpts, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "validatorAt", index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValidatorAt is a free data retrieval call binding the contract method 0x32e0aa1f.
//
// Solidity: function validatorAt(uint256 index) view returns(address)
func (_IValidatorManager *IValidatorManagerSession) ValidatorAt(index *big.Int) (common.Address, error) {
	return _IValidatorManager.Contract.ValidatorAt(&_IValidatorManager.CallOpts, index)
}

// ValidatorAt is a free data retrieval call binding the contract method 0x32e0aa1f.
//
// Solidity: function validatorAt(uint256 index) view returns(address)
func (_IValidatorManager *IValidatorManagerCallerSession) ValidatorAt(index *big.Int) (common.Address, error) {
	return _IValidatorManager.Contract.ValidatorAt(&_IValidatorManager.CallOpts, index)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "validatorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_IValidatorManager *IValidatorManagerSession) ValidatorCount() (*big.Int, error) {
	return _IValidatorManager.Contract.ValidatorCount(&_IValidatorManager.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_IValidatorManager *IValidatorManagerCallerSession) ValidatorCount() (*big.Int, error) {
	return _IValidatorManager.Contract.ValidatorCount(&_IValidatorManager.CallOpts)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo(address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerCaller) ValidatorInfo(opts *bind.CallOpts, valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "validatorInfo", valAddr)

	if err != nil {
		return *new(IValidatorManagerValidatorInfoResponse), err
	}

	out0 := *abi.ConvertType(out[0], new(IValidatorManagerValidatorInfoResponse)).(*IValidatorManagerValidatorInfoResponse)

	return out0, err

}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo(address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerSession) ValidatorInfo(valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	return _IValidatorManager.Contract.ValidatorInfo(&_IValidatorManager.CallOpts, valAddr)
}

// ValidatorInfo is a free data retrieval call binding the contract method 0x4f1811dd.
//
// Solidity: function validatorInfo(address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerCallerSession) ValidatorInfo(valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	return _IValidatorManager.Contract.ValidatorInfo(&_IValidatorManager.CallOpts, valAddr)
}

// ValidatorInfoAt is a free data retrieval call binding the contract method 0xf91f7e80.
//
// Solidity: function validatorInfoAt(uint256 epoch, address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerCaller) ValidatorInfoAt(opts *bind.CallOpts, epoch *big.Int, valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "validatorInfoAt", epoch, valAddr)

	if err != nil {
		return *new(IValidatorManagerValidatorInfoResponse), err
	}

	out0 := *abi.ConvertType(out[0], new(IValidatorManagerValidatorInfoResponse)).(*IValidatorManagerValidatorInfoResponse)

	return out0, err

}

// ValidatorInfoAt is a free data retrieval call binding the contract method 0xf91f7e80.
//
// Solidity: function validatorInfoAt(uint256 epoch, address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerSession) ValidatorInfoAt(epoch *big.Int, valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	return _IValidatorManager.Contract.ValidatorInfoAt(&_IValidatorManager.CallOpts, epoch, valAddr)
}

// ValidatorInfoAt is a free data retrieval call binding the contract method 0xf91f7e80.
//
// Solidity: function validatorInfoAt(uint256 epoch, address valAddr) view returns((address,bytes,address,address,address,uint256,uint256,uint256,bytes))
func (_IValidatorManager *IValidatorManagerCallerSession) ValidatorInfoAt(epoch *big.Int, valAddr common.Address) (IValidatorManagerValidatorInfoResponse, error) {
	return _IValidatorManager.Contract.ValidatorInfoAt(&_IValidatorManager.CallOpts, epoch, valAddr)
}

// ValidatorPubKeyToAddress is a free data retrieval call binding the contract method 0x97708c5d.
//
// Solidity: function validatorPubKeyToAddress(bytes pubKey) pure returns(address)
func (_IValidatorManager *IValidatorManagerCaller) ValidatorPubKeyToAddress(opts *bind.CallOpts, pubKey []byte) (common.Address, error) {
	var out []interface{}
	err := _IValidatorManager.contract.Call(opts, &out, "validatorPubKeyToAddress", pubKey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValidatorPubKeyToAddress is a free data retrieval call binding the contract method 0x97708c5d.
//
// Solidity: function validatorPubKeyToAddress(bytes pubKey) pure returns(address)
func (_IValidatorManager *IValidatorManagerSession) ValidatorPubKeyToAddress(pubKey []byte) (common.Address, error) {
	return _IValidatorManager.Contract.ValidatorPubKeyToAddress(&_IValidatorManager.CallOpts, pubKey)
}

// ValidatorPubKeyToAddress is a free data retrieval call binding the contract method 0x97708c5d.
//
// Solidity: function validatorPubKeyToAddress(bytes pubKey) pure returns(address)
func (_IValidatorManager *IValidatorManagerCallerSession) ValidatorPubKeyToAddress(pubKey []byte) (common.Address, error) {
	return _IValidatorManager.Contract.ValidatorPubKeyToAddress(&_IValidatorManager.CallOpts, pubKey)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xc883a0b7.
//
// Solidity: function createValidator(bytes pubKey, (address,address,address,uint256,bytes) request) payable returns()
func (_IValidatorManager *IValidatorManagerTransactor) CreateValidator(opts *bind.TransactOpts, pubKey []byte, request IValidatorManagerCreateValidatorRequest) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "createValidator", pubKey, request)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xc883a0b7.
//
// Solidity: function createValidator(bytes pubKey, (address,address,address,uint256,bytes) request) payable returns()
func (_IValidatorManager *IValidatorManagerSession) CreateValidator(pubKey []byte, request IValidatorManagerCreateValidatorRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.CreateValidator(&_IValidatorManager.TransactOpts, pubKey, request)
}

// CreateValidator is a paid mutator transaction binding the contract method 0xc883a0b7.
//
// Solidity: function createValidator(bytes pubKey, (address,address,address,uint256,bytes) request) payable returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) CreateValidator(pubKey []byte, request IValidatorManagerCreateValidatorRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.CreateValidator(&_IValidatorManager.TransactOpts, pubKey, request)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerTransactor) DepositCollateral(opts *bind.TransactOpts, valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "depositCollateral", valAddr)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerSession) DepositCollateral(valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.DepositCollateral(&_IValidatorManager.TransactOpts, valAddr)
}

// DepositCollateral is a paid mutator transaction binding the contract method 0x97d475d7.
//
// Solidity: function depositCollateral(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) DepositCollateral(valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.DepositCollateral(&_IValidatorManager.TransactOpts, valAddr)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_IValidatorManager *IValidatorManagerTransactor) SetFee(opts *bind.TransactOpts, fee *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "setFee", fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_IValidatorManager *IValidatorManagerSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.Contract.SetFee(&_IValidatorManager.TransactOpts, fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.Contract.SetFee(&_IValidatorManager.TransactOpts, fee)
}

// SetGlobalValidatorConfig is a paid mutator transaction binding the contract method 0x0201907f.
//
// Solidity: function setGlobalValidatorConfig((uint256,uint256,uint256,uint96) request) returns()
func (_IValidatorManager *IValidatorManagerTransactor) SetGlobalValidatorConfig(opts *bind.TransactOpts, request IValidatorManagerSetGlobalValidatorConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "setGlobalValidatorConfig", request)
}

// SetGlobalValidatorConfig is a paid mutator transaction binding the contract method 0x0201907f.
//
// Solidity: function setGlobalValidatorConfig((uint256,uint256,uint256,uint96) request) returns()
func (_IValidatorManager *IValidatorManagerSession) SetGlobalValidatorConfig(request IValidatorManagerSetGlobalValidatorConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.SetGlobalValidatorConfig(&_IValidatorManager.TransactOpts, request)
}

// SetGlobalValidatorConfig is a paid mutator transaction binding the contract method 0x0201907f.
//
// Solidity: function setGlobalValidatorConfig((uint256,uint256,uint256,uint96) request) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) SetGlobalValidatorConfig(request IValidatorManagerSetGlobalValidatorConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.SetGlobalValidatorConfig(&_IValidatorManager.TransactOpts, request)
}

// UnjailValidator is a paid mutator transaction binding the contract method 0x7cafdd79.
//
// Solidity: function unjailValidator(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerTransactor) UnjailValidator(opts *bind.TransactOpts, valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "unjailValidator", valAddr)
}

// UnjailValidator is a paid mutator transaction binding the contract method 0x7cafdd79.
//
// Solidity: function unjailValidator(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerSession) UnjailValidator(valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UnjailValidator(&_IValidatorManager.TransactOpts, valAddr)
}

// UnjailValidator is a paid mutator transaction binding the contract method 0x7cafdd79.
//
// Solidity: function unjailValidator(address valAddr) payable returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UnjailValidator(valAddr common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UnjailValidator(&_IValidatorManager.TransactOpts, valAddr)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x0e5268aa.
//
// Solidity: function updateMetadata(address valAddr, bytes metadata) returns()
func (_IValidatorManager *IValidatorManagerTransactor) UpdateMetadata(opts *bind.TransactOpts, valAddr common.Address, metadata []byte) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "updateMetadata", valAddr, metadata)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x0e5268aa.
//
// Solidity: function updateMetadata(address valAddr, bytes metadata) returns()
func (_IValidatorManager *IValidatorManagerSession) UpdateMetadata(valAddr common.Address, metadata []byte) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateMetadata(&_IValidatorManager.TransactOpts, valAddr, metadata)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0x0e5268aa.
//
// Solidity: function updateMetadata(address valAddr, bytes metadata) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UpdateMetadata(valAddr common.Address, metadata []byte) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateMetadata(&_IValidatorManager.TransactOpts, valAddr, metadata)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x8cd2d73e.
//
// Solidity: function updateOperator(address valAddr, address operator) returns()
func (_IValidatorManager *IValidatorManagerTransactor) UpdateOperator(opts *bind.TransactOpts, valAddr common.Address, operator common.Address) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "updateOperator", valAddr, operator)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x8cd2d73e.
//
// Solidity: function updateOperator(address valAddr, address operator) returns()
func (_IValidatorManager *IValidatorManagerSession) UpdateOperator(valAddr common.Address, operator common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateOperator(&_IValidatorManager.TransactOpts, valAddr, operator)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x8cd2d73e.
//
// Solidity: function updateOperator(address valAddr, address operator) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UpdateOperator(valAddr common.Address, operator common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateOperator(&_IValidatorManager.TransactOpts, valAddr, operator)
}

// UpdateRewardConfig is a paid mutator transaction binding the contract method 0x24cf2c54.
//
// Solidity: function updateRewardConfig(address valAddr, (uint256) request) returns()
func (_IValidatorManager *IValidatorManagerTransactor) UpdateRewardConfig(opts *bind.TransactOpts, valAddr common.Address, request IValidatorManagerUpdateRewardConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "updateRewardConfig", valAddr, request)
}

// UpdateRewardConfig is a paid mutator transaction binding the contract method 0x24cf2c54.
//
// Solidity: function updateRewardConfig(address valAddr, (uint256) request) returns()
func (_IValidatorManager *IValidatorManagerSession) UpdateRewardConfig(valAddr common.Address, request IValidatorManagerUpdateRewardConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateRewardConfig(&_IValidatorManager.TransactOpts, valAddr, request)
}

// UpdateRewardConfig is a paid mutator transaction binding the contract method 0x24cf2c54.
//
// Solidity: function updateRewardConfig(address valAddr, (uint256) request) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UpdateRewardConfig(valAddr common.Address, request IValidatorManagerUpdateRewardConfigRequest) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateRewardConfig(&_IValidatorManager.TransactOpts, valAddr, request)
}

// UpdateRewardManager is a paid mutator transaction binding the contract method 0x0b84c463.
//
// Solidity: function updateRewardManager(address valAddr, address rewardManager) returns()
func (_IValidatorManager *IValidatorManagerTransactor) UpdateRewardManager(opts *bind.TransactOpts, valAddr common.Address, rewardManager common.Address) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "updateRewardManager", valAddr, rewardManager)
}

// UpdateRewardManager is a paid mutator transaction binding the contract method 0x0b84c463.
//
// Solidity: function updateRewardManager(address valAddr, address rewardManager) returns()
func (_IValidatorManager *IValidatorManagerSession) UpdateRewardManager(valAddr common.Address, rewardManager common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateRewardManager(&_IValidatorManager.TransactOpts, valAddr, rewardManager)
}

// UpdateRewardManager is a paid mutator transaction binding the contract method 0x0b84c463.
//
// Solidity: function updateRewardManager(address valAddr, address rewardManager) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UpdateRewardManager(valAddr common.Address, rewardManager common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateRewardManager(&_IValidatorManager.TransactOpts, valAddr, rewardManager)
}

// UpdateWithdrawalRecipient is a paid mutator transaction binding the contract method 0x1c44e7ed.
//
// Solidity: function updateWithdrawalRecipient(address valAddr, address withdrawalRecipient) returns()
func (_IValidatorManager *IValidatorManagerTransactor) UpdateWithdrawalRecipient(opts *bind.TransactOpts, valAddr common.Address, withdrawalRecipient common.Address) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "updateWithdrawalRecipient", valAddr, withdrawalRecipient)
}

// UpdateWithdrawalRecipient is a paid mutator transaction binding the contract method 0x1c44e7ed.
//
// Solidity: function updateWithdrawalRecipient(address valAddr, address withdrawalRecipient) returns()
func (_IValidatorManager *IValidatorManagerSession) UpdateWithdrawalRecipient(valAddr common.Address, withdrawalRecipient common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateWithdrawalRecipient(&_IValidatorManager.TransactOpts, valAddr, withdrawalRecipient)
}

// UpdateWithdrawalRecipient is a paid mutator transaction binding the contract method 0x1c44e7ed.
//
// Solidity: function updateWithdrawalRecipient(address valAddr, address withdrawalRecipient) returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) UpdateWithdrawalRecipient(valAddr common.Address, withdrawalRecipient common.Address) (*types.Transaction, error) {
	return _IValidatorManager.Contract.UpdateWithdrawalRecipient(&_IValidatorManager.TransactOpts, valAddr, withdrawalRecipient)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount) payable returns()
func (_IValidatorManager *IValidatorManagerTransactor) WithdrawCollateral(opts *bind.TransactOpts, valAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.contract.Transact(opts, "withdrawCollateral", valAddr, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount) payable returns()
func (_IValidatorManager *IValidatorManagerSession) WithdrawCollateral(valAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.Contract.WithdrawCollateral(&_IValidatorManager.TransactOpts, valAddr, amount)
}

// WithdrawCollateral is a paid mutator transaction binding the contract method 0x350c35e9.
//
// Solidity: function withdrawCollateral(address valAddr, uint256 amount) payable returns()
func (_IValidatorManager *IValidatorManagerTransactorSession) WithdrawCollateral(valAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IValidatorManager.Contract.WithdrawCollateral(&_IValidatorManager.TransactOpts, valAddr, amount)
}

// IValidatorManagerCollateralDepositedIterator is returned from FilterCollateralDeposited and is used to iterate over the raw logs and unpacked data for CollateralDeposited events raised by the IValidatorManager contract.
type IValidatorManagerCollateralDepositedIterator struct {
	Event *IValidatorManagerCollateralDeposited // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerCollateralDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerCollateralDeposited)
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
		it.Event = new(IValidatorManagerCollateralDeposited)
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
func (it *IValidatorManagerCollateralDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerCollateralDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerCollateralDeposited represents a CollateralDeposited event raised by the IValidatorManager contract.
type IValidatorManagerCollateralDeposited struct {
	ValAddr   common.Address
	Depositor common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCollateralDeposited is a free log retrieval operation binding the contract event 0xf1c0dd7e9b98bbff859029005ef89b127af049cd18df1a8d79f0b7e019911e56.
//
// Solidity: event CollateralDeposited(address indexed valAddr, address indexed depositor, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) FilterCollateralDeposited(opts *bind.FilterOpts, valAddr []common.Address, depositor []common.Address) (*IValidatorManagerCollateralDepositedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "CollateralDeposited", valAddrRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerCollateralDepositedIterator{contract: _IValidatorManager.contract, event: "CollateralDeposited", logs: logs, sub: sub}, nil
}

// WatchCollateralDeposited is a free log subscription operation binding the contract event 0xf1c0dd7e9b98bbff859029005ef89b127af049cd18df1a8d79f0b7e019911e56.
//
// Solidity: event CollateralDeposited(address indexed valAddr, address indexed depositor, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) WatchCollateralDeposited(opts *bind.WatchOpts, sink chan<- *IValidatorManagerCollateralDeposited, valAddr []common.Address, depositor []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "CollateralDeposited", valAddrRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerCollateralDeposited)
				if err := _IValidatorManager.contract.UnpackLog(event, "CollateralDeposited", log); err != nil {
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

// ParseCollateralDeposited is a log parse operation binding the contract event 0xf1c0dd7e9b98bbff859029005ef89b127af049cd18df1a8d79f0b7e019911e56.
//
// Solidity: event CollateralDeposited(address indexed valAddr, address indexed depositor, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) ParseCollateralDeposited(log types.Log) (*IValidatorManagerCollateralDeposited, error) {
	event := new(IValidatorManagerCollateralDeposited)
	if err := _IValidatorManager.contract.UnpackLog(event, "CollateralDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerCollateralWithdrawnIterator is returned from FilterCollateralWithdrawn and is used to iterate over the raw logs and unpacked data for CollateralWithdrawn events raised by the IValidatorManager contract.
type IValidatorManagerCollateralWithdrawnIterator struct {
	Event *IValidatorManagerCollateralWithdrawn // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerCollateralWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerCollateralWithdrawn)
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
		it.Event = new(IValidatorManagerCollateralWithdrawn)
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
func (it *IValidatorManagerCollateralWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerCollateralWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerCollateralWithdrawn represents a CollateralWithdrawn event raised by the IValidatorManager contract.
type IValidatorManagerCollateralWithdrawn struct {
	ValAddr   common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCollateralWithdrawn is a free log retrieval operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed valAddr, address indexed recipient, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) FilterCollateralWithdrawn(opts *bind.FilterOpts, valAddr []common.Address, recipient []common.Address) (*IValidatorManagerCollateralWithdrawnIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "CollateralWithdrawn", valAddrRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerCollateralWithdrawnIterator{contract: _IValidatorManager.contract, event: "CollateralWithdrawn", logs: logs, sub: sub}, nil
}

// WatchCollateralWithdrawn is a free log subscription operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed valAddr, address indexed recipient, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) WatchCollateralWithdrawn(opts *bind.WatchOpts, sink chan<- *IValidatorManagerCollateralWithdrawn, valAddr []common.Address, recipient []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "CollateralWithdrawn", valAddrRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerCollateralWithdrawn)
				if err := _IValidatorManager.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
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

// ParseCollateralWithdrawn is a log parse operation binding the contract event 0x45892a46e6cef329bb642da6d69846d324db43d19008edc141ed82382eda1bee.
//
// Solidity: event CollateralWithdrawn(address indexed valAddr, address indexed recipient, uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) ParseCollateralWithdrawn(log types.Log) (*IValidatorManagerCollateralWithdrawn, error) {
	event := new(IValidatorManagerCollateralWithdrawn)
	if err := _IValidatorManager.contract.UnpackLog(event, "CollateralWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerEntrypointUpdatedIterator is returned from FilterEntrypointUpdated and is used to iterate over the raw logs and unpacked data for EntrypointUpdated events raised by the IValidatorManager contract.
type IValidatorManagerEntrypointUpdatedIterator struct {
	Event *IValidatorManagerEntrypointUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerEntrypointUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerEntrypointUpdated)
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
		it.Event = new(IValidatorManagerEntrypointUpdated)
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
func (it *IValidatorManagerEntrypointUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerEntrypointUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerEntrypointUpdated represents a EntrypointUpdated event raised by the IValidatorManager contract.
type IValidatorManagerEntrypointUpdated struct {
	Entrypoint common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEntrypointUpdated is a free log retrieval operation binding the contract event 0xad540c78e43fa817dbc0ff910b73ae9a8b5f35575a75aeeaeda7bd3b2def22f2.
//
// Solidity: event EntrypointUpdated(address indexed entrypoint)
func (_IValidatorManager *IValidatorManagerFilterer) FilterEntrypointUpdated(opts *bind.FilterOpts, entrypoint []common.Address) (*IValidatorManagerEntrypointUpdatedIterator, error) {

	var entrypointRule []interface{}
	for _, entrypointItem := range entrypoint {
		entrypointRule = append(entrypointRule, entrypointItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "EntrypointUpdated", entrypointRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerEntrypointUpdatedIterator{contract: _IValidatorManager.contract, event: "EntrypointUpdated", logs: logs, sub: sub}, nil
}

// WatchEntrypointUpdated is a free log subscription operation binding the contract event 0xad540c78e43fa817dbc0ff910b73ae9a8b5f35575a75aeeaeda7bd3b2def22f2.
//
// Solidity: event EntrypointUpdated(address indexed entrypoint)
func (_IValidatorManager *IValidatorManagerFilterer) WatchEntrypointUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerEntrypointUpdated, entrypoint []common.Address) (event.Subscription, error) {

	var entrypointRule []interface{}
	for _, entrypointItem := range entrypoint {
		entrypointRule = append(entrypointRule, entrypointItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "EntrypointUpdated", entrypointRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerEntrypointUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "EntrypointUpdated", log); err != nil {
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

// ParseEntrypointUpdated is a log parse operation binding the contract event 0xad540c78e43fa817dbc0ff910b73ae9a8b5f35575a75aeeaeda7bd3b2def22f2.
//
// Solidity: event EntrypointUpdated(address indexed entrypoint)
func (_IValidatorManager *IValidatorManagerFilterer) ParseEntrypointUpdated(log types.Log) (*IValidatorManagerEntrypointUpdated, error) {
	event := new(IValidatorManagerEntrypointUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "EntrypointUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerEpochFeederUpdatedIterator is returned from FilterEpochFeederUpdated and is used to iterate over the raw logs and unpacked data for EpochFeederUpdated events raised by the IValidatorManager contract.
type IValidatorManagerEpochFeederUpdatedIterator struct {
	Event *IValidatorManagerEpochFeederUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerEpochFeederUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerEpochFeederUpdated)
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
		it.Event = new(IValidatorManagerEpochFeederUpdated)
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
func (it *IValidatorManagerEpochFeederUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerEpochFeederUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerEpochFeederUpdated represents a EpochFeederUpdated event raised by the IValidatorManager contract.
type IValidatorManagerEpochFeederUpdated struct {
	EpochFeeder common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEpochFeederUpdated is a free log retrieval operation binding the contract event 0xfcfcf8d234ae7853aa4ac2ed78680b5f5f48a996115b800237981d0d9354c870.
//
// Solidity: event EpochFeederUpdated(address indexed epochFeeder)
func (_IValidatorManager *IValidatorManagerFilterer) FilterEpochFeederUpdated(opts *bind.FilterOpts, epochFeeder []common.Address) (*IValidatorManagerEpochFeederUpdatedIterator, error) {

	var epochFeederRule []interface{}
	for _, epochFeederItem := range epochFeeder {
		epochFeederRule = append(epochFeederRule, epochFeederItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "EpochFeederUpdated", epochFeederRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerEpochFeederUpdatedIterator{contract: _IValidatorManager.contract, event: "EpochFeederUpdated", logs: logs, sub: sub}, nil
}

// WatchEpochFeederUpdated is a free log subscription operation binding the contract event 0xfcfcf8d234ae7853aa4ac2ed78680b5f5f48a996115b800237981d0d9354c870.
//
// Solidity: event EpochFeederUpdated(address indexed epochFeeder)
func (_IValidatorManager *IValidatorManagerFilterer) WatchEpochFeederUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerEpochFeederUpdated, epochFeeder []common.Address) (event.Subscription, error) {

	var epochFeederRule []interface{}
	for _, epochFeederItem := range epochFeeder {
		epochFeederRule = append(epochFeederRule, epochFeederItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "EpochFeederUpdated", epochFeederRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerEpochFeederUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "EpochFeederUpdated", log); err != nil {
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

// ParseEpochFeederUpdated is a log parse operation binding the contract event 0xfcfcf8d234ae7853aa4ac2ed78680b5f5f48a996115b800237981d0d9354c870.
//
// Solidity: event EpochFeederUpdated(address indexed epochFeeder)
func (_IValidatorManager *IValidatorManagerFilterer) ParseEpochFeederUpdated(log types.Log) (*IValidatorManagerEpochFeederUpdated, error) {
	event := new(IValidatorManagerEpochFeederUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "EpochFeederUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerFeePaidIterator is returned from FilterFeePaid and is used to iterate over the raw logs and unpacked data for FeePaid events raised by the IValidatorManager contract.
type IValidatorManagerFeePaidIterator struct {
	Event *IValidatorManagerFeePaid // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerFeePaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerFeePaid)
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
		it.Event = new(IValidatorManagerFeePaid)
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
func (it *IValidatorManagerFeePaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerFeePaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerFeePaid represents a FeePaid event raised by the IValidatorManager contract.
type IValidatorManagerFeePaid struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeePaid is a free log retrieval operation binding the contract event 0x69e27f80547602d16208b028c44d20f25956e1fb7d0f51d62aa02f392426f371.
//
// Solidity: event FeePaid(uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) FilterFeePaid(opts *bind.FilterOpts) (*IValidatorManagerFeePaidIterator, error) {

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerFeePaidIterator{contract: _IValidatorManager.contract, event: "FeePaid", logs: logs, sub: sub}, nil
}

// WatchFeePaid is a free log subscription operation binding the contract event 0x69e27f80547602d16208b028c44d20f25956e1fb7d0f51d62aa02f392426f371.
//
// Solidity: event FeePaid(uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) WatchFeePaid(opts *bind.WatchOpts, sink chan<- *IValidatorManagerFeePaid) (event.Subscription, error) {

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "FeePaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerFeePaid)
				if err := _IValidatorManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
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

// ParseFeePaid is a log parse operation binding the contract event 0x69e27f80547602d16208b028c44d20f25956e1fb7d0f51d62aa02f392426f371.
//
// Solidity: event FeePaid(uint256 amount)
func (_IValidatorManager *IValidatorManagerFilterer) ParseFeePaid(log types.Log) (*IValidatorManagerFeePaid, error) {
	event := new(IValidatorManagerFeePaid)
	if err := _IValidatorManager.contract.UnpackLog(event, "FeePaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerFeeSetIterator is returned from FilterFeeSet and is used to iterate over the raw logs and unpacked data for FeeSet events raised by the IValidatorManager contract.
type IValidatorManagerFeeSetIterator struct {
	Event *IValidatorManagerFeeSet // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerFeeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerFeeSet)
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
		it.Event = new(IValidatorManagerFeeSet)
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
func (it *IValidatorManagerFeeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerFeeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerFeeSet represents a FeeSet event raised by the IValidatorManager contract.
type IValidatorManagerFeeSet struct {
	PreviousFee *big.Int
	NewFee      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFeeSet is a free log retrieval operation binding the contract event 0x74dbbbe280ef27b79a8a0c449d5ae2ba7a31849103241d0f98df70bbc9d03e37.
//
// Solidity: event FeeSet(uint256 previousFee, uint256 newFee)
func (_IValidatorManager *IValidatorManagerFilterer) FilterFeeSet(opts *bind.FilterOpts) (*IValidatorManagerFeeSetIterator, error) {

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerFeeSetIterator{contract: _IValidatorManager.contract, event: "FeeSet", logs: logs, sub: sub}, nil
}

// WatchFeeSet is a free log subscription operation binding the contract event 0x74dbbbe280ef27b79a8a0c449d5ae2ba7a31849103241d0f98df70bbc9d03e37.
//
// Solidity: event FeeSet(uint256 previousFee, uint256 newFee)
func (_IValidatorManager *IValidatorManagerFilterer) WatchFeeSet(opts *bind.WatchOpts, sink chan<- *IValidatorManagerFeeSet) (event.Subscription, error) {

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "FeeSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerFeeSet)
				if err := _IValidatorManager.contract.UnpackLog(event, "FeeSet", log); err != nil {
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

// ParseFeeSet is a log parse operation binding the contract event 0x74dbbbe280ef27b79a8a0c449d5ae2ba7a31849103241d0f98df70bbc9d03e37.
//
// Solidity: event FeeSet(uint256 previousFee, uint256 newFee)
func (_IValidatorManager *IValidatorManagerFilterer) ParseFeeSet(log types.Log) (*IValidatorManagerFeeSet, error) {
	event := new(IValidatorManagerFeeSet)
	if err := _IValidatorManager.contract.UnpackLog(event, "FeeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerGlobalValidatorConfigUpdatedIterator is returned from FilterGlobalValidatorConfigUpdated and is used to iterate over the raw logs and unpacked data for GlobalValidatorConfigUpdated events raised by the IValidatorManager contract.
type IValidatorManagerGlobalValidatorConfigUpdatedIterator struct {
	Event *IValidatorManagerGlobalValidatorConfigUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerGlobalValidatorConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerGlobalValidatorConfigUpdated)
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
		it.Event = new(IValidatorManagerGlobalValidatorConfigUpdated)
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
func (it *IValidatorManagerGlobalValidatorConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerGlobalValidatorConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerGlobalValidatorConfigUpdated represents a GlobalValidatorConfigUpdated event raised by the IValidatorManager contract.
type IValidatorManagerGlobalValidatorConfigUpdated struct {
	Request IValidatorManagerSetGlobalValidatorConfigRequest
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGlobalValidatorConfigUpdated is a free log retrieval operation binding the contract event 0x00b98761a7fbeabb6757dbe28a2d4458db1235509b31dca76cb9092472b4f873.
//
// Solidity: event GlobalValidatorConfigUpdated((uint256,uint256,uint256,uint96) request)
func (_IValidatorManager *IValidatorManagerFilterer) FilterGlobalValidatorConfigUpdated(opts *bind.FilterOpts) (*IValidatorManagerGlobalValidatorConfigUpdatedIterator, error) {

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "GlobalValidatorConfigUpdated")
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerGlobalValidatorConfigUpdatedIterator{contract: _IValidatorManager.contract, event: "GlobalValidatorConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchGlobalValidatorConfigUpdated is a free log subscription operation binding the contract event 0x00b98761a7fbeabb6757dbe28a2d4458db1235509b31dca76cb9092472b4f873.
//
// Solidity: event GlobalValidatorConfigUpdated((uint256,uint256,uint256,uint96) request)
func (_IValidatorManager *IValidatorManagerFilterer) WatchGlobalValidatorConfigUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerGlobalValidatorConfigUpdated) (event.Subscription, error) {

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "GlobalValidatorConfigUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerGlobalValidatorConfigUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "GlobalValidatorConfigUpdated", log); err != nil {
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

// ParseGlobalValidatorConfigUpdated is a log parse operation binding the contract event 0x00b98761a7fbeabb6757dbe28a2d4458db1235509b31dca76cb9092472b4f873.
//
// Solidity: event GlobalValidatorConfigUpdated((uint256,uint256,uint256,uint96) request)
func (_IValidatorManager *IValidatorManagerFilterer) ParseGlobalValidatorConfigUpdated(log types.Log) (*IValidatorManagerGlobalValidatorConfigUpdated, error) {
	event := new(IValidatorManagerGlobalValidatorConfigUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "GlobalValidatorConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerMetadataUpdatedIterator is returned from FilterMetadataUpdated and is used to iterate over the raw logs and unpacked data for MetadataUpdated events raised by the IValidatorManager contract.
type IValidatorManagerMetadataUpdatedIterator struct {
	Event *IValidatorManagerMetadataUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerMetadataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerMetadataUpdated)
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
		it.Event = new(IValidatorManagerMetadataUpdated)
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
func (it *IValidatorManagerMetadataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerMetadataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerMetadataUpdated represents a MetadataUpdated event raised by the IValidatorManager contract.
type IValidatorManagerMetadataUpdated struct {
	ValAddr  common.Address
	Operator common.Address
	Metadata []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdated is a free log retrieval operation binding the contract event 0xef5bd3dac1b82a93b66f1f3ee38a4a67da5b34fc1d6f38c73aac8919388e91cf.
//
// Solidity: event MetadataUpdated(address indexed valAddr, address indexed operator, bytes metadata)
func (_IValidatorManager *IValidatorManagerFilterer) FilterMetadataUpdated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address) (*IValidatorManagerMetadataUpdatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "MetadataUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerMetadataUpdatedIterator{contract: _IValidatorManager.contract, event: "MetadataUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdated is a free log subscription operation binding the contract event 0xef5bd3dac1b82a93b66f1f3ee38a4a67da5b34fc1d6f38c73aac8919388e91cf.
//
// Solidity: event MetadataUpdated(address indexed valAddr, address indexed operator, bytes metadata)
func (_IValidatorManager *IValidatorManagerFilterer) WatchMetadataUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerMetadataUpdated, valAddr []common.Address, operator []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "MetadataUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerMetadataUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
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

// ParseMetadataUpdated is a log parse operation binding the contract event 0xef5bd3dac1b82a93b66f1f3ee38a4a67da5b34fc1d6f38c73aac8919388e91cf.
//
// Solidity: event MetadataUpdated(address indexed valAddr, address indexed operator, bytes metadata)
func (_IValidatorManager *IValidatorManagerFilterer) ParseMetadataUpdated(log types.Log) (*IValidatorManagerMetadataUpdated, error) {
	event := new(IValidatorManagerMetadataUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerOperatorUpdatedIterator is returned from FilterOperatorUpdated and is used to iterate over the raw logs and unpacked data for OperatorUpdated events raised by the IValidatorManager contract.
type IValidatorManagerOperatorUpdatedIterator struct {
	Event *IValidatorManagerOperatorUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerOperatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerOperatorUpdated)
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
		it.Event = new(IValidatorManagerOperatorUpdated)
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
func (it *IValidatorManagerOperatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerOperatorUpdated represents a OperatorUpdated event raised by the IValidatorManager contract.
type IValidatorManagerOperatorUpdated struct {
	ValAddr  common.Address
	Operator common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorUpdated is a free log retrieval operation binding the contract event 0xfbe5b6cbafb274f445d7fed869dc77a838d8243a22c460de156560e8857cad03.
//
// Solidity: event OperatorUpdated(address indexed valAddr, address indexed operator)
func (_IValidatorManager *IValidatorManagerFilterer) FilterOperatorUpdated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address) (*IValidatorManagerOperatorUpdatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "OperatorUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerOperatorUpdatedIterator{contract: _IValidatorManager.contract, event: "OperatorUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorUpdated is a free log subscription operation binding the contract event 0xfbe5b6cbafb274f445d7fed869dc77a838d8243a22c460de156560e8857cad03.
//
// Solidity: event OperatorUpdated(address indexed valAddr, address indexed operator)
func (_IValidatorManager *IValidatorManagerFilterer) WatchOperatorUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerOperatorUpdated, valAddr []common.Address, operator []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "OperatorUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerOperatorUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
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

// ParseOperatorUpdated is a log parse operation binding the contract event 0xfbe5b6cbafb274f445d7fed869dc77a838d8243a22c460de156560e8857cad03.
//
// Solidity: event OperatorUpdated(address indexed valAddr, address indexed operator)
func (_IValidatorManager *IValidatorManagerFilterer) ParseOperatorUpdated(log types.Log) (*IValidatorManagerOperatorUpdated, error) {
	event := new(IValidatorManagerOperatorUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerRewardConfigUpdatedIterator is returned from FilterRewardConfigUpdated and is used to iterate over the raw logs and unpacked data for RewardConfigUpdated events raised by the IValidatorManager contract.
type IValidatorManagerRewardConfigUpdatedIterator struct {
	Event *IValidatorManagerRewardConfigUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerRewardConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerRewardConfigUpdated)
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
		it.Event = new(IValidatorManagerRewardConfigUpdated)
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
func (it *IValidatorManagerRewardConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerRewardConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerRewardConfigUpdated represents a RewardConfigUpdated event raised by the IValidatorManager contract.
type IValidatorManagerRewardConfigUpdated struct {
	ValAddr  common.Address
	Operator common.Address
	Request  IValidatorManagerUpdateRewardConfigRequest
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRewardConfigUpdated is a free log retrieval operation binding the contract event 0x58fdfec86a0389f1dc0c43064b5b8c63e3eebc9bab3e87031517e35b9b930030.
//
// Solidity: event RewardConfigUpdated(address indexed valAddr, address indexed operator, (uint256) request)
func (_IValidatorManager *IValidatorManagerFilterer) FilterRewardConfigUpdated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address) (*IValidatorManagerRewardConfigUpdatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "RewardConfigUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerRewardConfigUpdatedIterator{contract: _IValidatorManager.contract, event: "RewardConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardConfigUpdated is a free log subscription operation binding the contract event 0x58fdfec86a0389f1dc0c43064b5b8c63e3eebc9bab3e87031517e35b9b930030.
//
// Solidity: event RewardConfigUpdated(address indexed valAddr, address indexed operator, (uint256) request)
func (_IValidatorManager *IValidatorManagerFilterer) WatchRewardConfigUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerRewardConfigUpdated, valAddr []common.Address, operator []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "RewardConfigUpdated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerRewardConfigUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "RewardConfigUpdated", log); err != nil {
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

// ParseRewardConfigUpdated is a log parse operation binding the contract event 0x58fdfec86a0389f1dc0c43064b5b8c63e3eebc9bab3e87031517e35b9b930030.
//
// Solidity: event RewardConfigUpdated(address indexed valAddr, address indexed operator, (uint256) request)
func (_IValidatorManager *IValidatorManagerFilterer) ParseRewardConfigUpdated(log types.Log) (*IValidatorManagerRewardConfigUpdated, error) {
	event := new(IValidatorManagerRewardConfigUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "RewardConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerRewardManagerUpdatedIterator is returned from FilterRewardManagerUpdated and is used to iterate over the raw logs and unpacked data for RewardManagerUpdated events raised by the IValidatorManager contract.
type IValidatorManagerRewardManagerUpdatedIterator struct {
	Event *IValidatorManagerRewardManagerUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerRewardManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerRewardManagerUpdated)
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
		it.Event = new(IValidatorManagerRewardManagerUpdated)
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
func (it *IValidatorManagerRewardManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerRewardManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerRewardManagerUpdated represents a RewardManagerUpdated event raised by the IValidatorManager contract.
type IValidatorManagerRewardManagerUpdated struct {
	ValAddr       common.Address
	Operator      common.Address
	RewardManager common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRewardManagerUpdated is a free log retrieval operation binding the contract event 0xfde7a679f53c64d9a0ba22761e5c96c9aee5c6660a375715c9d5baf307f60bfd.
//
// Solidity: event RewardManagerUpdated(address indexed valAddr, address indexed operator, address indexed rewardManager)
func (_IValidatorManager *IValidatorManagerFilterer) FilterRewardManagerUpdated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address, rewardManager []common.Address) (*IValidatorManagerRewardManagerUpdatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var rewardManagerRule []interface{}
	for _, rewardManagerItem := range rewardManager {
		rewardManagerRule = append(rewardManagerRule, rewardManagerItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "RewardManagerUpdated", valAddrRule, operatorRule, rewardManagerRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerRewardManagerUpdatedIterator{contract: _IValidatorManager.contract, event: "RewardManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardManagerUpdated is a free log subscription operation binding the contract event 0xfde7a679f53c64d9a0ba22761e5c96c9aee5c6660a375715c9d5baf307f60bfd.
//
// Solidity: event RewardManagerUpdated(address indexed valAddr, address indexed operator, address indexed rewardManager)
func (_IValidatorManager *IValidatorManagerFilterer) WatchRewardManagerUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerRewardManagerUpdated, valAddr []common.Address, operator []common.Address, rewardManager []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var rewardManagerRule []interface{}
	for _, rewardManagerItem := range rewardManager {
		rewardManagerRule = append(rewardManagerRule, rewardManagerItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "RewardManagerUpdated", valAddrRule, operatorRule, rewardManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerRewardManagerUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "RewardManagerUpdated", log); err != nil {
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

// ParseRewardManagerUpdated is a log parse operation binding the contract event 0xfde7a679f53c64d9a0ba22761e5c96c9aee5c6660a375715c9d5baf307f60bfd.
//
// Solidity: event RewardManagerUpdated(address indexed valAddr, address indexed operator, address indexed rewardManager)
func (_IValidatorManager *IValidatorManagerFilterer) ParseRewardManagerUpdated(log types.Log) (*IValidatorManagerRewardManagerUpdated, error) {
	event := new(IValidatorManagerRewardManagerUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "RewardManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerValidatorCreatedIterator is returned from FilterValidatorCreated and is used to iterate over the raw logs and unpacked data for ValidatorCreated events raised by the IValidatorManager contract.
type IValidatorManagerValidatorCreatedIterator struct {
	Event *IValidatorManagerValidatorCreated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerValidatorCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerValidatorCreated)
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
		it.Event = new(IValidatorManagerValidatorCreated)
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
func (it *IValidatorManagerValidatorCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerValidatorCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerValidatorCreated represents a ValidatorCreated event raised by the IValidatorManager contract.
type IValidatorManagerValidatorCreated struct {
	ValAddr        common.Address
	Operator       common.Address
	PubKey         []byte
	InitialDeposit *big.Int
	Request        IValidatorManagerCreateValidatorRequest
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterValidatorCreated is a free log retrieval operation binding the contract event 0x2707028602f92fabeea22c9a85aacf912d326215752c36787a64b0eab8e8e8c8.
//
// Solidity: event ValidatorCreated(address indexed valAddr, address indexed operator, bytes pubKey, uint256 initialDeposit, (address,address,address,uint256,bytes) request)
func (_IValidatorManager *IValidatorManagerFilterer) FilterValidatorCreated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address) (*IValidatorManagerValidatorCreatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "ValidatorCreated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerValidatorCreatedIterator{contract: _IValidatorManager.contract, event: "ValidatorCreated", logs: logs, sub: sub}, nil
}

// WatchValidatorCreated is a free log subscription operation binding the contract event 0x2707028602f92fabeea22c9a85aacf912d326215752c36787a64b0eab8e8e8c8.
//
// Solidity: event ValidatorCreated(address indexed valAddr, address indexed operator, bytes pubKey, uint256 initialDeposit, (address,address,address,uint256,bytes) request)
func (_IValidatorManager *IValidatorManagerFilterer) WatchValidatorCreated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerValidatorCreated, valAddr []common.Address, operator []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "ValidatorCreated", valAddrRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerValidatorCreated)
				if err := _IValidatorManager.contract.UnpackLog(event, "ValidatorCreated", log); err != nil {
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

// ParseValidatorCreated is a log parse operation binding the contract event 0x2707028602f92fabeea22c9a85aacf912d326215752c36787a64b0eab8e8e8c8.
//
// Solidity: event ValidatorCreated(address indexed valAddr, address indexed operator, bytes pubKey, uint256 initialDeposit, (address,address,address,uint256,bytes) request)
func (_IValidatorManager *IValidatorManagerFilterer) ParseValidatorCreated(log types.Log) (*IValidatorManagerValidatorCreated, error) {
	event := new(IValidatorManagerValidatorCreated)
	if err := _IValidatorManager.contract.UnpackLog(event, "ValidatorCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerValidatorUnjailedIterator is returned from FilterValidatorUnjailed and is used to iterate over the raw logs and unpacked data for ValidatorUnjailed events raised by the IValidatorManager contract.
type IValidatorManagerValidatorUnjailedIterator struct {
	Event *IValidatorManagerValidatorUnjailed // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerValidatorUnjailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerValidatorUnjailed)
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
		it.Event = new(IValidatorManagerValidatorUnjailed)
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
func (it *IValidatorManagerValidatorUnjailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerValidatorUnjailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerValidatorUnjailed represents a ValidatorUnjailed event raised by the IValidatorManager contract.
type IValidatorManagerValidatorUnjailed struct {
	ValAddr common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorUnjailed is a free log retrieval operation binding the contract event 0x9390b453426557da5ebdc31f19a37753ca04addf656d32f35232211bb2af3f19.
//
// Solidity: event ValidatorUnjailed(address indexed valAddr)
func (_IValidatorManager *IValidatorManagerFilterer) FilterValidatorUnjailed(opts *bind.FilterOpts, valAddr []common.Address) (*IValidatorManagerValidatorUnjailedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "ValidatorUnjailed", valAddrRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerValidatorUnjailedIterator{contract: _IValidatorManager.contract, event: "ValidatorUnjailed", logs: logs, sub: sub}, nil
}

// WatchValidatorUnjailed is a free log subscription operation binding the contract event 0x9390b453426557da5ebdc31f19a37753ca04addf656d32f35232211bb2af3f19.
//
// Solidity: event ValidatorUnjailed(address indexed valAddr)
func (_IValidatorManager *IValidatorManagerFilterer) WatchValidatorUnjailed(opts *bind.WatchOpts, sink chan<- *IValidatorManagerValidatorUnjailed, valAddr []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "ValidatorUnjailed", valAddrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerValidatorUnjailed)
				if err := _IValidatorManager.contract.UnpackLog(event, "ValidatorUnjailed", log); err != nil {
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

// ParseValidatorUnjailed is a log parse operation binding the contract event 0x9390b453426557da5ebdc31f19a37753ca04addf656d32f35232211bb2af3f19.
//
// Solidity: event ValidatorUnjailed(address indexed valAddr)
func (_IValidatorManager *IValidatorManagerFilterer) ParseValidatorUnjailed(log types.Log) (*IValidatorManagerValidatorUnjailed, error) {
	event := new(IValidatorManagerValidatorUnjailed)
	if err := _IValidatorManager.contract.UnpackLog(event, "ValidatorUnjailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorManagerWithdrawalRecipientUpdatedIterator is returned from FilterWithdrawalRecipientUpdated and is used to iterate over the raw logs and unpacked data for WithdrawalRecipientUpdated events raised by the IValidatorManager contract.
type IValidatorManagerWithdrawalRecipientUpdatedIterator struct {
	Event *IValidatorManagerWithdrawalRecipientUpdated // Event containing the contract specifics and raw log

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
func (it *IValidatorManagerWithdrawalRecipientUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorManagerWithdrawalRecipientUpdated)
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
		it.Event = new(IValidatorManagerWithdrawalRecipientUpdated)
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
func (it *IValidatorManagerWithdrawalRecipientUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorManagerWithdrawalRecipientUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorManagerWithdrawalRecipientUpdated represents a WithdrawalRecipientUpdated event raised by the IValidatorManager contract.
type IValidatorManagerWithdrawalRecipientUpdated struct {
	ValAddr   common.Address
	Operator  common.Address
	Recipient common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalRecipientUpdated is a free log retrieval operation binding the contract event 0xc9e0c9fc260d3006a1443854782eaf0e82f7d2aaef5ac2dd10234843c97cccb0.
//
// Solidity: event WithdrawalRecipientUpdated(address indexed valAddr, address indexed operator, address indexed recipient)
func (_IValidatorManager *IValidatorManagerFilterer) FilterWithdrawalRecipientUpdated(opts *bind.FilterOpts, valAddr []common.Address, operator []common.Address, recipient []common.Address) (*IValidatorManagerWithdrawalRecipientUpdatedIterator, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IValidatorManager.contract.FilterLogs(opts, "WithdrawalRecipientUpdated", valAddrRule, operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &IValidatorManagerWithdrawalRecipientUpdatedIterator{contract: _IValidatorManager.contract, event: "WithdrawalRecipientUpdated", logs: logs, sub: sub}, nil
}

// WatchWithdrawalRecipientUpdated is a free log subscription operation binding the contract event 0xc9e0c9fc260d3006a1443854782eaf0e82f7d2aaef5ac2dd10234843c97cccb0.
//
// Solidity: event WithdrawalRecipientUpdated(address indexed valAddr, address indexed operator, address indexed recipient)
func (_IValidatorManager *IValidatorManagerFilterer) WatchWithdrawalRecipientUpdated(opts *bind.WatchOpts, sink chan<- *IValidatorManagerWithdrawalRecipientUpdated, valAddr []common.Address, operator []common.Address, recipient []common.Address) (event.Subscription, error) {

	var valAddrRule []interface{}
	for _, valAddrItem := range valAddr {
		valAddrRule = append(valAddrRule, valAddrItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IValidatorManager.contract.WatchLogs(opts, "WithdrawalRecipientUpdated", valAddrRule, operatorRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorManagerWithdrawalRecipientUpdated)
				if err := _IValidatorManager.contract.UnpackLog(event, "WithdrawalRecipientUpdated", log); err != nil {
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

// ParseWithdrawalRecipientUpdated is a log parse operation binding the contract event 0xc9e0c9fc260d3006a1443854782eaf0e82f7d2aaef5ac2dd10234843c97cccb0.
//
// Solidity: event WithdrawalRecipientUpdated(address indexed valAddr, address indexed operator, address indexed recipient)
func (_IValidatorManager *IValidatorManagerFilterer) ParseWithdrawalRecipientUpdated(log types.Log) (*IValidatorManagerWithdrawalRecipientUpdated, error) {
	event := new(IValidatorManagerWithdrawalRecipientUpdated)
	if err := _IValidatorManager.contract.UnpackLog(event, "WithdrawalRecipientUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
