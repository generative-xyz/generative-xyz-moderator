// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package auctionv2

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

// AuctionCollection2Bidder is an auto generated low-level Go binding around an user-defined struct.
type AuctionCollection2Bidder struct {
	Amount    *big.Int
	UnitPrice uint64
	Index     uint32
	Quantity  uint32
}

// AuctionCollection2BidderResponse is an auto generated low-level Go binding around an user-defined struct.
type AuctionCollection2BidderResponse struct {
	Bidder     common.Address
	IsWinner   bool
	BidderInfo AuctionCollection2Bidder
}

// AuctionMetaData contains all meta data concerning the Auction contract.
var AuctionMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_WINNERS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"unitPrice\",\"type\":\"uint64\"}],\"name\":\"bid\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bidMinimum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32[]\",\"name\":\"winnerList\",\"type\":\"uint32[]\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"name\":\"declareWinners\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"}],\"name\":\"getBidsByAddress\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"internalType\":\"uint64\",\"name\":\"unitPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"quantity\",\"type\":\"uint32\"}],\"internalType\":\"structAuctionCollection2.Bidder\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractAuctionCollection\",\"name\":\"singleAuction_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"}],\"name\":\"listBids\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"bidder\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isWinner\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"internalType\":\"uint64\",\"name\":\"unitPrice\",\"type\":\"uint64\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"quantity\",\"type\":\"uint32\"}],\"internalType\":\"structAuctionCollection2.Bidder\",\"name\":\"bidderInfo\",\"type\":\"tuple\"}],\"internalType\":\"structAuctionCollection2.BidderResponse[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"singleAuction\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalBids\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"winnerDeclared\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"withdrawPayment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AuctionABI is the input ABI used to generate the binding from.
// Deprecated: Use AuctionMetaData.ABI instead.
var AuctionABI = AuctionMetaData.ABI

// Auction is an auto generated Go binding around an Ethereum contract.
type Auction struct {
	AuctionCaller     // Read-only binding to the contract
	AuctionTransactor // Write-only binding to the contract
	AuctionFilterer   // Log filterer for contract events
}

// AuctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type AuctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AuctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AuctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AuctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AuctionSession struct {
	Contract     *Auction          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AuctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AuctionCallerSession struct {
	Contract *AuctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AuctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AuctionTransactorSession struct {
	Contract     *AuctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AuctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type AuctionRaw struct {
	Contract *Auction // Generic contract binding to access the raw methods on
}

// AuctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AuctionCallerRaw struct {
	Contract *AuctionCaller // Generic read-only contract binding to access the raw methods on
}

// AuctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AuctionTransactorRaw struct {
	Contract *AuctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAuction creates a new instance of Auction, bound to a specific deployed contract.
func NewAuction(address common.Address, backend bind.ContractBackend) (*Auction, error) {
	contract, err := bindAuction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Auction{AuctionCaller: AuctionCaller{contract: contract}, AuctionTransactor: AuctionTransactor{contract: contract}, AuctionFilterer: AuctionFilterer{contract: contract}}, nil
}

// NewAuctionCaller creates a new read-only instance of Auction, bound to a specific deployed contract.
func NewAuctionCaller(address common.Address, caller bind.ContractCaller) (*AuctionCaller, error) {
	contract, err := bindAuction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionCaller{contract: contract}, nil
}

// NewAuctionTransactor creates a new write-only instance of Auction, bound to a specific deployed contract.
func NewAuctionTransactor(address common.Address, transactor bind.ContractTransactor) (*AuctionTransactor, error) {
	contract, err := bindAuction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AuctionTransactor{contract: contract}, nil
}

// NewAuctionFilterer creates a new log filterer instance of Auction, bound to a specific deployed contract.
func NewAuctionFilterer(address common.Address, filterer bind.ContractFilterer) (*AuctionFilterer, error) {
	contract, err := bindAuction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AuctionFilterer{contract: contract}, nil
}

// bindAuction binds a generic wrapper to an already deployed contract.
func bindAuction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AuctionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auction *AuctionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auction.Contract.AuctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auction *AuctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.Contract.AuctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auction *AuctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auction.Contract.AuctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Auction *AuctionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Auction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Auction *AuctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Auction *AuctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Auction.Contract.contract.Transact(opts, method, params...)
}

// MAXWINNERS is a free data retrieval call binding the contract method 0x29a62a76.
//
// Solidity: function MAX_WINNERS() view returns(uint256)
func (_Auction *AuctionCaller) MAXWINNERS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "MAX_WINNERS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXWINNERS is a free data retrieval call binding the contract method 0x29a62a76.
//
// Solidity: function MAX_WINNERS() view returns(uint256)
func (_Auction *AuctionSession) MAXWINNERS() (*big.Int, error) {
	return _Auction.Contract.MAXWINNERS(&_Auction.CallOpts)
}

// MAXWINNERS is a free data retrieval call binding the contract method 0x29a62a76.
//
// Solidity: function MAX_WINNERS() view returns(uint256)
func (_Auction *AuctionCallerSession) MAXWINNERS() (*big.Int, error) {
	return _Auction.Contract.MAXWINNERS(&_Auction.CallOpts)
}

// BidMinimum is a free data retrieval call binding the contract method 0x0ad5c52a.
//
// Solidity: function bidMinimum() view returns(uint256)
func (_Auction *AuctionCaller) BidMinimum(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "bidMinimum")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BidMinimum is a free data retrieval call binding the contract method 0x0ad5c52a.
//
// Solidity: function bidMinimum() view returns(uint256)
func (_Auction *AuctionSession) BidMinimum() (*big.Int, error) {
	return _Auction.Contract.BidMinimum(&_Auction.CallOpts)
}

// BidMinimum is a free data retrieval call binding the contract method 0x0ad5c52a.
//
// Solidity: function bidMinimum() view returns(uint256)
func (_Auction *AuctionCallerSession) BidMinimum() (*big.Int, error) {
	return _Auction.Contract.BidMinimum(&_Auction.CallOpts)
}

// EndTime is a free data retrieval call binding the contract method 0x3197cbb6.
//
// Solidity: function endTime() view returns(uint256)
func (_Auction *AuctionCaller) EndTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "endTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EndTime is a free data retrieval call binding the contract method 0x3197cbb6.
//
// Solidity: function endTime() view returns(uint256)
func (_Auction *AuctionSession) EndTime() (*big.Int, error) {
	return _Auction.Contract.EndTime(&_Auction.CallOpts)
}

// EndTime is a free data retrieval call binding the contract method 0x3197cbb6.
//
// Solidity: function endTime() view returns(uint256)
func (_Auction *AuctionCallerSession) EndTime() (*big.Int, error) {
	return _Auction.Contract.EndTime(&_Auction.CallOpts)
}

// GetBidsByAddress is a free data retrieval call binding the contract method 0xee9b66ec.
//
// Solidity: function getBidsByAddress(address bidder) view returns(bool, (uint128,uint64,uint32,uint32))
func (_Auction *AuctionCaller) GetBidsByAddress(opts *bind.CallOpts, bidder common.Address) (bool, AuctionCollection2Bidder, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "getBidsByAddress", bidder)

	if err != nil {
		return *new(bool), *new(AuctionCollection2Bidder), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(AuctionCollection2Bidder)).(*AuctionCollection2Bidder)

	return out0, out1, err

}

// GetBidsByAddress is a free data retrieval call binding the contract method 0xee9b66ec.
//
// Solidity: function getBidsByAddress(address bidder) view returns(bool, (uint128,uint64,uint32,uint32))
func (_Auction *AuctionSession) GetBidsByAddress(bidder common.Address) (bool, AuctionCollection2Bidder, error) {
	return _Auction.Contract.GetBidsByAddress(&_Auction.CallOpts, bidder)
}

// GetBidsByAddress is a free data retrieval call binding the contract method 0xee9b66ec.
//
// Solidity: function getBidsByAddress(address bidder) view returns(bool, (uint128,uint64,uint32,uint32))
func (_Auction *AuctionCallerSession) GetBidsByAddress(bidder common.Address) (bool, AuctionCollection2Bidder, error) {
	return _Auction.Contract.GetBidsByAddress(&_Auction.CallOpts, bidder)
}

// ListBids is a free data retrieval call binding the contract method 0x331a3655.
//
// Solidity: function listBids(uint256 start, uint256 end) view returns((address,bool,(uint128,uint64,uint32,uint32))[])
func (_Auction *AuctionCaller) ListBids(opts *bind.CallOpts, start *big.Int, end *big.Int) ([]AuctionCollection2BidderResponse, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "listBids", start, end)

	if err != nil {
		return *new([]AuctionCollection2BidderResponse), err
	}

	out0 := *abi.ConvertType(out[0], new([]AuctionCollection2BidderResponse)).(*[]AuctionCollection2BidderResponse)

	return out0, err

}

// ListBids is a free data retrieval call binding the contract method 0x331a3655.
//
// Solidity: function listBids(uint256 start, uint256 end) view returns((address,bool,(uint128,uint64,uint32,uint32))[])
func (_Auction *AuctionSession) ListBids(start *big.Int, end *big.Int) ([]AuctionCollection2BidderResponse, error) {
	return _Auction.Contract.ListBids(&_Auction.CallOpts, start, end)
}

// ListBids is a free data retrieval call binding the contract method 0x331a3655.
//
// Solidity: function listBids(uint256 start, uint256 end) view returns((address,bool,(uint128,uint64,uint32,uint32))[])
func (_Auction *AuctionCallerSession) ListBids(start *big.Int, end *big.Int) ([]AuctionCollection2BidderResponse, error) {
	return _Auction.Contract.ListBids(&_Auction.CallOpts, start, end)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionSession) Owner() (common.Address, error) {
	return _Auction.Contract.Owner(&_Auction.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Auction *AuctionCallerSession) Owner() (common.Address, error) {
	return _Auction.Contract.Owner(&_Auction.CallOpts)
}

// SingleAuction is a free data retrieval call binding the contract method 0xd5a014f8.
//
// Solidity: function singleAuction() view returns(address)
func (_Auction *AuctionCaller) SingleAuction(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "singleAuction")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SingleAuction is a free data retrieval call binding the contract method 0xd5a014f8.
//
// Solidity: function singleAuction() view returns(address)
func (_Auction *AuctionSession) SingleAuction() (common.Address, error) {
	return _Auction.Contract.SingleAuction(&_Auction.CallOpts)
}

// SingleAuction is a free data retrieval call binding the contract method 0xd5a014f8.
//
// Solidity: function singleAuction() view returns(address)
func (_Auction *AuctionCallerSession) SingleAuction() (common.Address, error) {
	return _Auction.Contract.SingleAuction(&_Auction.CallOpts)
}

// TotalBids is a free data retrieval call binding the contract method 0x8b034136.
//
// Solidity: function totalBids() view returns(uint256)
func (_Auction *AuctionCaller) TotalBids(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "totalBids")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalBids is a free data retrieval call binding the contract method 0x8b034136.
//
// Solidity: function totalBids() view returns(uint256)
func (_Auction *AuctionSession) TotalBids() (*big.Int, error) {
	return _Auction.Contract.TotalBids(&_Auction.CallOpts)
}

// TotalBids is a free data retrieval call binding the contract method 0x8b034136.
//
// Solidity: function totalBids() view returns(uint256)
func (_Auction *AuctionCallerSession) TotalBids() (*big.Int, error) {
	return _Auction.Contract.TotalBids(&_Auction.CallOpts)
}

// WinnerDeclared is a free data retrieval call binding the contract method 0xff2d4812.
//
// Solidity: function winnerDeclared() view returns(bool)
func (_Auction *AuctionCaller) WinnerDeclared(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Auction.contract.Call(opts, &out, "winnerDeclared")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// WinnerDeclared is a free data retrieval call binding the contract method 0xff2d4812.
//
// Solidity: function winnerDeclared() view returns(bool)
func (_Auction *AuctionSession) WinnerDeclared() (bool, error) {
	return _Auction.Contract.WinnerDeclared(&_Auction.CallOpts)
}

// WinnerDeclared is a free data retrieval call binding the contract method 0xff2d4812.
//
// Solidity: function winnerDeclared() view returns(bool)
func (_Auction *AuctionCallerSession) WinnerDeclared() (bool, error) {
	return _Auction.Contract.WinnerDeclared(&_Auction.CallOpts)
}

// Bid is a paid mutator transaction binding the contract method 0xc261d456.
//
// Solidity: function bid(uint64 unitPrice) payable returns()
func (_Auction *AuctionTransactor) Bid(opts *bind.TransactOpts, unitPrice uint64) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "bid", unitPrice)
}

// Bid is a paid mutator transaction binding the contract method 0xc261d456.
//
// Solidity: function bid(uint64 unitPrice) payable returns()
func (_Auction *AuctionSession) Bid(unitPrice uint64) (*types.Transaction, error) {
	return _Auction.Contract.Bid(&_Auction.TransactOpts, unitPrice)
}

// Bid is a paid mutator transaction binding the contract method 0xc261d456.
//
// Solidity: function bid(uint64 unitPrice) payable returns()
func (_Auction *AuctionTransactorSession) Bid(unitPrice uint64) (*types.Transaction, error) {
	return _Auction.Contract.Bid(&_Auction.TransactOpts, unitPrice)
}

// DeclareWinners is a paid mutator transaction binding the contract method 0x82b55d81.
//
// Solidity: function declareWinners(uint32[] winnerList, bool isFinal) returns()
func (_Auction *AuctionTransactor) DeclareWinners(opts *bind.TransactOpts, winnerList []uint32, isFinal bool) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "declareWinners", winnerList, isFinal)
}

// DeclareWinners is a paid mutator transaction binding the contract method 0x82b55d81.
//
// Solidity: function declareWinners(uint32[] winnerList, bool isFinal) returns()
func (_Auction *AuctionSession) DeclareWinners(winnerList []uint32, isFinal bool) (*types.Transaction, error) {
	return _Auction.Contract.DeclareWinners(&_Auction.TransactOpts, winnerList, isFinal)
}

// DeclareWinners is a paid mutator transaction binding the contract method 0x82b55d81.
//
// Solidity: function declareWinners(uint32[] winnerList, bool isFinal) returns()
func (_Auction *AuctionTransactorSession) DeclareWinners(winnerList []uint32, isFinal bool) (*types.Transaction, error) {
	return _Auction.Contract.DeclareWinners(&_Auction.TransactOpts, winnerList, isFinal)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address singleAuction_) returns()
func (_Auction *AuctionTransactor) Initialize(opts *bind.TransactOpts, singleAuction_ common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "initialize", singleAuction_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address singleAuction_) returns()
func (_Auction *AuctionSession) Initialize(singleAuction_ common.Address) (*types.Transaction, error) {
	return _Auction.Contract.Initialize(&_Auction.TransactOpts, singleAuction_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address singleAuction_) returns()
func (_Auction *AuctionTransactorSession) Initialize(singleAuction_ common.Address) (*types.Transaction, error) {
	return _Auction.Contract.Initialize(&_Auction.TransactOpts, singleAuction_)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_Auction *AuctionTransactor) Refund(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "refund")
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_Auction *AuctionSession) Refund() (*types.Transaction, error) {
	return _Auction.Contract.Refund(&_Auction.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_Auction *AuctionTransactorSession) Refund() (*types.Transaction, error) {
	return _Auction.Contract.Refund(&_Auction.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionSession) RenounceOwnership() (*types.Transaction, error) {
	return _Auction.Contract.RenounceOwnership(&_Auction.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Auction *AuctionTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Auction.Contract.RenounceOwnership(&_Auction.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Auction.Contract.TransferOwnership(&_Auction.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Auction *AuctionTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Auction.Contract.TransferOwnership(&_Auction.TransactOpts, newOwner)
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x853828b6.
//
// Solidity: function withdrawAll() returns()
func (_Auction *AuctionTransactor) WithdrawAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "withdrawAll")
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x853828b6.
//
// Solidity: function withdrawAll() returns()
func (_Auction *AuctionSession) WithdrawAll() (*types.Transaction, error) {
	return _Auction.Contract.WithdrawAll(&_Auction.TransactOpts)
}

// WithdrawAll is a paid mutator transaction binding the contract method 0x853828b6.
//
// Solidity: function withdrawAll() returns()
func (_Auction *AuctionTransactorSession) WithdrawAll() (*types.Transaction, error) {
	return _Auction.Contract.WithdrawAll(&_Auction.TransactOpts)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address receiver) returns()
func (_Auction *AuctionTransactor) WithdrawPayment(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _Auction.contract.Transact(opts, "withdrawPayment", receiver)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address receiver) returns()
func (_Auction *AuctionSession) WithdrawPayment(receiver common.Address) (*types.Transaction, error) {
	return _Auction.Contract.WithdrawPayment(&_Auction.TransactOpts, receiver)
}

// WithdrawPayment is a paid mutator transaction binding the contract method 0x8ac28d5a.
//
// Solidity: function withdrawPayment(address receiver) returns()
func (_Auction *AuctionTransactorSession) WithdrawPayment(receiver common.Address) (*types.Transaction, error) {
	return _Auction.Contract.WithdrawPayment(&_Auction.TransactOpts, receiver)
}

// AuctionInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Auction contract.
type AuctionInitializedIterator struct {
	Event *AuctionInitialized // Event containing the contract specifics and raw log

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
func (it *AuctionInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionInitialized)
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
		it.Event = new(AuctionInitialized)
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
func (it *AuctionInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionInitialized represents a Initialized event raised by the Auction contract.
type AuctionInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Auction *AuctionFilterer) FilterInitialized(opts *bind.FilterOpts) (*AuctionInitializedIterator, error) {

	logs, sub, err := _Auction.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AuctionInitializedIterator{contract: _Auction.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Auction *AuctionFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AuctionInitialized) (event.Subscription, error) {

	logs, sub, err := _Auction.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionInitialized)
				if err := _Auction.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Auction *AuctionFilterer) ParseInitialized(log types.Log) (*AuctionInitialized, error) {
	event := new(AuctionInitialized)
	if err := _Auction.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AuctionOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Auction contract.
type AuctionOwnershipTransferredIterator struct {
	Event *AuctionOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AuctionOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AuctionOwnershipTransferred)
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
		it.Event = new(AuctionOwnershipTransferred)
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
func (it *AuctionOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AuctionOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AuctionOwnershipTransferred represents a OwnershipTransferred event raised by the Auction contract.
type AuctionOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Auction *AuctionFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AuctionOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Auction.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AuctionOwnershipTransferredIterator{contract: _Auction.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Auction *AuctionFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AuctionOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Auction.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AuctionOwnershipTransferred)
				if err := _Auction.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Auction *AuctionFilterer) ParseOwnershipTransferred(log types.Log) (*AuctionOwnershipTransferred, error) {
	event := new(AuctionOwnershipTransferred)
	if err := _Auction.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
