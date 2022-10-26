package adapter

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

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
	Video string `json:"video"`
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
	res, err := http.Post(a.Address, "application/json", bytes.NewBuffer(_bytes))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("call to render machine got error")
	}
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
