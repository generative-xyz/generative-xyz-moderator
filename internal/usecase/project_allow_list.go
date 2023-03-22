package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/helpers"
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
	if ! isWhitelist  {
		isWhitelist, _ = u.ProjectWhitelistERC20(*user)
		if ! isWhitelist {
			return nil, errors.New("User is not existed in whitelist")
		}
		allowedBy = entity.ERC20
	}

	_ = isWhitelist
	pe := &entity.ProjectAllowList{}
	pe.ProjectID = p.TokenID
	pe.UserWalletAddress = userAddress
	pe.AllowedBy = allowedBy
	err = u.Repo.CreateProjectAllowList(pe)
	if err != nil {
		err := fmt.Errorf("Error while create allow list: %v", err.Error())
		return nil, err
	}

	return pe, nil
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

	erc20WhiteList := make(map[string]uint64)
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