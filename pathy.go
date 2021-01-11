package main

import (
	"fmt"
	"os"
	// "io/ioutil"
	"strings"
	"path/filepath"
)

type PathyMode int
const (
	Draw PathyMode = iota
	BenchSingle
	BenchAndDrawSingle
	BenchMultiple
	BenchAndDrawMultiple
)

type PathyParameters struct {
	Mode     PathyMode
	InPath   string
	OutPath  string
	Scale    int
	Algo     func([][]bool, Node, Node) []Node
	N        int
	Trials   int
	StartX, StartY, GoalX, GoalY int
}

var counter = 0
func readNextArg() string {
	arg := os.Args[counter]
	counter++
	return arg 
}

func main() {
	// Print help
	if len(os.Args) < 2 {
		fmt.Printf("%s is a tool for visualization and benchmarking of pathfinding algorithms.\n\n", os.Args[0])
		fmt.Println("To draw a map:")
		fmt.Printf("    %s draw map_file output_jpg scale\n", os.Args[0])
		fmt.Println("To benchmark one scenario:")
		fmt.Printf("    %s single map_file start_x start_y goal_x goal_y algorithm trials\n", os.Args[0])
		fmt.Println("To benchmark one scenario and draw its path:")
		fmt.Printf("    %s single map_file start_x start_y goal_x goal_y algorithm trials output_jpg scale\n", os.Args[0])
		fmt.Println("To benchmark multiple scenarios:")
		fmt.Printf("    %s multiple scenarios_file algorithm n trials\n", os.Args[0])
		fmt.Println("To benchmark multiple scenarios and draw their paths:")
		fmt.Printf("    %s multiple scenarios_file algorithm n trials output_dir scale\n\n", os.Args[0])
		fmt.Println("Accepted algorithms are \"dijkstra\", \"astar\", \"astar-ps\" and \"thetastar\". N is the amount of scenarios to pick from the file. They are evenly spread out in terms of problem size.")
		os.Exit(0)
	}

	readNextArg() // Skip program name
	modeString := readNextArg()
	var p PathyParameters

	// Read command-line arguments
	switch (strings.ToLower(modeString)) {
		case "draw":
			p = getDrawModeParameters()
		case "single":
			p = getSingleModeParameters()
		case "multiple":
			p = getMultipleModeParameters()
		default:
			fmt.Printf("Unknown mode \"%s\", accepted modes are \"draw\", \"single\" and \"multiple\"\n", modeString)
			os.Exit(1)
	}

	// Check some of the arguments
	if (p.Mode == Draw || p.Mode == BenchAndDrawSingle || p.Mode == BenchAndDrawMultiple) && p.Scale < 1 {
		fmt.Println("Scale must be a positive integer.")
		os.Exit(1)
	}
	if (p.Mode == BenchMultiple || p.Mode == BenchAndDrawMultiple) && p.N < 1 {
		fmt.Println("N must be a positive integer.")
		os.Exit(1)
	}
	if (p.Mode == BenchSingle || p.Mode == BenchAndDrawSingle || p.Mode == BenchMultiple || p.Mode == BenchAndDrawMultiple) && p.Trials < 1 {
		fmt.Println("Trials must be a positive integer.")
		os.Exit(1)
	}

	// Run the appropriate mode
	switch (p.Mode) {
		case Draw:
			runDrawMode(p)
		case BenchSingle, BenchAndDrawSingle:
			runSingleMode(p)
		case BenchMultiple, BenchAndDrawMultiple:
			runMultipleMode(p)
		default:
			panic("Assertion failed: unexpected mode")
	}

	fmt.Println("Success")
}

func getDrawModeParameters() PathyParameters {
	if len(os.Args) != 5 {
		fmt.Printf("Wrong number of arguments. Run %s without parameters for more info.\n", os.Args[0])
		os.Exit(1)
	}
	p := PathyParameters{}
	p.Mode    = Draw
	p.InPath  = readNextArg()
	p.OutPath = readNextArg()
	p.Scale   = MustParseInt(readNextArg())
	return p
}

func getSingleModeParameters() PathyParameters {
	if len(os.Args) != 9 && len(os.Args) != 11 {
		fmt.Printf("Wrong number of arguments. Run %s without parameters for more info.\n", os.Args[0])
		os.Exit(1)
	}
	p := PathyParameters{}
	p.InPath = readNextArg()
	p.StartX = MustParseInt(readNextArg())
	p.StartY = MustParseInt(readNextArg())
	p.GoalX  = MustParseInt(readNextArg())
	p.GoalY  = MustParseInt(readNextArg())
	p.Algo   = MustParsePathfindingFunction(readNextArg())
	p.Trials = MustParseInt(readNextArg())
	if len(os.Args) == 11 {
		p.Mode    = BenchAndDrawSingle
		p.OutPath = readNextArg()
		p.Scale   = MustParseInt(readNextArg())
	} else {
		p.Mode = BenchSingle
	}
	return p
}

func getMultipleModeParameters() PathyParameters {
	if len(os.Args) != 6 && len(os.Args) != 8 {
		fmt.Printf("Wrong number of arguments. Run %s without parameters for more info.\n", os.Args[0])
		os.Exit(1)
	}
	p := PathyParameters{}
	p.InPath = readNextArg()
	p.Algo   = MustParsePathfindingFunction(readNextArg())
	p.N      = MustParseInt(readNextArg())
	p.Trials = MustParseInt(readNextArg())
	if len(os.Args) == 8 {
		p.Mode    = BenchAndDrawMultiple
		p.OutPath = readNextArg()
		p.Scale   = MustParseInt(readNextArg())
	} else {
		p.Mode = BenchMultiple
	}
	return p
}

func runDrawMode(p PathyParameters) {
	if p.Mode != Draw {
		panic("Assertion failed: unexpected mode")
	}
	m, err := LoadMap(p.InPath)
	if err != nil {
		fmt.Printf("Error reading file \"%s\": %s\n", p.InPath, err.Error())
		os.Exit(1)
	}
	img := MakeMapImage(m, p.Scale)
	err  = SaveImage(img, p.OutPath)
	if err != nil {
		fmt.Printf("Error writing image \"%s\": %s\n", p.OutPath, err.Error())
		os.Exit(1)
	}
}

func runSingleMode(p PathyParameters) {
	if p.Mode != BenchSingle && p.Mode != BenchAndDrawSingle {
		panic("Assertion failed: unexpected mode")
	}
	m, err := LoadMap(p.InPath)
	if err != nil {
		fmt.Printf("Error reading file \"%s\": %s\n", p.InPath, err.Error())
		os.Exit(1)
	}

	start := NewNode(p.StartX, p.StartY)
	goal  := NewNode(p.GoalX,  p.GoalY)
	path, turns, pathLen, avgAngle, avgRuntime := testOneScenario(start, goal, m, p.Algo, p.Trials)
	fmt.Printf("Stats: %d turn(s), length %.1f, avg angle %.1f rad (%.1f deg), runtime %dms\n", turns, pathLen, avgAngle, avgAngle*radToDeg, avgRuntime)

	if p.Mode == BenchAndDrawSingle {
		img := MakeMapImage(m, p.Scale)
		img  = DrawPath(img, path, p.Scale)
		err  = SaveImage(img, p.OutPath)
		if err != nil {
			fmt.Printf("Error writing image \"%s\": %s\n", p.OutPath, err.Error())
			os.Exit(1)
		}
	}
}

func runMultipleMode(p PathyParameters) {
	if p.Mode != BenchMultiple && p.Mode != BenchAndDrawMultiple {
		panic("Assertion failed: unexpected mode")
	}

	// Load scenarios
	scenarios, err := LoadScenarios(p.InPath)
	if err != nil {
		fmt.Printf("Error loading scenarios file \"%s\": %s\n", p.InPath, err.Error())
		os.Exit(1)
	}
	// Load map
	mapPath := filepath.Join(filepath.Dir(p.InPath), scenarios[0].MapName)
	m, err := LoadMap(mapPath)
	if err != nil {
		fmt.Printf("Error reading map file \"%s\": %s\n", mapPath, err.Error())
		os.Exit(1)
	}

	// If needed, create an output directory for images
	if p.Mode == BenchAndDrawMultiple {
		// Create the output directory if it doesn't exist
		_, err := os.Stat(p.OutPath)
		if os.IsNotExist(err) {
			err = os.Mkdir(p.OutPath, os.ModeDir)
			if err != nil {
				fmt.Printf("Error creating output directory: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}

	// Select n evenly spread out scenarios
	selectedScenarios := []Scenario{}
	var inc float64
	if p.N >= len(scenarios) {
		inc = 1
		p.N = len(scenarios)
	} else {
		inc = float64(len(scenarios)-1) / float64(p.N-1)
	}
	for i := 0.0; i < float64(len(scenarios)); i += inc {
		index := int(i)
		selectedScenarios = append(selectedScenarios, scenarios[index])
	}
	// Assertion
	if len(selectedScenarios) != p.N {
		panic("Assertion failed: unexpected number of selected scenarios")
	}

	// Benchmark and draw scenarios
	sumTurnCount  := 0.0
	sumPathLen    := 0.0
	sumAvgAngle   := 0.0
	sumAvgRuntime := 0
	for _, scenario := range selectedScenarios {
		// Assertion
		if scenario.MapName != scenarios[0].MapName {
			panic("Assertion failed: scenarios file referred to multiple map files")
		}
		sx, sy, gx, gy := scenario.Start.X, scenario.Start.Y, scenario.Goal.X, scenario.Goal.Y
		start := NewNode(sx,sy)
		goal  := NewNode(gx,gy)
		path, turns, pathLen, avgAngle, avgRuntime := testOneScenario(start, goal, m, p.Algo, p.Trials)
		fmt.Printf("(%d,%d) -> (%d,%d) stats: %d turn(s), length %.1f, avg angle %.1f rad (%.1f deg), runtime %dms\n", sx, sy, gx, gy, turns, pathLen, avgAngle, avgAngle*radToDeg, avgRuntime)

		sumTurnCount  += float64(turns)
		sumPathLen    += pathLen
		sumAvgAngle   += avgAngle
		sumAvgRuntime += avgRuntime

		if p.Mode == BenchAndDrawMultiple {
			// Create a nice name for the image
			ext   := filepath.Ext(scenario.MapName)
			fname := scenario.MapName[0:len(scenario.MapName)-len(ext)]
			fname  = fmt.Sprintf("%s_%d_%d_%d_%d.jpg", fname, sx, sy, gx, gy)
			out   := filepath.Join(p.OutPath, fname)

			img := MakeMapImage(m, p.Scale)
			img  = DrawPath(img, path, p.Scale)
			err  = SaveImage(img, out)
			if err != nil {
				fmt.Printf("Error writing image \"%s\": %s\n", out, err.Error())
				os.Exit(1)
			}
		}
	}

	// Stats are the average across all selected scenarios
	overallTurnCount  := sumTurnCount  / float64(p.N)
	overallPathLen    := sumPathLen    / float64(p.N)
	overallAvgAngle   := sumAvgAngle   / float64(p.N)
	overallAvgRuntime := sumAvgRuntime / p.N
	fmt.Printf("\nAvg stats: %f turn(s), length %f, avg angle %f rad (%.1f deg), runtime %dms\n", overallTurnCount, overallPathLen, overallAvgAngle, overallAvgAngle*radToDeg, overallAvgRuntime)
}

/*func oldmain() {
	if len(os.Args) != 5 {
		fmt.Println("Wrong number of arguments!")
		fmt.Printf("Usage: %s INPUT ALGORITHM N TRIALS\n", os.Args[0])
		fmt.Println("Algorithms: \"dijkstra\" \"astar\" \"astar-ps\" \"thetastar\"")
		os.Exit(1)
	}

	scenariosPath := os.Args[1]
	algoName := os.Args[2]
	algo     := MustParsePathfindingFunction(algoName)
	n        := MustParseInt(os.Args[3]) // Use this to select every nth scenario per file
	if n <= 0 {
		fmt.Println("Argument n must be a positive integer.")
		os.Exit(1)
	}
	trials := MustParseInt(os.Args[4]) // Number of times a scenario is to be repeated
	if trials <= 0 {
		fmt.Println("Trials must be a positive integer.")
		os.Exit(1)
	}

	// Load and store scenarios
	scenarios, err := LoadScenarios(scenariosPath)
	if err != nil {
		msg := fmt.Sprintf("Error loading scenarios file \"%s\": %s", scenariosPath, err.Error())
		fmt.Println(msg)
		os.Exit(1)
	}
	// Select n evenly spread out scenarios
	{
		inc := 1.0
		expectedIterations := len(scenarios)
		if n <= 1 {
			inc = float64(len(scenarios))
			expectedIterations = 1
		} else if n < len(scenarios) {
			inc = float64(len(scenarios)-1) / float64(n-1)
			expectedIterations = n
		}
		iterations := 0
		for i := 0.0; i < float64(len(scenarios)); i += inc {
			index := int(i)
			loadedScenarios = append(loadedScenarios, scenarios[index])
			iterations++
		}
		// Assertions
		if iterations != expectedIterations {
			msg := fmt.Sprintf("Unexpected number of picked scenarios (expected %d, got %d)\n", expectedIterations, iterations)
			panic(msg)
		}
	}
	fmt.Printf("Loaded %d scenario(s).\n", len(loadedScenarios))
	if len(loadedScenarios) == 0 {
		fmt.Println("Done.")
		os.Exit(0)
	}

	// Execute scenarios, loading new maps if the current scenario calls for it
	var mapName string
	m := [][]bool{}
	sumTurnCount  := 0.0
	sumPathLen    := 0.0
	sumAvgAngle   := 0.0
	sumAvgRuntime := 0
	for _, scenario := range loadedScenarios {
		scenariosPath := filepath.Join(path, scenario.Path)

		// Load the map if it's new
		if mapName != scenario.MapName {
			mapName = scenario.MapName
			fmt.Printf("Loading map file \"%s\"...\n", mapName)
			mapPath := filepath.Join(path, mapName)
			var err error
			m, err = LoadMap(mapPath)
			if err != nil {
				msg := fmt.Sprintf("Error loading map file \"%s\": %s", mapPath, err.Error())	
				fmt.Println(msg)
				os.Exit(1)
			}
		}

		// Assert that the map size of the scenarios file and its corresponding map match
		if len(m) != scenario.Height || len(m[0]) != scenario.Width {
			msg := fmt.Sprintf("Assertion failed: scenarios file \"%s\" refers to a different sized map", scenariosPath)
			panic(msg)
		}

		turns, pathLen, avgAngle, avgRuntime := testOneScenario(scenario.Start, scenario.Goal, m, algo, trials)
		sumTurnCount  += float64(turns)
		sumPathLen    += pathLen
		sumAvgAngle   += avgAngle
		sumAvgRuntime += avgRuntime
		
		fmt.Printf("    %s:%d (%d,%d) -> (%d,%d) tested.\n", scenariosPath, scenario.Bucket, scenario.Start.X, scenario.Start.Y, scenario.Goal.X, scenario.Goal.Y)

		if true {
			path := algo(m, scenario.Start, scenario.Goal)
			img  := DrawPath(MakeMapImage(m), path)
			out  := fmt.Sprintf("%s_%d_%d_%d_%d.jpg", strings.Split(scenario.MapName, ".")[0], scenario.Start.X, scenario.Start.Y, scenario.Goal.X, scenario.Goal.Y)
			err  := SaveImage(img, out)
			if err != nil {
				fmt.Printf("Error writing image \"%s\": %s\n", out, err.Error())
				os.Exit(1)
			}
		}

	}
	fmt.Printf("Ran %d scenario(s).\n", len(loadedScenarios))

	// Stats are the average across all selected scenarios
	testCount := float64(len(loadedScenarios))
	overallTurnCount  := sumTurnCount  / testCount
	overallPathLen    := sumPathLen    / testCount
	overallAvgAngle   := sumAvgAngle   / testCount
	overallAvgRuntime := sumAvgRuntime / int(testCount)
	fmt.Printf("Average stats: %f turn(s), length %f, average angle %f rad (%.1f deg), runtime %dms\n", overallTurnCount, overallPathLen, overallAvgAngle, overallAvgAngle*radToDeg, overallAvgRuntime)

	fmt.Println("Done.")
}*/
