package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"

	"github.com/abx123/library/constant"
	"github.com/abx123/library/entities"
)

// DBRepo defines a DBRepo object
type DBRepo struct {
	db *sqlx.DB
}

// NewDbRepo creates a new instance of DBRepo object
func NewDbRepo(db *sqlx.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}

// Upsert updates the record if a record is found, inserts a new record if no record is found.
func (r *DBRepo) Upsert(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	b, err := r.Get(ctx, book)
	if err != nil {
		if err == constant.ErrBookNotFound {
			//insert
			b, err = r.insert(ctx, book)
			if err != nil {
				return nil, err
			}
			return b, nil
		}

		return nil, err
	}
	book.BookID = b.BookID
	// update
	b, err = r.update(ctx, book)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Get get searches the database for a record match
func (r *DBRepo) Get(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	b := &entities.Book{}
	err := r.db.Get(b, "SELECT * FROM `books` WHERE isbn = ? AND userId = ?", book.ISBN, book.UserID)
	if err != nil {
		zap.L().Error(constant.ErrDBErr.Error(), zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constant.ErrBookNotFound
		}

		return nil, constant.ErrDBErr
	}
	return b, nil
}

func (r *DBRepo) insert(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	// Execute Statement
	res, err := r.db.Exec("INSERT INTO `books` (isbn, title, authors, imageUrl, smallImageUrl, publicationYear, publisher, userId, status, description, pageCount, categories, language, source) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", book.ISBN, book.Title, book.Authors, book.ImageURL, book.SmallImageURL, book.PublicationYear, book.Publisher, book.UserID, book.Status, book.Description, book.PageCount, book.Categories, book.Language, book.Source)
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
	res, err := r.db.Exec("UPDATE `books` SET isbn=?, title=?, authors=?, imageUrl=?, smallImageUrl=?, publicationYear=?, userId=?, status=?, description=?, pageCount=?, categories=?, language=?, source=?, publisher=? WHERE isbn = ? AND userId = ?", book.ISBN, book.Title, book.Authors, book.ImageURL, book.SmallImageURL, book.PublicationYear, book.UserID, book.Status, book.Description, book.PageCount, book.Categories, book.Language, book.Source, book.Publisher, book.ISBN, book.UserID)
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

// List returns list of records that matches the search criteria
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
