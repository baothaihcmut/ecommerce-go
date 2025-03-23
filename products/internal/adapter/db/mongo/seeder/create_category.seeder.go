package seeder

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var seederCategories = []interface{}{
	bson.M{"_id": primitive.NewObjectID(), "name": "Electronics"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Clothing"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Home & Kitchen"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Beauty & Personal Care"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Sports & Outdoors"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Books"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Toys & Games"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Automotive"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Health & Wellness"},
	bson.M{"_id": primitive.NewObjectID(), "name": "Office Supplies"},
}

func CreatCategorySeeders(db *mongo.Database) {
	collection := db.Collection("categories")

	// Prepare bulk operations
	operations := make([]mongo.WriteModel, 0, len(seederCategories))
	for _, category := range seederCategories {
		if cat, ok := category.(bson.M); ok {
			filter := bson.M{"name": cat["name"]}
			update := bson.M{"$setOnInsert": bson.M{
				"_id":  cat["_id"],
				"name": cat["name"],
			}}
			operations = append(operations, mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true))
		}

	}

	// Execute bulk write
	if len(operations) > 0 {
		result, err := collection.BulkWrite(context.TODO(), operations)
		if err != nil {
			log.Fatalf("Error performing bulk insert: %v", err)
		}
		log.Printf("Inserted %d new categories, skipped %d existing categories\n",
			result.UpsertedCount, len(seederCategories)-int(result.UpsertedCount))
	} else {
		log.Println("No categories to insert.")
	}
}
