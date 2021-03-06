package services

import (
	"context"
	"net/http"

	"github.com/abx123/library/entities"
)

// IdbService defines the interface for dbService
type IdbService interface {
	Upsert(context.Context, string, string, string, string, string, string, string, string, string, string, string, int64, int64, int64) (*entities.Book, error)
	Get(context.Context, string, string) (*entities.Book, error)
	List(context.Context, int64, int64, string) ([]*entities.Book, error)
}

// Ibooks defines the interface for bookService
type Ibooks interface {
	Get(context.Context, string) (*entities.Book, error)
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
