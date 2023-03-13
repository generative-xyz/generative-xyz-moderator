package request

type CreateDaoProjectRequest struct {
	ProjectId string `json:"project_id"`
	CreatedBy string `json:"-"`
}
