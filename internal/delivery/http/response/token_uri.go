package response

type TokenURIResp struct{
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
	AnimationURL string `json:"animation_url"`
	Attributes interface{} `json:"attributes"`
}
