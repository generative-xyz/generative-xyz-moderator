package repository

import "rederinghub.io/internal/entity"

func (r Repository) CreateDexBTCLog(log *entity.DexBTCLog) error {
	err := r.InsertOne(log.TableName(), log)
	if err != nil {
		return err
	}
	return nil
}
