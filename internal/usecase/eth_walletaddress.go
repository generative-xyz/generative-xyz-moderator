package usecase

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/davecgh/go-spew/spew"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/jinzhu/copier"

// 	"rederinghub.io/external/nfts"
// 	"rederinghub.io/external/ord_service"
// 	"rederinghub.io/internal/entity"
// 	"rederinghub.io/internal/usecase/structure"
// 	"rederinghub.io/utils/contracts/erc20"
// 	"rederinghub.io/utils/eth"
// 	"rederinghub.io/utils/helpers"
// )

// func (u Usecase) CreateETHWalletAddress(input structure.EthWalletAddressData) (*entity.ETHWalletAddress, error) {

// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))

// 	walletAddress := &entity.ETHWalletAddress{}
// 	err := copier.Copy(walletAddress, input)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.Copy", zap.Error(err))
// 		return nil, err
// 	}

// 	ethClient := eth.NewClient(nil)

// 	privKey, pubKey, address, err := ethClient.GenerateAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ethClient.GenerateAddress", zap.Error(err))
// 		return nil, err
// 	} else {
// 		walletAddress.Mnemonic = privKey
// 	}

// 	logger.AtLog.Logger.Info("CreateETHWalletAddress.createdWallet", fmt.Sprintf("%v %v %v", zap.Any("privKey, pubKey, address)", privKey, pubKey, address)))

// 	logger.AtLog.Logger.Info("CreateETHWalletAddress.receive", zap.Any("address", address))
// 	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.FindProjectByTokenID", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("found.Project", zap.Any("p.ID", p.ID))
// 	mintPriceInt, err := strconv.ParseInt(p.MintPrice, 10, 64)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("convertBTCToInt", zap.Error(err))
// 		return nil, err
// 	}
// 	if p.NetworkFee != "" {
// 		// extra network fee
// 		networkFee, err1 := strconv.ParseInt(p.NetworkFee, 10, 64)
// 		if err1 != nil {
// 			logger.AtLog.Logger.Error("convertBTCToInt", zap.Error(err))
// 		} else {
// 			mintPriceInt += networkFee
// 		}
// 	}
// 	mintPrice, _, _, err := u.convertBTCToETH(fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
// 	if err != nil {
// 		logger.AtLog.Logger.Error("convertBTCToETH", zap.Error(err))
// 		return nil, err
// 	}
// 	walletAddress.Amount = mintPrice // 0.023 * 1e18 eth
// 	walletAddress.UserAddress = input.WalletAddress
// 	walletAddress.OrdAddress = strings.ReplaceAll(address, "\n", "")
// 	walletAddress.IsConfirm = false
// 	walletAddress.IsMinted = false
// 	walletAddress.FileURI = ""       //find files from google store
// 	walletAddress.InscriptionID = "" //find files from google store
// 	walletAddress.ProjectID = input.ProjectID
// 	walletAddress.Balance = "0"
// 	walletAddress.BalanceCheckTime = 0

// 	err = u.Repo.InsertEthWalletAddress(walletAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.InsertEthWalletAddress", zap.Error(err))
// 		return nil, err
// 	}

// 	return walletAddress, nil
// }

// func (u Usecase) IsWhitelistedAddress(ctx context.Context, userAddr string, whitelistedAddrs []string) (bool, error) {

// 	logger.AtLog.Logger.Info("whitelistedAddrs", zap.Any("whitelistedAddrs", whitelistedAddrs))
// 	if len(whitelistedAddrs) == 0 {
// 		logger.AtLog.Logger.Info("whitelistedAddrs.Total", zap.Any("len(whitelistedAddrs)", len(whitelistedAddrs)))
// 		return false, nil
// 	}
// 	filter := nfts.MoralisFilter{}
// 	filter.Limit = new(int)
// 	*filter.Limit = 1
// 	filter.TokenAddresses = new([]string)
// 	*filter.TokenAddresses = whitelistedAddrs

// 	logger.AtLog.Logger.Info("filter.GetNftByWalletAddress", zap.Any("filter", filter))
// 	resp, err := u.MoralisNft.GetNftByWalletAddress(userAddr, filter)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.MoralisNft.GetNftByWalletAddress", zap.Error(err))
// 		return false, err
// 	}

// 	logger.AtLog.Logger.Info("resp", zap.Any("resp", resp))
// 	if len(resp.Result) > 0 {
// 		return true, nil
// 	}

// 	delegations, err := u.DelegateService.GetDelegationsByDelegate(ctx, userAddr)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.DelegateService.GetDelegationsByDelegate", zap.Error(err))
// 		return false, err
// 	}

// 	logger.AtLog.Logger.Info("delegations", zap.Any("delegations", delegations))
// 	for _, delegation := range delegations {

// 		if delegation.Type == 2 || delegation.Type == 3 {
// 			if containsIgnoreCase(whitelistedAddrs, delegation.Contract.String()) {
// 				return true, nil
// 			}
// 		}else if ( delegation.Type == 1) {
// 			resp, err := u.MoralisNft.GetNftByWalletAddress(delegation.Vault.Hex(), filter)
// 			if err != nil {
// 				logger.AtLog.Logger.Error("u.MoralisNft.GetNftByWalletAddress", zap.Error(err))
// 				continue
// 			}

// 			logger.AtLog.Logger.Info("resp", zap.Any("resp", resp))
// 			if len(resp.Result) > 0 {
// 				return true, nil
// 			}
// 		}

// 	}
// 	return false, nil
// }

// func (u Usecase) CreateWhitelistedETHWalletAddress(ctx context.Context, userAddr string, input structure.EthWalletAddressData) (*entity.ETHWalletAddress, error) {

// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))

// 	weth, err := u.Repo.FindDelegateEthWalletAddressByUserAddress(userAddr)
// 	if err == nil {
// 		if weth != nil {
// 			logger.AtLog.Logger.Info("ethWalletAddress", zap.Any("weth", weth))
// 			if weth.IsConfirm == true {
// 				err = errors.New("This account has applied discount")
// 				logger.AtLog.Logger.Error("applied.Discount", zap.Error(err))
// 				return nil, err
// 			}

// 			return weth, nil
// 		}
// 	} else {
// 		logger.AtLog.Logger.Error("u.Repo.FindEthWalletAddressByUserAddress", zap.Error(err))
// 	}

// 	walletAddress := &entity.ETHWalletAddress{}
// 	err = copier.Copy(walletAddress, input)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.Copy", zap.Error(err))
// 		return nil, err
// 	}

// 	// userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)
// 	// resp, err := u.OrdService.Exec(ord_service.ExecRequest{
// 	// 	Args: []string{
// 	// 		"--wallet",
// 	// 		userWallet,
// 	// 		"wallet",
// 	// 		"create",
// 	// 	},
// 	// })

// 	ethClient := eth.NewClient(nil)

// 	privKey, pubKey, address, err := ethClient.GenerateAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ethClient.GenerateAddress", zap.Error(err))
// 		return nil, err
// 	} else {
// 		walletAddress.Mnemonic = privKey
// 	}

// 	logger.AtLog.Logger.Info("CreateETHWalletAddress.createdWallet", fmt.Sprintf("%v %v %v", zap.Any("privKey, pubKey, address)", privKey, pubKey, address)))

// 	logger.AtLog.Logger.Info("CreateETHWalletAddress.receive", zap.Any("address", address))
// 	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.FindProjectByTokenID", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("found.Project", zap.Any("p.ID", p.ID))
// 	mintPriceInt, err := strconv.Atoi(p.MintPrice)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("convertBTCToInt", zap.Error(err))
// 		return nil, err
// 	}
// 	if p.NetworkFee != "" {
// 		// extra network fee
// 		networkFee, err1 := strconv.Atoi(p.NetworkFee)
// 		if err1 != nil {
// 			logger.AtLog.Logger.Error("convertBTCToInt", zap.Error(err))
// 		} else {
// 			mintPriceInt += networkFee
// 		}
// 	}
// 	mintPrice, _, _, err := u.convertBTCToETH(fmt.Sprintf("%f", float64(mintPriceInt)/1e8))
// 	if err != nil {
// 		logger.AtLog.Logger.Error("convertBTCToETH", zap.Error(err))
// 		return nil, err
// 	}

// 	walletAddress.Amount = mintPrice // 0.023 * 1e18 eth

// 	isWhitelist, err := u.IsWhitelistedAddress(ctx, userAddr, p.WhiteListEthContracts)

// 	if isWhitelist {
// 		whitelistedPrice := new(big.Float)
// 		ethPrice, _ := helpers.GetExternalPrice("ETH")
// 		if ethPrice == 0 {
// 			ethPrice = 1500
// 		}
// 		whitelistedPrice.SetFloat64(50.0 / ethPrice)
// 		whitelistedPrice = whitelistedPrice.Mul(whitelistedPrice, big.NewFloat(1e18))

// 		intPrice := new(big.Int)
// 		whitelistedPrice.Int(intPrice)

// 		walletAddress.Amount = intPrice.String()
// 	}

// 	walletAddress.UserAddress = input.WalletAddress
// 	walletAddress.OrdAddress = strings.ReplaceAll(address, "\n", "")
// 	walletAddress.IsConfirm = false
// 	walletAddress.IsMinted = false
// 	walletAddress.FileURI = ""       //find files from google store
// 	walletAddress.InscriptionID = "" //find files from google store
// 	walletAddress.ProjectID = input.ProjectID
// 	walletAddress.Balance = "0"
// 	walletAddress.BalanceCheckTime = 0
// 	walletAddress.DelegatedAddress = userAddr

// 	err = u.Repo.InsertEthWalletAddress(walletAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CreateETHWalletAddress.InsertEthWalletAddress", zap.Error(err))
// 		return nil, err
// 	}

// 	return walletAddress, nil
// }

// func (u Usecase) ETHMint(input structure.BctMintData) (*entity.ETHWalletAddress, error) {

// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))

// 	ethAddress, err := u.Repo.FindEthWalletAddressByOrd(input.Address)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.FindBtcWalletAddressByOrd", zap.Error(err))
// 		return nil, err
// 	}

// 	//mint logic
// 	ethAddress, err = u.MintLogicETH(ethAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.MintLogic", zap.Error(err))
// 		return nil, err
// 	}

// 	// get data from project
// 	p, err := u.Repo.FindProjectByTokenID(ethAddress.ProjectID)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.FindProjectByTokenID", zap.Error(err))
// 		return nil, err
// 	}
// 	//logger.AtLog.Logger.Info("found.Project", zap.Any("p", p))

// 	//prepare data for mint
// 	// - Get project.AnimationURL
// 	projectNftTokenUri := &structure.ProjectAnimationUrl{}
// 	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.helpers.Base64DecodeRaw", zap.Error(err))
// 		return nil, err
// 	}

// 	// - Upload the Animation URL to GCS
// 	animation := projectNftTokenUri.AnimationUrl
// 	animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")

// 	now := time.Now().UTC().Unix()
// 	uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", p.TokenID, now))
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.helpers.Base64DecodeRaw", zap.Error(err))
// 		return nil, err
// 	}

// 	fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
// 	ethAddress.FileURI = fileURI

// 	//TODO - enable this
// 	resp, err := u.OrdService.Mint(ord_service.MintRequest{
// 		WalletName: os.Getenv("ORD_MASTER_ADDRESS"),
// 		FileUrl:    fileURI,
// 		FeeRate:    entity.DEFAULT_FEE_RATE, //temp
// 		DryRun:     false,
// 	})
// 	u.Notify(fmt.Sprintf("[MintFor][projectID %s]", ethAddress.ProjectID), ethAddress.UserAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", ethAddress.UserAddress, resp.Stdout))
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.Mint", zap.Error(err))
// 		return nil, err
// 	}

// 	tmpText := resp.Stdout
// 	//tmpText := `{\n  \"commit\": \"7a47732d269d5c005c4df99f2e5cf1e268e217d331d175e445297b1d2991932f\",\n  \"inscription\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2afi0\",\n  \"reveal\": \"9925b5626058424d2fc93760fb3f86064615c184ac86b2d0c58180742683c2af\",\n  \"fees\": 185514\n}\n`
// 	jsonStr := strings.ReplaceAll(tmpText, `\n`, "")
// 	jsonStr = strings.ReplaceAll(jsonStr, "\\", "")
// 	btcMintResp := &ord_service.MintStdoputRespose{}

// 	bytes := []byte(jsonStr)
// 	err = json.Unmarshal(bytes, btcMintResp)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("BTCMint.helpers.JsonTransform", zap.Error(err))
// 		return nil, err
// 	}

// 	ethAddress.MintResponse = entity.MintStdoputResponse(*btcMintResp)
// 	ethAddress, err = u.UpdateEthMintedStatus(ethAddress)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("ETHMint.UpdateBtcMintedStatus", zap.Error(err))
// 		return nil, err
// 	}

// 	return ethAddress, nil
// }

// func (u Usecase) ReadGCSFolderETH(input structure.BctWalletAddressData) (*entity.ETHWalletAddress, error) {

// 	logger.AtLog.Logger.Info("input", zap.Any("input", input))
// 	u.GCS.ReadFolder("btc-projects/project-1/")
// 	return nil, nil
// }

// func (u Usecase) UpdateEthMintedStatus(ethWallet *entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {

// 	logger.AtLog.Logger.Info("input", zap.Any("ethWallet", ethWallet))
// 	ethWallet.IsMinted = true

// 	updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(ethWallet.OrdAddress, ethWallet)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))
// 	return ethWallet, nil
// }

// func (u Usecase) BalanceETHLogic(ethEntity entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {

// 	// check eth balance:
// 	ethClientWrap, err := ethclient.Dial(u.Config.BlockchainConfig.ETHEndpoint)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ethClient := eth.NewClient(ethClientWrap)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	balance, err := ethClient.GetBalance(ctx, ethEntity.OrdAddress)
// 	if err != nil {
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("balance", zap.Any("balance", balance))

// 	if balance == nil {
// 		err = errors.New("balance is nil")
// 		return nil, err
// 	}

// 	ethEntity.BalanceCheckTime = ethEntity.BalanceCheckTime + 1
// 	updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(ethEntity.OrdAddress, &ethEntity)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.Repo.UpdateBtcWalletAddressByOrdAddr", zap.Error(err))
// 		return nil, err
// 	}

// 	// check total amount = received amount?
// 	amount, ok := big.NewInt(0).SetString(ethEntity.Amount, 10)
// 	if !ok {
// 		err := errors.New("ethEntity.Amount.OK.False")
// 		return nil, err
// 	}

// 	if r := balance.Cmp(amount); r == -1 {
// 		err := errors.New("Not enough amount")
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("userWallet", zap.Any("ethEntity.UserAddress", ethEntity.UserAddress))
// 	logger.AtLog.Logger.Info("ordWalletAddress", zap.Any("ethEntity.OrdAddress", ethEntity.OrdAddress))

// 	ethEntity.IsConfirm = true
// 	ethEntity.Balance = balance.String()
// 	//TODO - save balance
// 	updated, err = u.Repo.UpdateEthWalletAddressByOrdAddr(ethEntity.OrdAddress, &ethEntity)
// 	if err != nil {
// 		logger.AtLog.Logger.Error("u.CheckBalance.updatedStatus", zap.Error(err))
// 		return nil, err
// 	}
// 	logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

// 	return &ethEntity, nil
// }

// func (u Usecase) MintLogicETH(ethEntity *entity.ETHWalletAddress) (*entity.ETHWalletAddress, error) {

// 	var err error

// 	//if this was minted, skip it
// 	if ethEntity.IsMinted {
// 		err = errors.New("This btc was minted")
// 		logger.AtLog.Logger.Error("ETHMint.Minted", zap.Error(err))
// 		return nil, err
// 	}

// 	if !ethEntity.IsConfirm {
// 		err = errors.New("This btc must be IsConfirmed")
// 		logger.AtLog.Logger.Error("ETHMint.IsConfirmed", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("ethEntity", zap.Any("ethEntity", ethEntity))
// 	return ethEntity, nil
// }

// // Mint flow
// func (u Usecase) WaitingForETHBalancing() ([]entity.ETHWalletAddress, error) {

// 	addreses, err := u.Repo.ListProcessingETHWalletAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("WaitingForETHBalancing.ListProcessingWalletAddress", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.ETHWalletAddress) {

// 			newItem, err := u.BalanceETHLogic(item)
// 			if err != nil {
// 				logger.AtLog.Logger.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), zap.Error(err))
// 				return
// 			}
// 			logger.AtLog.Logger.Info(fmt.Sprintf("WaitingForETHBalancing.BalanceLogic.%s", item.OrdAddress), newItem)
// 			u.Notify(fmt.Sprintf("[WaitingForBalance][projectID %s]", item.ProjectID), item.UserAddress, fmt.Sprintf("%s checking ETH %s from [user_address] %s", item.OrdAddress, newItem.Balance, item.UserAddress))
// 			updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, newItem)
// 			if err != nil {
// 				logger.AtLog.Logger.Error(fmt.Sprintf("WaitingForETHBalancing.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), zap.Error(err))
// 				return
// 			}
// 			logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

// 			u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
// 				MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
// 				Owner:         "",
// 				ProjectID:     item.ProjectID,
// 				Action:        entity.BLANCE,
// 				Type:          entity.ETH,
// 				ProccessID:    item.UUID,
// 			})
// 		}(item)

// 		time.Sleep(2 * time.Second)
// 	}

// 	return nil, nil
// }

// func (u Usecase) WaitingForETHMinting() ([]entity.ETHWalletAddress, error) {

// 	addreses, err := u.Repo.ListMintingETHWalletAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("WaitingForETHMinting.ListProcessingWalletAddress", zap.Error(err))
// 		return nil, err
// 	}

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.ETHWalletAddress) {

// 			if item.MintResponse.Inscription != "" {
// 				err = errors.New("Token is being minted")
// 				logger.AtLog.Logger.Error("Token.minted", zap.Error(err))
// 				return
// 			}

// 			_, _, err := u.BTCMint(structure.BctMintData{Address: item.OrdAddress})
// 			if err != nil {
// 				logger.AtLog.Logger.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), zap.Error(err))
// 				return
// 			}
// 		}(item)

// 		time.Sleep(2 * time.Second)
// 	}

// 	return nil, nil
// }

// func (u Usecase) WaitingForETHMinted() ([]entity.ETHWalletAddress, error) {

// 	addreses, err := u.Repo.ListETHAddress()
// 	if err != nil {
// 		logger.AtLog.Logger.Error("WillBeProcessWTC.ListETHAddress", zap.Error(err))
// 		return nil, err
// 	}

// 	_, bs, err := u.buildBTCClient()

// 	logger.AtLog.Logger.Info("addreses", zap.Any("addreses", addreses))
// 	for _, item := range addreses {
// 		func(item entity.ETHWalletAddress) {

// 			logger.AtLog.Logger.Info("userWallet", zap.Any("item.UserAddress", item.UserAddress))
// 			logger.AtLog.Logger.Info("ordWalletAddress", zap.Any("item.OrdAddress", item.OrdAddress))

// 			//check token is created or not via BlockcypherService
// 			txInfo, err := bs.CheckTx(item.MintResponse.Reveal)
// 			if err != nil {
// 				logger.AtLog.Logger.Error(" bs.CheckTx", zap.Error(err))
// 				u.Notify(fmt.Sprintf("[Error][ETH][SendToken.bs.CheckTx][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
// 				return
// 			}

// 			logger.AtLog.Logger.Info("txInfo", zap.Any("txInfo", txInfo))
// 			if txInfo.Confirmations > 1 {
// 				sentTokenResp, err := u.SendToken(item.UserAddress, item.MintResponse.Inscription) // TODO: BAO update this logic.
// 				if err != nil {
// 					u.Notify(fmt.Sprintf("[Error][ETH][SendToken][projectID %s]", item.ProjectID), item.InscriptionID, fmt.Sprintf("%s, object: %s", err.Error(), item.UUID))
// 					logger.AtLog.Logger.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), zap.Error(err))
// 					return
// 				}

// 				u.Repo.CreateTokenUriHistory(&entity.TokenUriHistories{
// 					TokenID:       item.MintResponse.Inscription,
// 					Commit:        item.MintResponse.Commit,
// 					Reveal:        item.MintResponse.Reveal,
// 					Fees:          item.MintResponse.Fees,
// 					MinterAddress: os.Getenv("ORD_MASTER_ADDRESS"),
// 					Owner:         item.UserAddress,
// 					Action:        entity.SENT,
// 					ProjectID:     item.ProjectID,
// 					Type:          entity.ETH,
// 					ProccessID:    item.UUID,
// 				})

// 				u.Notify(fmt.Sprintf("[SendToken][ProjectID: %s]", item.ProjectID), item.UserAddress, item.MintResponse.Inscription)

// 				logger.AtLog.Logger.Info(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)

// 				//TODO - fund via ETH

// 				item.MintResponse.IsSent = true
// 				updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, &item)
// 				if err != nil {
// 					logger.AtLog.Logger.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateEthWalletAddressByOrdAddr.Error", item.OrdAddress), zap.Error(err))
// 					return
// 				}
// 				logger.AtLog.Logger.Info("updated", zap.Any("updated", updated))

// 				//TODO: - create entity.TokenURI
// 				_, err = u.CreateBTCTokenURI(item.ProjectID, item.MintResponse.Inscription, item.FileURI, entity.ETH)
// 				if err != nil {
// 					logger.AtLog.Logger.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), zap.Error(err))
// 					return
// 				}

// 				err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
// 				if err != nil {
// 					logger.AtLog.Logger.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateTokenOnchainStatusByTokenId.Error", item.OrdAddress), zap.Error(err))
// 					return
// 				}
// 			} else {
// 				logger.AtLog.Logger.Info("checkTx.Inscription.Existed", zap.Any("false", false))
// 			}

// 		}(item)
// 		time.Sleep(5 * time.Second)
// 	}

// 	return nil, nil
// }

// // containsIgnoreCase ...
// // Todo: move to helper function
// func containsIgnoreCase(strSlice []string, item string) bool {
// 	for _, str := range strSlice {
// 		if strings.EqualFold(str, item) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func (u Usecase) IsWhitelistedAddressERC20(ctx context.Context, userAddr string, erc20WhiteList map[string]structure.Erc20Config) (bool, error) {
// 	client, err  := helpers.EthDialer()
// 	if err != nil {
// 		return false, err
// 	}

// 	for addr, whitelistedThres := range erc20WhiteList {
// 		erc20Client, err := erc20.NewErc20(common.HexToAddress(addr), client)
// 		if err != nil {
// 			continue
// 		}

// 		blance, err := erc20Client.BalanceOf(nil,  common.HexToAddress(userAddr))
// 		if err != nil {
// 			continue
// 		}

// 		pow := new(big.Int)
// 		pow = pow.Exp(big.NewInt(1), big.NewInt(whitelistedThres.Decimal), nil)
// 		confValue := big.NewInt(whitelistedThres.Value)

// 		confValue = confValue.Mul(confValue, pow)

// 		//bigInt64 := big.
// 		tmp := blance.Cmp(confValue)

// 		spew.Dump(whitelistedThres.Value, whitelistedThres.Decimal)
// 		spew.Dump(confValue.String())
// 		spew.Dump(blance.String())
// 		if tmp >= 0 {
// 			return true, nil
// 		}
// 	}

// 	return false, nil
// }

