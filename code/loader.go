package main

import (
    "bufio"
    "fmt"
    "os"
	"errors"
	"strings"
	"strconv"
)

/*
 * Reads a scenarios file according to this format:
 * https://movingai.com/benchmarks/formats.html
 * Returns a non-nil error if something goes wrong.
 */
func LoadScenarios(path string) ([]Scenario, error) {
	scenarios := []Scenario{}

	file, err := os.Open(path)
    if err != nil {
		return scenarios, errors.New("Could not open file "+err.Error())
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	if line != "version 1" && line != "version 1.0" {
		msg := fmt.Sprintf("Bad first line \"%s\"", line)
		return scenarios, errors.New(msg)
	}

	lineCounter := 2 // Starting at 2 because we have already read some lines
	for scanner.Scan() {
		line = scanner.Text()
		splits := strings.Fields(line)
		if len(splits) != 9 {
			msg := fmt.Sprintf("Expected 9 data points on line %d", lineCounter)
			return scenarios, errors.New(msg)
		}

		scenario := Scenario{}
		scenario.Bucket, err = strconv.Atoi(splits[0])
		if err != nil {
			msg := fmt.Sprintf("Non-int bucket \"%s\"", splits[0])
			return scenarios, errors.New(msg)
		}

		scenario.MapName = splits[1]

		scenario.Width, err = strconv.Atoi(splits[2])
		if err != nil {
			msg := fmt.Sprintf("Non-int width \"%s\"", splits[2])
			return scenarios, errors.New(msg)
		}

		scenario.Height, err = strconv.Atoi(splits[3])
		if err != nil {
			msg := fmt.Sprintf("Non-int height \"%s\"", splits[3])
			return scenarios, errors.New(msg)
		}

		startX, err := strconv.Atoi(splits[4])
		if err != nil {
			msg := fmt.Sprintf("Non-int start x-coordinate \"%s\"", splits[4])
			return scenarios, errors.New(msg)
		}

		startY, err := strconv.Atoi(splits[5])
		if err != nil {
			msg := fmt.Sprintf("Non-int start y-coordinate \"%s\"", splits[5])
			return scenarios, errors.New(msg)
		}

		goalX, err := strconv.Atoi(splits[6])
		if err != nil {
			msg := fmt.Sprintf("Non-int goal x-coordinate \"%s\"", splits[6])
			return scenarios, errors.New(msg)
		}

		goalY, err := strconv.Atoi(splits[7])
		if err != nil {
			msg := fmt.Sprintf("Non-int goal y-coordinate \"%s\"", splits[7])
			return scenarios, errors.New(msg)
		}

		scenario.OptimalLength, err = strconv.ParseFloat(splits[8], 64)
		if err != nil {
			msg := fmt.Sprintf("Non-float optimal length \"%s\"", splits[8])
			return scenarios, errors.New(msg)
		}

		scenario.Path  = path
		scenario.Start = NewNode(startX, startY)
		scenario.Goal  = NewNode(goalX,  goalY)

		scenarios = append(scenarios, scenario)
		lineCounter++
	}
	return scenarios, nil
}

/*
 * Reads a map file according to this format:
 * https://movingai.com/benchmarks/formats.html
 * Returns a bool matrix where true=blocked and false=traversable.
 * Returns a non-nil error if something goes wrong.
 * Only the '.' and '@' characters are accepted.
 */
func LoadMap(path string) ([][]bool, error) {
	m := [][]bool{}

	file, err := os.Open(path)
    if err != nil {
		return m, errors.New("Could not open the file")
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	if line != "type octile" {
		msg := fmt.Sprintf("Bad first line \"%s\"", line)
		return m, errors.New(msg)
	}

	scanner.Scan()
	line = scanner.Text()
	splits := strings.Split(line, " ")
	if len(splits) != 2 || splits[0] != "height" {
		msg := fmt.Sprintf("Bad second line \"%s\"", line)
		return m, errors.New(msg)
	}
	height, err := strconv.Atoi(splits[1])
	if err != nil {
		msg := fmt.Sprintf("Non-int height \"%s\"", splits[1])
		return m, errors.New(msg)
	}
	
	scanner.Scan()
	line = scanner.Text()
	splits = strings.Split(line, " ")
	if len(splits) != 2 || splits[0] != "width" {
		msg := fmt.Sprintf("Bad third line \"%s\"", line)
		return m, errors.New(msg)
	}
	width, err := strconv.Atoi(splits[1])
	if err != nil {
		msg := fmt.Sprintf("Non-int width \"%s\"", splits[1])
		return m, errors.New(msg)
	}

	scanner.Scan()
	line = scanner.Text()
	if line != "map" {
		msg := fmt.Sprintf("Bad fourth line \"%s\"", line)
		return m, errors.New(msg)
	}

	m = make([][]bool, height)
	for row := 0; row < height; row++ {
		m[row] = make([]bool, width)
    }

	row := 0
	lineNumber := 5 // We have already read some lines
    for scanner.Scan() {
		if row >= height {
			panic("Height mismatch in map file")
		}
		line := scanner.Text()
		if len(line) != width {
			panic("Width mismatch in map file")
		}
		for i, r := range(line) {
			var b bool
			switch r {
				case '.': // open space
					b = false
					break
				case '@': // wall
					b = true
					break
				default:
					msg := fmt.Sprintf("Bad rune '%c' on line %d", r, lineNumber)
					return m, errors.New(msg)
			}
			m[row][i] = b
		}
		row++
    }
	if height != row {
		panic("Height mismatch in map file")
	}
	return m, nil
}
