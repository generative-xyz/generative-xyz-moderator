package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
)

// api listing....
func (u Usecase) BTCMarketplaceListingNFT(listingInfo structure.MarketplaceBTC_ListingInfo) (*entity.MarketplaceBTCListing, error) {

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
	}

	listing := entity.MarketplaceBTCListing{

		HoldOrdAddress: "",
		SellOrdAddress: listingInfo.SellOrdAddress,
		Price:          listingInfo.Price,
		PayType:        listingInfo.PayType,

		ServiceFee:    listingInfo.ServiceFee,
		IsConfirm:     false,
		IsSold:        false,
		IsCancel:      false,
		ExpiredAt:     time.Now().Add(time.Hour * time.Duration(expiredTime)),
		Name:          listingInfo.Name,
		Description:   listingInfo.Description,
		InscriptionID: listingInfo.InscriptionID,
	}
	holdOrdAddress := ""
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_marketplace_master",
			"wallet",
			"receive",
		},
	})
	if err != nil {
		u.Logger.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return &listing, err
	}

	// parse json to get address:
	// ex: {"mnemonic": "chaos dawn between remember raw credit pluck acquire satoshi rain one valley","passphrase": ""}

	jsonStr := strings.ReplaceAll(resp.Stdout, "\n", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

	var receiveResp ord_service.ReceiveCmdStdoputRespose

	err = json.Unmarshal([]byte(jsonStr), &receiveResp)
	if err != nil {
		u.Logger.Error("BTCMarketplaceListingNFT.Unmarshal", err.Error(), err)
		return nil, err
	}

	holdOrdAddress = receiveResp.Address
	listing.HoldOrdAddress = holdOrdAddress

	// check if listing is created or not
	err = u.Repo.CreateMarketplaceListingBTC(&listing)
	if err != nil {
		u.Logger.Error("BTCMarketplaceListingNFT.Repo.CreateMarketplaceListingBTC", "", err)
		return &listing, err
	}
	return &listing, nil
}

// API list listing, support filter ...
func (u Usecase) BTCMarketplaceListNFT(filter *entity.FilterString, buyableOnly bool, limit, offset int64) ([]structure.MarketplaceNFTDetail, error) {

	result := []structure.MarketplaceNFTDetail{}
	var nftList []entity.MarketplaceBTCListingFilterPipeline
	var err error

	// if buyableOnly {
	nftList, err = u.Repo.RetrieveBTCNFTListingsUnsoldForSearch(filter, limit, offset)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	nftList, err = u.Repo.RetrieveBTCNFTListings(limit, offset)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	// get btc, btc rate:
	btcPrice, err := helpers.GetExternalPrice("BTC")
	if err != nil {
		u.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
		return nil, err
	}

	u.Logger.Info("btcPrice", btcPrice)
	ethPrice, err := helpers.GetExternalPrice("ETH")
	if err != nil {
		u.Logger.ErrorAny("convertBTCToETH", zap.Error(err))
		return nil, err
	}
	u.Logger.Info("btcPrice", btcPrice)

	for _, listing := range nftList {

		// get pyament listing info:
		paymentListingInfo, err := u.GetListingPaymentInfoWithEthBtcPrice(listing.PayType, listing.Price, btcPrice, ethPrice)
		if err != nil {
			return nil, err
		}

		// if listing.IsSold {
		// 	if !buyableOnly {
		// 		nftInfo := structure.MarketplaceNFTDetail{
		// 			InscriptionID: listing.InscriptionID,
		// 			Name:          listing.Name,
		// 			Description:   listing.Description,
		// 			Price:         listing.Price,
		// 			OrderID:       listing.UUID,
		// 			IsConfirmed:   listing.IsConfirm,
		// 			Buyable:       false,
		// 			IsCompleted:   listing.IsSold,
		// 			CreatedAt:     listing.CreatedAt,
		// 		}
		// 		result = append(result, nftInfo)
		// 	}
		// 	continue
		// }
		buyOrders, err := u.Repo.GetBTCListingHaveOngoingOrder(listing.UUID)
		if err != nil {

			if !buyableOnly {
				nftInfo := structure.MarketplaceNFTDetail{
					InscriptionID: listing.InscriptionID,
					Name:          listing.Name,
					Description:   listing.Description,
					Price:         listing.Price,
					OrderID:       listing.UUID,
					IsConfirmed:   listing.IsConfirm,
					Buyable:       false,
					IsCompleted:   listing.IsSold,
					CreatedAt:     listing.CreatedAt,

					Inscription:      listing.Inscription,
					InscriptionName:  listing.InscriptionName,
					InscriptionIndex: listing.InscriptionIndex,
					CollectionID:     listing.CollectionID,

					PaymentListingInfo: paymentListingInfo,
				}
				inscribeInfo, err := u.GetInscribeInfo(nftInfo.InscriptionID)
				if err != nil {
					u.Logger.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
				}
				if inscribeInfo != nil {
					nftInfo.InscriptionNumber = inscribeInfo.Index
					nftInfo.ContentType = inscribeInfo.ContentType
					nftInfo.ContentLength = inscribeInfo.ContentLength
				}
				result = append(result, nftInfo)
				continue
			}
		}
		currentTime := time.Now()
		isAvailable := true
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isAvailable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isAvailable = false
				break
			}
		}

		nftInfo := structure.MarketplaceNFTDetail{
			InscriptionID: listing.InscriptionID,
			Name:          listing.Name,
			Description:   listing.Description,
			Price:         listing.Price,
			OrderID:       listing.UUID,
			IsConfirmed:   listing.IsConfirm,
			Buyable:       isAvailable,
			IsCompleted:   listing.IsSold,
			CreatedAt:     listing.CreatedAt,

			Inscription: listing.Inscription,

			InscriptionName:  listing.InscriptionName,
			InscriptionIndex: listing.InscriptionIndex,
			CollectionID:     listing.CollectionID,

			PaymentListingInfo: paymentListingInfo,
		}
		inscribeInfo, err := u.GetInscribeInfo(nftInfo.InscriptionID)
		if err != nil {
			u.Logger.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
		}
		if inscribeInfo != nil {
			nftInfo.InscriptionNumber = inscribeInfo.Index
			nftInfo.ContentType = inscribeInfo.ContentType
			nftInfo.ContentLength = inscribeInfo.ContentLength
		}
		if buyableOnly && isAvailable {
			result = append(result, nftInfo)
		}
		if !buyableOnly {
			result = append(result, nftInfo)
		}
	}

	if !buyableOnly {
		sort.SliceStable(result, func(i, j int) bool {
			if result[i].Buyable && result[j].Buyable {
				if result[i].CreatedAt.After(result[j].CreatedAt) {
					return true
				}
			}
			return result[i].Buyable
		})
	}

	// result := []response.MarketplaceNFTDetail{}
	// for _, nft := range nfts {

	// 	result = append(result, nftInfo)
	// }
	return result, nil
}

func (u Usecase) BTCMarketplaceBuyOrder(orderInfo structure.MarketplaceBTC_BuyOrderInfo) (*entity.MarketplaceBTCBuyOrder, error) {

	order := entity.MarketplaceBTCBuyOrder{
		InscriptionID: orderInfo.InscriptionID,
		ItemID:        orderInfo.OrderID,
		OrdAddress:    orderInfo.BuyOrdAddress,
		ExpiredAt:     time.Now().Add(time.Minute * 30),
		PayType:       orderInfo.PayType,
	}

	// privKey, _, addressSegwit, err := btc.GenerateAddressSegwit()
	// if err != nil {
	// 	u.Logger.Error("u.OrdService.Exec.create.receive", err.Error(), err)
	// 	return nil, err
	// }
	// order.SegwitAddress = addressSegwit
	// order.SegwitKey = privKey

	// order.HoldOrdAddress = holdOrdAddress
	// sendMessage := func( offer entity.MarketplaceOffers) {
	// //

	// 	profile, err := u.Repo.FindUserByWalletAddress(offer.Buyer)
	// 	if err != nil {
	// 		u.Logger.Error("cancelListing.FindUserByWalletAddress", err.Error(), err)
	// 		return
	// 	}

	// 	token, err := u.Repo.FindTokenByGenNftAddr(offer.CollectionContract, offer.TokenId)
	// 	if err != nil {
	// 		u.Logger.Error("cancelListing.FindTokenByGenNftAddr", err.Error(), err)
	// 		return
	// 	}

	// 	preText := fmt.Sprintf("[OfferID %s] has been created by %s", offer.OfferingId, offer.Buyer)
	// 	content := fmt.Sprintf("TokenID: %s", helpers.CreateTokenLink(token.ProjectID, token.TokenID, token.Name))
	// 	title := fmt.Sprintf("User %s made offer with %s", helpers.CreateProfileLink(profile.WalletAddress, profile.DisplayName), offer.Price)

	// 	if _, _, err := u.Slack.SendMessageToSlack(preText, title, content); err != nil {
	// 		u.Logger.Error("s.Slack.SendMessageToSlack err", err.Error(), err)
	// 	}

	// }

	listing, _ := u.Repo.FindBtcNFTListingByOrderIDValid(orderInfo.OrderID)

	if listing == nil {
		err := errors.New("the listing is invalid")
		u.Logger.Error("u.FindBtcNFTListingByOrderIDValid.Check(listing)", err.Error(), err)
		return nil, err
	}

	// verify paytype:
	if orderInfo.PayType != utils.NETWORK_BTC && orderInfo.PayType != utils.NETWORK_ETH {
		err := errors.New("only support payType is eth or btc")
		u.Logger.Error("u.BTCMarketplaceListNFT.Check(payType)", err.Error(), err)
		return nil, err
	}

	// cal min price:
	priceStr := "0"
	priceInt, err := strconv.Atoi(listing.Price)
	if err != nil {
		u.Logger.Error("u.BTCMarketplaceListNFT.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}

	var btcRate, ethRate float64
	var privateKey, receiveAddress string

	// check type:
	if orderInfo.PayType == utils.NETWORK_BTC {
		privateKey, _, receiveAddress, err = btc.GenerateAddressSegwit()
		if err != nil {
			u.Logger.Error("u.BTCMarketplaceListNFT.GenerateAddressSegwit", err.Error(), err)
			return nil, err
		}
		priceStr = strconv.Itoa(priceInt)

		_, btcRate, ethRate, err = u.convertBTCToETH("1")
		if err != nil {
			u.Logger.Error("convertBTCToETH", err.Error(), err)
			return nil, err
		}

	} else if orderInfo.PayType == utils.NETWORK_ETH {
		ethClient := eth.NewClient(nil)

		privateKey, _, receiveAddress, err = ethClient.GenerateAddress()
		if err != nil {
			u.Logger.Error("BTCMarketplaceListNFT.ethClient.GenerateAddress", err.Error(), err)
			return nil, err
		}
		priceStr, btcRate, ethRate, err = u.convertBTCToETH(fmt.Sprintf("%f", float64(priceInt)/1e8))
		if err != nil {
			u.Logger.Error("convertBTCToETH", err.Error(), err)
			return nil, err
		}
		fmt.Println("priceStr ETH: ", priceStr)
	}

	if len(receiveAddress) == 0 || len(privateKey) == 0 {
		err = errors.New("can not create the wallet")
		u.Logger.Error("u.BTCMarketplaceListNFT.GenerateAddress", err.Error(), err)
		return nil, err
	}

	// set temp wallet info:
	order.PayType = orderInfo.PayType

	if len(os.Getenv("SECRET_KEY")) == 0 {
		err = errors.New("please config SECRET_KEY")
		u.Logger.Error("u.BTCMarketplaceListNFT.GenerateAddress", err.Error(), err)
		return nil, err
	}

	privateKeyEnCrypt, err := encrypt.EncryptToString(privateKey, os.Getenv("SECRET_KEY"))
	if err != nil {
		u.Logger.Error("u.BTCMarketplaceListNFT.Encrypt", err.Error(), err)
		return nil, err
	}

	order.PrivateKey = privateKeyEnCrypt
	order.ReceiveAddress = receiveAddress

	order.Price = priceStr
	order.EthRate = ethRate
	order.BtcRate = btcRate

	// check if listing is created or not
	err = u.Repo.CreateMarketplaceBuyOrder(&order)
	if err != nil {
		u.Logger.Error("BTCMarketplaceListNFT.Repo.CreateMarketplaceListingBTC", "", err)
		return nil, err
	}
	return &order, nil
}

// api listing detail:
func (u Usecase) BTCMarketplaceListingDetail(inscriptionID string) (*structure.MarketplaceNFTDetail, error) {

	var nft *entity.MarketplaceBTCListing
	var err error
	isBuyable := true
	isCompleted := false
	// lastPrice := int64(0)

	nft, err = u.Repo.FindBtcNFTListingUnsoldByNFTID(inscriptionID)
	if err != nil {
		isBuyable = false
		nft, err = u.Repo.FindBtcNFTListingLastSoldByNFTID(inscriptionID)

		if err != nil {
			u.Logger.Error("FindBtcNFTListingLastSoldByNFTID", err.Error(), err)
			return nil, err
		}
		isCompleted = true
	}

	if nft == nil {
		return nil, errors.New("nft not found")
	}

	if !nft.IsSold {
		buyOrders, err := u.Repo.GetBTCListingHaveOngoingOrder(nft.UUID)
		if err != nil {
			u.Logger.Error("h.Usecase.Repo.GetBTCListingHaveOngoingOrder", err.Error(), err)
		}
		currentTime := time.Now()
		for _, order := range buyOrders {
			expireTime := order.ExpiredAt
			// not expired yet still waiting for btc
			if currentTime.Before(expireTime) && (order.Status == entity.StatusBuy_Pending || order.Status == entity.StatusBuy_NotEnoughBalance) {
				isBuyable = false
				break
			}
			// could be expired but received btc
			if order.Status != entity.StatusBuy_Pending && order.Status != entity.StatusBuy_NotEnoughBalance {
				isBuyable = false
				break
			}
		}
	}

	// get pyament listing info:
	paymentListingInfo, err := u.getListingPaymentInfo(nft.PayType, nft.Price)
	if nft == nil {
		return nil, errors.New("nft not found")
	}

	nftInfo := &structure.MarketplaceNFTDetail{
		InscriptionID:      nft.InscriptionID,
		Name:               nft.Name,
		Description:        nft.Description,
		Price:              nft.Price,
		OrderID:            nft.UUID,
		IsConfirmed:        nft.IsConfirm,
		Buyable:            isBuyable,
		IsCompleted:        isCompleted,
		PaymentListingInfo: paymentListingInfo,
		// LastPrice:     lastPrice,
	}
	inscribeInfo, err := u.GetInscribeInfo(nftInfo.InscriptionID)
	if err != nil {
		u.Logger.Error("h.Usecase.GetInscribeInfo", err.Error(), err)
	}
	if inscribeInfo != nil {
		nftInfo.InscriptionNumber = inscribeInfo.Index
		nftInfo.ContentType = inscribeInfo.ContentType
		nftInfo.ContentLength = inscribeInfo.ContentLength
	}
	// h.Logger.Info("resp.Proposal", resp)

	return nftInfo, nil
}

// get filter info:
type CollectionFilter struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Value string `json:"value"`
	From  int64  `json:"from"`
	To    int64  `json:"to"`
}

func (u Usecase) JobMKP_CrawlToUpdateNftInfo() error {
	// create data for listing data:
	// get nft list + update collection:
	listingOrder, _ := u.Repo.RetrieveBTCNFTListingsUnsold(99999, 0)

	fmt.Println("len(listingOrder): ", len(listingOrder))

	for _, v := range listingOrder {
		// get nft collection info:
		nftCollectionInfo, _ := u.Repo.FindTokenByTokenID(v.InscriptionID)
		if nftCollectionInfo != nil {
			fmt.Println("can not get nftCollectionInfo with v.InscriptionID: ", v.InscriptionID)
			_, err := u.Repo.UpdateListingCollectionInfo(v.UUID, nftCollectionInfo)
			if err != nil {
				fmt.Println("can not UpdateListingCollectionInfo err: ", err)
			} else {
				fmt.Println("update done: ", nftCollectionInfo.TokenID)
			}
		} else {
			fmt.Println("can not found id", v.InscriptionID)
		}
	}
	return nil
}

func (u Usecase) BTCMarketplaceFilterInfo() (interface{}, error) {

	listCollectionFilterMap := map[string]CollectionFilter{}
	listPriceFilterMap := map[string]CollectionFilter{}
	listNftIDFilterMap := map[string]CollectionFilter{}

	listPriceFilterMap["range1"] = CollectionFilter{Name: "0 BTC - 2 BTC", Count: 0, Value: "0_200000000", From: 0, To: 200000000}
	listPriceFilterMap["range2"] = CollectionFilter{Name: "2 BTC - 5 BTC", Count: 0, Value: "200000000_500000000", From: 200000000, To: 500000000}
	listPriceFilterMap["range3"] = CollectionFilter{Name: "5 BTC - 7 BTC", Count: 0, Value: "500000000_700000000", From: 500000000, To: 700000000}
	listPriceFilterMap["range4"] = CollectionFilter{Name: "7 BTC - 9 BTC", Count: 0, Value: "700000000_900000000", From: 700000000, To: 900000000}
	listPriceFilterMap["range5"] = CollectionFilter{Name: "> 9 BTC", Count: 0, Value: "900000000_999999999999", From: 900000000, To: 999999999999}

	listNftIDFilterMap["range1"] = CollectionFilter{Name: "#0 - #1000", Count: 0, Value: "0_1000", From: 0, To: 1000}
	listNftIDFilterMap["range2"] = CollectionFilter{Name: "#1000 - #2000", Count: 0, Value: "1000_2000", From: 1000, To: 2000}
	listNftIDFilterMap["range3"] = CollectionFilter{Name: "#2000 - #3000", Count: 0, Value: "2000_3000", From: 2000, To: 3000}
	listNftIDFilterMap["range4"] = CollectionFilter{Name: "#3000 - #4000", Count: 0, Value: "3000_4000", From: 3000, To: 4000}
	listNftIDFilterMap["range5"] = CollectionFilter{Name: "> #4000", Count: 0, Value: "4000_9999999", From: 4000, To: 9999999}

	listCollectionFilterMapReturn := map[string]interface{}{}

	listingOrder, _ := u.Repo.RetrieveBTCNFTListingsUnsold(99999, 1) //TODO: @tri check this why offset is 1
	for _, v := range listingOrder {
		collectionFilter := CollectionFilter{
			Name:  v.CollectionName,
			Count: 1,
			Value: v.CollectionID,
		}

		val, ok := listCollectionFilterMap[v.CollectionID]
		// If the key exists
		if ok {
			collectionFilter.Count = val.Count + 1
			listCollectionFilterMap[v.CollectionID] = collectionFilter
		}
		listCollectionFilterMap[v.CollectionID] = collectionFilter

		var price int64
		if len(v.Price) > 0 {
			price, _ = strconv.ParseInt(v.Price, 10, 64)
		}

		var nftID int64
		if len(v.InscriptionIndex) > 0 {
			nftID, _ = strconv.ParseInt(v.InscriptionIndex, 10, 64)
		}

		// TODO important: if len(listPriceFilterMap) != len(listNftIDFilterMap) please use 2 for loop.
		for key, _ := range listPriceFilterMap {
			filterPrice := listPriceFilterMap[key]

			fmt.Println("filterPrice: ", filterPrice.From, filterPrice.To)

			// for price
			if listPriceFilterMap[key].From <= price && price < listPriceFilterMap[key].To {
				filterPrice.Count += 1
			}
			listPriceFilterMap[key] = filterPrice

			// for nftID:
			filterNftID := listNftIDFilterMap[key]
			// for price
			if listNftIDFilterMap[key].From < nftID && nftID <= listNftIDFilterMap[key].To {
				filterNftID.Count += 1
			}
			listNftIDFilterMap[key] = filterNftID
		}
	}

	listCollectionFilterMapReturn["collection"] = listCollectionFilterMap
	listCollectionFilterMapReturn["price"] = listPriceFilterMap
	listCollectionFilterMapReturn["inscriptionID"] = listNftIDFilterMap

	return listCollectionFilterMapReturn, nil
}

func (u Usecase) GetCollectionMarketplaceStats(collectionID string) (*structure.MarketplaceCollectionStats, error) {
	var result structure.MarketplaceCollectionStats
	floorPrice, err := u.Repo.RetrieveFloorPriceOfCollection(collectionID)
	if err != nil {
		return nil, err
	}
	result.FloorPrice = floorPrice
	return &result, nil
}

func (u Usecase) getListingPaymentInfo(payType map[string]string, btcPrice string) (map[string]structure.PaymentInfoForBuyOrder, error) {

	paymentListingInfo := map[string]structure.PaymentInfoForBuyOrder{}

	addressBtc, okBtc := payType["btc"]

	if okBtc {
		paymentListingInfo["btc"] = structure.PaymentInfoForBuyOrder{
			PaymentAddress: addressBtc,
			Price:          btcPrice,
		}
	}

	addressEth, okEth := payType["eth"]
	if okEth {
		priceStr := "0"
		priceInt, err := strconv.Atoi(btcPrice)
		if err != nil {
			u.Logger.Error("u.BTCMarketplaceListNFT.GetBTCListingHaveOngoingOrder", err.Error(), err)
			return nil, err
		}
		priceStr, _, _, err = u.convertBTCToETH(fmt.Sprintf("%f", float64(priceInt)/1e8))
		if err != nil {
			u.Logger.Error("convertBTCToETH", err.Error(), err)
			return nil, err
		}
		fmt.Println("priceStr ETH: ", priceStr)
		paymentListingInfo["eth"] = structure.PaymentInfoForBuyOrder{
			PaymentAddress: addressEth,
			Price:          priceStr,
		}
	}
	// end get pyament listing info

	return paymentListingInfo, nil

}

func (u Usecase) GetListingPaymentInfoWithEthBtcPrice(payType map[string]string, btcPrice string, btcUsdPrice, ethUsdPrice float64) (map[string]structure.PaymentInfoForBuyOrder, error) {

	paymentListingInfo := map[string]structure.PaymentInfoForBuyOrder{}

	addressBtc, okBtc := payType["btc"]

	// okBtc = true

	if okBtc {
		paymentListingInfo["btc"] = structure.PaymentInfoForBuyOrder{
			PaymentAddress: addressBtc,
			Price:          btcPrice,
		}
	}

	addressEth, okEth := payType["eth"]

	// okEth = true

	if okEth {
		priceStr := "0"
		priceInt, err := strconv.Atoi(btcPrice)
		if err != nil {
			u.Logger.Error("u.BTCMarketplaceListNFT.GetBTCListingHaveOngoingOrder", err.Error(), err)
			return nil, err
		}
		priceStr, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(priceInt)/1e8), btcUsdPrice, ethUsdPrice)
		if err != nil {
			u.Logger.Error("convertBTCToETH", err.Error(), err)
			return nil, err
		}
		fmt.Println("priceStr ETH: ", priceStr)
		paymentListingInfo["eth"] = structure.PaymentInfoForBuyOrder{
			PaymentAddress: addressEth,
			Price:          priceStr,
		}
	}
	// end get pyament listing info

	return paymentListingInfo, nil

}
