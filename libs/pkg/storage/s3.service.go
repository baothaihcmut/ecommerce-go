package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client *s3.Client
	cfg    *S3Config
	expiry time.Duration
}

func NewS3Service(s3Client *s3.Client, s3Config *S3Config) *S3Service {
	return &S3Service{
		client: s3Client,
		cfg:    s3Config,
		expiry: time.Minute * 3,
	}
}
func (s *S3Service) WithExpiry(expiry time.Duration) StorageService {
	s.expiry = expiry
	return s
}
func (s *S3Service) GetPresignUrl(ctx context.Context, args GetPresignUrlArg) (string, error) {
	presignner := s3.NewPresignClient(s.client)
	var res *v4.PresignedHTTPRequest
	var err error
	switch args.Method {
	case GetPresignUrlMethodGet:
		res, err = presignner.PresignGetObject(ctx, &s3.GetObjectInput{
			Key:    aws.String(args.Key),
			Bucket: aws.String(s.cfg.Bucket),
		})

	case GetPresignUrlMethodPut:
		res, err = presignner.PresignPutObject(ctx, &s3.PutObjectInput{
			Key:    aws.String(args.Key),
			Bucket: aws.String(s.cfg.Bucket),
		})
	}
	if err != nil {
		return "", err
	}
	return res.URL, err
}
