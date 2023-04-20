// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tcartifact

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

// NFT721MetaData contains all meta data concerning the NFT721 contract.
var NFT721MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_projectDataContextAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_randomizerAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_script\",\"type\":\"string\"}],\"name\":\"addProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeDataContextAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeRandomizerAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"completeProject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scriptIndex\",\"type\":\"uint256\"}],\"name\":\"deleteProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"randomizerAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"projectDataContextAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_creatorAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_license\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_desc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_image\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"_social\",\"type\":\"tuple\"},{\"internalType\":\"string[]\",\"name\":\"_scriptType\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_scripts\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"_styles\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_completeTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_genNFTAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_itemDesc\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"_reserves\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.Project\",\"name\":\"project\",\"type\":\"tuple\"},{\"internalType\":\"bool\",\"name\":\"disable\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"openingTime\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"projectDetails\",\"outputs\":[{\"components\":[{\"internalType\":\"uint24\",\"name\":\"_maxSupply\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"_limit\",\"type\":\"uint24\"},{\"internalType\":\"uint256\",\"name\":\"_mintPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintPriceAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_creator\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_creatorAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_license\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_desc\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_image\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"_social\",\"type\":\"tuple\"},{\"internalType\":\"string[]\",\"name\":\"_scriptType\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"_scripts\",\"type\":\"string[]\"},{\"internalType\":\"string\",\"name\":\"_styles\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_completeTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_genNFTAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_itemDesc\",\"type\":\"string\"},{\"internalType\":\"address[]\",\"name\":\"_reserves\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"internalType\":\"structNFTProject.Project\",\"name\":\"project\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"projectStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"}],\"name\":\"setProjectStatus\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"creatorName\",\"type\":\"string\"}],\"name\":\"updateProjectCreatorName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"license\",\"type\":\"string\"}],\"name\":\"updateProjectLicense\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"projectName\",\"type\":\"string\"}],\"name\":\"updateProjectName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"updateProjectPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scriptIndex\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"script\",\"type\":\"string\"}],\"name\":\"updateProjectScript\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"scriptType\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"i\",\"type\":\"uint256\"}],\"name\":\"updateProjectScriptType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"_web\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_twitter\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_discord\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_medium\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_instagram\",\"type\":\"string\"}],\"internalType\":\"structNFTProject.ProjectSocial\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"updateProjectSocial\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"erc20Addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// NFT721ABI is the input ABI used to generate the binding from.
// Deprecated: Use NFT721MetaData.ABI instead.
var NFT721ABI = NFT721MetaData.ABI

// NFT721 is an auto generated Go binding around an Ethereum contract.
type NFT721 struct {
	NFT721Caller     // Read-only binding to the contract
	NFT721Transactor // Write-only binding to the contract
	NFT721Filterer   // Log filterer for contract events
}

// NFT721Caller is an auto generated read-only Go binding around an Ethereum contract.
type NFT721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFT721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type NFT721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFT721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NFT721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NFT721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NFT721Session struct {
	Contract     *NFT721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NFT721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NFT721CallerSession struct {
	Contract *NFT721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// NFT721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NFT721TransactorSession struct {
	Contract     *NFT721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NFT721Raw is an auto generated low-level Go binding around an Ethereum contract.
type NFT721Raw struct {
	Contract *NFT721 // Generic contract binding to access the raw methods on
}

// NFT721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NFT721CallerRaw struct {
	Contract *NFT721Caller // Generic read-only contract binding to access the raw methods on
}

// NFT721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NFT721TransactorRaw struct {
	Contract *NFT721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewNFT721 creates a new instance of NFT721, bound to a specific deployed contract.
func NewNFT721(address common.Address, backend bind.ContractBackend) (*NFT721, error) {
	contract, err := bindNFT721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NFT721{NFT721Caller: NFT721Caller{contract: contract}, NFT721Transactor: NFT721Transactor{contract: contract}, NFT721Filterer: NFT721Filterer{contract: contract}}, nil
}

// NewNFT721Caller creates a new read-only instance of NFT721, bound to a specific deployed contract.
func NewNFT721Caller(address common.Address, caller bind.ContractCaller) (*NFT721Caller, error) {
	contract, err := bindNFT721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NFT721Caller{contract: contract}, nil
}

// NewNFT721Transactor creates a new write-only instance of NFT721, bound to a specific deployed contract.
func NewNFT721Transactor(address common.Address, transactor bind.ContractTransactor) (*NFT721Transactor, error) {
	contract, err := bindNFT721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NFT721Transactor{contract: contract}, nil
}

// NewNFT721Filterer creates a new log filterer instance of NFT721, bound to a specific deployed contract.
func NewNFT721Filterer(address common.Address, filterer bind.ContractFilterer) (*NFT721Filterer, error) {
	contract, err := bindNFT721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NFT721Filterer{contract: contract}, nil
}

// bindNFT721 binds a generic wrapper to an already deployed contract.
func bindNFT721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NFT721MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFT721 *NFT721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFT721.Contract.NFT721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFT721 *NFT721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFT721.Contract.NFT721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFT721 *NFT721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFT721.Contract.NFT721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NFT721 *NFT721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NFT721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NFT721 *NFT721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFT721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NFT721 *NFT721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NFT721.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_NFT721 *NFT721Caller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_NFT721 *NFT721Session) Admin() (common.Address, error) {
	return _NFT721.Contract.Admin(&_NFT721.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_NFT721 *NFT721CallerSession) Admin() (common.Address, error) {
	return _NFT721.Contract.Admin(&_NFT721.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_NFT721 *NFT721Caller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_NFT721 *NFT721Session) ParamsAddress() (common.Address, error) {
	return _NFT721.Contract.ParamsAddress(&_NFT721.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_NFT721 *NFT721CallerSession) ParamsAddress() (common.Address, error) {
	return _NFT721.Contract.ParamsAddress(&_NFT721.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_NFT721 *NFT721Caller) ProjectDataContextAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "_projectDataContextAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_NFT721 *NFT721Session) ProjectDataContextAddr() (common.Address, error) {
	return _NFT721.Contract.ProjectDataContextAddr(&_NFT721.CallOpts)
}

// ProjectDataContextAddr is a free data retrieval call binding the contract method 0x575b57ea.
//
// Solidity: function _projectDataContextAddr() view returns(address)
func (_NFT721 *NFT721CallerSession) ProjectDataContextAddr() (common.Address, error) {
	return _NFT721.Contract.ProjectDataContextAddr(&_NFT721.CallOpts)
}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_NFT721 *NFT721Caller) RandomizerAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "_randomizerAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_NFT721 *NFT721Session) RandomizerAddr() (common.Address, error) {
	return _NFT721.Contract.RandomizerAddr(&_NFT721.CallOpts)
}

// RandomizerAddr is a free data retrieval call binding the contract method 0x66215eb4.
//
// Solidity: function _randomizerAddr() view returns(address)
func (_NFT721 *NFT721CallerSession) RandomizerAddr() (common.Address, error) {
	return _NFT721.Contract.RandomizerAddr(&_NFT721.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NFT721 *NFT721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NFT721 *NFT721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NFT721.Contract.BalanceOf(&_NFT721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_NFT721 *NFT721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _NFT721.Contract.BalanceOf(&_NFT721.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NFT721.Contract.GetApproved(&_NFT721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _NFT721.Contract.GetApproved(&_NFT721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NFT721 *NFT721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NFT721 *NFT721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NFT721.Contract.IsApprovedForAll(&_NFT721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_NFT721 *NFT721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _NFT721.Contract.IsApprovedForAll(&_NFT721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NFT721 *NFT721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NFT721 *NFT721Session) Name() (string, error) {
	return _NFT721.Contract.Name(&_NFT721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_NFT721 *NFT721CallerSession) Name() (string, error) {
	return _NFT721.Contract.Name(&_NFT721.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFT721 *NFT721Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFT721 *NFT721Session) Owner() (common.Address, error) {
	return _NFT721.Contract.Owner(&_NFT721.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_NFT721 *NFT721CallerSession) Owner() (common.Address, error) {
	return _NFT721.Contract.Owner(&_NFT721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NFT721.Contract.OwnerOf(&_NFT721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_NFT721 *NFT721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _NFT721.Contract.OwnerOf(&_NFT721.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_NFT721 *NFT721Caller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_NFT721 *NFT721Session) Paused() (bool, error) {
	return _NFT721.Contract.Paused(&_NFT721.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_NFT721 *NFT721CallerSession) Paused() (bool, error) {
	return _NFT721.Contract.Paused(&_NFT721.CallOpts)
}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_NFT721 *NFT721Caller) ProjectDetails(opts *bind.CallOpts, projectId *big.Int) (NFTProjectProject, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "projectDetails", projectId)

	if err != nil {
		return *new(NFTProjectProject), err
	}

	out0 := *abi.ConvertType(out[0], new(NFTProjectProject)).(*NFTProjectProject)

	return out0, err

}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_NFT721 *NFT721Session) ProjectDetails(projectId *big.Int) (NFTProjectProject, error) {
	return _NFT721.Contract.ProjectDetails(&_NFT721.CallOpts, projectId)
}

// ProjectDetails is a free data retrieval call binding the contract method 0x8dd91a56.
//
// Solidity: function projectDetails(uint256 projectId) view returns((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project)
func (_NFT721 *NFT721CallerSession) ProjectDetails(projectId *big.Int) (NFTProjectProject, error) {
	return _NFT721.Contract.ProjectDetails(&_NFT721.CallOpts, projectId)
}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_NFT721 *NFT721Caller) ProjectStatus(opts *bind.CallOpts, projectId *big.Int) (bool, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "projectStatus", projectId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_NFT721 *NFT721Session) ProjectStatus(projectId *big.Int) (bool, error) {
	return _NFT721.Contract.ProjectStatus(&_NFT721.CallOpts, projectId)
}

// ProjectStatus is a free data retrieval call binding the contract method 0x50ac5892.
//
// Solidity: function projectStatus(uint256 projectId) view returns(bool enable)
func (_NFT721 *NFT721CallerSession) ProjectStatus(projectId *big.Int) (bool, error) {
	return _NFT721.Contract.ProjectStatus(&_NFT721.CallOpts, projectId)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 projectId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_NFT721 *NFT721Caller) RoyaltyInfo(opts *bind.CallOpts, projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "royaltyInfo", projectId, _salePrice)

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
func (_NFT721 *NFT721Session) RoyaltyInfo(projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _NFT721.Contract.RoyaltyInfo(&_NFT721.CallOpts, projectId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 projectId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_NFT721 *NFT721CallerSession) RoyaltyInfo(projectId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _NFT721.Contract.RoyaltyInfo(&_NFT721.CallOpts, projectId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NFT721 *NFT721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NFT721 *NFT721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NFT721.Contract.SupportsInterface(&_NFT721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_NFT721 *NFT721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _NFT721.Contract.SupportsInterface(&_NFT721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NFT721 *NFT721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NFT721 *NFT721Session) Symbol() (string, error) {
	return _NFT721.Contract.Symbol(&_NFT721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_NFT721 *NFT721CallerSession) Symbol() (string, error) {
	return _NFT721.Contract.Symbol(&_NFT721.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_NFT721 *NFT721Caller) TokenURI(opts *bind.CallOpts, projectId *big.Int) (string, error) {
	var out []interface{}
	err := _NFT721.contract.Call(opts, &out, "tokenURI", projectId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_NFT721 *NFT721Session) TokenURI(projectId *big.Int) (string, error) {
	return _NFT721.Contract.TokenURI(&_NFT721.CallOpts, projectId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 projectId) view returns(string result)
func (_NFT721 *NFT721CallerSession) TokenURI(projectId *big.Int) (string, error) {
	return _NFT721.Contract.TokenURI(&_NFT721.CallOpts, projectId)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_NFT721 *NFT721Transactor) AddProjectScript(opts *bind.TransactOpts, projectId *big.Int, _script string) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "addProjectScript", projectId, _script)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_NFT721 *NFT721Session) AddProjectScript(projectId *big.Int, _script string) (*types.Transaction, error) {
	return _NFT721.Contract.AddProjectScript(&_NFT721.TransactOpts, projectId, _script)
}

// AddProjectScript is a paid mutator transaction binding the contract method 0xacad0124.
//
// Solidity: function addProjectScript(uint256 projectId, string _script) returns()
func (_NFT721 *NFT721TransactorSession) AddProjectScript(projectId *big.Int, _script string) (*types.Transaction, error) {
	return _NFT721.Contract.AddProjectScript(&_NFT721.TransactOpts, projectId, _script)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Approve(&_NFT721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_NFT721 *NFT721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Approve(&_NFT721.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_NFT721 *NFT721Transactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_NFT721 *NFT721Session) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeAdmin(&_NFT721.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_NFT721 *NFT721TransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeAdmin(&_NFT721.TransactOpts, newAdm)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_NFT721 *NFT721Transactor) ChangeDataContextAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "changeDataContextAddr", newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_NFT721 *NFT721Session) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeDataContextAddr(&_NFT721.TransactOpts, newAddr)
}

// ChangeDataContextAddr is a paid mutator transaction binding the contract method 0x472f1e02.
//
// Solidity: function changeDataContextAddr(address newAddr) returns()
func (_NFT721 *NFT721TransactorSession) ChangeDataContextAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeDataContextAddr(&_NFT721.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_NFT721 *NFT721Transactor) ChangeParamAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "changeParamAddr", newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_NFT721 *NFT721Session) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeParamAddr(&_NFT721.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_NFT721 *NFT721TransactorSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeParamAddr(&_NFT721.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_NFT721 *NFT721Transactor) ChangeRandomizerAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "changeRandomizerAddr", newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_NFT721 *NFT721Session) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeRandomizerAddr(&_NFT721.TransactOpts, newAddr)
}

// ChangeRandomizerAddr is a paid mutator transaction binding the contract method 0x1ca9741b.
//
// Solidity: function changeRandomizerAddr(address newAddr) returns()
func (_NFT721 *NFT721TransactorSession) ChangeRandomizerAddr(newAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.ChangeRandomizerAddr(&_NFT721.TransactOpts, newAddr)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_NFT721 *NFT721Transactor) CompleteProject(opts *bind.TransactOpts, projectId *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "completeProject", projectId)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_NFT721 *NFT721Session) CompleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.CompleteProject(&_NFT721.TransactOpts, projectId)
}

// CompleteProject is a paid mutator transaction binding the contract method 0x2245f152.
//
// Solidity: function completeProject(uint256 projectId) returns()
func (_NFT721 *NFT721TransactorSession) CompleteProject(projectId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.CompleteProject(&_NFT721.TransactOpts, projectId)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_NFT721 *NFT721Transactor) DeleteProjectScript(opts *bind.TransactOpts, projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "deleteProjectScript", projectId, scriptIndex)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_NFT721 *NFT721Session) DeleteProjectScript(projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.DeleteProjectScript(&_NFT721.TransactOpts, projectId, scriptIndex)
}

// DeleteProjectScript is a paid mutator transaction binding the contract method 0x166b7469.
//
// Solidity: function deleteProjectScript(uint256 projectId, uint256 scriptIndex) returns()
func (_NFT721 *NFT721TransactorSession) DeleteProjectScript(projectId *big.Int, scriptIndex *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.DeleteProjectScript(&_NFT721.TransactOpts, projectId, scriptIndex)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_NFT721 *NFT721Transactor) Initialize(opts *bind.TransactOpts, name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "initialize", name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_NFT721 *NFT721Session) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.Initialize(&_NFT721.TransactOpts, name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xe56f2fe4.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress, address randomizerAddr, address projectDataContextAddr) returns()
func (_NFT721 *NFT721TransactorSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address, randomizerAddr common.Address, projectDataContextAddr common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.Initialize(&_NFT721.TransactOpts, name, symbol, admin, paramsAddress, randomizerAddr, projectDataContextAddr)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_NFT721 *NFT721Transactor) Mint(opts *bind.TransactOpts, project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "mint", project, disable, openingTime)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_NFT721 *NFT721Session) Mint(project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Mint(&_NFT721.TransactOpts, project, disable, openingTime)
}

// Mint is a paid mutator transaction binding the contract method 0x45bf4d08.
//
// Solidity: function mint((uint24,uint24,uint256,address,string,string,address,string,string,string,(string,string,string,string,string),string[],string[],string,uint256,address,string,address[],uint256) project, bool disable, uint256 openingTime) payable returns(uint256)
func (_NFT721 *NFT721TransactorSession) Mint(project NFTProjectProject, disable bool, openingTime *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Mint(&_NFT721.TransactOpts, project, disable, openingTime)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFT721 *NFT721Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFT721 *NFT721Session) RenounceOwnership() (*types.Transaction, error) {
	return _NFT721.Contract.RenounceOwnership(&_NFT721.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_NFT721 *NFT721TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _NFT721.Contract.RenounceOwnership(&_NFT721.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.SafeTransferFrom(&_NFT721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.SafeTransferFrom(&_NFT721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NFT721 *NFT721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NFT721 *NFT721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NFT721.Contract.SafeTransferFrom0(&_NFT721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_NFT721 *NFT721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _NFT721.Contract.SafeTransferFrom0(&_NFT721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NFT721 *NFT721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NFT721 *NFT721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NFT721.Contract.SetApprovalForAll(&_NFT721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_NFT721 *NFT721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _NFT721.Contract.SetApprovalForAll(&_NFT721.TransactOpts, operator, approved)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_NFT721 *NFT721Transactor) SetProjectStatus(opts *bind.TransactOpts, projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "setProjectStatus", projectId, enable)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_NFT721 *NFT721Session) SetProjectStatus(projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _NFT721.Contract.SetProjectStatus(&_NFT721.TransactOpts, projectId, enable)
}

// SetProjectStatus is a paid mutator transaction binding the contract method 0x3af6cd1e.
//
// Solidity: function setProjectStatus(uint256 projectId, bool enable) returns()
func (_NFT721 *NFT721TransactorSession) SetProjectStatus(projectId *big.Int, enable bool) (*types.Transaction, error) {
	return _NFT721.Contract.SetProjectStatus(&_NFT721.TransactOpts, projectId, enable)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.TransferFrom(&_NFT721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_NFT721 *NFT721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.TransferFrom(&_NFT721.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFT721 *NFT721Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFT721 *NFT721Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.TransferOwnership(&_NFT721.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_NFT721 *NFT721TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _NFT721.Contract.TransferOwnership(&_NFT721.TransactOpts, newOwner)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectCreatorName(opts *bind.TransactOpts, projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectCreatorName", projectId, creatorName)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_NFT721 *NFT721Session) UpdateProjectCreatorName(projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectCreatorName(&_NFT721.TransactOpts, projectId, creatorName)
}

// UpdateProjectCreatorName is a paid mutator transaction binding the contract method 0xe283d82f.
//
// Solidity: function updateProjectCreatorName(uint256 projectId, string creatorName) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectCreatorName(projectId *big.Int, creatorName string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectCreatorName(&_NFT721.TransactOpts, projectId, creatorName)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectLicense(opts *bind.TransactOpts, projectId *big.Int, license string) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectLicense", projectId, license)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_NFT721 *NFT721Session) UpdateProjectLicense(projectId *big.Int, license string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectLicense(&_NFT721.TransactOpts, projectId, license)
}

// UpdateProjectLicense is a paid mutator transaction binding the contract method 0x25b75d68.
//
// Solidity: function updateProjectLicense(uint256 projectId, string license) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectLicense(projectId *big.Int, license string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectLicense(&_NFT721.TransactOpts, projectId, license)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectName(opts *bind.TransactOpts, projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectName", projectId, projectName)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_NFT721 *NFT721Session) UpdateProjectName(projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectName(&_NFT721.TransactOpts, projectId, projectName)
}

// UpdateProjectName is a paid mutator transaction binding the contract method 0x0d170673.
//
// Solidity: function updateProjectName(uint256 projectId, string projectName) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectName(projectId *big.Int, projectName string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectName(&_NFT721.TransactOpts, projectId, projectName)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectPrice(opts *bind.TransactOpts, projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectPrice", projectId, price)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_NFT721 *NFT721Session) UpdateProjectPrice(projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectPrice(&_NFT721.TransactOpts, projectId, price)
}

// UpdateProjectPrice is a paid mutator transaction binding the contract method 0x92655336.
//
// Solidity: function updateProjectPrice(uint256 projectId, uint256 price) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectPrice(projectId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectPrice(&_NFT721.TransactOpts, projectId, price)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectScript(opts *bind.TransactOpts, projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectScript", projectId, scriptIndex, script)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_NFT721 *NFT721Session) UpdateProjectScript(projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectScript(&_NFT721.TransactOpts, projectId, scriptIndex, script)
}

// UpdateProjectScript is a paid mutator transaction binding the contract method 0xb1656ba3.
//
// Solidity: function updateProjectScript(uint256 projectId, uint256 scriptIndex, string script) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectScript(projectId *big.Int, scriptIndex *big.Int, script string) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectScript(&_NFT721.TransactOpts, projectId, scriptIndex, script)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectScriptType(opts *bind.TransactOpts, projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectScriptType", projectId, scriptType, i)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_NFT721 *NFT721Session) UpdateProjectScriptType(projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectScriptType(&_NFT721.TransactOpts, projectId, scriptType, i)
}

// UpdateProjectScriptType is a paid mutator transaction binding the contract method 0xdaf61800.
//
// Solidity: function updateProjectScriptType(uint256 projectId, string scriptType, uint256 i) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectScriptType(projectId *big.Int, scriptType string, i *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectScriptType(&_NFT721.TransactOpts, projectId, scriptType, i)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_NFT721 *NFT721Transactor) UpdateProjectSocial(opts *bind.TransactOpts, projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "updateProjectSocial", projectId, data)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_NFT721 *NFT721Session) UpdateProjectSocial(projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectSocial(&_NFT721.TransactOpts, projectId, data)
}

// UpdateProjectSocial is a paid mutator transaction binding the contract method 0x8e1ab777.
//
// Solidity: function updateProjectSocial(uint256 projectId, (string,string,string,string,string) data) returns()
func (_NFT721 *NFT721TransactorSession) UpdateProjectSocial(projectId *big.Int, data NFTProjectProjectSocial) (*types.Transaction, error) {
	return _NFT721.Contract.UpdateProjectSocial(&_NFT721.TransactOpts, projectId, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_NFT721 *NFT721Transactor) Withdraw(opts *bind.TransactOpts, receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFT721.contract.Transact(opts, "withdraw", receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_NFT721 *NFT721Session) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Withdraw(&_NFT721.TransactOpts, receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_NFT721 *NFT721TransactorSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _NFT721.Contract.Withdraw(&_NFT721.TransactOpts, receiver, erc20Addr, amount)
}

// NFT721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the NFT721 contract.
type NFT721ApprovalIterator struct {
	Event *NFT721Approval // Event containing the contract specifics and raw log

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
func (it *NFT721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721Approval)
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
		it.Event = new(NFT721Approval)
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
func (it *NFT721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721Approval represents a Approval event raised by the NFT721 contract.
type NFT721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NFT721 *NFT721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*NFT721ApprovalIterator, error) {

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

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NFT721ApprovalIterator{contract: _NFT721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_NFT721 *NFT721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *NFT721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721Approval)
				if err := _NFT721.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseApproval(log types.Log) (*NFT721Approval, error) {
	event := new(NFT721Approval)
	if err := _NFT721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the NFT721 contract.
type NFT721ApprovalForAllIterator struct {
	Event *NFT721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *NFT721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721ApprovalForAll)
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
		it.Event = new(NFT721ApprovalForAll)
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
func (it *NFT721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721ApprovalForAll represents a ApprovalForAll event raised by the NFT721 contract.
type NFT721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NFT721 *NFT721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*NFT721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &NFT721ApprovalForAllIterator{contract: _NFT721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_NFT721 *NFT721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *NFT721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721ApprovalForAll)
				if err := _NFT721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseApprovalForAll(log types.Log) (*NFT721ApprovalForAll, error) {
	event := new(NFT721ApprovalForAll)
	if err := _NFT721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the NFT721 contract.
type NFT721InitializedIterator struct {
	Event *NFT721Initialized // Event containing the contract specifics and raw log

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
func (it *NFT721InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721Initialized)
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
		it.Event = new(NFT721Initialized)
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
func (it *NFT721InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721Initialized represents a Initialized event raised by the NFT721 contract.
type NFT721Initialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_NFT721 *NFT721Filterer) FilterInitialized(opts *bind.FilterOpts) (*NFT721InitializedIterator, error) {

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &NFT721InitializedIterator{contract: _NFT721.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_NFT721 *NFT721Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *NFT721Initialized) (event.Subscription, error) {

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721Initialized)
				if err := _NFT721.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseInitialized(log types.Log) (*NFT721Initialized, error) {
	event := new(NFT721Initialized)
	if err := _NFT721.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the NFT721 contract.
type NFT721OwnershipTransferredIterator struct {
	Event *NFT721OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *NFT721OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721OwnershipTransferred)
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
		it.Event = new(NFT721OwnershipTransferred)
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
func (it *NFT721OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721OwnershipTransferred represents a OwnershipTransferred event raised by the NFT721 contract.
type NFT721OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NFT721 *NFT721Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*NFT721OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &NFT721OwnershipTransferredIterator{contract: _NFT721.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_NFT721 *NFT721Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *NFT721OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721OwnershipTransferred)
				if err := _NFT721.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseOwnershipTransferred(log types.Log) (*NFT721OwnershipTransferred, error) {
	event := new(NFT721OwnershipTransferred)
	if err := _NFT721.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721PausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the NFT721 contract.
type NFT721PausedIterator struct {
	Event *NFT721Paused // Event containing the contract specifics and raw log

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
func (it *NFT721PausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721Paused)
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
		it.Event = new(NFT721Paused)
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
func (it *NFT721PausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721PausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721Paused represents a Paused event raised by the NFT721 contract.
type NFT721Paused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_NFT721 *NFT721Filterer) FilterPaused(opts *bind.FilterOpts) (*NFT721PausedIterator, error) {

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &NFT721PausedIterator{contract: _NFT721.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_NFT721 *NFT721Filterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *NFT721Paused) (event.Subscription, error) {

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721Paused)
				if err := _NFT721.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParsePaused(log types.Log) (*NFT721Paused, error) {
	event := new(NFT721Paused)
	if err := _NFT721.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the NFT721 contract.
type NFT721TransferIterator struct {
	Event *NFT721Transfer // Event containing the contract specifics and raw log

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
func (it *NFT721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721Transfer)
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
		it.Event = new(NFT721Transfer)
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
func (it *NFT721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721Transfer represents a Transfer event raised by the NFT721 contract.
type NFT721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NFT721 *NFT721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*NFT721TransferIterator, error) {

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

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &NFT721TransferIterator{contract: _NFT721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_NFT721 *NFT721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *NFT721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721Transfer)
				if err := _NFT721.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseTransfer(log types.Log) (*NFT721Transfer, error) {
	event := new(NFT721Transfer)
	if err := _NFT721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NFT721UnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the NFT721 contract.
type NFT721UnpausedIterator struct {
	Event *NFT721Unpaused // Event containing the contract specifics and raw log

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
func (it *NFT721UnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NFT721Unpaused)
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
		it.Event = new(NFT721Unpaused)
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
func (it *NFT721UnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NFT721UnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NFT721Unpaused represents a Unpaused event raised by the NFT721 contract.
type NFT721Unpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_NFT721 *NFT721Filterer) FilterUnpaused(opts *bind.FilterOpts) (*NFT721UnpausedIterator, error) {

	logs, sub, err := _NFT721.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &NFT721UnpausedIterator{contract: _NFT721.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_NFT721 *NFT721Filterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *NFT721Unpaused) (event.Subscription, error) {

	logs, sub, err := _NFT721.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NFT721Unpaused)
				if err := _NFT721.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_NFT721 *NFT721Filterer) ParseUnpaused(log types.Log) (*NFT721Unpaused, error) {
	event := new(NFT721Unpaused)
	if err := _NFT721.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
