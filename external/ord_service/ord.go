package ord_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"rederinghub.io/utils/config"
	"rederinghub.io/utils/redis"
	"rederinghub.io/utils/tracer"

	"github.com/davecgh/go-spew/spew"
	"github.com/opentracing/opentracing-go"
)

type BtcOrd struct {
	conf            *config.Config
	tracer          tracer.ITracer
	rootSpan opentracing.Span
	serverURL string
	cache redis.IRedisCache
}

func NewBtcOrd(conf *config.Config, t tracer.ITracer, cache redis.IRedisCache) *BtcOrd {

	serverURL := os.Getenv("ORD_SERVER")
    return &BtcOrd{
		conf:            conf,
		tracer:          t,
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
	spew.Dump(string(data))
	resp := &ExecRespose{}
	err = json.Unmarshal(data, resp)
	if err != nil {
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
