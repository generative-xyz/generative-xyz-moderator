package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
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

func PostRequestWithHeaders(url string, headers map[string]string, body interface{}) (int, []byte, error) {

	data, err := json.Marshal(body)
	client := &http.Client{}

	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return http.StatusBadRequest, nil, err
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
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

func PostToRenderer(url string, body interface{}) (int, []byte, error) {
	headers := make(map[string]string)
	headers["Content-type"] = "application/json"
	headers["X-Password"] = os.Getenv("RENDER_API_KEY")

	return PostRequestWithHeaders(url, headers, body)
}

func GetRequest(url string) (int, []byte, error) {
	response, err := http.Get(url)
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
