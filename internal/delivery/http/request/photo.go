package request

type CaptureRequest struct {
	ID  string `json:"device_id"`
	Url string `json:"display_url"`
}

type ParseSvgRequest struct {
	Url       string `json:"display_url"`
	ID        string `json:"device_id"`
	DelayTime int    `json:"delay_time"`
}
