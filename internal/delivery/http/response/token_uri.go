package response

import "time"

type TokenURIResp struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	AnimationURL string `json:"animation_url"`
	Attributes interface{} `json:"attributes"`
}

type TokenTraitsResp struct{
	Attributes interface{} `json:"attributes"`
}


type InternalTokenURIResp struct{
	BaseResponse
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	AnimationURL string `json:"animationUrl"`
	Attributes interface{} `json:"attributes"`
	MintedTime time.Time `json:"mintedTime"`
	GenNFTAddr string `json:"genNFTAddr"`
	OwnerAddr string `json:"ownerAddr"`
	Owner *ProfileResponse `json:"owner"`
	Project *ProjectResp `json:"project"`
	Creator *ProfileResponse `json:"creator"`
}

type InternalTokenTraitsResp struct{
	Attributes interface{} `json:"attributes"`
}
