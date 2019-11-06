package main

import (
	"errors"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/vector"
)

// LineSegment3D is a series of ordered points that make up a line segment
// through 3D space.
type LineSegment3D []vector.Vector3

// CreatePipe draws a pipe using the line segment as guides
func (ls LineSegment3D) CreatePipe(pipeThickness float64) (mesh.Model, error) {
	if len(ls) < 2 {
		return mesh.Model{}, errors.New("Unable to create a pipe with less than 2 points")
	}

	points := make([][]vector.Vector3, len(ls))
	for i, p := range ls {
		var dir vector.Vector3

		if i == 0 {
			dir = ls[1].Sub(ls[0])
		} else {
			dir = ls[i].Sub(ls[i-1])
		}

		points[i] = GetPlaneOuterPoints(p, dir, pipeThickness, 64)
	}

	polygons := make([]mesh.Polygon, 0)

	for i := 1; i < len(points); i++ {
		for p := 0; p < len(points[i]); p++ {
			polygons = append(polygons, MakeSquare(
				points[i-1][p],
				points[i][p],
				points[i][(p+1)%len(points[i])],
				points[i-1][(p+1)%len(points[i])],
			)...)
		}
	}

	return mesh.NewModel(polygons)
}
