package main

import (
	"fmt"
	"time"
)

func main() {
	// get time it takes to run the program
	start := time.Now()
	get_all_pages("svelte")
	fmt.Printf("Time taken Golang: %v \n", time.Since(start))

}
