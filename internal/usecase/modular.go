package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"math"
	"net/url"
	"os"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"sort"
	"strconv"
	"strings"
	"time"
)

type CreatedTokenChan struct {
	Token *entity.TokenUri
	Err   error
}

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

func (u *Usecase) CreateModularTraits(attrs []entity.TokenUriAttrStr) string {
	t := ""

	sort.SliceStable(attrs, func(i, j int) bool {
		return attrs[i].TraitType > attrs[j].TraitType
	})

	for _, attr := range attrs {
		if !strings.EqualFold(attr.TraitType, "hash") {
			t += strings.ToLower(strings.ReplaceAll(fmt.Sprintf("%s.%s", attr.TraitType, attr.Value), " ", ""))
		}
	}

	return t
}

func (u Usecase) GroupListModulars(ctx context.Context, f structure.FilterTokens) (*entity.Pagination, error) {
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")
	checkOwner, _ := strconv.ParseBool(os.Getenv("MODULAR_WORKSHOP_CHECK_OWNER"))
	if checkOwner {
		f.GenNFTAddr = &genNFTAddr
	} else {
		f.GenNFTAddr = nil
	}
	inscriptions, count, err := u.Repo.GroupModularInscByAttr(ctx, f)
	if err != nil {
		return nil, err
	}

	inscriptionIDs := []string{}
	groupInscriptionIDs := make(map[string][]string)
	for _, i := range inscriptions {
		inscriptionIDs = append(inscriptionIDs, i.Attr[0].InscriptionID)

		ins := []string{}
		for _, i1 := range i.Attr {
			ins = append(ins, i1.InscriptionID)
		}

		groupInscriptionIDs[i.GroupID] = ins
	}

	for _, i := range inscriptions {
		inscriptionIDs = append(inscriptionIDs, i.Attr[0].InscriptionID)
	}

	allTokens, err := u.Repo.AggregateListModularInscriptionsByTokenIDs(ctx, inscriptionIDs)
	if err != nil {
		return nil, err
	}

	tokens := make(map[string]entity.ModularTokenUri)
	for _, i := range allTokens {
		tokens[i.TokenID] = i
	}

	for _, i := range inscriptions {
		t, ok := tokens[i.Attr[0].InscriptionID]
		if ok {
			gPi := entity.GroupedModularTokenUri{
				Image:               t.Image,
				AnimationURL:        t.AnimationURL,
				ParsedAttributesStr: t.ParsedAttributesStr,
				Thumbnail:           t.Thumbnail,
				//accept duplicated data to query more faster
				Project: t.Project,
			}

			i.GroupedModularTokenUri = gPi
		}

		t1, ok := groupInscriptionIDs[i.GroupID]
		if ok {
			i.Items = t1
		}
	}

	return &entity.Pagination{
		Result:    inscriptions,
		Page:      f.Page,
		PageSize:  f.Limit,
		Total:     int64(count),
		TotalPage: int64(math.Ceil(float64(count) / float64(f.Limit))),
	}, nil
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

// 3. Crontab update owner of modular inscriptions
func (u Usecase) CrontabUpdateModularInscOwners() error {
	ctx := context.Background()
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

// 4. Crontab update traits of modular inscriptions
func (u Usecase) CrontabUpdateModularInscTraits() error {
	f := structure.FilterTokens{}
	genNFTAddr := os.Getenv("MODULAR_PROJECT_ID")
	f.GenNFTAddr = &genNFTAddr
	inscriptions, err := u.Repo.AllModularInscriptions(context.Background(), f)
	if err != nil {
		return err
	}

	for _, token := range inscriptions {
		t := u.CreateModularTraits(token.ParsedAttributesStr)
		fmt.Println(fmt.Sprintf("[modular traits] - %s %s", token.TokenID, t))

		i, err := u.Repo.UpsertModularAttribute(token.TokenID, t)
		if err != nil {
			continue
		}

		spew.Dump(i)
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
	//if os.Getenv("ENV") != "mainnet" {
	//	return nil, nil
	//}

	result, err := u.Repo.DB.Collection(utils.COLLECTION_TOKEN_URI).UpdateOne(context.TODO(), f, update)
	if err != nil {
		return nil, err
	}

	return result, nil

}

// 1. Crontab add modular inscriptions
func (u Usecase) CrontabAddModularInscs() error {
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
	//logs = append(logs, zap.Any("resp", string(_b)))

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

// 2. Crontab Create token from modular inscriptions
func (u Usecase) CrontabCreateTokenFromInscriptions() error {
	ctx := context.Background()
	page := 1
	limit := 5

	for {

		offset := (page - 1) * limit
		inscriptions, err := u.Repo.UnCreatedModularInscriptions(ctx, offset, limit)
		if err != nil {
			return err
		}

		if len(inscriptions) == 0 {
			break
		}

		for _, i := range inscriptions {
			token, err := u.CreateTokenFromInscription(i)
			if err != nil {
				continue
			}

			u.Repo.SetCreatedTokenStatus(token.TokenID, true)
		}

		page++
	}

	return nil
}

func (u Usecase) CreateTokenFromInscription(input entity.ModularInscription) (*entity.TokenUri, error) {

	var err error
	token := &entity.TokenUri{}

	meta := entity.CollectionMeta{
		ProjectID: os.Getenv("MODULAR_PROJECT_ID"),
	}

	inscription := entity.CollectionInscription{
		ID: input.InscriptionID,
	}

	token, err = u.CreateBTCTokenURIFromModularCollectionInscription(meta, inscription, input.BlockHeight)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (u Usecase) CreateBTCTokenURIFromModularCollectionInscription(meta entity.CollectionMeta, inscription entity.CollectionInscription, blockHeight uint64) (*entity.TokenUri, error) {
	// find project by projectID
	project, err := u.Repo.FindProjectByTokenID(meta.ProjectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.AtLog.Logger.Error("CanNotFindProjectByInscriptionIcon", zap.Any("inscriptionIcon", meta.InscriptionIcon))
			return nil, repository.ErrNoProjectsFound
		} else {
			return nil, err
		}
	}

	tokenUri := entity.TokenUri{}
	tokenUri.ContractAddress = os.Getenv("GENERATIVE_BTC_PROJECT")
	tokenUri.TokenID = inscription.ID
	blockNumberMinted := fmt.Sprintf("%d", blockHeight)
	tokenUri.BlockNumberMinted = &blockNumberMinted
	tokenUri.Creator = &project.CreatorProfile
	tokenUri.CreatorAddr = project.CreatorAddrr
	tokenUri.Description = project.Description
	tokenUri.GenNFTAddr = project.GenNFTAddr

	mintedTime := time.Now()
	tokenUri.MintedTime = &mintedTime
	tokenUri.Name = inscription.Meta.Name
	tokenUri.Project = project
	tokenUri.ProjectID = project.TokenID
	tokenUri.ProjectIDInt = project.TokenIDInt
	tokenUri.IsOnchain = false
	tokenUri.CreatedByCollectionInscription = true
	count, err := u.Repo.CountTokenUriByProjectId(tokenUri.ProjectID)
	if err == nil && count != nil {
		tokenUri.OrderInscriptionIndex = int(*count + 1)
	}

	nftTokenUri := project.NftTokenUri
	logger.AtLog.Logger.Info("nftTokenUri", zap.Any("nftTokenUri", nftTokenUri))

	imageURI := fmt.Sprintf("https://generativeexplorer.com/content/%s", inscription.ID)

	//resp := &entity.InscriptionDetail{}
	//_, err = resty.New().R().
	//	SetResult(&resp).
	//	Get(fmt.Sprintf("%s/inscription/%s", u.Config.GenerativeExplorerApi, inscription.ID))
	//if err != nil {
	//	logger.AtLog.Logger.Error("err", zap.Error(err))
	//	return nil, err
	//}

	tokenUri.AnimationURL = imageURI
	tokenUri.Thumbnail = ""
	tokenUri.Image = imageURI
	tokenUri.ParsedImage = &imageURI
	tokenUri.Source = inscription.Source
	logger.AtLog.Logger.Info("mintedURL", zap.Any("imageURI", imageURI))

	_, err = u.Repo.UpdateOrInsertTokenUri(tokenUri.ContractAddress, tokenUri.TokenID, &tokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
		return nil, err
	}

	pTokenUri, err := u.Repo.FindTokenBy(tokenUri.ContractAddress, tokenUri.TokenID)
	if err != nil {
		return nil, err
	}

	// update project index and max supply
	index := project.MintingInfo.Index + 1
	maxSupply := project.MaxSupply
	if index > maxSupply {
		maxSupply = index
	}
	err = u.Repo.UpdateProjectIndexAndMaxSupply(project.TokenID, maxSupply, index)
	if err != nil {
		return nil, err
	}

	go u.TriggerPubsubTokenThumbnail(pTokenUri.TokenID)

	return pTokenUri, nil
}
