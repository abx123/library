package services

import (
	"context"

	"library/entities"
)

type IdbService interface {
	Upsert(context.Context, string, string, string, string, string, string, string, string, string, string, string, int64, int64, int64) (*entities.Book, error)
	Get(context.Context, string, string) (*entities.Book, error)
	List(context.Context, int64, int64, string) ([]*entities.Book, error)
}

type Ibooks interface {
	Get(context.Context, string) (*entities.Book, error)
}
