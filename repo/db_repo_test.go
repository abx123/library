package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/abx123/library/constant"
	"github.com/abx123/library/entities"
)

func NewMockDb() (*sqlx.DB, sqlxmock.Sqlmock) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGet(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `books` WHERE isbn = ? AND userId = ?")
	rows := sqlxmock.NewRows([]string{"id", "isbn", "title", "authors", "imageUrl", "smallImageUrl", "publicationYear", "publisher", "userId", "status", "description", "pageCount", "categories", "language", "source"}).AddRow(1, "9780751562774", "The Secrets She Keeps", "Michael Robotham", "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png", "", 0, "BB Publishing House", "8BeqLfieIiTOkruBBrQ6p8jOTsk2", 1, "", 0, "", "", "goodreads")
	type testCase struct {
		name   string
		desc   string
		err    error
		expRes *entities.Book
		expErr error
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				BookID:    1,
				ISBN:      "9780751562774",
				Title:     "The Secrets She Keeps",
				Authors:   "Michael Robotham",
				ImageURL:  "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
				Publisher: "BB Publishing House",
				UserID:    "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
				Status:    1,
				Source:    "goodreads",
			},
		},
		{
			name:   "Sad Case",
			desc:   "sql returns empty row",
			err:    sql.ErrNoRows,
			expErr: constant.ErrBookNotFound,
		},
		{
			name:   "Sad Case",
			desc:   "sql returns error",
			err:    fmt.Errorf("mock error"),
			expErr: constant.ErrDBErr,
		},
	}

	for _, v := range testCases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		}
		mock.ExpectQuery(query).WillReturnRows(rows)

		actRes, actErr := repo.Get(context.Background(), &entities.Book{ISBN: "9780751562774", UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2"})
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestInsert(t *testing.T) {
	query := regexp.QuoteMeta("INSERT INTO `books` (isbn, title, authors, imageUrl, smallImageUrl, publicationYear, publisher, userId, status, description, pageCount, categories, language, source) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	type testCase struct {
		name         string
		desc         string
		err          error
		dbErr        bool
		expRes       *entities.Book
		lastInertErr bool
		expErr       error
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				BookID: 99,
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
		},
		{
			name:         "Sad Case",
			desc:         "LastInsertId return error",
			err:          fmt.Errorf("lasatInsertId error"),
			expErr:       constant.ErrDBErr,
			lastInertErr: true,
		},
		{
			name:   "Sad Case",
			desc:   "db returns error",
			err:    fmt.Errorf("mock error"),
			expErr: constant.ErrDBErr,
			dbErr:  true,
		},
	}
	for _, v := range testCases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.dbErr {
			mock.ExpectExec(query).WillReturnError(v.err)
		}
		if v.lastInertErr {
			mock.ExpectExec(query).WillReturnResult(sqlxmock.NewErrorResult(v.err))
		}
		mock.ExpectExec(query).WillReturnResult(sqlxmock.NewResult(99, 1))

		actRes, actErr := repo.insert(context.Background(), &entities.Book{ISBN: "9780751562774", UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2"})
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestUpdate(t *testing.T) {
	query := regexp.QuoteMeta("UPDATE `books` SET isbn=?, title=?, authors=?, imageUrl=?, smallImageUrl=?, publicationYear=?, userId=?, status=?, description=?, pageCount=?, categories=?, language=?, source=? WHERE isbn = ? AND userId = ?")
	type testCase struct {
		name            string
		desc            string
		err             error
		expRes          *entities.Book
		expErr          error
		dbErr           bool
		rowsAffectedErr bool
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
		},
		{
			name:   "Sad Case",
			desc:   "db returns error",
			err:    fmt.Errorf("db returns error"),
			expErr: constant.ErrDBErr,
			dbErr:  true,
		},
		{
			name:            "Sad Case",
			desc:            "rows affected returns error",
			err:             fmt.Errorf("db returns error"),
			expErr:          constant.ErrDBErr,
			rowsAffectedErr: true,
		},
	}
	for _, v := range testCases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.dbErr {
			mock.ExpectExec(query).WillReturnError(v.err)
		}
		if v.rowsAffectedErr {
			mock.ExpectExec(query).WillReturnResult(sqlxmock.NewErrorResult(v.err))
		}
		mock.ExpectExec(query).WillReturnResult(sqlxmock.NewResult(1, 1))
		actRes, actErr := repo.update(context.Background(), &entities.Book{ISBN: "9780751562774", UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2"})
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestList(t *testing.T) {
	query := regexp.QuoteMeta("SELECT * FROM `books` WHERE userId=?  LIMIT ? OFFSET ?")
	row := sqlxmock.NewRows([]string{"id", "isbn", "title", "authors", "imageUrl", "smallImageUrl", "publicationYear", "publisher", "userId", "status", "description", "pageCount", "categories", "language", "source"}).AddRow(1, "9780751562774", "The Secrets She Keeps", "Michael Robotham", "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png", "", 0, "BB Publishing House", "8BeqLfieIiTOkruBBrQ6p8jOTsk2", 1, "", 0, "", "", "goodreads").AddRow(2, "9781407243207", "The Bourne Ultimatum", "", "https://images.isbndb.com/covers/32/07/9781407243207.jpg", "", 0, "BB Publishing House", "8BeqLfieIiTOkruBBrQ6p8jOTsk2", 1, "", 0, "", "en_US", "isbndb")
	type testCase struct {
		name   string
		desc   string
		expRes []*entities.Book
		expErr error
		err    error
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: []*entities.Book{
				{
					BookID:    1,
					ISBN:      "9780751562774",
					Title:     "The Secrets She Keeps",
					Authors:   "Michael Robotham",
					ImageURL:  "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png",
					Publisher: "BB Publishing House",
					UserID:    "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
					Status:    1,
					Source:    "goodreads",
				},
				{
					BookID:    2,
					ISBN:      "9781407243207",
					Title:     "The Bourne Ultimatum",
					ImageURL:  "https://images.isbndb.com/covers/32/07/9781407243207.jpg",
					Publisher: "BB Publishing House",
					UserID:    "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
					Status:    1,
					Language:  "en_US",
					Source:    "isbndb",
				},
			},
		},
		{
			name:   "Sad Case",
			desc:   "db returns empty row",
			err:    sql.ErrNoRows,
			expErr: constant.ErrBookNotFound,
		},
		{
			name:   "Sad Case",
			desc:   "db returns error",
			err:    fmt.Errorf("mock error"),
			expErr: constant.ErrDBErr,
		},
	}
	for _, v := range testCases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)
		if v.err != nil {
			mock.ExpectQuery(query).WillReturnError(v.err)
		}
		mock.ExpectQuery(query).WillReturnRows(row)

		actRes, actErr := repo.List(context.Background(), 10, 0, "8BeqLfieIiTOkruBBrQ6p8jOTsk2")
		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}

func TestUpsert(t *testing.T) {
	getQuery := regexp.QuoteMeta("SELECT * FROM `books` WHERE isbn = ? AND userId = ?")
	getRows := sqlxmock.NewRows([]string{"id", "isbn", "title", "authors", "imageUrl", "smallImageUrl", "publicationYear", "publisher", "userId", "status", "description", "pageCount", "categories", "language", "source"}).AddRow(1, "9780751562774", "The Secrets She Keeps", "Michael Robotham", "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png", "", 0, "BB Publishing House", "8BeqLfieIiTOkruBBrQ6p8jOTsk2", 1, "", 0, "", "", "goodreads")
	updateQuery := regexp.QuoteMeta("UPDATE `books` SET isbn=?, title=?, authors=?, imageUrl=?, smallImageUrl=?, publicationYear=?, userId=?, status=?, description=?, pageCount=?, categories=?, language=?, source=? WHERE isbn = ? AND userId = ?")
	insertQuery := regexp.QuoteMeta("INSERT INTO `books` (isbn, title, authors, imageUrl, smallImageUrl, publicationYear, publisher, userId, status, description, pageCount, categories, language, source) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	type testCase struct {
		name      string
		desc      string
		err       error
		expRes    *entities.Book
		expErr    error
		updateErr bool
		getErr    bool
		insertErr bool
		insert    bool
		update    bool
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "update",
			expRes: &entities.Book{
				BookID: 1,
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
			update: true,
		},
		{
			name: "Happy Case",
			desc: "insert",
			err:  sql.ErrNoRows,
			expRes: &entities.Book{
				BookID: 99,
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
			getErr: true,
			insert: true,
		},
		{
			name:      "Sad Case",
			desc:      "insert returns error",
			err:       sql.ErrNoRows,
			getErr:    true,
			insert:    true,
			insertErr: true,
			expErr:    constant.ErrDBErr,
		},
		{
			name:   "Sad Case",
			desc:   "get returns error",
			err:    fmt.Errorf("mock get error"),
			getErr: true,
			expErr: constant.ErrDBErr,
		},
		{
			name:      "Sad Case",
			desc:      "Update returns error",
			err:       fmt.Errorf("update returns error"),
			updateErr: true,
			expErr:    constant.ErrDBErr,
			update:    true,
		},
	}
	for _, v := range testCases {
		db, mock := NewMockDb()
		repo := NewDbRepo(db)

		if v.getErr {
			mock.ExpectQuery(getQuery).WillReturnError(v.err)
		} else {
			mock.ExpectQuery(getQuery).WillReturnRows(getRows)
		}
		if v.update {
			if v.updateErr {
				mock.ExpectExec(updateQuery).WillReturnError(v.err)
			} else {
				mock.ExpectExec(updateQuery).WillReturnResult(sqlxmock.NewResult(1, 1))
			}
		}

		if v.insert {
			if v.insertErr {
				mock.ExpectExec(insertQuery).WillReturnError(fmt.Errorf("mock insert error"))
			} else {
				mock.ExpectExec(insertQuery).WillReturnResult(sqlxmock.NewResult(99, 1))
			}
		}

		actRes, actErr := repo.Upsert(context.Background(), &entities.Book{ISBN: "9780751562774", UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2"})

		assert.Equal(t, v.expRes, actRes)
		assert.Equal(t, v.expErr, actErr)
	}
}
