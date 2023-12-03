package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Round struct {
	Balls []Ball
}

func NewRoundFromString(sRound string) Round {
	zap.L().Sugar().Debugw("Got round",
		"round", sRound,
	)
	sBalls := strings.Split(sRound, ",")
	var balls []Ball
	for _, sBall := range sBalls {
		ball := NewBallFromString(sBall)
		zap.L().Sugar().Debugw("received ball object",
			"ball", ball,
		)
		balls = append(balls, ball)
	}
	round := Round{
		Balls: balls,
	}
	return round
}

type Ball struct {
	Color string
	Count int
}

func NewBallFromString(sBall string) Ball {
	sBallParts := strings.Split(strings.Trim(sBall, " "), " ")
	sCount, parsedColor := sBallParts[0], sBallParts[1]
	parsedCount, _ := strconv.Atoi(sCount)

	var ball Ball
	ball.Count = parsedCount
	ball.Color = parsedColor

	zap.L().Sugar().Debugw("Got ball",
		"sBall", sBall,
		"ballPars", sBallParts,
		"parsedColor", parsedColor,
		"parsedCount", parsedCount,
		"ball", ball,
	)
	return ball
}

type Game struct {
	Id     int
	Rounds []Round
}

type GameEvaluator struct {
	Balls []Ball
}

type GameMinimizer struct {
	balls map[string]int
}

func (gameMinimizer *GameMinimizer) GetBallCount(color string) int {
	return gameMinimizer.balls[color]
}

func (gameMinimizer *GameMinimizer) Colors() []string {
	var colors []string
	for key, _ := range gameMinimizer.balls {
		colors = append(colors, key)
	}
	return colors
}

func (gameMinimizer *GameMinimizer) Minimize(game Game) {
	// For each round in the game
	if gameMinimizer.balls == nil {
		gameMinimizer.balls = map[string]int{}
	}
	for _, roundInGame := range game.Rounds {

		// Determine the maxSeen for each color ball in the round, store it in the game evaluator
		for _, ballInRound := range roundInGame.Balls {
			if ballInRound.Count > gameMinimizer.balls[ballInRound.Color] {
				gameMinimizer.balls[ballInRound.Color] = ballInRound.Count
			}
		}
	}
}

func (gameEvaluator *GameEvaluator) Evaluate(game Game) bool {
	// Iterate over games, summing ball colors per round
	for roundIndex, round := range game.Rounds {
		perColorSums := make(map[string]int)
		for _, ball := range round.Balls {
			perColorSums[ball.Color] = perColorSums[ball.Color] + ball.Count
		}

		// For each round, compare per-color sums to evaluator balls
		for _, evaluatorBall := range gameEvaluator.Balls {

			// If round sum is greater than any corresponding evaluator ball, return false
			if perColorSums[evaluatorBall.Color] > evaluatorBall.Count {
				zap.L().Sugar().Debugw("invalid game",
					"gameId", game.Id,
					"color", evaluatorBall.Color,
					"evaluatorCount", evaluatorBall.Count,
					"roundCount", perColorSums[evaluatorBall.Color],
					"roundIndex", roundIndex,
				)
				return false
			}
		}
	}
	zap.L().Sugar().Debugw("valid game",
		"gameId", game.Id)
	return true
}

func NewGameFromString(text string) Game {
	// Golang regex isn't suitable for deeply nested captures
	// Going with string splitting, quick and dirty because
	// these advent projects have to be fast
	parts := strings.Split(text, ":")

	// Parse out game header
	game, rest := parts[0], parts[1]
	id, _ := strconv.Atoi(strings.Split(game, " ")[1])
	zap.L().Sugar().Debugw("Got game Id",
		"Id", id,
	)

	sRounds := strings.Split(rest, ";")
	var rounds []Round
	for _, sRound := range sRounds {
		round := NewRoundFromString(sRound)
		zap.L().Sugar().Debugw("Got round",
			"round", round,
		)
		rounds = append(rounds, round)
	}

	return Game{
		Id:     id,
		Rounds: rounds,
	}
}

func parseGamesFile(file io.Reader) ([]Game, error) {
	var games []Game
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return []Game{}, err
		} else {
			theGame := NewGameFromString(scanner.Text())
			zap.L().Sugar().Debugw("Got game",
				"game", theGame,
			)
			games = append(games, theGame)
		}
	}
	return games, nil
}

func loadGames(inputFileName string) ([]Game, error) {
	file, err := os.Open(fmt.Sprintf("input/%s", inputFileName))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			zap.L().Sugar().Error("err", err)
		}
	}(file)

	return parseGamesFile(file)

}

func main() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer zap.L().Sync() // flush buffer on close

	//filename := os.Args[1]
	filename := "input.txt"
	zap.L().Sugar().Infow("Parsing file", "filename", filename)

	// Load games
	games, err := loadGames(filename)
	if err != nil {
		panic(err)
	}
	zap.L().Sugar().Infow("Loaded games",
		"count", len(games),
	)

	// 12 red cubes, 13 green cubes, and 14 blue cubes
	gameEvaluator := GameEvaluator{
		Balls: []Ball{
			{Count: 12,
				Color: "red",
			}, {
				Count: 13,
				Color: "green",
			}, {
				Count: 14,
				Color: "blue",
			},
		},
	}
	var validGameIdSum, validGameCount int
	for _, game := range games {
		if gameEvaluator.Evaluate(game) {
			validGameIdSum += game.Id
			validGameCount += 1
		}
	}

	logger.Sugar().Infow("summed valid game ids",
		"validGameIdSum", validGameIdSum,
		"validGameCount", validGameCount,
	)

	var minimumBallPowerSum int
	for _, game := range games {
		var minimizer GameMinimizer
		minimizer.Minimize(game)

		var powAccumulator int
		for _, color := range minimizer.Colors() {
			if powAccumulator == 0 {
				powAccumulator = minimizer.GetBallCount(color)
			} else {
				powAccumulator *= minimizer.GetBallCount(color)
			}
		}
		minimumBallPowerSum += powAccumulator
	}

	zap.L().Sugar().Infow("Minimized games",
		"minimizedGameSum", minimumBallPowerSum,
	)
}
