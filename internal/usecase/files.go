package usecase

import (
	"fmt"
	"net/http"
	"os"

	"github.com/opentracing/opentracing-go"

	"rederinghub.io/internal/entity"
	"rederinghub.io/utils/googlecloud"
)

func (u Usecase) UploadFile(rootSpan opentracing.Span, r *http.Request) (*entity.Files, error) {
	span, log := u.StartSpan("UploadFile", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	_, handler, err := r.FormFile("file")
	if err != nil {
		log.Error("r.FormFile.File", err.Error(), err)
		return nil, err
	}

	gf := googlecloud.GcsFile{
		FileHeader: handler,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		log.Error("u.GCS.FileUploadToBucke", err.Error(), err)
		return nil, err
	}

	log.SetData("uploaded", uploaded)
	

	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName:   uploaded.Name,
		FileSize:   int(uploaded.Size),
		MineType:   uploaded.Minetype,
		URL:  cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		log.Error("u.Repo.InsertOne", err.Error(), err)
		return nil, err
	}

	log.SetData("inserted.FileModel", fileModel)
	return fileModel, nil
}
