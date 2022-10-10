package covalenthq

import "rederinghub.io/pkg/third-party/crypto/web3/nftdata"

type (
	SolanaItem struct {
		Key             int    `json:"key"`
		UpdateAuthority string `json:"updateAuthority"`
		Mint            string `json:"mint"`
		Data            struct {
			Name                 string `json:"name"`
			Symbol               string `json:"symbol"`
			Uri                  string `json:"uri"`
			SellerFeeBasisPoints int    `json:"sellerFeeBasisPoints"`
			Creators             []struct {
				Address  string `json:"address"`
				Verified int    `json:"verified"`
				Share    int    `json:"share"`
			} `json:"creators"`
		} `json:"data"`
	}
	Solana struct {
		Items []*SolanaItem `json:"items"`
	}
)

func (s Solana) GetFirstNFTItem() *nftdata.Item {
	if len(s.Items) == 0 {
		return nil
	}

	nftItem := s.Items[0]
	return &nftdata.Item{
		Domain:      nftItem.GetNFTDomain(),
		CompanyName: nftItem.GetNFTCompanyName(),
	}
}

func (s *Solana) GetData() interface{} {
	return s
}

func (a *SolanaItem) GetNFTDomain() string {
	return nftdata.GetDomainFromURI(a.Data.Uri)
}

func (a *SolanaItem) GetNFTCompanyName() string {
	if a.Data.Name == "" {
		return nftdata.CompanyNameDefault
	}

	return a.Data.Name
}
