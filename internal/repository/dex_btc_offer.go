package repository

import "rederinghub.io/internal/entity"

func (r Repository) CreateDexBTCOffer(listing *entity.DexBTCOffer) error {
	err := r.InsertOne(listing.TableName(), listing)
	if err != nil {
		return err
	}
	return nil
}
