package helpers

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func EthDialer() (*ethclient.Client, error) {
	chainURL := os.Getenv("CHAIN_URL")
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}