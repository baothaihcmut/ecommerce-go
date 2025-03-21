package storage

import (
	"context"
	"time"
)

type GetPresigUrlMethod int

const (
	GetPresignUrlMethodGet GetPresigUrlMethod = iota
	GetPresignUrlMethodPut
)

type GetPresignUrlArg struct {
	Key    string
	Method GetPresigUrlMethod
}

type StorageService interface {
	GetPresignUrl(context.Context, GetPresignUrlArg) (string, error)
	WithExpiry(time.Duration) StorageService
}
