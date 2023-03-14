package usecase

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/btc"
)

func (u Usecase) CancelDexBTCListing(txhash string, seller_address string, inscription_id string, order_id string) error {
	orderInfo, err := u.Repo.GetDexBTCListingOrderByID(order_id)
	if err != nil {
		return err
	}
	if orderInfo.InscriptionID != inscription_id {
		return errors.New("invalid cancelling request")
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
	}
	if split_tx != "" {
		_, err = btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI)
		if err != nil {
			fmt.Printf("btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI) - with err: %v %v\n", err, split_tx)
			return nil, err
		}

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
func (u Usecase) JobWatchPendingDexBTCListing() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingDexBTCListing()
		if err != nil {
			log.Println("JobWatchPendingDexBTCListing watchPendingDexBTCListing err", err)
		}
	}(&wg)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingDexBTCBuyETH()
		if err != nil {
			log.Println("JobWatchPendingDexBTCListing watchPendingDexBTCListing err", err)
		}
	}(&wg)

	wg.Wait()

}

func (u Usecase) watchPendingDexBTCListing() error {
	pendingOrders, err := u.Repo.GetDexBTCListingOrderPending()
	if err != nil {
		return err
	}
	_, bs, err := u.buildBTCClient()
	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}
	for _, order := range pendingOrders {
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
			txDetail, err := btc.CheckTxFromBTC(inscriptionTx[0])
			if err == nil {
				if txDetail.Data.Outputs != nil {
					outputs := *txDetail.Data.Outputs
					if outputs[idx].SpentByTx != "" {
						spentTx = outputs[idx].SpentByTx
					}
				}
			} else {
				log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC %v\n", inscriptionTx[0])
				txStatus, err := bs.CheckTx(inscriptionTx[0])
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing bs.CheckTx(txhash) %v\n", order.Inputs)
					continue
				} else {
					if txStatus.Outputs[idx].SpentBy != "" {
						spentTx = txStatus.Outputs[idx].SpentBy
					}
				}
			}

			if spentTx != "" {
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
				// Discord Notify NEW SALE
				buyerAddress := order.Buyer
				go u.NotifyNewSale(order, buyerAddress)
			}
		} else {
			status, err := btc.GetBTCTxStatusExtensive(order.CancelTx, bs, u.Config.QuicknodeAPI)
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

func (u Usecase) watchPendingDexBTCBuyETH() error {
	pendingOrders, err := u.Repo.GetDexBTCBuyETHOrderByStatus([]entity.DexBTCETHBuyStatus{entity.StatusDEXBuy_Pending, entity.StatusDEXBuy_ReceivedFund, entity.StatusDEXBuy_Buying, entity.StatusDEXBuy_WaitingToRefund, entity.StatusDEXBuy_Refunding})
	if err != nil {
		return err
	}

	quickNodeAPI := u.Config.QuicknodeAPI

	for _, order := range pendingOrders {
		_, _, address, err := btc.GenerateAddressSegwit(order.TempBTCKey)
		if err != nil {
			log.Println("watchPendingDexBTCBuyETH GenerateAddressSegwit", err)
			continue
		}
		switch order.Status {
		case entity.StatusDEXBuy_Pending:
			// check wallet receive enough funds
			walletInfo, err := btc.GetBalanceFromQuickNode(address, quickNodeAPI)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GetBalanceFromQuickNode", order.ID, address, err)
				continue
			}
			if uint64(walletInfo.Balance) >= order.AmountBTC {
				order.Status = entity.StatusDEXBuy_ReceivedFund
				_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
				}
				continue
			} else {
				// not enough funds
			}
		case entity.StatusDEXBuy_ReceivedFund:
			// send tx buy update status to StatusDEXBuy_Buying
			listingOrder, err := u.Repo.GetDexBTCListingOrderByID(order.OrderID)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, err)
				continue
			}
			if listingOrder != nil {
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
				rawtx, _, err := btc.CreatePSBTToBuyInscription(listingOrder.RawPSBT, order.TempBTCKey, address, order.ReceiveAddress, listingOrder.Amount, utxos, 15)
				if err != nil {
					log.Println("watchPendingDexBTCBuyETH CreatePSBTToBuyInscription", order.ID, err)
					continue
				}

				_, err = btc.SendRawTxfromQuickNode(rawtx, quickNodeAPI)
				if err != nil {
					log.Println("watchPendingDexBTCBuyETH CreatePSBTToBuyInscription", order.ID, err)
					continue
				}
			} else {
				// ?? order not exist
			}
		case entity.StatusDEXBuy_Buying:
			// check tx buy if success => status = StatusDEXBuy_Bought else status = StatusDEXBuy_WaitingToRefund
			txStatus, err := btc.CheckTxfromQuickNode(order.BuyTx, quickNodeAPI)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH CheckTxfromQuickNode", order.ID, order.BuyTx, err)
				continue
			}
			if txStatus != nil {
				if txStatus.Result.Confirmations > 0 {
					order.Status = entity.StatusDEXBuy_Bought
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
					}
					continue
				}
			} else {
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
							log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
						}
						continue
					}
					if time.Since(*order.CreatedAt) >= 1*time.Hour {

					}
				} else {
					// ?? order not exist
				}
			}

		case entity.StatusDEXBuy_WaitingToRefund:
			//send tx refund and update status to StatusDEXBuy_Refunding
			_, bs, err := u.buildBTCClient()
			if err != nil {
				fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
				return err
			}
			// user, err := u.Repo.FindUserByID(order.UserID)
			// if err != nil {
			// 	log.Println("watchPendingDexBTCBuyETH FindUserByID", order.ID, order.UserID, err)
			// 	return err
			// }
			txID, err := bs.SendTransactionWithPreferenceFromSegwitAddress(
				order.TempBTCKey,
				address,
				order.ReceiveAddress,
				int(order.AmountBTC),
				btc.PreferenceMedium,
			)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH.SendTransactionWithPreferenceFromSegwitAddress", order.ID, order.UserID, err)
				return err
			}

			order.RefundTx = txID
			order.Status = entity.StatusDEXBuy_Refunding
			order.SetUpdatedAt()
			_, err = u.Repo.UpdateDexBTCBuyETHOrder(&order)
			if err != nil {
				log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
			}
			continue
		case entity.StatusDEXBuy_Refunding:
			// check tx refund if success => status = StatusDEXBuy_Refunded else status = StatusDEXBuy_WaitingToRefund
			txStatus, err := btc.CheckTxfromQuickNode(order.RefundTx, quickNodeAPI)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH CheckTxfromQuickNode", order.ID, order.RefundTx, err)
				continue
			}
			if txStatus != nil {
				if txStatus.Result.Confirmations > 0 {
					order.Status = entity.StatusDEXBuy_Refunded
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
					}
					continue
				}
			} else {

				order.Status = entity.StatusDEXBuy_WaitingToRefund
				_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
				}
			}
		}
	}
	return nil
}

func (u Usecase) GenBuyETHOrder(userID string, orderID string, amount uint64, feeRate uint64, receiveAddress string) (string, string, error) {
	order, err := u.Repo.GetDexBTCListingOrderByID(orderID)
	if err != nil {
		return "", "", err
	}
	expectedAmount := order.Amount + 1000

	switch len(order.Inputs) {
	case 1:
		expectedAmount += btc.EstimateTxFee(3, 2, uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))
	case 2:
		expectedAmount += btc.EstimateTxFee(4, 3, uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))
	case 3:
		expectedAmount += btc.EstimateTxFee(5, 4, uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))
	}

	if amount < expectedAmount {
		return "", "", fmt.Errorf("expected receieve amount to be %d, got %d", expectedAmount, amount)
	}

	var newOrder entity.DexBTCBuyWithETH
	var tempAddress string
	privKey, _, address, err := btc.GenerateAddressSegwit()
	if err != nil {
		return "", "", err
	}
	tempAddress = address
	newOrder.TempBTCKey = privKey
	newOrder.OrderID = orderID
	newOrder.AmountBTC = order.Amount
	newOrder.FeeRate = feeRate
	newOrder.Status = entity.StatusDEXBuy_Pending
	newOrder.ReceiveAddress = receiveAddress

	err = u.Repo.CreateDexBTCBuyWithETH(&newOrder)
	if err != nil {
		return "", "", err
	}

	return newOrder.UUID, tempAddress, nil
}
func (u Usecase) UpdateBuyETHOrderTx(buyOrderID string, userID string, txhash string) error {
	order, err := u.Repo.GetDexBTCBuyETHOrderByID(buyOrderID)
	if err != nil {
		return err
	}
	order.ETHTx = txhash
	_, err = u.Repo.UpdateDexBTCBuyETHOrderTx(order)
	if err != nil {
		return err
	}
	return nil
}
