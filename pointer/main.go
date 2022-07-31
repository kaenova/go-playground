package main

import "log"

func main() {
	var a *int64
	if a != nil && *a < 0 {
		log.Println(*a)
	}
}
