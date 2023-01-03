// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_marketplace

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
)

// MarketplaceListingTokenData is an auto generated low-level Go binding around an user-defined struct.
type MarketplaceListingTokenData struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}

// MarketplaceMakeOfferData is an auto generated low-level Go binding around an user-defined struct.
type MarketplaceMakeOfferData struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}

// GenerativeMarketplaceMetaData contains all meta data concerning the GenerativeMarketplace contract.
var GenerativeMarketplaceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_allowableERC20MakeListToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_allowableERC20MakeOffer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_arrayListingId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_arrayMakeOfferId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_listingTokenDataMapping\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_listingTokenIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_listingTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_makeOfferDataMapping\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"_makeOfferTokenIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"_makeOfferTokens\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_parameterAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"offerId\",\"type\":\"bytes32\"}],\"name\":\"acceptMakeOffer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_offeringId\",\"type\":\"bytes32\"}],\"name\":\"cancelListing\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"offerId\",\"type\":\"bytes32\"}],\"name\":\"cancelMakeOffer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddr\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"parameterControl\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"internalType\":\"structMarketplace.ListingTokenData\",\"name\":\"listingData\",\"type\":\"tuple\"}],\"name\":\"listToken\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"internalType\":\"structMarketplace.MakeOfferData\",\"name\":\"makeOfferData\",\"type\":\"tuple\"}],\"name\":\"makeOffer\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"}],\"name\":\"purchaseToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allow\",\"type\":\"bool\"}],\"name\":\"setApproveERC20ListToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc20\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allow\",\"type\":\"bool\"}],\"name\":\"setApproveERC20MakeOffer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"erc20Addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeMarketplaceABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeMarketplaceMetaData.ABI instead.
var GenerativeMarketplaceABI = GenerativeMarketplaceMetaData.ABI

// GenerativeMarketplace is an auto generated Go binding around an Ethereum contract.
type GenerativeMarketplace struct {
	GenerativeMarketplaceCaller     // Read-only binding to the contract
	GenerativeMarketplaceTransactor // Write-only binding to the contract
	GenerativeMarketplaceFilterer   // Log filterer for contract events
}

// GenerativeMarketplaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeMarketplaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeMarketplaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeMarketplaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeMarketplaceSession struct {
	Contract     *GenerativeMarketplace // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GenerativeMarketplaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeMarketplaceCallerSession struct {
	Contract *GenerativeMarketplaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GenerativeMarketplaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeMarketplaceTransactorSession struct {
	Contract     *GenerativeMarketplaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GenerativeMarketplaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeMarketplaceRaw struct {
	Contract *GenerativeMarketplace // Generic contract binding to access the raw methods on
}

// GenerativeMarketplaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeMarketplaceCallerRaw struct {
	Contract *GenerativeMarketplaceCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeMarketplaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeMarketplaceTransactorRaw struct {
	Contract *GenerativeMarketplaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeMarketplace creates a new instance of GenerativeMarketplace, bound to a specific deployed contract.
func NewGenerativeMarketplace(address common.Address, backend bind.ContractBackend) (*GenerativeMarketplace, error) {
	contract, err := bindGenerativeMarketplace(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplace{GenerativeMarketplaceCaller: GenerativeMarketplaceCaller{contract: contract}, GenerativeMarketplaceTransactor: GenerativeMarketplaceTransactor{contract: contract}, GenerativeMarketplaceFilterer: GenerativeMarketplaceFilterer{contract: contract}}, nil
}

// NewGenerativeMarketplaceCaller creates a new read-only instance of GenerativeMarketplace, bound to a specific deployed contract.
func NewGenerativeMarketplaceCaller(address common.Address, caller bind.ContractCaller) (*GenerativeMarketplaceCaller, error) {
	contract, err := bindGenerativeMarketplace(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceCaller{contract: contract}, nil
}

// NewGenerativeMarketplaceTransactor creates a new write-only instance of GenerativeMarketplace, bound to a specific deployed contract.
func NewGenerativeMarketplaceTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeMarketplaceTransactor, error) {
	contract, err := bindGenerativeMarketplace(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceTransactor{contract: contract}, nil
}

// NewGenerativeMarketplaceFilterer creates a new log filterer instance of GenerativeMarketplace, bound to a specific deployed contract.
func NewGenerativeMarketplaceFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeMarketplaceFilterer, error) {
	contract, err := bindGenerativeMarketplace(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceFilterer{contract: contract}, nil
}

// bindGenerativeMarketplace binds a generic wrapper to an already deployed contract.
func bindGenerativeMarketplace(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(GenerativeMarketplaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeMarketplace *GenerativeMarketplaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeMarketplace.Contract.GenerativeMarketplaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeMarketplace *GenerativeMarketplaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.GenerativeMarketplaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeMarketplace *GenerativeMarketplaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.GenerativeMarketplaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeMarketplace *GenerativeMarketplaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeMarketplace.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) Admin() (common.Address, error) {
	return _GenerativeMarketplace.Contract.Admin(&_GenerativeMarketplace.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) Admin() (common.Address, error) {
	return _GenerativeMarketplace.Contract.Admin(&_GenerativeMarketplace.CallOpts)
}

// AllowableERC20MakeListToken is a free data retrieval call binding the contract method 0x5b06e64a.
//
// Solidity: function _allowableERC20MakeListToken(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) AllowableERC20MakeListToken(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_allowableERC20MakeListToken", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowableERC20MakeListToken is a free data retrieval call binding the contract method 0x5b06e64a.
//
// Solidity: function _allowableERC20MakeListToken(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) AllowableERC20MakeListToken(arg0 common.Address) (bool, error) {
	return _GenerativeMarketplace.Contract.AllowableERC20MakeListToken(&_GenerativeMarketplace.CallOpts, arg0)
}

// AllowableERC20MakeListToken is a free data retrieval call binding the contract method 0x5b06e64a.
//
// Solidity: function _allowableERC20MakeListToken(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) AllowableERC20MakeListToken(arg0 common.Address) (bool, error) {
	return _GenerativeMarketplace.Contract.AllowableERC20MakeListToken(&_GenerativeMarketplace.CallOpts, arg0)
}

// AllowableERC20MakeOffer is a free data retrieval call binding the contract method 0xbd7ec0c2.
//
// Solidity: function _allowableERC20MakeOffer(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) AllowableERC20MakeOffer(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_allowableERC20MakeOffer", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowableERC20MakeOffer is a free data retrieval call binding the contract method 0xbd7ec0c2.
//
// Solidity: function _allowableERC20MakeOffer(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) AllowableERC20MakeOffer(arg0 common.Address) (bool, error) {
	return _GenerativeMarketplace.Contract.AllowableERC20MakeOffer(&_GenerativeMarketplace.CallOpts, arg0)
}

// AllowableERC20MakeOffer is a free data retrieval call binding the contract method 0xbd7ec0c2.
//
// Solidity: function _allowableERC20MakeOffer(address ) view returns(bool)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) AllowableERC20MakeOffer(arg0 common.Address) (bool, error) {
	return _GenerativeMarketplace.Contract.AllowableERC20MakeOffer(&_GenerativeMarketplace.CallOpts, arg0)
}

// ArrayListingId is a free data retrieval call binding the contract method 0xed4ec201.
//
// Solidity: function _arrayListingId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ArrayListingId(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_arrayListingId", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ArrayListingId is a free data retrieval call binding the contract method 0xed4ec201.
//
// Solidity: function _arrayListingId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ArrayListingId(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeMarketplace.Contract.ArrayListingId(&_GenerativeMarketplace.CallOpts, arg0)
}

// ArrayListingId is a free data retrieval call binding the contract method 0xed4ec201.
//
// Solidity: function _arrayListingId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ArrayListingId(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeMarketplace.Contract.ArrayListingId(&_GenerativeMarketplace.CallOpts, arg0)
}

// ArrayMakeOfferId is a free data retrieval call binding the contract method 0x1fd4d2ed.
//
// Solidity: function _arrayMakeOfferId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ArrayMakeOfferId(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_arrayMakeOfferId", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ArrayMakeOfferId is a free data retrieval call binding the contract method 0x1fd4d2ed.
//
// Solidity: function _arrayMakeOfferId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ArrayMakeOfferId(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeMarketplace.Contract.ArrayMakeOfferId(&_GenerativeMarketplace.CallOpts, arg0)
}

// ArrayMakeOfferId is a free data retrieval call binding the contract method 0x1fd4d2ed.
//
// Solidity: function _arrayMakeOfferId(uint256 ) view returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ArrayMakeOfferId(arg0 *big.Int) ([32]byte, error) {
	return _GenerativeMarketplace.Contract.ArrayMakeOfferId(&_GenerativeMarketplace.CallOpts, arg0)
}

// ListingTokenDataMapping is a free data retrieval call binding the contract method 0xd299a086.
//
// Solidity: function _listingTokenDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ListingTokenDataMapping(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_listingTokenDataMapping", arg0, arg1, arg2)

	outstruct := new(struct {
		CollectionContract common.Address
		TokenId            *big.Int
		Seller             common.Address
		Erc20Token         common.Address
		Price              *big.Int
		Closed             bool
		DurationTime       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CollectionContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Erc20Token = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Price = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Closed = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.DurationTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ListingTokenDataMapping is a free data retrieval call binding the contract method 0xd299a086.
//
// Solidity: function _listingTokenDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ListingTokenDataMapping(arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.ListingTokenDataMapping(&_GenerativeMarketplace.CallOpts, arg0, arg1, arg2)
}

// ListingTokenDataMapping is a free data retrieval call binding the contract method 0xd299a086.
//
// Solidity: function _listingTokenDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ListingTokenDataMapping(arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.ListingTokenDataMapping(&_GenerativeMarketplace.CallOpts, arg0, arg1, arg2)
}

// ListingTokenIds is a free data retrieval call binding the contract method 0xeab651d8.
//
// Solidity: function _listingTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ListingTokenIds(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_listingTokenIds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ListingTokenIds is a free data retrieval call binding the contract method 0xeab651d8.
//
// Solidity: function _listingTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ListingTokenIds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _GenerativeMarketplace.Contract.ListingTokenIds(&_GenerativeMarketplace.CallOpts, arg0, arg1)
}

// ListingTokenIds is a free data retrieval call binding the contract method 0xeab651d8.
//
// Solidity: function _listingTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ListingTokenIds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _GenerativeMarketplace.Contract.ListingTokenIds(&_GenerativeMarketplace.CallOpts, arg0, arg1)
}

// ListingTokens is a free data retrieval call binding the contract method 0xa3a1937c.
//
// Solidity: function _listingTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ListingTokens(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_listingTokens", arg0)

	outstruct := new(struct {
		CollectionContract common.Address
		TokenId            *big.Int
		Seller             common.Address
		Erc20Token         common.Address
		Price              *big.Int
		Closed             bool
		DurationTime       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CollectionContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Seller = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Erc20Token = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Price = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Closed = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.DurationTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ListingTokens is a free data retrieval call binding the contract method 0xa3a1937c.
//
// Solidity: function _listingTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ListingTokens(arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.ListingTokens(&_GenerativeMarketplace.CallOpts, arg0)
}

// ListingTokens is a free data retrieval call binding the contract method 0xa3a1937c.
//
// Solidity: function _listingTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _seller, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ListingTokens(arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Seller             common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.ListingTokens(&_GenerativeMarketplace.CallOpts, arg0)
}

// MakeOfferDataMapping is a free data retrieval call binding the contract method 0x14646f6d.
//
// Solidity: function _makeOfferDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) MakeOfferDataMapping(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_makeOfferDataMapping", arg0, arg1, arg2)

	outstruct := new(struct {
		CollectionContract common.Address
		TokenId            *big.Int
		Buyer              common.Address
		Erc20Token         common.Address
		Price              *big.Int
		Closed             bool
		DurationTime       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CollectionContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Buyer = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Erc20Token = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Price = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Closed = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.DurationTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// MakeOfferDataMapping is a free data retrieval call binding the contract method 0x14646f6d.
//
// Solidity: function _makeOfferDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) MakeOfferDataMapping(arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.MakeOfferDataMapping(&_GenerativeMarketplace.CallOpts, arg0, arg1, arg2)
}

// MakeOfferDataMapping is a free data retrieval call binding the contract method 0x14646f6d.
//
// Solidity: function _makeOfferDataMapping(address , uint256 , uint256 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) MakeOfferDataMapping(arg0 common.Address, arg1 *big.Int, arg2 *big.Int) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.MakeOfferDataMapping(&_GenerativeMarketplace.CallOpts, arg0, arg1, arg2)
}

// MakeOfferTokenIds is a free data retrieval call binding the contract method 0xa25d6dce.
//
// Solidity: function _makeOfferTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) MakeOfferTokenIds(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_makeOfferTokenIds", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MakeOfferTokenIds is a free data retrieval call binding the contract method 0xa25d6dce.
//
// Solidity: function _makeOfferTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) MakeOfferTokenIds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _GenerativeMarketplace.Contract.MakeOfferTokenIds(&_GenerativeMarketplace.CallOpts, arg0, arg1)
}

// MakeOfferTokenIds is a free data retrieval call binding the contract method 0xa25d6dce.
//
// Solidity: function _makeOfferTokenIds(address , uint256 ) view returns(uint256)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) MakeOfferTokenIds(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _GenerativeMarketplace.Contract.MakeOfferTokenIds(&_GenerativeMarketplace.CallOpts, arg0, arg1)
}

// MakeOfferTokens is a free data retrieval call binding the contract method 0x3dd72c40.
//
// Solidity: function _makeOfferTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) MakeOfferTokens(opts *bind.CallOpts, arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_makeOfferTokens", arg0)

	outstruct := new(struct {
		CollectionContract common.Address
		TokenId            *big.Int
		Buyer              common.Address
		Erc20Token         common.Address
		Price              *big.Int
		Closed             bool
		DurationTime       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CollectionContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Buyer = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Erc20Token = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Price = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Closed = *abi.ConvertType(out[5], new(bool)).(*bool)
	outstruct.DurationTime = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// MakeOfferTokens is a free data retrieval call binding the contract method 0x3dd72c40.
//
// Solidity: function _makeOfferTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) MakeOfferTokens(arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.MakeOfferTokens(&_GenerativeMarketplace.CallOpts, arg0)
}

// MakeOfferTokens is a free data retrieval call binding the contract method 0x3dd72c40.
//
// Solidity: function _makeOfferTokens(bytes32 ) view returns(address _collectionContract, uint256 _tokenId, address _buyer, address _erc20Token, uint256 _price, bool _closed, uint256 _durationTime)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) MakeOfferTokens(arg0 [32]byte) (struct {
	CollectionContract common.Address
	TokenId            *big.Int
	Buyer              common.Address
	Erc20Token         common.Address
	Price              *big.Int
	Closed             bool
	DurationTime       *big.Int
}, error) {
	return _GenerativeMarketplace.Contract.MakeOfferTokens(&_GenerativeMarketplace.CallOpts, arg0)
}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceCaller) ParameterAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeMarketplace.contract.Call(opts, &out, "_parameterAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ParameterAddr() (common.Address, error) {
	return _GenerativeMarketplace.Contract.ParameterAddr(&_GenerativeMarketplace.CallOpts)
}

// ParameterAddr is a free data retrieval call binding the contract method 0x72c035cf.
//
// Solidity: function _parameterAddr() view returns(address)
func (_GenerativeMarketplace *GenerativeMarketplaceCallerSession) ParameterAddr() (common.Address, error) {
	return _GenerativeMarketplace.Contract.ParameterAddr(&_GenerativeMarketplace.CallOpts)
}

// AcceptMakeOffer is a paid mutator transaction binding the contract method 0xd9e87e61.
//
// Solidity: function acceptMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) AcceptMakeOffer(opts *bind.TransactOpts, offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "acceptMakeOffer", offerId)
}

// AcceptMakeOffer is a paid mutator transaction binding the contract method 0xd9e87e61.
//
// Solidity: function acceptMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) AcceptMakeOffer(offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.AcceptMakeOffer(&_GenerativeMarketplace.TransactOpts, offerId)
}

// AcceptMakeOffer is a paid mutator transaction binding the contract method 0xd9e87e61.
//
// Solidity: function acceptMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) AcceptMakeOffer(offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.AcceptMakeOffer(&_GenerativeMarketplace.TransactOpts, offerId)
}

// CancelListing is a paid mutator transaction binding the contract method 0x9299e552.
//
// Solidity: function cancelListing(bytes32 _offeringId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) CancelListing(opts *bind.TransactOpts, _offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "cancelListing", _offeringId)
}

// CancelListing is a paid mutator transaction binding the contract method 0x9299e552.
//
// Solidity: function cancelListing(bytes32 _offeringId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) CancelListing(_offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.CancelListing(&_GenerativeMarketplace.TransactOpts, _offeringId)
}

// CancelListing is a paid mutator transaction binding the contract method 0x9299e552.
//
// Solidity: function cancelListing(bytes32 _offeringId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) CancelListing(_offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.CancelListing(&_GenerativeMarketplace.TransactOpts, _offeringId)
}

// CancelMakeOffer is a paid mutator transaction binding the contract method 0xf1f2df0e.
//
// Solidity: function cancelMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) CancelMakeOffer(opts *bind.TransactOpts, offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "cancelMakeOffer", offerId)
}

// CancelMakeOffer is a paid mutator transaction binding the contract method 0xf1f2df0e.
//
// Solidity: function cancelMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) CancelMakeOffer(offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.CancelMakeOffer(&_GenerativeMarketplace.TransactOpts, offerId)
}

// CancelMakeOffer is a paid mutator transaction binding the contract method 0xf1f2df0e.
//
// Solidity: function cancelMakeOffer(bytes32 offerId) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) CancelMakeOffer(offerId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.CancelMakeOffer(&_GenerativeMarketplace.TransactOpts, offerId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ChangeAdmin(&_GenerativeMarketplace.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ChangeAdmin(&_GenerativeMarketplace.TransactOpts, newAdm)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) ChangeParamAddr(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "changeParamAddr", newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ChangeParamAddr(&_GenerativeMarketplace.TransactOpts, newAddr)
}

// ChangeParamAddr is a paid mutator transaction binding the contract method 0xfebfd6c3.
//
// Solidity: function changeParamAddr(address newAddr) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) ChangeParamAddr(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ChangeParamAddr(&_GenerativeMarketplace.TransactOpts, newAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) Initialize(opts *bind.TransactOpts, admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "initialize", admin, parameterControl)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) Initialize(admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.Initialize(&_GenerativeMarketplace.TransactOpts, admin, parameterControl)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address admin, address parameterControl) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) Initialize(admin common.Address, parameterControl common.Address) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.Initialize(&_GenerativeMarketplace.TransactOpts, admin, parameterControl)
}

// ListToken is a paid mutator transaction binding the contract method 0xf40cb02d.
//
// Solidity: function listToken((address,uint256,address,address,uint256,bool,uint256) listingData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) ListToken(opts *bind.TransactOpts, listingData MarketplaceListingTokenData) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "listToken", listingData)
}

// ListToken is a paid mutator transaction binding the contract method 0xf40cb02d.
//
// Solidity: function listToken((address,uint256,address,address,uint256,bool,uint256) listingData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) ListToken(listingData MarketplaceListingTokenData) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ListToken(&_GenerativeMarketplace.TransactOpts, listingData)
}

// ListToken is a paid mutator transaction binding the contract method 0xf40cb02d.
//
// Solidity: function listToken((address,uint256,address,address,uint256,bool,uint256) listingData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) ListToken(listingData MarketplaceListingTokenData) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.ListToken(&_GenerativeMarketplace.TransactOpts, listingData)
}

// MakeOffer is a paid mutator transaction binding the contract method 0xd2abb69c.
//
// Solidity: function makeOffer((address,uint256,address,address,uint256,bool,uint256) makeOfferData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) MakeOffer(opts *bind.TransactOpts, makeOfferData MarketplaceMakeOfferData) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "makeOffer", makeOfferData)
}

// MakeOffer is a paid mutator transaction binding the contract method 0xd2abb69c.
//
// Solidity: function makeOffer((address,uint256,address,address,uint256,bool,uint256) makeOfferData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceSession) MakeOffer(makeOfferData MarketplaceMakeOfferData) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.MakeOffer(&_GenerativeMarketplace.TransactOpts, makeOfferData)
}

// MakeOffer is a paid mutator transaction binding the contract method 0xd2abb69c.
//
// Solidity: function makeOffer((address,uint256,address,address,uint256,bool,uint256) makeOfferData) returns(bytes32)
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) MakeOffer(makeOfferData MarketplaceMakeOfferData) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.MakeOffer(&_GenerativeMarketplace.TransactOpts, makeOfferData)
}

// PurchaseToken is a paid mutator transaction binding the contract method 0x3ac1089a.
//
// Solidity: function purchaseToken(bytes32 offeringId) payable returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) PurchaseToken(opts *bind.TransactOpts, offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "purchaseToken", offeringId)
}

// PurchaseToken is a paid mutator transaction binding the contract method 0x3ac1089a.
//
// Solidity: function purchaseToken(bytes32 offeringId) payable returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) PurchaseToken(offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.PurchaseToken(&_GenerativeMarketplace.TransactOpts, offeringId)
}

// PurchaseToken is a paid mutator transaction binding the contract method 0x3ac1089a.
//
// Solidity: function purchaseToken(bytes32 offeringId) payable returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) PurchaseToken(offeringId [32]byte) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.PurchaseToken(&_GenerativeMarketplace.TransactOpts, offeringId)
}

// SetApproveERC20ListToken is a paid mutator transaction binding the contract method 0xe5026b8b.
//
// Solidity: function setApproveERC20ListToken(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) SetApproveERC20ListToken(opts *bind.TransactOpts, erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "setApproveERC20ListToken", erc20, allow)
}

// SetApproveERC20ListToken is a paid mutator transaction binding the contract method 0xe5026b8b.
//
// Solidity: function setApproveERC20ListToken(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) SetApproveERC20ListToken(erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.SetApproveERC20ListToken(&_GenerativeMarketplace.TransactOpts, erc20, allow)
}

// SetApproveERC20ListToken is a paid mutator transaction binding the contract method 0xe5026b8b.
//
// Solidity: function setApproveERC20ListToken(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) SetApproveERC20ListToken(erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.SetApproveERC20ListToken(&_GenerativeMarketplace.TransactOpts, erc20, allow)
}

// SetApproveERC20MakeOffer is a paid mutator transaction binding the contract method 0x9142e6a8.
//
// Solidity: function setApproveERC20MakeOffer(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) SetApproveERC20MakeOffer(opts *bind.TransactOpts, erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "setApproveERC20MakeOffer", erc20, allow)
}

// SetApproveERC20MakeOffer is a paid mutator transaction binding the contract method 0x9142e6a8.
//
// Solidity: function setApproveERC20MakeOffer(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) SetApproveERC20MakeOffer(erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.SetApproveERC20MakeOffer(&_GenerativeMarketplace.TransactOpts, erc20, allow)
}

// SetApproveERC20MakeOffer is a paid mutator transaction binding the contract method 0x9142e6a8.
//
// Solidity: function setApproveERC20MakeOffer(address erc20, bool allow) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) SetApproveERC20MakeOffer(erc20 common.Address, allow bool) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.SetApproveERC20MakeOffer(&_GenerativeMarketplace.TransactOpts, erc20, allow)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactor) Withdraw(opts *bind.TransactOpts, receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeMarketplace.contract.Transact(opts, "withdraw", receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.Withdraw(&_GenerativeMarketplace.TransactOpts, receiver, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address receiver, address erc20Addr, uint256 amount) returns()
func (_GenerativeMarketplace *GenerativeMarketplaceTransactorSession) Withdraw(receiver common.Address, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeMarketplace.Contract.Withdraw(&_GenerativeMarketplace.TransactOpts, receiver, erc20Addr, amount)
}

// GenerativeMarketplaceInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GenerativeMarketplace contract.
type GenerativeMarketplaceInitializedIterator struct {
	Event *GenerativeMarketplaceInitialized // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceInitialized)
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
		it.Event = new(GenerativeMarketplaceInitialized)
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
func (it *GenerativeMarketplaceInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceInitialized represents a Initialized event raised by the GenerativeMarketplace contract.
type GenerativeMarketplaceInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeMarketplace *GenerativeMarketplaceFilterer) FilterInitialized(opts *bind.FilterOpts) (*GenerativeMarketplaceInitializedIterator, error) {

	logs, sub, err := _GenerativeMarketplace.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceInitializedIterator{contract: _GenerativeMarketplace.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeMarketplace *GenerativeMarketplaceFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceInitialized) (event.Subscription, error) {

	logs, sub, err := _GenerativeMarketplace.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceInitialized)
				if err := _GenerativeMarketplace.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GenerativeMarketplace *GenerativeMarketplaceFilterer) ParseInitialized(log types.Log) (*GenerativeMarketplaceInitialized, error) {
	event := new(GenerativeMarketplaceInitialized)
	if err := _GenerativeMarketplace.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
