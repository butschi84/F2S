package main

import "fmt"

func main() {
	strings := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	printSlice(strings)

	if len(strings) > 0 {
		firstNode := strings[0]
		s2 := make([]string, len(strings))
		copy(s2, strings[1:])
		s2[len(strings)-1] = firstNode
		strings = nil
		printSlice(s2)
	}

}

func printSlice(s []string) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
