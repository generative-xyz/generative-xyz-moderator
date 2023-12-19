package request

type CaptureRequest struct {
	ID  string `json:"device_id"`
	Url string `json:"display_url"`
}

type ParseSvgRequest struct {
	Url            string `json:"display_url"`
	HtmlContent    string `json:"html_content"`
	CaptureElement string `json:"capture_element"`
	ID             string `json:"app_id"`
	DelayTime      int    `json:"delay_time"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
}
