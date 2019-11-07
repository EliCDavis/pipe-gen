package main

import (
	"math"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/vector"
)

// Hopper those weird hopper things that have a hose coming in the top and a
// hose coming out the bottom
type Hopper struct {
	binHeight   float64
	taperHeight float64
	radius      float64
	position    vector.Vector3
	rotation    mesh.Quaternion
}

func (h Hopper) ToModel() mesh.Model {

	heightOffset := .1

	opening := GetClosestSpoutOpeningSize(.05)
	lipSize := opening + .1

	linkageHeight := .1
	linkageStartingRadius := h.radius * .95

	bottom := PipeSegment{
		positions: []vector.Vector3{
			h.position.Add(vector.NewVector3(0, heightOffset, 0)),      // 0
			h.position.Add(vector.NewVector3(0, heightOffset+.01, 0)),  // lip
			h.position.Add(vector.NewVector3(0, heightOffset+.1, 0)),   // lip
			h.position.Add(vector.NewVector3(0, heightOffset+.101, 0)), // opening - lip
			h.position.Add(vector.NewVector3(0, heightOffset+.201, 0)), // opening - lip
			h.position.Add(vector.NewVector3(0, heightOffset+h.taperHeight-linkageHeight, 0)),
		},
		thicknessess: []float64{
			0,
			lipSize,
			lipSize,
			lipSize - opening,
			lipSize - opening,
			linkageStartingRadius, // h.radius,
		},
	}

	topPortion := PipeSegment{
		positions: []vector.Vector3{
			h.position.Add(vector.NewVector3(0, heightOffset+h.taperHeight, 0)),
			h.position.Add(vector.NewVector3(0, heightOffset+h.taperHeight+h.binHeight, 0)),
		},
		thicknessess: []float64{
			h.radius,
			h.radius,
		},
	}

	bottomTopLinkage := GradientThickness(
		bottom.positions[len(bottom.positions)-1],
		topPortion.positions[0],
		bottom.thicknessess[len(bottom.thicknessess)-1],
		topPortion.thicknessess[0],
		10,
		func(percentDone, start, end float64) float64 {
			t :=math.Pow(percentDone-1.0, 3.0) + 1.0
			return start + ((end - start) * t)
		},
	)

	topCap := GradientThickness(
		topPortion.positions[len(topPortion.positions)-1],
		topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .5, 0)),
		topPortion.thicknessess[len(topPortion.thicknessess)-1],
		0,
		10,
		func(percentDone, start, end float64) float64 {
			t :=math.Pow(percentDone, 3.0)
			return start + ((end - start) * t)
		},
	)

	allSegments := bottom.Add(bottomTopLinkage).Add(topPortion).Add(topCap)

	mesh, _ := LineSegment3D(allSegments.positions).CreatePipeWithVarryingThickness(allSegments.thicknessess)
	return mesh
}
