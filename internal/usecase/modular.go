package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/url"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"strconv"
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

type InscriptionResp struct {
	Err    error         `json:"error"`
	Status bool          `json:"status"`
	Data   []Inscription `json:"data"`
}

type Inscription struct {
	InscID      string
	BlockHeight uint64
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

// Crontab update owner of modular inscriptions
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

// Crontab add modular inscriptions
func (u Usecase) CrontabAddModularInscs(ctx context.Context) error {
	fBlockKey := "from_ord_block"
	toBlockKey := "to_ord_block"
	processedBlockKey := "processed_ord_block"

	fBlock := uint64(0)
	toBlock := uint64(0)
	proccessedBlock := uint64(0)

	errCached := u.Cache.GetObjectData(fBlockKey, &fBlock)
	if errCached != nil {
		fInt, _ := strconv.Atoi(os.Getenv("MODULAR_FROM_BLOCK"))
		fBlock = uint64(fInt)
	}

	errCached2 := u.Cache.GetObjectData(toBlockKey, &toBlock)
	if errCached2 != nil {
		tInt, _ := strconv.Atoi(os.Getenv("MODULAR_TO_BLOCK"))
		toBlock = uint64(tInt)
	}

	logKey := "CrontabAddModularInscs"
	var err error
	logP := new([]zap.Field)
	logs := []zap.Field{}
	logP = &logs

	defer func() {
		if err != nil {
			logs = append(logs, zap.Error(err))
			logger.AtLog.Logger.Error(logKey, *logP...)
		} else {
			logger.AtLog.Logger.Info(logKey, *logP...)
		}
	}()

	u.Cache.GetObjectData(processedBlockKey, &proccessedBlock)
	logs = append(logs, zap.Uint64("from_block", fBlock))
	logs = append(logs, zap.Uint64("to_block", toBlock))
	logs = append(logs, zap.Uint64("processed_block", proccessedBlock))

	quickNode, err := btc.GetBlockCountfromQuickNode(u.Config.QuicknodeAPI)
	if err != nil {
		return err
	}

	fBlock = toBlock + 1
	toBlock += 1000

	if toBlock > quickNode.Result {
		toBlock = quickNode.Result
	}

	if fBlock > quickNode.Result {
		fBlock = quickNode.Result
	}

	if proccessedBlock == toBlock {
		logs = append(logs, zap.String("message", "processed"))
		return nil
	}

	queryParams := url.Values{}
	queryParams.Set("fromBlock", fmt.Sprintf("%d", fBlock))
	queryParams.Set("toBlock", fmt.Sprintf("%d", toBlock))

	_url := fmt.Sprintf("%s/bvm-insc/list", os.Getenv("MODULAR_BRIDGES_API"))
	_url += "?" + queryParams.Encode()
	logs = append(logs, zap.String("url", _url))
	_b, _, _, err := helpers.HttpRequest(_url, "GET", map[string]string{}, nil)
	if err != nil {
		return err
	}
	logs = append(logs, zap.Any("resp", _b))

	resp := InscriptionResp{}
	err = json.Unmarshal(_b, &resp)
	if err != nil {
		return err
	}

	logs = append(logs, zap.Int("total", len(resp.Data)))
	for i, item := range resp.Data {
		modulerObj := &entity.ModularInscription{
			InscriptionID:  item.InscID,
			BlockHeight:    item.BlockHeight,
			IsCreatedToken: false, // created token will be handled by the other crontab
		}

		logs = append(logs, zap.String(fmt.Sprintf("InscID.%d", i), item.InscID))

		//avoid duplicated by unique-index
		inserted, err1 := u.Repo.InsertModular(modulerObj)
		if err1 != nil {
			logs = append(logs, zap.String(fmt.Sprintf("%s.inserted", item.InscID), err1.Error()))
			continue
		}

		_ = inserted
	}

	u.Cache.SetData(processedBlockKey, toBlock)

	//set data
	u.Cache.SetData(fBlockKey, fBlock)
	u.Cache.SetData(toBlockKey, toBlock)
	return nil
}
