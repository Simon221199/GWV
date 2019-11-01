package main

import (
	"../queue"
	"../simple"
	"container/heap"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

type cell struct {
	inx      int
	iny      int
	blocked  bool
	priority int
	symbol   string
}

func (cell cell) coordinates() string {

	return fmt.Sprintf("(%d, %d)", cell.inx, cell.iny)
}

type grid [][]*cell

func DistanceEuclidean(goal, node *cell) int {

	tmp := math.Pow(float64(node.inx-goal.inx), 2) +
		math.Pow(float64(node.iny-goal.iny), 2)

	h := math.Sqrt(tmp)

	return -int(math.Round(h))
}

func findNeighbours(grid grid, node *cell) []*cell {

	startX := simple.Max(0, node.inx-1)
	endX := simple.Min(len(grid[ 0 ]), node.inx+1)

	startY := simple.Max(0, node.iny-1)
	endY := simple.Min(len(grid), node.iny+1)

	neighbours := make([]*cell, 0)

	// fmt.Printf("start: %s\n", node.coordinates())
	// fmt.Printf("startY: %d\n", startY)
	// fmt.Printf("endY: %d\n", endY)

	for iny := startY; iny <= endY; iny++ {
		for inx := startX; inx <= endX; inx++ {

			// fmt.Printf("(%d, %d)\n", iny, inx)

			if iny == node.iny && inx == node.inx {
				continue
			}

			cell := grid[ iny ][ inx ]

			if ! cell.blocked {
				neighbours = append(neighbours, cell)
			}
		}
	}

	return neighbours
}

func calculateDistances(grid grid, goal *cell) {

	for iny := range grid {
		for inx := range grid[ iny ] {

			cell := grid[ iny ][ inx ]
			// cell.priority = DistanceManhattan(goal, cell)
			cell.priority = DistanceEuclidean(goal, cell)
		}
	}
}

func printPriorities(grid grid) {
	for iny := range grid {
		for inx := range grid[ iny ] {

			cell := grid[ iny ][ inx ]

			if cell.blocked {
				fmt.Printf("%6s", "X")
			} else {
				fmt.Printf("%6d", cell.priority)
			}
			// fmt.Printf("(%d, %d) --> %d\n", cell.iny, cell.inx, cell.priority)
		}

		fmt.Println()
	}
}

func createHashMap(grid grid) map[string]*cell {

	hashMap := make(map[string]*cell)

	for iny := range grid {
		for inx := range grid[ iny ] {
			cell := grid[ iny ][ inx ]

			if ! cell.blocked {
				hashMap[ cell.coordinates() ] = cell
			}
		}
	}

	return hashMap
}

// func printPath(grid grid, path []*cell) {
func printPath(grid grid, path stack.Stack) {

	// for _, point := range path {
	//
	// 	tmp := grid[ point.iny ][ point.inx ].symbol + ""
	// 	// grid[ point.iny ][ point.inx ].symbol = "*"
	//
	// 	for iny := range grid {
	// 		for inx := range grid[ iny ] {
	// 			fmt.Printf("%s", grid[ iny ][ inx ].symbol)
	// 		}
	//
	// 		fmt.Println()
	// 	}
	//
	// 	grid[ point.iny ][ point.inx ].symbol = tmp
	// }

	fmt.Printf("path length: %d\n", path.Len())

	for path.Len() > 0 {

		point := path.Pop().(*cell)
		// tmp := grid[ point.iny ][ point.inx ].symbol + ""
		grid[ point.iny ][ point.inx ].symbol = "*"

		for iny := range grid {
			for inx := range grid[ iny ] {
				fmt.Printf("%s", grid[ iny ][ inx ].symbol)
			}

			fmt.Println()
		}

		// grid[ point.iny ][ point.inx ].symbol = tmp
	}
}

func main() {

	path := "/Users/patrick/Desktop/GWV/blatt3_environment.txt"

	if len(os.Args) > 2 {
		path = os.Args[ 1 ]
	}

	fmt.Printf("Sourcing %s\n", path)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	env := strings.TrimSpace(string(bytes))
	lines := strings.Split(env, "\n")

	grid := make(grid, len(lines))
	fmt.Printf("lines: %d\n", len(lines))

	var goal *cell
	var start *cell

	for iny, line := range lines {

		fields := strings.Split(line, "")
		grid[ iny ] = make([]*cell, len(fields))

		for inx, str := range fields {

			blocked := false
			if str == "x" {
				blocked = true
			}

			grid[ iny ][ inx ] = &cell{
				inx:     inx,
				iny:     iny,
				blocked: blocked,
				symbol:  str,
			}

			// fmt.Println("str=" + str)

			if str == "g" {
				// fmt.Printf("goal %d, %d\n", iny, inx)
				goal = grid[ iny ][ inx ]
			}

			if str == "s" {
				// fmt.Printf("start %d, %d\n", iny, inx)
				start = grid[ iny ][ inx ]
			}
		}

		fmt.Println(strings.Split(line, ""))
	}

	if start == nil || goal == nil {
		fmt.Println("Error: start == nil || goal == nil")
		os.Exit(1)
	}

	calculateDistances(grid, goal)
	printPriorities(grid)

	cells := createHashMap(grid)

	// neighbours := findNeighbours(grid, start)
	// fmt.Println(start.coordinates())
	// fmt.Println("#############################")
	//
	// for _, cell := range neighbours {
	// 	fmt.Println(cell.coordinates())
	// }

	pq := make(queue.PriorityQueue, 0)
	pq.Push(&queue.Item{
		Value:    start.coordinates(),
		Priority: 1,
	})
	heap.Init(&pq)

	done := make(map[string]bool)

	// path := make([]*cell, 0)
	pathCells := stack.New()

	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*queue.Item)
		cell := cells[ item.Value ]

		fmt.Printf("cell (%d) >> %s\n", cell.priority, cell.coordinates())

		done[ item.Value ] = true
		// path = append(path, cell)
		pathCells.Push(cell)

		if cell.coordinates() == goal.coordinates() {
			fmt.Println("Done")
			break
		}

		neighbours := findNeighbours(grid, cell)

		xxx := pq.Len()

		for _, neighbour := range neighbours {

			if done[ neighbour.coordinates() ] {
				continue
			}

			done[ neighbour.coordinates() ] = true

			// fmt.Printf("    neighbour (%d) >> %s\n", neighbour.priority, neighbour.coordinates())

			heap.Push(&pq, &queue.Item{
				Value:    neighbour.coordinates(),
				Priority: neighbour.priority,
			})
		}

		if xxx == pq.Len() {
			pathCells.Pop()
		}
	}

	printPath(grid, *pathCells)
}
