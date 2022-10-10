package googlesecret

import (
	"context"
	"google.golang.org/api/option"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type Client struct {
	googleSecretManagerClient *secretmanager.Client
}

func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	client, err := secretmanager.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{googleSecretManagerClient: client}, nil
}

func (g *Client) GetSecret(name string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	resp, err := g.googleSecretManagerClient.AccessSecretVersion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return string(resp.Payload.Data), nil
}
