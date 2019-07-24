package main

import "fmt"

func main()  {
	str1 := []string	{"hello", "world"}

	for i, st := range str1 {
		st = "HELLO"
		fmt.Printf("st[%d]: %s\n", i, st)
	}

	for i, st := range str1 {
		fmt.Printf("st[%d]: %s\n", i, st)
	}
}
