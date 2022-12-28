package usecase

import (
	"bufio"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
	"rederinghub.io/utils/googlecloud"
	"rederinghub.io/utils/helpers"
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
		FileName: uploaded.Name,
		FileSize: int(uploaded.Size),
		MineType: uploaded.Minetype,
		URL:      cdnURL,
	}

	err = u.Repo.InsertOne(fileModel.TableName(), fileModel)
	if err != nil {
		log.Error("u.Repo.InsertOne", err.Error(), err)
		return nil, err
	}

	log.SetData("inserted.FileModel", fileModel)
	return fileModel, nil
}

func (u Usecase) Deflate(inflated []byte) []byte {
	var b bytes.Buffer
	w, _ := flate.NewWriter(&b, flate.DefaultCompression)
	w.Write(inflated)
	w.Close()
	return b.Bytes()
}

func (u Usecase) MinifyFiles(rootSpan opentracing.Span, input structure.MinifyDataResp) (*structure.MinifyDataResp, error) {
	span, log := u.StartSpan("MinifyFiles", rootSpan)
	defer u.Tracer.FinishSpan(span, log)
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

	for fileName, fileInfo := range input.Files {
		bytes, err := helpers.Base64Decode(fileInfo.Content)
		if err != nil {
			log.Error("helpers.Base64Decode.fileInfo.Content", err.Error(), err)
			return nil, err
		}

		out, err := m.String(fileInfo.MediaType, string(bytes))
		if err != nil {
			log.Error("m.String", err.Error(), err)
			return nil, err
		}
		deflate := u.Deflate([]byte(out))
		resp[fileName] = structure.FileContentReq{MediaType: fileInfo.MediaType, Content: out, Deflate: helpers.Base64Eecode(deflate)}
	}

	return &structure.MinifyDataResp{Files: resp}, nil
}
