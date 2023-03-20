package usecase

import (
	"context"
	"log"
	"sync"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
)

func (u Usecase) SubmitBTCTransaction(list map[string]string) error {

	submiTxs := []entity.BTCTransactionSubmit{}
	for txhash, raw := range list {
		err := btc.SendTxBlockStream(raw)
		if err != nil {
			submiTxs = append(submiTxs, entity.BTCTransactionSubmit{
				Txhash: txhash,
				Raw:    raw,
				Status: 0,
			})
			continue
		}
		submiTxs = append(submiTxs, entity.BTCTransactionSubmit{
			Txhash: txhash,
			Raw:    raw,
			Status: 1,
		})
	}

	dbObjs := make([]interface{}, 0, len(submiTxs))
	for _, objectId := range submiTxs {
		objectId.SetID()
		objectId.SetCreatedAt()
		dbObjs = append(dbObjs, objectId)
	}

	_, err := u.Repo.CreateMany(context.Background(), utils.COLLECTION_BTC_TX_SUBMIT, dbObjs)
	if err != nil {
		return err
	}
	return nil
}

func (u Usecase) watchPendingBTCTxSubmit() error {
	pendingTxs, err := u.Repo.GetPendingBTCSubmitTx()
	if err != nil {
		return err
	}

	for _, tx := range pendingTxs {
		txDetail, err := btc.CheckTxfromQuickNode(tx.Txhash, u.Config.QuicknodeAPI)
		if err != nil {
			log.Printf("watchPendingBTCTxSubmit CheckTxfromQuickNode %v\n", err)
		} else {
			if txDetail.Result.Confirmations >= 0 {
				tx.Status = 1
				_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
				if err != nil {
					continue
				}
				continue
			}
		}

		err = btc.SendTxBlockStream(tx.Raw)
		if err != nil {
			continue
		}
		tx.Status = 1
		_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
		if err != nil {
			continue
		}
	}

	return nil
}

func (u Usecase) JobWatchPendingBTCTxSubmit() {
	var wg sync.WaitGroup

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err := u.watchPendingBTCTxSubmit()
		if err != nil {
			log.Println("JobWatchPendingBTCTxSubmit watchPendingBTCTxSubmit err", err)
		}
	}(&wg)

	wg.Wait()
}
