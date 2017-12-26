package main

import "fmt"

func main() {
	fmt.Println("hola asdfy n")
	s := "feliz avidad"
	r := []int32(s)
	fmt.Println(r)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	fmt.Print(string(r))
}
