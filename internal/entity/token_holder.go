package entity

import (
	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type FilterTokenHolders struct {
	BaseFilters
}

type TokenHolder struct {
	BaseEntity           `bson:",inline" json:"-"`
	ID                   string      `bson:"id" json:"-"`
	ContractDecimals     int         `bson:"contract_decimals" json:"contract_decimals"`
	ContractName         string      `bson:"contract_name" json:"contract_name"`
	ContractTickerSymbol string      `bson:"contract_ticker_symbol" json:"contract_ticker_symbol"`
	ContractAddress      string      `bson:"contract_address" json:"contract_address"`
	SupportsErc          interface{} `bson:"supports_erc" json:"supports_erc"`
	LogoURL              string      `bson:"logo_url" json:"logo_url"`
	Address              string      `bson:"address" json:"address"`
	Balance              string      `bson:"balance" json:"balance"`
	TotalSupply          string      `bson:"total_supply" json:"total_supply"`
	BlockHeight          int         `bson:"block_height" json:"block_height"`
	Profile              *Users      `bson:"profile" json:"profile"`
	CurrentRank          int32       `bson:"current_rank" json:"current_rank"`
	OldRank              *int32       `bson:"old_rank" json:"old_rank"`
}

func (u TokenHolder) TableName() string {
	return utils.COLLECTION_LEADERBOARD_TOKEN_HOLDER
}

func (u TokenHolder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
