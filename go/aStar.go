package main

import (
	"container/heap"
	"fmt"
)

type path struct {
	cells    []*cell
	// cellsMap map[*cell]bool
}

func (pathObj path) contains(cell *cell) bool {
	// return pathObj.cellsMap[ cell ]

	for _, elem := range pathObj.cells {
		if elem.coordinates() == cell.coordinates() {
			return true
		}
	}

	return false
}

func (pathObj path) append(cel *cell) path {

	newPath := path{
		cells: make([]*cell, 0),
		// cells:    append(pathObj.cells, cell),
		// cellsMap: pathObj.cellsMap,
	}

	newPath.cells = append(newPath.cells, pathObj.cells...)
	newPath.cells = append(newPath.cells, cel)
	// newPath.cells = append(newPath.cells, cell)
	// newPath.cellsMap[ cell ] = true

	return newPath
}

func newPath(cells ...*cell) path {

	path := path{
		cells:    make([]*cell, 0),
		// cellsMap: make(map[*cell]bool),
	}

	for _, cell := range cells {
		path.cells = append(path.cells, cell)
		// path.cellsMap[ cell ] = true
	}

	return path
}

func (env *field) searchAStar() []*cell {

	pq := make(priorityQueue, 0)
	pq.Push(&item{
		// Value:    newPath(env.start, env.cells[4][5], env.cells[3][5], env.cells[2][5], env.cells[2][6], env.cells[1][6], env.cells[1][7]),
		Value:    newPath(env.start),
		Priority: -env.start.distance,
	})
	heap.Init(&pq)

	for pq.Len() > 0 {

		popItem := heap.Pop(&pq).(*item)
		path := popItem.Value

		// fmt.Printf("cell: %s --> %1.f\n", cell.coordinates(), cell.distance)
		fmt.Printf("priority: %.0f\n", popItem.Priority)
		fmt.Print("path: ")

		for _, cell := range path.cells {
			fmt.Print(cell.coordinates() + " ")
		}
		fmt.Println()

		env.printFieldWithPath(path.cells)

		lastCell := path.cells[ len(path.cells) - 1 ]

		if lastCell == env.goal {
			return path.cells
		}

		neighbours := env.getNeighbours(lastCell)
		fmt.Printf("neighbours: %d\n", len(neighbours))
		fmt.Printf("queue size: %d\n", pq.Len())

		for _, neighbour := range neighbours {

			if path.contains(neighbour) {
				continue
			}

			// Priority is positive, because the calculateDistance queue
			// pops the highest Priority
			heap.Push(&pq, &item{
				Value:    path.append(neighbour),
				Priority: -(float64(len(path.cells) - 1) + neighbour.distance),
			})
		}

		fmt.Printf("queue size: %d\n", pq.Len())
	}

	return nil
}
