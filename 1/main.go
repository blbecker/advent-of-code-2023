package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const InputRegex = "(?:\\d|one|two|three|four|five|six|seven|eight|nine)"
const DigitRegex = "(?:\\d)"

var inputRegex = regexp.MustCompile(InputRegex)
var digitRegex = regexp.MustCompile(DigitRegex)

var digitStrings = []string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9",
	"one", "two", "three", "four", "five",
	"six", "seven", "eight", "nine",
}

func sumArray(array ...int) int {
	result := 0
	for _, value := range array {
		result += value
	}
	return result
}

func parseSymbol(symbol string) (int, error) {
	if digitRegex.Match([]byte(symbol)) {
		return strconv.Atoi(symbol)
	} else {
		switch symbol {
		case "one":
			return 1, nil
		case "two":
			return 2, nil
		case "three":
			return 3, nil
		case "four":
			return 4, nil
		case "five":
			return 5, nil
		case "six":
			return 6, nil
		case "seven":
			return 7, nil
		case "eight":
			return 8, nil
		case "nine":
			return 9, nil
		default:
			return 0, fmt.Errorf(
				"%s is not a valid digit word or digit", symbol,
			)
		}
	}
}

func findAllDigitsIteratively(input string, logger *zap.Logger) []int {
	var allDigits []int
	chars := []rune(input)
	for i := 0; i < len(chars); i++ {
		remainder := chars[i:]
		for _, digitString := range digitStrings {
			if strings.HasPrefix(string(remainder), digitString) {
				parsedDigit, _ := parseSymbol(digitString)

				allDigits = append(allDigits, parsedDigit)
			}
		}

	}
	return allDigits
}

func lineToCoordinate(line string, logger *zap.Logger) (int, error) {
	allDigits := findAllDigitsIteratively(line, logger)
	if len(allDigits) < 1 {
		return 0, fmt.Errorf("line '%s' had no digits in it", line)
	}
	firstDigit := allDigits[0]

	secondDigit := allDigits[len(allDigits)-1]

	coordinate := firstDigit*10 + secondDigit
	logger.Sugar().Infow("Parsing line...",
		"line", line,
		"allDigits", allDigits,
		"firstDigit", firstDigit,
		"secondDigit", secondDigit,
		"coordinate", coordinate,
	)
	return coordinate, nil
}

func scanToArray(inReader io.Reader, logger *zap.Logger) ([]int, error) {
	var parsedValues []int
	scanner := bufio.NewScanner(inReader)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return []int{}, err
		} else {
			coordinate, err := lineToCoordinate(scanner.Text(), logger)
			if err != nil {
				logger.Sugar().Error("err", err)
			}
			parsedValues = append(parsedValues, coordinate)
		}

	}
	return parsedValues, nil
}

func parseInput(inputFileName string, logger *zap.Logger) ([]int, error) {
	file, err := os.Open(fmt.Sprintf("input/%s", inputFileName))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Sugar().Error("err", err)
		}
	}(file)

	return scanToArray(file, logger)

}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flush buffer on close

	filename := os.Args[1]
	logger.Sugar().Info("Parsing file", "filename", filename)
	coordinates, err := parseInput(filename, logger)
	if err != nil {
		panic(err)
	}
	sumOfCoordinates := sumArray(coordinates...)
	fmt.Printf("Sum: %d\n", sumOfCoordinates)
}
