package main

import (
	"errors"
	"fmt"
	"reflect"
)

type anyPointer interface{}

var (
	ErrDataTypeNotDefined = errors.New("data type not defined")
)

const (
	TargetTag = "s3"
)

type A struct {
	F1    string `s3:"resolve;kasjhdkasjd;aksjdkasjcda"`
	F2    string
	F3    string
	Other *B
}

type B struct {
	F1 string
	F2 *[]string `s3:"resolve"`
	F3 string
	F4 []*string `s3:"resolve"`
}

func main() {
	x1 := "hey"
	x2 := "i just met you"
	x3 := "and this is crazy"
	b := B{"Halo", &[]string{"halo", "pakabar", "baik?"}, "Apakabar", []*string{&x1, &x2, &x3}}
	a := A{"Baik", "Banget", "Kok", &b}

	// a := "test"

	ResolveS3(&a)

	fmt.Println(a.F1)
	fmt.Println(a.F2)
	fmt.Println(a.F3)
	fmt.Println(a.Other.F1)
	fmt.Println(*a.Other.F2)
	fmt.Println(a.Other.F3)
	fmt.Println(*a.Other.F4[0])
	fmt.Println(*a.Other.F4[1])
	fmt.Println(*a.Other.F4[2])

}

func ResolveS3(d anyPointer) {

	t := reflect.TypeOf(d)
	v := reflect.ValueOf(d)

	// Check if its a pointer
	if t.Kind() != reflect.Ptr {
		panic("Cannot resolve because input is not a pointer")
	}

	RecurseResolve(v, t, false)
}

func RecurseResolve(v reflect.Value, t reflect.Type, resolve bool) {

	// Only resolve if its a string
	if t.Kind() == reflect.String && resolve {
		v.SetString("haha")
		return
	}

	// Change pointer to element
	if t.Kind() == reflect.Pointer {
		v = v.Elem()
		t = t.Elem()
		RecurseResolve(v, t, resolve)
		return
	}

	// Handle slice or array
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			RecurseResolve(v.Index(i), v.Index(i).Type(), resolve)
		}
		return
	}

	// Handle struct, resolve every field
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			isTag := t.Field(i).Tag.Get(TargetTag) != ""
			RecurseResolve(v.Field(i), t.Field(i).Type, isTag)
		}
		return
	}

}
