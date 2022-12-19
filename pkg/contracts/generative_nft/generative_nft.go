// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_nft

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

// GenerativeNftMetaData contains all meta data concerning the GenerativeNft contract.
var GenerativeNftMetaData = &bind.MetaData{
	ABI: `[{"inputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"address","name":"operator","type":"address"}],"name":"OperatorNotAllowed","type":"error"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"approved","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":false,"internalType":"bool","name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":true,"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"inputs":[],"name":"OPERATOR_FILTER_REGISTRY","outputs":[{"internalType":"contract IOperatorFilterRegistry","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_admin","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_nameCol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_paramsAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_project","outputs":[{"internalType":"address","name":"_projectAddr","type":"address"},{"internalType":"uint256","name":"_projectId","type":"uint256"},{"internalType":"uint24","name":"_maxSupply","type":"uint24"},{"internalType":"uint24","name":"_limit","type":"uint24"},{"internalType":"uint24","name":"_index","type":"uint24"},{"internalType":"uint24","name":"_indexReserve","type":"uint24"},{"internalType":"string","name":"_creator","type":"string"},{"internalType":"address","name":"_creatorAddr","type":"address"},{"internalType":"uint256","name":"_mintPrice","type":"uint256"},{"internalType":"address","name":"_mintPriceAddr","type":"address"},{"internalType":"string","name":"_name","type":"string"},{"components":[{"internalType":"uint256","name":"_initBlockTime","type":"uint256"},{"internalType":"uint256","name":"_openingTime","type":"uint256"}],"internalType":"struct NFTProject.ProjectMintingSchedule","name":"_mintingSchedule","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_projectDataContextAddr","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_randomizer","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"_reserves","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_royalty","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"_symbolCol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"newAdm","type":"address"}],"name":"changeAdmin","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newAddr","type":"address"}],"name":"changeDataContextAddr","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newAddr","type":"address"}],"name":"changeParamAddr","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newAddr","type":"address"}],"name":"changeRandomizerAddr","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getStatus","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"_projectAddr","type":"address"},{"internalType":"uint256","name":"_projectId","type":"uint256"},{"internalType":"uint24","name":"_maxSupply","type":"uint24"},{"internalType":"uint24","name":"_limit","type":"uint24"},{"internalType":"uint24","name":"_index","type":"uint24"},{"internalType":"uint24","name":"_indexReserve","type":"uint24"},{"internalType":"string","name":"_creator","type":"string"},{"internalType":"address","name":"_creatorAddr","type":"address"},{"internalType":"uint256","name":"_mintPrice","type":"uint256"},{"internalType":"address","name":"_mintPriceAddr","type":"address"},{"internalType":"string","name":"_name","type":"string"},{"components":[{"internalType":"uint256","name":"_initBlockTime","type":"uint256"},{"internalType":"uint256","name":"_openingTime","type":"uint256"}],"internalType":"struct NFTProject.ProjectMintingSchedule","name":"_mintingSchedule","type":"tuple"}],"internalType":"struct NFTProject.ProjectMinting","name":"project","type":"tuple"},{"internalType":"address","name":"admin","type":"address"},{"internalType":"address","name":"paramsAddr","type":"address"},{"internalType":"address","name":"randomizer","type":"address"},{"internalType":"address","name":"projectDataContextAddr","type":"address"},{"internalType":"address[]","name":"reserves","type":"address[]"},{"internalType":"bool","name":"disable","type":"bool"},{"internalType":"uint256","name":"royalty","type":"uint256"}],"name":"init","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"mint","outputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"reserveMint","outputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"uint256","name":"_tokenId","type":"uint256"},{"internalType":"uint256","name":"_salePrice","type":"uint256"}],"name":"royaltyInfo","outputs":[{"internalType":"address","name":"receiver","type":"address"},{"internalType":"uint256","name":"royaltyAmount","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"bool","name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bool","name":"enable","type":"bool"}],"name":"setStatus","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes4","name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenGenerativeURI","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenIdToHash","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`,
}

// GenerativeNftABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeNftMetaData.ABI instead.
var GenerativeNftABI = GenerativeNftMetaData.ABI

// GenerativeNft is an auto generated Go binding around an Ethereum contract.
type GenerativeNft struct {
	GenerativeNftCaller     // Read-only binding to the contract
	GenerativeNftTransactor // Write-only binding to the contract
	GenerativeNftFilterer   // Log filterer for contract events
}

// GenerativeNftCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeNftCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeNftTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeNftFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeNftSession struct {
	Contract     *GenerativeNft    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GenerativeNftCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeNftCallerSession struct {
	Contract *GenerativeNftCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GenerativeNftTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeNftTransactorSession struct {
	Contract     *GenerativeNftTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GenerativeNftRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeNftRaw struct {
	Contract *GenerativeNft // Generic contract binding to access the raw methods on
}

// GenerativeNftCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeNftCallerRaw struct {
	Contract *GenerativeNftCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeNftTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeNftTransactorRaw struct {
	Contract *GenerativeNftTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeNft creates a new instance of GenerativeNft, bound to a specific deployed contract.
func NewGenerativeNft(address common.Address, backend bind.ContractBackend) (*GenerativeNft, error) {
	contract, err := bindGenerativeNft(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeNft{GenerativeNftCaller: GenerativeNftCaller{contract: contract}, GenerativeNftTransactor: GenerativeNftTransactor{contract: contract}, GenerativeNftFilterer: GenerativeNftFilterer{contract: contract}}, nil
}

// NewGenerativeNftCaller creates a new read-only instance of GenerativeNft, bound to a specific deployed contract.
func NewGenerativeNftCaller(address common.Address, caller bind.ContractCaller) (*GenerativeNftCaller, error) {
	contract, err := bindGenerativeNft(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftCaller{contract: contract}, nil
}

// NewGenerativeNftTransactor creates a new write-only instance of GenerativeNft, bound to a specific deployed contract.
func NewGenerativeNftTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeNftTransactor, error) {
	contract, err := bindGenerativeNft(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftTransactor{contract: contract}, nil
}

// NewGenerativeNftFilterer creates a new log filterer instance of GenerativeNft, bound to a specific deployed contract.
func NewGenerativeNftFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeNftFilterer, error) {
	contract, err := bindGenerativeNft(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftFilterer{contract: contract}, nil
}

// bindGenerativeNft binds a generic wrapper to an already deployed contract.
func bindGenerativeNft(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeNftMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeNft *GenerativeNftRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeNft.Contract.GenerativeNftCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeNft *GenerativeNftRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNft.Contract.GenerativeNftTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeNft *GenerativeNftRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeNft.Contract.GenerativeNftTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeNft *GenerativeNftCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeNft.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeNft *GenerativeNftTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNft.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeNft *GenerativeNftTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeNft.Contract.contract.Transact(opts, method, params...)
}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_GenerativeNft *GenerativeNftCaller) Algorithm(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "algorithm")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_GenerativeNft *GenerativeNftSession) Algorithm() (string, error) {
	return _GenerativeNft.Contract.Algorithm(&_GenerativeNft.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) Algorithm() (string, error) {
	return _GenerativeNft.Contract.Algorithm(&_GenerativeNft.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNft *GenerativeNftCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNft *GenerativeNftSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeNft.Contract.BalanceOf(&_GenerativeNft.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNft *GenerativeNftCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeNft.Contract.BalanceOf(&_GenerativeNft.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNft.Contract.GetApproved(&_GenerativeNft.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNft.Contract.GetApproved(&_GenerativeNft.CallOpts, tokenId)
}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_GenerativeNft *GenerativeNftCaller) GetPalette(opts *bind.CallOpts, id *big.Int) ([4]string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getPalette", id)

	if err != nil {
		return *new([4]string), err
	}

	out0 := *abi.ConvertType(out[0], new([4]string)).(*[4]string)

	return out0, err

}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_GenerativeNft *GenerativeNftSession) GetPalette(id *big.Int) ([4]string, error) {
	return _GenerativeNft.Contract.GetPalette(&_GenerativeNft.CallOpts, id)
}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_GenerativeNft *GenerativeNftCallerSession) GetPalette(id *big.Int) ([4]string, error) {
	return _GenerativeNft.Contract.GetPalette(&_GenerativeNft.CallOpts, id)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_GenerativeNft *GenerativeNftCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) ([4]string, string, string, string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new([4]string), *new(string), *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([4]string)).(*[4]string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)

	return out0, out1, out2, out3, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_GenerativeNft *GenerativeNftSession) GetParamValues(tokenId *big.Int) ([4]string, string, string, string, error) {
	return _GenerativeNft.Contract.GetParamValues(&_GenerativeNft.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_GenerativeNft *GenerativeNftCallerSession) GetParamValues(tokenId *big.Int) ([4]string, string, string, string, error) {
	return _GenerativeNft.Contract.GetParamValues(&_GenerativeNft.CallOpts, tokenId)
}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCaller) GetShape(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getShape", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftSession) GetShape(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetShape(&_GenerativeNft.CallOpts, id)
}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) GetShape(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetShape(&_GenerativeNft.CallOpts, id)
}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCaller) GetSize(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getSize", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftSession) GetSize(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetSize(&_GenerativeNft.CallOpts, id)
}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) GetSize(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetSize(&_GenerativeNft.CallOpts, id)
}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCaller) GetSurface(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "getSurface", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftSession) GetSurface(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetSurface(&_GenerativeNft.CallOpts, id)
}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) GetSurface(id *big.Int) (string, error) {
	return _GenerativeNft.Contract.GetSurface(&_GenerativeNft.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNft *GenerativeNftCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNft *GenerativeNftSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeNft.Contract.IsApprovedForAll(&_GenerativeNft.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNft *GenerativeNftCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeNft.Contract.IsApprovedForAll(&_GenerativeNft.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNft *GenerativeNftCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNft *GenerativeNftSession) Name() (string, error) {
	return _GenerativeNft.Contract.Name(&_GenerativeNft.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) Name() (string, error) {
	return _GenerativeNft.Contract.Name(&_GenerativeNft.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNft *GenerativeNftCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNft *GenerativeNftSession) Owner() (common.Address, error) {
	return _GenerativeNft.Contract.Owner(&_GenerativeNft.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNft *GenerativeNftCallerSession) Owner() (common.Address, error) {
	return _GenerativeNft.Contract.Owner(&_GenerativeNft.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNft.Contract.OwnerOf(&_GenerativeNft.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNft *GenerativeNftCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNft.Contract.OwnerOf(&_GenerativeNft.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNft *GenerativeNftCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNft *GenerativeNftSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeNft.Contract.SupportsInterface(&_GenerativeNft.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNft *GenerativeNftCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeNft.Contract.SupportsInterface(&_GenerativeNft.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNft *GenerativeNftCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNft *GenerativeNftSession) Symbol() (string, error) {
	return _GenerativeNft.Contract.Symbol(&_GenerativeNft.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) Symbol() (string, error) {
	return _GenerativeNft.Contract.Symbol(&_GenerativeNft.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNft *GenerativeNftCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNft.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNft *GenerativeNftSession) TokenURI(tokenId *big.Int) (string, error) {
	return _GenerativeNft.Contract.TokenURI(&_GenerativeNft.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNft *GenerativeNftCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _GenerativeNft.Contract.TokenURI(&_GenerativeNft.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.Approve(&_GenerativeNft.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.Approve(&_GenerativeNft.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_GenerativeNft *GenerativeNftTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_GenerativeNft *GenerativeNftSession) Mint() (*types.Transaction, error) {
	return _GenerativeNft.Contract.Mint(&_GenerativeNft.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_GenerativeNft *GenerativeNftTransactorSession) Mint() (*types.Transaction, error) {
	return _GenerativeNft.Contract.Mint(&_GenerativeNft.TransactOpts)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_GenerativeNft *GenerativeNftTransactor) OwnerMint(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "ownerMint", id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_GenerativeNft *GenerativeNftSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.OwnerMint(&_GenerativeNft.TransactOpts, id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.OwnerMint(&_GenerativeNft.TransactOpts, id)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNft *GenerativeNftTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNft *GenerativeNftSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeNft.Contract.RenounceOwnership(&_GenerativeNft.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNft *GenerativeNftTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeNft.Contract.RenounceOwnership(&_GenerativeNft.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SafeTransferFrom(&_GenerativeNft.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SafeTransferFrom(&_GenerativeNft.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNft *GenerativeNftTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNft *GenerativeNftSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SafeTransferFrom0(&_GenerativeNft.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SafeTransferFrom0(&_GenerativeNft.TransactOpts, from, to, tokenId, data)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_GenerativeNft *GenerativeNftTransactor) SetAlgo(opts *bind.TransactOpts, algo string) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "setAlgo", algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_GenerativeNft *GenerativeNftSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SetAlgo(&_GenerativeNft.TransactOpts, algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SetAlgo(&_GenerativeNft.TransactOpts, algo)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNft *GenerativeNftTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNft *GenerativeNftSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SetApprovalForAll(&_GenerativeNft.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNft.Contract.SetApprovalForAll(&_GenerativeNft.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.TransferFrom(&_GenerativeNft.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNft.Contract.TransferFrom(&_GenerativeNft.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNft *GenerativeNftTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNft.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNft *GenerativeNftSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNft.Contract.TransferOwnership(&_GenerativeNft.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNft *GenerativeNftTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNft.Contract.TransferOwnership(&_GenerativeNft.TransactOpts, newOwner)
}

// GenerativeNftApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GenerativeNft contract.
type GenerativeNftApprovalIterator struct {
	Event *GenerativeNftApproval // Event containing the contract specifics and raw log

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
func (it *GenerativeNftApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftApproval)
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
		it.Event = new(GenerativeNftApproval)
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
func (it *GenerativeNftApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftApproval represents a Approval event raised by the GenerativeNft contract.
type GenerativeNftApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeNft *GenerativeNftFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GenerativeNftApprovalIterator, error) {

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

	logs, sub, err := _GenerativeNft.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftApprovalIterator{contract: _GenerativeNft.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeNft *GenerativeNftFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GenerativeNftApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeNft.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftApproval)
				if err := _GenerativeNft.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_GenerativeNft *GenerativeNftFilterer) ParseApproval(log types.Log) (*GenerativeNftApproval, error) {
	event := new(GenerativeNftApproval)
	if err := _GenerativeNft.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GenerativeNft contract.
type GenerativeNftApprovalForAllIterator struct {
	Event *GenerativeNftApprovalForAll // Event containing the contract specifics and raw log

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
func (it *GenerativeNftApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftApprovalForAll)
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
		it.Event = new(GenerativeNftApprovalForAll)
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
func (it *GenerativeNftApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftApprovalForAll represents a ApprovalForAll event raised by the GenerativeNft contract.
type GenerativeNftApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeNft *GenerativeNftFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GenerativeNftApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeNft.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftApprovalForAllIterator{contract: _GenerativeNft.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeNft *GenerativeNftFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GenerativeNftApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeNft.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftApprovalForAll)
				if err := _GenerativeNft.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_GenerativeNft *GenerativeNftFilterer) ParseApprovalForAll(log types.Log) (*GenerativeNftApprovalForAll, error) {
	event := new(GenerativeNftApprovalForAll)
	if err := _GenerativeNft.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GenerativeNft contract.
type GenerativeNftOwnershipTransferredIterator struct {
	Event *GenerativeNftOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GenerativeNftOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftOwnershipTransferred)
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
		it.Event = new(GenerativeNftOwnershipTransferred)
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
func (it *GenerativeNftOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftOwnershipTransferred represents a OwnershipTransferred event raised by the GenerativeNft contract.
type GenerativeNftOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeNft *GenerativeNftFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GenerativeNftOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeNft.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftOwnershipTransferredIterator{contract: _GenerativeNft.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeNft *GenerativeNftFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GenerativeNftOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeNft.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftOwnershipTransferred)
				if err := _GenerativeNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GenerativeNft *GenerativeNftFilterer) ParseOwnershipTransferred(log types.Log) (*GenerativeNftOwnershipTransferred, error) {
	event := new(GenerativeNftOwnershipTransferred)
	if err := _GenerativeNft.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GenerativeNft contract.
type GenerativeNftTransferIterator struct {
	Event *GenerativeNftTransfer // Event containing the contract specifics and raw log

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
func (it *GenerativeNftTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftTransfer)
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
		it.Event = new(GenerativeNftTransfer)
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
func (it *GenerativeNftTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftTransfer represents a Transfer event raised by the GenerativeNft contract.
type GenerativeNftTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeNft *GenerativeNftFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GenerativeNftTransferIterator, error) {

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

	logs, sub, err := _GenerativeNft.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftTransferIterator{contract: _GenerativeNft.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeNft *GenerativeNftFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GenerativeNftTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeNft.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftTransfer)
				if err := _GenerativeNft.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_GenerativeNft *GenerativeNftFilterer) ParseTransfer(log types.Log) (*GenerativeNftTransfer, error) {
	event := new(GenerativeNftTransfer)
	if err := _GenerativeNft.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
