package ord_service

type ExecRequest struct {
	Args []string `json:"args"`
}

type MintRequest struct {
	FileUrl            string `json:"fileUrl"`
	WalletName         string `json:"walletName"`
	FeeRate            int    `json:"feeRate"`
	DryRun             bool   `json:"dryRun"`
	AutoFeeRateSelect  bool   `json:"autoFeeRateSelect"`
	DestinationAddress string `json:"destinationAddress"`
	RequestId          string `json:"requestId"`
	ProjectID          string `json:"projectID"`
	FileUrlUnit        string `json:"fileUrlUnit"`
}

type ExecRespose struct {
	Message   string `json:"message"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ErrorCode string `json:"errorCode"`
	Error     string `json:"err"`
}

type MintRespose struct {
	Message   string `json:"message"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ErrorCode string `json:"errorCode"`
	Error     string `json:"err"`
}

type MintStdoputRespose struct {
	Commit      string `json:"commit"`
	Inscription string `json:"inscription"`
	Reveal      string `json:"reveal"`
	Fees        int    `json:"fees"`
	IsSent      bool   `json:"isSent"`
}

type ReceiveCmdStdoputRespose struct {
	Address string `json:"address"`
}
