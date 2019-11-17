package main

import (
	"math"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/meshedpotatoes/path"
	"github.com/EliCDavis/vector"
)

// Hopper those weird hopper things that have a hose coming in the top and a
// hose coming out the bottom
type Hopper struct {
	binHeight       float64
	taperHeight     float64
	radius          float64
	position        vector.Vector3
	rotation        mesh.Quaternion
	heightOffGround float64
}

func (h Hopper) GetSpouts() []Spout {
	return []Spout{
		Spout{
			opening:      GetClosestSpoutOpeningSize(.1),
			entrance:     h.position.Add(vector.NewVector3(0, h.heightOffGround+.01, 0)),
			outDirection: vector.NewVector3(0, -1, 0),
		},
	}
}

func (h Hopper) ToModel() mesh.Model {

	heightOffset := h.heightOffGround

	opening := GetClosestSpoutOpeningSize(.1)
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
			t := math.Pow(percentDone-1.0, 3.0) + 1.0
			return start + ((end - start) * t)
		},
	)

	topCap := GradientThickness(
		topPortion.positions[len(topPortion.positions)-1],
		topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .5, 0)),
		topPortion.thicknessess[len(topPortion.thicknessess)-1],
		opening,
		30,
		func(percentDone, start, end float64) float64 {
			t := math.Pow(percentDone, 3.0)
			return start + ((end - start) * t)
		},
	)

	topSpout := PipeSegment{
		positions: []vector.Vector3{
			topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .5, 0)),
			topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .51, 0)),
			topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .6, 0)),
			topPortion.positions[len(topPortion.positions)-1].Add(vector.NewVector3(0, .601, 0)),
		},
		thicknessess: []float64{
			opening,
			lipSize,
			lipSize,
			0,
		},
	}

	allSegments := bottom.
		Add(bottomTopLinkage).
		Add(topPortion).
		Add(topCap).
		Add(topSpout)

	var leg path.Path = []vector.Vector3{
		vector.Vector3Zero(),
		vector.Vector3Up().MultByConstant(heightOffset + h.taperHeight + .1),
	}
	legModel, err := leg.CreatePipe(.05, 32)

	if err != nil {
		panic(err)
	}

	hypotenuse := math.Sin(math.Pi/4.0)*h.radius - .05
	mesh, err := path.Path(allSegments.positions).CreatePipeWithVarryingThickness(allSegments.thicknessess, 32)

	if err != nil {
		panic(err)
	}

	return mesh.
		Merge(legModel.Translate(vector.NewVector3(hypotenuse, 0, hypotenuse).Add(h.position))).
		Merge(legModel.Translate(vector.NewVector3(hypotenuse, 0, -hypotenuse).Add(h.position))).
		Merge(legModel.Translate(vector.NewVector3(-hypotenuse, 0, -hypotenuse).Add(h.position))).
		Merge(legModel.Translate(vector.NewVector3(-hypotenuse, 0, hypotenuse).Add(h.position)))
}
