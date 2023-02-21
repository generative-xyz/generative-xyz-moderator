package request

type UpdateTokentReq struct {
	Priority *int `json:"priority"`
}


type UpdateTokenThumbnailReq struct {
	Thumbnail *string `json:"thumbnail"`
}
