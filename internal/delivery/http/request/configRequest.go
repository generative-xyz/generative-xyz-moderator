package request

import "github.com/pkg/errors"

type CreateConfigRequest struct {
	Key *string `json:"key" validate:"required"`
	Value *string `json:"value" validate:"required"`
}

func (r CreateConfigRequest) Validate() error {
	if r.Key == nil  {
		return errors.New("Key is required")
	}
	
	if *r.Key == ""  {
		return errors.New("Key is not empty")
	}
	
	if r.Value == nil  {
		return errors.New("Value is required")
	}

	if *r.Value == ""  {
		return errors.New("Value is not empty")
	}

	return nil
}