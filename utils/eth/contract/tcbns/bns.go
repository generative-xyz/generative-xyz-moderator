// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tcbns

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

// BNSMetaData contains all meta data concerning the BNS contract.
var BNSMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InsufficientRegistrationFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NameAlreadyRegistered\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"NameRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"map\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minRegistrationFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"names\",\"type\":\"bytes[]\"}],\"name\":\"registerBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"registered\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"resolver\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"setMinRegistrationFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BNSABI is the input ABI used to generate the binding from.
// Deprecated: Use BNSMetaData.ABI instead.
var BNSABI = BNSMetaData.ABI

// BNS is an auto generated Go binding around an Ethereum contract.
type BNS struct {
	BNSCaller     // Read-only binding to the contract
	BNSTransactor // Write-only binding to the contract
	BNSFilterer   // Log filterer for contract events
}

// BNSCaller is an auto generated read-only Go binding around an Ethereum contract.
type BNSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BNSTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BNSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BNSFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BNSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BNSSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BNSSession struct {
	Contract     *BNS              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BNSCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BNSCallerSession struct {
	Contract *BNSCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BNSTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BNSTransactorSession struct {
	Contract     *BNSTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BNSRaw is an auto generated low-level Go binding around an Ethereum contract.
type BNSRaw struct {
	Contract *BNS // Generic contract binding to access the raw methods on
}

// BNSCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BNSCallerRaw struct {
	Contract *BNSCaller // Generic read-only contract binding to access the raw methods on
}

// BNSTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BNSTransactorRaw struct {
	Contract *BNSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBNS creates a new instance of BNS, bound to a specific deployed contract.
func NewBNS(address common.Address, backend bind.ContractBackend) (*BNS, error) {
	contract, err := bindBNS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BNS{BNSCaller: BNSCaller{contract: contract}, BNSTransactor: BNSTransactor{contract: contract}, BNSFilterer: BNSFilterer{contract: contract}}, nil
}

// NewBNSCaller creates a new read-only instance of BNS, bound to a specific deployed contract.
func NewBNSCaller(address common.Address, caller bind.ContractCaller) (*BNSCaller, error) {
	contract, err := bindBNS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BNSCaller{contract: contract}, nil
}

// NewBNSTransactor creates a new write-only instance of BNS, bound to a specific deployed contract.
func NewBNSTransactor(address common.Address, transactor bind.ContractTransactor) (*BNSTransactor, error) {
	contract, err := bindBNS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BNSTransactor{contract: contract}, nil
}

// NewBNSFilterer creates a new log filterer instance of BNS, bound to a specific deployed contract.
func NewBNSFilterer(address common.Address, filterer bind.ContractFilterer) (*BNSFilterer, error) {
	contract, err := bindBNS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BNSFilterer{contract: contract}, nil
}

// bindBNS binds a generic wrapper to an already deployed contract.
func bindBNS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BNSMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BNS *BNSRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BNS.Contract.BNSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BNS *BNSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BNS.Contract.BNSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BNS *BNSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BNS.Contract.BNSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BNS *BNSCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BNS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BNS *BNSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BNS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BNS *BNSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BNS.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BNS *BNSCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BNS *BNSSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BNS.Contract.BalanceOf(&_BNS.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_BNS *BNSCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _BNS.Contract.BalanceOf(&_BNS.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BNS *BNSCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BNS *BNSSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BNS.Contract.GetApproved(&_BNS.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_BNS *BNSCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _BNS.Contract.GetApproved(&_BNS.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BNS *BNSCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BNS *BNSSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BNS.Contract.IsApprovedForAll(&_BNS.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_BNS *BNSCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _BNS.Contract.IsApprovedForAll(&_BNS.CallOpts, owner, operator)
}

// MinRegistrationFee is a free data retrieval call binding the contract method 0x9706828a.
//
// Solidity: function minRegistrationFee() view returns(uint256)
func (_BNS *BNSCaller) MinRegistrationFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "minRegistrationFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinRegistrationFee is a free data retrieval call binding the contract method 0x9706828a.
//
// Solidity: function minRegistrationFee() view returns(uint256)
func (_BNS *BNSSession) MinRegistrationFee() (*big.Int, error) {
	return _BNS.Contract.MinRegistrationFee(&_BNS.CallOpts)
}

// MinRegistrationFee is a free data retrieval call binding the contract method 0x9706828a.
//
// Solidity: function minRegistrationFee() view returns(uint256)
func (_BNS *BNSCallerSession) MinRegistrationFee() (*big.Int, error) {
	return _BNS.Contract.MinRegistrationFee(&_BNS.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BNS *BNSCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BNS *BNSSession) Name() (string, error) {
	return _BNS.Contract.Name(&_BNS.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_BNS *BNSCallerSession) Name() (string, error) {
	return _BNS.Contract.Name(&_BNS.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BNS *BNSCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BNS *BNSSession) Owner() (common.Address, error) {
	return _BNS.Contract.Owner(&_BNS.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BNS *BNSCallerSession) Owner() (common.Address, error) {
	return _BNS.Contract.Owner(&_BNS.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BNS *BNSCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BNS *BNSSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BNS.Contract.OwnerOf(&_BNS.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_BNS *BNSCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _BNS.Contract.OwnerOf(&_BNS.CallOpts, tokenId)
}

// Registered is a free data retrieval call binding the contract method 0x5aca952e.
//
// Solidity: function registered(bytes ) view returns(bool)
func (_BNS *BNSCaller) Registered(opts *bind.CallOpts, arg0 []byte) (bool, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "registered", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Registered is a free data retrieval call binding the contract method 0x5aca952e.
//
// Solidity: function registered(bytes ) view returns(bool)
func (_BNS *BNSSession) Registered(arg0 []byte) (bool, error) {
	return _BNS.Contract.Registered(&_BNS.CallOpts, arg0)
}

// Registered is a free data retrieval call binding the contract method 0x5aca952e.
//
// Solidity: function registered(bytes ) view returns(bool)
func (_BNS *BNSCallerSession) Registered(arg0 []byte) (bool, error) {
	return _BNS.Contract.Registered(&_BNS.CallOpts, arg0)
}

// Registry is a free data retrieval call binding the contract method 0xa15d581c.
//
// Solidity: function registry(bytes ) view returns(uint256)
func (_BNS *BNSCaller) Registry(opts *bind.CallOpts, arg0 []byte) (*big.Int, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "registry", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0xa15d581c.
//
// Solidity: function registry(bytes ) view returns(uint256)
func (_BNS *BNSSession) Registry(arg0 []byte) (*big.Int, error) {
	return _BNS.Contract.Registry(&_BNS.CallOpts, arg0)
}

// Registry is a free data retrieval call binding the contract method 0xa15d581c.
//
// Solidity: function registry(bytes ) view returns(uint256)
func (_BNS *BNSCallerSession) Registry(arg0 []byte) (*big.Int, error) {
	return _BNS.Contract.Registry(&_BNS.CallOpts, arg0)
}

// Resolver is a free data retrieval call binding the contract method 0x108eaa4e.
//
// Solidity: function resolver(uint256 ) view returns(address)
func (_BNS *BNSCaller) Resolver(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "resolver", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolver is a free data retrieval call binding the contract method 0x108eaa4e.
//
// Solidity: function resolver(uint256 ) view returns(address)
func (_BNS *BNSSession) Resolver(arg0 *big.Int) (common.Address, error) {
	return _BNS.Contract.Resolver(&_BNS.CallOpts, arg0)
}

// Resolver is a free data retrieval call binding the contract method 0x108eaa4e.
//
// Solidity: function resolver(uint256 ) view returns(address)
func (_BNS *BNSCallerSession) Resolver(arg0 *big.Int) (common.Address, error) {
	return _BNS.Contract.Resolver(&_BNS.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BNS *BNSCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BNS *BNSSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BNS.Contract.SupportsInterface(&_BNS.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BNS *BNSCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BNS.Contract.SupportsInterface(&_BNS.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BNS *BNSCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BNS *BNSSession) Symbol() (string, error) {
	return _BNS.Contract.Symbol(&_BNS.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_BNS *BNSCallerSession) Symbol() (string, error) {
	return _BNS.Contract.Symbol(&_BNS.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BNS *BNSCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _BNS.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BNS *BNSSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BNS.Contract.TokenURI(&_BNS.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_BNS *BNSCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _BNS.Contract.TokenURI(&_BNS.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BNS *BNSTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BNS *BNSSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.Approve(&_BNS.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_BNS *BNSTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.Approve(&_BNS.TransactOpts, to, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_BNS *BNSTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_BNS *BNSSession) Initialize() (*types.Transaction, error) {
	return _BNS.Contract.Initialize(&_BNS.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_BNS *BNSTransactorSession) Initialize() (*types.Transaction, error) {
	return _BNS.Contract.Initialize(&_BNS.TransactOpts)
}

// Map is a paid mutator transaction binding the contract method 0x376a2b69.
//
// Solidity: function map(uint256 tokenId, address to) returns()
func (_BNS *BNSTransactor) Map(opts *bind.TransactOpts, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "map", tokenId, to)
}

// Map is a paid mutator transaction binding the contract method 0x376a2b69.
//
// Solidity: function map(uint256 tokenId, address to) returns()
func (_BNS *BNSSession) Map(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _BNS.Contract.Map(&_BNS.TransactOpts, tokenId, to)
}

// Map is a paid mutator transaction binding the contract method 0x376a2b69.
//
// Solidity: function map(uint256 tokenId, address to) returns()
func (_BNS *BNSTransactorSession) Map(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _BNS.Contract.Map(&_BNS.TransactOpts, tokenId, to)
}

// Register is a paid mutator transaction binding the contract method 0x24b8fbf6.
//
// Solidity: function register(address owner, bytes name) payable returns(uint256)
func (_BNS *BNSTransactor) Register(opts *bind.TransactOpts, owner common.Address, name []byte) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "register", owner, name)
}

// Register is a paid mutator transaction binding the contract method 0x24b8fbf6.
//
// Solidity: function register(address owner, bytes name) payable returns(uint256)
func (_BNS *BNSSession) Register(owner common.Address, name []byte) (*types.Transaction, error) {
	return _BNS.Contract.Register(&_BNS.TransactOpts, owner, name)
}

// Register is a paid mutator transaction binding the contract method 0x24b8fbf6.
//
// Solidity: function register(address owner, bytes name) payable returns(uint256)
func (_BNS *BNSTransactorSession) Register(owner common.Address, name []byte) (*types.Transaction, error) {
	return _BNS.Contract.Register(&_BNS.TransactOpts, owner, name)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x27201c97.
//
// Solidity: function registerBatch(address owner, bytes[] names) returns()
func (_BNS *BNSTransactor) RegisterBatch(opts *bind.TransactOpts, owner common.Address, names [][]byte) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "registerBatch", owner, names)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x27201c97.
//
// Solidity: function registerBatch(address owner, bytes[] names) returns()
func (_BNS *BNSSession) RegisterBatch(owner common.Address, names [][]byte) (*types.Transaction, error) {
	return _BNS.Contract.RegisterBatch(&_BNS.TransactOpts, owner, names)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x27201c97.
//
// Solidity: function registerBatch(address owner, bytes[] names) returns()
func (_BNS *BNSTransactorSession) RegisterBatch(owner common.Address, names [][]byte) (*types.Transaction, error) {
	return _BNS.Contract.RegisterBatch(&_BNS.TransactOpts, owner, names)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BNS *BNSTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BNS *BNSSession) RenounceOwnership() (*types.Transaction, error) {
	return _BNS.Contract.RenounceOwnership(&_BNS.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BNS *BNSTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BNS.Contract.RenounceOwnership(&_BNS.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.SafeTransferFrom(&_BNS.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.SafeTransferFrom(&_BNS.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BNS *BNSTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BNS *BNSSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BNS.Contract.SafeTransferFrom0(&_BNS.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_BNS *BNSTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _BNS.Contract.SafeTransferFrom0(&_BNS.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BNS *BNSTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BNS *BNSSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BNS.Contract.SetApprovalForAll(&_BNS.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_BNS *BNSTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _BNS.Contract.SetApprovalForAll(&_BNS.TransactOpts, operator, approved)
}

// SetMinRegistrationFee is a paid mutator transaction binding the contract method 0x06301a9d.
//
// Solidity: function setMinRegistrationFee(uint256 fee) returns()
func (_BNS *BNSTransactor) SetMinRegistrationFee(opts *bind.TransactOpts, fee *big.Int) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "setMinRegistrationFee", fee)
}

// SetMinRegistrationFee is a paid mutator transaction binding the contract method 0x06301a9d.
//
// Solidity: function setMinRegistrationFee(uint256 fee) returns()
func (_BNS *BNSSession) SetMinRegistrationFee(fee *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.SetMinRegistrationFee(&_BNS.TransactOpts, fee)
}

// SetMinRegistrationFee is a paid mutator transaction binding the contract method 0x06301a9d.
//
// Solidity: function setMinRegistrationFee(uint256 fee) returns()
func (_BNS *BNSTransactorSession) SetMinRegistrationFee(fee *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.SetMinRegistrationFee(&_BNS.TransactOpts, fee)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.TransferFrom(&_BNS.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_BNS *BNSTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _BNS.Contract.TransferFrom(&_BNS.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BNS *BNSTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BNS.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BNS *BNSSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BNS.Contract.TransferOwnership(&_BNS.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BNS *BNSTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BNS.Contract.TransferOwnership(&_BNS.TransactOpts, newOwner)
}

// BNSApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the BNS contract.
type BNSApprovalIterator struct {
	Event *BNSApproval // Event containing the contract specifics and raw log

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
func (it *BNSApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSApproval)
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
		it.Event = new(BNSApproval)
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
func (it *BNSApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSApproval represents a Approval event raised by the BNS contract.
type BNSApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BNS *BNSFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*BNSApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BNS.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BNSApprovalIterator{contract: _BNS.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BNS *BNSFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *BNSApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BNS.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSApproval)
				if err := _BNS.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_BNS *BNSFilterer) ParseApproval(log types.Log) (*BNSApproval, error) {
	event := new(BNSApproval)
	if err := _BNS.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BNSApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the BNS contract.
type BNSApprovalForAllIterator struct {
	Event *BNSApprovalForAll // Event containing the contract specifics and raw log

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
func (it *BNSApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSApprovalForAll)
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
		it.Event = new(BNSApprovalForAll)
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
func (it *BNSApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSApprovalForAll represents a ApprovalForAll event raised by the BNS contract.
type BNSApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BNS *BNSFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*BNSApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BNS.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &BNSApprovalForAllIterator{contract: _BNS.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BNS *BNSFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *BNSApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _BNS.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSApprovalForAll)
				if err := _BNS.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_BNS *BNSFilterer) ParseApprovalForAll(log types.Log) (*BNSApprovalForAll, error) {
	event := new(BNSApprovalForAll)
	if err := _BNS.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BNSInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the BNS contract.
type BNSInitializedIterator struct {
	Event *BNSInitialized // Event containing the contract specifics and raw log

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
func (it *BNSInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSInitialized)
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
		it.Event = new(BNSInitialized)
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
func (it *BNSInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSInitialized represents a Initialized event raised by the BNS contract.
type BNSInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BNS *BNSFilterer) FilterInitialized(opts *bind.FilterOpts) (*BNSInitializedIterator, error) {

	logs, sub, err := _BNS.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BNSInitializedIterator{contract: _BNS.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BNS *BNSFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BNSInitialized) (event.Subscription, error) {

	logs, sub, err := _BNS.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSInitialized)
				if err := _BNS.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_BNS *BNSFilterer) ParseInitialized(log types.Log) (*BNSInitialized, error) {
	event := new(BNSInitialized)
	if err := _BNS.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BNSNameRegisteredIterator is returned from FilterNameRegistered and is used to iterate over the raw logs and unpacked data for NameRegistered events raised by the BNS contract.
type BNSNameRegisteredIterator struct {
	Event *BNSNameRegistered // Event containing the contract specifics and raw log

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
func (it *BNSNameRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSNameRegistered)
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
		it.Event = new(BNSNameRegistered)
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
func (it *BNSNameRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSNameRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSNameRegistered represents a NameRegistered event raised by the BNS contract.
type BNSNameRegistered struct {
	Name []byte
	Id   *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNameRegistered is a free log retrieval operation binding the contract event 0x804e89f7acfccaf194f6757b177acb1ab1828a7c01ef308f9f8e4bbd2fb63873.
//
// Solidity: event NameRegistered(bytes name, uint256 indexed id)
func (_BNS *BNSFilterer) FilterNameRegistered(opts *bind.FilterOpts, id []*big.Int) (*BNSNameRegisteredIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _BNS.contract.FilterLogs(opts, "NameRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return &BNSNameRegisteredIterator{contract: _BNS.contract, event: "NameRegistered", logs: logs, sub: sub}, nil
}

// WatchNameRegistered is a free log subscription operation binding the contract event 0x804e89f7acfccaf194f6757b177acb1ab1828a7c01ef308f9f8e4bbd2fb63873.
//
// Solidity: event NameRegistered(bytes name, uint256 indexed id)
func (_BNS *BNSFilterer) WatchNameRegistered(opts *bind.WatchOpts, sink chan<- *BNSNameRegistered, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _BNS.contract.WatchLogs(opts, "NameRegistered", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSNameRegistered)
				if err := _BNS.contract.UnpackLog(event, "NameRegistered", log); err != nil {
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

// ParseNameRegistered is a log parse operation binding the contract event 0x804e89f7acfccaf194f6757b177acb1ab1828a7c01ef308f9f8e4bbd2fb63873.
//
// Solidity: event NameRegistered(bytes name, uint256 indexed id)
func (_BNS *BNSFilterer) ParseNameRegistered(log types.Log) (*BNSNameRegistered, error) {
	event := new(BNSNameRegistered)
	if err := _BNS.contract.UnpackLog(event, "NameRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BNSOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BNS contract.
type BNSOwnershipTransferredIterator struct {
	Event *BNSOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BNSOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSOwnershipTransferred)
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
		it.Event = new(BNSOwnershipTransferred)
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
func (it *BNSOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSOwnershipTransferred represents a OwnershipTransferred event raised by the BNS contract.
type BNSOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BNS *BNSFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BNSOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BNS.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BNSOwnershipTransferredIterator{contract: _BNS.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BNS *BNSFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BNSOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BNS.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSOwnershipTransferred)
				if err := _BNS.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BNS *BNSFilterer) ParseOwnershipTransferred(log types.Log) (*BNSOwnershipTransferred, error) {
	event := new(BNSOwnershipTransferred)
	if err := _BNS.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BNSTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the BNS contract.
type BNSTransferIterator struct {
	Event *BNSTransfer // Event containing the contract specifics and raw log

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
func (it *BNSTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BNSTransfer)
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
		it.Event = new(BNSTransfer)
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
func (it *BNSTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BNSTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BNSTransfer represents a Transfer event raised by the BNS contract.
type BNSTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BNS *BNSFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*BNSTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BNS.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &BNSTransferIterator{contract: _BNS.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BNS *BNSFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BNSTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _BNS.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BNSTransfer)
				if err := _BNS.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_BNS *BNSFilterer) ParseTransfer(log types.Log) (*BNSTransfer, error) {
	event := new(BNSTransfer)
	if err := _BNS.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
