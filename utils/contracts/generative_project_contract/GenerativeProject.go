// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_project_contract

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

// NFTProjectProject is an auto generated low-level Go binding around an user-defined struct.
type NFTProjectProject struct {
	MaxSupply     *big.Int
	Limit         *big.Int
	MintPrice     *big.Int
	MintPriceAddr common.Address
	Name          string
	Creator       string
	CreatorAddr   common.Address
	License       string
	Desc          string
	Image         string
	Social        NFTProjectProjectSocial
	ScriptType    []string
	Scripts       []string
	Styles        string
	CompleteTime  *big.Int
	GenNFTAddr    common.Address
	ItemDesc      string
	Reserves      []common.Address
	Royalty       *big.Int
}

// NFTProjectProjectSocial is an auto generated low-level Go binding around an user-defined struct.
type NFTProjectProjectSocial struct {
	Web       string
	Twitter   string
	Discord   string
	Medium    string
	Instagram string
}

// GenerativeProjectContractMetaData contains all meta data concerning the GenerativeProjectContract contract.
var GenerativeProjectContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_projectDataContextAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_randomizerAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_script\",\"type\":\"string\"}],\"name\":\"addProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeDataContextAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeRandomizerAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"completeProject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scriptIndex\",\"type\":\"uint256\"}],\"name\":\"deleteProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"randomizerAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"projectDataContextAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_creatorAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_license\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_desc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_image\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"_social\",\"type\":\"tuple\"},{\"internalType\":\"string[]\",\"name\":\"_scriptType\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_scripts\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"_styles\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_completeTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_genNFTAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_itemDesc\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"_reserves\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.Project\",\"name\":\"project\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"disable\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"openingTime\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"projectDetails\",\"outputs\":[{\"components\":[{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_creatorAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_license\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_desc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_image\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"_social\",\"type\":\"tuple\"},{\"internalType\":\"string[]\",\"name\":\"_scriptType\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_scripts\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"_styles\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_completeTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_genNFTAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_itemDesc\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"_reserves\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.Project\",\"name\":\"project\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"projectStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"setProjectStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"creatorName\",\"type\":\"string\"}],\"name\":\"updateProjectCreatorName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"license\",\"type\":\"string\"}],\"name\":\"updateProjectLicense\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"projectName\",\"type\":\"string\"}],\"name\":\"updateProjectName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"updateProjectPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scriptIndex\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"script\",\"type\":\"string\"}],\"name\":\"updateProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"scriptType\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"updateProjectScriptType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"updateProjectSocial\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"erc20Addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeProjectContractABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeProjectContractMetaData.ABI instead.
var GenerativeProjectContractABI = GenerativeProjectContractMetaData.ABI

// GenerativeProjectContract is an auto generated Go binding around an Ethereum contract.
type GenerativeProjectContract struct {
	GenerativeProjectContractCaller     // Read-only binding to the contract
	GenerativeProjectContractTransactor // Write-only binding to the contract
	GenerativeProjectContractFilterer   // Log filterer for contract events
}

// GenerativeProjectContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeProjectContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeProjectContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeProjectContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeProjectContractSession struct {
	Contract     *GenerativeProjectContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// GenerativeProjectContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeProjectContractCallerSession struct {
	Contract *GenerativeProjectContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// GenerativeProjectContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeProjectContractTransactorSession struct {
	Contract     *GenerativeProjectContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// GenerativeProjectContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeProjectContractRaw struct {
	Contract *GenerativeProjectContract // Generic contract binding to access the raw methods on
}

// GenerativeProjectContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeProjectContractCallerRaw struct {
	Contract *GenerativeProjectContractCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeProjectContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeProjectContractTransactorRaw struct {
	Contract *GenerativeProjectContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeProjectContract creates a new instance of GenerativeProjectContract, bound to a specific deployed contract.
func NewGenerativeProjectContract(address common.Address, backend bind.ContractBackend) (*GenerativeProjectContract, error) {
	contract, err := bindGenerativeProjectContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContract{GenerativeProjectContractCaller: GenerativeProjectContractCaller{contract: contract}, GenerativeProjectContractTransactor: GenerativeProjectContractTransactor{contract: contract}, GenerativeProjectContractFilterer: GenerativeProjectContractFilterer{contract: contract}}, nil
}

// NewGenerativeProjectContractCaller creates a new read-only instance of GenerativeProjectContract, bound to a specific deployed contract.
func NewGenerativeProjectContractCaller(address common.Address, caller bind.ContractCaller) (*GenerativeProjectContractCaller, error) {
	contract, err := bindGenerativeProjectContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractCaller{contract: contract}, nil
}

// NewGenerativeProjectContractTransactor creates a new write-only instance of GenerativeProjectContract, bound to a specific deployed contract.
func NewGenerativeProjectContractTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeProjectContractTransactor, error) {
	contract, err := bindGenerativeProjectContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractTransactor{contract: contract}, nil
}

// NewGenerativeProjectContractFilterer creates a new log filterer instance of GenerativeProjectContract, bound to a specific deployed contract.
func NewGenerativeProjectContractFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeProjectContractFilterer, error) {
	contract, err := bindGenerativeProjectContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractFilterer{contract: contract}, nil
}

// bindGenerativeProjectContract binds a generic wrapper to an already deployed contract.
func bindGenerativeProjectContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeProjectContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeProjectContract *GenerativeProjectContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeProjectContract.Contract.GenerativeProjectContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeProjectContract *GenerativeProjectContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.GenerativeProjectContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeProjectContract *GenerativeProjectContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.GenerativeProjectContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeProjectContract *GenerativeProjectContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeProjectContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeProjectContract *GenerativeProjectContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeProjectContract *GenerativeProjectContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Admin() (common.Address, error) {
	return _GenerativeProjectContract.Contract.Admin(&_GenerativeProjectContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) Admin() (common.Address, error) {
	return _GenerativeProjectContract.Contract.Admin(&_GenerativeProjectContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) ParamsAddress() (common.Address, error) {
	return _GenerativeProjectContract.Contract.ParamsAddress(&_GenerativeProjectContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) ParamsAddress() (common.Address, error) {
	return _GenerativeProjectContract.Contract.ParamsAddress(&_GenerativeProjectContract.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) ProjectDataContextAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "_projectDataContextAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) ProjectDataContextAddr() (common.Address, error) {
	return _GenerativeProjectContract.Contract.ProjectDataContextAddr(&_GenerativeProjectContract.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) ProjectDataContextAddr() (common.Address, error) {
	return _GenerativeProjectContract.Contract.ProjectDataContextAddr(&_GenerativeProjectContract.CallOpts)
}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) RandomizerAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "_randomizerAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) RandomizerAddr() (common.Address, error) {
	return _GenerativeProjectContract.Contract.RandomizerAddr(&_GenerativeProjectContract.CallOpts)
}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) RandomizerAddr() (common.Address, error) {
	return _GenerativeProjectContract.Contract.RandomizerAddr(&_GenerativeProjectContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeProjectContract.Contract.BalanceOf(&_GenerativeProjectContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeProjectContract.Contract.BalanceOf(&_GenerativeProjectContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeProjectContract.Contract.GetApproved(&_GenerativeProjectContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeProjectContract.Contract.GetApproved(&_GenerativeProjectContract.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeProjectContract.Contract.IsApprovedForAll(&_GenerativeProjectContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeProjectContract.Contract.IsApprovedForAll(&_GenerativeProjectContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Name() (string, error) {
	return _GenerativeProjectContract.Contract.Name(&_GenerativeProjectContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) Name() (string, error) {
	return _GenerativeProjectContract.Contract.Name(&_GenerativeProjectContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Owner() (common.Address, error) {
	return _GenerativeProjectContract.Contract.Owner(&_GenerativeProjectContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) Owner() (common.Address, error) {
	return _GenerativeProjectContract.Contract.Owner(&_GenerativeProjectContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeProjectContract.Contract.OwnerOf(&_GenerativeProjectContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeProjectContract.Contract.OwnerOf(&_GenerativeProjectContract.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Paused() (bool, error) {
	return _GenerativeProjectContract.Contract.Paused(&_GenerativeProjectContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) Paused() (bool, error) {
	return _GenerativeProjectContract.Contract.Paused(&_GenerativeProjectContract.CallOpts)
}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) ProjectDetails(opts *bind.CallOpts, projectId *big.Int) (NFTProjectProject, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "projectDetails", projectId)

	if err != nil {
		return *new(NFTProjectProject), err
	}

	out0 := *abi.ConvertType(out[0], new(NFTProjectProject)).(*NFTProjectProject)

	return out0, err

}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_GenerativeProjectContract *GenerativeProjectContractSession) ProjectDetails(projectId *big.Int) (NFTProjectProject, error) {
	return _GenerativeProjectContract.Contract.ProjectDetails(&_GenerativeProjectContract.CallOpts, projectId)
}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) ProjectDetails(projectId *big.Int) (NFTProjectProject, error) {
	return _GenerativeProjectContract.Contract.ProjectDetails(&_GenerativeProjectContract.CallOpts, projectId)
}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) ProjectStatus(opts *bind.CallOpts, projectId *big.Int) (bool, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "projectStatus", projectId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_GenerativeProjectContract *GenerativeProjectContractSession) ProjectStatus(projectId *big.Int) (bool, error) {
	return _GenerativeProjectContract.Contract.ProjectStatus(&_GenerativeProjectContract.CallOpts, projectId)
}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) ProjectStatus(projectId *big.Int) (bool, error) {
	return _GenerativeProjectContract.Contract.ProjectStatus(&_GenerativeProjectContract.CallOpts, projectId)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 projectId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) RoyaltyInfo(opts *bind.CallOpts, projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "royaltyInfo", projectId, _salePrice)

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
// Solidity: function royaltyInfo(uint256 projectId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeProjectContract *GenerativeProjectContractSession) RoyaltyInfo(projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeProjectContract.Contract.RoyaltyInfo(&_GenerativeProjectContract.CallOpts, projectId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 projectId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) RoyaltyInfo(projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeProjectContract.Contract.RoyaltyInfo(&_GenerativeProjectContract.CallOpts, projectId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeProjectContract.Contract.SupportsInterface(&_GenerativeProjectContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeProjectContract.Contract.SupportsInterface(&_GenerativeProjectContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Symbol() (string, error) {
	return _GenerativeProjectContract.Contract.Symbol(&_GenerativeProjectContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) Symbol() (string, error) {
	return _GenerativeProjectContract.Contract.Symbol(&_GenerativeProjectContract.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectContract *GenerativeProjectContractCaller) TokenURI(opts *bind.CallOpts, projectId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeProjectContract.contract.Call(opts, &out, "tokenURI", projectId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectContract *GenerativeProjectContractSession) TokenURI(projectId *big.Int) (string, error) {
	return _GenerativeProjectContract.Contract.TokenURI(&_GenerativeProjectContract.CallOpts, projectId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectContract *GenerativeProjectContractCallerSession) TokenURI(projectId *big.Int) (string, error) {
	return _GenerativeProjectContract.Contract.TokenURI(&_GenerativeProjectContract.CallOpts, projectId)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) AddProjectScript(opts *bind.TransactOpts, projectId *big.Int, _script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "addProjectScript", projectId, _script)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) AddProjectScript(projectId *big.Int, _script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.AddProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, _script)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) AddProjectScript(projectId *big.Int, _script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.AddProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, _script)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Approve(&_GenerativeProjectContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Approve(&_GenerativeProjectContract.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeAdmin(&_GenerativeProjectContract.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeAdmin(&_GenerativeProjectContract.TransactOpts, newAdm)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) ChangeDataContextAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "changeDataContextAddr", newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeDataContextAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeDataContextAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) ChangeParamAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "changeParamAddr", newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeParamAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeParamAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) ChangeRandomizerAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "changeRandomizerAddr", newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeRandomizerAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.ChangeRandomizerAddr(&_GenerativeProjectContract.TransactOpts, newAddr)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) CompleteProject(opts *bind.TransactOpts, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "completeProject", projectId)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) CompleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.CompleteProject(&_GenerativeProjectContract.TransactOpts, projectId)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) CompleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.CompleteProject(&_GenerativeProjectContract.TransactOpts, projectId)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) DeleteProjectScript(opts *bind.TransactOpts, projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "deleteProjectScript", projectId, scriptIndex)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) DeleteProjectScript(projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.DeleteProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, scriptIndex)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) DeleteProjectScript(projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.DeleteProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, scriptIndex)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "initialize", name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Initialize(&_GenerativeProjectContract.TransactOpts, name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Initialize(&_GenerativeProjectContract.TransactOpts, name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) Mint(opts *bind.TransactOpts, project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "mint", project, disable, openingTime)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractSession) Mint(project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Mint(&_GenerativeProjectContract.TransactOpts, project, disable, openingTime)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) Mint(project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Mint(&_GenerativeProjectContract.TransactOpts, project, disable, openingTime)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.RenounceOwnership(&_GenerativeProjectContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.RenounceOwnership(&_GenerativeProjectContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SafeTransferFrom(&_GenerativeProjectContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SafeTransferFrom(&_GenerativeProjectContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SafeTransferFrom0(&_GenerativeProjectContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SafeTransferFrom0(&_GenerativeProjectContract.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SetApprovalForAll(&_GenerativeProjectContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SetApprovalForAll(&_GenerativeProjectContract.TransactOpts, operator, approved)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) SetProjectStatus(opts *bind.TransactOpts, projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "setProjectStatus", projectId, enable)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) SetProjectStatus(projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SetProjectStatus(&_GenerativeProjectContract.TransactOpts, projectId, enable)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) SetProjectStatus(projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.SetProjectStatus(&_GenerativeProjectContract.TransactOpts, projectId, enable)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.TransferFrom(&_GenerativeProjectContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.TransferFrom(&_GenerativeProjectContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.TransferOwnership(&_GenerativeProjectContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.TransferOwnership(&_GenerativeProjectContract.TransactOpts, newOwner)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectCreatorName(opts *bind.TransactOpts, projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectCreatorName", projectId, creatorName)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectCreatorName(projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectCreatorName(&_GenerativeProjectContract.TransactOpts, projectId, creatorName)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectCreatorName(projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectCreatorName(&_GenerativeProjectContract.TransactOpts, projectId, creatorName)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectLicense(opts *bind.TransactOpts, projectId *big.Int, license string) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectLicense", projectId, license)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectLicense(projectId *big.Int, license string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectLicense(&_GenerativeProjectContract.TransactOpts, projectId, license)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectLicense(projectId *big.Int, license string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectLicense(&_GenerativeProjectContract.TransactOpts, projectId, license)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectName(opts *bind.TransactOpts, projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectName", projectId, projectName)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectName(projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectName(&_GenerativeProjectContract.TransactOpts, projectId, projectName)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectName(projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectName(&_GenerativeProjectContract.TransactOpts, projectId, projectName)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectPrice(opts *bind.TransactOpts, projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectPrice", projectId, price)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectPrice(projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectPrice(&_GenerativeProjectContract.TransactOpts, projectId, price)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectPrice(projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectPrice(&_GenerativeProjectContract.TransactOpts, projectId, price)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectScript(opts *bind.TransactOpts, projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectScript", projectId, scriptIndex, script)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectScript(projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, scriptIndex, script)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectScript(projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectScript(&_GenerativeProjectContract.TransactOpts, projectId, scriptIndex, script)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectScriptType(opts *bind.TransactOpts, projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectScriptType", projectId, scriptType, i)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectScriptType(projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectScriptType(&_GenerativeProjectContract.TransactOpts, projectId, scriptType, i)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectScriptType(projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectScriptType(&_GenerativeProjectContract.TransactOpts, projectId, scriptType, i)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) UpdateProjectSocial(opts *bind.TransactOpts, projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "updateProjectSocial", projectId, data)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) UpdateProjectSocial(projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectSocial(&_GenerativeProjectContract.TransactOpts, projectId, data)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) UpdateProjectSocial(projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.UpdateProjectSocial(&_GenerativeProjectContract.TransactOpts, projectId, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactor) Withdraw(opts *bind.TransactOpts, receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.contract.Transact(opts, "withdraw", receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeProjectContract *GenerativeProjectContractSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Withdraw(&_GenerativeProjectContract.TransactOpts, receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeProjectContract *GenerativeProjectContractTransactorSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeProjectContract.Contract.Withdraw(&_GenerativeProjectContract.TransactOpts, receiver, erc20Addr, amount)
}

// GenerativeProjectContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractApprovalIterator struct {
	Event *GenerativeProjectContractApproval // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractApproval)
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
		it.Event = new(GenerativeProjectContractApproval)
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
func (it *GenerativeProjectContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractApproval represents a Approval event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GenerativeProjectContractApprovalIterator, error) {

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

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractApprovalIterator{contract: _GenerativeProjectContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractApproval)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseApproval(log types.Log) (*GenerativeProjectContractApproval, error) {
	event := new(GenerativeProjectContractApproval)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractApprovalForAllIterator struct {
	Event *GenerativeProjectContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractApprovalForAll)
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
		it.Event = new(GenerativeProjectContractApprovalForAll)
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
func (it *GenerativeProjectContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractApprovalForAll represents a ApprovalForAll event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GenerativeProjectContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractApprovalForAllIterator{contract: _GenerativeProjectContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractApprovalForAll)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseApprovalForAll(log types.Log) (*GenerativeProjectContractApprovalForAll, error) {
	event := new(GenerativeProjectContractApprovalForAll)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractInitializedIterator struct {
	Event *GenerativeProjectContractInitialized // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractInitialized)
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
		it.Event = new(GenerativeProjectContractInitialized)
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
func (it *GenerativeProjectContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractInitialized represents a Initialized event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*GenerativeProjectContractInitializedIterator, error) {

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractInitializedIterator{contract: _GenerativeProjectContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractInitialized) (event.Subscription, error) {

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractInitialized)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseInitialized(log types.Log) (*GenerativeProjectContractInitialized, error) {
	event := new(GenerativeProjectContractInitialized)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractOwnershipTransferredIterator struct {
	Event *GenerativeProjectContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractOwnershipTransferred)
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
		it.Event = new(GenerativeProjectContractOwnershipTransferred)
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
func (it *GenerativeProjectContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractOwnershipTransferred represents a OwnershipTransferred event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GenerativeProjectContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractOwnershipTransferredIterator{contract: _GenerativeProjectContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractOwnershipTransferred)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseOwnershipTransferred(log types.Log) (*GenerativeProjectContractOwnershipTransferred, error) {
	event := new(GenerativeProjectContractOwnershipTransferred)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractPausedIterator struct {
	Event *GenerativeProjectContractPaused // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractPaused)
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
		it.Event = new(GenerativeProjectContractPaused)
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
func (it *GenerativeProjectContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractPaused represents a Paused event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterPaused(opts *bind.FilterOpts) (*GenerativeProjectContractPausedIterator, error) {

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractPausedIterator{contract: _GenerativeProjectContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractPaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractPaused)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParsePaused(log types.Log) (*GenerativeProjectContractPaused, error) {
	event := new(GenerativeProjectContractPaused)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractTransferIterator struct {
	Event *GenerativeProjectContractTransfer // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractTransfer)
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
		it.Event = new(GenerativeProjectContractTransfer)
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
func (it *GenerativeProjectContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractTransfer represents a Transfer event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GenerativeProjectContractTransferIterator, error) {

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

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractTransferIterator{contract: _GenerativeProjectContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractTransfer)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseTransfer(log types.Log) (*GenerativeProjectContractTransfer, error) {
	event := new(GenerativeProjectContractTransfer)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GenerativeProjectContract contract.
type GenerativeProjectContractUnpausedIterator struct {
	Event *GenerativeProjectContractUnpaused // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectContractUnpaused)
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
		it.Event = new(GenerativeProjectContractUnpaused)
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
func (it *GenerativeProjectContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectContractUnpaused represents a Unpaused event raised by the GenerativeProjectContract contract.
type GenerativeProjectContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GenerativeProjectContractUnpausedIterator, error) {

	logs, sub, err := _GenerativeProjectContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectContractUnpausedIterator{contract: _GenerativeProjectContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GenerativeProjectContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeProjectContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectContractUnpaused)
				if err := _GenerativeProjectContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_GenerativeProjectContract *GenerativeProjectContractFilterer) ParseUnpaused(log types.Log) (*GenerativeProjectContractUnpaused, error) {
	event := new(GenerativeProjectContractUnpaused)
	if err := _GenerativeProjectContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
