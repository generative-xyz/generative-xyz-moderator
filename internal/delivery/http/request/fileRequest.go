package request

import "errors"

type CreateMultipartUploadRequest struct {
	FileName string `json:"fileName"`
	Group    string `json:"group"`
}

func (g CreateMultipartUploadRequest) SelfValidate() error {
	if g.FileName == "" {
		return errors.New("fileName is required")
	}

	if g.Group == "" {
		return errors.New("group should not be empty")
	}

	return nil
}

type FileResize struct {
	File     string `json:"file"`
	FileName string `json:"fileName"`
}
