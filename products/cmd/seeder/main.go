package main

import (
	"context"
	"log"
	"os"

	"github.com/baothaihcmut/Ecommerce-go/products/internal/adapter/db/mongo/seeder"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var seederFuncs = []func(*mongo.Database){
	seeder.CreatCategorySeeders,
}

func main() {
	if err := godotenv.Load(); err != nil {
		return
	}
	uri := os.Getenv("MONGO_URI")
	database := os.Getenv("MONGO_DB")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}
	for _, seederFunc := range seederFuncs {
		seederFunc(client.Database(database))
	}

}
