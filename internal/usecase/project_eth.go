package usecase

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
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

	//process ziplink
	if req.ZipLink != nil && *req.ZipLink != "" {
		imageLinks, maxSize, err := u.ProcessEthZip(*req.ZipLink)
		if err != nil {
			logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.ProcessEthZip.%s", pe.TokenId), zap.String("zipLink", *req.ZipLink), zap.Error(err))
			return nil, err
		}
		pe.Images = imageLinks
		pe.IsFullChain = true
		networkFee := big.NewInt(u.networkFeeBySize(int64(maxSize / 4))) // will update after unzip and check data
		pe.MaxFileSize = int64(maxSize)
		pe.NetworkFee = networkFee.String()
	}

	if req.CaptureImageTime == nil {
		cap := entity.DEFAULT_CAPTURE_TIME
		pe.CatureThumbnailDelayTime = &cap
	}

	pe.IsHidden = true
	pe.Status = false
	pe.IsSynced = false
	pe.TxHash = strings.ToLower(pe.TxHash)
	pe.TxHex = strings.ToLower(pe.TxHex)

	pe.TokenID = pe.TxHash
	pe.TokenId = pe.TxHash
	pe.ContractAddress = strings.ToLower(os.Getenv("GENERATIVE_PROJECT"))

	err = u.Repo.CreateProject(pe)
	if err != nil {
		logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Error(err))
		return nil, err
	}
	logger.AtLog.Logger.Error(fmt.Sprintf("CreateProject.%s", pe.TokenId), zap.Any("project", pe))
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
