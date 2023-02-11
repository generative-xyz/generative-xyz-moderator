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
	// resp, err = u.OrdService.Exec(ord_service.ExecRequest{
	// 	Args: []string{
	// 		"--wallet",
	// 		userWallet,
	// 		"wallet",
	// 		"receive",
	// 	},
	// })
	// if err != nil {
	// 	log.Error("u.OrdService.Exec.create.receive", err.Error(), err)
	// 	return nil, err
	// }

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
	mintPrice, err := u.convertBTCToETH(span, fmt.Sprintf("%f", float64(mintPriceInt) / 1e8))
	if err != nil {
		log.Error("convertBTCToETH", err.Error(), err)
		return nil, err
	}
	walletAddress.Amount = mintPrice
	walletAddress.UserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(address, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID

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
		WalletName: "ord_master",
		FileUrl:    fileURI,
		FeeRate:    15, //temp
		DryRun:     false,
	})
	u.Notify(rootSpan, "Mint for", ethAddress.UserAddress, fmt.Sprintf("Made mining transaction for %s, waiting network confirm %s", ethAddress.UserAddress, resp.Stdout))
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
	ethClientWrap, err := ethclient.Dial(u.Config.Moralis.URL)
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

	go func(rootSpan opentracing.Span, wallet *entity.ETHWalletAddress, balance *big.Int) {
		span, log := u.StartSpan("CheckBlance.RoutineUpdate", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		wallet.Balance = balance.String()
		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(wallet.OrdAddress, wallet)
		if err != nil {
			log.Error("u.Repo.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
			return
		}
		log.SetData("updated", updated)

	}(span, &ethEntity, balance)

	// check total amount = received amount?
	amount, _ := big.NewInt(0).SetString(ethEntity.Amount, 10)

	if r := balance.Cmp(amount); r == -1 {
		err := errors.New("Not enough amount")
		return nil, err

	}

	log.SetData("userWallet", ethEntity.UserAddress)
	log.SetData("ordWalletAddress", ethEntity.OrdAddress)

	ethEntity.IsConfirm = true

	//TODO - save balance

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

func (u Usecase) WaitingForETHBalancing(rootSpan opentracing.Span) ([]entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("WaitingForETHBalancing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	addreses, err := u.Repo.ListProcessingETHWalletAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		func (rootSpan opentracing.Span, item entity.ETHWalletAddress) { 
			span, log := u.StartSpan(fmt.Sprintf("WaitingForETHMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log) 
			log.SetTag(utils.WALLET_ADDRESS_TAG, item.UserAddress)
			log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
			newItem, err := u.BalanceETHLogic(span, item)
			if err != nil {
				//log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)
			u.Notify(rootSpan, "WaitingForBalancing", item.UserAddress, fmt.Sprintf("%s received ETH %s from [user_address] %s", item.OrdAddress, newItem.Balance, item.UserAddress))
			updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("updated", updated)
	
			mintReps, fileURI, err := u.BTCMint(span, structure.BctMintData{Address: newItem.OrdAddress})
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
				return
			}
	
			log.SetData("btc.Minted", mintReps)
	
			newItem.MintResponse = entity.MintStdoputResponse(*mintReps)
			newItem.IsMinted = true
			newItem.FileURI = *fileURI
			updated, err = u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, newItem)
			if err != nil {
				log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}

			//TODO: - create entity.TokenURI
			_, err = u.CreateBTCTokenURI(span, newItem.ProjectID, item.MintResponse.Inscription, newItem.FileURI, entity.ETH)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
				return
			}

		}(span, item)

		time.Sleep(5 * time.Second)
	}

	return nil, nil
}

func (u Usecase) WaitingForETHMinted(rootSpan opentracing.Span) ([]entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("WaitingForETHMinted", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListETHAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListETHAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		func (rootSpan opentracing.Span, item entity.ETHWalletAddress) { 
			span, log := u.StartSpan(fmt.Sprintf("WaitingForETHMinted.%s", item.UserAddress), rootSpan)
			defer u.Tracer.FinishSpan(span, log)

			log.SetData("userWallet", item.UserAddress)
			log.SetData("ordWalletAddress", item.OrdAddress)
			sentTokenResp, err := u.SendToken(span, item.UserAddress, item.MintResponse.Inscription) // TODO: BAO update this logic.
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
				return
			}

			log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)

			//TODO - fund via ETH

			item.MintResponse.IsSent = true
			updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, &item)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateEthWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
				return
			}
			log.SetData("updated", updated)

			err = u.Repo.UpdateTokenOnchainStatusByTokenId(item.MintResponse.Inscription)
			if err != nil {
				log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateTokenOnchainStatusByTokenId.Error", item.OrdAddress), err.Error(), err)
				return
			}

		}(span, item)
		time.Sleep(5 * time.Second)
	}

	return nil, nil
}

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
	btcToETH := btcPrice / ethPrice
	
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

