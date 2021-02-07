package main

import (
	"math"
	"fmt"
)

var _ = fmt.Sprint("") // To be able to keep the fmt import
var SQRT2 = math.Sqrt(2)
var radToDeg = 180/math.Pi

var grid       [][]bool
var open       map[Node]bool // These are maps just for quick membership testing
var closed     map[Node]bool // These are maps just for quick membership testing
var g          map[Node]float64
var f          map[Node]float64
var parent     map[Node]Node
var heuristic  func(Node, Node) float64
var timestamp  map[Node]int // Stores when a node had its f score updated last
var timestampCounter int

func resetPathfindingStructures() {
	open      = map[Node]bool{}
	closed    = map[Node]bool{}
	parent    = map[Node]Node{}
	timestamp = map[Node]int{} // Default value 0
	timestampCounter = 0

	heuristic = func(Node, Node) float64 {
		panic("Non-initialized heuristic function")
	}

	g = map[Node]float64{}
	f = map[Node]float64{}
	for y := 0; y <= len(grid); y++ {
		for x := 0; x <= len(grid[0]); x++ {
			node := NewNode(x, y)
			g[node] = math.Inf(1)
			f[node] = math.Inf(1)
		}
	}
}

func timestampNode(n Node) {
	timestamp[n] = timestampCounter
	timestampCounter++
}

// This function assumes that the nodes are neighbours
func costToNeighbour(n1, n2 Node) float64 {
	if n1.X != n2.X && n1.Y != n2.Y {
		return SQRT2
	} else {
		return 1
	}
}

func AStar(start, goal Node) []Node {
	resetPathfindingStructures()
	heuristic = func(from, to Node) float64 {
		dx := math.Abs(float64(from.X - to.X))
		dy := math.Abs(float64(from.Y - to.Y))
		return dx + dy + (SQRT2 - 2) * math.Min(dx, dy)
	}
	return findPath(start, goal)
}

func Dijkstra(start, goal Node) []Node {
	resetPathfindingStructures()
	heuristic = func(current Node, goal Node) float64 {
		return 0
	}
	return findPath(start, goal)
}

// Construct the path by starting at the goal and working backwards using the parent map.
func reconstructPath(start, goal Node) []Node {
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

func openNodeWithLowestF() Node {
	var lowestNode Node
	firstIter := true // used to initialize lowestNode
	for node, _ := range(open) {
		fScore, found := f[node]
		if !found {
			continue
		}

		if firstIter {
			lowestNode = node
			firstIter  = false
			continue
		}

		if fScore < f[lowestNode] {
			lowestNode = node
		} else if fScore == f[lowestNode] {
			if timestamp[lowestNode] == timestamp[node] {
				panic("Assertion error: unexpected equal timestamp")
			}
			// Break this tie by choosing the node
			// that was added last (LIFO)
			if timestamp[lowestNode] > timestamp[node] {
				lowestNode = node
			}
		}
	}
	return lowestNode
}

func findPath(start, goal Node) []Node {
	open[start] = true
	g[start]    = 0
	f[start]    = g[start] + heuristic(start, goal)

	for len(open) > 0 {
		node := openNodeWithLowestF()
		if node == goal {
			return reconstructPath(start, goal)
		}
		for _, neighbour := range(getTraversableNodes(node)) {
			_, found := closed[node]
			if found {
				continue // Closed node
			}
			tentativeG := g[node] + costToNeighbour(node, neighbour)
			if tentativeG < g[neighbour] {
				parent[neighbour] = node
				g[neighbour]      = tentativeG
				f[neighbour]      = g[neighbour] + heuristic(neighbour, goal)
				open[neighbour]   = true // Value doesn't matter
				timestampNode(neighbour)
			}
		}
		delete(open, node)
		closed[node] = true // Value doesn't matter
	}
	return []Node{}
}

func ThetaStar(start, goal Node) []Node {
	resetPathfindingStructures()
	open[start] = true
	heuristic   = StraightLineDist
	g[start]    = 0
	f[start]    = g[start] + heuristic(start, goal)

	for len(open) > 0 {
		node := openNodeWithLowestF()
		if node == goal {
			return reconstructPath(start, goal)
		}
		for _, neighbour := range(getTraversableNodes(node)) {
			_, found := closed[node]
			if found {
				continue // Closed node
			}
			par := parent[node]
			if lineOfSight(par, neighbour) {
				/* Path 2 */
				tentativeG := g[par] + StraightLineDist(par, neighbour)
				if tentativeG < g[neighbour] {
					parent[neighbour] = par
					g[neighbour]      = tentativeG
					f[neighbour]      = g[neighbour] + heuristic(neighbour, goal)
					open[neighbour]   = true // Value doesn't matter
					timestampNode(neighbour)
				}
			} else {
				/* Path 1 */
				tentativeG := g[node] + costToNeighbour(node, neighbour)
				if tentativeG < g[neighbour] {
					parent[neighbour] = node
					g[neighbour]      = tentativeG
					f[neighbour]      = g[neighbour] + heuristic(neighbour, goal)
					open[neighbour]   = true // Value doesn't matter
					timestampNode(neighbour)
				}
			}
		}
		delete(open, node)
		closed[node] = true // Value doesn't matter
	}
	return []Node{}
}

/*
 * Nodes outside the map are considered closed.
 */
func isOpen(x, y int) bool {
	w := len(grid[0])
	h := len(grid)
	return x >= 0 && x < w &&
	       y >= 0 && y < h &&
	       !grid[y][x]
}

// Small bug: If we begin in the corner of an L shape of blocked cells,
// we will only get diagonal neighbours.
// Possible solution: get neighbours based on the direction.
// The direction between node and parent[node] can easily be inferred.
func getTraversableNodes(node Node) []Node {
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
	nwOpen := isOpen(x-1, y-1)
	neOpen := isOpen(x,   y-1)
	seOpen := isOpen(x,   y)
	swOpen := isOpen(x-1, y)

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
func lineOfSight(start, end Node) bool {
	x0 := start.X
	y0 := start.Y
	x1 := end.X
	y1 := end.Y
	dx := x1-x0
	dy := y1-y0
	f  := 0

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
				if grid[y0 + (sy-1)/2][x0 + (sx-1)/2] {
					return false
				}
				y0 += sy
				f -= dx
			}
			if f != 0 && grid[y0 + (sy-1)/2][x0 + (sx-1)/2] {
				return false
			}
			if y0 == len(grid[0]) || // Otherwise we will go out of bounds
			   (dy == 0 && grid[y0][x0 + (sx-1)/2] && grid[y0-1][x0 + (sx-1)/2]) {
				return false
			}
			x0 += sx
		}
	} else {
		for y0 != y1 {
			f += dx
			if f >= dy {
				if grid[y0 + (sy-1)/2][x0 + (sx-1)/2] {
					return false
				}
				x0 += sx
				f -= dy
			}
			if f != 0 && grid[y0 + (sy-1)/2][x0 + (sx-1)/2] {
				return false
			}
			if x0 == len(grid[0]) || // Otherwise we will go out of bounds
			   (dx == 0 && grid[y0 + (sy-1)/2][x0] && grid[y0 + (sy-1)/2][x0-1]) {
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
func AStarPs(start, goal Node) []Node {
	path := AStar(start, goal)
    smoothPath := []Node{start}
    for i := 1; i < len(path)-1; i++ {
		last := smoothPath[len(smoothPath)-1]
        if !lineOfSight(last, path[i+1]) {
			smoothPath = append(smoothPath, path[i])
		}
	}
	smoothPath = append(smoothPath, goal)
    return smoothPath
}
