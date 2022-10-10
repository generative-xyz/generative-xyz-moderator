package web3

import (
	"go.uber.org/zap"
	"rederinghub.io/pkg/logger"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
)

type Client struct {
	autoREXNFTStrategy Strategy
	ethereumStrategy   Strategy
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) WithERC1155Strategy(providerURL, contractAddress string) *Client {
	c.autoREXNFTStrategy = NewClientERC1155(providerURL, contractAddress)
	return c
}

func (c *Client) WithEthereumStrategy(baseURL, apiKey string) *Client {
	c.ethereumStrategy = NewClientEthereum(baseURL, apiKey)
	return c
}

func (c *Client) BalanceOf(req *BalanceRequest) (*nftdata.NFTCustomerInfo, error) {
	needCheckFromEthereumStrategy := req.ChainID == 0
	if !needCheckFromEthereumStrategy {
		info, err := c.autoREXNFTStrategy.BalanceOf(req)
		if err != nil {
			logger.AtLog.Logger.Error("autoREXNFTStrategy.BalanceOf",
				zap.Error(err),
				zap.Any("rawData", req),
			)
			needCheckFromEthereumStrategy = true
		}

		if info.Balance == 0 {
			needCheckFromEthereumStrategy = true
		}

		if !needCheckFromEthereumStrategy {
			return info, nil
		}
	}

	return c.ethereumStrategy.BalanceOf(req)
}

func (c *Client) NFTFeeds(req *BalanceRequest) (interface{}, error) {
	return c.ethereumStrategy.NFTFeeds(req)
}
