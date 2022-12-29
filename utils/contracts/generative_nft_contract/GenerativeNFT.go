// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_nft_contract

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

// NFTProjectProjectMinting is an auto generated low-level Go binding around an user-defined struct.
type NFTProjectProjectMinting struct {
	ProjectAddr     common.Address
	ProjectId       *big.Int
	MaxSupply       *big.Int
	Limit           *big.Int
	Index           *big.Int
	IndexReserve    *big.Int
	Creator         string
	MintPrice       *big.Int
	MintPriceAddr   common.Address
	Name            string
	MintingSchedule NFTProjectProjectMintingSchedule
	Reserves        []common.Address
	Royalty         *big.Int
}

// NFTProjectProjectMintingSchedule is an auto generated low-level Go binding around an user-defined struct.
type NFTProjectProjectMintingSchedule struct {
	InitBlockTime *big.Int
	OpeningTime   *big.Int
}

// GenerativeNftContractMetaData contains all meta data concerning the GenerativeNftContract contract.
var GenerativeNftContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"OPERATOR_FILTER_REGISTRY\",\"outputs\":[{\"internalType\":\"contractIOperatorFilterRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_nameCol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_project\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_projectAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_index\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_indexReserve\",\"type\":\"uint24\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"_initBlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_openingTime\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.ProjectMintingSchedule\",\"name\":\"_mintingSchedule\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_projectDataContextAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_randomizer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_royalty\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeDataContextAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeRandomizerAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"_projectAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_index\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_indexReserve\",\"type\":\"uint24\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"_initBlockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_openingTime\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.ProjectMintingSchedule\",\"name\":\"_mintingSchedule\",\"type\":\"tuple\"},{\"internalType\":\"address[]\",\"name\":\"_reserves\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.ProjectMinting\",\"name\":\"project\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"randomizer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"projectDataContextAddr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"disable\",\"type\":\"bool\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserveMint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"setStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenGenerativeURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeNftContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeNftContractMetaData.ABI instead.
var GenerativeNftContractABI = GenerativeNftContractMetaData.ABI

// GenerativeNftContract is an auto generated Go binding around an Ethereum contract.
type GenerativeNftContract struct {
	GenerativeNftContractCaller     // Read-only binding to the contract
	GenerativeNftContractTransactor // Write-only binding to the contract
	GenerativeNftContractFilterer   // Log filterer for contract events
}

// GenerativeNftContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeNftContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeNftContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeNftContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeNftContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeNftContractSession struct {
	Contract     *GenerativeNftContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GenerativeNftContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeNftContractCallerSession struct {
	Contract *GenerativeNftContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GenerativeNftContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeNftContractTransactorSession struct {
	Contract     *GenerativeNftContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GenerativeNftContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeNftContractRaw struct {
	Contract *GenerativeNftContract // Generic contract binding to access the raw methods on
}

// GenerativeNftContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeNftContractCallerRaw struct {
	Contract *GenerativeNftContractCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeNftContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeNftContractTransactorRaw struct {
	Contract *GenerativeNftContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeNftContract creates a new instance of GenerativeNftContract, bound to a specific deployed contract.
func NewGenerativeNftContract(address common.Address, backend bind.ContractBackend) (*GenerativeNftContract, error) {
	contract, err := bindGenerativeNftContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContract{GenerativeNftContractCaller: GenerativeNftContractCaller{contract: contract}, GenerativeNftContractTransactor: GenerativeNftContractTransactor{contract: contract}, GenerativeNftContractFilterer: GenerativeNftContractFilterer{contract: contract}}, nil
}

// NewGenerativeNftContractCaller creates a new read-only instance of GenerativeNftContract, bound to a specific deployed contract.
func NewGenerativeNftContractCaller(address common.Address, caller bind.ContractCaller) (*GenerativeNftContractCaller, error) {
	contract, err := bindGenerativeNftContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractCaller{contract: contract}, nil
}

// NewGenerativeNftContractTransactor creates a new write-only instance of GenerativeNftContract, bound to a specific deployed contract.
func NewGenerativeNftContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeNftContractTransactor, error) {
	contract, err := bindGenerativeNftContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractTransactor{contract: contract}, nil
}

// NewGenerativeNftContractFilterer creates a new log filterer instance of GenerativeNftContract, bound to a specific deployed contract.
func NewGenerativeNftContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeNftContractFilterer, error) {
	contract, err := bindGenerativeNftContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractFilterer{contract: contract}, nil
}

// bindGenerativeNftContract binds a generic wrapper to an already deployed contract.
func bindGenerativeNftContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeNftContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeNftContract *GenerativeNftContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeNftContract.Contract.GenerativeNftContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeNftContract *GenerativeNftContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.GenerativeNftContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeNftContract *GenerativeNftContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.GenerativeNftContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeNftContract *GenerativeNftContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeNftContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeNftContract *GenerativeNftContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeNftContract *GenerativeNftContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.contract.Transact(opts, method, params...)
}

// OPERATORFILTERREGISTRY is a free data retrieval call binding the contract method 0x41f43434.
//
// Solidity: function OPERATOR_FILTER_REGISTRY() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) OPERATORFILTERREGISTRY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "OPERATOR_FILTER_REGISTRY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPERATORFILTERREGISTRY is a free data retrieval call binding the contract method 0x41f43434.
//
// Solidity: function OPERATOR_FILTER_REGISTRY() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) OPERATORFILTERREGISTRY() (common.Address, error) {
	return _GenerativeNftContract.Contract.OPERATORFILTERREGISTRY(&_GenerativeNftContract.CallOpts)
}

// OPERATORFILTERREGISTRY is a free data retrieval call binding the contract method 0x41f43434.
//
// Solidity: function OPERATOR_FILTER_REGISTRY() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) OPERATORFILTERREGISTRY() (common.Address, error) {
	return _GenerativeNftContract.Contract.OPERATORFILTERREGISTRY(&_GenerativeNftContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) Admin() (common.Address, error) {
	return _GenerativeNftContract.Contract.Admin(&_GenerativeNftContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Admin() (common.Address, error) {
	return _GenerativeNftContract.Contract.Admin(&_GenerativeNftContract.CallOpts)
}

// NameCol is a free data retrieval call binding the contract method 0x452fa3d3.
//
// Solidity: function _nameCol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCaller) NameCol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_nameCol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NameCol is a free data retrieval call binding the contract method 0x452fa3d3.
//
// Solidity: function _nameCol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractSession) NameCol() (string, error) {
	return _GenerativeNftContract.Contract.NameCol(&_GenerativeNftContract.CallOpts)
}

// NameCol is a free data retrieval call binding the contract method 0x452fa3d3.
//
// Solidity: function _nameCol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) NameCol() (string, error) {
	return _GenerativeNftContract.Contract.NameCol(&_GenerativeNftContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) ParamsAddress() (common.Address, error) {
	return _GenerativeNftContract.Contract.ParamsAddress(&_GenerativeNftContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) ParamsAddress() (common.Address, error) {
	return _GenerativeNftContract.Contract.ParamsAddress(&_GenerativeNftContract.CallOpts)
}

// Project is a free data retrieval call binding the contract method 0x775bff9f.
//
// Solidity: function _project() view returns(address _projectAddr, uint256 _projectId, uint24 _maxSupply, uint24 _limit, uint24 _index, uint24 _indexReserve, string _creator, uint256 _mintPrice, address _mintPriceAddr, string _name, (uint256,uint256) _mintingSchedule, uint256 _royalty)
func (_GenerativeNftContract *GenerativeNftContractCaller) Project(opts *bind.CallOpts) (struct {
	ProjectAddr     common.Address
	ProjectId       *big.Int
	MaxSupply       *big.Int
	Limit           *big.Int
	Index           *big.Int
	IndexReserve    *big.Int
	Creator         string
	MintPrice       *big.Int
	MintPriceAddr   common.Address
	Name            string
	MintingSchedule NFTProjectProjectMintingSchedule
	Royalty         *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_project")

	outstruct := new(struct {
		ProjectAddr     common.Address
		ProjectId       *big.Int
		MaxSupply       *big.Int
		Limit           *big.Int
		Index           *big.Int
		IndexReserve    *big.Int
		Creator         string
		MintPrice       *big.Int
		MintPriceAddr   common.Address
		Name            string
		MintingSchedule NFTProjectProjectMintingSchedule
		Royalty         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ProjectAddr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ProjectId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MaxSupply = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Limit = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Index = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IndexReserve = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Creator = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.MintPrice = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.MintPriceAddr = *abi.ConvertType(out[8], new(common.Address)).(*common.Address)
	outstruct.Name = *abi.ConvertType(out[9], new(string)).(*string)
	outstruct.MintingSchedule = *abi.ConvertType(out[10], new(NFTProjectProjectMintingSchedule)).(*NFTProjectProjectMintingSchedule)
	outstruct.Royalty = *abi.ConvertType(out[11], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Project is a free data retrieval call binding the contract method 0x775bff9f.
//
// Solidity: function _project() view returns(address _projectAddr, uint256 _projectId, uint24 _maxSupply, uint24 _limit, uint24 _index, uint24 _indexReserve, string _creator, uint256 _mintPrice, address _mintPriceAddr, string _name, (uint256,uint256) _mintingSchedule, uint256 _royalty)
func (_GenerativeNftContract *GenerativeNftContractSession) Project() (struct {
	ProjectAddr     common.Address
	ProjectId       *big.Int
	MaxSupply       *big.Int
	Limit           *big.Int
	Index           *big.Int
	IndexReserve    *big.Int
	Creator         string
	MintPrice       *big.Int
	MintPriceAddr   common.Address
	Name            string
	MintingSchedule NFTProjectProjectMintingSchedule
	Royalty         *big.Int
}, error) {
	return _GenerativeNftContract.Contract.Project(&_GenerativeNftContract.CallOpts)
}

// Project is a free data retrieval call binding the contract method 0x775bff9f.
//
// Solidity: function _project() view returns(address _projectAddr, uint256 _projectId, uint24 _maxSupply, uint24 _limit, uint24 _index, uint24 _indexReserve, string _creator, uint256 _mintPrice, address _mintPriceAddr, string _name, (uint256,uint256) _mintingSchedule, uint256 _royalty)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Project() (struct {
	ProjectAddr     common.Address
	ProjectId       *big.Int
	MaxSupply       *big.Int
	Limit           *big.Int
	Index           *big.Int
	IndexReserve    *big.Int
	Creator         string
	MintPrice       *big.Int
	MintPriceAddr   common.Address
	Name            string
	MintingSchedule NFTProjectProjectMintingSchedule
	Royalty         *big.Int
}, error) {
	return _GenerativeNftContract.Contract.Project(&_GenerativeNftContract.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) ProjectDataContextAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_projectDataContextAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) ProjectDataContextAddr() (common.Address, error) {
	return _GenerativeNftContract.Contract.ProjectDataContextAddr(&_GenerativeNftContract.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) ProjectDataContextAddr() (common.Address, error) {
	return _GenerativeNftContract.Contract.ProjectDataContextAddr(&_GenerativeNftContract.CallOpts)
}

// Randomizer is a free data retrieval call binding the contract method 0xffc27c0e.
//
// Solidity: function _randomizer() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) Randomizer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_randomizer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Randomizer is a free data retrieval call binding the contract method 0xffc27c0e.
//
// Solidity: function _randomizer() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) Randomizer() (common.Address, error) {
	return _GenerativeNftContract.Contract.Randomizer(&_GenerativeNftContract.CallOpts)
}

// Randomizer is a free data retrieval call binding the contract method 0xffc27c0e.
//
// Solidity: function _randomizer() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Randomizer() (common.Address, error) {
	return _GenerativeNftContract.Contract.Randomizer(&_GenerativeNftContract.CallOpts)
}

// Royalty is a free data retrieval call binding the contract method 0x3b66b00a.
//
// Solidity: function _royalty() view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractCaller) Royalty(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "_royalty")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Royalty is a free data retrieval call binding the contract method 0x3b66b00a.
//
// Solidity: function _royalty() view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractSession) Royalty() (*big.Int, error) {
	return _GenerativeNftContract.Contract.Royalty(&_GenerativeNftContract.CallOpts)
}

// Royalty is a free data retrieval call binding the contract method 0x3b66b00a.
//
// Solidity: function _royalty() view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Royalty() (*big.Int, error) {
	return _GenerativeNftContract.Contract.Royalty(&_GenerativeNftContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeNftContract.Contract.BalanceOf(&_GenerativeNftContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeNftContract.Contract.BalanceOf(&_GenerativeNftContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNftContract.Contract.GetApproved(&_GenerativeNftContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNftContract.Contract.GetApproved(&_GenerativeNftContract.CallOpts, tokenId)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCaller) GetStatus(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "getStatus")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractSession) GetStatus() (bool, error) {
	return _GenerativeNftContract.Contract.GetStatus(&_GenerativeNftContract.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) GetStatus() (bool, error) {
	return _GenerativeNftContract.Contract.GetStatus(&_GenerativeNftContract.CallOpts)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeNftContract.Contract.IsApprovedForAll(&_GenerativeNftContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeNftContract.Contract.IsApprovedForAll(&_GenerativeNftContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractSession) Name() (string, error) {
	return _GenerativeNftContract.Contract.Name(&_GenerativeNftContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Name() (string, error) {
	return _GenerativeNftContract.Contract.Name(&_GenerativeNftContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) Owner() (common.Address, error) {
	return _GenerativeNftContract.Contract.Owner(&_GenerativeNftContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Owner() (common.Address, error) {
	return _GenerativeNftContract.Contract.Owner(&_GenerativeNftContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNftContract.Contract.OwnerOf(&_GenerativeNftContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeNftContract.Contract.OwnerOf(&_GenerativeNftContract.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractSession) Paused() (bool, error) {
	return _GenerativeNftContract.Contract.Paused(&_GenerativeNftContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Paused() (bool, error) {
	return _GenerativeNftContract.Contract.Paused(&_GenerativeNftContract.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeNftContract *GenerativeNftContractCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

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
func (_GenerativeNftContract *GenerativeNftContractSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeNftContract.Contract.RoyaltyInfo(&_GenerativeNftContract.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeNftContract.Contract.RoyaltyInfo(&_GenerativeNftContract.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeNftContract.Contract.SupportsInterface(&_GenerativeNftContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeNftContract.Contract.SupportsInterface(&_GenerativeNftContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractSession) Symbol() (string, error) {
	return _GenerativeNftContract.Contract.Symbol(&_GenerativeNftContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) Symbol() (string, error) {
	return _GenerativeNftContract.Contract.Symbol(&_GenerativeNftContract.CallOpts)
}

// TokenGenerativeURI is a free data retrieval call binding the contract method 0x10da88c5.
//
// Solidity: function tokenGenerativeURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCaller) TokenGenerativeURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "tokenGenerativeURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenGenerativeURI is a free data retrieval call binding the contract method 0x10da88c5.
//
// Solidity: function tokenGenerativeURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractSession) TokenGenerativeURI(tokenId *big.Int) (string, error) {
	return _GenerativeNftContract.Contract.TokenGenerativeURI(&_GenerativeNftContract.CallOpts, tokenId)
}

// TokenGenerativeURI is a free data retrieval call binding the contract method 0x10da88c5.
//
// Solidity: function tokenGenerativeURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) TokenGenerativeURI(tokenId *big.Int) (string, error) {
	return _GenerativeNftContract.Contract.TokenGenerativeURI(&_GenerativeNftContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeNftContract.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractSession) TokenURI(tokenId *big.Int) (string, error) {
	return _GenerativeNftContract.Contract.TokenURI(&_GenerativeNftContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_GenerativeNftContract *GenerativeNftContractCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _GenerativeNftContract.Contract.TokenURI(&_GenerativeNftContract.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Approve(&_GenerativeNftContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Approve(&_GenerativeNftContract.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeAdmin(&_GenerativeNftContract.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeAdmin(&_GenerativeNftContract.TransactOpts, newAdm)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) ChangeDataContextAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "changeDataContextAddr", newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeDataContextAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeDataContextAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) ChangeParamAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "changeParamAddr", newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeParamAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeParamAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) ChangeRandomizerAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "changeRandomizerAddr", newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeRandomizerAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ChangeRandomizerAddr(&_GenerativeNftContract.TransactOpts, newAddr)
}

// Init is a paid mutator transaction binding the contract method 0xfc237452.
//
// Solidity: function init((address,uint256,uint24,uint24,uint24,uint24,string,uint256,address,string,(uint256,uint256),address[],uint256) project, address admin, address paramsAddr, address randomizer, address projectDataContextAddr, bool disable) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) Init(opts *bind.TransactOpts, project NFTProjectProjectMinting, admin common.Address, paramsAddr common.Address, randomizer common.Address, projectDataContextAddr common.Address, disable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "init", project, admin, paramsAddr, randomizer, projectDataContextAddr, disable)
}

// Init is a paid mutator transaction binding the contract method 0xfc237452.
//
// Solidity: function init((address,uint256,uint24,uint24,uint24,uint24,string,uint256,address,string,(uint256,uint256),address[],uint256) project, address admin, address paramsAddr, address randomizer, address projectDataContextAddr, bool disable) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) Init(project NFTProjectProjectMinting, admin common.Address, paramsAddr common.Address, randomizer common.Address, projectDataContextAddr common.Address, disable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Init(&_GenerativeNftContract.TransactOpts, project, admin, paramsAddr, randomizer, projectDataContextAddr, disable)
}

// Init is a paid mutator transaction binding the contract method 0xfc237452.
//
// Solidity: function init((address,uint256,uint24,uint24,uint24,uint24,string,uint256,address,string,(uint256,uint256),address[],uint256) project, address admin, address paramsAddr, address randomizer, address projectDataContextAddr, bool disable) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) Init(project NFTProjectProjectMinting, admin common.Address, paramsAddr common.Address, randomizer common.Address, projectDataContextAddr common.Address, disable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Init(&_GenerativeNftContract.TransactOpts, project, admin, paramsAddr, randomizer, projectDataContextAddr, disable)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractSession) Mint() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Mint(&_GenerativeNftContract.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) Mint() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.Mint(&_GenerativeNftContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNftContract *GenerativeNftContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.RenounceOwnership(&_GenerativeNftContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.RenounceOwnership(&_GenerativeNftContract.TransactOpts)
}

// ReserveMint is a paid mutator transaction binding the contract method 0x21c8d676.
//
// Solidity: function reserveMint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractTransactor) ReserveMint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "reserveMint")
}

// ReserveMint is a paid mutator transaction binding the contract method 0x21c8d676.
//
// Solidity: function reserveMint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractSession) ReserveMint() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ReserveMint(&_GenerativeNftContract.TransactOpts)
}

// ReserveMint is a paid mutator transaction binding the contract method 0x21c8d676.
//
// Solidity: function reserveMint() payable returns(uint256 tokenId)
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) ReserveMint() (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.ReserveMint(&_GenerativeNftContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SafeTransferFrom(&_GenerativeNftContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SafeTransferFrom(&_GenerativeNftContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SafeTransferFrom0(&_GenerativeNftContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SafeTransferFrom0(&_GenerativeNftContract.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SetApprovalForAll(&_GenerativeNftContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SetApprovalForAll(&_GenerativeNftContract.TransactOpts, operator, approved)
}

// SetStatus is a paid mutator transaction binding the contract method 0x5c40f6f4.
//
// Solidity: function setStatus(bool enable) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) SetStatus(opts *bind.TransactOpts, enable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "setStatus", enable)
}

// SetStatus is a paid mutator transaction binding the contract method 0x5c40f6f4.
//
// Solidity: function setStatus(bool enable) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) SetStatus(enable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SetStatus(&_GenerativeNftContract.TransactOpts, enable)
}

// SetStatus is a paid mutator transaction binding the contract method 0x5c40f6f4.
//
// Solidity: function setStatus(bool enable) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) SetStatus(enable bool) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.SetStatus(&_GenerativeNftContract.TransactOpts, enable)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.TransferFrom(&_GenerativeNftContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.TransferFrom(&_GenerativeNftContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNftContract *GenerativeNftContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.TransferOwnership(&_GenerativeNftContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeNftContract *GenerativeNftContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeNftContract.Contract.TransferOwnership(&_GenerativeNftContract.TransactOpts, newOwner)
}

// GenerativeNftContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GenerativeNftContract contract.
type GenerativeNftContractApprovalIterator struct {
	Event *GenerativeNftContractApproval // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractApproval)
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
		it.Event = new(GenerativeNftContractApproval)
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
func (it *GenerativeNftContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractApproval represents a Approval event raised by the GenerativeNftContract contract.
type GenerativeNftContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GenerativeNftContractApprovalIterator, error) {

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

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractApprovalIterator{contract: _GenerativeNftContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractApproval)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParseApproval(log types.Log) (*GenerativeNftContractApproval, error) {
	event := new(GenerativeNftContractApproval)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GenerativeNftContract contract.
type GenerativeNftContractApprovalForAllIterator struct {
	Event *GenerativeNftContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractApprovalForAll)
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
		it.Event = new(GenerativeNftContractApprovalForAll)
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
func (it *GenerativeNftContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractApprovalForAll represents a ApprovalForAll event raised by the GenerativeNftContract contract.
type GenerativeNftContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GenerativeNftContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractApprovalForAllIterator{contract: _GenerativeNftContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractApprovalForAll)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParseApprovalForAll(log types.Log) (*GenerativeNftContractApprovalForAll, error) {
	event := new(GenerativeNftContractApprovalForAll)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GenerativeNftContract contract.
type GenerativeNftContractOwnershipTransferredIterator struct {
	Event *GenerativeNftContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractOwnershipTransferred)
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
		it.Event = new(GenerativeNftContractOwnershipTransferred)
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
func (it *GenerativeNftContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractOwnershipTransferred represents a OwnershipTransferred event raised by the GenerativeNftContract contract.
type GenerativeNftContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GenerativeNftContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractOwnershipTransferredIterator{contract: _GenerativeNftContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractOwnershipTransferred)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParseOwnershipTransferred(log types.Log) (*GenerativeNftContractOwnershipTransferred, error) {
	event := new(GenerativeNftContractOwnershipTransferred)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GenerativeNftContract contract.
type GenerativeNftContractPausedIterator struct {
	Event *GenerativeNftContractPaused // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractPaused)
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
		it.Event = new(GenerativeNftContractPaused)
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
func (it *GenerativeNftContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractPaused represents a Paused event raised by the GenerativeNftContract contract.
type GenerativeNftContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterPaused(opts *bind.FilterOpts) (*GenerativeNftContractPausedIterator, error) {

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractPausedIterator{contract: _GenerativeNftContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractPaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractPaused)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParsePaused(log types.Log) (*GenerativeNftContractPaused, error) {
	event := new(GenerativeNftContractPaused)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GenerativeNftContract contract.
type GenerativeNftContractTransferIterator struct {
	Event *GenerativeNftContractTransfer // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractTransfer)
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
		it.Event = new(GenerativeNftContractTransfer)
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
func (it *GenerativeNftContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractTransfer represents a Transfer event raised by the GenerativeNftContract contract.
type GenerativeNftContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GenerativeNftContractTransferIterator, error) {

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

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractTransferIterator{contract: _GenerativeNftContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractTransfer)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParseTransfer(log types.Log) (*GenerativeNftContractTransfer, error) {
	event := new(GenerativeNftContractTransfer)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeNftContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GenerativeNftContract contract.
type GenerativeNftContractUnpausedIterator struct {
	Event *GenerativeNftContractUnpaused // Event containing the contract specifics and raw log

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
func (it *GenerativeNftContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeNftContractUnpaused)
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
		it.Event = new(GenerativeNftContractUnpaused)
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
func (it *GenerativeNftContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeNftContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeNftContractUnpaused represents a Unpaused event raised by the GenerativeNftContract contract.
type GenerativeNftContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeNftContract *GenerativeNftContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GenerativeNftContractUnpausedIterator, error) {

	logs, sub, err := _GenerativeNftContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GenerativeNftContractUnpausedIterator{contract: _GenerativeNftContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeNftContract *GenerativeNftContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GenerativeNftContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeNftContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeNftContractUnpaused)
				if err := _GenerativeNftContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_GenerativeNftContract *GenerativeNftContractFilterer) ParseUnpaused(log types.Log) (*GenerativeNftContractUnpaused, error) {
	event := new(GenerativeNftContractUnpaused)
	if err := _GenerativeNftContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
