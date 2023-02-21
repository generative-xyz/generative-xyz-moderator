package usecase

import (
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
)

func (u Usecase) ListArtist(rootSpan opentracing.Span, req entity.FilteArtist) (*entity.Pagination, error) {
	span, log := u.StartSpan("GetProjects", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	artists, err := u.Repo.ListArtist(req)
	if err != nil {
		log.Error("u.Repo.ListArtist", err.Error(), err)
		return nil, err
	}

	log.SetData("artists", artists.Total)
	return artists, nil
}
