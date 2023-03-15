package structure

type FilterCollectionListing struct {
	BaseFilters
	CollectionContract *string
	TokenId            *string
	Erc20Token         *string
	SellerAddress      *string
	Closed             *bool
	Finished           *bool
}
