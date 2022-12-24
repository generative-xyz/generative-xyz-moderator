package request

import "github.com/pkg/errors"

type CreateCategoryRequest struct {
	Name *string `json:"name" validate:"required"`
}

func (r CreateCategoryRequest) Validate() error {
	
	if r.Name == nil  {
		return errors.New("Name is required")
	}

	if *r.Name == ""  {
		return errors.New("Name is not empty")
	}

	return nil
}