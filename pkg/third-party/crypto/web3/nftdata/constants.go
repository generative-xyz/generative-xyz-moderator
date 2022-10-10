package nftdata

const (
	CompanyNameDefault   = "NFT Manifest"
	CompanyDomainDefault = "autonomousnft.manifest"
)

var (
	DefaultItem = &Item{
		Domain:      CompanyDomainDefault,
		CompanyName: CompanyNameDefault,
	}
)
