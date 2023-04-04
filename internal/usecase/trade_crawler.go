package usecase

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) crawlTokenTxFrom(tokenTx entity.TokenTx) ([]entity.TokenTx, error) {
	logger.AtLog.Logger.Info(
		"crawlTokenTxFrom.Start", 
		zap.String("inscriptionID", tokenTx.InscriptionID), 
		zap.String("startTx", tokenTx.Tx),
	)
	if tokenTx.NextTx != "" {
		return nil, errors.New("crawlTokenTxFrom.TailTxHasNotEmptyNextTx")
	}

	gotError := false

	txs := []string{}
	currentTx := tokenTx.Tx
	for {
		logger.AtLog.Logger.Info(
			"crawlTokenTxFrom.StartGetTransactionInfo",
			zap.String("inscriptionID", tokenTx.InscriptionID),
			zap.String("startTx", tokenTx.Tx),
		)
		resp, err := http.Get("https://api.blockchain.info/haskoin-store/btc/transaction/" + currentTx)
		if err != nil {
			logger.AtLog.Logger.Error(
				"crawlTokenTxFrom.ErrorWhenGetTransaction", 
				zap.Error(err),
				zap.String("currentTx", currentTx),
			)
			gotError = true
			break
		}
		tx := structure.Tx{}
		err = json.NewDecoder(resp.Body).Decode(&tx)
		if len(tx.Outputs) == 0 {
			break
		}
		if tx.Outputs[0].Spender.TxId == "" {
			break
		}
		currentTx = tx.Outputs[0].Spender.TxId
		txs = append(txs, currentTx)
		logger.AtLog.Logger.Info(
			"crawlTokenTxFrom.FindNewTx",
			zap.String("inscriptionID", tokenTx.InscriptionID), 
			zap.String("startTx", tokenTx.Tx),
			zap.String("newTx", currentTx),
		)
	}

	depth := tokenTx.Depth
	prevTx := tokenTx.Tx

	lastTimeCheck := tokenTx.LastTimeCheck
	if !gotError {
		x := time.Now()
		lastTimeCheck = &x 
	}

	tokenTx.LastTimeCheck = lastTimeCheck

	tokenTxs := []entity.TokenTx{
		tokenTx,
	}

	for _, txId := range txs {
		depth++
		newTokenTx := entity.TokenTx{
			InscriptionID: tokenTx.InscriptionID,
			Tx: txId,
			PrevTx: prevTx,
			Depth: depth,
			LastTimeCheck: lastTimeCheck,
			Priority: tokenTx.Priority,
		}
		tokenTxs = append(tokenTxs, newTokenTx)
	}
	
	for i := 0; i + 1 < len(tokenTxs); i++ {
		tokenTxs[i].NextTx = tokenTxs[i + 1].Tx
	}

	if gotError {
		tokenTxs[len(tokenTxs) - 1].NumFailed += 1
	}

	return tokenTxs, nil
}

func (u Usecase) fetchDataFromTx(tokenTx entity.TokenTx) error {
	u.Repo.AddTokenTxRetryResolve(tokenTx.InscriptionID, tokenTx.Tx)
	tx := structure.Tx{}
	txId := tokenTx.Tx
	logger.AtLog.Logger.Info(
		"fetchDataFromTx.Start", 
		zap.String("tx", txId),
	)
	resp, err := http.Get("https://api.blockchain.info/haskoin-store/btc/transaction/" + txId)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.NewDecoder(resp.Body).Decode(&tx)
	if err != nil {
		return errors.WithStack(err)
	}
	tradingTx := false
	for _, input := range tx.Inputs {
		if len(input.Witness) > 0 {
			temp, _ := hex.DecodeString(input.Witness[0])
			if len(temp) > 0 &&  temp[len(temp)-1] == 131 {
				tradingTx = true
				break
			}
		}
	}
	if tradingTx && len(tx.Outputs) >= 3 { 
		logger.AtLog.Logger.Info("fetchDataFromTx.MeetTradingTx", zap.String("tx", txId))
		buyer := tx.Outputs[0].Address
		seller := tx.Outputs[1].Address
		amount := tx.Outputs[1].Value
		txTime := time.Unix(tx.Time, 0)
		// Create listing
		// check existed
		existed, err := u.Repo.CheckMatchedTxExisted(tokenTx.Tx)
		if err != nil {
			return errors.WithStack(err)
		}
		
		if existed {
			logger.AtLog.Logger.Info("fetchDataFromTx.ListingExisted", zap.String("tx", tokenTx.Tx))
			u.Repo.UpdateResolvedTx(tokenTx.InscriptionID, tokenTx.Tx)
			return nil
		}

		newDexBTCListing := entity.DexBTCListing{
			InscriptionID: tokenTx.InscriptionID,
			Amount: amount,
			Matched: true,
			MatchedTx: tokenTx.Tx,
			MatchAt: &txTime,
			Verified: true,
			Cancelled: false,
			FromOtherMkp: true,
			SellerAddress: seller,
			Buyer: buyer,
		}
		err = u.Repo.CreateDexBTCListing(&newDexBTCListing)
		if err != nil {
			return errors.WithStack(err)
		}

		listing, err := u.Repo.GetDexBTCListingByMatchedTx(tokenTx.Tx)
		if err != nil {
			return errors.WithStack(err)
		}

		u.InsertDexVolumnInscription(*listing)
	}

	_, err = u.Repo.UpdateResolvedTx(tokenTx.InscriptionID, tokenTx.Tx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (u Usecase) GoFromTailTokenTx(tail entity.TokenTx) error {
	logger.AtLog.Logger.Info(
		"GoFromTailTokenTx.Start",
		zap.String("inscriptionID", tail.InscriptionID),
		zap.String("tx", tail.Tx),
	)
	tokenTxs, err := u.crawlTokenTxFrom(tail)
	if err != nil {
		logger.AtLog.Logger.Error("GoFromTailTokenTx.crawlTokenTxFrom", zap.Error(err))
		return errors.WithStack(err)
	}

	for i := len(tokenTxs) - 1; i >= 0; i-- {
		tokenTx := tokenTxs[i]
		logger.AtLog.Logger.Info("GoFromTailTokenTx.UpsertTokenTx", zap.Any("tokenTx", zap.Any("tokenTx)", tokenTx)))
		_, err = u.Repo.UpsertTokenTx(tokenTx.InscriptionID, tokenTx.Tx, &tokenTx)
		if err != nil {
			logger.AtLog.Logger.Error("GoFromTailTokenTx.UpsertTokenTxFailed", zap.Error(err), zap.Any("tokenTx", tokenTx))
			return errors.WithStack(err)
		}
	}
	return nil
}

func (u Usecase) JobCreateTokenTxFromTokenURI() error {
	logger.AtLog.Logger.Info("JobCreateTokenTxFromTokenURI.Start")
	startTime := time.Time{}
	for page := int64(1);; page++ {
		logger.AtLog.Logger.Info("JobCreateTokenTxFromTokenURI.GetPagingTokenUri", zap.Any("page", zap.Any("page)", page)))
		uTokens, err := u.Repo.GetNotCreatedTxToken(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokens := uTokens.Result.([]entity.TokenUri)
		if len(tokens) == 0 {
			break
		}
		logger.AtLog.Logger.Info(
			"JobCreateTokenTxFromTokenURI.DonePagingTokenUri", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokens)),
		)
		for _, token := range tokens {
			trendingScore, err := u.Repo.GetProjectTrendingScore(token.ProjectID) 
			if err != nil {
				logger.AtLog.Logger.Error(
					"JobCreateTokenTxFromTokenURI.ErrorGetProjectTrendingScore", 
					zap.Error(err), 
					zap.String("inscriptionID", token.TokenID),
				)
				continue
			}
			tokenTx := entity.TokenTx{
				InscriptionID: token.TokenID,
				Tx: token.TokenID[:len(token.TokenID) - 2],
				LastTimeCheck: &startTime,
				Priority: trendingScore,
			}
			
			if err := u.Repo.InsertTokenTx(&tokenTx); err != nil {
				logger.AtLog.Logger.Error(
					"JobCreateTokenTxFromTokenURI.InsertTokenTx", 
					zap.Error(err), 
					zap.String("token_id", token.TokenID),
				)
			} else {
				logger.AtLog.Logger.Info("JobCreateTokenTxFromTokenURI.InsertTokenTx", zap.Any("tokenTx", zap.Any("tokenTx)", tokenTx)))
				u.Repo.UpdateTokenCreatedTokenTx(token.TokenID)
			}
		}
	}
	return nil
}

func (u Usecase) JobContinueCrawlTxs() error {
	logger.AtLog.Logger.Info("JobContinueCrawlTxs.Start")
	var processed int64
	for page := int64(1);; page++ {
		logger.AtLog.Logger.Info("JobContinueCrawlTxs.GetPagingTokenTx", zap.Any("page", zap.Any("page)", page)))
		uTokenTxs, err := u.Repo.GetTailTokenTxs(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokenTxs := uTokenTxs.Result.([]entity.TokenTx)
		if len(tokenTxs) == 0 {
			break
		}
		logger.AtLog.Logger.Info(
			"JobContinueCrawlTxs.DonePagingTokenTx", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokenTxs)),
		)
		for _, tokenTx := range tokenTxs {
			if err := u.GoFromTailTokenTx(tokenTx); err != nil {
				logger.AtLog.Logger.Error(
					"JobContinueCrawlTxs.GoFromTailTokenTx",
					zap.Error(err),
					zap.String("inscriptionID", tokenTx. InscriptionID),
					zap.String("tx", tokenTx.Tx),
				)
			}
			processed++
			if processed % 5 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return nil
}

func (u Usecase) JobFetchUnresolvedTokenTxs() error {
	logger.AtLog.Logger.Info("JobFetchUnresolvedTokenTxs.Start")
	var processed int64
	for page := int64(1);; page++ {
		logger.AtLog.Logger.Info("JobFetchUnresolvedTokenTxs.GetPagingTokenTx", zap.Any("page", zap.Any("page)", page)))
		uTokenTxs, err := u.Repo.GetUnresolvedTokenTx(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokenTxs := uTokenTxs.Result.([]entity.TokenTx)
		if len(tokenTxs) == 0 {
			break
		}
		logger.AtLog.Logger.Info(
			"JobFetchUnresolvedTokenTxs.DonePagingTokenTx", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokenTxs)),
		)
		for _, tokenTx := range tokenTxs {
			if err := u.fetchDataFromTx(tokenTx); err != nil {
				logger.AtLog.Logger.Error(
					"JobFetchUnresolvedTokenTxs.fetchDataFromTx",
					zap.Error(err),
					zap.String("inscriptionID", tokenTx.InscriptionID),
					zap.String("tx", tokenTx.Tx),
				)
			}
			processed++
			if processed % 5 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}
	return nil
}

const LAST_CRAWL_INSCRIPTION_INDEX_KEY string = "last_crawl_inscription_index"

func (u Usecase) JobCrawlTokenTxNotFromTokenUri() error {
	startTime := time.Time{}
	var limit int64 = 50
	for {
		pLastCrawled, err := u.Repo.GetVariableInt(LAST_CRAWL_INSCRIPTION_INDEX_KEY)
		if err != nil {
			return errors.WithStack(err)
		}
		if pLastCrawled == nil {
			return errors.New("ErrGetVariableInt")
		}
		lastCrawled := *pLastCrawled
		fr := lastCrawled + 1
		to := fr + limit - 1
		logger.AtLog.Logger.Info(
			"JobCrawlTokenTxNotFromTokenUri.CrawlFromSearcher",
			zap.Int64("fr", fr),
			zap.Int64("to", to),
		)
		url := fmt.Sprintf("https://generative.xyz/generative/api/search?limit=%v&page=1&search=&type=inscription&fromNumber=%v&toNumber=%v", limit, fr, to)
		resp, err := http.Get(url)
		if err != nil {
			logger.AtLog.Logger.Error(
				"JobCrawlTokenTxNotFromTokenUri.ErrorCrawlFromSearcher", 
				zap.Error(err),
				zap.Int64("fr", fr),
				zap.Int64("to", to),
			)
			break
		}
		tx := structure.SearchInscriptionResult{}
		err = json.NewDecoder(resp.Body).Decode(&tx)
		if len(tx.Data.Result) == 0 {
			break
		}
		for _, inscription := range tx.Data.Result {
			if inscription.Inscription.ProjectTokenID != "" {
				continue
			}
			inscriptionID := inscription.Inscription.InscriptionID
			tokenTx := entity.TokenTx{
				InscriptionID: inscriptionID,
				Tx: inscriptionID[:len(inscriptionID) - 2],
				LastTimeCheck: &startTime,
				Priority: 0,
				Source: "NOT_GEN_TOKEN",
			}
			
			if err := u.Repo.InsertTokenTx(&tokenTx); err != nil {
				logger.AtLog.Logger.Error(
					"JobCrawlTokenTxNotFromTokenUri.InsertTokenTx", 
					zap.Error(err), 
					zap.String("token_id", inscriptionID),
				)
			} else {
				logger.AtLog.Logger.Info("JobCrawlTokenTxNotFromTokenUri.InsertTokenTx", zap.Any("tokenTx", zap.Any("tokenTx)", tokenTx)))
			}
		}
		minOf := func (vars ...int) int {
			min := vars[0]
		
			for _, i := range vars {
				if min > i {
					min = i
				}
			}
		
			return min
		}
		mxInscriptionIndex := 0
		for _, inscription := range tx.Data.Result {
			mxInscriptionIndex = minOf(inscription.Inscription.Number, mxInscriptionIndex)
		}
		u.Repo.UpdateGlobalActivity(LAST_CRAWL_INSCRIPTION_INDEX_KEY, mxInscriptionIndex)
		time.Sleep(time.Second)
	}
	return nil
}
