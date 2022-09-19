package main

import (
	"fmt"

	customerror "github.com/kaenova/go-playground/error_assert/custom_error"
	"github.com/kaenova/go-playground/error_assert/some"
)

func main() {
	// Try to call function with error
	err := some.OtherFunction()

	// assert error
	asErr, ok := err.(customerror.MyCustomError)

	if !ok {
		fmt.Println("this is normal error", err)
		return
	}

	fmt.Println("this is assert error", asErr)
	fmt.Println(asErr.Code)
	fmt.Println(asErr.Message)
	return
}
