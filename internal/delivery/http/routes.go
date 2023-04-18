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

	tokenActivities := api.PathPrefix("/token-activities").Subrouter()
	tokenActivities.HandleFunc("", h.getTokenActivities).Methods("GET")

	files := api.PathPrefix("/files").Subrouter()
	files.HandleFunc("", h.UploadFile).Methods("POST")
	files.HandleFunc("/minify", h.minifyFiles).Methods("POST")
	files.HandleFunc("/deflate", h.deflate).Methods("POST")
	files.HandleFunc("/image/resize", h.resizeImage).Methods("POST")

	files.HandleFunc("/multipart", h.CreateMultipartUpload).Methods("POST")
	files.HandleFunc("/multipart/{uploadID}", h.UploadPart).Methods("PUT")
	files.HandleFunc("/multipart/{uploadID}", h.CompleteMultipartUpload).Methods("POST")

	//profile
	profile := api.PathPrefix("/profile").Subrouter()
	profile.Use(h.MiddleWare.AuthorizationFunc)
	profile.HandleFunc("/wallet/{walletAddress}", h.profileByWallet).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/nfts", h.TokensOfAProfile).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/projects", h.getProjectsByWallet).Methods("GET")
	profile.HandleFunc("/wallet/{walletAddress}/volumn", h.getVolumnByWallet).Methods("GET")

	singedIn := api.PathPrefix("/profile").Subrouter()
	singedIn.Use(h.MiddleWare.AccessToken)
	singedIn.HandleFunc("/withdraw", h.withdraw).Methods("POST")
	singedIn.HandleFunc("", h.profile).Methods("GET")
	singedIn.HandleFunc("/projects", h.getUserProjects).Methods("GET")
	singedIn.HandleFunc("", h.updateProfile).Methods("PUT")
	singedIn.HandleFunc("/logout", h.logout).Methods("PUT")

	//project
	project := api.PathPrefix("/project").Subrouter()
	project.Use(h.MiddleWare.AuthorizationFunc)
	project.HandleFunc("", h.getProjects).Methods("GET")
	project.HandleFunc("/all", h.getAllProjects).Methods("GET")
	project.HandleFunc("", h.createEthProjects).Methods("POST")

	project.HandleFunc("/random", h.getRandomProject).Methods("GET")
	project.HandleFunc("/upcomming", h.getUpcommingProjects).Methods("GET")
	project.HandleFunc("/minted-out", h.getMintedOutProjects).Methods("GET")
	project.HandleFunc("/recent-works", h.getRecentWorksProjects).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}", h.projectDetail).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}/marketplace-data", h.projectMarketplaceData).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}/volumn", h.projectVolumn).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}/random-images", h.projectRandomImages).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}/token-traits", h.tokenTraits).Methods("GET")
	project.HandleFunc("/{contractAddress}/tokens/{projectID}/token-traits", h.uploadTokenTraits).Methods("POST")
	project.HandleFunc("/{contractAddress}/{projectID}", h.updateProject).Methods("PUT")
	project.HandleFunc("/{contractAddress}/tx-hash/{txHash}", h.updateProjectHash).Methods("PUT")
	project.HandleFunc("/{contractAddress}/{projectID}/allow-list", h.createProjectAllowList).Methods("POST")
	project.HandleFunc("/{contractAddress}/{projectID}/allow-list", h.getProjectAllowList).Methods("GET")
	project.HandleFunc("/{contractAddress}/{projectID}/counting-allow-list", h.getCountingAllowList).Methods("GET")
	project.HandleFunc("/{contractAddress}/{projectID}/token", h.searchToken).Methods("GET")

	project.HandleFunc("/{contractAddress}/{projectID}/categories", h.updateBTCProjectcategories).Methods("PUT")
	// project.HandleFunc("/{genNFTAddr}/tokens", h.TokensOfAProject).Methods("GET")
	project.HandleFunc("/{genNFTAddr}/tokens", h.TokensOfAProjectNew).Methods("GET")

	projectAuth := api.PathPrefix("/project").Subrouter()
	projectAuth.Use(h.MiddleWare.AccessToken)
	projectAuth.HandleFunc("/{projectID}/report", h.reportProject).Methods("POST")
	projectAuth.HandleFunc("/btc", h.createBTCProject).Methods("POST")
	projectAuth.HandleFunc("/btc/files", h.UploadProjectFiles).Methods("POST")
	projectAuth.HandleFunc("/{contractAddress}/tokens/{projectID}", h.updateBTCProject).Methods("PUT")
	projectAuth.HandleFunc("/{contractAddress}/{projectID}", h.deleteBTCProject).Methods("DELETE")
	projectAuth.HandleFunc("/{contractAddress}/tokens/{projectID}/token-thumbnail", h.triggerPubsubTokenThumbnail).Methods("GET")

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
	admin.HandleFunc("/check-refund", h.checkRefundMintBtc).Methods("POST")
	admin.HandleFunc("/gen-temp-address", h.getMintFreeTemAddress).Methods("POST")

	admin.HandleFunc("/auction/update-declared-now", h.updateDeclaredNow).Methods("POST")          // auction
	admin.HandleFunc("/auction/crawl-list-winner-now", h.updateWinnerFromContract).Methods("POST") // auction
	admin.HandleFunc("/evm/update-enabled-job", h.updateEnabledJob).Methods("POST")                // auction

	adminTest := api.PathPrefix("/admin-test").Subrouter()
	adminTest.HandleFunc("", h.adminTest).Methods("GET")

	//Marketplace
	marketplace := api.PathPrefix("/marketplace").Subrouter()
	marketplace.HandleFunc("/listing/{genNFTAddr}/token/{tokenID}", h.getListingViaGenAddressTokenID).Methods("GET")
	marketplace.HandleFunc("/offers/{genNFTAddr}/token/{tokenID}", h.getOffersViaGenAddressTokenID).Methods("GET")
	marketplace.HandleFunc("/wallet/{walletAddress}/listing", h.ListingOfAProfile).Methods("GET")
	marketplace.HandleFunc("/wallet/{walletAddress}/offer", h.OfferOfAProfile).Methods("GET")
	marketplace.HandleFunc("/stats/{genNFTAddr}", h.getCollectionStats).Methods("GET")
	marketplace.HandleFunc("/stats/{genNFTAddr}/first-sale", h.getCollectionStatsFirstSale).Methods("GET")

	// New Marketplace
	collection := api.PathPrefix("/collections").Subrouter()
	collection.HandleFunc("", h.getCollectionListing).Methods("GET")
	collection.HandleFunc("/items", h.getItemListing).Methods("GET")
	collection.HandleFunc("/sub-collection-items", h.getSubCollectionItemListing).Methods("GET")
	collection.HandleFunc("/on-sale-items", h.getItemListingOnSale).Methods("GET")

	charts := api.PathPrefix("/charts").Subrouter()
	charts.HandleFunc("/collections/{projectID}", h.getChartDataForCollection).Methods("GET")
	charts.HandleFunc("/tokens/{tokenID}", h.getChartDataFoTokenURI).Methods("GET")

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
	// inscribe.Use(h.MiddleWare.AccessToken)
	inscribe.HandleFunc("/receive-address", h.btcCreateInscribeBTC).Methods("POST")
	inscribe.HandleFunc("/list", h.btcListInscribeBTC).Methods("GET")
	inscribe.HandleFunc("/nft-detail/{ID}", h.btcDetailInscribeBTC).Methods("GET")
	inscribe.HandleFunc("/retry/{ID}", h.btcRetryInscribeBTC).Methods("POST")
	inscribe.HandleFunc("/info/{ID}", h.getInscribeInfo).Methods("GET")
	inscribe.HandleFunc("/compress-image", h.compressImage).Methods("POST")

	inscribeAuth := inscribe.PathPrefix("/auth").Subrouter()
	inscribeAuth.Use(h.MiddleWare.AccessToken)
	inscribeAuth.HandleFunc("/receive-address", h.btcCreateInscribeBTC).Methods("POST")

	tokenMoralis := api.PathPrefix("/token-moralis").Subrouter()
	tokenMoralis.Use(h.MiddleWare.AccessToken)
	tokenMoralis.HandleFunc("/nfts", h.listNftFromMoralis).Methods("GET")
	tokenMoralis.HandleFunc("/nfts/{tokenAddress}", h.nftFromMoralis).Methods("GET")

	//btc
	eth := api.PathPrefix("/eth").Subrouter()
	eth.HandleFunc("/receive-address", h.ethGetReceiveWalletAddress).Methods("POST")
	signedEth := api.PathPrefix("/eth").Subrouter()
	signedEth.Use(h.MiddleWare.AccessToken)
	signedEth.HandleFunc("/receive-address/whitelist", h.ethGetReceiveWhitelistedWalletAddress).Methods("POST")
	btc.HandleFunc("/balance", h.checkBalance).Methods("POST")

	// request-mint (new flow)
	mintNftBtcAuth := api.PathPrefix("/mint-nft-btc").Subrouter()
	mintNftBtcAuth.Use(h.MiddleWare.AccessToken)
	mintNftBtcAuth.HandleFunc("/receive-address", h.createMintReceiveAddress).Methods("POST")
	mintNftBtcAuth.HandleFunc("/receive-address/{uuid}", h.getDetailMintNftBtc).Methods("GET")
	mintNftBtcAuth.HandleFunc("/receive-address/{uuid}", h.cancelMintNftBt).Methods("DELETE")

	mintNftBtc := api.PathPrefix("/mint-nft-btc").Subrouter()
	mintNftBtc.HandleFunc("/get-mint-fee-rate-info/{fileSize}/{customRate}/{mintPrice}", h.getMintFeeRateInfos).Methods("GET")

	// faucet:
	faucet := api.PathPrefix("/faucet").Subrouter()
	faucet.HandleFunc("/request", h.requestFaucet).Methods("POST")
	faucet.HandleFunc("/list", h.listFaucet).Methods("GET")

	marketplaceBTC := api.PathPrefix("/marketplace-btc").Subrouter()
	// marketplaceBTC.HandleFunc("/listing", h.btcMarketplaceListing).Methods("POST")
	// marketplaceBTC.HandleFunc("/list", h.btcMarketplaceListNFTs).Methods("GET")
	// marketplaceBTC.HandleFunc("/nft-detail/{ID}", h.btcMarketplaceNFTDetail).Methods("GET")
	// marketplaceBTC.HandleFunc("/nft-gen-order", h.btcMarketplaceCreateBuyOrder).Methods("POST")
	// marketplaceBTC.HandleFunc("/listing-fee", h.btcMarketplaceListingFee).Methods("POST")
	// marketplaceBTC.HandleFunc("/filter-info", h.btcMarketplaceFilterInfo).Methods("GET")
	// marketplaceBTC.HandleFunc("/run-filter-info", h.btcMarketplaceRunFilterInfo).Methods("GET")
	marketplaceBTC.HandleFunc("/collection-stats", h.btcMarketplaceCollectionStats).Methods("GET")

	referral := api.PathPrefix("/referrals").Subrouter()
	referral.HandleFunc("", h.getReferrals).Methods("GET")

	referralAuth := api.PathPrefix("/referrals").Subrouter()
	referralAuth.Use(h.MiddleWare.AccessToken)
	referralAuth.HandleFunc("/{referrerID}", h.createReferral).Methods("POST")

	// marketplaceBTC.HandleFunc("/search", h.btcMarketplaceSearch).Methods("GET") //TODO: implement

	marketplaceBTC.HandleFunc("/test-listen", h.btcTestListen).Methods("GET")

	// marketplaceBTC.HandleFunc("/test-transfer", h.btcTestTransfer).Methods("POST")

	wallet := api.PathPrefix("/wallet").Subrouter()
	// wallet.HandleFunc("/inscription-by-output", h.inscriptionByOutput).Methods("POST")
	wallet.Use(h.MiddleWare.AccessTokenPassThrough)
	wallet.HandleFunc("/wallet-info", h.walletInfo).Methods("GET")
	wallet.HandleFunc("/mint-status", h.mintStatus).Methods("GET")
	wallet.HandleFunc("/track-tx", h.trackTx).Methods("POST")
	wallet.HandleFunc("/txs", h.walletTrackedTx).Methods("GET")
	wallet.HandleFunc("/submit-tx", h.submitTx).Methods("POST")

	tc := api.PathPrefix("/tc").Subrouter()
	tc.HandleFunc("/nfts/{address}", h.getNftsByAddress).Methods("GET")
	tc.HandleFunc("/nfts-owner/{address}", h.getNftsByAddressFromTokenUri).Methods("GET")

	inscriptionDex := api.PathPrefix("/dex").Subrouter()
	inscriptionDex.Use(h.MiddleWare.AccessTokenPassThrough)
	inscriptionDex.HandleFunc("/listing", h.dexBTCListing).Methods("POST")
	inscriptionDex.HandleFunc("/listing-fee", h.dexBTCListingFee).Methods("POST")
	inscriptionDex.HandleFunc("/cancel", h.cancelBTCListing).Methods("POST")
	inscriptionDex.HandleFunc("/retrieve-order", h.retrieveBTCListingOrderInfo).Methods("GET")
	inscriptionDex.HandleFunc("/retrieve-orders", h.retrieveBTCListingOrdersInfo).Methods("POST")
	inscriptionDex.HandleFunc("/history", h.historyBTCListing).Methods("GET")
	// inscriptionDex.HandleFunc("/submit-buy", h.submitDexBTCBuy).Methods("GET")
	//buy with eth
	inscriptionDexAuth := api.PathPrefix("/dex").Subrouter()
	inscriptionDexAuth.Use(h.MiddleWare.AccessToken)
	inscriptionDexAuth.HandleFunc("/gen-eth-order", h.genDexBTCBuyETHOrder).Methods("POST")
	inscriptionDexAuth.HandleFunc("/buy-eth-history", h.dexBTCBuyETHHistory).Methods("GET")
	// inscriptionDex.HandleFunc("/update-eth-order-tx", h.updateDexBTCBuyETHOrderTx).Methods("POST")
	// inscriptionDex.HandleFunc("/submit-buy-eth", h.submitDexBTCBuyETHTx).Methods("POST")

	// dex := api.PathPrefix("/dex-buy").Subrouter()
	// dex.HandleFunc("/list-buy-address", h.ListBuyAddress).Methods("GET")

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("", h.getUsers).Methods("GET")
	user.HandleFunc("/artist", h.listArtist).Methods("GET")

	tokenUri := api.PathPrefix("/token-uri").Subrouter()
	tokenUri.HandleFunc("", h.getTokenUris).Methods("GET")

	search := api.PathPrefix("/search").Subrouter()
	search.HandleFunc("", h.search).Methods("GET")

	// for dev:
	developerAuth := api.PathPrefix("/developer").Subrouter()
	developerAuth.Use(h.MiddleWare.AccessToken)
	developerAuth.HandleFunc("/api-key", h.apiDeveloper_GenApiKey).Methods("POST")
	developerAuth.HandleFunc("/api-key", h.apiDeveloper_GetApiKey).Methods("GET")

	// public but required api key:
	developer := api.PathPrefix("/developer").Subrouter()
	developer.HandleFunc("/inscribe", h.developerCreateInscribe).Methods("POST")
	developer.HandleFunc("/inscribe", h.developerInscribeList).Methods("GET")
	developer.HandleFunc("/inscribe/{ID}", h.developerDetailInscribe).Methods("GET")
	// inscribe.HandleFunc("/retry/{ID}", h.btcRetryInscribeBTC).Methods("POST")
	// inscribe.HandleFunc("/info/{ID}", h.getInscribeInfo).Methods("GET")

	// Firebase FCM registration token management
	fcm := api.PathPrefix("/fcm").Subrouter()
	fcm.Use(h.MiddleWare.AuthorizationFunc)
	fcm.HandleFunc("/token", h.getFcmToken).Methods("GET")
	fcm.HandleFunc("/token", h.createFcmToken).Methods("POST")
	// For test, will remove
	fcm.HandleFunc("/token/data", h.createFcmTestData).Methods("POST")

	// DAO Project
	daoProject := api.PathPrefix("/dao-project").Subrouter()
	daoProject.Use(h.MiddleWare.AuthorizationFunc)
	daoProject.HandleFunc("", h.listDaoProject).Methods("GET")
	daoProject.HandleFunc("", h.createDaoProject).Methods("POST")
	daoProject.HandleFunc("/{id}", h.getDaoProject).Methods("GET")
	daoProject.HandleFunc("/{id}", h.voteDaoProject).Methods("PUT")
	daoProject.HandleFunc("/me/projects-hidden", h.listYourProjectsIsHidden).Methods("GET")

	// DAO Artist
	daoArtist := api.PathPrefix("/dao-artist").Subrouter()
	daoArtist.Use(h.MiddleWare.AuthorizationFunc)
	daoArtist.HandleFunc("", h.listDaoArtist).Methods("GET")
	daoArtist.HandleFunc("", h.createDaoArtist).Methods("POST")
	daoArtist.HandleFunc("/{id}", h.getDaoArtist).Methods("GET")
	daoArtist.HandleFunc("/{id}", h.voteDaoArtist).Methods("PUT")

	//ordinal collections
	orCollections := api.PathPrefix("/ordinal").Subrouter()
	//orCollections.Use(h.MiddleWare.AuthorizeFunc)
	orCollections.HandleFunc("/collections/template", h.getOrdinalTemplate).Methods("GET")
	orCollections.HandleFunc("/collections", h.uploadOrdinalTemplate).Methods("POST")

	auction := api.PathPrefix("/auction").Subrouter()
	auction.HandleFunc("/list", h.getListAuction).Methods("GET")
	auction.HandleFunc("/list-snapshot", h.listSnapshot).Methods("GET")
	auction.HandleFunc("/check-declared", h.checkDeclared).Methods("GET")

	auction.HandleFunc("/list-winner-config", h.listBidWinnerConfig).Methods("GET")

	auction.HandleFunc("/list-winner-bid-list", h.listBidWinnerBigList).Methods("GET")

	auctionAuth := api.PathPrefix("/auction").Subrouter()
	auctionAuth.Use(h.MiddleWare.AccessToken)
	auctionAuth.HandleFunc("/shard-now", h.shareNow).Methods("POST")

	discord := api.PathPrefix("/discord").Subrouter()
	discord.Use(h.MiddleWare.AuthorizationFunc)
	discord.HandleFunc("/new-bid", h.sendDiscordNewBid).Methods("POST")

	aiSchool := api.PathPrefix("/ai-school").Subrouter()
	aiSchool.HandleFunc("/check-progress", h.schoolCheckProgress).Methods("GET")
	aiSchool.HandleFunc("/upload", h.schoolUpload).Methods("POST")
	aiSchool.HandleFunc("/download", h.schoolDownload).Methods("GET")

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
	resp.Currsor = data.Cursor
	resp.Total = data.Total
	resp.Page = data.Page
	resp.PageSize = data.PageSize
	resp.TotalPage = data.TotalPage
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
	sortObject := utils.ParseSortNew(sortQuery)

	f.SortBy = sortObject.SortBy
	f.Sort = sortObject.Sort
	f.Page = int64(pageInt)
	f.Limit = int64(limitInt)

	return f, nil
}

func (h *httpDelivery) BaseAlgoliaFilters(r *http.Request) (*structure.BaseFilters, error) {
	f := &structure.BaseFilters{}
	limitInt := 10
	pageInt := 0
	var err error
	limit := r.URL.Query().Get("limit")
	if limit != "" {
		if limitInt, err = strconv.Atoi(limit); err != nil {
			return nil, err
		}
	}

	page := r.URL.Query().Get("page")
	if page != "" {
		if pageInt, err = strconv.Atoi(page); err != nil {
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
