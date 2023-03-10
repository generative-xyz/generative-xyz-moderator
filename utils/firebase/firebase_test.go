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
	err = client.SendMessagesToSpecificDevices(context.Background(), "c0fGiOSasYEI3FSAz9KpU2:APA91bEEnY0ZfNKjz2dE67M6W2ofpgeHmDbLb8hQj2cy61vr8ZuUrHPtYy0N5IbLKbR8o0E7DWpxzjn4-F9o5S2e-F0ZzU2YRAfCTEUiqSEudglpyP4EAlcwkc85GRsABv2fEPcgY2wJ",
		map[string]string{
			"test": "test",
			"data": "test",
		},
	)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}
