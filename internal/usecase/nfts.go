package usecase

import (
	"errors"
	"math/big"
	"time"

	"github.com/opentracing/opentracing-go"
	nftStructure "rederinghub.io/external/nfts"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetNftTransactions(rootSpan opentracing.Span, req structure.GetNftTransactionsReq) (*nftStructure.CovalentGetNftTransactionResponse, error) {
	span, log := u.StartSpan("GetNftTransactions", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req);
	resp, err := u.CovalentNft.GetNftTransactions(nftStructure.CovalentNftTransactionFilter{
		Chain : req.Chain,
		ContractAddress: req.ContractAddress,
		TokenID: req.TokenID,
	})
	return resp, err
}

func (u Usecase) GetNftMintedTime(rootSpan opentracing.Span, req structure.GetNftMintedTimeReq) (*structure.NftMintedTime, error) {
	span, log := u.StartSpan("GetNftMintedTime", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	log.SetData("req", req);
	
	// try to get block number minted and minted time from moralis
	nft, err := u.MoralisNft.GetNftByContractAndTokenID(req.ContractAddress, req.TokenID)
	if err != nil {
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
		MintedTime: &mintedTime,
	}, nil
}
