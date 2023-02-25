package errno

import "github.com/pkg/errors"

var (
	// OK represents a successful request.
	OK = &Errno{Code: 0, Message: "OK"}

	// InternalServerError represents all unknown server-side errors.
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}

	// ErrDatabase represents a database error.
	ErrDatabase = &Errno{Code: 20002, Message: "Database error."}

	ErrCreateDocument = errors.New("failed to create document")
)
