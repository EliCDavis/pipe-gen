package main

import (
	"github.com/EliCDavis/vector"
)

func GetClosestSpoutOpeningSize(size float64) float64 {
	return size
}

func GradientThickness(
	startingPoint, endingPoint vector.Vector3,
	startingWidth, endingWidth float64,
	steps int,
	f func(percentDone, start, end float64) float64,
) PipeSegment {
	points := make([]vector.Vector3, steps-2)

	thicknesses := make([]float64, steps-2)

	dir := endingPoint.Sub(startingPoint)

	for i := 1; i < steps-1; i++ {
		percentDone := float64(i) / float64(steps)
		thicknesses[i-1] = f(percentDone, startingWidth, endingWidth)
		points[i-1] = startingPoint.Add(dir.MultByConstant(percentDone))
	}

	return PipeSegment{positions: points, thicknessess: thicknesses}
}
