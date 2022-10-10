package nftdata

type Item struct {
	Domain      string
	CompanyName string
}

type Client interface {
	GetFirstNFTItem() *Item
	GetData() interface{}
}
