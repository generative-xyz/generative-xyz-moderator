package usecase

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) crawlTokenTxFrom(tokenTx entity.TokenTx) ([]entity.TokenTx, error) {
	u.Logger.LogAny(
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
		u.Logger.LogAny(
			"crawlTokenTxFrom.StartGetTransactionInfo",
			zap.String("inscriptionID", tokenTx.InscriptionID),
			zap.String("startTx", tokenTx.Tx),
		)
		resp, err := http.Get("https://api.blockchain.info/haskoin-store/btc/transaction/" + currentTx)
		if err != nil {
			u.Logger.ErrorAny(
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
		u.Logger.LogAny(
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
	tx := structure.Tx{}
	txId := tokenTx.Tx
	u.Logger.LogAny(
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
		temp, _ := hex.DecodeString(input.Witness[0])
		if temp[len(temp)-1] == 131 {
			tradingTx = true
			break
		}
	}
	if tradingTx { 
		u.Logger.LogAny("fetchDataFromTx.MeetTradingTx", zap.String("tx", txId))
		if len(tx.Outputs) < 3 {
			return errors.New("trading tx must havee at least 3 items")
		} 
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
			u.Logger.LogAny("fetchDataFromTx.ListingExisted", zap.String("tx", tokenTx.Tx))
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

		_, err = u.Repo.UpdateResolvedTx(tokenTx.InscriptionID, tokenTx.Tx)
		if err != nil {
			return errors.WithStack(err)
		}

		u.InsertDexVolumnInscription(*listing)
	}
	return nil
}

func (u Usecase) GoFromTailTokenTx(tail entity.TokenTx) error {
	u.Logger.LogAny(
		"GoFromTailTokenTx.Start",
		zap.String("inscriptionID", tail.InscriptionID),
		zap.String("tx", tail.Tx),
	)
	tokenTxs, err := u.crawlTokenTxFrom(tail)
	if err != nil {
		u.Logger.ErrorAny("GoFromTailTokenTx.crawlTokenTxFrom", zap.Error(err))
		return errors.WithStack(err)
	}

	for i := len(tokenTxs) - 1; i >= 0; i-- {
		tokenTx := tokenTxs[i]
		u.Logger.LogAny("GoFromTailTokenTx.UpsertTokenTx", zap.Any("tokenTx", tokenTx))
		_, err = u.Repo.UpsertTokenTx(tokenTx.InscriptionID, tokenTx.Tx, &tokenTx)
		if err != nil {
			u.Logger.ErrorAny("GoFromTailTokenTx.UpsertTokenTxFailed", zap.Error(err), zap.Any("tokenTx", tokenTx))
			return errors.WithStack(err)
		}
	}
	return nil
}

func (u Usecase) JobCreateTokenTxFromTokenURI() error {
	u.Logger.LogAny("JobCreateTokenTxFromTokenURI.Start")
	startTime := time.Time{}
	for page := int64(1);; page++ {
		u.Logger.LogAny("JobCreateTokenTxFromTokenURI.GetPagingTokenUri", zap.Any("page", page))
		uTokens, err := u.Repo.GetNotCreatedTxToken(page, 1)
		if err != nil {
			return errors.WithStack(err)
		}
		tokens := uTokens.Result.([]entity.TokenUri)
		if len(tokens) == 0 {
			break
		}
		u.Logger.Info(
			"JobCreateTokenTxFromTokenURI.DonePagingTokenUri", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokens)),
		)
		for _, token := range tokens {
			trendingScore, err := u.Repo.GetProjectTrendingScore(token.ProjectID) 
			if err != nil {
				u.Logger.ErrorAny(
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
			
			if _, err := u.Repo.UpsertTokenTx(tokenTx.InscriptionID, tokenTx.Tx, &tokenTx); err != nil {
				u.Logger.ErrorAny(
					"JobCreateTokenTxFromTokenURI.UpsertTokenTx", 
					zap.Error(err), 
					zap.String("token_id", token.TokenID),
				)
			} else {
				u.Logger.Info("JobCreateTokenTxFromTokenURI.UpsertTokenTx", zap.Any("tokenTx", tokenTx))
				u.Repo.UpdateTokenCreatedTokenTx(token.TokenID)
			}
		}
	}
	return nil
}

func (u Usecase) JobContinueCrawlTxs() error {
	u.Logger.LogAny("JobContinueCrawlTxs.Start")
	var processed int64
	for page := int64(1);; page++ {
		u.Logger.LogAny("JobContinueCrawlTxs.GetPagingTokenTx", zap.Any("page", page))
		uTokenTxs, err := u.Repo.GetTailTokenTxs(page, 100)
		if err != nil {
			return errors.WithStack(err)
		}
		tokenTxs := uTokenTxs.Result.([]entity.TokenTx)
		if len(tokenTxs) == 0 {
			break
		}
		u.Logger.Info(
			"JobContinueCrawlTxs.DonePagingTokenTx", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokenTxs)),
		)
		for _, tokenTx := range tokenTxs {
			if err := u.GoFromTailTokenTx(tokenTx); err != nil {
				u.Logger.ErrorAny(
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
	u.Logger.LogAny("JobFetchUnresolvedTokenTxs.Start")
	var processed int64
	for page := int64(1); page < 2; page++ {
		u.Logger.LogAny("JobFetchUnresolvedTokenTxs.GetPagingTokenTx", zap.Any("page", page))
		uTokenTxs, err := u.Repo.GetUnresolvedTokenTx(page, 1)
		if err != nil {
			return errors.WithStack(err)
		}
		tokenTxs := uTokenTxs.Result.([]entity.TokenTx)
		if len(tokenTxs) == 0 {
			break
		}
		u.Logger.Info(
			"JobFetchUnresolvedTokenTxs.DonePagingTokenTx", 
			zap.Any("page", page), 
			zap.Any("numItem", len(tokenTxs)),
		)
		for _, tokenTx := range tokenTxs {
			tokenTx.Tx = "1d0b0a3560a92ea0ed075bac90fd2359cab6d0743e44c1b7640e3c4b9a40cc13"
			if err := u.fetchDataFromTx(tokenTx); err != nil {
				u.Logger.ErrorAny(
					"JobFetchUnresolvedTokenTxs.fetchDataFromTx",
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
