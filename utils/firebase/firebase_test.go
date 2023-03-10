package firebase_test

import (
	"context"
	"testing"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/firebase"
)

var client firebase.Service
var err error

func init() {
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	client, err = firebase.NewService(c.Gcs.Auth)
	if err != nil {
		panic(err)
	}
}
func Test_SendMessagesToSpecificDevices(t *testing.T) {
	err = client.SendMessagesToSpecificDevices(context.Background(), "",
		map[string]string{
			"test": "test",
		},
	)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}
