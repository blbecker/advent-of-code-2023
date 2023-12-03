package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func TestLoadGames(t *testing.T) {
	Convey("Load games correctly loads a games file", t, func() {
		logger, _ := zap.NewDevelopment()
		zap.ReplaceGlobals(logger)
		games, err := loadGames("input_test.txt")
		if err != nil {
			panic(err)
		}
		So(games, ShouldNotBeNil)
		So(len(games), ShouldEqual, 5)
	})
}

func TestGameEvaluator_Evaluate(t *testing.T) {
	Convey("Game validator", t, func() {
		evaluator := GameEvaluator{
			Balls: []Ball{
				{
					Color: "red",
					Count: 3,
				},
				{
					Color: "blue",
					Count: 3,
				},
				{
					Color: "green",
					Count: 3,
				},
			},
		}
		Convey("fails invalid games", func() {
			game := Game{
				Id: 0,
				Rounds: []Round{
					{
						Balls: []Ball{
							{
								Count: 2,
								Color: "red",
							},
							{
								Count: 5,
								Color: "blue",
							},
						},
					},
				},
			}
			So(evaluator.Evaluate(game), ShouldBeFalse)
		})
		Convey("passes valid games", func() {
			game := Game{
				Id: 0,
				Rounds: []Round{
					{
						Balls: []Ball{
							{
								Count: 2,
								Color: "red",
							},
							{
								Count: 3,
								Color: "blue",
							},
						},
					},
				},
			}
			So(evaluator.Evaluate(game), ShouldBeTrue)
		})
	})
}

func TestGameMinimizer_Minimize(t *testing.T) {
	Convey("Minimize should determine the minimum number of "+
		"per-color balls required for a game to be valid", t, func() {
		game := Game{
			Id: 0,
			Rounds: []Round{
				{
					Balls: []Ball{
						{
							Color: "blue",
							Count: 3,
						},
						{
							Color: "red",
							Count: 4,
						},
						{
							Color: "green",
							Count: 3,
						},
					},
				},
				{
					Balls: []Ball{
						{
							Color: "blue",
							Count: 7,
						},
						{
							Color: "green",
							Count: 1,
						},
					},
				},
			},
		}
		minimizer := GameMinimizer{}
		minimizer.Minimize(game)
		So(minimizer.GetBallCount("red"), ShouldEqual, 4)
		So(minimizer.GetBallCount("green"), ShouldEqual, 3)
		So(minimizer.GetBallCount("blue"), ShouldEqual, 7)
	})
}
