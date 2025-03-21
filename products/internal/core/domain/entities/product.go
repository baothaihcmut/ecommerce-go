package entities

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID   `bson:"_id"`
	Name        string               `bson:"name"`
	Description string               `bson:"description"`
	Images      []string             `bson:"images"`
	Thumbnail   *string              `bson:"thumbnail"`
	CategoryIds []primitive.ObjectID `bson:"category_ids"`
	ShopId      primitive.ObjectID   `bson:"shop_id"`
	Variations  []string             `bson:"variations"`
	SoldTotal   int                  `bson:"sold_total"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}

func NewProduct(
	name, description string,
	shopId primitive.ObjectID,
	categoryIds []primitive.ObjectID,
	variations []string,
	hasThumbnail bool,
	numsOfImage int,
) *Product {
	images := make([]string, 0, numsOfImage)
	for i := 0; i < numsOfImage; i++ {
		images = append(images, uuid.NewString())
	}
	var thumbnail *string
	if hasThumbnail {
		thumbnailKey := uuid.NewString()
		thumbnail = &thumbnailKey
	}
	return &Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Images:      images,
		Thumbnail:   thumbnail,
		CategoryIds: categoryIds,
		Variations:  variations,
		SoldTotal:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
