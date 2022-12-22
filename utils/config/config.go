package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

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

func NewConfig() (*Config, error) {
	godotenv.Load()
	services := make(map[string]string)
	isDebug,  _ := strconv.ParseBool(os.Getenv("DEBUG"))
	timeOut, err := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	if err != nil {
		panic(err)
	}
	
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
	}

	c, _ := json.MarshalIndent(conf, "", "\t")
	fmt.Println("Config:", string(c))

	return conf, nil
}
