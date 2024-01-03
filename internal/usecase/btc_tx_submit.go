package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils"
	"rederinghub.io/utils/btc"
)

func (u Usecase) SubmitBTCTransaction(list map[string]string) error {
	submiTxs := []entity.BTCTransactionSubmit{}
	for txhash, raw := range list {
		relatedInscription := make(map[string]struct{})
		relatedInscriptionArray := []string{}
		txDetail, err := btc.ParseTx(raw)
		if err != nil {
			return fmt.Errorf("can't parse transaction")
		}
		outPointList := []string{}
		for _, txIn := range txDetail.TxIn {
			outPoint := fmt.Sprintf("%v:%v", txIn.PreviousOutPoint.Hash.String(), txIn.PreviousOutPoint.Index)
			outPointList = append(outPointList, outPoint)
		}

		listingOrder, err := u.Repo.GetDexBTCListingOrderPendingByInputs(outPointList)
		if err != nil {
			return fmt.Errorf("GetDexBTCListingOrderPendingByInputs err %v", err)
		}

		for _, listing := range listingOrder {
			relatedInscription[listing.InscriptionID] = struct{}{}
		}

		for v, _ := range relatedInscription {
			relatedInscriptionArray = append(relatedInscriptionArray, v)
		}

		err = btc.SendTxBlockStream(raw)
		if err != nil {
			submiTxs = append(submiTxs, entity.BTCTransactionSubmit{
				Txhash:              txhash,
				Raw:                 raw,
				RelatedInscriptions: relatedInscriptionArray,
				Status:              0,
				Error1:              err.Error(),
			})
			continue
		}
		submiTxs = append(submiTxs, entity.BTCTransactionSubmit{
			Txhash:              txhash,
			Raw:                 raw,
			RelatedInscriptions: relatedInscriptionArray,
			Status:              1,
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
		log.Printf("SubmitBTCTransaction CreateMany err %v\n", err)
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
		switch tx.Status {
		case entity.StatusBTCTransactionSubmit_Waiting:
			txDetail, err := btc.CheckTxfromQuickNode(tx.Txhash, u.Config.QuicknodeAPI)
			if err != nil {
				log.Printf("watchPendingBTCTxSubmit CheckTxfromQuickNode %v\n", err)
			} else {
				if txDetail.Result.Confirmations >= 0 {
					tx.Status = 1
					_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
					if err != nil {
						log.Printf("watchPendingBTCTxSubmit UpdateBTCTxSubmitStatus err %v\n", err)
						continue
					}
					continue
				}
			}

			err = btc.SendTxBlockStream(tx.Raw)
			if err != nil {
				log.Printf("watchPendingBTCTxSubmit SendTxBlockStream err %v\n", err)
				tx.Status = entity.StatusBTCTransactionSubmit_Failed
				_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
				if err != nil {
					log.Printf("watchPendingBTCTxSubmit UpdateBTCTxSubmitStatus err %v\n", err)
					continue
				}
				continue
			}
			tx.Status = entity.StatusBTCTransactionSubmit_Pending
			_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
			if err != nil {
				log.Printf("watchPendingBTCTxSubmit UpdateBTCTxSubmitStatus err %v\n", err)
				continue
			}
			continue
		case entity.StatusBTCTransactionSubmit_Pending:
			txDetail, err := btc.CheckTxfromQuickNode(tx.Txhash, u.Config.QuicknodeAPI)
			if err != nil {
				log.Printf("watchPendingBTCTxSubmit CheckTxfromQuickNode %v\n", err)
				tx.Status = entity.StatusBTCTransactionSubmit_Failed
				_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
				if err != nil {
					log.Printf("watchPendingBTCTxSubmit UpdateBTCTxSubmitStatus err %v\n", err)
					continue
				}
				continue
			} else {
				if txDetail.Result.Confirmations >= 1 {
					tx.Status = entity.StatusBTCTransactionSubmit_Success
					_, err = u.Repo.UpdateBTCTxSubmitStatus(&tx)
					if err != nil {
						log.Printf("watchPendingBTCTxSubmit UpdateBTCTxSubmitStatus err %v\n", err)
						continue
					}
					continue
				}
			}
		}
	}

	return nil
}

func (u Usecase) JobWatchPendingBTCTxSubmit() error {
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
	return nil
}
