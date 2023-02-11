package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
)

func (u Usecase) buildBTCClient() (*rpcclient.Client, *btc.BlockcypherService, error) {
	host := u.Config.BTC_FULLNODE
	user := u.Config.BTC_RPCUSER
	pass := u.Config.BTC_RPCPASSWORD

	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   false, //!(os.Getenv("BTC_NODE_HTTPS") == "true"), // Bitcoin core does not provide TLS by default
	}

	rpcclient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, nil, err
	}

	bs := btc.NewBlockcypherService(u.Config.BlockcypherAPI, "", u.Config.BlockcypherToken, &chaincfg.MainNetParams)

	return rpcclient, bs, nil
}

func (u Usecase) loopGetTx(btcClient *rpcclient.Client, tx string, item *entity.MarketplaceBTCListing) (string, bool, error) {

	txOut := ""

	detail, err := chainhash.NewHashFromStr(tx)
	if err != nil {
		fmt.Println("can not NewHashFromStr with err:", err)
		go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "chainhash.NewHashFromStr- with err", err.Error())
		return txOut, false, err
	}
	result, err := btcClient.GetRawTransactionVerboseAsync(detail).Receive()
	if err != nil {
		fmt.Println("can not GetRawTransactionVerboseAsync with err:", err)
		go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "GetRawTransactionVerboseAsync- with err", err.Error())
		return txOut, false, err
	}

	for _, vin := range result.Vin {
		fmt.Println("vin==>", vin.Txid)
		fmt.Println("item.InscriptionID==>", item.InscriptionID)

		txOut = vin.Txid

		if strings.Contains(item.InscriptionID, txOut) {
			return txOut, true, nil
		}
	}

	return txOut, false, nil
}

// check receive of the nft:
func (u Usecase) BtcChecktListNft(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcChecktListNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackHistory("", "BtcChecktListNft", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingListings()
	if len(listPending) == 0 {
		// go u.trackHistory("", "BtcChecktListNft", "", "", "RetrieveBTCNFTPendingListings", "[]")
		return nil
	}

	for _, item := range listPending {

		if len(item.InscriptionID) == 0 {
			continue
		}

		txs, _ := bs.GetLastTxs(item.HoldOrdAddress)

		if len(txs) == 0 {
			go u.trackHistory("", "GetLastTxs", "", "", "len(txs) ", "[]")
			continue
		}

		found := false

		for _, tx := range txs {
			detail, err := chainhash.NewHashFromStr(tx.Tx)
			if err != nil {
				fmt.Println("can not NewHashFromStr with err:", err)
				go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "chainhash.NewHashFromStr- with err", err.Error())
				continue
			}
			result, err := btcClient.GetRawTransactionVerboseAsync(detail).Receive()
			if err != nil {
				fmt.Println("can not GetRawTransactionVerboseAsync with err:", err)
				go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "GetRawTransactionVerboseAsync- with err", err.Error())
				continue
			}

			for _, vin := range result.Vin {
				fmt.Println("vin==>", vin.Txid)
				fmt.Println("item.InscriptionID==>", item.InscriptionID)
				if strings.Contains(item.InscriptionID, vin.Txid) {
					found = true
					item.IsConfirm = true
					item.TxNFT = vin.Txid
					_, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
					if err != nil {
						go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "UpdateBTCNFTConfirmListings - with err", err.Error())
					}
					break
				} else {

					tx := vin.Txid

					for i := 0; i < 20; i++ {

						fmt.Println("count: ", i+1, "tx: ", tx)

						hash, f, err := u.loopGetTx(btcClient, tx, &item)

						if err != nil {
							go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "loopGetTx - with err", err.Error())
							break
						}
						if f {
							fmt.Println("found nft for listing: ", "tx: ", hash)
							found = true
							item.IsConfirm = true
							item.TxNFT = hash
							updated, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
							fmt.Println("UpdateBTCNFTConfirmListings updated: ", updated)
							if err != nil {
								fmt.Println("can not UpdateBTCNFTConfirmListings err ", err)
								go u.trackHistory(item.ID.String(), "BtcChecktListNft", item.TableName(), item.IsConfirm, "UpdateBTCNFTConfirmListings - with err", err.Error())
							}
							break
						}
						if len(hash) == 0 {
							break
						}
						tx = hash
					}
				}
			}
			if found {
				break
			}
		}
	}

	return nil
}

// check receive buy the nft:
func (u Usecase) BtcCheckReceivedBuyingNft(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcCheckReceivedBuyingNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingBuyOrders()
	if len(listPending) == 0 {
		return nil
	}

	for _, item := range listPending {

		// check balance:

		balance, err := bs.GetBalance(item.SegwitAddress)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get amount nft:
		nftListing, err := u.Repo.FindBtcNFTListingByNFTID(item.InscriptionID)
		if err != nil {
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "FindBtcNFTListingByNFTID err", err.Error())
			continue
		}
		if nftListing == nil {

			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "FindBtcNFTListingByNFTID nil", "updated need to refund now")

			// update StatusBuy_NeedToRefund now for listing:
			item.Status = entity.StatusBuy_NeedToRefund
			log.SetData(fmt.Sprintf("BtcCheckBuyingNft.CheckReceiveNFT.%s", item.SegwitAddress), item)
			u.Notify(rootSpan, "WaitingForBTCBalancingOfBuyOrder", item.SegwitAddress, fmt.Sprintf("%s Need to refund BTC %s from [InscriptionID] %s", item.SegwitAddress, item.ReceivedBalance, item.InscriptionID))

			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "UpdateBTCNFTBuyOrder err", err.Error())
			}
			continue
		}
		// update isSold
		nftListing.IsSold = true
		_, err = u.Repo.UpdateBTCNFTConfirmListings(nftListing)
		if err != nil {
			fmt.Printf("Could not UpdateBTCNFTConfirmListings id %s - with err: %v", item.ID, err)
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "UpdateBTCNFTConfirmListings IsSold = true err", err.Error())
		}

		amount, _ := big.NewInt(0).SetString(nftListing.Price, 10)

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if r := balance.Cmp(amount); r == -1 {
			err := errors.New("Not enough amount")
			fmt.Printf("buy order id: %s err: %v", item.InscriptionID, err)

			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "Receive balance err", err.Error())

			item.Status = entity.StatusBuy_NotEnoughBalance
			u.Repo.UpdateBTCNFTBuyOrder(&item)
			continue
		}

		item.Status = entity.StatusBuy_ReceivedFund

		_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
			continue
		}

		go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "Updated StatusBuy_ReceivedFund", "ok")
		log.SetData(fmt.Sprintf("BtcCheckBuyingNft.CheckReceiveNFT.%s", item.SegwitAddress), item)
		u.Notify(rootSpan, "WaitingForBTCBalancingOfBuyOrder", item.SegwitAddress, fmt.Sprintf("%s received BTC %s from [InscriptionID] %s", item.SegwitAddress, item.ReceivedBalance, item.InscriptionID))

	}

	return nil
}

// send btc for buy order records:
func (u Usecase) BtcSendBTCForBuyOrder(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("BtcSendBTCForBuyOrder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingNFT)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SendingNFT {

			// get amount nft:
			nftListing, err := u.Repo.FindBtcNFTListingByNFTIDValid(item.InscriptionID)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "FindBtcNFTListingByNFTIDValid err", err.Error())
				continue
			}
			if nftListing == nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "FindBtcNFTListingByNFTIDValid nil", "[]")
				continue
			}

			// Todo cal amount to send user and master
			// send user first:
			receiveAmount, _ := big.NewInt(0).SetString(item.ReceivedBalance, 10)
			// charge 10% total amount:
			amountWithChargee := int(receiveAmount.Uint64()) - int(receiveAmount.Uint64()*utils.BUY_NFT_CHARGE/100)

			// transfer now:
			txID, err := bs.SendTransactionWithPreferenceFromSegwitAddress(
				item.SegwitKey,
				nftListing.SellOrdAddress,
				item.SegwitAddress,
				amountWithChargee,
				btc.PreferenceMedium,
			)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
				continue
			}

			item.TxSendBTC = txID
			item.ErrCount = 0 // reset error count!
			// TODO: update item
			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
			}

		}
	}
	return nil
}

func (u Usecase) BtcCheckSendBTCForBuyOrder(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcCheckSendBTCForBuyOrder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, _, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingNFT)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SendingBTC {
			txHash, err := chainhash.NewHashFromStr(item.TxSendBTC)
			if err != nil {
				fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
				continue
			}

			txResponse, err := btcClient.GetTransaction(txHash)

			if err != nil {
				fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
				go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendBTC, err.Error())
				continue
			}
			go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+item.TxSendBTC, txResponse.Confirmations)
			if txResponse.Confirmations >= 1 {
				// send btc ok now:
				item.Status = entity.StatusBuy_SentBTC
				_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
				if err != nil {
					fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				}
			}
		}
	}

	return nil
}

// send nft for buy order records:
func (u Usecase) BtcSendNFTForBuyOrder(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("BtcSendBTCForBuyOrder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	// get list buy order status = StatusBuy_ReceivedFund:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_ReceivedFund)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_ReceivedFund {

			// transfer now:
			sentTokenResp, err := u.SendToken(rootSpan, item.OrdAddress, item.InscriptionID)
			if err != nil {
				log.Error(fmt.Sprintf("BtcSendNFTForBuyOrder.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
				go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "SendToken", err.Error())
				continue
			}
			tmpText := sentTokenResp.Stdout
			//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
			jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
			jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
			btcMintResp := &ord_service.MintStdoputRespose{}

			bytes := []byte(jsonStr)
			err = json.Unmarshal(bytes, btcMintResp)
			if err != nil {
				log.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", err.Error(), err)
				continue
			}

			log.SetData(fmt.Sprintf("BtcSendNFTForBuyOrder.execResp.%s", item.OrdAddress), sentTokenResp)
			item.TxSendNFT = btcMintResp.Commit
			item.Status = entity.StatusBuy_SendingNFT
			item.ErrCount = 0 // reset error count!
			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				errPack := fmt.Errorf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				log.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", errPack.Error(), errPack)
			}
		}
	}
	return nil
}

func (u Usecase) BtcCheckSendNFTForBuyOrder(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcCheckSendNFTForBuyOrder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, _, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingNFT)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SendingNFT {
			txHash, err := chainhash.NewHashFromStr(item.TxSendNFT)
			if err != nil {
				fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
				continue
			}

			txResponse, err := btcClient.GetTransaction(txHash)

			if err != nil {
				fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
				go u.trackHistory(item.ID.String(), "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendBTC, err.Error())
				continue
			}
			if txResponse.Confirmations >= 1 {
				// send nft ok now:
				item.Status = entity.StatusBuy_SentNFT
				_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
				if err != nil {
					fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				}
			}

		}
	}

	return nil
}

func (u *Usecase) trackHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}) {
	trackData := &entity.MarketplaceBTCLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.CreateMarketplaceBTCLog(trackData)
	if err != nil {
		fmt.Printf("trackHistory.%s.Error:%s", name, err.Error())
	}

}
