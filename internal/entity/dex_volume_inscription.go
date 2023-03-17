package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
	"time"
)

type DexVolumeInscription struct {
	BaseEntity `bson:",inline"`
	Timestamp  *time.Time                   `bson:"timestamp"`
	Metadata   DexVolumeInscriptionMetadata `bson:"metadata"`
	Amount     uint64                       `bson:"amount"`
}

type DexVolumeInscriptionMetadata struct {
	InscriptionId string `bson:"inscription_id"`
	MatchedTx     string `bson:"matched_tx"`
}

func (u DexVolumeInscription) TableName() string {
	return "dex_volume_inscription"
}

func (u DexVolumeInscription) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
