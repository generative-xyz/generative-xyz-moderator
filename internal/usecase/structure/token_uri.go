package structure

import "rederinghub.io/utils/contracts/generative_project_contract"


type GetTokenMessageReq struct {
	ContractAddress string
	TokenID string
}

type GetProjectDetailMessageReq struct {
	ContractAddress string
	ProjectID string
}

type GetTokenMessageResp struct {
	ContractAddress string
	TokenID string
}

type ProjectDetail struct {
	ProjectDetail *generative_project_contract.NFTProjectProject
	Status bool
	NftTokenUri string
}

type ProjectDetailChan struct {
	ProjectDetail *generative_project_contract.NFTProjectProject
	Err error
}

type ProjectStatusChan struct {
	Status *bool
	Err error
}

type ProjectNftTokenUriChan struct {
	TokenURI *string
	Err error
}

