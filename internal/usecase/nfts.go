package usecase

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	nftStructure "rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_nft_contract"
	"time"
)

func (u Usecase) GetNftTransactions(req structure.GetNftTransactionsReq) (*nftStructure.CovalentGetNftTransactionResponse, error) {

	resp, err := u.CovalentNft.GetNftTransactions(nftStructure.CovalentNftTransactionFilter{
		Chain:           req.Chain,
		ContractAddress: req.ContractAddress,
		TokenID:         req.TokenID,
	})
	return resp, err
}

func (u Usecase) GetAllTokenHolder() ([]structure.TokenHolder, error) {
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

func (u Usecase) GetTokenHolders(req structure.GetTokenHolderRequest) (*entity.Pagination, error) {
	bf := entity.BaseFilters{
		Page:   int64(req.Page),
		Limit:  int64(req.Limit),
		SortBy: "current_rank",
		Sort:   entity.SORT_ASC,
	}

	resp, err := u.Repo.FilterTokenHolders(entity.FilterTokenHolders{
		BaseFilters: bf,
	})

	return resp, err
}

func (u Usecase) GetNftMintedTime(client *ethclient.Client, req structure.GetNftMintedTimeReq) (*structure.NftMintedTime, error) {
	// try to get block number minted and minted time from moralis
	/*nft, err := u.MoralisNft.GetNftByContractAndTokenID(req.ContractAddress, req.TokenID)
	if err != nil {
		logger.AtLog.Logger.Error("err", zap.Error(err))
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
	}, nil*/

	contractAddr := common.HexToAddress(req.ContractAddress)
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(req.TokenID, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	value, err := gNft.OwnerOf(nil, tokenID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	return &structure.NftMintedTime{
		BlockNumberMinted: nil,
		MintedTime:        &now,
		Nft: &nftStructure.MoralisToken{
			TokenID: req.TokenID,
			Owner:   value.String(),
		},
	}, nil
}

func (u Usecase) ownerOf(client *ethclient.Client, contractAddr common.Address, tokenIDStr string) (*common.Address, error) {
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		return nil, err
	}
	tokenID := new(big.Int)
	tokenID, ok := tokenID.SetString(tokenIDStr, 10)
	if !ok {
		return nil, errors.New("cannot convert tokenID")
	}
	value, err := gNft.OwnerOf(nil, tokenID)
	if err != nil {
		return nil, err
	}
	return &value, err
}
