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

	// testnet
	if u.Config.ENV == "develop" {
		contractV1 = "0x367504f3d304c39154acafb769ad25d861fb78fb"
		contractV2 = "0x3b724a99c9d427d0793b63088a39c19735208900"
	}

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
					Contract:  "v1",
				})
				fmt.Println("inserted address1 db, err:", err)
			} else {
				fmt.Println("address1 exit db, update now")
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount.String()
				item.Quantity = v.Quantity
				item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
				item.Ens = v.Ens
				item.Contract = "v1"
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

		contract := ""

		if !ok {
			contract = "v2"
		} else {
			// merge:
			fmt.Println("merge v2 vs v1 now ")
			fmt.Println("v2 amount: ", v.Amount)
			fmt.Println("v2-v1 amount: ", v1.Amount)

			fmt.Println("v2 Quantity: ", v.Quantity)
			fmt.Println("v2-v1 Quantity: ", v1.Quantity)
			v.Amount = big.NewInt(0).Add(v.Amount, v1.Amount)
			contract = "both"
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
				Contract:  contract,
			})
			fmt.Println("inserted address2 db, err:", err)
		} else {
			fmt.Println("address2 exit db, update now")
			item.IsWinner = v.IsWinner
			item.Amount = v.Amount.String()
			item.Quantity = v.Quantity
			item.UnitPrice = big.NewInt(0).SetUint64(v.UnitPrice).String()
			item.Ens = v.Ens
			item.Contract = contract
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

func (u Usecase) APIAuctionCheckDeclared() bool {
	key := "auction-declared"
	config, _ := u.Repo.FindConfig(key)
	if config != nil {
		return config.Data.(bool)
	} else {
		config = &entity.Configs{
			Key:  key,
			Data: false,
		}
		u.Repo.InsertConfig(config)
	}
	return false

}

func (u Usecase) APIAuctionDeclaredNow() error {
	key := "auction-declared"
	var err error
	config, _ := u.Repo.FindConfig(key)
	if config != nil {
		config.Data = true
		_, err = u.Repo.UpdateConfig(config.UUID, config)
	} else {
		config = &entity.Configs{
			Key:  key,
			Data: true,
		}
		err = u.Repo.InsertConfig(config)
		fmt.Println("err InsertConfig declared: ", err)
	}

	if err == nil {
		keySnapShot := "auction-list-snapshot"
		configListSnapshot, _ := u.Repo.FindConfig(keySnapShot)
		if configListSnapshot == nil {
			configListSnapshot = &entity.Configs{
				Key: keySnapShot,
			}
			err = u.Repo.InsertConfig(configListSnapshot)
			fmt.Println("err InsertConfig keySnapShot: ", err)
		}
		if configListSnapshot != nil {
			listAuctionBid, err := u.Repo.ListAuctionCollectionBidderShort()
			fmt.Println("err ListAuctionCollectionBidder: ", err)
			if listAuctionBid != nil {

				// listAuctionBidJson, err := json.Marshal(listAuctionBid)
				// if err == nil {
				configListSnapshot.Data = listAuctionBid
				// update:
				_, err = u.Repo.UpdateConfig(configListSnapshot.UUID, configListSnapshot)
				return err
				// }

			}

		}

	}

	return err
}

func (u Usecase) APIAuctionListSnapshot() interface{} {

	var result struct {
		Key   string                                `bson:"key"`
		Value string                                `bson:"value"`
		Data  []entity.AuctionCollectionBidderShort `bson:"data"`
	}

	// var result:
	err := u.Repo.FindConfigCustom("auction-list-snapshot", &result)
	if err != nil {
		return err
	}
	return result.Data

}

func (u Usecase) GetAuctionListWinnerAddress() ([]entity.AuctionWinnerList, error) {

	var auctionWinnerList []entity.AuctionWinnerList

	var result struct {
		Key  string                     `bson:"key"`
		Data []entity.AuctionWinnerList `bson:"data"`
	}

	// var result:
	err := u.Repo.FindConfigCustom("auction-list-winner-address", &result)
	if err != nil {
		return auctionWinnerList, err
	}
	return result.Data, nil

}
