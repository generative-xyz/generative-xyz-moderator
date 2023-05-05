package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type StoreFileWallet struct {
	BaseEntity    `bson:",inline"`
	WalletAddress string `bson:"wallet_address"` // the wallet address
	PrivateKey    string `bson:"private_key"`    // private key (has encrypt).
}

func (u StoreFileWallet) TableName() string {
	return "store_file_wallets"
}
func (u StoreFileWallet) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
