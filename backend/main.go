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
	"github.com/gin-gonic/contrib/static"
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

func GetPodiumFromCalc(scores map[int]int) []PodiumElement {
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

func BestScore() []PodiumElement {
	frames := GetAllFrames()
	scores := make(map[int]int)
	for _, frames := range frames {
		scores[frames.PlayerId] += frames.Total
	}

	return GetPodiumFromCalc(scores)
}

func HighScore() []PodiumElement {
	scores := DoQueryGetFrames(`SELECT * FROM frames ORDER BY total DESC LIMIT 3;`)
	out := make([]PodiumElement, 0)
	for _, frame := range scores {
		fmt.Println(frame.Scorecard)
		out = append(out, PodiumElement{
			Id:    frame.PlayerId,
			Name:  GetNameFromId(frame.PlayerId),
			Score: frame.Total,
		})
	}
	return out
}

func BestLeagueScore() []PodiumElement {
	frames := DoQueryGetFrames("SELECT frames.id, frames.playerId, frames.gameId, frames.total, frames.scorecard, frames.imgPath FROM frames JOIN games ON games.id = frames.gameId WHERE games.league = 1;")

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

func MostGutters() []PodiumElement {
	frames := GetAllFrames()
	scores := make(map[int]int)
	for _, frames := range frames {
		scores[frames.PlayerId] += strings.Count(frames.Scorecard, "-")
	}

	return GetPodiumFromCalc(scores)
}

func MostStrikes() []PodiumElement {
	frames := GetAllFrames()
	scores := make(map[int]int)
	for _, frames := range frames {
		scores[frames.PlayerId] += strings.Count(frames.Scorecard, "X")
	}

	return GetPodiumFromCalc(scores)
}

func MostSpares() []PodiumElement {
	frames := GetAllFrames()
	scores := make(map[int]int)
	for _, frames := range frames {
		scores[frames.PlayerId] += strings.Count(frames.Scorecard, "/")
	}

	return GetPodiumFromCalc(scores)
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
		Title:  "Total Score",
		Podium: BestScore(),
	})
	stats = append(stats, StatList{
		Title:  "Highest Score",
		Podium: HighScore(),
	})
	stats = append(stats, StatList{
		Title:  "Total Leauge Score",
		Podium: BestLeagueScore(),
	})
	stats = append(stats, StatList{
		Title:  "Most Gutters",
		Podium: MostGutters(),
	})
	stats = append(stats, StatList{
		Title:  "Most Strikes",
		Podium: MostStrikes(),
	})
	stats = append(stats, StatList{
		Title:  "Most Spares",
		Podium: MostSpares(),
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

func GetAllFramesByGame(c *gin.Context) {
	type FrameGame struct {
		*Frame
		Date string `json:"date"`
	}

	rows := DoQuery(`
		SELECT frames.id, frames.playerId, frames.gameId, frames.total, frames.scorecard, frames.imgPath, games.date 
		FROM frames
		JOIN games ON frames.gameId = games.id
		ORDER BY games.date DESC;
	`)
	var frumes []FrameGame
	for rows.Next() {
		var f FrameGame
		rows.Scan(&f.Id, &f.PlayerId, &f.GameId, &f.Total, &f.Scorecard, &f.ImgPath, &f.Date)
		frumes = append(frumes, f)
	}

	c.JSON(200, frumes)
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

	r.Use(static.Serve("/", static.LocalFile("../frontend", true)))

	r.GET("/stats", Stats)
	r.POST("/game", AddGame)
	r.GET("/allgames", GetAllFramesByGame)

	r.Run(":8888")
	fmt.Println(BestScore())
}
