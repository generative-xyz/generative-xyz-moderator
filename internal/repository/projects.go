package repository

import (
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"

	"go.mongodb.org/mongo-driver/bson"
)

func (r Repository) FindProject( projectID string) (*entity.Projects, error) {
	resp := &entity.Projects{}

	usr, err := r.FilterOne(entity.Projects{}.TableName(), bson.D{{utils.KEY_UUID, projectID}})
	if err != nil {
		return nil, err
	}

	err = helpers.Transform(usr, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r Repository) CreateProject(data *entity.Projects) error {

	err := r.InsertOne(data.TableName(), data)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) ListProjects(filter entity.FilterProjects) (*entity.Pagination, error)  {
	confs := []entity.Configs{}
	resp := &entity.Pagination{}
	f := bson.M{}

	p, err := r.Paginate(utils.COLLECTION_PROJECTS, filter.Page, filter.Limit, f, filter.SortBy, filter.Sort, &confs)
	if err != nil {
		return nil, err
	}
	
	resp.Result = confs
	resp.Page = p.Pagination.Page
	resp.Total = p.Pagination.Total
	return resp, nil
}
