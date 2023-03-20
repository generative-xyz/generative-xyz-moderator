package usecase

import (
	"github.com/jinzhu/copier"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) GetChartDataOFTokens(req structure.AggerateChartForProject) (*structure.AggragetedInscriptionVolumnResp, error) {

	pe := &entity.AggerateChartForProject{}
	err := copier.Copy(pe, req)
	if err != nil {
		return nil, err
	}

	res, err := u.Repo.AggregateVolumeInscription(pe)
	if err != nil {
		return nil, err
	}

	resp := []structure.AggragetedInscription{}
	for _, item := range res{
		tmp := structure.AggragetedInscription{
			ProjectID: item.ID.ProjectID,
			ProjectName: item.ID.ProjectName,
			Timestamp: item.ID.Timestamp,
			Amount: item.Amount,

		}

		resp = append(resp, tmp)
	}
	
	return &structure.AggragetedInscriptionVolumnResp{Volumns: resp}, nil
}