package nfts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/davecgh/go-spew/spew"
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
			if  *filters.Limit != 0 {
				params[KeyLimit] = []string{
					strconv.Itoa(*filters.Limit),
				}
			}	
		}
		
		if filters.TotalRanges != nil { 
			if *filters.TotalRanges != 0 {
				params[KeyTotalRanges] = []string{
					strconv.Itoa(*filters.TotalRanges),
				}
			}
		}
		
		if filters.Range != nil { 
			if *filters.Range != 0 {
				params[KeyRange] = []string{
					strconv.Itoa(*filters.Range),
				}
			}
		}
		
		if filters.Cursor != nil { 
			if *filters.Cursor != "" {
				params[KeyCurrsor] = []string{
					*filters.Cursor,
				}
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
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-API-Key", m.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (m MoralisNfts) GetNftByContract(contractAddr string,f MoralisFilter) (*MoralisTokensResp, error){
	url := fmt.Sprintf("%s/%s", URLNft, contractAddr )
	fullUrl := m.generateUrl(url, &f)
	spew.Dump(fullUrl)

	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	spew.Dump(fullUrl)
	resp := &MoralisTokensResp{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m MoralisNfts) GetMultipleNfts(f MoralisGetMultipleNftsFilter) ([]MoralisToken, error) {
	url := fmt.Sprintf("%s/%s", URLNft, "getMultipleNFTs")
	fullUrl := m.generateUrl(url, &MoralisFilter{Chain: f.Chain})
	spew.Dump(fullUrl)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f.ReqBody)
	if err != nil {
			return nil, err
	}
	data, err := m.request(fullUrl, "POST", nil, &buf)
	if err != nil {
		return nil, err
	}
	spew.Dump(fullUrl)
	resp := make([]MoralisToken, 0);
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m MoralisNfts) GetNftByContractAndTokenID(contractAddr string, tokenID string) (*MoralisToken, error) {
	key := fmt.Sprintf("NFT_C_%s_T_%s", contractAddr, tokenID)

	cached, err := m.cache.GetData(key)
	if cached == nil || err != nil {
		nfts, err := m.GetMultipleNfts(MoralisGetMultipleNftsFilter{
			Chain: nil,
			ReqBody: MoralisGetMultipleNftsReqBody{
				Tokens: []NftFilter{
					{
						TokenAddress: contractAddr,
						TokenId: tokenID,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		if len(nfts) != 1 {
			return nil, errors.New("cannot find moralis token") 
		} else {

			nft := nfts[0]
			m.cache.SetData(key, nft)
			return &nft, nil
		}
	}
	
	resp := &MoralisToken{}
	bytes := []byte(*cached)
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
