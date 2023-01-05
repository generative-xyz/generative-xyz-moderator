package usecase

import "github.com/opentracing/opentracing-go"

func (u Usecase) GetRedis(rootSpan opentracing.Span, key string) (*string, error) {
	span, log := u.StartSpan("GetRedis", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var res *string
	var err error

	res, err = u.Cache.GetData(key)

	return res, err
}

func (u Usecase) UpsertRedis(rootSpan opentracing.Span, key string, value string) (*string, error) {
	span, log := u.StartSpan("UpsertRedis", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var res *string
	var err error

	err = u.Cache.SetStringData(key, value)
	if err != nil {
		return nil, err
	}

	res, err = u.Cache.GetData(key)

	return res, err
}

func (u Usecase) DeleteRedis(rootSpan opentracing.Span, key string) error {
	span, log := u.StartSpan("DeleteRedis", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
	var err error
	err = u.Cache.Delete(key)
	return err
}
