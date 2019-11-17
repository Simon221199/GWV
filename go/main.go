package main

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

// Cell in the field
type cell struct {
	inx      int
	iny      int
	blocked  bool
	distance float64
	symbol   string
}

// String representation of cell coordinates
func (cell cell) coordinates() string {
	return fmt.Sprintf("(%d, %d)", cell.inx, cell.iny)
}

func (cell cell) isPortal() bool {
	return isNumber(cell.symbol)
}

// Wrapper struct for environment
type field struct {
	cells       [][]*cell
	coordinates map[string]*cell
	start       *cell
	goals       map[*cell]bool
}

func (env field) getPortalCell(cell1 *cell) *cell {

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell2 := env.cells[iny][inx]

			if cell1 != cell2 && cell1.symbol == cell2.symbol {
				return cell2
			}
		}
	}

	return nil
}

// Get neighbours for a cell
func (env field) getNeighbours(node *cell) []*cell {

	if node == nil {
		return nil
	}

	neighbours := make([]*cell, 0)

	findNeighbours := func(inx, iny int) {

		// Exclude source node form neighbours
		if inx == node.inx && iny == node.iny {
			return
		}

		cell := env.cells[iny][inx]

		if cell.blocked {
			return
		}

		// Question: Is to use portal mandatory?
		neighbours = append(neighbours, cell)

		if cell.isPortal() {

			portal := env.getPortalCell(cell)
			neighbours = append(neighbours, env.getNeighbours(portal)...)
		}
	}

	// Calculate coordinate range/vector for y
	startY := max(0, node.iny-1)
	endY := min(len(env.cells), node.iny+1)

	// Calculate coordinate range/vector for x
	startX := max(0, node.inx-1)
	endX := min(len(env.cells[0]), node.inx+1)

	for iny := startY; iny <= endY; iny++ {
		findNeighbours(node.inx, iny)
	}

	for inx := startX; inx <= endX; inx++ {
		findNeighbours(inx, node.iny)
	}

	return neighbours
}

// Prints the field as matrix of distances to goal cell
func (env field) printPriorityMatrix() {

	matrix := ""

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell := env.cells[iny][inx]

			if cell.blocked {
				matrix += fmt.Sprintf("%4s", "X")
			} else {
				matrix += fmt.Sprintf("%4d", int(math.Round(cell.distance)))
			}
		}

		matrix += "\n"
	}

	fmt.Println(matrix)
}

// Prints field with path in it
func (env field) printFieldWithPath(path []*cell) {

	tmp := make(map[string]bool)
	for _, cell := range path {
		tmp[cell.coordinates()] = true
	}

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell := env.cells[iny][inx]

			if tmp[cell.coordinates()] {
				fmt.Printf("*")
			} else {
				fmt.Printf("%s", cell.symbol)
			}
		}

		fmt.Println()
	}
}

// Prints the field
func (env field) printField() {
	env.printFieldWithPath(nil)
}

func queueInfos(pq priorityQueue) (nodes, maxPathLength int) {

	for _, item := range pq {

		size := len(item.Value.cells)
		nodes += size

		if maxPathLength < size {
			maxPathLength = size
		}
	}

	return nodes, maxPathLength
}

func (env field) genericSearch(count int, priority func(path) float64) []path {

	startPath := newPath(env.start)

	pq := make(priorityQueue, 0)
	pq.Push(&item{
		Value:    startPath,
		Priority: priority(startPath),
	})
	heap.Init(&pq)

	pathsToGoal := make([]path, 0)
	expandOperations := 0

	for pq.Len() > 0 {

		popItem := heap.Pop(&pq).(*item)
		path := popItem.Value
		expandOperations++

		queueNodes, maxPathLength := queueInfos(pq)

		fmt.Printf("expanded path: %s\n", path.toString())
		env.printFieldWithPath(path.cells)
		fmt.Printf("expansion operations:   %d\n", expandOperations)
		fmt.Printf("enqueued paths:         %d\n", pq.Len())
		fmt.Printf("nodes in queue:         %d\n", queueNodes)
		fmt.Printf("longest path in queue:  %d\n", maxPathLength)
		fmt.Printf("path priority:          %.f\n", popItem.Priority)
		fmt.Println()

		lastCell := path.cells[len(path.cells)-1]

		if env.goals[lastCell] {
			pathsToGoal = append(pathsToGoal, path)

			if len(pathsToGoal) >= count {
				return pathsToGoal
			}
		}

		neighbours := env.getNeighbours(lastCell)
		// fmt.Printf("neighbours: %d\n", len(neighbours))

		for _, neighbour := range neighbours {

			if path.contains(neighbour) {
				continue
			}

			newPath := path.append(neighbour)

			// Priority is negative, because the calculateDistance queue
			// pops the highest Priority
			heap.Push(&pq, &item{
				Value:    newPath,
				Priority: priority(newPath),
			})
		}
	}

	return pathsToGoal
}

// Best-First-Search
func (env field) searchBestFirst() []path {

	// return negative value, because prioQuere picks highest value
	h := func(path path) float64 {

		last := path.cells[len(path.cells)-1]
		return -last.distance
	}

	return env.genericSearch(1, h)
}

// A* Search
func (env *field) searchAStar() []path {

	// return negative value, because prioQuere picks highest value
	h := func(path path) float64 {

		last := path.cells[len(path.cells)-1]
		return -(float64(len(path.cells)-1) + last.distance)
		// return -(float64(len(path.cells)) + last.distance)
	}

	return env.genericSearch(1, h)
}

// Breadth-First-Search
func (env *field) searchBreadthFirst() []path {

	// return negative value, because prioQuere picks highest value
	h := func(path path) float64 {

		return -float64(len(path.cells))
	}

	return env.genericSearch(1, h)
}

// Depth-First-Search
func (env *field) searchDepthFirst(count int) []path {

	// return negative value, because prioQuere picks highest value
	h := func(path path) float64 {

		return float64(len(path.cells))
	}

	return env.genericSearch(count, h)
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

	var start *cell
	goals := make(map[*cell]bool)

	for iny, line := range lines {

		fields := strings.Split(line, "")
		cells[iny] = make([]*cell, len(fields))

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

			cells[iny][inx] = cell
			coordinates[cell.coordinates()] = cell

			if str == "g" {
				goals[cell] = true
			}

			if str == "s" {
				start = cell
			}
		}
	}

	if start == nil || len(goals) == 0 {
		return nil, fmt.Errorf("error: start == nil || goal == nil")
	}

	grid := &field{
		cells:       cells,
		coordinates: coordinates,
		start:       start,
		goals:       goals,
	}

	return grid, nil
}

func main() {

	// createEnv(60, 30)
	// os.Exit(0)

	// src := "./environment/stupid.txt"
	// src := "./environment/blatt3_environment.txt"
	src := "./environment/blatt3_environment_portal.txt"
	// src := "./environment/test_env.txt"
	// src := "./environment/test_env_2.txt"
	// src := "./environment/blatt3_environment-2.txt"

	bestFirst := "best-first"
	aStar := "aStar"
	breadthFirst := "breadth-first"
	depthFirst := "depth-first"
	findKPaths := "find-NUMBER"

	search := bestFirst

	if len(os.Args) > 2 {
		search = os.Args[1]
		src = os.Args[2]
	} else {
		fmt.Printf("How to use: go run ./go [%s, %s, %s, %s, %s] PATH_TO_ENV_TXT\n",
			bestFirst, aStar, breadthFirst, depthFirst, findKPaths)
		return
	}

	fmt.Printf("sourcing  %s\n", src)

	env, err := Init(src)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Searching (%s) path form %s\n", search, env.start.coordinates())
	// env.calculateDistances()
	env.calculateDistancesPortal()
	env.printPriorityMatrix()

	var pathsToGoal []path

	if search == bestFirst {
		pathsToGoal = env.searchBestFirst()
	}

	if search == aStar {
		pathsToGoal = env.searchAStar()
	}

	if search == breadthFirst {
		pathsToGoal = env.searchBreadthFirst()
	}

	if search == depthFirst {
		pathsToGoal = env.searchDepthFirst(1)
	}

	if strings.HasPrefix(search, "find-") {
		count, err := strconv.Atoi(strings.ReplaceAll(search, "find-", ""))
		if err != nil {
			panic(err)
		}

		pathsToGoal = env.searchDepthFirst(count)
	}

	fmt.Printf("\n############ Found %d paths to goal ############\n\n", len(pathsToGoal))

	env.printField()

	for _, path := range pathsToGoal {
		fmt.Println()
		env.printFieldWithPath(path.cells)
		fmt.Printf("Path: %s\n", path.toString())
		fmt.Printf("Path length: %d\n", len(path.cells))
	}
}
