package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Lib "github.com/baothaihcmut/Ecommerce-Go/libs/pkg/s3"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/core/command/port/outbound/storage"
)

type S3StorageService struct {
	s3Client *s3.Client
}

func NewS3StorageService(s3Client *s3.Client) storage.StorageService {
	return &S3StorageService{
		s3Client: s3Client,
	}
}
func (s *S3StorageService) GetPresignUrl(ctx context.Context, args storage.GetPresignUrlArgs) (string, error) {
	if args.Method == storage.GET {
		return s3Lib.GetPresignLinkForGet(ctx, s.s3Client, args.Link.Bucket, args.Link.Key, time.Hour*3)
	} else {
		return s3Lib.GetPresignLinkForPut(ctx, s.s3Client, args.Link.Bucket, args.Link.Key, time.Hour*3)
	}
}
