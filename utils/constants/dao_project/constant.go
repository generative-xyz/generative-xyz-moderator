package dao_project

type Status int64

const (
	Voting Status = iota
	Executed
	Defeated
)
