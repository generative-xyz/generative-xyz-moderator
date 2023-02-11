package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type BuyStatus int

const (
	StatusBuy_Pending          BuyStatus = iota // 0: pending: waiting for fund
	StatusBuy_ReceivedFund                      // 1: received fund from user (buyer)
	StatusBuy_SendingNFT                        // 2: sending nft
	StatusBuy_SentNFT                           // 3: send nft success
	StatusBuy_SendingBTC                        // 4: send nft to buyer success
	StatusBuy_SentBTC                           // 5: send btc to seller success
	StatusBuy_TxSendNFTFailed                   // 6: tx send nft to buyer failed
	StatusBuy_TxSendBTCFailed                   // 7: tx send btc to seller failed
	StatusBuy_NotEnoughBalance                  // 8: balance not enough
	StatusBuy_NeedToRefund                      // 9: Need to refund BTC
)

type MarketplaceBTCBuyOrder struct {
	BaseEntity      `bson:",inline"`
	OrdAddress      string    `bson:"ord_address"`
	ItemID          string    `bson:"item_id"`
	InscriptionID   string    `bson:"inscriptionID"` // tokenID in btc
	Status          BuyStatus `bson:"status"`
	ErrCount        int       `bson:"err_count"`
	SegwitAddress   string    `bson:"segwit_address"`
	ReceivedBalance string    `bson:"received_balance"`
	SegwitKey       string    `bson:"segwit_key"`
	ExpiredAt       time.Time `bson:"expired_at"`
	RawData         string    `bson:"raw_data"` // raw data of sending btc (for retry)
	UTXO            string    `bson:"utxo"`
	TxSendNFT       string    `bson:"tx_send_nft"`
	TxSendBTC       string    `bson:"tx_send_btc"`
}

func (u MarketplaceBTCBuyOrder) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_BUY
}

func (u MarketplaceBTCBuyOrder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
