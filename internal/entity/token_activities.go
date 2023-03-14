package entity

import "time"

type TokenActivityType int
const (
	TokenMint TokenActivityType = 0
	TokenListing TokenActivityType = 1
	TokenCancelListing TokenActivityType = 2
	TokenMatched TokenActivityType = 3
)

type TokenActivity struct {
	Type TokenActivityType
	Title string
	UserAAddress string
	UserA *Users
	UserBAddress string
	UserB *Users
	Amount int64
	Time  *time.Time
}
