package main

import (
	"math"
	"fmt"
)

var _ = fmt.Sprint("") // To be able to keep the fmt import
var SQRT2 = math.Sqrt(2)
var radToDeg = 180/math.Pi

func StraightLineDist(n1, n2 Node) float64 {
	dx := math.Abs(float64(n1.X - n2.X))
	dy := math.Abs(float64(n1.Y - n2.Y))
	return math.Sqrt(dx*dx + dy*dy)
}

func octileDistance(from, to Node) float64 {
	dx := math.Abs(float64(from.X - to.X))
    dy := math.Abs(float64(from.Y - to.Y))
    return dx + dy + (SQRT2 - 2) * math.Min(dx, dy)
}

func AStar(m [][]bool, start, goal Node) []Node {
	return shortestPath(m, start, goal, octileDistance)
}

func Dijkstra(m [][]bool, start, goal Node) []Node {
	h := func(current Node, goal Node) float64 {
		return 0
	}
	return shortestPath(m, start, goal, h)
}

// Construct the path by starting at the goal and working backwards using the parent map.
func reconstructPath(parent map[Node]Node, start, goal Node) []Node {
	node := goal
	path := []Node{}
	for node != start {
		path = append(path, node)
		var found bool
		node, found = parent[node]
		if !found {
			panic("Unexpected child node")
		}
	}
	path = append(path, start)
	// Reverse
	reversed := []Node{}
	for i := len(path)-1; i >= 0; i-- {
		reversed = append(reversed, path[i])
	}
	if reversed[0] != start {
		panic("First path node was not start")
	}
	if reversed[len(reversed)-1] != goal {
		panic("Last path node was not goal")
	}
	return reversed
}

func openNodeWithLowestF(open map[Node]float64, g map[Node]float64) Node {
	var lowestNode Node
	firstIter := true // used to initialize lowestNode
	for node, f := range(open) {
		if firstIter {
			lowestNode = node
			firstIter  = false
		} else {
			if f <= open[lowestNode] {
				if f == open[lowestNode] {
					// break this tie by favouring the node of highest g score
					if g[node] > g[lowestNode] {
						lowestNode = node
					}
				} else {
					lowestNode = node
				}
			}
		}
	}
	return lowestNode
}

func ThetaStar(m [][]bool, start, goal Node) []Node {
	h      := StraightLineDist
	g      := map[Node]float64{}
	open   := map[Node]float64{} // Map between node and f score
	parent := map[Node]Node{}

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			node := NewNode(x, y)
			g[node] = math.Inf(1)
		}
	}
	g[start]      = 0
	open[start]   = g[start] + h(start, goal)
	parent[start] = start

	for len(open) > 0 {
		node := openNodeWithLowestF(open, g)
		delete(open, node)
		if node == goal {
			return reconstructPath(parent, start, goal)
		}
		for _, neighbour := range(getTraversableNodes(m, node)) {
			par := parent[node]
			if lineOfSight(m, par, neighbour) {
				/* Path 2 */
				tentativeG := g[par] + StraightLineDist(par, neighbour)
				if tentativeG < g[neighbour] {
					g[neighbour] = tentativeG
					parent[neighbour] = par
					open[neighbour] = g[neighbour] + h(neighbour, goal)
				}
			} else {
				/* Path 1 */
				tentativeG := g[node] + StraightLineDist(node, neighbour)
				if tentativeG < g[neighbour] {
					g[neighbour] = tentativeG
					parent[neighbour] = node
					open[neighbour] = g[neighbour] + h(neighbour, goal)
				}
			}
		}
	}
	return []Node{}
}

func shortestPath(m [][]bool, start, goal Node, h func(Node, Node) float64) []Node {
	g      := map[Node]float64{}
	open   := map[Node]float64{} // Map between node and f score
	parent := map[Node]Node{}

	for y := 0; y < len(m); y++ {
		for x := 0; x < len(m[0]); x++ {
			node := NewNode(x, y)
			g[node] = math.Inf(1)
		}
	}
	g[start]      = 0
	open[start]   = g[start] + h(start, goal)
	parent[start] = start

	for len(open) > 0 {
		node := openNodeWithLowestF(open, g)
		delete(open, node)
		if node == goal {
			return reconstructPath(parent, start, goal)
		}
		for _, neighbour := range(getTraversableNodes(m, node)) {
			tentativeG := g[node] + StraightLineDist(node, neighbour)
			if tentativeG < g[neighbour] {
				parent[neighbour] = node
				g[neighbour]      = tentativeG
				open[neighbour]   = g[neighbour] + h(neighbour, goal)
			}
		}
	}
	return []Node{}
}

/*
 * Nodes outside the map are considered closed.
 */
func isOpen(m [][]bool, x, y int) bool {
	w := len(m[0])
	h := len(m)
	return x >= 0 && x < w &&
	       y >= 0 && y < h &&
	       !m[y][x]
}

func getTraversableNodes(m [][]bool, node Node) []Node {
    /*
	Traversal is done between nodes (ie grid edges) and open/closed
	spaces are entire cells. Because of this difference the figure
	below clarifies how indexing should be done.

    -------------
    | x-1 |  x  |
    |     |     |
    | y-1 | y-1 |
    ------N------
    | x-1 |  x  |
    |     |     |
    |  y  |  y  |
    -------------

	N is the given node at (x,y). Dashed lines are grid edges. The
	cells surrounding N contain their respective indices to check for
	open/closed cells in the map.
    */
	neighbours := []Node{}
	x := node.X
	y := node.Y
	nwOpen := isOpen(m, x-1, y-1)
	neOpen := isOpen(m, x,   y-1)
	seOpen := isOpen(m, x,   y)
	swOpen := isOpen(m, x-1, y)

	if nwOpen || neOpen { // We can traverse north
		neighbours = append(neighbours, NewNode(x, y-1))
	}
	if neOpen || seOpen { // We can traverse east
		neighbours = append(neighbours, NewNode(x+1, y))
	}
	if seOpen || swOpen { // We can traverse south
		neighbours = append(neighbours, NewNode(x, y+1))
	}
	if swOpen || nwOpen { // We can traverse west
		neighbours = append(neighbours, NewNode(x-1, y))
	}

	if nwOpen && (neOpen || swOpen) { // We can traverse north-west
		neighbours = append(neighbours, NewNode(x-1, y-1))
	}
	if neOpen && (nwOpen || seOpen) { // We can traverse north-east
		neighbours = append(neighbours, NewNode(x+1, y-1))
	}
	if seOpen && (neOpen || swOpen) { // We can traverse south-east
		neighbours = append(neighbours, NewNode(x+1, y+1))
	}
	if swOpen && (nwOpen || seOpen) { // We can traverse south-west
		neighbours = append(neighbours, NewNode(x-1, y+1))
	}

	return neighbours
}

// Adapted Bresenham's Line Algorithm from link below
// https://web.archive.org/web/20190717211246/http://aigamedev.com/open/tutorials/theta-star-any-angle-paths/
func lineOfSight(m [][]bool, start, end Node) bool {
	x0 := start.X
	y0 := start.Y
	x1 := end.X
	y1 := end.Y
	dx := x1-x0
	dy := y1-y0
	f := 0

	var sx, sy int
	if dx < 0 {
		dx = -dx
		sx = -1
	} else {
		sx = 1
	}
	if dy < 0 {
		dy = -dy
		sy = -1
	} else {
		sy = 1
	}

	if dx >= dy {
		for x0 != x1 {
			f += dy
			if f >= dx {
				if m[y0 + (sy-1)/2][x0 + (sx-1)/2] {
					return false
				}
				y0 += sy
				f -= dx
			}
			if f != 0 && m[y0 + (sy-1)/2][x0 + (sx-1)/2] {
				return false
			}
			if dy == 0 && m[y0][x0 + (sx-1)/2] && m[y0-1][x0 + (sx-1)/2] {
				return false
			}
			x0 += sx
		}
	} else {
		for y0 != y1 {
			f += dx
			if f >= dy {
				if m[y0 + (sy-1)/2][x0 + (sx-1)/2] {
					return false
				}
				x0 += sx
				f -= dy
			}
			if f != 0 && m[y0 + (sy-1)/2][x0 + (sx-1)/2] {
				return false
			}
			if dx == 0 && m[y0 + (sy-1)/2][x0] && m[y0 + (sy-1)/2][x0-1] {
				return false
			}
			y0 += sy
		}
	}
	return true
}

/*
 * A* with post-smoothing
 */
func AStarPs(m [][]bool, start, goal Node) []Node {
	path := AStar(m, start, goal)

	if len(path) <= 2 {
		return path
	}
	newPath := []Node{path[0]}
	firstIter := true
	for i := 0; i < len(path); i++ {
		for j := i+1; j < len(path); j++ {
			los := lineOfSight(m, path[i], path[j])
			if !los {
				if firstIter {
					panic("Assertion failed: LOS broken between the first two path nodes")
				}
				newPath = append(newPath, path[j-1])
				i = j-2
				break
			}
			firstIter = false
		}
	}
	// Add the goal
	newPath = append(newPath, path[len(path)-1])
	return newPath
}
