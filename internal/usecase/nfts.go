package usecase

import (
	"errors"
	"math/big"
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

func (u Usecase) GetAllTokenHolder(rootSpan opentracing.Span) ([]structure.TokenHolder, error) {
	span, log := u.StartSpan("GetAllTokenHolder", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	covalentResps, err := u.CovalentNft.GetAllTokenHolder(nftStructure.CovalentGetAllTokenHolderRequest{
		ContractAddress: u.Config.GENToken.Contract,
		Limit:           100,
	})
	if err != nil {
		return nil, err
	}

	tokenHolders := make([]structure.TokenHolder, 0)
	for _, resp := range covalentResps {
		for _, item := range resp.Data.Items {
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
		}
	}

	return tokenHolders, nil
}

func (u Usecase) GetTokenHolders(rootSpan opentracing.Span, req structure.GetTokenHolderRequest) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetTokenHolders", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req)

	bf := entity.BaseFilters{
		Page: int64(req.Page),
		Limit: int64(req.Limit),
		SortBy: "current_rank",
		Sort: entity.SORT_ASC,
	}

	resp, err := u.Repo.FilterTokenHolders(entity.FilterTokenHolders{
		BaseFilters: bf,
	})

	return resp, err
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
