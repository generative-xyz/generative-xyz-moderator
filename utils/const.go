package utils

type QuerySort struct {
	Sort   int
	SortBy string
}

func ParseSort(key string) QuerySort {
	sortParams := make(map[string]QuerySort)
	sortParams["custom"] = QuerySort{Sort: -1, SortBy: ""}
	sortParams["newest"] = QuerySort{Sort: -1, SortBy: "created_at"}
	sortParams["oldest"] = QuerySort{Sort: 1, SortBy: "created_at"}
	sortParams["priority-asc"] = QuerySort{Sort: 1, SortBy: "priority"}
	sortParams["priority-desc"] = QuerySort{Sort: -1, SortBy: "priority"}
	sortParams["price-asc"] = QuerySort{Sort: 1, SortBy: "price"}
	sortParams["price-desc"] = QuerySort{Sort: -1, SortBy: "price"}
	sortParams["minted-newest"] = QuerySort{Sort: -1, SortBy: "minted_time"}
	sortParams["token-price-desc"] = QuerySort{Sort: -1, SortBy: "stats.price_int"}
	sortParams["token-price-asc"] = QuerySort{Sort: 1, SortBy: "stats.price_int"}
	sortParams["trending-score"] = QuerySort{Sort: -1, SortBy: "stats.trending_score"}

	sort, ok := sortParams[key]
	if !ok {
		return sortParams["custom"]
	}

	return sort
}

const (
	MAX_CHECK_BALANCE int    = 15
	PERCENT_EARNING   int    = 900
	API_KEY           string = "Api-Key"
	//AUTH_TOKEN           string = "Authorization" //token will be save in this variable
	AUTH_TOKEN                string = "Authorization" //token will be save in this variable
	REDIS_VERIFIED_TOKEN      string = "verified_token"
	REDIS_PROFILE             string = "profile"
	REDIS_NFT_METADATA_KEY    string = "nfts_metadata_%s_%s"
	REDIS_PAGINATION_KEY      string = "pagination_%s"
	REDIS_PAGINATION_DATA_KEY string = "pagination_data_%s"
	SIGNED_USER_ID            string = "signed_user_id"
	SIGNED_ADMIN_KEY          string = "admin_user_id_%s"
	SIGNED_WALLET_ADDRESS     string = "signed_wallet_address"
	SIGNED_EMAIL              string = "signed_email"
	SERVICE_API_KEY           string = "service_key"
	TRACER_EMAIL              string = "email"

	CODE_LENGTH         int    = 3
	WORK_STATION        string = "working_place"
	WORK_STATION_PREFIX        = "SD4"
	OTHER_TYPE          string = "other"
	MODIFIED_TOKEN      string = "modified-token-%s"
	VERIFY_TOKEN        string = "verify-token-%s"

	EMAIL_TAG              string = "email"
	TOKEN_ID_TAG           string = "tokenID"
	PROJECT_ID_TAG         string = "projectID"
	WALLET_ADDRESS_TAG     string = "wallet_address"
	ORD_WALLET_ADDRESS_TAG string = "ord_wallet_address"
	GEN_NFT_ADDRESS_TAG    string = "gen_nft_address"

	PubsubCreateDeviceType           string = "Device:PubsubCreateDeviceType"
	PubsubUpdateDeviceType           string = "Device:PubsubUpdateDeviceType"
	PubsubDeleteDeviceType           string = "Device:PubsubDeleteDeviceType"
	PubsubSendMessageToSlack         string = "Device:PubsubSendMessageToSlack"
	PUBSUB_SEND_OTP                  string = "Hybrid:SendOtp"
	PUBSUB_REGISTER                  string = "WorkspaceGateway::PubsubRegister"
	PUBSUB_FORGOT_PASSWORD           string = "Hybrid:ResetPasswordEmail"
	NFT_CACHE_EXPIRED_TIME           int    = 86400
	TOKEN_CACHE_EXPIRED_TIME         int    = 86400 * 30       //a month (second)
	REFRESH_TOKEN_CACHE_EXPIRED_TIME int    = 86400 * 360      //a year (second)
	DB_CACHE_EXPIRED_TIME            int    = 86400            //a week
	DB_CACHE_KEY                     string = "db.cache.%s.%s" //a week
	NONCE_MESSAGE_FORMAT             string = "Welcome %s to Generative"

	KEY_UUID                       string = "uuid"
	KEY_BASE_PRODUCT_KEY           string = "product_key"
	KEY_ORDER_ID                   string = "order_id"
	KEY_AUTO_USERID                string = "user_id"
	KEY_WALLET_ADDRESS             string = "wallet_address"
	KEY_WALLET_ADDRESS_BTC         string = "wallet_address_btc"
	KEY_WALLET_ADDRESS_BTC_TAPROOT string = "wallet_address_btc_taproot"
	KEY_DELETED_AT                 string = "deleted_at"
	KEY_PROJECT_ID                 string = "project_id"
	KEY_LISTING_CONTRACT           string = "collection_contract"
	KEY_BTC_WALLET_INFO            string = "btc_wallet_info"

	COLLECTION_USERS                    string = "users"
	COLLECTION_USER_VOLUMN              string = "user_volumn"
	COLLECTION_WITHDRAW                 string = "withdraw"
	COLLECTION_TOKEN_URI                string = "token_uri"
	COLLECTION_TOKEN_URI_HISTORIES      string = "token_uri_histories"
	COLLECTION_FILES                    string = "files"
	COLLECTION_PROJECTS                 string = "projects"
	COLLECTION_CONFIGS                  string = "configs"
	COLLECTION_CATEGORIES               string = "categories"
	COLLECTION_ACTIVITIES               string = "activities"
	COLLECTION_REFERRALS                string = "referrals"
	COLLECTION_MARKETPLACE_LISTINGS     string = "marketplace_listings"
	COLLECTION_MARKETPLACE_OFFERS       string = "marketplace_offers"
	COLLECTION_DAO_PROPOSAL             string = "proposals"
	COLLECTION_DAO_PROPOSAL_DETAIL      string = "proposal_detail"
	COLLECTION_LEADERBOARD_TOKEN_HOLDER string = "token_holders"
	COLLECTION_DAO_PROPOSAL_VOTES       string = "proposal_votes"
	COLLECTION_BTC_WALLET_ADDRESS       string = "btc_wallet_address"
	INSCRIBE_BTC                        string = "inscribe_btc"
	INSCRIBE_INFO                       string = "inscribe_infos"
	COLLECTION_ETH_WALLET_ADDRESS       string = "eth_wallet_address"
	COLLECTION_MARKETPLACE_BTC_LISTING  string = "marketplace_btc_listing"
	COLLECTION_MARKETPLACE_BTC_BUY      string = "marketplace_btc_buy"
	COLLECTION_MARKETPLACE_BTC_LOGS     string = "marketplace_btc_logs"
	COLLECTION_COLLECTION_META          string = "collection_metas"
	COLLECTION_COLLECTION_INSCRIPTION   string = "collection_inscriptions"
	WALLET_TRACK_TX                     string = "wallet_track_txs"
	COLLECTION_AIRDROP                  string = "airdrop"

	MINT_NFT_BTC string = "mint_nft_btc"

	REDIS_KEY_LOCK_TX_CONSUMER_CONSUMER_BLOCK string = "lock-tx-consumer-update-last-processed-block"
	EVM_NULL_ADDRESS                          string = "0x0000000000000000000000000000000000000000"
	PUBSUB_TOKEN_THUMBNAIL                    string = "token_thumbnail"
	PUBSUB_PROJECT_UNZIP                      string = "project_unzip"

	BTCConfirmationThreshold = 6
	FirstScannedBTCBlkHeight = 697200
	BUY_NFT_CHARGE           = 0      // 0%
	MIN_BTC_TO_LIST_BTC      = 500000 // 0.005 btc

	FEE_BTC_SEND_AGV = 8000 // fee send btc
	MIN_FILE_SIZE    = 4096 // min file size (for linux system)

	INSCRIBE_TIMEOUT = 6

	MASTER_ADDRESS = "bc1p8ts7h86jgduat5v98cwlurngeyasqrd5c6ch2my8qwen3ykpagyswv2sy8"

	NETWORK_BTC = "btc"
	NETWORK_ETH = "eth"

	AIRDROP_MAGIC  = "https://storage.googleapis.com/generative-static-prod/airdrop/1.txt"
	AIRDROP_GOLDEN = "https://storage.googleapis.com/generative-static-prod/airdrop/1.txt"
	AIRDROP_SILVER = "https://storage.googleapis.com/generative-static-prod/airdrop/1.txt"
)

type PubSubSendOtp struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	AppName string `json:"app_name"`
}
