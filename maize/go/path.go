package main

import "fmt"

type path struct {
	cells []*cell
}

func (pathObj path) contains(cell *cell) bool {

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
	}

	newPath.cells = append(newPath.cells, pathObj.cells...)
	newPath.cells = append(newPath.cells, cel)

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
		cells: make([]*cell, 0),
	}

	for _, cell := range cells {
		path.cells = append(path.cells, cell)
	}

	return path
}
