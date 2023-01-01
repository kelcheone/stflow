package main

import (
	"log"
	"strconv"
	"time"
)

func convertDateToUnix(s string) int64 {

	layout := "2006-01-02 15:04:05Z"
	t, err := time.Parse(layout, s)
	if err != nil {
		log.Fatal(err)
	}
	return t.Unix()
}

// convert 1.234 to 1234
func convertViewsToNumber(s string) float64 {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func convertToInt(s string) int {
	// include those with decimal points
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func checkError(s string, err error) {
	if err != nil {
		log.Fatal(s, err)
	}
}
