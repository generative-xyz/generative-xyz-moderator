package covalenthq

import "rederinghub.io/pkg/third-party/crypto/web3/nftdata"

type (
	MetamaskItem struct {
		TokenId           string      `json:"token_id"`
		TokenBalance      interface{} `json:"token_balance"`
		TokenUrl          string      `json:"token_url"`
		SupportsErc       []string    `json:"supports_erc"`
		TokenPriceWei     interface{} `json:"token_price_wei"`
		TokenQuoteRateEth interface{} `json:"token_quote_rate_eth"`
		OriginalOwner     string      `json:"original_owner"`
		ExternalData      struct {
			Name         string      `json:"name"`
			Description  string      `json:"description"`
			Image        string      `json:"image"`
			Image256     string      `json:"image_256"`
			Image512     string      `json:"image_512"`
			Image1024    string      `json:"image_1024"`
			AnimationUrl interface{} `json:"animation_url"`
			ExternalUrl  string      `json:"external_url"`
			Attributes   []struct {
				TraitType string `json:"trait_type"`
				Value     string `json:"value"`
			} `json:"attributes"`
			Owner string `json:"owner"`
		} `json:"external_data"`
		Owner        string      `json:"owner"`
		OwnerAddress interface{} `json:"owner_address"`
		Burned       interface{} `json:"burned"`
	}
	Metamask struct {
		Items []*MetamaskItem `json:"items"`
	}
)

func (s Metamask) GetFirstNFTItem() *nftdata.Item {
	if len(s.Items) == 0 {
		return nil
	}

	nftItem := s.Items[0]
	return &nftdata.Item{
		Domain:      nftItem.GetNFTDomain(),
		CompanyName: nftItem.GetNFTCompanyName(),
	}
}

func (s *Metamask) GetData() interface{} {
	return s
}

func (a *MetamaskItem) GetNFTCompanyName() string {
	if a.ExternalData.Name == "" {
		return nftdata.CompanyNameDefault
	}

	return a.ExternalData.Name
}

func (a *MetamaskItem) GetNFTDomain() string {
	url := a.TokenUrl
	if url == "" {
		url = a.ExternalData.ExternalUrl
	}
	return nftdata.GetDomainFromURI(url)
}
