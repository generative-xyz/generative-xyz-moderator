package artblock

type GetArtists struct {
	Data ArtistArray `json:"data"`
}

type ArtistArray struct {
	Artists []Artist `json:"artists"`
}

type Artist struct {
	PublicAddress string `json:"public_address"`
}
