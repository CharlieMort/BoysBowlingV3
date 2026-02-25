package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func CheckNilError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func CalculateScore(scoreCard string) int {
	// -- | -- | 9/ | -/ | 72 | -1 | -- | 8- | 7- | 72  = 61
	scoreCard = strings.TrimPrefix(scoreCard, " ")
	scoreCheck, err := strconv.ParseInt(strings.Split(scoreCard, "= ")[1], 10, 16)
	CheckNilError(err, "error parseing score")
	scoreCard = strings.TrimRight(scoreCard, "  = "+strings.Split(scoreCard, "= ")[1])
	frames := strings.Split(scoreCard, " | ")
	framesNums := make([]int, 0)
	modifier := 1
	for i, frame := range frames {
		if frame == "X" {
			if len(framesNums) > 0 {
				fr
			}
		}
		frameInt1, err := strconv.ParseInt(string(frame[0]), 10, 8)
		CheckNilError(err, "Failed parsing frame1")
		frameInt2, err := strconv.ParseInt(string(frame[1]), 10, 8)
		CheckNilError(err, "Failed parsing frame2")

	}
}

func main() {
	fmt.Println("Hello Bowlers")
	SetupDatabase()
}
