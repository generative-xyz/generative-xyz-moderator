package usecase

import (
	"github.com/opentracing/opentracing-go"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) CreateProject(rootSpan opentracing.Span,  req structure.CreateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("CreateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	resp := &entity.Projects{}
	

	return resp, nil
}

func (u Usecase) UpdateProject(rootSpan opentracing.Span,  req structure.UpdateProjectReq) (*entity.Projects, error) {
	span, log := u.StartSpan("UpdateProject", rootSpan)
	defer u.Tracer.FinishSpan(span, log )
	resp := &entity.Projects{}
	

	return resp, nil
}
