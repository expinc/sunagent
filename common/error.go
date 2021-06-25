package common

import "fmt"

const (
	ErrorUnknown              = -1
	ErrorNone                 = 0
	ErrorUnexpected           = 1
	ErrorInvalidParameter     = 2
	ErrorNotFound             = 3
	ErrorUnauthorized         = 4
	ErrorNotAllowed           = 5
	ErrorInsufficientResource = 6
)

type Error struct {
	code int
	msg  string
}

func NewError(code int, msg string) error {
	return Error{
		code: code,
		msg:  msg,
	}
}

func (err Error) Error() string {
	return fmt.Sprintf("%d - %s", err.code, err.msg)
}

func (err Error) Code() int {
	return err.code
}
