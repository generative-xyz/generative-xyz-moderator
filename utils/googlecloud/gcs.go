package googlecloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"rederinghub.io/utils/config"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

const (
	GcloudStorePath string = "https://storage.googleapis.com"
)

type IGcstorage interface {
	FileUploadToBucket(file GcsFile) (*GcsUploadedObject, error)
	ReadFileFromBucket(fileName string) ([]byte, error)
}

type GcsUploadedObject struct {
	Name     string
	FullName string
	Path     string
	Minetype string
	Size     int64
	FullPath string
}

type GcsFile struct {
	FileHeader *multipart.FileHeader
}
type gcstorage struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	projectId  string
	ctx        context.Context
	formatType string
}

func NewDataGCStorage(config config.Config) (*gcstorage, error) {
	// Creates a Google Cloud client from config GC Auth key
	jsonKey, _ := base64.StdEncoding.DecodeString(config.Gcs.Auth)
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(jsonKey)))
	if err != nil {
		return nil, err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(config.Gcs.Bucket)

	// Init our GCStorage object
	gcStorage := gcstorage{
		bucketName: config.Gcs.Bucket,    // get bucket name from config
		bucket:     bucket,               // assign bucket object
		client:     client,               // assign client object
		ctx:        ctx,                  // assign context object
		projectId:  config.Gcs.ProjectId, // assign project id, not required
	}

	return &gcStorage, nil
}

func (g gcstorage) FileUploadToBucket(file GcsFile) (*GcsUploadedObject, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*60)
	defer cancel()

	now := time.Now().Unix()
	fname := strings.ToLower(file.FileHeader.Filename)
	fname = strings.ReplaceAll(fname, " ", "_")
	fname = strings.TrimSpace(fname)

	strLen := len(fname)
	if strLen >= 10 {
		fname = fname[strLen-8 : strLen]
	}

	fname = fmt.Sprintf("%d-%s", now, fname)

	fmt.Printf("Uploaded File: %+v\n", file.FileHeader.Filename)
	fmt.Printf("File Size: %+v\n", file.FileHeader.Size)
	fmt.Printf("MIME Header: %+v\n", file.FileHeader.Header)

	header := file.FileHeader.Header
	_ = header

	contentType := header.Get("Content-Type")

	// create writer
	sw := g.bucket.Object(fname).NewWriter(ctx)
	sw.ContentType = contentType

	openedFile, err := file.FileHeader.Open()
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(sw, openedFile); err != nil {
		return nil, err
	}

	if err := sw.Close(); err != nil {
		return nil, err
	}

	attrs := sw.Attrs()
	u, err := url.Parse(g.bucketName + "/" + attrs.Name)
	if err != nil {
		return nil, err
	}
	filePath := u.EscapedPath()

	result := GcsUploadedObject{
		Name:     attrs.Name,
		Minetype: attrs.ContentType,
		Size:     attrs.Size,
		Path:     filePath,
		FullPath: attrs.MediaLink,
	}
	return &result, nil
}

func (g gcstorage) ReadFileFromBucket(fileName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*60)
	defer cancel()

	// create reader
	r, err := g.bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil

}
