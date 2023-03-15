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

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/eth"
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
				spentTxDetail, err := btc.CheckTxfromQuickNode(spentTx, u.Config.QuicknodeAPI)
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC(spentTx) %v %v\n", order.Inputs, err)
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
					if totalOutvalue >= order.Amount && matchIns == len(order.Inputs) && matchOuts == len(psbtData.UnsignedTx.TxOut) {
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
					// Discord Notify NEW SALE
					buyerAddress := order.Buyer
					go u.NotifyNewSale(order, buyerAddress)
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
	pendingOrders, err := u.Repo.GetDexBTCBuyETHOrderByStatus([]entity.DexBTCETHBuyStatus{entity.StatusDEXBuy_Pending, entity.StatusDEXBuy_ReceivedFund, entity.StatusDEXBuy_Buying, entity.StatusDEXBuy_WaitingToRefund, entity.StatusDEXBuy_Refunding, entity.StatusDEXBuy_Bought, entity.StatusDEXBuy_SendingMaster})
	if err != nil {
		return err
	}

	quickNodeAPI := u.Config.QuicknodeAPI
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		log.Printf("watchPendingDexBTCBuyETH ethclientDial err %v\n", err)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)
	currentBlockHeight, err := ethClient.GetClient().BlockNumber(context.Background())
	if err != nil {
		log.Printf("watchPendingDexBTCBuyETH BlockNumber err %v\n", err)
		return err
	}
	for _, order := range pendingOrders {
		ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
		if err != nil {
			log.Printf("watchPendingDexBTCBuyETH ethclientDial err %v\n", err)
			continue
		}
		ethClient := eth.NewClient(ethClientWrap)

		switch order.Status {
		case entity.StatusDEXBuy_Pending:
			// check wallet receive enough funds
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, ethAddress, err := ethClient.GenerateAddressFromPrivKey(order.ETHKey)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GenerateAddressFromPrivKey", order.ID, err)
				continue
			}
			amountRequired, ok := new(big.Int).SetString(order.AmountETH, 10)
			if !ok {
				log.Println("watchPendingDexBTCBuyETH new(bigInt)", order.ID, order.AmountETH, err)
				continue
			}
			balance, err := ethClient.GetBalance(ctx, ethAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GetBalanceFromQuickNode", order.ID, ethAddress, err)
				continue
			}
			if balance.Cmp(amountRequired) > -1 {
				if order.Confirmation > 1 && time.Since(*order.CreatedAt) > 3*time.Minute {
					order.Status = entity.StatusDEXBuy_ReceivedFund
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
					}
					continue
				} else {
					order.Confirmation += 1
					_, err := u.Repo.UpdateDexBTCBuyETHOrderConfirmation(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
					}
					continue
				}
			} else {
				// not enough funds
				if time.Since(order.ExpiredAt) > 1*time.Second {
					order.Status = entity.StatusDEXBuy_Expired
					_, err := u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
					if err != nil {
						log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
					}
					continue
				}
			}
		case entity.StatusDEXBuy_ReceivedFund:
			// send tx buy update status to StatusDEXBuy_Buying
			listingOrder, err := u.Repo.GetDexBTCListingOrderByID(order.OrderID)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH GetDexBTCListingOrderByID", order.ID, err)
				continue
			}
			if listingOrder != nil {

				// privKeyDecrypt, err := encrypt.DecryptToString(u.Config.DexBTCKey, os.Getenv("SECRET_KEY"))
				// if err != nil {
				// 	log.Println("watchPendingDexBTCBuyETH DecryptToString", order.ID, err)
				// 	continue
				// }
				// _, _, address, err := btc.GenerateAddressSegwit(privKeyDecrypt)
				// if err != nil {
				// 	log.Println("watchPendingDexBTCBuyETH GenerateAddressSegwit", order.ID, err)
				// 	continue
				// }
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
				psbt, err := btc.ParsePSBTFromBase64(listingOrder.RawPSBT)
				if err != nil {
					log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
					continue
				}

				feeRate := order.FeeRate
				amountBTCFee := uint64(0)
				amountBTCFee = btc.EstimateTxFee(uint(len(listingOrder.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))

				respondData, err := btc.CreatePSBTToBuyInscriptionViaAPI(u.Config.DexBTCBuyService, address, listingOrder.RawPSBT, order.ReceiveAddress, listingOrder.Amount, utxos, 15, amountBTCFee)
				if err != nil {
					log.Println("watchPendingDexBTCBuyETH CreatePSBTToBuyInscription", order.ID, err)
					continue
				}
				if respondData.SplitTxRaw != "" {
					_, err = btc.SendRawTxfromQuickNode(respondData.SplitTxRaw, quickNodeAPI)
					if err != nil {
						dataBytes, _ := json.Marshal(respondData)
						log.Println("watchPendingDexBTCBuyETH SendRawTxfromQuickNode SplitTxHex", order.ID, string(dataBytes), err)
						continue
					}
				}
				_, err = btc.SendRawTxfromQuickNode(respondData.TxHex, quickNodeAPI)
				if err != nil {
					dataBytes, _ := json.Marshal(respondData)
					log.Println("watchPendingDexBTCBuyETH SendRawTxfromQuickNode TxHex", order.ID, string(dataBytes), err)
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
			txID, _, err := ethClient.TransferMax(order.ETHKey, order.RefundAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransferMax", order.ID, order.RefundAddress, err)
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			txhash := common.HexToHash(order.RefundTx)
			receipt, err := ethClient.GetClient().TransactionReceipt(ctx, txhash)
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
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
				}
			}
			// order.Status = entity.StatusDEXBuy_WaitingToRefund
			// _, err = u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
			// if err != nil {
			// 	log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
			// }
		case entity.StatusDEXBuy_Bought:
			// send eth to master and update status to StatusDEXBuy_SENDING_MASTER
			txID, _, err := ethClient.TransferMax(order.ETHKey, u.Config.DexBTCMasterETHAddress)
			if err != nil {
				log.Println("watchPendingDexBTCBuyETH TransferMax", order.ID, order.RefundAddress, err)
				return err
			}

			order.MasterTx = txID
			order.Status = entity.StatusDEXBuy_SendingMaster
			order.SetUpdatedAt()
			_, err = u.Repo.UpdateDexBTCBuyETHOrder(&order)
			if err != nil {
				log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
			}
			continue
		case entity.StatusDEXBuy_SendingMaster:
			// monitor tx and update status to StatusDEXBuy_SENT_MASTER
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			txhash := common.HexToHash(order.MasterTx)
			receipt, err := ethClient.GetClient().TransactionReceipt(ctx, txhash)
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
				order.Status = entity.StatusDEXBuy_SentMaster
				_, err = u.Repo.UpdateDexBTCBuyETHOrderStatus(&order)
				if err != nil {
					log.Printf("watchPendingDexBTCBuyETH UpdateDexBTCBuyETHOrderStatus err %v\n", err)
				}
			}
		}
	}
	return nil
}

func (u Usecase) GenBuyETHOrder(userID string, orderID string, feeRate uint64, receiveAddress, refundAddress string) (string, string, string, int64, string, string, bool, error) {
	order, err := u.Repo.GetDexBTCListingOrderByID(orderID)
	if err != nil {
		return "", "", "", 0, "", "", false, err
	}
	if order.Cancelled || order.Matched || order.InvalidMatch {
		return "", "", "", 0, "", "", false, errors.New("order no longer valid")
	}

	if !order.Verified {
		return "", "", "", 0, "", "", false, errors.New("order not ready yet")
	}
	psbt, err := btc.ParsePSBTFromBase64(order.RawPSBT)
	if err != nil {
		log.Println("watchPendingDexBTCBuyETH ParsePSBTFromBase64", order.ID, err)
		return "", "", "", 0, "", "", false, err
	}

	amountBTCRequired := order.Amount + 1000
	amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%
	amountBTCNoFee := amountBTCRequired
	amountBTCFee := btc.EstimateTxFee(uint(len(order.Inputs)+3), uint(len(psbt.UnsignedTx.TxOut)+2), uint(feeRate)) + btc.EstimateTxFee(1, 2, uint(feeRate))

	btcRate, ethRate, _ := u.GetBTCToETHRate()

	amountETHFee, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCFee)/1e8), btcRate, ethRate)
	if err != nil {
		u.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
		return "", "", "", 0, "", "", false, err
	}
	amountETHOriginal, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCNoFee)/1e8), btcRate, ethRate)
	if err != nil {
		u.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
		return "", "", "", 0, "", "", false, err
	}

	amountBTCRequired += amountBTCFee
	amountETH, _, _, err := u.convertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
	if err != nil {
		u.Logger.Error("GenBuyETHOrder convertBTCToETH", err.Error(), err)
		return "", "", "", 0, "", "", false, err
	}

	ethClient := eth.NewClient(nil)
	privKey, _, tempETHAddress, err := ethClient.GenerateAddress()
	if err != nil {
		u.Logger.Error("GenBuyETHOrder GenerateAddress", err.Error(), err)
		return "", "", "", 0, "", "", false, err
	}
	hasRoyalty := false
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

	var newOrder entity.DexBTCBuyWithETH
	newOrder.UserID = userID
	newOrder.OrderID = orderID
	newOrder.AmountETH = amountETH
	newOrder.FeeRate = feeRate
	newOrder.Status = entity.StatusDEXBuy_Pending
	newOrder.ReceiveAddress = receiveAddress
	newOrder.RefundAddress = refundAddress
	newOrder.ETHKey = privKey
	newOrder.ExpiredAt = time.Now().Add(2 * time.Hour)
	newOrder.InscriptionID = order.InscriptionID
	newOrder.AmountBTC = order.Amount
	expiredAt := newOrder.ExpiredAt.Unix()
	err = u.Repo.CreateDexBTCBuyWithETH(&newOrder)
	if err != nil {
		return "", "", "", expiredAt, "", "", false, err
	}

	return newOrder.UUID, tempETHAddress, amountETH, expiredAt, amountETHOriginal, amountETHFee, hasRoyalty, nil
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
