package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"

	"rederinghub.io/external/nfts"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/eth"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateETHWalletAddress(rootSpan opentracing.Span, input structure.EthWalletAddressData) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("CreateETHWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag("btcUserWallet", input.WalletAddress)

	walletAddress := &entity.ETHWalletAddress{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.Copy", err.Error(), err)
		return nil, err
	}

	ethClient := eth.NewClient(nil)

	privKey, pubKey, address, err := ethClient.GenerateAddress()
	if err != nil {
		log.Error("ethClient.GenerateAddress", err.Error(), err)
		return nil, err
	} else {
		walletAddress.Mnemonic = privKey
	}

	log.SetData("CreateETHWalletAddress.createdWallet", fmt.Sprintf("%v %v %v", privKey, pubKey, address))

	log.SetData("CreateETHWalletAddress.receive", address)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}

	log.SetData("found.Project", p.ID)
	mintPriceInt, err := strconv.ParseInt(p.MintPrice, 10, 64)
	if err != nil {
		log.Error("convertBTCToInt", err.Error(), err)
		return nil, err
	}
	if p.NetworkFee != "" {
		// extra network fee
		networkFee, err1 := strconv.ParseInt(p.NetworkFee, 10, 64)
		if err1 != nil {
			log.Error("convertBTCToInt", err.Error(), err)
		} else {
			mintPriceInt += networkFee
		}
	}
	mintPrice, err := u.convertBTCToETH(span, fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
	if err != nil {
		log.Error("convertBTCToETH", err.Error(), err)
		return nil, err
	}
	walletAddress.Amount = mintPrice // 0.023 * 1e18 eth
	walletAddress.UserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(address, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.BalanceCheckTime = 0

	log.SetTag("ordAddress", walletAddress.OrdAddress)
	err = u.Repo.InsertEthWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.InsertEthWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) IsWhitelistedAddress(ctx context.Context, rootSpan opentracing.Span, userAddr string, whitelistedAddrs []string) (bool, error) {
	span, log := u.StartSpan("IsWhitelistedAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("whitelistedAddrs", whitelistedAddrs)
	if len(whitelistedAddrs) == 0 {
		log.SetData("whitelistedAddrs.Total", len(whitelistedAddrs))
		return false, nil
	}
	filter := nfts.MoralisFilter{}
	filter.Limit = new(int)
	*filter.Limit = 1
	filter.TokenAddresses = new([]string)
	*filter.TokenAddresses = whitelistedAddrs

	log.SetData("filter.GetNftByWalletAddress", filter)
	resp, err := u.MoralisNft.GetNftByWalletAddress(userAddr, filter)
	if err != nil {
		log.Error("u.MoralisNft.GetNftByWalletAddress", err.Error(), err)
		return false, err
	}

	log.SetData("resp", resp)
	if len(resp.Result) > 0 {
		return true, nil
	}

	delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userAddr)
	if err != nil {
		log.Error("u.DelegateService.GetDelegationsByDelegate", err.Error(), err)
		return false, err
	}

	log.SetData("delegations", delegations)
	for _, delegation := range delegations {
		if containsIgnoreCase(whitelistedAddrs, delegation.Contract.String()) {
			return true, nil
		}
	}
	return false, nil
}

func (u Usecase) CreateWhitelistedETHWalletAddress(ctx context.Context, rootSpan opentracing.Span, userAddr string, input structure.EthWalletAddressData) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("CreateETHWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag("btcUserWallet", input.WalletAddress)

	weth, err := u.Repo.FindDelegateEthWalletAddressByUserAddress(userAddr)
	if err == nil {
		if weth != nil {
			log.SetData("ethWalletAddress", weth)
			if weth.IsConfirm == true {
				err = errors.New("This account has applied discount")
				log.Error("applied.Discount", err.Error(), err)
				return nil, err
			}

			return weth, nil
		}
	} else {
		log.Error("u.Repo.FindEthWalletAddressByUserAddress", err.Error(), err)
	}

	walletAddress := &entity.ETHWalletAddress{}
	err = copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.Copy", err.Error(), err)
		return nil, err
	}

	// userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)
	// resp, err := u.OrdService.Exec(ord_service.ExecRequest{
	// 	Args: []string{
	// 		"--wallet",
	// 		userWallet,
	// 		"wallet",
	// 		"create",
	// 	},
	// })

	ethClient := eth.NewClient(nil)

	privKey, pubKey, address, err := ethClient.GenerateAddress()
	if err != nil {
		log.Error("ethClient.GenerateAddress", err.Error(), err)
		return nil, err
	} else {
		walletAddress.Mnemonic = privKey
	}

	log.SetData("CreateETHWalletAddress.createdWallet", fmt.Sprintf("%v %v %v", privKey, pubKey, address))

	log.SetData("CreateETHWalletAddress.receive", address)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}

	log.SetData("found.Project", p.ID)
	mintPriceInt, err := strconv.Atoi(p.MintPrice)
	if err != nil {
		log.Error("convertBTCToInt", err.Error(), err)
		return nil, err
	}
	if p.NetworkFee != "" {
		// extra network fee
		networkFee, err1 := strconv.Atoi(p.NetworkFee)
		if err1 != nil {
			log.Error("convertBTCToInt", err.Error(), err)
		} else {
			mintPriceInt += networkFee
		}
	}
	mintPrice, err := u.convertBTCToETH(span, fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
	if err != nil {
		log.Error("convertBTCToETH", err.Error(), err)
		return nil, err
	}

	walletAddress.Amount = mintPrice // 0.023 * 1e18 eth

	isWhitelist, err := u.IsWhitelistedAddress(ctx, span, userAddr, p.WhiteListEthContracts)

	if isWhitelist {
		whitelistedPrice := new(big.Float)
		ethPrice, _ := helpers.GetExternalPrice("ETH")
		if ethPrice == 0 {
			ethPrice = 1500
		}
		whitelistedPrice.SetFloat64(50.0 / ethPrice)
		whitelistedPrice = whitelistedPrice.Mul(whitelistedPrice, big.NewFloat(1e18))

		intPrice := new(big.Int)
		whitelistedPrice.Int(intPrice)

		walletAddress.Amount = intPrice.String()
	}

	walletAddress.UserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(address, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID
	walletAddress.Balance = "0"
	walletAddress.BalanceCheckTime = 0
	walletAddress.DelegatedAddress = userAddr

	log.SetTag("ordAddress", walletAddress.OrdAddress)
	err = u.Repo.InsertEthWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.InsertEthWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ETHMint(rootSpan opentracing.Span, input structure.BctMintData) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("ETHMint", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag("ordWalletaddress", input.Address)

	ethAddress, err := u.Repo.FindEthWalletAddressByOrd(input.Address)
	if err != nil {
		log.Error("ETHMint.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	//mint logic
	ethAddress, err = u.MintLogicETH(span, ethAddress)
	if err != nil {
		log.Error("ETHMint.MintLogic", err.Error(), err)
		return nil, err
	}

	// get data from project
	p, err := u.Repo.FindProjectByTokenID(ethAddress.ProjectID)
	if err != nil {
		log.Error("ETHMint.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	//log.SetData("found.Project", p)
	log.SetTag("projectID", p.TokenID)

	//prepare data for mint
	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		log.Error("ETHMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")

	now := time.Now().UTC().Unix()
	uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
	if err != nil {
		log.Error("ETHMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	ethAddress.FileURI = fileURI

	//TODO - enable this
	resp, err := u.OrdService.Mint(ord_service.MintRequest{
		WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
		FileUrl:    fileURI,
		FeeRate:    entity.DEFAULT_FEE_RATE, //temp
		DryRun:     false,
	})
	u.Notify(rootSpan, fmt.Sprintf("[MintFor][projectID %s]", ethAddress.ProjectID), ethAddress.UserAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", ethAddress.UserAddress, resp.Stdout))
	if err != nil {
		log.Error("ETHMint.Mint", err.Error(), err)
		return nil, err
	}

	tmpText := resp.Stdout
	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
	btcMintResp := &ord_service.MintStdoputRespose{}

	bytes := []byte(jsonStr)
	err = json.Unmarshal(bytes, btcMintResp)
	if err != nil {
		log.Error("BTCMint.helpers.JsonTransform", err.Error(), err)
		return nil, err
	}

	ethAddress.MintResponse = entity.MintStdoputResponse(*btcMintResp)
	ethAddress, err = u.UpdateEthMintedStatus(span, ethAddress)
	if err != nil {
		log.Error("ETHMint.UpdateBtcMintedStatus", err.Error(), err)
		return nil, err
	}

	return ethAddress, nil
}

func (u Usecase) ReadGCSFolderETH(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("ReadGCSFolderETH", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateEthMintedStatus(rootSpan opentracing.Span, ethWallet *entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatus", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", ethWallet)
	ethWallet.IsMinted = true

	updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(ethWallet.OrdAddress, ethWallet)
	if err != nil {
		log.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return ethWallet, nil
}

func (u Usecase) BalanceETHLogic(rootSpan opentracing.Span, ethEntity entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("BalanceETHLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	// check eth balance:
	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
	if err != nil {
		return nil, err
	}
	ethClient := eth.NewClient(ethClientWrap)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, err := ethClient.GetBalance(ctx, ethEntity.OrdAddress)
	if err != nil {
		return nil, err
	}

	log.SetData("balance", balance)

	if balance == nil {
		err = errors.New("balance is nil")
		return nil, err
	}

	ethEntity.BalanceCheckTime = ethEntity.BalanceCheckTime + 1
	updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(ethEntity.OrdAddress, &ethEntity)
	if err != nil {
		log.Error("u.Repo.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	// check total amount = received amount?
	amount, ok := big.NewInt(0).SetString(ethEntity.Amount, 10)
	if !ok {
		err := errors.New("ethEntity.Amount.OK.False")
		return nil, err
	}

	if r := balance.Cmp(amount); r == -1 {
		err := errors.New("Not enough amount")
		return nil, err
	}

	log.SetData("userWallet", ethEntity.UserAddress)
	log.SetData("ordWalletAddress", ethEntity.OrdAddress)

	ethEntity.IsConfirm = true
	ethEntity.Balance = balance.String()
	//TODO - save balance
	updated, err = u.Repo.UpdateEthWalletAddressByOrdAddr(ethEntity.OrdAddress, &ethEntity)
	if err != nil {
		log.Error("u.CheckBalance.updatedStatus", err.Error(), err)
		return nil, err
	}
	log.SetData("updated", updated)

	return &ethEntity, nil
}

func (u Usecase) MintLogicETH(rootSpan opentracing.Span, ethEntity *entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("MintLogicETH", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var err error

	//if this was minted, skip it
	if ethEntity.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("ETHMint.Minted", err.Error(), err)
		return nil, err
	}

	if !ethEntity.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		log.Error("ETHMint.IsConfirmed", err.Error(), err)
		return nil, err
	}

	log.SetData("ethEntity", ethEntity)
	return ethEntity, nil
}

//Mint flow
func (u Usecase) WaitingForETHBalancing() ([]entity.ETHWalletAddress, error) {
	span, log := u.StartSpanWithoutRoot("WaitingForETHBalancing")
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListProcessingETHWalletAddress()
	if err != nil {
		log.Error("WaitingForETHBalancing.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	log.SetData("addreses", addreses)
	for _, item := range addreses {
		func(rootSpan opentracing.Span, item entity.ETHWalletAddress) {
			span, log := u.StartSpan(fmt.Sprintf("WaitingForETHMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)
			log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
			log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
			newItem, err := u.BalanceETHLogic(span, item)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData(fmt.Sprintf("WaitingForETHBalancing.BalanceLogic.%s", item.OrdAddress), newItem)
			u.Notify(rootSpan, fmt.Sprintf("[WaitingForBalance][projectID %s]", item.ProjectID), item.UserAddress, fmt.Sprintf("%s checking ETH %s from [user_address] %s", item.OrdAddress, newItem.Balance, item.UserAddress))
			updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				log.Error(fmt.Sprintf("WaitingForETHBalancing.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("updated", updated)

			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
				MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
				Owner:         "",
				ProjectID:     item.ProjectID,
				Action:        entity.BLANCE,
				Type:          entity.ETH,
				TraceID:       u.Tracer.TraceID(span),
				ProccessID: item.UUID,
			})
		}(span, item)

		time.Sleep(2 * time.Second)
	}

	return nil, nil
}


func (u Usecase) WaitingForETHMinting() ([]entity.ETHWalletAddress, error) {
	span, log := u.StartSpanWithoutRoot("WaitingForETHMinting")
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListMintingETHWalletAddress()
	if err != nil {
		log.Error("WaitingForETHMinting.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	log.SetData("addreses", addreses)
	for _, item := range addreses {
		func(rootSpan opentracing.Span, item entity.ETHWalletAddress) {
			span, log := u.StartSpan(fmt.Sprintf("WaitingForETHMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)
			log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
			log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
			
			if item.MintResponse.Inscription != "" {
				err = errors.New("Token is being minted")
				log.Error("Token.minted", err.Error(), err)
				return
			}

			mintReps, fileURI, err := u.BTCMint(span, structure.BctMintData{Address: item.OrdAddress})
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}

			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
				TokenID:       mintReps.Inscription,
				Commit:        mintReps.Commit,
				Reveal:        mintReps.Reveal,
				Fees:          mintReps.Fees,
				MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
				Owner:         "",
				ProjectID:     item.ProjectID,
				Action:        entity.MINT,
				Type:          entity.ETH,
				TraceID:       u.Tracer.TraceID(span),
				ProccessID: item.UUID,
			})

			log.SetData("btc.Minted", mintReps)
			item.MintResponse = entity.MintStdoputResponse(*mintReps)
			item.IsMinted = true
			item.FileURI = *fileURI
			updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, &item)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}

			log.SetData("updated", updated)

		}(span, item)

		time.Sleep(2 * time.Second)
	}

	return nil, nil
}

func (u Usecase) WaitingForETHMinted() ([]entity.ETHWalletAddress, error) {
	span, log := u.StartSpanWithoutRoot("WaitingForETHMinted")
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListETHAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListETHAddress", err.Error(), err)
		return nil, err
	}

	_, bs, err := u.buildBTCClient()

	log.SetData("addreses", addreses)
	for _, item := range addreses {
		func(rootSpan opentracing.Span, item entity.ETHWalletAddress) {
			span, log := u.StartSpan(fmt.Sprintf("WaitingForETHMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			log.SetData("userWallet", item.UserAddress)
			log.SetData("ordWalletAddress", item.OrdAddress)

			//check token is created or not via BlockcypherService
			txInfo, err := bs.CheckTx(item.MintResponse.Reveal)
			if err != nil {
				log.Error(" bs.CheckTx", err.Error(), err)
				u.Notify(rootSpan, fmt.Sprintf("[Error][ETH][SendToken.bs.CheckTx][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %v", err.Error(), item))
				return
			}

			log.SetData("txInfo", txInfo)
			if txInfo.Confirmations > 1 { 
				sentTokenResp, err := u.SendToken(span, item.UserAddress, item.MintResponse.Inscription) // TODO: BAO update this logic.
				if err != nil {
					u.Notify(rootSpan, fmt.Sprintf("[Error][ETH][SendToken][projectID %s]", item.ProjectID), item.InscriptionID,  fmt.Sprintf("%s, object: %v", err.Error(), item))
					log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
					return
				}

				u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
					TokenID:       item.MintResponse.Inscription,
					Commit:        item.MintResponse.Commit,
					Reveal:        item.MintResponse.Reveal,
					Fees:          item.MintResponse.Fees,
					MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
					Owner:         item.UserAddress,
					Action:        entity.SENT,
					ProjectID:     item.ProjectID,
					Type:          entity.ETH,
					TraceID:       u.Tracer.TraceID(span),
					ProccessID: item.UUID,
				})

				u.Notify(rootSpan, fmt.Sprintf("[SendToken][ProjectID: %s]", item.ProjectID), item.UserAddress, item.MintResponse.Inscription)

				log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)

				//TODO - fund via ETH

				item.MintResponse.IsSent = true
				updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, &item)
				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateEthWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
					return
				}
				log.SetData("updated", updated)

				//TODO: - create entity.TokenURI
				_, err = u.CreateBTCTokenURI(span, item.ProjectID, item.MintResponse.Inscription, item.FileURI, entity.ETH)
				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
					return
				}

				err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateTokenOnchainStatusByTokenId.Error", item.OrdAddress), err.Error(), err)
					return
				}
			}else{
				log.SetData("checkTx.Inscription.Existed", false)
			}

			

		}(span, item)
		time.Sleep(5 * time.Second)
	}

	return nil, nil
}

//Mint flow
func (u Usecase) convertBTCToETH(rootSpan opentracing.Span, amount string) (string, error) {

	span, log := u.StartSpan("convertBTCToETH", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	//amount = "0.1"
	powIntput := math.Pow10(8)
	powIntputBig := new(big.Float)
	powIntputBig.SetFloat64(powIntput)

	log.SetData("amount", amount)
	amountMintBTC, _ := big.NewFloat(0).SetString(amount)
	amountMintBTC.Mul(amountMintBTC, powIntputBig)
	// if err != nil {
	// 	log.Error("strconv.ParseFloat", err.Error(), err)
	// 	return "", err
	// }

	_ = amountMintBTC
	btcPrice, err := helpers.GetExternalPrice("BTC")
	if err != nil {
		log.Error("strconv.getExternalPrice", err.Error(), err)
		return "", err
	}

	log.SetData("btcPrice", btcPrice)
	ethPrice, err := helpers.GetExternalPrice("ETH")
	if err != nil {
		log.Error("strconv.getExternalPrice", err.Error(), err)
		return "", err
	}
	log.SetData("ethPrice", ethPrice)

	// amountMintBTCBigInt := web3.FloatAsInt(amountMintBTC, 8)

	log.SetData("amountMintBTC", amountMintBTC.String())
	//btcToETH := btcPrice / ethPrice
	btcToETH := 14.27

	rate := new(big.Float)
	rate.SetFloat64(btcToETH)
	log.SetData("rate", rate.String())

	amountMintBTC.Mul(amountMintBTC, rate)
	log.SetData("btcToETH", btcToETH)

	pow := math.Pow10(10)
	powBig := new(big.Float)
	powBig.SetFloat64(pow)

	amountMintBTC.Mul(amountMintBTC, powBig)
	log.SetData("amountMintBTC.Mul", btcToETH)

	result := new(big.Int)
	amountMintBTC.Int(result)

	return result.String(), nil
}

// containsIgnoreCase ...
// Todo: move to helper function
func containsIgnoreCase(strSlice []string, item string) bool {
	for _, str := range strSlice {
		if strings.EqualFold(str, item) {
			return true
		}
	}

	return false
}
