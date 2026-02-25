package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	fmt.Println("  = " + scoreTotal)
	scoreString = strings.TrimSuffix(scoreString, "  = "+scoreTotal)
	totalNum, err := strconv.ParseInt(scoreTotal, 10, 16)
	CheckNilError(err, "error parsing score total for frame")

	InsertFrame(playerId, gameId, int(totalNum), scoreString, "null")
	LastModifiedStats = time.Now()
}

func GetAllFrames() []Frame {
	if LastFetchedFrames.Before(LastModifiedFrames) {
		newFrames := GetFrames()
		framesCache = newFrames
	} else {
		fmt.Println("Using frame cache")
	}
	return framesCache
}

func BestScore() []PodiumElement {
	frames := GetAllFrames()
	if len(frames) < 3 {
		fmt.Println("bad cache")
		LastModifiedFrames = time.Now()
		return []PodiumElement{}
	}
	scores := make(map[int]int)
	for _, frames := range frames {
		scores[frames.PlayerId] += frames.Total
	}

	scoreArr := make([]PodiumElement, 0)
	for player := range scores {
		scoreArr = append(scoreArr, PodiumElement{
			Id:    player,
			Name:  GetNameFromId(player),
			Score: scores[player],
		})
	}

	slices.SortFunc(scoreArr, func(a, b PodiumElement) int {
		return cmp.Compare(b.Score, a.Score)
	})

	return scoreArr[:3]
}

type StatList struct {
	Title  string          `json:"title"`
	Podium []PodiumElement `json:"podium"`
}

func Stats(c *gin.Context) {
	if LastFetchedStats.After(LastModifiedStats) {
		fmt.Println("Using Stat Cache")
		LastFetchedStats = time.Now()
		c.JSON(200, statsCache)
		return
	}
	var stats []StatList
	stats = append(stats, StatList{
		Title:  "Best Score",
		Podium: BestScore(),
	})

	LastFetchedStats = time.Now()
	statsCache = stats
	c.JSON(200, stats)
}

type FrameSimple struct {
	Name      string `json:"name"`
	Scorecard string `json:"scorecard"`
}

type AddGameBody struct {
	DatePlayed string        `json:"datePlayed"`
	Frames     []FrameSimple `json:"frames"`
}

type Reponse struct {
	Msg string `json:"msg"`
}

func AddGame(c *gin.Context) {
	dat, err := io.ReadAll(c.Request.Body)
	CheckNilError(err, "erroring reading add game request body")
	var body AddGameBody
	json.Unmarshal(dat, &body)

	gameId, err := InsertGame("Untitled Game", body.DatePlayed)
	CheckNilError(err, "error inserting game")

	for _, frame := range body.Frames {
		playerId := GetUserIdFromName(frame.Name)
		AddScoreCard(frame.Scorecard, playerId, int(gameId))
	}

	c.JSON(200, Reponse{Msg: "Accepted"})
}

var LastModifiedFrames time.Time
var LastFetchedFrames time.Time
var framesCache []Frame

var LastModifiedStats time.Time
var LastFetchedStats time.Time
var statsCache []StatList

type PodiumElement struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func main() {
	fmt.Println("Hello Bowlers")
	LastModifiedFrames = time.Now()
	LastModifiedStats = time.Now()

	SetupDatabase()
	defer db.Close()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/stats", Stats)
	r.POST("/game", AddGame)

	r.Run(":8888")
	fmt.Println(BestScore())
}
