package contractaddress

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/laziercoder/go-web3"
	"github.com/laziercoder/go-web3/eth"
	"golang.org/x/crypto/sha3"
	"rederinghub.io/pkg/third-party/crypto/constants/abi"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptotransactionstatus"
)

type ethereumContractImpl struct {
	contract          *eth.Contract
	ethClient         *ethclient.Client
	ethNativeStrategy Strategy
	coin              string
	contractAddress   common.Address
}

func (e *ethereumContractImpl) GetNativeStrategy() Strategy {
	return newEthereumChainByClient(e.ethClient)
}

func (e ethereumContractImpl) GenerateContractAddress() (*ContractAddress, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hexutil.Encode(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		// log.Fatal("error casting public key to ECDSA")
		return nil, errors.New("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// check address is valid or not
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(address) {
		return nil, fmt.Errorf("invalid address")
	}

	cAddress := NewContractAddress().WithAddress(address).WithPrivateKey(privateKeyHex)
	return cAddress, nil
}

func (e ethereumContractImpl) CheckBalance(address string) (float64, error) {
	ownerAddress := common.HexToAddress(address)
	b, err := e.contract.Call("balanceOf", ownerAddress)
	if err != nil {
		return 0, err
	}
	balance, ok := b.(*big.Int)
	if !ok {
		return 0, errors.New("invalid type")
	}

	balanceFloat64 := balanceToFloat64(balance, e.coin)

	return balanceFloat64, nil
}

func (e ethereumContractImpl) Transfer(req *TransferRequest) (string, error) {
	fromPrivateKey := strings.TrimPrefix(req.FromAddress.GetPrivateKey(), "0x") // remove 0x prefix

	privateKey, err := crypto.HexToECDSA(fromPrivateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fromAddress.Hex()

	// ref: https://ethereum.stackexchange.com/questions/27256/error-replacement-transaction-underpriced
	nonce, err := e.ethClient.NonceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", err
	}

	if req.GasPrice == nil {
		gasPrice, err := e.SuggestGasPrice(context.Background())
		if err != nil {
			return "", err
		}
		req.GasPrice = gasPrice
	}

	value := big.NewInt(0) // in wei (0 eth)
	if req.GasLimit == 0 {
		gasLimit, err := e.EstimateGas(req)
		if err != nil {
			return "", err
		}
		req.GasLimit = gasLimit
	}

	ethBalance, err := e.ethClient.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", err
	}

	minETHBalance := req.GasLimit * req.GasPrice.Uint64()
	if ethBalance.Uint64() < minETHBalance {
		return "", errors.New("insufficient balance")
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      req.GasLimit,
		GasPrice: req.GasPrice,
		To:       &e.contractAddress,
		Value:    value,
		Data:     req.Data,
	})

	chainID, err := e.ethClient.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = e.ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

func (e ethereumContractImpl) SuggestGasPrice(context context.Context) (*big.Int, error) {
	return e.ethClient.SuggestGasPrice(context)
}

// EstimateGas estimates the gas cost of a transaction.
// Ref: https://goethereumbook.org/transfer-tokens/
func (e ethereumContractImpl) EstimateGas(req *TransferRequest) (uint64, error) {
	toAddress := common.HexToAddress(req.ToAddress)
	fromAddress := common.HexToAddress(req.FromAddress.GetWalletAddress())

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	amount := big.NewInt(int64(req.Amount))
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Data:     data,
		GasPrice: req.GasPrice,
		Gas:      100000,
	}

	gasLimit, err := e.ethClient.EstimateGas(context.Background(), msg) // gas limit is estimated
	if err != nil {
		return 0, err
	}

	req.Data = data
	return gasLimit, nil
}

func (e ethereumContractImpl) GetTransactionReceipt(ctx context.Context, transactionID string) (*TransactionReceipt, error) {
	hash := common.HexToHash(transactionID)
	ethReceipt, err := e.ethClient.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	rawData, _ := ethReceipt.MarshalJSON()
	receipt := &TransactionReceipt{
		Status:        cryptotransactionstatus.Success,
		RawData:       string(rawData),
		Coin:          e.coin,
		TransactionID: transactionID,
	}

	if ethReceipt.Status == types.ReceiptStatusFailed {
		receipt.Status = cryptotransactionstatus.Failure
	}

	return receipt, nil
}

func NewEthereumContractClient(cfg *TokenClientConfig) (Strategy, error) {
	web3Client, err := web3.NewWeb3(cfg.ProviderURL)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethclient.Dial(cfg.ProviderURL)
	if err != nil {
		return nil, err
	}

	contract, err := web3Client.Eth.NewContract(abi.EthereumERC20, cfg.ContractAddress)
	if err != nil {
		return nil, err
	}

	return &ethereumContractImpl{
		contract:        contract,
		coin:            cfg.Coin,
		ethClient:       ethClient,
		contractAddress: common.HexToAddress(cfg.ContractAddress),
	}, nil
}
