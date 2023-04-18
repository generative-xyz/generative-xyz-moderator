package structure

import (
	"errors"
	"strings"
)

type AISchoolModelParams struct {
	Name               string                     `json:"model_name"`
	InputDimension     []int                      `json:"input_dim"`
	Structure          []int                      `json:"structure"`
	ActivationFunction string                     `json:"activation_name"`
	Epoch              int                        `json:"epoch_num"`
	BatchSize          int                        `json:"batch_size"`
	ValidationPercent  float64                    `json:"val_percent"`
	Augmentation       *AISchoolDataAugmentConfig `json:"data_augmentation_config"`
}

type AISchoolDataAugmentConfig struct {
	AugmentRandomFlip       string  `json:"random_flip"`
	AugmentRandomRotate     float64 `json:"random_rotation"`
	AugmentRandomZoom       float64 `json:"random_zoom"`
	AugmentRandomContrast   float64 `json:"random_contrast"`
	AugmentRandomBrightness float64 `json:"random_brightness"`
}

func (params AISchoolModelParams) SelfValidate() error {
	if params.Name == "" {
		return errors.New("Name must be filled")
	}

	if len(params.InputDimension) != 2 {
		return errors.New("InputDimension must be 2")
	}

	if params.InputDimension[0] < 10 || params.InputDimension[0] > 32 {
		return errors.New("DimensionWidth must be between 10 and 32")
	}

	if params.InputDimension[1] < 10 || params.InputDimension[1] > 32 {
		return errors.New("DimensionHeight must be between 10 and 32")
	}

	if len(params.Structure) > 10 {
		return errors.New("HiddenLayer must be less than 10")
	} else {
		for _, layer := range params.Structure {
			if layer < 1 || layer > 20 {
				return errors.New("Structure must be between 1 and 20")
			}
		}
	}

	if params.ActivationFunction == "" {
		return errors.New("ActivationFunction must be filled")
	}

	switch strings.ToLower(params.ActivationFunction) {
	case "relu", "sigmoid", "tanh", "leaky_relu":
	default:
		return errors.New("ActivationFunction must be relu, sigmoid, tanh, or leaky_relu")
	}

	if params.ValidationPercent < 1 || params.ValidationPercent > 20 {
		return errors.New("ValidationPercent must be between 1 and 20")
	}

	if params.Epoch < 1 || params.Epoch > 30 {
		return errors.New("Epoch must be between 1 and 30")
	}
	if params.BatchSize < 1 || params.BatchSize > 512 {
		return errors.New("BatchSize must be between 1 and 512")
	}

	if params.Augmentation != nil {
		switch strings.ToLower(params.Augmentation.AugmentRandomFlip) {
		case "horizontal", "vertical", "both", "none":
		default:
			return errors.New("AugmentRandomFlip must be horizontal, vertical, both, or none")
		}
		if params.Augmentation.AugmentRandomRotate < 0 || params.Augmentation.AugmentRandomRotate > 0.5 {
			return errors.New("AugmentRandomRotate must be between 1 and 0.5")
		}
		if params.Augmentation.AugmentRandomZoom < 0 || params.Augmentation.AugmentRandomZoom > 0.5 {
			return errors.New("AugmentRandomZoom must be between 1 and 0.5")
		}
		if params.Augmentation.AugmentRandomContrast < 0 || params.Augmentation.AugmentRandomContrast > 0.5 {
			return errors.New("AugmentRandomContrast must be between 1 and 0.5")
		}
		if params.Augmentation.AugmentRandomBrightness < 0 || params.Augmentation.AugmentRandomBrightness > 0.5 {
			return errors.New("AugmentRandomBrightness must be between 1 and 0.5")
		}
	}

	return nil
}
