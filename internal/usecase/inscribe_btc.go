package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
	"rederinghub.io/utils/contracts/ordinals"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/fileutil"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type BitcoinTokenMintFee struct {
	Amount       string
	MintFee      string
	SentTokenFee string
	Size         int
}

func decodeFileBase64(file string) (string, string, error) {
	i := strings.Index(file, ",")
	if i < 0 {
		return "", "", errors.New("no comma")
	}

	dec, err := base64.StdEncoding.DecodeString(file[i+1:])
	if err != nil {
		return "", "", err
	}
	return string(dec), file[i+1:], nil
}

func calculateMintPrice(input structure.InscribeBtcReceiveAddrRespReq) (*BitcoinTokenMintFee, error) {
	// base64String := input.File
	// base64String = strings.ReplaceAll(base64String, "data:text/html;base64,", "")
	// base64String = strings.ReplaceAll(base64String, "data:image/png;base64,", "")
	// dec, err := base64.StdEncoding.DecodeString(base64String)
	// if err != nil {
	// 	return nil, err
	// }

	// need to encode file: phuong viet lai:
	fileDecode, _, err := decodeFileBase64(input.File)
	if err != nil {
		return nil, err
	}

	fileSize := len([]byte(fileDecode))

	fmt.Println("fileSize===>", fileSize)
	fmt.Println("input.FeeRate===>", input.FeeRate)

	if fileSize < utils.MIN_FILE_SIZE {
		fileSize = utils.MIN_FILE_SIZE
	}
	fmt.Println("new fileSize===>", fileSize)

	fileSize += utils.MIN_FILE_SIZE // add 4kb

	mintFee := int32(fileSize) / 4 * input.FeeRate

	fmt.Println("mintFee===>", mintFee)

	sentTokenFee := utils.FEE_BTC_SEND_AGV + utils.FEE_BTC_SEND_NFT
	totalFee := int(mintFee) + sentTokenFee

	fmt.Println("total fee ===>", totalFee)

	return &BitcoinTokenMintFee{
		Amount:       strconv.FormatInt(int64(totalFee), 10),
		MintFee:      strconv.FormatInt(int64(mintFee), 10),
		SentTokenFee: strconv.FormatInt(int64(sentTokenFee), 10),
		Size:         fileSize,
	}, nil
}

func (u Usecase) CreateInscribeBTC(ctx context.Context, input structure.InscribeBtcReceiveAddrRespReq) (*entity.InscribeBTC, error) {

	u.Logger.LogAny("CreateInscribeBTC", zap.Any("input", input))

	// todo remove:
	// _, base64Str, err := decodeFileBase64(input.File)
	// if err != nil {
	// 	u.Logger.Error("JobInscribeMintNft.decodeFileBase64", err.Error(), err)
	// 	return nil, err
	// }

	// now := time.Now().UTC().Unix()
	// uploaded, err := u.GCS.UploadBaseToBucket(base64Str, fmt.Sprintf("btc-projects/%s/%d.%s", "bc1p3lh2xp8a63rlwpk8zkxrwhhzwqgskfr9el3lmhceu3atyam4rvmshf24vt", now, "txt"))
	// if err != nil {
	// 	return nil, err
	// }

	// fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	// fmt.Println("fileURI===> ", fileURI)

	// end remove

	if input.FeeRate <= 3 {
		err := errors.New("fee rate must be > 3")
		u.Logger.ErrorAny("u.CreateInscribeBTC.Copy", zap.Error(err))
		return nil, err
	}

	walletAddress := &entity.InscribeBTC{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		u.Logger.ErrorAny("u.CreateInscribeBTC.Copy", zap.Error(err))
		return nil, err
	}

	// need function get size only:
	mintFee, err := calculateMintPrice(input)
	if err != nil {
		u.Logger.ErrorAny("u.CreateSegwitBTCWalletAddress.calculateMintPrice", zap.Error(err))
		return nil, err
	}

	var mfTotal, mfMintFee, mfSentTokenFee string

	// cal fee again:
	feeInfos, err := u.calMintFeeInfo(0, int64(mintFee.Size), int64(input.FeeRate), 0, 0)
	if err != nil {
		u.Logger.Error("u.calMintFeeInfo.Err", err.Error(), err)
		return nil, err
	}

	mfTotal = big.NewInt(0).Add(feeInfos[input.PayType].MintFeeBigInt, feeInfos[input.PayType].SendNftFeeBigInt).String()

	fmt.Println("mfTotal eth 0==>", mfTotal)

	mfMintFee = feeInfos[input.PayType].MintFee
	mfSentTokenFee = feeInfos[input.PayType].SendNftFee

	if input.PayType == utils.NETWORK_ETH {

		mfTotal = feeInfos[input.PayType].NetworkFee
		fmt.Println("mfTotal eth 1===>", mfTotal)

		mfMintFee = feeInfos[input.PayType].MintFee
		mfSentTokenFee = big.NewInt(0).Add(feeInfos[input.PayType].SendNftFeeBigInt, feeInfos[input.PayType].SendFundFeeBigInt).String()
	}

	privKey := ""
	addressSegwit := ""
	payType := input.PayType

	fmt.Println("payType: ", payType)

	if len(payType) == 0 {
		payType = utils.NETWORK_BTC
	}

	if strings.ToLower(payType) == strings.ToLower(utils.NETWORK_ETH) {

		ethClient := eth.NewClient(nil)

		// create segwit address
		privKey, _, addressSegwit, err = ethClient.GenerateAddress()
		if err != nil {
			u.Logger.ErrorAny("CreateInscribeBTC.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}

		privKey = strings.ToLower(privKey)
		addressSegwit = strings.ToLower(addressSegwit)

	} else {
		// just create ord wallet for btc payment:
		// create wallet name
		userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)

		// create master wallet:
		resp, err := u.OrdService.Exec(ord_service.ExecRequest{
			Args: []string{
				"--wallet",
				userWallet,
				"wallet",
				"create",
			},
		})

		if err != nil {
			u.Logger.ErrorAny("u.OrdService.Exec.create.Wallet", zap.Error(err))
			return nil, err
		}
		walletAddress.Mnemonic = resp.Stdout

		u.Logger.Info("CreateOrdBTCWalletAddress.createdWallet", resp)
		resp, err = u.OrdService.Exec(ord_service.ExecRequest{
			Args: []string{
				"--wallet",
				userWallet,
				"wallet",
				"receive",
			},
		})

		if err != nil {
			u.Logger.ErrorAny("u.OrdService.Exec.create.receive", zap.Error(err))
			return nil, err
		}
		u.Logger.Info("CreateInscribeBTC.calculateMintPrice", resp)

		// parse json to get address:
		// ex: {"mnemonic": "chaos dawn between remember raw credit pluck acquire satoshi rain one valley","passphrase": ""}

		jsonStr := strings.ReplaceAll(resp.Stdout, "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

		var receiveResp ord_service.ReceiveCmdStdoputRespose

		err = json.Unmarshal([]byte(jsonStr), &receiveResp)
		if err != nil {
			u.Logger.ErrorAny("CreateInscribeBTC.Unmarshal", zap.Error(err))
			return nil, err
		}

		walletAddress.UserAddress = userWallet // name
		walletAddress.OrdAddress = receiveResp.Address

		// create segwit address
		privKey, _, addressSegwit, err = btc.GenerateAddressSegwit()
		if err != nil {
			u.Logger.ErrorAny("CreateInscribeBTC.GenerateAddressSegwit", zap.Error(err))
			return nil, err
		}

	}

	if privKey == "" {
		err := errors.New("Cannot create privKey")
		u.Logger.ErrorAny("CreateInscribeBTC.privKey", zap.Error(err))
		return nil, err
	}

	if addressSegwit == "" {
		err := errors.New("Cannot create addressSegwit")
		u.Logger.ErrorAny("CreateInscribeBTC.addressSegwit", zap.Error(err))
		return nil, err
	}

	walletAddress.SegwitKey = privKey
	walletAddress.SegwitAddress = addressSegwit
	walletAddress.PayType = payType

	expiredTime := utils.INSCRIBE_TIMEOUT
	if u.Config.ENV == "develop" {
		expiredTime = 1
	}

	walletAddress.Amount = mfTotal
	walletAddress.MintFee = mfMintFee
	walletAddress.SentTokenFee = mfSentTokenFee
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = input.File
	walletAddress.InscriptionID = ""
	walletAddress.FeeRate = input.FeeRate
	walletAddress.ExpiredAt = time.Now().Add(time.Hour * time.Duration(expiredTime))
	walletAddress.FileName = input.FileName
	walletAddress.UserUuid = input.UserUuid
	walletAddress.UserWalletAddress = input.UserWallerAddress
	walletAddress.BTCRate = feeInfos[payType].BtcPrice
	walletAddress.ETHRate = feeInfos[payType].EthPrice
	walletAddress.EstFeeInfo = feeInfos

	if input.NeedVerifyAuthentic() {
		// inscribeBtc := &entity.InscribeBTC{}
		// opt := &options.FindOneOptions{}
		// opt.SetSort(bson.M{"_id": -1})
		// err := u.Repo.FindOneBy(ctx,
		// 	inscribeBtc.TableName(),
		// 	bson.M{
		// 		"user_uuid":     input.UserUuid,
		// 		"token_address": input.TokenAddress,
		// 		"token_id":      input.TokenId,
		// 	},
		// 	inscribeBtc,
		// 	opt)
		// if err != nil {
		// 	if !errors.Is(err, mongo.ErrNoDocuments) {
		// 		return nil, err
		// 	}
		// } else {
		// 	if inscribeBtc.Status == entity.StatusInscribe_Pending {
		// 		if !inscribeBtc.Expired() {

		// 			return inscribeBtc, nil
		// 		}
		// 	} else if inscribeBtc.Status != entity.StatusInscribe_TxMintFailed {
		// 		return inscribeBtc, nil
		// 	}
		// }
		if nft, err := u.MoralisNft.GetNftByContractAndTokenID(input.TokenAddress, input.TokenId); err == nil {
			logger.AtLog.Logger.Info("MoralisNft.GetNftByContractAndTokenID",
				zap.Any("raw_data", nft))
			walletAddress.IsAuthentic = true
			walletAddress.TokenAddress = nft.TokenAddress
			walletAddress.TokenId = nft.TokenID
			walletAddress.OwnerOf = nft.Owner
		}
	}

	fmt.Println("walletAddress.Amount===>", walletAddress.Amount)

	err = u.Repo.InsertInscribeBTC(walletAddress)
	if err != nil {
		u.Logger.ErrorAny("u.CreateInscribeBTC.InsertInscribeBTC", zap.Error(err))
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ListInscribeBTC(req *entity.FilterInscribeBT) (*entity.Pagination, error) {
	return u.Repo.ListInscribeBTC(req)
}

func (u Usecase) DetailInscribeBTC(inscriptionID string) (*entity.InscribeBTCResp, error) {
	return u.Repo.FindInscribeBTCByNftID(inscriptionID)
}

func (u Usecase) RetryInscribeBTC(id string) error {
	item, _ := u.Repo.FindInscribeBTC(id)
	u.Logger.Info("item: ", item, id)
	if item != nil {
		if item.Status == entity.StatusInscribe_NotEnoughBalance {
			item.Status = entity.StatusInscribe_Pending
			_, err := u.Repo.UpdateBtcInscribe(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// JOBs:
// step 1: job check balance for list inscribe
func (u Usecase) JobInscribeWaitingBalance() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackInscribeHistory("", "JobInscribeWaitingBalance", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}
	listPending, _ := u.Repo.ListBTCInscribePending()
	if len(listPending) == 0 {
		// go u.trackInscribeHistory("", "JobInscribeWaitingBalance", "", "", "ListBTCInscribePending", "[]")
		return nil
	}
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		go u.trackInscribeHistory("", "JobInscribeWaitingBalance", "", "", "Could not initialize Ether RPCClient - with err", err.Error())
		return err
	}
	ethClient := eth.NewClient(ethClientWrap)

	// get list btc to check a Batch
	var batchBTCBalance []string
	for _, item := range listPending {
		if item.PayType == utils.NETWORK_BTC {
			batchBTCBalance = append(batchBTCBalance, item.SegwitAddress)
		}
	}

	isRateLimitErr := false
	balanceMaps, err := bs.BTCGetAddrInfoMulti(batchBTCBalance)
	if err != nil && strings.Contains(err.Error(), "rate_limit") {
		isRateLimitErr = true
	}

	for _, item := range listPending {

		// check balance:
		balance := big.NewInt(0)
		confirm := -1

		if item.PayType == utils.NETWORK_BTC {

			// check balance:
			// balance, confirm, err = bs.GetBalance(item.SegwitAddress)
			// fmt.Println("GetBalance response: ", balance, confirm, err)
			if !isRateLimitErr {
				balanceInfo, ok := balanceMaps[item.SegwitAddress]
				// If the key exists
				if ok {
					balance = big.NewInt(0).SetUint64(balanceInfo.Balance)
					if len(balanceInfo.TxRefs) > 0 {
						confirm = balanceInfo.TxRefs[0].Confirmations
					}
				}
			} else if isRateLimitErr {
				// get balance from quicknode:
				var balanceQuickNode *structure.BlockCypherWalletInfo
				balanceQuickNode, err = btc.GetBalanceFromQuickNode(item.SegwitAddress, u.Config.QuicknodeAPI)
				if err == nil {
					if balanceQuickNode != nil {
						balance = big.NewInt(int64(balanceQuickNode.Balance))
						// check confirm:
						if len(balanceQuickNode.Txrefs) > 0 {
							var txInfo *btc.QuickNodeTx
							txInfo, err = btc.CheckTxfromQuickNode(balanceQuickNode.Txrefs[0].TxHash, u.Config.QuicknodeAPI)
							if err == nil {
								if txInfo != nil {
									confirm = txInfo.Result.Confirmations
								}

							} else {
								go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "CheckTxfromQuickNode from quicknode - with err", err.Error())
							}
						}
					}

				} else {
					go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance from quicknode - with err", err.Error())
				}
			}

		} else if item.PayType == utils.NETWORK_ETH {

			// check eth balance:

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			balance, err = ethClient.GetBalance(ctx, item.SegwitAddress)
			fmt.Println("GetBalance eth response: ", balance, err)

			confirm = 1
		}

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance - with err", err.Error())
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", err.Error())
			continue
		}

		if balance.Uint64() == 0 {
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "GetBalance", "0")
			continue
		}

		// get required amount to check vs temp wallet balance:
		amount, ok := big.NewInt(0).SetString(item.Amount, 10)
		if !ok {
			err := errors.New("cannot parse amount")
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "SetString(amount) err", err.Error())
			continue
		}

		if amount.Uint64() == 0 {
			err := errors.New("balance is zero")
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "amount.Uint64() err", err.Error())
			continue
		}

		if balance.Uint64() < amount.Uint64() {
			err := fmt.Errorf("Not enough amount %d < %d ", balance.Uint64(), amount.Uint64())
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "compare balance err", err.Error())

			item.Status = entity.StatusInscribe_NotEnoughBalance
			u.Repo.UpdateBtcInscribe(&item)
			continue
		}

		if confirm <= 0 {
			go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "Updated StatusMint_WaitingForConfirms", confirm)
			continue
		}

		// received fund:
		item.Status = entity.StatusInscribe_ReceivedFund
		if item.PayType == utils.NETWORK_ETH {
			item.Status = entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd // next step to mint (ready to mint)
			item.IsMergeMint = true                                           // jus for eth now
		}

		item.IsConfirm = true

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			continue
		}

		go u.trackInscribeHistory(item.UUID, "JobInscribeWaitingBalance", item.TableName(), item.Status, "Updated StatusInscribe_ReceivedFund", "ok")
		u.Logger.Info(fmt.Sprintf("JobInscribeWaitingBalance.CheckReceiveBTC.%s", item.SegwitAddress), item)
		u.Notify("JobInscribeWaitingBalance", item.SegwitAddress, fmt.Sprintf("%s received BTC %d from [InscriptionID] %s", item.SegwitAddress, item.Status, item.InscriptionID))

	}

	return nil
}

// step 2: job send all fund from segwit address to ord wallet:
func (u Usecase) JobInscribeSendBTCToOrdWallet() error {

	_, bs, err := u.buildBTCClient()

	if err != nil {
		go u.trackInscribeHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "Could not initialize Bitcoin RPCClient - with err", err.Error())
		return err
	}

	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_ReceivedFund})
	if len(listTosendBtc) == 0 {
		// go u.trackInscribeHistory("", "JobInscribeSendBTCToOrdWallet", "", "", "ListBTCInscribeByStatus", "[]")
		return nil
	}

	for _, item := range listTosendBtc {
		if item.Status == entity.StatusInscribe_ReceivedFund {

			if item.PayType == utils.NETWORK_ETH {
				continue
			}

			// send all amount:
			fmt.Println("send all btc from", item.SegwitAddress, "to: ", item.OrdAddress)

			// transfer now:
			txID, err := bs.SendTransactionWithPreferenceFromSegwitAddress(
				item.SegwitKey,
				item.SegwitAddress,
				item.OrdAddress,
				-1,
				btc.PreferenceHigh,
			)
			if err != nil {
				go u.trackInscribeHistory(item.UUID, "JobInscribeSendBTCToOrdWallet", item.TableName(), item.Status, "SendTransactionWithPreferenceFromSegwitAddress err", err.Error())
				continue
			}

			item.TxSendBTC = txID
			item.Status = entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr
			// item.ErrCount = 0 // reset error count!
			// TODO: update item
			_, err = u.Repo.UpdateBtcInscribe(&item)
			if err != nil {
				fmt.Printf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			}

		}
	}

	return nil
}

// job check 3 tx send: tx user send to temp wallet, tx mint, tx send nft to user
func (u Usecase) JobInscribeCheckTxSend() error {
	btcClient, bs, err := u.buildBTCClient()
	if err != nil {
		logger.AtLog.Logger.Error("Could not initialize Bitcoin RPCClient failed", zap.Error(err))
		return err
	}
	ordinalsSrv, _ := ordinals.NewService(
		u.Config.Ordinals.OrdinalsContract,
		u.Config.Ordinals.CallerOrdinalsPrivateKey,
		int64(u.Config.ChainId),
	)

	// get list sending tx:
	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_Minting, entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr, entity.StatusInscribe_SendingNFTToUser})
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		fields := []zapcore.Field{
			zap.String("id", item.ID.Hex()),
			zap.String("file_name", item.FileName),
		}

		statusSuccess := entity.StatusInscribe_Minted
		txHashDb := item.TxMintNft

		if item.Status == entity.StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr {
			statusSuccess = entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd
			txHashDb = item.TxSendBTC
		}
		if item.Status == entity.StatusInscribe_SendingNFTToUser {
			statusSuccess = entity.StatusInscribe_SentNFTToUser
			txHashDb = item.TxSendNft
		}
		if item.Status == entity.StatusInscribe_Minting {
			item.IsMinted = true
		}

		txHash, err := chainhash.NewHashFromStr(txHashDb)
		if err != nil {
			logger.AtLog.Logger.With(fields...).Error("Could not NewHashFromStr Bitcoin RPCClient ")
			continue
		}
		txConfirmation := false
		txResponse, err := btcClient.GetTransaction(txHash)
		if err == nil {
			go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.txResponse.Confirmations: "+txHashDb, txResponse.Confirmations)
			if txResponse.Confirmations >= 1 {
				txConfirmation = true
				// send btc ok now:
				item.Status = statusSuccess
				_, err = u.Repo.UpdateBtcInscribe(&item)
				if err != nil {
					logger.AtLog.Logger.With(fields...).Error("Could not JobInscribeCheckTxSend")
				}
			}
		} else {
			logger.AtLog.Logger.With(fields...).Error("Could not GetTransaction Bitcoin RPCClient")
			go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "btcClient.GetTransaction: "+txHashDb, err.Error())

			go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, "Begin check tx via api.")

			// check with api:
			txInfo, err := bs.CheckTx(txHashDb)
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.AtLog.Logger.With(fields...).Error("Could not CheckTx")
				go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx: "+txHashDb, err.Error())
			}

			if txInfo.Confirmations >= 1 {
				txConfirmation = true
				go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "bs.CheckTx.txInfo.Confirmations: "+txHashDb, txInfo.Confirmations)
				// send nft ok now:
				item.Status = statusSuccess
				item.IsSuccess = statusSuccess == entity.StatusInscribe_SentNFTToUser
				_, err = u.Repo.UpdateBtcInscribe(&item)
				if err != nil {
					fields = append(fields, zap.Error(err))
					logger.AtLog.Logger.With(fields...).Error("Could not UpdateBtcInscribe")
				}
				/* remove this feature
				if item.Status == entity.StatusInscribe_SentNFTToUser {
					go func(u Usecase, item entity.InscribeBTC) {
						owner, err := u.Repo.FindUserByBtcAddressTaproot(item.OriginUserAddress)
						if err != nil || owner == nil {
							return
						}
						u.AirdropCollector("0000000", item.InscriptionID, os.Getenv("AIRDROP_WALLET"), *owner, 3)
					}(u, item)
				}*/
			}
		}

		// add contract
		if ordinalsSrv != nil && txConfirmation && item.NeedAddContractToOrdinalsContract() {
			err = u.AddContractToOrdinalsContract(context.Background(), ordinalsSrv, item)
			if err != nil {
				go u.trackInscribeHistory(item.UUID, "JobInscribeCheckTxSend", item.TableName(), item.Status, "JobInscribeCheckTxSend.AddContractToOrdinalsContract", err.Error())
				fields = append(fields, zap.Error(err))
				logger.AtLog.Logger.With(fields...).Error("AddContractToOrdinalsContract failed")
				continue
			}
		}
	}

	return nil
}

// job 4: mint nft:
func (u Usecase) JobInscribeMintNft() error {
	listTosendBtc, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_SentBTCFromSegwitAddrToOrdAdd})
	if len(listTosendBtc) == 0 {
		return nil
	}

	for _, item := range listTosendBtc {
		fields := []zapcore.Field{
			zap.String("id", item.ID.Hex()),
			zap.String("file_name", item.FileName),
		}

		logger.AtLog.Logger.With(fields...).Info("Mint nft now...")

		// - Upload the Animation URL to GCS
		typeFile := ""

		if len(item.FileName) == 0 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobInscribeMintNft.len(Filename)", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFiles := strings.Split(item.FileName, ".")
		if len(typeFiles) < 2 {
			err := errors.New("File name invalid")
			u.Logger.Error("JobInscribeMintNft.len(Filename)", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "CheckFileName", err.Error())
			continue
		}

		typeFile = typeFiles[len(typeFiles)-1]
		fields = append(fields, zap.String("type_file", typeFile))
		logger.AtLog.Logger.Info("TypeFile", fields...)

		// update google clound: TODO need to move into api to avoid create file many time.
		_, base64Str, err := decodeFileBase64(item.FileURI)
		if err != nil {
			u.Logger.Error("JobInscribeMintNft.decodeFileBase64", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "helpers.decodeFileBase64", err.Error())
			continue
		}

		now := time.Now().UTC().Unix()
		uploaded, err := u.GCS.UploadBaseToBucket(base64Str, fmt.Sprintf("btc-projects/%s/%d.%s", item.SegwitAddress, now, typeFile))
		if err != nil {
			u.Logger.Error("JobInscribeMintNft.helpers.UploadBaseToBucket.Base64DecodeRaw", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "helpers.BUploadBaseToBucket.ase64DecodeRaw", err.Error())
			continue
		}
		item.LocalLink = uploaded.FullPath

		fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
		item.FileURI = fileURI

		ordWalletNameToMint := item.UserAddress

		if item.PayType == utils.NETWORK_ETH {
			ordWalletNameToMint = "ord_master_eth" // TODO: move to env
		}

		//TODO - enable this
		mintData := ord_service.MintRequest{
			WalletName:        ordWalletNameToMint,
			FileUrl:           fileURI,
			FeeRate:           int(item.FeeRate),
			DryRun:            false,
			AutoFeeRateSelect: false,

			RequestId: item.UUID, // for tracking log

			// new key for ord v5.1, support mint + send in 1 tx:
			DestinationAddress: item.OriginUserAddress, // the address mint to.
		}
		resp, err := u.OrdService.Mint(mintData)

		if err != nil {
			u.Logger.Error("OrdService.Mint", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, mintData, err.Error())
			continue
		}
		// if not err => update status ok now:
		//TODO: handle log err: Database already open. Cannot acquire lock

		item.Status = entity.StatusInscribe_Minting
		item.IsMergeMint = true // used ord 5.1, mint+send in 1 tx
		// item.ErrCount = 0 // reset error count!

		item.OutputMintNFT = resp

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.UpdateBtcInscribe", err.Error())
			continue
		}

		tmpText := resp.Stdout
		//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
		jsonStr := strings.ReplaceAll(tmpText, "\n", "")
		jsonStr = strings.ReplaceAll(jsonStr, "\\", "")

		var btcMintResp ord_service.MintStdoputRespose

		err = json.Unmarshal([]byte(jsonStr), &btcMintResp)
		if err != nil {
			u.Logger.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.Unmarshal(btcMintResp)", err.Error())
			continue
		}

		item.TxMintNft = btcMintResp.Reveal
		item.InscriptionID = btcMintResp.Inscription
		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			fields = append(fields, zap.Error(err))
			logger.AtLog.Logger.With(fields...).Error("Could not UpdateBtcInscribe")
			go u.trackInscribeHistory(item.UUID, "JobInscribeMintNft", item.TableName(), item.Status, "JobInscribeMintNft.UpdateBtcInscribe", err.Error())
		}
	}

	return nil
}

// job 5: send nft:
// send nft for buy order records:
func (u Usecase) JobInscribeSendNft() error {

	// get list buy order status = StatusInscribe_Minted:
	listTosendNft, _ := u.Repo.ListBTCInscribeByStatus([]entity.StatusInscribe{entity.StatusInscribe_Minted})
	if len(listTosendNft) == 0 {
		return nil
	}

	for _, item := range listTosendNft {

		// update for ord v5.1: is merged tx
		if item.IsMergeMint {
			// don't send, update isSent = true
			item.Status = entity.StatusInscribe_SentNFTToUser
			item.IsSuccess = true
			u.Repo.UpdateBtcInscribe(&item)
			continue

		}

		// check nft in master wallet or not:
		listNFTsRep, err := u.GetNftsOwnerOf(item.UserAddress)
		if err != nil {
			go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Error", err.Error())
			continue
		}

		go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.listNFTsRep", listNFTsRep)

		// parse nft data:
		var resp []struct {
			Inscription string `json:"inscription"`
			Location    string `json:"location"`
			Explorer    string `json:"explorer"`
		}

		err = json.Unmarshal([]byte(listNFTsRep.Stdout), &resp)
		if err != nil {
			go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.Unmarshal(listNFTsRep)", err.Error())
			continue
		}
		owner := false
		for _, nft := range resp {
			if strings.EqualFold(nft.Inscription, item.InscriptionID) {
				owner = true
				break
			}

		}
		go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "GetNftsOwnerOf.CheckNFTOwner", owner)
		if !owner {
			continue
		}

		// transfer now:
		sentTokenResp, err := u.SendTokenByWallet(item.OriginUserAddress, item.InscriptionID, item.UserAddress, int(item.FeeRate))

		go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "SendTokenByWallet.sentTokenResp", sentTokenResp)

		if err != nil {
			u.Logger.Error(fmt.Sprintf("JobInscribeSendNft.SendTokenMKP.%s.Error", item.OrdAddress), err.Error(), err)
			go u.trackInscribeHistory(item.UUID, "JobInscribeSendNft", item.TableName(), item.Status, "SendTokenByWallet.err", err.Error())
			continue
		}

		//TODO: handle log err: Database already open. Cannot acquire lock

		// Update status first if none err:
		item.Status = entity.StatusInscribe_SendingNFTToUser
		// item.ErrCount = 0 // reset error count!

		item.OutputSendNFT = sentTokenResp

		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err.Error())
			u.Logger.Error("JobMKP_SendNftToBuyer.helpers.JsonTransform", errPack.Error(), errPack)
			go u.trackInscribeHistory(item.UUID, "UpdateBtcInscribe", item.TableName(), item.Status, "SendTokenMKP.UpdateBtcInscribe", err.Error())
			continue
		}

		txResp := sentTokenResp.Stdout
		//txResp := `fd31946b855cbaaf91df4b2c432f9b173e053e65a9879ac909bad028e21b950e\n`
		txResp = strings.ReplaceAll(txResp, "\n", "")

		// update tx:
		item.TxSendNft = txResp
		// item.ErrCount = 0 // reset error count!
		_, err = u.Repo.UpdateBtcInscribe(&item)
		if err != nil {
			errPack := fmt.Errorf("Could not UpdateBtcInscribe id %s - with err: %v", item.ID, err)
			u.Logger.Error("UpdateBtcInscribe.UpdateBtcInscribe", errPack.Error(), errPack)
			go u.trackInscribeHistory(item.UUID, "UpdateBtcInscribe", item.TableName(), item.Status, "u.UpdateBtcInscribe.UpdateBTCNFTBuyOrder", err.Error())
		}
		// save log:
		u.Logger.Info(fmt.Sprintf("UpdateBtcInscribe.execResp.%s", item.OrdAddress), sentTokenResp)

	}
	return nil
}

func (u Usecase) SendTokenByWallet(receiveAddr, inscriptionID, walletAddressName string, rate int) (*ord_service.ExecRespose, error) {

	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletAddressName,
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			fmt.Sprintf("%d", rate),
		}}

	resp, err := u.OrdService.Exec(sendTokenReq)
	return resp, err
}

func (u Usecase) GetNftsOwnerOf(walletName string) (*ord_service.ExecRespose, error) {

	listNFTsReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			walletName,
			"wallet",
			"inscriptions",
		}}

	u.Logger.Info("listNFTsReq", listNFTsReq)
	resp, err := u.OrdService.Exec(listNFTsReq)
	defer u.Notify("GetNftsOwnerOf", "ord_marketplace_master", "inscriptions")
	if err != nil {
		u.Logger.Info("u.OrdService.Exec.Error", err.Error())
		u.Logger.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}
	u.Logger.Info("listNFTsRep", resp)
	return resp, err
}

func (u *Usecase) trackInscribeHistory(id, name, table string, status interface{}, requestMsg interface{}, responseMsg interface{}) {
	trackData := &entity.InscribeBTCLogs{
		RecordID:    id,
		Name:        name,
		Table:       table,
		Status:      status,
		RequestMsg:  requestMsg,
		ResponseMsg: responseMsg,
	}
	err := u.Repo.CreateInscribeBTCLog(trackData)
	if err != nil {
		fmt.Printf("trackInscribeHistory.%s.Error:%s", name, err.Error())
	}

}

func (u Usecase) ApiCheckListTempAddress() error {
	var autoGenerated []struct {
		SegwitAddress string `json:"segwit_address"`
	}
	listBtc := `[{}]`

	err := json.Unmarshal([]byte(listBtc), &autoGenerated)
	if err != nil {
		fmt.Println("err")
		return nil
	}

	_, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}
	fmt.Println("len(autoGenerated)", len((autoGenerated)))

	for _, btc := range autoGenerated {

		fmt.Println("check address: ", btc.SegwitAddress)

		balance, confirm, err := bs.GetBalance(btc.SegwitAddress)

		fmt.Println("GetBalance response: ", balance, confirm, err)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}
		if balance.Uint64() > 0 {
			fmt.Println("Balance OK now====>", btc.SegwitAddress)
		}
		time.Sleep(time.Second * 1)

	}

	return nil
}

func (u Usecase) ListNftFromMoralis(ctx context.Context, userId, userWallet, delegateWallet string, pag *entity.Pagination) (map[string]*entity.Pagination, error) {
	var (
		pageSize              = int(pag.PageSize)
		cursor        *string = nil
		resp                  = make(map[string]*entity.Pagination)
		walletAddress string
	)
	if len(pag.Cursor) > 0 {
		cursor = &pag.Cursor
	}
	reqMoralisFilter := nfts.MoralisFilter{
		Limit:  &pageSize,
		Cursor: cursor,
	}
	var (
		pageListInscribe  = int64(1)
		limitListInscribe = int64(100)
	)
	mapNftMinted := make(map[string]entity.InscribeBTCResp)
	for {
		resp, err := u.Repo.ListInscribeBTC(&entity.FilterInscribeBT{
			BaseFilters: entity.BaseFilters{
				Page:  pageListInscribe,
				Limit: limitListInscribe,
			},
			NeStatuses: []entity.StatusInscribe{entity.StatusInscribe_TxMintFailed},
			UserUuid:   &userId,
		})
		if err != nil {
			return nil, err
		}
		inscribes := resp.Result.([]entity.InscribeBTCResp)
		if len(inscribes) <= 0 {
			break
		}
		for _, inscribe := range inscribes {
			if inscribe.TokenAddress == "" || inscribe.TokenId == "" {
				continue
			}
			if inscribe.Status == entity.StatusInscribe_TxMintFailed ||
				inscribe.Status == entity.StatusInscribe_Pending {
				continue
			}
			if inscribe.InscriptionID != "" {
				tokenUri := &entity.TokenUri{}
				if err := u.Repo.FindOneBy(ctx, tokenUri.TableName(), bson.M{"token_id": inscribe.InscriptionID}, tokenUri); err == nil {
					inscribe.ProjectTokenId = tokenUri.ProjectID
				}
			}
			mapNftMinted[fmt.Sprintf("%s_%s", inscribe.TokenAddress, inscribe.TokenId)] = inscribe
		}
		if len(inscribes) < int(limitListInscribe) {
			break
		}
		pageListInscribe += 1
	}

	filterMoralisTokens := func(datas []nfts.MoralisToken) []nfts.MoralisToken {
		results := make([]nfts.MoralisToken, 0, len(datas))
		for _, data := range datas {
			if data.IsERC1155Type() {
				continue
			}
			if val, ok := mapNftMinted[fmt.Sprintf("%s_%s", data.TokenAddress, data.TokenID)]; ok {
				data.IsMinted = true
				data.InscribeBTC = &nfts.InscribeBTC{
					Status:         val.Status.Ordinal(),
					ProjectTokenId: val.ProjectTokenId,
					InscriptionID:  val.InscriptionID,
				}
			}
			results = append(results, data)
		}
		return results
	}

	if delegateWallet == "" {
		delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userWallet)
		if err != nil {
			return nil, err
		}
		if len(delegations) > 0 {
			for i := range delegations {
				delegateWalletAddress := delegations[i].Vault.String()
				resp[delegateWalletAddress] = &entity.Pagination{
					PageSize: int64(*reqMoralisFilter.Limit),
				}
				nft, err := u.MoralisNft.GetNftByWalletAddress(delegateWalletAddress, reqMoralisFilter)
				if err != nil {
					return nil, err
				}
				resp[delegateWalletAddress].Result = filterMoralisTokens(nft.Result)
				resp[delegateWalletAddress].Cursor = nft.Cursor
			}
		} else {
			walletAddress = userWallet
		}
	} else {
		walletAddress = delegateWallet
	}

	if walletAddress != "" {
		resp[walletAddress] = pag
		nft, err := u.MoralisNft.GetNftByWalletAddress(walletAddress, reqMoralisFilter)
		if err != nil {
			return nil, err
		}
		pag.Result = filterMoralisTokens(nft.Result)
		pag.Cursor = nft.Cursor
	}

	return resp, nil
}

func (u Usecase) CompressNftImageFromMoralis(ctx context.Context, urlStr string, compressPercents []int) (interface{}, error) {

	type CompressInfo struct {
		ImageUrl        string `json:"imageUrl"`
		CompressPercent int    `json:"compressPercent"`
		FileSize        int    `json:"fileSize"`
	}

	var compressInfoArray []*CompressInfo

	if strings.HasPrefix(urlStr, "http") {

		compressAndUploadImage := func(urlStr string, quality int) (*CompressInfo, error) {
			client := http.Client{}
			r, err := client.Get(urlStr)
			if err != nil {
				return nil, err
			}
			if r.StatusCode > http.StatusNoContent {
				return nil, err
			}
			defer r.Body.Close()
			imgByte, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}

			byteSize := len(imgByte)
			// if byteSize > fileutil.MaxImageByteSize || quality != -1 {
			if quality != -1 {

				// ext, err := utils.GetFileExtensionFromUrl(urlStr)
				// if err != nil {
				// 	contentType := r.Header.Get("content-type")
				// 	arr := strings.Split(contentType, "/")
				// 	if len(arr) <= 1 {
				// 		return "", err
				// 	}
				// 	ext = arr[1]
				// }

				// newImgByte, err := fileutil.ResizeImage(imgByte, ext, fileutil.MaxImageByteSize)
				// if err == nil {
				// 	imgByte = newImgByte
				// }
				linkImage, err := fileutil.ImageCompress(urlStr, quality, u.Config.THUMBOR_SECRET_KEY)
				if err != nil {
					return nil, err
				}

				rsp, err := client.Get(linkImage)
				if err != nil {
					return nil, err
				}
				if rsp.StatusCode > http.StatusNoContent {
					return nil, err
				}
				defer rsp.Body.Close()
				imgNewByte, err := io.ReadAll(rsp.Body)
				if err != nil {
					return nil, err
				}

				byteSizeNew := len(imgNewByte)

				return &CompressInfo{
					ImageUrl:        linkImage,
					CompressPercent: quality,
					FileSize:        byteSizeNew,
				}, nil

				// name := fmt.Sprintf("authentic/%v.%s", uuid.New().String(), ext)
				// _, err = u.GCS.UploadBaseToBucket(helpers.Base64Encode(newImgByte), name)
				// if err != nil {
				// 	return "", err
				// }

				// return fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name), nil
			}
			return &CompressInfo{
				ImageUrl:        urlStr,
				CompressPercent: quality,
				FileSize:        byteSize,
			}, nil
		}
		for _, compressPercent := range compressPercents {

			compressInfo, err := compressAndUploadImage(urlStr, compressPercent)
			if err != nil {
				return nil, err
			}
			compressInfoArray = append(compressInfoArray, compressInfo)

		}

	} else {
		return nil, errors.New("url invalid")
	}
	return compressInfoArray, nil
}

func (u Usecase) NftFromMoralis(ctx context.Context, tokenAddress, tokenId string) (*nfts.MoralisToken, error) {
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(tokenAddress, tokenId)
	if err != nil {
		return nil, err
	}
	metaData := &nfts.MoralisTokenMetadata{}
	if nft.MetadataString != nil {
		if err := json.Unmarshal([]byte(*nft.MetadataString), metaData); err != nil {
			return nil, err
		}
	}
	nft.Metadata = metaData
	if metaData.Image == "" {
		return nft, nil
	}
	urlStr := utils.ConvertIpfsToHttp(metaData.Image)
	if strings.HasPrefix(urlStr, "http") {
		resizeAndUploadImage := func(urlStr string) string {
			client := http.Client{}
			r, err := client.Get(urlStr)
			if err != nil {
				return urlStr
			}
			if r.StatusCode > http.StatusNoContent {
				return urlStr
			}
			defer r.Body.Close()
			imgByte, err := io.ReadAll(r.Body)
			if err != nil {
				return urlStr
			}
			ext, err := utils.GetFileExtensionFromUrl(urlStr)
			if err != nil {
				contentType := r.Header.Get("content-type")
				arr := strings.Split(contentType, "/")
				if len(arr) <= 1 {
					return urlStr
				}
				ext = arr[1]
			}
			// maybe use for thumb?
			// newImgByte, err := fileutil.ResizeImage(imgByte, ext, fileutil.MaxImageByteSize)
			// if err == nil {
			// 	imgByte = newImgByte
			// }
			name := fmt.Sprintf("authentic/%v.%s", uuid.New().String(), ext)
			_, err = u.GCS.UploadBaseToBucket(helpers.Base64Encode(imgByte), name)
			if err != nil {
				return urlStr
			}

			return fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
		}
		nft.Metadata.Image = resizeAndUploadImage(urlStr)

	} else if strings.HasPrefix(urlStr, ";base64,") {
		resizeImage := func(imageStr string) string {
			coI := strings.Index(imageStr, ",")
			dec, err := base64.StdEncoding.DecodeString(imageStr[coI+1:])
			if err != nil {
				return imageStr
			}
			exts := strings.Split(strings.TrimSuffix(imageStr[5:coI], ";base64"), "/")
			if len(exts) < 2 {
				return imageStr
			}
			imgByte := dec
			return imageStr[:coI+1] + base64.StdEncoding.EncodeToString(imgByte)
		}
		nft.Metadata.Image = resizeImage(urlStr)
	}
	return nft, nil
}

func (u Usecase) AddContractToOrdinalsContract(ctx context.Context, ordinalsSrv *ordinals.Service, item entity.InscribeBTC) error {
	txId, status, err := ordinalsSrv.AddContractToOrdinalsContract(ctx, item.TokenAddress, item.TokenId, item.InscriptionID)
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("AddContractToOrdinalsContract successfully",
		zap.String("id", item.ID.Hex()),
		zap.String("tx_id", txId),
	)
	item.OrdinalsTx = txId
	item.OrdinalsTxStatus = status
	_, err = u.Repo.UpdateBtcInscribe(&item)
	if err != nil {
		return err
	}
	if err := u.CreateProjectsAndTokenUriFromInscribeAuthentic(ctx, item); err != nil {
		return err
	}
	return nil
}
