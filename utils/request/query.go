package request

import (
	"encoding/json"
	"net/http"
)

func Queries(r *http.Request) map[string]string {
	query := r.URL.Query()
	result := make(map[string]string)
	for key, values := range query {
		result[key] = values[0]
	}
	return result
}

func Query(r *http.Request, key, defaultValue string) string {
	query := r.URL.Query()
	if value, ok := query[key]; ok {
		return value[0]
	}
	return defaultValue
}

func BindJson(r *http.Request, payload interface{}) error {

	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		return err
	}
	return nil
}
