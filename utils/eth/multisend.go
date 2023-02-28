package eth

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

// MultisendMetaData contains all meta data concerning the Multisend contract.
var MultisendMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"emergencyStop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"_escapeHatchDestination\",\"type\":\"address\"}],\"name\":\"escapeHatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EscapeHatchCalled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"addresspayable[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"}],\"name\":\"multiTransfer_OST\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"multiTransferEqual_L1R\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_amountSum\",\"type\":\"uint256\"}],\"name\":\"multiTransferToken_a4A\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"multiTransferTokenEqual_71p\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"addresspayable[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_amountSum\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"_amountsEther\",\"type\":\"uint256[]\"}],\"name\":\"multiTransferTokenEther\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"addresspayable[]\",\"name\":\"_addresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amountEther\",\"type\":\"uint256\"}],\"name\":\"multiTransferTokenEtherEqual\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_address1\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount1\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"_address2\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount2\",\"type\":\"uint256\"}],\"name\":\"transfer2\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// MultisendABI is the input ABI used to generate the binding from.
// Deprecated: Use MultisendMetaData.ABI instead.
var MultisendABI = MultisendMetaData.ABI

// Multisend is an auto generated Go binding around an Ethereum contract.
type Multisend struct {
	MultisendCaller     // Read-only binding to the contract
	MultisendTransactor // Write-only binding to the contract
	MultisendFilterer   // Log filterer for contract events
}

// MultisendCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultisendCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisendTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultisendTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisendFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultisendFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisendSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultisendSession struct {
	Contract     *Multisend        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultisendCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultisendCallerSession struct {
	Contract *MultisendCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// MultisendTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultisendTransactorSession struct {
	Contract     *MultisendTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// MultisendRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultisendRaw struct {
	Contract *Multisend // Generic contract binding to access the raw methods on
}

// MultisendCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultisendCallerRaw struct {
	Contract *MultisendCaller // Generic read-only contract binding to access the raw methods on
}

// MultisendTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultisendTransactorRaw struct {
	Contract *MultisendTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultisend creates a new instance of Multisend, bound to a specific deployed contract.
func NewMultisend(address common.Address, backend bind.ContractBackend) (*Multisend, error) {
	contract, err := bindMultisend(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multisend{MultisendCaller: MultisendCaller{contract: contract}, MultisendTransactor: MultisendTransactor{contract: contract}, MultisendFilterer: MultisendFilterer{contract: contract}}, nil
}

// NewMultisendCaller creates a new read-only instance of Multisend, bound to a specific deployed contract.
func NewMultisendCaller(address common.Address, caller bind.ContractCaller) (*MultisendCaller, error) {
	contract, err := bindMultisend(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultisendCaller{contract: contract}, nil
}

// NewMultisendTransactor creates a new write-only instance of Multisend, bound to a specific deployed contract.
func NewMultisendTransactor(address common.Address, transactor bind.ContractTransactor) (*MultisendTransactor, error) {
	contract, err := bindMultisend(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultisendTransactor{contract: contract}, nil
}

// NewMultisendFilterer creates a new log filterer instance of Multisend, bound to a specific deployed contract.
func NewMultisendFilterer(address common.Address, filterer bind.ContractFilterer) (*MultisendFilterer, error) {
	contract, err := bindMultisend(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultisendFilterer{contract: contract}, nil
}

// bindMultisend binds a generic wrapper to an already deployed contract.
func bindMultisend(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisendABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multisend *MultisendRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multisend.Contract.MultisendCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multisend *MultisendRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisend.Contract.MultisendTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multisend *MultisendRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multisend.Contract.MultisendTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Multisend *MultisendCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Multisend.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Multisend *MultisendTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisend.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Multisend *MultisendTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Multisend.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Multisend *MultisendCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Multisend.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Multisend *MultisendSession) Owner() (common.Address, error) {
	return _Multisend.Contract.Owner(&_Multisend.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Multisend *MultisendCallerSession) Owner() (common.Address, error) {
	return _Multisend.Contract.Owner(&_Multisend.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Multisend *MultisendCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Multisend.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Multisend *MultisendSession) Paused() (bool, error) {
	return _Multisend.Contract.Paused(&_Multisend.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Multisend *MultisendCallerSession) Paused() (bool, error) {
	return _Multisend.Contract.Paused(&_Multisend.CallOpts)
}

// EmergencyStop is a paid mutator transaction binding the contract method 0x63a599a4.
//
// Solidity: function emergencyStop() returns()
func (_Multisend *MultisendTransactor) EmergencyStop(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "emergencyStop")
}

// EmergencyStop is a paid mutator transaction binding the contract method 0x63a599a4.
//
// Solidity: function emergencyStop() returns()
func (_Multisend *MultisendSession) EmergencyStop() (*types.Transaction, error) {
	return _Multisend.Contract.EmergencyStop(&_Multisend.TransactOpts)
}

// EmergencyStop is a paid mutator transaction binding the contract method 0x63a599a4.
//
// Solidity: function emergencyStop() returns()
func (_Multisend *MultisendTransactorSession) EmergencyStop() (*types.Transaction, error) {
	return _Multisend.Contract.EmergencyStop(&_Multisend.TransactOpts)
}

// EscapeHatch is a paid mutator transaction binding the contract method 0x1c65a898.
//
// Solidity: function escapeHatch(address _token, address _escapeHatchDestination) returns()
func (_Multisend *MultisendTransactor) EscapeHatch(opts *bind.TransactOpts, _token common.Address, _escapeHatchDestination common.Address) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "escapeHatch", _token, _escapeHatchDestination)
}

// EscapeHatch is a paid mutator transaction binding the contract method 0x1c65a898.
//
// Solidity: function escapeHatch(address _token, address _escapeHatchDestination) returns()
func (_Multisend *MultisendSession) EscapeHatch(_token common.Address, _escapeHatchDestination common.Address) (*types.Transaction, error) {
	return _Multisend.Contract.EscapeHatch(&_Multisend.TransactOpts, _token, _escapeHatchDestination)
}

// EscapeHatch is a paid mutator transaction binding the contract method 0x1c65a898.
//
// Solidity: function escapeHatch(address _token, address _escapeHatchDestination) returns()
func (_Multisend *MultisendTransactorSession) EscapeHatch(_token common.Address, _escapeHatchDestination common.Address) (*types.Transaction, error) {
	return _Multisend.Contract.EscapeHatch(&_Multisend.TransactOpts, _token, _escapeHatchDestination)
}

// MultiTransferEqualL1R is a paid mutator transaction binding the contract method 0x00000983.
//
// Solidity: function multiTransferEqual_L1R(address[] _addresses, uint256 _amount) payable returns(bool)
func (_Multisend *MultisendTransactor) MultiTransferEqualL1R(opts *bind.TransactOpts, _addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransferEqual_L1R", _addresses, _amount)
}

// MultiTransferEqualL1R is a paid mutator transaction binding the contract method 0x00000983.
//
// Solidity: function multiTransferEqual_L1R(address[] _addresses, uint256 _amount) payable returns(bool)
func (_Multisend *MultisendSession) MultiTransferEqualL1R(_addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferEqualL1R(&_Multisend.TransactOpts, _addresses, _amount)
}

// MultiTransferEqualL1R is a paid mutator transaction binding the contract method 0x00000983.
//
// Solidity: function multiTransferEqual_L1R(address[] _addresses, uint256 _amount) payable returns(bool)
func (_Multisend *MultisendTransactorSession) MultiTransferEqualL1R(_addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferEqualL1R(&_Multisend.TransactOpts, _addresses, _amount)
}

// MultiTransferTokenEqual71p is a paid mutator transaction binding the contract method 0x00004f8a.
//
// Solidity: function multiTransferTokenEqual_71p(address _token, address[] _addresses, uint256 _amount) payable returns()
func (_Multisend *MultisendTransactor) MultiTransferTokenEqual71p(opts *bind.TransactOpts, _token common.Address, _addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransferTokenEqual_71p", _token, _addresses, _amount)
}

// MultiTransferTokenEqual71p is a paid mutator transaction binding the contract method 0x00004f8a.
//
// Solidity: function multiTransferTokenEqual_71p(address _token, address[] _addresses, uint256 _amount) payable returns()
func (_Multisend *MultisendSession) MultiTransferTokenEqual71p(_token common.Address, _addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEqual71p(&_Multisend.TransactOpts, _token, _addresses, _amount)
}

// MultiTransferTokenEqual71p is a paid mutator transaction binding the contract method 0x00004f8a.
//
// Solidity: function multiTransferTokenEqual_71p(address _token, address[] _addresses, uint256 _amount) payable returns()
func (_Multisend *MultisendTransactorSession) MultiTransferTokenEqual71p(_token common.Address, _addresses []common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEqual71p(&_Multisend.TransactOpts, _token, _addresses, _amount)
}

// MultiTransferTokenEther is a paid mutator transaction binding the contract method 0xcc5dcd11.
//
// Solidity: function multiTransferTokenEther(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum, uint256[] _amountsEther) payable returns()
func (_Multisend *MultisendTransactor) MultiTransferTokenEther(opts *bind.TransactOpts, _token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int, _amountsEther []*big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransferTokenEther", _token, _addresses, _amounts, _amountSum, _amountsEther)
}

// MultiTransferTokenEther is a paid mutator transaction binding the contract method 0xcc5dcd11.
//
// Solidity: function multiTransferTokenEther(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum, uint256[] _amountsEther) payable returns()
func (_Multisend *MultisendSession) MultiTransferTokenEther(_token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int, _amountsEther []*big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEther(&_Multisend.TransactOpts, _token, _addresses, _amounts, _amountSum, _amountsEther)
}

// MultiTransferTokenEther is a paid mutator transaction binding the contract method 0xcc5dcd11.
//
// Solidity: function multiTransferTokenEther(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum, uint256[] _amountsEther) payable returns()
func (_Multisend *MultisendTransactorSession) MultiTransferTokenEther(_token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int, _amountsEther []*big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEther(&_Multisend.TransactOpts, _token, _addresses, _amounts, _amountSum, _amountsEther)
}

// MultiTransferTokenEtherEqual is a paid mutator transaction binding the contract method 0x738d127a.
//
// Solidity: function multiTransferTokenEtherEqual(address _token, address[] _addresses, uint256 _amount, uint256 _amountEther) payable returns()
func (_Multisend *MultisendTransactor) MultiTransferTokenEtherEqual(opts *bind.TransactOpts, _token common.Address, _addresses []common.Address, _amount *big.Int, _amountEther *big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransferTokenEtherEqual", _token, _addresses, _amount, _amountEther)
}

// MultiTransferTokenEtherEqual is a paid mutator transaction binding the contract method 0x738d127a.
//
// Solidity: function multiTransferTokenEtherEqual(address _token, address[] _addresses, uint256 _amount, uint256 _amountEther) payable returns()
func (_Multisend *MultisendSession) MultiTransferTokenEtherEqual(_token common.Address, _addresses []common.Address, _amount *big.Int, _amountEther *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEtherEqual(&_Multisend.TransactOpts, _token, _addresses, _amount, _amountEther)
}

// MultiTransferTokenEtherEqual is a paid mutator transaction binding the contract method 0x738d127a.
//
// Solidity: function multiTransferTokenEtherEqual(address _token, address[] _addresses, uint256 _amount, uint256 _amountEther) payable returns()
func (_Multisend *MultisendTransactorSession) MultiTransferTokenEtherEqual(_token common.Address, _addresses []common.Address, _amount *big.Int, _amountEther *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenEtherEqual(&_Multisend.TransactOpts, _token, _addresses, _amount, _amountEther)
}

// MultiTransferTokenA4A is a paid mutator transaction binding the contract method 0x0000e2a7.
//
// Solidity: function multiTransferToken_a4A(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum) payable returns()
func (_Multisend *MultisendTransactor) MultiTransferTokenA4A(opts *bind.TransactOpts, _token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransferToken_a4A", _token, _addresses, _amounts, _amountSum)
}

// MultiTransferTokenA4A is a paid mutator transaction binding the contract method 0x0000e2a7.
//
// Solidity: function multiTransferToken_a4A(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum) payable returns()
func (_Multisend *MultisendSession) MultiTransferTokenA4A(_token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenA4A(&_Multisend.TransactOpts, _token, _addresses, _amounts, _amountSum)
}

// MultiTransferTokenA4A is a paid mutator transaction binding the contract method 0x0000e2a7.
//
// Solidity: function multiTransferToken_a4A(address _token, address[] _addresses, uint256[] _amounts, uint256 _amountSum) payable returns()
func (_Multisend *MultisendTransactorSession) MultiTransferTokenA4A(_token common.Address, _addresses []common.Address, _amounts []*big.Int, _amountSum *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferTokenA4A(&_Multisend.TransactOpts, _token, _addresses, _amounts, _amountSum)
}

// MultiTransferOST is a paid mutator transaction binding the contract method 0x000055be.
//
// Solidity: function multiTransfer_OST(address[] _addresses, uint256[] _amounts) payable returns(bool)
func (_Multisend *MultisendTransactor) MultiTransferOST(opts *bind.TransactOpts, _addresses []common.Address, _amounts []*big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "multiTransfer_OST", _addresses, _amounts)
}

// MultiTransferOST is a paid mutator transaction binding the contract method 0x000055be.
//
// Solidity: function multiTransfer_OST(address[] _addresses, uint256[] _amounts) payable returns(bool)
func (_Multisend *MultisendSession) MultiTransferOST(_addresses []common.Address, _amounts []*big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferOST(&_Multisend.TransactOpts, _addresses, _amounts)
}

// MultiTransferOST is a paid mutator transaction binding the contract method 0x000055be.
//
// Solidity: function multiTransfer_OST(address[] _addresses, uint256[] _amounts) payable returns(bool)
func (_Multisend *MultisendTransactorSession) MultiTransferOST(_addresses []common.Address, _amounts []*big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.MultiTransferOST(&_Multisend.TransactOpts, _addresses, _amounts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Multisend *MultisendTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Multisend *MultisendSession) RenounceOwnership() (*types.Transaction, error) {
	return _Multisend.Contract.RenounceOwnership(&_Multisend.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Multisend *MultisendTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Multisend.Contract.RenounceOwnership(&_Multisend.TransactOpts)
}

// Transfer2 is a paid mutator transaction binding the contract method 0xe4b8cb61.
//
// Solidity: function transfer2(address _address1, uint256 _amount1, address _address2, uint256 _amount2) payable returns(bool)
func (_Multisend *MultisendTransactor) Transfer2(opts *bind.TransactOpts, _address1 common.Address, _amount1 *big.Int, _address2 common.Address, _amount2 *big.Int) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "transfer2", _address1, _amount1, _address2, _amount2)
}

// Transfer2 is a paid mutator transaction binding the contract method 0xe4b8cb61.
//
// Solidity: function transfer2(address _address1, uint256 _amount1, address _address2, uint256 _amount2) payable returns(bool)
func (_Multisend *MultisendSession) Transfer2(_address1 common.Address, _amount1 *big.Int, _address2 common.Address, _amount2 *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.Transfer2(&_Multisend.TransactOpts, _address1, _amount1, _address2, _amount2)
}

// Transfer2 is a paid mutator transaction binding the contract method 0xe4b8cb61.
//
// Solidity: function transfer2(address _address1, uint256 _amount1, address _address2, uint256 _amount2) payable returns(bool)
func (_Multisend *MultisendTransactorSession) Transfer2(_address1 common.Address, _amount1 *big.Int, _address2 common.Address, _amount2 *big.Int) (*types.Transaction, error) {
	return _Multisend.Contract.Transfer2(&_Multisend.TransactOpts, _address1, _amount1, _address2, _amount2)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Multisend *MultisendTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Multisend.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Multisend *MultisendSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Multisend.Contract.TransferOwnership(&_Multisend.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Multisend *MultisendTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Multisend.Contract.TransferOwnership(&_Multisend.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Multisend *MultisendTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Multisend.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Multisend *MultisendSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Multisend.Contract.Fallback(&_Multisend.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Multisend *MultisendTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Multisend.Contract.Fallback(&_Multisend.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Multisend *MultisendTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Multisend.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Multisend *MultisendSession) Receive() (*types.Transaction, error) {
	return _Multisend.Contract.Receive(&_Multisend.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Multisend *MultisendTransactorSession) Receive() (*types.Transaction, error) {
	return _Multisend.Contract.Receive(&_Multisend.TransactOpts)
}

// MultisendEscapeHatchCalledIterator is returned from FilterEscapeHatchCalled and is used to iterate over the raw logs and unpacked data for EscapeHatchCalled events raised by the Multisend contract.
type MultisendEscapeHatchCalledIterator struct {
	Event *MultisendEscapeHatchCalled // Event containing the contract specifics and raw log

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
func (it *MultisendEscapeHatchCalledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisendEscapeHatchCalled)
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
		it.Event = new(MultisendEscapeHatchCalled)
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
func (it *MultisendEscapeHatchCalledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisendEscapeHatchCalledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisendEscapeHatchCalled represents a EscapeHatchCalled event raised by the Multisend contract.
type MultisendEscapeHatchCalled struct {
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEscapeHatchCalled is a free log retrieval operation binding the contract event 0xa50dde912fa22ea0d215a0236093ac45b4d55d6ef0c604c319f900029c5d10f2.
//
// Solidity: event EscapeHatchCalled(address token, uint256 amount)
func (_Multisend *MultisendFilterer) FilterEscapeHatchCalled(opts *bind.FilterOpts) (*MultisendEscapeHatchCalledIterator, error) {

	logs, sub, err := _Multisend.contract.FilterLogs(opts, "EscapeHatchCalled")
	if err != nil {
		return nil, err
	}
	return &MultisendEscapeHatchCalledIterator{contract: _Multisend.contract, event: "EscapeHatchCalled", logs: logs, sub: sub}, nil
}

// WatchEscapeHatchCalled is a free log subscription operation binding the contract event 0xa50dde912fa22ea0d215a0236093ac45b4d55d6ef0c604c319f900029c5d10f2.
//
// Solidity: event EscapeHatchCalled(address token, uint256 amount)
func (_Multisend *MultisendFilterer) WatchEscapeHatchCalled(opts *bind.WatchOpts, sink chan<- *MultisendEscapeHatchCalled) (event.Subscription, error) {

	logs, sub, err := _Multisend.contract.WatchLogs(opts, "EscapeHatchCalled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisendEscapeHatchCalled)
				if err := _Multisend.contract.UnpackLog(event, "EscapeHatchCalled", log); err != nil {
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

// ParseEscapeHatchCalled is a log parse operation binding the contract event 0xa50dde912fa22ea0d215a0236093ac45b4d55d6ef0c604c319f900029c5d10f2.
//
// Solidity: event EscapeHatchCalled(address token, uint256 amount)
func (_Multisend *MultisendFilterer) ParseEscapeHatchCalled(log types.Log) (*MultisendEscapeHatchCalled, error) {
	event := new(MultisendEscapeHatchCalled)
	if err := _Multisend.contract.UnpackLog(event, "EscapeHatchCalled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisendOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Multisend contract.
type MultisendOwnershipTransferredIterator struct {
	Event *MultisendOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MultisendOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisendOwnershipTransferred)
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
		it.Event = new(MultisendOwnershipTransferred)
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
func (it *MultisendOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisendOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisendOwnershipTransferred represents a OwnershipTransferred event raised by the Multisend contract.
type MultisendOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Multisend *MultisendFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultisendOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Multisend.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultisendOwnershipTransferredIterator{contract: _Multisend.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Multisend *MultisendFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultisendOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Multisend.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisendOwnershipTransferred)
				if err := _Multisend.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Multisend *MultisendFilterer) ParseOwnershipTransferred(log types.Log) (*MultisendOwnershipTransferred, error) {
	event := new(MultisendOwnershipTransferred)
	if err := _Multisend.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisendPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Multisend contract.
type MultisendPausedIterator struct {
	Event *MultisendPaused // Event containing the contract specifics and raw log

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
func (it *MultisendPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisendPaused)
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
		it.Event = new(MultisendPaused)
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
func (it *MultisendPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisendPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisendPaused represents a Paused event raised by the Multisend contract.
type MultisendPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Multisend *MultisendFilterer) FilterPaused(opts *bind.FilterOpts) (*MultisendPausedIterator, error) {

	logs, sub, err := _Multisend.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &MultisendPausedIterator{contract: _Multisend.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Multisend *MultisendFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *MultisendPaused) (event.Subscription, error) {

	logs, sub, err := _Multisend.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisendPaused)
				if err := _Multisend.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Multisend *MultisendFilterer) ParsePaused(log types.Log) (*MultisendPaused, error) {
	event := new(MultisendPaused)
	if err := _Multisend.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MultisendUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Multisend contract.
type MultisendUnpausedIterator struct {
	Event *MultisendUnpaused // Event containing the contract specifics and raw log

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
func (it *MultisendUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisendUnpaused)
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
		it.Event = new(MultisendUnpaused)
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
func (it *MultisendUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisendUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisendUnpaused represents a Unpaused event raised by the Multisend contract.
type MultisendUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Multisend *MultisendFilterer) FilterUnpaused(opts *bind.FilterOpts) (*MultisendUnpausedIterator, error) {

	logs, sub, err := _Multisend.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &MultisendUnpausedIterator{contract: _Multisend.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Multisend *MultisendFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *MultisendUnpaused) (event.Subscription, error) {

	logs, sub, err := _Multisend.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisendUnpaused)
				if err := _Multisend.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Multisend *MultisendFilterer) ParseUnpaused(log types.Log) (*MultisendUnpaused, error) {
	event := new(MultisendUnpaused)
	if err := _Multisend.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
