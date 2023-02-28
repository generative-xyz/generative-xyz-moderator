package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/logger"
)

func (u Usecase) GetBTCWalletInfo(address string) (*structure.WalletInfo, error) {
	cacheKey := utils.KEY_BTC_WALLET_INFO + "_" + address
	var result structure.WalletInfo
	// exist, err := u.Repo.Cache.Exists(cacheKey)
	// if err == nil && *exist {
	// 	data, err := u.Repo.Cache.GetData(cacheKey)
	// 	if err == nil && data != nil {
	// 		err := json.Unmarshal([]byte(*data), &result)
	// 		if err != nil {
	// 			u.Logger.Error("GetBTCWalletInfo json.Unmarshal", address, err)
	// 		}
	// 		return &result, nil
	// 	}
	// }

	apiToken := u.Config.BlockcypherToken
	u.Logger.Info("GetBTCWalletInfo apiToken debug", apiToken)
	walletBasicInfo, err := getWalletInfo(address, apiToken, u.Logger)
	if err != nil {
		u.Logger.Info("GetBTCWalletInfo apiToken debug err", err)
		return nil, err
	}

	result.BlockCypherWalletInfo = *walletBasicInfo
	outcoins := []string{}
	for _, outcoin := range result.Txrefs {
		o := fmt.Sprintf("%s:%v", outcoin.TxHash, outcoin.TxOutputN)
		outcoins = append(outcoins, o)
	}
	inscriptions, outputInscMap, outputSatRanges, err := u.InscriptionsByOutputs(outcoins)
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

	resultBytes, err := json.Marshal(result)
	if err != nil {
		u.Logger.Error("GetBTCWalletInfo json.Marshal", address, err)
	} else {
		err = u.Repo.Cache.SetDataWithExpireTime(cacheKey, string(resultBytes), 60)
		if err != nil {
			u.Logger.Error("GetBTCWalletInfo CreateCache", address, err)
		}
	}

	return &result, nil
}

func (u Usecase) InscriptionsByOutputs(outputs []string) (map[string][]structure.WalletInscriptionInfo, map[string][]structure.WalletInscriptionByOutput, map[string][][]uint64, error) {
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

type BTCTxInfo struct {
	Data struct {
		BlockHeight   int         `json:"block_height"`
		BlockHash     string      `json:"block_hash"`
		BlockTime     int         `json:"block_time"`
		CreatedAt     int         `json:"created_at"`
		Confirmations int         `json:"confirmations"`
		Fee           int         `json:"fee"`
		Hash          string      `json:"hash"`
		InputsCount   int         `json:"inputs_count"`
		InputsValue   int         `json:"inputs_value"`
		IsCoinbase    bool        `json:"is_coinbase"`
		IsDoubleSpend bool        `json:"is_double_spend"`
		IsSwTx        bool        `json:"is_sw_tx"`
		LockTime      int         `json:"lock_time"`
		OutputsCount  int         `json:"outputs_count"`
		OutputsValue  int64       `json:"outputs_value"`
		Sigops        int         `json:"sigops"`
		Size          int         `json:"size"`
		Version       int         `json:"version"`
		Vsize         int         `json:"vsize"`
		Weight        int         `json:"weight"`
		WitnessHash   string      `json:"witness_hash"`
		Inputs        interface{} `json:"inputs"`
		Outputs       interface{} `json:"outputs"`
	} `json:"data"`
	ErrCode int    `json:"err_code"`
	ErrNo   int    `json:"err_no"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

var btcRateLock sync.Mutex

func checkTxFromBTC(txhash string) (*BTCTxInfo, error) {
	btcRateLock.Lock()
	defer func() {
		time.Sleep(100 * time.Millisecond)
		btcRateLock.Unlock()
	}()
	url := fmt.Sprintf("https://chain.api.btc.com/v3/tx/%s?verbose=1", txhash)
	fmt.Println("url", url)
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

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getInscriptionByOutput Response status != 200")
	}
	var result BTCTxInfo

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(string(body))
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("getWalletInfo Response status != 200 " + result.Message + " " + url)
	}

	return &result, nil
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
		if createdAt != 0 && time.Since(*tx.CreatedAt) >= 1*time.Hour {
			if err := checkTxInBlockFromOrd(ordServer, trackTx.Txhash); err == nil {
				trackTx.Status = "Success"
			} else {
				err = getBTCTxStatusExtensive(&trackTx, &u)
				if err != nil {
					return nil, err
				}
				result = append(result, trackTx)
			}
		} else {
			err = getBTCTxStatusExtensive(&trackTx, &u)
			if err != nil {
				return nil, err
			}
			result = append(result, trackTx)
		}
	}
	return result, nil
}

func getBTCTxStatusExtensive(trackTx *structure.WalletTrackTx, u *Usecase) error {
	_, bs, err := u.buildBTCClient()
	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	txStatus, err := bs.CheckTx(trackTx.Txhash)
	if err != nil {
		txInfo, err := checkTxFromBTC(trackTx.Txhash)
		if err != nil {
			fmt.Printf("checkTxFromBTC err: %v", err)
			trackTx.Status = "Failed"
		} else {
			if txInfo.Data.Confirmations > 0 {
				trackTx.Status = "Success"
			} else {
				trackTx.Status = "Pending"
			}
		}
	} else {
		if txStatus.Confirmations > 0 {
			trackTx.Status = "Success"
		} else {
			trackTx.Status = "Pending"
		}
	}
	return nil
}
