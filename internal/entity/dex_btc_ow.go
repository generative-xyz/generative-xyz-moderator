package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type DexBTCOWSubmitTx struct {
	BaseEntity  `bson:",inline"`
	OrderID     string `bson:"order_id" json:"order_id"`
	PurchaseRaw string `bson:"purchase_raw" json:"purchase_raw"`
	SetupRaw    string `bson:"setup_raw" json:"setup_raw"`
	PurchaseTx  string `bson:"purchase_tx" json:"purchase_tx"`
	SetupTx     string `bson:"setup_tx" json:"setup_tx"`
}

func (u DexBTCOWSubmitTx) TableName() string {
	return utils.COLLECTION_DEX_BTC_OW_SUBMIT
}

func (u DexBTCOWSubmitTx) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type DexBTCOWInscription struct {
	BaseEntity     `bson:",inline"`
	CollectionSlug string `bson:"collection_slug" json:"collection_slug"`
	InscriptionID  string `bson:"inscription_id" json:"inscription_id"`
	Price          int    `bson:"price" json:"price"`
	SellerAddress  string `bson:"seller_address" json:"seller_address"`
}

func (u DexBTCOWInscription) TableName() string {
	return utils.COLLECTION_DEX_BTC_OW_INSCRIPTION
}

func (u DexBTCOWInscription) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
