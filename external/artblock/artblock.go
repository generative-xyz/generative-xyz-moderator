package artblock

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"rederinghub.io/utils/config"
)

type ArtBlockService struct {
	conf *config.Config

	serverURL string
}

func NewArtBlockService(conf *config.Config, serverURL string) *ArtBlockService {

	if len(serverURL) == 0 {
		serverURL = os.Getenv("ARTBLOCK_SERVER")
	}

	return &ArtBlockService{
		conf:      conf,
		serverURL: serverURL,
	}
}

func (m ArtBlockService) GetArtist(publicAddress string) (*GetArtists, error) {
	format := `{
		"query": "query GetArtists($limit: Int!, $where: artists_bool_exp!, $orderBy: [artists_order_by!], $offset: Int!, $ttl: Int!) @cached(ttl: $ttl) {\n  artists(limit: $limit, where: $where, order_by: $orderBy, offset: $offset) {\n    ...ArtistItemInfo\n    __typename\n  }\n}\nfragment ArtistItemInfo on artists {\n  public_address\n  user {\n    display_name\n    is_curated\n    profile {\n      name\n      username\n      bio\n      profile_picture {\n        file_path\n        metadata\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n  projects_aggregate(where: {vertical: {category: {hosted: {_eq: true}}}}) {\n    aggregate {\n      count\n      __typename\n    }\n    __typename\n  }\n  most_recent_hosted_project {\n    id\n    project_id\n    artist_name\n    disable_auto_image_format\n    tokens(\n      order_by: {invocation: desc}\n      where: {image_id: {_is_null: false}}\n      limit: 1\n    ) {\n      id\n      token_id\n      image {\n        url\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}",
		"operationName": "GetArtists",
		"variables": {
			"orderBy": [
				{
					"most_recent_hosted_project": {
						"start_datetime": "desc_nulls_last"
					}
				}
			],
			"offset": 0,
			"where": {
				"most_recent_hosted_project": {
					"name": {
						"_is_null": false
					}
				},
				"public_address": {"_eq": "%s"}
			},
			"limit": 100,
			"ttl": 300
		}
	}`
	f := fmt.Sprintf(format, publicAddress)
	fBytes := []byte(f)
	temp := make(map[string]interface{})
	err2 := json.Unmarshal(fBytes, &temp)
	if err2 != nil {
		return nil, err2
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(temp)
	if err != nil {
		return nil, err
	}
	url := m.generateUrl("/v1/graphql")
	data, err := m.request(url, "POST", nil, &buf)
	if err != nil {
		return nil, err
	}
	resp := &GetArtists{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m ArtBlockService) generateUrl(path string) string {
	fullUrl := fmt.Sprintf("%s/%s", m.serverURL, path)
	return fullUrl
}

func (m ArtBlockService) request(fullUrl string, method string, headers map[string]string, reqBody io.Reader) ([]byte, error) {
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
