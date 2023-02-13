// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegate

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

// IDelegationRegistryContractDelegation is an auto generated low-level Go binding around an user-defined struct.
type IDelegationRegistryContractDelegation struct {
	Contract common.Address
	Delegate common.Address
}

// IDelegationRegistryDelegationInfo is an auto generated low-level Go binding around an user-defined struct.
type IDelegationRegistryDelegationInfo struct {
	Type     uint8
	Vault    common.Address
	Delegate common.Address
	Contract common.Address
	TokenId  *big.Int
}

// IDelegationRegistryTokenDelegation is an auto generated low-level Go binding around an user-defined struct.
type IDelegationRegistryTokenDelegation struct {
	Contract common.Address
	TokenId  *big.Int
	Delegate common.Address
}

// DelegateMetaData contains all meta data concerning the Delegate contract.
var DelegateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"DelegateForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"DelegateForContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"DelegateForToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"RevokeAllDelegates\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"RevokeDelegate\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"checkDelegateForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"}],\"name\":\"checkDelegateForContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"checkDelegateForToken\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"delegateForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"delegateForContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"value\",\"type\":\"bool\"}],\"name\":\"delegateForToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"getContractLevelDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"internalType\":\"structIDelegationRegistry.ContractDelegation[]\",\"name\":\"contractDelegations\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"getDelegatesForAll\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"delegates\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"}],\"name\":\"getDelegatesForContract\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"delegates\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getDelegatesForToken\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"delegates\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"getDelegationsByDelegate\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIDelegationRegistry.DelegationType\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"internalType\":\"structIDelegationRegistry.DelegationInfo[]\",\"name\":\"info\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"getTokenLevelDelegations\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"contract_\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"internalType\":\"structIDelegationRegistry.TokenDelegation[]\",\"name\":\"tokenDelegations\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"revokeAllDelegates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"delegate\",\"type\":\"address\"}],\"name\":\"revokeDelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"vault\",\"type\":\"address\"}],\"name\":\"revokeSelf\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DelegateABI is the input ABI used to generate the binding from.
// Deprecated: Use DelegateMetaData.ABI instead.
var DelegateABI = DelegateMetaData.ABI

// Delegate is an auto generated Go binding around an Ethereum contract.
type Delegate struct {
	DelegateCaller     // Read-only binding to the contract
	DelegateTransactor // Write-only binding to the contract
	DelegateFilterer   // Log filterer for contract events
}

// DelegateCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegateSession struct {
	Contract     *Delegate         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DelegateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegateCallerSession struct {
	Contract *DelegateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DelegateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegateTransactorSession struct {
	Contract     *DelegateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DelegateRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegateRaw struct {
	Contract *Delegate // Generic contract binding to access the raw methods on
}

// DelegateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegateCallerRaw struct {
	Contract *DelegateCaller // Generic read-only contract binding to access the raw methods on
}

// DelegateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegateTransactorRaw struct {
	Contract *DelegateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegate creates a new instance of Delegate, bound to a specific deployed contract.
func NewDelegate(address common.Address, backend bind.ContractBackend) (*Delegate, error) {
	contract, err := bindDelegate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegate{DelegateCaller: DelegateCaller{contract: contract}, DelegateTransactor: DelegateTransactor{contract: contract}, DelegateFilterer: DelegateFilterer{contract: contract}}, nil
}

// NewDelegateCaller creates a new read-only instance of Delegate, bound to a specific deployed contract.
func NewDelegateCaller(address common.Address, caller bind.ContractCaller) (*DelegateCaller, error) {
	contract, err := bindDelegate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateCaller{contract: contract}, nil
}

// NewDelegateTransactor creates a new write-only instance of Delegate, bound to a specific deployed contract.
func NewDelegateTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegateTransactor, error) {
	contract, err := bindDelegate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegateTransactor{contract: contract}, nil
}

// NewDelegateFilterer creates a new log filterer instance of Delegate, bound to a specific deployed contract.
func NewDelegateFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegateFilterer, error) {
	contract, err := bindDelegate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegateFilterer{contract: contract}, nil
}

// bindDelegate binds a generic wrapper to an already deployed contract.
func bindDelegate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DelegateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegate *DelegateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegate.Contract.DelegateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegate *DelegateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegate *DelegateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegate *DelegateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Delegate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegate *DelegateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegate *DelegateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegate.Contract.contract.Transact(opts, method, params...)
}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0x9c395bc2.
//
// Solidity: function checkDelegateForAll(address delegate, address vault) view returns(bool)
func (_Delegate *DelegateCaller) CheckDelegateForAll(opts *bind.CallOpts, delegate common.Address, vault common.Address) (bool, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "checkDelegateForAll", delegate, vault)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0x9c395bc2.
//
// Solidity: function checkDelegateForAll(address delegate, address vault) view returns(bool)
func (_Delegate *DelegateSession) CheckDelegateForAll(delegate common.Address, vault common.Address) (bool, error) {
	return _Delegate.Contract.CheckDelegateForAll(&_Delegate.CallOpts, delegate, vault)
}

// CheckDelegateForAll is a free data retrieval call binding the contract method 0x9c395bc2.
//
// Solidity: function checkDelegateForAll(address delegate, address vault) view returns(bool)
func (_Delegate *DelegateCallerSession) CheckDelegateForAll(delegate common.Address, vault common.Address) (bool, error) {
	return _Delegate.Contract.CheckDelegateForAll(&_Delegate.CallOpts, delegate, vault)
}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x90c9a2d0.
//
// Solidity: function checkDelegateForContract(address delegate, address vault, address contract_) view returns(bool)
func (_Delegate *DelegateCaller) CheckDelegateForContract(opts *bind.CallOpts, delegate common.Address, vault common.Address, contract_ common.Address) (bool, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "checkDelegateForContract", delegate, vault, contract_)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x90c9a2d0.
//
// Solidity: function checkDelegateForContract(address delegate, address vault, address contract_) view returns(bool)
func (_Delegate *DelegateSession) CheckDelegateForContract(delegate common.Address, vault common.Address, contract_ common.Address) (bool, error) {
	return _Delegate.Contract.CheckDelegateForContract(&_Delegate.CallOpts, delegate, vault, contract_)
}

// CheckDelegateForContract is a free data retrieval call binding the contract method 0x90c9a2d0.
//
// Solidity: function checkDelegateForContract(address delegate, address vault, address contract_) view returns(bool)
func (_Delegate *DelegateCallerSession) CheckDelegateForContract(delegate common.Address, vault common.Address, contract_ common.Address) (bool, error) {
	return _Delegate.Contract.CheckDelegateForContract(&_Delegate.CallOpts, delegate, vault, contract_)
}

// CheckDelegateForToken is a free data retrieval call binding the contract method 0xaba69cf8.
//
// Solidity: function checkDelegateForToken(address delegate, address vault, address contract_, uint256 tokenId) view returns(bool)
func (_Delegate *DelegateCaller) CheckDelegateForToken(opts *bind.CallOpts, delegate common.Address, vault common.Address, contract_ common.Address, tokenId *big.Int) (bool, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "checkDelegateForToken", delegate, vault, contract_, tokenId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CheckDelegateForToken is a free data retrieval call binding the contract method 0xaba69cf8.
//
// Solidity: function checkDelegateForToken(address delegate, address vault, address contract_, uint256 tokenId) view returns(bool)
func (_Delegate *DelegateSession) CheckDelegateForToken(delegate common.Address, vault common.Address, contract_ common.Address, tokenId *big.Int) (bool, error) {
	return _Delegate.Contract.CheckDelegateForToken(&_Delegate.CallOpts, delegate, vault, contract_, tokenId)
}

// CheckDelegateForToken is a free data retrieval call binding the contract method 0xaba69cf8.
//
// Solidity: function checkDelegateForToken(address delegate, address vault, address contract_, uint256 tokenId) view returns(bool)
func (_Delegate *DelegateCallerSession) CheckDelegateForToken(delegate common.Address, vault common.Address, contract_ common.Address, tokenId *big.Int) (bool, error) {
	return _Delegate.Contract.CheckDelegateForToken(&_Delegate.CallOpts, delegate, vault, contract_, tokenId)
}

// GetContractLevelDelegations is a free data retrieval call binding the contract method 0xf956cf94.
//
// Solidity: function getContractLevelDelegations(address vault) view returns((address,address)[] contractDelegations)
func (_Delegate *DelegateCaller) GetContractLevelDelegations(opts *bind.CallOpts, vault common.Address) ([]IDelegationRegistryContractDelegation, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getContractLevelDelegations", vault)

	if err != nil {
		return *new([]IDelegationRegistryContractDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegationRegistryContractDelegation)).(*[]IDelegationRegistryContractDelegation)

	return out0, err

}

// GetContractLevelDelegations is a free data retrieval call binding the contract method 0xf956cf94.
//
// Solidity: function getContractLevelDelegations(address vault) view returns((address,address)[] contractDelegations)
func (_Delegate *DelegateSession) GetContractLevelDelegations(vault common.Address) ([]IDelegationRegistryContractDelegation, error) {
	return _Delegate.Contract.GetContractLevelDelegations(&_Delegate.CallOpts, vault)
}

// GetContractLevelDelegations is a free data retrieval call binding the contract method 0xf956cf94.
//
// Solidity: function getContractLevelDelegations(address vault) view returns((address,address)[] contractDelegations)
func (_Delegate *DelegateCallerSession) GetContractLevelDelegations(vault common.Address) ([]IDelegationRegistryContractDelegation, error) {
	return _Delegate.Contract.GetContractLevelDelegations(&_Delegate.CallOpts, vault)
}

// GetDelegatesForAll is a free data retrieval call binding the contract method 0x1b61f675.
//
// Solidity: function getDelegatesForAll(address vault) view returns(address[] delegates)
func (_Delegate *DelegateCaller) GetDelegatesForAll(opts *bind.CallOpts, vault common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getDelegatesForAll", vault)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetDelegatesForAll is a free data retrieval call binding the contract method 0x1b61f675.
//
// Solidity: function getDelegatesForAll(address vault) view returns(address[] delegates)
func (_Delegate *DelegateSession) GetDelegatesForAll(vault common.Address) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForAll(&_Delegate.CallOpts, vault)
}

// GetDelegatesForAll is a free data retrieval call binding the contract method 0x1b61f675.
//
// Solidity: function getDelegatesForAll(address vault) view returns(address[] delegates)
func (_Delegate *DelegateCallerSession) GetDelegatesForAll(vault common.Address) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForAll(&_Delegate.CallOpts, vault)
}

// GetDelegatesForContract is a free data retrieval call binding the contract method 0xed4b878e.
//
// Solidity: function getDelegatesForContract(address vault, address contract_) view returns(address[] delegates)
func (_Delegate *DelegateCaller) GetDelegatesForContract(opts *bind.CallOpts, vault common.Address, contract_ common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getDelegatesForContract", vault, contract_)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetDelegatesForContract is a free data retrieval call binding the contract method 0xed4b878e.
//
// Solidity: function getDelegatesForContract(address vault, address contract_) view returns(address[] delegates)
func (_Delegate *DelegateSession) GetDelegatesForContract(vault common.Address, contract_ common.Address) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForContract(&_Delegate.CallOpts, vault, contract_)
}

// GetDelegatesForContract is a free data retrieval call binding the contract method 0xed4b878e.
//
// Solidity: function getDelegatesForContract(address vault, address contract_) view returns(address[] delegates)
func (_Delegate *DelegateCallerSession) GetDelegatesForContract(vault common.Address, contract_ common.Address) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForContract(&_Delegate.CallOpts, vault, contract_)
}

// GetDelegatesForToken is a free data retrieval call binding the contract method 0x1221156b.
//
// Solidity: function getDelegatesForToken(address vault, address contract_, uint256 tokenId) view returns(address[] delegates)
func (_Delegate *DelegateCaller) GetDelegatesForToken(opts *bind.CallOpts, vault common.Address, contract_ common.Address, tokenId *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getDelegatesForToken", vault, contract_, tokenId)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetDelegatesForToken is a free data retrieval call binding the contract method 0x1221156b.
//
// Solidity: function getDelegatesForToken(address vault, address contract_, uint256 tokenId) view returns(address[] delegates)
func (_Delegate *DelegateSession) GetDelegatesForToken(vault common.Address, contract_ common.Address, tokenId *big.Int) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForToken(&_Delegate.CallOpts, vault, contract_, tokenId)
}

// GetDelegatesForToken is a free data retrieval call binding the contract method 0x1221156b.
//
// Solidity: function getDelegatesForToken(address vault, address contract_, uint256 tokenId) view returns(address[] delegates)
func (_Delegate *DelegateCallerSession) GetDelegatesForToken(vault common.Address, contract_ common.Address, tokenId *big.Int) ([]common.Address, error) {
	return _Delegate.Contract.GetDelegatesForToken(&_Delegate.CallOpts, vault, contract_, tokenId)
}

// GetDelegationsByDelegate is a free data retrieval call binding the contract method 0x4fc69282.
//
// Solidity: function getDelegationsByDelegate(address delegate) view returns((uint8,address,address,address,uint256)[] info)
func (_Delegate *DelegateCaller) GetDelegationsByDelegate(opts *bind.CallOpts, delegate common.Address) ([]IDelegationRegistryDelegationInfo, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getDelegationsByDelegate", delegate)

	if err != nil {
		return *new([]IDelegationRegistryDelegationInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegationRegistryDelegationInfo)).(*[]IDelegationRegistryDelegationInfo)

	return out0, err

}

// GetDelegationsByDelegate is a free data retrieval call binding the contract method 0x4fc69282.
//
// Solidity: function getDelegationsByDelegate(address delegate) view returns((uint8,address,address,address,uint256)[] info)
func (_Delegate *DelegateSession) GetDelegationsByDelegate(delegate common.Address) ([]IDelegationRegistryDelegationInfo, error) {
	return _Delegate.Contract.GetDelegationsByDelegate(&_Delegate.CallOpts, delegate)
}

// GetDelegationsByDelegate is a free data retrieval call binding the contract method 0x4fc69282.
//
// Solidity: function getDelegationsByDelegate(address delegate) view returns((uint8,address,address,address,uint256)[] info)
func (_Delegate *DelegateCallerSession) GetDelegationsByDelegate(delegate common.Address) ([]IDelegationRegistryDelegationInfo, error) {
	return _Delegate.Contract.GetDelegationsByDelegate(&_Delegate.CallOpts, delegate)
}

// GetTokenLevelDelegations is a free data retrieval call binding the contract method 0x6f007d87.
//
// Solidity: function getTokenLevelDelegations(address vault) view returns((address,uint256,address)[] tokenDelegations)
func (_Delegate *DelegateCaller) GetTokenLevelDelegations(opts *bind.CallOpts, vault common.Address) ([]IDelegationRegistryTokenDelegation, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "getTokenLevelDelegations", vault)

	if err != nil {
		return *new([]IDelegationRegistryTokenDelegation), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDelegationRegistryTokenDelegation)).(*[]IDelegationRegistryTokenDelegation)

	return out0, err

}

// GetTokenLevelDelegations is a free data retrieval call binding the contract method 0x6f007d87.
//
// Solidity: function getTokenLevelDelegations(address vault) view returns((address,uint256,address)[] tokenDelegations)
func (_Delegate *DelegateSession) GetTokenLevelDelegations(vault common.Address) ([]IDelegationRegistryTokenDelegation, error) {
	return _Delegate.Contract.GetTokenLevelDelegations(&_Delegate.CallOpts, vault)
}

// GetTokenLevelDelegations is a free data retrieval call binding the contract method 0x6f007d87.
//
// Solidity: function getTokenLevelDelegations(address vault) view returns((address,uint256,address)[] tokenDelegations)
func (_Delegate *DelegateCallerSession) GetTokenLevelDelegations(vault common.Address) ([]IDelegationRegistryTokenDelegation, error) {
	return _Delegate.Contract.GetTokenLevelDelegations(&_Delegate.CallOpts, vault)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegate *DelegateCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Delegate.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegate *DelegateSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegate.Contract.SupportsInterface(&_Delegate.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Delegate *DelegateCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Delegate.Contract.SupportsInterface(&_Delegate.CallOpts, interfaceId)
}

// DelegateForAll is a paid mutator transaction binding the contract method 0x685ee3e8.
//
// Solidity: function delegateForAll(address delegate, bool value) returns()
func (_Delegate *DelegateTransactor) DelegateForAll(opts *bind.TransactOpts, delegate common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "delegateForAll", delegate, value)
}

// DelegateForAll is a paid mutator transaction binding the contract method 0x685ee3e8.
//
// Solidity: function delegateForAll(address delegate, bool value) returns()
func (_Delegate *DelegateSession) DelegateForAll(delegate common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForAll(&_Delegate.TransactOpts, delegate, value)
}

// DelegateForAll is a paid mutator transaction binding the contract method 0x685ee3e8.
//
// Solidity: function delegateForAll(address delegate, bool value) returns()
func (_Delegate *DelegateTransactorSession) DelegateForAll(delegate common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForAll(&_Delegate.TransactOpts, delegate, value)
}

// DelegateForContract is a paid mutator transaction binding the contract method 0x49c95d29.
//
// Solidity: function delegateForContract(address delegate, address contract_, bool value) returns()
func (_Delegate *DelegateTransactor) DelegateForContract(opts *bind.TransactOpts, delegate common.Address, contract_ common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "delegateForContract", delegate, contract_, value)
}

// DelegateForContract is a paid mutator transaction binding the contract method 0x49c95d29.
//
// Solidity: function delegateForContract(address delegate, address contract_, bool value) returns()
func (_Delegate *DelegateSession) DelegateForContract(delegate common.Address, contract_ common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForContract(&_Delegate.TransactOpts, delegate, contract_, value)
}

// DelegateForContract is a paid mutator transaction binding the contract method 0x49c95d29.
//
// Solidity: function delegateForContract(address delegate, address contract_, bool value) returns()
func (_Delegate *DelegateTransactorSession) DelegateForContract(delegate common.Address, contract_ common.Address, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForContract(&_Delegate.TransactOpts, delegate, contract_, value)
}

// DelegateForToken is a paid mutator transaction binding the contract method 0x537a5c3d.
//
// Solidity: function delegateForToken(address delegate, address contract_, uint256 tokenId, bool value) returns()
func (_Delegate *DelegateTransactor) DelegateForToken(opts *bind.TransactOpts, delegate common.Address, contract_ common.Address, tokenId *big.Int, value bool) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "delegateForToken", delegate, contract_, tokenId, value)
}

// DelegateForToken is a paid mutator transaction binding the contract method 0x537a5c3d.
//
// Solidity: function delegateForToken(address delegate, address contract_, uint256 tokenId, bool value) returns()
func (_Delegate *DelegateSession) DelegateForToken(delegate common.Address, contract_ common.Address, tokenId *big.Int, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForToken(&_Delegate.TransactOpts, delegate, contract_, tokenId, value)
}

// DelegateForToken is a paid mutator transaction binding the contract method 0x537a5c3d.
//
// Solidity: function delegateForToken(address delegate, address contract_, uint256 tokenId, bool value) returns()
func (_Delegate *DelegateTransactorSession) DelegateForToken(delegate common.Address, contract_ common.Address, tokenId *big.Int, value bool) (*types.Transaction, error) {
	return _Delegate.Contract.DelegateForToken(&_Delegate.TransactOpts, delegate, contract_, tokenId, value)
}

// RevokeAllDelegates is a paid mutator transaction binding the contract method 0x36137872.
//
// Solidity: function revokeAllDelegates() returns()
func (_Delegate *DelegateTransactor) RevokeAllDelegates(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "revokeAllDelegates")
}

// RevokeAllDelegates is a paid mutator transaction binding the contract method 0x36137872.
//
// Solidity: function revokeAllDelegates() returns()
func (_Delegate *DelegateSession) RevokeAllDelegates() (*types.Transaction, error) {
	return _Delegate.Contract.RevokeAllDelegates(&_Delegate.TransactOpts)
}

// RevokeAllDelegates is a paid mutator transaction binding the contract method 0x36137872.
//
// Solidity: function revokeAllDelegates() returns()
func (_Delegate *DelegateTransactorSession) RevokeAllDelegates() (*types.Transaction, error) {
	return _Delegate.Contract.RevokeAllDelegates(&_Delegate.TransactOpts)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_Delegate *DelegateTransactor) RevokeDelegate(opts *bind.TransactOpts, delegate common.Address) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "revokeDelegate", delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_Delegate *DelegateSession) RevokeDelegate(delegate common.Address) (*types.Transaction, error) {
	return _Delegate.Contract.RevokeDelegate(&_Delegate.TransactOpts, delegate)
}

// RevokeDelegate is a paid mutator transaction binding the contract method 0xfa352c00.
//
// Solidity: function revokeDelegate(address delegate) returns()
func (_Delegate *DelegateTransactorSession) RevokeDelegate(delegate common.Address) (*types.Transaction, error) {
	return _Delegate.Contract.RevokeDelegate(&_Delegate.TransactOpts, delegate)
}

// RevokeSelf is a paid mutator transaction binding the contract method 0x219044b0.
//
// Solidity: function revokeSelf(address vault) returns()
func (_Delegate *DelegateTransactor) RevokeSelf(opts *bind.TransactOpts, vault common.Address) (*types.Transaction, error) {
	return _Delegate.contract.Transact(opts, "revokeSelf", vault)
}

// RevokeSelf is a paid mutator transaction binding the contract method 0x219044b0.
//
// Solidity: function revokeSelf(address vault) returns()
func (_Delegate *DelegateSession) RevokeSelf(vault common.Address) (*types.Transaction, error) {
	return _Delegate.Contract.RevokeSelf(&_Delegate.TransactOpts, vault)
}

// RevokeSelf is a paid mutator transaction binding the contract method 0x219044b0.
//
// Solidity: function revokeSelf(address vault) returns()
func (_Delegate *DelegateTransactorSession) RevokeSelf(vault common.Address) (*types.Transaction, error) {
	return _Delegate.Contract.RevokeSelf(&_Delegate.TransactOpts, vault)
}

// DelegateDelegateForAllIterator is returned from FilterDelegateForAll and is used to iterate over the raw logs and unpacked data for DelegateForAll events raised by the Delegate contract.
type DelegateDelegateForAllIterator struct {
	Event *DelegateDelegateForAll // Event containing the contract specifics and raw log

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
func (it *DelegateDelegateForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateDelegateForAll)
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
		it.Event = new(DelegateDelegateForAll)
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
func (it *DelegateDelegateForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateDelegateForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateDelegateForAll represents a DelegateForAll event raised by the Delegate contract.
type DelegateDelegateForAll struct {
	Vault    common.Address
	Delegate common.Address
	Value    bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateForAll is a free log retrieval operation binding the contract event 0x58781eab4a0743ab1c285a238be846a235f06cdb5b968030573a635e5f8c92fa.
//
// Solidity: event DelegateForAll(address vault, address delegate, bool value)
func (_Delegate *DelegateFilterer) FilterDelegateForAll(opts *bind.FilterOpts) (*DelegateDelegateForAllIterator, error) {

	logs, sub, err := _Delegate.contract.FilterLogs(opts, "DelegateForAll")
	if err != nil {
		return nil, err
	}
	return &DelegateDelegateForAllIterator{contract: _Delegate.contract, event: "DelegateForAll", logs: logs, sub: sub}, nil
}

// WatchDelegateForAll is a free log subscription operation binding the contract event 0x58781eab4a0743ab1c285a238be846a235f06cdb5b968030573a635e5f8c92fa.
//
// Solidity: event DelegateForAll(address vault, address delegate, bool value)
func (_Delegate *DelegateFilterer) WatchDelegateForAll(opts *bind.WatchOpts, sink chan<- *DelegateDelegateForAll) (event.Subscription, error) {

	logs, sub, err := _Delegate.contract.WatchLogs(opts, "DelegateForAll")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateDelegateForAll)
				if err := _Delegate.contract.UnpackLog(event, "DelegateForAll", log); err != nil {
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

// ParseDelegateForAll is a log parse operation binding the contract event 0x58781eab4a0743ab1c285a238be846a235f06cdb5b968030573a635e5f8c92fa.
//
// Solidity: event DelegateForAll(address vault, address delegate, bool value)
func (_Delegate *DelegateFilterer) ParseDelegateForAll(log types.Log) (*DelegateDelegateForAll, error) {
	event := new(DelegateDelegateForAll)
	if err := _Delegate.contract.UnpackLog(event, "DelegateForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateDelegateForContractIterator is returned from FilterDelegateForContract and is used to iterate over the raw logs and unpacked data for DelegateForContract events raised by the Delegate contract.
type DelegateDelegateForContractIterator struct {
	Event *DelegateDelegateForContract // Event containing the contract specifics and raw log

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
func (it *DelegateDelegateForContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateDelegateForContract)
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
		it.Event = new(DelegateDelegateForContract)
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
func (it *DelegateDelegateForContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateDelegateForContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateDelegateForContract represents a DelegateForContract event raised by the Delegate contract.
type DelegateDelegateForContract struct {
	Vault    common.Address
	Delegate common.Address
	Contract common.Address
	Value    bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateForContract is a free log retrieval operation binding the contract event 0x8d6b2f5255b8d815cc368855b2251146e003bf4e2fcccaec66145fff5c174b4f.
//
// Solidity: event DelegateForContract(address vault, address delegate, address contract_, bool value)
func (_Delegate *DelegateFilterer) FilterDelegateForContract(opts *bind.FilterOpts) (*DelegateDelegateForContractIterator, error) {

	logs, sub, err := _Delegate.contract.FilterLogs(opts, "DelegateForContract")
	if err != nil {
		return nil, err
	}
	return &DelegateDelegateForContractIterator{contract: _Delegate.contract, event: "DelegateForContract", logs: logs, sub: sub}, nil
}

// WatchDelegateForContract is a free log subscription operation binding the contract event 0x8d6b2f5255b8d815cc368855b2251146e003bf4e2fcccaec66145fff5c174b4f.
//
// Solidity: event DelegateForContract(address vault, address delegate, address contract_, bool value)
func (_Delegate *DelegateFilterer) WatchDelegateForContract(opts *bind.WatchOpts, sink chan<- *DelegateDelegateForContract) (event.Subscription, error) {

	logs, sub, err := _Delegate.contract.WatchLogs(opts, "DelegateForContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateDelegateForContract)
				if err := _Delegate.contract.UnpackLog(event, "DelegateForContract", log); err != nil {
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

// ParseDelegateForContract is a log parse operation binding the contract event 0x8d6b2f5255b8d815cc368855b2251146e003bf4e2fcccaec66145fff5c174b4f.
//
// Solidity: event DelegateForContract(address vault, address delegate, address contract_, bool value)
func (_Delegate *DelegateFilterer) ParseDelegateForContract(log types.Log) (*DelegateDelegateForContract, error) {
	event := new(DelegateDelegateForContract)
	if err := _Delegate.contract.UnpackLog(event, "DelegateForContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateDelegateForTokenIterator is returned from FilterDelegateForToken and is used to iterate over the raw logs and unpacked data for DelegateForToken events raised by the Delegate contract.
type DelegateDelegateForTokenIterator struct {
	Event *DelegateDelegateForToken // Event containing the contract specifics and raw log

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
func (it *DelegateDelegateForTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateDelegateForToken)
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
		it.Event = new(DelegateDelegateForToken)
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
func (it *DelegateDelegateForTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateDelegateForTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateDelegateForToken represents a DelegateForToken event raised by the Delegate contract.
type DelegateDelegateForToken struct {
	Vault    common.Address
	Delegate common.Address
	Contract common.Address
	TokenId  *big.Int
	Value    bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDelegateForToken is a free log retrieval operation binding the contract event 0xe89c6ba1e8957285aed22618f52aa1dcb9d5bb64e1533d8b55136c72fcf5aa5d.
//
// Solidity: event DelegateForToken(address vault, address delegate, address contract_, uint256 tokenId, bool value)
func (_Delegate *DelegateFilterer) FilterDelegateForToken(opts *bind.FilterOpts) (*DelegateDelegateForTokenIterator, error) {

	logs, sub, err := _Delegate.contract.FilterLogs(opts, "DelegateForToken")
	if err != nil {
		return nil, err
	}
	return &DelegateDelegateForTokenIterator{contract: _Delegate.contract, event: "DelegateForToken", logs: logs, sub: sub}, nil
}

// WatchDelegateForToken is a free log subscription operation binding the contract event 0xe89c6ba1e8957285aed22618f52aa1dcb9d5bb64e1533d8b55136c72fcf5aa5d.
//
// Solidity: event DelegateForToken(address vault, address delegate, address contract_, uint256 tokenId, bool value)
func (_Delegate *DelegateFilterer) WatchDelegateForToken(opts *bind.WatchOpts, sink chan<- *DelegateDelegateForToken) (event.Subscription, error) {

	logs, sub, err := _Delegate.contract.WatchLogs(opts, "DelegateForToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateDelegateForToken)
				if err := _Delegate.contract.UnpackLog(event, "DelegateForToken", log); err != nil {
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

// ParseDelegateForToken is a log parse operation binding the contract event 0xe89c6ba1e8957285aed22618f52aa1dcb9d5bb64e1533d8b55136c72fcf5aa5d.
//
// Solidity: event DelegateForToken(address vault, address delegate, address contract_, uint256 tokenId, bool value)
func (_Delegate *DelegateFilterer) ParseDelegateForToken(log types.Log) (*DelegateDelegateForToken, error) {
	event := new(DelegateDelegateForToken)
	if err := _Delegate.contract.UnpackLog(event, "DelegateForToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateRevokeAllDelegatesIterator is returned from FilterRevokeAllDelegates and is used to iterate over the raw logs and unpacked data for RevokeAllDelegates events raised by the Delegate contract.
type DelegateRevokeAllDelegatesIterator struct {
	Event *DelegateRevokeAllDelegates // Event containing the contract specifics and raw log

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
func (it *DelegateRevokeAllDelegatesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateRevokeAllDelegates)
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
		it.Event = new(DelegateRevokeAllDelegates)
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
func (it *DelegateRevokeAllDelegatesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateRevokeAllDelegatesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateRevokeAllDelegates represents a RevokeAllDelegates event raised by the Delegate contract.
type DelegateRevokeAllDelegates struct {
	Vault common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRevokeAllDelegates is a free log retrieval operation binding the contract event 0x32d74befd0b842e19694e3e3af46263e18bcce41352c8b600ff0002b49edf662.
//
// Solidity: event RevokeAllDelegates(address vault)
func (_Delegate *DelegateFilterer) FilterRevokeAllDelegates(opts *bind.FilterOpts) (*DelegateRevokeAllDelegatesIterator, error) {

	logs, sub, err := _Delegate.contract.FilterLogs(opts, "RevokeAllDelegates")
	if err != nil {
		return nil, err
	}
	return &DelegateRevokeAllDelegatesIterator{contract: _Delegate.contract, event: "RevokeAllDelegates", logs: logs, sub: sub}, nil
}

// WatchRevokeAllDelegates is a free log subscription operation binding the contract event 0x32d74befd0b842e19694e3e3af46263e18bcce41352c8b600ff0002b49edf662.
//
// Solidity: event RevokeAllDelegates(address vault)
func (_Delegate *DelegateFilterer) WatchRevokeAllDelegates(opts *bind.WatchOpts, sink chan<- *DelegateRevokeAllDelegates) (event.Subscription, error) {

	logs, sub, err := _Delegate.contract.WatchLogs(opts, "RevokeAllDelegates")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateRevokeAllDelegates)
				if err := _Delegate.contract.UnpackLog(event, "RevokeAllDelegates", log); err != nil {
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

// ParseRevokeAllDelegates is a log parse operation binding the contract event 0x32d74befd0b842e19694e3e3af46263e18bcce41352c8b600ff0002b49edf662.
//
// Solidity: event RevokeAllDelegates(address vault)
func (_Delegate *DelegateFilterer) ParseRevokeAllDelegates(log types.Log) (*DelegateRevokeAllDelegates, error) {
	event := new(DelegateRevokeAllDelegates)
	if err := _Delegate.contract.UnpackLog(event, "RevokeAllDelegates", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DelegateRevokeDelegateIterator is returned from FilterRevokeDelegate and is used to iterate over the raw logs and unpacked data for RevokeDelegate events raised by the Delegate contract.
type DelegateRevokeDelegateIterator struct {
	Event *DelegateRevokeDelegate // Event containing the contract specifics and raw log

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
func (it *DelegateRevokeDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegateRevokeDelegate)
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
		it.Event = new(DelegateRevokeDelegate)
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
func (it *DelegateRevokeDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegateRevokeDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegateRevokeDelegate represents a RevokeDelegate event raised by the Delegate contract.
type DelegateRevokeDelegate struct {
	Vault    common.Address
	Delegate common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRevokeDelegate is a free log retrieval operation binding the contract event 0x3e34a3ee53064fb79c0ee57448f03774a627a9270b0c41286efb7d8e32dcde93.
//
// Solidity: event RevokeDelegate(address vault, address delegate)
func (_Delegate *DelegateFilterer) FilterRevokeDelegate(opts *bind.FilterOpts) (*DelegateRevokeDelegateIterator, error) {

	logs, sub, err := _Delegate.contract.FilterLogs(opts, "RevokeDelegate")
	if err != nil {
		return nil, err
	}
	return &DelegateRevokeDelegateIterator{contract: _Delegate.contract, event: "RevokeDelegate", logs: logs, sub: sub}, nil
}

// WatchRevokeDelegate is a free log subscription operation binding the contract event 0x3e34a3ee53064fb79c0ee57448f03774a627a9270b0c41286efb7d8e32dcde93.
//
// Solidity: event RevokeDelegate(address vault, address delegate)
func (_Delegate *DelegateFilterer) WatchRevokeDelegate(opts *bind.WatchOpts, sink chan<- *DelegateRevokeDelegate) (event.Subscription, error) {

	logs, sub, err := _Delegate.contract.WatchLogs(opts, "RevokeDelegate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegateRevokeDelegate)
				if err := _Delegate.contract.UnpackLog(event, "RevokeDelegate", log); err != nil {
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

// ParseRevokeDelegate is a log parse operation binding the contract event 0x3e34a3ee53064fb79c0ee57448f03774a627a9270b0c41286efb7d8e32dcde93.
//
// Solidity: event RevokeDelegate(address vault, address delegate)
func (_Delegate *DelegateFilterer) ParseRevokeDelegate(log types.Log) (*DelegateRevokeDelegate, error) {
	event := new(DelegateRevokeDelegate)
	if err := _Delegate.contract.UnpackLog(event, "RevokeDelegate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
