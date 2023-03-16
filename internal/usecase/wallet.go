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
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/logger"
)

func (u Usecase) GetInscriptionByIDFromOrd(id string) (*structure.InscriptionOrdInfoByID, error) {
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}
	return getInscriptionByID(ordServer, id)
}

func (u Usecase) GetBTCWalletInfo(address string) (*structure.WalletInfo, error) {
	// cacheKey := utils.KEY_BTC_WALLET_INFO + "_" + address
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
	t := time.Now()
	apiToken := u.Config.BlockcypherToken
	u.Logger.Info("GetBTCWalletInfo apiToken debug", apiToken)
	quickNode := u.Config.QuicknodeAPI

	walletBasicInfo, err := btc.GetBalanceFromQuickNode(address, quickNode)
	if err != nil {
		var err2 error
		walletBasicInfo, err2 = getWalletInfo(address, apiToken, u.Logger)
		if err != nil {
			u.Logger.Info("GetBTCWalletInfo apiToken debug err", err2, err)
			return nil, err2
		}
	}
	trackT1 := time.Since(t)

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
	trackT2 := time.Since(t)

	inscriptions, outputInscMap, err := u.InscriptionsByOutputs(outcoins, currentListing)
	if err != nil {
		return nil, err
	}

	dupInscMap := make(map[string]struct{})
	result.InscriptionsByOutputs = outputInscMap

	for _, items := range inscriptions {
		for _, item := range items {
			if _, ok := dupInscMap[item.InscriptionID]; ok {
				continue
			}
			dupInscMap[item.InscriptionID] = struct{}{}
			result.Inscriptions = append(result.Inscriptions, item)
		}

	}
	trackT3 := time.Since(t)
	// newTxrefs := []structure.TxRef{}
	// for _, outcoin := range result.Txrefs {
	// 	o := fmt.Sprintf("%s:%v", outcoin.TxHash, outcoin.TxOutputN)
	// 	satRanges, ok := outputSatRanges[o]
	// 	if ok {
	// 		outcoin.SatRanges = satRanges
	// 		newTxrefs = append(newTxrefs, outcoin)
	// 	}
	// }
	// result.Txrefs = newTxrefs

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
				newTxrefsFiltered = append(newTxrefsFiltered, output)
			}
		}
		result.Txrefs = newTxrefsFiltered
	}
	result.Loadtime = make(map[string]string)
	result.Loadtime["trackT1"] = trackT1.String()
	result.Loadtime["trackT2"] = trackT2.String()
	result.Loadtime["trackT3"] = trackT3.String()

	// err = u.Repo.Cache.SetDataWithExpireTime(cacheKey, result, 10)
	// if err != nil {
	// 	u.Logger.Error("GetBTCWalletInfo CreateCache", address, err)
	// }
	// }

	return &result, nil
}

func (u Usecase) InscriptionsByOutputs(outputs []string, currentListing []entity.DexBTCListing) (map[string][]structure.WalletInscriptionInfo, map[string][]structure.WalletInscriptionByOutput, error) {
	result := make(map[string][]structure.WalletInscriptionInfo)
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}
	// outputSatRanges := make(map[string][][]uint64)
	outputInscMap := make(map[string][]structure.WalletInscriptionByOutput)
	var wg sync.WaitGroup
	var lock sync.Mutex
	waitChan := make(chan struct{}, 10)
	btcRate, ethRate, err := u.GetBTCToETHRate()
	if err != nil {
		log.Println("GenBuyETHOrder GetBTCToETHRate", err.Error(), err)
	}
	for _, output := range outputs {
		wg.Add(1)
		waitChan <- struct{}{}
		go func(op string) {
			defer func() {
				wg.Done()
				<-waitChan
			}()
			lock.Lock()
			if _, ok := result[op]; ok {
				lock.Unlock()
				return
			}
			lock.Unlock()

			inscriptions, err := getInscriptionByOutput(ordServer, op)
			if err != nil {
				return
			}
			if len(inscriptions.Inscriptions) > 0 {
				for _, insc := range inscriptions.Inscriptions {
					data, err := getInscriptionByID(ordServer, insc)
					if err != nil {
						return
					}
					tokenURI, err := u.Repo.FindTokenByTokenID(insc)
					if err != nil {
						// fmt.Errorf("FindTokenByTokenID error", err)
					}
					offset, err := strconv.ParseInt(strings.Split(data.Satpoint, ":")[2], 10, 64)
					if err != nil {
						return
					}
					inscWalletInfo := structure.WalletInscriptionInfo{
						InscriptionID: data.InscriptionID,
						Number:        data.Number,
						ContentType:   data.ContentType,
						Offset:        offset,
					}
					if tokenURI != nil {
						inscWalletInfo.TokenNumber = tokenURI.OrderInscriptionIndex
					}
					inscWalletByOutput := structure.WalletInscriptionByOutput{
						InscriptionID: data.InscriptionID,
						Offset:        offset,
						Sat:           data.Sat,
					}
					internalInfo, _ := u.Repo.FindTokenByTokenIDCustomField(insc, []string{"token_id", "project_id", "project.name", "thumbnail", "creator_address"})
					if internalInfo != nil {
						inscWalletInfo.ProjectID = internalInfo.ProjectID
						inscWalletInfo.Thumbnail = internalInfo.Thumbnail
						project, err := u.Repo.FindProjectByTokenIDCustomField(internalInfo.ProjectID, []string{"tokenid", "name"})
						if err != nil {
							log.Println("InscriptionsByOutputs.FindProjectByTokenIDCustomField", err)
						} else {
							inscWalletInfo.ProjectName = project.Name
						}
						creator, err := u.Repo.FindUserByAddress(internalInfo.CreatorAddr)
						if err != nil {
							log.Println("InscriptionsByOutputs.FindUserByAddress", err)
						} else {
							if creator != nil {
								inscWalletInfo.ArtistID = creator.UUID
								inscWalletInfo.ArtistName = creator.DisplayName
							}
						}
					}
					for _, listing := range currentListing {
						if listing.InscriptionID == data.InscriptionID {
							if listing.CancelTx == "" {
								inscWalletInfo.Buyable = true
							} else {
								inscWalletInfo.Cancelling = true
							}
							inscWalletInfo.SellVerified = listing.Verified
							inscWalletInfo.OrderID = listing.UUID
							inscWalletInfo.PriceBTC = fmt.Sprintf("%v", listing.Amount)

							amountBTCRequired := uint64(listing.Amount) + 1000
							amountBTCRequired += (amountBTCRequired / 10000) * 15 // + 0,15%

							amountETH, _, _, err := u.ConvertBTCToETHWithPriceEthBtc(fmt.Sprintf("%f", float64(amountBTCRequired)/1e8), btcRate, ethRate)
							if err != nil {
								log.Println("GenBuyETHOrder convertBTCToETH", err.Error(), err)
							}
							inscWalletInfo.PriceETH = amountETH
						}
					}
					lock.Lock()
					result[op] = append(result[op], inscWalletInfo)
					outputInscMap[op] = append(outputInscMap[op], inscWalletByOutput)
					lock.Unlock()
				}
			}
		}(output)
	}
	wg.Wait()
	return result, outputInscMap, nil
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
	// fmt.Println("url", url)
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

	// fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

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
	// fmt.Println("url", url)
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
		return nil, errors.New("getWalletInfo Response status != 200 " + result.Error)
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
	// t := time.Now()
	txList, err := u.Repo.GetTrackTxs(address, limit, offset)
	if err != nil {
		return nil, err
	}

	// t2 := time.Since(t)
	// log.Println("t2", t2)
	ordServer := os.Getenv("CUSTOM_ORD_SERVER")
	if ordServer == "" {
		ordServer = "https://dev-v5.generativeexplorer.com"
	}
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, item := range txList {
		wg.Add(1)
		// time.Sleep(10 * time.Millisecond)
		go func(tx entity.WalletTrackTx) {
			defer wg.Done()
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
				// return nil, err
				return
			}
			if createdAt != 0 && time.Since(*tx.CreatedAt) >= 24*time.Hour {
				if err := checkTxInBlockFromOrd(ordServer, trackTx.Txhash); err == nil {
					trackTx.Status = "Success"
				} else {
					// status, err := btc.GetBTCTxStatusExtensive(trackTx.Txhash, bs, u.Config.QuicknodeAPI)
					// if err != nil {
					// 	// return nil, err
					// 	return
					// }
					trackTx.Status = "Failed"
				}
				lock.Lock()
				result = append(result, trackTx)
				lock.Unlock()
			} else {
				status, err := btc.GetBTCTxStatusExtensive(trackTx.Txhash, bs, u.Config.QuicknodeAPI)
				if err != nil {
					// return nil, err
					return
				}
				trackTx.Status = status
				lock.Lock()
				result = append(result, trackTx)
				lock.Unlock()
			}
			// t3 := time.Since(t)
			// log.Println("t3 tx.Txhash", tx.Txhash, t3)
		}(item)
	}
	wg.Wait()
	// t3 := time.Since(t)
	// log.Println("t3", t3)
	return result, nil
}
