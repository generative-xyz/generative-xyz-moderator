package usecase

import (
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetChartDataOFProject(req structure.AggerateChartForProject) (*structure.AggragetedCollectionVolumnResp, error) {

	pe := &entity.AggerateChartForProject{}
	err := copier.Copy(pe, req)
	if err != nil {
		return nil, err
	}

	res, err := u.Repo.AggregateVolumnCollection(pe)
	if err != nil {
		return nil, err
	}

	resp := []structure.AggragetedCollection{}
	for _, item := range res{
		tmp := structure.AggragetedCollection{
			ProjectID: item.ID.ProjectID,
			ProjectName: item.ID.ProjectName,
			Timestamp: item.ID.Timestamp,
			Amount: item.Amount,

		}

		resp = append(resp, tmp)
	}
	
	return &structure.AggragetedCollectionVolumnResp{Volumns: resp}, nil
}

func (u Usecase) GetChartDataOFTokens(req structure.AggerateChartForToken) (*structure.AggragetedTokenVolumnResp, error) {

	pe := &entity.AggerateChartForToken{}
	err := copier.Copy(pe, req)
	if err != nil {
		return nil, err
	}

	res, err := u.Repo.AggregateVolumnToken(pe)
	if err != nil {
		return nil, err
	}

	resp := []structure.AggragetedTokenURI{}
	for _, item := range res{
		tmp := structure.AggragetedTokenURI{
			TokenID: item.ID.TokenID,
			Timestamp: item.ID.Timestamp,
			Amount: item.Amount,

		}

		resp = append(resp, tmp)
	}
	
	return &structure.AggragetedTokenVolumnResp{Volumns: resp}, nil
}