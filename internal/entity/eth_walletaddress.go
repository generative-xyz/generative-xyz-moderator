package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type ETHWalletAddress struct {
	BaseEntity       `bson:",inline"`
	UserAddress      string              `bson:"user_address"` //user's wallet address from FE
	Amount           string              `bson:"amount"`
	OrdAddress       string              `bson:"ordAddress"`    // address is generated from ORD service, which receive all amount
	FileURI          string              `bson:"fileURI"`       // FileURI will be mount if OrdAddress get all amount
	IsConfirm        bool                `bson:"isConfirm"`     //default: false, if OrdAddress get all amount it will be set true
	InscriptionID    string              `bson:"inscriptionID"` // tokenID in ETH
	Mnemonic         string              `bson:"mnemonic"`
	IsMinted         bool                `bson:"isMinted"`         //default: false. If InscriptionID exist which means token is minted, it's true
	ProjectID        string              `bson:"projectID"`        //projectID
	MintResponse     MintStdoputResponse `bson:"mintResponse"`     // after token has been mint
	Balance          string              `bson:"balance"`          // balance after check
	BalanceCheckTime int                 `bson:"balanceCheckTime"` // Total balance check time
	DelegatedAddress string              `bson:"delegatedAddress"` //applied the whitelist contract for discount: default false
}

type FilterETHWalletAddress struct {
	BaseFilters
	UserAddress   *string
	Amount        *string
	OrdAddress    *string
	IsConfirm     *string
	InscriptionID *string
}

func (u ETHWalletAddress) TableName() string {
	return utils.COLLECTION_ETH_WALLET_ADDRESS
}

func (u ETHWalletAddress) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
