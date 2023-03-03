package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"rederinghub.io/utils/helpers"
)

type DeveloperKey struct {
	BaseEntity     `bson:",inline"`
	UserUuid       string `bson:"user_uuid"`
	ApiKey         string `bson:"api_key"`
	ApiName        string `bson:"api_name"`
	ApiEmail       string `bson:"api_email"`
	ApiCompany     string `bson:"api_company"`
	ApiDescription string `bson:"api_description"`
	Status         int    `bson:"status"`
}

func (u DeveloperKey) TableName() string {
	return "developer_key"
}

func (u DeveloperKey) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

//
type DeveloperKeyRequests struct {
	BaseEntity `bson:",inline"`
	RecordId   string `bson:"record_id"`
	ApiKey     string `bson:"api_key"`

	EndpointName string `bson:"endpoint_name"`
	EndpointUrl  string `bson:"endpoint_url"`

	Status interface{} `bson:"status"`

	DayReqLastTime *time.Time `bson:"day_req_last_time"`
	DayReqCounter  int        `bson:"day_req_counter"`
}

func (u DeveloperKeyRequests) TableName() string {
	return "developer_key_requests"
}

func (u DeveloperKeyRequests) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}

//
type DeveloperKeyReqLogs struct {
	BaseEntity  `bson:",inline"`
	RecordId    string      `bson:"record_id"`
	Endpoint    string      `bson:"endpoint"`
	Status      int         `bson:"status"`
	RequestMsg  interface{} `bson:"request_msg"`
	ResponseMsg interface{} `bson:"response_msg"`
}

func (u DeveloperKeyReqLogs) TableName() string {
	return "developer_key_req_logs"
}
func (u DeveloperKeyReqLogs) ToBson() (*bson.D, error) {
	return helpers.ToDoc(u)
}
