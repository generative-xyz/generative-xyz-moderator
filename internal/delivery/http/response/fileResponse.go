package response

type FileResponse struct {
	UploadID string `json:"uploadId"`
}

type MultipartUploadResponse struct {
	FileURL string `json:"fileUrl"`
}
