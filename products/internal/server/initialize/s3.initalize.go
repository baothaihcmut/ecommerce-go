package initialize

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/baothaihcmut/Ecommerce-go/products/internal/config"
)

func InitalizeS3(cfg *config.Config) (*s3.Client, error) {
	s3Cfg, err := awsConfig.LoadDefaultConfig(
		context.Background(),
		awsConfig.WithRegion(cfg.S3.Region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(s3Cfg), nil
}
