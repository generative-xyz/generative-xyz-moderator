package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"rederinghub.io/external/artblock"
	"rederinghub.io/external/nfts"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/external/ord_service"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"

	"go.uber.org/zap"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"rederinghub.io/utils/contracts/generative_project_contract"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/rediskey"
)

type uploadFileChan struct {
	FileURL *string
	Err     error
}

func (u Usecase) CreateBTCProject(req structure.CreateBtcProjectReq) (*entity.Projects, error) {
	logger.AtLog.Logger.Info("CreateBTCProject", zap.Any("CreateBtcProjectReq", zap.Any("req)", req)))

	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("CreateBTCProject", zap.Any("copier.Copy", err))
		return nil, err
	}

	mReserveMintPrice := helpers.StringToBTCAmount(req.ReserveMintPrice)

	mPrice := helpers.StringToBTCAmount(req.MintPrice)
	maxID, err := u.Repo.GetMaxBtcProjectID()
	if err != nil {
		logger.AtLog.Error("CreateBTCProject", zap.Any("err.GetMaxBtcProjectID", err))
		return nil, err
	}
	maxID = maxID + 1
	pe.TokenIDInt = maxID
	pe.TokenID = fmt.Sprintf("%d", maxID)
	pe.ContractAddress = os.Getenv("GENERATIVE_BTC_PROJECT")
	pe.MintPrice = mPrice
	pe.ReserveMintPrice = mReserveMintPrice
	pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(300000 / 4))).String() // will update after unzip and check data or check from animation url
	pe.IsHidden = true
	if req.IsHidden != nil {
		pe.IsHidden = *req.IsHidden
	}
	pe.Status = true
	pe.IsSynced = true
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}
	creatorAddrr, err := u.Repo.FindUserByWalletAddress(req.CreatorAddrr)
	if err != nil {
		logger.AtLog.Error("CreateBTCProject", zap.Any("err.FindUserByWalletAddress", err))
		return nil, err
	}

	if creatorAddrr.WalletAddressBTC == "" {
		creatorAddrr.WalletAddressBTC = req.CreatorAddrrBTC
		updated, err := u.Repo.UpdateUserByID(creatorAddrr.UUID, creatorAddrr)
		if err != nil {
			logger.AtLog.Error("CreateBTCProject", zap.Any("err.UpdateUserByID", err))

		} else {
			logger.AtLog.Info("updated.creatorAddrr", creatorAddrr)
			logger.AtLog.Info("updated", updated)
		}
	}

	isPubsub := false
	animationURL := ""
	zipLink := req.ZipLink
	if zipLink != nil && *zipLink != "" {
		now := time.Now().UTC()
		pe.IsHidden = true
		isPubsub = true
		pe.Status = false
		pe.IsSynced = false

		//create unzip's log here
		unzipLog := &entity.ProjectZipLinks{
			ProjectID: pe.TokenID,
			ZipLink:   *zipLink,
			Status:    entity.UzipStatusFail,
			Message:   "Create project",
			ReTries:   0,
		}

		unzipLog.Logs = []entity.ProjectZipLinkLog{}
		unzipLog.Logs = append(unzipLog.Logs, entity.ProjectZipLinkLog{
			Message:     "Create project",
			Status:      entity.UzipStatusFail,
			CreatedTime: &now,
		})

		unzipLog.UpdatedAt = &now
		unzipLog.CreatedAt = &now
		err = u.Repo.CreateProjectUnzip(unzipLog)
		if err != nil {
			logger.AtLog.Error("UnzipProjectFile.defer", zap.Any("projectID", pe.TokenID), zap.Error(err))
		}

	} else {
		if req.AnimationURL != nil {
			animationURL = *req.AnimationURL
			maxSize := helpers.CalcOrigBinaryLength(animationURL)
			pe.MaxFileSize = int64(maxSize)
			pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(maxSize / 4))).String()
			htmlContent, err := helpers.Base64Decode(strings.ReplaceAll(animationURL, "data:text/html;base64,", ""))
			if err == nil {
				isFullChain, err := helpers.IsFullChain(string(htmlContent))
				if err == nil {
					pe.IsFullChain = isFullChain
					logger.AtLog.Logger.Info("CreateBTCProject", zap.Any("isFullChain", zap.Any("isFullChain)", isFullChain)))
				} else {
					logger.AtLog.Error("CreateBTCProject", zap.Any("isFullChain", err))
				}
			} else {
				logger.AtLog.Error("CreateBTCProject", zap.Any("isFullChain", err))
			}
			nftTokenURI["animation_url"] = animationURL

			//Html
			htmlUrl, err := u.parseAnimationURL(*pe)
			if err == nil {
				animationHtml := fmt.Sprintf("%s", *htmlUrl)
				pe.AnimationHtml = &animationHtml
			}

		}
	}

	bytes, err := json.Marshal(nftTokenURI)
	if err != nil {
		logger.AtLog.Error("CreateBTCProject", zap.Any("marshal", err))
		return nil, err
	}
	nftToken := helpers.Base64Encode(bytes)
	now := time.Now().UTC()

	pe.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s", nftToken)
	pe.ProcessingImages = []string{}
	pe.MintedImages = nil
	pe.MintedTime = &now
	pe.CreatorProfile = *creatorAddrr
	pe.CreatorAddrrBTC = req.CreatorAddrrBTC
	pe.LimitSupply = 0
	pe.GenNFTAddr = pe.TokenID

	captureTime := entity.DEFAULT_CAPTURE_TIME
	if req.CaptureImageTime != nil && *req.CaptureImageTime != 0 {
		captureTime = *req.CaptureImageTime
	}

	pe.CatureThumbnailDelayTime = &captureTime
	if len(req.Categories) != 0 {
		pe.Categories = []string{req.Categories[0]}
	}

	if pe.Categories == nil || len(pe.Categories) == 0 {
		pe.Categories = []string{u.Config.OtherCategoryID}
	}

	pe.OpenMintUnixTimestamp = int(time.Now().Add(time.Hour * entity.DEFAULT_DELAY_OPEN_MINT_TIME_IN_HOUR).Unix())

	logger.AtLog.Logger.Info("CreateBTCProject", zap.Any("project", zap.Any("pe)", pe)))
	err = u.Repo.CreateProject(pe)
	if err != nil {
		logger.AtLog.Error("CreateBTCProject", zap.Any("CreateProject", err))
		return nil, err
	}

	if isPubsub {

		err = u.PubSub.Producer(utils.PUBSUB_PROJECT_UNZIP, redis.PubSubPayload{Data: structure.ProjectUnzipPayload{ProjectID: pe.TokenID, ZipLink: *zipLink}})
		if err != nil {
			logger.AtLog.Logger.Error("u.Repo.CreateProject", zap.Error(err))
			//return nil, err
		}

	} else {
		u.AirdropArtist(pe.TokenID, os.Getenv("AIRDROP_WALLET"), pe.CreatorProfile, 3)
	}

	go u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Project is created][project %s]", helpers.CreateProjectLink(pe.TokenID, pe.Name)), fmt.Sprintf("TraceID: %s", pe.TraceID), fmt.Sprintf("Project %s has been created by user %s", helpers.CreateProjectLink(pe.TokenID, pe.Name), helpers.CreateProfileLink(pe.CreatorAddrr, pe.CreatorName)))

	if pe.IsHidden && pe.IsSynced {
		ids, err := u.CreateDAOProject(context.Background(), &request.CreateDaoProjectRequest{
			ProjectIds: []string{pe.ID.Hex()},
			CreatedBy:  pe.CreatorAddrr,
		})
		if err != nil {
			logger.AtLog.Logger.Error("CreateDAOProject failed by", zap.Error(err))
		} else {
			logger.AtLog.Logger.Info("CreateDAOProject success",
				zap.String("project_id", pe.ID.Hex()),
				zap.Strings("ids", ids),
			)
		}
		if len(ids) > 0 {
			u.NotifyNewProject(pe, creatorAddrr, true, ids[0])
		}
	}

	return pe, nil
}

func (u Usecase) JobCheckAirdrop() error {
	airdrops, err := u.Repo.FindAirdropByStatus(0)
	if err != nil {
		fmt.Printf("JobCheckAirdrop - with err: %v", err)
		return err
	}
	logger.AtLog.Info(fmt.Sprintf("Start check airdrops len %d", len(airdrops)))
	for _, airdrop := range airdrops {
		if airdrop.Tx != "" {
			logger.AtLog.Info(fmt.Sprintf("Start check airdrop %s", airdrop.UUID), zap.Any("airdrop", airdrop))
			_, bs, err := u.buildBTCClient()

			if err != nil {
				fmt.Printf("JobCheckAirdrop - with err: %v", err)
				continue
			}
			// check with api:
			txInfo, err := bs.CheckTx(airdrop.Tx)
			if err != nil {
				fmt.Printf("JobCheckAirdrop - with err: %v", err)
				u.Repo.UpdateAirdropStatusByTx(airdrop.Tx, 2, "")
				continue
			}
			if txInfo.Confirmations > 0 {
				fmt.Printf("JobCheckAirdrop success - %v", txInfo)
				data, err := json.Marshal(txInfo)
				temp := ""
				if err == nil {
					temp = string(data)
				}
				u.Repo.UpdateAirdropStatusByTx(airdrop.Tx, 1, temp)
				go u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"),
					"Airdrop success",
					airdrop.ReceiverBtcAddressTaproot,
					fmt.Sprintf("Type: %d - file %s airdrop tx %s for userUUid %s", airdrop.Type, airdrop.File, airdrop.Tx, airdrop.Receiver))
				go u.NotifyNewAirdrop(airdrop)
			}
		}
	}
	return nil
}

func (u Usecase) JobCheckAirdropInit() error {
	airdrops, err := u.Repo.FindAirdropByStatus(-1)
	if err != nil {
		fmt.Printf("JobCheckAirdropInit - with err: %v", err)
		return err
	}
	logger.AtLog.Info(fmt.Sprintf("Start check JobCheckAirdropInit len %d", len(airdrops)))
	for _, airdrop := range airdrops {
		if airdrop.Type == 0 {
			// for airdrop artist
			// check something like
			projectId := airdrop.ProjectId
			project, err := u.Repo.FindProjectByTokenID(projectId)
			if err != nil {
				logger.AtLog.Error("JobCheckAirdropInit project not found", zap.Any("projectID", projectId))
				continue
			}
			if project.MintingInfo.Index == 0 {
				logger.AtLog.Error("JobCheckAirdropInit project still not mint", zap.Any("project", project))
				continue
			}
			mintPrice, e := strconv.Atoi(project.MintPrice)
			if e != nil {
				logger.AtLog.Error("JobCheckAirdropInit project get mint price", zap.Any("project", project))
				continue
			}
			if project.MintingInfo.Index*int64(mintPrice) < 430000 {
				logger.AtLog.Error("JobCheckAirdropInit project still not mint volum reach ~100usd", zap.Any("project", project))
				continue
			}
		}
		u.AirdropUpdateMintInfo(airdrop, os.Getenv("AIRDROP_WALLET"), 3)
	}
	return nil
}

func (u Usecase) AirdropUpdateMintInfo(airDrop *entity.Airdrop, from string, feerate int) (*entity.Airdrop, error) {
	mintReq := ord_service.MintRequest{
		WalletName:         from,
		ProjectID:          airDrop.ProjectId,
		DryRun:             false,
		AutoFeeRateSelect:  false,
		FeeRate:            feerate,
		FileUrl:            airDrop.File,
		DestinationAddress: airDrop.ReceiverBtcAddressTaproot,
	}
	logger.AtLog.Logger.Info(fmt.Sprintf("Mint airdrop request %v", mintReq), zap.Any("mintReq", mintReq))

	resp, respStr, err := u.OrdService.Mint(mintReq)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("OrdService.Mint airdrop %v %v", err, respStr), zap.Any("Error", err))
		return nil, err
	}
	logger.AtLog.Logger.Info("OrdService.Mint resp", zap.Any("Resp", zap.Any("resp)", resp)))

	_, err = u.Repo.UpdateAirdropMintInfoByUUid(airDrop.UUID, resp)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UpdateAirdropMintInfo airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	tmpText := resp.Stdout
	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
	btcMintResp := &ord_service.MintStdoputRespose{}
	bytes := []byte(jsonStr)
	err = json.Unmarshal(bytes, btcMintResp)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UpdateAirdropMintInfo Unmarshal airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}
	_, err = u.Repo.UpdateAirdropInscriptionByUUid(airDrop.UUID, btcMintResp.Reveal, btcMintResp.Inscription)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UpdateAirdrop Unmarshal airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}
	return airDrop, nil
}

func (u Usecase) AirdropArtist(projectid string, from string, receiver entity.Users, feerate int) (*entity.Airdrop, error) {
	if !helpers.IsOrdinalProject(projectid) {
		return nil, nil
	}
	if os.Getenv("ENV") != "mainnet" {
		return nil, nil
	}
	feerate = 3
	// get file
	random := rand.Intn(100)
	file := utils.AIRDROP_MAGIC
	if random >= 30 {
		file = utils.AIRDROP_SILVER
	} else if random < 30 && random >= 5 {
		file = utils.AIRDROP_GOLDEN
	}

	airDrop := &entity.Airdrop{
		File:                      file,
		Receiver:                  receiver.UUID,
		ReceiverBtcAddressTaproot: receiver.WalletAddressBTCTaproot,
		Type:                      0,
		ProjectId:                 projectid,
		OrdinalResponseAction:     nil,
		Status:                    -1,
		MintedInscriptionId:       "",
		Tx:                        "",
		InscriptionId:             "",
	}
	err := u.Repo.InsertAirdrop(airDrop)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("InsertAirdrop airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	/* not airdrop immediately anymore
	airDrop, err = u.AirdropUpdateMintInfo(airDrop, from, feerate)
	if err != nil {
		return nil, err
	}*/

	return airDrop, nil
}

func (u Usecase) AirdropCollector(projectid string, mintedInscriptionId string, from string, receiver entity.Users, feerate int) (*entity.Airdrop, error) {
	if !helpers.IsOrdinalProject(projectid) {
		return nil, nil
	}
	if os.Getenv("ENV") != "mainnet" {
		return nil, nil
	}
	// get file
	feerate = 3
	random := rand.Intn(100)
	file := utils.AIRDROP_MAGIC
	if random >= 13 {
		file = utils.AIRDROP_SILVER
	} else if random < 13 && random >= 3 {
		file = utils.AIRDROP_GOLDEN
	}

	airDrop := &entity.Airdrop{
		File:                      file,
		Receiver:                  receiver.UUID,
		ReceiverBtcAddressTaproot: receiver.WalletAddressBTCTaproot,
		Type:                      1,
		ProjectId:                 projectid,
		OrdinalResponseAction:     nil,
		Status:                    -1,
		MintedInscriptionId:       mintedInscriptionId,
		InscriptionId:             "",
		Tx:                        "",
	}
	err := u.Repo.InsertAirdrop(airDrop)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("InsertAirdrop airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	airDrop, err = u.AirdropUpdateMintInfo(airDrop, from, feerate)
	if err != nil {
		return nil, err
	}

	return airDrop, nil
}

func (u Usecase) IsTokenGatedNewUserAirdrop(user *entity.Users, whiteListEthContracts []string) (bool, error) {
	if len(whiteListEthContracts) == 0 {
		return false, nil
	}
	airdrop, err := u.Repo.FindAirdropByTokenGatedNewUser(user.UUID)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("ERROR AirdropTokenGatedNewUser"), zap.Any("error", err))
		return u.IsWhitelistedAddress(context.Background(), user.WalletAddress, whiteListEthContracts)
	} else {
		if airdrop != nil {
			logger.AtLog.Error(fmt.Sprintf("ERROR Exist AirdropTokenGatedNewUser"), zap.Any("airdrop", airdrop))
			return false, err
		}
		return u.IsWhitelistedAddress(context.Background(), user.WalletAddress, whiteListEthContracts)
	}
	return false, nil
}

func (u Usecase) IsArtistABNewUserAirdrop(user *entity.Users) (bool, error) {
	airdrop, err := u.Repo.FindAirdropByTokenGatedNewUser(user.UUID)
	flag := false
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("ERROR IsArtistABNewUserAirdrop"), zap.Any("error", err))
		flag = true
	} else {
		if airdrop != nil {
			logger.AtLog.Error(fmt.Sprintf("ERROR Exist IsArtistABNewUserAirdrop"), zap.Any("airdrop", airdrop))
			return false, err
		}
		flag = true
	}
	if flag {
		artblockService := artblock.NewArtBlockService(nil, "https://artblocks-mainnet.hasura.app")
		data, err := artblockService.GetArtist(strings.ToLower(user.WalletAddress))
		if err != nil {
			logger.AtLog.Error(fmt.Sprintf("Error IsArtistABNewUserAirdrop"), zap.Any("error", err))
			return false, err
		}
		if len(data.Data.Artists) != 1 {
			logger.AtLog.Error(fmt.Sprintf("Error IsArtistABNewUserAirdrop"), zap.Any("data.Data.Artists", data.Data.Artists))
			return false, nil
		}
		if strings.ToLower(data.Data.Artists[0].PublicAddress) != strings.ToLower(user.WalletAddress) {
			logger.AtLog.Error(fmt.Sprintf("Error IsArtistABNewUserAirdrop"), zap.Any("data.Data.Artists", data.Data.Artists))
			return false, nil
		}
		return true, nil
	}
	return false, nil
}

func (u Usecase) AirdropArtistABNewUser(from string, receiver entity.Users, feerate int) (*entity.Airdrop, error) {
	if os.Getenv("ENV") != "mainnet" && true {
		return nil, nil
	}
	if receiver.UUID == "" || receiver.WalletAddressBTCTaproot == "" {
		return nil, nil
	}

	isArtistABNewUserAirdrop, err := u.IsArtistABNewUserAirdrop(&receiver)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("Error AirdropArtistABNewUser"), zap.Any("error", err))
	}
	if !isArtistABNewUserAirdrop {
		return nil, nil
	}
	// get file
	feerate = 3
	random := rand.Intn(100)
	file := utils.AIRDROP_MAGIC
	if random >= 13 {
		file = utils.AIRDROP_SILVER
	} else if random < 13 && random >= 3 {
		file = utils.AIRDROP_GOLDEN
	}

	airDrop := &entity.Airdrop{
		File:                      file,
		Receiver:                  receiver.UUID,
		ReceiverBtcAddressTaproot: receiver.WalletAddressBTCTaproot,
		Type:                      3,
		ProjectId:                 "",
		OrdinalResponseAction:     nil,
		Status:                    -1,
		MintedInscriptionId:       "",
		InscriptionId:             "",
		Tx:                        "",
	}
	err = u.Repo.InsertAirdrop(airDrop)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("AirdropArtistABNewUser InsertAirdrop airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	airDrop, err = u.AirdropUpdateMintInfo(airDrop, from, feerate)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("AirdropArtistABNewUser AirdropUpdateMintInfo airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	return airDrop, nil

}

func (u Usecase) AirdropTokenGatedNewUser(from string, receiver entity.Users, feerate int) (*entity.Airdrop, error) {
	if os.Getenv("ENV") != "mainnet" && true {
		return nil, nil
	}
	if receiver.UUID == "" || receiver.WalletAddressBTCTaproot == "" {
		return nil, nil
	}
	whitelist := os.Getenv("WHITELIST_AIRDROP_TOKENGATED")
	if len(strings.TrimSpace(whitelist)) == 0 {
		return nil, nil
	}
	whitelistArr := strings.Split(whitelist, ",")
	isTokenGated, err := u.IsTokenGatedNewUserAirdrop(&receiver, whitelistArr)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("Error AirdropTokenGatedNewUser"), zap.Any("error", err))
	}
	if !isTokenGated {
		return nil, nil
	}

	// get file
	feerate = 3
	random := rand.Intn(100)
	file := utils.AIRDROP_MAGIC
	if random >= 13 {
		file = utils.AIRDROP_SILVER
	} else if random < 13 && random >= 3 {
		file = utils.AIRDROP_GOLDEN
	}

	airDrop := &entity.Airdrop{
		File:                      file,
		Receiver:                  receiver.UUID,
		ReceiverBtcAddressTaproot: receiver.WalletAddressBTCTaproot,
		Type:                      2,
		ProjectId:                 "",
		OrdinalResponseAction:     nil,
		Status:                    -1,
		MintedInscriptionId:       "",
		InscriptionId:             "",
		Tx:                        "",
	}
	err = u.Repo.InsertAirdrop(airDrop)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("AirdropTokenGatedNewUser InsertAirdrop airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	airDrop, err = u.AirdropUpdateMintInfo(airDrop, from, feerate)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("AirdropTokenGatedNewUser AirdropUpdateMintInfo airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	return airDrop, nil
}

func (u Usecase) AirdropPerceptronAuctionList() {
	//temp := `[{"address":"0x290cf6c2277f953f70b4fa32a01f1140f629acd6","price":0.123,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xa24b965d7ce7b20cb4696dd76eccdcf7f6001215","price":0.17,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x9a3514439aed2aa3c608fa3defcc16fe900efd1a","price":0.12,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x6b73f21f89c7149689dc34e03ed7b559e8396f86","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x9f8f3d7023e57cc94bf118713913bbf99edda34a","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5b3a3e754c2a138a80b44ea961bae33c5315b515","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x283c52890a6e05f948817daf56c252e67d1aa17b","price":0.1256,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xb34e7f7c33d1520027e880b490b5e1e3be473ce9","price":0.0666666,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x05eb7f0ebcfc8bee7d5283521a08ebef149569ed","price":0.124,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x81bc9a4600470346b1306e13843871a1624afac0","price":0.126,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xc93c7f71581dfeaab59bed908888dac5689f312a","price":0.0666,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xca8e7c89d791abcf4ac9b3481679e70d2c2db796","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xdae6ca75bb2afd213e5887513d8b1789122eaaea","price":0.065,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf8a065f287d91d77cd626af38ffa220d9b552a2b","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6232d7a6085d0ab8f885292078eeb723064a376b","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4a46a8aed4aff4eeafbba10deed91d2450233c20","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x112bb5dd677ec044f24c045db1a0ab3b167841ea","price":0.32,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0x68bca5a8bdebe05fb8a6648c7316b4eb7e19a064","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xdd3e9d0ee979e5c1689a18992647312b42d6d8f3","price":0.111111,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xafbc3f98eedb5f9a25a4ab2232d1346612efe77c","price":0.2069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x204a731a626438a0b2a657378668cfbfe70a90f5","price":0.1469,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xec56139a5b2b8958d927a5b7986de7010ce3e699","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x688266aca50e59016ef912e2e7ff61f70e70bdac","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xaa6c28a775ca64c4104057eff7d2f0a07ce89b29","price":0.22,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xa7a3a06e9a649939f60be309831b5e0ea6cc2513","price":0.2,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x7729a5cfe2b008b7b19525a10420e6f53941d2a4","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xe90da7ea578fef3365631aaf7b9d5094152bfae9","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x229bf0492a6c1b3ccb82bbdc4dc488124ddee8c8","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x083ae6f19a964d0339c91e808b66e060f9d8d9ff","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x2b614bf0cff1dff5d099f88a98ae449604379bc2","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4f7f3606bcf2a4f805e2f1ce85d9a6d95fa85c6d","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1ea65a3cbf739564dad1f073f12dd6c1673c1db6","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x08c265911bd4639e72bb547348b65f9cee0b1c66","price":0.071,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd45a29008cec83c58cb332ef2d68d4dafede1635","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6f418683ffc57ca7bd93b3c8a595de2d8035df77","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x9d0da7fa00bed2cfa6f1dec9d9c16f075a589673","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6186290b28d511bff971631c916244a9fc539cfe","price":0.3,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0xd649946076d587af74cad517734935018a54cc2a","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5957f814aafe77f669ff7ec6230a76a639425e1e","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x480c68f002da6bf9a932814737f0317ef58a49fb","price":0.0769,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x58d49377c74fe5aa1c098d9ed4161248b73faa30","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xfb83388cb78a07bd7e163e70aecd1d1e2a19d81e","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x52e7dcdba39f37683bbdd83b8c32be44809ec3c8","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x84300dcc7ca9cf447e886fa17c11fa22557d1af0","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x88d3574660711e03196af8a96f268697590000fa","price":0.11,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xddf06174511f1467811aa55cd6eb4efe0dffc2e8","price":0.056,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x82b332fdd56d480a33b4da58d83d5e0e432f1032","price":0.062,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xba8e220834c32fac376bbfb33820c001f022d72e","price":0.052,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x271d52a030b09755390e997b3052ceccc0920008","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xb630abd9a5367763b7cba316e870c4a54064cc9f","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xb3a898e8e48ed3fb1b6e147448de0f84be977876","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xbe13c5e443f04b8925f5467f579b942f38671e54","price":0.21,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x24a905d3388bbc88d9df5ac60ac8a44c70f973ce","price":0.3,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0x68bb4c94609afe07bf5f8f1e14a0289beef3577d","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf6a9dc9a41e8817fbc86532d53d89ca80d9aa46b","price":0.0666,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x098066e8530eec4c4dede0f1b5a8035706bbc7d0","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x84cd27c4d94588f09a201b0ce6c07945d67b2a2c","price":1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0xb0985d0583b8de16393d22249fba0351992c190b","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf338afe0b7d435a59e57e450d408eb4be3b62e99","price":0.115,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x9ff21d0c36f8f871275272bd417f6b748ed13ffb","price":0.05742,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd561ba5bdbfea7a39bf073b7520a7273bc767131","price":0.0511,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x852548338cd5a8de1384fde3bfc678a8669c0f4d","price":0.071,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x302da6acbd97c425192af64b9ae739ca5c2b52bc","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7915e43086cd78be341df73726c0947b6334b978","price":0.051009,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xa3c277b8f35881cbdb017e52bcc376b3ce8f21da","price":0.071,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x296e0c21db4061ebf971e55d5db85011e7ff9797","price":0.11,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x000b75fcdc15d41277deb033c72d2c8d774ccced","price":0.3,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0x41955ab7d12f9f6c03de972b91d9b895d9c2eaf8","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc5ebfede936447ef2f1b2de2ebbb0601f58a29a5","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xe6e4d92009406d08851c2e65ce6dd324ad76a87e","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7f04c4387423c5460f0a797b79b7de2a4769567a","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x48bb463710c5a42ef01e1fdbeae879758b597f4c","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7b3dbd9fe6ebc6bd76dbb198c65691fc6c1165ce","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x171446b041c6683e5f138b8a3f263bbffb8ee74a","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1b13cc4c11666a83518df8df5af2976e8be4d7fa","price":0.1069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xbc2f3873cd6650474d7153f3ecc0a4ac23968fff","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6c85c13d32082da3b36c2aa1247e698fde529aea","price":0.1001,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xebfd774c1c2008e56ce40e0a4504ebecc81b1921","price":0.0621,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x0a690b298f84d12414f5c8db7de1ece5a4605877","price":0.075,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xa63c2a96b84b73867c0c6a89331907ebb4c94d56","price":0.0514,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xfd5e9bb7e45fa0eb44ac7ec72643a1cc5f185dbc","price":0.054,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x478087e12db15302a364c64cdb79f14ae6c5c9b7","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd5abe2b92817be749ccc077ae19a4cffc0f47c16","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc589630a48c26920bd1deca9dd522aa547380f4e","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xb8adb4ad39f4892ee0c79f9e3dde5904ac6ae291","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x0eb9a7ff5cbf719251989caf1599c1270eafb531","price":0.0555,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xb47666a6d0096d5871c51d9bee04b6f97f773f6b","price":0.0600015,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x02c2ca0c5e140d82cd10bf5641515f10b17a2b9b","price":0.1003,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x1ea75eb1459c11256a9ceaffb03ebd5caa8a52c8","price":0.11002,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x1341df844780b66af4ccc98ae0f34be87eabe1d5","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd8f27e7055d854e660f65d6a2960f867252a8bc0","price":0.0623,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x19748446a67b690ef1dd13ee61a615e9028bc6e0","price":0.075,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1b042ab6bcba48b18eb24e4e23424bf4c10bae90","price":0.101,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x29d3c95d4d8a42ad598112a1c6f909e010ad37dc","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xa88e4a192f3ff5e46dcc96efefb38dfec7bb250c","price":0.065,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7ae2c7ca28575b3225dacd91242dae4420d17323","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc3a03a0cb785441943c7c7fcb461809b841fb124","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xa351ec9d2bc46fa9030eedb069a96757d406e06a","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x54021e58af1756dc70ce7034d7636de2d2f1fa74","price":0.056,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6f808a8be84dc3aa67f075feaf39db43082c005c","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x3a6e84c67dab20dc32400ab98c95c91ce2d721d8","price":0.27,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x5b5fc02d41eaafa7ecde3c02c3e5c59110a77d99","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x727fb80e6ab480c1cd34d55a63fce6395df7c642","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf57a9b1f574b1f80c8cebe252706bb8b4d783d21","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xcfb098c1d44eb12f93f9aaece5d6054e2a2240ab","price":0.102,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x37eb49074b9c7061a51017321ef8ba9ae3fdc708","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x9aa824a302597e3477d0435cfc5b1e1b9eb23449","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xa4e7e0dca4393cbfefb91d4094530dbc0bda1fbc","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x94803d6c3cb805a308b9814202d7d01dcb26049d","price":0.0506969,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc7734f68f937715a640f65e1a71b316f83ee31eb","price":0.0569,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xbf5aa87878c968c9d062f51cabd04b66023c8e36","price":0.08,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd9a04c14c31be4f764b111d0b147bb332fbc2fbc","price":0.2,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xba246a7df93bf8c7009dba1fee636ac40c210f00","price":0.079,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x0f1cd3003768759c899a026d6ccf0b62b041164e","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4c70afea2ec4fe9a42d35396806b7529f62a4f89","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5e5c35235b7e288ba0547c622be93c39058538a0","price":0.064,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4a39ae58b605102913ac19b7c071da75b55b2674","price":0.101,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xdacddea78d9b9c7eba66e4da7930e92b357425ce","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xbfdf3266847b0cc9cf9bdc626bef48ff9c46e9cd","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4914888353151c450fbd65bb943767e5b45d1050","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7c03f5cc04e5a5b81956557b9e7429e7ab65884b","price":0.052,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x8f1213383bf38bdc2370ee89d219fe141b184530","price":0.052,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x44f0630f0be9c7b6a5188b9b289f2752f56b443b","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1a1f337faa595a3603b4a3e67776150e7883954d","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xfe57fdcd9d0df47356ccdebe08f3dc454e472f9b","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4a1a0aeea2a03ed451b577f51f1d7e5568f29736","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xb7ef2a5637021424026edb948bb8779ffd76c190","price":0.11,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x1bec0e6266e97f0c7c42aa720959e46598dd8ce1","price":0.065,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5801646569c6dbae73432d85979aa57bacc07aa5","price":0.072,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5225787a5d9fbe412ff9a4cae43c1c3f3e026003","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x34a32ad4ba1ea1eb02ccd3ed5b9af9a8d8ea07a8","price":0.12,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x2bc52900a118390ea8a97b29eb5e7e75f1cd912c","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x162031b2f953d7ba0a5638b02f910af1c5279754","price":0.0601,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5a7283ce9c4b72026c3d07ee9d604103fb351771","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd61b77bc330ee594712a83fe3254b37e540d89c8","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc858c02f6e8e55530646ce6ecdb5b6b0817102c8","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc3771f93fd002fa3fd195fb175507bb701428b85","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xba768efae34d2880dfcccdf3b3c433c7a19ae102","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x77e143163904960528c60c619c58f73718b0a098","price":0.0711,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc20ca4126d455d1156da91408ba48d05b0687fd0","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7736e5af60c6535c3df1ed4f9c196dde1b5e2a4a","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc77890b96acf67e185474d422216fb546f990aff","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1889120e81099e7de90a3a746ccee1cadc6e8d05","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x969eac041cfa4bac53109bffd20268293a68d805","price":0.051002,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x3ba2dea8301aaf9536f2bc6233004628fed566a4","price":0.25,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x493304495b1d6708c40d47f056547eceb8834373","price":0.066,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xe22692b30a35fa6efe572909d875173d4bc21eae","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc886db8b8cd260f5ee38ba3d8f8e9324ee27ea33","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x368da465f18715367b31bc2d9162742ed5353539","price":0.0712,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x651da61ec8c9512b816086835c7f90c485850b62","price":0.1,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4cfd427ac7217ab1768f410efc33a37132b8f3c9","price":0.063,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc1f44c447307daaf106bdec3db4a5d660798d81d","price":0.063,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x9bdc866b4b03452df00c8c67a4a215c104dc8d41","price":0.053,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xbc7da0f4d3681d33c8be5eb3405be83eb34d697e","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x11c87f72e0d0990ac73fe4fd3be3f26adfa0e607","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xe4797e3e062b799a203fda426f916878a9f669cb","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x2336e78ad16ca2899c01e55e6a98dacd9f56d877","price":0.33,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"},{"address":"0x6873a346e0dddbe8fc7b83d6bbd713044d7f090f","price":0.075,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4819951ef30d768046212e705b92993a91e906d1","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x99d529f064cb498ef7b3798042194175da1f4d8a","price":0.1005,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0xc0452507b41d2da08a1630394d064b45617017c6","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x4970d00357ff515858b3191ed577bead738acf3d","price":0.061,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xc591674216324dc6f5496be098dfb52b674cbaca","price":0.063,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xce19e315f1abb6e921b057f050d67526b9a2f0ed","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x884ff907d5fb8bae239b64aa8ad18ba3f8196038","price":0.0551,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xaad0ac2ef7ec521d0990ca83586239dd5e915688","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x03b56b2157174c7a063d9c2f64955bdba24abbc9","price":0.075,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7d513fd6562ce83ba828e2d0aba957b763484cb9","price":0.0514,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x3a9eb2d9ef30e121f6fb4a0e4d3df3175381d2eb","price":0.057,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x9a8992f0db3e10d0ae4934df06a92a0a1d961d0c","price":0.0615,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x46034f239680a07937ad8389c23ab3bb9afba3fd","price":0.125,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x5f7c27402603d2607b07c19fcd41c84a1f5557ac","price":0.0601,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x3d2f9441fde0db524a77bf9fdda610538fcccd16","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xedd1f30d69898e4cb710cfb47c6114d31e6fed06","price":0.054321,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf2439241881964006369c0e2377d45f3740f48a0","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x02d78c75532a490dae2596e6b07a75857815cd70","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x8264e9e0f4cbcbbbb3f8ecaec0a625b590ae790e","price":0.0713,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x12ae2b6f5bb10ff14db81c1083b33dd4dd0f82be","price":0.067,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xbc5dbeb86d062ecaee41d5007ad6533f63711545","price":0.0611,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x96fd61202a698ee3eac21e247a6b209ea5ffeb91","price":0.079,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1ff2ca546389e6c0161b07f938de5a2b9f0265a5","price":0.08,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x35f2752974fad42cfc5d28d2c84e3503017a826b","price":0.057,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x94db5f225a1f6968cd33c84580c0adae52a04edf","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7edb0cbfe1217fe53af0fb871d92c1c627c6673c","price":0.07,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xae33e6d251d8669e512d071482cbb8a9ca50a25a","price":0.18,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x9674ede5eb609f9483dd1682f2160ec45105250d","price":0.0615,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x17cbd516166720b8c75fe75235ed0218d9cfdfbe","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x17e31bf839acb700e0f584797574a2c1fde46d0b","price":0.0520069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x6853a596d6d7264d3622546da3b891b6fe17eb82","price":0.0654321,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xba1bfcc44843babbab4a7e3ecacff886cc8c62ea","price":0.067,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xd210a01936a901c79254de8fc693f743e5757f8f","price":0.05555,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x01c2e7b1de06da53bd0ec82fdb59e5767b8c6da1","price":0.08,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf96d35106618e07a79febf2d346fcd14ac1e241c","price":0.13,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"},{"address":"0x330bce303d27df0eb1b856da3464278db1db1ac5","price":0.052,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x3016dd9dd812290122ec453cb01b48d5a4aff602","price":0.053,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x0b7576a64a0f4b4924d55ed328ede4979446521b","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x845a23ff9f59f8f788f9b94181e5326fcb8c9f6a","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xf611318ea51897ae2a98e537e55112c062febf3b","price":0.051,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x76da715b266323f4eb9c9ade2127e0611f9f6c30","price":0.053,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x0a52a605927bfaab5d781eb27dbf44534d9601db","price":0.053,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7ad785503a44c786e8ba25354109c230c3e04621","price":0.06,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xea6955894ba0953fa8e18da9d3d432e3c6b8b06e","price":0.0511,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xac693e7e769facc6f238306e10a865438a0a0e24","price":0.052,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x01c66e6ed4d287ba0fbaea1b6327b2335b7df931","price":0.056,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x5470c5a6fce7447afd2c9be3a0f25e362c093661","price":0.055,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xfd63d9503e4f0d538b48c163c224d93e0cc9537c","price":0.065,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0xe9ed3ad8e68b3925a33cab867a29c73e8357cfc4","price":0.0521,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x7daac88c492c641a0ac8d08420b6a7d78764615d","price":0.0611,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x1526474a52e17078dd13700e9d4b10eddee7e536","price":0.069,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"},{"address":"0x196a3dc8446920cef0f0d1f6bf7ba5b40702c79f","price":0.059,"key":"https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"}]`
	temp := ``
	var data []map[string]interface{}
	e := json.Unmarshal([]byte(temp), &data)
	if e == nil {
		for _, v := range data {
			receiver, e := u.Repo.FindUserByAddress(fmt.Sprintf("%v", v["address"]))
			if e == nil && receiver.UUID != "" {
				u.AirdropPerceptronAuction(os.Getenv("AIRDROP_WALLET"), *receiver, 3, fmt.Sprintf("%v", v["key"]))
			}
		}
	}

}

func (u Usecase) AirdropPerceptronAuction(from string, receiver entity.Users, feerate int, file string) (*entity.Airdrop, error) {
	if os.Getenv("ENV") != "local" && true {
		return nil, nil
	}
	if receiver.UUID == "" || receiver.WalletAddressBTCTaproot == "" {
		return nil, nil
	}

	airDrop := &entity.Airdrop{
		File:                      file,
		Receiver:                  receiver.UUID,
		ReceiverBtcAddressTaproot: receiver.WalletAddressBTCTaproot,
		Type:                      4,
		ProjectId:                 "",
		OrdinalResponseAction:     nil,
		Status:                    -1,
		MintedInscriptionId:       "",
		InscriptionId:             "",
		Tx:                        "",
	}
	err := u.Repo.InsertAirdrop(airDrop)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("AirdropPerceptronAuction InsertAirdrop airdrop %v %v", err, airDrop), zap.Any("Error", err))
		return nil, err
	}

	//airDrop, err = u.AirdropUpdateMintInfo(airDrop, from, feerate)
	//if err != nil {
	//	logger.AtLog.Error(fmt.Sprintf("AirdropPerceptronAuction AirdropUpdateMintInfo airdrop %v %v", err, airDrop), zap.Any("Error", err))
	//	return nil, err
	//}

	return airDrop, nil
}

func (u Usecase) resolveMintPriceBTC(priceStr string) string {
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return priceStr
	}
	return strconv.FormatFloat(float64(price)/1e8, 'f', -1, 64) + " BTC"

}

func (u Usecase) resolveShortName(userName string, userAddr string) string {
	if userName != "" {
		return userName
	}
	end := 10
	if end > len(userAddr) {
		end = len(userAddr)
	}
	return userAddr[:end]
}

func (u Usecase) resolveShortDescription(description string) string {
	if len(description) > 300 {
		return description[:250] + "..."
	}
	return description
}

func (u Usecase) UpdateBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	if req.ProjectID == nil {
		err := errors.New("ProjectID is requeried")
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if req.CreatetorAddress == nil {
		err := errors.New("CreatorAddress is requeried")
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if strings.ToLower(p.CreatorAddrr) != strings.ToLower(*req.CreatetorAddress) {
		err := errors.New("Only owner can update this project")
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if req.Name != nil && *req.Name != "" {
		p.Name = *req.Name
	}

	if req.Description != nil && *req.Description != "" {
		p.Description = *req.Description
	}

	if req.Thumbnail != nil && *req.Thumbnail != "" {
		bas64Data := strings.ReplaceAll(p.NftTokenUri, "data:application/json;base64,", "")
		bytes, err := helpers.Base64Decode(bas64Data)

		nftTokenURI := make(map[string]interface{})
		err = json.Unmarshal(bytes, &nftTokenURI)
		if err == nil {
			nftTokenURI["image"] = *req.Thumbnail
			bytes, err := json.Marshal(nftTokenURI)
			if err == nil {
				nftToken := helpers.Base64Encode(bytes)
				spew.Dump(fmt.Sprintf("data:application/json;base64,%s", nftToken))
				p.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s", nftToken)
			}
		}
		p.Thumbnail = *req.Thumbnail

	}
	needSetExpireAvailableDaoProject := false
	if req.IsHidden != nil && *req.IsHidden != p.IsHidden {
		if !*req.IsHidden {
			if u.IsProjectReviewing(context.Background(), p.ID.Hex()) {
				return nil, errors.New("Collection is reviewing")
			}
		} else {
			needSetExpireAvailableDaoProject = true
		}
		p.IsHidden = *req.IsHidden
	}

	if len(req.Categories) > 0 {
		p.Categories = []string{req.Categories[0]}
	}

	p.Reservers = req.Reservers
	if req.ReserveMintLimit != nil {
		p.ReserveMintLimit = *req.ReserveMintLimit
	}

	if req.ReserveMintPrice != nil && *req.ReserveMintPrice != "" {
		mReserveMintPrice := helpers.StringToBTCAmount(*req.ReserveMintPrice)
		p.ReserveMintPrice = mReserveMintPrice
	}

	if req.MaxSupply != nil && *req.MaxSupply != 0 && *req.MaxSupply != p.MaxSupply {
		// if p.MintingInfo.Index > 0 {
		// 	err := errors.New("Project is minted, cannot update max supply")
		// 	logger.AtLog.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }

		p.MaxSupply = *req.MaxSupply
	}

	if req.Royalty != nil {
		// if *req.Royalty > 2500 {
		// 	err := errors.New("Royalty must be less than 25")
		// 	logger.AtLog.Error("pjID.empty", err.Error(), err)
		// 	return nil, err
		// }

		// if *req.Royalty != p.Royalty && p.MintingInfo.Index > 0 {
		// 	err := errors.New("Project is minted, cannot update max supply")
		// 	logger.AtLog.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }

		p.Royalty = *req.Royalty
	}

	if req.MintPrice != nil {
		// mFStr := p.MintPrice
		reqMfFStr := helpers.StringToBTCAmount(*req.MintPrice)
		// if p.MintingInfo.Index > 0 && mFStr != reqMfFStr.String() {
		// 	err := errors.New("Project is minted, cannot update mint price")
		// 	logger.AtLog.Error("pjID.minted", err.Error(), err)
		// 	return nil, err
		// }
		p.MintPrice = reqMfFStr
	}

	if req.CaptureImageTime != nil && *req.CaptureImageTime != 0 {
		if p.CatureThumbnailDelayTime != nil && *p.CatureThumbnailDelayTime != *req.CaptureImageTime {
			p.CatureThumbnailDelayTime = req.CaptureImageTime
		}

		if p.CatureThumbnailDelayTime == nil {
			p.CatureThumbnailDelayTime = req.CaptureImageTime
		}
	}

	if req.LimitMintPerProcess != nil {
		if p.LimitMintPerProcess != *req.LimitMintPerProcess {
			p.LimitMintPerProcess = *req.LimitMintPerProcess
		}
	}

	if req.Index != nil {
		p.MintingInfo.Index = *req.Index
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Error("updated", err.Error(), err)
		return nil, err
	}
	if needSetExpireAvailableDaoProject {
		go u.SetExpireAvailableDAOProject(context.TODO(), p.ID)
	}
	logger.AtLog.Info("updated", updated)
	return p, nil
}

func (u Usecase) DeleteBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		logger.AtLog.Error("DeleteProject", zap.Any("err.FindProjectBy", err))
		return nil, err
	}
	if strings.ToLower(p.CreatorAddrr) != strings.ToLower(*req.CreatetorAddress) {
		logger.AtLog.Error("DeleteProject", zap.Any("err.CreatorAddrr", err))
		return nil, err
	}

	p.IsSynced = false
	p.Status = false
	p.IsHidden = true

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Error("UpdateProject", zap.Any("err.UpdateProject", err))
		return nil, err

	}

	_ = u.RedisV9.DelPrefix(context.TODO(), rediskey.Beauty(entity.DaoProject{}.TableName()).WithParams("list").String())

	logger.AtLog.Info("updated", updated)
	logger.AtLog.Logger.Info("UpdateProject", zap.Any("project", zap.Any("p)", p)))
	return p, nil
}

func (u Usecase) SetCategoriesForBTCProject(req structure.UpdateBTCProjectReq) (*entity.Projects, error) {

	if req.ProjectID == nil {
		err := errors.New("ProjectID is requeried")
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	p, err := u.Repo.FindProjectByTokenID(*req.ProjectID)
	if err != nil {
		logger.AtLog.Error("pjID.empty", err.Error(), err)
		return nil, err
	}

	if len(req.Categories) > 0 {
		p.Categories = []string{req.Categories[0]}
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Error("updated", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("updated", updated)
	return p, nil
}

func (u Usecase) UpdateProject(req structure.UpdateProjectReq) (*entity.Projects, error) {
	p, err := u.Repo.FindProjectBy(req.ContracAddress, req.TokenID)
	if err != nil {
		logger.AtLog.Error("UpdateProject", zap.Any("err.FindProjectBy", err))
		return nil, err
	}

	if req.Priority != nil {
		priority := 0
		p.Priority = &priority
	}

	if len(p.ReportUsers) >= u.Config.MaxReportCount {
		p.IsHidden = true
	}
	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Error("UpdateProject", zap.Any("err.UpdateProject", err))
		return nil, err
	}

	logger.AtLog.Info("updated", updated)
	logger.AtLog.Logger.Info("UpdateProject", zap.Any("project", zap.Any("p)", p)))
	return p, nil
}

func (u Usecase) ReportProject(tokenId, iWalletAddress, originalLink string) (*entity.Projects, error) {
	p, err := u.Repo.FindProjectByTokenID(tokenId)
	if err != nil {
		logger.AtLog.Error("ReportProject.FindProjectBy", err.Error(), err)
		return nil, err
	}

	for _, r := range p.ReportUsers {
		if r.ReportUserAddress == iWalletAddress {
			return nil, errors.New("You have already reported before.")
		}
	}

	rep := &entity.ReportProject{
		ReportUserAddress: iWalletAddress,
		OriginalLink:      originalLink,
	}

	p.ReportUsers = append(p.ReportUsers, rep)
	if len(p.ReportUsers) >= u.Config.MaxReportCount {
		p.IsHidden = true
		p.Status = false
		u.NotifiNewProjectHidden(p)
	}

	updated, err := u.Repo.UpdateProjectFields(p.UUID, map[string]interface{}{
		"isHidden": p.IsHidden,
		"status":   p.Status,
	})

	if err != nil {
		logger.AtLog.Error("UpdateProject.ReportProject", err.Error(), err)
		return nil, err
	}
	logger.AtLog.Info("updated", updated)

	u.NotifyWithChannel(
		os.Getenv("SLACK_PROJECT_CHANNEL_ID"),
		fmt.Sprintf("[Project is reported][projectID %s]", p.TokenID),
		"",
		fmt.Sprintf("Project %s has been report by user %s - original link: %s", p.Name, iWalletAddress, originalLink),
	)
	u.NotifiNewProjectReport(p, originalLink, iWalletAddress)

	return p, nil
}

func (u Usecase) GetProjectByGenNFTAddr(genNFTAddr string) (*entity.Projects, error) {
	project, err := u.Repo.FindProjectByGenNFTAddr(genNFTAddr)
	return project, err
}

func (u Usecase) GetProjects(req structure.FilterProjects) (*entity.Pagination, error) {
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	//if !u.CheckExisted("6406e7abb90f7fc13f55490c", pe.CategoryIds) && !u.CheckExisted("63f8325a1460b1502544101b", pe.CategoryIds) {
	//	pe.CustomQueries = make(map[string]primitive.M)
	//	pe.CustomQueries["$expr"] = bson.M{"$lt": bson.A{"$index", "$maxSupply"}}
	//}

	projects, err := u.Repo.GetProjects(*pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetProjects", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetAllProjects(req structure.FilterProjects) (*entity.Pagination, error) {
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	projects, err := u.Repo.GetAllRawProjects(*pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetProjects", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) CheckExisted(s string, arr []string) bool {
	for _, item := range arr {
		if item == s {
			return true
		}
	}

	return false
}

func (u Usecase) GetUpcommingProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}

	now := time.Now().UTC().Unix()
	pe.CustomQueries = make(map[string]primitive.M)
	pe.CustomQueries["openMintUnixTimestamp"] = bson.M{"$gt": now}

	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	projects, err := u.Repo.GetProjects(*pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetProjects", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetRandomProject() (*entity.Projects, error) {

	caddr := os.Getenv("RANDOM_PR_CONTRACT")
	pID := os.Getenv("RANDOM_PR_PROJECT")

	if caddr != "" && pID != "" {
		return u.GetProjectDetail(structure.GetProjectDetailMessageReq{
			ContractAddress: caddr,
			ProjectID:       pID,
		})
	}

	key := helpers.ProjectRandomKey()

	//always reload data
	go func() {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			return
		}
		u.Cache.SetData(key, p)
	}()

	cached, err := u.Cache.GetData(key)
	if err != nil {
		p, err := u.Repo.GetAllProjects(entity.FilterProjects{})
		if err != nil {
			logger.AtLog.Error("u.Repo.GetProjects", err.Error(), err)
			return nil, err
		}
		u.Cache.SetData(key, p)
	}

	cached, err = u.Cache.GetData(key)
	projects := []entity.Projects{}
	bytes := []byte(*cached)
	err = json.Unmarshal(bytes, &projects)
	if err != nil {
		logger.AtLog.Error("json.Unmarshal", err.Error(), err)
		return nil, err
	}

	if len(projects) == 0 {
		err := errors.New("Project are not found")
		logger.AtLog.Error("Projects.are.not.found", err.Error(), err)
		return nil, err
	}

	timeNow := time.Now().UTC().Nanosecond()
	rand := int(timeNow) % len(projects)

	//TODO - cache will be applied here

	projectRand := projects[rand]
	return u.GetProjectDetail(structure.GetProjectDetailMessageReq{
		ContractAddress: projectRand.ContractAddress,
		ProjectID:       projectRand.TokenID,
	})
}

func (u Usecase) GetMintedOutProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetMintedOutProjects(*pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetMintedOutProjects", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetProjectDetail(req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	logger.AtLog.Logger.Info("GetProjectDetail", zap.Any("req", req))
	c, _ := u.Repo.FindProjectByProjectIdWithoutCache(req.ProjectID)

	if (c == nil) || (c != nil && !c.IsSynced) || c.MintedTime == nil {
		return nil, errors.New("project is not found")
	}
	logger.AtLog.Logger.Info("GetProjectDetail", zap.Any("project", zap.Any("c)", c)))
	return c, nil
}

// only using for project detail api, support est fee:
func (u Usecase) GetProjectDetailWithFeeInfo(req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	logger.AtLog.Logger.Info("GetProjectDetail", zap.Any("req", req))
	c, err := u.Repo.FindProjectByProjectIdWithoutCache(req.ProjectID)

	if err != nil {
		return nil, err
	}

	// fmt.Println("c.MintedTime", c)

	if (c == nil) || (c != nil && !c.IsSynced) || c.MintedTime == nil {
		return nil, errors.New("project is not found")
	}
	if c.MintingInfo.Index < c.MaxSupply {

		if len(req.UserAddressToCheckDiscount) > 0 {
			if len(c.Reservers) > 0 {
				for _, address := range c.Reservers {
					if strings.EqualFold(address, req.UserAddressToCheckDiscount) {

						// get list item mint:
						countMinted := 0
						mintReadyList, _ := u.Repo.GetLimitWhiteList(req.UserAddressToCheckDiscount, req.ProjectID)

						for _, mItem := range mintReadyList {

							if mItem.IsConfirm {
								if mItem.Status == entity.StatusMint_Minting || mItem.Status == entity.StatusMint_Minted || mItem.IsMinted {
									countMinted += 1
								}

							} else if mItem.Status == entity.StatusMint_Pending || mItem.Status == entity.StatusMint_WaitingForConfirms {
								if time.Since(mItem.ExpiredAt) < 1*time.Second {
									countMinted += mItem.Quantity
								}
							}
						}
						maxSlot := c.ReserveMintLimit - countMinted

						if maxSlot <= 0 {
							c.Reservers = []string{}
						}

						break
					}
				}

			}
		}

		mintPrice, ok := big.NewInt(0).SetString(c.MintPrice, 10)
		if !ok {
			mintPrice = big.NewInt(0)
		}

		// cal fee:
		feeInfos, err := u.calMintFeeInfo(mintPrice.Int64(), c.MaxFileSize, entity.DEFAULT_FEE_RATE, 0, 0)
		if err != nil {
			logger.AtLog.Error("u.calMintFeeInfo.Err", err.Error(), err)
			return nil, err
		}

		// set price, fee:
		c.NetworkFee = feeInfos["btc"].NetworkFee
		c.NetworkFeeEth = feeInfos["eth"].NetworkFee

		c.MintPriceEth = feeInfos["eth"].MintPrice

		bidProjectID := utils.BidProjectIDProd
		if u.Config.ENV == "develop" {
			bidProjectID = utils.BidProjectIDDev
		}

		// return list winner:
		if strings.EqualFold(c.TokenID, bidProjectID) {
			auctionWinnerList, err := u.GetAuctionListWinnerAddressFromConfig()
			if err == nil {
				c.AuctionWinnerList = auctionWinnerList
			}
		}
	}

	go func() {
		//upload animation URL
		if c.AnimationHtml == nil {

			htmlUrl, err := u.parseAnimationURL(*c)
			if err != nil {
				return
			}

			animationHtml := fmt.Sprintf("%s", *htmlUrl)
			c.AnimationHtml = &animationHtml

			// _, err = u.Repo.UpdateProject(c.UUID, c) // remove for safe...
			_, err = u.Repo.UpdateProjectAnimationHtml(c.UUID, animationHtml)
			if err != nil {
				return
			}
		}

	}()

	// get index if project in TC:
	if c.IsMintTC() {
		contract, _ := generative_nft_contract.NewGenerativeNftContract(common.HexToAddress(c.GenNFTAddr), u.TcClient.GetClient())
		if contract != nil {
			projectContract, err := contract.Project(nil)
			if err == nil {
				c.MintingInfo.Index = projectContract.Index.Int64()
			}
		}
	}

	// logger.AtLog.Logger.Info("GetProjectDetail", zap.Any("project", zap.Any("c)", c)))
	return c, nil
}

func (u Usecase) GetProjectVolumn(req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	logger.AtLog.Logger.Info("GetProjectVolumn", zap.Any("req", zap.Any("req)", req)))
	c, err := u.Repo.FindProjectWithoutCache(req.ContractAddress, req.ProjectID)
	if err != nil {
		logger.AtLog.Error("GetProjectVolumn", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("GetProjectDetail", zap.Any("project", zap.Any("c)", c)))
	return c, nil
}

func (u Usecase) GetRecentWorksProjects(req structure.FilterProjects) (*entity.Pagination, error) {

	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	pe.WalletAddress = req.WalletAddress
	projects, err := u.Repo.GetRecentWorksProjects(*pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetRecentWorksProjects", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info("projects", projects.Total)
	return projects, nil
}

func (u Usecase) GetUpdatedProjectStats(req structure.GetProjectReq) (*entity.ProjectStat, []entity.TraitStat, error) {

	project, err := u.Repo.FindProjectByProjectIdWithoutCache(req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	// do not resync
	if project.Stats.LastTimeSynced != nil && project.Stats.LastTimeSynced.Unix()+int64(u.Config.TimeResyncProjectStat) > time.Now().Unix() {
		return &project.Stats, project.TraitsStat, nil
	}

	allTokenFromDb, err := u.Repo.GetAllTokensByProjectID(project.TokenID)
	if err != nil {
		return nil, nil, err
	}
	owners := make(map[string]bool)
	for _, token := range allTokenFromDb {
		owners[token.OwnerAddr] = true
	}

	var allListings []entity.MarketplaceListings
	var allOffers []entity.MarketplaceOffers

	allListings, err = u.Repo.GetAllListingByCollectionContract(project.GenNFTAddr)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetAllListingByCollectionContract", err.Error(), err)
		return nil, nil, err
	}

	allOffers, err = u.Repo.GetAllOfferByCollectionContract(project.GenNFTAddr)
	if err != nil {
		logger.AtLog.Error("u.Repo.GetAllOfferByCollectionContract", err.Error(), err)
		return nil, nil, err
	}

	var totalTradingVolumn *big.Int
	var floorPrice *big.Int
	var bestMakeOfferPrice *big.Int
	var listedPercent int32
	listingSet := make(map[string]bool)

	for _, listing := range allListings {
		if listing.Erc20Token != utils.EVM_NULL_ADDRESS {
			continue
		}
		price := new(big.Int)
		price, ok := price.SetString(listing.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			logger.AtLog.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(listing.DurationTime, 10, 64)
		if err != nil {
			logger.AtLog.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if listing.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}
		// update listing percent
		if !listing.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			listingSet[listing.TokenId] = true
		}

		// update floor price
		if listing.Finished {
			if floorPrice == nil {
				floorPrice = price
			} else {
				if floorPrice.Cmp(price) > 0 {
					floorPrice = price
				}
			}
		}
	}

	for _, offer := range allOffers {
		price := new(big.Int)
		price, ok := price.SetString(offer.Price, 10)
		if !ok {
			err := errors.New("fail to convert price to big int")
			logger.AtLog.Error("fail to convert price to big int", err.Error(), err)
			continue
		}
		durationTime, err := strconv.ParseInt(offer.DurationTime, 10, 64)
		if err != nil {
			logger.AtLog.Error("fail to parse duration time", err.Error(), err)
			continue
		}

		// update total volumn trading
		if offer.Finished {
			if totalTradingVolumn == nil {
				totalTradingVolumn = new(big.Int)
			}
			totalTradingVolumn.Add(totalTradingVolumn, price)
		}

		// update floor price
		if !offer.Closed && (time.Now().Unix() < durationTime || durationTime == 0) {
			if bestMakeOfferPrice == nil {
				bestMakeOfferPrice = price
			} else {
				if bestMakeOfferPrice.Cmp(price) < 0 {
					bestMakeOfferPrice = price
				}
			}
		}

		// update floor price
		if offer.Finished {
			if floorPrice == nil {
				floorPrice = price
			} else {
				if floorPrice.Cmp(price) > 0 {
					floorPrice = price
				}
			}
		}
	}

	if len(allTokenFromDb) > 0 {
		listedPercent = int32(len(listingSet) * 100 / len(allTokenFromDb))
	} else {
		listedPercent = 0
	}

	if totalTradingVolumn == nil {
		totalTradingVolumn = new(big.Int)
	}
	if floorPrice == nil {
		floorPrice = new(big.Int)
	}
	if bestMakeOfferPrice == nil {
		bestMakeOfferPrice = new(big.Int)
	}

	// update trait stats
	traitToCnt := make(map[string]int32)
	traitValueToCnt := make(map[string]map[string]int32)
	for _, token := range allTokenFromDb {
		for _, attribute := range token.ParsedAttributes {
			traitToCnt[attribute.TraitType] += 1
			if traitValueToCnt[attribute.TraitType] == nil {
				traitValueToCnt[attribute.TraitType] = make(map[string]int32)
			}
			traitValueToCnt[attribute.TraitType][fmt.Sprintf("%v", attribute.Value)] += 1
		}
	}

	traitsStat := make([]entity.TraitStat, 0)
	for k, cnt := range traitToCnt {
		traitValueStat := make([]entity.TraitValueStat, 0)
		for value, cntValue := range traitValueToCnt[k] {
			traitValueStat = append(traitValueStat, entity.TraitValueStat{
				Value:  value,
				Rarity: int32(cntValue * 100 / cnt),
			})
		}
		traitsStat = append(traitsStat, entity.TraitStat{
			TraitName:       k,
			TraitValuesStat: traitValueStat,
		})
	}

	now := time.Now()

	project, err = u.Repo.FindProjectBy(req.ContractAddr, req.TokenID)
	if err != nil {
		return nil, nil, err
	}

	return &entity.ProjectStat{
		LastTimeSynced:     &now,
		UniqueOwnerCount:   uint32(len(owners)),
		TotalTradingVolumn: totalTradingVolumn.String(),
		FloorPrice:         floorPrice.String(),
		BestMakeOfferPrice: bestMakeOfferPrice.String(),
		ListedPercent:      listedPercent,
		MintedCount:        project.Stats.MintedCount,
	}, traitsStat, nil
}

func (u Usecase) getProjectDetailFromChainWithoutCache(req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {

	contractDataKey := fmt.Sprintf("detail.%s.%s", req.ContractAddress, req.ProjectID)

	logger.AtLog.Info("req", req)

	addr := common.HexToAddress(req.ContractAddress)
	// call to contract to get emotion
	client, err := helpers.TCDialer()
	if err != nil {
		logger.AtLog.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	projectID := new(big.Int)
	projectID, ok := projectID.SetString(req.ProjectID, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
	if err != nil {
		logger.AtLog.Error("u.getNftContractDetailInternal", err.Error(), err)
		return nil, err
	}
	//logger.AtLog.Info("contractDetail", contractDetail)
	u.Cache.SetData(contractDataKey, contractDetail)
	return contractDetail, nil
}

// Get from chain with cache
func (u Usecase) getProjectDetailFromChain(req structure.GetProjectDetailMessageReq) (*structure.ProjectDetail, error) {

	contractDataKey := helpers.ProjectDetailKey(req.ContractAddress, req.ProjectID)

	//u.Cache.Delete(contractDataKey)
	data, err := u.Cache.GetData(contractDataKey)
	if err != nil {
		logger.AtLog.Info("req", req)

		addr := common.HexToAddress(req.ContractAddress)
		// call to contract to get emotion
		client, err := helpers.TCDialer()
		if err != nil {
			logger.AtLog.Error("ethclient.Dial", err.Error(), err)
			return nil, err
		}

		projectID := new(big.Int)
		projectID, ok := projectID.SetString(req.ProjectID, 10)
		if !ok {
			return nil, errors.New("cannot convert tokenID")
		}
		contractDetail, err := u.getNftContractDetailInternal(client, addr, *projectID)
		if err != nil {
			logger.AtLog.Error("u.getNftContractDetail", err.Error(), err)
			return nil, err
		}
		logger.AtLog.Info("contractDetail", contractDetail)
		u.Cache.SetData(contractDataKey, contractDetail)
		return contractDetail, nil
	}

	contractDetail := &structure.ProjectDetail{}
	err = helpers.ParseCache(data, contractDetail)
	if err != nil {
		logger.AtLog.Error("helpers.ParseCache", err.Error(), err)
		return nil, err
	}

	return contractDetail, nil
}

// Internal get project detail
func (u Usecase) getNftContractDetailInternal(client *ethclient.Client, contractAddr common.Address, projectID big.Int) (*structure.ProjectDetail, error) {

	gProject, err := generative_project_contract.NewGenerativeProjectContract(contractAddr, client)
	if err != nil {
		logger.AtLog.Error("generative_project_contract.NewGenerativeProjectContract", err.Error(), err)
		return nil, err
	}

	pDchan := make(chan structure.ProjectDetailChan, 1)
	pStatuschan := make(chan structure.ProjectStatusChan, 1)
	pTokenURIchan := make(chan structure.ProjectNftTokenUriChan, 1)

	go func(pDchan chan structure.ProjectDetailChan, projectID *big.Int) {
		proDetail := &generative_project_contract.NFTProjectProject{}
		var err error

		defer func() {
			pDchan <- structure.ProjectDetailChan{
				ProjectDetail: proDetail,
				Err:           err,
			}
		}()

		proDetailReps, err := gProject.ProjectDetails(nil, projectID)
		if err != nil {
			return
		}

		proDetail = &proDetailReps

	}(pDchan, &projectID)

	go func(pDchan chan structure.ProjectStatusChan, projectID *big.Int) {
		var status *bool
		var err error

		defer func() {
			pDchan <- structure.ProjectStatusChan{
				Status: status,
				Err:    err,
			}
		}()

		pStatus, err := gProject.ProjectStatus(nil, projectID)
		if err != nil {
			return
		}

		status = &pStatus

	}(pStatuschan, &projectID)

	go func(pDchan chan structure.ProjectNftTokenUriChan, projectID *big.Int) {
		var tokenURI *string
		var err error

		defer func() {
			pDchan <- structure.ProjectNftTokenUriChan{
				TokenURI: tokenURI,
				Err:      err,
			}
		}()

		pTokenUri, err := gProject.TokenURI(nil, projectID)
		if err != nil {
			return
		}

		tokenURI = &pTokenUri

	}(pTokenURIchan, &projectID)

	detailFromChain := <-pDchan
	statusFromChain := <-pStatuschan
	tokenFromChain := <-pTokenURIchan

	if detailFromChain.Err != nil {
		return nil, detailFromChain.Err
	}

	if statusFromChain.Err != nil {
		logger.AtLog.Error("statusFromChain.Err", statusFromChain.Err.Error(), statusFromChain.Err)
		return nil, statusFromChain.Err
	}

	if tokenFromChain.Err != nil {
		logger.AtLog.Error("tokenFromChain.Err", tokenFromChain.Err.Error(), tokenFromChain.Err)
		return nil, tokenFromChain.Err
	}

	gNftProject, err := generative_nft_contract.NewGenerativeNftContract(detailFromChain.ProjectDetail.GenNFTAddr, client)
	if err != nil {
		logger.AtLog.Error("generative_nft_contract.NewGenerativeNftContract", err.Error(), err)
		return nil, err
	}

	//nft project detail chain
	nftProjectDchan := make(chan structure.NftProjectDetailChan, 1)
	go func(nftProjectDchan chan structure.NftProjectDetailChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		data := &structure.NftProjectDetail{}
		var err error

		defer func() {
			nftProjectDchan <- structure.NftProjectDetailChan{
				Data: data,
				Err:  err,
			}
		}()

		respData, err := gNftProject.Project(nil)
		err = copier.Copy(data, respData)

	}(nftProjectDchan, gNftProject)

	nftRoyaltychan := make(chan structure.RoyaltyChan, 1)
	go func(nftRoyaltychan chan structure.RoyaltyChan, gNftProject *generative_nft_contract.GenerativeNftContract) {
		var data *big.Int
		var err error

		defer func() {
			nftRoyaltychan <- structure.RoyaltyChan{
				Data: data,
				Err:  err,
			}
		}()

		data, err = gNftProject.Royalty(nil)

	}(nftRoyaltychan, gNftProject)

	dataFromNftPChan := <-nftProjectDchan
	dataFromRoyaltyPChan := <-nftRoyaltychan

	resp := &structure.ProjectDetail{
		ProjectDetail: detailFromChain.ProjectDetail,
		Status:        *statusFromChain.Status,
		NftTokenUri:   *tokenFromChain.TokenURI,
	}

	logger.AtLog.Info("resp", resp)
	if dataFromNftPChan.Err == nil && dataFromNftPChan.Data != nil {
		resp.NftProjectDetail = *dataFromNftPChan.Data
	} else {
		resp.NftProjectDetail = structure.NftProjectDetail{}
	}

	if dataFromRoyaltyPChan.Err == nil && dataFromRoyaltyPChan.Data != nil {
		resp.Royalty = structure.ProjectRoyalty{
			Data: *dataFromRoyaltyPChan.Data,
		}
	}

	logger.AtLog.Info("resp", resp)
	return resp, nil
}

func (u Usecase) UnzipProjectFile(zipPayload *structure.ProjectUnzipPayload) (*entity.Projects, error) {
	var err error
	pe := &entity.Projects{}
	zipLink := zipPayload.ZipLink

	defer func() {
		now := time.Now().UTC()
		status := entity.UzipStatusSuccess
		message := ""
		if err == nil {
			u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Project images are Unzipped][project %s]", helpers.CreateProjectLink(pe.TokenID, pe.Name)), "", fmt.Sprintf("Project's images have been unzipped with %d files, zipLink: %s", len(pe.Images), helpers.CreateTokenImageLink(zipLink)))
		} else {
			status = entity.UzipStatusFail
			message = err.Error()
			u.NotifyWithChannel(os.Getenv("SLACK_PROJECT_CHANNEL_ID"), fmt.Sprintf("[Error while unzip][project %s]", helpers.CreateProjectLink(pe.TokenID, pe.Name)), "", fmt.Sprintf("Project's images have been unzipped with %d files, zipLink: %s, error: %s", len(pe.Images), helpers.CreateTokenImageLink(zipLink), message))
		}

		up, err := u.Repo.GetProjectUnzip(zipPayload.ProjectID)
		if err == nil && up != nil {

			up.Logs = append(up.Logs, entity.ProjectZipLinkLog{
				Message:     message,
				Status:      status,
				CreatedTime: &now,
			})

			up.Status = status
			up.Message = message
			up.ReTries = up.ReTries + 1
			updated, err := u.Repo.UpdateProjectUnzip(zipPayload.ProjectID, up)
			if err != nil {
				logger.AtLog.Error("UnzipProjectFile.defer", zap.Any("projectID", zipPayload.ProjectID), zap.Error(err))
			}
			logger.AtLog.Logger.Info("UnzipProjectFile.defer", zap.Any("projectID", zipPayload.ProjectID), zap.Any("updated", updated))

		} else {
			unzipLog := &entity.ProjectZipLinks{
				ProjectID: zipPayload.ProjectID,
				ZipLink:   zipPayload.ZipLink,
				Status:    status,
				Message:   message,
				ReTries:   1,
			}

			unzipLog.Logs = []entity.ProjectZipLinkLog{}
			unzipLog.Logs = append(unzipLog.Logs, entity.ProjectZipLinkLog{
				Message:     message,
				Status:      status,
				CreatedTime: &now,
			})

			err = u.Repo.CreateProjectUnzip(unzipLog)
			if err != nil {
				logger.AtLog.Error("UnzipProjectFile.defer", zap.Any("projectID", zipPayload.ProjectID), zap.Error(err))
			}
			logger.AtLog.Logger.Info("UnzipProjectFile.defer", zap.Any("projectID", zipPayload.ProjectID), zap.Bool("created", true))
		}
	}()

	pe, err = u.Repo.FindProjectByTokenID(zipPayload.ProjectID)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), err.Error(), zap.Error(err), zap.String("projectID", pe.TokenID))
		return nil, err
	}

	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = pe.Name
	nftTokenURI["description"] = pe.Description
	nftTokenURI["image"] = pe.Thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}
	logger.AtLog.Logger.Info(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("zipPayload", zipPayload), zap.String("projectID", pe.TokenID))

	images := []string{}

	//spew.Dump(os.Getenv("GCS_DOMAIN"))
	groupIndex := strings.Index(zipLink, "btc-projects/")
	strLen := len(zipLink)
	zipLink = zipLink[groupIndex:strLen]
	//spew.Dump(zipLink)
	err = u.GCS.UnzipFile(zipLink)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("UnzipFile", zipLink), zap.String("projectID", pe.TokenID), zap.Error(err))
		return nil, err
	}

	unzipFoler := zipLink + "_unzip"
	files, err := u.GCS.ReadFolder(unzipFoler)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("ReadFolder", unzipFoler), zap.String("projectID", pe.TokenID), zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("ReadFolder", unzipFoler), zap.String("projectID", pe.TokenID), zap.Error(err))
	maxSize := uint64(0)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	for _, f := range files {
		if strings.Index(strings.ToLower(f.Name), strings.ToLower("__MACOSX")) > -1 {
			continue
		}
		if strings.Index(strings.ToLower(f.Name), strings.ToLower(".DS_Store")) > -1 {
			continue
		}

		temp := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), f.Name)
		images = append(images, temp)
		nftTokenURI["image"] = temp
		if uint64(f.Size) > maxSize {
			maxSize = uint64(f.Size)
		}
	}
	//

	logger.AtLog.Logger.Info(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("zipPayload", zipPayload), zap.Any("projecID", pe.TokenID), zap.Int("images", len(pe.Images)))
	pe.Images = images
	if len(images) > 0 {
		pe.IsFullChain = true
		if len(pe.Images) == 1 {
			// edition
			if pe.MaxSupply > 1 {
				//-> clone for maxsupply - 1 files
				for i := 1; i < int(pe.MaxSupply); i++ {
					pe.Images = append(pe.Images, pe.Images[0])
				}
			}
		} else {
			// list file
			if len(pe.Images) < int(pe.MaxSupply) {
				// max supply need to equal max files
				pe.MaxSupply = int64(len(pe.Images))
			}
		}

	}
	pe.IsHidden = true
	pe.Status = true
	pe.IsSynced = true

	networkFee := big.NewInt(u.networkFeeBySize(int64(maxSize / 4))) // will update after unzip and check data
	pe.MaxFileSize = int64(maxSize)
	pe.NetworkFee = networkFee.String()

	updated, err := u.Repo.UpdateProject(pe.UUID, pe)
	if err != nil {
		logger.AtLog.Error(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("ReadFolder", unzipFoler), zap.String("projectID", pe.TokenID), zap.Error(err))
		return nil, err
	}

	ids, err1 := u.CreateDAOProject(context.TODO(), &request.CreateDaoProjectRequest{
		ProjectIds: []string{pe.ID.Hex()},
		CreatedBy:  pe.CreatorAddrr,
	})
	if err1 != nil {
		logger.AtLog.Logger.Error("CreateDAOProject failed", zap.String("projectID", pe.TokenID), zap.Error(err1))
	} else {
		logger.AtLog.Logger.Info("CreateDAOProject success",
			zap.String("project_id", pe.ID.Hex()),
			zap.Strings("ids", ids),
		)
	}

	logger.AtLog.Logger.Info(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("zipPayload", zipPayload), zap.Any("updated", updated), zap.Any("projectID", pe.TokenID), zap.Int("images", len(images)))
	go func() {
		owner, err := u.Repo.FindUserByWalletAddress(pe.CreatorAddrr)
		if err != nil {
			logger.AtLog.Error("UnzipProjectFile.FindUserByWalletAddress failed", zap.Error(err))
			return
		}
		if len(ids) > 0 {
			u.NotifyNewProject(pe, owner, true, ids[0])
		}
		u.AirdropArtist(pe.TokenID, os.Getenv("AIRDROP_WALLET"), *owner, 3)
	}()

	logger.AtLog.Logger.Info(fmt.Sprintf("UnzipProjectFile.%s", pe.TokenID), zap.Any("updated", updated), zap.String("projectID", pe.TokenID))
	return pe, nil
}

func (u Usecase) UploadFileZip(fc []byte, uploadChan chan uploadFileChan, peName string, fileName string, wg *sync.WaitGroup) {

	var err error
	var uploadedUrl *string

	defer func() {
		uploadChan <- uploadFileChan{
			FileURL: uploadedUrl,
			Err:     err,
		}
		wg.Done()
	}()

	base64Data := helpers.Base64Encode(fc)

	key := helpers.GenerateSlug(peName)
	key = fmt.Sprintf("btc-projects/%s/unzip", key)

	uploadFileName := fmt.Sprintf("%s/%s", key, fileName)
	uploaded, err := u.GCS.UploadBaseToBucket(base64Data, uploadFileName)
	if err != nil {
		logger.AtLog.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
		return
	}

	logger.AtLog.Info("uploaded", uploaded)
	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	uploadedUrl = &cdnURL

}

func (u Usecase) CreateProjectFromCollectionMeta(meta entity.CollectionMeta) (*entity.Projects, error) {
	logger.AtLog.Info(fmt.Sprintf("Start create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))
	pe := &entity.Projects{}

	mPrice := helpers.StringToBTCAmount("0")

	thumbnail := fmt.Sprintf("https://generativeexplorer.com/content/%s", meta.InscriptionIcon)

	pe.ContractAddress = os.Getenv("GENERATIVE_BTC_PROJECT")
	pe.MintPrice = mPrice
	pe.NetworkFee = big.NewInt(u.networkFeeBySize(int64(300000 / 4))).String() // will update after unzip and check data or check from animation url
	pe.IsHidden = false
	pe.Status = false
	pe.IsSynced = true
	nftTokenURI := make(map[string]interface{})
	nftTokenURI["name"] = meta.Name
	nftTokenURI["description"] = meta.Description
	nftTokenURI["image"] = thumbnail
	nftTokenURI["animation_url"] = ""
	nftTokenURI["attributes"] = []string{}

	pe.CreatorAddrr = "0x0000000000000000000000000000000000000000"
	if meta.WalletAddress != "" {
		pe.CreatorAddrr = meta.WalletAddress
	}
	creatorAddrr, err := u.Repo.FindUserByWalletAddress(pe.CreatorAddrr)
	if err != nil {
		logger.AtLog.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
		pe.CreatorAddrr = "0x0000000000000000000000000000000000000000"
		creatorAddrr, err = u.Repo.FindUserByWalletAddress(pe.CreatorAddrr)
		if err != nil {
			logger.AtLog.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
			return nil, err
		}
	}

	pe.CreatorName = creatorAddrr.DisplayName

	bytes, err := json.Marshal(nftTokenURI)
	if err != nil {
		logger.AtLog.Error("json.Marshal.nftTokenURI", err.Error(), err)
		return nil, err
	}
	nftToken := helpers.Base64Encode(bytes)
	now := time.Now().UTC()

	pe.NftTokenUri = fmt.Sprintf("data:application/json;base64,%s", nftToken)
	pe.ProcessingImages = []string{}
	pe.MintedImages = nil
	pe.MintedTime = &now
	pe.CreatorProfile = *creatorAddrr
	pe.CreatorAddrrBTC = creatorAddrr.WalletAddressBTC
	pe.LimitSupply = 0
	pe.GenNFTAddr = pe.TokenID

	pe.Name = meta.Name
	pe.Description = meta.Description
	maxSupply, err := strconv.ParseInt(meta.Supply, 10, 64)
	if err != nil {
		maxSupply = 0
	}
	pe.MaxSupply = maxSupply
	pe.MintingInfo.Index = 0

	if pe.Categories == nil || len(pe.Categories) == 0 {
		pe.Categories = []string{u.Config.UnverifiedCategoryID}
	}

	royalty, err := strconv.Atoi(meta.Royalty)
	if err != nil {
		royalty = 0
	}
	pe.Royalty = royalty
	pe.SocialTwitter = meta.TwitterLink
	pe.SocialDiscord = meta.DiscordLink
	pe.SocialWeb = meta.WebsiteLink
	pe.Thumbnail = thumbnail

	maxID, err := u.Repo.GetMaxBtcProjectID()
	if err != nil {
		logger.AtLog.Error("u.Repo.GetMaxBtcProjectID", err.Error(), err)
		return nil, err
	}
	maxID = maxID + 1
	pe.TokenIDInt = maxID
	pe.TokenID = fmt.Sprintf("%d", maxID)
	pe.GenNFTAddr = pe.TokenID
	pe.InscriptionIcon = meta.InscriptionIcon
	pe.CreatedByCollectionMeta = true
	blockNumberMinted := "0"
	pe.BlockNumberMinted = &blockNumberMinted
	pe.Source = meta.Source + "$$" + meta.From

	err = u.Repo.CreateProject(pe)
	if err != nil {
		logger.AtLog.Error("u.Repo.CreateProjectFromInscription", err.Error(), err)
		return nil, err
	}

	logger.AtLog.Info(fmt.Sprintf("Done create project from collection meta %s %s", meta.Name, meta.InscriptionIcon))

	return pe, nil
}

type Volumes struct {
	Items    []Volume `json:"items"`
	TotalBTC float64  `json:"totalAmountBTC"`
	TotalETH float64  `json:"totalAmountETH"`
}

type Volume struct {
	ProjectID string `json:"projectID"`
	PayType   string `json:"payType"`
	Amount    string `json:"amount"`
	Earning   string `json:"earning"`
	Withdraw  string `json:"withdraw"`
	Available string `json:"available"`
	Status    int    `json:"status"`
}

func (u Usecase) CreatorVolume(creatoreAddress string, paytype string) (*Volume, error) {
	data, err := u.GetVolumeOfUser(creatoreAddress, &paytype)
	if err != nil {
		logger.AtLog.Error("CollectorVolume", zap.Any("err", err))
		return nil, err
	}

	tmp := Volume{
		ProjectID: data.ID.ProjectID,
		PayType:   data.ID.Paytype,
		Amount:    fmt.Sprintf("%d", int(data.Amount)),
	}

	return &tmp, nil
}

func (u Usecase) ProjectVolume(projectID string, paytype string) (*Volume, error) {
	data, err := u.GetVolumeOfProject(projectID, &paytype)
	if err != nil {
		logger.AtLog.Error("CollectorVolume", zap.Any("err", err))
		tmp := Volume{
			ProjectID: projectID,
			PayType:   paytype,
			Amount:    "0",
			Earning:   "0",
			Withdraw:  "0",
			Available: "0",
			Status:    entity.StatusWithdraw_Available,
		}

		return &tmp, nil
	}

	latestWd, err := u.Repo.GetLastWithdraw(entity.FilterWithdraw{
		WithdrawItemID: &projectID,
		PaymentType:    &paytype,
	})

	wdraw := 0.0
	w, err := u.Repo.AggregateWithDrawByUser(&entity.FilterWithdraw{
		WithdrawItemID: &projectID,
		PaymentType:    &paytype,
		Statuses: []int{
			entity.StatusWithdraw_Pending,
			entity.StatusWithdraw_Approve,
		},
	})

	status := entity.StatusWithdraw_Available
	if err == nil && len(w) > 0 {
		wdraw = w[0].Amount
	}

	//the status show int FE, that allows user can click withdraw button
	if latestWd != nil {
		status = latestWd.Status
		if status == entity.StatusWithdraw_Approve {
			status = entity.StatusWithdraw_Available
		}

		if status == entity.StatusWithdraw_Reject {
			status = entity.StatusWithdraw_Available
		}
	}

	available := data.Earning - wdraw
	if available < 0 {
		available = 0
	}
	tmp := Volume{
		ProjectID: data.ID.ProjectID,
		PayType:   data.ID.Paytype,
		Amount:    fmt.Sprintf("%d", int(data.Amount)),
		Earning:   fmt.Sprintf("%d", int(data.Earning)),
		Withdraw:  fmt.Sprintf("%d", int(wdraw)),
		Available: fmt.Sprintf("%d", int(available)),
		Status:    status,
	}

	logger.AtLog.Logger.Info("ProjectVolume ...", zap.String("projectID", projectID), zap.Any("volume", tmp), zap.Any("AggregateWithDrawByUser", w), zap.Any("latestWd", zap.Any("latestWd)", latestWd)))
	return &tmp, nil
}

func (u Usecase) CreateProjectsAndTokenUriFromInscribeAuthentic(ctx context.Context, item entity.InscribeBTC) error {
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(item.TokenAddress, item.TokenId)
	if err != nil {
		return err
	}
	project := &entity.Projects{}
	if err := u.Repo.FindOneBy(ctx, entity.Projects{}.TableName(), bson.M{
		"fromAuthentic": true,
		"tokenAddress":  item.TokenAddress,
	}, project); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		creator := &entity.Users{}
		if err := u.Repo.FindOneBy(ctx, entity.Users{}.TableName(),
			bson.M{"wallet_address": "0x1111111111111111111111111111111111111111"},
			creator); err != nil {
			return err
		}
		isFalse := false
		reqBtcProject := structure.CreateBtcProjectReq{
			Name:            nft.Name,
			MaxSupply:       1,
			Index:           1,
			CreatorName:     creator.DisplayName,
			CreatorAddrr:    creator.WalletAddress,
			CreatorAddrrBTC: creator.WalletAddressBTC,
			FromAuthentic:   true,
			TokenAddress:    item.TokenAddress,
			TokenId:         item.TokenId,
			OwnerOf:         item.OwnerOf,
			OrdinalsTx:      item.OrdinalsTx,
			Thumbnail:       item.FileURI,
			InscribedBy:     item.UserWalletAddress,
			IsHidden:        &isFalse,
		}
		if nft.MetadataString != nil && *nft.MetadataString != "" {
			metadata := &nfts.MoralisTokenMetadata{}
			if err := json.Unmarshal([]byte(*nft.MetadataString), metadata); err == nil {
				reqBtcProject.AnimationURL = &metadata.AnimationUrl
			}
		}
		reqBtcProject.Categories = []string{u.Config.EthereumCategoryID}
		project, err = u.CreateBTCProject(reqBtcProject)
		if err != nil {
			return err
		}
	} else {
		maxSupply := project.MaxSupply + 1
		index := project.MintingInfo.Index + 1
		project, err = u.UpdateBTCProject(structure.UpdateBTCProjectReq{
			CreatetorAddress: &project.CreatorAddrr,
			ProjectID:        &project.TokenID,
			MaxSupply:        &maxSupply,
			Index:            &index,
		})
		if err != nil {
			return err
		}
	}

	_, err = u.CreateBTCTokenURI(item.OriginUserAddress, project.TokenID, item.InscriptionID, item.FileURI, entity.BIT, item.TokenId, item.UserWalletAddress)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) ProjectRandomImages(projectID string) ([]string, error) {
	max := 10
	p, err := u.Repo.FindProjectByTokenID(projectID)
	if err != nil {
		return nil, err
	}
	totalImages := len(p.Images)
	totalProcessingImages := len(p.ProcessingImages)

	if totalImages == 0 && totalProcessingImages == 0 {
		return nil, errors.New("Project doesn's have any images")
	}

	returnImages := []string{}
	for _, item := range p.Images {
		if len(returnImages) >= max {
			break
		}
		returnImages = append(returnImages, item)
	}

	for _, item := range p.ProcessingImages {
		if len(returnImages) >= max {
			break
		}
		returnImages = append(returnImages, item)
	}

	return returnImages, nil

}

func (u Usecase) ProjectTokenTraits(projectID string) ([]structure.TokenTraits, error) {
	resp := []structure.TokenTraits{}
	tokens, err := u.Repo.GetAllTokenTraitsByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		attrs := []structure.TraitAttribute{}
		tmp := structure.TokenTraits{}
		tmp.ID = token.TokenID
		tmp.Atrributes = attrs

		for _, attr := range token.ParsedAttributesStr {
			attrsTmp := structure.TraitAttribute{
				TraitType: attr.TraitType,
				Value:     attr.Value,
			}

			attrs = append(attrs, attrsTmp)
		}

		tmp.Atrributes = attrs
		resp = append(resp, tmp)
	}
	return resp, nil
}

func (u Usecase) UploadTokenTraits(projectID string, r *http.Request) (*entity.TokenUriMetadata, error) {
	p, err := u.Repo.FindProjectByTokenID(projectID)
	if err != nil {
		logger.AtLog.Errorf("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	/*totalImages := len(p.Images)
	totalProcessingImages := len(p.ProcessingImages)
	if totalImages == 0 && totalProcessingImages == 0 {
		err = errors.New("Project doesn's have any files")
		logger.AtLog.Error(zap.String("projectID", projectID), err.Error())
		return nil, err
	}*/
	isGenerative := true
	if p.Source == "" {
		if len(p.Images)+len(p.ProcessingImages) > 0 {
			// -> from generative.xyz
			if len(p.Images) > 0 {
				if !strings.HasSuffix(p.Images[0], ".html") {
					isGenerative = false
				}
			}

			if isGenerative && len(p.ProcessingImages) > 0 {
				if !strings.HasSuffix(p.ProcessingImages[0], ".html") {
					isGenerative = false
				}
			}
		}
	} else {
		// crawler
		isGenerative = false
	}
	if isGenerative {
		err = errors.New("Project is generative")
		logger.AtLog.Error(zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	_, handler, err := r.FormFile("file")
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	key := helpers.GenerateSlug(projectID)
	key = fmt.Sprintf("btc-projects/%s/json", key)
	gf := googlecloud.GcsFile{
		FileHeader: handler,
		Path:       &key,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	content, err := u.GCS.ReadFile(uploaded.Name)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	data := []entity.TokenTraits{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	h := &entity.TokenUriMetadata{
		ProjectID:    projectID,
		UploadedFile: uploaded.FullPath,
		Content:      data,
	}

	err = u.Repo.CreateTokenUriMetadata(h)
	if err != nil {
		logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
		return nil, err
	}

	for _, item := range data {
		tokenID := item.ID
		token, err := u.Repo.FindTokenByTokenID(tokenID)
		if err != nil {
			err = fmt.Errorf("token %s was not found: %v", tokenID, err)
			logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
			return nil, err
		}

		if token.ProjectID != p.TokenID {
			err = fmt.Errorf("token %s is not belong to this project %s", tokenID, projectID)
			logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
			return nil, err
		}

		attrs := []entity.TokenUriAttr{}
		attrStrs := []entity.TokenUriAttrStr{}

		for _, itemAttr := range item.Attributes {
			attr := entity.TokenUriAttr{
				TraitType: itemAttr.TraitType,
				Value:     itemAttr.Value,
			}

			attrStr := entity.TokenUriAttrStr{
				TraitType: itemAttr.TraitType,
				Value:     itemAttr.Value,
			}

			attrs = append(attrs, attr)
			attrStrs = append(attrStrs, attrStr)
		}

		token.ParsedAttributes = attrs
		token.ParsedAttributesStr = attrStrs
		if len(strings.TrimSpace(item.Name)) > 0 && item.Name != token.Name {
			token.Name = item.Name
		}

		//spew.Dump(token.TokenID, token.ParsedAttributes)
		_, err = u.Repo.UpdateOrInsertTokenUri(token.ContractAddress, tokenID, token)
		if err != nil {
			err = fmt.Errorf("Cannot update token %s - %v", tokenID, err)
			logger.AtLog.Error("UploadTokenTraits", zap.String("projectID", projectID), err.Error())
			return nil, err
		}
	}

	return h, nil
}

func (u Usecase) GetProjectFirstSale(genNFTAddr string) string {
	totalAmount := "0"

	//u.Cache.Delete(helpers.ProjectFirstSaleKey(genNFTAddr))
	cached, err := u.Cache.GetData(helpers.ProjectFirstSaleKey(genNFTAddr))
	if err != nil || cached == nil {
		newAmount := 0.0
		data, err := u.Repo.AggregateBTCVolumn(genNFTAddr)
		if err == nil && data != nil {
			if len(data) > 0 {
				newAmount = data[0].Amount
			}
		}

		oldAmount := 0.0
		paytypes := []string{string(entity.BIT), string(entity.ETH)}
		for _, paytype := range paytypes {
			oldBTCData, err := u.AggregateOldData(genNFTAddr, paytype)
			if err != nil {
				continue
			}

			amount := oldBTCData.Amount
			if paytype == string(entity.ETH) {
				amount = amount / float64(oldBTCData.BtcRate/oldBTCData.EthRate)
			}

			oldAmount += amount
		}

		total := newAmount + oldAmount
		totalAmount = fmt.Sprintf("%d", int(total))
		u.Cache.SetStringData(helpers.ProjectFirstSaleKey(genNFTAddr), totalAmount)
		return totalAmount
	}
	return *cached

}

func (u Usecase) GetProjectsFloorPrice(projects []string) (map[string]uint64, error) {
	result := make(map[string]uint64)
	data, err := u.Repo.AggregateProjectsFloorPrice(projects)
	if err != nil {
		return nil, err
	}
	for _, v := range data {
		result[v.ID] = v.Floor
	}
	return result, nil
}

func (u Usecase) UpdateProjectHash(req structure.UpdateProjectHash) (*entity.Projects, error) {
	p, err := u.Repo.FindProjectByTxHash(*req.TxHash)
	if err != nil {
		logger.AtLog.Error("UpdateProjectHash", zap.Any("err.FindProjectBy", err))
		return nil, err
	}

	if req.CommitTxHash != nil {
		p.CommitTxHash = *req.CommitTxHash
	}

	if req.RevealTxHash != nil {
		p.RevealTxHash = *req.RevealTxHash
	}

	updated, err := u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Error("UpdateProject", zap.Any("err.UpdateProject", err))
		return nil, err
	}

	logger.AtLog.Logger.Info("UpdateProject", zap.Any("project", p), zap.Any("updated", updated))
	return p, nil
}
