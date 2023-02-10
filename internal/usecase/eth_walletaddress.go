package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
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

	userWallet := helpers.CreateBTCOrdWallet(input.WalletAddress)
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

	log.SetData("CreateETHWalletAddress.createdWallet", resp)
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

	log.SetData("CreateETHWalletAddress.receive", resp)
	p, err := u.Repo.FindProjectByTokenID(input.ProjectID)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}

	log.SetData("found.Project", p.ID)
	walletAddress.Amount = p.MintPrice
	walletAddress.UserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = ""       //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID

	log.SetTag("ordAddress", walletAddress.OrdAddress)
	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateETHWalletAddress.InsertBtcWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) ETHMint(rootSpan opentracing.Span, input structure.EthMintData) (*entity.ETHWalletAddress, error) {
	span, log := u.StartSpan("ETHMint", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("input", input)
	log.SetTag("ordWalletaddress", input.Address)

	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {
		log.Error("ETHMint.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	//mint logic
	btc, err = u.MintLogic(span, btc)
	if err != nil {
		log.Error("ETHMint.MintLogic", err.Error(), err)
		return nil, err
	}

	// get data from project
	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
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
		log.Error("ETHMint.helpers.JsonTransform", err.Error(), err)
		return nil, err
	}

	btc.MintResponse = entity.MintStdoputResponse(*btcMintResp)
	btc, err = u.UpdateBtcMintedStatus(span, btc)
	if err != nil {
		log.Error("ETHMint.UpdateBtcMintedStatus", err.Error(), err)
		return nil, err
	}

	return btc, nil
}

func (u Usecase) ReadGCSFolder(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("ReadGCSFolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateBtcMintedStatus(rootSpan opentracing.Span, btcWallet *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatus", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	log.SetData("input", btcWallet)
	btcWallet.IsMinted = true

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		log.Error("ETHMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return btcWallet, nil
}

func (u Usecase) BalanceLogic(rootSpan opentracing.Span, btc entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("BalanceLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"balance",
		},
	})

	log.SetData("userWallet", btc.UserAddress)
	log.SetData("ordWalletAddress", btc.OrdAddress)
	if err != nil {
		log.Error("ETHMint.Exec.balance", err.Error(), err)
		return nil, err
	}

	log.SetData("resp", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	log.SetData("balance", balance)

	//TODO logic of the checked balance here
	if balance < btc.Amount {
		err := errors.New("Not enough amount")
		return nil, err
	}

	btc.IsConfirm = true
	return &btc, nil
}

func (u Usecase) MintLogic(rootSpan opentracing.Span, btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("MintLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var err error

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("ETHMint.Minted", err.Error(), err)
		return nil, err
	}

	if !btc.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		log.Error("ETHMint.IsConfirmed", err.Error(), err)
		return nil, err
	}

	log.SetData("btc", btc)
	return btc, nil
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
		newItem, err := u.BalanceLogic(span, item)
		if err != nil {
			//log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), newItem)

		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, newItem)
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
			continue
		}
		log.SetData("updated", updated)

		btc, err := u.BTCMint(span, structure.BctMintData{Address: newItem.OrdAddress})
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", newItem.OrdAddress), err.Error(), err)
			continue
		}

		_ = btc
	}

	return nil, nil
}

func (u Usecase) WaitingForMinted(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("WaitingForMinted", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	addreses, err := u.Repo.ListBTCAddress()
	if err != nil {
		log.Error("WillBeProcessWTC.ListBTCAddress", err.Error(), err)
		return nil, err
	}

	for _, item := range addreses {
		log.SetData("userWallet", item.UserAddress)
		log.SetData("ordWalletAddress", item.OrdAddress)
		sentTokenResp, err := u.SendToken(item.UserAddress, item.MintResponse.Inscription)
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
		updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, &item)
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateBtcWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
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

func (u Usecase) SendToken(receiveAddr string, inscriptionID string) (*ord_service.ExecRespose, error) {

	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_master",
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
			"--fee-rate",
			"15",
		}})

	if err != nil {
		return nil, err
	}

	return resp, err

}
