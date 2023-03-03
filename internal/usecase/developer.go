package usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
)

func (u Usecase) ApiDevelop_GenApiKey(userAddr string, req *request.GetApiKeyReq) (*entity.DeveloperKey, error) {

	// check admin user:
	profile, err := u.GetUserProfileByWalletAddress(userAddr)
	if err != nil {
		u.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
		return nil, err
	}

	// check exists:
	apiKeyObj, _ := u.Repo.FindIDeveloperKeyByUserID(profile.UUID)

	if apiKeyObj != nil {
		return apiKeyObj, nil
	}

	fmt.Println("profile: ", profile)

	if req.RecaptchaResponse == "" {
		return nil, errors.New("the recaptcha is required.")
	}

	if len(u.Config.CaptcharSecret) > 0 {
		captcha, _ := recaptcha.NewReCAPTCHA(u.Config.CaptcharSecret, recaptcha.V3, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin

		err = captcha.Verify(req.RecaptchaResponse)
		if err != nil {
			return nil, err
		}
	}
	key := make([]byte, 16)
	_, err = rand.Read(key)
	if err != nil {
		return nil, err
	}
	apiKey := strings.ToUpper(fmt.Sprintf("%x", key))
	fmt.Println("apiKey: ", apiKey)

	// insert:
	apiKeyObj = &entity.DeveloperKey{
		ApiKey:         apiKey,
		ApiDescription: req.Description,
		ApiName:        req.Name,
		ApiEmail:       req.Email,
		ApiCompany:     req.Company,
		UserUuid:       profile.UUID,
		Status:         1,
	}

	err = u.Repo.InsertDeveloperKey(apiKeyObj)

	if err != nil {
		return nil, err
	}

	// gen key now:
	return apiKeyObj, nil
}

func (u Usecase) ApiDevelop_GetApiKey(userAddr string) (*entity.DeveloperKey, error) {

	// check admin user:
	profile, err := u.GetUserProfileByWalletAddress(userAddr)
	if err != nil {
		u.Logger.Error("h.Usecase.GetUserProfileByWalletAddress(", err.Error(), err)
		return nil, err
	}

	// check exists:
	apiKeyObj, _ := u.Repo.FindIDeveloperKeyByUserID(profile.UUID)

	if apiKeyObj != nil {
		return apiKeyObj, nil
	}

	return nil, errors.New("api key not found")
}
