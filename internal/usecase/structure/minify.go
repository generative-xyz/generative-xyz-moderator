package structure

type MinifyDataResp struct {
	Files map[string]FileContentReq `json:"files"`
}

type DeflateDataResp struct {
	Data string  `json:"data"`
}

type FileContentReq struct {
	MediaType string `json:"mediaType"`
	Content   string `json:"content"`
	Deflate   string `json:"deflate"`
}
