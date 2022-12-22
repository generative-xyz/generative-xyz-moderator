package structure

type MinifyDataResp struct {
	Files map[string]FileContentReq `json:"files"`
}

type FileContentReq struct {
	MediaType string `json:"mediaType"`
	Content string `json:"content"`
}
