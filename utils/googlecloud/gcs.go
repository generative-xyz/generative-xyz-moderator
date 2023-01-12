package googlecloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
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
	UploadBaseToBucket(base64Srting string, name string) (*GcsUploadedObject, error)
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

	fname = fmt.Sprintf("upload/%d-%s", now, fname)

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
	r, err := g.bucket.Object(fmt.Sprintf("upload/%s",fileName)).NewReader(ctx)
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

func (g *gcstorage) UploadBaseToBucket(base64Srting string, name string) (*GcsUploadedObject, error) {
	return g.writer(base64Srting, name)
}

type ImageConfig struct {
	Width int64
	Height int64
	Ratio string
	RatioWidth int
	RatioHeight int
}


type uploadGcsChannel struct {
	Attrs *storage.ObjectAttrs
	Err error
	FilePath string
}

type detectImageSizeChannel struct {
	Size *ImageConfig
	Err error
}

func (g *gcstorage) writer(base64Image string, objectName string) (*GcsUploadedObject, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*60)
	defer cancel()

	gcsChannel := make(chan *uploadGcsChannel, 1)
	detectSizeChannel := make(chan *detectImageSizeChannel, 1)

	//upload to GCS routine
	go func(gcsChannel chan *uploadGcsChannel, base64Image string, objectName string) {

		decode, err := base64.StdEncoding.DecodeString(base64Image)
		// create writer
		sw := g.bucket.Object(objectName).NewWriter(ctx)

		channel := &uploadGcsChannel{}

		defer func ()  {
			gcsChannel <- channel
		}()

		//bytesData := []byte(file.ImageData)

		_, err = sw.Write(decode)
		if err != nil {
			channel.Err = err
			return 
		}

		if err = sw.Close(); err != nil {
			channel.Err = err
			return 
		}

		attrs := sw.Attrs()
		u, err := url.Parse("/" + g.bucketName + "/" + attrs.Name)
		if err != nil {
			channel.Err = err
			return 
		}
		filePath := u.EscapedPath()
		//fullPath := fmt.Sprintf("%s%s", GcloudStorePath, filePath)
		
		channel.Attrs = attrs
		channel.FilePath = filePath

	}(gcsChannel, base64Image, objectName)
	
	
	go func(detectSizeChannel chan *detectImageSizeChannel, base64Image string, objectName string) {
		channel := &detectImageSizeChannel{}
		dec, err := base64.StdEncoding.DecodeString(base64Image)

		defer func ()  {
			detectSizeChannel <- channel
		}()
		
		if err != nil {
			channel.Err = err
			return
		}

		f, err := os.Create(objectName)
		if err != nil {
			channel.Err = err
			return
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			channel.Err = err
			return
		}
		if err := f.Sync(); err != nil {
			channel.Err = err
			return
		}

		//Detect image size & ratio
		size, err := g.detectImageSize(objectName)
		if err != nil {
			channel.Err = err
			return 
		}

		channel.Size = size

		//Delete the redundant files after info has been detected.
		g.deleFile(objectName)

	}(detectSizeChannel, base64Image, objectName)

	uploadedInfo := <- gcsChannel
	if uploadedInfo.Err != nil {
		return nil, uploadedInfo.Err
	}

	attrs := uploadedInfo.Attrs
	filePath := uploadedInfo.FilePath
	

	result := GcsUploadedObject{
		Name:     attrs.Name,
		Minetype: attrs.ContentType,
		Size:     attrs.Size,
		Path:     filePath,
		FullPath: attrs.MediaLink,
		
	}
	return &result, nil
}

func (g *gcstorage) detectImageSize(fileName string)  (*ImageConfig, error) {
	reader, err := os.Open(fileName); 
	if err != nil {
		fmt.Println("Impossible to open the file:", err)
		return nil, err		
	} 

	defer reader.Close()
	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return nil, err	
	}

	detectedRation := g.detectRatio(&im)
	return &detectedRation, nil
}


func (g *gcstorage) detectRatio(size *image.Config) ImageConfig {
	width := size.Width
	height := size.Height
	returnData := ImageConfig{
		Width:  int64(width),
		Height:  int64(height),
		Ratio:   "1:1",
		RatioWidth: 1,
		RatioHeight: 1,
	}

	if width == height {
		 returnData.Ratio = "1:1"
		 return returnData
	}

	number := g.findDeviedNumber(width, height)
	ratioW := width
	ratioH := height
	for {
		if ratioW % number != 0 || ratioH % number != 0 {
			break
		}
		ratioW = ratioW / number
		ratioH = ratioH / number
	}

	returnData.Ratio = fmt.Sprintf("%d:%d",ratioW, ratioH)
	returnData.RatioWidth = ratioW
	returnData.RatioHeight = ratioH
	return returnData
}


func (g *gcstorage) findDeviedNumber(with int, height int) int {
	i := 2
	for {
		if with % i == 0 && height %i == 0 {
			break
		}
		i ++
	}
	return i
} 

func (g *gcstorage) Delete(objectName string) error {
	// [START delete_file]
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	o := g.client.Bucket(g.bucketName).Object(objectName)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

func (g *gcstorage) Read(objectName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(g.ctx, time.Second*60)
	defer cancel()
	rc, err := g.bucket.Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return data, nil
}

func (g *gcstorage) deleFile(tmpFileName string) error {
	// Removing file from the directory
	// Using Remove() function
	e := os.Remove(tmpFileName)
	if e != nil {
		return e
	}
	return nil
}