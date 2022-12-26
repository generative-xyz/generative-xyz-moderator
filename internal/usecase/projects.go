package usecase

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_nft_contract"
)

func (u Usecase) CreateProject(rootSpan opentracing.Span,  req structure.CreateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	err = u.Repo.CreateProject(pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("pe",pe)
	return pe, nil
}

func (u Usecase) GetTokensByContract(rootSpan opentracing.Span,  contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	chainURL := os.Getenv("CHAIN_URL")
	log.SetData("chainURL", chainURL)
	// call to contract to get emotion
	client, err := ethclient.Dial(chainURL)
	if err != nil {
		log.Error("ethclient.Dial", err.Error(), err)
		return nil, err
	}

	contractAddr :=  common.HexToAddress(contractAddress)
	gNft, err := generative_nft_contract.NewGenerativeNftContract(contractAddr, client)
	if err != nil {
		log.Error("generative_nft_contract.NewGenerativeNftContract", err.Error(), err)
		return nil, err
	}

	project, err := gNft.Project(nil)
	if err != nil {
		log.Error("gNft.Project", err.Error(), err)
		return nil, err
	}
	parentAddr := project.ProjectAddr

	resp, err := u.MoralisNft.GetNftByContract(contractAddress, filter)
	if err != nil {
		log.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
		return nil, err
	}
	parentAddrStr :=  parentAddr.String()
	result := []entity.TokenUri{}
	for _, item := range resp.Result {
		tokenID := item.TokenID
		token, err := u.GetToken(span, structure.GetTokenMessageReq{ContractAddress: parentAddrStr, TokenID: tokenID })
		if err != nil {
			log.Error("u.getTokenInfo", err.Error(), err)
			return nil, err
		}
		result = append(result, *token)
	}
	
	p := &entity.Pagination{}
	p.Result = result
	p.Currsor = resp.Cursor
	p.Total = int64(resp.Total)
	p.Page = int64(resp.Page)
	p.PageSize = int64(resp.PageSize)
	return p, nil
}

func (u Usecase) UpdateProject(rootSpan opentracing.Span,  req structure.UpdateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("UpdateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	resp := &entity.Projects{}
	

	return resp, nil
}

func (u Usecase) GetProjects(rootSpan opentracing.Span,  req structure.FilterProjects) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	pe := &entity.FilterProjects{}
	err := copier.Copy(pe, req)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	projects, err := u.Repo.GetProjects(*pe)
	if err != nil {
		log.Error("u.Repo.CreateProject", err.Error(), err)
		return nil, err
	}

	log.SetData("projects",projects)
	return projects, nil
}

func (u Usecase) GetProjectDetail(rootSpan opentracing.Span,  req structure.GetProjectDetailMessageReq) (*entity.Projects, error) {
	span, log := u.StartSpan("GetProjectDetail", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	//alway update project in a separated process
	go func (rootSpan opentracing.Span)  {
		span, log := u.StartSpan("GetProjectDetail.GetProjectFromChain", rootSpan)
		defer u.Tracer.FinishSpan(span, log )

		_, err :=  u.UpdateProjectFromChain(req.ContractAddress, req.ProjectID)
		if err != nil {
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return
		}

	}(span)
	
	c, _ := u.Repo.FindProjectBy(req.ContractAddress, req.ProjectID) 
	if  (c == nil) || (c != nil && !c.IsSynced ) {
		p, err :=  u.UpdateProjectFromChain(req.ContractAddress, req.ProjectID)
		if err != nil {
			log.Error("u.Repo.FindProjectBy", err.Error(), err)
			return nil, err
		}
		return p, nil
	}
	return c, nil
}