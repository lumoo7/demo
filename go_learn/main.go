package main

import "fmt"

func main() {
	var array1 = [][]int{{1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}}
	for i := 0; i < len(array1); i++ {
		for j := range 4 {
			fmt.Print(array1[i][j], " ")
		}
		fmt.Print("\n")
	}
}
