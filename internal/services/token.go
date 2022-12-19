package services

import (
	"context"

	"rederinghub.io/api"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetToken(ctx context.Context, req *api.GetTokenMessageReq) (*api.GetTokenMessageResp, error) {
	logger.AtLog.Infof("Handle [GetToken] %s %s", req.ContractAddr, req.TokenId)

	resp := &api.GetTokenMessageResp{}
	return resp, nil
}
