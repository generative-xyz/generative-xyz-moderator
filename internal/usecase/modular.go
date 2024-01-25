package usecase

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
	"strings"
)

type InsOwner struct {
	InscriptionID string
	OwnerAddress  string
	IsUpdated     bool
	Err           error
}

type Ins struct {
	InscriptionID string
	OwnerAddress  string
}

func (u Usecase) ListModulars(ctx context.Context, f structure.FilterTokens) (*entity.Pagination, error) {
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")
	f.GenNFTAddr = &genNFTAddr
	inscriptions, err := u.Repo.AggregateListModularInscriptions(ctx, f)
	if err != nil {
		return nil, err
	}

	return inscriptions, nil
}

func (u Usecase) CrontabUpdateModularInscOwners(ctx context.Context) error {

	page := 1
	limit := 100
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")

	for {
		offset := (page - 1) * limit
		inscriptions, err := u.Repo.AggregateModularInscriptions(ctx, genNFTAddr, offset, limit)
		if err != nil {
			return err
		}

		if len(inscriptions) == 0 {
			break
		}

		inChan := make(chan Ins, len(inscriptions))
		outChan := make(chan InsOwner, len(inscriptions))

		for range inscriptions {
			go u.FindModularInscOwner(inChan, outChan)
		}

		for _, i := range inscriptions {
			inChan <- Ins{
				InscriptionID: i.TokenID,
				OwnerAddress:  i.OwnerAddr,
			}
		}

		for range inscriptions {
			outFChan := <-outChan
			if outFChan.Err != nil {
				continue
			}

			if outFChan.IsUpdated {
				//TODO - update owner
				fmt.Println(fmt.Sprintf("[ins] %s-%s-%v", outFChan.InscriptionID, outFChan.OwnerAddress, outFChan.IsUpdated))
				_, err := u.UpdateModularInscOwner(outFChan.InscriptionID, outFChan.OwnerAddress)
				if err != nil {
					logger.AtLog.Logger.Error("CrontabUpdateModularInscOwners", zap.Error(err), zap.String("token_id", outFChan.InscriptionID), zap.String("owner_address", outFChan.OwnerAddress))
				}
			}
		}

		page++
	}

	return nil
}

func (u Usecase) FindModularInscOwner(in chan Ins, out chan InsOwner) {
	var err error
	addr := ""
	inscID := <-in
	info := &structure.InscriptionOrdInfoByID{}
	isUpdate := false

	defer func() {
		isUpdate = !strings.EqualFold(inscID.OwnerAddress, addr)
		out <- InsOwner{
			Err:           err,
			InscriptionID: inscID.InscriptionID,
			OwnerAddress:  addr,
			IsUpdated:     isUpdate,
		}
	}()

	info, err = u.GetInscriptionByIDFromOrd(inscID.InscriptionID)
	if err != nil {
		return
	}

	addr = info.Address
}

func (u Usecase) UpdateModularInscOwner(insID string, ownerAddress string) (*mongo.UpdateResult, error) {
	f := bson.D{
		{"token_id", insID},
		{"project_id", os.Getenv("MODULAR_PROJECT_ID")},
	}

	uupdate := bson.D{
		{"owner_addrress", ownerAddress},
	}

	update := bson.D{{"$set", uupdate}}

	//prevent update from local
	if os.Getenv("ENV") != "mainnet" {
		return nil, nil
	}

	result, err := u.Repo.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return nil, err
	}

	return result, nil

}
