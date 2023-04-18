package structure

import "errors"

type AISchoolModelParams struct {
	Name                    string  `json:"name"`
	DimensionWidth          int     `json:"dimension_width"`
	DimensionHeight         int     `json:"dimension_height"`
	HiddenLayer             []int   `json:"hidden_layer"`
	ActivationFunction      string  `json:"activation_function"`
	Epoch                   int     `json:"epoch"`
	BatchSize               int     `json:"batch_size"`
	ValidationSet           float64 `json:"validation_set"`
	Augmentation            bool    `json:"augmentation"`
	AugmentRandomFlip       string  `json:"augment_random_flip"`
	AugmentRandomRotate     float64 `json:"augment_random_rotate"`
	AugmentRandomZoom       float64 `json:"augment_random_zoom"`
	AugmentRandomContrast   float64 `json:"augment_random_contrast"`
	AugmentRandomBrightness float64 `json:"augment_random_brightness"`
}

func (params AISchoolModelParams) SelfValidate() error {
	if params.Name == "" {
		return errors.New("Name must be filled")
	}

	if params.DimensionWidth < 10 || params.DimensionWidth > 32 {
		return errors.New("DimensionWidth must be between 10 and 32")
	}

	if params.DimensionHeight < 10 || params.DimensionHeight > 32 {
		return errors.New("DimensionHeight must be between 10 and 32")
	}

	if len(params.HiddenLayer) > 10 {
		return errors.New("HiddenLayer must be less than 10")
	} else {
		for _, hiddenLayer := range params.HiddenLayer {
			if hiddenLayer < 1 || hiddenLayer > 20 {
				return errors.New("HiddenLayer must be between 1 and 20")
			}
		}
	}

	if params.ActivationFunction == "" {
		return errors.New("ActivationFunction must be filled")
	}

	switch params.ActivationFunction {
	case "relu", "sigmoid", "tanh", "leaky_relu":
	default:
		return errors.New("ActivationFunction must be relu, sigmoid, tanh, or leaky_relu")
	}

	if params.ValidationSet < 1 || params.ValidationSet > 20 {
		return errors.New("ValidationSet must be between 1 and 20")
	}

	if params.Epoch < 1 || params.Epoch > 30 {
		return errors.New("Epoch must be between 1 and 30")
	}
	if params.BatchSize < 1 || params.BatchSize > 512 {
		return errors.New("BatchSize must be between 1 and 512")
	}

	if params.Augmentation {
		switch params.AugmentRandomFlip {
		case "horizontal", "vertical", "both", "none":
		default:
			return errors.New("AugmentRandomFlip must be horizontal, vertical, both, or none")
		}
		if params.AugmentRandomRotate < 0 || params.AugmentRandomRotate > 0.5 {
			return errors.New("AugmentRandomRotate must be between 1 and 0.5")
		}
		if params.AugmentRandomZoom < 0 || params.AugmentRandomZoom > 0.5 {
			return errors.New("AugmentRandomZoom must be between 1 and 0.5")
		}
		if params.AugmentRandomContrast < 0 || params.AugmentRandomContrast > 0.5 {
			return errors.New("AugmentRandomContrast must be between 1 and 0.5")
		}
		if params.AugmentRandomBrightness < 0 || params.AugmentRandomBrightness > 0.5 {
			return errors.New("AugmentRandomBrightness must be between 1 and 0.5")
		}
	}

	return nil
}

type AISchoolModelParamsExec struct {
	Name                    string  `json:"name"`
	DimensionWidth          int     `json:"dimension_width"`
	DimensionHeight         int     `json:"dimension_height"`
	HiddenLayer             []int   `json:"hidden_layer"`
	ActivationFunction      string  `json:"activation_function"`
	Epoch                   int     `json:"epoch"`
	BatchSize               int     `json:"batch_size"`
	ValidationSet           float64 `json:"validation_set"`
	Augmentation            bool    `json:"augmentation"`
	AugmentRandomFlip       string  `json:"augment_random_flip"`
	AugmentRandomRotate     float64 `json:"augment_random_rotate"`
	AugmentRandomZoom       float64 `json:"augment_random_zoom"`
	AugmentRandomContrast   float64 `json:"augment_random_contrast"`
	AugmentRandomBrightness float64 `json:"augment_random_brightness"`
}

func (params AISchoolModelParams) TransformToExec() *AISchoolModelParamsExec {
	return &AISchoolModelParamsExec{
		Name:                    params.Name,
		DimensionWidth:          params.DimensionWidth,
		DimensionHeight:         params.DimensionHeight,
		HiddenLayer:             params.HiddenLayer,
		ActivationFunction:      params.ActivationFunction,
		Epoch:                   params.Epoch,
		BatchSize:               params.BatchSize,
		ValidationSet:           params.ValidationSet,
		Augmentation:            params.Augmentation,
		AugmentRandomFlip:       params.AugmentRandomFlip,
		AugmentRandomRotate:     params.AugmentRandomRotate,
		AugmentRandomZoom:       params.AugmentRandomZoom,
		AugmentRandomContrast:   params.AugmentRandomContrast,
		AugmentRandomBrightness: params.AugmentRandomBrightness,
	}
}
