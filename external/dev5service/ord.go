package dev5service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
)

type Dev5Service struct {
	conf *config.Config

	serverURL string
	cache     redis.IRedisCache
}

func NewDev5Service(conf *config.Config, cache redis.IRedisCache) *Dev5Service {

	return &Dev5Service{
		conf:      conf,
		cache:     cache,
	}
}

type metadataChan struct {
	Key int
	Err error
}

func (m Dev5Service) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", URL, path)
	return fullUrl
}

func (m Dev5Service) Inscriptions(limit string) (*InscriptionsResp, error) {
	url := fmt.Sprintf("%s/%s", INSCRIPTIONS, limit)
	fullUrl := m.generateUrl(url)
	spew.Dump(fullUrl)
	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &InscriptionsResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m Dev5Service) Inscription(inscriptionID string) (*InscriptionResp, error) {
	url := fmt.Sprintf("%s/%s", INSCRIPTION, inscriptionID)
	fullUrl := m.generateUrl(url)
	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	resp := &InscriptionResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m Dev5Service) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, fullUrl, reqBody)
	if err != nil {
		return nil, err
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
