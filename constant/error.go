package constant

import "errors"

var (
	// ErrDBErr ...
	ErrDBErr = errors.New("database returns error")

	// ErrBookNotFound ...
	ErrBookNotFound = errors.New("book not found")

	// ErrInvalidRequest ...
	ErrInvalidRequest = errors.New("invalid request parameter")

	// ErrRetrievingBookDetails ...
	ErrRetrievingBookDetails = errors.New("error retrieving book details")
)
