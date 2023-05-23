package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/repository"
	erc202 "rederinghub.io/utils/contracts/erc20"
	"sort"
	"strings"
	"time"
)

type respStruct struct {
	Time             time.Time `json:"time"`
	Timestamp        int       `json:"timestamp"`
	VolumeFrom       string    `json:"volume_from"`
	VolumeTo         string    `json:"volume_to"`
	BtcPrice         float64   `json:"btc_price"`
	UsdPrice         string    `json:"usd_price"`
	Low              string    `json:"low"`
	Open             string    `json:"open"`
	Close            string    `json:"close"`
	High             string    `json:"high"`
	VolumeFromUsd    string    `json:"volume_from_usd"`
	VolumeToUsd      string    `json:"volume_to_usd"`
	TotalVolumeUsd   string    `json:"total_volume_usd"`
	TotalVolume      float64   `json:"total_volume"`
	LowUsd           string    `json:"low_usd"`
	OpenUsd          string    `json:"open_usd"`
	CloseUsd         string    `json:"close_usd"`
	HighUsd          string    `json:"high_usd"`
	ConversionType   string    `json:"conversion_type"`
	ConversionSymbol string    `json:"conversion_symbol"`
}

func (u Usecase) GetNftsByAddress(address string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/nft-explorer/owner-address/%s/nfts", address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var result struct {
		Data []*struct {
			Collection        string `json:"collection"`
			CollectionAddress string `json:"collection_address"`
			TokenID           string `json:"token_id"`
			Name              string `json:"name"`
			ContentType       string `json:"content_type"`
			Image             string `json:"image"`
			Explorer          string `json:"explorer"`
			ArtistName        string `json:"artist_name"`
		} `json:"data"`
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// var listContractID []string

	for _, nft := range result.Data {

		// listContractID = append(listContractID, nft.CollectionAddress)

		// if len(nft.Image) > 0 {
		// 	nft.Image += "/content"
		// }
		nft.Explorer = fmt.Sprintf("https://trustless.computer/inscription?contract=%s&id=%s", nft.CollectionAddress, nft.TokenID)

		p, _ := u.Repo.FindProjectByGenNFTAddr(nft.CollectionAddress)
		if p != nil {
			nft.ArtistName = p.Name
		}

	}
	return result.Data, err
}

func (u Usecase) GetNftsByAddressFromTokenUri(address string) (interface{}, error) {

	type Data struct {
		Collection        string               `json:"collection"`
		CollectionAddress string               `json:"collection_address"`
		TokenID           string               `json:"token_id"`
		ProjectID         string               `json:"project_id"`
		ProjectName       string               `json:"project_name"`
		TokenNumber       *int                 `json:"token_number"`
		Name              string               `json:"name"`
		ContentType       string               `json:"content_type"`
		Image             string               `json:"image"`
		Explorer          string               `json:"explorer"`
		Buyable           bool                 `json:"buyable"`
		ArtistName        string               `json:"artist_name"`
		GenNftAddress     string               `json:"gen_nft_addrress"`
		Royalty           int                  `json:"royalty"`
		PriceBRC20        entity.PriceBRC20Obj `json:"priceBrc20"`

		//Used for big file
		MintingInfo repository.AggregateTokenMintingInfo `json:"minting_info"`
		IsMinting   bool                                 `json:"is_minting"`
	}

	ctx := context.Background()
	var dataList []*Data
	listToken, _ := u.Repo.GetOwnerTokens(address)

	fmt.Println("len(listToken) > 0", len(listToken) > 0)

	if len(listToken) > 0 {
		for _, nft := range listToken {

			data := &Data{
				Collection:        "",
				CollectionAddress: nft.ContractAddress,
				TokenID:           nft.TokenID,
				TokenNumber:       &nft.TokenIDMini,
				ProjectID:         nft.ProjectID,
				Name:              nft.Name,
				Image:             nft.Thumbnail,
				Explorer:          fmt.Sprintf("https://trustless.computer/inscription?contract=%s&id=%s", nft.ContractAddress, nft.TokenID),
				ArtistName:        nft.CreatorName,
				ProjectName:       nft.ProjectName,
				GenNftAddress:     nft.GenNFTAddr,
				Buyable:           nft.Buyable,
				Royalty:           int(nft.Royalty),
				PriceBRC20: entity.PriceBRC20Obj{
					Value:      nft.PriceBRC20,
					Address:    nft.PriceBRC20Address,
					OfferingID: nft.OfferingID,
				},
				IsMinting: false,
				MintingInfo: repository.AggregateTokenMintingInfo{
					All:     0,
					Done:    0,
					Pending: 0,
				},
			}

			mintingInfo, err := u.Repo.AggregateMintingInfo(ctx, nft.TokenID)
			if err == nil {
				if len(mintingInfo) >= 1 {
					mtinfo := mintingInfo[0]
					data.MintingInfo = mtinfo
					if mtinfo.Done < mtinfo.All {
						data.IsMinting = true
					}
				}
			}

			dataList = append(dataList, data)
		}
	}
	return dataList, nil

}

func (u Usecase) GetTokenSwapChartMinMax(address string, chartType string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/swap/token/price?chart_type=%s&contract_address=%s", chartType, address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Err    error         `json:"error"`
		Status bool          `json:"status"`
		Data   []*respStruct `json:"data"`
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.Status && result.Err != nil {
		return nil, result.Err
	}

	sort.SliceStable(result.Data, func(i, j int) bool {
		closedUsdI := result.Data[i].CloseUsd
		closedUsdJ := result.Data[j].CloseUsd

		closedUsdIBig, _ := new(big.Float).SetString(closedUsdI)
		closedUsdJBig, _ := new(big.Float).SetString(closedUsdJ)
		cmp := closedUsdIBig.Cmp(closedUsdJBig)
		if cmp >= 0 {
			return true
		}

		return false
	})

	max := result.Data[0]
	min := result.Data[len(result.Data)-1]

	type respResult struct {
		Min *respStruct `json:"min"`
		Max *respStruct `json:"max"`
	}

	return respResult{
		Min: min,
		Max: max,
	}, err
}

func (u Usecase) GetTokenSwapChart(address string, chartType string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/swap/token/price?chart_type=%s&contract_address=%s", chartType, address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Err    error         `json:"error"`
		Status bool          `json:"status"`
		Data   []*respStruct `json:"data"`
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.Status && result.Err != nil {
		return nil, result.Err
	}

	sort.SliceStable(result.Data, func(i, j int) bool {
		closedUsdI := result.Data[i].CloseUsd
		closedUsdJ := result.Data[j].CloseUsd

		closedUsdIBig, _ := new(big.Float).SetString(closedUsdI)
		closedUsdJBig, _ := new(big.Float).SetString(closedUsdJ)
		cmp := closedUsdIBig.Cmp(closedUsdJBig)
		if cmp >= 0 {
			return true
		}

		return false
	})

	max := result.Data[0]
	min := result.Data[len(result.Data)-1]

	type respResult struct {
		Min *respStruct `json:"min"`
		Max *respStruct `json:"max"`
	}

	return respResult{
		Min: min,
		Max: max,
	}, err
}

func (u Usecase) GetTokenReport(address string) (interface{}, error) {
	url := fmt.Sprintf("https://dapp.trustless.computer/dapp/api/swap/token/report?address=%s", address)
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Err    error `json:"error"`
		Status bool  `json:"status"`
		Data   []*struct {
			Address           string `json:"address"`
			TotalSupply       string `json:"total_supply"`
			TotalSupplyNumber string `json:"total_supply_number"`
			Owner             string `json:"owner"`
			Decimal           int    `json:"decimal"`
			DeployedAtBlock   int    `json:"deployed_at_block"`
			Slug              string `json:"slug"`
			Symbol            string `json:"symbol"`
			Name              string `json:"name"`
			Thumbnail         string `json:"thumbnail"`
			Description       string `json:"description"`
			Social            struct {
				Website   string `json:"website"`
				Discord   string `json:"discord"`
				Twitter   string `json:"twitter"`
				Telegram  string `json:"telegram"`
				Medium    string `json:"medium"`
				Instagram string `json:"instagram"`
			} `json:"social"`
			Index           int     `json:"index"`
			Volume          string  `json:"volume"`
			TotalVolume     string  `json:"total_volume"`
			BtcVolume       float64 `json:"btc_volume"`
			UsdVolume       float64 `json:"usd_volume"`
			BtcTotalVolume  float64 `json:"btc_total_volume"`
			UsdTotalVolume  float64 `json:"usd_total_volume"`
			MarketCap       string  `json:"market_cap"`
			UsdMarketCap    float64 `json:"usd_market_cap"`
			Price           string  `json:"price"`
			BtcPrice        float64 `json:"btc_price"`
			UsdPrice        float64 `json:"usd_price"`
			Percent         string  `json:"percent"`
			Percent7Day     string  `json:"percent_7day"`
			Network         string  `json:"network"`
			Priority        int     `json:"priority"`
			BaseTokenSymbol string  `json:"base_token_symbol"`
			Status          string  `json:"status"`
		} `json:"data"`
	}

	// parse:
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if !result.Status && result.Err != nil {
		return nil, result.Err
	}

	return result.Data, err
}

func (u Usecase) tokenBalance(userWalletAddress string, tokenAddress string) (*big.Int, int, error) {
	client := u.TcClientPublicNode.GetClient()

	hexAddress := common.HexToAddress(tokenAddress)
	erc20, err := erc202.NewErc20(hexAddress, client)
	if err != nil {
		return nil, 0, err
	}

	blance, err := erc20.BalanceOf(nil, common.HexToAddress(userWalletAddress))
	if err != nil {
		return nil, 0, err
	}

	decimal := 18
	pow := math.Pow10(decimal)
	powBig := big.NewInt(int64(pow))

	blanceNew := blance.Quo(blance, powBig)
	return blanceNew, decimal, nil
}

func (u Usecase) TokenHolderBalance(userWalletAddress string, tokenAddress string) (interface{}, error) {
	balanceNew, _, err := u.tokenBalance(userWalletAddress, tokenAddress)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", balanceNew.Uint64()), nil
}

func (u Usecase) SoralisSnapShotUserTokenBalance(userWalletAddress string, tokenAddress string) (interface{}, error) {
	balanceNew, decimal, err := u.tokenBalance(userWalletAddress, tokenAddress)
	if err != nil {
		return nil, err
	}

	obj := &entity.SoralisSnapShotBalance{
		Balance:       fmt.Sprintf("%d", balanceNew.Uint64()),
		Decimal:       decimal,
		WalletAddress: strings.ToLower(userWalletAddress),
		TokenAddress:  strings.ToLower(tokenAddress),
	}

	err = u.Repo.InsertSnapshot(obj)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("%d", balanceNew.Uint64()), nil
}

func (u Usecase) SoralisGetSnapShotUserTokenBalance(userWalletAddress string, tokenAddress string) ([]entity.FilteredSoralisSnapShotBalance, error) {
	data, err := u.Repo.GetSnapshotByWalletAddress(userWalletAddress, tokenAddress)
	if err != nil {
		return nil, err
	}

	return data, nil
}
