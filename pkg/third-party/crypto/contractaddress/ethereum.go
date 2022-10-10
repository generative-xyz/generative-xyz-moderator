package contractaddress

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptotransactionstatus"

	"math/big"
	"regexp"
	"strings"
)

type ethereumImpl struct {
	client *ethclient.Client
}

func (e *ethereumImpl) GetNativeStrategy() Strategy {
	return e
}

func (e ethereumImpl) GetTransactionReceipt(ctx context.Context, transactionID string) (*TransactionReceipt, error) {
	hash := common.HexToHash(transactionID)
	ethReceipt, err := e.client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, err
	}

	rawData, _ := ethReceipt.MarshalJSON()
	receipt := &TransactionReceipt{
		Status:        cryptotransactionstatus.Success,
		RawData:       string(rawData),
		Coin:          cryptocurrency.Ethereum,
		TransactionID: transactionID,
	}

	if ethReceipt.Status == types.ReceiptStatusFailed {
		receipt.Status = cryptotransactionstatus.Failure
	}

	return receipt, nil
}

// EstimateGas estimates the gas needed to transfer the funds from one address to another
func (e ethereumImpl) EstimateGas(req *TransferRequest) (uint64, error) {
	fromAddress := common.HexToAddress(req.FromAddress.GetWalletAddress())
	toAddress := common.HexToAddress(req.ToAddress)
	msg := ethereum.CallMsg{
		From:     fromAddress,
		To:       &toAddress,
		Value:    big.NewInt(int64(req.Amount)),
		GasPrice: req.GasPrice,
	}
	gasLimit, err := e.client.EstimateGas(context.Background(), msg) // gas limit is estimated
	if err != nil {
		return 0, err
	}
	return gasLimit, nil
}

func (e ethereumImpl) SuggestGasPrice(context context.Context) (*big.Int, error) {
	return e.client.SuggestGasPrice(context)
}

func (e ethereumImpl) CheckBalance(address string) (float64, error) {
	account := common.HexToAddress(address)
	balance, err := e.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return 0, err
	}

	value := float64(balance.Uint64()) / math.Pow10(cryptocurrency.EthereumRoundPlaces)
	return value, nil
}

// GenerateContractAddress generates a new contract address
func (e ethereumImpl) GenerateContractAddress() (*ContractAddress, error) {
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

// Transfer transfers the funds from one address to another
// Test receiver account:
//   address: 0xb0cda09aBcc2DA7760AE3862a9204401721c9bB1
//   pri: 0x695b1510bc4088fe3c34fda3ef439a8ec8473c2702462b45ae743829ce95e38d
// Test sender account:
//   address: 0x976D5565927cF44Ee19c346F61FCB37238B426D1
//   pri: 0x319717766d44c592d9971d8c595c95caea59017755933deb850b37f0b7941dc4
func (e ethereumImpl) Transfer(req *TransferRequest) (string, error) {
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
	nonce, err := e.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(int64(req.Amount))
	toAddressETH := common.HexToAddress(req.ToAddress)

	if req.GasLimit == 0 {
		gasLimit, err := e.EstimateGas(req)
		if err != nil {
			return "", err
		}
		req.GasLimit = gasLimit
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		Gas:      req.GasLimit,
		GasPrice: req.GasPrice,
		To:       &toAddressETH,
		Value:    value,
	})

	chainID, err := e.client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return "", err
	}

	err = e.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// NewEthereumChain returns a new ethereumImpl instance
// testnet goerli:= https://goerli.infura.io/v3/1db8e9d0f92d4c6d88fa7d15da5dcefb
func NewEthereumChain(clientProxy string) (Strategy, error) {
	client, err := ethclient.Dial(clientProxy)
	if err != nil {
		return nil, err
	}
	return &ethereumImpl{
		client: client,
	}, nil
}

func newEthereumChainByClient(c *ethclient.Client) Strategy {
	return &ethereumImpl{
		client: c,
	}
}
