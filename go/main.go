package main

import (
	"container/heap"
	"container/list"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

// Cell in the field
type cell struct {
	inx         int
	iny         int
	blocked     bool
	distance    float64
	symbol      string
}

// String representation of cell coordinates
func (cell cell) coordinates() string {
	return fmt.Sprintf("(%d, %d)", cell.inx, cell.iny)
}

func (cell cell) isPortal() bool {
	return isNumber(cell.symbol)
}

// Euclidean Distance for two points
// Source https://en.wikipedia.org/wiki/Euclidean_distance
func distanceEuclidean(goal, node *cell) float64 {

	tmp := math.Pow(float64(node.inx-goal.inx), 2) +
		math.Pow(float64(node.iny-goal.iny), 2)

	h := math.Sqrt(tmp)

	return h
}

// Manhattan Distance for two points
// Source https://en.wikipedia.org/wiki/Euclidean_distance
func distanceManhattan(goal, node *cell) float64 {

	h := math.Abs(float64(node.inx-goal.inx)) +
		math.Abs(float64(node.iny-goal.iny))

	return h
}

func calculateDistance(goal, node *cell) float64 {

	// return distanceEuclidean(goal, node)
	return distanceManhattan(goal, node)
}

// Wrapper struct for environment
type field struct {
	cells       [][]*cell
	coordinates map[string]*cell
	start       *cell
	goal        *cell
}

func (env field) getPortalCell(cell1 *cell) *cell {

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell2 := env.cells[ iny ][ inx ]

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

	findNeighbours := func (inx, iny int) {

		// Exclude source node form neighbours
		if inx == node.inx && iny == node.iny {
			return
		}

		cell := env.cells[ iny ][ inx ]

		if cell.blocked {
			return
		}

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
	endX := min(len(env.cells[ 0 ]), node.inx+1)

	for iny := startY; iny <= endY; iny++ {
		findNeighbours(node.inx, iny)
	}

	for inx := startX; inx <= endX; inx++ {
		findNeighbours(inx, node.iny)
	}

	return neighbours
}

func (env field) portalsDistance2Goal() map[*cell]float64 {

	values := make(map[*cell]float64)

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell := env.cells[ iny ][ inx ]

			if cell.isPortal() {

				portalCell := env.getPortalCell(cell)
				values[ cell ] = calculateDistance(env.goal, portalCell)
			}
		}
	}

	return values
}

// Calculates/sets distances for each cell to goal cell
func (env *field) calculateDistancesPortal() {

	distancesPortals := env.portalsDistance2Goal()

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell := env.cells[ iny ][ inx ]
			distance := calculateDistance(env.goal, cell)

			for portal, portalDistance := range distancesPortals {

				portalDistance := calculateDistance(portal, cell) + portalDistance

				if distance > portalDistance {
					distance = portalDistance
				}
			}

			cell.distance = distance
		}
	}
}

func (env *field) calculateDistances() {

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell := env.cells[ iny ][ inx ]
			cell.distance = calculateDistance(env.goal, cell)
		}
	}
}

// Prints the field as matrix of distances to goal cell
func (env field) printPriorityMatrix() {

	matrix := ""

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell := env.cells[ iny ][ inx ]

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
		tmp[ cell.coordinates() ] = true
	}

	for iny := range env.cells {
		for inx := range env.cells[ iny ] {

			cell := env.cells[ iny ][ inx ]

			if tmp[ cell.coordinates() ] {
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

// Get calculated path form ....
func (env field) getPathToGoal() []*cell {
	node := env.goal

	path := make([]*cell, 0)
	// path = append(path, env.goal)

	for node != nil {
		path = append(path, node)
		// node = node.predecessor
	}

	// path = append(path, env.start)

	// Reverse order
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path
}

// Print path from start to goal
func (env field) printPathToGoal() {

	path := env.getPathToGoal()

	fmt.Printf("Steps: %d\n", len(path))
	fmt.Printf("Coordinates:\n")

	for _, cell := range path {
		fmt.Printf("%s\n", cell.coordinates())
	}
}

// Prints field with path to goal
func (env field) printFieldWithPathToGoal() {
	env.printFieldWithPath(env.getPathToGoal())
}

// Calculate path form start to goal
// Here happens the important stuff
func (env field) searchBestFirst() []*cell {

	pq := make(priorityQueue, 0)
	pq.Push(&item{
		Value:    newPath(env.start),
		Priority: -env.start.distance,
	})
	heap.Init(&pq)

	for pq.Len() > 0 {

		popItem := heap.Pop(&pq).(*item)
		path := popItem.Value

		// fmt.Printf("cell: %s --> %1.f\n", cell.coordinates(), cell.distance)
		fmt.Print("path: ")

		for _, cell := range path.cells {
			fmt.Print(cell.coordinates() + " ")
		}
		fmt.Println()
		// fmt.Printf("queue size: %d\n", pq.Len())
		env.printFieldWithPath(path.cells)

		lastCell := path.cells[ len(path.cells) - 1 ]

		if lastCell == env.goal {
			return path.cells
		}

		neighbours := env.getNeighbours(lastCell)
		fmt.Printf("neighbours: %d\n", len(neighbours))

		for _, neighbour := range neighbours {

			if path.contains(neighbour) {
				continue
			}

			// Priority is negative, because the calculateDistance queue
			// pops the highest Priority
			heap.Push(&pq, &item{
				Value:    path.append(neighbour),
				Priority: -neighbour.distance,
			})
		}
	}

	return nil
}

func (env *field) searchBreadthFirst() {

	cellQueue := list.New()
	cellQueue.PushBack(env.start)

	for cellQueue.Len() > 0 {

		elem := cellQueue.Front()
		cellQueue.Remove(elem)

		cell := elem.Value.(*cell)

		fmt.Printf("Cell: %s\n", cell.coordinates())

		neighbours := env.getNeighbours(cell)

		for _, neighbour := range neighbours {

			// if neighbour.predecessor != nil {
			// 	continue
			// }

			if neighbour == env.start {
				continue
			}

			// neighbour.predecessor = cell
			cellQueue.PushBack(neighbour)

			if neighbour == env.goal {
				break
			}
		}
	}
}

func (env *field) searchDepthFirst() {

	cellQueue := list.New()
	cellQueue.PushFront(env.start)

	for cellQueue.Len() > 0 {

		elem := cellQueue.Front()
		cellQueue.Remove(elem)

		cell := elem.Value.(*cell)

		fmt.Printf("Cell: %s\n", cell.coordinates())

		neighbours := env.getNeighbours(cell)

		for _, neighbour := range neighbours {

			// if neighbour.predecessor != nil {
			// 	continue
			// }

			if neighbour == env.start {
				continue
			}

			// neighbour.predecessor = cell
			cellQueue.PushFront(neighbour)

			if neighbour == env.goal {
				break
			}
		}
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
	}

	if start == nil || goal == nil {
		return nil, fmt.Errorf("error: start == nil || goal == nil")
	}

	grid := &field{
		cells:       cells,
		coordinates: coordinates,
		start:       start,
		goal:        goal,
	}

	return grid, nil
}

func main() {

	// createEnv(60, 30)
	// os.Exit(0)

	// path := "./environment/stupid.txt"
	path := "./environment/blatt3_environment.txt"
	// path := "./environment/blatt3_environment_portal.txt"
	// path := "./environment/test_env.txt"
	// path := "./environment/test_env_2.txt"
	// path := "./environment/blatt3_environment-2.txt"

	if len(os.Args) > 1 {
		path = os.Args[ 1 ]
	}

	fmt.Printf("sourcing  %s\n", path)

	env, err := Init(path)
	if err != nil {
		panic(err)
	}

	fmt.Printf("######## Finding path form %s to %s\n", env.start.coordinates(), env.goal.coordinates())
	// env.calculateDistances()
	env.calculateDistancesPortal()
	env.printPriorityMatrix()
	// pathToGoal := env.searchBestFirst()
	pathToGoal := env.searchAStar()
	env.searchBreadthFirst()
	// env.searchDepthFirst()

	fmt.Printf("######## Path form %s to %s\n", env.start.coordinates(), env.goal.coordinates())
	env.printFieldWithPath(pathToGoal)
	// env.printPathToGoal()
	env.printField()
	// env.printFieldWithPathToGoal()
}
