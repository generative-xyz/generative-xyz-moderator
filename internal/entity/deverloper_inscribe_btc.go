package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type DeveloperInscribe struct {
	BaseEntity `bson:",inline"`

	DeveloperKeyUuid string `bson:"developer_key_uuid"`

	UserAddress       string `bson:"user_address"`        //user's wallet address from FE
	OriginUserAddress string `bson:"origin_user_address"` //user's wallet address from FE
	Amount            string `bson:"amount"`
	MintFee           string `bson:"mint_fee"`
	SentTokenFee      string `bson:"sent_token_fee"`
	OrdAddress        string `bson:"ordAddress"` // address is generated from ORD service, which receive all amount
	SegwitAddress     string `bson:"segwit_address"`
	FileURI           string `bson:"fileURI"` // FileURI will be mount if OrdAddress get all amount
	LocalLink         string `bson:"local_link"`
	FileName          string `bson:"file_name"`
	IsConfirm         bool   `bson:"isConfirm"`     //default: false, if OrdAddress get all amount it will be set true
	InscriptionID     string `bson:"inscriptionID"` // tokenID in ETH
	Mnemonic          string `bson:"mnemonic"`
	SegwitKey         string `bson:"segwit_key"`
	IsMinted          bool   `bson:"isMinted"` //default: false. If InscriptionID exist which means token is minted, it's true

	Balance   string    `bson:"balance"` // balance after check
	FeeRate   int32     `bson:"fee_rate"`
	ExpiredAt time.Time `bson:"expired_at"`
	IsSuccess bool      `bson:"isSuccess"`

	Status    StatusDeveloperInscribe `bson:"status"` // status for record
	TxSendBTC string                  `bson:"tx_send_btc"`
	TxSendNft string                  `bson:"tx_send_nft"`
	TxMintNft string                  `bson:"tx_mint_nft"`

	OutputMintNFT    interface{} `bson:"output_mint_nft"`
	OutputSendNFT    interface{} `bson:"output_send_nft"`
	UserUuid         string      `bson:"user_uuid"`
	TokenAddress     string      `bson:"token_address"`
	TokenId          string      `bson:"token_id"`
	IsAuthentic      bool        `bson:"is_authentic"`
	OrdinalsTx       string      `bson:"ordinals_tx"`
	OrdinalsTxStatus uint64      `bson:"ordinals_tx_status"`
}

func (u DeveloperInscribe) TableName() string {
	return "developer_inscribes"
}

func (u DeveloperInscribe) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type StatusDeveloperInscribe int

const (
	StatusDeveloperInscribe_Pending      StatusDeveloperInscribe = iota // 0: pending: waiting for fund
	StatusDeveloperInscribe_ReceivedFund                                // 1: received fund from user (buyer)

	StatusDeveloperInscribe_SendingBTCFromSegwitAddrToOrdAddr // 2: sending btc from segwit address to ord address
	StatusDeveloperInscribe_SentBTCFromSegwitAddrToOrdAdd     // 3: send btc from segwit address to ord address success

	StatusDeveloperInscribe_Minting // 4: minting
	StatusDeveloperInscribe_Minted  // 5: mint success

	StatusDeveloperInscribe_TxSendBTCFromSegwitAddrToOrdAddrFailed // 6: send btc from segwit address to ord address failed
	StatusDeveloperInscribe_TxSendBTCToUserFailed                  // 7: send nft to user failed
	StatusDeveloperInscribe_TxMintFailed                           // 8: tx mint failed

	StatusDeveloperInscribe_NotEnoughBalance // 9: balance not enough
	StatusDeveloperInscribe_NeedToRefund     // 10: Need to refund BTC
)

func (s StatusDeveloperInscribe) Ordinal() int {
	return int(s)
}

type DeveloperInscribeBTCLogs struct {
	BaseEntity  `bson:",inline"`
	RecordID    string      `bson:"record_id"`
	Table       string      `bson:"table"`
	Name        string      `bson:"name"`
	Status      interface{} `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u DeveloperInscribeBTCLogs) TableName() string {
	return "inscribe_btc_logs"
}

func (u DeveloperInscribeBTCLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type DeveloperInscribeBTCResp struct {
	UUID          string                  `bson:"uuid" json:"uuid"`
	Amount        string                  `bson:"amount"  json:"amount"`
	IsConfirm     bool                    `bson:"isConfirm" json:"isConfirm"`         //default: false, if OrdAddress get all amount it will be set true
	IsMinted      bool                    `bson:"isMinted" json:"isMinted"`           //default: false. If InscriptionID exist which means token is minted, it's true
	InscriptionID string                  `bson:"inscriptionID" json:"inscriptionID"` // tokenID in ETH
	FeeRate       int32                   `bson:"fee_rate" json:"feeRate"`
	ExpiredAt     time.Time               `bson:"expired_at" json:"expiredAt"`
	Status        StatusDeveloperInscribe `bson:"status" json:"status"` // status for record
	TxSendBTC     string                  `bson:"tx_send_btc" json:"txSendBtc"`
	TxSendNft     string                  `bson:"tx_send_nft" json:"txSendNft"`
	TxMintNft     string                  `bson:"tx_mint_nft" json:"txMintNft"`
	UserUuid      string                  `bson:"user_uuid" json:"userUuid"`
	IsAuthentic   bool                    `bson:"is_authentic" json:"isAuthentic"`
	TokenAddress  string                  `bson:"token_address" json:"tokenAddress"`
	TokenId       string                  `bson:"token_id" json:"tokenId"`
}

type FilterDeveloperInscribeBT struct {
	BaseFilters
	UserAddress   *string
	Amount        *string
	OrdAddress    *string
	IsConfirm     *string
	InscriptionID *string
	UserUuid      *string
	TokenAddress  *string
	TokenId       *string
	NeStatuses    []StatusDeveloperInscribe
}
