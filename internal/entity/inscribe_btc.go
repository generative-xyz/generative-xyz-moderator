package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type InscribeBTC struct {
	BaseEntity        `bson:",inline"`
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
	// MintResponse      MintStdoputResponse `bson:"mintResponse"` // after token has been mint
	Balance   string    `bson:"balance"` // balance after check
	FeeRate   int32     `bson:"fee_rate"`
	ExpiredAt time.Time `bson:"expired_at"`
	IsSuccess bool      `bson:"isSuccess"`

	Status    StatusInscribe `bson:"status"` // status for record
	TxSendBTC string         `bson:"tx_send_btc"`
	TxSendNft string         `bson:"tx_send_nft"`
	TxMintNft string         `bson:"tx_mint_nft"`

	OutputMintNFT interface{} `bson:"output_mint_nft"`
	OutputSendNFT interface{} `bson:"output_send_nft"`
	UserUuid      string      `bson:"user_uuid"`
}

func (u InscribeBTC) TableName() string {
	return utils.INSCRIBE_BTC
}

func (u InscribeBTC) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type StatusInscribe int

const (
	StatusInscribe_Pending      StatusInscribe = iota // 0: pending: waiting for fund
	StatusInscribe_ReceivedFund                       // 1: received fund from user (buyer)

	StatusInscribe_SendingBTCFromSegwitAddrToOrdAddr // 2: sending btc from segwit address to ord address
	StatusInscribe_SentBTCFromSegwitAddrToOrdAdd     // 3: send btc from segwit address to ord address success

	StatusInscribe_Minting // 4: minting
	StatusInscribe_Minted  // 5: mint success

	StatusInscribe_SendingNFTToUser // 6: sending nft to user
	StatusInscribe_SentNFTToUser    // 7: send nft to user success: flow DONE

	StatusInscribe_TxSendBTCFromSegwitAddrToOrdAddrFailed // 8: send btc from segwit address to ord address failed
	StatusInscribe_TxSendBTCToUserFailed                  // 9: send nft to user failed
	StatusInscribe_TxMintFailed                           // 10: tx mint failed

	StatusInscribe_NotEnoughBalance // 11: balance not enough
	StatusInscribe_NeedToRefund     // 12: Need to refund BTC
)

type InscribeBTCLogs struct {
	BaseEntity  `bson:",inline"`
	RecordID    string      `bson:"record_id"`
	Table       string      `bson:"table"`
	Name        string      `bson:"name"`
	Status      interface{} `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u InscribeBTCLogs) TableName() string {
	return "inscribe_btc_logs"
}

func (u InscribeBTCLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type FilterInscribeBT struct {
	BaseFilters
	UserAddress   *string
	Amount        *string
	OrdAddress    *string
	IsConfirm     *string
	InscriptionID *string
	UserUuid      *string
}

type InscribeBTCResp struct {
	UUID string `bson:"uuid" json:"uuid,omitempty"`
	// UserAddress string `bson:"user_address"` //user's wallet address from FE
	// OriginUserAddress string `bson:"origin_user_address"` //user's wallet address from FE
	// Amount            string `bson:"amount"`
	// MintFee           string `bson:"mint_fee"`
	// SentTokenFee      string `bson:"sent_token_fee"`
	// OrdAddress        string `bson:"ordAddress"` // address is generated from ORD service, which receive all amount
	// SegwitAddress     string `bson:"segwit_address"`
	// FileURI       string `bson:"fileURI"`       // FileURI will be mount if OrdAddress get all amount
	IsConfirm bool `bson:"isConfirm" json:"isConfirm,omitempty"` //default: false, if OrdAddress get all amount it will be set true
	IsMinted  bool `bson:"isMinted" json:"isMinted,omitempty"`   //default: false. If InscriptionID exist which means token is minted, it's true
	IsSuccess bool `bson:"isSuccess" json:"isSuccess,omitempty"` //default: false. If InscriptionID was sent to user, it's true

	InscriptionID string `bson:"inscriptionID" json:"inscriptionID,omitempty"` // tokenID in ETH
	// Mnemonic          string `bson:"mnemonic"`
	// SegwitKey         string `bson:"segwit_key"`
	// MintResponse MintStdoputResponse `bson:"mintResponse"` // after token has been mint
	// Balance   string    `bson:"balance"` // balance after check
	FeeRate   int32     `bson:"fee_rate" json:"feeRate,omitempty"`
	ExpiredAt time.Time `bson:"expired_at" json:"expiredAt,omitempty"`

	Status    StatusInscribe `bson:"status" json:"status,omitempty"` // status for record
	TxSendBTC string         `bson:"tx_send_btc" json:"txSendBtc,omitempty"`
	TxSendNft string         `bson:"tx_send_nft" json:"txSendNft,omitempty"`
	TxMintNft string         `bson:"tx_mint_nft" json:"txMintNft,omitempty"`
	UserUuid  string         `bson:"user_uuid" json:"userUuid,omitempty"`
}
