package services

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TemplateNotFoundError struct {
	TokenID string
	ChainID string
}

func (e TemplateNotFoundError) Error() string {
	return fmt.Sprintf("template with tokenId %s and chainId %s not found", e.TokenID, e.ChainID)
}

func (e TemplateNotFoundError) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

type InvalidTokenIDError struct {
	TokenID string
}

func (e InvalidTokenIDError) Error() string {
	return fmt.Sprintf("invalid token id: %s", e.TokenID)
}

func (e InvalidTokenIDError) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}
