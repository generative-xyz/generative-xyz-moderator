package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CancelDexBTCListing(txhash string, seller_address string, inscription_id string, order_id string) error {
	orderInfo, err := u.Repo.GetDexBTCListingOrderByID(order_id)
	if err != nil {
		return err
	}
	if orderInfo.InscriptionID != inscription_id {
		return errors.New("invalid cancelling request")
	}

	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}

	inscriptionInfo, err := getInscriptionByID(ordServer, inscription_id)
	if err != nil {
		fmt.Printf("Could not get inscription info - with err: %v", err)
		return err
	}

	if !strings.EqualFold(inscriptionInfo.Address, seller_address) {
		return fmt.Errorf("seller address not match inscription owner")
	}

	if !orderInfo.Cancelled && orderInfo.CancelTx == "" {
		currentTime := time.Now()
		orderInfo.CancelTx = txhash
		orderInfo.CancelAt = &currentTime
	} else {
		return errors.New("order already cancelling/cancelled")
	}
	_, err = u.Repo.UpdateDexBTCListingOrderCancelTx(orderInfo)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) DexBTCListing(seller_address string, raw_psbt string, inscription_id string, split_tx string) (*entity.DexBTCListing, error) {
	newListing := entity.DexBTCListing{
		RawPSBT:       raw_psbt,
		InscriptionID: inscription_id,
		SellerAddress: seller_address,
		Cancelled:     false,
		CancelTx:      "",
		SplitTx:       split_tx,
	}

	psbtData, err := btc.ParsePSBTFromBase64(raw_psbt)
	if err != nil {
		return nil, err
	}
	var splitTxData *wire.MsgTx
	if split_tx != "" {
		splitTxData, err = btc.ParseTx(split_tx)
		if err != nil {
			return nil, err
		}
	}

	outputList, err := extractAllOutputFromPSBT(psbtData)
	if err != nil {
		return nil, err
	}

	totalOuputValue := uint64(0)
	for _, output := range psbtData.UnsignedTx.TxOut {
		totalOuputValue += uint64(output.Value)
	}

	if len(psbtData.Inputs) == 1 {
		newListing.Amount = totalOuputValue
	} else {
		newListing.Amount = totalOuputValue - uint64(psbtData.Inputs[1].WitnessUtxo.Value)
	}

	txInputs := []string{}
	for _, input := range psbtData.UnsignedTx.TxIn {
		i := fmt.Sprintf("%v:%v", input.PreviousOutPoint.Hash.String(), input.PreviousOutPoint.Index)
		txInputs = append(txInputs, i)
	}
	newListing.Inputs = txInputs

	artistAddress := ""
	royaltyFeePercent := float64(0)
	internalInfo, _ := u.Repo.FindTokenByTokenID(inscription_id)
	if internalInfo != nil {
		projectDetail, _ := u.Repo.FindProjectByTokenID(internalInfo.ProjectID)
		creator, err := u.GetUserProfileByWalletAddress(projectDetail.CreatorAddrr)
		if err == nil {
			if creator.WalletAddressBTC != "" || creator.WalletAddressBTCTaproot != "" {
				royaltyFeePercent = float64(projectDetail.Royalty) / 10000
				// prioritize WalletAddressBTC address
				if creator.WalletAddressBTC != "" {
					artistAddress = creator.WalletAddressBTC
				} else {
					artistAddress = creator.WalletAddressBTCTaproot
				}
			}
		}
	}

	if artistAddress != "" && royaltyFeePercent > 0 {
		if len(psbtData.UnsignedTx.TxOut) == 1 {
			//force receiver == artistAddress when only one output
			for receiver, _ := range outputList {
				if receiver != artistAddress {
					return nil, fmt.Errorf("expected to pay royalty fees to %v", artistAddress)
				}
			}
		} else {
			royaltyFeeExpected := int64(float64(psbtData.UnsignedTx.TxOut[0].Value) * royaltyFeePercent)
			for receiver, outputs := range outputList {
				if receiver == artistAddress {
					totalValue := int64(0)
					for _, output := range outputs {
						totalValue += output.Value
					}
					if totalValue < royaltyFeeExpected {
						return nil, fmt.Errorf("expected royalty fees of artist %v to be %v, got %v", artistAddress, royaltyFeeExpected, totalValue)
					}
				}
			}
		}
	}

	previousTxs, err := retrievePreviousTxFromPSBT(psbtData)
	if err != nil {
		return nil, err
	}

	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}

	inscriptionInfo, err := getInscriptionByID(ordServer, inscription_id)
	if err != nil {
		fmt.Printf("Could not get inscription info - with err: %v", err)
		return nil, err
	}

	if !strings.EqualFold(inscriptionInfo.Address, seller_address) {
		return nil, fmt.Errorf("seller address not match inscription owner")
	}

	inscriptionTx := strings.Split(inscriptionInfo.Satpoint, ":")[0]

	if inscriptionTx != previousTxs[0] {
		found := false
		if splitTxData != nil {
			for _, input := range splitTxData.TxIn {
				if inscriptionTx == input.PreviousOutPoint.Hash.String() {
					found = true
					break
				}
			}
		}
		if !found {
			return nil, errors.New("can't found inscription in split tx")
		}
	} else {
		newListing.Verified = true
	}

	if split_tx != "" {
		// _, err = btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI)
		// if err != nil {
		// 	fmt.Printf("btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI) - with err: %v %v\n", err, split_tx)
		// 	return nil, err
		// }
		txMap := make(map[string]string)
		txMap[splitTxData.TxHash().String()] = split_tx
		err = u.SubmitBTCTransaction(txMap)
		if err != nil {
			log.Println("httpDelivery.submitTx.SubmitBTCTransaction", err.Error())
			return nil, err
		}
		newListing.Verified = false
	}

	return &newListing, u.Repo.CreateDexBTCListing(&newListing)
}

func retrievePreviousTxFromPSBT(psbtData *psbt.Packet) ([]string, error) {
	result := []string{}
	for _, input := range psbtData.UnsignedTx.TxIn {
		result = append(result, input.PreviousOutPoint.Hash.String())
	}
	return result, nil
}

func extractAllOutputFromPSBT(psbtData *psbt.Packet) (map[string][]*wire.TxOut, error) {
	result := make(map[string][]*wire.TxOut)
	for _, output := range psbtData.UnsignedTx.TxOut {
		address, err := btc.GetAddressFromPKScript(output.PkScript)
		if err != nil {
			return nil, err
		}
		result[address] = append(result[address], output)
	}
	return result, nil
}
func (u Usecase) JobWatchPendingDexBTCListing() error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingDexBTCListing()
		if err != nil {
			log.Println("JobWatchPendingDexBTCListing watchPendingDexBTCListing err", err)
		}
	}(&wg)

	wg.Wait()
	return nil
}

func (u Usecase) JobWatchPendingDexBTCBuyETH() error {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingDexBTCBuyETH()
		if err != nil {
			log.Println("JobWatchPendingDexBTCListing watchPendingDexBTCBuyETH err", err)
		}
	}(&wg)

	wg.Wait()
	return nil
}

func (u Usecase) InsertDexVolumnInscription(o entity.DexBTCListing) {
	logger.AtLog.Logger.Info("DexVolumeInscription Insert to time series data %s", zap.Any("o.InscriptionID", o.InscriptionID))
	data := entity.DexVolumeInscription{
		Amount:    o.Amount,
		Timestamp: o.MatchAt,
		Metadata: entity.DexVolumeInscriptionMetadata{
			InscriptionId: o.InscriptionID,
			MatchedTx:     o.MatchedTx,
		},
	}
	err := u.Repo.InsertDexVolumeInscription(&data)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("DexVolumeInscription Error Insert %s to time series data", o.InscriptionID), zap.Any("error", err))
	} else {
		o.IsTimeSeriesData = true
		_, err = u.Repo.UpdateDexBTCListingTimeseriesData(&o)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("DexVolumeInscription Error Insert %s to time series data - UpdateDexBTCListingTimeseriesData", o.InscriptionID), zap.Any("error", err))
		}
	}
}

func (u Usecase) preCheckPendingDexBTCListingTx(pendingOrders []entity.DexBTCListing) (map[string]*btc.GoBCYMultiTx, error) {
	txNeedToCheck := []string{}
	for _, order := range pendingOrders {
		if order.CancelTx == "" {
			inscriptionTx := strings.Split(order.Inputs[0], ":")
			txNeedToCheck = append(txNeedToCheck, inscriptionTx[0])
			if len(order.Inputs) > 1 {
				for idx, v := range order.Inputs {
					if idx == 0 {
						continue
					}
					tx := strings.Split(v, ":")
					txNeedToCheck = append(txNeedToCheck, tx[0])
				}
			}
		}
	}
	log.Println("preCheckPendingDexBTCListingTx len(txNeedToCheck)", len(txNeedToCheck))
	result, _, err := btc.CheckTxMultiBlockcypher(txNeedToCheck, u.Config.DEXBTCBlockcypherToken)
	if err != nil {
		log.Println("preCheckPendingDexBTCListingTx CheckTxMultiBlockcypher", err.Error())
		return nil, err
	}
	if len(result) == 0 {
		result = make(map[string]*btc.GoBCYMultiTx)
	}
	log.Println("preCheckPendingDexBTCListingTx len(result)", len(result))
	return result, nil
}

func (u Usecase) watchPendingDexBTCListing() error {
	pendingOrders, err := u.Repo.GetDexBTCListingOrderPending()
	if err != nil {
		return err
	}

	preCheckTxs, err := u.preCheckPendingDexBTCListingTx(pendingOrders)
	if err != nil {
		log.Println("JobWatchPendingDexBTCListing preCheckPendingDexBTCListingTx err", err.Error())
	}
	for _, order := range pendingOrders {
		if len(order.Inputs) == 0 {
			continue
		}
		inscriptionTx := strings.Split(order.Inputs[0], ":")
		idx, err := strconv.Atoi(inscriptionTx[1])
		if err != nil {
			log.Printf("JobWatchPendingDexBTCListing strconv.Atoi(inscriptionTx[1]) %v\n", order.Inputs)
			continue
		}
		if !order.Verified {
			txDetail, err := btc.CheckTxfromQuickNode(inscriptionTx[0], u.Config.QuicknodeAPI)
			if err != nil {
				log.Printf("btc.GetBTCTxStatusExtensive %v\n", err)
			} else {
				if txDetail.Result.Confirmations > 0 {
					order.Verified = true
					_, err = u.Repo.UpdateDexBTCListingOrderConfirm(&order)
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderConfirm err %v\n", err)
					}
				}
			}
		}
		if order.CancelTx == "" {
			spentTx := ""
			// txStatus, err := bs.CheckTx(inscriptionTx[0])
			txStatus, exist := preCheckTxs[inscriptionTx[0]]
			if !exist {
				//log.Printf("JobWatchPendingDexBTCListing bs.CheckTx(txhash) %v\n", order.Inputs)
				spentTx, err = btc.CheckOutcoinSpentBlockStream(inscriptionTx[0], uint(idx))
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream %v\n", order.Inputs)
					continue
				}
				if spentTx != "" {
					//log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream success\n")
				}
				if spentTx == "" {
					//log.Printf("JobWatchPendingDexBTCListing GetInscriptionByIDFromOrd %v\n", order.InscriptionID)
					inscriptionInfo, err := u.GetInscriptionByIDFromOrd(order.InscriptionID)
					if err != nil {
						//log.Printf("JobWatchPendingDexBTCListing GetInscriptionByIDFromOrd %v\n", order.Inputs)
						continue
					}
					if inscriptionInfo != nil {
						found := false
						inscPoint := strings.Split(inscriptionInfo.Satpoint, ":")
						if strings.EqualFold(strings.Join([]string{inscPoint[0], inscPoint[1]}, "i"), order.InscriptionID) {
							continue
						}

						curInscTx := inscPoint[0]
						for _, vIn := range order.Inputs {
							vInTx := strings.Split(vIn, ":")[0]
							if curInscTx == vInTx {
								found = true
								break
							}
						}
						if !found {
							spentTx = curInscTx
						}
					}
					if spentTx == "" {
						for _, vIn := range order.Inputs {
							vinParts := strings.Split(vIn, ":")
							vinIdx, _ := strconv.Atoi(vinParts[1])
							spentTx, err = btc.CheckOutcoinSpentBlockStream(vinParts[0], uint(vinIdx))
							if err != nil {
								log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream %v\n", order.Inputs)
								continue
							}

							if spentTx != "" {
								log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream success2\n")
							}
						}
					}
				}
			} else {
				if len(txStatus.Outputs) == 0 {
					continue
				}

				if txStatus.Outputs[idx].SpentBy != "" {
					spentTx = txStatus.Outputs[idx].SpentBy
				}
			}
			if strings.Contains(order.InscriptionID, spentTx) {
				continue
			}

			if spentTx != "" {
				spentTxDetail, err := btc.CheckTxfromQuickNode(spentTx, u.Config.QuicknodeAPI)
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckTxfromQuickNode(spentTx) %v %v\n", order.Inputs, err)
					continue
				}

				inputTxDetail, err := btc.CheckTxfromQuickNode(inscriptionTx[0], u.Config.QuicknodeAPI)
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckTxfromQuickNode(spentTx) %v %v\n", order.Inputs, err)
					continue
				}

				if inputTxDetail.Result.Confirmations <= 0 {
					continue
				}

				if spentTxDetail.Result.Blocktime <= inputTxDetail.Result.Blocktime {
					log.Printf("JobWatchPendingDexBTCListing blocktime not valid %v %v\n", spentTx, inscriptionTx[0])
					continue
				}

				isValidMatch := false
				if spentTxDetail != nil {
					psbtData, _ := btc.ParsePSBTFromBase64(order.RawPSBT)
					matchIns := 0
					matchOuts := 0
					for _, vin := range spentTxDetail.Result.Vin {
						input := fmt.Sprintf("%v:%v", vin.Txid, vin.Vout)
						for _, oin := range order.Inputs {
							if input == oin {
								matchIns += 1
								break
							}
						}
					}
					totalOutvalue := uint64(0)
					for _, vout := range spentTxDetail.Result.Vout {
						voutAddress := vout.ScriptPubKey.Address
						totalOutvalue += uint64(vout.Value * 1e8)
						for _, out := range psbtData.UnsignedTx.TxOut {
							pkAddress, _ := btc.GetAddressFromPKScript(out.PkScript)
							if pkAddress == voutAddress {
								matchOuts += 1
								break
							}
						}
					}
					if totalOutvalue >= order.Amount && matchIns == len(order.Inputs) && matchOuts >= len(psbtData.UnsignedTx.TxOut) {
						isValidMatch = true
					}
				} else {
					continue
				}

				if isValidMatch {
					currentTime := time.Now()
					order.MatchedTx = spentTx
					order.MatchAt = &currentTime
					order.Matched = true
					txDetail, err := btc.CheckTxfromQuickNode(spentTx, u.Config.QuicknodeAPI)
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC(spentTx) %v %v\n", order.Inputs, err)
					}
					output := txDetail.Result.Vout[0]
					order.Buyer = output.ScriptPubKey.Address

					_, err = u.Repo.UpdateDexBTCListingOrderMatchTx(&order)
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderMatchTx err %v\n", err)
						continue
					}

					// Insert time series data
					go func(o entity.DexBTCListing, userCase Usecase) {
						userCase.Logger.Info("DexVolumeInscription Insert to time series data %s", o.InscriptionID)
						data := entity.DexVolumeInscription{
							Amount:    o.Amount,
							Timestamp: o.MatchAt,
							Metadata: entity.DexVolumeInscriptionMetadata{
								InscriptionId: o.InscriptionID,
								MatchedTx:     o.MatchedTx,
							},
						}
						err := userCase.Repo.InsertDexVolumeInscription(&data)
						if err != nil {
							userCase.Logger.ErrorAny(fmt.Sprintf("DexVolumeInscription Error Insert %s to time series data", o.InscriptionID), zap.Any("error", err))
						} else {
							order.IsTimeSeriesData = true
							_, err = userCase.Repo.UpdateDexBTCListingTimeseriesData(&order)
							if err != nil {
								userCase.Logger.ErrorAny(fmt.Sprintf("DexVolumeInscription Error Insert %s to time series data - UpdateDexBTCListingTimeseriesData", o.InscriptionID), zap.Any("error", err))
							}
						}
					}(order, u)

					// Discord Notify NEW SALE
					go u.NotifyNewSale(order)
				} else {
					log.Printf("JobWatchPendingDexBTCListing not valid match err %v\n", err)
					txDetail, err := btc.CheckTxfromQuickNode(spentTx, u.Config.QuicknodeAPI)
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC(spentTx) %v %v\n", order.Inputs, err)
					}
					output := txDetail.Result.Vout[0]
					if output.ScriptPubKey.Address == order.SellerAddress {
						currentTime := time.Now()
						order.CancelAt = &currentTime
						order.CancelTx = spentTx
						order.Cancelled = true
						_, err = u.Repo.UpdateDexBTCListingOrderCancelTx(&order)
						if err != nil {
							log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderCancelTx err %v\n", err)
							continue
						}
					} else {
						currentTime := time.Now()
						order.CancelAt = &currentTime
						order.Cancelled = true
						order.InvalidMatch = true
						order.InvalidMatchTx = spentTx
						_, err = u.Repo.UpdateDexBTCListingOrderInvalidMatch(&order)
						if err != nil {
							log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderCancelTx err %v\n", err)
							continue
						}
					}
				}
			}
			if len(order.Inputs) > 1 {
				for idx, vin := range order.Inputs {
					if idx == 0 {
						continue
					}
					tx := strings.Split(vin, ":")
					idx, err := strconv.Atoi(tx[1])
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing2 strconv.Atoi(tx[1]) %v\n", order.Inputs)
						continue
					}
					txStatus, exist := preCheckTxs[tx[0]]
					if exist {
						if len(txStatus.Outputs) == 0 {
							continue
						}

						if txStatus.Outputs[idx].SpentBy != "" {
							spentTx = txStatus.Outputs[idx].SpentBy

							currentTime := time.Now()
							order.CancelAt = &currentTime
							order.Cancelled = true
							order.InvalidMatch = true
							order.InvalidMatchTx = spentTx
							_, err = u.Repo.UpdateDexBTCListingOrderInvalidMatch(&order)
							if err != nil {
								log.Printf("JobWatchPendingDexBTCListing2 UpdateDexBTCListingOrderCancelTx err %v\n", err)
								continue
							}
						}
					}

				}
			}
		} else {
			status, err := btc.GetBTCTxStatusExtensive(order.CancelTx, nil, u.Config.QuicknodeAPI)
			if err != nil {
				log.Printf("JobWatchPendingDexBTCListing btc.GetBTCTxStatusExtensive err %v\n", err)
				continue
			}
			if status == "Pending" {
				continue
			}
			if status == "Success" {
				order.Cancelled = true
			}
			if status == "Failed" {
				order.CancelAt = nil
				order.CancelTx = ""
				order.Cancelled = false
			}
			_, err = u.Repo.UpdateDexBTCListingOrderCancelTx(&order)
			if err != nil {
				log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderCancelTx err %v\n", err)
				continue
			}
		}
	}
	return nil
}

// func (u Usecase) DexBTCBuyWithETH(userID string, orderID string, txhash string, feeRate uint64) error {
// 	newListing := entity.DexBTCBuyWithETH{
// 		OrderID: orderID,
// 		Txhash:  txhash,
// 		FeeRate: feeRate,
// 		UserID:  userID,
// 		Status:  entity.StatusDEXBuy_Pending,
// 	}

// 	return u.Repo.CreateDexBTCBuyWithETH(&newListing)
// }

// list address by:
func (u Usecase) ListBuyAddress() (interface{}, error) {

	type AddressObject struct {
		Uuid, Address string
		Status        int
	}

	var listAddrssObject []AddressObject

	listItem, err := u.Repo.ListAllDexBTCBuyETH()
	if err != nil {
		return nil, err
	}
	if len(listItem) == 0 {
		return nil, fmt.Errorf("listItem is empty")
	}
	for _, v := range listItem {

		_, ethAddress, err := u.EthClientDex.GenerateAddressFromPrivKey(v.ETHKey)
		if err != nil {
			continue
		}

		listAddrssObject = append(listAddrssObject, AddressObject{
			Uuid:    v.UUID,
			Address: ethAddress,
			Status:  int(v.Status),
		})
	}
	return listAddrssObject, nil
}

func (u Usecase) watchPendingDexBTCBuyETH() error {
	pendingOrders, err := u.Repo.GetDexBTCBuyETHOrderByStatus([]entity.DexBTCETHBuyStatus{entity.StatusDEXBuy_Pending, entity.StatusDEXBuy_ReceivedFund, entity.StatusDEXBuy_Buying, entity.StatusDEXBuy_WaitingToRefund, entity.StatusDEXBuy_Refunding, entity.StatusDEXBuy_Bought, entity.StatusDEXBuy_SendingMaster})
	if err != nil {
		return err
	}

	quickNodeAPI := u.Config.QuicknodeAPI
	currentBlockHeight, err := u.EthClientDex.GetClient().BlockNumber(context.Background())
	if err != nil {
		log.Printf("watchPendingDexBTCBuyETH BlockNumber err %v\n", err)
		return err
	}
	for _, order := range pendingOrders {
		switch order.Status {
		case entity.StatusDEXBuy_Pending:
			// check wallet receive enough funds
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, ethAddress, err := u.EthClientDex.GenerateAddressFromPrivKey(order.ETHKey)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GenerateAddressFromPrivKey", order.ID, err)
				continue
			}
			amountRequired, ok := new(big.Int).SetString(order.AmountETH, 10)
			if !ok {
				log.Println("watchPendingDexBTCBuyETH new(bigInt)", order.ID, order.AmountETH, err)
				continue
			}
			balance, err := u.EthClientDex.GetBalance(ctx, ethAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GetBalanceFromQuickNode", order.ID, ethAddress, err)
				continue
			}
			if balance.Cmp(amountRequired) > -1 {
				order.Status = entity.StatusDEXBuy_ReceivedFund
				_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
				}
				continue
				// if order.Confirmation >= 1 {
				// 	order.Status = entity.StatusDEXBuy_ReceivedFund
				// 	_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				// 	if err != nil {
				// 		log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
				// 	}
				// 	continue
				// } else {
				// 	order.Confirmation += 1
				// 	_, err := u.Repo.UpdateDexBTCBuyETHOrderConfirmation(&order)
				// 	if err != nil {
				// 		log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderConfirmation err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
				// 	}
				// 	continue
				// }
			} else {
				// not enough funds
				if time.Since(order.ExpiredAt) > 1*time.Second {
					order.Status = entity.StatusDEXBuy_Expired
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
					}
					continue
				}
			}
		case entity.StatusDEXBuy_ReceivedFund:
			// send tx buy update status to StatusDEXBuy_Buying
			//TODO: 2077 remove this in the future
			if len(order.SellOrderList) == 0 {
				listingOrder, err := u.Repo.GetDexBTCListingOrderByID(order.OrderID)
				if err != nil {
					log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, err)
					continue
				}
				if listingOrder != nil {
					if listingOrder.Cancelled || listingOrder.Matched {
						order.Status = entity.StatusDEXBuy_WaitingToRefund
						_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
						if err != nil {
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
						}
						continue
					}
					address := u.Config.DexBTCWalletAddress

					walletInfo, err := btc.GetBalanceFromQuickNode(address, quickNodeAPI)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH GetBalanceFromQuickNode", order.ID, address, err)
						continue
					}
					utxos, err := btc.ConvertToUTXOType(walletInfo.Txrefs)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH ConvertToUTXOType", order.ID, err)
						continue
					}

					filteredUTXOs, err := u.filterUTXOInscription(utxos)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH filterUTXOInscription", order.ID, err)
						continue
					}

					psbt, err := btc.ParsePSBTFromBase64(listingOrder.RawPSBT)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
						continue
					}

					feeRate := order.FeeRate
					amountBTCFee := uint64(0)
					amountBTCFee = btc.EstimateTxFee(uint(len(listingOrder.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))
					filteredBalance := uint64(0)
					for _, v := range filteredUTXOs {
						filteredBalance += v.Value
					}
					if filteredBalance <= amountBTCFee+listingOrder.Amount {
						go u.NotifyWithChannel("C052CAWFB0D", "Insufficient fund", address, fmt.Sprintf("filteredBalance %v <= amountBTCFee %v + listingOrder.Amount %v", filteredBalance, amountBTCFee, listingOrder.Amount))
						time.Sleep(300 * time.Millisecond)
						continue
					}

					respondData, err := btc.CreatePSBTToBuyInscriptionViaAPI(u.Config.DexBTCBuyService, address, listingOrder.RawPSBT, order.ReceiveAddress, listingOrder.Amount, filteredUTXOs, order.FeeRate, amountBTCFee)
					if err != nil {
						go u.NotifyWithChannel("C052CAWFB0D", "Create buy ", address, fmt.Sprintf("filteredBalance %v <= amountBTCFee %v + listingOrder.Amount %v", filteredBalance, amountBTCFee, listingOrder.Amount))

						logData := make(map[string]interface{})
						logData["u.Config.DexBTCBuyService"] = u.Config.DexBTCBuyService
						logData["address"] = address
						logData["listingOrder.RawPSBT"] = listingOrder.RawPSBT
						logData["order.ReceiveAddress"] = order.ReceiveAddress
						logData["listingOrder.Amount"] = listingOrder.Amount
						logData["filteredUTXOs"] = filteredUTXOs
						logData["order.FeeRate"] = order.FeeRate
						logData["amountBTCFee"] = amountBTCFee
						logData["respondData"] = respondData
						logData["err"] = err.Error()

						u.Repo.CreateDexBTCLog(&entity.DexBTCLog{Function: "CreatePSBTToBuyInscriptionViaAPI", Data: logData})
						log.Println("watchPendingDexBTCBuyETH CreatePSBTToBuyInscription", order.ID, err)
						time.Sleep(300 * time.Millisecond)
						continue
					}
					if respondData.SplitTxRaw != "" {
						err = btc.SendTxBlockStream(respondData.SplitTxRaw)
						if err != nil {
							dataBytes, _ := json.Marshal(respondData)
							log.Println("watchPendingDexBTCBuyETH SendTxBlockStream SplitTxHex", order.ID, string(dataBytes), err)
							continue
						}
						time.Sleep(5 * time.Second)
					}
					err = btc.SendTxBlockStream(respondData.TxHex)
					if err != nil {
						dataBytes, _ := json.Marshal(respondData)
						log.Println("watchPendingDexBTCBuyETH SendTxBlockStream TxHex", order.ID, string(dataBytes), err)
						continue
					}
					order.BuyTx = respondData.TxID
					order.SplitTx = respondData.SplitTxID
					order.Status = entity.StatusDEXBuy_Buying
					_, err = u.Repo.UpdateDexBTCBuyETHOrderBuy(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderBuy err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
					}
					continue
				} else {
					// ?? order not exist
					order.Status = entity.StatusDEXBuy_WaitingToRefund
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
					}
					continue
				}
			} else {
				newStatus := order.Status
				feeRate := order.FeeRate

				amountBTCFee := uint64(0)
				amountBTC := uint64(0)
				var buyReqInfos []btc.BuyReqInfo
				for _, listingOrderID := range order.SellOrderList {
					listingOrder, err := u.Repo.GetDexBTCListingOrderByID(listingOrderID)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, listingOrderID, err)
						break
					}
					if listingOrder != nil {
						if listingOrder.Cancelled || listingOrder.Matched {
							newStatus = entity.StatusDEXBuy_WaitingToRefund
							break
						}
						psbt, err := btc.ParsePSBTFromBase64(listingOrder.RawPSBT)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
							break
						}
						amountBTCFee += btc.EstimateTxFee(uint(len(listingOrder.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))
						amountBTC += listingOrder.Amount
						newStatus = entity.StatusDEXBuy_Buying

						buyReqInfos = append(buyReqInfos, btc.BuyReqInfo{
							Price:                      int(listingOrder.Amount),
							ReceiverInscriptionAddress: order.ReceiveAddress,
							SellerSignedPsbtB64:        listingOrder.RawPSBT,
						})

					} else {
						// ?? order not exist
						newStatus = entity.StatusDEXBuy_WaitingToRefund
						break
					}
				}

				if newStatus != order.Status {
					log.Println("watchPendingDexBTCBuyETH update buy multi-status", order.ID, newStatus, err)
					if newStatus == entity.StatusDEXBuy_Buying {
						address := u.Config.DexBTCWalletAddress

						walletInfo, err := btc.GetBalanceFromQuickNode(address, quickNodeAPI)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH GetBalanceFromQuickNode", order.ID, address, err)
							continue
						}
						utxos, err := btc.ConvertToUTXOType(walletInfo.Txrefs)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH ConvertToUTXOType", order.ID, err)
							continue
						}

						filteredUTXOs, err := u.filterUTXOInscription(utxos)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH filterUTXOInscription", order.ID, err)
							continue
						}
						filteredBalance := uint64(0)
						for _, v := range filteredUTXOs {
							filteredBalance += v.Value
						}

						if filteredBalance <= amountBTC+amountBTCFee {
							go u.NotifyWithChannel("C052CAWFB0D", "Multi Insufficient fund", address, fmt.Sprintf("filteredBalance %v <= amountBTCFee %v + listingOrder.Amount %v", filteredBalance, amountBTCFee, amountBTC))
							time.Sleep(300 * time.Millisecond)
							continue
						}

						// respondData, err := btc.CreatePSBTToBuyInscriptionMultiViaAPI(u.Config.DexBTCBuyService, address, psbtList, order.ReceiveAddress, amountBTC, filteredUTXOs, order.FeeRate, amountBTCFee)
						dataBytes, _ := json.Marshal(buyReqInfos)
						log.Printf("watchPendingDexBTCBuyETH sending multi--buy %v %v %v\n", order.ID.Hex(), order.ToJsonString(), string(dataBytes))
						respondData, err := btc.CreatePSBTToBuyInscriptionMultiViaAPI(u.Config.DexBTCBuyService, address, buyReqInfos, filteredUTXOs, order.FeeRate)
						if err != nil {
							go u.NotifyWithChannel("C052CAWFB0D", "Create multi buy ", address, fmt.Sprintf("filteredBalance %v <= amountBTCFee %v + listingOrder.Amount %v", filteredBalance, amountBTCFee, amountBTC))

							logData := make(map[string]interface{})
							logData["u.Config.DexBTCBuyService"] = u.Config.DexBTCBuyService
							logData["address"] = address
							logData["order.ReceiveAddress"] = order.ReceiveAddress
							logData["amountBTC"] = amountBTC
							logData["filteredUTXOs"] = filteredUTXOs
							logData["order.FeeRate"] = order.FeeRate
							logData["amountBTCFee"] = amountBTCFee
							logData["respondData"] = respondData
							logData["err"] = err.Error()

							u.Repo.CreateDexBTCLog(&entity.DexBTCLog{Function: "CreatePSBTToBuyInscriptionMultiViaAPI", Data: logData})
							log.Println("watchPendingDexBTCBuyETH CreatePSBTToBuyInscriptionMultiViaAPI", order.ID, err)
							time.Sleep(300 * time.Millisecond)
							continue
						}
						if respondData.SplitTxRaw != "" {
							err = btc.SendTxBlockStream(respondData.SplitTxRaw)
							if err != nil {
								dataBytes, _ := json.Marshal(respondData)
								log.Println("watchPendingDexBTCBuyETH SendTxBlockStream SplitTxHex", order.ID, string(dataBytes), err)
								continue
							}
							time.Sleep(5 * time.Second)
						}
						err = btc.SendTxBlockStream(respondData.TxHex)
						if err != nil {
							dataBytes, _ := json.Marshal(respondData)
							log.Println("watchPendingDexBTCBuyETH SendTxBlockStream TxHex", order.ID, string(dataBytes), err)
							continue
						}
						order.BuyTx = respondData.TxID
						order.SplitTx = respondData.SplitTxID
						order.Status = entity.StatusDEXBuy_Buying
						_, err = u.Repo.UpdateDexBTCBuyETHOrderBuy(&order)
						if err != nil {
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderBuy err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
						}
						continue
					} else {
						order.Status = newStatus
						_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
						if err != nil {
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
						}
						continue
					}
				}
			}
		case entity.StatusDEXBuy_Buying:
			// check tx buy if success => status = StatusDEXBuy_Bought else status = StatusDEXBuy_WaitingToRefund
			txStatus, err := btc.CheckTxfromQuickNode(order.BuyTx, quickNodeAPI)
			if err != nil {
				if strings.Contains(err.Error(), "tx not found") {
					if order.SplitTx != "" {
						txStatusSplit, err := btc.CheckTxfromQuickNode(order.SplitTx, quickNodeAPI)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH CheckTxfromQuickNode split", order.ID, order.SplitTx, err)
							//TODO: 2077 remove this in the future
							if len(order.SellOrderList) == 0 {
								listingOrder, err := u.Repo.GetDexBTCListingOrderByID(order.OrderID)
								if err != nil {
									log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, order.OrderID, err)
									continue
								}
								if listingOrder != nil {
									if listingOrder.Cancelled || listingOrder.Matched {
										order.Status = entity.StatusDEXBuy_WaitingToRefund
										_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
										if err != nil {
											log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
										}
										continue
									}
									//retry send buy tx
									order.Status = entity.StatusDEXBuy_ReceivedFund
									_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
									if err != nil {
										log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
									}
									continue
								}
							} else {
								newStatus := order.Status
								for _, listingOrderID := range order.SellOrderList {
									listingOrder, err := u.Repo.GetDexBTCListingOrderByID(listingOrderID)
									if err != nil {
										log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, listingOrderID, err)
										break
									}
									if listingOrder != nil {
										if listingOrder.Cancelled || listingOrder.Matched {
											newStatus = entity.StatusDEXBuy_WaitingToRefund
											break
										}
										//retry send buy tx
										newStatus = entity.StatusDEXBuy_ReceivedFund
									}
								}

								if newStatus != order.Status {
									order.Status = newStatus
									_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
									if err != nil {
										log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
									}
									continue
								}
							}
						}
						if txStatusSplit != nil {
							if txStatusSplit.Result.Confirmations >= 1 {
								//retry send buy tx
								order.Status = entity.StatusDEXBuy_ReceivedFund
								_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
								if err != nil {
									log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
								}
								continue
							}
							continue
						}
					}
				} else {
					log.Println("watchPendingDexBTCBuyETH CheckTxfromQuickNode", order.ID, order.BuyTx, err)
					continue
				}

			}
			if txStatus != nil {
				if txStatus.Result.Confirmations > 0 {
					order.Status = entity.StatusDEXBuy_Bought
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
					}
					continue
				}
			} else {
				//TODO: 2077 remove this in the future
				if len(order.SellOrderList) == 0 {
					listingOrder, err := u.Repo.GetDexBTCListingOrderByID(order.OrderID)
					if err != nil {
						log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, err)
						continue
					}
					if listingOrder != nil {
						if listingOrder.Cancelled || listingOrder.Matched {
							order.Status = entity.StatusDEXBuy_WaitingToRefund
							_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
							if err != nil {
								log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
							}
							continue
						} else {
							if time.Since(*order.CreatedAt) <= 2*time.Hour {
								//retry send buy tx
								order.Status = entity.StatusDEXBuy_ReceivedFund
								_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
								if err != nil {
									log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
								}
								continue
							} else {
								order.Status = entity.StatusDEXBuy_WaitingToRefund
								_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
								if err != nil {
									log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
								}
								continue
							}
						}
					} else {
						// ?? order not exist
						order.Status = entity.StatusDEXBuy_WaitingToRefund
						_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
						if err != nil {
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
						}
						continue
					}
				} else {
					newStatus := order.Status
					for _, listingOrderID := range order.SellOrderList {
						listingOrder, err := u.Repo.GetDexBTCListingOrderByID(listingOrderID)
						if err != nil {
							log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, listingOrderID, err)
							break
						}
						if listingOrder != nil {
							if listingOrder.Cancelled || listingOrder.Matched {
								newStatus = entity.StatusDEXBuy_WaitingToRefund
								break
							} else {
								if time.Since(*order.CreatedAt) <= 2*time.Hour {
									//retry send buy tx
									newStatus = entity.StatusDEXBuy_ReceivedFund
									continue
								} else {
									newStatus = entity.StatusDEXBuy_WaitingToRefund
									break
								}
							}
						}
					}

					if newStatus != order.Status {
						order.Status = newStatus
						_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
						if err != nil {
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
						}
						continue
					}
				}
			}

		case entity.StatusDEXBuy_WaitingToRefund:
			//send tx refund and update status to StatusDEXBuy_Refunding
			txID, _, err := u.EthClientDex.TransferMax(order.ETHKey, order.RefundAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransferMax", order.ID, order.RefundAddress, err)
				return err
			}

			order.RefundTx = txID
			order.Status = entity.StatusDEXBuy_Refunding
			order.SetUpdatedAt()
			_, err = u.Repo.UpdateDexBTCBuyETHOrderRefund(&order)
			if err != nil {
				log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderRefund err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
			}
			continue
		case entity.StatusDEXBuy_Refunding:
			// check tx refund if success => status = StatusDEXBuy_Refunded else status = StatusDEXBuy_WaitingToRefund
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			txhash := common.HexToHash(order.RefundTx)
			receipt, err := u.EthClientDex.GetClient().TransactionReceipt(ctx, txhash)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransactionReceipt", order.ID, order.RefundTx, err)
				continue
			}
			if receipt == nil {
				log.Println("watchPendingDexBTCBuyETH receipt is empty", order.ID, order.RefundTx, err)
				continue
			}
			if receipt.BlockNumber.Uint64()-currentBlockHeight < 15 {
				continue
			} else {
				order.Status = entity.StatusDEXBuy_Refunded
				_, err = u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
				}
			}
		case entity.StatusDEXBuy_Bought:
			// send eth to master and update status to StatusDEXBuy_SENDING_MASTER
			txID, _, err := u.EthClientDex.TransferMax(order.ETHKey, u.Config.DexBTCMasterETHAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransferMax", order.ID, u.Config.DexBTCMasterETHAddress, err)
				return err
			}

			order.MasterTx = txID
			order.Status = entity.StatusDEXBuy_SendingMaster
			order.SetUpdatedAt()
			_, err = u.Repo.UpdateDexBTCBuyETHOrderSendMaster(&order)
			if err != nil {
				log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderSendMaster err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
			}
			continue
		case entity.StatusDEXBuy_SendingMaster:
			// monitor tx and update status to StatusDEXBuy_SENT_MASTER
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			txhash := common.HexToHash(order.MasterTx)
			receipt, err := u.EthClientDex.GetClient().TransactionReceipt(ctx, txhash)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransactionReceipt", order.ID, order.MasterTx, err)
				continue
			}
			if receipt == nil {
				log.Println("watchPendingDexBTCBuyETH receipt is empty", order.ID, order.MasterTx, err)
				continue
			}
			if receipt.BlockNumber.Uint64()-currentBlockHeight < 15 {
				continue
			} else {
				order.Status = entity.StatusDEXBuy_SentMaster
				_, err = u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v %v %v\n", order.ID.Hex(), order.ToJsonString(), err)
				}
			}
		}
	}
	return nil
}

func (u Usecase) GenBuyETHOrder(isEstimate bool, userID string, orderID string, orderList []string, feeRate uint64, receiveAddress, refundAddress string) (string, string, string, int64, string, string, []string, bool, error) {
	var newOrder entity.DexBTCBuyWithETH
	var tempETHAddress string
	hasRoyalty := false
	var amountETH string
	var amountETHOriginal string
	var amountETHFee string
	btcRate, ethRate, _ := u.GetBTCToETHRate()
	if btcRate == 0 || ethRate == 0 {
		return "", "", "", 0, "", "", []string{}, false, errors.New("can't get exchange rate")
	}
	if len(orderList) == 0 {
		order, err := u.Repo.GetDexBTCListingOrderByID(orderID)
		if err != nil {
			return "", "", "", 0, "", "", []string{}, false, err
		}
		if order.Cancelled || order.Matched || order.InvalidMatch {
			return "", "", "", 0, "", "", []string{}, false, errors.New("order no longer valid")
		}

		if !order.Verified {
			return "", "", "", 0, "", "", []string{}, false, errors.New("order not ready yet")
		}
		psbt, err := btc.ParsePSBTFromBase64(order.RawPSBT)
		if err != nil {
			log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
			return "", "", "", 0, "", "", []string{}, false, err
		}

		amountBTCRequired := order.Amount + 1000
		amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
		amountBTCNoFee := amountBTCRequired
		amountBTCFee := btc.EstimateTxFee(uint(len(order.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))

		amountETHFee, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCFee)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}
		amountETHOriginal, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCNoFee)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}

		amountBTCRequired += amountBTCFee
		amountETH, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}

		var privKey, address string
		if !isEstimate {
			privKey, _, address, err = eth.NewClient(nil).GenerateAddress()
			if err != nil {
				logger.AtLog.Logger.Error("GenBuyETHOrder GenerateAddress", zap.Error(err))
				return "", "", "", 0, "", "", []string{}, false, err
			}
			tempETHAddress = address
		}

		tokenUri, err := u.GetTokenByTokenID(order.InscriptionID, 0)
		if err == nil {
			projectDetail, err := u.GetProjectDetail(structure.GetProjectDetailMessageReq{
				ContractAddress: tokenUri.ContractAddress,
				ProjectID:       tokenUri.ProjectID,
			})
			if err == nil {
				if projectDetail.Royalty > 0 {
					hasRoyalty = true
				}
			}

		}
		newOrder.UserID = userID
		newOrder.OrderID = orderID
		newOrder.AmountETH = amountETH
		newOrder.FeeRate = feeRate
		newOrder.Status = entity.StatusDEXBuy_Pending
		newOrder.ReceiveAddress = receiveAddress
		newOrder.RefundAddress = refundAddress
		newOrder.ETHKey = privKey
		newOrder.ETHAddress = address
		newOrder.ExpiredAt = time.Now().Add(2 * time.Hour)
		newOrder.InscriptionID = order.InscriptionID
		newOrder.AmountBTC = order.Amount

		fmt.Println("newOrder UUID: ", newOrder.UUID)

		expiredAt := time.Now().Add(1 * time.Hour).Unix()
		if !isEstimate {
			err := u.Repo.CreateDexBTCBuyWithETH(&newOrder)
			if err != nil {
				return "", "", "", expiredAt, "", "", []string{}, false, err
			}
		}

		return newOrder.UUID, tempETHAddress, amountETH, expiredAt, amountETHOriginal, amountETHFee, []string{}, hasRoyalty, nil

	} else {

		// get list order from db:
		listOrder, err := u.Repo.GetDexBTListingOrderByListID(orderList)
		if err != nil {
			return "", "", "", 0, "", "", []string{}, false, err
		}

		var inscriptionList []string
		var orderListFinal []string
		var orderListInvalid []string

		var amountBtcSum, feeAmountBtcSum, totalAmountBtcSum uint64

		for _, order := range listOrder {

			if order.Cancelled || order.Matched || order.InvalidMatch {
				orderListInvalid = append(orderListInvalid, order.UUID)
				continue
			}

			if !order.Verified {
				orderListInvalid = append(orderListInvalid, order.UUID)
				continue
			}
			psbt, err := btc.ParsePSBTFromBase64(order.RawPSBT)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
				orderListInvalid = append(orderListInvalid, order.UUID)
				continue
			}

			// amount btc with no fee:
			amountBtc := order.Amount + 1000
			amountBtc += (amountBtc / 10000) * 15 // + 0,15%

			// btc fee:
			feeAmountBtc := btc.EstimateTxFee(uint(len(order.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))

			// total btc amount:
			totalAmountBtc := amountBtc + feeAmountBtc

			// sum:
			amountBtcSum += amountBtc
			feeAmountBtcSum += feeAmountBtc
			totalAmountBtcSum += totalAmountBtc

			// add list:
			inscriptionList = append(inscriptionList, order.InscriptionID)
			orderListFinal = append(orderListFinal, order.UUID)
		}
		// check list:
		if len(orderListFinal) == 0 {
			return "", "", "", 0, "", "", []string{}, false, err
		}
		// convert to eth:

		// amount eth:
		amountETHOriginal, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBtcSum)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}

		// fee eth
		amountETHFee, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(feeAmountBtcSum)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}

		// total amount eth:
		amountETH, _, _, err = u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(totalAmountBtcSum)/1e8), btcRate, ethRate)
		if err != nil {
			logger.AtLog.Logger.Error("GenBuyETHOrder convertBTCToETH", zap.Error(err))
			return "", "", "", 0, "", "", []string{}, false, err
		}

		var privKey, address string
		if !isEstimate {
			privKey, _, address, err = eth.NewClient(nil).GenerateAddress()
			if err != nil {
				logger.AtLog.Logger.Error("GenBuyETHOrder GenerateAddress", zap.Error(err))
				return "", "", "", 0, "", "", []string{}, false, err
			}
			tempETHAddress = address
		}

		newOrder.UserID = userID
		newOrder.OrderID = orderID

		newOrder.AmountETH = amountETH

		newOrder.FeeRate = feeRate

		newOrder.Status = entity.StatusDEXBuy_Pending
		newOrder.ReceiveAddress = receiveAddress
		newOrder.RefundAddress = refundAddress
		newOrder.ETHKey = privKey
		newOrder.ETHAddress = address
		newOrder.ExpiredAt = time.Now().Add(2 * time.Hour)
		newOrder.InscriptionList = inscriptionList
		newOrder.SellOrderList = orderListFinal

		newOrder.AmountBTC = totalAmountBtcSum

		expiredAt := time.Now().Add(1 * time.Hour).Unix()
		if !isEstimate {

			// check valid:
			if len(orderListInvalid) > 0 {
				return "", "", "", expiredAt, "", "", orderListInvalid, false, errors.New("orderList invalid")
			}

			err := u.Repo.CreateDexBTCBuyWithETH(&newOrder)
			if err != nil {
				return "", "", "", expiredAt, "", "", orderListInvalid, false, err
			}
		}
		return newOrder.UUID, tempETHAddress, amountETH, expiredAt, amountETHOriginal, amountETHFee, orderListInvalid, hasRoyalty, nil

	}
}

// func (u Usecase) UpdateBuyETHOrderTx(buyOrderID string, userID string, txhash string) error {
// 	order, err := u.Repo.GetDexBTCBuyETHOrderByID(buyOrderID)
// 	if err != nil {
// 		return err
// 	}
// 	order.ETHTx = txhash
// 	_, err = u.Repo.UpdateDexBTCBuyETHOrderTx(order)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (u Usecase) filterUTXOInscription(utxos []btc.UTXOType) ([]btc.UTXOType, error) {
	var result []btc.UTXOType
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	waitChan := make(chan struct{}, 10)
	for _, output := range utxos {
		wg.Add(1)
		time.Sleep(10 * time.Millisecond)
		waitChan <- struct{}{}
		go func(op btc.UTXOType) {
			defer func() {
				wg.Done()
				<-waitChan
			}()
			outStr := fmt.Sprintf("%v:%v", op.TxHash, op.TxOutIndex)
			inscriptions, err := getInscriptionByOutput(ordServer, outStr)
			if err != nil {
				return
			}
			if len(inscriptions.Inscriptions) == 0 {
				lock.Lock()
				result = append(result, op)
				lock.Unlock()
			}
		}(output)
	}
	wg.Wait()
	return result, nil
}
