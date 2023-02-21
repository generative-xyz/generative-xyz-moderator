package usecase

import (
	"github.com/opentracing/opentracing-go"
	"rederinghub.io/internal/entity"
)

func (u Usecase) ListArtist(rootSpan opentracing.Span, req entity.FilteArtist) (*entity.Pagination, error) {

	artists, err := u.Repo.ListArtist(req)
	if err != nil {
		u.Logger.Error("u.Repo.ListArtist", err.Error(), err)
		return nil, err
	}

	u.Logger.Info("artists", artists.Total)
	return artists, nil
}
