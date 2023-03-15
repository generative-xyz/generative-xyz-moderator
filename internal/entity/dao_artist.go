package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/utils/constants/dao_artist"
	"rederinghub.io/utils/constants/dao_artist_voted"
)

type DaoArtist struct {
	BaseEntity     `json:",inline" bson:",inline"`
	SeqId          uint              `json:"seq_id" bson:"seq_id"`
	CreatedBy      string            `json:"created_by,omitempty" bson:"created_by"`
	User           *Users            `json:"user,omitempty" bson:"user,omitempty"`
	ExpiredAt      time.Time         `json:"expired_at,omitempty" bson:"expired_at"`
	Status         dao_artist.Status `json:"status,omitempty" bson:"status"`
	DaoArtistVoted []*DaoArtistVoted `json:"dao_artist_voted,omitempty" bson:"dao_artist_voted,omitempty"`
	TotalReport    *int64            `json:"total_report,omitempty" bson:"total_report,omitempty"`
	TotalVerify    *int64            `json:"total_verify,omitempty" bson:"total_verify,omitempty"`
}

func (m DaoArtist) Expired() bool {
	return m.ExpiredAt.UTC().Unix() < time.Now().UTC().Unix()
}
func (m DaoArtist) TableName() string {
	return "dao_artist"
}

type DaoArtistVoted struct {
	BaseEntity  `json:",inline" bson:",inline"`
	CreatedBy   string                  `json:"created_by,omitempty" bson:"created_by"`
	DaoArtistId primitive.ObjectID      `json:"dao_artist_id,omitempty" bson:"dao_artist_id"`
	Status      dao_artist_voted.Status `json:"status,omitempty" bson:"status"`
}

func (m DaoArtistVoted) TableName() string {
	return "dao_artist_voted"
}
