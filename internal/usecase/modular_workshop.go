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
	if len(entity.Thumbnail) == 0 {
		return errors.New(fmt.Sprintf("No thumbnail found in the metadata."))
	}
	if len(entity.OwnerAddr) == 0 {
		return errors.New(fmt.Sprintf("No Owner Addr found in the metadata."))
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
			inscriptionOwnerAddress := ""
			tokenURI, err := u.GetTokenByTokenID(id, 0)
			if err == nil && tokenURI != nil {
				inscriptionOwnerAddress = tokenURI.OwnerAddr
			} else {
				info, err := u.GetInscriptionByIDFromOrd(id)
				if err == nil && info != nil {
					inscriptionOwnerAddress = info.Address
				}
			}
			if len(inscriptionOwnerAddress) == 0 {
				return errors.New(fmt.Sprintf("Can not check owner of inscription id %s due to Oridinal & DB . Try again", id))
			}
			if entity.OwnerAddr != inscriptionOwnerAddress {
				return errors.New(fmt.Sprintf("Error occurred when checking the owner of inscription id %s from Original, sender: %s, owner: %s.", id, entity.OwnerAddr, inscriptionOwnerAddress))
			}
		}
	}
	return nil
}
