package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/logger"
)

func (u Usecase) GetInscriptionByIDFromOrd(id string) (*structure.InscriptionOrdInfoByID, error) {
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev.generativeexplorer.com"
	}
	return getInscriptionByID(ordServer, id)
}

func (u Usecase) GetBTCWalletInfo(address string) (*structure.WalletInfo, error) {
	cacheKey := utils.KEY_BTC_WALLET_INFO + "_" + address
	var result structure.WalletInfo
	exist, err := u.Repo.Cache.Exists(cacheKey)
	if err == nil && *exist {
		data, err := u.Repo.Cache.GetData(cacheKey)
		if err == nil && data != nil {
			err := json.Unmarshal([]byte(*data), &result)
			if err != nil {
				u.Logger.Error("GetBTCWalletInfo json.Unmarshal", address, err)
			}
			return &result, nil
		}
	}

	apiToken := u.Config.BlockcypherToken
	u.Logger.Info("GetBTCWalletInfo apiToken debug", apiToken)
	quickNode := u.Config.QuicknodeAPI

	walletBasicInfo, err := btc.GetBalanceFromQuickNode(address, quickNode)
	if err != nil {
		walletBasicInfo, err = getWalletInfo(address, apiToken, u.Logger)
		if err != nil {
			u.Logger.Info("GetBTCWalletInfo apiToken debug err", err)
			return nil, err
		}
	}

	result.BlockCypherWalletInfo = *walletBasicInfo
	outcoins := []string{}
	for _, outcoin := range result.Txrefs {
		o := fmt.Sprintf("%s:%v", outcoin.TxHash, outcoin.TxOutputN)
		outcoins = append(outcoins, o)
	}
	currentListing, err := u.Repo.GetDexBTCListingOrderUserPending(address)
	if err != nil {
		u.Logger.Error("u.Repo.GetDexBTCListingOrderUserPending", address, err)
	}

	inscriptions, outputInscMap, outputSatRanges, err := u.InscriptionsByOutputs(outcoins, currentListing)
	if err != nil {
		return nil, err
	}
	result.InscriptionsByOutputs = outputInscMap
	for _, item := range inscriptions {
		result.Inscriptions = append(result.Inscriptions, item...)
	}
	newTxrefs := []structure.TxRef{}
	for _, outcoin := range result.Txrefs {
		o := fmt.Sprintf("%s:%v", outcoin.TxHash, outcoin.TxOutputN)
		satRanges, ok := outputSatRanges[o]
		if ok {
			outcoin.SatRanges = satRanges
			newTxrefs = append(newTxrefs, outcoin)
		}
	}
	result.Txrefs = newTxrefs

	newTxrefsFiltered := []structure.TxRef{}
	if len(currentListing) > 0 {
		pendingUTXO := make(map[string]struct{})
		for _, listing := range currentListing {
			for _, input := range listing.Inputs {
				// only filter out non-inscription utxo
				if _, ok := result.InscriptionsByOutputs[input]; !ok {
					pendingUTXO[input] = struct{}{}
				}
			}
		}

		for _, output := range result.Txrefs {
			voutStr := fmt.Sprintf("%v:%v", output.TxHash, output.TxOutputN)
			if _, ok := pendingUTXO[voutStr]; !ok {
				newTxrefsFiltered = append(newTxrefs, output)
			}
		}
		result.Txrefs = newTxrefsFiltered
	}

	err = u.Repo.Cache.SetDataWithExpireTime(cacheKey, result, 10)
	if err != nil {
		u.Logger.Error("GetBTCWalletInfo CreateCache", address, err)
	}

	return &result, nil
}

func (u Usecase) InscriptionsByOutputs(outputs []string, currentListing []entity.DexBTCListing) (map[string][]structure.WalletInscriptionInfo, map[string][]structure.WalletInscriptionByOutput, map[string][][]uint64, error) {
	result := make(map[string][]structure.WalletInscriptionInfo)
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev.generativeexplorer.com"
	}
	outputSatRanges := make(map[string][][]uint64)
	outputInscMap := make(map[string][]structure.WalletInscriptionByOutput)
	for _, output := range outputs {
		if _, ok := result[output]; ok {
			continue
		}
		inscriptions, err := getInscriptionByOutput(ordServer, output)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(inscriptions.Inscriptions) > 0 {
			for _, insc := range inscriptions.Inscriptions {
				data, err := getInscriptionByID(ordServer, insc)
				if err != nil {
					return nil, nil, nil, err
				}
				offset, err := strconv.ParseInt(strings.Split(data.Satpoint, ":")[2], 10, 64)
				if err != nil {
					return nil, nil, nil, err
				}
				inscWalletInfo := structure.WalletInscriptionInfo{
					InscriptionID: data.InscriptionID,
					Number:        data.Number,
					ContentType:   data.ContentType,
					Offset:        offset,
				}
				inscWalletByOutput := structure.WalletInscriptionByOutput{
					InscriptionID: data.InscriptionID,
					Offset:        offset,
					Sat:           data.Sat,
				}
				internalInfo, _ := u.Repo.FindTokenByTokenID(insc)
				if internalInfo != nil {
					inscWalletInfo.ProjectID = internalInfo.ProjectID
					inscWalletInfo.ProjecName = internalInfo.Project.Name
					inscWalletInfo.Thumbnail = internalInfo.Thumbnail
				}
				for _, listing := range currentListing {
					if listing.InscriptionID == data.InscriptionID {
						if listing.CancelTx == "" {
							inscWalletInfo.Buyable = true
						} else {
							inscWalletInfo.Cancelling = true
						}
						inscWalletInfo.OrderID = listing.UUID
						inscWalletInfo.PriceBTC = fmt.Sprintf("%v", listing.Amount)
					}
				}
				result[output] = append(result[output], inscWalletInfo)
				outputInscMap[output] = append(outputInscMap[output], inscWalletByOutput)
			}
		}
		outputSatRanges[output] = inscriptions.List.Unspent
	}
	// if len(outputSatRanges) != len(outputs) {
	// 	return nil, nil, nil, errors.New("")
	// }
	return result, outputInscMap, outputSatRanges, nil
}

func getInscriptionByOutput(ordServer, output string) (*structure.InscriptionOrdInfoByOutput, error) {
	url := fmt.Sprintf("%s/api/output/%s", ordServer, output)
	fmt.Println("url", url)
	var result structure.InscriptionOrdInfoByOutput
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getInscriptionByID(ordServer, id string) (*structure.InscriptionOrdInfoByID, error) {
	url := fmt.Sprintf("%s/api/inscription/%s", ordServer, id)
	fmt.Println("url", url)
	var result structure.InscriptionOrdInfoByID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func checkTxInBlockFromOrd(ordServer, txhash string) error {
	url := fmt.Sprintf("%s/tx/%s", ordServer, txhash)
	fmt.Println("url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	if res.StatusCode != http.StatusOK {
		return errors.New("getInscriptionByOutput Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Read body failed")
	}
	if strings.Contains(string(body), "not found") {
		return errors.New("tx not found")
	}

	return nil
}

func getWalletInfo(address string, apiToken string, logger logger.Ilogger) (*structure.BlockCypherWalletInfo, error) {
	url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s?unspentOnly=true&includeScript=false&token=%s", address, apiToken)
	var result structure.BlockCypherWalletInfo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getWalletInfo Response status != 200 " + result.Error + " " + url)
	}

	return &result, nil

}

func (u Usecase) TrackWalletTx(address string, tx structure.WalletTrackTx) error {
	trackTx := entity.WalletTrackTx{
		Address:           address,
		Txhash:            tx.Txhash,
		Type:              tx.Type,
		Amount:            tx.Amount,
		InscriptionID:     tx.InscriptionID,
		InscriptionNumber: tx.InscriptionNumber,
		Receiver:          tx.Receiver,
	}
	return u.Repo.CreateTrackTx(&trackTx)
}

func (u Usecase) GetWalletTrackTxs(address string, limit, offset int64) ([]structure.WalletTrackTx, error) {
	var result []structure.WalletTrackTx
	txList, err := u.Repo.GetTrackTxs(address, limit, offset)
	if err != nil {
		return nil, err
	}
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev.generativeexplorer.com"
	}
	for _, tx := range txList {
		createdAt := uint64(0)
		if tx.CreatedAt != nil {
			createdAt = uint64(tx.CreatedAt.Unix())
		}
		trackTx := structure.WalletTrackTx{
			Txhash:            tx.Txhash,
			Type:              tx.Type,
			Amount:            tx.Amount,
			InscriptionID:     tx.InscriptionID,
			InscriptionNumber: tx.InscriptionNumber,
			Receiver:          tx.Receiver,
			CreatedAt:         createdAt,
		}
		_, bs, err := u.buildBTCClient()
		if err != nil {
			fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
			return nil, err
		}
		if createdAt != 0 && time.Since(*tx.CreatedAt) >= 24*time.Hour {
			if err := checkTxInBlockFromOrd(ordServer, trackTx.Txhash); err == nil {
				trackTx.Status = "Success"
			} else {
				status, err := btc.GetBTCTxStatusExtensive(trackTx.Txhash, bs)
				if err != nil {
					return nil, err
				}
				trackTx.Status = status
			}
			result = append(result, trackTx)
		} else {
			status, err := btc.GetBTCTxStatusExtensive(trackTx.Txhash, bs)
			if err != nil {
				return nil, err
			}
			trackTx.Status = status
			result = append(result, trackTx)
		}
	}
	return result, nil
}
