package usecase

import (
	"bufio"
	"bytes"
	"compress/flate"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"go.uber.org/zap"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/contracts/generative_project_data"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/helpers"
	"rederinghub.io/utils/logger"
)

type File interface {
	io.ReadSeeker
}

func (u Usecase) CreateMultipartUpload(ctx context.Context, group string, fileName string) (*string, error) {
	//TODO
	group = helpers.GenerateSlug(group)
	group = fmt.Sprintf("%s-%d", group, time.Now().UTC().Nanosecond())

	fileName = helpers.GenerateSlug(fileName)
	uploadID, err := u.S3Adapter.CreateMultiplePartsUpload(ctx, "btc-projects/"+group, fileName)
	return uploadID, err
}

func (u Usecase) UploadPart(ctx context.Context, uploadID string, file File, fileSize int64, partNumber int) error {

	if err := u.S3Adapter.UploadPart(uploadID, file, fileSize, partNumber); err != nil {
		return err
	}
	return nil
}

func (u Usecase) CompleteMultipartUpload(ctx context.Context, uploadID string) (*string, error) {
	data, err := u.S3Adapter.CompleteMultipartUpload(ctx, uploadID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return data, nil
}

func (u Usecase) UploadFile(r *http.Request) (*entity.Files, error) {

	_, handler, err := r.FormFile("file")
	if err != nil {
		logger.AtLog.Logger.Error("r.FormFile.File", zap.Error(err))
		return nil, err
	}

	gf := googlecloud.GcsFile{
		FileHeader: handler,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Logger.Error("u.GCS.FileUploadToBucke", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("uploaded", zap.Any("uploaded", uploaded))

	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.InsertOne", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("inserted.FileModel", zap.Any("fileModel", fileModel))
	return fileModel, nil
}

func (u Usecase) Deflate(inflated []byte) []byte {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, flate.DefaultCompression)
	w.Write(inflated)
	w.Close()
	return b.Bytes()
}

func (u Usecase) MinifyFiles(input structure.MinifyDataResp) (*structure.MinifyDataResp, error) {

	resp := make(map[string]structure.FileContentReq)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	m.AddFunc("text/plain", func(m *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
		// remove newlines and spaces
		rb := bufio.NewReader(r)
		for {
			line, err := rb.ReadString('\n')
			if err != nil && err != io.EOF {
				return err
			}
			if _, errws := io.WriteString(w, strings.Replace(line, " ", "", -1)); errws != nil {
				return errws
			}
			if err == io.EOF {
				break
			}
		}
		return nil
	})

	client, err := helpers.TCDialer()
	if err != nil {
		logger.AtLog.Logger.Error("ethclient.Dial", zap.Error(err))
		return nil, err
	}

	addr := common.HexToAddress(os.Getenv("GENERATIVE_PROJECT_DATA"))
	gDataNft, err := generative_project_data.NewGenerativeProjectData(addr, client)
	if err != nil {
		logger.AtLog.Logger.Error("generative_project_data.NewGenerativeProjectData", zap.Error(err))
		return nil, err
	}

	for fileName, fileInfo := range input.Files {
		bytes, err := helpers.Base64Decode(fileInfo.Content)
		if err != nil {
			logger.AtLog.Logger.Error("helpers.Base64Decode.fileInfo.Content", zap.Error(err))
			return nil, err
		}

		/*out, err := m.String(fileInfo.MediaType, string(bytes))
		if err != nil {
			logger.AtLog.Logger.Error("m.String", zap.Error(err))
			return nil, err
		}*/
		out := string(bytes)
		deflate := u.Deflate([]byte(out))

		script := helpers.Base64Encode(deflate)
		inflate, _ := gDataNft.InflateString(nil, script)
		if inflate.Err != 0 || inflate.Result != out {
			script = ""
		}

		logger.AtLog.Logger.Info("inflate", zap.Any("inflate", inflate))
		resp[fileName] = structure.FileContentReq{MediaType: fileInfo.MediaType, Content: out, Deflate: script}
	}

	return &structure.MinifyDataResp{Files: resp}, nil
}

func (u Usecase) DeflateString(input *structure.DeflateDataResp) error {

	//TODO implement here

	client, err := helpers.TCDialer()
	if err != nil {
		logger.AtLog.Logger.Error("ethclient.Dial", zap.Error(err))
		return err
	}

	addr := common.HexToAddress(os.Getenv("GENERATIVE_PROJECT_DATA"))
	gDataNft, err := generative_project_data.NewGenerativeProjectData(addr, client)
	if err != nil {
		logger.AtLog.Logger.Error("generative_project_data.NewGenerativeProjectData", zap.Error(err))
		return err
	}
	inputByte := []byte(input.Data)
	deflate := u.Deflate(inputByte)
	script := helpers.Base64Encode(deflate)
	logger.AtLog.Logger.Info("len(deflate)", zap.Any("len(deflate)", len(deflate)))
	logger.AtLog.Logger.Info("len(inputByte)", zap.Any("len(inputByte)", len(inputByte)))
	if len(deflate) > len(inputByte) {
		input.Data = ""
		return nil
	}
	inflate, _ := gDataNft.InflateString(nil, script)
	if inflate.Err != 0 || inflate.Result != input.Data {
		logger.AtLog.Logger.Info("inflate.Err", zap.Any("inflate.Err", inflate.Err))
		input.Data = ""
		return nil
	}
	logger.AtLog.Logger.Info("inflate", zap.Any("inflate", inflate))
	input.Data = script
	return nil
}

func (u Usecase) UploadProjectFiles(r *http.Request) (*entity.Files, error) {

	projectName := r.FormValue("projectName")
	_, handler, err := r.FormFile("file")
	if err != nil {
		logger.AtLog.Error("UploadProjectFiles", zap.String("action", "FormFile"), zap.String("err", err.Error()))
		return nil, err
	}

	if handler.Size <= 0 {
		err := errors.New("The uploaded file is empty")
		logger.AtLog.Error("UploadProjectFiles", zap.String("action", "Checkfilesize"), zap.String("err", err.Error()))
		return nil, err
	}

	key := helpers.GenerateSlug(projectName)
	key = fmt.Sprintf("btc-projects/%s", key)
	gf := googlecloud.GcsFile{
		FileHeader: handler,
		Path:       &key,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Error("UploadProjectFiles", zap.String("action", "FileUploadToBucket"), zap.String("err", err.Error()))
		return nil, err
	}

	logger.AtLog.Info("uploaded", uploaded)
	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		logger.AtLog.Error("UploadProjectFiles", zap.String("action", "InsertOne"), zap.String("err", err.Error()))
		return nil, err
	}

	logger.AtLog.Info("UploadProjectFiles", zap.Any("fileModel", fileModel))
	return fileModel, nil
}

func (u Usecase) UploadDatasetFile(r *http.Request, requestID string, datasetName string) (*entity.Files, error) {

	_, handler, err := r.FormFile("file")
	if err != nil {
		logger.AtLog.Logger.Error("r.FormFile.File", zap.Error(err))
		return nil, err
	}

	key := fmt.Sprintf("ai-school/%s", requestID)
	handler.Filename = datasetName + ".zip"
	gf := googlecloud.GcsFile{
		FileHeader: handler,
		Path:       &key,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Logger.Error("u.GCS.FileUploadToBucke", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("uploaded", zap.Any("uploaded", uploaded))

	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.InsertOne", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("inserted.FileModel", zap.Any("fileModel", fileModel))
	return fileModel, nil
}

func (u Usecase) UploadFileInternal(r *http.Request, path string) (*entity.Files, error) {

	_, handler, err := r.FormFile("file")
	if err != nil {
		logger.AtLog.Logger.Error("r.FormFile.File", zap.Error(err))
		return nil, err
	}

	gf := googlecloud.GcsFile{
		FileHeader: handler,
		Path:       &path,
	}

	uploaded, err := u.GCS.FileUploadToBucket(gf)
	if err != nil {
		logger.AtLog.Logger.Error("u.GCS.FileUploadToBucke", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("uploaded", zap.Any("uploaded", uploaded))

	cdnURL := fmt.Sprintf("%s/%s", os.Getenv("GCS_DOMAIN"), uploaded.Name)
	fileModel := &entity.Files{
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		logger.AtLog.Logger.Error("u.Repo.InsertOne", zap.Error(err))
		return nil, err
	}

	logger.AtLog.Logger.Info("inserted.FileModel", zap.Any("fileModel", fileModel))
	return fileModel, nil
}

func (u Usecase) DeleteFile(uuid string) error {
	file, err := u.Repo.GetFileByUUID(uuid)
	if err != nil {
		logger.AtLog.Logger.Error("DeleteFile", zap.Error(err))
		return err
	}
	_, err = u.Repo.SoftDelete(file)
	if err != nil {
		logger.AtLog.Logger.Error("DeleteFile", zap.Error(err))
		return err
	}
	return nil
}
