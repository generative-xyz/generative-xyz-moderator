package usecase

import (
	"go.uber.org/zap"
	"rederinghub.io/external/token_explorer"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/utils/logger"
)

func (u *Usecase) SoralisAllTokens(query request.PaginationReq) ([]token_explorer.Token, error) {
	params := query.ToNFTServiceUrlQuery()
	tokens, err := u.TokenExplorer.Tokens(params)
	if err != nil {
		logger.AtLog.Logger.Error("Tokens() failed", zap.Error(err))
		return nil, err
	}
	return tokens, nil
}
