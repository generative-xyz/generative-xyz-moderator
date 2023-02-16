package usecase

import (
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/oauth2service"
)


func (u Usecase) GenerateMessage(rootSpan opentracing.Span, data structure.GenerateMessage) (*string, error) {
	span, log := u.StartSpan("GenerateMessage", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	addrr := data.Address
	addrr = strings.ToLower(addrr)
	log.SetTag(utils.WALLET_ADDRESS_TAG, addrr)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Error("rand.Read", err.Error(), err)
		return nil, err
	}
	message := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])


	message = fmt.Sprintf(utils.NONCE_MESSAGE_FORMAT, message)
	log.SetData("message", message)
	
	now := time.Now().UTC()
	user, err := u.Repo.FindUserByWalletAddress(addrr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			//insert
			user := &entity.Users{}
			user.WalletAddress =  addrr
			user.Message = message
			user.CreatedAt = &now
	
			log.SetData("inserted.User", user)
			err = u.Repo.CreateUser(user)
			if err != nil {
				log.Error("u.Repo.CreateUser", err.Error(), err)
				return nil, err
			}

			return &message, nil

		}else{
			log.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
			return nil, err
		}
	}

	log.SetData("user", user)
	user.Message = message
	user.UpdatedAt = &now
	user.IsVerified = false
	updated, err := u.Repo.UpdateUserByWalletAddress(addrr, user)
	if err != nil {
		return nil, err
	}
	
	log.SetData("updated",updated)
	log.SetData("updated.User",message)
	return &message, nil
}

func (u Usecase) VerifyMessage(rootSpan opentracing.Span, data structure.VerifyMessage) (*structure.VerifyResponse, error) {
	span, log := u.StartSpan("VerifyMessage", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	
	log.SetData("input", data)
	addrr := strings.ToLower(data.Address) 
	signature := data.Signature
	log.SetData("wallet_address", addrr)

	user, err := u.Repo.FindUserByWalletAddress(addrr)
	if err != nil {
		log.Error("u.Repo.FindUserByWalletAddress", err.Error(), err)
		return nil, err
	}
	log.SetData("user", user)

	isVeried, err :=   u.verify(span, signature, user.WalletAddress, user.Message)
	if err != nil {
		log.Error("u.verify", err.Error(), err)
		return nil, err
	}
	log.SetData("isVeried", isVeried)

	//isVeried := true
	if !isVeried {
		err := errors.New("Cannot verify wallet address")
		log.Error("u.verify", err.Error(), err)
		return nil, err
	}

	now := time.Now()
	user.IsVerified = isVeried
	user.VerifiedAt = &now
	user.UpdatedAt = &now

	userID := user.UUID
	token, refreshToken, err := u.Auth2.GenerateAllTokens(user.WalletAddress, "", "", "", userID)
	if err != nil {
		log.Error("u.Auth2.GenerateAllTokens", err.Error(), err)
		return nil, err
	}

	log.SetData("token", token)
	tokenMd5 := helpers.GenerateMd5String(token)
	log.SetData("tokenMd5", tokenMd5)
	err = u.Cache.SetDataWithExpireTime(tokenMd5, userID, int(utils.TOKEN_CACHE_EXPIRED_TIME))
	if  err != nil {
		log.Error("Login.Redis.SetData", err.Error(), err)
		return nil, err
	}

	//TODO - Tri temporay comment
	// if data.AddressBTC !=  nil {
	// 	if *data.AddressBTC  != "" {
	// 		if user.WalletAddressBTC == ""  {
	// 			user.WalletAddressBTC = *data.AddressBTC
	// 			log.SetData("user.WalletAddressBTC.Updated", true)
	// 		}
	// 	}
	// }

	updated, err := u.Repo.UpdateUserByWalletAddress(user.WalletAddress, user)
	if err != nil {
		log.Error("u.Repo.UpdateUserByWalletAddress", err.Error(), err)
		return nil, err
	}

	log.SetData("updated.Info", updated)
	log.SetData("generated.Token", token)
	log.SetData("generated.refreshToken", refreshToken)

	verified := structure.VerifyResponse{
		Token:  token,
		RefreshToken:  refreshToken,
		IsVerified: isVeried,
	}

	return &verified, nil
}

func (u Usecase) verify(rootSpan opentracing.Span, signatureHex string, signer string, msgStr string) (bool, error) {
	span, log := u.StartSpan("verify", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	
	log.SetData("input.signatureHex", signatureHex)
	log.SetData("input.signer", signer)
	log.SetData("input.msgStr", msgStr)
	
	log.SetTag(utils.WALLET_ADDRESS_TAG, signer)
	
	sig := hexutil.MustDecode(signatureHex)

	msgBytes := []byte(msgStr)
	msgHash := accounts.TextHash(msgBytes)
	
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		log.Error("crypto.SigToPub", err.Error(), err)
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	signerHex := recoveredAddr.Hex()
	isVerified := strings.ToLower(signer) ==  strings.ToLower(signerHex)

	log.SetData("recoveredAddr", recoveredAddr)
	log.SetData("signerHex", signerHex)
	log.SetData("isVerified", isVerified)
	return isVerified, nil
}

func  (u Usecase) UserProfile(rootSpan opentracing.Span, userID string) (*entity.Users, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input.userID", userID)
	user, err := u.Repo.FindUserByID(userID)
	if err != nil {
		log.Error("u.Auth2.ValidateToken", err.Error(), err)
		return nil, err
	}

	log.SetTag(utils.WALLET_ADDRESS_TAG, user.WalletAddress)
	return user, nil
}

func (u Usecase) GetUserProfileByWalletAddress(rootSpan opentracing.Span, userAddr string) (*entity.Users, error) {
	span, log := u.StartSpan("GetUserProfileByWalletAddress", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input.userAddr", userAddr)
	user, err := u.Repo.FindUserByWalletAddress(userAddr)
	if err != nil {
		log.Error("u.GetUserProfileByWalletAddress", err.Error(), err)
		return nil, err
	}

	log.SetTag(utils.WALLET_ADDRESS_TAG, user.WalletAddress)
	return user, nil
}

func  (u Usecase) UpdateUserProfile(rootSpan opentracing.Span, userID string, data structure.UpdateProfile) (*entity.Users, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input.UserID", userID)
	log.SetData("input.data", data)

	user, err := u.Repo.FindUserByID(userID)
	if err != nil {
		log.Error("u.Repo.FindUserByID", err.Error(), err)
		return nil, err
	}

	log.SetTag(utils.WALLET_ADDRESS_TAG, user.WalletAddress)
	if data.DisplayName != nil {
		user.DisplayName = *data.DisplayName
	}
	
	if data.Avatar != nil && *data.Avatar != "" {
		user.Avatar = *data.Avatar
		uploaded, err := u.UploadUserAvatar(span, *user)
		if err != nil {
			log.Error("u.UploadUserAvatar", err.Error(), err)
		}else{
			user.Avatar = *uploaded
		}
		
	}
	
	if data.Bio != nil {
		user.Bio = *data.Bio
	}
	
	if data.ProfileSocial.Discord != nil {
		user.ProfileSocial.Discord = *data.ProfileSocial.Discord
	}
	
	if data.ProfileSocial.Web != nil {
		user.ProfileSocial.Web = *data.ProfileSocial.Web
	}
	
	if data.ProfileSocial.Twitter != nil {
		user.ProfileSocial.Twitter = *data.ProfileSocial.Twitter
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
	
	updated, err := u.Repo.UpdateUserByID(userID, user)
	if err != nil {
		log.Error("u.Repo.UpdateUserByID", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)
	return user, nil
}

func  (u Usecase) Logout(rootSpan opentracing.Span, accessToken string) (bool, error) {
	span, log := u.StartSpan("Logout", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	tokenMd5 := helpers.GenerateMd5String(accessToken)
	err := u.Cache.Delete(tokenMd5)
	if err != nil {
		log.Error("u.Cache.Delete", err.Error(), err)
		return false, err
	}
	
	return true, nil
}

func  (u Usecase) ValidateAccessToken(rootSpan opentracing.Span, accessToken string) (*oauth2service.SignedDetails, error) {
	span, log := u.StartSpan("ValidateAccessToken", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	tokenMd5 := helpers.GenerateMd5String(accessToken)
	log.SetData("tokenMd5", tokenMd5)

	userID, err := u.Cache.GetData(tokenMd5)
	if err != nil {
		err = errors.New("Access token is invaild")
		log.Error("u.Cache.GetData", err.Error(), err)
		return nil, err
	}

	log.SetData("cached.UserID", userID)
	
	//Claim wallet Address
	claim, err := u.Auth2.ValidateToken(accessToken)
	if err != nil {
		log.Error("u.Auth2.ValidateToken", err.Error(), err)
		return nil, err
	}

	if userID == nil {
		err := errors.New("Cannot find userID")
		log.Error("userID.Empty",err.Error(), err)
		return nil, err
	}

	timeT := time.Unix(claim.ExpiresAt, 0)
	log.SetData("claim.Exp", timeT)
	log.SetData("claim", claim)
	return claim, err
}

func  (u Usecase) UserProfileByWallet(rootSpan opentracing.Span, walletAddress string) (*entity.Users, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input.walletAddress", walletAddress)
	user, err := u.Repo.FindUserByWalletAddress(walletAddress)
	if err != nil {
		log.Error("u.Auth2.ValidateToken", err.Error(), err)
		return nil, err
	}


	log.SetTag(utils.WALLET_ADDRESS_TAG, user.WalletAddress)
	return user, nil
}

func (u Usecase) UploadUserAvatar(rootSpan opentracing.Span, user entity.Users) (*string, error) {
	span, log := u.StartSpan("UploadUserAvatar", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	thumbnail := ""
	base64Image := user.Avatar
	i := strings.Index(base64Image, ",")
	if i > -1 {
		base64Image = base64Image[i+1:]
		name := fmt.Sprintf("thumb/%s.png", user.WalletAddress)
		uploaded, err := u.GCS.UploadBaseToBucket(base64Image,  name)
		if err != nil {
			log.Error("u.GCS.UploadBaseToBucket", err.Error(), err)
			return nil, err
		}else{
			log.SetData("uploaded", uploaded)
			thumbnail = fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), name)
		}
	
		return &thumbnail, nil
	}
	return &user.Avatar, nil
}

func (u Usecase) UpdateUserAvatars(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("UpdateUserAvatars", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	users, err := u.Repo.GetAllUsers(entity.FilterUsers{IsUpdatedAvatar: nil})
	if err != nil {
		log.Error("UpdateUserAvatars.GetAllUsers", err.Error(), err)
		return  err
	}

	for _, user := range users {
		if user.Avatar == "" {
			user.Avatar = helpers.CreateIcon(&user.WalletAddress)
		}

		if true {
			uploadedAvatar, err := u.UploadUserAvatar(span, user)
			if err != nil {
				log.Error("UpdateUserAvatars.UploadUserAvatar", err.Error(), err)
				continue
			}

			aUpdated := true
			user.Avatar = *uploadedAvatar
			user.IsUpdatedAvatar = &aUpdated
			updated, err := u.Repo.UpdateUserByWalletAddress(user.WalletAddress, &user)
			if err != nil {
				log.Error("UpdateUserAvatars.UpdateUserByWalletAddress", err.Error(), err)
				continue
			}
			log.SetData("updated", updated)
		}
	}
	return nil
}
