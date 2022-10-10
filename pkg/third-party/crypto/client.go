package crypto

import (
	"context"

	"go.uber.org/zap"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/third-party/crypto/constants/cryptocurrency"
	"rederinghub.io/pkg/third-party/crypto/contractaddress"
	"rederinghub.io/pkg/third-party/crypto/exchange"
	"rederinghub.io/pkg/third-party/googlesecret"

	"math/big"
)

type Client struct {
	// exchangeClient Client
	exchangeClient        exchange.Client
	contractAddressClient *contractaddress.Client
}

// GetMetadata returns metadata for a given cryptocurrency
func (c *Client) GetMetadata(currency string) *Metadata {
	rate, err := c.exchangeClient.GetRate(currency)
	if err != nil {
		logger.AtLog.Logger.Error(
			"[crypto] GetMetadata fail to get rate",
			zap.Error(err),
			zap.String("coin", currency),
		)
		return nil
	}

	currencySymbol, hasSymbol := cryptocurrency.CurrencySymbolByCurrency[currency]
	if !hasSymbol {
		logger.AtLog.Logger.Error(
			"[crypto] GetMetadata fail to get currency symbol",
			zap.Error(exchange.ErrCurrencyNotSupported),
			zap.String("coin", currency),
		)
		return nil
	}
	result := NewMetadata().WithCurrency(currency).WithRate(rate).WithCurrencySymbol(currencySymbol)
	return result
}

// GenerateTokenAddress generate token address and return public key
func (c *Client) GenerateTokenAddress(coin string) (*contractaddress.ContractAddress, error) {
	cAddress, err := c.contractAddressClient.GenerateTokenAddress(coin)
	if err != nil {
		return nil, err
	}
	return cAddress, nil
}

// RevokeTokenAddress revoke token address and return private key
func (c *Client) RevokeTokenAddress(address *contractaddress.ContractAddress) (*contractaddress.ContractAddress, error) {
	return c.contractAddressClient.RevokeTokenAddress(address)
}

// Transfer this func transfer crypto from one address to another
func (c *Client) Transfer(coin string, req *contractaddress.TransferRequest) (string, error) {
	revokedPayer, err := c.RevokeTokenAddress(req.Payer)
	if err != nil {
		return "", err
	}
	req.Payer = revokedPayer

	revokedFromAddress, err := c.RevokeTokenAddress(req.FromAddress)
	if err != nil {
		return "", err
	}
	req.FromAddress = revokedFromAddress

	return c.contractAddressClient.Transfer(
		coin,
		req,
	)
}

// CheckAddressBalance check balance of token address
func (c *Client) CheckAddressBalance(coin, address string, useNativeStrategy bool) (float64, error) {
	return c.contractAddressClient.CheckBalance(coin, address, useNativeStrategy)
}

// SuggestGasPrice returns suggested gas price
func (c *Client) SuggestGasPrice(coin string, ctx context.Context, useNativeStrategy bool) (*big.Int, error) {
	return c.contractAddressClient.SuggestGasPrice(coin, ctx, useNativeStrategy)
}

// EstimateGas estimate gas for transfer
func (c *Client) EstimateGas(coin string, req *contractaddress.TransferRequest) (uint64, error) {
	return c.contractAddressClient.EstimateGas(coin, req)
}

// GetTransactionReceipt returns transaction receipt
func (c *Client) GetTransactionReceipt(coin string, ctx context.Context, transactionID string, useNativeStrategy bool) (*contractaddress.TransactionReceipt, error) {
	return c.contractAddressClient.GetTransactionReceipt(coin, ctx, transactionID, useNativeStrategy)
}

// RegisterCoins register coin with strategy
func (c *Client) RegisterCoins(s contractaddress.Strategy, coins ...string) {
	c.contractAddressClient.RegisterCoins(s, coins...)
}

// WithPublicKeySecretName set public key secret name
func (c *Client) WithPublicKeySecretName(name string) *Client {
	c.contractAddressClient.WithPublicKeySecretName(name)
	return c
}

func (c *Client) WithGoogleSecretClient(client *googlesecret.Client) *Client {
	c.contractAddressClient.WithGoogleSecretClient(client)
	return c
}

// WithPrivateKeySecretName set private key secret name
func (c *Client) WithPrivateKeySecretName(name string) *Client {
	c.contractAddressClient.WithPrivateKeySecretName(name)
	return c
}

func NewClient() *Client {
	return &Client{
		exchangeClient:        exchange.NewCoingeckoClient(),
		contractAddressClient: contractaddress.NewClient(),
	}
}
