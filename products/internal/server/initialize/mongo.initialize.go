package initialize

import (
	"context"
	"time"

	"github.com/baothaihcmut/Ecommerce-Go/libs/pkg/logger"
	"github.com/baothaihcmut/Ecommerce-Go/products/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeMongo(cfg *config.Config, logger logger.ILogger) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.Mongo.URI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	//ping for mongo
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}
