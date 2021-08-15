package services

import (
	"context"

	"github.com/abx123/library/entities"
	"github.com/abx123/library/repo"
)

type DBService struct {
	repo repo.IdbRepo
}

func NewDbService(r repo.IdbRepo) *DBService {
	return &DBService{
		repo: r,
	}
}

func (svc *DBService) Upsert(ctx context.Context, isbn, title, authors, imageURL, smallImageURL, publisher, userId, description, categories, language, source string, publicationYear, status, pageCount int64) (*entities.Book, error) {
	book := &entities.Book{
		ISBN:            isbn,
		Title:           title,
		Authors:         authors,
		ImageURL:        imageURL,
		SmallImageURL:   smallImageURL,
		PublicationYear: publicationYear,
		Publisher:       publisher,
		UserID:          userId,
		Status:          status,
		Description:     description,
		PageCount:       pageCount,
		Categories:      categories,
		Language:        language,
		Source:          source,
	}
	book, err := svc.repo.Upsert(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (svc *DBService) Get(ctx context.Context, isbn string, userId string) (*entities.Book, error) {
	book := &entities.Book{ISBN: isbn, UserID: userId}
	book, err := svc.repo.Get(ctx, book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (svc *DBService) List(ctx context.Context, limit, offset int64, userId string) ([]*entities.Book, error) {
	books, err := svc.repo.List(ctx, limit, offset, userId)
	if err != nil {
		return nil, err
	}
	return books, err
}
