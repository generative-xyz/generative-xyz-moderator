package covalenthq

import (
	"errors"

	json "github.com/json-iterator/go"
	"rederinghub.io/pkg/third-party/crypto/web3/nftdata"
	"rederinghub.io/pkg/utils"

	"time"
)

type (
	Balance struct {
		Data         *Data  `json:"data"`
		Error        bool   `json:"error"`
		ErrorMessage string `json:"error_message"`
		ErrorCode    int    `json:"error_code"`
	}
	Data struct {
		Address       string      `json:"address,omitempty"`
		UpdatedAt     time.Time   `json:"updated_at,omitempty"`
		NextUpdateAt  time.Time   `json:"next_update_at,omitempty"`
		QuoteCurrency string      `json:"quote_currency,omitempty"`
		ChainId       int         `json:"chain_id,omitempty"`
		Items         []*Item     `json:"items,omitempty"`
		Pagination    interface{} `json:"pagination,omitempty"`
	}

	Item struct {
		ContractDecimals     int            `json:"contract_decimals,omitempty"`
		ContractName         string         `json:"contract_name,omitempty"`
		ContractTickerSymbol string         `json:"contract_ticker_symbol,omitempty"`
		ContractAddress      string         `json:"contract_address,omitempty"`
		SupportsErc          []string       `json:"supports_erc,omitempty"`
		LogoUrl              string         `json:"logo_url,omitempty"`
		LastTransferredAt    *time.Time     `json:"last_transferred_at,omitempty"`
		Type                 string         `json:"type,omitempty"`
		Balance              string         `json:"balance,omitempty"`
		Balance24H           *string        `json:"balance_24h,omitempty"`
		QuoteRate            float64        `json:"quote_rate,omitempty"`
		QuoteRate24H         *float64       `json:"quote_rate_24h,omitempty"`
		Quote                float64        `json:"quote,omitempty"`
		Quote24H             *float64       `json:"quote_24h,omitempty"`
		NftDataInterface     interface{}    `json:"nft_data,omitempty"`
		NFTData              nftdata.Client `json:"-"`
	}
)

func NewNFTData(data interface{}, chanID int) (nftdata.Client, error) {
	// include interface to map for unmarshalling
	nftDataInterface := map[string]interface{}{
		"items": data,
	}
	switch chanID {
	case chainIDSolana:
		solonaNFTData := &Solana{}
		err := utils.InterfaceToStruct(nftDataInterface, solonaNFTData)
		if err != nil {
			return nil, err
		}
		return solonaNFTData, nil
	}

	// default for metamask
	metamaskNFTData := &Metamask{}
	err := utils.InterfaceToStruct(nftDataInterface, metamaskNFTData)
	if err != nil {
		return nil, err
	}
	return metamaskNFTData, nil
}

func (a *Balance) UnmarshalJSON(data []byte) error {
	type Alias Balance
	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}
	*a = Balance(alias)

	if a.Error == true {
		return errors.New(a.ErrorMessage)
	}

	if a.Data == nil {
		return errors.New("no items found")
	}

	for i := range a.Data.Items {
		item := a.Data.Items[i]
		if item.NftDataInterface == nil {
			continue
		}

		// only unmarshal if it's a nft
		if item.Type != typeNFT {
			continue
		}

		nftData, err := NewNFTData(item.NftDataInterface, a.Data.ChainId)
		if err != nil {
			return err
		}
		item.NFTData = nftData
		item.NftDataInterface = nil
		a.Data.Items[i] = item
	}

	return nil
}

// HasNFT returns true if the balance has NFTs
func (a *Balance) HasNFT() bool {
	return a.Data != nil
}

func (a *Balance) GetNFTCustomerInfo() *nftdata.NFTCustomerInfo {
	nftItem := a.GetFirstNFTItem()
	if nftItem == nil {
		return nil
	}

	return nftdata.NewNFTCustomerInfo(nftItem, a.Data.Address)
}

func (a *Balance) GetFirstNFTItem() *nftdata.Item {
	if a == nil || a.Data == nil {
		return nil
	}

	if len(a.Data.Items) == 0 {
		return nftdata.DefaultItem
	}

	for i := range a.Data.Items {
		item := a.Data.Items[i]
		if item.NFTData == nil {
			continue
		}

		return item.NFTData.GetFirstNFTItem()
	}

	// get default item in case no nft data is found
	return nftdata.DefaultItem
}

// GetNFTItems returns all nft items
func (a *Balance) GetNFTItems() []interface{} {
	if a == nil || a.Data == nil {
		return nil
	}

	if len(a.Data.Items) == 0 {
		return nil
	}

	var items []interface{}
	for i := range a.Data.Items {
		item := a.Data.Items[i]
		if item.NFTData == nil {
			continue
		}

		items = append(items, item.NFTData.GetData())
	}

	return items
}
