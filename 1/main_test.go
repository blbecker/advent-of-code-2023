package main

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

var EXPECTED_CALIBRATION_SUM = 142
var INPUT_FILE_NAME = "test_input.txt"
var INPUT_WORDS_FILE_NAME = "test_input_with_words.txt"

func TestParse(t *testing.T) {
	Convey("Parses the input correctly", t, func() {
		logger, err := zap.NewDevelopment()
		Convey("With only ints", func() {
			if err != nil {
				panic(err)
			}
			parsedValues, err := parseInput(INPUT_FILE_NAME, logger)

			if err != nil {
				fmt.Println(err.Error())

			}
			So(parsedValues, ShouldResemble, []int{12, 38, 15, 77, 52})
		})

		Convey("with words and ints", func() {
			parsedValues, err := parseInput(INPUT_WORDS_FILE_NAME, logger)

			if err != nil {
				fmt.Println(err.Error())

			}
			So(parsedValues, ShouldResemble, []int{29, 83, 13, 24, 42, 14, 76})
		})

		Convey("with no digits", func() {
			parsedValues, err := parseInput("test_input_no_digits.txt", logger)
			So(err, ShouldBeNil)
			So(parsedValues, ShouldResemble, []int{0})
		})

	})
}

func TestSumArray(t *testing.T) {
	Convey("Summing an array gives the right result", t, func() {
		testArray := []int{1, 2, 3, 4, 5}
		So(sumArray(testArray...), ShouldEqual, 15)
	})
}

func TestParseSymbol(t *testing.T) {
	Convey("Parse symbol correctly parses", t, func() {
		Convey("a symbol", func() {
			parsedValue, err := parseSymbol("five")
			if err != nil {
				panic(err)
			}
			So(parsedValue, ShouldEqual, 5)
		})
		Convey("a bad string", func() {
			parsedValue, err := parseSymbol("flurp")
			So(err, ShouldNotBeNil)
			So(parsedValue, ShouldEqual, 0)
		})
	})
}
