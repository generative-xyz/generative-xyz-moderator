package generativeexplorer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

type GenerativeExplorer struct {
	serverURL string
	cache     redis.IRedisCache
}

func NewGenerativeExplorer(cache redis.IRedisCache) *GenerativeExplorer {

	serverURL := os.Getenv("GENERATIVE_EXPLORER_API")
	return &GenerativeExplorer{
		serverURL: serverURL,
		cache:     cache,
	}
}

type metadataChan struct {
	Key int
	Err error
}

func (m GenerativeExplorer) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}


func (m GenerativeExplorer) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, error) {
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

func (m GenerativeExplorer) Inscription(inscription string) (*InscriptioResponse, error) {
	url := fmt.Sprintf("%s/%s", INSCRIPTION, inscription)
	fullUrl := m.generateUrl(url)
	//m.cache.Delete(helpers.GenerateCachedInscriptionOnweKey(inscription))

	cached, err := m.cache.GetData(helpers.GenerateCachedInscriptionOnweKey(inscription))
	if err != nil || cached == nil {
		data, err := m.request(fullUrl, "GET", nil, nil)
		if err != nil {
			return nil, err
		}
		resp := &InscriptioResponse{}
		err = json.Unmarshal(data, resp)
		if err != nil {
			return nil, err
		}
	
		m.cache.SetStringData(helpers.GenerateCachedInscriptionOnweKey(resp.InscriptionID), resp.Address)
	
		return resp, nil
	}
	
	resp := &InscriptioResponse{
		InscriptionID: inscription,
		Address:  *cached,
	}
	return resp, nil
}
