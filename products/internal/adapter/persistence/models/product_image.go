package models

type ProductImage struct {
	Url             string `bson:"url"`
	StorageProvider string `bson:"storage_provider"`
	Size            int    `bson:"size"`
	Type            string `bson:"type"`
	Width           int    `bson:"width"`
	Height          int    `bson:"height"`
}
