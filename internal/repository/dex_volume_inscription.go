package repository

import (
	"errors"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (r Repository) FindDexVolumeInscription(filter *structure.DexVolumeInscritionFilter) ([]entity.DexVolumeInscription, error) {
	return nil, nil
}

func (r Repository) InsertDexVolumeInscription(data *entity.DexVolumeInscription) error {
	if data == nil {
		return errors.New("insertDexVolumeInscription Invalid data")
	}
	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}
	return nil
}