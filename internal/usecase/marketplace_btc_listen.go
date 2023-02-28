package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/encrypt"
	"rederinghub.io/utils/eth"
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
		go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "chainhash.NewHashFromStr- with err", err.Error())
		return txOut, false, err
	}
	result, err := btcClient.GetRawTransactionVerboseAsync(detail).Receive()
	if err != nil {
		fmt.Println("can not GetRawTransactionVerboseAsync with err:", err)
		go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "GetRawTransactionVerboseAsync- with err", err.Error())
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

//JOBS:
// job1:  check receive the nft from seller:
func (u Usecase) JobMKP_BtcChecktListNft() error {

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
				go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "chainhash.NewHashFromStr- with err", err.Error())
				continue
			}
			result, err := btcClient.GetRawTransactionVerboseAsync(detail).Receive()
			if err != nil {
				fmt.Println("can not GetRawTransactionVerboseAsync with err:", err)
				go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "GetRawTransactionVerboseAsync- with err", err.Error())
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
						go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "UpdateBTCNFTConfirmListings - with err", err.Error())
					}
					break
				} else {

					tx := vin.Txid

					for i := 0; i < 20; i++ {

						fmt.Println("count: ", i+1, "tx: ", tx)

						hash, f, err := u.loopGetTx(btcClient, tx, &item)

						if err != nil {
							go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "loopGetTx - with err", err.Error())
							break
						}
						if f {
							fmt.Println("found nft for listing: ", "tx: ", hash)
							found = true
							item.IsConfirm = true
							item.Status = entity.StatusListing_ReceivedNft
							item.TxNFT = hash
							updated, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
							fmt.Println("UpdateBTCNFTConfirmListings updated: ", updated)
							if err != nil {
								fmt.Println("can not UpdateBTCNFTConfirmListings err ", err)
								go u.trackHistory(item.UUID, "BtcChecktListNft", item.TableName(), item.IsConfirm, "UpdateBTCNFTConfirmListings - with err", err.Error())
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

// job2: check receive fund (eth/btc) from buyer:
func (u Usecase) JobMKP_WaitForBalanceFromBuyer() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMKP_WaitForBalanceFromBuyer", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	listPending, _ := u.Repo.RetrieveBTCNFTPendingBuyOrders()
	if len(listPending) == 0 {
		fmt.Printf("RetrieveBTCNFTPendingBuyOrders list empty")
		return nil
	}

	for _, item := range listPending {

		// check balance:

		// check balance:
		balance := big.NewInt(0)
		confirm := -1

		if len(item.PayType) == 0 {
			item.PayType = utils.NETWORK_BTC
		}

		if item.PayType == utils.NETWORK_BTC {
			balance, confirm, err = bs.GetBalance(item.ReceiveAddress)

			fmt.Println("GetBalance response: ", balance, confirm, err)

		} else if item.PayType == utils.NETWORK_ETH {

			// check eth balance:

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			balance, err = ethClient.GetBalance(ctx, item.ReceiveAddress)
			fmt.Println("GetBalance eth response: ", balance, err)

			confirm = 1
		}

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin/Ether - with err: %v", err)
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get amount nft:
		nftListing, err := u.Repo.FindBtcNFTListingByOrderID(item.ItemID)
		if err != nil {
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "FindBtcNFTListingByOrderID err", err.Error())
			continue
		}
		if nftListing == nil {

			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "FindBtcNFTListingByOrderID nil", "updated need to refund now")

			// update StatusBuy_NeedToRefund now for listing:
			item.Status = entity.StatusBuy_NeedToRefund
			u.Logger.Info(fmt.Sprintf("JobMKP_WaitForBalanceFromBuyer.CheckReceiveNFT.%s", item.ReceiveAddress), item)
			u.Notify("WaitingForBTCBalancingOfBuyOrder", item.ReceiveAddress, fmt.Sprintf("%s Need to refund BTC %s from [InscriptionID] %s", item.ReceiveAddress, item.ReceivedBalance, item.InscriptionID))

			_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
			if err != nil {
				go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "UpdateBTCNFTBuyOrder err", err.Error())
			}
			continue
		}

		amount, ok := big.NewInt(0).SetString(nftListing.Price, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackInscribeHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "SetString(amount) err", err.Error())
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if r := balance.Cmp(amount); r == -1 {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			fmt.Printf("buy order id: %s err: %v", item.InscriptionID, err)

			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "Receive balance err", err.Error())

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
			go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "UpdateBTCNFTConfirmListings IsSold = true err", err.Error())
		}

		go u.trackHistory(item.UUID, "JobMKP_WaitForBalanceFromBuyer", item.TableName(), item.Status, "Updated StatusBuy_ReceivedFund", "ok")
		u.Logger.Info(fmt.Sprintf("JobMKP_WaitForBalanceFromBuyer.CheckReceiveNFT.%s", item.ReceiveAddress), item)
		u.Notify("JobMKP_WaitForBalanceFromBuyer", item.ReceiveAddress, fmt.Sprintf("%s received BTC %s from [InscriptionID] %s", item.ReceiveAddress, item.ReceivedBalance, item.InscriptionID))

	}

	return nil
}

// send btc/eth for buy order records:
func (u Usecase) JobMKP_Payment() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMKP_WaitForBalanceFromBuyer", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list buy order status = sent nft:
	listTosendBtc, err := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SentNFT)

	fmt.Println("len(listTosendBtc)", len(listTosendBtc), err)

	if len(listTosendBtc) == 0 {
		return nil
	}
	serviceFeeAddress := "bc1q2a7j7zxqc0l43xd9urahxywqt7zl462hgpm0wh"
	if u.Config.MarketBTCServiceFeeAddress != "" {
		serviceFeeAddress = u.Config.MarketBTCServiceFeeAddress
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SentNFT {

			nftListing, err := u.Repo.FindBtcNFTListingByOrderIDValid(item.ItemID)
			if err != nil {
				go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "FindBtcNFTListingByOrderIDValid err", err.Error())
				continue
			}
			if nftListing == nil {
				go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "FindBtcNFTListingByOrderIDValid nil", "[]")
				continue
			}

			if len(item.PayType) == 0 {
				item.PayType = utils.NETWORK_BTC
			}

			// check paid type:
			if item.PayType == utils.NETWORK_BTC {

				// get amount nft:

				// Todo cal amount to send user and master
				// send user first:
				totalAmount, ok := big.NewInt(0).SetString(nftListing.Price, 10)
				if !ok {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "SetString(nftListing.Price)", err.Error())
					continue
				}
				// charge x% total amount:
				serviceFee := int(float64(totalAmount.Int64()) * float64(utils.BUY_NFT_CHARGE) / 100)

				royaltyFee := int(0)
				artistAddress := ""
				tokenUri, err := u.GetTokenByTokenID(item.InscriptionID, 0)
				if err == nil {
					projectDetail, err := u.GetProjectDetail(structure.GetProjectDetailMessageReq{
						ContractAddress: tokenUri.ContractAddress,
						ProjectID:       tokenUri.ProjectID,
					})
					if err == nil {
						if projectDetail.Royalty > 0 {
							creator, err := u.GetUserProfileByWalletAddress(projectDetail.CreatorAddrr)
							if err == nil {
								if creator.WalletAddressBTC != "" {
									royaltyFeePercent := float64(projectDetail.Royalty) / 10000
									royaltyFee = int(float64(totalAmount.Int64()) * royaltyFeePercent)
									artistAddress = creator.WalletAddressBTC
								}
							}
						}
					}
				}

				amountWithChargee := int(totalAmount.Uint64()) - serviceFee - royaltyFee

				fmt.Println("send btc from", item.ReceiveAddress, "to: ", nftListing.PayType["btc"])

				destinations := make(map[string]int)

				destinations[nftListing.PayType["btc"]] = amountWithChargee
				if artistAddress != "" && royaltyFee > 0 {
					if artistAddress == nftListing.PayType["btc"] {
						amountWithChargee = amountWithChargee + royaltyFee
					} else {
						destinations[artistAddress] = royaltyFee
					}
				}

				if serviceFee > 0 {
					destinations[serviceFeeAddress] = serviceFee
				}

				txFee, err := bs.EstimateFeeTransactionWithPreferenceFromSegwitAddressMultiAddress(item.SegwitKey, item.ReceiveAddress, destinations, btc.PreferenceMedium)
				if err != nil {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "EstimateFeeTransactionWithPreferenceFromSegwitAddressMultiAddress err", err.Error())
					continue
				}
				amountWithChargee = amountWithChargee - int(txFee.Int64())
				destinations[nftListing.PayType["btc"]] = amountWithChargee

				privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
				if err != nil {
					u.Logger.Error(fmt.Sprintf("JobMKP_Payment.Decrypt.%s.Error", item.ReceiveAddress), err.Error(), err)
					go u.trackMintNftBtcHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "JobMKP_Payment.DecryptToString", err.Error(), true)
					continue
				}

				txID, err := bs.SendTransactionWithPreferenceFromSegwitAddressMultiAddress(
					privateKeyDeCrypt,
					item.ReceiveAddress,
					destinations,
					btc.PreferenceMedium,
				)
				if err != nil {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "SendTransactionWithPreferenceFromReceiveAddressMultiAddress err", err.Error())
					continue
				}
				// // transfer now:
				// txID, err := bs.SendTransactionWithPreferenceFromReceiveAddress(
				// 	item.SegwitKey,
				// 	item.ReceiveAddress,
				// 	nftListing.SellOrdAddress,
				// 	amountWithChargee,
				// 	btc.PreferenceMedium,
				// )
				// if err != nil {
				// 	go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "SendTransactionWithPreferenceFromReceiveAddress err", err.Error())
				// 	continue
				// }
				item.FeeChargeBuyer = fmt.Sprintf("%d", serviceFee)
				item.RoyaltyChargeBuyer = fmt.Sprintf("%d", royaltyFee)
				item.AmountSentSeller = fmt.Sprintf("%d", amountWithChargee)
				item.TxSendPayment = txID
				item.Status = entity.StatusBuy_SendingPayment
				item.ErrCount = 0 // reset error count!
				// TODO: update item
				_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
				if err != nil {
					fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				}
			} else if item.PayType == utils.NETWORK_ETH {
				_ = ethClient

				totalAmountByEth, ok := big.NewInt(0).SetString(item.Price, 10)
				if !ok {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "SetString(nftListing.Price)", err.Error())
					continue
				}
				// charge x% total amount:
				serviceFeeByEth := uint64(float64(totalAmountByEth.Uint64()) * float64(utils.BUY_NFT_CHARGE) / 100)

				royaltyFee := uint64(0)
				artistAddress := ""
				tokenUri, err := u.GetTokenByTokenID(item.InscriptionID, 0)
				if err == nil {
					projectDetail, err := u.GetProjectDetail(structure.GetProjectDetailMessageReq{
						ContractAddress: tokenUri.ContractAddress,
						ProjectID:       tokenUri.ProjectID,
					})
					if err == nil {
						if projectDetail.Royalty > 0 {
							creator, err := u.GetUserProfileByWalletAddress(projectDetail.CreatorAddrr)
							if err == nil {
								if creator.WalletAddressPayment != "" {
									royaltyFeePercent := float64(projectDetail.Royalty) / 10000
									royaltyFee = uint64(float64(totalAmountByEth.Uint64()) * royaltyFeePercent)
									artistAddress = creator.WalletAddressPayment
								}
							}
						}
					}
				}

				fmt.Println("serviceFeeByEth:", serviceFeeByEth)
				fmt.Println("royaltyFee:", royaltyFee)

				amountWithChargee := uint64(totalAmountByEth.Uint64()) - serviceFeeByEth - royaltyFee
				fmt.Println("send eth from", item.ReceiveAddress, "to: ", nftListing.PayType["eth"], "amountWithChargee: ", amountWithChargee)

				destinations := make(map[string]*big.Int)

				destinations[nftListing.PayType["eth"]] = big.NewInt(int64(amountWithChargee))
				if artistAddress != "" && royaltyFee > 0 {
					destinations[artistAddress] = big.NewInt(int64(royaltyFee))
				}

				if serviceFeeByEth > 0 {
					destinations[serviceFeeAddress] = big.NewInt(int64(serviceFeeByEth))
					fmt.Println("serviceFeeByEth: ", serviceFeeByEth)
				}

				balance, err := ethClient.GetClient().BalanceAt(context.Background(), common.HexToAddress(item.ReceiveAddress), nil)

				fmt.Println("eth balance: ", balance)

				if err != nil {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "BalanceAt err", err.Error())
					continue
				}

				if balance != nil && balance.Uint64() == 0 {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "Balance", "0")
					continue
				}

				//TODO: send to 3 output
				gasPrice, err := ethClient.GetClient().SuggestGasPrice(context.Background())
				if err != nil {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "SuggestGasPrice err", err.Error())
					continue
				}

				fmt.Println("gasPrice: ", gasPrice, len(destinations))

				txFee := new(big.Int).Mul(new(big.Int).SetUint64(gasPrice.Uint64()), new(big.Int).SetUint64(21000*uint64(len(destinations))))

				fmt.Println("txFee: ", txFee) // ~ 0.0004 eth

				if txFee.Uint64() > balance.Uint64() {
					fmt.Println("not enough balance: ", txFee.Uint64(), balance.Uint64())
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "check fee and balance", "txFee > balance")
					continue
				}

				if amountWithChargee < txFee.Uint64() {
					fmt.Println("not enough amountWithChargee: ", txFee.Uint64(), balance.Uint64())
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "check fee and balance", "txamountWithChargeeFee < txFee")
					continue
				}

				amountWithChargee = amountWithChargee - txFee.Uint64()
				destinations[nftListing.PayType["eth"]] = big.NewInt(int64(amountWithChargee))

				fmt.Println("destinations: ", destinations)

				// start send now:
				privateKeyDeCrypt, err := encrypt.DecryptToString(item.PrivateKey, os.Getenv("SECRET_KEY"))
				if err != nil {
					u.Logger.Error(fmt.Sprintf("JobMKP_Payment.Decrypt.%s.Error", item.ReceiveAddress), err.Error(), err)
					go u.trackMintNftBtcHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "JobMKP_Payment.DecryptToString", err.Error(), true)
					continue
				}

				txID, err := ethClient.SendMulti(
					"0xcd5485b34c9902527bbee21f69312fe2a73bc802",
					privateKeyDeCrypt,
					destinations,
				)
				if err != nil {
					go u.trackHistory(item.UUID, "JobMKP_Payment", item.TableName(), item.Status, "ethClientSendMulti", err.Error())
					continue
				}
				item.FeeChargeBuyer = fmt.Sprintf("%d", serviceFeeByEth)
				item.RoyaltyChargeBuyer = fmt.Sprintf("%d", royaltyFee)
				item.AmountSentSeller = fmt.Sprintf("%d", amountWithChargee)
				item.TxSendPayment = txID
				item.Status = entity.StatusBuy_SendingPayment
				item.ErrCount = 0 // reset error count!
				// TODO: update item
				_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
				if err != nil {
					fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
				}

			}

		}
	}
	return nil
}

func (u Usecase) JpbMKP_CheckTxSendPayment() error {

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackMintNftBtcHistory("", "JobMKP_WaitForBalanceFromBuyer", "", "", "Could not initialize Ether RPCClient - with err", err.Error(), true)
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list buy order status = sent nft:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_SendingPayment)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_SendingPayment {

			if item.PayType == utils.NETWORK_BTC {

				txHash, err := chainhash.NewHashFromStr(item.TxSendPayment)
				if err != nil {
					fmt.Printf("Could not NewHashFromStr Bitcoin RPCClient - with err: %v", err)
					continue
				}

				txResponse, err := btcClient.GetTransaction(txHash)

				if err == nil {
					go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+item.TxSendPayment, txResponse.Confirmations)
					if txResponse.Confirmations >= 1 {
						// send btc ok now:
						item.Status = entity.StatusBuy_SentPayment
						_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
						if err != nil {
							fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
						}
					}
				} else {
					fmt.Printf("Could not GetTransaction Bitcoin RPCClient - with err: %v", err)
					go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendPayment, err.Error())

					go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendPayment, "Begin check tx via api.")

					// check with api:
					txInfo, err := bs.CheckTx(item.TxSendPayment)
					if err != nil {
						fmt.Printf("Could not bs - with err: %v", err)
						go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendPayment, err.Error())
					}
					if txInfo.Confirmations >= 1 {
						go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+item.TxSendPayment, txInfo.Confirmations)
						// send nft ok now:
						item.Status = entity.StatusBuy_SentPayment
						_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
						if err != nil {
							fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
						}
					}
				}
			} else if item.PayType == utils.NETWORK_ETH {
				context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				status, err := ethClient.GetTransaction(context, item.TxSendPayment)
				if err == nil {
					if status > 0 {
						go u.trackHistory(item.UUID, "JpbMKP_CheckTxSendPaymen", item.TableName(), item.Status, "ethClient.GetTransaction: "+item.TxSendPayment, status)
						// send nft ok now:
						item.Status = entity.StatusBuy_SentPayment
						_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
						if err != nil {
							fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
						}

					} else {
						continue
					}
				} else {
					// if error maybe tx is pending or rejected
					// TODO check timeout to detect tx is rejected or not.
				}
			}

		}
	}

	return nil
}

// send nft for buy order records:
func (u Usecase) BtcSendNFTForBuyOrder() error {

	// get list buy order status = StatusBuy_ReceivedFund:
	listTosendBtc, _ := u.Repo.RetrieveBTCNFTBuyOrdersByStatus(entity.StatusBuy_ReceivedFund)
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusBuy_ReceivedFund {

			// check nft in master wallet or not:
			listNFTsRep, err := u.GetMasterNfts()
			if err != nil {
				go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.Error", err.Error())
				continue
			}

			go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.listNFTsRep", listNFTsRep)

			// parse nft data:
			var resp []struct {
				Inscription string `json:"inscription"`
				Location    string `json:"location"`
				Explorer    string `json:"explorer"`
			}

			err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
			if err != nil {
				go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.Unmarshal(listNFTsRep)", err.Error())
				continue
			}
			owner := false
			for _, nft := range resp {
				if strings.EqualFold(nft.Inscription, item.InscriptionID) {
					owner = true
					break
				}

			}
			go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "GetMasterNfts.CheckNFTOwner", owner)
			if !owner {
				continue
			}

			// transfer now:
			sentTokenResp, err := u.SendTokenMKP(item.OrdAddress, item.InscriptionID)

			go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "SendTokenMKP.sentTokenResp", sentTokenResp)

			if err != nil {
				u.Logger.Error(fmt.Sprintf("BtcSendNFTForBuyOrder.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
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
				u.Logger.Error("BtcSendNFTForBuyOrder.helpers.JsonTransform", errPack.Error(), errPack)
				go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "SendTokenMKP.UpdateBTCNFTBuyOrder", err.Error())
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
				u.Logger.Error("BtcSendNFTForBuyOrder.Repo.UpdateBTCNFTBuyOrder", errPack.Error(), errPack)
				go u.trackHistory(item.UUID, "BtcSendNFTForBuyOrder", item.TableName(), item.Status, "u.Repo.UpdateBTCNFTBuyOrder", err.Error())
			}
			// save log:
			u.Logger.Info(fmt.Sprintf("BtcSendNFTForBuyOrder.execResp.%s", item.OrdAddress), sentTokenResp)
		}
	}
	return nil
}

func (u Usecase) BtcCheckSendNFTForBuyOrder() error {

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
				go u.trackHistory(item.UUID, "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "btcClient.GetTransaction: "+item.TxSendNFT, err.Error())

				go u.trackHistory(item.UUID, "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendNFT, "Begin check tx via api.")

				// check with api:
				txInfo, err := bs.CheckTx(item.TxSendNFT)
				if err != nil {
					fmt.Printf("Could not bs - with err: %v", err)
					go u.trackHistory(item.UUID, "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx: "+item.TxSendNFT, err.Error())
				}
				if txInfo.Confirmations >= 1 {
					go u.trackHistory(item.UUID, "BtcCheckSendNFTForBuyOrder", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+item.TxSendNFT, txInfo.Confirmations)
					// send nft ok now:
					item.Status = entity.StatusBuy_SentNFT
					_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
					if err != nil {
						fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
					}

					// Add successfully buy activity.
					go func(item entity.MarketplaceBTCBuyOrder) {
						nftListing, err := u.Repo.FindBtcNFTListingByOrderID(item.ItemID)
						if err != nil {
							fmt.Println("can not FindBtcNFTListingByOrderID with err:", err)
							return
						}
						u.CreateBuyActivity(item.InscriptionID, nftListing.Price)
					}(item)
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

func (u Usecase) GetMasterNfts() (*ord_service.ExecRespose, error) {

	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_marketplace_master",
			"wallet",
			"inscriptions",
		}}

	u.Logger.Info("listNFTsReq", listNFTsReq)
	resp, err := u.OrdService.Exec(listNFTsReq)
	defer u.Notify("GetMasterNfts", "ord_marketplace_master", "inscriptions")
	if err != nil {
		u.Logger.Info("u.OrdService.Exec.Error", err.Error())
		u.Logger.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("listNFTsRep", resp)
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
func (u Usecase) SendTokenMKPTest(walletName, receiveAddr, inscriptionID string) (*ord_service.ExecRespose, error) {

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

	u.Logger.Info("sendTokenReq", sendTokenReq)

	resp, err := u.OrdService.Exec(sendTokenReq)

	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "", "after call OrdService.Exec")
	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "SendTokenMKP.JsonTransform", resp)

	defer u.Notify("SendTokenMKPTest", receiveAddr, inscriptionID)
	if err != nil {
		u.Logger.Info("u.OrdService.Exec.Error", err.Error())
		u.Logger.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("sendTokenRes", resp)

	go u.trackHistory("test_send_nft", "SendTokenMKPTest", "", 0, "", "return now...")

	return resp, err
}

// admin
// check receive of the nft:
func (u Usecase) AutoListing(reqs *request.ListNftIdsReq) interface{} {
	var listIdSuccess []string

	if reqs != nil {
		for _, v := range reqs.InscriptionID {
			//v.Inscription
			listing := entity.MarketplaceBTCListing{
				SellOrdAddress: reqs.SellOrdAddress,
				SellerAddress:  reqs.SellerAddress,
				HoldOrdAddress: "",
				ServiceFee:     "0",
				Price:          reqs.Price,
				IsConfirm:      true,
				IsSold:         false,
				ExpiredAt:      time.Now().Add(time.Hour * 1),
				Name:           "",
				Description:    "",
				InscriptionID:  v,
			}
			// get first:
			nftList, _ := u.Repo.FindBtcNFTListingByNFTID(v)
			if nftList != nil && nftList.IsConfirm && !nftList.IsSold {
				u.Logger.Error("AutoListing.Repo.FindBtcNFTListingByNFTID", "", errors.New("item exist"))
				continue
			}

			// check if listing is created or not
			err := u.Repo.CreateMarketplaceListingBTC(&listing)
			if err != nil {
				u.Logger.Error("AutoListing.Repo.CreateMarketplaceBTCListing", "", err)
				continue
			}
			listIdSuccess = append(listIdSuccess, v)
		}
	}

	return listIdSuccess
}
