package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetPresignLinkForPut(ctx context.Context, s3Client *s3.Client, bucket string, key string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s3Client)
	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func GetPresignLinkForGet(ctx context.Context, s3Client *s3.Client, bucket string, key string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(s3Client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", nil
	}
	return req.URL, nil
}
