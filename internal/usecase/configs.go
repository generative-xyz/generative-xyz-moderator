package usecase

import (
	"errors"

	"github.com/jinzhu/copier"
	"github.com/opentracing/opentracing-go"
	"go.mongodb.org/mongo-driver/mongo"
	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateConfig(rootSpan opentracing.Span, input structure.ConfigData) (*entity.Configs, error) {
	span, log := u.StartSpan("CreateConfig", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)
	config := &entity.Configs{
		Key: input.Key,
		Value: input.Value,
	}

	conf, err := u.Repo.FindConfig(input.Key)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err := u.Repo.InsertConfig(config)
			if err != nil {
				log.Error(" u.Repo.InsertConfig", err.Error(), err)
				return nil, err
			}
		}
	}

	conf.Value = input.Value
	updated, err := u.Repo.UpdateConfig(input.Key, conf)
	if err != nil {
		log.Error(" u.Repo.UpdateConfig", err.Error(), err)
		return nil, err
	}

	log.SetData("updated",updated)
	return config, nil
}

func (u Usecase) UpdateConfig(rootSpan opentracing.Span, input structure.ConfigData) (*entity.Configs, error) {
	return nil, nil
}

func (u Usecase) DeleteConfig(rootSpan opentracing.Span, input string) error {
	span, log := u.StartSpan("GetConfig", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)
	deleted, err := u.Repo.DeleteConfig(input)
	if err != nil {
		log.Error(" u.Repo.DeleteConfig", err.Error(), err)
		return err
	}
	log.SetData("deleted",deleted)

	return nil
}

func (u Usecase) GetConfig(rootSpan opentracing.Span, input string) (*entity.Configs, error) {
	span, log := u.StartSpan("GetConfig", rootSpan)
	defer u.Tracer.FinishSpan(span, log )

	log.SetData("input", input)

	config, err := u.Repo.FindConfig(input)
	if err != nil {
		log.Error(" u.Repo.FindConfig", err.Error(), err)
		return nil, err
	}

	return config, nil
}

func (u Usecase) GetConfigs(rootSpan opentracing.Span, input structure.FilterConfigs) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetConfigs", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	f := &entity.FilterConfigs{}
	err := copier.Copy(f, input)
	if err != nil {
		log.Error("copier.Copy", err.Error(), err)
		return nil, err
	}

	confs,  err := u.Repo.ListConfigs(*f)
	if err != nil {
		log.Error(" u.Repo.FindConfig", err.Error(), err)
		return nil, err
	}

	return confs, nil

}