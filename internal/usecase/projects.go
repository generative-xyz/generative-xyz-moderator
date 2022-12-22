package usecase

import (
	"github.com/opentracing/opentracing-go"

	"rederinghub.io/external/nfts"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateProject(rootSpan opentracing.Span,  req structure.CreateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	resp := &entity.Projects{}
	

	return resp, nil
}


func (u Usecase) GetTokensByContract(rootSpan opentracing.Span,  contractAddress string, filter nfts.MoralisFilter) (*entity.Pagination, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	resp, err := u.MoralisNft.GetNftByContract(contractAddress, filter)
	if err != nil {
		log.Error("u.MoralisNft.GetNftByContract", err.Error(), err)
		return nil, err
	}

	result := []structure.ProjectDetail{}
	for _, item := range resp.Result {
		parsed, err := u.GetProjectDetail(span, structure.GetProjectDetailMessageReq{ContractAddress: item.TokenAddress, ProjectID: item.TokenID })
		if err != nil {
			log.Error("u.MoralisNft.GetProjectDetail", err.Error(), err)
			return nil, err
		}
		result = append(result, *parsed)
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
