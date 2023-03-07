package entity

import (
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
)

type StatusMint int

const (
	StatusMint_Pending            StatusMint = iota // 0: pending: waiting for payment
	StatusMint_WaitingForConfirms                   // 1: Waiting for confirmations 0 of 6
	StatusMint_ReceivedFund                         // 2: received fund from user

	StatusMint_Minting // 3: minting
	StatusMint_Minted  // 4: mint success

	StatusMint_SendingNFTToUser // 5: sending nft to user
	StatusMint_SentNFTToUser    // 6: send nft to user success:

	StatusMint_SendingFundToMaster // 7: sending nft to user
	StatusMint_SentFundToMaster    // 8: send btc to master:

	StatusMint_TxMintFailed // 9: tx mint failed

	StatusMint_NeedToRefund // 10: balance not enough or mint out...

	StatusMint_Refunding // 11: refunding
	StatusMint_Refunded  // 12: refunding

	StatusMint_TxRefundFailed // 13: tx refund fund failed

	StatusMint_NotEnoughBalanceToSendMaster // 14: tx refund fund failed
	StatusMint_AlreadySentMaster            // 15: sent but can't save the tx, need to get it later.
)

var StatusMintToText = map[StatusMint]string{

	StatusMint_Pending: "Waiting for payment",

	StatusMint_WaitingForConfirms: "Waiting for payment confirmation",

	StatusMint_ReceivedFund: "Transferring", // //"Minting",

	StatusMint_Minting: "Transferring", // "Minting",
	StatusMint_Minted:  "Transferred",  // "Minted",

	StatusMint_SendingNFTToUser: "Transferring",
	StatusMint_SentNFTToUser:    "Transferred",

	StatusMint_SendingFundToMaster: "Transferred", //"Sending funds to master",
	StatusMint_SentFundToMaster:    "Transferred", //"Sent funds to master",

	StatusMint_TxMintFailed: "Mint failed",

	StatusMint_NeedToRefund: "Waiting to refund",

	StatusMint_Refunding: "Refunding",
	StatusMint_Refunded:  "Refunded",

	StatusMint_TxRefundFailed: "Refunding",
}

type MintNftBtc struct {
	BaseEntity  `bson:",inline"`
	UserAddress string `bson:"user_address"` //user's wallet address from FE

	UserID string `bson:"user_id"`

	OriginUserAddress string `bson:"origin_user_address"` //user's wallet address from FE
	RefundUserAdress  string `bson:"refund_user_address"`

	Amount string `bson:"amount"` // amount required

	PayType string `bson:"payType"` // eth/btc...

	ReceiveAddress string `bson:"receiveAddress"` // address generated to receive coin from users.
	PrivateKey     string `bson:"privateKey"`     // private key of the receive wallet.

	Balance string `bson:"balance"` // balance after check

	ExpiredAt time.Time `bson:"expired_at"`

	Status StatusMint `bson:"status"` // status for record

	TxReceived   string `bson:"tx_received"` // tx received fund from user.
	TxMintNft    string `bson:"tx_mint_nft"`
	TxSendNft    string `bson:"tx_send_nft"`
	TxSendMaster string `bson:"tx_send_master"`
	TxRefund     string `bson:"tx_refund"`

	FileURI string `bson:"fileURI"` // FileURI will be mount if OrdAddress get all amount

	InscriptionID string `bson:"inscriptionID"`

	ProjectID string `bson:"projectID"` //projectID

	// just log for users, not using for the job checking.
	IsConfirm        bool `bson:"isConfirm"`        // rereive fund
	IsMinted         bool `bson:"isMinted"`         // minted
	IsSentUser       bool `bson:"isSentUser"`       // sent nft to user
	IsSentMaster     bool `bson:"isSentMaster"`     // withdrawn to master wallet
	IsRefund         bool `bson:"isRefund"`         // refund btc to btc
	IsUpdatedNftInfo bool `bson:"isUpdatedNftInfo"` // update project info

	OutputMintNFT interface{} `bson:"output_mint_nft"` // output from mint nft execute
	OutputSendNFT interface{} `bson:"output_send_nft"` // output from send nft execute

	ReasonRefund string `bson:"reason_refund"` // the rason of refund

	AmountSentMaster string `bson:"amount_sent_master"` // amount withdrawn to the master wallet
	AmountRefundUser string `bson:"amount_refund_user"` // amount refund eth/btc user

	// for analytics:

	BtcRate           float64 `bson:"btc_rate"`
	EthRate           float64 `bson:"eth_rate"`
	ProjectMintPrice  int     `bson:"project_mint_price"`  // btc for 1 item
	ProjectNetworkFee int     `bson:"project_network_fee"` // btc for 1 item

	MintPriceByPayType  string `bson:"mint_price_by_pay_type"`  // by pay type, for 1 item
	NetworkFeeByPayType string `bson:"network_fee_by_pay_type"` // by pay type, for 1 item

	MintPriceByPayTypeTotal  string `bson:"mint_price_by_pay_type_total"`  // by pay type, for n item
	NetworkFeeByPayTypeTotal string `bson:"network_fee_by_pay_type_total"` // by pay type, for n item

	EstFeeInfo map[string]MintFeeInfo `bson:"mintFeeInfo"` // 1 item

	FeeSendMaster string `bson:"fee_send_master"` // maybe for eth only

	MintFee int `bson:"mint_fee"` // real mint fee

	IsMerged bool `bson:"isMerged"` // with ord v5.1: mint = mint + send, 1 tx

	// for mint batch:
	Quantity      int `bson:"quantity"`
	BatchParentId int `bson:"patch_parent_id"`
}

func (u MintNftBtc) TableName() string {
	return utils.MINT_NFT_BTC
}

func (u MintNftBtc) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

type MintNftBtcLogs struct {
	BaseEntity  `bson:",inline"`
	RecordID    string      `bson:"record_id"`
	Table       string      `bson:"table"`
	Name        string      `bson:"name"`
	Status      interface{} `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u MintNftBtcLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

func (u MintNftBtcLogs) TableName() string {
	return "mint_nft_btc_logs"
}

type MintNftBtcResp struct {
	UUID string `bson:"uuid"`

	InscriptionID  string `bson:"inscriptionID"`  // tokenID in ETH
	ReceiveAddress string `bson:"receiveAddress"` // address is generated for receive user fund.

	Amount string `bson:"amount"`

	PayType string `bson:"payType"` // eth/btc...

	ExpiredAt time.Time `bson:"expired_at"`

	Status StatusMint `bson:"status"` // status for record

	TxSendNft string `bson:"tx_send_nft"`
	TxMintNft string `bson:"tx_mint_nft"`

	OriginUserAddress string `bson:"origin_user_address"`
}

type MintFeeInfo struct {

	//string
	MintPrice   string `json:"mintPrice"`
	NetworkFee  string `json:"networkFee"`
	MintFee     string `json:"mintFee"`
	SendNftFee  string `json:"sendNftFee"`
	SendFundFee string `json:"sendFundFee"`
	TotalAmount string `json:"totalAmount"`

	// big number
	MintPriceBigInt  *big.Int `json:"mintPriceBigInt"`
	NetworkFeeBigInt *big.Int `json:"networkFeeBigInt"`

	MintFeeBigInt     *big.Int `json:"mintFeeBigInt"`
	SendNftFeeBigInt  *big.Int `json:"sendNftFeeBigInt"`
	SendFundFeeBigInt *big.Int `json:"sendFundFeeBigInt"`

	TotalAmountBigInt *big.Int `json:"totalAmountBigInt"`

	EthPrice float64 `json:"ethPrice"`
	BtcPrice float64 `json:"btcPrice"`
	Decimal  int     `json:"decimal"`
}
