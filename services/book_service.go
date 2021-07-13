package services

import (
	"context"

	"library/entities"
)

type BookService struct {
	gapi Igapi
	gr   Igoodread
}

func NewBookService(gapi *GapiService, gr *GoodreadService) *BookService {
	return &BookService{
		gapi: gapi,
		gr:   gr,
	}
}

func (svc *BookService) Get(ctx context.Context, isbn string) ([]*entities.Book, error) {
	gapi, err := svc.gapi.GetFromGoogleAPI(ctx, isbn)

	if err != nil {
		return nil, err
	}

	gr, err := svc.gr.GetFromGoodread(ctx, isbn)

	if err != nil {
		return nil, err
	}

	res := []*entities.Book{gapi}

	if gapi.Title != gr.Title {
		res = append(res, gr)
	}

	return res, nil
}
