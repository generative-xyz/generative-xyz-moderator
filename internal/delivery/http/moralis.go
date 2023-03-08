package http

import (
	"context"
	"net/http"

	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
)

// @Summary List NFT from Moralis
// @Description List NFT from Moralis
// @Tags Token Moralis
// @Accept json
// @Produce json
// @Param walletAddress query string false "Wallet Address"
// @Param cursor query string false "Last Id"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.Pagination{}
// @Router /token-moralis/nfts [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) listNftFromMoralis(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			userWallet := ctx.Value(utils.SIGNED_WALLET_ADDRESS).(string)
			userId := ctx.Value(utils.SIGNED_USER_ID).(string)
			pag := entity.GetPagination(r)
			// TODO: 0x2525
			return h.Usecase.ListNftFromMoralis(ctx, userId, userWallet, r.URL.Query().Get("walletAddress"), pag)
		},
	).ServeHTTP(w, r)
}

// @Summary NFT from Moralis
// @Description NFT from Moralis
// @Tags Token Moralis
// @Accept json
// @Produce json
// @Param tokenAddress path string false "Token Address"
// @Param tokenId query string false "Token Id"
// @Success 200 {object} nfts.MoralisToken{}
// @Router /token-moralis/nfts/{tokenAddress} [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) nftFromMoralis(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			tokenAddress := muxVars["tokenAddress"]
			return h.Usecase.NftFromMoralis(ctx, tokenAddress, r.URL.Query().Get("tokenId"))
		},
	).ServeHTTP(w, r)
}
