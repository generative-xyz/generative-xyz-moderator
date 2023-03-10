package firebase

import (
	"context"
	"encoding/base64"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"rederinghub.io/utils/logger"
)

type Service interface {
	SendMessagesToSpecificDevices(ctx context.Context, registrationToken string, msg map[string]string) error
	SendMessagesToMultipleDevices(ctx context.Context, registrationTokens []string, msg map[string]string) error
	SendMessageToTopic(ctx context.Context, topic string, msg map[string]string) error
}

type serviceImpl struct {
	app *firebase.App
}

func NewService(authKey string) (Service, error) {
	jsonKey, _ := base64.StdEncoding.DecodeString(authKey)
	opt := option.WithCredentialsJSON([]byte(jsonKey))
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return &serviceImpl{app}, nil
}
func (s *serviceImpl) SendMessagesToSpecificDevices(ctx context.Context, registrationToken string, msg map[string]string) error {
	msgClient, err := s.app.Messaging(ctx)
	if err != nil {
		return err
	}
	message := &messaging.Message{
		Data:  msg,
		Token: registrationToken,
	}
	response, err := msgClient.Send(ctx, message)
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("Successfully sent message", zap.String("response", response))
	return nil
}
func (s *serviceImpl) SendMessagesToMultipleDevices(ctx context.Context, registrationTokens []string, msg map[string]string) error {
	msgClient, err := s.app.Messaging(ctx)
	if err != nil {
		return err
	}
	message := &messaging.MulticastMessage{
		Data:   msg,
		Tokens: registrationTokens,
	}
	br, err := msgClient.SendMulticast(ctx, message)
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("Successfully sent message", zap.Any("response", br))
	return nil
}
func (s *serviceImpl) SendMessageToTopic(ctx context.Context, topic string, msg map[string]string) error {
	msgClient, err := s.app.Messaging(ctx)
	if err != nil {
		return err
	}
	message := &messaging.Message{
		Data:  msg,
		Topic: topic,
	}
	response, err := msgClient.Send(ctx, message)
	if err != nil {
		return err
	}
	logger.AtLog.Logger.Info("Successfully sent message", zap.String("response", response))
	return nil
}
