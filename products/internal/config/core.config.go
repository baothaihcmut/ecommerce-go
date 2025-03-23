package config

import (
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/mongo"
	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/storage"
)

type CoreConfig struct {
	Server *ServerConfig        `mapstructure:"server"`
	Mongo  *mongo.MongoConfig   `mapstructure:"mongo"`
	S3     *storage.S3Config    `mapstructure:"s3"`
	Logger *logger.LoggerConfig `mapstructure:"logger"`
}
