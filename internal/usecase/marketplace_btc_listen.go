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
	"rederinghub.io/internal/usecase/structure"
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

		txs, err := bs.GetLastTxs(item.HoldOrdAddress)

		if err != nil {
			go u.trackHistory("", "BtcChecktListNft", "", "", "bs.GetLastTxs at "+bs.GetEnpointURL(), err.Error())
		}

		if len(txs) == 0 {
			go u.trackHistory("", "GetLastTxs", "", "", "len(txs) from "+bs.GetEnpointURL(), "[]")
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

// check receive BTC for buying the nft:
func (u Usecase) BtcCheckReceivedBuyingNft(rootSpan opentracing.Span) error {

	fmt.Printf("go BtcCheckReceivedBuyingNft....")

	span, log := u.StartSpan("BtcCheckReceivedBuyingNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingBuyOrders()
	if len(listPending) == 0 {
		fmt.Printf("RetrieveBTCNFTPendingBuyOrders list empty")
		return nil
	}

	for _, item := range listPending {

		// check balance:

		balance, confirm, err := bs.GetBalance(item.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

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
		nftListing, err := u.Repo.FindBtcNFTListingByOrderID(item.ItemID)
		if err != nil {
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "FindBtcNFTListingByOrderID err", err.Error())
			continue
		}
		if nftListing == nil {

			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "FindBtcNFTListingByOrderID nil", "updated need to refund now")

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

		// update isSold
		nftListing.IsSold = true
		_, err = u.Repo.UpdateBTCNFTConfirmListings(nftListing)
		if err != nil {
			fmt.Printf("Could not UpdateBTCNFTConfirmListings id %s - with err: %v", item.ID, err)
			go u.trackHistory(item.ID.String(), "BtcCheckReceivedBuyingNft", item.TableName(), item.Status, "UpdateBTCNFTConfirmListings IsSold = true err", err.Error())
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
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SentNFT)

	fmt.Println("len(listTosendBtc)", len(listTosendBtc))

	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SentNFT {

			// get amount nft:
			nftListing, err := u.Repo.FindBtcNFTListingByOrderIDValid(item.ItemID)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "FindBtcNFTListingByOrderIDValid err", err.Error())
				continue
			}
			if nftListing == nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "FindBtcNFTListingByOrderIDValid nil", "[]")
				continue
			}

			// Todo cal amount to send user and master
			// send user first:
			totalAmount, ok := big.NewInt(0).SetString(nftListing.Price, 10)
			if !ok {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "SetString(nftListing.Price)", err.Error())
				continue
			}
			// charge x% total amount:
			serviceFee := int(float64(totalAmount.Int64()) * float64(utils.BUY_NFT_CHARGE) / 100)

			royaltyFee := int(0)

			tokenUri, err := u.GetTokenByTokenID(span, item.InscriptionID, 0)
			if err == nil {
				projectDetail, err := u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{
					ContractAddress: tokenUri.ContractAddress,
					ProjectID:       tokenUri.ProjectID,
				})
				if err != nil {
					log.Error("u.GetProjectDetail", err.Error(), err)
				}
				royaltyFeePercent := float64(projectDetail.Royalty / 10000)
				royaltyFee = int(float64(totalAmount.Int64()) * royaltyFeePercent)
				royaltyFee = 0 //TODO: lam
			}

			amountWithChargee := int(totalAmount.Uint64()) - serviceFee - royaltyFee
			fmt.Println("send btc from", item.SegwitAddress, "to: ", nftListing.SellerAddress)

			// transfer now:
			txID, err := bs.SendTransactionWithPreferenceFromSegwitAddress(
				item.SegwitKey,
				item.SegwitAddress,
				nftListing.SellerAddress,
				amountWithChargee,
				btc.PreferenceMedium,
			)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendBTCForBuyOrder", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
				continue
			}
			item.FeeChargeBTCBuyer = serviceFee
			item.RoyaltyChargeBTCBuyer = royaltyFee
			item.AmountBTCSentSeller = amountWithChargee
			item.TxSendBTC = txID
			item.Status = entity.StatusBuy_SendingBTC
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

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingBTC)
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

			if err == nil {
				go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+item.TxSendBTC, txResponse.Confirmations)
				if txResponse.Confirmations >= 1 {
					// send btc ok now:
					item.Status = entity.StatusBuy_SentBTC
					_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
					if err != nil {
						fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
					}
				}
			} else {
				fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
				go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendBTC, err.Error())

				go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendBTC, "Begin check tx via api.")

				// check with api:
				txInfo, err := bs.CheckTx(item.TxSendBTC)
				if err != nil {
					fmt.Printf("Could not bs - with err: %v", err)
					go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendBTC, err.Error())
				}
				if txInfo.Confirmations >= 1 {
					go u.trackHistory(item.ID.String(), "BtcCheckSendBTCForBuyOrder", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+item.TxSendBTC, txInfo.Confirmations)
					// send nft ok now:
					item.Status = entity.StatusBuy_SentBTC
					_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
					if err != nil {
						fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
					}
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

			// check nft in master wallet or not:
			listNFTsRep, err := u.GetMasterNfts(span)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.Error", err.Error())
				continue
			}

			go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.listNFTsRep", listNFTsRep)

			// parse nft data:
			var resp []struct {
				Inscription string `json:"inscription"`
				Location    string `json:"location"`
				Explorer    string `json:"explorer"`
			}

			err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
			if err != nil {
				go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.Unmarshal(listNFTsRep)", err.Error())
				continue
			}
			owner := false
			for _, nft := range resp {
				if strings.EqualFold(nft.Inscription, item.InscriptionID) {
					owner = true
					break
				}

			}
			go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.CheckNFTOwner", owner)
			if !owner {
				continue
			}

			// transfer now:
			sentTokenResp, err := u.SendTokenMKP(item.OrdAddress, item.InscriptionID)

			go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "SendTokenMKP.sentTokenResp", sentTokenResp)

			if err != nil {
				log.Error(fmt.Sprintf("BtcSendNFTForBuyOrder.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
				continue
			}

			//TODO: handle log err: Database already open. Cannot acquire lock

			// Update status first if none err:
			item.Status = entity.StatusBuy_SendingNFT
			item.ErrCount = 0 // reset error count!

			item.OutputSendNFT = sentTokenResp

			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				errPack := fmt.Errorf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err.Error())
				log.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", errPack.Error(), errPack)
				go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "SendTokenMKP.UpdateBTCNFTBuyOrder", err.Error())
				continue
			}

			txResp := sentTokenResp.Stdout
			//txResp := `fd31946b855cbaaf91df4b2c432f9b173e053e65a9879ac909bad028e21b950e\n`
			txResp = strings.ReplaceAll(txResp, "\n", "")

			// update tx:
			item.TxSendNFT = txResp
			item.ErrCount = 0 // reset error count!
			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				errPack := fmt.Errorf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				log.Error("BtcSendNFTForBuyOrder.Repo.UpdateBTCNFTBuyOrder", errPack.Error(), errPack)
				go u.trackHistory(item.ID.String(), "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "u.Repo.UpdateBTCNFTBuyOrder", err.Error())
			}
			// save log:
			log.SetData(fmt.Sprintf("BtcSendNFTForBuyOrder.execResp.%s", item.OrdAddress), sentTokenResp)
		}
	}
	return nil
}

func (u Usecase) BtcCheckSendNFTForBuyOrder(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcCheckSendNFTForBuyOrder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingNFT)
	if len(listTosendBtc) == 0 {
		fmt.Printf("BtcCheckSendNFTForBuyOrder empty")
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SendingNFT {

			if len(item.TxSendNFT) == 0 {
				// TODO ....
				continue
			}

			txHash, err := chainhash.NewHashFromStr(item.TxSendNFT)
			if err != nil {
				fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with tx: %v err: %v", item.TxSendNFT, err)
				continue
			}

			fmt.Println("txHash: ", txHash)

			txResponse, err := btcClient.GetTransaction(txHash)

			fmt.Println("txResponse of GetTransaction: ", txResponse)

			if err == nil {
				if txResponse.Confirmations >= 1 {
					// send nft ok now:
					item.Status = entity.StatusBuy_SentNFT
					_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
					if err != nil {
						fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
					}
				}
			} else {
				fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
				go u.trackHistory(item.ID.String(), "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendNFT, err.Error())

				go u.trackHistory(item.ID.String(), "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendNFT, "Begin check tx via api.")

				// check with api:
				txInfo, err := bs.CheckTx(item.TxSendNFT)
				if err != nil {
					fmt.Printf("Could not bs - with err: %v", err)
					go u.trackHistory(item.ID.String(), "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendNFT, err.Error())
				}
				if txInfo.Confirmations >= 1 {
					go u.trackHistory(item.ID.String(), "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+item.TxSendNFT, txInfo.Confirmations)
					// send nft ok now:
					item.Status = entity.StatusBuy_SentNFT
					_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
					if err != nil {
						fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
					}
				}
			}

		}
	}

	return nil
}

func (u Usecase) SendTokenMKP(receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_marketplace_master",
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
		}}

	resp, err := u.OrdService.Exec(sendTokenReq)
	return resp, err
}

func (u Usecase) GetMasterNfts(rootSpan opentracing.Span) (*ord_service.ExecRespose, error) {
	span, log := u.StartSpan(fmt.Sprintf("GetMasterNfts.%s", "inscriptions"), rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.TOKEN_ID_TAG, "inscriptions")
	log.SetTag(utils.WALLET_ADDRESS_TAG, "ord_marketplace_master")
	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_marketplace_master",
			"wallet",
			"inscriptions",
		}}

	log.SetData("listNFTsReq", listNFTsReq)
	resp, err := u.OrdService.Exec(listNFTsReq)
	defer u.Notify(rootSpan, "GetMasterNfts", "ord_marketplace_master", "inscriptions")
	if err != nil {
		log.SetData("u.OrdService.Exec.Error", err.Error())
		log.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	log.SetData("listNFTsRep", resp)
	return resp, err
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

// tesst:
func (u Usecase) SendTokenMKPTest(rootSpan opentracing.Span, walletName, receiveAddr, inscriptionID string) (*ord_service.ExecRespose, error) {
	span, log := u.StartSpan(fmt.Sprintf("SendTokenMKPTest.%s", inscriptionID), rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetTag(utils.TOKEN_ID_TAG, inscriptionID)
	log.SetTag(utils.WALLET_ADDRESS_TAG, receiveAddr)

	go u.trackHistory("test_send_nft", "SendTokenMKPTest", inscriptionID, receiveAddr, walletName, "before call ord_service.ExecRequest")

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletName,
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
		}}

	log.SetData("sendTokenReq", sendTokenReq)

	resp, err := u.OrdService.Exec(sendTokenReq)

	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "", "after call OrdService.Exec")
	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "SendTokenMKP.JsonTransform", resp)

	defer u.Notify(rootSpan, "SendTokenMKPTest", receiveAddr, inscriptionID)
	if err != nil {
		log.SetData("u.OrdService.Exec.Error", err.Error())
		log.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	log.SetData("sendTokenRes", resp)

	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "", "return now...")

	return resp, err
}
