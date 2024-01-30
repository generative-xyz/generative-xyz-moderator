package usecase

import (
	"errors"
	"fmt"
	"rederinghub.io/internal/entity"
)

func (u Usecase) ValidateModularWorkshopEntity(entity *entity.ModularWorkshopEntity, checkEmpty bool, checkOwner bool) error {
	inscriptionIds := entity.GetListInscriptionIds()
	mapInscriptionIds := map[string]bool{}
	if len(inscriptionIds) == 0 && checkEmpty {
		return errors.New(fmt.Sprintf("No inscription found in the metadata."))
	}
	for i, id := range inscriptionIds {
		if len(id) == 0 {
			return errors.New(fmt.Sprintf("Empty inscription id at index %d", i))
		}
		if mapInscriptionIds[id] {
			return errors.New(fmt.Sprintf("Duplicate inscription id %s at index %d", id, i))
		}
		mapInscriptionIds[id] = true
		if checkOwner {
			info, err := u.GetInscriptionByIDFromOrd(id)
			if err != nil {
				return errors.New(fmt.Sprintf("Cannot check owner of inscription id %s due to Oridinal error: %s", id, err.Error()))
			}
			if info == nil {
				return errors.New(fmt.Sprintf("Can not check owner of inscription id %s due to Oridinal info == nil", id))
			}
			if entity.OwnerAddr != info.Address {
				return errors.New(fmt.Sprintf("Error occurred when checking the owner of inscription id %s from Original, sender: %s, owner: %s.", id, entity.OwnerAddr, info.Address))
			}
		}
	}
	return nil
}
