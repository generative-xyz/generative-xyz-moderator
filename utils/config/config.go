package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Debug     bool
	Context   *Context
	Databases *Databases
	Sentry    *Sentry
	Redis     RedisConfig
	ENV string
	ServicePort string
	SigningKey string
	Services map[string]string
	MQTTConfig     MQTTConfig
	Gcs       *GCS
	Moralis MoralisConfig
	BlockchainConfig BlockchainConfig
	TxConsumerConfig TxConsumerConfig
}

type MQTTConfig struct {
	Address  string
	Port string
	UserName       string
	Password      string
}

type MoralisConfig struct {
	Key  string
	URL string
	Chain string
}


type Context struct {
	TimeOut int
}

type Databases struct {
	Postgres *DBConnection
	Mongo *DBConnection
}

type DBConnection struct {
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	Sslmode string
}

type Mongo struct {
	DBConnection
}

type Sentry struct {
	Dsn   string
	Env   string
	Debug bool
}

type Services struct {
	Name string
	Url  string
}

type GCS struct {
	ProjectId string
	Bucket    string
	Auth      string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       string
	ENV      string
}


type Chain struct {
	ID  int
	Name string
	FullName string
	Currency       string
	CurrencyLogo       string
}

type BlockchainConfig struct {
	ETHEndpoint string
}

type TxConsumerConfig struct {
	Enabled bool
	StartBlock int64
	CronJobPeriod int32
	BatchLogSize int32
	Addresses []string
}

func NewConfig() (*Config, error) {
	godotenv.Load()
	services := make(map[string]string)
	isDebug,  _ := strconv.ParseBool(os.Getenv("DEBUG"))
	timeOut, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	if err != nil {
		panic(err)
	}

	// tx consumer config
	enabled, _ := strconv.ParseBool(os.Getenv("TX_CONSUMER_ENABLED"))
	startBlock, _ := strconv.Atoi(os.Getenv("TX_CONSUMER_START_BLOCK"))
	cronJobPeriod, _ := strconv.Atoi(os.Getenv("TX_CONSUMER_CRON_JOB_PERIOD"))
	batchLogSize, _ := strconv.Atoi(os.Getenv("TX_CONSUMER_BATCH_LOG_SIZE"))
	addresses := strings.Split(os.Getenv("TX_CONSUMER_ADDRESSES"), ",")

	services["og"] = os.Getenv("OG_SERVICE_URL")
	conf := &Config{
		ENV:         os.Getenv("ENV"),
		Context:         &Context{
			TimeOut: timeOut,
		},
		Debug:        isDebug ,
		ServicePort: os.Getenv("SERVICE_PORT"),
		Databases: &Databases{
			Mongo: &DBConnection{
				Host:     os.Getenv("MONGO_HOST"),
				Port:     os.Getenv("MONGO_PORT"),
				User:     os.Getenv("MONGO_USER"), 
				Pass:     os.Getenv("MONGO_PASSWORD"),    
				Name :     os.Getenv("MONGO_DB"),   
			},
		},
		Redis: RedisConfig{
			Address:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       os.Getenv("REDIS_DB"),
			ENV:       os.Getenv("REDIS_ENV"),
		},
		SigningKey: os.Getenv("AUTH_SECRET_KEY"),
		Services: services,
		MQTTConfig: MQTTConfig{
			Address:     os.Getenv("MQTT_ADDR"),
			Port:     os.Getenv("MQTT_PORT"),
			UserName:     os.Getenv("MQTT_USERNAME"),
			Password:     os.Getenv("MQTT_PASSWORD"),
		},
		Gcs: &GCS{
			ProjectId: os.Getenv("GCS_PROJECT_ID"),
			Bucket: os.Getenv("GCS_BUCKET"),
			Auth: os.Getenv("GCS_AUTH"),
		},
		Moralis: MoralisConfig{
			Key: os.Getenv("MORALIS_KEY"),
			URL: os.Getenv("MORALIS_API_URL"),
			Chain: os.Getenv("MORALIS_CHAIN"),
		},
		BlockchainConfig: BlockchainConfig{
			ETHEndpoint: os.Getenv("ETH_ENDPOINT"),
		},
		TxConsumerConfig: TxConsumerConfig{
			Enabled: enabled,
			StartBlock: int64(startBlock),
			CronJobPeriod: int32(cronJobPeriod),
			BatchLogSize: int32(batchLogSize),
			Addresses: addresses,
		},
	}

	c, _ := json.MarshalIndent(conf, "", "\t")
	fmt.Println("Config:", string(c))

	return conf, nil
}
