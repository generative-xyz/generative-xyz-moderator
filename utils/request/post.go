package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func PostRequest(url string, body interface{}) (int, []byte, error) {

	data, err := json.Marshal(body)

	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, nil, err
	}

	return response.StatusCode, responseBody, nil
}
