// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ordinals

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

// OrdinalsMetaData contains all meta data concerning the Ordinals contract.
var OrdinalsMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_caller\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_inscription\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_parameterAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"parameterControl\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setCaller\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"coll\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"inscriptionId\",\"type\":\"string\"}],\"name\":\"setInscription\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// OrdinalsABI is the input ABI used to generate the binding from.
// Deprecated: Use OrdinalsMetaData.ABI instead.
var OrdinalsABI = OrdinalsMetaData.ABI

// Ordinals is an auto generated Go binding around an Ethereum contract.
type Ordinals struct {
	OrdinalsCaller     // Read-only binding to the contract
	OrdinalsTransactor // Write-only binding to the contract
	OrdinalsFilterer   // Log filterer for contract events
}

// OrdinalsCaller is an auto generated read-only Go binding around an Ethereum contract.
type OrdinalsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrdinalsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OrdinalsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrdinalsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OrdinalsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrdinalsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OrdinalsSession struct {
	Contract     *Ordinals         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OrdinalsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OrdinalsCallerSession struct {
	Contract *OrdinalsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// OrdinalsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OrdinalsTransactorSession struct {
	Contract     *OrdinalsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OrdinalsRaw is an auto generated low-level Go binding around an Ethereum contract.
type OrdinalsRaw struct {
	Contract *Ordinals // Generic contract binding to access the raw methods on
}

// OrdinalsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OrdinalsCallerRaw struct {
	Contract *OrdinalsCaller // Generic read-only contract binding to access the raw methods on
}

// OrdinalsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OrdinalsTransactorRaw struct {
	Contract *OrdinalsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOrdinals creates a new instance of Ordinals, bound to a specific deployed contract.
func NewOrdinals(address common.Address, backend bind.ContractBackend) (*Ordinals, error) {
	contract, err := bindOrdinals(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ordinals{OrdinalsCaller: OrdinalsCaller{contract: contract}, OrdinalsTransactor: OrdinalsTransactor{contract: contract}, OrdinalsFilterer: OrdinalsFilterer{contract: contract}}, nil
}

// NewOrdinalsCaller creates a new read-only instance of Ordinals, bound to a specific deployed contract.
func NewOrdinalsCaller(address common.Address, caller bind.ContractCaller) (*OrdinalsCaller, error) {
	contract, err := bindOrdinals(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OrdinalsCaller{contract: contract}, nil
}

// NewOrdinalsTransactor creates a new write-only instance of Ordinals, bound to a specific deployed contract.
func NewOrdinalsTransactor(address common.Address, transactor bind.ContractTransactor) (*OrdinalsTransactor, error) {
	contract, err := bindOrdinals(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OrdinalsTransactor{contract: contract}, nil
}

// NewOrdinalsFilterer creates a new log filterer instance of Ordinals, bound to a specific deployed contract.
func NewOrdinalsFilterer(address common.Address, filterer bind.ContractFilterer) (*OrdinalsFilterer, error) {
	contract, err := bindOrdinals(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OrdinalsFilterer{contract: contract}, nil
}

// bindOrdinals binds a generic wrapper to an already deployed contract.
func bindOrdinals(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OrdinalsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ordinals *OrdinalsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ordinals.Contract.OrdinalsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ordinals *OrdinalsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ordinals.Contract.OrdinalsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ordinals *OrdinalsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ordinals.Contract.OrdinalsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ordinals *OrdinalsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ordinals.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ordinals *OrdinalsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ordinals.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ordinals *OrdinalsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ordinals.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_Ordinals *OrdinalsCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ordinals.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_Ordinals *OrdinalsSession) Admin() (common.Address, error) {
	return _Ordinals.Contract.Admin(&_Ordinals.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_Ordinals *OrdinalsCallerSession) Admin() (common.Address, error) {
	return _Ordinals.Contract.Admin(&_Ordinals.CallOpts)
}

// Caller is a free data retrieval call binding the contract method 0x72c20ba4.
//
// Solidity: function _caller(address ) view returns(bool)
func (_Ordinals *OrdinalsCaller) Caller(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Ordinals.contract.Call(opts, &out, "_caller", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Caller is a free data retrieval call binding the contract method 0x72c20ba4.
//
// Solidity: function _caller(address ) view returns(bool)
func (_Ordinals *OrdinalsSession) Caller(arg0 common.Address) (bool, error) {
	return _Ordinals.Contract.Caller(&_Ordinals.CallOpts, arg0)
}

// Caller is a free data retrieval call binding the contract method 0x72c20ba4.
//
// Solidity: function _caller(address ) view returns(bool)
func (_Ordinals *OrdinalsCallerSession) Caller(arg0 common.Address) (bool, error) {
	return _Ordinals.Contract.Caller(&_Ordinals.CallOpts, arg0)
}

// Inscription is a free data retrieval call binding the contract method 0x32524dce.
//
// Solidity: function _inscription(address , uint256 ) view returns(string)
func (_Ordinals *OrdinalsCaller) Inscription(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _Ordinals.contract.Call(opts, &out, "_inscription", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Inscription is a free data retrieval call binding the contract method 0x32524dce.
//
// Solidity: function _inscription(address , uint256 ) view returns(string)
func (_Ordinals *OrdinalsSession) Inscription(arg0 common.Address, arg1 *big.Int) (string, error) {
	return _Ordinals.Contract.Inscription(&_Ordinals.CallOpts, arg0, arg1)
}

// Inscription is a free data retrieval call binding the contract method 0x32524dce.
//
// Solidity: function _inscription(address , uint256 ) view returns(string)
func (_Ordinals *OrdinalsCallerSession) Inscription(arg0 common.Address, arg1 *big.Int) (string, error) {
	return _Ordinals.Contract.Inscription(&_Ordinals.CallOpts, arg0, arg1)
}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_Ordinals *OrdinalsCaller) ParameterAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ordinals.contract.Call(opts, &out, "_parameterAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_Ordinals *OrdinalsSession) ParameterAddr() (common.Address, error) {
	return _Ordinals.Contract.ParameterAddr(&_Ordinals.CallOpts)
}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_Ordinals *OrdinalsCallerSession) ParameterAddr() (common.Address, error) {
	return _Ordinals.Contract.ParameterAddr(&_Ordinals.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ordinals *OrdinalsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ordinals.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ordinals *OrdinalsSession) Owner() (common.Address, error) {
	return _Ordinals.Contract.Owner(&_Ordinals.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ordinals *OrdinalsCallerSession) Owner() (common.Address, error) {
	return _Ordinals.Contract.Owner(&_Ordinals.CallOpts)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_Ordinals *OrdinalsTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_Ordinals *OrdinalsSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.ChangeAdmin(&_Ordinals.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_Ordinals *OrdinalsTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.ChangeAdmin(&_Ordinals.TransactOpts, newAdm)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newAdm) returns()
func (_Ordinals *OrdinalsTransactor) ChangeParam(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "changeParam", newAdm)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newAdm) returns()
func (_Ordinals *OrdinalsSession) ChangeParam(newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.ChangeParam(&_Ordinals.TransactOpts, newAdm)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newAdm) returns()
func (_Ordinals *OrdinalsTransactorSession) ChangeParam(newAdm common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.ChangeParam(&_Ordinals.TransactOpts, newAdm)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_Ordinals *OrdinalsTransactor) Initialize(opts *bind.TransactOpts, admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "initialize", admin, parameterControl)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_Ordinals *OrdinalsSession) Initialize(admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.Initialize(&_Ordinals.TransactOpts, admin, parameterControl)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_Ordinals *OrdinalsTransactorSession) Initialize(admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.Initialize(&_Ordinals.TransactOpts, admin, parameterControl)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ordinals *OrdinalsTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ordinals *OrdinalsSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ordinals.Contract.RenounceOwnership(&_Ordinals.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ordinals *OrdinalsTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ordinals.Contract.RenounceOwnership(&_Ordinals.TransactOpts)
}

// SetCaller is a paid mutator transaction binding the contract method 0x9cae6eae.
//
// Solidity: function setCaller(address caller, bool approved) returns()
func (_Ordinals *OrdinalsTransactor) SetCaller(opts *bind.TransactOpts, caller common.Address, approved bool) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "setCaller", caller, approved)
}

// SetCaller is a paid mutator transaction binding the contract method 0x9cae6eae.
//
// Solidity: function setCaller(address caller, bool approved) returns()
func (_Ordinals *OrdinalsSession) SetCaller(caller common.Address, approved bool) (*types.Transaction, error) {
	return _Ordinals.Contract.SetCaller(&_Ordinals.TransactOpts, caller, approved)
}

// SetCaller is a paid mutator transaction binding the contract method 0x9cae6eae.
//
// Solidity: function setCaller(address caller, bool approved) returns()
func (_Ordinals *OrdinalsTransactorSession) SetCaller(caller common.Address, approved bool) (*types.Transaction, error) {
	return _Ordinals.Contract.SetCaller(&_Ordinals.TransactOpts, caller, approved)
}

// SetInscription is a paid mutator transaction binding the contract method 0x1389d3b4.
//
// Solidity: function setInscription(address coll, uint256 tokenId, string inscriptionId) returns()
func (_Ordinals *OrdinalsTransactor) SetInscription(opts *bind.TransactOpts, coll common.Address, tokenId *big.Int, inscriptionId string) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "setInscription", coll, tokenId, inscriptionId)
}

// SetInscription is a paid mutator transaction binding the contract method 0x1389d3b4.
//
// Solidity: function setInscription(address coll, uint256 tokenId, string inscriptionId) returns()
func (_Ordinals *OrdinalsSession) SetInscription(coll common.Address, tokenId *big.Int, inscriptionId string) (*types.Transaction, error) {
	return _Ordinals.Contract.SetInscription(&_Ordinals.TransactOpts, coll, tokenId, inscriptionId)
}

// SetInscription is a paid mutator transaction binding the contract method 0x1389d3b4.
//
// Solidity: function setInscription(address coll, uint256 tokenId, string inscriptionId) returns()
func (_Ordinals *OrdinalsTransactorSession) SetInscription(coll common.Address, tokenId *big.Int, inscriptionId string) (*types.Transaction, error) {
	return _Ordinals.Contract.SetInscription(&_Ordinals.TransactOpts, coll, tokenId, inscriptionId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ordinals *OrdinalsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ordinals.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ordinals *OrdinalsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.TransferOwnership(&_Ordinals.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ordinals *OrdinalsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ordinals.Contract.TransferOwnership(&_Ordinals.TransactOpts, newOwner)
}

// OrdinalsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Ordinals contract.
type OrdinalsInitializedIterator struct {
	Event *OrdinalsInitialized // Event containing the contract specifics and raw log

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
func (it *OrdinalsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrdinalsInitialized)
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
		it.Event = new(OrdinalsInitialized)
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
func (it *OrdinalsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrdinalsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrdinalsInitialized represents a Initialized event raised by the Ordinals contract.
type OrdinalsInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ordinals *OrdinalsFilterer) FilterInitialized(opts *bind.FilterOpts) (*OrdinalsInitializedIterator, error) {

	logs, sub, err := _Ordinals.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OrdinalsInitializedIterator{contract: _Ordinals.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ordinals *OrdinalsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OrdinalsInitialized) (event.Subscription, error) {

	logs, sub, err := _Ordinals.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrdinalsInitialized)
				if err := _Ordinals.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ordinals *OrdinalsFilterer) ParseInitialized(log types.Log) (*OrdinalsInitialized, error) {
	event := new(OrdinalsInitialized)
	if err := _Ordinals.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OrdinalsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ordinals contract.
type OrdinalsOwnershipTransferredIterator struct {
	Event *OrdinalsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OrdinalsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OrdinalsOwnershipTransferred)
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
		it.Event = new(OrdinalsOwnershipTransferred)
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
func (it *OrdinalsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OrdinalsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OrdinalsOwnershipTransferred represents a OwnershipTransferred event raised by the Ordinals contract.
type OrdinalsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ordinals *OrdinalsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OrdinalsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ordinals.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OrdinalsOwnershipTransferredIterator{contract: _Ordinals.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ordinals *OrdinalsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OrdinalsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ordinals.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OrdinalsOwnershipTransferred)
				if err := _Ordinals.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ordinals *OrdinalsFilterer) ParseOwnershipTransferred(log types.Log) (*OrdinalsOwnershipTransferred, error) {
	event := new(OrdinalsOwnershipTransferred)
	if err := _Ordinals.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
