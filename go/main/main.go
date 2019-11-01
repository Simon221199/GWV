package main

import (
	"../queue"
	"../simple"
	"container/heap"
	"fmt"
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
	predecessor *cell
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

	for iny := startY; iny <= endY; iny++ {
		for inx := startX; inx <= endX; inx++ {
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

func (grid field) printFieldPath(path []*cell) {

	for _, cell := range path {
		cell.symbol = "*"
	}

	for iny := range grid.cells {
		for inx := range grid.cells[ iny ] {
			fmt.Printf("%s", grid.cells[ iny ][ inx ].symbol)
		}

		fmt.Println()
	}
}

func (grid field) printPathFromGoal() {

	steps := 0
	node := grid.goal

	path := make([]*cell, 0)

	for node != nil {
		fmt.Println(node.coordinates())
		path = append(path, node)
		node = node.predecessor
		steps++
	}

	fmt.Printf("Steps: %d\n", steps)

	grid.printFieldPath(path)
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

	// path := "/Users/patrick/Desktop/GWV/blatt3_environment.txt"
	path := "/Users/patrick/Desktop/GWV/blatt3_environment-2.txt"

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

	for pq.Len() > 0 {

		item := heap.Pop(&pq).(*queue.Item)
		cell := grid.coordinates[ item.Value ]

		fmt.Printf("cell (%d) >> %s\n", cell.priority, cell.coordinates())

		done[ item.Value ] = true

		if cell.coordinates() == grid.goal.coordinates() {
			fmt.Println("Done")
			break
		}

		neighbours := grid.findNeighbours(cell)

		for _, neighbour := range neighbours {

			if done[ neighbour.coordinates() ] {
				continue
			}

			done[ neighbour.coordinates() ] = true

			neighbour.predecessor = cell
			// fmt.Printf("    neighbour (%d) >> %s\n", neighbour.priority, neighbour.coordinates())

			heap.Push(&pq, &queue.Item{
				Value:    neighbour.coordinates(),
				Priority: neighbour.priority,
			})
		}
	}

	// grid.printPath(pathCells)
	grid.printPathFromGoal()
}
