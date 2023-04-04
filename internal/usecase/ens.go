package usecase

import (
	"os"

	"rederinghub.io/internal/entity"
)

func (u Usecase) JobAuction_GetListAuction() error {
	listMap, _ := u.EthClient.GetListDomainName(os.Getenv("AUCTION_CONTRACT")) // 0xB57BF9669186dCe8FbCe7E5EbeE41f210bb6a7Eb
	if listMap != nil {
		for _, v := range listMap {
			// get item from db:
			item, _ := u.Repo.FindAuctionCollectionBidderByAddress(v.Bidder)
			if item == nil {
				// insert:
				u.Repo.InsertAuctionCollectionBidder(&entity.AuctionCollectionBidder{
					Bidder:   v.Bidder,
					IsWinner: v.IsWinner,
					Amount:   v.Amount,
					Ens:      v.Ens,
				})
			} else {
				item.IsWinner = v.IsWinner
				item.Amount = v.Amount
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
