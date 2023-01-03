package utils

const (
	API_KEY string = "Api-Key"
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

	EMAIL_TAG          string = "email"
	WALLET_ADDRESS_TAG string = "wallet_address"

	PubsubCreateDeviceType           string = "Device:PubsubCreateDeviceType"
	PubsubUpdateDeviceType           string = "Device:PubsubUpdateDeviceType"
	PubsubDeleteDeviceType           string = "Device:PubsubDeleteDeviceType"
	PubsubSendMessageToSlack         string = "Device:PubsubSendMessageToSlack"
	PUBSUB_SEND_OTP                  string = "Hybrid:SendOtp"
	PUBSUB_REGISTER                  string = "WorkspaceGateway::PubsubRegister"
	PUBSUB_FORGOT_PASSWORD           string = "Hybrid:ResetPasswordEmail"
	NFT_CACHE_EXPIRED_TIME           int    = 86400
	TOKEN_CACHE_EXPIRED_TIME         int    = 86400 * 30           //a month (second)
	REFRESH_TOKEN_CACHE_EXPIRED_TIME int    = 86400 * 360          //a year (second)
	DB_CACHE_EXPIRED_TIME            int    = 86400                //a week
	DB_CACHE_KEY                     string = "object_cache_%s_%s" //a week
	NONCE_MESSAGE_FORMAT             string = "Welcome %s to Generative"

	KEY_UUID             string = "uuid"
	KEY_BASE_PRODUCT_KEY string = "product_key"
	KEY_ORDER_ID         string = "order_id"
	KEY_AUTO_USERID      string = "user_id"
	KEY_WALLET_ADDRESS   string = "wallet_address"
	KEY_DELETED_AT       string = "deleted_at"

	COLLECTION_USERS string = "users"
	COLLECTION_TOKEN_URI string = "token_uri"
	COLLECTION_FILES string = "files"
	COLLECTION_PROJECTS string = "projects"
	COLLECTION_CONFIGS string = "configs"
	COLLECTION_CATEGORIES string = "categories"
)

type PubSubSendOtp struct {
	Email   string `json:"email"`
	Code    string `json:"code"`
	AppName string `json:"app_name"`
}
