package apperror

type Error string

func (e Error) String() string {
	return string(e)
}

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNotFound      Error = "Not found"
	ErrBadRequest    Error = "Bad request"
	ErrAlreadyExists Error = "Already exists"
	ErrInternal      Error = "Internal server error"
	ErrData          Error = "Data error"
)

func NewError(msg string) *Error {
	return (*Error)(&msg)
}
