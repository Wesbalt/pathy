package main

type Node struct {
	X int
	Y int
}

func NewNode(x, y int) Node {
	n := Node{}
	n.X = x
	n.Y = y
	return n
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
