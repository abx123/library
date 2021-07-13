package services

import (
	"context"

	"library/entities"
	"library/repo"
)

type DBService struct {
	repo repo.IdbRepo
}

func NewDbService(r *repo.DBRepo) *DBService {
	return &DBService{
		repo: r,
	}
}

func (svc *DBService) Upsert(ctx context.Context, isbn, title, author, imageURL, smallImageURL, userId string, averageRating float64, publicationYear, status int64) (*entities.Book, error) {
	book := &entities.Book{
		ISBN:            isbn,
		Title:           title,
		Author:          author,
		ImageURL:        imageURL,
		SmallImageURL:   smallImageURL,
		PublicationYear: publicationYear,
		AverageRating:   averageRating,
		UserID:          userId,
		Status:          status,
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
