package main

import (
	"errors"
	"log"
	"math"

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

// CreatePipeWithVarryingThickness draws a pipe using the line segment as guides
func (ls LineSegment3D) CreatePipeWithVarryingThickness(thicknesses []float64) (mesh.Model, error) {
	if len(ls) < 2 {
		return mesh.Model{}, errors.New("Unable to create a pipe with less than 2 points")
	}

	if len(ls) != len(thicknesses) {
		return mesh.Model{}, errors.New("We need a thickness value per point on the line")
	}

	points := make([][]vector.Vector3, len(ls))
	for i, p := range ls {
		var dir vector.Vector3

		if i == 0 {
			dir = ls[1].Sub(ls[0])
		} else {
			dir = ls[i].Sub(ls[i-1])
		}

		points[i] = GetPlaneOuterPoints(p, dir, thicknesses[i], 64)
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

// CreateLoopingPlatform creates a walkway that will close in on itself to create an
// enclosed track.
func (ls LineSegment3D) CreateLoopingPlatform(width float64) (mesh.Model, error) {
	if len(ls) < 3 {
		return mesh.Model{}, errors.New("Unable to create a looping platform with less than 3 points")
	}

	leftSidesDirty := make([]LineSegment3D, 0)
	rightSidesDirty := make([]LineSegment3D, 0)
	for i := 0; i < len(ls); i++ {
		previous := i - 1
		if previous < 0 {
			previous = len(ls) - 1
		}
		dir := ls[i].Sub(ls[previous])

		rightSide := dir.Cross(vector.Vector3Up()).Normalized().MultByConstant(width / 2.0)

		leftSidesDirty = append(leftSidesDirty, []vector.Vector3{
			ls[previous].Sub(rightSide),
			ls[i].Sub(rightSide),
		})

		rightSidesDirty = append(rightSidesDirty, []vector.Vector3{
			ls[previous].Add(rightSide),
			ls[i].Add(rightSide),
		})
	}

	polys := make([]mesh.Polygon, 0)
	cleanLeftSide, err := properlyAlignSegments(leftSidesDirty)
	if err != nil {
		return mesh.Model{}, err
	}

	cleanRightSide, err := properlyAlignSegments(rightSidesDirty)
	if err != nil {
		return mesh.Model{}, err
	}

	railing := mesh.Model{}
	for i := 0; i < len(cleanLeftSide); i++ {
		polys = append(polys, MakeSquare(
			cleanLeftSide[i][0],
			cleanLeftSide[i][1],
			cleanRightSide[i][1],
			cleanRightSide[i][0],
		)...)

		leftRailing, _ := LineSegment3D([]vector.Vector3{
			cleanLeftSide[i][0].Add(vector.Vector3Up()),
			cleanLeftSide[i][1].Add(vector.Vector3Up()),
		}).CreatePipe(.05)

		rightRailing, _ := LineSegment3D([]vector.Vector3{
			cleanRightSide[i][0].Add(vector.Vector3Up()),
			cleanRightSide[i][1].Add(vector.Vector3Up()),
		}).CreatePipe(.05)

		railing = railing.
			Merge(leftRailing).
			Merge(rightRailing).
			Merge(leftRailing.Translate(vector.Vector3Up().MultByConstant(-.5))).
			Merge(rightRailing.Translate(vector.Vector3Up().MultByConstant(-.5))).
			Merge(railingPostsBetweenTwoPoints(cleanLeftSide[i][0], cleanLeftSide[i][1])).
			Merge(railingPostsBetweenTwoPoints(cleanRightSide[i][0], cleanRightSide[i][1]))
	}

	pathwayMesh, err := mesh.NewModel(polys)

	return pathwayMesh.Merge(railing), nil
}

// Can be optimized
func railingPostsBetweenTwoPoints(start vector.Vector3, end vector.Vector3) mesh.Model {

	maxDistBetweenRailing := 2.0

	numOfPosts := 2

	for end.Distance(start)/float64(numOfPosts-1) > maxDistBetweenRailing {
		numOfPosts++
	}

	distBetweenPoints := end.Distance(start) / float64(numOfPosts-1)
	dir := end.Sub(start).Normalized()

	out := mesh.Model{}

	for i := 0; i < numOfPosts; i++ {
		out = out.Merge(railingPost(start.Add(dir.MultByConstant(distBetweenPoints * float64(i)))))
	}

	return out
}

func railingPost(position vector.Vector3) mesh.Model {
	pipe, err := LineSegment3D([]vector.Vector3{
		position,
		position.Add(vector.Vector3Up()),
	}).CreatePipe(.05)
	if err != nil {
		panic(err)
	}
	return pipe
}

// YIntersection assumes line segment 3d is only two points
// Uses the paremetric equations of the line
func (ls LineSegment3D) YIntersection(x float64, z float64) float64 {
	v := ls[1].Sub(ls[0])
	t := (x - ls[0].X()) / v.X()

	// This would mean that the X direction is 0 (x never changes) so we'll
	// have to figure out where we are on the line using the z axis.
	if math.IsNaN(t) {
		t = (z - ls[0].Z()) / v.Z()
	}

	// Well then uh... return y slope I guess? Maybe I should throw NaN?
	// Ima throw a NaN
	if math.IsNaN(t) {
		log.Printf("getting nan")
		return math.NaN()
	}

	return ls[0].Y() + (v.Y() * t)
}

// ScaleOutwards assumes line segment 3d is only two points
// multiplies the current length of the line by extending it out
// further in the two different directions it's heading
func (ls LineSegment3D) ScaleOutwards(amount float64) LineSegment3D {
	dirAndMag := ls[1].Sub(ls[0]).DivByConstant(2.0)
	center := dirAndMag.Add(ls[0])
	return []vector.Vector3{
		center.Add(dirAndMag.MultByConstant(amount)),
		center.Add(dirAndMag.MultByConstant(-amount)),
	}
}

// Can be optimized
func properlyAlignSegments(dirty []LineSegment3D) ([]LineSegment3D, error) {

	clean := make([]LineSegment3D, len(dirty))

	for i := 0; i < len(clean); i++ {
		previous := i - 1
		next := i + 1
		if previous < 0 {
			previous = len(dirty) - 1
		}

		if next >= len(dirty) {
			next = 0
		}

		scaledLastLine := dirty[previous].ScaleOutwards(100000)
		lastLine2D := mesh.NewLine(
			vector.NewVector2(scaledLastLine[0].X(), scaledLastLine[0].Z()),
			vector.NewVector2(scaledLastLine[1].X(), scaledLastLine[1].Z()),
		)

		scaledCurLine := dirty[i].ScaleOutwards(100000)
		curLine2D := mesh.NewLine(
			vector.NewVector2(scaledCurLine[0].X(), scaledCurLine[0].Z()),
			vector.NewVector2(scaledCurLine[1].X(), scaledCurLine[1].Z()),
		)

		scaledNextLine := dirty[next].ScaleOutwards(100000)
		nextLine2D := mesh.NewLine(
			vector.NewVector2(scaledNextLine[0].X(), scaledNextLine[0].Z()),
			vector.NewVector2(scaledNextLine[1].X(), scaledNextLine[1].Z()),
		)

		// I need to do something special if  these lines are going in the
		// same direction...
		samePreviousDirection := lastLine2D.Dir().
			Normalized().
			Sub(curLine2D.Dir().Normalized())

		sameNextDirection := curLine2D.Dir().
			Normalized().
			Sub(nextLine2D.Dir().Normalized())

		if samePreviousDirection.X() == 0 && samePreviousDirection.Y() == 0 {
			if clean[i] == nil {
				clean[i] = make([]vector.Vector3, 2)
			}

			clean[i][0] = vector.NewVector3(
				dirty[previous][1].X(),
				math.Min(dirty[previous][1].Y(), dirty[i][0].Y()),
				dirty[previous][1].Z(),
			)

		} else if lastLine2D.Intersects(curLine2D) {
			intersection, _ := lastLine2D.Intersection(curLine2D)
			if clean[i] == nil {
				clean[i] = make([]vector.Vector3, 2)

			}
			clean[i][0] = vector.NewVector3(
				intersection.X(),
				math.Min(
					scaledLastLine.YIntersection(intersection.X(), intersection.Y()),
					scaledCurLine.YIntersection(intersection.X(), intersection.Y()),
				),
				intersection.Y(),
			)
		} else {
			return []LineSegment3D{}, errors.New("Unable to achieve intersection before to current. Guess I'm the big stupid")
		}

		if sameNextDirection.X() == 0 && sameNextDirection.Y() == 0 {
			if clean[i] == nil {
				clean[i] = make([]vector.Vector3, 2)
			}

			clean[i][1] = vector.NewVector3(
				dirty[i][1].X(),
				math.Min(dirty[i][1].Y(), dirty[next][0].Y()),
				dirty[i][1].Z(),
			)

		} else if curLine2D.Intersects(nextLine2D) {
			intersection, _ := curLine2D.Intersection(nextLine2D)
			if clean[i] == nil {
				clean[i] = make([]vector.Vector3, 2)
			}
			clean[i][1] = vector.NewVector3(
				intersection.X(),
				math.Min(
					scaledNextLine.YIntersection(intersection.X(), intersection.Y()),
					scaledCurLine.YIntersection(intersection.X(), intersection.Y()),
				),
				intersection.Y(),
			)
		} else {
			return []LineSegment3D{}, errors.New("Unable to achieve intersection current to next. Guess I'm the big stupid")
		}
	}

	return clean, nil
}
