package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type FirebaseRegistrationToken struct {
	BaseEntity        `json:"base_entity,omitempty"`
	UserWallet        string    `json:"user_wallet,omitempty" bson:"user_wallet"`
	RegistrationToken string    `json:"registration_token,omitempty" bson:"registration_token"`
	DeviceType        string    `json:"device_type,omitempty" bson:"device_type"`
	CreatedAt         time.Time `json:"created_at,omitempty" bson:"created_at"`
}

func (m FirebaseRegistrationToken) TableName() string {
	return "firebase_registration_token"
}
func (m FirebaseRegistrationToken) ToBson() (*bson.D, error) {
	return helpers.ToDoc(m)
}
