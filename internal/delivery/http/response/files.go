package response

type FileRes struct{
	BaseResponse
	FileName string `json:"file_name"`
	UploadedBy string `json:"uploaded_by"`
	URL string `json:"url"`
	MineType     string       `json:"mime_type"`
	FileSize     int       `json:"file_size"`
}
