package request

type CreateFcmRequest struct {
	RegistrationToken string `json:"registration_token"  validate:"required"`
	DeviceType        string `json:"device_type"  validate:"required"`
	UserWallet        string `json:"-"`
}

type CreateFcmDataTest struct {
	Data              map[string]string `json:"data" validate:"required"`
	DeviceType        string            `json:"device_type" validate:"required"`
	RegistrationToken string            `json:"registration_token"`
}
