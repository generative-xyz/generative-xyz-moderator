package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

// for api create a new mint:

func (u Usecase) CreateMintReceiveAddress(input structure.MintNftBtcData) (*entity.MintNftBtc, error) {

	if len(input.ProjectID) == 0 || len(input.WalletAddress) == 0 || len(input.RefundUserAddress) == 0 {
		return nil, errors.New("data invalid")
	}

	walletAddress := &entity.MintNftBtc{}

	receiveAddress := ""
	privateKey := ""
	feeSendMaster := big.NewInt(0)
	var err error

	if input.Quantity <= 0 {
		err = errors.New("quantity invalid")

		return nil, err
	}

	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.FindProjectByTokenID", zap.Error(err))
		return nil, errors.New("project not found")
	}

	// find Project and make sure index < max supply
	if p.MintingInfo.Index >= p.MaxSupply {
		err = fmt.Errorf("project %s is minted out", input.ProjectID)
		logger.AtLog.Logger.Error("projectIsMintedOut", zap.Error(err))
		return nil, err
	}

	if p.MintingInfo.Index+int64(input.Quantity) > p.MaxSupply {
		err = fmt.Errorf("not enough quantity %s", input.ProjectID)
		logger.AtLog.Logger.Error("projectIsMintedOut", zap.Error(err))
		return nil, err
	}

	// verify paytype:
	if input.PayType != utils.NETWORK_BTC && input.PayType != utils.NETWORK_ETH {
		err = errors.New("only support payType is eth or btc")
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Check(payType)", zap.Error(err))
		return nil, err
	}

	// check type:
	if input.PayType == utils.NETWORK_BTC {
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}

	} else if input.PayType == utils.NETWORK_ETH {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			logger.AtLog.Logger.Error("CreateMintReceiveAddress.ethClient.GenerateAddress", zap.Error(err))
			return nil, err
		}
	}

	if len(receiveAddress) == 0 || len(privateKey) == 0 {
		err = errors.New("can not create the wallet")
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.GenerateAddress", zap.Error(err))
		return nil, err
	}

	// set temp wallet info:
	walletAddress.PayType = input.PayType

	if len(os.Getenv("SECRET_KEY")) == 0 {
		err = errors.New("please config SECRET_KEY")
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.GenerateAddress", zap.Error(err))
		return nil, err
	}

	privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
		return nil, err
	}

	fmt.Println("input.FeeRate", input.FeeRate)

	walletAddress.FeeRate = int64(input.FeeRate)

	walletAddress.UserID = input.UserID
	walletAddress.UserAddress = input.UserAddress

	walletAddress.PrivateKey = privateKeyEnCrypt
	walletAddress.ReceiveAddress = receiveAddress
	walletAddress.RefundUserAdress = input.RefundUserAddress

	walletAddress.EstMintFeeInfoFe = input.EstMintFeeInfo
	walletAddress.IsCustomFeeRate = input.IsCustomFeeRate

	mintPrice, ok := big.NewInt(0).SetString(p.MintPrice, 10)
	if !ok {
		mintPrice = big.NewInt(0)
	}

	// check discount now:

	discountFlagReserve := false
	discountFlagBidWinner := false
	discountLimit := -1
	discountMintPrice := big.NewInt(0)

	// check discount whitelist price in project.reservers to get reserveMintPrice
	if len(p.Reservers) > 0 {

		for _, address := range p.Reservers {
			if strings.EqualFold(address, walletAddress.UserAddress) {

				discountMintPrice, ok = big.NewInt(0).SetString(p.ReserveMintPrice, 10)
				if ok {
					discountFlagReserve = true
					discountLimit = p.ReserveMintLimit
				}
				break
			}
		}

	}
	// check discount in bidder list:
	bidProjectID := utils.BidProjectIDProd
	if u.Config.ENV == "develop" {
		bidProjectID = utils.BidProjectIDDev
	}

	// return list winner:
	if strings.EqualFold(p.TokenID, bidProjectID) {
		auctionWinnerList, _ := u.GetAuctionListWinnerAddressFromConfig()
		if auctionWinnerList != nil && len(auctionWinnerList) > 0 {
			for _, auctionWinner := range auctionWinnerList {
				if strings.EqualFold(auctionWinner.Address, walletAddress.UserAddress) {
					discountFlagBidWinner = true
					discountMintPrice = big.NewInt(0) // mint free now, set it = auctionWinner.MintPrice if you need
					discountLimit = auctionWinner.Quantity
					break
				}
			}
		}
		if !discountFlagBidWinner {
			// not allowed buy for this item:
			err = errors.New("You are not allowed to purchase this item.")
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.CheckBidBuy", zap.Error(err))
			return nil, err
		}
	}

	// have a good mint price:
	if discountFlagReserve || discountFlagBidWinner {
		// get list item mint:
		countMinted := 0
		mintReadyList, _ := u.Repo.GetLimitWhiteList(input.UserAddress, input.ProjectID)

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
		maxSlot := discountLimit - countMinted

		fmt.Println("discountLimit: ", discountLimit)
		fmt.Println("countMinted: ", countMinted)
		fmt.Println("maxSlot: ", maxSlot)

		if maxSlot > 0 {
			if input.Quantity > maxSlot {
				return nil, errors.New(fmt.Sprintf("You can mint up to %d items at the price of %.6f BTC.", maxSlot, float64(discountMintPrice.Int64())/1e8))
			}

			mintPrice = big.NewInt(discountMintPrice.Int64())
			walletAddress.IsDiscount = true
			logger.AtLog.Logger.Info("CreateMintReceiveAddress.walletAddress.IsDiscount", zap.Any("true", true))

		} else {
			if discountFlagBidWinner {
				return nil, errors.New(fmt.Sprintf("You can mint up to %d items at the price of %.6f BTC.", maxSlot, float64(discountMintPrice.Int64())/1e8))
			}
		}
	}

	// cal fee:
	// todo: cal fee for minting on TC:
	feeInfos, err := u.calMintFeeInfo(mintPrice.Int64(), p.MaxFileSize, int64(input.FeeRate), 0, 0)
	if err != nil {
		logger.AtLog.Logger.Error("u.calMintFeeInfo.Err", zap.Error(err))
		return nil, err
	}

	fmt.Println("feeInfos: ", feeInfos)

	walletAddress.Platform = utils.PLATFORM_ORDINAL
	if p.IsMintTC() {
		// is mint on TC network
		walletAddress.Platform = utils.PLATFORM_TC
	}

	walletAddress.ProjectNetworkFee = int(feeInfos["btc"].NetworkFeeBigInt.Int64()) // btc value
	walletAddress.ProjectMintPrice = int(feeInfos["btc"].MintPriceBigInt.Int64())   // btc value

	walletAddress.MintPriceByPayType = feeInfos[input.PayType].MintPrice   // 1 item
	walletAddress.NetworkFeeByPayType = feeInfos[input.PayType].NetworkFee // 1 item

	walletAddress.BtcRate = feeInfos[input.PayType].BtcPrice
	walletAddress.EthRate = feeInfos[input.PayType].EthPrice

	walletAddress.EstFeeInfo = feeInfos

	walletAddress.FeeSendMaster = feeSendMaster.String()

	logger.AtLog.Logger.Info("CreateMintReceiveAddress.receive", zap.Any("receiveAddress", receiveAddress))

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
	}
	if input.PayType == utils.NETWORK_ETH {
		expiredTime = 2 // just 1h for checking eth balance
	}

	walletAddress.Amount = feeInfos[input.PayType].TotalAmount
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.Status = entity.StatusMint_Pending
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * time.Duration(expiredTime))

	fmt.Println("feeInfos[eth].MintPriceBigIn", feeInfos["eth"].MintPriceBigInt)
	fmt.Println("feeInfos[btc].MintPriceBigIn", feeInfos["btc"].MintPriceBigInt)

	// for batch mint, update here:
	walletAddress.Quantity = input.Quantity
	walletAddress.MintPriceByPayTypeTotal = big.NewInt(0).Mul(feeInfos[input.PayType].MintPriceBigInt, big.NewInt(int64(input.Quantity))).String()   // total
	walletAddress.NetworkFeeByPayTypeTotal = big.NewInt(0).Mul(feeInfos[input.PayType].NetworkFeeBigInt, big.NewInt(int64(input.Quantity))).String() // total

	// update total amount:
	walletAddress.Amount = big.NewInt(0).Mul(feeInfos[input.PayType].TotalAmountBigInt, big.NewInt(int64(input.Quantity))).String()

	// insert now:
	err = u.Repo.InsertMintNftBtc(walletAddress)
	if err != nil {
		logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.InsertMintNftBtc", zap.Error(err))
		return nil, err
	}

	return walletAddress, nil
}

// api cancel mint
func (u Usecase) CancelMintNftBtc(wallet, uuid string) error {
	mintItem, _ := u.Repo.FindMintNftBtcByNftID(uuid)
	if mintItem == nil {
		return errors.New("item not found")
	}

	if int(mintItem.Status) == -1 {
		return nil
	}

	fmt.Println("mintItem.OriginUserAddress", mintItem.OriginUserAddress)

	if !strings.EqualFold(wallet, mintItem.OriginUserAddress) {

		if !strings.EqualFold(wallet, mintItem.UserAddress) {
			return errors.New("permission denied")
		}

	}
	if mintItem.Status != entity.StatusMint_Pending {
		return errors.New("Can not cancel this, the item is in progress.")
	}
	return u.Repo.UpdateCancelMintNftBtc(mintItem.UUID)
}

// api get list mint:
func (u Usecase) GetCurrentMintingByWalletAddress(address string) ([]structure.MintingInscription, error) {
	result := []structure.MintingInscription{}

	// listMintV2, err := u.Repo.ListMintNftBtcByStatusAndAddress(address, []entity.StatusMint{entity.StatusMint_Pending, entity.StatusMint_WaitingForConfirms, entity.StatusMint_ReceivedFund, entity.StatusMint_Minting, entity.StatusMint_Minted, entity.StatusMint_SendingNFTToUser, entity.StatusMint_NeedToRefund, entity.StatusMint_Refunding, entity.StatusMint_TxRefundFailed, entity.StatusMint_TxMintFailed})

	listMintV2, err := u.Repo.ListMintNftBtcByStatusAndUserAddress(address, []entity.StatusMint{entity.StatusMint_Pending, entity.StatusMint_WaitingForConfirms, entity.StatusMint_ReceivedFund, entity.StatusMint_Minting, entity.StatusMint_Minted, entity.StatusMint_SendingNFTToUser, entity.StatusMint_NeedToRefund, entity.StatusMint_Refunding, entity.StatusMint_TxRefundFailed, entity.StatusMint_TxMintFailed})
	if err != nil {
		go u.trackMintNftBtcHistory("", "GetCurrentMintingByWalletAddress", "", 0, "ListMintNftBtcByStatusAndAddress", err.Error(), true)
		return nil, err
	}

	for _, item := range listMintV2 {
		projectInfo, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			log.Println("FindProjectByTokenID", item.ProjectID)
			go u.trackMintNftBtcHistory(item.UUID, "GetCurrentMintingByWalletAddress", item.TableName(), item.Status, "FindProjectByTokenID", err.Error(), true)
			// return nil, err
			continue
		}
		creator, err := u.Repo.FindUserByAddress(projectInfo.CreatorAddrr)
		if err != nil {
			log.Println("InscriptionsByOutputs.FindUserByAddress", err)
			go u.trackMintNftBtcHistory(item.UUID, "GetCurrentMintingByWalletAddress", item.TableName(), item.Status, "FindUserByAddress", err.Error(), true)
		}

		status := ""
		if time.Since(item.ExpiredAt) >= 1*time.Second && item.Status == entity.StatusMint_Pending {
			continue
		}
		if (item.Status) == -1 {
			continue
		}
		switch item.Status {
		case entity.StatusMint_NeedToRefund, entity.StatusMint_TxRefundFailed:
			status = entity.StatusMintToText[entity.StatusMint_Refunding]
		case entity.StatusMint_TxMintFailed:
			status = entity.StatusMintToText[entity.StatusMint_Minting]
		default:
			status = entity.StatusMintToText[item.Status]
		}

		if item.PayType == "eth" {
			if (item.Status == entity.StatusMint_Refunding || item.Status == entity.StatusMint_NeedToRefund) && item.ProjectID == "1001311" {
				status = entity.StatusMintToText[entity.StatusMint_Refunded]
				continue
			}

			if item.Status == entity.StatusMint_Refunded && item.ProjectID != "1001311" {
				status = entity.StatusMintToText[entity.StatusMint_Refunding]
			}
		}

		minting := structure.MintingInscription{
			ID:            item.UUID,
			CreatedAt:     item.CreatedAt,
			Status:        status,
			StatusIndex:   int(item.Status),
			FileURI:       item.FileURI,
			ProjectID:     item.ProjectID,
			ProjectImage:  projectInfo.Thumbnail,
			ProjectName:   projectInfo.Name,
			InscriptionID: item.InscriptionID,
			IsCancel:      int(item.Status) == 0,
			Quantity:      item.Quantity,
		}
		if creator != nil {
			minting.ArtistID = creator.UUID
			minting.ArtistName = creator.DisplayName
		}

		if minting.StatusIndex != 0 {
			minting.Quantity = 1
		}

		result = append(result, minting)
	}

	return result, nil
}

// api get mint detail + step status:
func (u Usecase) GetDetalMintNftBtc(uuid string) (*structure.MintingInscription, error) {
	mintItem, _ := u.Repo.FindMintNftBtcByNftID(uuid)
	if mintItem == nil {
		return nil, errors.New("item not found")
	}

	projectInfo, _ := u.Repo.FindProjectByTokenID(mintItem.ProjectID)
	if projectInfo == nil {
		return nil, errors.New("item not found")
	}

	type statusprogressStruct struct {
		Message string `json:"message"`
		Status  bool   `json:"status"`
		Title   string `json:"title"`
		Tx      string `json:"tx"`
	}

	// fix for project 1001311
	if mintItem.PayType == "eth" {

		if (mintItem.Status == entity.StatusMint_Refunding || mintItem.Status == entity.StatusMint_NeedToRefund) && mintItem.ProjectID == "1001311" {
			mintItem.Status = entity.StatusMint_Refunded
			mintItem.Amount = "0"
			mintItem.ReceiveAddress = ""
		}

		if mintItem.Status == entity.StatusMint_Refunded && mintItem.ProjectID != "1001311" {
			mintItem.Status = entity.StatusMint_Refunding
		}
	}

	statusMap := make(map[string]statusprogressStruct)

	statusMap["1"] = statusprogressStruct{
		Title:  entity.StatusMintToText[entity.StatusMint_Pending],
		Status: int(mintItem.Status) > 0,
	}
	statusMap["2"] = statusprogressStruct{
		Title:  entity.StatusMintToText[entity.StatusMint_WaitingForConfirms],
		Status: int(mintItem.Status) > 1,
	}

	if mintItem.Status == entity.StatusMint_NeedToRefund || mintItem.Status == entity.StatusMint_Refunding || mintItem.Status == entity.StatusMint_Refunded || mintItem.Status == entity.StatusMint_TxRefundFailed {
		statusMap["3"] = statusprogressStruct{
			Title:   entity.StatusMintToText[entity.StatusMint_Refunding],
			Status:  mintItem.Status == entity.StatusMint_Refunding || mintItem.Status == entity.StatusMint_Refunded || mintItem.Status == entity.StatusMint_TxRefundFailed,
			Tx:      mintItem.TxRefund,
			Message: mintItem.ReasonRefund,
		}

		statusMap["4"] = statusprogressStruct{
			Title:   entity.StatusMintToText[entity.StatusMint_Refunded],
			Status:  mintItem.Status == entity.StatusMint_Refunded,
			Tx:      mintItem.TxRefund,
			Message: mintItem.ReasonRefund,
		}

	} else {

		statusMap["3"] = statusprogressStruct{
			Title:   entity.StatusMintToText[entity.StatusMint_Minting],
			Status:  mintItem.IsMinted || mintItem.Status == entity.StatusMint_Minting,
			Tx:      mintItem.TxMintNft,
			Message: mintItem.MintMessage,
		}

		message := ""
		if mintItem.Status == entity.StatusMint_Minting {
			message = "Waiting for minting confirmation."
		}
		statusMap["4"] = statusprogressStruct{
			Title:   entity.StatusMintToText[entity.StatusMint_Minted],
			Status:  mintItem.IsMinted,
			Tx:      mintItem.TxMintNft,
			Message: message,
		}

	}

	isCancel := int(mintItem.Status) == 0

	minting := &structure.MintingInscription{
		ID:            mintItem.UUID,
		CreatedAt:     mintItem.CreatedAt,
		Status:        entity.StatusMintToText[mintItem.Status],
		StatusIndex:   int(mintItem.Status),
		FileURI:       mintItem.FileURI,
		ProjectID:     mintItem.ProjectID,
		ProjectImage:  projectInfo.Thumbnail,
		ProjectName:   projectInfo.Name,
		InscriptionID: mintItem.InscriptionID,

		ReceiveAddress:    mintItem.ReceiveAddress,
		OriginUserAddress: mintItem.OriginUserAddress,
		IsCancel:          isCancel,
		TxMint:            mintItem.TxMintNft,
		TxSendNft:         mintItem.TxSendNft,

		Amount:  mintItem.Amount,
		PayType: mintItem.PayType,

		ProgressStatus: statusMap,

		UserID: mintItem.UserID,
	}
	return minting, nil
}

// JOBs mint begin:
// job 1: job check balance for list mint_nft_btc
func (u Usecase) JobMint_CheckBalance() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list mint pending to check balance:
	listPending, _ := u.Repo.ListMintNftBtcPending()
	if len(listPending) == 0 {
		// go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "ListMintNftBtcPending", "[]", false)
		return nil
	}

	// get list btc to check a Batch
	// var batchBTCBalance []string
	// for _, item := range listPending {
	// 	if item.PayType == utils.NETWORK_BTC {
	// 		batchBTCBalance = append(batchBTCBalance, item.ReceiveAddress)
	// 	}
	// }

	// isRateLimitErr := false
	// balanceMaps, err := bs.BTCGetAddrInfoMulti(batchBTCBalance)
	// if err != nil && strings.Contains(err.Error(), "rate_limit") {
	// 	isRateLimitErr = true
	// }

	for _, item := range listPending {

		if len(item.ReceiveAddress) == 0 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "the receive address empty", "", false)
			continue
		}

		time.Sleep(1 * time.Second)

		// check balance:
		balance := big.NewInt(0)
		confirm := -1

		if item.PayType == utils.NETWORK_BTC {

			// remove this:
			balance, confirm, err = bs.GetBalance(item.ReceiveAddress)
			fmt.Println("GetBalance btc response: ", balance, confirm, err)

			// if !isRateLimitErr {
			// 	balanceInfo, ok := balanceMaps[item.ReceiveAddress]
			// 	// If the key exists
			// 	if ok {
			// 		balance = big.NewInt(0).SetUint64(balanceInfo.Balance)
			// 		if len(balanceInfo.TxRefs) > 0 {
			// 			confirm = balanceInfo.TxRefs[0].Confirmations
			// 		}
			// 	}
			// }
			if err != nil {
				// get balance from quicknode:
				var balanceQuickNode *structure.BlockCypherWalletInfo
				balanceQuickNode, err = btc.GetBalanceFromQuickNode(item.ReceiveAddress, u.Config.QuicknodeAPI)
				if err == nil {
					if balanceQuickNode != nil {
						balance = big.NewInt(int64(balanceQuickNode.Balance))
						// check confirm:
						if len(balanceQuickNode.Txrefs) > 0 {
							var txInfo *btc.QuickNodeTx
							txInfo, err = btc.CheckTxfromQuickNode(balanceQuickNode.Txrefs[0].TxHash, u.Config.QuicknodeAPI)
							if err == nil {
								if txInfo != nil {
									confirm = txInfo.Result.Confirmations
								}

							} else {
								go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "CheckTxfromQuickNode from quicknode - with err", err.Error(), true)
							}
						}
					}

				} else {
					go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance from quicknode - with err", err.Error(), true)
				}
			}

		} else if item.PayType == utils.NETWORK_ETH {
			// check eth balance:

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			balance, err = ethClient.GetBalance(ctx, item.ReceiveAddress)
			fmt.Println("GetBalance eth response: ", balance, err)

			confirm = 1
		}

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error(), true)
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance", err.Error(), false)
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "GetBalance", "0", false)
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "SetString(amount) err", err.Error(), true)
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error(), true)
			continue
		}

		// set receive balance:
		item.Balance = amount.String()

		if balance.Uint64() < amount.Uint64() {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "compare balance err", err.Error(), true)

			item.Status = entity.StatusMint_NeedToRefund
			item.ReasonRefund = "Not enough balance."
			u.Repo.UpdateMintNftBtc(&item)
			continue
		}

		if confirm == 0 {
			item.Status = entity.StatusMint_WaitingForConfirms
			u.Repo.UpdateMintNftBtc(&item)
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Updated StatusMint_WaitingForConfirms", "0", true)
		}
		if confirm >= 1 {
			// received fund:
			item.Status = entity.StatusMint_ReceivedFund
			item.IsConfirm = true

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Updated StatusMint_ReceivedFund", "ok", true)
			logger.AtLog.Logger.Info(fmt.Sprintf("JobMint_CheckBalance.CheckReceiveFund.%s", item.ReceiveAddress), zap.Any("item", item))
			go u.Notify("JobMint_CheckBalance", item.ReceiveAddress, fmt.Sprintf("%s received %s %d from [UUID] %s", item.ReceiveAddress, item.PayType, item.Status, item.UUID))
		}

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			fmt.Printf("Could not UpdateMintNftBtc uuid %s - with err: %v", item.UUID, err)
			continue
		}

		// check Project and make sure index < max supply
		p, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			logger.AtLog.Logger.Error("JobMint_CheckBalance.FindProjectByTokenID", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "FindProjectByTokenID", err.Error(), true)
			continue
		}
		if p.MintingInfo.Index >= p.MaxSupply {
			// update need to return:
			u.updateProjectMintedOut(&item, "JobMint_CheckBalance")
			continue
		}

		// create batch record:
		if item.Status == entity.StatusMint_ReceivedFund && item.Quantity > 1 {
			// create
			// verify item
			listPath, _ := u.Repo.CountBatchRecordOfItems(item.UUID)

			totaltem := item.Quantity - len(listPath)

			for i := 0; i < totaltem-1; i++ { // n - item
				batchItem := entity.MintNftBtc{
					BatchParentId:     item.UUID,
					ProjectID:         item.ProjectID,
					Status:            entity.StatusMint_ReceivedFund, // wait for mint.
					PayType:           item.PayType,
					IsConfirm:         true,
					UserAddress:       item.UserAddress,
					OriginUserAddress: item.OriginUserAddress,
					RefundUserAdress:  item.RefundUserAdress,
					IsSubItem:         true,
					ProjectMintPrice:  item.ProjectMintPrice,
					ProjectNetworkFee: item.ProjectNetworkFee,

					NetworkFeeByPayType: item.NetworkFeeByPayType,
					MintPriceByPayType:  item.MintPriceByPayType,
					Amount:              item.EstFeeInfo[item.PayType].TotalAmount,
					Quantity:            1,
					UserID:              item.UserID,
					IsMerged:            item.IsMerged,
					EthRate:             item.EthRate,
					BtcRate:             item.BtcRate,

					EstFeeInfo: item.EstFeeInfo,
					IsDiscount: item.IsDiscount,

					FeeRate:  item.FeeRate,
					Platform: item.Platform,
				}
				// insert now:
				err = u.Repo.InsertMintNftBtc(&batchItem)
				if err != nil {
					go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckBalance", item.TableName(), item.Status, "Can not InsertMintNftBtc sub item", i, true)
					logger.AtLog.Logger.Error("u.CheckReceiveFund.InsertMintNftBtc", zap.Error(err))
					continue
				}

			}
		}

	}

	return nil
}

// job 2: mint nft now:
func (u Usecase) JobMint_MintNftBtc() error {

	listToMint, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_ReceivedFund)})
	if len(listToMint) == 0 {
		return nil
	}

	feeRateCurrent, _ := u.getFeeRateFromChain()

	for _, item := range listToMint {

		// check if it is a child item but its parent does not have mint yet, then continue:
		if len(item.BatchParentId) > 0 {
			parentItem, _ := u.Repo.FindMintNftBtc(item.BatchParentId)
			if !(parentItem.Status == entity.StatusMint_Minting || parentItem.IsMinted) {
				continue
			}
		}

		// get data from project
		p, err := u.Repo.FindProjectByTokenID(item.ProjectID)
		if err != nil {
			logger.AtLog.Logger.Error("JobMint_MintNftBtc.FindProjectByTokenID", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "FindProjectByTokenID", err.Error(), true)
			continue
		}

		// check Project and make sure index < max supply
		if p.MintingInfo.Index >= p.MaxSupply {
			// update need to return:
			u.updateProjectMintedOut(&item, "JobMint_MintNftBtc")
			continue
		}

		// "smart fee": check fee rate with current fee:
		if feeRateCurrent != nil {
			mintNetworkFeeRate := int64(feeRateCurrent.EconomyFee) //safe to thr
			if int64(item.FeeRate) < mintNetworkFeeRate {
				message := fmt.Sprintf("Your transaction will be processed once the network fee is reduced to %d sat/vB. The current minimum network rate is %d sat/vB. Check the current rate here <a href='https://mempool.space/'>https://mempool.space</a>", item.FeeRate, mintNetworkFeeRate)
				item.MintMessage = message
				u.Repo.UpdateMintNftBtc(&item)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "SmartFee.Wait", message, true)
				continue
			}
		}

		// check platform to mint:
		if item.Platform == utils.PLATFORM_ORDINAL {

			u.MintNftViaOrdinal(&item, p)

		} else if item.Platform == utils.PLATFORM_TC {

			u.MintNftViaTrustlessComputer(&item, p)
		}
	}

	return nil
}

// function for job mint:
func (u Usecase) MintNftViaOrdinal(item *entity.MintNftBtc, p *entity.Projects) error {
	feeRate := item.FeeRate
	if feeRate == 0 {
		feeRate = entity.DEFAULT_FEE_RATE
	}

	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err := helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("JobMint_MintNftBtc.Base64DecodeRaw", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Base64DecodeRaw", err.Error(), true)
		return nil
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	logger.AtLog.Logger.Info("animation", zap.Any("animation", animation))

	// for html type:
	if animation != "" {
		animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%s-%d.html", p.TokenID, item.UUID, now))
		if err != nil {
			logger.AtLog.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "UploadBaseToBucket", err.Error(), true)
			return nil
		}
		item.FileURI = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)

	} else {
		// for image type:
		images := p.Images
		logger.AtLog.Logger.Info("images", zap.Any("len(images)", len(images)))
		if len(images) > 0 {
			item.FileURI = images[0]
			newImages := []string{}
			processingImages := p.ProcessingImages

			//remove the project's image out of the current projects
			for i := 1; i < len(images); i++ {
				newImages = append(newImages, images[i])
			}
			processingImages = append(p.ProcessingImages, item.FileURI)
			p.Images = newImages
			p.ProcessingImages = processingImages
		}
	}
	//end Animation URL
	if len(item.FileURI) == 0 {
		err = errors.New("There is no file uri to mint")
		logger.AtLog.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "UploadBaseToBucket", err.Error(), true)
		return nil
	}

	baseUrl, err := url.Parse(item.FileURI)
	if err != nil {
		logger.AtLog.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Parse(FileURI)", err.Error(), true)
		return nil
	}

	// start call rpc mint nft now:
	mintData := ord_service.MintRequest{
		WalletName:  os.Getenv("ORD_MASTER_ADDRESS"),
		FileUrl:     baseUrl.String(),
		FeeRate:     int(feeRate),
		DryRun:      false,
		RequestId:   item.UUID,      // for tracking log
		ProjectID:   item.ProjectID, // for tracking log
		FileUrlUnit: item.FileURI,   // for tracking log

		AutoFeeRateSelect: false, // not auto

		// new key for ord v5.1, support mint + send in 1 tx:
		DestinationAddress: item.OriginUserAddress,
	}

	// logger.AtLog.Logger.Info("mintData", zap.Any("mintData", mintData))
	// execute mint:
	resp, respStr, err := u.OrdService.Mint(mintData)
	if err != nil {

		tagMention := ""
		// check fee:
		if item.ProjectNetworkFee > 500000 {
			tagMention = "<@phuong> <@yen> please check: "
		}

		logger.AtLog.Logger.Error("JobMint_MintNftBtc.OrdService", zap.Error(err))
		messageError := tagMention + respStr + "|" + err.Error() + fmt.Sprintf("| network fe: %d", item.ProjectNetworkFee)
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.Mint", item.TableName(), item.Status, mintData, messageError, true)
		return nil
	}
	// logger.AtLog.Logger.Info("mint.resp", zap.Any("resp", resp))

	go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, mintData, resp, false)

	//update
	// if not err => update status ok now:
	//TODO: handle log err: Database already open. Cannot acquire lock

	item.Status = entity.StatusMint_Minting
	item.MintMessage = ""
	item.IsMerged = true // new key for ord v5.1, support mint + send in 1 tx.

	// item.ErrCount = 0 // reset error count!

	item.OutputMintNFT = resp

	_, err = u.Repo.UpdateMintNftBtc(item)
	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateMintNftBtc", err.Error(), true)
		return nil
	}

	tmpText := resp.Stdout
	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
	jsonStr := strings.ReplaceAll(tmpText, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

	var btcMintResp ord_service.MintStdoputRespose

	err = json.Unmarshal([]byte(jsonStr), &btcMintResp)
	if err != nil {
		logger.AtLog.Logger.Error("BTCMint.helpers.JsonTransform", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.Unmarshal(btcMintResp)", err.Error(), true)
		return nil
	}

	item.TxMintNft = btcMintResp.Reveal
	item.InscriptionID = btcMintResp.Inscription
	item.MintFee = btcMintResp.Fees
	// TODO: update item
	_, err = u.Repo.UpdateMintNftBtc(item)
	if err != nil {
		fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err.Error())
	}

	// update project:
	_, err = u.Repo.UpdateProject(p.UUID, p)
	if err != nil {
		logger.AtLog.Logger.Error("JobMint_MintNftBtc.UpdateProject", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateProject", err.Error(), true)
	}
	// logger.AtLog.Logger.Info("project.Updated", zap.Any("updated", updated))

	fmt.Println("update project, token info when minting ...")

	// create entity.TokenURI
	_, err = u.CreateBTCTokenURI(item.OriginUserAddress, item.ProjectID, item.InscriptionID, item.FileURI, entity.TokenPaidType(item.PayType))
	if err != nil {
		fmt.Printf("Could CreateBTCTokenURI - with err: %v", err)
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "u.CreateBTCTokenURI()", err.Error(), true)
		return nil
	}
	_, err = u.Repo.UpdateMintNftBtcByFilter(item.UUID, bson.M{"$set": bson.M{"isUpdatedNftInfo": true}})
	if err != nil {
		fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "UpdateMintNftBtc", err.Error(), true)
	}

	go u.NotifyWithChannel(os.Getenv("SLACK_MINT_CREATED_NFT_CHANNEL_ID"), fmt.Sprintf("[MintWith][%s][uuid: %s][projectID %s]", item.PayType, item.UUID, item.ProjectID), item.ReceiveAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", item.UserAddress, resp.Stdout))
	go u.NotifyNFTMinted(item.InscriptionID)

	return nil
}

// function for job mint:
func (u Usecase) MintNftViaTrustlessComputer(item *entity.MintNftBtc, p *entity.Projects) error {

	// if not had tx, then contract:
	if !(item.IsCalledMintTc && len(item.TxMintNft) > 0) {

		u.MintNftViaTrustlessComputer_CallContract(item, p)

		if item.IsCalledMintTc && len(item.TxMintNft) > 0 {
			// call immediately
			u.MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate(item, p)
		}
	} else {
		// call inscribe via pc:
		u.MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate(item, p)

	}

	return nil
}

// function for job mint:
func (u Usecase) MintNftViaTrustlessComputer_CallContract(item *entity.MintNftBtc, p *entity.Projects) error {

	if len(p.GenNFTAddr) == 0 {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "GenNFTAddr", "empty", true)
		return nil
	}

	isMintOut, _ := u.TcClient.CheckProjectIsMintedOut(p.GenNFTAddr)
	if isMintOut {
		u.updateProjectMintedOut(item, "JobMint_MintNftBtc.MintNftViaTrustlessComputer_CallContract")
		return nil
	}

	feeRate := item.FeeRate
	if feeRate == 0 {
		feeRate = entity.DEFAULT_FEE_RATE
	}

	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err := helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		logger.AtLog.Logger.Error("JobMint_MintNftBtc.Base64DecodeRaw", zap.Error(err))
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Base64DecodeRaw", err.Error(), true)
		return nil
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	logger.AtLog.Logger.Info("animation", zap.Any("animation", animation))

	urlToMint := ""

	// for html type:
	if len(animation) == 0 {
		// for image type:
		images := p.Images
		logger.AtLog.Logger.Info("images", zap.Any("len(images)", len(images)))
		if len(images) > 0 {
			item.FileURI = images[0]
			newImages := []string{}
			processingImages := p.ProcessingImages

			//remove the project's image out of the current projects
			for i := 1; i < len(images); i++ {
				newImages = append(newImages, images[i])
			}
			processingImages = append(p.ProcessingImages, item.FileURI)
			p.Images = newImages
			p.ProcessingImages = processingImages
		}
		//end Animation URL
		if len(item.FileURI) == 0 {
			err = errors.New("There is no file uri to mint")
			logger.AtLog.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "UploadBaseToBucket", err.Error(), true)
			return nil
		}
		baseUrl, err := url.Parse(item.FileURI)
		if err != nil {
			logger.AtLog.Logger.Error("JobMint_MintNftBtc.UploadBaseToBucket", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "Parse(FileURI)", err.Error(), true)
			return nil
		}
		urlToMint = baseUrl.String()
	}

	// create byte data:
	var byteData [][]byte
	if len(urlToMint) > 0 {
		byteData, err = u.ConvertImageToByteArrayToMintTC(urlToMint)
		if err != nil {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.MintTC", item.TableName(), item.Status, "ConvertImageToByteArrayToMintTC", err.Error(), true)
			return nil
		}
	}
	// fmt.Println("byteData", byteData)

	// get free temp wallet:
	tempWallet := u.GetMintFreeTemAddress()
	if tempWallet == nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.MintTC", item.TableName(), item.Status, "GetMintFreeTemAddress", "can not get temp free wallet", true)
		return nil
	}
	fmt.Println("found temp wallet: ", tempWallet.WalletAddress)
	// encrypt:
	privateKeyDeCrypt, err := encrypt.DecryptToString(tempWallet.PrivateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		u.Logger.Error(fmt.Sprintf("JobMint_MintNftBtc.Decrypt.%s.Error", "decrypt tc privKey"), err.Error(), err)
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.DecryptToString", err.Error(), true)
		return nil
	}

	tx, err := u.TcClient.MintTC(p.GenNFTAddr, privateKeyDeCrypt, item.OriginUserAddress, byteData)
	if err != nil {

		if strings.Contains(err.Error(), "minted_out") {
			u.updateProjectMintedOut(item, "JobMint_MintNftBtc")
			return nil
		}
		fmt.Println("can not mint MintTC: ", err)
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.MintTC", item.TableName(), item.Status, tx, err.Error(), true)
		return nil
	}

	// update make busy temp wallet:
	u.Repo.UpdateTcTempWalletAddress(tempWallet.WalletAddress, entity.StatusEvmTempWallets_Busy)

	item.TxMintNft = tx
	item.MintMessage = ""
	item.IsCalledMintTc = true
	item.TcTempWallet = tempWallet.WalletAddress

	_, err = u.Repo.UpdateMintNftBtc(item)
	if err != nil {
		fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err.Error())
	}

	// update project:
	if len(urlToMint) > 0 {
		_, err = u.Repo.UpdateProject(p.UUID, p)
		if err != nil {
			logger.AtLog.Logger.Error("JobMint_MintNftBtc.UpdateProject", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "JobMint_MintNftBtc.UpdateProject", err.Error(), true)
		}
	}

	err = u.Repo.IncreaseProjectIndex(item.ProjectID)
	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc.MintTC", item.TableName(), item.Status, tx, err.Error(), true)
		return nil
	}
	return nil
}

func (u Usecase) updateProjectMintedOut(item *entity.MintNftBtc, jobName string) {
	// update need to return:
	item.ReasonRefund = "Project is minted out."
	item.Status = entity.StatusMint_NeedToRefund

	_, err := u.Repo.UpdateMintNftBtc(item)
	if err != nil {
		fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
		go u.trackMintNftBtcHistory(item.UUID, jobName, item.TableName(), item.Status, "can not update need to refund for minted out", err.Error(), true)
	}
	err = fmt.Errorf("project %s is minted out", item.ProjectID)
	logger.AtLog.Logger.Error("projectIsMintedOut", zap.Error(err))
	go u.trackMintNftBtcHistory(item.UUID, jobName, item.TableName(), item.Status, "Updated to minted out", err.Error(), true)
}

// function for job mint:
func (u Usecase) MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate(item *entity.MintNftBtc, p *entity.Projects) error {

	var resp struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	payloadStr := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"method": "eth_inscribeTxWithTargetFeeRate",
			"params": [
				"%s",%d
			],
			"id": 1
		}`, item.TxMintNft, item.FeeRate)

	payload := strings.NewReader(payloadStr)

	fmt.Println("payloadStr: ", payloadStr)

	// context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// status, err := u.TcClient.GetTransaction(context, item.TxMintNft)

	// fmt.Println("GetTransaction tc tx: ", status, err)

	// balance, err := u.TcClient.GetBalance(context, "0x232FdCd3a77A21F3C8b50F64ba56daFF80bBfA97")

	// fmt.Println("balance, err: ", balance, err)

	client := &http.Client{}
	req, err := http.NewRequest("POST", u.Config.BlockchainConfig.TCEndpoint, payload)

	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate.http.NewRequest", err.Error(), true)
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate.client.Do", err.Error(), true)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate.ioutil.ReadAll", err.Error(), true)
		return err
	}

	fmt.Println("body", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, "MintNftViaTrustlessComputer_CallRPCEthInscribeTxWithTargetFeeRate.Unmarshal", err.Error(), true)
		return err
	}
	if len(resp.Result) == 0 && resp.Error != nil {
		// error:
		go u.trackMintNftBtcHistory(item.UUID, "JobMint_MintNftBtc", item.TableName(), item.Status, payloadStr, resp.Error, true)
		return err
	}

	// inscribe ok now:
	btcTx := resp.Result
	item.TxSendNft = btcTx
	item.MintMessage = ""
	item.IsMerged = true
	item.Status = entity.StatusMint_Minting // wait for minting confirmation...
	_, err = u.Repo.UpdateMintNftBtc(item)
	if err != nil {
		fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err.Error())
	}

	return nil
}

// job check 3 tx mint:
func (u Usecase) JobMint_CheckTxMintSend() error {

	// check tx mint from ord:
	u.checkTxMintSend_ForOrdinal()

	// check tx mint from tc:
	u.checkTxMintSend_ForTc()

	return nil
}

func (u Usecase) checkTxMintSend_ForOrdinal() error {

	// get list pending tx:
	// todo: need update platform for old records.
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatusAndPlatform([]entity.StatusMint{entity.StatusMint_Minting}, utils.PLATFORM_ORDINAL)

	if len(listTxToCheck) == 0 {
		return nil
	}

	// get list btc to check a Batch
	var batchBTCTx []string
	for _, item := range listTxToCheck {
		batchBTCTx = append(batchBTCTx, item.TxMintNft)
	}

	var err error
	isRateLimitErr := false
	txInfoMaps, _, errFromCheckBatch := btc.CheckTxMultiBlockcypher(batchBTCTx, u.Config.BlockcypherToken)

	fmt.Println("isRateLimitErr, errFromCheckBatch", txInfoMaps, errFromCheckBatch)

	if errFromCheckBatch != nil {
		if strings.Contains(errFromCheckBatch.Error(), "rate_limit") {
			isRateLimitErr = true
		}
		go u.trackMintNftBtcHistory("", "JobMint_CheckTxMintSend", entity.MintNftBtc{}.TableName(), 0, "check Batch txs err", errFromCheckBatch.Error(), true)
	} else {
		go u.trackMintNftBtcHistory("", "JobMint_CheckTxMintSend", entity.MintNftBtc{}.TableName(), 0, "check Batch txs ok", len(batchBTCTx), true)
	}

	for _, item := range listTxToCheck {

		txToCheck := item.TxMintNft
		var confirm int64 = -1

		if !isRateLimitErr && txInfoMaps != nil {
			txInfo, ok := txInfoMaps[txToCheck]

			fmt.Println("txInfo, ok", txInfo, ok)

			// If the key exists
			if ok {
				confirm = int64(txInfo.Confirmations)
			} else {
				// err = errors.New("tx invalid")
				// if errFromCheckBatch != nil {
				// 	err = errors.New(errFromCheckBatch.Error())
				// }
				// go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx: "+txToCheck, err.Error(), true)
				txInfoQn, err := btc.CheckTxfromQuickNode(txToCheck, u.Config.QuicknodeAPI)
				if err == nil {
					if txInfoQn != nil {
						confirm = int64(txInfoQn.Result.Confirmations)
					}

				} else {
					go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "CheckTxfromQuickNode from quicknode - with err", err.Error(), true)
				}
			}
		} else {
			txInfoQn, err := btc.CheckTxfromQuickNode(txToCheck, u.Config.QuicknodeAPI)
			if err == nil {
				if txInfoQn != nil {
					confirm = int64(txInfoQn.Result.Confirmations)
				}

			} else {
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "CheckTxfromQuickNode from quicknode - with err", err.Error(), true)
			}
		}

		fmt.Println("confirm===>", confirm)

		// check confirm >= 1
		if confirm >= 1 {

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, true)
			// tx ok now:

			// update for ord5
			item.Status = entity.StatusMint_Minted
			item.IsMinted = true

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
				continue
			}

			err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.InscriptionID)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_CheckTxMintSend.%s.UpdateTokenOnchainStatusByTokenId.Error", item.InscriptionID), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "UpdateTokenOnchainStatusByTokenId()", err.Error(), true)
			}
			// update inscription_index for token uri
			go u.getInscribeInfoForMintSuccessToUpdate(item.InscriptionID)
			go u.CreateMintActivity(item.InscriptionID, item.Amount)
			if item.ProjectMintPrice >= 100000 {
				go func(u Usecase, item entity.MintNftBtc) {
					owner, err := u.Repo.FindUserByBtcAddressTaproot(item.OriginUserAddress)
					if err != nil || owner == nil {
						return
					}
					u.AirdropCollector(item.ProjectID, item.InscriptionID, os.Getenv("AIRDROP_WALLET"), *owner, 3)
				}(u, item)
			}

		} else {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "check tx confirm", 0, false)
			continue
		}

	}

	return nil
}
func (u Usecase) checkTxMintSend_ForTc() error {

	// get list pending tx:
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatusAndPlatform([]entity.StatusMint{entity.StatusMint_Minting}, utils.PLATFORM_TC)

	if len(listTxToCheck) == 0 {
		return nil
	}

	for _, item := range listTxToCheck {

		txToCheck := item.TxMintNft
		var confirm int64 = -1

		context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		status, err := u.TcClient.GetTransaction(context, txToCheck)

		fmt.Println("GetTransaction status, err ", txToCheck, status, err)

		if err == nil {
			if status > 0 {
				confirm = 1

			} else {
				return nil
			}
		} else {
			// if error maybe tx is pending or rejected
			// TODO check timeout to detect tx is rejected or not.
		}

		// check confirm >= 1
		if confirm >= 1 {

			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, true)
			// tx ok now:

			// update for ord5
			item.Status = entity.StatusMint_Minted
			item.IsMinted = true

			//get token ID from the tx (log event)
			nftID, err := u.TcClient.GetNftIDFromTx(item.TxMintNft, os.Getenv("TRANSFER_NFT_SIGNATURE"))
			if err != nil {
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "GetNftIDFromTx()", err.Error(), true)
			}
			// save nft_id info:
			item.InscriptionID = nftID.String()

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
				continue
			}

			// update make free temp wallet:
			u.Repo.UpdateTcTempWalletAddress(item.TcTempWallet, entity.StatusEvmTempWallets_Free)

			// update token uri auto:
			p, err := u.Repo.FindProjectByTokenID(item.ProjectID)
			if err != nil {
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "project not found", "", true)
				return err
			}

			projectIndex, err := u.TcClient.GetProjectIndex(p.GenNFTAddr)

			fmt.Println("projectIndex, p.MintingInfo.Index", projectIndex, p.MintingInfo.Index)

			if err == nil {
				if p.MintingInfo.Index != int64(projectIndex) {
					// update project:
					// u.Repo.SetProjectIndex(p.TokenId, int(projectIndex))
				}
			}

			go u.getTokenInfo(structure.GetTokenMessageReq{
				ContractAddress: p.ContractAddress,
				TokenID:         item.InscriptionID,
			})
			// create mint activity:
			go u.CreateMintActivity(item.InscriptionID, item.Amount)

		} else {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMintSend", item.TableName(), item.Status, "check tx confirm", 0, false)
			continue
		}

	}

	return nil
}

// job 4: send nft:
func (u Usecase) JobMint_SendNftToUser() error {

	// get list buy order status = StatusInscribe_Minted:
	listTosendNft, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Minted})
	if len(listTosendNft) == 0 {
		return nil
	}

	for _, item := range listTosendNft {

		// check nft in master wallet or not:
		if len(item.InscriptionID) == 0 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "checkEmpty(nftID)", "Nft id empty", true)
			continue
		}

		// update for ord v5.1: is merged tx
		if item.IsMerged {
			// don't send, update isSent = true
			item.Status = entity.StatusMint_SentNFTToUser
			item.IsSentUser = true
			u.Repo.UpdateMintNftBtc(&item)
			continue

		}

		if false {
			listNFTsRep, err := u.GetNftsOwnerOf(os.Getenv("ORD_MASTER_ADDRESS"))
			if err != nil {
				go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.Error", err.Error(), true)
				continue
			}

			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.listNFTsRep", listNFTsRep, false)

			// parse nft data:
			var resp []struct {
				Inscription string `json:"inscription"`
				Location    string `json:"location"`
				Explorer    string `json:"explorer"`
			}

			err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
			if err != nil {
				go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.Unmarshal(listNFTsRep)", err.Error(), true)
				continue
			}
			owner := false
			for _, nft := range resp {
				if strings.EqualFold(nft.Inscription, item.InscriptionID) {
					owner = true
					break
				}

			}

			if !owner {
				go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "GetNftsOwnerOf.CheckNFTOwner", owner, true)
				continue
			}
		}

		// transfer now:
		sendTokenReq := ord_service.ExecRequest{
			Args: []string{
				"--wallet",
				os.Getenv("ORD_MASTER_ADDRESS"),
				"wallet",
				"send",
				item.OriginUserAddress,
				item.InscriptionID,
				"--fee-rate",
				fmt.Sprintf("%d", entity.DEFAULT_FEE_RATE),
			}}

		logger.AtLog.Logger.Info("sendTokenReq", zap.Any("sendTokenReq", sendTokenReq))
		mintResp, err := u.OrdService.Exec(sendTokenReq)

		go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "SendTokenByWallet.ExecRequest.SendNft()", mintResp, true)

		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("JobMin_SendNftToUser.SendTokenMKP.%s.Error", item.OriginUserAddress), zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "SendTokenByWallet.err", err.Error(), true)
			continue
		}

		//TODO: handle log err: Database already open. Cannot acquire lock

		// Update status first if none err:
		item.Status = entity.StatusMint_SendingNFTToUser
		// item.ErrCount = 0 // reset error count!

		item.OutputSendNFT = mintResp

		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			err := fmt.Errorf("Could not UpdateMintNftBtc id %s - with err: %v", item.UUID, err.Error())
			logger.AtLog.Logger.Error("JobMin_SendNftToUser.UpdateMintNftBtc", zap.Error(err))
			go u.trackMintNftBtcHistory(item.UUID, "UpdateMintNftBtc", item.TableName(), item.Status, "SendTokenMKP.UpdateMintNftBtc", err.Error(), true)
			continue
		}

		txResp := mintResp.Stdout
		//txResp := `fd31946b855cbaaf91df4b2c432f9b173e053e65a9879ac909bad028e21b950e\n`
		txResp = strings.ReplaceAll(txResp, "\n", "")

		// update tx:
		item.TxSendNft = txResp
		// item.ErrCount = 0 // reset error count!
		_, err = u.Repo.UpdateMintNftBtc(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateMintNftBtc id %s - with err: %v", item.ID, err)
			logger.AtLog.Logger.Error("JobMin_SendNftToUser.UpdateMintNftBtc", zap.Error(errPack))
			go u.trackMintNftBtcHistory(item.UUID, "JobMin_SendNftToUser", item.TableName(), item.Status, "u.UpdateMintNftBtc.JobMin_SendNftToUser", err.Error(), true)
		}

		logger.AtLog.Logger.Info(fmt.Sprintf("JobMin_SendNftToUser.SendNft.execResp.%s", item.OriginUserAddress), zap.Any("mintResp", mintResp))

	}
	return nil
}

// job 6:
// refund btc to users:
func (u Usecase) JobMint_RefundBtc() error {

	listToRefund, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_NeedToRefund)})

	if len(listToRefund) == 0 {
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_RefundBtc", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	// eth:
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_RefundBtc", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	for _, item := range listToRefund {

		if len(item.RefundUserAdress) == 0 {
			continue
		}

		// if parent item:
		if item.Quantity > 1 {
			// get list of sub-items, if all have minted then refund all:
			childItems, _ := u.Repo.CountBatchRecordOfItems(item.UUID)

			fmt.Println("childItems: ", len(childItems))

			minedItems := 0
			needRefundItems := 0
			if len(childItems) > 0 {
				for _, childItem := range childItems {
					if childItem.IsMinted {
						minedItems++
					} else if childItem.Status == entity.StatusMint_NeedToRefund {
						needRefundItems++
					}
				}
				// if not enough need-to-refund item then wait or refund&fund ...
				if !(needRefundItems == item.Quantity-1) {
					if minedItems+needRefundItems == item.Quantity-1 {
						// refund + fund now:
						// this function send+refund vs 1 tx.
						err := u.SendMasterAndRefund(item.UUID, bs, ethClient)
						if err != nil {
							go u.trackMintNftBtcHistory("", "JobMint_RefundBtc", "", "", "u.SendMasterAndRefund", err.Error(), true)
						}

					}
					continue
				} // all subitem need to refund (no sub item minted) => refund max...
			}
			// no sub item => refund max...

		} else if item.IsSubItem {
			continue
		}

		if item.PayType == utils.NETWORK_BTC {

			if len(os.Getenv("SECRET_KEY")) == 0 {
				err = errors.New("please config SECRET_KEY")
				logger.AtLog.Logger.Error("u.JobMint_RefundBtc.GenerateAddress", zap.Error(err))
				continue
			}

			// the user address to refund:
			btcRefundAddress := item.RefundUserAdress

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.Decrypt.%s.Error", btcRefundAddress), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}

			// send user now:
			tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.ReceiveAddress, btcRefundAddress, -1, btc.PreferenceMedium)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", btcRefundAddress), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
				continue
			}
			// save tx:
			item.TxRefund = tx
			item.Status = entity.StatusMint_Refunding
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.UpdateMintNftBtc.%s.Error", btcRefundAddress), zap.Error(err))
				continue
			}
		} else if item.PayType == utils.NETWORK_ETH {

			ethAdressRefund := item.RefundUserAdress

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.Decrypt.%s.Error", ethAdressRefund), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}
			tx, value, err := ethClient.TransferMax(privateKeyDeCrypt, ethAdressRefund)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.ethClient.TransferMax.%s.Error", ethAdressRefund), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_RefundBtc", item.TableName(), item.Status, "JobMint_RefundBtc.ethClient.TransferMax", err.Error(), true)
				continue
			}
			// save tx:
			item.TxRefund = tx
			item.AmountRefundUser = value
			item.Status = entity.StatusMint_Refunding
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_RefundBtc.UpdateMintNftBtc.%s.Error", ethAdressRefund), zap.Error(err))
				continue
			}

		}
		time.Sleep(3 * time.Second)
	}

	return nil
}

// job 7:
// send send max fund to master address:
func (u Usecase) JobMint_SendFundToMaster() error {

	listToSentMaster, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint(entity.StatusMint_SentNFTToUser)})

	if len(listToSentMaster) == 0 {
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	// eth:
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	for _, item := range listToSentMaster {

		// if parent item:
		if item.Quantity > 1 {
			// get list of sub-items, if all have minted then send all funds to master:
			childItems, _ := u.Repo.CountBatchRecordOfItems(item.UUID)

			fmt.Println("childItems: ", len(childItems))

			minedItems := 0
			needRefundItems := 0
			if len(childItems) > 0 {
				for _, childItem := range childItems {
					if childItem.IsMinted {
						minedItems++
					} else if childItem.Status == entity.StatusMint_NeedToRefund {
						needRefundItems++
					}
				}
			}
			// if not enough mint then wait or refund&fund ...
			if !(minedItems == item.Quantity-1) {
				if minedItems+needRefundItems == item.Quantity-1 {
					// refund + fund now:
					err = u.SendMasterAndRefund(item.UUID, bs, ethClient)
					if err != nil {
						go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "u.SendMasterAndRefund", err.Error(), true)
					}
				}
				continue
			} else {
				testCronTab, _ := u.Repo.FindCronJobManagerByUUID("64071ce60ae9297684ebc528_1")
				if testCronTab == nil || testCronTab.Enabled {
					go u.trackMintNftBtcHistory("", "JobMint_SendFundToMaster", "", "", "u.SendMasterAndRefund.pause for test", item.UUID, true)
					continue
				}
			}
			// send master all.

		} else if item.IsSubItem {
			continue
		}

		if item.PayType == utils.NETWORK_BTC {

			if len(os.Getenv("SECRET_KEY")) == 0 {
				err = errors.New("please config SECRET_KEY")
				logger.AtLog.Logger.Error("u.JobMint_SendFundToMaster.GenerateAddress", zap.Error(err))
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SECRET_KEY.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), zap.Error(err))
				continue
			}

			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.Decrypt.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
				continue
			}

			// send master now:
			tx, err := bs.SendTransactionWithPreferenceFromSegwitAddress(privateKeyDeCrypt, item.ReceiveAddress, u.Config.MASTER_ADDRESS_CLAIM_BTC, -1, btc.PreferenceMedium)
			if err != nil {

				// check if not enough balance:
				if strings.Contains(err.Error(), "insufficient priority and fee for relay") {
					item.Status = entity.StatusMint_NotEnoughBalanceToSendMaster
					u.Repo.UpdateMintNftBtc(&item)

				}

				if strings.Contains(err.Error(), "already exists") {
					item.Status = entity.StatusMint_AlreadySentMaster
					item.TxSendMaster = err.Error()
					u.Repo.UpdateMintNftBtc(&item)

				}
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.SendTransactionWithPreferenceFromSegwitAddress", err.Error(), true)
				time.Sleep(1 * time.Second)
				continue
			}
			// save tx:
			item.TxSendMaster = tx
			item.AmountSentMaster = item.Amount
			item.Status = entity.StatusMint_SendingFundToMaster // TODO: need to a job to check tx.
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.UpdateBtcWalletAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_BTC), zap.Error(err))
				continue
			}
		} else if item.PayType == utils.NETWORK_ETH {
			privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.Decrypt.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.DecryptToString", err.Error(), true)
				continue
			}
			tx, amount, err := ethClient.TransferMax(privateKeyDeCrypt, u.Config.MASTER_ADDRESS_CLAIM_ETH)
			if err != nil {

				// check if not enough balance:
				if strings.Contains(err.Error(), "rlp: cannot encode negative big.Int") {
					item.Status = entity.StatusMint_NotEnoughBalanceToSendMaster
					u.Repo.UpdateMintNftBtc(&item)
				}

				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.ethClient.TransferMax.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), zap.Error(err))
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_SendFundToMaster", item.TableName(), item.Status, "JobMint_SendFundToMaster.ethClient.TransferMax", err.Error(), true)
				time.Sleep(1 * time.Second)
				continue
			}
			// save tx:
			item.TxSendMaster = tx
			item.AmountSentMaster = amount
			item.Status = entity.StatusMint_SendingFundToMaster
			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				logger.AtLog.Logger.Error(fmt.Sprintf("JobMint_SendFundToMaster.UpdateBtcWalletAddress.%s.Error", u.Config.MASTER_ADDRESS_CLAIM_ETH), zap.Error(err))
				continue
			}

		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// job 8: check tx send master, refund user:
func (u Usecase) JobMint_CheckTxMasterAndRefund() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		go u.trackMintNftBtcHistory("", "JobMint_CheckTxMasterAndRefund", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error(), true)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMint_CheckBalance", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list pending tx:
	listTxToCheck, _ := u.Repo.ListMintNftBtcByStatus([]entity.StatusMint{entity.StatusMint_Refunding, entity.StatusMint_SendingFundToMaster})
	if len(listTxToCheck) == 0 {
		return nil
	}

	for _, item := range listTxToCheck {

		if item.IsSubItem {
			continue
		}

		var txToCheck string
		var confirm int64 = -1

		if item.Status == entity.StatusMint_Refunding {
			txToCheck = item.TxRefund
		} else if item.Status == entity.StatusMint_SendingFundToMaster {
			txToCheck = item.TxSendMaster
		}

		// amountSent := ""

		if item.PayType == utils.NETWORK_ETH {
			context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			status, err := ethClient.GetTransaction(context, txToCheck)
			if err == nil {
				if status > 0 {
					confirm = 1

				} else {
					continue
				}
			} else {
				// if error maybe tx is pending or rejected
				// TODO check timeout to detect tx is rejected or not.
			}
		} else {
			// check with api btc:
			txInfo, err := bs.CheckTx(txToCheck)
			if err != nil {
				fmt.Printf("Could not bs - with err: %v", err)
				go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx: "+txToCheck, err.Error(), true)
				continue
			}

			confirm = int64(txInfo.Confirmations)
			// amountSent = txInfo.Total.String()

		}

		go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, false)

		// just check 1 confirm:
		if confirm >= 1 {
			go u.trackMintNftBtcHistory(item.UUID, "JobMint_CheckTxMasterAndRefund", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txToCheck, confirm, true)
			// tx ok now:

			if item.Status == entity.StatusMint_Refunding {
				item.Status = entity.StatusMint_Refunded
				item.IsRefund = true
				if item.PayType == utils.NETWORK_BTC {
					// item.AmountRefundUser = amountSent
				}
			} else if item.Status == entity.StatusMint_SendingFundToMaster {
				item.Status = entity.StatusMint_SentFundToMaster
				item.IsSentMaster = true
			}
			if item.Quantity <= 1 {
				if item.PayType == utils.NETWORK_BTC {
					// item.AmountSentMaster = amountSent
				}
			} else {
				// loop on sub items:
				listSubItem, _ := u.Repo.CountBatchRecordOfItems(item.UUID)
				for _, sub := range listSubItem {
					if sub.Status == entity.StatusMint_Refunding || sub.Status == entity.StatusMint_NeedToRefund {
						sub.Status = entity.StatusMint_Refunded
						if item.Status == entity.StatusMint_Refunded {
							sub.TxRefund = item.TxRefund
						}
						if item.Status == entity.StatusMint_SentFundToMaster {
							sub.TxRefund = item.TxSendMaster
						}
						sub.IsRefund = true
						u.Repo.UpdateMintNftBtc(&sub)
					}
					if sub.Status == entity.StatusMint_SendingFundToMaster || sub.Status == entity.StatusMint_SentNFTToUser {
						sub.Status = entity.StatusMint_SentFundToMaster
						sub.IsSentMaster = true
						if item.Status == entity.StatusMint_Refunded {
							sub.TxSendMaster = item.TxRefund
						}
						if item.Status == entity.StatusMint_SentFundToMaster {
							sub.TxSendMaster = item.TxSendMaster
						}
						u.Repo.UpdateMintNftBtc(&sub)
					}
				}
			}

			_, err = u.Repo.UpdateMintNftBtc(&item)
			if err != nil {
				fmt.Printf("Could not JobMint_CheckTxMasterAndRefund id %s - with err: %v", item.ID, err)
				continue
			}
		}

	}

	return nil
}

func (u Usecase) SendMasterAndRefund(uuid string, bs *btc.BlockcypherService, ethClient *eth.Client) error {
	mintItem, _ := u.Repo.FindMintNftBtcByNftID(uuid)
	if mintItem == nil {
		return errors.New("item not found")
	}
	// check valid (just for parent items):
	if mintItem.Quantity <= 1 || len(mintItem.BatchParentId) > 0 || mintItem.IsSubItem {
		return errors.New("item invalid")
	}

	if !(mintItem.Status == entity.StatusMint_SentNFTToUser || mintItem.Status == entity.StatusMint_NeedToRefund) {
		return errors.New("item invalid")
	}

	// get list child item:
	childItems, _ := u.Repo.CountBatchRecordOfItems(uuid)

	if len(childItems) == 0 {
		return errors.New("list sub item not found")
	}
	minedItems := 0
	needRefundItems := 0

	totalMintedAmount := big.NewInt(0)
	totalRefundAmount := big.NewInt(0)

	var listSubIdMinted []string
	var listSubIdRefund []string

	amountPerItem := big.NewInt(0)

	if len(childItems) > 0 {
		for _, childItem := range childItems {
			childAmount, _ := big.NewInt(0).SetString(childItem.Amount, 10)
			amountPerItem = childAmount
			if childItem.IsMinted {
				minedItems++
				totalMintedAmount = totalMintedAmount.Add(totalMintedAmount, childAmount)
				listSubIdMinted = append(listSubIdMinted, childItem.UUID)
			} else if childItem.Status == entity.StatusMint_NeedToRefund {
				needRefundItems++
				totalRefundAmount = totalRefundAmount.Add(totalRefundAmount, childAmount)
				listSubIdRefund = append(listSubIdRefund, childItem.UUID)
			}
		}
	}
	// add amount of parent item:
	if mintItem.Status == entity.StatusMint_NeedToRefund {
		totalRefundAmount = totalRefundAmount.Add(totalRefundAmount, amountPerItem)
		needRefundItems++
	} else if mintItem.Status == entity.StatusMint_SentNFTToUser {
		totalMintedAmount = totalMintedAmount.Add(totalMintedAmount, amountPerItem)
		minedItems++
	}

	// if !(needRefundItems == mintItem.Quantity-1)
	{
		if minedItems+needRefundItems == mintItem.Quantity {
			// refund + fund now:
			// code this function send+refund vs 1 tx.
			if mintItem.PayType == utils.NETWORK_BTC {

				destinations := make(map[string]int)
				// refund info:
				destinations[mintItem.RefundUserAdress] = int(totalRefundAmount.Int64())

				// send fund to master:
				destinations[u.Config.MASTER_ADDRESS_CLAIM_BTC] = int(totalMintedAmount.Int64())

				// log destinations:
				go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "destinations est to send", destinations, true)

				privateKeyDeCrypt, err := encrypt.DecryptToString(mintItem.PrivateKey, os.Getenv("SECRET_KEY"))
				if err != nil {
					logger.AtLog.Logger.Error(fmt.Sprintf("SendMasterAndRefund.Decrypt.%s.Error", mintItem.ReceiveAddress), zap.Error(err))
					go u.trackMintNftBtcHistory(mintItem.UUID, "JobMint_RefundBtc", mintItem.TableName(), mintItem.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
					return err
				}

				txFee, err := bs.EstimateFeeTransactionWithPreferenceFromSegwitAddressMultiAddress(privateKeyDeCrypt, mintItem.ReceiveAddress, destinations, btc.PreferenceMedium)
				if err != nil {
					// check if not enough balance:
					if strings.Contains(err.Error(), "insufficient priority and fee for relay") {
						mintItem.Status = entity.StatusMint_NotEnoughBalanceToSendMaster
						u.Repo.UpdateMintNftBtc(mintItem)
					}

					if strings.Contains(err.Error(), "already exists") {
						mintItem.Status = entity.StatusMint_AlreadySentMaster
						mintItem.IsSentMaster = true
						mintItem.TxSendMaster = err.Error()
						u.Repo.UpdateMintNftBtc(mintItem)

					}
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "EstimateFeeTransactionWithPreferenceFromSegwitAddressMultiAddress err", err.Error(), true)
					return err
				}

				// update refund info with -fee:
				amountToRefundWithFee := big.NewInt(0).Sub(totalRefundAmount, txFee)
				destinations[mintItem.RefundUserAdress] = int(amountToRefundWithFee.Int64())

				// log destinations:
				go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "destinations final to send", destinations, true)

				// pause to test first
				testCronTab, _ := u.Repo.FindCronJobManagerByUUID("64071ce60ae9297684ebc528_1")
				if testCronTab == nil || testCronTab.Enabled {
					return errors.New("pause for test -> SendMasterAndRefund" + mintItem.UUID)
				}

				txID, err := bs.SendTransactionWithPreferenceFromSegwitAddressMultiAddress(
					privateKeyDeCrypt,
					mintItem.ReceiveAddress,
					destinations,
					btc.PreferenceMedium,
				)
				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "SendTransactionWithPreferenceFromReceiveAddressMultiAddress err", err.Error(), true)
					return err
				}
				// update now:
				// update parent item:
				mintItem.Status = entity.StatusMint_SendingFundToMaster
				mintItem.AmountRefundUser = totalRefundAmount.String()
				mintItem.AmountSentMaster = totalMintedAmount.String()
				mintItem.IsRefund = true
				mintItem.TxSendMaster = txID
				mintItem.TxRefund = txID
				mintItem.FeeSendMaster = txFee.String()

				_, err = u.Repo.UpdateMintNftBtc(mintItem)
				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtc.Done err", err.Error(), true)
					return err
				}

				// update sub item:
				if len(listSubIdMinted) > 0 {
					_, err = u.Repo.UpdateMintNftBtcSubItemRefundOrDone(listSubIdMinted, entity.StatusMint_SendingFundToMaster, txID, false)
					if err != nil {
						go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtcSubItemRefundOrDone.UpdateMaster err", err.Error(), true)
					}
				}
				if len(listSubIdRefund) > 0 {
					_, err = u.Repo.UpdateMintNftBtcSubItemRefundOrDone(listSubIdRefund, entity.StatusMint_Refunding, txID, true)
					if err != nil {
						go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtcSubItemRefundOrDone.UpdateRefund err", err.Error(), true)
					}
				}
				return nil // all is ok now!

			} else if mintItem.PayType == utils.NETWORK_ETH {

				destinations := make(map[string]*big.Int)
				// refund info:
				destinations[mintItem.RefundUserAdress] = totalRefundAmount
				// send fund to master:
				destinations[u.Config.MASTER_ADDRESS_CLAIM_ETH] = totalMintedAmount

				fmt.Println("destinations 1: ", destinations)

				// log destinations:
				go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "destinations eth est to send", destinations, true)

				// cal tx fee:
				gasPrice, err := ethClient.GetClient().SuggestGasPrice(context.Background())
				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "SuggestGasPrice err", err.Error(), true)
					return err
				}

				fmt.Println("SendMasterAndRefund gasPrice: ", gasPrice, len(destinations))

				gasLimit := 25000 * (len(destinations))

				txFee := new(big.Int).Mul(new(big.Int).SetUint64(gasPrice.Uint64()), new(big.Int).SetInt64(int64(gasLimit)))

				fmt.Println("txFee: ", txFee)

				balance, err := ethClient.GetClient().BalanceAt(context.Background(), common.HexToAddress(mintItem.ReceiveAddress), nil)

				fmt.Println("SendMasterAndRefund eth balance: ", balance)

				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "BalanceAt err", err.Error(), true)
					return err
				}

				if txFee.Uint64() > balance.Uint64() {

					mintItem.Status = entity.StatusMint_NotEnoughBalanceToSendMaster
					u.Repo.UpdateMintNftBtc(mintItem)

					fmt.Println("not enough balance: ", txFee.Uint64(), balance.Uint64())
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "check fee and balance", "txFee > balance", true)
					return errors.New("not enough balance < tx fee")
				}

				// update refund info with -fee:
				amountToRefundWithFee := big.NewInt(0).Sub(totalRefundAmount, txFee)
				destinations[mintItem.RefundUserAdress] = amountToRefundWithFee

				if amountToRefundWithFee.Uint64() < txFee.Uint64() {

					mintItem.Status = entity.StatusMint_NotEnoughBalanceToSendMaster
					u.Repo.UpdateMintNftBtc(mintItem)

					fmt.Println("not enough amountToRefundWithFee: ", txFee.Uint64(), balance.Uint64())
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "check fee and balance", "amountToRefundWithFee < txFee", true)
					return err
				}

				// log destinations:
				go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "destinations eth final to send", destinations, true)

				// check test first:
				testCronTab, _ := u.Repo.FindCronJobManagerByUUID("64071ce60ae9297684ebc528_1")
				if testCronTab == nil || testCronTab.Enabled {
					return errors.New("pause for test -> SendMasterAndRefund: " + mintItem.UUID)
				}

				privateKeyDeCrypt, err := encrypt.DecryptToString(mintItem.PrivateKey, os.Getenv("SECRET_KEY"))
				if err != nil {
					logger.AtLog.Logger.Error(fmt.Sprintf("SendMasterAndRefund.Decrypt.%s.Error", mintItem.ReceiveAddress), zap.Error(err))
					go u.trackMintNftBtcHistory(mintItem.UUID, "JobMint_RefundBtc", mintItem.TableName(), mintItem.Status, "JobMint_RefundBtc.DecryptToString", err.Error(), true)
					return err
				}
				txID, err := ethClient.SendMulti(
					"0xcd5485b34c9902527bbee21f69312fe2a73bc802",
					privateKeyDeCrypt,
					destinations,
					gasPrice,
					uint64(gasLimit),
				)

				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "SendMulti err", err.Error(), true)
					return err
				}
				// update now:
				// update parent item:
				mintItem.Status = entity.StatusMint_SendingFundToMaster
				mintItem.AmountRefundUser = totalRefundAmount.String()
				mintItem.AmountSentMaster = totalMintedAmount.String()
				mintItem.IsRefund = true
				mintItem.TxSendMaster = txID
				mintItem.TxRefund = txID
				mintItem.FeeSendMaster = txFee.String()
				_, err = u.Repo.UpdateMintNftBtc(mintItem)
				if err != nil {
					go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtc.Done err", err.Error(), true)
					return err
				}

				// update sub item:
				if len(listSubIdMinted) > 0 {
					_, err = u.Repo.UpdateMintNftBtcSubItemRefundOrDone(listSubIdMinted, entity.StatusMint_SendingFundToMaster, txID, false)
					if err != nil {
						go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtcSubItemRefundOrDone.UpdateMaster err", err.Error(), true)
					}
				}
				if len(listSubIdRefund) > 0 {
					_, err = u.Repo.UpdateMintNftBtcSubItemRefundOrDone(listSubIdRefund, entity.StatusMint_Refunding, txID, true)
					if err != nil {
						go u.trackMintNftBtcHistory(mintItem.UUID, "SendMasterAndRefund", mintItem.TableName(), mintItem.Status, "UpdateMintNftBtcSubItemRefundOrDone.UpdateRefund err", err.Error(), true)
					}
				}
				return nil // all is ok now!
			}

		}
	}
	return errors.New("don't need refund for this: " + uuid)

}

func (u *Usecase) trackMintNftBtcHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}, notify bool) {

	trackData := &entity.MintNftBtcLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.CreateMintNftBtcLog(trackData)
	if err != nil {
		fmt.Printf("trackMintNftBtcHistory.%s.Error:%s", name, err.Error())
	}

	if notify && requestMsg != nil && responseMsg != nil {
		requestMsgStr := fmt.Sprintf("%v", requestMsg)
		responseMsgStr := fmt.Sprintf("%v", responseMsg)

		preText := fmt.Sprintf("[App: %s][recordID %s] - %s", os.Getenv("JAEGER_SERVICE_NAME"), id, requestMsgStr)

		if _, _, err := u.Slack.SendMessageToSlackWithChannel(os.Getenv("SLACK_MINT_NFT_CHANNEL_ID"), preText, name, responseMsgStr); err != nil {
			fmt.Println("s.Slack.SendMessageToSlack err", err)
		}
	}

}
func (u Usecase) getInscribeInfoForMintSuccessToUpdate(inscriptionID string) error {
	inscribeInfo, err := u.GetInscribeInfo(inscriptionID)
	if err != nil {
		return err
	}
	u.Repo.UpdateTokenInscriptionIndexForMint(inscriptionID, inscribeInfo.Index)

	return nil
}

// Mint flow
func (u Usecase) convertBTCToETH(amount string) (string, float64, float64, error) {

	//amount = "0.1"
	powIntput := math.Pow10(8)
	powIntputBig := new(big.Float)
	powIntputBig.SetFloat64(powIntput)
	amountMintBTC, _ := big.NewFloat(0).SetString(amount)
	amountMintBTC.Mul(amountMintBTC, powIntputBig)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("strconv.ParseFloat", zap.Error(err))
	// 	return "", err
	// }

	_ = amountMintBTC
	btcPrice, err := helpers.GetExternalPrice("BTC")
	if err != nil {
		logger.AtLog.Logger.Error("convertBTCToETH", zap.Error(err))
		return "", 0, 0, err
	}

	logger.AtLog.Logger.Info("btcPrice", zap.Any("btcPrice", btcPrice))
	ethPrice, err := helpers.GetExternalPrice("ETH")
	if err != nil {
		logger.AtLog.Logger.Error("convertBTCToETH", zap.Error(err))
		return "", 0, 0, err
	}

	btcToETH := btcPrice / ethPrice
	// btcToETH := 14.27 // remove hardcode, why tri hardcode this??

	rate := new(big.Float)
	rate.SetFloat64(btcToETH)
	amountMintBTC.Mul(amountMintBTC, rate)

	pow := math.Pow10(10)
	powBig := new(big.Float)
	powBig.SetFloat64(pow)

	amountMintBTC.Mul(amountMintBTC, powBig)
	result := new(big.Int)
	amountMintBTC.Int(result)

	logger.AtLog.Logger.Info("convertBTCToETH", zap.String("amount", amount), zap.Float64("btcPrice", btcPrice), zap.Float64("ethPrice", ethPrice))
	return result.String(), btcPrice, ethPrice, nil
}

func (u Usecase) ConvertBTCToETHWithPriceEthBtc(amount string, btcPrice, ethPrice float64) (string, float64, float64, error) {
	return u.convertBTCToETHWithPriceEthBtc(amount, btcPrice, ethPrice)
}

func (u Usecase) convertBTCToETHWithPriceEthBtc(amount string, btcPrice, ethPrice float64) (string, float64, float64, error) {

	//amount = "0.1"
	powIntput := math.Pow10(8)
	powIntputBig := new(big.Float)
	powIntputBig.SetFloat64(powIntput)
	amountMintBTC, _ := big.NewFloat(0).SetString(amount)
	amountMintBTC.Mul(amountMintBTC, powIntputBig)
	// if err != nil {
	// 	logger.AtLog.Logger.Error("strconv.ParseFloat", zap.Error(err))
	// 	return "", err
	// }

	_ = amountMintBTC
	btcToETH := btcPrice / ethPrice
	// btcToETH := 14.27 // remove hardcode, why tri hardcode this??

	rate := new(big.Float)
	rate.SetFloat64(btcToETH)
	amountMintBTC.Mul(amountMintBTC, rate)

	pow := math.Pow10(10)
	powBig := new(big.Float)
	powBig.SetFloat64(pow)

	amountMintBTC.Mul(amountMintBTC, powBig)
	result := new(big.Int)
	amountMintBTC.Int(result)

	logger.AtLog.Logger.Info("convertBTCToETH", zap.String("amount", amount), zap.Float64("btcPrice", btcPrice), zap.Float64("ethPrice", ethPrice))
	return result.String(), btcPrice, ethPrice, nil
}

// please donate P some money:
func (u Usecase) calMintFeeInfo(mintBtcPrice, fileSize, feeRate int64, btcRate, ethRate float64) (map[string]entity.MintFeeInfo, error) {

	fmt.Println("fileSize, feeRate: ", fileSize, feeRate)

	listMintFeeInfo := make(map[string]entity.MintFeeInfo)

	mintPrice := big.NewInt(0)
	feeSendFund := big.NewInt(utils.FEE_BTC_SEND_AGV)
	feeSendNft := big.NewInt(utils.FEE_BTC_SEND_NFT)
	feeMintNft := big.NewInt(0)

	totalAmountToMint := big.NewInt(0)
	netWorkFee := big.NewInt(0)

	var err error

	// cal min price:
	mintPrice = mintPrice.SetUint64(uint64(mintBtcPrice))

	if fileSize > 0 {

		fileSize += utils.MIN_FILE_SIZE // auto add 4kb

		// auto fee if feeRate <= 0:
		if feeRate <= 0 {
			calNetworkFee := u.networkFeeBySize(int64(fileSize / 4))
			if calNetworkFee == -1 {
				err = errors.New("can not cal networkFeeBySize")
				logger.AtLog.Logger.Error("u.calMintFeeInfo.networkFeeBySize", zap.Error(err))
				return nil, err
			}
			feeMintNft = big.NewInt(calNetworkFee)
		} else {
			calNetworkFee := int64(fileSize/4) * feeRate
			// fee mint:
			feeMintNft = big.NewInt(calNetworkFee)
		}

	}

	// default feeMintNft if 0:
	if feeMintNft.Uint64() == 0 {
		feeMintNft = big.NewInt(0).SetUint64(feeSendNft.Uint64())
	}

	fmt.Println("feeMintNft: ", feeMintNft)

	if btcRate <= 0 {
		btcRate, err = helpers.GetExternalPrice("BTC")
		if err != nil {
			logger.AtLog.Logger.Error("getExternalPrice", zap.Error(err))
			return nil, err
		}

		ethRate, err = helpers.GetExternalPrice("ETH")
		if err != nil {
			logger.AtLog.Logger.Error("helpers.GetExternalPrice", zap.Error(err))
			return nil, err
		}
	}

	fmt.Println("btcRate, ethRate", btcRate, ethRate)

	// total amount by BTC:
	netWorkFee = netWorkFee.Add(feeMintNft, feeSendNft)  // + feeMintNft	+ feeSendNft
	netWorkFee = netWorkFee.Add(netWorkFee, feeSendFund) // + feeSendFund

	totalAmountToMint = totalAmountToMint.Add(mintPrice, netWorkFee) // mintPrice, netWorkFee

	listMintFeeInfo["btc"] = entity.MintFeeInfo{

		MintPrice:   mintPrice.String(),
		MintFee:     feeMintNft.String(),
		NetworkFee:  netWorkFee.String(),
		TotalAmount: totalAmountToMint.String(),
		SendNftFee:  feeSendNft.String(),
		SendFundFee: feeSendFund.String(),

		MintPriceBigInt:   mintPrice,
		MintFeeBigInt:     feeMintNft,
		SendNftFeeBigInt:  feeSendNft,
		SendFundFeeBigInt: feeSendFund,
		NetworkFeeBigInt:  netWorkFee,
		TotalAmountBigInt: totalAmountToMint,

		EthPrice: ethRate,
		BtcPrice: btcRate,

		Decimal: 8,
	}

	fmt.Println("feeInfos[btc].MintPriceBigIn1", listMintFeeInfo["btc"].MintPriceBigInt)

	// 1. convert mint price btc to eth  ==========
	mintPriceByEth, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(mintPrice.Uint64())/1e8), btcRate, ethRate)
	if err != nil {
		logger.AtLog.Logger.Error("calMintFeeInfo.convertBTCToETHWithPriceEthBtc", zap.Error(err))
		return nil, err
	}
	// 1. set mint price by eth
	mintPriceEth, ok := big.NewInt(0).SetString(mintPriceByEth, 10)
	if !ok {
		err = errors.New("can not set mintPriceByEth")
		logger.AtLog.Logger.Error("u.calMintFeeInfo.Set(mintPriceByEth)", zap.Error(err))
		return nil, err
	}

	// 2. convert mint fee btc to eth  ==========
	feeMintNftByEth, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(feeMintNft.Uint64())/1e8), btcRate, ethRate)
	if err != nil {
		logger.AtLog.Logger.Error("calMintFeeInfo.convertBTCToETHWithPriceEthBtc", zap.Error(err))
		return nil, err
	}
	// 2. set mint fee by eth
	feeMintNftEth, ok := big.NewInt(0).SetString(feeMintNftByEth, 10)
	if !ok {
		err = errors.New("can not set feeMintNftByEth")
		logger.AtLog.Logger.Error("u.calMintFeeInfo.Set(feeMintNftByEth)", zap.Error(err))
		return nil, err
	}

	// 3. convert mint fee btc to eth ==========
	feeSendNftByEth, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(feeSendNft.Uint64())/1e8), btcRate, ethRate)
	if err != nil {
		logger.AtLog.Logger.Error("calMintFeeInfo.convertBTCToETHWithPriceEthBtc", zap.Error(err))
		return nil, err
	}
	// 3. set mint fee by eth
	feeSendNftEth, ok := big.NewInt(0).SetString(feeSendNftByEth, 10)
	if !ok {
		err = errors.New("can not set feeMintNftByEth")
		logger.AtLog.Logger.Error("u.calMintFeeInfo.Set(feeMintNftByEth)", zap.Error(err))
		return nil, err
	}

	// 4. fee send master by eth:
	feeSendFundEth := big.NewInt(utils.FEE_ETH_SEND_MASTER * 1e18)

	// total amount by ETH:
	netWorkFeeEth := big.NewInt(0).Add(feeMintNftEth, feeSendNftEth) // + feeMintNft	+ feeSendNft
	netWorkFeeEth = big.NewInt(0).Add(netWorkFeeEth, feeSendFundEth) // + feeSendFund

	totalAmountToMintEth := big.NewInt(0).Add(mintPriceEth, netWorkFeeEth) // mintPrice, netWorkFee

	listMintFeeInfo["eth"] = entity.MintFeeInfo{
		MintPrice:   mintPriceEth.String(),
		MintFee:     feeMintNftEth.String(),
		NetworkFee:  netWorkFeeEth.String(),
		TotalAmount: totalAmountToMintEth.String(),
		SendNftFee:  feeSendNftEth.String(),
		SendFundFee: feeSendFundEth.String(),

		MintPriceBigInt:   mintPriceEth,
		MintFeeBigInt:     feeMintNftEth,
		SendNftFeeBigInt:  feeSendNftEth,
		SendFundFeeBigInt: feeSendFundEth,

		NetworkFeeBigInt:  netWorkFeeEth,
		TotalAmountBigInt: totalAmountToMintEth,

		EthPrice: ethRate,
		BtcPrice: btcRate,

		Decimal: 18,
	}

	return listMintFeeInfo, err
}

// Mint flow
func (u Usecase) GetBTCToETHRate() (float64, float64, error) {
	key := "btc-eth-rate"
	exist, err := u.Cache.Exists(key)
	if err == nil {
		if *exist {
			value, err := u.Cache.GetData(key)
			if err == nil && value != nil {
				values := strings.Split(*value, "|")
				btcPrice, _ := strconv.ParseFloat(values[0], 10)
				ethPrice, _ := strconv.ParseFloat(values[1], 10)
				return btcPrice, ethPrice, nil
			}
		}
	}

	btcPrice, err := helpers.GetExternalPrice("BTC")
	if err != nil {
		logger.AtLog.Error("convertBTCToETH", zap.Error(err))
		return 0, 0, err
	}

	ethPrice, err := helpers.GetExternalPrice("ETH")
	if err != nil {
		logger.AtLog.Error("convertBTCToETH", zap.Error(err))
		return 0, 0, err
	}

	value := fmt.Sprintf("%f|%f", btcPrice, ethPrice)
	u.Cache.SetStringDataWithExpTime(key, value, 60)
	return btcPrice, ethPrice, nil
}

func (u Usecase) ConvertImageToByteArrayToMintTC(imageURL string) ([][]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("err1:", err)
		return nil, err
	}
	defer resp.Body.Close()

	var imgData [][]byte

	// split item:
	buffer := make([]byte, 350000)
	for {
		n, err := io.ReadFull(resp.Body, buffer)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			//read all image:
			imgData = append(imgData, buffer[:n])
			break
		} else if err != nil {
			fmt.Println("err2:", err)
			return nil, err
		}
		imgData = append(imgData, buffer[:n])
	}

	fmt.Printf("len: %d\n", len(imgData))

	return imgData, nil

}

func (u Usecase) GetMintFreeTemAddress() *entity.EvmTempWallets {
	// mutex := u.RedisV9.GetRedSyncClient().NewMutex("GetMintFreeTemAddress")

	var freeWallet *entity.EvmTempWallets

	// if err := mutex.Lock(); err != nil {
	// 	fmt.Println("can not lock")
	// 	return nil
	// }

	freeWallet, _ = u.Repo.GetMintFreeTempAddress()

	// if ok, err := mutex.Unlock(); !ok || err != nil {
	// 	fmt.Println("can not unlock")
	// }
	return freeWallet
}

func (u Usecase) GenMintFreeTemAddress() (string, error) {

	fmt.Println("start------")

	// mutex := u.RedisV9.GetRedSyncClient().NewMutex("GetMintFreeTemAddress")

	// if err := mutex.Lock(); err != nil {
	// 	fmt.Println("can not lock", err)
	// 	return
	// }
	// time.Sleep(1 * time.Second)

	maxItem := 200

	if len(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET")) == 0 {
		return "", errors.New("PRIVATE_KEY_FEE_TC_WALLET empty")
	}
	if len(os.Getenv("SECRET_KEY")) == 0 {
		return "", errors.New("SECRET_KEY empty")
	}

	if len(os.Getenv("TC_MULTI_CONTRACT")) == 0 {
		return "", errors.New("TC_MULTI_CONTRACT empty")
	}

	// get data:
	list, _ := u.Repo.ListEvmTempWallets()

	if len(list) == 0 {
		for i := 0; i < maxItem; i++ {
			fmt.Println("xxxxxx i =>", i)

			privateKey, _, receiveAddress, err := u.TcClient.GenerateAddress()
			if err != nil {
				logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
				return "", err
			}

			privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
			if err != nil {
				logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
				return "", err
			}

			err = u.Repo.InsertEvmTempWallets(&entity.EvmTempWallets{
				WalletAddress: receiveAddress,
				PrivateKey:    privateKeyEnCrypt,
				Status:        0,
			})
			logger.AtLog.Logger.Error("u.CreateMintReceiveAddress.Encrypt", zap.Error(err))
		}
	}

	// send JUICE:
	destinations := make(map[string]*big.Int)

	// get list again:
	list, _ = u.Repo.ListEvmTempWallets()
	for _, item := range list {
		destinations[item.WalletAddress] = big.NewInt(0.5 * 1e18)
	}

	fmt.Println("destinations: ", destinations)

	privateKeyDeCrypt, err := encrypt.DecryptToString(os.Getenv("PRIVATE_KEY_FEE_TC_WALLET"), os.Getenv("SECRET_KEY"))
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("GenMintFreeTemAddress.Decrypt.%s.Error", "can decrypt"), zap.Error(err))
		return "", err
	}

	txID, err := u.TcClient.SendMulti(
		os.Getenv("TC_MULTI_CONTRACT"),
		privateKeyDeCrypt,
		destinations,
		nil,
		0,
	)
	fmt.Println("txID, err ", txID, err)

	if err != nil {
		return "", err
	}

	// if ok, err := mutex.Unlock(); !ok || err != nil {
	// 	fmt.Println("can not unlock", err)
	// }

	fmt.Println("done------")

	return txID, err
}

func (u Usecase) SubmitTCToBtcChain(tx string, feeRate int) (string, error) {

	var resp struct {
		Result string `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	payloadStr := fmt.Sprintf(`{
			"jsonrpc": "2.0",
			"method": "eth_inscribeTxWithTargetFeeRate",
			"params": [
				"%s",%d
			],
			"id": 1
		}`, tx, feeRate)

	payload := strings.NewReader(payloadStr)

	fmt.Println("payloadStr: ", payloadStr)

	// context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// status, err := u.TcClient.GetTransaction(context, item.TxMintNft)

	// fmt.Println("GetTransaction tc tx: ", status, err)

	// balance, err := u.TcClient.GetBalance(context, "0x232FdCd3a77A21F3C8b50F64ba56daFF80bBfA97")

	// fmt.Println("balance, err: ", balance, err)

	client := &http.Client{}
	req, err := http.NewRequest("POST", u.Config.BlockchainConfig.TCEndpoint, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("body", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}
	if len(resp.Result) == 0 && resp.Error != nil {
		// error:
		return "", errors.New(resp.Error.Message)
	}

	// inscribe ok now:

	return resp.Result, nil
}
