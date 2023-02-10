package usecase

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
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

	rpcclient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, nil, err
	}

	bs := btc.NewBlockcypherService(u.Config.BlockcypherAPI, "", u.Config.BlockcypherToken, &chaincfg.MainNetParams)

	return rpcclient, bs, nil
}

// check nft of the nft:
func (u Usecase) BtcChecktListNft(rootSpan opentracing.Span) error {

	span, log := u.StartSpan("BtcChecktListNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	btcClient, bs, err := u.buildBTCClient()

	if err != nil {
		fmt.Printf("Could not initialize Bitcoin RPCClient - with err: %v", err)
		return err
	}

	listPending, _ := u.Repo.RetrieveBTCNFTPendingListings()
	if len(listPending) == 0 {
		return nil
	}

	for _, item := range listPending {

		txs, _ := bs.GetLastTxs(item.HoldOrdAddress)

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

	span, log := u.StartSpan("BtcCheckBuyingNft", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, bs, err := u.buildBTCClient()

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

		balance, err := bs.GetBalance(item.SegwitAddress)

		if err != nil {
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}
		if balance == nil {
			err = errors.New("balance is nil")
			fmt.Printf("Could not GetBalance Bitcoin - with err: %v", err)
			continue
		}

		if balance.Uint64() == 0 {
			continue
		}

		// get amount nft:
		nftListing, err := u.Repo.FindBtcNFTListingByNFTID(item.InscriptionID)
		if err != nil {
			fmt.Printf("Could not FindBtcNFTListingByNFTID nftID: %s - with err: %v", item.InscriptionID, err)
			continue
		}
		if nftListing == nil {
			fmt.Printf("Could not FindBtcNFTListingByNFTID nftID: %s - item nil", item.InscriptionID)
			continue
		}

		amount, _ := big.NewInt(0).SetString(nftListing.Price, 10)

		if r := balance.Cmp(amount); r == -1 {
			err := errors.New("Not enough amount")
			return err

		}

		item.Status = entity.StatusBuy_ReceivedFund
		_, err = u.Repo.UpdateBTCNFTBuyOrder(&item)
		if err != nil {
			fmt.Printf("Could not UpdateBTCNFTBuyOrder id %s - with err: %v", item.ID, err)
			continue
		}

	}

	return nil
}
