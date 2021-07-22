package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"library/constant"
	"library/entities"
)

type DBRepo struct {
	db *sqlx.DB
}

func NewDbRepo(db *sqlx.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// GetTable returns detail of a single table.
func (r *DBRepo) Upsert(ctx context.Context, book *entities.Book) (*entities.Book, error) {

	b, err := r.Get(ctx, book)
	if err != nil {
		return nil, err
	}

	if b.BookID != 0 {
		// update
		b, err = r.update(ctx, book)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	b, err = r.insert(ctx, book)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *DBRepo) Get(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	b := entities.Book{}
	err := r.db.Get(b, "SELECT * FROM `books` WHERE isbn = ? AND userId = ?", book.ISBN, book.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constant.ErrBookNotFound
		}
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	return &b, nil
}

func (r *DBRepo) insert(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	// Execute Statement
	res, err := r.db.Exec("INSERT INTO `books` (isbn, title, author, imageUrl, userId) VALUES(?, ?, ?, ?, ?, ?, ?, ?)", book.ISBN, book.Title, book.Author, book.ImageURL, book.UserID)
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	id, err := res.LastInsertId()
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		// Error getting ID of newly created record
		return nil, constant.ErrDBErr
	}
	book.BookID = id

	return book, nil
}

func (r *DBRepo) update(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	// Execute Statement
	res, err := r.db.Exec("UPDATE `books` SET isbn=?, title=?, author=?, imageUrl=?, userId=?", book.ISBN, book.Title, book.Author, book.ImageURL, book.UserID)
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		// Error paring statement result into struct
		return nil, constant.ErrDBErr
	}
	_, err = res.RowsAffected()
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		// Error getting ID of newly created record
		return nil, constant.ErrDBErr
	}

	return book, nil
}

func (r *DBRepo) List(ctx context.Context, limit, offset int64, userId string) ([]*entities.Book, error) {
	books := []*entities.Book{}
	err := r.db.Select(&books, "SELECT * FROM `books` WHERE userId=?  LIMIT ? OFFSET ?", userId, limit, offset)
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrBookNotFound
		}
		return nil, constant.ErrDBErr
	}
	return books, nil
}
