package utils

import (
	"time"

	"github.com/jinzhu/now"
)

type QuerySort struct {
	Sort   int
	SortBy string
}

type AggregateDexBTCListing struct {
	FromDate time.Time
	ToDate   time.Time
}

func ParseAggregation(key string) AggregateDexBTCListing {
	sortParams := make(map[string]AggregateDexBTCListing)
	to := time.Now().UTC()
	sortParams["week"] = AggregateDexBTCListing{FromDate: now.BeginningOfDay().AddDate(0, 0, -7), ToDate: to}
	sortParams["month"] = AggregateDexBTCListing{FromDate: now.BeginningOfDay().AddDate(0, 0, -30), ToDate: to}
	filter, ok := sortParams[key]
	if !ok {
		return sortParams["custom"]
	}
	return filter
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

func ParseSortNew(key string) QuerySort {
	sortParams := make(map[string]QuerySort)
	sortParams["custom"] = QuerySort{Sort: -1, SortBy: ""}
	sortParams["newest"] = QuerySort{Sort: -1, SortBy: "created_at"}
	sortParams["oldest"] = QuerySort{Sort: 1, SortBy: "created_at"}
	sortParams["priority-asc"] = QuerySort{Sort: 1, SortBy: "priority"}
	sortParams["priority-desc"] = QuerySort{Sort: -1, SortBy: "priority"}
	sortParams["price-asc"] = QuerySort{Sort: 1, SortBy: "priceBTC"}
	sortParams["price-desc"] = QuerySort{Sort: -1, SortBy: "priceBTC"}
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
	REDIS_INSCRIPTION         string = "inscription"
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
	DB_CACHE_EXPIRED_TIME            int    = 86400            //a day
	REDIS_CACHE_EXPIRED_TIME         int    = 86400            //a day
	DB_CACHE_KEY                     string = "db.cache.%s.%s" //a week
	NONCE_MESSAGE_FORMAT             string = "%s"

	KEY_UUID                       string = "uuid"
	KEY_BASE_PRODUCT_KEY           string = "product_key"
	KEY_ORDER_ID                   string = "order_id"
	KEY_AUTO_USERID                string = "user_id"
	KEY_WALLET_ADDRESS             string = "wallet_address"
	KEY_WALLET_ADDRESS_BTC         string = "wallet_address_btc"
	KEY_WALLET_ADDRESS_BTC_TAPROOT string = "wallet_address_btc_taproot"
	KEY_WALLET_SLUG                string = "slug"
	KEY_DELETED_AT                 string = "deleted_at"
	KEY_PROJECT_ID                 string = "project_id"
	KEY_LISTING_CONTRACT           string = "collection_contract"
	KEY_BTC_WALLET_INFO            string = "btc_wallet_info"

	COLLECTION_USERS                           string = "users"
	COLLECTION_USER_VOLUMN                     string = "user_volumn"
	COLLECTION_WITHDRAW                        string = "withdraw"
	COLLECTION_TOKEN_URI                       string = "token_uri"
	COLLECTION_MODULAR_INSCRIPTION             string = "modular_inscription"
	COLLECTION_TOKEN_URI_HISTORIES             string = "token_uri_histories"
	COLLECTION_TOKEN_URI_METADATA              string = "token_uri_metadata"
	COLLECTION_FILES                           string = "files"
	COLLECTION_PROJECTS                        string = "projects"
	COLLECTION_PROJECT_PROTAB                  string = "projects_protab" // protab of generative
	COLLECTION_PROJECT_ALLOW_LIST              string = "project_allow_list"
	COLLECTION_PROJECT_ZIPLINKS                string = "project_zip_links"
	COLLECTION_CONFIGS                         string = "configs"
	COLLECTION_ACTIONS                         string = "actions"
	COLLECTION_CATEGORIES                      string = "categories"
	COLLECTION_ACTIVITIES                      string = "activities"
	COLLECTION_REFERRALS                       string = "referrals"
	COLLECTION_MARKETPLACE_LISTINGS            string = "marketplace_listings"
	COLLECTION_CACHED_GM_DASHBOARD             string = "cached_gm_dashboard"
	COLLECTION_CACHED_GM_DASHBOARD_NEW         string = "cached_gm_dashboard_new"
	COLLECTION_CACHED_REALLOCATED_GM_DASHBOARD string = "cached_gm_reallowcacted_dashboard"
	COLLECTION_MARKETPLACE_OFFERS              string = "marketplace_offers"
	COLLECTION_DAO_PROPOSAL                    string = "proposals"
	COLLECTION_DAO_PROPOSAL_DETAIL             string = "proposal_detail"
	COLLECTION_LEADERBOARD_TOKEN_HOLDER        string = "token_holders"
	COLLECTION_DAO_PROPOSAL_VOTES              string = "proposal_votes"
	COLLECTION_BTC_WALLET_ADDRESS              string = "btc_wallet_address"
	INSCRIBE_BTC                               string = "inscribe_btc"
	INSCRIBE_INFO                              string = "inscribe_infos"
	COLLECTION_ETH_WALLET_ADDRESS              string = "eth_wallet_address"
	COLLECTION_MARKETPLACE_BTC_LISTING         string = "marketplace_btc_listing"
	COLLECTION_MARKETPLACE_BTC_BUY             string = "marketplace_btc_buy"
	COLLECTION_MARKETPLACE_BTC_LOGS            string = "marketplace_btc_logs"
	COLLECTION_COLLECTION_META                 string = "collection_metas"
	COLLECTION_COLLECTION_INSCRIPTION          string = "collection_inscriptions"
	WALLET_TRACK_TX                            string = "wallet_track_txs"
	COLLECTION_AIRDROP                         string = "airdrop"
	COLLECTION_DEX_BTC_LISTING                 string = "dex_btc_listing"
	COLLECTION_DISCORD_NOTI                    string = "discord_notis"
	COLLECTION_DEX_BTC_BUY_ETH                 string = "dex_btc_buy_eth"
	COLLECTION_BTC_TX_SUBMIT                   string = "btc_tx_submit"
	COLLECTION_TOKEN_ACTIVITY                  string = "token_activities"
	COLLECTION_DISCORD_PARTNER                 string = "discord_partners"
	COLLECTION_DISCORD_CUTOM_TYPES             string = "discord_custom_types"
	COLLECTION_TOKEN_TX                        string = "token_txs"
	COLLECTION_DEX_BTC_TRACKING_INTERNAL       string = "dex_btc_tracking_internal"
	COLLECTION_GLOBAL_VARIABLE                 string = "global_variables"
	COLLECTION_SORALIS_SNAPSHOT_BALANCE        string = "soralis_snapshot_balance"
	AI_SCHOOL_JOB                              string = "ai_school_jobs"
	AI_SCHOOL_DATASET                          string = "ai_school_dataset"
	TOKEN_FILE_FRAGMENT                        string = "token_file_fragments"
	TOKEN_FILE_FRAGMENT_JOB                    string = "token_file_fragment_jobs"

	MINT_NFT_BTC string = "mint_nft_btc"

	REDIS_KEY_LOCK_TX_CONSUMER_CONSUMER_BLOCK string = "lock-tx-consumer-update-last-processed-block"
	EVM_NULL_ADDRESS                          string = "0x0000000000000000000000000000000000000000"
	PUBSUB_TOKEN_THUMBNAIL                    string = "token_thumbnail"
	PUBSUB_CAPTURE_THUMBNAIL                  string = "capture_thumbnail"
	PUBSUB_PROJECT_UNZIP                      string = "project_unzip"
	PUBSUB_ETH_PROJECT_UNZIP                  string = "eth_project_unzip"

	BTCConfirmationThreshold = 1
	FirstScannedBTCBlkHeight = 697200
	BUY_NFT_CHARGE           = 0      // 0%
	MIN_BTC_TO_LIST_BTC      = 500000 // 0.005 btc

	FEE_BTC_SEND_AGV = 8000 // fee send btc
	MIN_FILE_SIZE    = 4096 // min file size (for linux system)

	FEE_ETH_SEND_MASTER = 0.0007
	FEE_BTC_SEND_NFT    = 10000

	DEVELOPER_INSCRIBE_MAX_REQUEST = 200

	INSCRIBE_TIMEOUT = 3

	MASTER_ADDRESS = "bc1p8ts7h86jgduat5v98cwlurngeyasqrd5c6ch2my8qwen3ykpagyswv2sy8"

	NETWORK_BTC = "btc"
	NETWORK_ETH = "eth"

	PLATFORM_ORDINAL = "ordinal"
	PLATFORM_TC      = "tc"

	GM_CRONTAB_PROCESSING_KEY = "gm_crontab_running"

	AIRDROP_MAGIC  = "https://storage.googleapis.com/generative-static-prod/airdrop/magickey.html"
	AIRDROP_GOLDEN = "https://storage.googleapis.com/generative-static-prod/airdrop/goldenkey.html"
	AIRDROP_SILVER = "https://storage.googleapis.com/generative-static-prod/airdrop/silverkey.html"
)

type PubSubSendOtp struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	AppName string `json:"app_name"`
}

const HttpContextTimeOut = time.Second * 15

const (
	BidProjectIDProd = "1002573"
	BidProjectIDDev  = "1000362"
)

var ExceptionProjectContract = map[string]string{
	"0x9841faa1133da03b9ae09e8daa1a725bc15575f0": "999998",
	"0xda00b6a8b521113501bb98fd0a7ffcfe756d9962": "999997",
}
