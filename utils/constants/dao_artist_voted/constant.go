package dao_artist_voted

type Status int64

const (
	Report Status = iota
	Verify
)
