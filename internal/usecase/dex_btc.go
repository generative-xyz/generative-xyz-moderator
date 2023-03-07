package usecase

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
			if creator.WalletAddressBTC != "" && creator.WalletAddressBTCTaproot != "" {
				royaltyFeePercent = float64(projectDetail.Royalty) / 10000
				if creator.WalletAddressBTCTaproot != "" {
					artistAddress = creator.WalletAddressBTCTaproot
				} else {
					artistAddress = creator.WalletAddressBTC
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

func (u Usecase) JobWatchPendingDexBTCListing() error {
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
					log.Printf("JobWatchPendingDexBTCListing btc.CheckTxFromBTC(spentTx) u.Config.QuicknodeAPI %v %v %v\n", u.Config.QuicknodeAPI, order.Inputs, err)
				}
				_ = txDetail
				// output := txDetail.Result.Vout[0]
				// order.Buyer = output.ScriptPubKey.Address

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
		if !order.Verified {
			txDetail, err := btc.CheckTxfromQuickNode(inscriptionTx[0], u.Config.QuicknodeAPI)
			if err != nil {
				fmt.Errorf("btc.GetBTCTxStatusExtensive %v\n", err)
			} else {
				if txDetail.Result.Confirmations > 0 {
					order.Verified = true
					_, err = u.Repo.UpdateDexBTCListingOrderMatchTx(&order)
					if err != nil {
						log.Printf("JobWatchPendingDexBTCListing UpdateDexBTCListingOrderMatchTx err %v\n", err)
						continue
					}
				}
			}
		}
	}
	return nil
}
