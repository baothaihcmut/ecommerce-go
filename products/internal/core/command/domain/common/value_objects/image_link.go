package valueobjects

import "time"

type ImageLink struct {
	Bucket string
	Key    string
}

func NewImageLink(bucket string, key string) ImageLink {
	return ImageLink{
		Bucket: bucket,
		Key:    time.Now().String() + key,
	}
}
