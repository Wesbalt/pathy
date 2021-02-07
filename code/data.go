package main

import (
	"math"
)

type Node struct {
	X int
	Y int
	Timestamp int
}

func NewNode(x, y int) Node {
	n := Node{}
	n.X = x
	n.Y = y
	n.Timestamp = -1
	return n
}

func StraightLineDist(n1, n2 Node) float64 {
	dx := math.Abs(float64(n1.X - n2.X))
	dy := math.Abs(float64(n1.Y - n2.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

type Scenario struct {
	Path     string // Filepath to its belonging scenarios file
	Bucket   int
	MapName  string
	Width    int
	Height   int
	Start    Node
	Goal     Node
	OptimalLength  float64
}
