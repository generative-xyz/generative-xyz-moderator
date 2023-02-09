package ord_service

type ExecRequest struct {
	Args []string `json:"args"`
}

type ExecRespose struct {
	Message string `json:"message"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	ErrorCode string `json:"errorCode"`
	Error string `json:"err"`
}
