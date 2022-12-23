package response


type ProjectResp struct{
	BaseResponse
	ContractAddress string `json:"contractAddress"`
	TokenID string `json:"tokenID"`
	MaxSupply int `json:"maxSupply"`
	Limit int `json:"limit"`
	MintPrice string `json:"mintPrice"`
	MintPriceAddr string `json:"mintPriceAddr"`
	Name string `json:"name"`
	Creator string `json:"creator"`
	CreatorAddr string `json:"creatorAddr"`
	License string `json:"license"`
	Desc string `json:"desc"`
	Image string `json:"image"`
	ScriptType []string `json:"scriptType"`
	Social interface{} `json:"social"`
	Scripts []string `json:"scripts"`
	Styles string `json:"styles"`
	CompleteTime int  `json:"completeTime"`
	GenNFTAddr string `json:"genNFTAddr"`
	ItemDesc string `json:"itemDesc"`
	Status bool `json:"status"`
	NftTokenURI string `json:"projectURI"`
}
