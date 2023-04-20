package request

type CaptureRequest struct {
	ID  string `json:"device_id"`
	Url string `json:"display_url"`
}
