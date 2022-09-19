package errordocs

import "fmt"

type ErrorType string

const (
	System   ErrorType = "System"
	Internal ErrorType = "System"
	User     ErrorType = "System"
)

type ErrDocs struct {
	Code    int
	Type    ErrorType
	Message string
}

func (e ErrDocs) Error() string {
	return fmt.Sprintf("[%d][%s] %s", e.Code, e.Type, e.Message)
}
