package handler

import (
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

// Handler defines a handler struct
type Handler struct {
	bookSvc services.Ibooks
	dbSvc   services.IdbService
}

// NewHandler returns a new instance of Handler
func NewHandler(dbSvc services.IdbService, bSvc services.Ibooks) *Handler {
	return &Handler{
		dbSvc:   dbSvc,
		bookSvc: bSvc,
	}
}

// GetBook resolves GET /{userID}/book/{isbn}, retreives details of a book from database
func (h *Handler) GetBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	isbn := c.Param("isbn")
	userId := c.Param("userId")

	data, err := h.dbSvc.Get(c.Request().Context(), isbn, userId)
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

// GetNewBook resolves GET /book/:isbn, retreives details of a book from providers.
func (h *Handler) GetNewBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	isbn := c.Param("isbn")
	if _, err := strconv.ParseInt(isbn, 10, 64); err != nil {
		// Invalid request parameter
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	data, err := h.bookSvc.Get(c.Request().Context(), isbn)
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

// ListBook resolves GET /{userID}/books, retreives the list of books related to the userID
func (h *Handler) ListBook(c echo.Context) (err error) {
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	limit, offset, err := getLimitAndOffest(c)
	userId := c.Param("userId")
	if err != nil {
		zap.L().Error(constant.ErrInvalidRequest.Error(), zap.Error(err))
		return c.JSON(http.StatusBadRequest, presenter.ErrResp(reqID, err))
	}

	data, err := h.dbSvc.List(c.Request().Context(), limit, offset, userId)
	if err != nil {
		// Error while querying database
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}
	books := []*presenter.Book{}
	for _, d := range data {
		books = append(books, &presenter.Book{
			ISBN:            d.ISBN,
			Title:           d.Title,
			Author:          d.Authors,
			ImageURL:        d.ImageURL,
			SmallImageURL:   d.SmallImageURL,
			Publisher:       d.Publisher,
			Description:     d.Description,
			PageCount:       d.PageCount,
			Categories:      d.Categories,
			Language:        d.Language,
			PublicationYear: d.PublicationYear,
			UserID:          d.UserID,
			Status:          d.Status,
			Source:          d.Source,
		})
	}

	return c.JSON(http.StatusOK, books)
}

// UpsertBook resolves POST /{userID}/book, updates a database book record if record is found, creates a new record if no record found.
func (h *Handler) UpsertBook(c echo.Context) (err error) {
	r := &postUpsertBookRequest{}
	userId := c.Param("userId")
	// isbn := c.Param("isbn")
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
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
	book, err := h.dbSvc.Upsert(c.Request().Context(), r.ISBN, r.Title, r.Author, r.ImageURL, r.SmallImageURL, r.Publisher, userId, r.Description, r.Categories, r.Language, r.Source, r.PublicationYear, r.Status, r.PageCount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, presenter.ErrResp(reqID, err))
	}

	return c.JSON(http.StatusOK, &presenter.Book{
		ISBN:            book.ISBN,
		Title:           book.Title,
		Author:          book.Authors,
		ImageURL:        book.ImageURL,
		SmallImageURL:   book.SmallImageURL,
		Publisher:       book.Publisher,
		Description:     book.Description,
		PageCount:       book.PageCount,
		Categories:      book.Categories,
		Language:        book.Language,
		PublicationYear: book.PublicationYear,
		UserID:          book.UserID,
		Status:          book.Status,
		Source:          book.Source,
	})
}

// Ping resolves GET /ping, returns "Pong", used for healthcheck.
func (h *Handler) Ping(c echo.Context) (err error) {
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
