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
	BaseEntity    `bson:",inline"`
	OrdAddress    string    `bson:"ord_address"`
	ItemID        string    `bson:"item_id"`
	InscriptionID string    `bson:"inscriptionID"` // nftID in btc
	Status        BuyStatus `bson:"status"`

	PayType string `bson:"pay_type"`

	Price   string  `bson:"price"` //maybe eth/btc
	BtcRate float64 `bson:"btc_rate"`
	EthRate float64 `bson:"eth_rate"`

	ErrCount int `bson:"err_count"`

	ReceiveAddress  string `bson:"receive_address"` // address generated to receive coin from users.
	PrivateKey      string `bson:"privateKey"`      // private key of the receive wallet.
	ReceivedBalance string `bson:"received_balance"`

	ExpiredAt             time.Time   `bson:"expired_at"`
	RawData               string      `bson:"raw_data"` // raw data of sending btc (for retry)
	UTXO                  string      `bson:"utxo"`
	TxSendNFT             string      `bson:"tx_send_nft"`
	TxSendBTC             string      `bson:"tx_send_btc"`
	OutputSendNFT         interface{} `bson:"output_send_nft"`
	FeeChargeBTCBuyer     int         `bson:"fee_charge_btc_buyer"`
	RoyaltyChargeBTCBuyer int         `bson:"royalty_charge_btc_buyer"`
	AmountBTCSentSeller   int         `bson:"amount_btc_send_seller"`

	SegwitAddress string `bson:"segwit_address"` // remove for new flow
	SegwitKey     string `bson:"segwit_key"`     // remove for new flow
}

func (u MarketplaceBTCBuyOrder) TableName() string {
	return utils.COLLECTION_MARKETPLACE_BTC_BUY
}

func (u MarketplaceBTCBuyOrder) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
