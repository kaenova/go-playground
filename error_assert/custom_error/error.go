package customerror

import "fmt"

type MyCustomError struct {
	Code    int
	Message string
}

func (c MyCustomError) Error() string {
	return fmt.Sprintf("[%d] %s", c.Code, c.Message)
}

func MakeCustomError(code int, message string) MyCustomError {
	return MyCustomError{
		Code:    code,
		Message: message,
	}
}
