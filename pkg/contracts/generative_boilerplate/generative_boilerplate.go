// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_boilerplate

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

// BoilerplateParamParamsOfProject is an auto generated low-level Go binding around an user-defined struct.
type BoilerplateParamParamsOfProject struct {
	Seed   [32]byte
	Params []BoilerplateParamParamTemplate
}

// IGenerativeBoilerplateNFTMintRequest is an auto generated low-level Go binding around an user-defined struct.
type IGenerativeBoilerplateNFTMintRequest struct {
	FromProjectId *big.Int
	MintTo        common.Address
	UriBatch      []string
	ParamsBatch   []BoilerplateParamParamsOfProject
}

// GenerativeBoilerplateMetaData contains all meta data concerning the GenerativeBoilerplate contract.
var GenerativeBoilerplateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"seeds\",\"type\":\"bytes32[]\"}],\"name\":\"GenerateSeeds\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"_fromProjectId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintTo\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_uriBatch\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_typeValue\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_editable\",\"type\":\"bool\"}],\"internalType\":\"structBoilerplateParam.ParamTemplate[]\",\"name\":\"_params\",\"type\":\"tuple[]\"}],\"internalType\":\"structBoilerplateParam.ParamsOfProject[]\",\"name\":\"_paramsBatch\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structIGenerativeBoilerplateNFT.MintRequest\",\"name\":\"request\",\"type\":\"tuple\"}],\"name\":\"MintBatchNFT\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINTER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_approvalForAllSeeds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_minterNFTInfos\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_projects\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_mintMaxSupply\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_mintTotalSupply\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_script\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"_scriptType\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_creator\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_customUri\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_projectName\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"_clientSeed\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_typeValue\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_editable\",\"type\":\"bool\"}],\"internalType\":\"structBoilerplateParam.ParamTemplate[]\",\"name\":\"_params\",\"type\":\"tuple[]\"}],\"internalType\":\"structBoilerplateParam.ParamsOfProject\",\"name\":\"_paramsTemplate\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"approveForAllSeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseTokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"generateSeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getRoleMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleMemberCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseTokenURI\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"baseUri\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"isApprovedOrOwnerForSeed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"_fromProjectId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_mintTo\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"_uriBatch\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_typeValue\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_editable\",\"type\":\"bool\"}],\"internalType\":\"structBoilerplateParam.ParamTemplate[]\",\"name\":\"_params\",\"type\":\"tuple[]\"}],\"internalType\":\"structBoilerplateParam.ParamsOfProject[]\",\"name\":\"_paramsBatch\",\"type\":\"tuple[]\"}],\"internalType\":\"structIGenerativeBoilerplateNFT.MintRequest\",\"name\":\"mintBatch\",\"type\":\"tuple\"}],\"name\":\"mintBatchUniqueNFT\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenID\",\"type\":\"uint256\"}],\"name\":\"mintMaxSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"projectName\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"maxSupply\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"script\",\"type\":\"string\"},{\"internalType\":\"uint32\",\"name\":\"scriptType\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"clientSeed\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"feeAdd\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_seed\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"_typeValue\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_max\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_min\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"_decimal\",\"type\":\"uint8\"},{\"internalType\":\"string[]\",\"name\":\"_availableValues\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_editable\",\"type\":\"bool\"}],\"internalType\":\"structBoilerplateParam.ParamTemplate[]\",\"name\":\"_params\",\"type\":\"tuple[]\"}],\"internalType\":\"structBoilerplateParam.ParamsOfProject\",\"name\":\"paramsTemplate\",\"type\":\"tuple\"}],\"name\":\"mintProject\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenID\",\"type\":\"uint256\"}],\"name\":\"mintTotalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"ownerOfSeed\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"seeds\",\"type\":\"bytes32[]\"}],\"name\":\"registerSeeds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"royalties\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"amount\",\"type\":\"uint24\"},{\"internalType\":\"bool\",\"name\":\"isValue\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"setCreator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"_newURI\",\"type\":\"string\"}],\"name\":\"setCustomURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setTokenRoyalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"transferSeed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"erc20Addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeBoilerplateABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeBoilerplateMetaData.ABI instead.
var GenerativeBoilerplateABI = GenerativeBoilerplateMetaData.ABI

// GenerativeBoilerplate is an auto generated Go binding around an Ethereum contract.
type GenerativeBoilerplate struct {
	GenerativeBoilerplateCaller     // Read-only binding to the contract
	GenerativeBoilerplateTransactor // Write-only binding to the contract
	GenerativeBoilerplateFilterer   // Log filterer for contract events
}

// GenerativeBoilerplateCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeBoilerplateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeBoilerplateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeBoilerplateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeBoilerplateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeBoilerplateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeBoilerplateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeBoilerplateSession struct {
	Contract     *GenerativeBoilerplate // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GenerativeBoilerplateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeBoilerplateCallerSession struct {
	Contract *GenerativeBoilerplateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GenerativeBoilerplateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeBoilerplateTransactorSession struct {
	Contract     *GenerativeBoilerplateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GenerativeBoilerplateRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeBoilerplateRaw struct {
	Contract *GenerativeBoilerplate // Generic contract binding to access the raw methods on
}

// GenerativeBoilerplateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeBoilerplateCallerRaw struct {
	Contract *GenerativeBoilerplateCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeBoilerplateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeBoilerplateTransactorRaw struct {
	Contract *GenerativeBoilerplateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeBoilerplate creates a new instance of GenerativeBoilerplate, bound to a specific deployed contract.
func NewGenerativeBoilerplate(address common.Address, backend bind.ContractBackend) (*GenerativeBoilerplate, error) {
	contract, err := bindGenerativeBoilerplate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplate{GenerativeBoilerplateCaller: GenerativeBoilerplateCaller{contract: contract}, GenerativeBoilerplateTransactor: GenerativeBoilerplateTransactor{contract: contract}, GenerativeBoilerplateFilterer: GenerativeBoilerplateFilterer{contract: contract}}, nil
}

// NewGenerativeBoilerplateCaller creates a new read-only instance of GenerativeBoilerplate, bound to a specific deployed contract.
func NewGenerativeBoilerplateCaller(address common.Address, caller bind.ContractCaller) (*GenerativeBoilerplateCaller, error) {
	contract, err := bindGenerativeBoilerplate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateCaller{contract: contract}, nil
}

// NewGenerativeBoilerplateTransactor creates a new write-only instance of GenerativeBoilerplate, bound to a specific deployed contract.
func NewGenerativeBoilerplateTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeBoilerplateTransactor, error) {
	contract, err := bindGenerativeBoilerplate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateTransactor{contract: contract}, nil
}

// NewGenerativeBoilerplateFilterer creates a new log filterer instance of GenerativeBoilerplate, bound to a specific deployed contract.
func NewGenerativeBoilerplateFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeBoilerplateFilterer, error) {
	contract, err := bindGenerativeBoilerplate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateFilterer{contract: contract}, nil
}

// bindGenerativeBoilerplate binds a generic wrapper to an already deployed contract.
func bindGenerativeBoilerplate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeBoilerplateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeBoilerplate *GenerativeBoilerplateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeBoilerplate.Contract.GenerativeBoilerplateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeBoilerplate *GenerativeBoilerplateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GenerativeBoilerplateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeBoilerplate *GenerativeBoilerplateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GenerativeBoilerplateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeBoilerplate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.DEFAULTADMINROLE(&_GenerativeBoilerplate.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.DEFAULTADMINROLE(&_GenerativeBoilerplate.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) MINTERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "MINTER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MINTERROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.MINTERROLE(&_GenerativeBoilerplate.CallOpts)
}

// MINTERROLE is a free data retrieval call binding the contract method 0xd5391393.
//
// Solidity: function MINTER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) MINTERROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.MINTERROLE(&_GenerativeBoilerplate.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) PAUSERROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.PAUSERROLE(&_GenerativeBoilerplate.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) PAUSERROLE() ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.PAUSERROLE(&_GenerativeBoilerplate.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Admin() (common.Address, error) {
	return _GenerativeBoilerplate.Contract.Admin(&_GenerativeBoilerplate.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Admin() (common.Address, error) {
	return _GenerativeBoilerplate.Contract.Admin(&_GenerativeBoilerplate.CallOpts)
}

// ApprovalForAllSeeds is a free data retrieval call binding the contract method 0x36bd9e63.
//
// Solidity: function _approvalForAllSeeds(address , address ) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) ApprovalForAllSeeds(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "_approvalForAllSeeds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ApprovalForAllSeeds is a free data retrieval call binding the contract method 0x36bd9e63.
//
// Solidity: function _approvalForAllSeeds(address , address ) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) ApprovalForAllSeeds(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.ApprovalForAllSeeds(&_GenerativeBoilerplate.CallOpts, arg0, arg1)
}

// ApprovalForAllSeeds is a free data retrieval call binding the contract method 0x36bd9e63.
//
// Solidity: function _approvalForAllSeeds(address , address ) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) ApprovalForAllSeeds(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.ApprovalForAllSeeds(&_GenerativeBoilerplate.CallOpts, arg0, arg1)
}

// MinterNFTInfos is a free data retrieval call binding the contract method 0x8f1637d7.
//
// Solidity: function _minterNFTInfos(uint256 ) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) MinterNFTInfos(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "_minterNFTInfos", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MinterNFTInfos is a free data retrieval call binding the contract method 0x8f1637d7.
//
// Solidity: function _minterNFTInfos(uint256 ) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MinterNFTInfos(arg0 *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.MinterNFTInfos(&_GenerativeBoilerplate.CallOpts, arg0)
}

// MinterNFTInfos is a free data retrieval call binding the contract method 0x8f1637d7.
//
// Solidity: function _minterNFTInfos(uint256 ) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) MinterNFTInfos(arg0 *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.MinterNFTInfos(&_GenerativeBoilerplate.CallOpts, arg0)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) ParamsAddress() (common.Address, error) {
	return _GenerativeBoilerplate.Contract.ParamsAddress(&_GenerativeBoilerplate.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) ParamsAddress() (common.Address, error) {
	return _GenerativeBoilerplate.Contract.ParamsAddress(&_GenerativeBoilerplate.CallOpts)
}

// Projects is a free data retrieval call binding the contract method 0x821f9dc9.
//
// Solidity: function _projects(uint256 ) view returns(uint256 _fee, address _feeToken, uint256 _mintMaxSupply, uint256 _mintTotalSupply, string _script, uint32 _scriptType, address _creator, string _customUri, string _projectName, bool _clientSeed, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) _paramsTemplate)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Projects(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Fee             *big.Int
	FeeToken        common.Address
	MintMaxSupply   *big.Int
	MintTotalSupply *big.Int
	Script          string
	ScriptType      uint32
	Creator         common.Address
	CustomUri       string
	ProjectName     string
	ClientSeed      bool
	ParamsTemplate  BoilerplateParamParamsOfProject
}, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "_projects", arg0)

	outstruct := new(struct {
		Fee             *big.Int
		FeeToken        common.Address
		MintMaxSupply   *big.Int
		MintTotalSupply *big.Int
		Script          string
		ScriptType      uint32
		Creator         common.Address
		CustomUri       string
		ProjectName     string
		ClientSeed      bool
		ParamsTemplate  BoilerplateParamParamsOfProject
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fee = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeToken = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.MintMaxSupply = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MintTotalSupply = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Script = *abi.ConvertType(out[4], new(string)).(*string)
	outstruct.ScriptType = *abi.ConvertType(out[5], new(uint32)).(*uint32)
	outstruct.Creator = *abi.ConvertType(out[6], new(common.Address)).(*common.Address)
	outstruct.CustomUri = *abi.ConvertType(out[7], new(string)).(*string)
	outstruct.ProjectName = *abi.ConvertType(out[8], new(string)).(*string)
	outstruct.ClientSeed = *abi.ConvertType(out[9], new(bool)).(*bool)
	outstruct.ParamsTemplate = *abi.ConvertType(out[10], new(BoilerplateParamParamsOfProject)).(*BoilerplateParamParamsOfProject)

	return *outstruct, err

}

// Projects is a free data retrieval call binding the contract method 0x821f9dc9.
//
// Solidity: function _projects(uint256 ) view returns(uint256 _fee, address _feeToken, uint256 _mintMaxSupply, uint256 _mintTotalSupply, string _script, uint32 _scriptType, address _creator, string _customUri, string _projectName, bool _clientSeed, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) _paramsTemplate)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Projects(arg0 *big.Int) (struct {
	Fee             *big.Int
	FeeToken        common.Address
	MintMaxSupply   *big.Int
	MintTotalSupply *big.Int
	Script          string
	ScriptType      uint32
	Creator         common.Address
	CustomUri       string
	ProjectName     string
	ClientSeed      bool
	ParamsTemplate  BoilerplateParamParamsOfProject
}, error) {
	return _GenerativeBoilerplate.Contract.Projects(&_GenerativeBoilerplate.CallOpts, arg0)
}

// Projects is a free data retrieval call binding the contract method 0x821f9dc9.
//
// Solidity: function _projects(uint256 ) view returns(uint256 _fee, address _feeToken, uint256 _mintMaxSupply, uint256 _mintTotalSupply, string _script, uint32 _scriptType, address _creator, string _customUri, string _projectName, bool _clientSeed, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) _paramsTemplate)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Projects(arg0 *big.Int) (struct {
	Fee             *big.Int
	FeeToken        common.Address
	MintMaxSupply   *big.Int
	MintTotalSupply *big.Int
	Script          string
	ScriptType      uint32
	Creator         common.Address
	CustomUri       string
	ProjectName     string
	ClientSeed      bool
	ParamsTemplate  BoilerplateParamParamsOfProject
}, error) {
	return _GenerativeBoilerplate.Contract.Projects(&_GenerativeBoilerplate.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.BalanceOf(&_GenerativeBoilerplate.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.BalanceOf(&_GenerativeBoilerplate.CallOpts, owner)
}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) BaseTokenURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "baseTokenURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) BaseTokenURI() (string, error) {
	return _GenerativeBoilerplate.Contract.BaseTokenURI(&_GenerativeBoilerplate.CallOpts)
}

// BaseTokenURI is a free data retrieval call binding the contract method 0xd547cfb7.
//
// Solidity: function baseTokenURI() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) BaseTokenURI() (string, error) {
	return _GenerativeBoilerplate.Contract.BaseTokenURI(&_GenerativeBoilerplate.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 _id) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Exists(opts *bind.CallOpts, _id *big.Int) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "exists", _id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 _id) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Exists(_id *big.Int) (bool, error) {
	return _GenerativeBoilerplate.Contract.Exists(&_GenerativeBoilerplate.CallOpts, _id)
}

// Exists is a free data retrieval call binding the contract method 0x4f558e79.
//
// Solidity: function exists(uint256 _id) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Exists(_id *big.Int) (bool, error) {
	return _GenerativeBoilerplate.Contract.Exists(&_GenerativeBoilerplate.CallOpts, _id)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.GetApproved(&_GenerativeBoilerplate.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.GetApproved(&_GenerativeBoilerplate.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.GetRoleAdmin(&_GenerativeBoilerplate.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _GenerativeBoilerplate.Contract.GetRoleAdmin(&_GenerativeBoilerplate.CallOpts, role)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) GetRoleMember(opts *bind.CallOpts, role [32]byte, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "getRoleMember", role, index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.GetRoleMember(&_GenerativeBoilerplate.CallOpts, role, index)
}

// GetRoleMember is a free data retrieval call binding the contract method 0x9010d07c.
//
// Solidity: function getRoleMember(bytes32 role, uint256 index) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) GetRoleMember(role [32]byte, index *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.GetRoleMember(&_GenerativeBoilerplate.CallOpts, role, index)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) GetRoleMemberCount(opts *bind.CallOpts, role [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "getRoleMemberCount", role)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.GetRoleMemberCount(&_GenerativeBoilerplate.CallOpts, role)
}

// GetRoleMemberCount is a free data retrieval call binding the contract method 0xca15c873.
//
// Solidity: function getRoleMemberCount(bytes32 role) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) GetRoleMemberCount(role [32]byte) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.GetRoleMemberCount(&_GenerativeBoilerplate.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GenerativeBoilerplate.Contract.HasRole(&_GenerativeBoilerplate.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _GenerativeBoilerplate.Contract.HasRole(&_GenerativeBoilerplate.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeBoilerplate.Contract.IsApprovedForAll(&_GenerativeBoilerplate.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _GenerativeBoilerplate.Contract.IsApprovedForAll(&_GenerativeBoilerplate.CallOpts, owner, operator)
}

// IsApprovedOrOwnerForSeed is a free data retrieval call binding the contract method 0xb0f4cf93.
//
// Solidity: function isApprovedOrOwnerForSeed(address operator, bytes32 seed, uint256 projectId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) IsApprovedOrOwnerForSeed(opts *bind.CallOpts, operator common.Address, seed [32]byte, projectId *big.Int) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "isApprovedOrOwnerForSeed", operator, seed, projectId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedOrOwnerForSeed is a free data retrieval call binding the contract method 0xb0f4cf93.
//
// Solidity: function isApprovedOrOwnerForSeed(address operator, bytes32 seed, uint256 projectId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) IsApprovedOrOwnerForSeed(operator common.Address, seed [32]byte, projectId *big.Int) (bool, error) {
	return _GenerativeBoilerplate.Contract.IsApprovedOrOwnerForSeed(&_GenerativeBoilerplate.CallOpts, operator, seed, projectId)
}

// IsApprovedOrOwnerForSeed is a free data retrieval call binding the contract method 0xb0f4cf93.
//
// Solidity: function isApprovedOrOwnerForSeed(address operator, bytes32 seed, uint256 projectId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) IsApprovedOrOwnerForSeed(operator common.Address, seed [32]byte, projectId *big.Int) (bool, error) {
	return _GenerativeBoilerplate.Contract.IsApprovedOrOwnerForSeed(&_GenerativeBoilerplate.CallOpts, operator, seed, projectId)
}

// MintMaxSupply is a free data retrieval call binding the contract method 0x030595d7.
//
// Solidity: function mintMaxSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) MintMaxSupply(opts *bind.CallOpts, _tokenID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "mintMaxSupply", _tokenID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintMaxSupply is a free data retrieval call binding the contract method 0x030595d7.
//
// Solidity: function mintMaxSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MintMaxSupply(_tokenID *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.MintMaxSupply(&_GenerativeBoilerplate.CallOpts, _tokenID)
}

// MintMaxSupply is a free data retrieval call binding the contract method 0x030595d7.
//
// Solidity: function mintMaxSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) MintMaxSupply(_tokenID *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.MintMaxSupply(&_GenerativeBoilerplate.CallOpts, _tokenID)
}

// MintTotalSupply is a free data retrieval call binding the contract method 0x3fbdb0c2.
//
// Solidity: function mintTotalSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) MintTotalSupply(opts *bind.CallOpts, _tokenID *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "mintTotalSupply", _tokenID)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MintTotalSupply is a free data retrieval call binding the contract method 0x3fbdb0c2.
//
// Solidity: function mintTotalSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MintTotalSupply(_tokenID *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.MintTotalSupply(&_GenerativeBoilerplate.CallOpts, _tokenID)
}

// MintTotalSupply is a free data retrieval call binding the contract method 0x3fbdb0c2.
//
// Solidity: function mintTotalSupply(uint256 _tokenID) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) MintTotalSupply(_tokenID *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.MintTotalSupply(&_GenerativeBoilerplate.CallOpts, _tokenID)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Name() (string, error) {
	return _GenerativeBoilerplate.Contract.Name(&_GenerativeBoilerplate.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Name() (string, error) {
	return _GenerativeBoilerplate.Contract.Name(&_GenerativeBoilerplate.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.OwnerOf(&_GenerativeBoilerplate.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.OwnerOf(&_GenerativeBoilerplate.CallOpts, tokenId)
}

// OwnerOfSeed is a free data retrieval call binding the contract method 0x1697eccf.
//
// Solidity: function ownerOfSeed(bytes32 seed, uint256 projectId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) OwnerOfSeed(opts *bind.CallOpts, seed [32]byte, projectId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "ownerOfSeed", seed, projectId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOfSeed is a free data retrieval call binding the contract method 0x1697eccf.
//
// Solidity: function ownerOfSeed(bytes32 seed, uint256 projectId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) OwnerOfSeed(seed [32]byte, projectId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.OwnerOfSeed(&_GenerativeBoilerplate.CallOpts, seed, projectId)
}

// OwnerOfSeed is a free data retrieval call binding the contract method 0x1697eccf.
//
// Solidity: function ownerOfSeed(bytes32 seed, uint256 projectId) view returns(address)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) OwnerOfSeed(seed [32]byte, projectId *big.Int) (common.Address, error) {
	return _GenerativeBoilerplate.Contract.OwnerOfSeed(&_GenerativeBoilerplate.CallOpts, seed, projectId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Paused() (bool, error) {
	return _GenerativeBoilerplate.Contract.Paused(&_GenerativeBoilerplate.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Paused() (bool, error) {
	return _GenerativeBoilerplate.Contract.Paused(&_GenerativeBoilerplate.CallOpts)
}

// Royalties is a free data retrieval call binding the contract method 0x7f77f574.
//
// Solidity: function royalties(uint256 ) view returns(address recipient, uint24 amount, bool isValue)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Royalties(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "royalties", arg0)

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
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Royalties(arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	return _GenerativeBoilerplate.Contract.Royalties(&_GenerativeBoilerplate.CallOpts, arg0)
}

// Royalties is a free data retrieval call binding the contract method 0x7f77f574.
//
// Solidity: function royalties(uint256 ) view returns(address recipient, uint24 amount, bool isValue)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Royalties(arg0 *big.Int) (struct {
	Recipient common.Address
	Amount    *big.Int
	IsValue   bool
}, error) {
	return _GenerativeBoilerplate.Contract.Royalties(&_GenerativeBoilerplate.CallOpts, arg0)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

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
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeBoilerplate.Contract.RoyaltyInfo(&_GenerativeBoilerplate.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _GenerativeBoilerplate.Contract.RoyaltyInfo(&_GenerativeBoilerplate.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeBoilerplate.Contract.SupportsInterface(&_GenerativeBoilerplate.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeBoilerplate.Contract.SupportsInterface(&_GenerativeBoilerplate.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Symbol() (string, error) {
	return _GenerativeBoilerplate.Contract.Symbol(&_GenerativeBoilerplate.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) Symbol() (string, error) {
	return _GenerativeBoilerplate.Contract.Symbol(&_GenerativeBoilerplate.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TokenByIndex(&_GenerativeBoilerplate.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TokenByIndex(&_GenerativeBoilerplate.CallOpts, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TokenOfOwnerByIndex(&_GenerativeBoilerplate.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TokenOfOwnerByIndex(&_GenerativeBoilerplate.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) TokenURI(opts *bind.CallOpts, _tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "tokenURI", _tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _GenerativeBoilerplate.Contract.TokenURI(&_GenerativeBoilerplate.CallOpts, _tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 _tokenId) view returns(string)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) TokenURI(_tokenId *big.Int) (string, error) {
	return _GenerativeBoilerplate.Contract.TokenURI(&_GenerativeBoilerplate.CallOpts, _tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeBoilerplate.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TotalSupply() (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TotalSupply(&_GenerativeBoilerplate.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateCallerSession) TotalSupply() (*big.Int, error) {
	return _GenerativeBoilerplate.Contract.TotalSupply(&_GenerativeBoilerplate.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Approve(&_GenerativeBoilerplate.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Approve(&_GenerativeBoilerplate.TransactOpts, to, tokenId)
}

// ApproveForAllSeeds is a paid mutator transaction binding the contract method 0x62fe4e4a.
//
// Solidity: function approveForAllSeeds(address operator, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) ApproveForAllSeeds(opts *bind.TransactOpts, operator common.Address, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "approveForAllSeeds", operator, projectId)
}

// ApproveForAllSeeds is a paid mutator transaction binding the contract method 0x62fe4e4a.
//
// Solidity: function approveForAllSeeds(address operator, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) ApproveForAllSeeds(operator common.Address, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.ApproveForAllSeeds(&_GenerativeBoilerplate.TransactOpts, operator, projectId)
}

// ApproveForAllSeeds is a paid mutator transaction binding the contract method 0x62fe4e4a.
//
// Solidity: function approveForAllSeeds(address operator, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) ApproveForAllSeeds(operator common.Address, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.ApproveForAllSeeds(&_GenerativeBoilerplate.TransactOpts, operator, projectId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Burn(&_GenerativeBoilerplate.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Burn(&_GenerativeBoilerplate.TransactOpts, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.ChangeAdmin(&_GenerativeBoilerplate.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.ChangeAdmin(&_GenerativeBoilerplate.TransactOpts, newAdm)
}

// GenerateSeeds is a paid mutator transaction binding the contract method 0x46398139.
//
// Solidity: function generateSeeds(uint256 projectId, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) GenerateSeeds(opts *bind.TransactOpts, projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "generateSeeds", projectId, amount)
}

// GenerateSeeds is a paid mutator transaction binding the contract method 0x46398139.
//
// Solidity: function generateSeeds(uint256 projectId, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GenerateSeeds(projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GenerateSeeds(&_GenerativeBoilerplate.TransactOpts, projectId, amount)
}

// GenerateSeeds is a paid mutator transaction binding the contract method 0x46398139.
//
// Solidity: function generateSeeds(uint256 projectId, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) GenerateSeeds(projectId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GenerateSeeds(&_GenerativeBoilerplate.TransactOpts, projectId, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GrantRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.GrantRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6487c53.
//
// Solidity: function initialize(string name, string symbol, string baseTokenURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, baseTokenURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "initialize", name, symbol, baseTokenURI)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6487c53.
//
// Solidity: function initialize(string name, string symbol, string baseTokenURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Initialize(name string, symbol string, baseTokenURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Initialize(&_GenerativeBoilerplate.TransactOpts, name, symbol, baseTokenURI)
}

// Initialize is a paid mutator transaction binding the contract method 0xa6487c53.
//
// Solidity: function initialize(string name, string symbol, string baseTokenURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Initialize(name string, symbol string, baseTokenURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Initialize(&_GenerativeBoilerplate.TransactOpts, name, symbol, baseTokenURI)
}

// Initialize0 is a paid mutator transaction binding the contract method 0xd6d0faee.
//
// Solidity: function initialize(string name, string symbol, string baseUri, address admin, address paramsAddress) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Initialize0(opts *bind.TransactOpts, name string, symbol string, baseUri string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "initialize0", name, symbol, baseUri, admin, paramsAddress)
}

// Initialize0 is a paid mutator transaction binding the contract method 0xd6d0faee.
//
// Solidity: function initialize(string name, string symbol, string baseUri, address admin, address paramsAddress) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Initialize0(name string, symbol string, baseUri string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Initialize0(&_GenerativeBoilerplate.TransactOpts, name, symbol, baseUri, admin, paramsAddress)
}

// Initialize0 is a paid mutator transaction binding the contract method 0xd6d0faee.
//
// Solidity: function initialize(string name, string symbol, string baseUri, address admin, address paramsAddress) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Initialize0(name string, symbol string, baseUri string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Initialize0(&_GenerativeBoilerplate.TransactOpts, name, symbol, baseUri, admin, paramsAddress)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Mint(to common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Mint(&_GenerativeBoilerplate.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Mint(&_GenerativeBoilerplate.TransactOpts, to)
}

// MintBatchUniqueNFT is a paid mutator transaction binding the contract method 0x6006e9ef.
//
// Solidity: function mintBatchUniqueNFT((uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) mintBatch) payable returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) MintBatchUniqueNFT(opts *bind.TransactOpts, mintBatch IGenerativeBoilerplateNFTMintRequest) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "mintBatchUniqueNFT", mintBatch)
}

// MintBatchUniqueNFT is a paid mutator transaction binding the contract method 0x6006e9ef.
//
// Solidity: function mintBatchUniqueNFT((uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) mintBatch) payable returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MintBatchUniqueNFT(mintBatch IGenerativeBoilerplateNFTMintRequest) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.MintBatchUniqueNFT(&_GenerativeBoilerplate.TransactOpts, mintBatch)
}

// MintBatchUniqueNFT is a paid mutator transaction binding the contract method 0x6006e9ef.
//
// Solidity: function mintBatchUniqueNFT((uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) mintBatch) payable returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) MintBatchUniqueNFT(mintBatch IGenerativeBoilerplateNFTMintRequest) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.MintBatchUniqueNFT(&_GenerativeBoilerplate.TransactOpts, mintBatch)
}

// MintProject is a paid mutator transaction binding the contract method 0x716013fe.
//
// Solidity: function mintProject(address to, string projectName, uint256 maxSupply, string script, uint32 scriptType, bool clientSeed, string uri, uint256 fee, address feeAdd, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) paramsTemplate) payable returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) MintProject(opts *bind.TransactOpts, to common.Address, projectName string, maxSupply *big.Int, script string, scriptType uint32, clientSeed bool, uri string, fee *big.Int, feeAdd common.Address, paramsTemplate BoilerplateParamParamsOfProject) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "mintProject", to, projectName, maxSupply, script, scriptType, clientSeed, uri, fee, feeAdd, paramsTemplate)
}

// MintProject is a paid mutator transaction binding the contract method 0x716013fe.
//
// Solidity: function mintProject(address to, string projectName, uint256 maxSupply, string script, uint32 scriptType, bool clientSeed, string uri, uint256 fee, address feeAdd, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) paramsTemplate) payable returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) MintProject(to common.Address, projectName string, maxSupply *big.Int, script string, scriptType uint32, clientSeed bool, uri string, fee *big.Int, feeAdd common.Address, paramsTemplate BoilerplateParamParamsOfProject) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.MintProject(&_GenerativeBoilerplate.TransactOpts, to, projectName, maxSupply, script, scriptType, clientSeed, uri, fee, feeAdd, paramsTemplate)
}

// MintProject is a paid mutator transaction binding the contract method 0x716013fe.
//
// Solidity: function mintProject(address to, string projectName, uint256 maxSupply, string script, uint32 scriptType, bool clientSeed, string uri, uint256 fee, address feeAdd, (bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[]) paramsTemplate) payable returns(uint256)
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) MintProject(to common.Address, projectName string, maxSupply *big.Int, script string, scriptType uint32, clientSeed bool, uri string, fee *big.Int, feeAdd common.Address, paramsTemplate BoilerplateParamParamsOfProject) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.MintProject(&_GenerativeBoilerplate.TransactOpts, to, projectName, maxSupply, script, scriptType, clientSeed, uri, fee, feeAdd, paramsTemplate)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Pause() (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Pause(&_GenerativeBoilerplate.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Pause() (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Pause(&_GenerativeBoilerplate.TransactOpts)
}

// RegisterSeeds is a paid mutator transaction binding the contract method 0x189c168d.
//
// Solidity: function registerSeeds(uint256 projectId, bytes32[] seeds) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) RegisterSeeds(opts *bind.TransactOpts, projectId *big.Int, seeds [][32]byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "registerSeeds", projectId, seeds)
}

// RegisterSeeds is a paid mutator transaction binding the contract method 0x189c168d.
//
// Solidity: function registerSeeds(uint256 projectId, bytes32[] seeds) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) RegisterSeeds(projectId *big.Int, seeds [][32]byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RegisterSeeds(&_GenerativeBoilerplate.TransactOpts, projectId, seeds)
}

// RegisterSeeds is a paid mutator transaction binding the contract method 0x189c168d.
//
// Solidity: function registerSeeds(uint256 projectId, bytes32[] seeds) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) RegisterSeeds(projectId *big.Int, seeds [][32]byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RegisterSeeds(&_GenerativeBoilerplate.TransactOpts, projectId, seeds)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RenounceRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RenounceRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RevokeRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.RevokeRole(&_GenerativeBoilerplate.TransactOpts, role, account)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SafeTransferFrom(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SafeTransferFrom(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SafeTransferFrom0(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SafeTransferFrom0(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetApprovalForAll(&_GenerativeBoilerplate.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetApprovalForAll(&_GenerativeBoilerplate.TransactOpts, operator, approved)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SetCreator(opts *bind.TransactOpts, _to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "setCreator", _to, _ids)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SetCreator(_to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetCreator(&_GenerativeBoilerplate.TransactOpts, _to, _ids)
}

// SetCreator is a paid mutator transaction binding the contract method 0xd2a6b51a.
//
// Solidity: function setCreator(address _to, uint256[] _ids) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SetCreator(_to common.Address, _ids []*big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetCreator(&_GenerativeBoilerplate.TransactOpts, _to, _ids)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _id, string _newURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SetCustomURI(opts *bind.TransactOpts, _id *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "setCustomURI", _id, _newURI)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _id, string _newURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SetCustomURI(_id *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetCustomURI(&_GenerativeBoilerplate.TransactOpts, _id, _newURI)
}

// SetCustomURI is a paid mutator transaction binding the contract method 0x3adf80b4.
//
// Solidity: function setCustomURI(uint256 _id, string _newURI) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SetCustomURI(_id *big.Int, _newURI string) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetCustomURI(&_GenerativeBoilerplate.TransactOpts, _id, _newURI)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) SetTokenRoyalty(opts *bind.TransactOpts, _tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "setTokenRoyalty", _tokenId, _recipient, _value)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) SetTokenRoyalty(_tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetTokenRoyalty(&_GenerativeBoilerplate.TransactOpts, _tokenId, _recipient, _value)
}

// SetTokenRoyalty is a paid mutator transaction binding the contract method 0x9713c807.
//
// Solidity: function setTokenRoyalty(uint256 _tokenId, address _recipient, uint256 _value) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) SetTokenRoyalty(_tokenId *big.Int, _recipient common.Address, _value *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.SetTokenRoyalty(&_GenerativeBoilerplate.TransactOpts, _tokenId, _recipient, _value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.TransferFrom(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.TransferFrom(&_GenerativeBoilerplate.TransactOpts, from, to, tokenId)
}

// TransferSeed is a paid mutator transaction binding the contract method 0x6e4d0fe6.
//
// Solidity: function transferSeed(address from, address to, bytes32 seed, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) TransferSeed(opts *bind.TransactOpts, from common.Address, to common.Address, seed [32]byte, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "transferSeed", from, to, seed, projectId)
}

// TransferSeed is a paid mutator transaction binding the contract method 0x6e4d0fe6.
//
// Solidity: function transferSeed(address from, address to, bytes32 seed, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) TransferSeed(from common.Address, to common.Address, seed [32]byte, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.TransferSeed(&_GenerativeBoilerplate.TransactOpts, from, to, seed, projectId)
}

// TransferSeed is a paid mutator transaction binding the contract method 0x6e4d0fe6.
//
// Solidity: function transferSeed(address from, address to, bytes32 seed, uint256 projectId) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) TransferSeed(from common.Address, to common.Address, seed [32]byte, projectId *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.TransferSeed(&_GenerativeBoilerplate.TransactOpts, from, to, seed, projectId)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Unpause() (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Unpause(&_GenerativeBoilerplate.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Unpause() (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Unpause(&_GenerativeBoilerplate.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactor) Withdraw(opts *bind.TransactOpts, receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.contract.Transact(opts, "withdraw", receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Withdraw(&_GenerativeBoilerplate.TransactOpts, receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeBoilerplate *GenerativeBoilerplateTransactorSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeBoilerplate.Contract.Withdraw(&_GenerativeBoilerplate.TransactOpts, receiver, erc20Addr, amount)
}

// GenerativeBoilerplateApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateApprovalIterator struct {
	Event *GenerativeBoilerplateApproval // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateApproval)
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
		it.Event = new(GenerativeBoilerplateApproval)
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
func (it *GenerativeBoilerplateApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateApproval represents a Approval event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*GenerativeBoilerplateApprovalIterator, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateApprovalIterator{contract: _GenerativeBoilerplate.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateApproval)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseApproval(log types.Log) (*GenerativeBoilerplateApproval, error) {
	event := new(GenerativeBoilerplateApproval)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateApprovalForAllIterator struct {
	Event *GenerativeBoilerplateApprovalForAll // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateApprovalForAll)
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
		it.Event = new(GenerativeBoilerplateApprovalForAll)
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
func (it *GenerativeBoilerplateApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateApprovalForAll represents a ApprovalForAll event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*GenerativeBoilerplateApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateApprovalForAllIterator{contract: _GenerativeBoilerplate.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateApprovalForAll)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseApprovalForAll(log types.Log) (*GenerativeBoilerplateApprovalForAll, error) {
	event := new(GenerativeBoilerplateApprovalForAll)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateGenerateSeedsIterator is returned from FilterGenerateSeeds and is used to iterate over the raw logs and unpacked data for GenerateSeeds events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateGenerateSeedsIterator struct {
	Event *GenerativeBoilerplateGenerateSeeds // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateGenerateSeedsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateGenerateSeeds)
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
		it.Event = new(GenerativeBoilerplateGenerateSeeds)
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
func (it *GenerativeBoilerplateGenerateSeedsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateGenerateSeedsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateGenerateSeeds represents a GenerateSeeds event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateGenerateSeeds struct {
	Sender    common.Address
	ProjectId *big.Int
	Seeds     [][32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterGenerateSeeds is a free log retrieval operation binding the contract event 0xeac055ed8b84f6ed7b96c917e1a500965c90467998c41630ebabcc5bc99b306b.
//
// Solidity: event GenerateSeeds(address sender, uint256 projectId, bytes32[] seeds)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterGenerateSeeds(opts *bind.FilterOpts) (*GenerativeBoilerplateGenerateSeedsIterator, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "GenerateSeeds")
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateGenerateSeedsIterator{contract: _GenerativeBoilerplate.contract, event: "GenerateSeeds", logs: logs, sub: sub}, nil
}

// WatchGenerateSeeds is a free log subscription operation binding the contract event 0xeac055ed8b84f6ed7b96c917e1a500965c90467998c41630ebabcc5bc99b306b.
//
// Solidity: event GenerateSeeds(address sender, uint256 projectId, bytes32[] seeds)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchGenerateSeeds(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateGenerateSeeds) (event.Subscription, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "GenerateSeeds")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateGenerateSeeds)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "GenerateSeeds", log); err != nil {
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

// ParseGenerateSeeds is a log parse operation binding the contract event 0xeac055ed8b84f6ed7b96c917e1a500965c90467998c41630ebabcc5bc99b306b.
//
// Solidity: event GenerateSeeds(address sender, uint256 projectId, bytes32[] seeds)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseGenerateSeeds(log types.Log) (*GenerativeBoilerplateGenerateSeeds, error) {
	event := new(GenerativeBoilerplateGenerateSeeds)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "GenerateSeeds", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateInitializedIterator struct {
	Event *GenerativeBoilerplateInitialized // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateInitialized)
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
		it.Event = new(GenerativeBoilerplateInitialized)
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
func (it *GenerativeBoilerplateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateInitialized represents a Initialized event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterInitialized(opts *bind.FilterOpts) (*GenerativeBoilerplateInitializedIterator, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateInitializedIterator{contract: _GenerativeBoilerplate.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateInitialized) (event.Subscription, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateInitialized)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseInitialized(log types.Log) (*GenerativeBoilerplateInitialized, error) {
	event := new(GenerativeBoilerplateInitialized)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateMintBatchNFTIterator is returned from FilterMintBatchNFT and is used to iterate over the raw logs and unpacked data for MintBatchNFT events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateMintBatchNFTIterator struct {
	Event *GenerativeBoilerplateMintBatchNFT // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateMintBatchNFTIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateMintBatchNFT)
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
		it.Event = new(GenerativeBoilerplateMintBatchNFT)
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
func (it *GenerativeBoilerplateMintBatchNFTIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateMintBatchNFTIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateMintBatchNFT represents a MintBatchNFT event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateMintBatchNFT struct {
	Sender  common.Address
	Request IGenerativeBoilerplateNFTMintRequest
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMintBatchNFT is a free log retrieval operation binding the contract event 0x67175103d209cf669be17e253aba378479477b5a1ef3b1e9eaedce52f0ebab54.
//
// Solidity: event MintBatchNFT(address sender, (uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) request)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterMintBatchNFT(opts *bind.FilterOpts) (*GenerativeBoilerplateMintBatchNFTIterator, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "MintBatchNFT")
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateMintBatchNFTIterator{contract: _GenerativeBoilerplate.contract, event: "MintBatchNFT", logs: logs, sub: sub}, nil
}

// WatchMintBatchNFT is a free log subscription operation binding the contract event 0x67175103d209cf669be17e253aba378479477b5a1ef3b1e9eaedce52f0ebab54.
//
// Solidity: event MintBatchNFT(address sender, (uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) request)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchMintBatchNFT(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateMintBatchNFT) (event.Subscription, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "MintBatchNFT")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateMintBatchNFT)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "MintBatchNFT", log); err != nil {
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

// ParseMintBatchNFT is a log parse operation binding the contract event 0x67175103d209cf669be17e253aba378479477b5a1ef3b1e9eaedce52f0ebab54.
//
// Solidity: event MintBatchNFT(address sender, (uint256,address,string[],(bytes32,(uint8,uint256,uint256,uint8,string[],uint256,bool)[])[]) request)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseMintBatchNFT(log types.Log) (*GenerativeBoilerplateMintBatchNFT, error) {
	event := new(GenerativeBoilerplateMintBatchNFT)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "MintBatchNFT", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplatePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplatePausedIterator struct {
	Event *GenerativeBoilerplatePaused // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplatePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplatePaused)
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
		it.Event = new(GenerativeBoilerplatePaused)
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
func (it *GenerativeBoilerplatePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplatePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplatePaused represents a Paused event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplatePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterPaused(opts *bind.FilterOpts) (*GenerativeBoilerplatePausedIterator, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplatePausedIterator{contract: _GenerativeBoilerplate.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplatePaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplatePaused)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParsePaused(log types.Log) (*GenerativeBoilerplatePaused, error) {
	event := new(GenerativeBoilerplatePaused)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleAdminChangedIterator struct {
	Event *GenerativeBoilerplateRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateRoleAdminChanged)
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
		it.Event = new(GenerativeBoilerplateRoleAdminChanged)
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
func (it *GenerativeBoilerplateRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateRoleAdminChanged represents a RoleAdminChanged event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*GenerativeBoilerplateRoleAdminChangedIterator, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateRoleAdminChangedIterator{contract: _GenerativeBoilerplate.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateRoleAdminChanged)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseRoleAdminChanged(log types.Log) (*GenerativeBoilerplateRoleAdminChanged, error) {
	event := new(GenerativeBoilerplateRoleAdminChanged)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleGrantedIterator struct {
	Event *GenerativeBoilerplateRoleGranted // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateRoleGranted)
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
		it.Event = new(GenerativeBoilerplateRoleGranted)
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
func (it *GenerativeBoilerplateRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateRoleGranted represents a RoleGranted event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GenerativeBoilerplateRoleGrantedIterator, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateRoleGrantedIterator{contract: _GenerativeBoilerplate.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateRoleGranted)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseRoleGranted(log types.Log) (*GenerativeBoilerplateRoleGranted, error) {
	event := new(GenerativeBoilerplateRoleGranted)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleRevokedIterator struct {
	Event *GenerativeBoilerplateRoleRevoked // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateRoleRevoked)
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
		it.Event = new(GenerativeBoilerplateRoleRevoked)
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
func (it *GenerativeBoilerplateRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateRoleRevoked represents a RoleRevoked event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*GenerativeBoilerplateRoleRevokedIterator, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateRoleRevokedIterator{contract: _GenerativeBoilerplate.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateRoleRevoked)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseRoleRevoked(log types.Log) (*GenerativeBoilerplateRoleRevoked, error) {
	event := new(GenerativeBoilerplateRoleRevoked)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateTransferIterator struct {
	Event *GenerativeBoilerplateTransfer // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateTransfer)
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
		it.Event = new(GenerativeBoilerplateTransfer)
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
func (it *GenerativeBoilerplateTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateTransfer represents a Transfer event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*GenerativeBoilerplateTransferIterator, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateTransferIterator{contract: _GenerativeBoilerplate.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateTransfer)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseTransfer(log types.Log) (*GenerativeBoilerplateTransfer, error) {
	event := new(GenerativeBoilerplateTransfer)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeBoilerplateUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateUnpausedIterator struct {
	Event *GenerativeBoilerplateUnpaused // Event containing the contract specifics and raw log

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
func (it *GenerativeBoilerplateUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeBoilerplateUnpaused)
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
		it.Event = new(GenerativeBoilerplateUnpaused)
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
func (it *GenerativeBoilerplateUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeBoilerplateUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeBoilerplateUnpaused represents a Unpaused event raised by the GenerativeBoilerplate contract.
type GenerativeBoilerplateUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) FilterUnpaused(opts *bind.FilterOpts) (*GenerativeBoilerplateUnpausedIterator, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &GenerativeBoilerplateUnpausedIterator{contract: _GenerativeBoilerplate.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *GenerativeBoilerplateUnpaused) (event.Subscription, error) {

	logs, sub, err := _GenerativeBoilerplate.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeBoilerplateUnpaused)
				if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_GenerativeBoilerplate *GenerativeBoilerplateFilterer) ParseUnpaused(log types.Log) (*GenerativeBoilerplateUnpaused, error) {
	event := new(GenerativeBoilerplateUnpaused)
	if err := _GenerativeBoilerplate.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
