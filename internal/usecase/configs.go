package usecase

import (
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateConfig(rootSpan opentracing.Span, input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) UpdateConfig(rootSpan opentracing.Span, input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) DeleteConfig(rootSpan opentracing.Span, input structure.ConfigData) error {
	return nil
}

func (u Usecase) GetConfig(rootSpan opentracing.Span, input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) GetConfigs(rootSpan opentracing.Span, input structure.FilterCoonfigs) (*entity.Pagination, error) {
	return nil, nil
}