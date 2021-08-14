package handler

import (
	"context"
	"fmt"
	"library/constant"
	"library/entities"
	"library/services/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetBook(t *testing.T) {
	type testCase struct {
		name     string
		desc     string
		url      string
		err      error
		expRes   *entities.Book
		httpCode int
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: &entities.Book{
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/book/9780751562774",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad Case",
			desc:     "svc return error",
			err:      fmt.Errorf("mock error"),
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/book/9780751562774",
			httpCode: http.StatusInternalServerError,
		},
	}

	for _, v := range testCases {

		dbSvc := mocks.IdbService{}
		bSvc := mocks.Ibooks{}
		h := NewHandler(&dbSvc, &bSvc)
		dbSvc.On("Get", context.WithValue(context.Background(), constant.ContextKeyRequestID, ""), "9780751562774", "8BeqLfieIiTOkruBBrQ6p8jOTsk2").Return(v.expRes, v.err)
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/:userId/book/:isbn", h.GetBook)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestGetNewBook(t *testing.T) {
	type testCase struct {
		name     string
		desc     string
		url      string
		err      error
		expRes   *entities.Book
		httpCode int
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			url:  "http://localhost:1323/book/9780751562774",
			expRes: &entities.Book{
				ISBN:   "9780751562774",
				UserID: "8BeqLfieIiTOkruBBrQ6p8jOTsk2",
			},
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad Case",
			desc:     "invalid request param",
			url:      "http://localhost:1323/book/97807dsa51562774",
			httpCode: http.StatusBadRequest,
		},
		{
			name:     "Sad Case",
			desc:     "svc return error",
			url:      "http://localhost:1323/book/9780751562774",
			err:      fmt.Errorf("mock error"),
			httpCode: http.StatusInternalServerError,
		},
	}
	for _, v := range testCases {
		dbSvc := mocks.IdbService{}
		bSvc := mocks.Ibooks{}
		h := NewHandler(&dbSvc, &bSvc)
		bSvc.On("Get", context.WithValue(context.Background(), constant.ContextKeyRequestID, ""), "9780751562774").Return(v.expRes, v.err)
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/book/:isbn", h.GetNewBook)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestListBook(t *testing.T) {
	type testCase struct {
		name     string
		desc     string
		url      string
		err      error
		expRes   []*entities.Book
		httpCode int
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			expRes: []*entities.Book{
				{
					ISBN: "9780751562774",
				},
				{
					ISBN: "9780751562774",
				},
			},
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/books?limit=10&offset=0",
			httpCode: http.StatusOK,
		},
		{
			name:     "Sad Case",
			desc:     "invalid request param",
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/books?limit=1d0&offset=0",
			httpCode: http.StatusBadRequest,
		},
		{
			name:     "Sad Case",
			desc:     "invalid request param",
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/books?limit=10&offset=a0",
			httpCode: http.StatusBadRequest,
		},
		{
			name:     "Sad Case",
			desc:     "svc return err",
			err:      fmt.Errorf("mock error"),
			url:      "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/books?limit=10&offset=0",
			httpCode: http.StatusInternalServerError,
		},
	}
	for _, v := range testCases {
		dbSvc := mocks.IdbService{}
		bSvc := mocks.Ibooks{}
		h := NewHandler(&dbSvc, &bSvc)
		dbSvc.On("List", context.WithValue(context.Background(), constant.ContextKeyRequestID, ""), int64(10), int64(0), "8BeqLfieIiTOkruBBrQ6p8jOTsk2").Return(v.expRes, v.err)
		req := httptest.NewRequest(http.MethodGet, v.url, nil)
		w := httptest.NewRecorder()
		r := echo.New()
		r.GET("/:userId/books", h.ListBook)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestUpsert(t *testing.T) {
	type testCase struct {
		name     string
		desc     string
		err      error
		expRes   *entities.Book
		httpCode int
		form     url.Values
	}
	testCases := []testCase{
		{
			name: "Happy Case",
			desc: "all ok",
			form: map[string][]string{
				"isbn":     {"9780751562774"},
				"title":    {"The Secrets She Keeps"},
				"author":   {"Michael Robotham"},
				"imageURL": {"https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png"},
				"status":   {"1"},
				"source":   {"goodreads"},
			},
			httpCode: http.StatusOK,
			expRes: &entities.Book{
				ISBN: "isbn",
			},
		},
		{
			name: "Sad Case",
			desc: "invalid request param",
			form: map[string][]string{
				"isbn":     {"9780751562774"},
				"title":    {"The Secrets She Keeps"},
				"author":   {"Michael Robotham"},
				"imageURL": {"https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png"},
				"status":   {"1dsa"},
				"source":   {"goodreads"},
			},
			httpCode: http.StatusBadRequest,
			expRes: &entities.Book{
				ISBN: "isbn",
			},
		},
		{
			name: "Sad Case",
			desc: "invalid request param",
			form: map[string][]string{
				"isbn":     {"978ddsa0751562774"},
				"title":    {"The Secrets She Keeps"},
				"author":   {"Michael Robotham"},
				"imageURL": {"https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png"},
				"status":   {"1"},
				"source":   {"goodreads"},
			},
			httpCode: http.StatusBadRequest,
			expRes: &entities.Book{
				ISBN: "isbn",
			},
		},
		{
			name: "Sad Case",
			desc: "svc return error",
			err:  fmt.Errorf("mock error"),
			form: map[string][]string{
				"isbn":     {"9780751562774"},
				"title":    {"The Secrets She Keeps"},
				"author":   {"Michael Robotham"},
				"imageURL": {"https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png"},
				"status":   {"1"},
				"source":   {"goodreads"},
			},
			httpCode: http.StatusInternalServerError,
		},
	}
	for _, v := range testCases {
		dbSvc := mocks.IdbService{}
		bSvc := mocks.Ibooks{}
		h := NewHandler(&dbSvc, &bSvc)
		dbSvc.On("Upsert", context.WithValue(context.Background(), constant.ContextKeyRequestID, ""), "9780751562774", "The Secrets She Keeps", "Michael Robotham", "https://s.gr-assets.com/assets/nophoto/book/111x148-bcc042a9c91a29c1d680899eff700a03.png", "", "", "8BeqLfieIiTOkruBBrQ6p8jOTsk2", "", "", "", "goodreads", int64(0), int64(1), int64(0)).Return(v.expRes, v.err)
		req := httptest.NewRequest(http.MethodPost, "http://localhost:1323/8BeqLfieIiTOkruBBrQ6p8jOTsk2/book", strings.NewReader(v.form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		r := echo.New()
		r.POST("/:userId/book", h.UpsertBook)
		r.ServeHTTP(w, req)
		assert.Equal(t, v.httpCode, w.Code)
	}
}

func TestPing(t *testing.T) {
	dbSvc := mocks.IdbService{}
	bSvc := mocks.Ibooks{}
	h := NewHandler(&dbSvc, &bSvc)
	req := httptest.NewRequest(http.MethodGet, "http://localhost:1323/ping", nil)
	w := httptest.NewRecorder()
	r := echo.New()
	r.GET("/ping", h.Ping)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
