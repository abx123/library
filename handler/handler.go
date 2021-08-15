package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go.uber.org/zap"

	"github.com/abx123/library/constant"
	"github.com/abx123/library/handler/presenter"
	"github.com/abx123/library/services"
)

type postUpsertBookRequest struct {
	BookID          int64   `json:"id" form:"id"`
	ISBN            string  `json:"isbn" form:"isbn"`
	Title           string  `json:"title" form:"title"`
	Author          string  `json:"author" form:"author"`
	ImageURL        string  `json:"imageUrl" form:"imageUrl"`
	SmallImageURL   string  `json:"smallImageUrl" form:"smallImageUrl"`
	PublicationYear int64   `json:"publicationYear" form:"publicationYear"`
	AverageRating   float64 `json:"averageRating" form:"averageRating"`
	Status          int64   `json:"status" form:"status"`
	Publisher       string  `json:"publisher" form:"publisher"`
	Description     string  `json:"description" form:"description"`
	Categories      string  `json:"categories" form:"categories"`
	Language        string  `json:"language" form:"language"`
	Source          string  `json:"source" form:"source"`
	PageCount       int64   `json:"pageCount" form:"pageCount"`
}

type Handler struct {
	bookSvc services.Ibooks
	dbSvc   services.IdbService
}

func NewHandler(dbSvc services.IdbService, bSvc services.Ibooks) *Handler {
	return &Handler{
		dbSvc:   dbSvc,
		bookSvc: bSvc,
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
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}

	// Return ok
	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:        isbn,
		Title:       data.Title,
		Author:      data.Authors,
		ImageURL:    data.ImageURL,
		UserID:      data.UserID,
		Status:      data.Status,
		Source:      data.Source,
		Language:    data.Language,
		PageCount:   data.PageCount,
		Description: data.Description,
		Publisher:   data.Publisher,
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
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}

	// Return ok
	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:        data.ISBN,
		Title:       data.Title,
		Author:      data.Authors,
		ImageURL:    data.ImageURL,
		Status:      data.Status,
		Source:      data.Source,
		Language:    data.Language,
		PageCount:   data.PageCount,
		Description: data.Description,
		Publisher:   data.Publisher,
	})
}

func (h *Handler) ListBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	limit, offset, err := getLimitAndOffest(c)
	userId := c.Param("userId")
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
			ISBN:     d.ISBN,
			Title:    d.Title,
			Author:   d.Authors,
			ImageURL: d.ImageURL,
			UserID:   d.UserID,
			Status:   d.Status,
		})
	}

	return c.JSON(http.StatusOK, books)
}

func (h *Handler) UpsertBook(c echo.Context) (err error) {
	r := &postUpsertBookRequest{}
	userId := c.Param("userId")
	// isbn := c.Param("isbn")
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	ctx := context.WithValue(c.Request().Context(), constant.ContextKeyRequestID, c.Response().Header().Get(echo.HeaderXRequestID))
	if err = c.Bind(r); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, constant.ErrInvalidRequest))
	}
	if _, err := strconv.ParseInt(r.ISBN, 10, 64); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}
	book, err := h.dbSvc.Upsert(ctx, r.ISBN, r.Title, r.Author, r.ImageURL, r.SmallImageURL, r.Publisher, userId, r.Description, r.Categories, r.Language, r.Source, r.PublicationYear, r.Status, r.PageCount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}

	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:     book.ISBN,
		Title:    book.Title,
		Author:   book.Authors,
		ImageURL: book.ImageURL,
		UserID:   book.UserID,
		Status:   book.Status,
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
