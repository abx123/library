package repo

import (
	"context"

	"github.com/abx123/library/entities"
)

// IdcRepo defines a dbRepo interface
type IdbRepo interface {
	Get(context.Context, *entities.Book) (*entities.Book, error)
	Upsert(context.Context, *entities.Book) (*entities.Book, error)
	List(context.Context, int64, int64, string) ([]*entities.Book, error)
}
