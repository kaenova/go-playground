package some

import (
	"errors"

	customerror "github.com/kaenova/go-playground/error_assert/custom_error"
)

func OtherFunction() error {
	return errors.New("this is just a normal error")
	return customerror.MakeCustomError(1, "This is an error")
}
