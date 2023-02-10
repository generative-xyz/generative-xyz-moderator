package usecase

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/utils/btc"
)

func (u Usecase) buildBTCClient() (*rpcclient.Client, *btc.BlockcypherService, error) {
	host := u.Config.BTC_FULLNODE
	user := u.Config.BTC_RPCUSER
	pass := u.Config.BTC_RPCPASSWORD

	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true,  // Bitcoin core only supports HTTP POST mode
		DisableTLS:   false, //!(os.Getenv("BTC_NODE_HTTPS") == "true"), // Bitcoin core does not provide TLS by default
	}

	bs := btc.NewBlockcypherService(u.Config.BlockcypherAPI, "", u.Config.BlockcypherToken, &chaincfg.MainNetParams)

	return rpcclient.New(connCfg, nil), bs, nil
}

// check nft of the nft:
func (u Usecase) BtcChecktListNft(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcChecktListNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingListings()
	if len(listPending) == 0 {
		return nil
	}

	for _, item := range listPending {

		txs, _ := u.getLastTxs(item.HoldOrdAddress)

		if len(txs) == 0 {
			continue
		}
		found := false
		for _, tx := range txs {
			detail, err := chainhash.NewHashFromStr(tx.Tx)
			if err != nil {
				fmt.Println("can not NewHashFromStr with err:", err)
				continue
			}
			result, _ := btcClient.GetRawTransactionVerboseAsync(detail).Receive()

			for _, vin := range result.Vin {
				if strings.Contains(vin.Txid, item.InscriptionID) {
					found = true
					item.IsConfirm = true
					_, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
					if err != nil {
						fmt.Println("UpdateBTCNFTConfirmListings", err.Error(), err)
					}
					break
				}
			}
			if found {
				break
			}
		}
	}

	return nil
}

// check nft of the nft:
func (u Usecase) BtcCheckBuyingNft(rootSpan opentracing.Span) error {

	bs := btc.NewBlockcypherService()

	span, log := u.StartSpan("BtcCheckBuyingNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingBuyOrders()
	if len(listPending) == 0 {
		return nil
	}

	for _, item := range listPending {

		// check balance:

		txs, _ := u.getLastTxs(item.HoldOrdAddress)

		if len(txs) == 0 {
			continue
		}
		found := false
		for _, tx := range txs {
			detail, err := chainhash.NewHashFromStr(tx.Tx)
			if err != nil {
				fmt.Println("can not NewHashFromStr with err:", err)
				continue
			}
			result, _ := btcClient.GetRawTransactionVerboseAsync(detail).Receive()

			for _, vin := range result.Vin {
				if strings.Contains(vin.Txid, item.InscriptionID) {
					found = true
					item.IsConfirm = true
					_, err := u.Repo.UpdateBTCNFTConfirmListings(&item)
					if err != nil {
						fmt.Println("UpdateBTCNFTConfirmListings", err.Error(), err)
					}
					break
				}
			}
			if found {
				break
			}
		}
	}

	return nil
}
