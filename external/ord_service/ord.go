package ord_service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
)

type BtcOrd struct {
	conf            *config.Config
	
	serverURL string
	cache redis.IRedisCache
}

func NewBtcOrd(conf *config.Config,  cache redis.IRedisCache) *BtcOrd {

	serverURL := os.Getenv("ORD_SERVER")
    return &BtcOrd{
		conf:            conf,
		serverURL: serverURL,
		cache: cache,
	}
}

type metadataChan struct {
	Key int
	Err error
}

func (m BtcOrd) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}

func (m BtcOrd) Exec(f ExecRequest) (*ExecRespose, error){
	url := fmt.Sprintf("%s", Exec)
	fullUrl := m.generateUrl(url)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f)
	if err != nil {
			return nil, err
	}

	data, err := m.request(fullUrl, "POST", nil, &buf)
	if err != nil {
		return nil, err
	}
resp := &ExecRespose{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		err = errors.New(resp.Error)
		return nil, err
	}

	return resp, nil
}


func (m BtcOrd) Mint(f MintRequest) (*MintRespose, error){
	url := fmt.Sprintf("%s", Inscribe)
	fullUrl := m.generateUrl(url)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f)
	if err != nil {
			return nil, err
	}

	data, err := m.request(fullUrl, "POST", nil, &buf)
	if err != nil {
		return nil, err
	}
resp := &MintRespose{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		err = errors.New(resp.Error)
		return nil, err
	}

	return resp, nil
}

func (m BtcOrd) request(fullUrl string, method string, headers map[string]string , reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, err
	}

	if len(headers) > 0 {
		for key, val := range headers{
			req.Header.Add(key,  val)
		}
	}
req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Response with status %d", res.StatusCode))
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

