# Pathy

![](./maps_with_paths.png)

A tool for visualization and benchmarking of pathfinding algorithms (Dijkstra, A, Post-Smoothed A* and Theta*).
The program operates on map and scenarios files from [movingai.com/benchmarks/grids.html](https://www.movingai.com/benchmarks/grids.html), some of which are available in the `maps` directory.

## Building

Run `go build pathy.go data.go pathfinding.go mapimage.go loader.go` in the `code` directory.
[draw2d](https://godoc.org/github.com/llgcode/draw2d) is required to build this project.

## Using the CLI

Run `pathy` without parameters to view the available commands.

Drawing an image based on a map file where each cell is 16x16 pixels: `pathy draw mapfile.map image.jpg 16`

Benchmarking Dijkstra using start and goal coordinates in 10 trials: `pathy single mapfile.map 5 5 100 250 dijkstra 10`

Benchmarking Post-Smoothed A* in 5 scenarios in 10 trials: `pathy multiple scenariosfile.scen astar 5 10`

## Licenses

The files under the `maps` directory are under the Open Data Commons Attribution License.
The rest of the repository is under the MIT license.
