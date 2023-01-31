package usecase

import (
	"errors"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	nftStructure "rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetNftTransactions(rootSpan opentracing.Span, req structure.GetNftTransactionsReq) (*nftStructure.CovalentGetNftTransactionResponse, error) {
	span, log := u.StartSpan("GetNftTransactions", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)
	resp, err := u.CovalentNft.GetNftTransactions(nftStructure.CovalentNftTransactionFilter{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		TokenID:         req.TokenID,
	})
	return resp, err
}

func (u Usecase) GetTokenHolders(rootSpan opentracing.Span, req structure.GetTokenHolderRequest) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetTokenHolders", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)
	covalentResp, err := u.CovalentNft.GetTokenHolder(nftStructure.CovalentGetTokenHolderRequest{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		Page:            req.Page,
		Limit:           req.Limit,
	})
	if err != nil {
		return nil, err
	}

	tokenHolders := make([]structure.TokenHolder, 0)

	getProfile := func(profileChan chan structure.ProfileChan, address string) {
		var user *entity.Users
		var err error

		defer func() {
			profileChan <- structure.ProfileChan{
				Data: user,
				Err: err,
			}
		}()

		user, err = u.GetUserProfileByWalletAddress(span, strings.ToLower(address))
		if err != nil {
			return
		}
	}

	profileChans := make([]chan structure.ProfileChan, 0)

	for _, item := range covalentResp.Data.Items {
		tokenHolders = append(tokenHolders, structure.TokenHolder{
			ContractDecimals:     item.ContractDecimals,
			ContractName:         item.ContractName,
			ContractTickerSymbol: item.ContractTickerSymbol,
			ContractAddress:      item.ContractAddress,
			SupportsErc:          item.SupportsErc,
			LogoURL:              item.LogoURL,
			Address:              item.Address,
			Balance:              item.Balance,
			TotalSupply:          item.TotalSupply,
			BlockHeight:          item.BlockHeight,
		})
		profileChan := make(chan structure.ProfileChan, 1)
		go getProfile(profileChan, item.Address)
		profileChans = append(profileChans, profileChan)
	}

	for i := 0; i < len(tokenHolders); i++ {
		usrFromChan := <-profileChans[i]
		tokenHolders[i].Profile = usrFromChan.Data
	}

	sort.SliceStable(tokenHolders, func(i, j int) bool {
		lhs := new(big.Int)
		lhs, ok := lhs.SetString(tokenHolders[i].Balance, 10)
		if !ok {
			lhs = big.NewInt(0)
		}
		rhs := new(big.Int)
		rhs, ok = rhs.SetString(tokenHolders[j].Balance, 10)
		if !ok {
			rhs = big.NewInt(0)
		}
		return lhs.Cmp(rhs) < 0
	})
	total, ok := covalentResp.Data.Pagination.TotalCount.(int)
	if !ok {
		total = len(tokenHolders)
	}
	resp := &entity.Pagination{
		Page:   int64(req.Page),
		PageSize:  int64(req.Limit),
		Result: tokenHolders,
		Total: int64(total),
	}

	return resp, nil
}

func (u Usecase) GetNftMintedTime(rootSpan opentracing.Span, req structure.GetNftMintedTimeReq) (*structure.NftMintedTime, error) {
	span, log := u.StartSpan("GetNftMintedTime", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)
	log.SetTag("tokenID", req.TokenID)
	log.SetTag("contractAddress", req.ContractAddress)

	// try to get block number minted and minted time from moralis
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(req.ContractAddress, req.TokenID)
	if err != nil {
		log.Error("u.GetNftMintedTime.MoralisNft.GetNftByContractAndTokenID", err.Error(), err)
		return nil, err
	}
	blockNumber := nft.BlockNumberMinted
	blockNumberBigInt := new(big.Int)
	blockNumberBigInt, ok := blockNumberBigInt.SetString(blockNumber, 10)
	if !ok {
		return nil, errors.New("cannot convert blockNumber to bigint")
	}
	// get time by block number
	block, err := u.Blockchain.GetBlockByNumber(*blockNumberBigInt)
	if err != nil {
		return nil, err
	}
	// get time from block
	mintedTime := time.Unix(int64(block.Time()), 0)
	return &structure.NftMintedTime{
		BlockNumberMinted: &blockNumber,
		MintedTime:        &mintedTime,
		Nft:               nft,
	}, nil
}
