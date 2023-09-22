package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/utils/constants/dao_project"
	"rederinghub.io/utils/constants/dao_project_voted"
)

type DaoProject struct {
	BaseEntity      `json:",inline" bson:",inline"`
	SeqId           uint               `json:"seq_id" bson:"seq_id"`
	CreatedBy       string             `json:"created_by,omitempty" bson:"created_by"`
	User            *Users             `json:"user,omitempty" bson:"user,omitempty"`
	ProjectId       primitive.ObjectID `json:"project_id,omitempty" bson:"project_id"`
	Project         *Projects          `json:"project,omitempty" bson:"project,omitempty"`
	ExpiredAt       time.Time          `json:"expired_at,omitempty" bson:"expired_at"`
	UpdatedAt       *time.Time         `json:"updated_at" bson:"updated_at"`
	Status          dao_project.Status `json:"status,omitempty" bson:"status"`
	DaoProjectVoted []*DaoProjectVoted `json:"dao_project_voted,omitempty" bson:"dao_project_voted,omitempty"`
	TotalAgainst    *int64             `json:"total_against,omitempty" bson:"total_against,omitempty"`
	TotalVote       *int64             `json:"total_vote,omitempty" bson:"total_vote,omitempty"`
}

func (m DaoProject) TableName() string {
	return "dao_project"
}

func (m DaoProject) Expired() bool {
	return m.ExpiredAt.UTC().Unix() < time.Now().UTC().Unix()
}

type DaoProjectVoted struct {
	BaseEntity   `json:",inline" bson:",inline"`
	CreatedBy    string                   `json:"created_by,omitempty" bson:"created_by"`
	DaoProjectId primitive.ObjectID       `json:"dao_project_id,omitempty" bson:"dao_project_id"`
	Status       dao_project_voted.Status `json:"status,omitempty" bson:"status"`
}

func (m DaoProjectVoted) TableName() string {
	return "dao_project_voted"
}
