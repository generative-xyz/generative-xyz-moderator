package contractaddress

import (
	"context"
	"math/big"

	"rederinghub.io/pkg/third-party/googlesecret"
)

type Client struct {
	coinStrategy         map[string]Strategy
	googleSecretClient   *googlesecret.Client
	publicKeySecretName  string
	privateKeySecretName string
}

// GenerateTokenAddress generate token address and return public key
func (c *Client) GenerateTokenAddress(coin string) (*ContractAddress, error) {
	cAddress, err := c.coinStrategy[coin].GenerateContractAddress()
	if err != nil {
		return nil, err
	}

	publicKey, err := c.googleSecretClient.GetSecret(c.publicKeySecretName)
	if err != nil {
		return nil, err
	}
	err = cAddress.EncryptPrivateKey(publicKey)
	if err != nil {
		return nil, err
	}

	return cAddress, nil
}

// RevokeTokenAddress revoke token address and return private key
func (c *Client) RevokeTokenAddress(address *ContractAddress) (*ContractAddress, error) {
	privateKey, err := c.googleSecretClient.GetSecret(c.privateKeySecretName)
	if err != nil {
		return nil, err
	}
	err = address.DecryptPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return address, nil
}

// CheckBalance check balance of token address
func (c *Client) CheckBalance(coin, address string, useNativeStrategy bool) (float64, error) {
	return c.getStrategyByCoin(coin, useNativeStrategy).
		CheckBalance(address)
}

// Transfer this func transfer crypto from one address to another
func (c *Client) Transfer(coin string, req *TransferRequest) (string, error) {
	return c.getStrategyByCoin(coin, req.UseNativeStrategy).
		Transfer(req)
}

// SuggestGasPrice this func suggest gas price for transfer
func (c *Client) SuggestGasPrice(coin string, ctx context.Context, useNativeStrategy bool) (*big.Int, error) {
	return c.getStrategyByCoin(coin, useNativeStrategy).SuggestGasPrice(ctx)
}

// EstimateGas this func estimate gas for transfer
func (c *Client) EstimateGas(coin string, req *TransferRequest) (uint64, error) {
	return c.getStrategyByCoin(coin, req.UseNativeStrategy).EstimateGas(req)
}

// GetTransactionReceipt this func get transaction receipt
func (c *Client) GetTransactionReceipt(
	coin string,
	ctx context.Context,
	transactionID string,
	useNativeStrategy bool,
) (*TransactionReceipt, error) {
	return c.getStrategyByCoin(coin, useNativeStrategy).GetTransactionReceipt(ctx, transactionID)
}

func (c *Client) getStrategyByCoin(coin string, useNativeStrategy bool) Strategy {
	if useNativeStrategy {
		return c.coinStrategy[coin].GetNativeStrategy()
	}
	return c.coinStrategy[coin]
}

// RegisterCoins register coins for client by str
func (c *Client) RegisterCoins(s Strategy, coins ...string) {
	for _, coin := range coins {
		c.coinStrategy[coin] = s
	}
}

func (c *Client) WithGoogleSecretClient(client *googlesecret.Client) *Client {
	c.googleSecretClient = client
	return c
}

func (c *Client) WithPublicKeySecretName(name string) *Client {
	c.publicKeySecretName = name
	return c
}

func (c *Client) WithPrivateKeySecretName(name string) *Client {
	c.privateKeySecretName = name
	return c
}

func NewClient() *Client {
	result := &Client{
		coinStrategy: map[string]Strategy{},
	}

	return result
}
