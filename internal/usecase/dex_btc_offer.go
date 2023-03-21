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

	"github.com/btcsuite/btcd/wire"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/btc"
)

func (u Usecase) JobWatchPendingDexBTCOffer() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingDexBTCOffering()
		if err != nil {
			log.Println("JobWatchPendingDexBTCListing watchPendingDexBTCListing err", err)
		}
	}(&wg)

	wg.Wait()
}

func (u Usecase) CancelDexBTCOffer(txhash string, offerer_address string, inscription_id string, order_id string) error {
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

func (u Usecase) DexBTCOffering(offerer_address string, raw_psbt string, inscription_id string, split_tx string, inscription_owner string) (*entity.DexBTCOffer, error) {
	newOffer := entity.DexBTCOffer{
		RawPSBT:        raw_psbt,
		InscriptionID:  inscription_id,
		OffererAddress: offerer_address,
		Seller:         inscription_owner,
		Cancelled:      false,
		CancelTx:       "",
		SplitTx:        split_tx,
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
		newOffer.Amount = totalOuputValue
	} else {
		newOffer.Amount = totalOuputValue - uint64(psbtData.Inputs[1].WitnessUtxo.Value)
	}

	txInputs := []string{}
	for _, input := range psbtData.UnsignedTx.TxIn {
		i := fmt.Sprintf("%v:%v", input.PreviousOutPoint.Hash.String(), input.PreviousOutPoint.Index)
		txInputs = append(txInputs, i)
	}
	newOffer.Inputs = txInputs

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

	previousVins, err := retrieveAllVinFromPSBT(psbtData)
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

	//require found current inscription outpoint
	// inscriptionTx := strings.Split(inscriptionInfo.Satpoint, ":")[0]
	found := false
	for _, tx := range previousVins {
		if inscriptionInfo.Satpoint == tx {
			found = true
			break
		}
	}
	if !found {
		return nil, errors.New("can't found inscription in offer request")
	}

	if split_tx != "" {
		_, err = btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI)
		if err != nil {
			fmt.Printf("btc.SendRawTxfromQuickNode(split_tx, u.Config.QuicknodeAPI) - with err: %v %v\n", err, split_tx)
			return nil, err
		}
		newOffer.Verified = false
	}

	return &newOffer, u.Repo.CreateDexBTCOffer(&newOffer)
}

func (u Usecase) watchPendingDexBTCOffering() error {
	pendingOrders, err := u.Repo.GetDexBTCListingOrderPending()
	if err != nil {
		return err
	}
	_, bs, err := u.buildBTCClientCustomToken(u.Config.DEXBTCBlockcypherToken)
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

			log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC %v\n", inscriptionTx[0])
			txStatus, err := bs.CheckTx(inscriptionTx[0])
			if err != nil {
				log.Printf("JobWatchPendingDexBTCListing bs.CheckTx(txhash) %v\n", order.Inputs)
				spentTx, err = btc.CheckOutcoinSpentBlockStream(inscriptionTx[0], uint(idx))
				if err != nil {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream %v\n", order.Inputs)
					continue
				}
				if spentTx != "" {
					log.Printf("JobWatchPendingDexBTCListing btc.CheckOutcoinSpentBlockStream success\n")
				}
			} else {
				if txStatus.Outputs[idx].SpentBy != "" {
					spentTx = txStatus.Outputs[idx].SpentBy
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
					// txDetail, err := btc.CheckTxfromQuickNode(spentTx, u.Config.QuicknodeAPI)
					// if err != nil {
					// 	log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC(spentTx) %v %v\n", order.Inputs, err)
					// }
					// output := txDetail.Result.Vout[0]
					// order.Seller = output.ScriptPubKey.Address

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

					// Discord Notify NEW SALE or NEW OFFER ACCEPT @dac
					// order.OffererAddress
					// order.Seller
					// buyerAddress := order.Buyer
					// go u.NotifyNewSale(order, buyerAddress)
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
