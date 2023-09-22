package opensea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rederinghub.io/utils/redis"
)

type OpenseaService struct {
	//apiKey    string
	cache redis.IRedisCache
}

func (openseaService OpenseaService) GetProfileAvatar(addr string) (string, error) {
	data, _, err := openseaService.request(fmt.Sprintf("https://api.opensea.io/api/v1/user/%s", addr), "GET", nil, nil)
	if err != nil {
		return "", err
	}
	resp := User{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return "", err
	}

	return resp.Account.ProfileImgUrl, nil
}

func (openseaService OpenseaService) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, int, error) {

	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, 0, err
	}

	if len(headers) > 0 {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	// remove this for error response:
	// if res.StatusCode != http.StatusOK {
	// 	err = errors.New(fmt.Sprintf("Response with status %d", res.StatusCode))
	// 	return nil, statusCode, err
	// }

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return body, res.StatusCode, nil
}
