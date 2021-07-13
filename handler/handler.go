package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"library/constant"
	"library/handler/presenter"
	"library/repo"
	"library/services"
)

type postUpsertBookRequest struct {
	BookID          int64   `json:"id" form:"id"`
	Title           string  `json:"title" form:"title"`
	Author          string  `json:"author" form:"author"`
	ImageURL        string  `json:"imageUrl" form:"imageUrl"`
	SmallImageURL   string  `json:"smallImageUrl" form:"smallImageUrl"`
	PublicationYear int64   `json:"publicationYear" form:"publicationYear"`
	AverageRating   float64 `json:"averageRating" form:"averageRating"`
	Status          int64   `json:"status" form:"status"`
}

type Handler struct {
	bookSvc services.Ibooks
	gapiSvc services.Igapi
	grSvc   services.Igoodread
	dbSvc   services.IdbService
}

func NewHandler(conn *sqlx.DB) *Handler {
	dbRepo := repo.NewDbRepo(conn)
	dbSvc := services.NewDbService(dbRepo)
	gapiSvc := services.NewGapiService()
	grSvc := services.NewGoodreadService()
	bookSvc := services.NewBookService(gapiSvc, grSvc)

	return &Handler{
		dbSvc:   dbSvc,
		gapiSvc: gapiSvc,
		grSvc:   grSvc,
		bookSvc: bookSvc,
	}
}

func (h *Handler) GetBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	isbn := c.Param("isbn")
	userId := c.Param("userId")

	data, err := h.dbSvc.Get(ctx, isbn, userId)
	if err != nil {
		zap.L().Error(err.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	// Return ok
	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:            isbn,
		Title:           data.Title,
		Author:          data.Author,
		ImageURL:        data.ImageURL,
		SmallImageURL:   data.SmallImageURL,
		PublicationYear: data.PublicationYear,
		AverageRating:   data.AverageRating,
		UserID:          data.UserID,
	})
}

// GetBook handles GET /book/:isbn"
func (h *Handler) GetNewBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	isbn := c.Param("isbn")
	if _, err := strconv.ParseInt(isbn, 10, 64); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	data, err := h.bookSvc.Get(ctx, isbn)
	if err != nil {
		zap.L().Error(err.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	// Return ok
	books := []*presenter.Book{}
	for _, d := range data {
		books = append(books, &presenter.Book{
			ISBN:            d.ISBN,
			Title:           d.Title,
			Author:          d.Author,
			ImageURL:        d.ImageURL,
			SmallImageURL:   d.SmallImageURL,
			PublicationYear: d.PublicationYear,
			AverageRating:   d.AverageRating,
			UserID:          d.UserID,
			Status:          d.Status,
		})
	}
	return c.JSON(http.StatusOK, books)
}

func (h *Handler) ListBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	limit, offset, err := getLimitAndOffest(c)
	userId := c.Param("userid")
	if err != nil {
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	data, err := h.dbSvc.List(ctx, limit, offset, userId)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	books := []*presenter.Book{}
	for _, d := range data {
		books = append(books, &presenter.Book{
			ISBN:            d.ISBN,
			Title:           d.Title,
			Author:          d.Author,
			ImageURL:        d.ImageURL,
			SmallImageURL:   d.SmallImageURL,
			PublicationYear: d.PublicationYear,
			AverageRating:   d.AverageRating,
			UserID:          d.UserID,
			Status:          d.Status,
		})
	}

	return c.JSON(http.StatusOK, books)
}

func (h *Handler) UpsertBook(c echo.Context) (err error) {
	r := &postUpsertBookRequest{}
	userId := c.Param("userid")
	isbn := c.Param("isbn")
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	if _, err := strconv.ParseInt(isbn, 10, 64); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, constant.ErrInvalidRequest))
	}

	book, err := h.dbSvc.Upsert(ctx, isbn, r.Title, r.Author, r.ImageURL, r.SmallImageURL, userId, r.AverageRating, r.PublicationYear, r.Status)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}

	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:            book.ISBN,
		Title:           book.Title,
		Author:          book.Author,
		ImageURL:        book.ImageURL,
		SmallImageURL:   book.SmallImageURL,
		PublicationYear: book.PublicationYear,
		AverageRating:   book.AverageRating,
		UserID:          book.UserID,
		Status:          book.Status,
	})
}

func (con *Handler) Ping(c echo.Context) (err error) {
	// Server is up and running, return OK!
	return c.String(http.StatusOK, "Pong")
}

func getLimitAndOffest(c echo.Context) (int64, int64, error) {
	strlimit := c.QueryParam("limit")
	stroffset := c.QueryParam("offset")
	var limit int64 = 10
	var offset int64
	var err error
	if strlimit != "" {
		limit, err = strconv.ParseInt(strlimit, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	if stroffset != "" {
		offset, err = strconv.ParseInt(stroffset, 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	return limit, offset, nil
}
