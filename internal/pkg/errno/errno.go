package errno

type Errno struct {
	Code    int
	Message string
}

// Error implement the `Error` method in error interface.
func (err Errno) Error() string {
	return err.Message
}

func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, err.Error()
}
