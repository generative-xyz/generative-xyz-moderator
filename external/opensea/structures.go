package opensea

type User struct {
	Username string  `json:"username"`
	Account  Account `json:"account"`
}

type Account struct {
	ProfileImgUrl string `json:"profile_img_url"`
}
