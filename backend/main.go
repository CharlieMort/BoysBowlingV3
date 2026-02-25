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

// func CalculateScore(scoreCard string) int {
// 	// -- | -- | 9/ | -/ | 72 | -1 | -- | 8- | 7- | 72  = 61
// 	scoreCard = strings.TrimPrefix(scoreCard, " ")
// 	scoreCheck, err := strconv.ParseInt(strings.Split(scoreCard, "= ")[1], 10, 16)
// 	CheckNilError(err, "error parseing score")
// 	scoreCard = strings.TrimRight(scoreCard, "  = "+strings.Split(scoreCard, "= ")[1])
// }

func AddScoreCard(scoreCard string, playerId int, gameId int) {
	scoreString := strings.TrimPrefix(scoreCard, " ")
	scoreTotal := strings.Split(scoreCard, "= ")[1]
	scoreString = strings.TrimSuffix(scoreString, "  = "+scoreTotal)
	totalNum, err := strconv.ParseInt(scoreTotal, 10, 16)
	CheckNilError(err, "error parsing score total for frame")

	InsertFrame(playerId, gameId, int(totalNum), scoreCard, "null")
}

func main() {
	fmt.Println("Hello Bowlers")
	SetupDatabase()
	defer db.Close()
	AddScoreCard(" -- | -- | 9/ | -/ | 72 | -1 | -- | 8- | 7- | 72  = 61", 0, 0)
}
