package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5}
	arr2 := []int{4, 5, 6, 7}

	arr = append(arr, arr2...)

	fmt.Println(arr)
}
