package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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
		orderInfo.CancelTx = txhash
	} else {
		return errors.New("order already cancelling/cancelled")
	}
	_, err = u.Repo.UpdateDexBTCListingOrder(orderInfo)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) DexBTCListing(seller_address string, raw_psbt string, inscription_id string) error {
	newListing := entity.DexBTCListing{
		RawPSBT:       raw_psbt,
		InscriptionID: inscription_id,
		SellerAddress: seller_address,
		Cancelled:     false,
		CancelTx:      "",
	}
	var psbtTx structure.PSBTData

	err := json.Unmarshal([]byte(raw_psbt), &psbtTx)
	if err != nil {
		return err
	}

	psbtData, err := btc.ParsePSBTFromBase64(raw_psbt)
	if err != nil {
		return err
	}

	outputList, err := extractAllOutputFromPSBT(psbtData)
	if err != nil {
		return err
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
			if creator.WalletAddressBTCTaproot != "" {
				royaltyFeePercent = float64(projectDetail.Royalty) / 10000
				// royaltyFee = int(float64(totalAmount.Int64()) * royaltyFeePercent)
				artistAddress = creator.WalletAddressBTC
			}
		}
	}

	if artistAddress != "" && royaltyFeePercent > 0 {
		if len(psbtData.UnsignedTx.TxOut) == 1 {
			//force receiver == artistAddress when only one output
			for receiver, _ := range outputList {
				if receiver != artistAddress {
					return fmt.Errorf("expected to paid royalty fee to %v", artistAddress)
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
					if totalValue >= royaltyFeeExpected {
						return fmt.Errorf("expected royalty fee of artist %v to be %v, got %v", artistAddress, royaltyFeeExpected, totalValue)
					}
				}
			}
		}
	}

	// previousTxs, err := retrievePreviousTxFromPSBT(psbtData)
	// if err != nil {
	// 	return err
	// }

	// _, bs, err := u.buildBTCClient()
	// if err != nil {
	// 	fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
	// 	return err
	// }

	//TODO: check previous tx
	// for tx, _ := range previousTxs {
	// 	status, err := btc.GetBTCTxStatusExtensive(tx, bs)
	// 	if err != nil {
	// 		fmt.Errorf("btc.GetBTCTxStatusExtensive %v\n", err)
	// 	}
	// 	switch status {
	// 	case "Failed":

	// 	case "Success":
	// 		newListing.Verified = true
	// 	case "Pending":

	// 	}
	// }

	return u.Repo.CreateDexBTCListing(&newListing)
}

func retrievePreviousTxFromPSBT(psbtData *psbt.Packet) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	for _, input := range psbtData.UnsignedTx.TxIn {
		result[input.PreviousOutPoint.Hash.String()] = struct{}{}
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

	return nil
}
