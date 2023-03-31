package helpers

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"rederinghub.io/utils/logger"
)

func EthDialer() (*ethclient.Client, error) {
	chainURL := os.Getenv("CHAIN_URL")
	return ChainDialer(chainURL)
}

func TCDialer() (*ethclient.Client, error) {
	chainURL := os.Getenv("TC_ENDPOINT")
	return ChainDialer(chainURL)
}

func ChainDialer(chainURL string) (*ethclient.Client, error) {
	logger.AtLog.Logger.Info("ChainDialer",zap.String("chainURL", chainURL))
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}