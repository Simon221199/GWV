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

func DistanceEuclidean(goal, node *cell) int {

	tmp := math.Pow(float64(node.inx-goal.inx), 2) +
		math.Pow(float64(node.iny-goal.iny), 2)

	h := math.Sqrt(tmp)

	return -int(math.Round(h))
}

type field struct {
	cells       [][]*cell
	coordinates map[string]*cell
	start       *cell
	goal        *cell
}

func (grid field) findNeighbours(node *cell) []*cell {

	startX := simple.Max(0, node.inx-1)
	endX := simple.Min(len(grid.cells[ 0 ]), node.inx+1)

	startY := simple.Max(0, node.iny-1)
	endY := simple.Min(len(grid.cells), node.iny+1)

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

			cell := grid.cells[ iny ][ inx ]

			if ! cell.blocked {
				neighbours = append(neighbours, cell)
			}
		}
	}

	return neighbours
}

func (grid *field) calculateDistances() {

	for iny := range grid.cells {
		for inx := range grid.cells[ iny ] {

			cell := grid.cells[ iny ][ inx ]
			// cell.priority = DistanceManhattan(goal, cell)
			cell.priority = DistanceEuclidean(grid.goal, cell)
		}
	}
}

func (grid field) priorityMatrix() string {

	matrix := ""

	for iny := range grid.cells {
		for inx := range grid.cells[ iny ] {

			cell := grid.cells[ iny ][ inx ]

			if cell.blocked {
				matrix += fmt.Sprintf("%3s", "X")
			} else {
				matrix += fmt.Sprintf("%3d", cell.priority)
			}
		}

		matrix += "\n"
	}

	return matrix
}

// func (grid field) coordinatesMap() map[string]*cell {
//
// 	hashMap := make(map[string]*cell)
//
// 	for iny := range grid.cells {
// 		for inx := range grid.cells[ iny ] {
// 			cell := grid.cells[ iny ][ inx ]
//
// 			if ! cell.blocked {
// 				hashMap[ cell.coordinates() ] = cell
// 			}
// 		}
// 	}
//
// 	return hashMap
// }

// func printPath(field field, path []*cell) {
func (grid field) printPath(path *stack.Stack) {

	// for _, point := range path {
	//
	// 	tmp := field[ point.iny ][ point.inx ].symbol + ""
	// 	// field[ point.iny ][ point.inx ].symbol = "*"
	//
	// 	for iny := range field {
	// 		for inx := range field[ iny ] {
	// 			fmt.Printf("%s", field[ iny ][ inx ].symbol)
	// 		}
	//
	// 		fmt.Println()
	// 	}
	//
	// 	field[ point.iny ][ point.inx ].symbol = tmp
	// }

	fmt.Printf("path length: %d\n", path.Len())

	for path.Len() > 0 {

		point := path.Pop().(*cell)
		// tmp := field[ point.iny ][ point.inx ].symbol + ""
		grid.cells[ point.iny ][ point.inx ].symbol = "*"

		for iny := range grid.cells {
			for inx := range grid.cells[ iny ] {
				fmt.Printf("%s", grid.cells[ iny ][ inx ].symbol)
			}

			fmt.Println()
		}

		// field[ point.iny ][ point.inx ].symbol = tmp
	}
}

func Init(path string) (*field, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	env := strings.TrimSpace(string(bytes))
	lines := strings.Split(env, "\n")

	cells := make([][]*cell, len(lines))
	coordinates := make(map[string]*cell)

	var goal *cell
	var start *cell

	for iny, line := range lines {

		fields := strings.Split(line, "")
		cells[ iny ] = make([]*cell, len(fields))

		for inx, str := range fields {

			blocked := false
			if str == "x" {
				blocked = true
			}

			cell := &cell{
				inx:     inx,
				iny:     iny,
				blocked: blocked,
				symbol:  str,
			}

			cells[ iny ][ inx ] = cell
			coordinates[ cell.coordinates() ] = cell

			if str == "g" {
				goal = cell
			}

			if str == "s" {
				start = cell
			}
		}

		fmt.Println(strings.Split(line, ""))
	}

	if start == nil || goal == nil {
		return nil, fmt.Errorf("error: start == nil || goal == nil")
	}

	grid := &field{
		cells: cells,
		coordinates: coordinates,
		start: start,
		goal:  goal,
	}

	grid.calculateDistances()

	return grid, nil
}

func main() {

	path := "/Users/patrick/Desktop/GWV/blatt3_environment.txt"

	if len(os.Args) > 2 {
		path = os.Args[ 1 ]
	}

	fmt.Printf("sourcing %s\n", path)

	grid, err := Init(path)
	if err != nil {
		panic(err)
	}

	pq := make(queue.PriorityQueue, 0)
	pq.Push(&queue.Item{
		Value:    grid.start.coordinates(),
		Priority: 1,
	})
	heap.Init(&pq)

	done := make(map[string]bool)

	// path := make([]*cell, 0)
	pathCells := stack.New()

	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*queue.Item)
		cell := grid.coordinates[ item.Value ]

		fmt.Printf("cell (%d) >> %s\n", cell.priority, cell.coordinates())

		done[ item.Value ] = true
		// path = append(path, cell)
		pathCells.Push(cell)

		if cell.coordinates() == grid.goal.coordinates() {
			fmt.Println("Done")
			break
		}

		neighbours := grid.findNeighbours(cell)

		pop := true

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

			pop = false
		}

		if pop {
			pathCells.Pop()
		}
	}

	grid.printPath(pathCells)
}
