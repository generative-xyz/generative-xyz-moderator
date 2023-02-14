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
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/opentracing/opentracing-go"
)

type MoralisNfts struct {
	conf      *config.Config
	tracer    tracer.ITracer
	rootSpan  opentracing.Span
	serverURL string
	apiKey    string
	//client forwarder.IForwarder
	cache redis.IRedisCache
}

func NewMoralisNfts(conf *config.Config, t tracer.ITracer, cache redis.IRedisCache) *MoralisNfts {

	apiKey := conf.Moralis.Key
	serverURL := conf.Moralis.URL
	return &MoralisNfts{
		conf:      conf,
		tracer:    t,
		serverURL: serverURL,
		apiKey:    apiKey,
		cache:     cache,
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
		} else {
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
			if *filters.Limit != 0 {
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

		if filters.TokenAddresses != nil {
			tokenAddresses := *filters.TokenAddresses
			if len(tokenAddresses) > 0 {
				params[KeyTokenAddresses] = tokenAddresses
			}
		}

		fullUrl = fullUrl + "?" + params.Encode()
	}

	return fullUrl
}

func (m MoralisNfts) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, error) {
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

func (m MoralisNfts) GetNftByContract(contractAddr string, f MoralisFilter) (*MoralisTokensResp, error) {
	url := fmt.Sprintf("%s/%s", URLNft, contractAddr) // Todo: review this url
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

func (m MoralisNfts) GetNftByWalletAddress(wallletAddress string, filter MoralisFilter) (*MoralisTokensResp, error) {
	url := fmt.Sprintf("%s/%s", wallletAddress, URLNft)
	fullUrl := m.generateUrl(url, &filter)

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

func (m MoralisNfts) GetMultipleNfts(f MoralisGetMultipleNftsFilter) ([]MoralisToken, error) {
	url := fmt.Sprintf("%s/%s", URLNft, "getMultipleNFTs")
	fullUrl := m.generateUrl(url, &MoralisFilter{Chain: f.Chain})
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(f.ReqBody)
	if err != nil {
		return nil, err
	}

	data, err := m.request(fullUrl, "POST", nil, &buf)
	if err != nil {
		return nil, err
	}

	resp := []MoralisToken{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		messageResp := &MoralisMessage{}
		err = json.Unmarshal(data, &messageResp)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(messageResp.Message)
	}

	return resp, nil
}

func (m MoralisNfts) GetNftByContractAndTokenID(contractAddr string, tokenID string) (*MoralisToken, error) {
	key := helpers.NftFromMoralisKey(contractAddr, tokenID)
	liveReload := func(contractAddr string, tokenID string) (*MoralisToken, error) {
		nfts, err := m.GetMultipleNfts(MoralisGetMultipleNftsFilter{
			Chain: nil,
			ReqBody: MoralisGetMultipleNftsReqBody{
				Tokens: []NftFilter{
					{
						TokenAddress: contractAddr,
						TokenId:      tokenID,
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

	go liveReload(contractAddr, tokenID)

	cached, err := m.cache.GetData(key)
	if cached == nil || err != nil {
		return liveReload(contractAddr, tokenID)
	}

	resp := &MoralisToken{}
	err = helpers.ParseCache(cached, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m MoralisNfts) GetNftByContractAndTokenIDNoCahe(contractAddr string, tokenID string) (*MoralisToken, error) {
	nfts, err := m.GetMultipleNfts(MoralisGetMultipleNftsFilter{
		Chain: nil,
		ReqBody: MoralisGetMultipleNftsReqBody{
			Tokens: []NftFilter{
				{
					TokenAddress: contractAddr,
					TokenId:      tokenID,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(nfts) != 1 {
		return nil, errors.New("cannot find moralis token")
	}

	nft := nfts[0]
	return &nft, nil
}
