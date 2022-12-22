package nfts

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
)

type MoralisNfts struct {
	conf            *config.Config
	tracer          tracer.ITracer
	rootSpan opentracing.Span
	serverURL string
	apiKey string
	//client forwarder.IForwarder
	cache redis.IRedisCache
}

func NewMoralisNfts(conf *config.Config, t tracer.ITracer, cache redis.IRedisCache) *MoralisNfts {
	
	apiKey := conf.Moralis.Key
	serverURL := conf.Moralis.URL
    return &MoralisNfts{
		conf:            conf,
		tracer:          t,
		serverURL: serverURL,
		apiKey: apiKey,
		cache: cache,
	}
}

type metadataChan struct {
	Key int
	Err error
}

func (m MoralisNfts) generateUrl(path string, filters *MoralisFilter) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	if filters != nil {
		params := url.Values{}
		
		if filters.Chain != nil {
			params[KeyChain] = []string{
				*filters.Chain,
			}
		}else{
			params[KeyChain] = []string{
				m.conf.Moralis.Chain,
			}
		}
		
		if filters.Format != nil { 
			params[KeyFormat] = []string{
				*filters.Format,
			}
		}
		
		if filters.Limit != nil { 
			params[KeyLimit] = []string{
				strconv.Itoa(*filters.Limit),
			}
		}
		
		if filters.TotalRanges != nil { 
			params[KeyTotalRanges] = []string{
				strconv.Itoa(*filters.TotalRanges),
			}
		}
		
		if filters.Range != nil { 
			params[KeyRange] = []string{
				strconv.Itoa(*filters.Range),
			}
		}
		
		if filters.Cursor != nil { 
			params[KeyCurrsor] = []string{
				*filters.Cursor,
			}
		}

		fullUrl = fullUrl + "?" + params.Encode()
	}

	return fullUrl
}

func (m MoralisNfts) request(fullUrl string, method string, headers map[string]string , reqBody io.Reader) ([]byte, error) {
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
	req.Header.Add("X-API-Key", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (m MoralisNfts) GetNftByContract(contractAddr string,f MoralisFilter) (*MoralisTokensResp, error){
	url := fmt.Sprintf("%s/%s", URLNft, contractAddr )
	fullUrl := m.generateUrl(url, &f)

	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	resp := &MoralisTokensResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}