package main

import (
	"fmt"

	errordocs "github.com/kaenova/go-playground/error_docs/error_docs"
)

func main() {
	err := functionThatFail()
	fmt.Println(err.Error())
}

func functionThatFail() error {
	return errordocs.ErrorDocs(1)
}
