package usecase

import (
	"fmt"

	"go.uber.org/zap"
	"rederinghub.io/utils/logger"
)



const (
	API_URL string = `https://dev-v5.generativeexplorer.com/api/`
)

func (u Usecase) FindInscriptions() {
	totalItem := 100
	max := 200
	inscriptions := []string{}
	for {
		if totalItem > max {
			break
		}

		data, err := u.Dev5service.Inscriptions(fmt.Sprintf("%d",totalItem))
		if err != nil {
			logger.AtLog.Logger.Error("u.Dev5service.Inscriptions", zap.Error(err))
		}
		logger.AtLog.Logger.Info("u.Dev5service.Inscriptions", zap.Any("next", data.Next))
		totalItem = data.Next
		//inscriptions = append(inscriptions, data.Inscriptions)
	}
	
}