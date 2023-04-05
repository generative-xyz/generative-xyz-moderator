package usecase

import (
	"fmt"
	"math/big"
	"os"

	"rederinghub.io/internal/entity"
)

func (u Usecase) JobAuction_GetListAuction() error {

	contractV1 := os.Getenv("AUCTION_CONTRACT")
	contractV2 := os.Getenv("AUCTION_CONTRACT_v2")

	// contractV1 = "0x678B7313E34350Ec233Df5Ee0F25EFEa5C88B29f"
	// contractV2 = "0x90047bc21b0cf748507551fa1a29a40e912ce088"

	listMap1, _ := u.EthClient.GetListBidV1(contractV1)
	listMap2, _ := u.EthClient.GetListBidV2(contractV2)

	for address, v := range listMap1 {

		fmt.Println("address v1", address)

		_, ok := listMap2[address]

		fmt.Println("check exits in v2: ", ok)

		if !ok {
			// get item from db:
			item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)
			if item == nil {
				fmt.Println("address1 not exit db")
				// insert:
				err := u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
					Bidder:    v.Bidder,
					IsWinner:  v.IsWinner,
					Amount:    v.Amount.String(),
					Quantity:  v.Quantity,
					UnitPrice: big.NewInt(0).SetUint64(v.UnitPrice).String(),
					Ens:       v.Ens,
				})
				fmt.Println("inserted address1 db, err:", err)
			} else {
				fmt.Println("address1 exit db, update now")
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount.String()
				item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
				item.Ens = v.Ens
				// update:
				_, err := u.Repo.UpdateAuctionCollectionBidder(item)
				fmt.Println("update address1 db, err:", err)
			}
		}

	}
	for address, v := range listMap2 {

		fmt.Println("address v2", address)

		v1, ok := listMap1[address]

		fmt.Println("check exits in v1: ", ok)

		if !ok {

		} else {
			// merge:
			fmt.Println("merge v2 vs v1 now ")
			v.Amount = big.NewInt(0).Add(v.Amount, v1.Amount)
		}
		// get item from db:
		item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)

		if item == nil {
			fmt.Println("address2 not exit db")
			// insert:
			err := u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
				Bidder:    v.Bidder,
				IsWinner:  v.IsWinner,
				Amount:    v.Amount.String(),
				Quantity:  v.Quantity,
				UnitPrice: big.NewInt(0).SetUint64(v.UnitPrice).String(),
				Ens:       v.Ens,
			})
			fmt.Println("inserted address2 db, err:", err)
		} else {
			fmt.Println("address2 exit db, update now")
			item.IsWinner = v.IsWinner
			item.Amount = v.Amount.String()
			item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
			item.Ens = v.Ens
			// update:
			_, err := u.Repo.UpdateAuctionCollectionBidder(item)
			fmt.Println("update address2 db, err:", err)
		}

	}

	return nil
}

func (u Usecase) APIGetListAuction() ([]entity.AuctionCollectionBidder, error) {
	return u.Repo.ListAuctionCollectionBidder()
}
