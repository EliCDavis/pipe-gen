package main

import (
	"github.com/EliCDavis/vector"
)

// PipeSegment used for sotring data to 3D model a pipe
type PipeSegment struct {
	positions    []vector.Vector3
	thicknessess []float64
}

// Add is not communitive, since the nature of pipes implies a direction of
// flow
func (ps PipeSegment) Add(other PipeSegment) PipeSegment {
	return PipeSegment{
		positions:    append(ps.positions, other.positions...),
		thicknessess: append(ps.thicknessess, other.thicknessess...),
	}
}
