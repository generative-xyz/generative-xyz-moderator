package usecase

import (
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
		return nil, err
	}
	walletAddress.Mnemonic = resp.Stdout

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

	//if this was minted, skip it
	if btc.IsMinted {
		err = errors.New("This btc was minted")
		log.Error("BTCMint.Minted", err.Error(), err)
		return nil, err
	}


	// get data from project
	p, err := u.Repo.FindProjectByTokenID(btc.ProjectID)
	if err != nil {
		log.Error("BTCMint.FindProjectByTokenID", err.Error(), err)
		return nil, err
	}
	log.SetData("found.Project", p)


	//TODO - Check balance
	// userWallet := helpers.CreateBTCOrdWallet(btc.UserAddress)
	// resp, err := u.OrdService.Exec(ord_service.ExecRequest{
	// 	Args: []string{
	// 		"--wallet",
	// 		userWallet,
	// 		"wallet",
	// 		"balance",
	// 	},
	// })

	// if err != nil {
	// 	log.Error("BTCMint.Exec.balance", err.Error(), err)
	// 	return nil, err
	// }
	// log.SetData("resp", resp)

	//TODO logic of the checked balance here


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
	spew.Dump(fileURI)

	//mint
	// resp, err = u.OrdService.Exec(ord_service.ExecRequest{
	// 	Args: []string{
	// 		"--wallet",
	// 		userWallet,
	// 		"wallet",
	// 		"inscriptions",
	// 		"", //FILE URL
	// 	},GCS_DOMAIN
	// })


	return btc, nil
}

func (u Usecase) ReadGCSFolder(rootSpan opentracing.Span, input structure.BctWalletAddressData) (*entity.BTCWalletAddress, error) {
	span, log := u.StartSpan("ReadGCSFolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	log.SetData("input", input)
	u.GCS.ReadFolder("btc-projects/project-1/")
	return nil, nil
}