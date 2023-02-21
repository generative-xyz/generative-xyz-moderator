package googlecloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const (
	maxRetries = 3
)

type File interface {
	io.ReadSeeker
}

type CompletedPart struct {
	ETag       *string `json:"etag,omitempty"`
	PartNumber *int64  `json:"part_number,omitempty"`
}

type S3Adapter struct {
	bucketName  string
	s3Client    *s3.S3
	redisClient *redis.Client
}

type S3AdapterConfig struct {
	BucketName string
	Endpoint   string
	Region     string
	AccessKey  string // Hmac access key
	SecretKey  string // Hmac secret key

}

func NewS3Adapter(config S3AdapterConfig, redisClient *redis.Client) S3Adapter {
	fileSession := session.Must(session.NewSession())
	s3Client := s3.New(fileSession, aws.NewConfig().
		WithRegion(config.Region).
		WithEndpoint(config.Endpoint).
		WithCredentials(credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")))
	return S3Adapter{
		bucketName:  config.BucketName,
		s3Client:    s3Client,
		redisClient: redisClient,
	}
}

func (a S3Adapter) SetUploadIDToKey(uploadID string, key string) error {
	redisKey := fmt.Sprintf("uploadId:key:%s", uploadID)
	return a.redisClient.Set(redisKey, key, time.Hour*10).Err()
}

func (a S3Adapter) GetKeyFromUploadID(uploadID string) (key string, err error) {
	redisKey := fmt.Sprintf("uploadId:key:%s", uploadID)
	resp, err := a.redisClient.Get(redisKey).Result()
	if err != nil {
		return
	}
	if resp == "" {
		return "", errors.New(fmt.Sprintf("upload id %s not found", uploadID))

	}
	return resp, nil
}

func (a S3Adapter) AddS3CompletedPart(uploadID string, part CompletedPart) error {
	data, err := json.Marshal(part)
	if err != nil {
		return errors.WithStack(err)
	}

	redisKey := fmt.Sprintf("uploadId:completedPart:%s", uploadID)
	return a.redisClient.RPush(redisKey, string(data)).Err()
}

func (a S3Adapter) GetS3CompletedPart(uploadID string) (parts []*CompletedPart, err error) {
	redisKey := fmt.Sprintf("uploadId:completedPart:%s", uploadID)
	resp, err := a.redisClient.LRange(redisKey, 0, -1).Result()
	if err != nil {
		return
	}

	parts = make([]*CompletedPart, 0, len(resp))
	for _, item := range resp {
		var part CompletedPart
		if err = json.Unmarshal([]byte(item), &part); err != nil {
			return nil, errors.WithStack(err)
		}
		parts = append(parts, &part)

	}

	return
}

func (a S3Adapter) CreateMultiplePartsUpload(ctx context.Context, group string, fileName string) (*string, error) {
	fileName = NormalizeFileName(fileName)
	uploadPath := fmt.Sprintf("multipart/%s/%s", group, fileName)
	var resp, err = a.s3Client.CreateMultipartUploadWithContext(ctx, &s3.CreateMultipartUploadInput{
		Bucket: &a.bucketName,
		Key:    &uploadPath,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err = a.SetUploadIDToKey(*resp.UploadId, *resp.Key); err != nil {
		return nil, errors.WithStack(err)
	}
	return resp.UploadId, nil

}

func (a S3Adapter) CompleteMultipartUpload(ctx context.Context, uploadID string) (*string, error) {
	completedParts, err := a.GetS3CompletedPart(uploadID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	s3CompletedPart := make([]*s3.CompletedPart, 0, len(completedParts))
	for _, item := range completedParts {
		s3CompletedPart = append(s3CompletedPart, &s3.CompletedPart{
			ETag:       item.ETag,
			PartNumber: item.PartNumber,
		})
	}

	key, err := a.GetKeyFromUploadID(uploadID)
	if err != nil {
		return nil, err
	}

	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   &a.bucketName,
		Key:      &key,
		UploadId: &uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: s3CompletedPart,
		},
	}
	resp, err := a.s3Client.CompleteMultipartUploadWithContext(ctx, completeInput)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Successfully uploaded file: %s\n", resp.String())
	return resp.Location, nil
}

func (a S3Adapter) UploadPart(uploadID string, file File, fileSize int64, partNumber int) error {
	tryNum := 1
	key, err := a.GetKeyFromUploadID(uploadID)
	if err != nil {
		return errors.WithStack(err)
	}
	partInput := &s3.UploadPartInput{
		Body:          file,
		Bucket:        &a.bucketName,
		Key:           &key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      &uploadID,
		ContentLength: aws.Int64(fileSize),
	}

	for tryNum <= maxRetries {
		uploadResult, err := a.s3Client.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return errors.WithStack(aerr)
				}
				return errors.WithStack(err)
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			err = a.AddS3CompletedPart(uploadID, CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			})
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	}
	return errors.New("can not upload after reach match retry")
}

func (a S3Adapter) abortMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := svc.AbortMultipartUpload(abortInput)
	return err
}
