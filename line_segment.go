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
func (ls LineSegment3D) CreatePipe() (mesh.Model, error) {
	if len(ls) < 2 {
		return mesh.Model{}, errors.New("Unable to create a pipe with less than 2 points")
	}

	return mesh.Model{}, nil
}
