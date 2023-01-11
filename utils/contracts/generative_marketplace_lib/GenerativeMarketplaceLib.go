// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_marketplace_lib

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

// GenerativeMarketplaceLibMetaData contains all meta data concerning the GenerativeMarketplaceLib contract.
var GenerativeMarketplaceLibMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.MakeOfferData\",\"name\":\"data\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"}],\"name\":\"AcceptMakeOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.ListingTokenData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"CancelListing\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.MakeOfferData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"CancelMakeOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.ListingTokenData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"ListingToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"result\",\"type\":\"bytes32[]\"}],\"name\":\"MakeCollectionOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.MakeOfferData\",\"name\":\"data\",\"type\":\"tuple\"}],\"name\":\"MakeOffer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"_collectionContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_erc20Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"_closed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"_durationTime\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structMarketplace.ListingTokenData\",\"name\":\"data\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"}],\"name\":\"PurchaseToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"result\",\"type\":\"bytes32[]\"}],\"name\":\"Sweep\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"UpdateListingPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"offeringId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"UpdateMakeOfferPrice\",\"type\":\"event\"}]",
}

// GenerativeMarketplaceLibABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeMarketplaceLibMetaData.ABI instead.
var GenerativeMarketplaceLibABI = GenerativeMarketplaceLibMetaData.ABI

// GenerativeMarketplaceLib is an auto generated Go binding around an Ethereum contract.
type GenerativeMarketplaceLib struct {
	GenerativeMarketplaceLibCaller     // Read-only binding to the contract
	GenerativeMarketplaceLibTransactor // Write-only binding to the contract
	GenerativeMarketplaceLibFilterer   // Log filterer for contract events
}

// GenerativeMarketplaceLibCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeMarketplaceLibCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceLibTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeMarketplaceLibTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceLibFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeMarketplaceLibFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeMarketplaceLibSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeMarketplaceLibSession struct {
	Contract     *GenerativeMarketplaceLib // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// GenerativeMarketplaceLibCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeMarketplaceLibCallerSession struct {
	Contract *GenerativeMarketplaceLibCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// GenerativeMarketplaceLibTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeMarketplaceLibTransactorSession struct {
	Contract     *GenerativeMarketplaceLibTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// GenerativeMarketplaceLibRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeMarketplaceLibRaw struct {
	Contract *GenerativeMarketplaceLib // Generic contract binding to access the raw methods on
}

// GenerativeMarketplaceLibCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeMarketplaceLibCallerRaw struct {
	Contract *GenerativeMarketplaceLibCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeMarketplaceLibTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeMarketplaceLibTransactorRaw struct {
	Contract *GenerativeMarketplaceLibTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeMarketplaceLib creates a new instance of GenerativeMarketplaceLib, bound to a specific deployed contract.
func NewGenerativeMarketplaceLib(address common.Address, backend bind.ContractBackend) (*GenerativeMarketplaceLib, error) {
	contract, err := bindGenerativeMarketplaceLib(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLib{GenerativeMarketplaceLibCaller: GenerativeMarketplaceLibCaller{contract: contract}, GenerativeMarketplaceLibTransactor: GenerativeMarketplaceLibTransactor{contract: contract}, GenerativeMarketplaceLibFilterer: GenerativeMarketplaceLibFilterer{contract: contract}}, nil
}

// NewGenerativeMarketplaceLibCaller creates a new read-only instance of GenerativeMarketplaceLib, bound to a specific deployed contract.
func NewGenerativeMarketplaceLibCaller(address common.Address, caller bind.ContractCaller) (*GenerativeMarketplaceLibCaller, error) {
	contract, err := bindGenerativeMarketplaceLib(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibCaller{contract: contract}, nil
}

// NewGenerativeMarketplaceLibTransactor creates a new write-only instance of GenerativeMarketplaceLib, bound to a specific deployed contract.
func NewGenerativeMarketplaceLibTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeMarketplaceLibTransactor, error) {
	contract, err := bindGenerativeMarketplaceLib(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibTransactor{contract: contract}, nil
}

// NewGenerativeMarketplaceLibFilterer creates a new log filterer instance of GenerativeMarketplaceLib, bound to a specific deployed contract.
func NewGenerativeMarketplaceLibFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeMarketplaceLibFilterer, error) {
	contract, err := bindGenerativeMarketplaceLib(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibFilterer{contract: contract}, nil
}

// bindGenerativeMarketplaceLib binds a generic wrapper to an already deployed contract.
func bindGenerativeMarketplaceLib(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeMarketplaceLibMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeMarketplaceLib.Contract.GenerativeMarketplaceLibCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeMarketplaceLib.Contract.GenerativeMarketplaceLibTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeMarketplaceLib.Contract.GenerativeMarketplaceLibTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeMarketplaceLib.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeMarketplaceLib.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeMarketplaceLib.Contract.contract.Transact(opts, method, params...)
}

// GenerativeMarketplaceLibAcceptMakeOfferIterator is returned from FilterAcceptMakeOffer and is used to iterate over the raw logs and unpacked data for AcceptMakeOffer events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibAcceptMakeOfferIterator struct {
	Event *GenerativeMarketplaceLibAcceptMakeOffer // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibAcceptMakeOfferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibAcceptMakeOffer)
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
		it.Event = new(GenerativeMarketplaceLibAcceptMakeOffer)
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
func (it *GenerativeMarketplaceLibAcceptMakeOfferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibAcceptMakeOfferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibAcceptMakeOffer represents a AcceptMakeOffer event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibAcceptMakeOffer struct {
	OfferingId [32]byte
	Data       MarketplaceMakeOfferData
	Buyer      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAcceptMakeOffer is a free log retrieval operation binding the contract event 0xe0787f0a475d5edb93acefb462596edcf3ab4d4949407b732fec63d7e5afc60c.
//
// Solidity: event AcceptMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterAcceptMakeOffer(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibAcceptMakeOfferIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "AcceptMakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibAcceptMakeOfferIterator{contract: _GenerativeMarketplaceLib.contract, event: "AcceptMakeOffer", logs: logs, sub: sub}, nil
}

// WatchAcceptMakeOffer is a free log subscription operation binding the contract event 0xe0787f0a475d5edb93acefb462596edcf3ab4d4949407b732fec63d7e5afc60c.
//
// Solidity: event AcceptMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchAcceptMakeOffer(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibAcceptMakeOffer, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "AcceptMakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibAcceptMakeOffer)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "AcceptMakeOffer", log); err != nil {
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

// ParseAcceptMakeOffer is a log parse operation binding the contract event 0xe0787f0a475d5edb93acefb462596edcf3ab4d4949407b732fec63d7e5afc60c.
//
// Solidity: event AcceptMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseAcceptMakeOffer(log types.Log) (*GenerativeMarketplaceLibAcceptMakeOffer, error) {
	event := new(GenerativeMarketplaceLibAcceptMakeOffer)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "AcceptMakeOffer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibCancelListingIterator is returned from FilterCancelListing and is used to iterate over the raw logs and unpacked data for CancelListing events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibCancelListingIterator struct {
	Event *GenerativeMarketplaceLibCancelListing // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibCancelListingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibCancelListing)
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
		it.Event = new(GenerativeMarketplaceLibCancelListing)
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
func (it *GenerativeMarketplaceLibCancelListingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibCancelListingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibCancelListing represents a CancelListing event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibCancelListing struct {
	OfferingId [32]byte
	Data       MarketplaceListingTokenData
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelListing is a free log retrieval operation binding the contract event 0x0f9936e2210c5cb61170d9da568110dc1426511cd2a1ffabe518cb3f12791c32.
//
// Solidity: event CancelListing(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterCancelListing(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibCancelListingIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "CancelListing", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibCancelListingIterator{contract: _GenerativeMarketplaceLib.contract, event: "CancelListing", logs: logs, sub: sub}, nil
}

// WatchCancelListing is a free log subscription operation binding the contract event 0x0f9936e2210c5cb61170d9da568110dc1426511cd2a1ffabe518cb3f12791c32.
//
// Solidity: event CancelListing(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchCancelListing(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibCancelListing, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "CancelListing", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibCancelListing)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "CancelListing", log); err != nil {
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

// ParseCancelListing is a log parse operation binding the contract event 0x0f9936e2210c5cb61170d9da568110dc1426511cd2a1ffabe518cb3f12791c32.
//
// Solidity: event CancelListing(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseCancelListing(log types.Log) (*GenerativeMarketplaceLibCancelListing, error) {
	event := new(GenerativeMarketplaceLibCancelListing)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "CancelListing", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibCancelMakeOfferIterator is returned from FilterCancelMakeOffer and is used to iterate over the raw logs and unpacked data for CancelMakeOffer events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibCancelMakeOfferIterator struct {
	Event *GenerativeMarketplaceLibCancelMakeOffer // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibCancelMakeOfferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibCancelMakeOffer)
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
		it.Event = new(GenerativeMarketplaceLibCancelMakeOffer)
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
func (it *GenerativeMarketplaceLibCancelMakeOfferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibCancelMakeOfferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibCancelMakeOffer represents a CancelMakeOffer event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibCancelMakeOffer struct {
	OfferingId [32]byte
	Data       MarketplaceMakeOfferData
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCancelMakeOffer is a free log retrieval operation binding the contract event 0xd1e7e193455ba4241e816f1b6fa787b59dad27c1279c69179f6cb4b4387bf543.
//
// Solidity: event CancelMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterCancelMakeOffer(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibCancelMakeOfferIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "CancelMakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibCancelMakeOfferIterator{contract: _GenerativeMarketplaceLib.contract, event: "CancelMakeOffer", logs: logs, sub: sub}, nil
}

// WatchCancelMakeOffer is a free log subscription operation binding the contract event 0xd1e7e193455ba4241e816f1b6fa787b59dad27c1279c69179f6cb4b4387bf543.
//
// Solidity: event CancelMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchCancelMakeOffer(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibCancelMakeOffer, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "CancelMakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibCancelMakeOffer)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "CancelMakeOffer", log); err != nil {
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

// ParseCancelMakeOffer is a log parse operation binding the contract event 0xd1e7e193455ba4241e816f1b6fa787b59dad27c1279c69179f6cb4b4387bf543.
//
// Solidity: event CancelMakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseCancelMakeOffer(log types.Log) (*GenerativeMarketplaceLibCancelMakeOffer, error) {
	event := new(GenerativeMarketplaceLibCancelMakeOffer)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "CancelMakeOffer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibListingTokenIterator is returned from FilterListingToken and is used to iterate over the raw logs and unpacked data for ListingToken events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibListingTokenIterator struct {
	Event *GenerativeMarketplaceLibListingToken // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibListingTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibListingToken)
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
		it.Event = new(GenerativeMarketplaceLibListingToken)
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
func (it *GenerativeMarketplaceLibListingTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibListingTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibListingToken represents a ListingToken event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibListingToken struct {
	OfferingId [32]byte
	Data       MarketplaceListingTokenData
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterListingToken is a free log retrieval operation binding the contract event 0x36baa67ff8b95fa14aac46cea75517dc7f57ebebc1b90d48388c8c085f5cbfe3.
//
// Solidity: event ListingToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterListingToken(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibListingTokenIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "ListingToken", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibListingTokenIterator{contract: _GenerativeMarketplaceLib.contract, event: "ListingToken", logs: logs, sub: sub}, nil
}

// WatchListingToken is a free log subscription operation binding the contract event 0x36baa67ff8b95fa14aac46cea75517dc7f57ebebc1b90d48388c8c085f5cbfe3.
//
// Solidity: event ListingToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchListingToken(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibListingToken, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "ListingToken", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibListingToken)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "ListingToken", log); err != nil {
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

// ParseListingToken is a log parse operation binding the contract event 0x36baa67ff8b95fa14aac46cea75517dc7f57ebebc1b90d48388c8c085f5cbfe3.
//
// Solidity: event ListingToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseListingToken(log types.Log) (*GenerativeMarketplaceLibListingToken, error) {
	event := new(GenerativeMarketplaceLibListingToken)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "ListingToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibMakeCollectionOfferIterator is returned from FilterMakeCollectionOffer and is used to iterate over the raw logs and unpacked data for MakeCollectionOffer events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibMakeCollectionOfferIterator struct {
	Event *GenerativeMarketplaceLibMakeCollectionOffer // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibMakeCollectionOfferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibMakeCollectionOffer)
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
		it.Event = new(GenerativeMarketplaceLibMakeCollectionOffer)
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
func (it *GenerativeMarketplaceLibMakeCollectionOfferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibMakeCollectionOfferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibMakeCollectionOffer represents a MakeCollectionOffer event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibMakeCollectionOffer struct {
	Result [][32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMakeCollectionOffer is a free log retrieval operation binding the contract event 0xc123571abb841e956264da4696ac135e6e5ac168d2b17cfdf2b1fb4dd1e30ce1.
//
// Solidity: event MakeCollectionOffer(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterMakeCollectionOffer(opts *bind.FilterOpts) (*GenerativeMarketplaceLibMakeCollectionOfferIterator, error) {

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "MakeCollectionOffer")
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibMakeCollectionOfferIterator{contract: _GenerativeMarketplaceLib.contract, event: "MakeCollectionOffer", logs: logs, sub: sub}, nil
}

// WatchMakeCollectionOffer is a free log subscription operation binding the contract event 0xc123571abb841e956264da4696ac135e6e5ac168d2b17cfdf2b1fb4dd1e30ce1.
//
// Solidity: event MakeCollectionOffer(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchMakeCollectionOffer(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibMakeCollectionOffer) (event.Subscription, error) {

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "MakeCollectionOffer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibMakeCollectionOffer)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "MakeCollectionOffer", log); err != nil {
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

// ParseMakeCollectionOffer is a log parse operation binding the contract event 0xc123571abb841e956264da4696ac135e6e5ac168d2b17cfdf2b1fb4dd1e30ce1.
//
// Solidity: event MakeCollectionOffer(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseMakeCollectionOffer(log types.Log) (*GenerativeMarketplaceLibMakeCollectionOffer, error) {
	event := new(GenerativeMarketplaceLibMakeCollectionOffer)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "MakeCollectionOffer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibMakeOfferIterator is returned from FilterMakeOffer and is used to iterate over the raw logs and unpacked data for MakeOffer events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibMakeOfferIterator struct {
	Event *GenerativeMarketplaceLibMakeOffer // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibMakeOfferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibMakeOffer)
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
		it.Event = new(GenerativeMarketplaceLibMakeOffer)
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
func (it *GenerativeMarketplaceLibMakeOfferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibMakeOfferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibMakeOffer represents a MakeOffer event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibMakeOffer struct {
	OfferingId [32]byte
	Data       MarketplaceMakeOfferData
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterMakeOffer is a free log retrieval operation binding the contract event 0x1101875263abf1b7c3fc9cf51275180975f69f768a3c09e1b453fa0cad6f29b1.
//
// Solidity: event MakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterMakeOffer(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibMakeOfferIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "MakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibMakeOfferIterator{contract: _GenerativeMarketplaceLib.contract, event: "MakeOffer", logs: logs, sub: sub}, nil
}

// WatchMakeOffer is a free log subscription operation binding the contract event 0x1101875263abf1b7c3fc9cf51275180975f69f768a3c09e1b453fa0cad6f29b1.
//
// Solidity: event MakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchMakeOffer(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibMakeOffer, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "MakeOffer", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibMakeOffer)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "MakeOffer", log); err != nil {
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

// ParseMakeOffer is a log parse operation binding the contract event 0x1101875263abf1b7c3fc9cf51275180975f69f768a3c09e1b453fa0cad6f29b1.
//
// Solidity: event MakeOffer(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseMakeOffer(log types.Log) (*GenerativeMarketplaceLibMakeOffer, error) {
	event := new(GenerativeMarketplaceLibMakeOffer)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "MakeOffer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibPurchaseTokenIterator is returned from FilterPurchaseToken and is used to iterate over the raw logs and unpacked data for PurchaseToken events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibPurchaseTokenIterator struct {
	Event *GenerativeMarketplaceLibPurchaseToken // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibPurchaseTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibPurchaseToken)
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
		it.Event = new(GenerativeMarketplaceLibPurchaseToken)
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
func (it *GenerativeMarketplaceLibPurchaseTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibPurchaseTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibPurchaseToken represents a PurchaseToken event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibPurchaseToken struct {
	OfferingId [32]byte
	Data       MarketplaceListingTokenData
	Buyer      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPurchaseToken is a free log retrieval operation binding the contract event 0x744eb8cec82caa596aa6c23e80f64b53c47d2139a53db9b3bbd7eb0a0a4b7992.
//
// Solidity: event PurchaseToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterPurchaseToken(opts *bind.FilterOpts, offeringId [][32]byte) (*GenerativeMarketplaceLibPurchaseTokenIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "PurchaseToken", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibPurchaseTokenIterator{contract: _GenerativeMarketplaceLib.contract, event: "PurchaseToken", logs: logs, sub: sub}, nil
}

// WatchPurchaseToken is a free log subscription operation binding the contract event 0x744eb8cec82caa596aa6c23e80f64b53c47d2139a53db9b3bbd7eb0a0a4b7992.
//
// Solidity: event PurchaseToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchPurchaseToken(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibPurchaseToken, offeringId [][32]byte) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "PurchaseToken", offeringIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibPurchaseToken)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "PurchaseToken", log); err != nil {
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

// ParsePurchaseToken is a log parse operation binding the contract event 0x744eb8cec82caa596aa6c23e80f64b53c47d2139a53db9b3bbd7eb0a0a4b7992.
//
// Solidity: event PurchaseToken(bytes32 indexed offeringId, (address,uint256,address,address,uint256,bool,uint256) data, address buyer)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParsePurchaseToken(log types.Log) (*GenerativeMarketplaceLibPurchaseToken, error) {
	event := new(GenerativeMarketplaceLibPurchaseToken)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "PurchaseToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibSweepIterator is returned from FilterSweep and is used to iterate over the raw logs and unpacked data for Sweep events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibSweepIterator struct {
	Event *GenerativeMarketplaceLibSweep // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibSweepIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibSweep)
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
		it.Event = new(GenerativeMarketplaceLibSweep)
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
func (it *GenerativeMarketplaceLibSweepIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibSweepIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibSweep represents a Sweep event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibSweep struct {
	Result [][32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSweep is a free log retrieval operation binding the contract event 0x158dd3c50fd0d637d65e59a4fe2a179d6ec485730ab6a010b2329e8f0484c8b8.
//
// Solidity: event Sweep(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterSweep(opts *bind.FilterOpts) (*GenerativeMarketplaceLibSweepIterator, error) {

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "Sweep")
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibSweepIterator{contract: _GenerativeMarketplaceLib.contract, event: "Sweep", logs: logs, sub: sub}, nil
}

// WatchSweep is a free log subscription operation binding the contract event 0x158dd3c50fd0d637d65e59a4fe2a179d6ec485730ab6a010b2329e8f0484c8b8.
//
// Solidity: event Sweep(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchSweep(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibSweep) (event.Subscription, error) {

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "Sweep")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibSweep)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "Sweep", log); err != nil {
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

// ParseSweep is a log parse operation binding the contract event 0x158dd3c50fd0d637d65e59a4fe2a179d6ec485730ab6a010b2329e8f0484c8b8.
//
// Solidity: event Sweep(bytes32[] result)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseSweep(log types.Log) (*GenerativeMarketplaceLibSweep, error) {
	event := new(GenerativeMarketplaceLibSweep)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "Sweep", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibUpdateListingPriceIterator is returned from FilterUpdateListingPrice and is used to iterate over the raw logs and unpacked data for UpdateListingPrice events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibUpdateListingPriceIterator struct {
	Event *GenerativeMarketplaceLibUpdateListingPrice // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibUpdateListingPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibUpdateListingPrice)
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
		it.Event = new(GenerativeMarketplaceLibUpdateListingPrice)
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
func (it *GenerativeMarketplaceLibUpdateListingPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibUpdateListingPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibUpdateListingPrice represents a UpdateListingPrice event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibUpdateListingPrice struct {
	OfferingId [32]byte
	Price      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdateListingPrice is a free log retrieval operation binding the contract event 0xd427e420f535ff7b1631b634183f475083a8a062ca123f6633ae118eafd329a4.
//
// Solidity: event UpdateListingPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterUpdateListingPrice(opts *bind.FilterOpts, offeringId [][32]byte, price []*big.Int) (*GenerativeMarketplaceLibUpdateListingPriceIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "UpdateListingPrice", offeringIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibUpdateListingPriceIterator{contract: _GenerativeMarketplaceLib.contract, event: "UpdateListingPrice", logs: logs, sub: sub}, nil
}

// WatchUpdateListingPrice is a free log subscription operation binding the contract event 0xd427e420f535ff7b1631b634183f475083a8a062ca123f6633ae118eafd329a4.
//
// Solidity: event UpdateListingPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchUpdateListingPrice(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibUpdateListingPrice, offeringId [][32]byte, price []*big.Int) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "UpdateListingPrice", offeringIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibUpdateListingPrice)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "UpdateListingPrice", log); err != nil {
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

// ParseUpdateListingPrice is a log parse operation binding the contract event 0xd427e420f535ff7b1631b634183f475083a8a062ca123f6633ae118eafd329a4.
//
// Solidity: event UpdateListingPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseUpdateListingPrice(log types.Log) (*GenerativeMarketplaceLibUpdateListingPrice, error) {
	event := new(GenerativeMarketplaceLibUpdateListingPrice)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "UpdateListingPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeMarketplaceLibUpdateMakeOfferPriceIterator is returned from FilterUpdateMakeOfferPrice and is used to iterate over the raw logs and unpacked data for UpdateMakeOfferPrice events raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibUpdateMakeOfferPriceIterator struct {
	Event *GenerativeMarketplaceLibUpdateMakeOfferPrice // Event containing the contract specifics and raw log

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
func (it *GenerativeMarketplaceLibUpdateMakeOfferPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeMarketplaceLibUpdateMakeOfferPrice)
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
		it.Event = new(GenerativeMarketplaceLibUpdateMakeOfferPrice)
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
func (it *GenerativeMarketplaceLibUpdateMakeOfferPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeMarketplaceLibUpdateMakeOfferPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeMarketplaceLibUpdateMakeOfferPrice represents a UpdateMakeOfferPrice event raised by the GenerativeMarketplaceLib contract.
type GenerativeMarketplaceLibUpdateMakeOfferPrice struct {
	OfferingId [32]byte
	Price      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdateMakeOfferPrice is a free log retrieval operation binding the contract event 0x3c1b9515c85fbdd17ae5c958c9173645c3acd8f86d5d73f6ab9ea580c2a5b311.
//
// Solidity: event UpdateMakeOfferPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) FilterUpdateMakeOfferPrice(opts *bind.FilterOpts, offeringId [][32]byte, price []*big.Int) (*GenerativeMarketplaceLibUpdateMakeOfferPriceIterator, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.FilterLogs(opts, "UpdateMakeOfferPrice", offeringIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeMarketplaceLibUpdateMakeOfferPriceIterator{contract: _GenerativeMarketplaceLib.contract, event: "UpdateMakeOfferPrice", logs: logs, sub: sub}, nil
}

// WatchUpdateMakeOfferPrice is a free log subscription operation binding the contract event 0x3c1b9515c85fbdd17ae5c958c9173645c3acd8f86d5d73f6ab9ea580c2a5b311.
//
// Solidity: event UpdateMakeOfferPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) WatchUpdateMakeOfferPrice(opts *bind.WatchOpts, sink chan<- *GenerativeMarketplaceLibUpdateMakeOfferPrice, offeringId [][32]byte, price []*big.Int) (event.Subscription, error) {

	var offeringIdRule []interface{}
	for _, offeringIdItem := range offeringId {
		offeringIdRule = append(offeringIdRule, offeringIdItem)
	}
	var priceRule []interface{}
	for _, priceItem := range price {
		priceRule = append(priceRule, priceItem)
	}

	logs, sub, err := _GenerativeMarketplaceLib.contract.WatchLogs(opts, "UpdateMakeOfferPrice", offeringIdRule, priceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeMarketplaceLibUpdateMakeOfferPrice)
				if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "UpdateMakeOfferPrice", log); err != nil {
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

// ParseUpdateMakeOfferPrice is a log parse operation binding the contract event 0x3c1b9515c85fbdd17ae5c958c9173645c3acd8f86d5d73f6ab9ea580c2a5b311.
//
// Solidity: event UpdateMakeOfferPrice(bytes32 indexed offeringId, uint256 indexed price)
func (_GenerativeMarketplaceLib *GenerativeMarketplaceLibFilterer) ParseUpdateMakeOfferPrice(log types.Log) (*GenerativeMarketplaceLibUpdateMakeOfferPrice, error) {
	event := new(GenerativeMarketplaceLibUpdateMakeOfferPrice)
	if err := _GenerativeMarketplaceLib.contract.UnpackLog(event, "UpdateMakeOfferPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
