package usecase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BitcoinTokenMintFee struct {
	Amount string
	MintFee string
	SentTokenFee string
}

const (
	SENT_TOKEN_FEE = 0
)

func calculateMintPrice(input structure.BctWalletAddressDataV2) BitcoinTokenMintFee {
	base64String := input.File
	base64String = strings.ReplaceAll(base64String, "data:text/html;base64,", "")
	base64String = strings.ReplaceAll(base64String, "data:image/png;base64,", "")
	dec, _ := base64.StdEncoding.DecodeString(base64String)
	fileSize := len([]byte(dec))
	mintFee := int32(fileSize) / 4 * input.FeeRate
	return BitcoinTokenMintFee{
		Amount: strconv.FormatInt(int64(mintFee + SENT_TOKEN_FEE), 10),
		MintFee: strconv.FormatInt(int64(mintFee), 10),
		SentTokenFee: strconv.FormatInt(int64(SENT_TOKEN_FEE), 10),
	}
}

func (u Usecase) CreateBTCWalletAddressV2(rootSpan opentracing.Span, input structure.BctWalletAddressDataV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("CreateBTCWalletAddressV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)

	walletAddress := &entity.BTCWalletAddressV2{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateBTCWalletAddressV2.Copy", err.Error(), err)
		return nil, err
	}

	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)
	log.SetTag(utils.WALLET_ADDRESS_TAG, userWallet)
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"create",
		},
	})

	if err != nil {
		log.Error("u.OrdService.Exec.create.Wallet", err.Error(), err)
		//return nil, err
	} else {
		walletAddress.Mnemonic = resp.Stdout
	}

	log.SetData("CreateBTCWalletAddress.createdWallet", resp)
	resp, err = u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"receive",
		},
	})
	if err != nil {
		log.Error("u.OrdService.Exec.create.receive", err.Error(), err)
		return nil, err
	}

	log.SetData("CreateBTCWalletAddressV2.receive", resp)
	mintFee := calculateMintPrice(input)
	walletAddress.Amount = mintFee.Amount
	walletAddress.MintFee = mintFee.MintFee
	walletAddress.SentTokenFee = mintFee.SentTokenFee
	walletAddress.UserAddress = userWallet
	walletAddress.OriginUserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = input.File
	walletAddress.InscriptionID = ""

	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, walletAddress.OrdAddress)
	err = u.Repo.InsertBtcWalletAddressV2(walletAddress)
	if err != nil {
		log.Error("u.CreateBTCWalletAddressV2.InsertBtcWalletAddressV2", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) CheckbalanceWalletAddressV2(rootSpan opentracing.Span, input structure.CheckBalance) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("CheckbalanceWalletAddressV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, input.Address)
	btc, err := u.Repo.FindBtcWalletAddressByOrdV2(input.Address)
	if err != nil {
		log.Error("u.Repo.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}
	
	balance, err := u.CheckBalanceV2(span, *btc)
	if err != nil {
		log.Error("u.BalanceLogic", err.Error(), err)
		return nil, err
	}

	return balance, nil
}

func (u Usecase) BTCMintV2(rootSpan opentracing.Span, input structure.BctMintData) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("BTCMintV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag(utils.WALLET_ADDRESS_TAG, input.Address)

	btc, err := u.Repo.FindBtcWalletAddressByOrdV2(input.Address)
	if err != nil {
		log.Error("BTCMint.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	//mint logic
	btc, err = u.MintLogicV2(span, btc)
	if err != nil {
		log.Error("BTCMint.MintLogic", err.Error(), err)
		return nil, err
	}

	// - Upload the Animation URL to GCS
	animation := btc.FileURI
	animation = strings.ReplaceAll(animation, "data:text/html;base64,", "")
	animation = strings.ReplaceAll(animation, "data:image/png;base64,", "")

	now := time.Now().UTC().Unix()
	uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html", btc.OrdAddress, now))
	if err != nil {
		log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	btc.FileURI = fileURI

	//TODO - enable this
	spew.Dump(fileURI)
	resp, err := u.OrdService.Mint(ord_service.MintRequest{
		WalletName: "ord_master",
		FileUrl:    fileURI,
		FeeRate:    15, //temp
		DryRun:     false,
	})

	if err != nil {
		log.Error("BTCMint.Mint", err.Error(), err)
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

	btc.MintResponse = entity.MintStdoputResponse(*btcMintResp)
	btc, err = u.UpdateBtcMintedStatusV2(span, btc)
	if err != nil {
		log.Error("BTCMint.UpdateBtcMintedStatus", err.Error(), err)
		return nil, err
	}

	return btc, nil
}

func (u Usecase) ReadGCSFolderV2(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("ReadGCSFolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateBtcMintedStatusV2(rootSpan opentracing.Span, btcWallet *entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatusV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", btcWallet)

	btcWallet.IsMinted = true
	log.SetTag(utils.WALLET_ADDRESS_TAG, btcWallet.UserAddress)
	log.SetTag(utils.ORD_WALLET_ADDRESS_TAG, btcWallet.OrdAddress)
	log.SetTag(utils.TOKEN_ID_TAG, btcWallet.InscriptionID)

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		log.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return btcWallet, nil
}
 
func (u Usecase) CheckBalanceV2(rootSpan opentracing.Span, btc entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error)  {
	span, log := u.StartSpan("CheckBlanceV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	log.SetData(utils.WALLET_ADDRESS_TAG, btc.UserAddress)
	log.SetData(utils.ORD_WALLET_ADDRESS_TAG, btc.OrdAddress)
	
	balanceRequest := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			btc.UserAddress,
			"wallet",
			"balance",
		},
	}

	log.SetData("balanceRequest", balanceRequest)
	//userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(balanceRequest)
	if err != nil {
		log.Error("BTCMint.Exec.balance", err.Error(), err)
		return nil, err
	}

	log.SetData("balanceResponse", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	log.SetData("balance", balance)

	btc.Balance = balance

	go func(rootSpan opentracing.Span, balance *entity.BTCWalletAddressV2) {
		span, log := u.StartSpan("CheckBlance.RoutineUpdate", rootSpan)
		defer u.Tracer.FinishSpan(span, log )

		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(balance.OrdAddress, balance)
		if err != nil {
			log.Error("u.Repo.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
			return
		}
		log.SetData("updated", updated)

	}(span, &btc)

	return &btc, nil
}

func (u Usecase) BalanceLogicV2(rootSpan opentracing.Span, btc entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("BalanceLogicV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	balance, err := u.CheckBalanceV2(span, btc)
	if err != nil {
		log.Error("u.CheckBlance", err.Error(), err)
		return nil, err
	}
	
	//TODO logic of the checked balance here
	if balance.Balance < btc.Amount {
		err := errors.New("Not enough amount")
		return nil, err
	}
	btc.IsConfirm = true
	return &btc, nil
}

func (u Usecase) MintLogicV2(rootSpan opentracing.Span, btc *entity.BTCWalletAddressV2) (*entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("MintLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var err error

	log.SetData(utils.WALLET_ADDRESS_TAG, btc.UserAddress)
	log.SetData(utils.ORD_WALLET_ADDRESS_TAG, btc.OrdAddress)

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("BTCMint.Minted", err.Error(), err)
		return nil, err
	}

	if !btc.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		log.Error("BTCMint.IsConfirmed", err.Error(), err)
		return nil, err
	}

	log.SetData("btc", btc)
	return btc, nil
}

func (u Usecase) WaitingForBalancingV2(rootSpan opentracing.Span) ([]entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("WaitingForBalancing", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListProcessingWalletAddressV2()
	if err != nil {
		log.Error("WillBeProcessWTC.ListProcessingWalletAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		log.SetData(utils.WALLET_ADDRESS_TAG, item.UserAddress)
		log.SetData(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
		newItem, err := u.BalanceLogicV2(span, item)
		if err != nil {
			//log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)

		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(item.OrdAddress, newItem)
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		btc, err := u.BTCMintV2(span, structure.BctMintData{Address: newItem.OrdAddress})
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
			continue
		}

		log.SetData("btc.Minted", btc)
	}

	return nil, nil
}

func (u Usecase) WaitingForMintedV2(rootSpan opentracing.Span) ([]entity.BTCWalletAddressV2, error) {
	span, log := u.StartSpan("WaitingForMintedV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListBTCAddressV2()
	if err != nil {
		log.Error("WillBeProcessWTC.ListBTCAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		log.SetData(utils.WALLET_ADDRESS_TAG, item.UserAddress)
		log.SetData(utils.ORD_WALLET_ADDRESS_TAG, item.OrdAddress)
		
		addr := item.OriginUserAddress
		if addr == "" {
			addr = item.UserAddress
		}
	
		sentTokenResp, err := u.SendToken(rootSpan, addr, item.MintResponse.Inscription)
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
		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddrV2(item.OrdAddress, &item)
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateBtcWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		//TODO: - create bitcoin token here
		// _, err = u.CreateBTCTokenURI(span, item.ProjectID, item.MintResponse.Inscription)
		// if err != nil {
		// 	log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.CreateBTCTokenURI.Error", item.OrdAddress), err.Error(), err)
		// 	continue
		// }
	}

	return nil, nil
}

func (u Usecase) SendTokenV2(rootSpan opentracing.Span, receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {
	span, log := u.StartSpan("SendTokenV2", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData(utils.TOKEN_ID_TAG, inscriptionID)
	sendTokenReq := ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_master",
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
		}}

	log.SetData("sendTokenReq", sendTokenReq)
	resp, err := u.OrdService.Exec(sendTokenReq)

	if err != nil {
		log.Error("u.OrdService.Exec", err.Error(), err)
		return nil, err
	}

	log.SetData("sendTokenRes", resp)
	return resp, err

}
