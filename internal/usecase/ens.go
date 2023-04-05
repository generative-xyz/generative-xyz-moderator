package usecase

import (
	"math/big"
	"os"

	"rederinghub.io/internal/entity"
)

func (u Usecase) JobAuction_GetListAuction() error {

	contractV1 := os.Getenv("AUCTION_CONTRACT")
	contractV2 := os.Getenv("AUCTION_CONTRACT_v2")

	contractV1 = "0xd178cC5A4001fDecc6cBc0293C5Bed5e3887351D"
	contractV2 = "0xdfd1EdD11748E14620EaE6F45093d5A25434e07a"

	listMap1, _ := u.EthClient.GetListBidV1(contractV1)
	listMap2, _ := u.EthClient.GetListBidV1(contractV2)

	for address, v := range listMap1 {

		_, ok := listMap2[address]
		if !ok {
			// get item from db:
			item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)
			if item == nil {
				// insert:
				u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
					Bidder:    v.Bidder,
					IsWinner:  v.IsWinner,
					Amount:    v.Amount.String(),
					Quantity:  v.Quantity,
					UnitPrice: big.NewInt(0).SetUint64(v.UnitPrice).String(),
					Ens:       v.Ens,
				})
			} else {
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount.String()
				item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
				item.Ens = v.Ens
				// update:
				u.Repo.UpdateAuctionCollectionBidder(item)
			}
		}

	}
	for address, v := range listMap2 {

		v1, ok := listMap1[address]
		if !ok {
			// get item from db:
			item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)
			if item == nil {
				// insert:
				u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
					Bidder:    v.Bidder,
					IsWinner:  v.IsWinner,
					Amount:    v.Amount.String(),
					Quantity:  v.Quantity,
					UnitPrice: big.NewInt(0).SetUint64(v.UnitPrice).String(),
					Ens:       v.Ens,
				})
			} else {
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount.String()
				item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
				item.Ens = v.Ens
				// update:
				u.Repo.UpdateAuctionCollectionBidder(item)
			}
		} else {
			// merge:
			v.Amount = big.NewInt(0).Add(v.Amount, v1.Amount)
			// get item from db:
			item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)
			if item == nil {
				// insert:
				u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
					Bidder:    v.Bidder,
					IsWinner:  v.IsWinner,
					Amount:    v.Amount.String(),
					Quantity:  v.Quantity,
					UnitPrice: big.NewInt(0).SetUint64(v.UnitPrice).String(),
					Ens:       v.Ens,
				})
			} else {
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount.String()
				item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
				item.Ens = v.Ens
				// update:
				u.Repo.UpdateAuctionCollectionBidder(item)
			}

		}

	}

	return nil
}

func (u Usecase) APIGetListAuction() ([]entity.AuctionCollectionBidder, error) {
	return u.Repo.ListAuctionCollectionBidder()
}
