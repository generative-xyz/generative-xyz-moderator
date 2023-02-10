package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/eth"
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
	mintPrice, err := convertBTCToETH(p.MintPrice)
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

	// check total amount = received amount?
	amount, _ := big.NewInt(0).SetString(ethEntity.Amount, 10)

	if r := balance.Cmp(amount); r == -1 {
		err := errors.New("Not enough amount")
		return nil, err

	}

	log.SetData("userWallet", ethEntity.UserAddress)
	log.SetData("ordWalletAddress", ethEntity.OrdAddress)

	ethEntity.IsConfirm = true

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
		log.SetData("userWallet", item.UserAddress)
		log.SetData("ordWalletAddress", item.OrdAddress)
		newItem, err := u.BalanceETHLogic(span, item)
		if err != nil {
			//log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)

		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, newItem)
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		btc, err := u.BTCMint(span, structure.BctMintData{Address: newItem.OrdAddress})
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateEthWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
			continue
		}

		_ = btc
	}

	return nil, nil
}

func (u Usecase) WaitingForETHMinted(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("WaitingForMinted", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListETHAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListETHAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		log.SetData("userWallet", item.UserAddress)
		log.SetData("ordWalletAddress", item.OrdAddress)
		sentTokenResp, err := u.SendToken(item.UserAddress, item.MintResponse.Inscription) // TODO: BAO update this logic.
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}

		log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)
		// amout, err := strconv.ParseFloat(item.Amount, 10)
		// if err != nil {
		// 	log.Error("ListenTheMintedBTC.%s. strconv.ParseFloa.Error", err.Error(), err)
		// 	continue
		// }
		// amout = amout * 0.9

		// fundData := ord_service.ExecRequest{
		// 	Args: []string{
		// 		"--wallet",
		// 		item.OrdAddress,
		// 		"send",
		// 		"ord_master",
		// 		fmt.Sprintf("%f", amout),
		// 		"--fee-rate",
		// 		"15",
		// 	},
		// }

		// log.SetData("fundData", fundData)
		// fundResp, err := u.OrdService.Exec(fundData)

		// if err != nil {
		// 	log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.ReFund.Error", item.OrdAddress), err.Error(), err)
		// 	continue
		// }

		// log.SetData("fundResp", fundResp)

		item.MintResponse.IsSent = true
		updated, err := u.Repo.UpdateEthWalletAddressByOrdAddr(item.OrdAddress, &item)
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateEthWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		//TODO: - create entity.TokenURI
		_, err = u.CreateBTCTokenURI(span, item.ProjectID, item.MintResponse.Inscription)
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
			continue
		}
	}

	return nil, nil
}

func convertBTCToETH(amount string) (string, error) {
	amountMintBTC, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		return "", err
	}
	btcPrice, err := getExternalPrice("BTC")
	if err != nil {
		return "", err
	}
	ethPrice, err := getExternalPrice("ETH")
	if err != nil {
		return "", err
	}
	amountMintBTC = amountMintBTC / math.Pow10(8)

	btcToETH := btcPrice / ethPrice

	amountMintETH := amountMintBTC * btcToETH
	amountMintETHInt := uint64(math.Floor(amountMintETH * math.Pow10(18)))
	return fmt.Sprintf("%d", amountMintETHInt), nil
}

func getExternalPrice(tokenSymbol string) (float64, error) {
	binancePriceURL := "https://api.binance.com/api/v3/ticker/price?symbol="
	var price struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
	var jsonErr struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	retryTimes := 0
retry:
	retryTimes++
	if retryTimes > 2 {
		return 0, nil
	}
	tk := strings.ToUpper(tokenSymbol)
	resp, err := http.Get(binancePriceURL + tk + "USDT")
	if err != nil {
		log.Println(err)
		goto retry
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &price)
	if err != nil {
		err = json.Unmarshal(body, &jsonErr)
		if err != nil {
			log.Println(err)
			goto retry
		}
	}
	resp.Body.Close()
	value, err := strconv.ParseFloat(price.Price, 32)
	if err != nil {
		log.Println("getExternalPrice", tokenSymbol, err)
		return 0, nil
	}
	return value, nil
}
