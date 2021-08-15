package presenter

// Error defines an error object
type Error struct {
	ReqId  string `json:"requestID"`
	ErrMsg string `json:"message"`
}

// ErrResp wraps the requestID and error into a single Error object
func ErrResp(reqID string, err error) *Error {
	return &Error{
		ReqId:  reqID,
		ErrMsg: err.Error(),
	}
}
