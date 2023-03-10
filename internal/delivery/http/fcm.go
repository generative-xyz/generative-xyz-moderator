package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/delivery/http/response"
	"rederinghub.io/utils"
	"rederinghub.io/utils/firebase"
)

// @Summary FCM Token
// @Description FCM Token
// @Tags FCM
// @Accept json
// @Produce json
// @Param device_type query string true "Device Type"
// @Success 200 {object} entity.FirebaseRegistrationToken{}
// @Router /fcm/token [GET]
// @Security ApiKeyAuth
func (h *httpDelivery) getFcmToken(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			deviceType := muxVars["device_type"]
			if deviceType == "" {
				return nil, errors.New("device type can not empty")
			}
			if userWallet != "" {
				fcm, err := h.Usecase.GetFcmByUserWalletAndDeviceType(ctx, userWallet, deviceType)
				if err != nil {
					return nil, err
				}
				return fcm, nil
			}
			return nil, nil
		},
	).ServeHTTP(w, r)
}

// @Summary Create FCM Token
// @Description Create FCM Token
// @Tags FCM
// @Accept json
// @Produce json
// @Param request body request.CreateFcmRequest true "Create fcm request"
// @Success 200
// @Router /fcm/token [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createFcmToken(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			reqBody := &request.CreateFcmRequest{}
			err := json.NewDecoder(r.Body).Decode(reqBody)
			if err != nil {
				return nil, err
			}
			if err := h.Validator.Struct(reqBody); err != nil {
				return nil, err
			}
			reqBody.UserWallet = muxVars[utils.SIGNED_WALLET_ADDRESS]
			return h.Usecase.CreateFcm(ctx, reqBody), nil
		},
	).ServeHTTP(w, r)
}

// @Summary Create FCM Token Test Data
// @Description Create FCM Token Test Data
// @Tags FCM
// @Accept json
// @Produce json
// @Param request body request.CreateFcmDataTest true "Create fcm test data request"
// @Success 200
// @Router /fcm/token/data [POST]
// @Security ApiKeyAuth
func (h *httpDelivery) createFcmTestData(w http.ResponseWriter, r *http.Request) {
	response.NewRESTHandlerTemplate(
		func(ctx context.Context, r *http.Request, muxVars map[string]string) (interface{}, error) {
			userWallet := muxVars[utils.SIGNED_WALLET_ADDRESS]
			reqBody := &request.CreateFcmDataTest{}
			err := json.NewDecoder(r.Body).Decode(reqBody)
			if err != nil {
				return nil, err
			}
			app, err := firebase.NewService(h.Config.Gcs.Auth)
			if err != nil {
				return nil, err
			}
			if reqBody.RegistrationToken == "" {
				if userWallet == "" {
					return nil, errors.New("registration token is empty")
				}
				fcm, err := h.Usecase.GetFcmByUserWalletAndDeviceType(ctx, userWallet, reqBody.DeviceType)
				if err != nil {
					return nil, err
				}
				reqBody.RegistrationToken = fcm.RegistrationToken
			}
			err = app.SendMessagesToSpecificDevices(ctx, reqBody.RegistrationToken, reqBody.Data)
			return nil, err
		},
	).ServeHTTP(w, r)
}
