package usecase

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"rederinghub.io/utils"
	"rederinghub.io/utils/redis"
	"strings"
	"time"

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

	if req.CaptureImageTime == nil {
		cap := entity.DEFAULT_CAPTURE_TIME
		pe.CatureThumbnailDelayTime = &cap
	}

	pe.IsBigFile = false //wil be updated by pubsub - soon
	pe.IsHidden = true
	pe.Status = false
	pe.IsSynced = false
	pe.TxHash = strings.ToLower(pe.TxHash)
	pe.TxHex = strings.ToLower(pe.TxHex)
	pe.CommitTxHash = strings.ToLower(pe.CommitTxHash)
	pe.RevealTxHash = strings.ToLower(pe.RevealTxHash)

	pe.TokenID = pe.TxHash
	pe.TokenId = pe.TxHash
	pe.ContractAddress = strings.ToLower(os.Getenv("GENERATIVE_PROJECT"))

	err = u.Repo.CreateProject(pe)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
		return nil, err
	}
	logger.AtLog.Logger.Info(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Any("project", pe))

	//process ziplink
	if req.ZipLink != nil && *req.ZipLink != "" {
		//move them to pubsub to prevent 502 error
		err := u.PubSub.Producer(utils.PUBSUB_ETH_PROJECT_UNZIP,
			redis.PubSubPayload{
				Data: structure.ProjectUnzipPayload{
					ProjectID: pe.TxHash,
					ZipLink:   *req.ZipLink}},
		)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
			return nil, err
		}
		pe.HasZipFile = true

		updatedField := make(map[string]interface{})
		updatedField["isFullChain"] = true
		_, err = u.Repo.UpdateProjectFields(pe.UUID, updatedField)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
		}
	}

	return pe, nil
}

func (u Usecase) ProcessEthZip(zipLink string) ([]string, uint64, error) {
	maxSize := uint64(0)
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
		return nil, maxSize, err
	}

	unzippedFolder := strings.TrimSuffix(contentPath+"_unzip", filepath.Ext(zipLink))
	files, err := u.GCS.ReadFolder(unzippedFolder)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("ProcessEthZip.ReadFolder.%s", zipLink), zap.String("zipLink", zipLink), zap.String("unzippedFolder", unzippedFolder))
		return nil, maxSize, err
	}

	for _, file := range files {
		if strings.Index(file.Name, ".json") == -1 {
			continue
		}

		if uint64(file.Size) > maxSize {
			maxSize = uint64(file.Size)
		}

		path := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), file.Name)
		resp = append(resp, path)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(resp), func(i, j int) { resp[i], resp[j] = resp[j], resp[i] })

	return resp, maxSize, nil
}
