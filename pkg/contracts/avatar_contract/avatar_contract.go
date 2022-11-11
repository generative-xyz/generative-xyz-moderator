// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package avatar_contract

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

// AVATARSPlayer is an auto generated low-level Go binding around an user-defined struct.
type AVATARSPlayer struct {
	Emotion     string
	EmotionTime string
	Nation      string
	Dna         string
	Beard       string
	Hair        string
	Undershirt  string
	Shoes       string
	Top         string
	Bottom      string
	Number      *big.Int
	Tatoo       string
	Glasses     string
	Captain     string
}

// AvatarContractMetaData contains all meta data concerning the AvatarContract contract.
var AvatarContractMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"RequestFulfilledData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_algorithm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_requestIdData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_tokenAddrErc721\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_whitelistFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"count\",\"type\":\"uint256\"}],\"name\":\"addWhitelist\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"changeBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"changeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sweet\",\"type\":\"address\"}],\"name\":\"changeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"requestId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"gameData\",\"type\":\"bytes\"}],\"name\":\"fulfill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getParamValues\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"_emotion\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_emotionTime\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_nation\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_dna\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_beard\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_hair\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_undershirt\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_shoes\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_top\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_bottom\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_number\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_tatoo\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_glasses\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_captain\",\"type\":\"string\"}],\"internalType\":\"structAVATARS.Player\",\"name\":\"player\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenIdGated\",\"type\":\"uint256\"}],\"name\":\"mintByToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mintWhitelist\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"algo\",\"type\":\"string\"}],\"name\":\"setAlgo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"whitelistFee\",\"type\":\"uint256\"}],\"name\":\"setWhitelistFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AvatarContractABI is the input ABI used to generate the binding from.
// Deprecated: Use AvatarContractMetaData.ABI instead.
var AvatarContractABI = AvatarContractMetaData.ABI

// AvatarContract is an auto generated Go binding around an Ethereum contract.
type AvatarContract struct {
	AvatarContractCaller     // Read-only binding to the contract
	AvatarContractTransactor // Write-only binding to the contract
	AvatarContractFilterer   // Log filterer for contract events
}

// AvatarContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type AvatarContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AvatarContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AvatarContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AvatarContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AvatarContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AvatarContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AvatarContractSession struct {
	Contract     *AvatarContract   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AvatarContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AvatarContractCallerSession struct {
	Contract *AvatarContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AvatarContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AvatarContractTransactorSession struct {
	Contract     *AvatarContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AvatarContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type AvatarContractRaw struct {
	Contract *AvatarContract // Generic contract binding to access the raw methods on
}

// AvatarContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AvatarContractCallerRaw struct {
	Contract *AvatarContractCaller // Generic read-only contract binding to access the raw methods on
}

// AvatarContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AvatarContractTransactorRaw struct {
	Contract *AvatarContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAvatarContract creates a new instance of AvatarContract, bound to a specific deployed contract.
func NewAvatarContract(address common.Address, backend bind.ContractBackend) (*AvatarContract, error) {
	contract, err := bindAvatarContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AvatarContract{AvatarContractCaller: AvatarContractCaller{contract: contract}, AvatarContractTransactor: AvatarContractTransactor{contract: contract}, AvatarContractFilterer: AvatarContractFilterer{contract: contract}}, nil
}

// NewAvatarContractCaller creates a new read-only instance of AvatarContract, bound to a specific deployed contract.
func NewAvatarContractCaller(address common.Address, caller bind.ContractCaller) (*AvatarContractCaller, error) {
	contract, err := bindAvatarContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AvatarContractCaller{contract: contract}, nil
}

// NewAvatarContractTransactor creates a new write-only instance of AvatarContract, bound to a specific deployed contract.
func NewAvatarContractTransactor(address common.Address, transactor bind.ContractTransactor) (*AvatarContractTransactor, error) {
	contract, err := bindAvatarContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AvatarContractTransactor{contract: contract}, nil
}

// NewAvatarContractFilterer creates a new log filterer instance of AvatarContract, bound to a specific deployed contract.
func NewAvatarContractFilterer(address common.Address, filterer bind.ContractFilterer) (*AvatarContractFilterer, error) {
	contract, err := bindAvatarContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AvatarContractFilterer{contract: contract}, nil
}

// bindAvatarContract binds a generic wrapper to an already deployed contract.
func bindAvatarContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AvatarContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AvatarContract *AvatarContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AvatarContract.Contract.AvatarContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AvatarContract *AvatarContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.Contract.AvatarContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AvatarContract *AvatarContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AvatarContract.Contract.AvatarContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AvatarContract *AvatarContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AvatarContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AvatarContract *AvatarContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AvatarContract *AvatarContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AvatarContract.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_AvatarContract *AvatarContractCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_AvatarContract *AvatarContractSession) Admin() (common.Address, error) {
	return _AvatarContract.Contract.Admin(&_AvatarContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_AvatarContract *AvatarContractCallerSession) Admin() (common.Address, error) {
	return _AvatarContract.Contract.Admin(&_AvatarContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_AvatarContract *AvatarContractCaller) Algorithm(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_algorithm")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_AvatarContract *AvatarContractSession) Algorithm() (string, error) {
	return _AvatarContract.Contract.Algorithm(&_AvatarContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_AvatarContract *AvatarContractCallerSession) Algorithm() (string, error) {
	return _AvatarContract.Contract.Algorithm(&_AvatarContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_AvatarContract *AvatarContractCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_AvatarContract *AvatarContractSession) Counter() (*big.Int, error) {
	return _AvatarContract.Contract.Counter(&_AvatarContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_AvatarContract *AvatarContractCallerSession) Counter() (*big.Int, error) {
	return _AvatarContract.Contract.Counter(&_AvatarContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_AvatarContract *AvatarContractCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_AvatarContract *AvatarContractSession) Fee() (*big.Int, error) {
	return _AvatarContract.Contract.Fee(&_AvatarContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_AvatarContract *AvatarContractCallerSession) Fee() (*big.Int, error) {
	return _AvatarContract.Contract.Fee(&_AvatarContract.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x247ac063.
//
// Solidity: function _oracle() view returns(address)
func (_AvatarContract *AvatarContractCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x247ac063.
//
// Solidity: function _oracle() view returns(address)
func (_AvatarContract *AvatarContractSession) Oracle() (common.Address, error) {
	return _AvatarContract.Contract.Oracle(&_AvatarContract.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x247ac063.
//
// Solidity: function _oracle() view returns(address)
func (_AvatarContract *AvatarContractCallerSession) Oracle() (common.Address, error) {
	return _AvatarContract.Contract.Oracle(&_AvatarContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_AvatarContract *AvatarContractCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_AvatarContract *AvatarContractSession) ParamsAddress() (common.Address, error) {
	return _AvatarContract.Contract.ParamsAddress(&_AvatarContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_AvatarContract *AvatarContractCallerSession) ParamsAddress() (common.Address, error) {
	return _AvatarContract.Contract.ParamsAddress(&_AvatarContract.CallOpts)
}

// RequestIdData is a free data retrieval call binding the contract method 0xd5978314.
//
// Solidity: function _requestIdData(bytes32 ) view returns(bytes)
func (_AvatarContract *AvatarContractCaller) RequestIdData(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_requestIdData", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RequestIdData is a free data retrieval call binding the contract method 0xd5978314.
//
// Solidity: function _requestIdData(bytes32 ) view returns(bytes)
func (_AvatarContract *AvatarContractSession) RequestIdData(arg0 [32]byte) ([]byte, error) {
	return _AvatarContract.Contract.RequestIdData(&_AvatarContract.CallOpts, arg0)
}

// RequestIdData is a free data retrieval call binding the contract method 0xd5978314.
//
// Solidity: function _requestIdData(bytes32 ) view returns(bytes)
func (_AvatarContract *AvatarContractCallerSession) RequestIdData(arg0 [32]byte) ([]byte, error) {
	return _AvatarContract.Contract.RequestIdData(&_AvatarContract.CallOpts, arg0)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_AvatarContract *AvatarContractCaller) TokenAddrErc721(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_tokenAddrErc721")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_AvatarContract *AvatarContractSession) TokenAddrErc721() (common.Address, error) {
	return _AvatarContract.Contract.TokenAddrErc721(&_AvatarContract.CallOpts)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_AvatarContract *AvatarContractCallerSession) TokenAddrErc721() (common.Address, error) {
	return _AvatarContract.Contract.TokenAddrErc721(&_AvatarContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_AvatarContract *AvatarContractCaller) Uri(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_uri")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_AvatarContract *AvatarContractSession) Uri() (string, error) {
	return _AvatarContract.Contract.Uri(&_AvatarContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_AvatarContract *AvatarContractCallerSession) Uri() (string, error) {
	return _AvatarContract.Contract.Uri(&_AvatarContract.CallOpts)
}

// WhitelistFee is a free data retrieval call binding the contract method 0x5275f37a.
//
// Solidity: function _whitelistFee() view returns(uint256)
func (_AvatarContract *AvatarContractCaller) WhitelistFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "_whitelistFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WhitelistFee is a free data retrieval call binding the contract method 0x5275f37a.
//
// Solidity: function _whitelistFee() view returns(uint256)
func (_AvatarContract *AvatarContractSession) WhitelistFee() (*big.Int, error) {
	return _AvatarContract.Contract.WhitelistFee(&_AvatarContract.CallOpts)
}

// WhitelistFee is a free data retrieval call binding the contract method 0x5275f37a.
//
// Solidity: function _whitelistFee() view returns(uint256)
func (_AvatarContract *AvatarContractCallerSession) WhitelistFee() (*big.Int, error) {
	return _AvatarContract.Contract.WhitelistFee(&_AvatarContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AvatarContract *AvatarContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AvatarContract *AvatarContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AvatarContract.Contract.BalanceOf(&_AvatarContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_AvatarContract *AvatarContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _AvatarContract.Contract.BalanceOf(&_AvatarContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AvatarContract.Contract.GetApproved(&_AvatarContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _AvatarContract.Contract.GetApproved(&_AvatarContract.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string,string,string,string,string,string,string,string,uint256,string,string,string) player)
func (_AvatarContract *AvatarContractCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) (AVATARSPlayer, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new(AVATARSPlayer), err
	}

	out0 := *abi.ConvertType(out[0], new(AVATARSPlayer)).(*AVATARSPlayer)

	return out0, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string,string,string,string,string,string,string,string,uint256,string,string,string) player)
func (_AvatarContract *AvatarContractSession) GetParamValues(tokenId *big.Int) (AVATARSPlayer, error) {
	return _AvatarContract.Contract.GetParamValues(&_AvatarContract.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string,string,string,string,string,string,string,string,uint256,string,string,string) player)
func (_AvatarContract *AvatarContractCallerSession) GetParamValues(tokenId *big.Int) (AVATARSPlayer, error) {
	return _AvatarContract.Contract.GetParamValues(&_AvatarContract.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AvatarContract *AvatarContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AvatarContract *AvatarContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AvatarContract.Contract.IsApprovedForAll(&_AvatarContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_AvatarContract *AvatarContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _AvatarContract.Contract.IsApprovedForAll(&_AvatarContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AvatarContract *AvatarContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AvatarContract *AvatarContractSession) Name() (string, error) {
	return _AvatarContract.Contract.Name(&_AvatarContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AvatarContract *AvatarContractCallerSession) Name() (string, error) {
	return _AvatarContract.Contract.Name(&_AvatarContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AvatarContract *AvatarContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AvatarContract *AvatarContractSession) Owner() (common.Address, error) {
	return _AvatarContract.Contract.Owner(&_AvatarContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AvatarContract *AvatarContractCallerSession) Owner() (common.Address, error) {
	return _AvatarContract.Contract.Owner(&_AvatarContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AvatarContract.Contract.OwnerOf(&_AvatarContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_AvatarContract *AvatarContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _AvatarContract.Contract.OwnerOf(&_AvatarContract.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AvatarContract *AvatarContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AvatarContract *AvatarContractSession) Paused() (bool, error) {
	return _AvatarContract.Contract.Paused(&_AvatarContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AvatarContract *AvatarContractCallerSession) Paused() (bool, error) {
	return _AvatarContract.Contract.Paused(&_AvatarContract.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_AvatarContract *AvatarContractCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

	outstruct := new(struct {
		Receiver      common.Address
		RoyaltyAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.RoyaltyAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_AvatarContract *AvatarContractSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _AvatarContract.Contract.RoyaltyInfo(&_AvatarContract.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_AvatarContract *AvatarContractCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _AvatarContract.Contract.RoyaltyInfo(&_AvatarContract.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AvatarContract *AvatarContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AvatarContract *AvatarContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AvatarContract.Contract.SupportsInterface(&_AvatarContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AvatarContract *AvatarContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AvatarContract.Contract.SupportsInterface(&_AvatarContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AvatarContract *AvatarContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AvatarContract *AvatarContractSession) Symbol() (string, error) {
	return _AvatarContract.Contract.Symbol(&_AvatarContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AvatarContract *AvatarContractCallerSession) Symbol() (string, error) {
	return _AvatarContract.Contract.Symbol(&_AvatarContract.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AvatarContract *AvatarContractCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _AvatarContract.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AvatarContract *AvatarContractSession) TokenURI(tokenId *big.Int) (string, error) {
	return _AvatarContract.Contract.TokenURI(&_AvatarContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_AvatarContract *AvatarContractCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _AvatarContract.Contract.TokenURI(&_AvatarContract.CallOpts, tokenId)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0x645b6701.
//
// Solidity: function addWhitelist(address[] addrs, uint256 count) returns()
func (_AvatarContract *AvatarContractTransactor) AddWhitelist(opts *bind.TransactOpts, addrs []common.Address, count *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "addWhitelist", addrs, count)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0x645b6701.
//
// Solidity: function addWhitelist(address[] addrs, uint256 count) returns()
func (_AvatarContract *AvatarContractSession) AddWhitelist(addrs []common.Address, count *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.AddWhitelist(&_AvatarContract.TransactOpts, addrs, count)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0x645b6701.
//
// Solidity: function addWhitelist(address[] addrs, uint256 count) returns()
func (_AvatarContract *AvatarContractTransactorSession) AddWhitelist(addrs []common.Address, count *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.AddWhitelist(&_AvatarContract.TransactOpts, addrs, count)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.Approve(&_AvatarContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.Approve(&_AvatarContract.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_AvatarContract *AvatarContractTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_AvatarContract *AvatarContractSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeAdmin(&_AvatarContract.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_AvatarContract *AvatarContractTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeAdmin(&_AvatarContract.TransactOpts, newAdm)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_AvatarContract *AvatarContractTransactor) ChangeBaseURI(opts *bind.TransactOpts, baseURI string) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "changeBaseURI", baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_AvatarContract *AvatarContractSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeBaseURI(&_AvatarContract.TransactOpts, baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_AvatarContract *AvatarContractTransactorSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeBaseURI(&_AvatarContract.TransactOpts, baseURI)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(address oracle) returns()
func (_AvatarContract *AvatarContractTransactor) ChangeOracle(opts *bind.TransactOpts, oracle common.Address) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "changeOracle", oracle)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(address oracle) returns()
func (_AvatarContract *AvatarContractSession) ChangeOracle(oracle common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeOracle(&_AvatarContract.TransactOpts, oracle)
}

// ChangeOracle is a paid mutator transaction binding the contract method 0x47c421b5.
//
// Solidity: function changeOracle(address oracle) returns()
func (_AvatarContract *AvatarContractTransactorSession) ChangeOracle(oracle common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeOracle(&_AvatarContract.TransactOpts, oracle)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_AvatarContract *AvatarContractTransactor) ChangeToken(opts *bind.TransactOpts, sweet common.Address) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "changeToken", sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_AvatarContract *AvatarContractSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeToken(&_AvatarContract.TransactOpts, sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_AvatarContract *AvatarContractTransactorSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.ChangeToken(&_AvatarContract.TransactOpts, sweet)
}

// Fulfill is a paid mutator transaction binding the contract method 0x7c1de7e1.
//
// Solidity: function fulfill(bytes32 requestId, bytes gameData) returns()
func (_AvatarContract *AvatarContractTransactor) Fulfill(opts *bind.TransactOpts, requestId [32]byte, gameData []byte) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "fulfill", requestId, gameData)
}

// Fulfill is a paid mutator transaction binding the contract method 0x7c1de7e1.
//
// Solidity: function fulfill(bytes32 requestId, bytes gameData) returns()
func (_AvatarContract *AvatarContractSession) Fulfill(requestId [32]byte, gameData []byte) (*types.Transaction, error) {
	return _AvatarContract.Contract.Fulfill(&_AvatarContract.TransactOpts, requestId, gameData)
}

// Fulfill is a paid mutator transaction binding the contract method 0x7c1de7e1.
//
// Solidity: function fulfill(bytes32 requestId, bytes gameData) returns()
func (_AvatarContract *AvatarContractTransactorSession) Fulfill(requestId [32]byte, gameData []byte) (*types.Transaction, error) {
	return _AvatarContract.Contract.Fulfill(&_AvatarContract.TransactOpts, requestId, gameData)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_AvatarContract *AvatarContractTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "initialize", name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_AvatarContract *AvatarContractSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.Initialize(&_AvatarContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_AvatarContract *AvatarContractTransactorSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.Initialize(&_AvatarContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_AvatarContract *AvatarContractTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_AvatarContract *AvatarContractSession) Mint() (*types.Transaction, error) {
	return _AvatarContract.Contract.Mint(&_AvatarContract.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_AvatarContract *AvatarContractTransactorSession) Mint() (*types.Transaction, error) {
	return _AvatarContract.Contract.Mint(&_AvatarContract.TransactOpts)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_AvatarContract *AvatarContractTransactor) MintByToken(opts *bind.TransactOpts, tokenIdGated *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "mintByToken", tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_AvatarContract *AvatarContractSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.MintByToken(&_AvatarContract.TransactOpts, tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_AvatarContract *AvatarContractTransactorSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.MintByToken(&_AvatarContract.TransactOpts, tokenIdGated)
}

// MintWhitelist is a paid mutator transaction binding the contract method 0x2d3df31f.
//
// Solidity: function mintWhitelist() payable returns()
func (_AvatarContract *AvatarContractTransactor) MintWhitelist(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "mintWhitelist")
}

// MintWhitelist is a paid mutator transaction binding the contract method 0x2d3df31f.
//
// Solidity: function mintWhitelist() payable returns()
func (_AvatarContract *AvatarContractSession) MintWhitelist() (*types.Transaction, error) {
	return _AvatarContract.Contract.MintWhitelist(&_AvatarContract.TransactOpts)
}

// MintWhitelist is a paid mutator transaction binding the contract method 0x2d3df31f.
//
// Solidity: function mintWhitelist() payable returns()
func (_AvatarContract *AvatarContractTransactorSession) MintWhitelist() (*types.Transaction, error) {
	return _AvatarContract.Contract.MintWhitelist(&_AvatarContract.TransactOpts)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_AvatarContract *AvatarContractTransactor) OwnerMint(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "ownerMint", id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_AvatarContract *AvatarContractSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.OwnerMint(&_AvatarContract.TransactOpts, id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_AvatarContract *AvatarContractTransactorSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.OwnerMint(&_AvatarContract.TransactOpts, id)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AvatarContract *AvatarContractTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AvatarContract *AvatarContractSession) Pause() (*types.Transaction, error) {
	return _AvatarContract.Contract.Pause(&_AvatarContract.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AvatarContract *AvatarContractTransactorSession) Pause() (*types.Transaction, error) {
	return _AvatarContract.Contract.Pause(&_AvatarContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AvatarContract *AvatarContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AvatarContract *AvatarContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _AvatarContract.Contract.RenounceOwnership(&_AvatarContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AvatarContract *AvatarContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AvatarContract.Contract.RenounceOwnership(&_AvatarContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SafeTransferFrom(&_AvatarContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SafeTransferFrom(&_AvatarContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AvatarContract *AvatarContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AvatarContract *AvatarContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AvatarContract.Contract.SafeTransferFrom0(&_AvatarContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_AvatarContract *AvatarContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _AvatarContract.Contract.SafeTransferFrom0(&_AvatarContract.TransactOpts, from, to, tokenId, data)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_AvatarContract *AvatarContractTransactor) SetAlgo(opts *bind.TransactOpts, algo string) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "setAlgo", algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_AvatarContract *AvatarContractSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetAlgo(&_AvatarContract.TransactOpts, algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_AvatarContract *AvatarContractTransactorSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetAlgo(&_AvatarContract.TransactOpts, algo)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AvatarContract *AvatarContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AvatarContract *AvatarContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetApprovalForAll(&_AvatarContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_AvatarContract *AvatarContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetApprovalForAll(&_AvatarContract.TransactOpts, operator, approved)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_AvatarContract *AvatarContractTransactor) SetFee(opts *bind.TransactOpts, fee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "setFee", fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_AvatarContract *AvatarContractSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetFee(&_AvatarContract.TransactOpts, fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_AvatarContract *AvatarContractTransactorSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetFee(&_AvatarContract.TransactOpts, fee)
}

// SetWhitelistFee is a paid mutator transaction binding the contract method 0xc6283c38.
//
// Solidity: function setWhitelistFee(uint256 whitelistFee) returns()
func (_AvatarContract *AvatarContractTransactor) SetWhitelistFee(opts *bind.TransactOpts, whitelistFee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "setWhitelistFee", whitelistFee)
}

// SetWhitelistFee is a paid mutator transaction binding the contract method 0xc6283c38.
//
// Solidity: function setWhitelistFee(uint256 whitelistFee) returns()
func (_AvatarContract *AvatarContractSession) SetWhitelistFee(whitelistFee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetWhitelistFee(&_AvatarContract.TransactOpts, whitelistFee)
}

// SetWhitelistFee is a paid mutator transaction binding the contract method 0xc6283c38.
//
// Solidity: function setWhitelistFee(uint256 whitelistFee) returns()
func (_AvatarContract *AvatarContractTransactorSession) SetWhitelistFee(whitelistFee *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.SetWhitelistFee(&_AvatarContract.TransactOpts, whitelistFee)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.TransferFrom(&_AvatarContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_AvatarContract *AvatarContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.TransferFrom(&_AvatarContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AvatarContract *AvatarContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AvatarContract *AvatarContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.TransferOwnership(&_AvatarContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AvatarContract *AvatarContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AvatarContract.Contract.TransferOwnership(&_AvatarContract.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AvatarContract *AvatarContractTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AvatarContract *AvatarContractSession) Unpause() (*types.Transaction, error) {
	return _AvatarContract.Contract.Unpause(&_AvatarContract.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AvatarContract *AvatarContractTransactorSession) Unpause() (*types.Transaction, error) {
	return _AvatarContract.Contract.Unpause(&_AvatarContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_AvatarContract *AvatarContractTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _AvatarContract.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_AvatarContract *AvatarContractSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.Withdraw(&_AvatarContract.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_AvatarContract *AvatarContractTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _AvatarContract.Contract.Withdraw(&_AvatarContract.TransactOpts, amount)
}

// AvatarContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AvatarContract contract.
type AvatarContractApprovalIterator struct {
	Event *AvatarContractApproval // Event containing the contract specifics and raw log

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
func (it *AvatarContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractApproval)
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
		it.Event = new(AvatarContractApproval)
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
func (it *AvatarContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractApproval represents a Approval event raised by the AvatarContract contract.
type AvatarContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AvatarContract *AvatarContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*AvatarContractApprovalIterator, error) {

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

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AvatarContractApprovalIterator{contract: _AvatarContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_AvatarContract *AvatarContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AvatarContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractApproval)
				if err := _AvatarContract.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_AvatarContract *AvatarContractFilterer) ParseApproval(log types.Log) (*AvatarContractApproval, error) {
	event := new(AvatarContractApproval)
	if err := _AvatarContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the AvatarContract contract.
type AvatarContractApprovalForAllIterator struct {
	Event *AvatarContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *AvatarContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractApprovalForAll)
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
		it.Event = new(AvatarContractApprovalForAll)
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
func (it *AvatarContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractApprovalForAll represents a ApprovalForAll event raised by the AvatarContract contract.
type AvatarContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AvatarContract *AvatarContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*AvatarContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &AvatarContractApprovalForAllIterator{contract: _AvatarContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_AvatarContract *AvatarContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *AvatarContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractApprovalForAll)
				if err := _AvatarContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_AvatarContract *AvatarContractFilterer) ParseApprovalForAll(log types.Log) (*AvatarContractApprovalForAll, error) {
	event := new(AvatarContractApprovalForAll)
	if err := _AvatarContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AvatarContract contract.
type AvatarContractInitializedIterator struct {
	Event *AvatarContractInitialized // Event containing the contract specifics and raw log

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
func (it *AvatarContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractInitialized)
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
		it.Event = new(AvatarContractInitialized)
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
func (it *AvatarContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractInitialized represents a Initialized event raised by the AvatarContract contract.
type AvatarContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AvatarContract *AvatarContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*AvatarContractInitializedIterator, error) {

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AvatarContractInitializedIterator{contract: _AvatarContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AvatarContract *AvatarContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AvatarContractInitialized) (event.Subscription, error) {

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractInitialized)
				if err := _AvatarContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AvatarContract *AvatarContractFilterer) ParseInitialized(log types.Log) (*AvatarContractInitialized, error) {
	event := new(AvatarContractInitialized)
	if err := _AvatarContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AvatarContract contract.
type AvatarContractOwnershipTransferredIterator struct {
	Event *AvatarContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AvatarContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractOwnershipTransferred)
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
		it.Event = new(AvatarContractOwnershipTransferred)
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
func (it *AvatarContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractOwnershipTransferred represents a OwnershipTransferred event raised by the AvatarContract contract.
type AvatarContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AvatarContract *AvatarContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AvatarContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AvatarContractOwnershipTransferredIterator{contract: _AvatarContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AvatarContract *AvatarContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AvatarContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractOwnershipTransferred)
				if err := _AvatarContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AvatarContract *AvatarContractFilterer) ParseOwnershipTransferred(log types.Log) (*AvatarContractOwnershipTransferred, error) {
	event := new(AvatarContractOwnershipTransferred)
	if err := _AvatarContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the AvatarContract contract.
type AvatarContractPausedIterator struct {
	Event *AvatarContractPaused // Event containing the contract specifics and raw log

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
func (it *AvatarContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractPaused)
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
		it.Event = new(AvatarContractPaused)
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
func (it *AvatarContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractPaused represents a Paused event raised by the AvatarContract contract.
type AvatarContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AvatarContract *AvatarContractFilterer) FilterPaused(opts *bind.FilterOpts) (*AvatarContractPausedIterator, error) {

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AvatarContractPausedIterator{contract: _AvatarContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AvatarContract *AvatarContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AvatarContractPaused) (event.Subscription, error) {

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractPaused)
				if err := _AvatarContract.contract.UnpackLog(event, "Paused", log); err != nil {
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
// Solidity: event Paused(address account)
func (_AvatarContract *AvatarContractFilterer) ParsePaused(log types.Log) (*AvatarContractPaused, error) {
	event := new(AvatarContractPaused)
	if err := _AvatarContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractRequestFulfilledDataIterator is returned from FilterRequestFulfilledData and is used to iterate over the raw logs and unpacked data for RequestFulfilledData events raised by the AvatarContract contract.
type AvatarContractRequestFulfilledDataIterator struct {
	Event *AvatarContractRequestFulfilledData // Event containing the contract specifics and raw log

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
func (it *AvatarContractRequestFulfilledDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractRequestFulfilledData)
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
		it.Event = new(AvatarContractRequestFulfilledData)
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
func (it *AvatarContractRequestFulfilledDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractRequestFulfilledDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractRequestFulfilledData represents a RequestFulfilledData event raised by the AvatarContract contract.
type AvatarContractRequestFulfilledData struct {
	RequestId [32]byte
	Data      common.Hash
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestFulfilledData is a free log retrieval operation binding the contract event 0xc6ef1c5f8f954d0251e9620ce4908018f1c31f432eb9b44fd9efb2ca569453f1.
//
// Solidity: event RequestFulfilledData(bytes32 indexed requestId, bytes indexed data)
func (_AvatarContract *AvatarContractFilterer) FilterRequestFulfilledData(opts *bind.FilterOpts, requestId [][32]byte, data [][]byte) (*AvatarContractRequestFulfilledDataIterator, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "RequestFulfilledData", requestIdRule, dataRule)
	if err != nil {
		return nil, err
	}
	return &AvatarContractRequestFulfilledDataIterator{contract: _AvatarContract.contract, event: "RequestFulfilledData", logs: logs, sub: sub}, nil
}

// WatchRequestFulfilledData is a free log subscription operation binding the contract event 0xc6ef1c5f8f954d0251e9620ce4908018f1c31f432eb9b44fd9efb2ca569453f1.
//
// Solidity: event RequestFulfilledData(bytes32 indexed requestId, bytes indexed data)
func (_AvatarContract *AvatarContractFilterer) WatchRequestFulfilledData(opts *bind.WatchOpts, sink chan<- *AvatarContractRequestFulfilledData, requestId [][32]byte, data [][]byte) (event.Subscription, error) {

	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}
	var dataRule []interface{}
	for _, dataItem := range data {
		dataRule = append(dataRule, dataItem)
	}

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "RequestFulfilledData", requestIdRule, dataRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractRequestFulfilledData)
				if err := _AvatarContract.contract.UnpackLog(event, "RequestFulfilledData", log); err != nil {
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

// ParseRequestFulfilledData is a log parse operation binding the contract event 0xc6ef1c5f8f954d0251e9620ce4908018f1c31f432eb9b44fd9efb2ca569453f1.
//
// Solidity: event RequestFulfilledData(bytes32 indexed requestId, bytes indexed data)
func (_AvatarContract *AvatarContractFilterer) ParseRequestFulfilledData(log types.Log) (*AvatarContractRequestFulfilledData, error) {
	event := new(AvatarContractRequestFulfilledData)
	if err := _AvatarContract.contract.UnpackLog(event, "RequestFulfilledData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AvatarContract contract.
type AvatarContractTransferIterator struct {
	Event *AvatarContractTransfer // Event containing the contract specifics and raw log

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
func (it *AvatarContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractTransfer)
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
		it.Event = new(AvatarContractTransfer)
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
func (it *AvatarContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractTransfer represents a Transfer event raised by the AvatarContract contract.
type AvatarContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AvatarContract *AvatarContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*AvatarContractTransferIterator, error) {

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

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &AvatarContractTransferIterator{contract: _AvatarContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_AvatarContract *AvatarContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AvatarContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractTransfer)
				if err := _AvatarContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_AvatarContract *AvatarContractFilterer) ParseTransfer(log types.Log) (*AvatarContractTransfer, error) {
	event := new(AvatarContractTransfer)
	if err := _AvatarContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AvatarContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the AvatarContract contract.
type AvatarContractUnpausedIterator struct {
	Event *AvatarContractUnpaused // Event containing the contract specifics and raw log

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
func (it *AvatarContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AvatarContractUnpaused)
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
		it.Event = new(AvatarContractUnpaused)
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
func (it *AvatarContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AvatarContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AvatarContractUnpaused represents a Unpaused event raised by the AvatarContract contract.
type AvatarContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AvatarContract *AvatarContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AvatarContractUnpausedIterator, error) {

	logs, sub, err := _AvatarContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AvatarContractUnpausedIterator{contract: _AvatarContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AvatarContract *AvatarContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AvatarContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _AvatarContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AvatarContractUnpaused)
				if err := _AvatarContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
// Solidity: event Unpaused(address account)
func (_AvatarContract *AvatarContractFilterer) ParseUnpaused(log types.Log) (*AvatarContractUnpaused, error) {
	event := new(AvatarContractUnpaused)
	if err := _AvatarContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
