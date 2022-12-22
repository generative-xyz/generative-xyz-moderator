package response

type FileRes struct{
	BaseResponse
	FileName string `json:"fileName"`
	UploadedBy string `json:"uploadedBy"`
	URL string `json:"url"`
	MineType     string       `json:"mimeType"`
	FileSize     int       `json:"fileSize"`
}
