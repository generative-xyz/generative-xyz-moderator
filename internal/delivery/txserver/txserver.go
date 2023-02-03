package txserver

import (
	"rederinghub.io/internal/txconsumer"
	"rederinghub.io/internal/usecase"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/global"
)

type TxServer struct {
	Txconsumer txconsumer.HttpTxConsumer
}

func NewTxServer(global *global.Global,uc usecase.Usecase, cfg config.Config) (*TxServer, error) {
	t := &TxServer{}
	txConsumer, err := txconsumer.NewHttpTxConsumer(global, uc, cfg)
	if err != nil {
		return  nil, err
	}
	t.Txconsumer = *txConsumer
	return t, nil
}

func (tx TxServer) StartServer() {
	tx.Txconsumer.StartServer()
}