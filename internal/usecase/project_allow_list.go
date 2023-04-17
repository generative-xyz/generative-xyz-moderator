package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateProjectAllowList(req structure.CreateProjectAllowListReq) (*entity.ProjectAllowList, error) {
	userAddress := strings.ToLower(*req.UserWalletAddress)
	projectID := strings.ToLower(*req.ProjectID)
	allowedBy := entity.ERC721

	p, err := u.Repo.FindProjectByTokenID(projectID)
	if err != nil {
		return nil, err
	}

	user, err := u.Repo.FindUserByAddress(userAddress)
	if err != nil {
		return nil, err
	}

	isWhitelist, _ := u.ProjectWhitelistERC721(*user)
	if !isWhitelist {
		isWhitelist, _ = u.ProjectWhitelistERC20(*user)
		if !isWhitelist {
			//return nil, errors.New("User is not existed in allowlist")
			allowedBy = entity.PUBLIC
		} else {
			allowedBy = entity.ERC20
		}
	}

	_ = isWhitelist
	pe := &entity.ProjectAllowList{}
	pe.ProjectID = p.TokenID
	pe.UserWalletAddress = userAddress
	pe.UserWalletAddressBTCTaproot = user.WalletAddressBTCTaproot
	pe.AllowedBy = allowedBy
	err = u.Repo.CreateProjectAllowList(pe)
	if err != nil {
		//err := errors.New("Error while create allow list")
		logger.AtLog.Logger.Error("Error while create allow list", zap.Any("error", err))
		return pe, nil
	}

	//SLACK_ALLOW_LIST_CHANNEL
	go func(u Usecase, user entity.Users, p entity.Projects, allowBy entity.AllowedByType) {
		totalCount, _ := u.Repo.GetProjectAllowListTotal(p.TokenID)
		publicCount, _ := u.Repo.GetProjectAllowListTotalByTyppe(p.TokenID, "public")
		u.NotifyWithChannel(os.Getenv("SLACK_ALLOW_LIST_CHANNEL"),
			"",
			"",
			fmt.Sprintf("%s registered to Metaverse as %s PUBLIC: %d AL: %d",
				helpers.CreateProfileLink(user.WalletAddressBTCTaproot, user.DisplayName), allowedBy, publicCount, totalCount-publicCount))
	}(u, *user, *p, allowedBy)
	return pe, nil
}

func (u Usecase) CountingProjectAllowList(projectId string) (int64, int64, error) {
	totalCount, err := u.Repo.GetProjectAllowListTotal(projectId)
	if err != nil {
		return 0, 0, err
	}
	publicCount, err := u.Repo.GetProjectAllowListTotalByTyppe(projectId, "public")
	if err != nil {
		return 0, 0, err
	}
	return publicCount, totalCount - publicCount, nil
}

func (u Usecase) GetProjectAllowList(req structure.CreateProjectAllowListReq) (*entity.ProjectAllowList, error) {
	userAddress := strings.ToLower(*req.UserWalletAddress)
	projectID := strings.ToLower(*req.ProjectID)

	allowed, err := u.Repo.GetProjectAllowList(projectID, userAddress)
	if err != nil {
		err := errors.New("Error while create allow list")
		return nil, err
	}

	return allowed, nil
}

func (u Usecase) CheckExistedProjectAllowList(req structure.CreateProjectAllowListReq) (bool, string) {
	userAddress := strings.ToLower(*req.UserWalletAddress)
	projectID := strings.ToLower(*req.ProjectID)

	allowed, err := u.Repo.GetProjectAllowList(projectID, userAddress)
	if err != nil {
		return false, ""
	}

	if allowed == nil {
		return false, ""
	}

	return true, string(allowed.AllowedBy)
}

func (u Usecase) ProjectWhitelistERC721(user entity.Users) (bool, error) {
	if user.UUID == "" || user.WalletAddressBTCTaproot == "" {
		return false, errors.New("User is not found")
	}

	whitelist := os.Getenv("WHITELIST_PROJECT_ALLOWED_LIST")
	if len(strings.TrimSpace(whitelist)) == 0 {
		return false, errors.New("Error while get whitelist")
	}

	whitelistArreses := strings.Split(whitelist, ",")

	isWhiteList, err := u.IsWhitelistedAddress(context.Background(), user.WalletAddress, whitelistArreses)
	if err != nil {
		return false, err
	}

	return isWhiteList, nil
}

func (u Usecase) ProjectWhitelistERC20(user entity.Users) (bool, error) {
	if user.UUID == "" || user.WalletAddressBTCTaproot == "" {
		return false, errors.New("User is not found")
	}

	whitelist := os.Getenv("WHITELIST_ERC20_PROJECT_ALLOWED_LIST")
	bytes, err := helpers.Base64Decode(whitelist)
	if err != nil {
		return false, err
	}

	erc20WhiteList := make(map[string]structure.Erc20Config)
	err = json.Unmarshal(bytes, &erc20WhiteList)
	if err != nil {
		return false, err
	}

	isWhiteList, err := u.IsWhitelistedAddressERC20(context.Background(), user.WalletAddress, erc20WhiteList)
	if err != nil {
		return false, err
	}

	return isWhiteList, nil
}
