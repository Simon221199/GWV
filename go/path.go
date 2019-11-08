package main

import "fmt"

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

func (pathObj path) toString() string {

	str := ""

	for _, cell := range pathObj.cells {
		str += fmt.Sprint(cell.coordinates() + " ")
	}

	return str
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
