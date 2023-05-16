package nfts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/redis"
)

type MoralisNfts struct {
	conf      *config.Config
	serverURL string
	apiKey    string
	//client forwarder.IForwarder
	cache redis.IRedisCache
}

func NewMoralisNfts(conf *config.Config, cache redis.IRedisCache) *MoralisNfts {

	apiKey := conf.Moralis.Key
	serverURL := conf.Moralis.URL
	serverURL = "https://deep-index.moralis.io/api/v2"
	return &MoralisNfts{
		conf:      conf,
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

	apiKeys := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6Ijk1MjI3ZDJiLWRmNDktNDZiOS1iZGMxLTdkZDMyYjMyZGZhMyIsIm9yZ0lkIjoiMzI3NDI5IiwidXNlcklkIjoiMzM2NjQyIiwidHlwZUlkIjoiNGEyZTNhZTQtZDAxNy00ZTYzLWFlODgtZmE0ZWQyZGJhNDEwIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2ODM4OTgwNzIsImV4cCI6NDgzOTY1ODA3Mn0.b2xPRQJJyMd1nJog6mkoUve-S4Mh2C_tlJuw55yxPZc",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6Ijk2YzM4OTczLWZkNzUtNDRlZC1hMzg1LTQ2MjZkNjlkZDg1OSIsIm9yZ0lkIjoiMzI3NDI5IiwidXNlcklkIjoiMzM2NjQyIiwidHlwZUlkIjoiNGEyZTNhZTQtZDAxNy00ZTYzLWFlODgtZmE0ZWQyZGJhNDEwIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2ODM5NDc5NzMsImV4cCI6NDgzOTcwNzk3M30.cy77GPkYLY6-XpdIUdkv_SlxZe8Whw4ftaHPpP-j-F8",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjQxZmNhY2M5LTA3ZDctNGE3My1hN2EyLWYzNDg2MWNkNjNmZSIsIm9yZ0lkIjoiMzI3NDI5IiwidXNlcklkIjoiMzM2NjQyIiwidHlwZUlkIjoiNGEyZTNhZTQtZDAxNy00ZTYzLWFlODgtZmE0ZWQyZGJhNDEwIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2ODM5NDgwMDgsImV4cCI6NDgzOTcwODAwOH0.7oiCoODECGfvyXlpvJ8_ykryrYrj_DXVmgENhEUHFKI",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjBlYTY3YmQyLTk3M2QtNDhmYi1iZmQ4LTYyNjU5MDE3NGY3MSIsIm9yZ0lkIjoiMzI3NDI5IiwidXNlcklkIjoiMzM2NjQyIiwidHlwZUlkIjoiNGEyZTNhZTQtZDAxNy00ZTYzLWFlODgtZmE0ZWQyZGJhNDEwIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2ODQxOTY0OTAsImV4cCI6NDgzOTk1NjQ5MH0.ZS1Enk_ns8bxUI-10bHABQ8BRRAthyr-O4QRccouyXQ",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjBmYzRmZTdjLWVkMDktNGYyYy05MjFlLTgwNzg2MTVhZDQyMSIsIm9yZ0lkIjoiMzI3NDI5IiwidXNlcklkIjoiMzM2NjQyIiwidHlwZUlkIjoiNGEyZTNhZTQtZDAxNy00ZTYzLWFlODgtZmE0ZWQyZGJhNDEwIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2ODQxOTY1MTksImV4cCI6NDgzOTk1NjUxOX0.YKKXfgDevKS6skTikA5VSzxK5sgfadzwnXj6gFv0RF4",
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	req.Header.Add("X-API-Key", apiKeys[r.Intn(len(apiKeys))])

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

func (m MoralisNfts) AddressBalance(walletAddress string) (*MoralisBalanceResp, error) {
	fullUrl := m.generateUrl(fmt.Sprintf("%s/%s", walletAddress, WalletAddressBalance), &MoralisFilter{})

	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	resp := &MoralisBalanceResp{}
	if string(data) == `{"message":"Invalid key"}` {
		return nil, errors.New("invalid key")
	}
	if strings.Contains(string(data), "limit") {
		return nil, errors.New("rate limit")
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m MoralisNfts) TokenBalanceByWalletAddress(walletAddress string, tAddresses []string) (map[string]MoralisBalanceResp, error) {
	f := &MoralisFilter{}
	f.TokenAddresses = new([]string)

	urls := url.Values{}
	urls.Add("chain", m.conf.Moralis.Chain)
	for key, tAddress := range tAddresses {
		urls.Add(fmt.Sprintf("token_addresses[%d]", key), tAddress)
	}

	path := fmt.Sprintf("%s/%s?%s", walletAddress, WalletAddressTokenBalance, urls.Encode())
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)

	data, err := m.request(fullUrl, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	resp := []MoralisBalanceResp{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	result := make(map[string]MoralisBalanceResp)
	for _, i := range resp {
		result[strings.ToLower(i.TokenAddress)] = i
	}

	return result, nil
}
