package main

import "math"

// Euclidean Distance for two points
// Source https://en.wikipedia.org/wiki/Euclidean_distance
func distanceEuclidean(goal, node *cell) float64 {

	tmp := math.Pow(float64(node.inx-goal.inx), 2) +
		math.Pow(float64(node.iny-goal.iny), 2)

	h := math.Sqrt(tmp)

	return h
}

// Manhattan Distance for two points
func distanceManhattan(goal, node *cell) float64 {

	h := math.Abs(float64(node.inx-goal.inx)) +
		math.Abs(float64(node.iny-goal.iny))

	return h
}

// Helper function to easily switch between heuristic distances functions
func calculateDistance(goal, node *cell) float64 {

	// return distanceEuclidean(goal, node)
	return distanceManhattan(goal, node)
}

// Calculate distances to goal for each cell
func (env *field) calculateDistances() {

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell := env.cells[iny][inx]
			cell.distance = math.MaxFloat64

			for goal, _ := range env.goals {

				distance := calculateDistance(goal, cell)

				if cell.distance > distance {
					cell.distance = distance
				}
			}
		}
	}
}

// Calculate the nearest distance to multiple goals
func (env field) calculateMinDistanceToGoals(cell *cell) float64 {

	distance := math.MaxFloat64

	for goal := range env.goals {

		value := calculateDistance(goal, cell)

		if distance > value {
			distance = value
		}
	}

	return distance
}

// Calculates/sets distances for each cell to goal cell
func (env *field) calculateDistancesPortal() {

	distancesPortals := env.portalsDistance2Goal()

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell := env.cells[iny][inx]
			distance := env.calculateMinDistanceToGoals(cell)

			// find shortest path to goal with portals as possible shortcuts
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

// Calculate distance for one portals counterpart to goal
func (env field) portalsDistance2Goal() map[*cell]float64 {

	values := make(map[*cell]float64)

	for iny := range env.cells {
		for inx := range env.cells[iny] {

			cell := env.cells[iny][inx]

			if cell.isPortal() {

				portalCell := env.getPortalCell(cell)
				values[cell] = env.calculateMinDistanceToGoals(portalCell)
			}
		}
	}

	return values
}
