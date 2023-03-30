package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/copier"

	"go.uber.org/zap"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/logger"
)

func (u Usecase) CreateProject(req structure.CreateProjectReq) (*entity.Projects, error) {
	pe := &entity.Projects{}
	err := copier.Copy(pe, req)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
		return nil, err
	}

	//process ziplink
	if req.ZipLink != nil && *req.ZipLink != "" {
		imageLinks, err := u.ProcessEthZip(*req.ZipLink)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.ProcessEthZip.%s", pe.TokenId), zap.String("zipLink", *req.ZipLink), zap.Error(err))
			return nil, err
		}
		pe.Images = imageLinks
	}

	pe.IsHidden = true
	pe.Status = false
	pe.IsSynced = false
	pe.TxHash = strings.ToLower(pe.TxHash)

	pe.TokenID = pe.TxHash
	pe.TokenId = pe.TxHash
	pe.ContractAddress = os.Getenv("GENERATIVE_PROJECT")

	err = u.Repo.CreateProject(pe)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
		return nil, err
	}
	logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Any("project", pe))
	return pe, nil
}

func (u Usecase) ProcessEthZip(zipLink string) ([]string, error) {
	resp := []string{}
	linksArr := strings.Split(zipLink, "/")

	contentPath := ""
	for _, path := range linksArr {
		if strings.Index(path, "http") != -1 {
			continue
		}
		
		if strings.Index(path, "storage.googleapis.com") != -1 {
			continue
		}
		
		if strings.Index(path, os.Getenv("GCS_BUCKET")) != -1 {
			continue
		}

		if strings.Index(path, os.Getenv("GCS_DOMAIN")) != -1 || strings.Index(os.Getenv("GCS_DOMAIN"), path) != -1 {
			continue
		}

		if path == "" {
			continue
		}

		prefix := ""
		if contentPath != "" {
			prefix = "/"
		}

		contentPath += prefix + path
	}

	err := u.GCS.UnzipFile(contentPath)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ProcessEthZip.UnzipFile.%s", zipLink), zap.String("zipLink", zipLink), zap.String("contentPath", contentPath))
		return nil, err
	}

	unzippedFolder := strings.TrimSuffix(contentPath+"_unzip", filepath.Ext(zipLink))
	files, err := u.GCS.ReadFolder(unzippedFolder)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ProcessEthZip.ReadFolder.%s", zipLink), zap.String("zipLink", zipLink), zap.String("unzippedFolder", unzippedFolder))
		return nil, err
	}

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), file.Name)
		resp = append(resp, path)
	}

	return resp, nil
}
