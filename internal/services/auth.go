package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/api"
	"rederinghub.io/internal/model"
	utils "rederinghub.io/pkg"
	"rederinghub.io/pkg/helpers"
	"rederinghub.io/pkg/interceptor"
	"rederinghub.io/pkg/logger"
)

func (s *service) GetAuthNonce(ctx context.Context, req *api.GetNonceMessageReq) (*api.GetNonceMessageResp, error) {
	addr := strings.ToLower(req.GetAddress())
	logger.AtLog.Infof("Handle [GetAuthNonce] wallet_address: %s", addr)
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		logger.AtLog.Errorf("Handle [GetAuthNonce]- rand.Read %s", err)
		return nil, err
	}

	message := fmt.Sprintf("%x-%x-%x-%x-%x",b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	message = fmt.Sprintf(utils.NONCE_MESSAGE_FORMAT, message)

	logger.AtLog.Infof("Handle [GetAuthNonce]- message %s", message)
	// find in mongo
	user, err := s.userRepository.FindUserByWalletAddress(ctx, addr)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			//insert
			user := &model.Users{}
			user.WalletAddress =  addr
			user.Message =  message
	
			id, err := s.userRepository.CreateOne(context.Background(), user )
			if err != nil {
				logger.AtLog.Errorf("Handle [GetAuthNonce]- s.userRepository.CreateOne %s", err)
				return nil, err
			}
			_ = id
			logger.AtLog.Infof("Handle [GetAuthNonce]- s.userRepository.CreateOne %s", user.UUID)

		}else{
			logger.AtLog.Errorf("Handle [GetAuthNonce]- s.userRepository.FindUserByWalletAddress %s", err)
			return nil, err
		}
	}

	user.Message = message
	err = s.userRepository.UpdateOneByID(ctx, user, user.ID)
	if err != nil {
		logger.AtLog.Errorf("Handle [GetAuthNonce]- s.userRepository.UpdateOneByID %s", err)
		return nil, err
	}

	resp := &api.GetNonceMessageResp{Message:  message}
	return resp, nil
}

func (s *service) VerifyAuthNounce(ctx context.Context,  req *api.VerifyNonceMessageReq) (*api.VerifyNonceMessageResp, error) {
	logger.AtLog.Infof("Handle [VerifyAuthNounce] - wallet address: %s - signature: %s", req.GetAddress(), req.GetSignature())
	signature := req.GetSignature()
	addrr := strings.ToLower(req.GetAddress()) 
	resp := &api.VerifyNonceMessageResp{}
	
	// find in mongo
	user, err := s.userRepository.FindUserByWalletAddress(ctx, addrr)
	if err != nil { 
		logger.AtLog.Errorf("Handle [GetAuthNonce]- userRepository.FindUserByWalletAddress %s", err)
		return resp, err
	}

	logger.AtLog.Infof("Handle [GetAuthNonce]- userRepository.FindUserByWalletAddress - USER_ID %s", user.UUID)
	isVeried, err :=   s.verify(signature, user.WalletAddress, user.Message)
	if err != nil {
		logger.AtLog.Errorf("Handle [GetAuthNonce]- s.verify: user.WalletAddress: %s, user.Message: %s - %s",user.WalletAddress, user.Message, err)
		return resp, err
	}

	//isVeried := true
	if !isVeried {
		err := errors.New("Cannot verify wallet address")
		logger.AtLog.Errorf("Handle [GetAuthNonce]- isVeried false",  err)
		return resp, err
	}

	userID := user.UUID
	token, refreshToken, err := s.auth2Service.GenerateAllTokens(user.WalletAddress, "", "", "",userID)
	if err != nil {
		logger.AtLog.Errorf("Handle [GetAuthNonce]- auth2Service.GenerateAllTokens",  err)
		return resp, err
	}

	tokenMd5 := helpers.GenerateMd5String(token)
	err = s.redisClient.Set(ctx, tokenMd5, userID, 0).Err()
	if  err != nil {
		logger.AtLog.Errorf("Handle [GetAuthNonce]- s.redisClient.Set",  err)
		return resp, err
	}

	resp.AccessToken =  token
	resp.RefreshToken =  refreshToken
	return resp, nil
}

func (s service) verify(signatureHex string, signer string, msgStr string) (bool, error) {
	logger.AtLog.Infof("Handle [VerifyAuthNounce][verify] - signatureHex: %s, signer: %s, msgStr: %s", signatureHex, signer, msgStr)
	sig := hexutil.MustDecode(signatureHex)

	msgBytes := []byte(msgStr)
	msgHash := accounts.TextHash(msgBytes)
	
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1
	}

	recovered, err := crypto.SigToPub(msgHash, sig)
	if err != nil {
		logger.AtLog.Errorf("Handle [VerifyAuthNounce][verify] -  crypto.SigToPub",  err)
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*recovered)
	signerHex := recoveredAddr.Hex()
	isVerified := strings.ToLower(signer) ==  strings.ToLower(signerHex)
	logger.AtLog.Infof("Handle [VerifyAuthNounce][verify] - signer: %s, signerHex: %s", strings.ToLower(signer), strings.ToLower(signerHex))
	return isVerified, nil
}

func (s service) GetProfile(ctx context.Context,  req *api.UserProfileReq) (*api.UserProfileResp, error) {
	userID := ctx.Value(interceptor.ContextKey("id"))
	logger.AtLog.Infof("Handle [GetProfile] - %s", userID)

	token := ctx.Value(interceptor.ContextKey("token"))
	spew.Dump(ctx)
	spew.Dump(token)
	objectID, err :=  primitive.ObjectIDFromHex(userID.(string))
	if  err != nil {
		logger.AtLog.Errorf("Handle [GetProfile]- primitive.ObjectIDFromHex %s", err)
		return nil, err
	}

	user := &model.Users{}
 	err = s.userRepository.FindById(ctx,objectID, user)
	 if  err != nil {
		logger.AtLog.Errorf("Handle [GetProfile]- .userRepository.FindById %s", err)
		return nil, err
	}

	resp := &api.UserProfileResp{
		Id: user.UUID,
		WalletAddress: user.WalletAddress,
		AvatarURL: user.AvatarURL,
		Bio: user.Bio,
		DisplayName: user.DisplayName,
	}
	return resp, nil
}

func (s service) UpdateProfile(ctx context.Context,  req *api.UpdateUserProfileReq) (*api.UserProfileResp, error) {
	userID := ctx.Value(interceptor.ContextKey("id"))
	logger.AtLog.Infof("Handle [UpdateProfile] - %s", userID)	
	objectID, err :=  primitive.ObjectIDFromHex(userID.(string))
	if  err != nil {
		logger.AtLog.Errorf("Handle [UpdateProfile]- .userRepository.FindById %s", err)
		return nil, err
	}

	user := &model.Users{}
 	err = s.userRepository.FindById(ctx,objectID, user)
	 if  err != nil {
		logger.AtLog.Errorf("Handle [UpdateProfile]- .userRepository.FindById %s", err)
		return nil, err
	}

	if req.GetAvatarURL() != "" {
		user.DisplayName = req.GetAvatarURL()
	}
	
	if req.GetBio() != "" {
		user.Bio = req.GetBio()
	}
	
	if req.GetDisplayName() != "" {
		user.DisplayName = req.GetDisplayName()
	}

	err = s.userRepository.UpdateOneByID(ctx, user, user.ID)
	if  err != nil {
		logger.AtLog.Errorf("Handle [UpdateProfile]- .userRepository.UpdateOneByID %s", err)
		return nil, err
	}

	resp := &api.UserProfileResp{
		Id: user.UUID,
		WalletAddress: user.WalletAddress,
		AvatarURL: user.AvatarURL,
		Bio: user.Bio,
		DisplayName: user.DisplayName,
	}
	return resp, nil
}