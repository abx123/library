package constant

import "errors"

var (
	ErrDBErr = errors.New("database returns error")

	ErrBookNotFound = errors.New("book not found")

	ErrInvalidRequest = errors.New("invalid request parameter")

	ErrGapiError = errors.New("google api return error")

	ErrGoodreadError = errors.New("goodread api return error")

	ErrRetrievingBookDetails = errors.New("error retrieving book details")
)
