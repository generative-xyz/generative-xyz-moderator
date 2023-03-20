package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type TokenTx struct {
	BaseEntity               `bson:",inline"`
	InscriptionID string     `bson:"inscription_id"`
	Tx            string     `bson:"tx"`
	PrevTx        string     `bson:"prev_tx"`
	NextTx        string     `bson:"next_tx"`
	LastTimeCheck *time.Time `bson:"last_time_check"`
	Depth         int        `bson:"depth"`
	NumFailed     int        `bson:"num_failed"`
	Priority      int64      `bson:"pritority"`
}

type UpdateTokenTxRequest struct {
	InscriptionID string     `bson:"inscription_id"`
	Tx            string     `bson:"tx"`
	PrevTx        *string    `bson:"prev_tx"`
	NextTx        *string    `bson:"next_tx"`
	LastTimeCheck  *time.Time `bson:"lst_time_check"`
	Depth         *int       `bson:"depth"`
}

func (t TokenTx) TableName() string {
	return utils.COLLECTION_TOKEN_TX
}

func (t TokenTx) ToBson() (*bson.D, error) {
	return helpers.ToDoc(t)
}

func (t *TokenTx) SetTokenTxID() {
	if t.UUID == "" {
		t.ID = primitive.NewObjectID()
		t.UUID = t.ID.Hex()
	}
}
