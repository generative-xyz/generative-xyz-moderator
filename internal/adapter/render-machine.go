package adapter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"rederinghub.io/pkg/config"
)

type RenderMachineAdapter interface {
	Render(ctx context.Context, req *RenderRequest) (*RenderResponse, error)
}

type renderMachineAdapter struct {
	Address string
}

type RenderRequest struct {
	Script string   `json:"script"`
	Params []string `json:"params"`
	Seed   string   `json:"seed"`
}

type RenderResponse struct {
	Glb   string `json:"glb"`
	Image string `json:"image"`
}

func NewRenderMachineAdapter() RenderMachineAdapter {
	appConfig := config.AppConfig()
	return &renderMachineAdapter{
		Address: appConfig.RenderMachineAddr,
	}
}

func (a *renderMachineAdapter) Render(ctx context.Context, request *RenderRequest) (*RenderResponse, error) {
	request.Script = base64.StdEncoding.EncodeToString([]byte(request.Script))
	_bytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	req, _ := http.NewRequest("POST", a.Address, bytes.NewBuffer(_bytes))
	req.Header.Add("accept", "application/json")
	req = req.WithContext(ctxCancel)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resp RenderResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
