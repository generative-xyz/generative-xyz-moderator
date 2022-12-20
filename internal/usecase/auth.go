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
)


func (u Usecase) GenerateMessage(rootSpan opentracing.Span, data structure.GenerateMessage) (*string, error) {
	span, log := u.StartSpan("GenerateMessage", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	addrr := data.Address
	addrr = strings.ToLower(addrr)
	log.SetTag("wallet_address", addrr)

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

	tokenMd5 := helpers.GenerateMd5String(token)
	err = u.Cache.SetDataWithExpireTime(tokenMd5, userID, int(utils.TOKEN_CACHE_EXPIRED_TIME))
	if  err != nil {
		log.Error("Login.Redis.SetData", err.Error(), err)
		return nil, err
	}

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

func  (u Usecase) UserProfile(rootSpan opentracing.Span, userID string) (*structure.ProfileResponse, error) {
	span, log := u.StartSpan("UserProfile", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input.userID", userID)
	user, err := u.Repo.FindUserByID(userID)
	if err != nil {
		log.Error("u.Auth2.ValidateToken", err.Error(), err)
		return nil, err
	}


	log.SetTag(utils.WALLET_ADDRESS_TAG, user.WalletAddress)
	resp :=  u.profileToResp(user)
	return &resp, nil
}

func  (u Usecase) UpdateUserProfile(rootSpan opentracing.Span, userID string, data structure.UpdateProfile) (*structure.ProfileResponse, error) {
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
	
	if data.Bio != nil {
		user.Bio = *data.Bio
	}
	
	updated, err := u.Repo.UpdateUserByID(userID, user)
	if err != nil {
		log.Error("u.Repo.UpdateUserByID", err.Error(), err)
		return nil, err
	}

	log.SetData("updated", updated)

	resp :=  u.profileToResp(user)
	return &resp, nil
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

func (u Usecase) profileToResp(profile *entity.Users) structure.ProfileResponse {
	domain := os.Getenv("API_DOMAIN")

	profileAvatar := os.Getenv("DEFAUTL_AVATAR")
	if profile.Avatar != "" {
		profileAvatar = profile.Avatar
	}
	avatarURL := fmt.Sprintf("%s/files/%s", domain, profileAvatar)

	
	addr := profile.WalletAddress
	walletAddresses := []string{}
	walletAddresses = append(walletAddresses, addr)
	
	resp := structure.ProfileResponse{
		ID: profile.UUID,
		DisplayName: profile.DisplayName,
		Bio: profile.Bio,
		Avatar: avatarURL,
	}

	
	return resp
}