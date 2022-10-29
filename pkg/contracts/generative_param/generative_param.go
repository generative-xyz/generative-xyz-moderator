// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_param

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

// BoilerplateParamParamTemplate is an auto generated low-level Go binding around an user-defined struct.
type BoilerplateParamParamTemplate struct {
	TypeValue       uint8
	Max             *big.Int
	Min             *big.Int
	Decimal         uint8
	AvailableValues []string
	Value           *big.Int
	Editable        bool
}

// BoilerplateParamParamsOfNFT is an auto generated low-level Go binding around an user-defined struct.
type BoilerplateParamParamsOfNFT struct {
	Seed  [32]byte
	Value []*big.Int
}

// TraitInfoTrait is an auto generated low-level Go binding around an user-defined struct.
type TraitInfoTrait struct {
	Name            string
	AvailableValues []string
	Value           *big.Int
	ValueStr        string
}

// TraitInfoTraits is an auto generated low-level Go binding around an user-defined struct.
type TraitInfoTraits struct {
	Traits []TraitInfoTrait
}

// GenerativeParamMetaData contains all meta data concerning the GenerativeParam contract.
var GenerativeParamMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_baseuri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"mintTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"MintGenerativeNFT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_boilerplateAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_boilerplateId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_creators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_paramsValues\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseTokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newBoilerplate\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"newBoilerplateId\",\"type\":\"uint256\"}],\"name\":\"changeBoilerplate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getParamValues\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_typeValue\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_editable\",\"type\":\"bool\"}],\"internalType\":\"structBoilerplateParam.ParamTemplate[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getTokenTraits\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_valueStr\",\"type\":\"string\"}],\"internalType\":\"structTraitInfo.Trait[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTraits\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_valueStr\",\"type\":\"string\"}],\"internalType\":\"structTraitInfo.Trait[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"boilerplateAdd\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"boilerplateId\",\"type\":\"uint256\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"mintTo\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"_value\",\"type\":\"uint256[]\"}],\"internalType\":\"structBoilerplateParam.ParamsOfNFT\",\"name\":\"_paramsTemplateValue\",\"type\":\"tuple\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"royalties\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"amount\",\"type\":\"uint24\"},{\"internalType\":\"bool\",\"name\":\"isValue\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"setCreator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_newURI\",\"type\":\"string\"}],\"name\":\"setCustomURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setTokenRoyalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_valueStr\",\"type\":\"string\"}],\"internalType\":\"structTraitInfo.Trait[]\",\"name\":\"_traits\",\"type\":\"tuple[]\"}],\"internalType\":\"structTraitInfo.Traits\",\"name\":\"traits\",\"type\":\"tuple\"}],\"name\":\"updateTraits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeParamABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeParamMetaData.ABI instead.
var GenerativeParamABI = GenerativeParamMetaData.ABI

// GenerativeParam is an auto generated Go binding around an Ethereum contract.
type GenerativeParam struct {
	GenerativeParamCaller     // Read-only binding to the contract
	GenerativeParamTransactor // Write-only binding to the contract
	GenerativeParamFilterer   // Log filterer for contract events
}

// GenerativeParamCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeParamCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeParamTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeParamTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeParamFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeParamFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeParamSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeParamSession struct {
	Contract     *GenerativeParam  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GenerativeParamCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeParamCallerSession struct {
	Contract *GenerativeParamCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// GenerativeParamTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeParamTransactorSession struct {
	Contract     *GenerativeParamTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// GenerativeParamRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeParamRaw struct {
	Contract *GenerativeParam // Generic contract binding to access the raw methods on
}

// GenerativeParamCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeParamCallerRaw struct {
	Contract *GenerativeParamCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeParamTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeParamTransactorRaw struct {
	Contract *GenerativeParamTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeParam creates a new instance of GenerativeParam, bound to a specific deployed contract.
func NewGenerativeParam(address common.Address, backend bind.ContractBackend) (*GenerativeParam, error) {
	contract, err := bindGenerativeParam(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeParam{GenerativeParamCaller: GenerativeParamCaller{contract: contract}, GenerativeParamTransactor: GenerativeParamTransactor{contract: contract}, GenerativeParamFilterer: GenerativeParamFilterer{contract: contract}}, nil
}

// NewGenerativeParamCaller creates a new read-only instance of GenerativeParam, bound to a specific deployed contract.
func NewGenerativeParamCaller(address common.Address, caller bind.ContractCaller) (*GenerativeParamCaller, error) {
	contract, err := bindGenerativeParam(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamCaller{contract: contract}, nil
}

// NewGenerativeParamTransactor creates a new write-only instance of GenerativeParam, bound to a specific deployed contract.
func NewGenerativeParamTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeParamTransactor, error) {
	contract, err := bindGenerativeParam(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamTransactor{contract: contract}, nil
}

// NewGenerativeParamFilterer creates a new log filterer instance of GenerativeParam, bound to a specific deployed contract.
func NewGenerativeParamFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeParamFilterer, error) {
	contract, err := bindGenerativeParam(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamFilterer{contract: contract}, nil
}

// bindGenerativeParam binds a generic wrapper to an already deployed contract.
func bindGenerativeParam(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeParamMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeParam *GenerativeParamRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeParam.Contract.GenerativeParamCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeParam *GenerativeParamRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeParam.Contract.GenerativeParamTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeParam *GenerativeParamRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeParam.Contract.GenerativeParamTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeParam *GenerativeParamCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeParam.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeParam *GenerativeParamTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeParam.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeParam *GenerativeParamTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeParam.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.DEFAULTADMINROLE(&_GenerativeParam.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.DEFAULTADMINROLE(&_GenerativeParam.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamSession) MINTERROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.MINTERROLE(&_GenerativeParam.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCallerSession) MINTERROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.MINTERROLE(&_GenerativeParam.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamSession) PAUSERROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.PAUSERROLE(&_GenerativeParam.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeParam *GenerativeParamCallerSession) PAUSERROLE() ([32]byte, error) {
	return _GenerativeParam.Contract.PAUSERROLE(&_GenerativeParam.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeParam *GenerativeParamCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeParam *GenerativeParamSession) Admin() (common.Address, error) {
	return _GenerativeParam.Contract.Admin(&_GenerativeParam.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) Admin() (common.Address, error) {
	return _GenerativeParam.Contract.Admin(&_GenerativeParam.CallOpts)
}

// BoilerplateAddr is a free data retrieval call binding the contract method 0xf9ecc22c.
//
// Solidity: function _boilerplateAddr() view returns(address)
func (_GenerativeParam *GenerativeParamCaller) BoilerplateAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "_boilerplateAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BoilerplateAddr is a free data retrieval call binding the contract method 0xf9ecc22c.
//
// Solidity: function _boilerplateAddr() view returns(address)
func (_GenerativeParam *GenerativeParamSession) BoilerplateAddr() (common.Address, error) {
	return _GenerativeParam.Contract.BoilerplateAddr(&_GenerativeParam.CallOpts)
}

// BoilerplateAddr is a free data retrieval call binding the contract method 0xf9ecc22c.
//
// Solidity: function _boilerplateAddr() view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) BoilerplateAddr() (common.Address, error) {
	return _GenerativeParam.Contract.BoilerplateAddr(&_GenerativeParam.CallOpts)
}

// BoilerplateId is a free data retrieval call binding the contract method 0x531f7b53.
//
// Solidity: function _boilerplateId() view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) BoilerplateId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "_boilerplateId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BoilerplateId is a free data retrieval call binding the contract method 0x531f7b53.
//
// Solidity: function _boilerplateId() view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) BoilerplateId() (*big.Int, error) {
	return _GenerativeParam.Contract.BoilerplateId(&_GenerativeParam.CallOpts)
}

// BoilerplateId is a free data retrieval call binding the contract method 0x531f7b53.
//
// Solidity: function _boilerplateId() view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) BoilerplateId() (*big.Int, error) {
	return _GenerativeParam.Contract.BoilerplateId(&_GenerativeParam.CallOpts)
}

// Creators is a free data retrieval call binding the contract method 0x4816edf4.
//
// Solidity: function _creators(uint256 ) view returns(address)
func (_GenerativeParam *GenerativeParamCaller) Creators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "_creators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creators is a free data retrieval call binding the contract method 0x4816edf4.
//
// Solidity: function _creators(uint256 ) view returns(address)
func (_GenerativeParam *GenerativeParamSession) Creators(arg0 *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.Creators(&_GenerativeParam.CallOpts, arg0)
}

// Creators is a free data retrieval call binding the contract method 0x4816edf4.
//
// Solidity: function _creators(uint256 ) view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) Creators(arg0 *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.Creators(&_GenerativeParam.CallOpts, arg0)
}

// ParamsValues is a free data retrieval call binding the contract method 0x953cd246.
//
// Solidity: function _paramsValues(uint256 ) view returns(bytes32 _seed)
func (_GenerativeParam *GenerativeParamCaller) ParamsValues(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "_paramsValues", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ParamsValues is a free data retrieval call binding the contract method 0x953cd246.
//
// Solidity: function _paramsValues(uint256 ) view returns(bytes32 _seed)
func (_GenerativeParam *GenerativeParamSession) ParamsValues(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeParam.Contract.ParamsValues(&_GenerativeParam.CallOpts, arg0)
}

// ParamsValues is a free data retrieval call binding the contract method 0x953cd246.
//
// Solidity: function _paramsValues(uint256 ) view returns(bytes32 _seed)
func (_GenerativeParam *GenerativeParamCallerSession) ParamsValues(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeParam.Contract.ParamsValues(&_GenerativeParam.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeParam.Contract.BalanceOf(&_GenerativeParam.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeParam.Contract.BalanceOf(&_GenerativeParam.CallOpts, owner)
}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeParam *GenerativeParamCaller) BaseTokenURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "baseTokenURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeParam *GenerativeParamSession) BaseTokenURI() (string, error) {
	return _GenerativeParam.Contract.BaseTokenURI(&_GenerativeParam.CallOpts)
}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeParam *GenerativeParamCallerSession) BaseTokenURI() (string, error) {
	return _GenerativeParam.Contract.BaseTokenURI(&_GenerativeParam.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.GetApproved(&_GenerativeParam.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.GetApproved(&_GenerativeParam.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(bytes32, (uint8,uint256,uint256,uint8,string[],uint256,bool)[])
func (_GenerativeParam *GenerativeParamCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) ([32]byte, []BoilerplateParamParamTemplate, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new([32]byte), *new([]BoilerplateParamParamTemplate), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new([]BoilerplateParamParamTemplate)).(*[]BoilerplateParamParamTemplate)

	return out0, out1, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(bytes32, (uint8,uint256,uint256,uint8,string[],uint256,bool)[])
func (_GenerativeParam *GenerativeParamSession) GetParamValues(tokenId *big.Int) ([32]byte, []BoilerplateParamParamTemplate, error) {
	return _GenerativeParam.Contract.GetParamValues(&_GenerativeParam.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(bytes32, (uint8,uint256,uint256,uint8,string[],uint256,bool)[])
func (_GenerativeParam *GenerativeParamCallerSession) GetParamValues(tokenId *big.Int) ([32]byte, []BoilerplateParamParamTemplate, error) {
	return _GenerativeParam.Contract.GetParamValues(&_GenerativeParam.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeParam *GenerativeParamCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeParam *GenerativeParamSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GenerativeParam.Contract.GetRoleAdmin(&_GenerativeParam.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeParam *GenerativeParamCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GenerativeParam.Contract.GetRoleAdmin(&_GenerativeParam.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeParam *GenerativeParamCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeParam *GenerativeParamSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.GetRoleMember(&_GenerativeParam.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.GetRoleMember(&_GenerativeParam.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _GenerativeParam.Contract.GetRoleMemberCount(&_GenerativeParam.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _GenerativeParam.Contract.GetRoleMemberCount(&_GenerativeParam.CallOpts, role)
}

// GetTokenTraits is a free data retrieval call binding the contract method 0x94e56847.
//
// Solidity: function getTokenTraits(uint256 tokenId) view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamCaller) GetTokenTraits(opts *bind.CallOpts, tokenId *big.Int) ([]TraitInfoTrait, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getTokenTraits", tokenId)

	if err != nil {
		return *new([]TraitInfoTrait), err
	}

	out0 := *abi.ConvertType(out[0], new([]TraitInfoTrait)).(*[]TraitInfoTrait)

	return out0, err

}

// GetTokenTraits is a free data retrieval call binding the contract method 0x94e56847.
//
// Solidity: function getTokenTraits(uint256 tokenId) view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamSession) GetTokenTraits(tokenId *big.Int) ([]TraitInfoTrait, error) {
	return _GenerativeParam.Contract.GetTokenTraits(&_GenerativeParam.CallOpts, tokenId)
}

// GetTokenTraits is a free data retrieval call binding the contract method 0x94e56847.
//
// Solidity: function getTokenTraits(uint256 tokenId) view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamCallerSession) GetTokenTraits(tokenId *big.Int) ([]TraitInfoTrait, error) {
	return _GenerativeParam.Contract.GetTokenTraits(&_GenerativeParam.CallOpts, tokenId)
}

// GetTraits is a free data retrieval call binding the contract method 0x7f0429ef.
//
// Solidity: function getTraits() view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamCaller) GetTraits(opts *bind.CallOpts) ([]TraitInfoTrait, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "getTraits")

	if err != nil {
		return *new([]TraitInfoTrait), err
	}

	out0 := *abi.ConvertType(out[0], new([]TraitInfoTrait)).(*[]TraitInfoTrait)

	return out0, err

}

// GetTraits is a free data retrieval call binding the contract method 0x7f0429ef.
//
// Solidity: function getTraits() view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamSession) GetTraits() ([]TraitInfoTrait, error) {
	return _GenerativeParam.Contract.GetTraits(&_GenerativeParam.CallOpts)
}

// GetTraits is a free data retrieval call binding the contract method 0x7f0429ef.
//
// Solidity: function getTraits() view returns((string,string[],uint256,string)[])
func (_GenerativeParam *GenerativeParamCallerSession) GetTraits() ([]TraitInfoTrait, error) {
	return _GenerativeParam.Contract.GetTraits(&_GenerativeParam.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeParam *GenerativeParamCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeParam *GenerativeParamSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GenerativeParam.Contract.HasRole(&_GenerativeParam.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeParam *GenerativeParamCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GenerativeParam.Contract.HasRole(&_GenerativeParam.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeParam *GenerativeParamCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeParam *GenerativeParamSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeParam.Contract.IsApprovedForAll(&_GenerativeParam.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeParam *GenerativeParamCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeParam.Contract.IsApprovedForAll(&_GenerativeParam.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeParam *GenerativeParamCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeParam *GenerativeParamSession) Name() (string, error) {
	return _GenerativeParam.Contract.Name(&_GenerativeParam.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeParam *GenerativeParamCallerSession) Name() (string, error) {
	return _GenerativeParam.Contract.Name(&_GenerativeParam.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeParam *GenerativeParamCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeParam *GenerativeParamSession) Owner() (common.Address, error) {
	return _GenerativeParam.Contract.Owner(&_GenerativeParam.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) Owner() (common.Address, error) {
	return _GenerativeParam.Contract.Owner(&_GenerativeParam.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.OwnerOf(&_GenerativeParam.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeParam *GenerativeParamCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeParam.Contract.OwnerOf(&_GenerativeParam.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeParam *GenerativeParamCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeParam *GenerativeParamSession) Paused() (bool, error) {
	return _GenerativeParam.Contract.Paused(&_GenerativeParam.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeParam *GenerativeParamCallerSession) Paused() (bool, error) {
	return _GenerativeParam.Contract.Paused(&_GenerativeParam.CallOpts)
}

// Royalties is a free data retrieval call binding the contract method 0x7f77f574.
//
// Solidity: function royalties(uint256 ) view returns(address recipient, uint24 amount, bool isValue)
func (_GenerativeParam *GenerativeParamCaller) Royalties(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "royalties", arg0)

	outstruct := new(struct {
		Recipient common.Address
		Amount    *big.Int
		IsValue   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Recipient = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.IsValue = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Royalties is a free data retrieval call binding the contract method 0x7f77f574.
//
// Solidity: function royalties(uint256 ) view returns(address recipient, uint24 amount, bool isValue)
func (_GenerativeParam *GenerativeParamSession) Royalties(arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	return _GenerativeParam.Contract.Royalties(&_GenerativeParam.CallOpts, arg0)
}

// Royalties is a free data retrieval call binding the contract method 0x7f77f574.
//
// Solidity: function royalties(uint256 ) view returns(address recipient, uint24 amount, bool isValue)
func (_GenerativeParam *GenerativeParamCallerSession) Royalties(arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	return _GenerativeParam.Contract.Royalties(&_GenerativeParam.CallOpts, arg0)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeParam *GenerativeParamCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

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
func (_GenerativeParam *GenerativeParamSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeParam.Contract.RoyaltyInfo(&_GenerativeParam.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeParam *GenerativeParamCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeParam.Contract.RoyaltyInfo(&_GenerativeParam.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeParam *GenerativeParamCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeParam *GenerativeParamSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeParam.Contract.SupportsInterface(&_GenerativeParam.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeParam *GenerativeParamCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeParam.Contract.SupportsInterface(&_GenerativeParam.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeParam *GenerativeParamCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeParam *GenerativeParamSession) Symbol() (string, error) {
	return _GenerativeParam.Contract.Symbol(&_GenerativeParam.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeParam *GenerativeParamCallerSession) Symbol() (string, error) {
	return _GenerativeParam.Contract.Symbol(&_GenerativeParam.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _GenerativeParam.Contract.TokenByIndex(&_GenerativeParam.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _GenerativeParam.Contract.TokenByIndex(&_GenerativeParam.CallOpts, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _GenerativeParam.Contract.TokenOfOwnerByIndex(&_GenerativeParam.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _GenerativeParam.Contract.TokenOfOwnerByIndex(&_GenerativeParam.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeParam *GenerativeParamCaller) TokenURI(opts *bind.CallOpts, _tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "tokenURI", _tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeParam *GenerativeParamSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _GenerativeParam.Contract.TokenURI(&_GenerativeParam.CallOpts, _tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeParam *GenerativeParamCallerSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _GenerativeParam.Contract.TokenURI(&_GenerativeParam.CallOpts, _tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeParam *GenerativeParamCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeParam.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeParam *GenerativeParamSession) TotalSupply() (*big.Int, error) {
	return _GenerativeParam.Contract.TotalSupply(&_GenerativeParam.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeParam *GenerativeParamCallerSession) TotalSupply() (*big.Int, error) {
	return _GenerativeParam.Contract.TotalSupply(&_GenerativeParam.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Approve(&_GenerativeParam.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Approve(&_GenerativeParam.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Burn(&_GenerativeParam.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Burn(&_GenerativeParam.TransactOpts, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _newAdmin) returns()
func (_GenerativeParam *GenerativeParamTransactor) ChangeAdmin(opts *bind.TransactOpts, _newAdmin common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "changeAdmin", _newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _newAdmin) returns()
func (_GenerativeParam *GenerativeParamSession) ChangeAdmin(_newAdmin common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.ChangeAdmin(&_GenerativeParam.TransactOpts, _newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address _newAdmin) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) ChangeAdmin(_newAdmin common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.ChangeAdmin(&_GenerativeParam.TransactOpts, _newAdmin)
}

// ChangeBoilerplate is a paid mutator transaction binding the contract method 0x456c9504.
//
// Solidity: function changeBoilerplate(address newBoilerplate, uint256 newBoilerplateId) returns()
func (_GenerativeParam *GenerativeParamTransactor) ChangeBoilerplate(opts *bind.TransactOpts, newBoilerplate common.Address, newBoilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "changeBoilerplate", newBoilerplate, newBoilerplateId)
}

// ChangeBoilerplate is a paid mutator transaction binding the contract method 0x456c9504.
//
// Solidity: function changeBoilerplate(address newBoilerplate, uint256 newBoilerplateId) returns()
func (_GenerativeParam *GenerativeParamSession) ChangeBoilerplate(newBoilerplate common.Address, newBoilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.ChangeBoilerplate(&_GenerativeParam.TransactOpts, newBoilerplate, newBoilerplateId)
}

// ChangeBoilerplate is a paid mutator transaction binding the contract method 0x456c9504.
//
// Solidity: function changeBoilerplate(address newBoilerplate, uint256 newBoilerplateId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) ChangeBoilerplate(newBoilerplate common.Address, newBoilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.ChangeBoilerplate(&_GenerativeParam.TransactOpts, newBoilerplate, newBoilerplateId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.GrantRole(&_GenerativeParam.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.GrantRole(&_GenerativeParam.TransactOpts, role, account)
}

// Init is a paid mutator transaction binding the contract method 0xae3dd095.
//
// Solidity: function init(string name, string symbol, address admin, address boilerplateAdd, uint256 boilerplateId) returns()
func (_GenerativeParam *GenerativeParamTransactor) Init(opts *bind.TransactOpts, name string, symbol string, admin common.Address, boilerplateAdd common.Address, boilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "init", name, symbol, admin, boilerplateAdd, boilerplateId)
}

// Init is a paid mutator transaction binding the contract method 0xae3dd095.
//
// Solidity: function init(string name, string symbol, address admin, address boilerplateAdd, uint256 boilerplateId) returns()
func (_GenerativeParam *GenerativeParamSession) Init(name string, symbol string, admin common.Address, boilerplateAdd common.Address, boilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Init(&_GenerativeParam.TransactOpts, name, symbol, admin, boilerplateAdd, boilerplateId)
}

// Init is a paid mutator transaction binding the contract method 0xae3dd095.
//
// Solidity: function init(string name, string symbol, address admin, address boilerplateAdd, uint256 boilerplateId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Init(name string, symbol string, admin common.Address, boilerplateAdd common.Address, boilerplateId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Init(&_GenerativeParam.TransactOpts, name, symbol, admin, boilerplateAdd, boilerplateId)
}

// Mint is a paid mutator transaction binding the contract method 0x2f19212f.
//
// Solidity: function mint(address mintTo, address creator, string uri, (bytes32,uint256[]) _paramsTemplateValue) returns()
func (_GenerativeParam *GenerativeParamTransactor) Mint(opts *bind.TransactOpts, mintTo common.Address, creator common.Address, uri string, _paramsTemplateValue BoilerplateParamParamsOfNFT) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "mint", mintTo, creator, uri, _paramsTemplateValue)
}

// Mint is a paid mutator transaction binding the contract method 0x2f19212f.
//
// Solidity: function mint(address mintTo, address creator, string uri, (bytes32,uint256[]) _paramsTemplateValue) returns()
func (_GenerativeParam *GenerativeParamSession) Mint(mintTo common.Address, creator common.Address, uri string, _paramsTemplateValue BoilerplateParamParamsOfNFT) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Mint(&_GenerativeParam.TransactOpts, mintTo, creator, uri, _paramsTemplateValue)
}

// Mint is a paid mutator transaction binding the contract method 0x2f19212f.
//
// Solidity: function mint(address mintTo, address creator, string uri, (bytes32,uint256[]) _paramsTemplateValue) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Mint(mintTo common.Address, creator common.Address, uri string, _paramsTemplateValue BoilerplateParamParamsOfNFT) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Mint(&_GenerativeParam.TransactOpts, mintTo, creator, uri, _paramsTemplateValue)
}

// Mint0 is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeParam *GenerativeParamTransactor) Mint0(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "mint0", to)
}

// Mint0 is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeParam *GenerativeParamSession) Mint0(to common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Mint0(&_GenerativeParam.TransactOpts, to)
}

// Mint0 is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Mint0(to common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.Mint0(&_GenerativeParam.TransactOpts, to)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeParam *GenerativeParamTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeParam *GenerativeParamSession) Pause() (*types.Transaction, error) {
	return _GenerativeParam.Contract.Pause(&_GenerativeParam.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Pause() (*types.Transaction, error) {
	return _GenerativeParam.Contract.Pause(&_GenerativeParam.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeParam *GenerativeParamTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeParam *GenerativeParamSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeParam.Contract.RenounceOwnership(&_GenerativeParam.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeParam *GenerativeParamTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeParam.Contract.RenounceOwnership(&_GenerativeParam.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.RenounceRole(&_GenerativeParam.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.RenounceRole(&_GenerativeParam.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.RevokeRole(&_GenerativeParam.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.RevokeRole(&_GenerativeParam.TransactOpts, role, account)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SafeTransferFrom(&_GenerativeParam.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SafeTransferFrom(&_GenerativeParam.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeParam *GenerativeParamTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeParam *GenerativeParamSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SafeTransferFrom0(&_GenerativeParam.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SafeTransferFrom0(&_GenerativeParam.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeParam *GenerativeParamTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeParam *GenerativeParamSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetApprovalForAll(&_GenerativeParam.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetApprovalForAll(&_GenerativeParam.TransactOpts, operator, approved)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeParam *GenerativeParamTransactor) SetCreator(opts *bind.TransactOpts, _to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "setCreator", _to, _ids)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeParam *GenerativeParamSession) SetCreator(_to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetCreator(&_GenerativeParam.TransactOpts, _to, _ids)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SetCreator(_to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetCreator(&_GenerativeParam.TransactOpts, _to, _ids)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _tokenId, string _newURI) returns()
func (_GenerativeParam *GenerativeParamTransactor) SetCustomURI(opts *bind.TransactOpts, _tokenId *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "setCustomURI", _tokenId, _newURI)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _tokenId, string _newURI) returns()
func (_GenerativeParam *GenerativeParamSession) SetCustomURI(_tokenId *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetCustomURI(&_GenerativeParam.TransactOpts, _tokenId, _newURI)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _tokenId, string _newURI) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SetCustomURI(_tokenId *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetCustomURI(&_GenerativeParam.TransactOpts, _tokenId, _newURI)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeParam *GenerativeParamTransactor) SetTokenRoyalty(opts *bind.TransactOpts, _tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "setTokenRoyalty", _tokenId, _recipient, _value)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeParam *GenerativeParamSession) SetTokenRoyalty(_tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetTokenRoyalty(&_GenerativeParam.TransactOpts, _tokenId, _recipient, _value)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) SetTokenRoyalty(_tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.SetTokenRoyalty(&_GenerativeParam.TransactOpts, _tokenId, _recipient, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.TransferFrom(&_GenerativeParam.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeParam.Contract.TransferFrom(&_GenerativeParam.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeParam *GenerativeParamTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeParam *GenerativeParamSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.TransferOwnership(&_GenerativeParam.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeParam.Contract.TransferOwnership(&_GenerativeParam.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeParam *GenerativeParamTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeParam *GenerativeParamSession) Unpause() (*types.Transaction, error) {
	return _GenerativeParam.Contract.Unpause(&_GenerativeParam.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeParam *GenerativeParamTransactorSession) Unpause() (*types.Transaction, error) {
	return _GenerativeParam.Contract.Unpause(&_GenerativeParam.TransactOpts)
}

// UpdateTraits is a paid mutator transaction binding the contract method 0x81d22564.
//
// Solidity: function updateTraits(((string,string[],uint256,string)[]) traits) returns()
func (_GenerativeParam *GenerativeParamTransactor) UpdateTraits(opts *bind.TransactOpts, traits TraitInfoTraits) (*types.Transaction, error) {
	return _GenerativeParam.contract.Transact(opts, "updateTraits", traits)
}

// UpdateTraits is a paid mutator transaction binding the contract method 0x81d22564.
//
// Solidity: function updateTraits(((string,string[],uint256,string)[]) traits) returns()
func (_GenerativeParam *GenerativeParamSession) UpdateTraits(traits TraitInfoTraits) (*types.Transaction, error) {
	return _GenerativeParam.Contract.UpdateTraits(&_GenerativeParam.TransactOpts, traits)
}

// UpdateTraits is a paid mutator transaction binding the contract method 0x81d22564.
//
// Solidity: function updateTraits(((string,string[],uint256,string)[]) traits) returns()
func (_GenerativeParam *GenerativeParamTransactorSession) UpdateTraits(traits TraitInfoTraits) (*types.Transaction, error) {
	return _GenerativeParam.Contract.UpdateTraits(&_GenerativeParam.TransactOpts, traits)
}

// GenerativeParamApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GenerativeParam contract.
type GenerativeParamApprovalIterator struct {
	Event *GenerativeParamApproval // Event containing the contract specifics and raw log

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
func (it *GenerativeParamApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamApproval)
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
		it.Event = new(GenerativeParamApproval)
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
func (it *GenerativeParamApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamApproval represents a Approval event raised by the GenerativeParam contract.
type GenerativeParamApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeParam *GenerativeParamFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GenerativeParamApprovalIterator, error) {

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

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamApprovalIterator{contract: _GenerativeParam.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeParam *GenerativeParamFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GenerativeParamApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamApproval)
				if err := _GenerativeParam.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParseApproval(log types.Log) (*GenerativeParamApproval, error) {
	event := new(GenerativeParamApproval)
	if err := _GenerativeParam.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GenerativeParam contract.
type GenerativeParamApprovalForAllIterator struct {
	Event *GenerativeParamApprovalForAll // Event containing the contract specifics and raw log

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
func (it *GenerativeParamApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamApprovalForAll)
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
		it.Event = new(GenerativeParamApprovalForAll)
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
func (it *GenerativeParamApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamApprovalForAll represents a ApprovalForAll event raised by the GenerativeParam contract.
type GenerativeParamApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeParam *GenerativeParamFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GenerativeParamApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamApprovalForAllIterator{contract: _GenerativeParam.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeParam *GenerativeParamFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GenerativeParamApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamApprovalForAll)
				if err := _GenerativeParam.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParseApprovalForAll(log types.Log) (*GenerativeParamApprovalForAll, error) {
	event := new(GenerativeParamApprovalForAll)
	if err := _GenerativeParam.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamMintGenerativeNFTIterator is returned from FilterMintGenerativeNFT and is used to iterate over the raw logs and unpacked data for MintGenerativeNFT events raised by the GenerativeParam contract.
type GenerativeParamMintGenerativeNFTIterator struct {
	Event *GenerativeParamMintGenerativeNFT // Event containing the contract specifics and raw log

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
func (it *GenerativeParamMintGenerativeNFTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamMintGenerativeNFT)
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
		it.Event = new(GenerativeParamMintGenerativeNFT)
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
func (it *GenerativeParamMintGenerativeNFTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamMintGenerativeNFTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamMintGenerativeNFT represents a MintGenerativeNFT event raised by the GenerativeParam contract.
type GenerativeParamMintGenerativeNFT struct {
	MintTo  common.Address
	Creator common.Address
	Uri     string
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMintGenerativeNFT is a free log retrieval operation binding the contract event 0xbf83d3c95e3a6b6b6cac45594a672c1d703818ada1304cf6e235fdc0bd6e4371.
//
// Solidity: event MintGenerativeNFT(address mintTo, address creator, string uri, uint256 tokenId)
func (_GenerativeParam *GenerativeParamFilterer) FilterMintGenerativeNFT(opts *bind.FilterOpts) (*GenerativeParamMintGenerativeNFTIterator, error) {

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "MintGenerativeNFT")
	if err != nil {
		return nil, err
	}
	return &GenerativeParamMintGenerativeNFTIterator{contract: _GenerativeParam.contract, event: "MintGenerativeNFT", logs: logs, sub: sub}, nil
}

// WatchMintGenerativeNFT is a free log subscription operation binding the contract event 0xbf83d3c95e3a6b6b6cac45594a672c1d703818ada1304cf6e235fdc0bd6e4371.
//
// Solidity: event MintGenerativeNFT(address mintTo, address creator, string uri, uint256 tokenId)
func (_GenerativeParam *GenerativeParamFilterer) WatchMintGenerativeNFT(opts *bind.WatchOpts, sink chan<- *GenerativeParamMintGenerativeNFT) (event.Subscription, error) {

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "MintGenerativeNFT")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamMintGenerativeNFT)
				if err := _GenerativeParam.contract.UnpackLog(event, "MintGenerativeNFT", log); err != nil {
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

// ParseMintGenerativeNFT is a log parse operation binding the contract event 0xbf83d3c95e3a6b6b6cac45594a672c1d703818ada1304cf6e235fdc0bd6e4371.
//
// Solidity: event MintGenerativeNFT(address mintTo, address creator, string uri, uint256 tokenId)
func (_GenerativeParam *GenerativeParamFilterer) ParseMintGenerativeNFT(log types.Log) (*GenerativeParamMintGenerativeNFT, error) {
	event := new(GenerativeParamMintGenerativeNFT)
	if err := _GenerativeParam.contract.UnpackLog(event, "MintGenerativeNFT", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GenerativeParam contract.
type GenerativeParamOwnershipTransferredIterator struct {
	Event *GenerativeParamOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GenerativeParamOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamOwnershipTransferred)
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
		it.Event = new(GenerativeParamOwnershipTransferred)
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
func (it *GenerativeParamOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamOwnershipTransferred represents a OwnershipTransferred event raised by the GenerativeParam contract.
type GenerativeParamOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeParam *GenerativeParamFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GenerativeParamOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamOwnershipTransferredIterator{contract: _GenerativeParam.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeParam *GenerativeParamFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GenerativeParamOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamOwnershipTransferred)
				if err := _GenerativeParam.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParseOwnershipTransferred(log types.Log) (*GenerativeParamOwnershipTransferred, error) {
	event := new(GenerativeParamOwnershipTransferred)
	if err := _GenerativeParam.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GenerativeParam contract.
type GenerativeParamPausedIterator struct {
	Event *GenerativeParamPaused // Event containing the contract specifics and raw log

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
func (it *GenerativeParamPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamPaused)
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
		it.Event = new(GenerativeParamPaused)
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
func (it *GenerativeParamPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamPaused represents a Paused event raised by the GenerativeParam contract.
type GenerativeParamPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeParam *GenerativeParamFilterer) FilterPaused(opts *bind.FilterOpts) (*GenerativeParamPausedIterator, error) {

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GenerativeParamPausedIterator{contract: _GenerativeParam.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeParam *GenerativeParamFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GenerativeParamPaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamPaused)
				if err := _GenerativeParam.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParsePaused(log types.Log) (*GenerativeParamPaused, error) {
	event := new(GenerativeParamPaused)
	if err := _GenerativeParam.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GenerativeParam contract.
type GenerativeParamRoleAdminChangedIterator struct {
	Event *GenerativeParamRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *GenerativeParamRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamRoleAdminChanged)
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
		it.Event = new(GenerativeParamRoleAdminChanged)
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
func (it *GenerativeParamRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamRoleAdminChanged represents a RoleAdminChanged event raised by the GenerativeParam contract.
type GenerativeParamRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GenerativeParam *GenerativeParamFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GenerativeParamRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamRoleAdminChangedIterator{contract: _GenerativeParam.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GenerativeParam *GenerativeParamFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GenerativeParamRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamRoleAdminChanged)
				if err := _GenerativeParam.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GenerativeParam *GenerativeParamFilterer) ParseRoleAdminChanged(log types.Log) (*GenerativeParamRoleAdminChanged, error) {
	event := new(GenerativeParamRoleAdminChanged)
	if err := _GenerativeParam.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GenerativeParam contract.
type GenerativeParamRoleGrantedIterator struct {
	Event *GenerativeParamRoleGranted // Event containing the contract specifics and raw log

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
func (it *GenerativeParamRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamRoleGranted)
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
		it.Event = new(GenerativeParamRoleGranted)
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
func (it *GenerativeParamRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamRoleGranted represents a RoleGranted event raised by the GenerativeParam contract.
type GenerativeParamRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GenerativeParamRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamRoleGrantedIterator{contract: _GenerativeParam.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GenerativeParamRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamRoleGranted)
				if err := _GenerativeParam.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) ParseRoleGranted(log types.Log) (*GenerativeParamRoleGranted, error) {
	event := new(GenerativeParamRoleGranted)
	if err := _GenerativeParam.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GenerativeParam contract.
type GenerativeParamRoleRevokedIterator struct {
	Event *GenerativeParamRoleRevoked // Event containing the contract specifics and raw log

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
func (it *GenerativeParamRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamRoleRevoked)
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
		it.Event = new(GenerativeParamRoleRevoked)
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
func (it *GenerativeParamRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamRoleRevoked represents a RoleRevoked event raised by the GenerativeParam contract.
type GenerativeParamRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GenerativeParamRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamRoleRevokedIterator{contract: _GenerativeParam.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GenerativeParamRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamRoleRevoked)
				if err := _GenerativeParam.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeParam *GenerativeParamFilterer) ParseRoleRevoked(log types.Log) (*GenerativeParamRoleRevoked, error) {
	event := new(GenerativeParamRoleRevoked)
	if err := _GenerativeParam.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GenerativeParam contract.
type GenerativeParamTransferIterator struct {
	Event *GenerativeParamTransfer // Event containing the contract specifics and raw log

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
func (it *GenerativeParamTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamTransfer)
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
		it.Event = new(GenerativeParamTransfer)
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
func (it *GenerativeParamTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamTransfer represents a Transfer event raised by the GenerativeParam contract.
type GenerativeParamTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeParam *GenerativeParamFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GenerativeParamTransferIterator, error) {

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

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeParamTransferIterator{contract: _GenerativeParam.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeParam *GenerativeParamFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GenerativeParamTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamTransfer)
				if err := _GenerativeParam.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParseTransfer(log types.Log) (*GenerativeParamTransfer, error) {
	event := new(GenerativeParamTransfer)
	if err := _GenerativeParam.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeParamUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GenerativeParam contract.
type GenerativeParamUnpausedIterator struct {
	Event *GenerativeParamUnpaused // Event containing the contract specifics and raw log

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
func (it *GenerativeParamUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeParamUnpaused)
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
		it.Event = new(GenerativeParamUnpaused)
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
func (it *GenerativeParamUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeParamUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeParamUnpaused represents a Unpaused event raised by the GenerativeParam contract.
type GenerativeParamUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeParam *GenerativeParamFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GenerativeParamUnpausedIterator, error) {

	logs, sub, err := _GenerativeParam.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GenerativeParamUnpausedIterator{contract: _GenerativeParam.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeParam *GenerativeParamFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GenerativeParamUnpaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeParam.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeParamUnpaused)
				if err := _GenerativeParam.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_GenerativeParam *GenerativeParamFilterer) ParseUnpaused(log types.Log) (*GenerativeParamUnpaused, error) {
	event := new(GenerativeParamUnpaused)
	if err := _GenerativeParam.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
