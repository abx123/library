package repo

import (
	"context"

	"library/entities"
)

type IdbRepo interface {
	Get(context.Context, *entities.Book) (*entities.Book, error)
	Upsert(context.Context, *entities.Book) (*entities.Book, error)
	List(context.Context, int64, int64, string) ([]*entities.Book, error)
}
