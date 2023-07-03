package usecase

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"

	// "github.com/btcsuite/btcd/btcec/v2"
	// "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
	"rederinghub.io/utils/oauth2service"
	"rederinghub.io/utils/rediskey"
)

func (u Usecase) GenerateMessage(data structure.GenerateMessage) (*string, error) {
	logger.AtLog.Info("GenerateMessage", zap.Any("data", data))

	addrr := data.Address
	addrr = strings.ToLower(addrr)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}
	message := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	message = fmt.Sprintf(utils.NONCE_MESSAGE_FORMAT, message)
	logger.AtLog.Info("GenerateMessage", zap.Any("message", message))

	now := time.Now().UTC()
	user, err := u.Repo.FindUserByWalletAddress(addrr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//insert
			user := &entity.Users{}
			user.WalletType = data.WalletType
			user.WalletAddress = addrr
			user.Message = message
			user.CreatedAt = &now

			logger.AtLog.Info("inserted.User", user)
			err = u.Repo.CreateUser(user)
			if err != nil {
				logger.AtLog.Error(err)
				return nil, err
			}

			return &message, nil

		} else {
			logger.AtLog.Error(err)
			return nil, err
		}
	}

	logger.AtLog.Info("GenerateMessage", zap.Any("user", user))
	user.Message = message
	user.UpdatedAt = &now
	user.IsVerified = false
	_, err = u.Repo.UpdateUserByWalletAddress(addrr, user)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (u Usecase) VerifyMessage(data structure.VerifyMessage) (*structure.VerifyResponse, error) {
	logger.AtLog.Info("VerifyMessage", zap.Any("data", data))
	logger.AtLog.Logger.Error("user try trace log data", zap.Any("data", data))

	// validate data
	if data.ETHSignature == "" || data.Signature == "" ||
		data.Address == "" || data.AddressBTC == nil || *data.AddressBTC == "" || data.AddressBTCSegwit == nil || *data.AddressBTCSegwit == "" ||
		data.MessagePrefix == nil || *data.MessagePrefix == "" {
		return nil, errors.New("invalid params")
	}

	addrr := strings.ToLower(data.Address)
	logger.AtLog.Info("wallet_address", addrr)

	user, err := u.Repo.FindUserByWalletAddress(addrr)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}
	logger.AtLog.Logger.Error("user try trace log user", zap.Any("user", user))

	isVeried, err := u.verifyBTCSegwit(user.Message, data)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}
	logger.AtLog.Info("isVeried", isVeried)
	if !isVeried {
		err := errors.New("Cannot verify wallet address")
		logger.AtLog.Error(err)
		return nil, err
	}

	if *data.AddressBTC != "" {
		user2, _ := u.Repo.FindUserByAddress(*data.AddressBTC)
		if user2 != nil {
			if user2.WalletAddressBTCTaproot == *data.AddressBTC {
				if data.Address != user2.WalletAddress {
					return nil, errors.New("invalid wallet address")
				}
			}
		}
	}

	now := time.Now()
	user.IsVerified = isVeried
	user.VerifiedAt = &now
	user.UpdatedAt = &now

	userID := user.UUID
	token, refreshToken, err := u.Auth2.GenerateAllTokens(user.WalletAddress, "", "", "", userID)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	logger.AtLog.Info("token", token)
	tokenMd5 := helpers.GenerateMd5String(token)
	logger.AtLog.Info("tokenMd5", tokenMd5)
	err = u.Cache.SetDataWithExpireTime(tokenMd5, userID, int(utils.TOKEN_CACHE_EXPIRED_TIME))
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	if data.AddressBTC != nil && *data.AddressBTC != "" {
		if user.WalletAddressBTCTaproot == "" {
			user.WalletAddressBTCTaproot = *data.AddressBTC
			logger.AtLog.Info("user.WalletAddressBTCTaproot.Updated", true)
		}
		if user.WalletAddressBTC == "" {
			user.WalletAddressBTC = *data.AddressBTC
			logger.AtLog.Info("user.WalletAddressBTC.Updated", true)
		}
	}

	if user.WalletAddressPayment == "" {
		if data.AddressPayment == "" {
			if user.WalletType != entity.WalletType_BTC_PRVKEY {
				user.WalletAddressPayment = user.WalletAddress
				logger.AtLog.Info("user.WalletAddressPayment.Updated", true)
			}
		} else {
			user.WalletAddressPayment = data.AddressPayment
			logger.AtLog.Info("user.WalletAddressPayment.Updated", true)
		}
	}

	updated, err := u.Repo.UpdateUserByWalletAddress(user.WalletAddress, user)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	logger.AtLog.Info("updated.Info", updated)
	logger.AtLog.Info("generated.Token", token)
	logger.AtLog.Info("generated.refreshToken", refreshToken)

	verified := structure.VerifyResponse{
		Token:        token,
		RefreshToken: refreshToken,
		IsVerified:   isVeried,
	}

	go u.AirdropTokenGatedNewUser(os.Getenv("AIRDROP_WALLET"), *user, 3)
	go u.AirdropArtistABNewUser(os.Getenv("AIRDROP_WALLET"), *user, 3)

	return &verified, nil
}

func buildMsgETH(taprootAddress, segwitAddress, nonceMessage string) string {
	msg := "GM.\n\nPlease sign this message to confirm your Generative wallet addresses generated by your Ethereum address.\n\n"
	msg += "Taproot address:\n" + taprootAddress
	msg += "\n\nSegwit address:\n" + segwitAddress
	msg += "\n\nNonce:\n" + nonceMessage
	msg += "\n\nThe Generative Core Team"
	return msg
}

func (u Usecase) verifyBTCSegwit(msgStr string, data structure.VerifyMessage) (bool, error) {
	// verify BTC segwit signature
	// Reconstruct the pubkey
	publicKey, wasCompressed, err := helpers.PubKeyFromSignature(data.Signature, msgStr, *data.MessagePrefix)
	if err != nil {
		return false, err
	}

	// Get the address
	var addressWitnessPubKeyHash *btcutil.AddressPubKeyHash
	if addressWitnessPubKeyHash, err = helpers.GetAddressFromPubKey(publicKey, wasCompressed); err != nil {
		return false, err
	}

	// Return nil if addresses match.
	temp := addressWitnessPubKeyHash.String()
	if temp != *data.AddressBTCSegwit {
		return false, fmt.Errorf(
			"address (%s) not found - compressed: %t\n%s was found instead",
			*data.AddressBTCSegwit,
			wasCompressed,
			addressWitnessPubKeyHash.String(),
		)
	}

	// verify ETH signature
	msg2 := buildMsgETH(*data.AddressBTC, *data.AddressBTCSegwit, msgStr)
	return u.verify(data.ETHSignature, data.Address, msg2)
}

func (u Usecase) verify(signatureHex string, signer string, msgStr string) (bool, error) {
	logger.AtLog.Info("verify", zap.String("signatureHex", signatureHex), zap.String("signer", signer), zap.String("msgStr", msgStr))
	sig := hexutil.MustDecode(signatureHex)

	msgBytes := []byte(msgStr)
	msgHash := accounts.TextHash(msgBytes)

	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		logger.AtLog.Error(err)
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	signerHex := recoveredAddr.Hex()
	isVerified := strings.ToLower(signer) == strings.ToLower(signerHex)

	logger.AtLog.Info("verify", zap.Bool("isVerified", isVerified), zap.String("signerHex", signerHex), zap.String("signatureHex", signatureHex), zap.String("signer", signer), zap.String("msgStr", msgStr), zap.Any("recoveredAddr", recoveredAddr))
	return isVerified, nil
}

func (u Usecase) UserProfile(userID string) (*entity.Users, error) {

	logger.AtLog.Info("UserProfile", zap.String("userID", userID))
	user, err := u.Repo.FindUserByID(userID)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	return user, nil
}

func (u Usecase) GetUserProfileByWalletAddress(userAddr string) (*entity.Users, error) {

	logger.AtLog.Info("GetUserProfileByWalletAddress", zap.String("userAddr", userAddr))
	user, err := u.Repo.FindUserByWalletAddress(userAddr)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	return user, nil
}

func (u Usecase) GetUserProfileByBtcAddressTaproot(userAddr string) (*entity.Users, error) {

	logger.AtLog.Info("GetUserProfileByBtcAddressTaproot", zap.String("userAddr", userAddr))
	user, err := u.Repo.FindUserByBtcAddressTaproot(userAddr)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	return user, nil
}

func (u Usecase) GetUserProfileBySlug(slug string) (*entity.Users, error) {

	logger.AtLog.Info("GetUserProfileBySlug", zap.String("slug", slug))
	user, err := u.Repo.FindUserBySlug(slug)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	return user, nil
}

func (u Usecase) GetUserProfileByBtcAddress(userAddr string) (*entity.Users, error) {

	logger.AtLog.Info("GetUserProfileByBtcAddress", zap.String("userAddr", userAddr))
	user, err := u.Repo.FindUserByBtcAddress(userAddr)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	return user, nil
}

func (u Usecase) UpdateUserProfile(userID string, data structure.UpdateProfile) (*entity.Users, error) {

	isUpdateWalletAddress := false
	oldBtcAdress := ""
	isUpdateWalletAddressPayment := false
	oldAdressPayment := ""
	user, err := u.Repo.FindUserByID(userID)
	if err != nil {
		logger.AtLog.Logger.Error("UpdateUserProfile", zap.String("action", "FindUserByID"), zap.String("userID", userID), zap.Any("data", data), zap.Error(err))
		return nil, err
	}

	if data.DisplayName != nil {
		user.DisplayName = *data.DisplayName
		user.Slug = helpers.ConvertSlug(user.DisplayName)
	}

	if data.Avatar != nil && *data.Avatar != "" {
		user.Avatar = *data.Avatar
		uploaded, err := u.UploadUserAvatar(*user)
		if err != nil {
			logger.AtLog.Logger.Error("UpdateUserProfile", zap.String("action", "UploadUserAvatar"), zap.String("userID", userID), zap.Any("data", data), zap.Error(err))
		} else {
			user.Avatar = *uploaded
		}

	}

	if data.Bio != nil {
		user.Bio = *data.Bio
	}

	if data.WalletAddressBTC != nil && *data.WalletAddressBTC != "" && strings.ToLower(user.WalletAddressBTC) != strings.ToLower(*data.WalletAddressBTC) {
		isUpdateWalletAddress = true
		oldBtcAdress = user.WalletAddressBTC
		user.WalletAddressBTC = *data.WalletAddressBTC
	}

	if data.WalletAddressPayment != nil && *data.WalletAddressPayment != "" && strings.ToLower(user.WalletAddressPayment) != strings.ToLower(*data.WalletAddressPayment) {
		isUpdateWalletAddressPayment = true
		oldAdressPayment = user.WalletAddressPayment
		user.WalletAddressPayment = *data.WalletAddressPayment
	}

	if data.ProfileSocial.Discord != nil {
		user.ProfileSocial.Discord = *data.ProfileSocial.Discord
	}

	if data.ProfileSocial.Web != nil {
		user.ProfileSocial.Web = *data.ProfileSocial.Web
	}

	needUpdateProposalForCreateNew := false
	if data.ProfileSocial.Twitter != nil && user.ProfileSocial.Twitter != *data.ProfileSocial.Twitter {
		user.ProfileSocial.Twitter = *data.ProfileSocial.Twitter
		user.ProfileSocial.TwitterVerified = false
		needUpdateProposalForCreateNew = true
	}

	if data.ProfileSocial.Medium != nil {
		user.ProfileSocial.Medium = *data.ProfileSocial.Medium
	}

	if data.ProfileSocial.Web != nil {
		user.ProfileSocial.Web = *data.ProfileSocial.Web
	}

	if data.ProfileSocial.Instagram != nil {
		user.ProfileSocial.Instagram = *data.ProfileSocial.Instagram
	}

	if data.ProfileSocial.EtherScan != nil {
		user.ProfileSocial.EtherScan = *data.ProfileSocial.EtherScan
	}

	isDisableNotification := false
	if data.EnableNotification != nil && user.EnableNotification != *data.EnableNotification {
		if !*data.EnableNotification {
			isDisableNotification = true
		}
		user.EnableNotification = *data.EnableNotification
	}
	if data.Banner != nil && *data.Banner != "" {
		user.Banner = *data.Banner
	}

	_, err = u.Repo.UpdateUserByID(userID, user)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	//update project's creator profile
	go func(user entity.Users) {

		projects, err := u.Repo.GetAllProjects(entity.FilterProjects{
			WalletAddress: &user.WalletAddress,
		})

		logger.AtLog.Logger.Info("UpdateUserProfile", zap.Any("projects", zap.Any("projects)", projects)))
		if err != nil {
			logger.AtLog.Logger.Error("UpdateUserProfile", zap.String("action", "GetAllProjects"), zap.String("userID", userID), zap.Any("data", data), zap.Error(err))
			return
		}

		for _, p := range projects {
			if p.CreatorAddrr != user.WalletAddress {
				continue
			}
			p.CreatorProfile = user

			_, err := u.Repo.UpdateProjectFields(p.UUID, map[string]interface{}{
				"creatorProfile": p.CreatorProfile,
			})

			if err != nil {
				logger.AtLog.Logger.Error("UpdateUserProfile", zap.String("action", "GetAllProjects"), zap.String("userID", userID), zap.Any("data", data), zap.Error(err))
				continue
			}

		}

	}(*user)

	logger.AtLog.Logger.Info("UpdateUserProfile", zap.String("userID", userID), zap.Any("input", data), zap.Any("user", zap.Any("user)", user)))
	if isUpdateWalletAddress {
		if user.Stats.CollectionCreated > 0 {
			go u.NotifyWithChannel(os.Getenv("SLACK_USER_CHANNEL"), fmt.Sprintf("[User BTC wallet address payment has been updated][User %s][%s]", helpers.CreateProfileLink(user.WalletAddress, user.DisplayName), user.WalletAddress), "", fmt.Sprintf("BTC wallet address payment was changed from %s to %s", oldBtcAdress, *data.WalletAddressBTC))
		}
	}

	if isUpdateWalletAddressPayment {
		if user.Stats.CollectionCreated > 0 {
			go u.NotifyWithChannel(os.Getenv("SLACK_USER_CHANNEL"), fmt.Sprintf("[User ETH wallet address payment has been updated][User %s][%s]", helpers.CreateProfileLink(user.WalletAddress, user.DisplayName), user.WalletAddress), "", fmt.Sprintf("ETH wallet address payment was changed from %s to %s", oldAdressPayment, *data.WalletAddressPayment))
		}
	}

	if isDisableNotification {
		go func() {
			_, err = u.Repo.DeleteMany(context.TODO(), entity.FirebaseRegistrationToken{}.TableName(), bson.M{"user_wallet": user.WalletAddress})
			if err != nil {
				logger.AtLog.Logger.Error("Delete FirebaseRegistrationToken failed", zap.Error(err))
			}
		}()
	}
	if needUpdateProposalForCreateNew {
		_ = u.RedisV9.DelPrefix(context.TODO(), rediskey.Beauty(entity.DaoArtist{}.TableName()).WithParams("list", user.WalletAddress).String())
		_ = u.RedisV9.DelPrefix(context.TODO(), rediskey.Beauty(entity.DaoProject{}.TableName()).WithParams("list", user.WalletAddress).String())
		go u.SetExpireYourProposalArtist(context.TODO(), user.WalletAddress)
	}

	u.Cache.SetData(helpers.GenerateUserWalletAddressKey(user.WalletAddress), user)
	return user, nil
}

func (u Usecase) Logout(accessToken string) (bool, error) {

	tokenMd5 := helpers.GenerateMd5String(accessToken)
	err := u.Cache.Delete(tokenMd5)
	if err != nil {
		logger.AtLog.Error(err)
		return false, err
	}

	return true, nil
}

func (u Usecase) ValidateAccessToken(accessToken string) (*oauth2service.SignedDetails, error) {

	//tokenMd5 := helpers.GenerateMd5String(accessToken)
	//logger.AtLog.Logger.Info("ValidateAccessToken", zap.String("ValidateAccessToken", zap.Any("accessToken)", accessToken)))

	// userID, err := u.Cache.GetData(tokenMd5)
	// if err != nil {
	// 	err = errors.New("Access token is invaild")
	// 	logger.AtLog.Logger.Error("ValidateAccessToken", zap.String("GetData", accessToken), zap.Error(err))
	// 	return nil, err

	// }

	//Claim wallet Address
	claim, err := u.Auth2.ValidateToken(accessToken)
	if err != nil {
		logger.AtLog.Logger.Error("ValidateAccessToken", zap.String("ValidateToken", accessToken), zap.Error(err))
		return nil, err
	}

	userID := &claim.Uid
	if userID == nil {
		err := errors.New("Cannot find userID")
		logger.AtLog.Logger.Error("ValidateAccessToken", zap.String("userID", accessToken), zap.Error(err))
		return nil, err
	}

	//timeT := time.Unix(claim.ExpiresAt, 0)
	return claim, err
}

func (u Usecase) UserProfileByWallet(walletAddress string) (*entity.Users, error) {

	//logger.AtLog.Info("input.walletAddress", walletAddress)
	user, err := u.Repo.FindUserByWalletAddress(walletAddress)
	if err != nil {
		logger.AtLog.Error(err)
		return nil, err
	}

	u.Cache.SetData(helpers.GenerateUserWalletAddressKey(walletAddress), user)
	return user, nil
}

func (u Usecase) UserProfileByWalletWithCache(walletAddress string) (*entity.Users, error) {
	go u.UserProfileByWallet(walletAddress)

	userCache, err := u.Cache.GetData(helpers.GenerateUserWalletAddressKey(walletAddress))
	if err != nil && userCache == nil {
		return u.UserProfileByWallet(walletAddress)
	}
	user := &entity.Users{}

	bytes := []byte(*userCache)
	err = json.Unmarshal(bytes, user)
	if err != nil {
		return u.UserProfileByWallet(walletAddress)
	}

	return user, nil
}

func (u Usecase) UploadUserAvatar(user entity.Users) (*string, error) {
	thumbnail := ""
	base64Image := user.Avatar
	i := strings.Index(base64Image, ",")
	if i > -1 {
		base64Image = base64Image[i+1:]
		name := fmt.Sprintf("thumb/%s.png", user.WalletAddress)
		uploaded, err := u.GCS.UploadBaseToBucket(base64Image, name)
		if err != nil {
			logger.AtLog.Error(err)
			return nil, err
		} else {
			logger.AtLog.Info("uploaded", uploaded)
			thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
		}

		return &thumbnail, nil
	}
	return &user.Avatar, nil
}

func (u Usecase) UpdateUserAvatars() error {
	users, err := u.Repo.GetAllUsers(entity.FilterUsers{IsUpdatedAvatar: nil})
	if err != nil {
		logger.AtLog.Error(err)
		return err
	}

	for _, user := range users {
		if user.Avatar == "" {
			user.Avatar = helpers.CreateIcon(&user.WalletAddress)
		}

		if true {
			uploadedAvatar, err := u.UploadUserAvatar(user)
			if err != nil {
				logger.AtLog.Error(err)
				continue
			}

			aUpdated := true
			user.Avatar = *uploadedAvatar
			user.IsUpdatedAvatar = &aUpdated
			updated, err := u.Repo.UpdateUserByWalletAddress(user.WalletAddress, &user)
			if err != nil {
				logger.AtLog.Logger.Error("UpdateUserByWalletAddress", zap.Error(err))
				continue
			}
			logger.AtLog.Info("updated", updated)
		}
	}
	return nil
}
