package main

import (
	"time"
	"math"
	"fmt"
	"strconv"
	"strings"
	"os"
)

var _ = fmt.Sprint("") // To be able to keep the fmt import

/*
 * Performs test runs of the scenario. Returns the following things:
 * turn count
 * path length
 * average angle of turns (radians)
 * average runtime (ms)
 */
func testOneScenario(start, goal Node, m [][]bool, algo func (m [][]bool, start, goal Node) []Node, trials int) ([]Node, int, float64, float64, int) {
	var path []Node

	// Get path and average runtime
	totalRuntime := 0.0
	for i := 0; i < trials; i++ {
		before  := time.Now()
		path     = algo(m, start, goal)
		after   := time.Now()
		elapsed := after.Sub(before)
		totalRuntime += float64(elapsed.Milliseconds())
	}
	avgRuntime := int(math.Round(totalRuntime/float64(trials)))

	// Calculate path length
	pathLen := 0.0
	for i := 0; i < len(path)-1; i++ {
		n1 := path[i]
		n2 := path[i+1]
		pathLen += StraightLineDist(n1, n2)
	}

	// Calculate turn count and average angle of turns
	turns    := 0
	avgAngle := 0.0
	for i := 0; i < len(path)-2; i++ {
		n1 := path[i]
		n2 := path[i+1]
		n3 := path[i+2]
		// There are two vectors (n1,n2) and (n2,n3)
		v1_x, v1_y := float64(n2.X - n1.X), float64(n2.Y - n1.Y)
		v2_x, v2_y := float64(n3.X - n2.X), float64(n3.Y - n2.Y)
		dot    := v1_x * v2_x + v1_y * v2_y
		v1_len := math.Sqrt(v1_x * v1_x + v1_y * v1_y)
		v2_len := math.Sqrt(v2_x * v2_x + v2_y * v2_y)
		a := dot / (v1_len * v2_len)
		a  = math.Max(-1, math.Min(1, a)) // Rounding errors may produce values outside [-1,1] so clamp it.
		angle := math.Acos(a)
		if angle >= 0.001 {
			// We are turning at node n1
			avgAngle += angle
			turns++
		}
	}

	if turns > 0 {
		avgAngle /= float64(turns)
	}

	return path, turns, pathLen, avgAngle, avgRuntime
}

func MustParsePathfindingFunction(algoName string) func([][]bool, Node, Node) []Node {
	switch strings.ToLower(algoName) {
		case "dijkstra":
			return Dijkstra
		case "astar":
			return AStar
		case "astar-ps":
			return AStarPs
		case "thetastar":
			return ThetaStar
		// case "ap-thetastar":
			// return true, ApThetaStar
	}
	fmt.Printf("Unknown algorithm \"%s\"\n", algoName)
	os.Exit(1)
	return func([][]bool, Node, Node) []Node { return []Node{} }
}

func MustParseInt(arg string) int {
	n, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("Non-int argument \"%s\"\n", arg)
		os.Exit(1)
	}
	return n
}
