package usecase

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/delivery/http/request"
	"rederinghub.io/internal/entity"
)

func (u Usecase) GetFcmByUserWalletAndDeviceType(ctx context.Context, userWallet, deviceType string) (*entity.FirebaseRegistrationToken, error) {
	fcm := &entity.FirebaseRegistrationToken{}
	if err := u.Repo.FindOneBy(ctx, fcm.TableName(), bson.M{"user_wallet": userWallet, "device_type": deviceType}, fcm); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return nil, nil
	}
	return fcm, nil
}

func (u Usecase) CreateFcm(ctx context.Context, data *request.CreateFcmRequest) error {
	fcm := &entity.FirebaseRegistrationToken{
		RegistrationToken: data.RegistrationToken,
		DeviceType:        data.DeviceType,
		UserWallet:        data.UserWallet,
		CreatedAt:         time.Now(),
	}
	if err := u.Repo.InsertOne(fcm.TableName(), fcm); err != nil {
		return err
	}
	return nil
}

func (u Usecase) GetFcmByUserWallet(ctx context.Context, userWallet string) ([]*entity.FirebaseRegistrationToken, error) {
	fcms := []*entity.FirebaseRegistrationToken{}
	if err := u.Repo.Find(ctx, entity.FirebaseRegistrationToken{}.TableName(), bson.M{"user_wallet": userWallet}, &fcms); err != nil {
		return nil, err
	}
	return fcms, nil
}
