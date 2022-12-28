package response

type ConfigResp struct {
	BaseResponse
	Key string `json:"key"`
	Value string `json:"value"`
}


