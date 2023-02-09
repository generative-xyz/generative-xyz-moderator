package usecase

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/external/ord_service"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
)

func (u Usecase) CreateBTCWalletAddress(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("CreateConfig", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("input", input)
	log.SetTag("btcUserWallet", input.WalletAddress)

	walletAddress := &entity.BTCWalletAddress{}
	err := copier.Copy(walletAddress, input)
	if err != nil {
		log.Error("u.CreateBTCWalletAddress.Copy", err.Error(), err)
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
	}else{
		walletAddress.Mnemonic = resp.Stdout
	}
	

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

	p, err := u.Repo.FindProjectByTokenID(input.ProjectID )
	if err != nil {
		log.Error("u.CreateBTCWalletAddress.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	log.SetData("found.Project", p)

	walletAddress.Amount = p.MintPrice
	walletAddress.UserAddress = input.WalletAddress
	walletAddress.OrdAddress = strings.ReplaceAll(resp.Stdout, "\n", "")
	walletAddress.IsConfirm = false
	walletAddress.IsMinted = false
	walletAddress.FileURI = "" //find files from google store
	walletAddress.InscriptionID = "" //find files from google store
	walletAddress.ProjectID = input.ProjectID 

	err = u.Repo.InsertBtcWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.CreateBTCWalletAddress.InsertBtcWalletAddress", err.Error(), err)
		return nil, err
	}

	return walletAddress, nil
}

func (u Usecase) BTCMint(rootSpan opentracing.Span, input structure.BctMintData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("BTCMint", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	log.SetData("input", input)
	log.SetTag("ordWalletaddress", input.Address)

	btc, err := u.Repo.FindBtcWalletAddressByOrd(input.Address)
	if err != nil {
		log.Error("BTCMint.FindBtcWalletAddressByOrd", err.Error(), err)
		return nil, err
	}

	//mint logic
	btc, err = u.MintLogic(span, btc)
	if err != nil {
		log.Error("BTCMint.MintLogic", err.Error(), err)
		return nil, err
	}

	// get data from project
	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
	if err != nil {
		log.Error("BTCMint.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	log.SetData("found.Project", p)
	log.SetTag("projectID", p.TokenID)

	//prepare data for mint
	// - Get project.AnimationURL
	projectNftTokenUri := &structure.ProjectAnimationUrl{}
	err = helpers.Base64DecodeRaw(p.NftTokenUri, projectNftTokenUri)
	if err != nil {
		log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	// - Upload the Animation URL to GCS
	animation := projectNftTokenUri.AnimationUrl
	animation = strings.ReplaceAll(animation,"data:text/html;base64,", "")

	now := time.Now().UTC().Unix()
	uploaded, err := u.GCS.UploadBaseToBucket(animation, fmt.Sprintf("btc-projects/%s/%d.html",p.TokenID, now))
	if err != nil {
		log.Error("BTCMint.helpers.Base64DecodeRaw", err.Error(), err)
		return nil, err
	}

	fileURI := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	btc.FileURI = fileURI
	spew.Dump(fileURI)

	//TODO - enable this
	resp, err := u.OrdService.Mint(ord_service.MintRequest{
		WalletName: "ord_master",
		FileUrl: fileURI,
		FeeRate: 7,//temp
		DryRun: true,
	})

	btc.MintResponse = entity.MintStdoputResponse(resp.Stdout)
	btc, err = u.UpdateBtcMintedStatus(span, btc)
	if err != nil {
		log.Error("BTCMint.UpdateBtcMintedStatus", err.Error(), err)
		return nil, err
	}

	return btc, nil
}

func (u Usecase) ReadGCSFolder(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("ReadGCSFolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}

func (u Usecase) UpdateBtcMintedStatus(rootSpan opentracing.Span, btcWallet *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("UpdateBtcMintedStatus", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("input", btcWallet)
	btcWallet.IsMinted = true

	updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(btcWallet.OrdAddress, btcWallet)
	if err != nil {
		log.Error("BTCMint.helpers.UpdateBtcWalletAddressByOrdAddr", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return btcWallet, nil
}

func (u Usecase) BalanceLogic(rootSpan opentracing.Span, btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("BalanceLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			userWallet,
			"wallet",
			"balance",
		},
	})

	if err != nil {
		log.Error("BTCMint.Exec.balance", err.Error(), err)
		return nil, err
	}

	log.SetData("resp", resp)
	balance := strings.ReplaceAll(resp.Stdout, "\n", "")
	log.SetData("balance", balance)

	//TODO logic of the checked balance here
	btc.IsConfirm = true
	return btc, nil
}

func (u Usecase) MintLogic(rootSpan opentracing.Span, btc *entity.BTCWalletAddress) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("MintLogic", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	var err error

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("BTCMint.Minted", err.Error(), err)
		return nil, err
	}

	btc, err = u.BalanceLogic(span, btc)
	if err != nil {
		log.Error("MintLogic.BalanceLogic", err.Error(), err)
		return nil, err
	}
	
	if !btc.IsConfirm {
		err = errors.New("This btc must be IsConfirmed")
		log.Error("BTCMint.IsConfirmed", err.Error(), err)
		return  nil, err
	}

	log.SetData("btc", btc)
	return btc, nil
}

func (u Usecase) WillBeProcessWTC(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("WillBeProcessWTC", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	resp := []entity.BTCWalletAddress{}
	page := 1
	limit := 10

	getAddrr :=  func(rootSpan opentracing.Span, page int, limit int) ( []entity.BTCWalletAddress, error) {
		span, log := u.StartSpan("WillBeProcessWTC.GetAddresses", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		log.SetData("page", page)
		log.SetData("limit", limit)
		
		addreses, err := u.Repo.ListProcessingWalletAddress(page, limit)
		if err != nil {
			log.Error("WillBeProcessWTC.ListProcessingWalletAddress", err.Error(), err)
			return nil, err
		}

		iAddreses := addreses.Result 
		ad := iAddreses.([]entity.BTCWalletAddress)
		return ad, nil
	}
	
	for {
		var wg sync.WaitGroup
		resp, err := getAddrr(span, page, limit)
		if err != nil {
			log.Error(fmt.Sprintf("WillBeProcessWTC.page.%d.Error", page), err.Error(), err)
			break
		}
		wg.Add(limit)
		for _, item := range resp {
			go func(wg *sync.WaitGroup, rootSpan opentracing.Span, item *entity.BTCWalletAddress) {
				defer wg.Done()

				item, err := u.BalanceLogic(span, item)
				if err != nil {
					log.Error(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s.Error", item.OrdAddress), err.Error(), err)
					return
				}
				log.SetData(fmt.Sprintf("WillBeProcessWTC.BalanceLogic.%s", item.OrdAddress), item)
		
				updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, item)
				if err != nil {
					log.Error(fmt.Sprintf("WillBeProcessWTC.UpdateBtcWalletAddressByOrdAddr.%s.Error", item.OrdAddress), err.Error(), err)
					return
				}
				log.SetData("updated", updated)

		
				topicName := helpers.CreateMqttTopic(item.OrdAddress)
				log.SetData("topicName", topicName)
				err = u.MqttClient.Publish(topicName, item)
				if err != nil {
					log.Error(fmt.Sprintf("WillBeProcessWTC.Mqtt.%s.Error", item.OrdAddress), err.Error(), err)
					//return
				}
				

			}(&wg, span, &item)
		}

		wg.Wait()
		time.Sleep(2 * time.Minute)
	}

	return resp, nil
}

func (u Usecase) ListenTheMintedBTC(rootSpan opentracing.Span) ([]entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("ListenBTC", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	resp := []entity.BTCWalletAddress{}
	page := 1
	limit := 10

	getAddrr :=  func(rootSpan opentracing.Span, page int, limit int) ( []entity.BTCWalletAddress, error) {
		span, log := u.StartSpan("WillBeProcessWTC.GetAddresses", rootSpan)
		defer u.Tracer.FinishSpan(span, log)

		log.SetData("page", page)
		log.SetData("limit", limit)
		
		addreses, err := u.Repo.ListBTCAddress(page, limit)
		if err != nil {
			log.Error("WillBeProcessWTC.ListBTCAddress", err.Error(), err)
			return nil, err
		}

		iAddreses := addreses.Result 
		ad := iAddreses.([]entity.BTCWalletAddress)
		return ad, nil
	}
	
	for {
		var wg sync.WaitGroup
		resp, err := getAddrr(span, page, limit)
		if err != nil {
			log.Error(fmt.Sprintf("ListenTheMintedBTC.page.%d.Error", page), err.Error(), err)
			break
		}
		wg.Add(limit)
		for _, item := range resp {
			go func(wg *sync.WaitGroup, rootSpan opentracing.Span, item *entity.BTCWalletAddress) {
				defer wg.Done()

				sentTokenResp, err := u.SendToken(item.UserAddress, item.MintResponse.Inscription)
				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.sentToken.%s.Error", item.OrdAddress), err.Error(), err)
					return
				}

				log.SetData(fmt.Sprintf("ListenTheMintedBTC.execResp.%s", item.OrdAddress), sentTokenResp)
				item.MintResponse.IsSent = true
				updated, err := u.Repo.UpdateBtcWalletAddressByOrdAddr(item.OrdAddress, item)
				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.UpdateBtcWalletAddressByOrdAddr.Error", item.OrdAddress), err.Error(), err)
					return
				}
				log.SetData("updated", updated)
				
				fundResp, err := u.OrdService.Exec(
					ord_service.ExecRequest{
						Args: []string{
							"--wallet",
							item.OrdAddress,
							"send",
							"--fee-rate",
							"7",
							"ord_master",
							item.Amount,
						},
					})

				if err != nil {
					log.Error(fmt.Sprintf("ListenTheMintedBTC.%s.ReFund.Error", item.OrdAddress), err.Error(), err)
					return
				}

				log.SetData("fundResp", fundResp)
				
			}(&wg, span, &item)
		}

		wg.Wait()
		time.Sleep(2 * time.Minute)
	}

	return resp, nil
}

func (u Usecase) SendToken(receiveAddr string, inscriptionID string)  (*ord_service.ExecRespose, error) {

	resp, err := u.OrdService.Exec(ord_service.ExecRequest{
		Args: []string{
			"--wallet",
			"ord_master",
			"wallet",
			"send",
			receiveAddr,
			inscriptionID,
		}})

	if err != nil {
		return nil, err
	}


	return resp, err
	
}