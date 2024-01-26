package entity

import (
	"encoding/json"
	"testing"
)

func Test_GetRunningActivitiesByAccessToken(t *testing.T) {
	a := ModularWorkshopEntity{}
	b, _ := json.Marshal(a)
	println(string(b))
}
