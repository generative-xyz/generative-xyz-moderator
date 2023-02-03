// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_dao

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

// IGovernorCompatibilityBravoUpgradeableReceipt is an auto generated low-level Go binding around an user-defined struct.
type IGovernorCompatibilityBravoUpgradeableReceipt struct {
	HasVoted bool
	Support  uint8
	Votes    *big.Int
}

// GenerativeDaoMetaData contains all meta data concerning the GenerativeDao contract.
var GenerativeDaoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"Empty\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalCanceled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"indexed\":false,\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"ProposalCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"ProposalExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"}],\"name\":\"ProposalQueued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldQuorumNumerator\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newQuorumNumerator\",\"type\":\"uint256\"}],\"name\":\"QuorumNumeratorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"VoteCast\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"weight\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"VoteCastWithParams\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BALLOT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"COUNTING_MODE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"EXTENDED_BALLOT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_proposalThresholdPercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_quorumVotePercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_votingDelays\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_votingPeriods\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_votingToken\",\"outputs\":[{\"internalType\":\"contractIGENToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"}],\"name\":\"castVote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"castVoteBySig\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"}],\"name\":\"castVoteWithReason\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"castVoteWithReasonAndParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"reason\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"castVoteWithReasonAndParamsBySig\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new\",\"type\":\"uint256\"}],\"name\":\"changeProposalThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new\",\"type\":\"uint256\"}],\"name\":\"changeQuorumVotes\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new\",\"type\":\"uint256\"}],\"name\":\"changeVoteDelay\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_new\",\"type\":\"uint256\"}],\"name\":\"changeVotePeriod\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIGENToken\",\"name\":\"_new\",\"type\":\"address\"}],\"name\":\"changeVotingToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"descriptionHash\",\"type\":\"bytes32\"}],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"getActions\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"voter\",\"type\":\"address\"}],\"name\":\"getReceipt\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"hasVoted\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"support\",\"type\":\"uint8\"},{\"internalType\":\"uint96\",\"name\":\"votes\",\"type\":\"uint96\"}],\"internalType\":\"structIGovernorCompatibilityBravoUpgradeable.Receipt\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"getVotes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"params\",\"type\":\"bytes\"}],\"name\":\"getVotesWithParams\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasVoted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"descriptionHash\",\"type\":\"bytes32\"}],\"name\":\"hashProposal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramAddr\",\"type\":\"address\"},{\"internalType\":\"contractIGENToken\",\"name\":\"votingToken\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155BatchReceived\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC1155Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"proposalDeadline\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"proposalEta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"proposalSnapshot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposalThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"proposals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"eta\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"endBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"forVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"againstVotes\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"abstainVotes\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"canceled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"executed\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"signatures\",\"type\":\"string[]\"},{\"internalType\":\"bytes[]\",\"name\":\"calldatas\",\"type\":\"bytes[]\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"propose\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"queue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"quorum\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quorumDenominator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"quorumNumerator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quorumNumerator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"quorumVotes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"relay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"proposalId\",\"type\":\"uint256\"}],\"name\":\"state\",\"outputs\":[{\"internalType\":\"enumIGovernorUpgradeable.ProposalState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timelock\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIVotesUpgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newQuorumNumerator\",\"type\":\"uint256\"}],\"name\":\"updateQuorumNumerator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"votingDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"votingPeriod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"erc20Addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// GenerativeDaoABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeDaoMetaData.ABI instead.
var GenerativeDaoABI = GenerativeDaoMetaData.ABI

// GenerativeDao is an auto generated Go binding around an Ethereum contract.
type GenerativeDao struct {
	GenerativeDaoCaller     // Read-only binding to the contract
	GenerativeDaoTransactor // Write-only binding to the contract
	GenerativeDaoFilterer   // Log filterer for contract events
}

// GenerativeDaoCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeDaoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeDaoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeDaoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeDaoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeDaoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeDaoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeDaoSession struct {
	Contract     *GenerativeDao    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GenerativeDaoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeDaoCallerSession struct {
	Contract *GenerativeDaoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// GenerativeDaoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeDaoTransactorSession struct {
	Contract     *GenerativeDaoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// GenerativeDaoRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeDaoRaw struct {
	Contract *GenerativeDao // Generic contract binding to access the raw methods on
}

// GenerativeDaoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeDaoCallerRaw struct {
	Contract *GenerativeDaoCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeDaoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeDaoTransactorRaw struct {
	Contract *GenerativeDaoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeDao creates a new instance of GenerativeDao, bound to a specific deployed contract.
func NewGenerativeDao(address common.Address, backend bind.ContractBackend) (*GenerativeDao, error) {
	contract, err := bindGenerativeDao(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeDao{GenerativeDaoCaller: GenerativeDaoCaller{contract: contract}, GenerativeDaoTransactor: GenerativeDaoTransactor{contract: contract}, GenerativeDaoFilterer: GenerativeDaoFilterer{contract: contract}}, nil
}

// NewGenerativeDaoCaller creates a new read-only instance of GenerativeDao, bound to a specific deployed contract.
func NewGenerativeDaoCaller(address common.Address, caller bind.ContractCaller) (*GenerativeDaoCaller, error) {
	contract, err := bindGenerativeDao(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoCaller{contract: contract}, nil
}

// NewGenerativeDaoTransactor creates a new write-only instance of GenerativeDao, bound to a specific deployed contract.
func NewGenerativeDaoTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeDaoTransactor, error) {
	contract, err := bindGenerativeDao(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoTransactor{contract: contract}, nil
}

// NewGenerativeDaoFilterer creates a new log filterer instance of GenerativeDao, bound to a specific deployed contract.
func NewGenerativeDaoFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeDaoFilterer, error) {
	contract, err := bindGenerativeDao(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoFilterer{contract: contract}, nil
}

// bindGenerativeDao binds a generic wrapper to an already deployed contract.
func bindGenerativeDao(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeDaoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeDao *GenerativeDaoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeDao.Contract.GenerativeDaoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeDao *GenerativeDaoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeDao.Contract.GenerativeDaoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeDao *GenerativeDaoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeDao.Contract.GenerativeDaoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeDao *GenerativeDaoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeDao.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeDao *GenerativeDaoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeDao.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeDao *GenerativeDaoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeDao.Contract.contract.Transact(opts, method, params...)
}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoCaller) BALLOTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "BALLOT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoSession) BALLOTTYPEHASH() ([32]byte, error) {
	return _GenerativeDao.Contract.BALLOTTYPEHASH(&_GenerativeDao.CallOpts)
}

// BALLOTTYPEHASH is a free data retrieval call binding the contract method 0xdeaaa7cc.
//
// Solidity: function BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoCallerSession) BALLOTTYPEHASH() ([32]byte, error) {
	return _GenerativeDao.Contract.BALLOTTYPEHASH(&_GenerativeDao.CallOpts)
}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_GenerativeDao *GenerativeDaoCaller) COUNTINGMODE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "COUNTING_MODE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_GenerativeDao *GenerativeDaoSession) COUNTINGMODE() (string, error) {
	return _GenerativeDao.Contract.COUNTINGMODE(&_GenerativeDao.CallOpts)
}

// COUNTINGMODE is a free data retrieval call binding the contract method 0xdd4e2ba5.
//
// Solidity: function COUNTING_MODE() pure returns(string)
func (_GenerativeDao *GenerativeDaoCallerSession) COUNTINGMODE() (string, error) {
	return _GenerativeDao.Contract.COUNTINGMODE(&_GenerativeDao.CallOpts)
}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoCaller) EXTENDEDBALLOTTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "EXTENDED_BALLOT_TYPEHASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoSession) EXTENDEDBALLOTTYPEHASH() ([32]byte, error) {
	return _GenerativeDao.Contract.EXTENDEDBALLOTTYPEHASH(&_GenerativeDao.CallOpts)
}

// EXTENDEDBALLOTTYPEHASH is a free data retrieval call binding the contract method 0x2fe3e261.
//
// Solidity: function EXTENDED_BALLOT_TYPEHASH() view returns(bytes32)
func (_GenerativeDao *GenerativeDaoCallerSession) EXTENDEDBALLOTTYPEHASH() ([32]byte, error) {
	return _GenerativeDao.Contract.EXTENDEDBALLOTTYPEHASH(&_GenerativeDao.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeDao *GenerativeDaoCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeDao *GenerativeDaoSession) Admin() (common.Address, error) {
	return _GenerativeDao.Contract.Admin(&_GenerativeDao.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeDao *GenerativeDaoCallerSession) Admin() (common.Address, error) {
	return _GenerativeDao.Contract.Admin(&_GenerativeDao.CallOpts)
}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeDao *GenerativeDaoCaller) ParamAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_paramAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeDao *GenerativeDaoSession) ParamAddr() (common.Address, error) {
	return _GenerativeDao.Contract.ParamAddr(&_GenerativeDao.CallOpts)
}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeDao *GenerativeDaoCallerSession) ParamAddr() (common.Address, error) {
	return _GenerativeDao.Contract.ParamAddr(&_GenerativeDao.CallOpts)
}

// ProposalThresholdPercent is a free data retrieval call binding the contract method 0x0f1fe6ec.
//
// Solidity: function _proposalThresholdPercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) ProposalThresholdPercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_proposalThresholdPercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalThresholdPercent is a free data retrieval call binding the contract method 0x0f1fe6ec.
//
// Solidity: function _proposalThresholdPercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) ProposalThresholdPercent() (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalThresholdPercent(&_GenerativeDao.CallOpts)
}

// ProposalThresholdPercent is a free data retrieval call binding the contract method 0x0f1fe6ec.
//
// Solidity: function _proposalThresholdPercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) ProposalThresholdPercent() (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalThresholdPercent(&_GenerativeDao.CallOpts)
}

// QuorumVotePercent is a free data retrieval call binding the contract method 0xaa662214.
//
// Solidity: function _quorumVotePercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) QuorumVotePercent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_quorumVotePercent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumVotePercent is a free data retrieval call binding the contract method 0xaa662214.
//
// Solidity: function _quorumVotePercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) QuorumVotePercent() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumVotePercent(&_GenerativeDao.CallOpts)
}

// QuorumVotePercent is a free data retrieval call binding the contract method 0xaa662214.
//
// Solidity: function _quorumVotePercent() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) QuorumVotePercent() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumVotePercent(&_GenerativeDao.CallOpts)
}

// VotingDelays is a free data retrieval call binding the contract method 0x050c0299.
//
// Solidity: function _votingDelays() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) VotingDelays(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_votingDelays")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingDelays is a free data retrieval call binding the contract method 0x050c0299.
//
// Solidity: function _votingDelays() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) VotingDelays() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingDelays(&_GenerativeDao.CallOpts)
}

// VotingDelays is a free data retrieval call binding the contract method 0x050c0299.
//
// Solidity: function _votingDelays() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) VotingDelays() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingDelays(&_GenerativeDao.CallOpts)
}

// VotingPeriods is a free data retrieval call binding the contract method 0x58394825.
//
// Solidity: function _votingPeriods() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) VotingPeriods(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_votingPeriods")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingPeriods is a free data retrieval call binding the contract method 0x58394825.
//
// Solidity: function _votingPeriods() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) VotingPeriods() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingPeriods(&_GenerativeDao.CallOpts)
}

// VotingPeriods is a free data retrieval call binding the contract method 0x58394825.
//
// Solidity: function _votingPeriods() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) VotingPeriods() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingPeriods(&_GenerativeDao.CallOpts)
}

// VotingToken is a free data retrieval call binding the contract method 0xa4d82225.
//
// Solidity: function _votingToken() view returns(address)
func (_GenerativeDao *GenerativeDaoCaller) VotingToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "_votingToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VotingToken is a free data retrieval call binding the contract method 0xa4d82225.
//
// Solidity: function _votingToken() view returns(address)
func (_GenerativeDao *GenerativeDaoSession) VotingToken() (common.Address, error) {
	return _GenerativeDao.Contract.VotingToken(&_GenerativeDao.CallOpts)
}

// VotingToken is a free data retrieval call binding the contract method 0xa4d82225.
//
// Solidity: function _votingToken() view returns(address)
func (_GenerativeDao *GenerativeDaoCallerSession) VotingToken() (common.Address, error) {
	return _GenerativeDao.Contract.VotingToken(&_GenerativeDao.CallOpts)
}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_GenerativeDao *GenerativeDaoCaller) GetActions(opts *bind.CallOpts, proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "getActions", proposalId)

	outstruct := new(struct {
		Targets    []common.Address
		Values     []*big.Int
		Signatures []string
		Calldatas  [][]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Targets = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Values = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.Signatures = *abi.ConvertType(out[2], new([]string)).(*[]string)
	outstruct.Calldatas = *abi.ConvertType(out[3], new([][]byte)).(*[][]byte)

	return *outstruct, err

}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_GenerativeDao *GenerativeDaoSession) GetActions(proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	return _GenerativeDao.Contract.GetActions(&_GenerativeDao.CallOpts, proposalId)
}

// GetActions is a free data retrieval call binding the contract method 0x328dd982.
//
// Solidity: function getActions(uint256 proposalId) view returns(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas)
func (_GenerativeDao *GenerativeDaoCallerSession) GetActions(proposalId *big.Int) (struct {
	Targets    []common.Address
	Values     []*big.Int
	Signatures []string
	Calldatas  [][]byte
}, error) {
	return _GenerativeDao.Contract.GetActions(&_GenerativeDao.CallOpts, proposalId)
}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint96))
func (_GenerativeDao *GenerativeDaoCaller) GetReceipt(opts *bind.CallOpts, proposalId *big.Int, voter common.Address) (IGovernorCompatibilityBravoUpgradeableReceipt, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "getReceipt", proposalId, voter)

	if err != nil {
		return *new(IGovernorCompatibilityBravoUpgradeableReceipt), err
	}

	out0 := *abi.ConvertType(out[0], new(IGovernorCompatibilityBravoUpgradeableReceipt)).(*IGovernorCompatibilityBravoUpgradeableReceipt)

	return out0, err

}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint96))
func (_GenerativeDao *GenerativeDaoSession) GetReceipt(proposalId *big.Int, voter common.Address) (IGovernorCompatibilityBravoUpgradeableReceipt, error) {
	return _GenerativeDao.Contract.GetReceipt(&_GenerativeDao.CallOpts, proposalId, voter)
}

// GetReceipt is a free data retrieval call binding the contract method 0xe23a9a52.
//
// Solidity: function getReceipt(uint256 proposalId, address voter) view returns((bool,uint8,uint96))
func (_GenerativeDao *GenerativeDaoCallerSession) GetReceipt(proposalId *big.Int, voter common.Address) (IGovernorCompatibilityBravoUpgradeableReceipt, error) {
	return _GenerativeDao.Contract.GetReceipt(&_GenerativeDao.CallOpts, proposalId, voter)
}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) GetVotes(opts *bind.CallOpts, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "getVotes", account, blockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) GetVotes(account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.GetVotes(&_GenerativeDao.CallOpts, account, blockNumber)
}

// GetVotes is a free data retrieval call binding the contract method 0xeb9019d4.
//
// Solidity: function getVotes(address account, uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) GetVotes(account common.Address, blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.GetVotes(&_GenerativeDao.CallOpts, account, blockNumber)
}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 blockNumber, bytes params) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) GetVotesWithParams(opts *bind.CallOpts, account common.Address, blockNumber *big.Int, params []byte) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "getVotesWithParams", account, blockNumber, params)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 blockNumber, bytes params) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) GetVotesWithParams(account common.Address, blockNumber *big.Int, params []byte) (*big.Int, error) {
	return _GenerativeDao.Contract.GetVotesWithParams(&_GenerativeDao.CallOpts, account, blockNumber, params)
}

// GetVotesWithParams is a free data retrieval call binding the contract method 0x9a802a6d.
//
// Solidity: function getVotesWithParams(address account, uint256 blockNumber, bytes params) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) GetVotesWithParams(account common.Address, blockNumber *big.Int, params []byte) (*big.Int, error) {
	return _GenerativeDao.Contract.GetVotesWithParams(&_GenerativeDao.CallOpts, account, blockNumber, params)
}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 proposalId, address account) view returns(bool)
func (_GenerativeDao *GenerativeDaoCaller) HasVoted(opts *bind.CallOpts, proposalId *big.Int, account common.Address) (bool, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "hasVoted", proposalId, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 proposalId, address account) view returns(bool)
func (_GenerativeDao *GenerativeDaoSession) HasVoted(proposalId *big.Int, account common.Address) (bool, error) {
	return _GenerativeDao.Contract.HasVoted(&_GenerativeDao.CallOpts, proposalId, account)
}

// HasVoted is a free data retrieval call binding the contract method 0x43859632.
//
// Solidity: function hasVoted(uint256 proposalId, address account) view returns(bool)
func (_GenerativeDao *GenerativeDaoCallerSession) HasVoted(proposalId *big.Int, account common.Address) (bool, error) {
	return _GenerativeDao.Contract.HasVoted(&_GenerativeDao.CallOpts, proposalId, account)
}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) HashProposal(opts *bind.CallOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "hashProposal", targets, values, calldatas, descriptionHash)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) HashProposal(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	return _GenerativeDao.Contract.HashProposal(&_GenerativeDao.CallOpts, targets, values, calldatas, descriptionHash)
}

// HashProposal is a free data retrieval call binding the contract method 0xc59057e4.
//
// Solidity: function hashProposal(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) HashProposal(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*big.Int, error) {
	return _GenerativeDao.Contract.HashProposal(&_GenerativeDao.CallOpts, targets, values, calldatas, descriptionHash)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeDao *GenerativeDaoCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeDao *GenerativeDaoSession) Name() (string, error) {
	return _GenerativeDao.Contract.Name(&_GenerativeDao.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GenerativeDao *GenerativeDaoCallerSession) Name() (string, error) {
	return _GenerativeDao.Contract.Name(&_GenerativeDao.CallOpts)
}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) ProposalDeadline(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "proposalDeadline", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) ProposalDeadline(proposalId *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalDeadline(&_GenerativeDao.CallOpts, proposalId)
}

// ProposalDeadline is a free data retrieval call binding the contract method 0xc01f9e37.
//
// Solidity: function proposalDeadline(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) ProposalDeadline(proposalId *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalDeadline(&_GenerativeDao.CallOpts, proposalId)
}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) ProposalEta(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "proposalEta", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) ProposalEta(arg0 *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalEta(&_GenerativeDao.CallOpts, arg0)
}

// ProposalEta is a free data retrieval call binding the contract method 0xab58fb8e.
//
// Solidity: function proposalEta(uint256 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) ProposalEta(arg0 *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalEta(&_GenerativeDao.CallOpts, arg0)
}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) ProposalSnapshot(opts *bind.CallOpts, proposalId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "proposalSnapshot", proposalId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) ProposalSnapshot(proposalId *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalSnapshot(&_GenerativeDao.CallOpts, proposalId)
}

// ProposalSnapshot is a free data retrieval call binding the contract method 0x2d63f693.
//
// Solidity: function proposalSnapshot(uint256 proposalId) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) ProposalSnapshot(proposalId *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalSnapshot(&_GenerativeDao.CallOpts, proposalId)
}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) ProposalThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "proposalThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) ProposalThreshold() (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalThreshold(&_GenerativeDao.CallOpts)
}

// ProposalThreshold is a free data retrieval call binding the contract method 0xb58131b0.
//
// Solidity: function proposalThreshold() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) ProposalThreshold() (*big.Int, error) {
	return _GenerativeDao.Contract.ProposalThreshold(&_GenerativeDao.CallOpts)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 proposalId) view returns(uint256 id, address proposer, uint256 eta, uint256 startBlock, uint256 endBlock, uint256 forVotes, uint256 againstVotes, uint256 abstainVotes, bool canceled, bool executed)
func (_GenerativeDao *GenerativeDaoCaller) Proposals(opts *bind.CallOpts, proposalId *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Eta          *big.Int
	StartBlock   *big.Int
	EndBlock     *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	AbstainVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "proposals", proposalId)

	outstruct := new(struct {
		Id           *big.Int
		Proposer     common.Address
		Eta          *big.Int
		StartBlock   *big.Int
		EndBlock     *big.Int
		ForVotes     *big.Int
		AgainstVotes *big.Int
		AbstainVotes *big.Int
		Canceled     bool
		Executed     bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Proposer = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Eta = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.StartBlock = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.EndBlock = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.ForVotes = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.AgainstVotes = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.AbstainVotes = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Canceled = *abi.ConvertType(out[8], new(bool)).(*bool)
	outstruct.Executed = *abi.ConvertType(out[9], new(bool)).(*bool)

	return *outstruct, err

}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 proposalId) view returns(uint256 id, address proposer, uint256 eta, uint256 startBlock, uint256 endBlock, uint256 forVotes, uint256 againstVotes, uint256 abstainVotes, bool canceled, bool executed)
func (_GenerativeDao *GenerativeDaoSession) Proposals(proposalId *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Eta          *big.Int
	StartBlock   *big.Int
	EndBlock     *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	AbstainVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	return _GenerativeDao.Contract.Proposals(&_GenerativeDao.CallOpts, proposalId)
}

// Proposals is a free data retrieval call binding the contract method 0x013cf08b.
//
// Solidity: function proposals(uint256 proposalId) view returns(uint256 id, address proposer, uint256 eta, uint256 startBlock, uint256 endBlock, uint256 forVotes, uint256 againstVotes, uint256 abstainVotes, bool canceled, bool executed)
func (_GenerativeDao *GenerativeDaoCallerSession) Proposals(proposalId *big.Int) (struct {
	Id           *big.Int
	Proposer     common.Address
	Eta          *big.Int
	StartBlock   *big.Int
	EndBlock     *big.Int
	ForVotes     *big.Int
	AgainstVotes *big.Int
	AbstainVotes *big.Int
	Canceled     bool
	Executed     bool
}, error) {
	return _GenerativeDao.Contract.Proposals(&_GenerativeDao.CallOpts, proposalId)
}

// Queue is a free data retrieval call binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] , uint256[] , bytes[] , bytes32 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) Queue(opts *bind.CallOpts, arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "queue", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Queue is a free data retrieval call binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] , uint256[] , bytes[] , bytes32 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) Queue(arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 [32]byte) (*big.Int, error) {
	return _GenerativeDao.Contract.Queue(&_GenerativeDao.CallOpts, arg0, arg1, arg2, arg3)
}

// Queue is a free data retrieval call binding the contract method 0x160cbed7.
//
// Solidity: function queue(address[] , uint256[] , bytes[] , bytes32 ) pure returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) Queue(arg0 []common.Address, arg1 []*big.Int, arg2 [][]byte, arg3 [32]byte) (*big.Int, error) {
	return _GenerativeDao.Contract.Queue(&_GenerativeDao.CallOpts, arg0, arg1, arg2, arg3)
}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) Quorum(opts *bind.CallOpts, blockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "quorum", blockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) Quorum(blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.Quorum(&_GenerativeDao.CallOpts, blockNumber)
}

// Quorum is a free data retrieval call binding the contract method 0xf8ce560a.
//
// Solidity: function quorum(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) Quorum(blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.Quorum(&_GenerativeDao.CallOpts, blockNumber)
}

// QuorumDenominator is a free data retrieval call binding the contract method 0x97c3d334.
//
// Solidity: function quorumDenominator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) QuorumDenominator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "quorumDenominator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumDenominator is a free data retrieval call binding the contract method 0x97c3d334.
//
// Solidity: function quorumDenominator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) QuorumDenominator() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumDenominator(&_GenerativeDao.CallOpts)
}

// QuorumDenominator is a free data retrieval call binding the contract method 0x97c3d334.
//
// Solidity: function quorumDenominator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) QuorumDenominator() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumDenominator(&_GenerativeDao.CallOpts)
}

// QuorumNumerator is a free data retrieval call binding the contract method 0x60c4247f.
//
// Solidity: function quorumNumerator(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) QuorumNumerator(opts *bind.CallOpts, blockNumber *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "quorumNumerator", blockNumber)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumNumerator is a free data retrieval call binding the contract method 0x60c4247f.
//
// Solidity: function quorumNumerator(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) QuorumNumerator(blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumNumerator(&_GenerativeDao.CallOpts, blockNumber)
}

// QuorumNumerator is a free data retrieval call binding the contract method 0x60c4247f.
//
// Solidity: function quorumNumerator(uint256 blockNumber) view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) QuorumNumerator(blockNumber *big.Int) (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumNumerator(&_GenerativeDao.CallOpts, blockNumber)
}

// QuorumNumerator0 is a free data retrieval call binding the contract method 0xa7713a70.
//
// Solidity: function quorumNumerator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) QuorumNumerator0(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "quorumNumerator0")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumNumerator0 is a free data retrieval call binding the contract method 0xa7713a70.
//
// Solidity: function quorumNumerator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) QuorumNumerator0() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumNumerator0(&_GenerativeDao.CallOpts)
}

// QuorumNumerator0 is a free data retrieval call binding the contract method 0xa7713a70.
//
// Solidity: function quorumNumerator() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) QuorumNumerator0() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumNumerator0(&_GenerativeDao.CallOpts)
}

// QuorumVotes is a free data retrieval call binding the contract method 0x24bc1a64.
//
// Solidity: function quorumVotes() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) QuorumVotes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "quorumVotes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// QuorumVotes is a free data retrieval call binding the contract method 0x24bc1a64.
//
// Solidity: function quorumVotes() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) QuorumVotes() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumVotes(&_GenerativeDao.CallOpts)
}

// QuorumVotes is a free data retrieval call binding the contract method 0x24bc1a64.
//
// Solidity: function quorumVotes() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) QuorumVotes() (*big.Int, error) {
	return _GenerativeDao.Contract.QuorumVotes(&_GenerativeDao.CallOpts)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_GenerativeDao *GenerativeDaoCaller) State(opts *bind.CallOpts, proposalId *big.Int) (uint8, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "state", proposalId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_GenerativeDao *GenerativeDaoSession) State(proposalId *big.Int) (uint8, error) {
	return _GenerativeDao.Contract.State(&_GenerativeDao.CallOpts, proposalId)
}

// State is a free data retrieval call binding the contract method 0x3e4f49e6.
//
// Solidity: function state(uint256 proposalId) view returns(uint8)
func (_GenerativeDao *GenerativeDaoCallerSession) State(proposalId *big.Int) (uint8, error) {
	return _GenerativeDao.Contract.State(&_GenerativeDao.CallOpts, proposalId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeDao *GenerativeDaoCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeDao *GenerativeDaoSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeDao.Contract.SupportsInterface(&_GenerativeDao.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_GenerativeDao *GenerativeDaoCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _GenerativeDao.Contract.SupportsInterface(&_GenerativeDao.CallOpts, interfaceId)
}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() pure returns(address)
func (_GenerativeDao *GenerativeDaoCaller) Timelock(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "timelock")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() pure returns(address)
func (_GenerativeDao *GenerativeDaoSession) Timelock() (common.Address, error) {
	return _GenerativeDao.Contract.Timelock(&_GenerativeDao.CallOpts)
}

// Timelock is a free data retrieval call binding the contract method 0xd33219b4.
//
// Solidity: function timelock() pure returns(address)
func (_GenerativeDao *GenerativeDaoCallerSession) Timelock() (common.Address, error) {
	return _GenerativeDao.Contract.Timelock(&_GenerativeDao.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GenerativeDao *GenerativeDaoCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GenerativeDao *GenerativeDaoSession) Token() (common.Address, error) {
	return _GenerativeDao.Contract.Token(&_GenerativeDao.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_GenerativeDao *GenerativeDaoCallerSession) Token() (common.Address, error) {
	return _GenerativeDao.Contract.Token(&_GenerativeDao.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GenerativeDao *GenerativeDaoCaller) Version(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GenerativeDao *GenerativeDaoSession) Version() (string, error) {
	return _GenerativeDao.Contract.Version(&_GenerativeDao.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(string)
func (_GenerativeDao *GenerativeDaoCallerSession) Version() (string, error) {
	return _GenerativeDao.Contract.Version(&_GenerativeDao.CallOpts)
}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) VotingDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "votingDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) VotingDelay() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingDelay(&_GenerativeDao.CallOpts)
}

// VotingDelay is a free data retrieval call binding the contract method 0x3932abb1.
//
// Solidity: function votingDelay() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) VotingDelay() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingDelay(&_GenerativeDao.CallOpts)
}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCaller) VotingPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GenerativeDao.contract.Call(opts, &out, "votingPeriod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) VotingPeriod() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingPeriod(&_GenerativeDao.CallOpts)
}

// VotingPeriod is a free data retrieval call binding the contract method 0x02a251a3.
//
// Solidity: function votingPeriod() view returns(uint256)
func (_GenerativeDao *GenerativeDaoCallerSession) VotingPeriod() (*big.Int, error) {
	return _GenerativeDao.Contract.VotingPeriod(&_GenerativeDao.CallOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoTransactor) Cancel(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "cancel", proposalId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoSession) Cancel(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Cancel(&_GenerativeDao.TransactOpts, proposalId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Cancel(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Cancel(&_GenerativeDao.TransactOpts, proposalId)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 support) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) CastVote(opts *bind.TransactOpts, proposalId *big.Int, support uint8) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "castVote", proposalId, support)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 support) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) CastVote(proposalId *big.Int, support uint8) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVote(&_GenerativeDao.TransactOpts, proposalId, support)
}

// CastVote is a paid mutator transaction binding the contract method 0x56781388.
//
// Solidity: function castVote(uint256 proposalId, uint8 support) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) CastVote(proposalId *big.Int, support uint8) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVote(&_GenerativeDao.TransactOpts, proposalId, support)
}

// CastVoteBySig is a paid mutator transaction binding the contract method 0x3bccf4fd.
//
// Solidity: function castVoteBySig(uint256 proposalId, uint8 support, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) CastVoteBySig(opts *bind.TransactOpts, proposalId *big.Int, support uint8, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "castVoteBySig", proposalId, support, v, r, s)
}

// CastVoteBySig is a paid mutator transaction binding the contract method 0x3bccf4fd.
//
// Solidity: function castVoteBySig(uint256 proposalId, uint8 support, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) CastVoteBySig(proposalId *big.Int, support uint8, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteBySig(&_GenerativeDao.TransactOpts, proposalId, support, v, r, s)
}

// CastVoteBySig is a paid mutator transaction binding the contract method 0x3bccf4fd.
//
// Solidity: function castVoteBySig(uint256 proposalId, uint8 support, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) CastVoteBySig(proposalId *big.Int, support uint8, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteBySig(&_GenerativeDao.TransactOpts, proposalId, support, v, r, s)
}

// CastVoteWithReason is a paid mutator transaction binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 proposalId, uint8 support, string reason) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) CastVoteWithReason(opts *bind.TransactOpts, proposalId *big.Int, support uint8, reason string) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "castVoteWithReason", proposalId, support, reason)
}

// CastVoteWithReason is a paid mutator transaction binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 proposalId, uint8 support, string reason) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) CastVoteWithReason(proposalId *big.Int, support uint8, reason string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReason(&_GenerativeDao.TransactOpts, proposalId, support, reason)
}

// CastVoteWithReason is a paid mutator transaction binding the contract method 0x7b3c71d3.
//
// Solidity: function castVoteWithReason(uint256 proposalId, uint8 support, string reason) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) CastVoteWithReason(proposalId *big.Int, support uint8, reason string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReason(&_GenerativeDao.TransactOpts, proposalId, support, reason)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) CastVoteWithReasonAndParams(opts *bind.TransactOpts, proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "castVoteWithReasonAndParams", proposalId, support, reason, params)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) CastVoteWithReasonAndParams(proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReasonAndParams(&_GenerativeDao.TransactOpts, proposalId, support, reason, params)
}

// CastVoteWithReasonAndParams is a paid mutator transaction binding the contract method 0x5f398a14.
//
// Solidity: function castVoteWithReasonAndParams(uint256 proposalId, uint8 support, string reason, bytes params) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) CastVoteWithReasonAndParams(proposalId *big.Int, support uint8, reason string, params []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReasonAndParams(&_GenerativeDao.TransactOpts, proposalId, support, reason, params)
}

// CastVoteWithReasonAndParamsBySig is a paid mutator transaction binding the contract method 0x03420181.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 proposalId, uint8 support, string reason, bytes params, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) CastVoteWithReasonAndParamsBySig(opts *bind.TransactOpts, proposalId *big.Int, support uint8, reason string, params []byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "castVoteWithReasonAndParamsBySig", proposalId, support, reason, params, v, r, s)
}

// CastVoteWithReasonAndParamsBySig is a paid mutator transaction binding the contract method 0x03420181.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 proposalId, uint8 support, string reason, bytes params, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) CastVoteWithReasonAndParamsBySig(proposalId *big.Int, support uint8, reason string, params []byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReasonAndParamsBySig(&_GenerativeDao.TransactOpts, proposalId, support, reason, params, v, r, s)
}

// CastVoteWithReasonAndParamsBySig is a paid mutator transaction binding the contract method 0x03420181.
//
// Solidity: function castVoteWithReasonAndParamsBySig(uint256 proposalId, uint8 support, string reason, bytes params, uint8 v, bytes32 r, bytes32 s) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) CastVoteWithReasonAndParamsBySig(proposalId *big.Int, support uint8, reason string, params []byte, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.CastVoteWithReasonAndParamsBySig(&_GenerativeDao.TransactOpts, proposalId, support, reason, params, v, r, s)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeAdmin(&_GenerativeDao.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeAdmin(&_GenerativeDao.TransactOpts, newAdm)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeParamAddress(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeParamAddress", newAddr)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeParamAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeParamAddress(&_GenerativeDao.TransactOpts, newAddr)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeParamAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeParamAddress(&_GenerativeDao.TransactOpts, newAddr)
}

// ChangeProposalThreshold is a paid mutator transaction binding the contract method 0xf9dabc55.
//
// Solidity: function changeProposalThreshold(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeProposalThreshold(opts *bind.TransactOpts, _new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeProposalThreshold", _new)
}

// ChangeProposalThreshold is a paid mutator transaction binding the contract method 0xf9dabc55.
//
// Solidity: function changeProposalThreshold(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeProposalThreshold(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeProposalThreshold(&_GenerativeDao.TransactOpts, _new)
}

// ChangeProposalThreshold is a paid mutator transaction binding the contract method 0xf9dabc55.
//
// Solidity: function changeProposalThreshold(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeProposalThreshold(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeProposalThreshold(&_GenerativeDao.TransactOpts, _new)
}

// ChangeQuorumVotes is a paid mutator transaction binding the contract method 0xa8f0a82e.
//
// Solidity: function changeQuorumVotes(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeQuorumVotes(opts *bind.TransactOpts, _new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeQuorumVotes", _new)
}

// ChangeQuorumVotes is a paid mutator transaction binding the contract method 0xa8f0a82e.
//
// Solidity: function changeQuorumVotes(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeQuorumVotes(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeQuorumVotes(&_GenerativeDao.TransactOpts, _new)
}

// ChangeQuorumVotes is a paid mutator transaction binding the contract method 0xa8f0a82e.
//
// Solidity: function changeQuorumVotes(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeQuorumVotes(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeQuorumVotes(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVoteDelay is a paid mutator transaction binding the contract method 0xb3af0777.
//
// Solidity: function changeVoteDelay(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeVoteDelay(opts *bind.TransactOpts, _new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeVoteDelay", _new)
}

// ChangeVoteDelay is a paid mutator transaction binding the contract method 0xb3af0777.
//
// Solidity: function changeVoteDelay(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeVoteDelay(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVoteDelay(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVoteDelay is a paid mutator transaction binding the contract method 0xb3af0777.
//
// Solidity: function changeVoteDelay(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeVoteDelay(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVoteDelay(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVotePeriod is a paid mutator transaction binding the contract method 0xb15c4a42.
//
// Solidity: function changeVotePeriod(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeVotePeriod(opts *bind.TransactOpts, _new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeVotePeriod", _new)
}

// ChangeVotePeriod is a paid mutator transaction binding the contract method 0xb15c4a42.
//
// Solidity: function changeVotePeriod(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeVotePeriod(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVotePeriod(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVotePeriod is a paid mutator transaction binding the contract method 0xb15c4a42.
//
// Solidity: function changeVotePeriod(uint256 _new) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeVotePeriod(_new *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVotePeriod(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVotingToken is a paid mutator transaction binding the contract method 0x54dfc2ed.
//
// Solidity: function changeVotingToken(address _new) returns()
func (_GenerativeDao *GenerativeDaoTransactor) ChangeVotingToken(opts *bind.TransactOpts, _new common.Address) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "changeVotingToken", _new)
}

// ChangeVotingToken is a paid mutator transaction binding the contract method 0x54dfc2ed.
//
// Solidity: function changeVotingToken(address _new) returns()
func (_GenerativeDao *GenerativeDaoSession) ChangeVotingToken(_new common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVotingToken(&_GenerativeDao.TransactOpts, _new)
}

// ChangeVotingToken is a paid mutator transaction binding the contract method 0x54dfc2ed.
//
// Solidity: function changeVotingToken(address _new) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) ChangeVotingToken(_new common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.ChangeVotingToken(&_GenerativeDao.TransactOpts, _new)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) Execute(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "execute", targets, values, calldatas, descriptionHash)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) Execute(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Execute(&_GenerativeDao.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Execute is a paid mutator transaction binding the contract method 0x2656227d.
//
// Solidity: function execute(address[] targets, uint256[] values, bytes[] calldatas, bytes32 descriptionHash) payable returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) Execute(targets []common.Address, values []*big.Int, calldatas [][]byte, descriptionHash [32]byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Execute(&_GenerativeDao.TransactOpts, targets, values, calldatas, descriptionHash)
}

// Execute0 is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_GenerativeDao *GenerativeDaoTransactor) Execute0(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "execute0", proposalId)
}

// Execute0 is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_GenerativeDao *GenerativeDaoSession) Execute0(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Execute0(&_GenerativeDao.TransactOpts, proposalId)
}

// Execute0 is a paid mutator transaction binding the contract method 0xfe0d94c1.
//
// Solidity: function execute(uint256 proposalId) payable returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Execute0(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Execute0(&_GenerativeDao.TransactOpts, proposalId)
}

// Initialize is a paid mutator transaction binding the contract method 0xf34822b4.
//
// Solidity: function initialize(string name, address admin, address paramAddr, address votingToken) returns()
func (_GenerativeDao *GenerativeDaoTransactor) Initialize(opts *bind.TransactOpts, name string, admin common.Address, paramAddr common.Address, votingToken common.Address) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "initialize", name, admin, paramAddr, votingToken)
}

// Initialize is a paid mutator transaction binding the contract method 0xf34822b4.
//
// Solidity: function initialize(string name, address admin, address paramAddr, address votingToken) returns()
func (_GenerativeDao *GenerativeDaoSession) Initialize(name string, admin common.Address, paramAddr common.Address, votingToken common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Initialize(&_GenerativeDao.TransactOpts, name, admin, paramAddr, votingToken)
}

// Initialize is a paid mutator transaction binding the contract method 0xf34822b4.
//
// Solidity: function initialize(string name, address admin, address paramAddr, address votingToken) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Initialize(name string, admin common.Address, paramAddr common.Address, votingToken common.Address) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Initialize(&_GenerativeDao.TransactOpts, name, admin, paramAddr, votingToken)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactor) OnERC1155BatchReceived(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "onERC1155BatchReceived", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC1155BatchReceived(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155BatchReceived is a paid mutator transaction binding the contract method 0xbc197c81.
//
// Solidity: function onERC1155BatchReceived(address , address , uint256[] , uint256[] , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactorSession) OnERC1155BatchReceived(arg0 common.Address, arg1 common.Address, arg2 []*big.Int, arg3 []*big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC1155BatchReceived(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactor) OnERC1155Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "onERC1155Received", arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC1155Received(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC1155Received is a paid mutator transaction binding the contract method 0xf23a6e61.
//
// Solidity: function onERC1155Received(address , address , uint256 , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactorSession) OnERC1155Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 *big.Int, arg4 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC1155Received(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3, arg4)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactor) OnERC721Received(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "onERC721Received", arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC721Received(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a paid mutator transaction binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) returns(bytes4)
func (_GenerativeDao *GenerativeDaoTransactorSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.OnERC721Received(&_GenerativeDao.TransactOpts, arg0, arg1, arg2, arg3)
}

// Propose is a paid mutator transaction binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] targets, uint256[] values, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) Propose(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "propose", targets, values, calldatas, description)
}

// Propose is a paid mutator transaction binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] targets, uint256[] values, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) Propose(targets []common.Address, values []*big.Int, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Propose(&_GenerativeDao.TransactOpts, targets, values, calldatas, description)
}

// Propose is a paid mutator transaction binding the contract method 0x7d5e81e2.
//
// Solidity: function propose(address[] targets, uint256[] values, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) Propose(targets []common.Address, values []*big.Int, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Propose(&_GenerativeDao.TransactOpts, targets, values, calldatas, description)
}

// Propose0 is a paid mutator transaction binding the contract method 0xda95691a.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactor) Propose0(opts *bind.TransactOpts, targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "propose0", targets, values, signatures, calldatas, description)
}

// Propose0 is a paid mutator transaction binding the contract method 0xda95691a.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoSession) Propose0(targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Propose0(&_GenerativeDao.TransactOpts, targets, values, signatures, calldatas, description)
}

// Propose0 is a paid mutator transaction binding the contract method 0xda95691a.
//
// Solidity: function propose(address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, string description) returns(uint256)
func (_GenerativeDao *GenerativeDaoTransactorSession) Propose0(targets []common.Address, values []*big.Int, signatures []string, calldatas [][]byte, description string) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Propose0(&_GenerativeDao.TransactOpts, targets, values, signatures, calldatas, description)
}

// Queue0 is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoTransactor) Queue0(opts *bind.TransactOpts, proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "queue0", proposalId)
}

// Queue0 is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoSession) Queue0(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Queue0(&_GenerativeDao.TransactOpts, proposalId)
}

// Queue0 is a paid mutator transaction binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 proposalId) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Queue0(proposalId *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Queue0(&_GenerativeDao.TransactOpts, proposalId)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_GenerativeDao *GenerativeDaoTransactor) Relay(opts *bind.TransactOpts, target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "relay", target, value, data)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_GenerativeDao *GenerativeDaoSession) Relay(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Relay(&_GenerativeDao.TransactOpts, target, value, data)
}

// Relay is a paid mutator transaction binding the contract method 0xc28bc2fa.
//
// Solidity: function relay(address target, uint256 value, bytes data) payable returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Relay(target common.Address, value *big.Int, data []byte) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Relay(&_GenerativeDao.TransactOpts, target, value, data)
}

// UpdateQuorumNumerator is a paid mutator transaction binding the contract method 0x06f3f9e6.
//
// Solidity: function updateQuorumNumerator(uint256 newQuorumNumerator) returns()
func (_GenerativeDao *GenerativeDaoTransactor) UpdateQuorumNumerator(opts *bind.TransactOpts, newQuorumNumerator *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "updateQuorumNumerator", newQuorumNumerator)
}

// UpdateQuorumNumerator is a paid mutator transaction binding the contract method 0x06f3f9e6.
//
// Solidity: function updateQuorumNumerator(uint256 newQuorumNumerator) returns()
func (_GenerativeDao *GenerativeDaoSession) UpdateQuorumNumerator(newQuorumNumerator *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.UpdateQuorumNumerator(&_GenerativeDao.TransactOpts, newQuorumNumerator)
}

// UpdateQuorumNumerator is a paid mutator transaction binding the contract method 0x06f3f9e6.
//
// Solidity: function updateQuorumNumerator(uint256 newQuorumNumerator) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) UpdateQuorumNumerator(newQuorumNumerator *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.UpdateQuorumNumerator(&_GenerativeDao.TransactOpts, newQuorumNumerator)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address erc20Addr, uint256 amount) returns()
func (_GenerativeDao *GenerativeDaoTransactor) Withdraw(opts *bind.TransactOpts, erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.contract.Transact(opts, "withdraw", erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address erc20Addr, uint256 amount) returns()
func (_GenerativeDao *GenerativeDaoSession) Withdraw(erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Withdraw(&_GenerativeDao.TransactOpts, erc20Addr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address erc20Addr, uint256 amount) returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Withdraw(erc20Addr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _GenerativeDao.Contract.Withdraw(&_GenerativeDao.TransactOpts, erc20Addr, amount)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GenerativeDao *GenerativeDaoTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeDao.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GenerativeDao *GenerativeDaoSession) Receive() (*types.Transaction, error) {
	return _GenerativeDao.Contract.Receive(&_GenerativeDao.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_GenerativeDao *GenerativeDaoTransactorSession) Receive() (*types.Transaction, error) {
	return _GenerativeDao.Contract.Receive(&_GenerativeDao.TransactOpts)
}

// GenerativeDaoInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GenerativeDao contract.
type GenerativeDaoInitializedIterator struct {
	Event *GenerativeDaoInitialized // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoInitialized)
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
		it.Event = new(GenerativeDaoInitialized)
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
func (it *GenerativeDaoInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoInitialized represents a Initialized event raised by the GenerativeDao contract.
type GenerativeDaoInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeDao *GenerativeDaoFilterer) FilterInitialized(opts *bind.FilterOpts) (*GenerativeDaoInitializedIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoInitializedIterator{contract: _GenerativeDao.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeDao *GenerativeDaoFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GenerativeDaoInitialized) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoInitialized)
				if err := _GenerativeDao.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GenerativeDao *GenerativeDaoFilterer) ParseInitialized(log types.Log) (*GenerativeDaoInitialized, error) {
	event := new(GenerativeDaoInitialized)
	if err := _GenerativeDao.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoProposalCanceledIterator is returned from FilterProposalCanceled and is used to iterate over the raw logs and unpacked data for ProposalCanceled events raised by the GenerativeDao contract.
type GenerativeDaoProposalCanceledIterator struct {
	Event *GenerativeDaoProposalCanceled // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoProposalCanceledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoProposalCanceled)
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
		it.Event = new(GenerativeDaoProposalCanceled)
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
func (it *GenerativeDaoProposalCanceledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoProposalCanceledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoProposalCanceled represents a ProposalCanceled event raised by the GenerativeDao contract.
type GenerativeDaoProposalCanceled struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalCanceled is a free log retrieval operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) FilterProposalCanceled(opts *bind.FilterOpts) (*GenerativeDaoProposalCanceledIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "ProposalCanceled")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoProposalCanceledIterator{contract: _GenerativeDao.contract, event: "ProposalCanceled", logs: logs, sub: sub}, nil
}

// WatchProposalCanceled is a free log subscription operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) WatchProposalCanceled(opts *bind.WatchOpts, sink chan<- *GenerativeDaoProposalCanceled) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "ProposalCanceled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoProposalCanceled)
				if err := _GenerativeDao.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
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

// ParseProposalCanceled is a log parse operation binding the contract event 0x789cf55be980739dad1d0699b93b58e806b51c9d96619bfa8fe0a28abaa7b30c.
//
// Solidity: event ProposalCanceled(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) ParseProposalCanceled(log types.Log) (*GenerativeDaoProposalCanceled, error) {
	event := new(GenerativeDaoProposalCanceled)
	if err := _GenerativeDao.contract.UnpackLog(event, "ProposalCanceled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoProposalCreatedIterator is returned from FilterProposalCreated and is used to iterate over the raw logs and unpacked data for ProposalCreated events raised by the GenerativeDao contract.
type GenerativeDaoProposalCreatedIterator struct {
	Event *GenerativeDaoProposalCreated // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoProposalCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoProposalCreated)
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
		it.Event = new(GenerativeDaoProposalCreated)
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
func (it *GenerativeDaoProposalCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoProposalCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoProposalCreated represents a ProposalCreated event raised by the GenerativeDao contract.
type GenerativeDaoProposalCreated struct {
	ProposalId  *big.Int
	Proposer    common.Address
	Targets     []common.Address
	Values      []*big.Int
	Signatures  []string
	Calldatas   [][]byte
	StartBlock  *big.Int
	EndBlock    *big.Int
	Description string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterProposalCreated is a free log retrieval operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 startBlock, uint256 endBlock, string description)
func (_GenerativeDao *GenerativeDaoFilterer) FilterProposalCreated(opts *bind.FilterOpts) (*GenerativeDaoProposalCreatedIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "ProposalCreated")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoProposalCreatedIterator{contract: _GenerativeDao.contract, event: "ProposalCreated", logs: logs, sub: sub}, nil
}

// WatchProposalCreated is a free log subscription operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 startBlock, uint256 endBlock, string description)
func (_GenerativeDao *GenerativeDaoFilterer) WatchProposalCreated(opts *bind.WatchOpts, sink chan<- *GenerativeDaoProposalCreated) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "ProposalCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoProposalCreated)
				if err := _GenerativeDao.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
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

// ParseProposalCreated is a log parse operation binding the contract event 0x7d84a6263ae0d98d3329bd7b46bb4e8d6f98cd35a7adb45c274c8b7fd5ebd5e0.
//
// Solidity: event ProposalCreated(uint256 proposalId, address proposer, address[] targets, uint256[] values, string[] signatures, bytes[] calldatas, uint256 startBlock, uint256 endBlock, string description)
func (_GenerativeDao *GenerativeDaoFilterer) ParseProposalCreated(log types.Log) (*GenerativeDaoProposalCreated, error) {
	event := new(GenerativeDaoProposalCreated)
	if err := _GenerativeDao.contract.UnpackLog(event, "ProposalCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoProposalExecutedIterator is returned from FilterProposalExecuted and is used to iterate over the raw logs and unpacked data for ProposalExecuted events raised by the GenerativeDao contract.
type GenerativeDaoProposalExecutedIterator struct {
	Event *GenerativeDaoProposalExecuted // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoProposalExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoProposalExecuted)
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
		it.Event = new(GenerativeDaoProposalExecuted)
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
func (it *GenerativeDaoProposalExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoProposalExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoProposalExecuted represents a ProposalExecuted event raised by the GenerativeDao contract.
type GenerativeDaoProposalExecuted struct {
	ProposalId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalExecuted is a free log retrieval operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) FilterProposalExecuted(opts *bind.FilterOpts) (*GenerativeDaoProposalExecutedIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "ProposalExecuted")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoProposalExecutedIterator{contract: _GenerativeDao.contract, event: "ProposalExecuted", logs: logs, sub: sub}, nil
}

// WatchProposalExecuted is a free log subscription operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) WatchProposalExecuted(opts *bind.WatchOpts, sink chan<- *GenerativeDaoProposalExecuted) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "ProposalExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoProposalExecuted)
				if err := _GenerativeDao.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
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

// ParseProposalExecuted is a log parse operation binding the contract event 0x712ae1383f79ac853f8d882153778e0260ef8f03b504e2866e0593e04d2b291f.
//
// Solidity: event ProposalExecuted(uint256 proposalId)
func (_GenerativeDao *GenerativeDaoFilterer) ParseProposalExecuted(log types.Log) (*GenerativeDaoProposalExecuted, error) {
	event := new(GenerativeDaoProposalExecuted)
	if err := _GenerativeDao.contract.UnpackLog(event, "ProposalExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoProposalQueuedIterator is returned from FilterProposalQueued and is used to iterate over the raw logs and unpacked data for ProposalQueued events raised by the GenerativeDao contract.
type GenerativeDaoProposalQueuedIterator struct {
	Event *GenerativeDaoProposalQueued // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoProposalQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoProposalQueued)
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
		it.Event = new(GenerativeDaoProposalQueued)
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
func (it *GenerativeDaoProposalQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoProposalQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoProposalQueued represents a ProposalQueued event raised by the GenerativeDao contract.
type GenerativeDaoProposalQueued struct {
	ProposalId *big.Int
	Eta        *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterProposalQueued is a free log retrieval operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 eta)
func (_GenerativeDao *GenerativeDaoFilterer) FilterProposalQueued(opts *bind.FilterOpts) (*GenerativeDaoProposalQueuedIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "ProposalQueued")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoProposalQueuedIterator{contract: _GenerativeDao.contract, event: "ProposalQueued", logs: logs, sub: sub}, nil
}

// WatchProposalQueued is a free log subscription operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 eta)
func (_GenerativeDao *GenerativeDaoFilterer) WatchProposalQueued(opts *bind.WatchOpts, sink chan<- *GenerativeDaoProposalQueued) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "ProposalQueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoProposalQueued)
				if err := _GenerativeDao.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
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

// ParseProposalQueued is a log parse operation binding the contract event 0x9a2e42fd6722813d69113e7d0079d3d940171428df7373df9c7f7617cfda2892.
//
// Solidity: event ProposalQueued(uint256 proposalId, uint256 eta)
func (_GenerativeDao *GenerativeDaoFilterer) ParseProposalQueued(log types.Log) (*GenerativeDaoProposalQueued, error) {
	event := new(GenerativeDaoProposalQueued)
	if err := _GenerativeDao.contract.UnpackLog(event, "ProposalQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoQuorumNumeratorUpdatedIterator is returned from FilterQuorumNumeratorUpdated and is used to iterate over the raw logs and unpacked data for QuorumNumeratorUpdated events raised by the GenerativeDao contract.
type GenerativeDaoQuorumNumeratorUpdatedIterator struct {
	Event *GenerativeDaoQuorumNumeratorUpdated // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoQuorumNumeratorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoQuorumNumeratorUpdated)
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
		it.Event = new(GenerativeDaoQuorumNumeratorUpdated)
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
func (it *GenerativeDaoQuorumNumeratorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoQuorumNumeratorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoQuorumNumeratorUpdated represents a QuorumNumeratorUpdated event raised by the GenerativeDao contract.
type GenerativeDaoQuorumNumeratorUpdated struct {
	OldQuorumNumerator *big.Int
	NewQuorumNumerator *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterQuorumNumeratorUpdated is a free log retrieval operation binding the contract event 0x0553476bf02ef2726e8ce5ced78d63e26e602e4a2257b1f559418e24b4633997.
//
// Solidity: event QuorumNumeratorUpdated(uint256 oldQuorumNumerator, uint256 newQuorumNumerator)
func (_GenerativeDao *GenerativeDaoFilterer) FilterQuorumNumeratorUpdated(opts *bind.FilterOpts) (*GenerativeDaoQuorumNumeratorUpdatedIterator, error) {

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "QuorumNumeratorUpdated")
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoQuorumNumeratorUpdatedIterator{contract: _GenerativeDao.contract, event: "QuorumNumeratorUpdated", logs: logs, sub: sub}, nil
}

// WatchQuorumNumeratorUpdated is a free log subscription operation binding the contract event 0x0553476bf02ef2726e8ce5ced78d63e26e602e4a2257b1f559418e24b4633997.
//
// Solidity: event QuorumNumeratorUpdated(uint256 oldQuorumNumerator, uint256 newQuorumNumerator)
func (_GenerativeDao *GenerativeDaoFilterer) WatchQuorumNumeratorUpdated(opts *bind.WatchOpts, sink chan<- *GenerativeDaoQuorumNumeratorUpdated) (event.Subscription, error) {

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "QuorumNumeratorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoQuorumNumeratorUpdated)
				if err := _GenerativeDao.contract.UnpackLog(event, "QuorumNumeratorUpdated", log); err != nil {
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

// ParseQuorumNumeratorUpdated is a log parse operation binding the contract event 0x0553476bf02ef2726e8ce5ced78d63e26e602e4a2257b1f559418e24b4633997.
//
// Solidity: event QuorumNumeratorUpdated(uint256 oldQuorumNumerator, uint256 newQuorumNumerator)
func (_GenerativeDao *GenerativeDaoFilterer) ParseQuorumNumeratorUpdated(log types.Log) (*GenerativeDaoQuorumNumeratorUpdated, error) {
	event := new(GenerativeDaoQuorumNumeratorUpdated)
	if err := _GenerativeDao.contract.UnpackLog(event, "QuorumNumeratorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoVoteCastIterator is returned from FilterVoteCast and is used to iterate over the raw logs and unpacked data for VoteCast events raised by the GenerativeDao contract.
type GenerativeDaoVoteCastIterator struct {
	Event *GenerativeDaoVoteCast // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoVoteCastIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoVoteCast)
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
		it.Event = new(GenerativeDaoVoteCast)
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
func (it *GenerativeDaoVoteCastIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoVoteCastIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoVoteCast represents a VoteCast event raised by the GenerativeDao contract.
type GenerativeDaoVoteCast struct {
	Voter      common.Address
	ProposalId *big.Int
	Support    uint8
	Weight     *big.Int
	Reason     string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCast is a free log retrieval operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_GenerativeDao *GenerativeDaoFilterer) FilterVoteCast(opts *bind.FilterOpts, voter []common.Address) (*GenerativeDaoVoteCastIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "VoteCast", voterRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoVoteCastIterator{contract: _GenerativeDao.contract, event: "VoteCast", logs: logs, sub: sub}, nil
}

// WatchVoteCast is a free log subscription operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_GenerativeDao *GenerativeDaoFilterer) WatchVoteCast(opts *bind.WatchOpts, sink chan<- *GenerativeDaoVoteCast, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "VoteCast", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoVoteCast)
				if err := _GenerativeDao.contract.UnpackLog(event, "VoteCast", log); err != nil {
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

// ParseVoteCast is a log parse operation binding the contract event 0xb8e138887d0aa13bab447e82de9d5c1777041ecd21ca36ba824ff1e6c07ddda4.
//
// Solidity: event VoteCast(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason)
func (_GenerativeDao *GenerativeDaoFilterer) ParseVoteCast(log types.Log) (*GenerativeDaoVoteCast, error) {
	event := new(GenerativeDaoVoteCast)
	if err := _GenerativeDao.contract.UnpackLog(event, "VoteCast", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeDaoVoteCastWithParamsIterator is returned from FilterVoteCastWithParams and is used to iterate over the raw logs and unpacked data for VoteCastWithParams events raised by the GenerativeDao contract.
type GenerativeDaoVoteCastWithParamsIterator struct {
	Event *GenerativeDaoVoteCastWithParams // Event containing the contract specifics and raw log

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
func (it *GenerativeDaoVoteCastWithParamsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeDaoVoteCastWithParams)
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
		it.Event = new(GenerativeDaoVoteCastWithParams)
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
func (it *GenerativeDaoVoteCastWithParamsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeDaoVoteCastWithParamsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeDaoVoteCastWithParams represents a VoteCastWithParams event raised by the GenerativeDao contract.
type GenerativeDaoVoteCastWithParams struct {
	Voter      common.Address
	ProposalId *big.Int
	Support    uint8
	Weight     *big.Int
	Reason     string
	Params     []byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteCastWithParams is a free log retrieval operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_GenerativeDao *GenerativeDaoFilterer) FilterVoteCastWithParams(opts *bind.FilterOpts, voter []common.Address) (*GenerativeDaoVoteCastWithParamsIterator, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _GenerativeDao.contract.FilterLogs(opts, "VoteCastWithParams", voterRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeDaoVoteCastWithParamsIterator{contract: _GenerativeDao.contract, event: "VoteCastWithParams", logs: logs, sub: sub}, nil
}

// WatchVoteCastWithParams is a free log subscription operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_GenerativeDao *GenerativeDaoFilterer) WatchVoteCastWithParams(opts *bind.WatchOpts, sink chan<- *GenerativeDaoVoteCastWithParams, voter []common.Address) (event.Subscription, error) {

	var voterRule []interface{}
	for _, voterItem := range voter {
		voterRule = append(voterRule, voterItem)
	}

	logs, sub, err := _GenerativeDao.contract.WatchLogs(opts, "VoteCastWithParams", voterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeDaoVoteCastWithParams)
				if err := _GenerativeDao.contract.UnpackLog(event, "VoteCastWithParams", log); err != nil {
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

// ParseVoteCastWithParams is a log parse operation binding the contract event 0xe2babfbac5889a709b63bb7f598b324e08bc5a4fb9ec647fb3cbc9ec07eb8712.
//
// Solidity: event VoteCastWithParams(address indexed voter, uint256 proposalId, uint8 support, uint256 weight, string reason, bytes params)
func (_GenerativeDao *GenerativeDaoFilterer) ParseVoteCastWithParams(log types.Log) (*GenerativeDaoVoteCastWithParams, error) {
	event := new(GenerativeDaoVoteCastWithParams)
	if err := _GenerativeDao.contract.UnpackLog(event, "VoteCastWithParams", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
