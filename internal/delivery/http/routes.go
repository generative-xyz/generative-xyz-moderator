package http

import (
	"net/http"
	"os"
	"strconv"

	"rederinghub.io/docs"
	_ "rederinghub.io/docs"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"

	httpSwagger "github.com/swaggo/http-swagger"
)

func (h *httpDelivery) registerRoutes() {
	h.RegisterDocumentRoutes()
	h.RegisterV1Routes()
}

func (h *httpDelivery) RegisterV1Routes() {
	h.Handler.Use(h.MiddleWare.LoggingMiddleware)
	h.Handler.HandleFunc("/", h.healthCheck).Methods("GET")

	//api
	api := h.Handler.PathPrefix("/generative/api").Subrouter()
	api.HandleFunc("/token/{contractAddress}/{tokenID}", h.tokenURI).Methods("GET")
	api.HandleFunc("/trait/{contractAddress}/{tokenID}", h.tokenTrait).Methods("GET")

	//api
	tokens := api.PathPrefix("/tokens").Subrouter()
	tokens.HandleFunc("", h.Tokens).Methods("GET")
	tokens.HandleFunc("/{tokenID}/thumbnail", h.updateTokenThumbnail).Methods("POST")
	tokens.HandleFunc("/{contractAddress}/{tokenID}", h.tokenURIWithResp).Methods("GET")
	tokens.HandleFunc("/{contractAddress}/{tokenID}", h.tokenURIWithResp).Methods("PUT")
	tokens.HandleFunc("/traits/{contractAddress}/{tokenID}", h.tokenTraitWithResp).Methods("GET")

	//v1 := api.PathPrefix("/v1").Subrouter()
	api.HandleFunc("/", h.healthCheck).Methods("GET")

	//auth
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/nonce", h.generateMessage).Methods("POST")
	auth.HandleFunc("/nonce/verify", h.verifyMessage).Methods("POST")

	files := api.PathPrefix("/files").Subrouter()
	// files.Use(h.MiddleWare.AccessToken)
	files.HandleFunc("", h.UploadFile).Methods("POST")
	files.HandleFunc("/minify", h.minifyFiles).Methods("POST")
	files.HandleFunc("/deflate", h.deflate).Methods("POST")

	files.HandleFunc("/multipart", h.CreateMultipartUpload).Methods("POST")
	files.HandleFunc("/multipart/{uploadID}", h.UploadPart).Methods("PUT")
	files.HandleFunc("/multipart/{uploadID}", h.CompleteMultipartUpload).Methods("POST")

	//profile
	profile := api.PathPrefix("/profile").Subrouter()
	profile.Use(h.MiddleWare.UserToken)
	profile.HandleFunc("/wallet/{walletAddress}", h.profileByWallet).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/nfts", h.TokensOfAProfile).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/projects", h.getProjectsByWallet).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/volume", h.getVolumeByWallet).Methods("GET")

	singedIn := api.PathPrefix("/profile").Subrouter()
	singedIn.Use(h.MiddleWare.AccessToken)
	singedIn.HandleFunc("", h.profile).Methods("GET")
	singedIn.HandleFunc("/projects", h.getUserProjects).Methods("GET")
	singedIn.HandleFunc("", h.updateProfile).Methods("PUT")
	singedIn.HandleFunc("/logout", h.logout).Methods("PUT")

	//project
	project := api.PathPrefix("/project").Subrouter()
	project.HandleFunc("", h.getProjects).Methods("GET")
	project.HandleFunc("", h.createProjects).Methods("POST")

	project.HandleFunc("/random", h.getRandomProject).Methods("GET")
	project.HandleFunc("/minted-out", h.getMintedOutProjects).Methods("GET")
	project.HandleFunc("/recent-works", h.getRecentWorksProjects).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}", h.projectDetail).Methods("GET")

	project.HandleFunc("/{contractAddress}/{projectID}", h.updateProject).Methods("PUT")

	project.HandleFunc("/{contractAddress}/{projectID}/categories", h.updateBTCProjectcategories).Methods("PUT")
	project.HandleFunc("/{genNFTAddr}/tokens", h.TokensOfAProject).Methods("GET")

	projectAuth := api.PathPrefix("/project").Subrouter()
	projectAuth.Use(h.MiddleWare.AccessToken)
	projectAuth.HandleFunc("/{projectID}/report", h.reportProject).Methods("POST")
	projectAuth.HandleFunc("/btc", h.createBTCProject).Methods("POST")
	projectAuth.HandleFunc("/btc/files", h.UploadProjectFiles).Methods("POST")
	projectAuth.HandleFunc("/{contractAddress}/tokens/{projectID}", h.updateBTCProject).Methods("PUT")
	projectAuth.HandleFunc("/{contractAddress}/{projectID}", h.deleteBTCProject).Methods("DELETE")

	//configs
	config := api.PathPrefix("/configs").Subrouter()
	config.HandleFunc("", h.getConfigs).Methods("GET")
	config.HandleFunc("", h.createConfig).Methods("POST")
	config.HandleFunc("/{key}", h.getConfig).Methods("GET")
	config.HandleFunc("/{key}", h.deleteConfig).Methods("DELETE")

	//categories
	categories := api.PathPrefix("/categories").Subrouter()
	categories.HandleFunc("", h.getCategories).Methods("GET")
	categories.HandleFunc("", h.createCategory).Methods("POST")
	categories.HandleFunc("/{id}", h.getCategory).Methods("GET")
	categories.HandleFunc("/{id}", h.updateCategory).Methods("PUT")
	categories.HandleFunc("/{id}", h.deleteCategory).Methods("DELETE")

	//nfts
	nfts := api.PathPrefix("/nfts").Subrouter()
	nfts.HandleFunc("/{contractAddress}/transactions/{tokenID}", h.getNftTransactions).Methods("GET")
	nfts.HandleFunc("/{contractAddress}/nft_holders", h.getTokenHolder).Methods("GET")

	//admin
	admin := api.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/redis", h.getRedisKeys).Methods("GET")
	admin.HandleFunc("/redis/{key}", h.getRedis).Methods("GET")
	admin.HandleFunc("/redis", h.upsertRedis).Methods("POST")
	admin.HandleFunc("/redis", h.deleteAllRedis).Methods("DELETE")
	admin.HandleFunc("/redis/{key}", h.deleteRedis).Methods("DELETE")

	admin.Use(h.MiddleWare.AccessToken)
	admin.HandleFunc("/auto-listing", h.autoListing).Methods("POST")

	//Marketplace
	marketplace := api.PathPrefix("/marketplace").Subrouter()
	marketplace.HandleFunc("/listing/{genNFTAddr}/token/{tokenID}", h.getListingViaGenAddressTokenID).Methods("GET")
	marketplace.HandleFunc("/offers/{genNFTAddr}/token/{tokenID}", h.getOffersViaGenAddressTokenID).Methods("GET")
	marketplace.HandleFunc("/wallet/{walletAddress}/listing", h.ListingOfAProfile).Methods("GET")
	marketplace.HandleFunc("/wallet/{walletAddress}/offer", h.OfferOfAProfile).Methods("GET")
	marketplace.HandleFunc("/stats/{genNFTAddr}", h.getCollectionStats).Methods("GET")

	//dao
	dao := api.PathPrefix("/dao").Subrouter()
	dao.HandleFunc("/proposals", h.proposals).Methods("GET")
	dao.HandleFunc("/proposals", h.createDraftProposals).Methods("POST")
	dao.HandleFunc("/proposals/{proposalID}", h.getProposal).Methods("GET")
	dao.HandleFunc("/proposals/{proposalID}/votes", h.getProposalVotes).Methods("GET")
	dao.HandleFunc("/proposals/{ID}/{proposalID}", h.mapOffAndOnChainProposal).Methods("PUT")

	//btc
	btc := api.PathPrefix("/btc").Subrouter()
	btc.HandleFunc("/receive-address", h.btcGetReceiveWalletAddress).Methods("POST")
	btc.HandleFunc("/balance", h.checkBalance).Methods("POST")

	// btcV2 := api.PathPrefix("/btc-v2").Subrouter()
	// btcV2.HandleFunc("/receive-address", h.btcGetReceiveWalletAddressV2).Methods("POST")

	inscribe := api.PathPrefix("/inscribe").Subrouter()
	inscribe.HandleFunc("/receive-address", h.btcCreateInscribeBTC).Methods("POST")
	inscribe.HandleFunc("/list", h.btcListInscribeBTC).Methods("GET")
	inscribe.HandleFunc("/nft-detail/{ID}", h.btcDetailInscribeBTC).Methods("GET")
	inscribe.HandleFunc("/retry/{ID}", h.btcRetryInscribeBTC).Methods("POST")
	inscribe.HandleFunc("/info/{ID}", h.getInscribeInfo).Methods("GET")

	//btc
	eth := api.PathPrefix("/eth").Subrouter()
	eth.HandleFunc("/receive-address", h.ethGetReceiveWalletAddress).Methods("POST")
	signedEth := api.PathPrefix("/eth").Subrouter()
	signedEth.Use(h.MiddleWare.AccessToken)
	signedEth.HandleFunc("/receive-address/whitelist", h.ethGetReceiveWhitelistedWalletAddress).Methods("POST")
	btc.HandleFunc("/balance", h.checkBalance).Methods("POST")

	// request-mint (new flow)
	mintNftBtc := api.PathPrefix("/mint-nft-btc").Subrouter()
	mintNftBtc.HandleFunc("/receive-address", h.createMintReceiveAddress).Methods("POST")
	mintNftBtc.HandleFunc("/receive-address/{uuid}", h.getDetailMintNftBtc).Methods("GET")

	mintNftBtcAuth := api.PathPrefix("/mint-nft-btc").Subrouter()
	mintNftBtcAuth.Use(h.MiddleWare.AccessToken)
	mintNftBtcAuth.HandleFunc("/receive-address/{uuid}", h.cancelMintNftBt).Methods("DELETE")

	marketplaceBTC := api.PathPrefix("/marketplace-btc").Subrouter()
	marketplaceBTC.HandleFunc("/listing", h.btcMarketplaceListing).Methods("POST")
	marketplaceBTC.HandleFunc("/list", h.btcMarketplaceListNFTs).Methods("GET")
	marketplaceBTC.HandleFunc("/nft-detail/{ID}", h.btcMarketplaceNFTDetail).Methods("GET")
	marketplaceBTC.HandleFunc("/nft-gen-order", h.btcMarketplaceCreateBuyOrder).Methods("POST")
	marketplaceBTC.HandleFunc("/listing-fee", h.btcMarketplaceListingFee).Methods("POST")
	marketplaceBTC.HandleFunc("/filter-info", h.btcMarketplaceFilterInfo).Methods("GET")
	marketplaceBTC.HandleFunc("/run-filter-info", h.btcMarketplaceRunFilterInfo).Methods("GET")
	marketplaceBTC.HandleFunc("/collection-stats", h.btcMarketplaceCollectionStats).Methods("GET")

	referral := api.PathPrefix("/referrals").Subrouter()
	referral.Use(h.MiddleWare.AccessToken)
	referral.HandleFunc("/{referrerID}", h.createReferral).Methods("POST")
	referral.HandleFunc("", h.getReferrals).Methods("GET")

	// marketplaceBTC.HandleFunc("/search", h.btcMarketplaceSearch).Methods("GET") //TODO: implement

	// marketplaceBTC.HandleFunc("/test-listen", h.btcTestListen).Methods("GET")

	// marketplaceBTC.HandleFunc("/test-transfer", h.btcTestTransfer).Methods("POST")

	wallet := api.PathPrefix("/wallet").Subrouter()
	// wallet.Use(h.MiddleWare.AccessToken)
	// wallet.HandleFunc("/inscription-by-output", h.inscriptionByOutput).Methods("POST")
	wallet.HandleFunc("/wallet-info", h.walletInfo).Methods("GET")
	wallet.HandleFunc("/mint-status", h.mintStatus).Methods("GET")
	// wallet.HandleFunc("/")

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("", h.getUsers).Methods("GET")
	user.HandleFunc("/artist", h.listArtist).Methods("GET")

	tokenUri := api.PathPrefix("/token-uri").Subrouter()
	tokenUri.HandleFunc("", h.getTokenUris).Methods("GET")
}

func (h *httpDelivery) RegisterDocumentRoutes() {
	documentUrl := `/generative/swagger/`
	domain := os.Getenv("swagger_domain")
	docs.SwaggerInfo.Host = domain
	docs.SwaggerInfo.BasePath = "/generative/api"
	swaggerURL := documentUrl + "swagger/doc.json"
	h.Handler.PathPrefix(documentUrl).Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerURL), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		//httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
}

func (h *httpDelivery) healthCheck(w http.ResponseWriter, r *http.Request) {
	h.Response.RespondSuccess(w, http.StatusOK, response.Success, "It work!", "")
}

func (h *httpDelivery) PaginationResp(data *entity.Pagination, items interface{}) response.PaginationResponse {
	resp := response.PaginationResponse{}
	resp.Result = items
	resp.Currsor = data.Currsor
	resp.Total = data.Total
	resp.Page = data.Page
	resp.PageSize = data.PageSize
	return resp
}

func (h *httpDelivery) BaseFilters(r *http.Request) (*structure.BaseFilters, error) {
	f := &structure.BaseFilters{}

	limitInt := 10
	pageInt := 1
	var err error

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			return nil, err
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			return nil, err
		}
	}

	sortQuery := r.URL.Query().Get("sort")
	sortObject := utils.ParseSort(sortQuery)

	f.SortBy = sortObject.SortBy
	f.Sort = sortObject.Sort
	f.Page = int64(pageInt)
	f.Limit = int64(limitInt)

	return f, nil
}
