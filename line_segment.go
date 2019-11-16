package main

import (
	"errors"
	"math"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/meshedpotatoes"
	"github.com/EliCDavis/meshedpotatoes/path"
	"github.com/EliCDavis/vector"
)

// CreateLoopingPlatform creates a walkway that will close in on itself to create an
// enclosed track.
func CreateLoopingPlatform(pathOfPlatform path.Path, width float64) (mesh.Model, error) {
	if len(pathOfPlatform) < 3 {
		return mesh.Model{}, errors.New("Unable to create a looping platform with less than 3 points")
	}

	leftSidesDirty := make([]mesh.Line3D, 0)
	rightSidesDirty := make([]mesh.Line3D, 0)
	for i := 0; i < len(pathOfPlatform); i++ {
		previous := i - 1
		if previous < 0 {
			previous = len(pathOfPlatform) - 1
		}
		dir := pathOfPlatform[i].Sub(pathOfPlatform[previous])

		rightSide := dir.Cross(vector.Vector3Up()).Normalized().MultByConstant(width / 2.0)

		leftSidesDirty = append(leftSidesDirty, mesh.NewLine3D(
			pathOfPlatform[previous].Sub(rightSide),
			pathOfPlatform[i].Sub(rightSide),
		))

		rightSidesDirty = append(rightSidesDirty, mesh.NewLine3D(
			pathOfPlatform[previous].Add(rightSide),
			pathOfPlatform[i].Add(rightSide),
		))
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
		polys = append(polys, meshedpotatoes.MakeSquare(
			cleanLeftSide[i].GetStartPoint(),
			cleanLeftSide[i].GetEndPoint(),
			cleanRightSide[i].GetEndPoint(),
			cleanRightSide[i].GetStartPoint(),
		)...)

		leftRailing, _ := path.Path([]vector.Vector3{
			cleanLeftSide[i].GetStartPoint().Add(vector.Vector3Up()),
			cleanLeftSide[i].GetEndPoint().Add(vector.Vector3Up()),
		}).CreatePipe(.05)

		//.CreatePipe(.05)

		rightRailing, _ := path.Path([]vector.Vector3{
			cleanRightSide[i].GetStartPoint().Add(vector.Vector3Up()),
			cleanRightSide[i].GetEndPoint().Add(vector.Vector3Up()),
		}).CreatePipe(.05)

		railing = railing.
			Merge(leftRailing).
			Merge(rightRailing).
			Merge(leftRailing.Translate(vector.Vector3Up().MultByConstant(-.5))).
			Merge(rightRailing.Translate(vector.Vector3Up().MultByConstant(-.5))).
			Merge(railingPostsBetweenTwoPoints(cleanLeftSide[i].GetStartPoint(), cleanLeftSide[i].GetEndPoint())).
			Merge(railingPostsBetweenTwoPoints(cleanRightSide[i].GetStartPoint(), cleanRightSide[i].GetEndPoint()))
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
	pipe, err := path.Path([]vector.Vector3{
		position,
		position.Add(vector.Vector3Up()),
	}).CreatePipe(.05)
	if err != nil {
		panic(err)
	}
	return pipe
}

// Can be optimized
func properlyAlignSegments(dirty []mesh.Line3D) ([]mesh.Line3D, error) {

	clean := make([]mesh.Line3D, len(dirty))

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
		lastLine2D := mesh.NewLine2D(
			vector.NewVector2(scaledLastLine.GetStartPoint().X(), scaledLastLine.GetStartPoint().Z()),
			vector.NewVector2(scaledLastLine.GetEndPoint().X(), scaledLastLine.GetEndPoint().Z()),
		)

		scaledCurLine := dirty[i].ScaleOutwards(100000)
		curLine2D := mesh.NewLine2D(
			vector.NewVector2(scaledCurLine.GetStartPoint().X(), scaledCurLine.GetStartPoint().Z()),
			vector.NewVector2(scaledCurLine.GetEndPoint().X(), scaledCurLine.GetEndPoint().Z()),
		)

		scaledNextLine := dirty[next].ScaleOutwards(100000)
		nextLine2D := mesh.NewLine2D(
			vector.NewVector2(scaledNextLine.GetStartPoint().X(), scaledNextLine.GetStartPoint().Z()),
			vector.NewVector2(scaledNextLine.GetEndPoint().X(), scaledNextLine.GetEndPoint().Z()),
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

			clean[i] = clean[i].SetStartPoint(vector.NewVector3(
				dirty[previous].GetEndPoint().X(),
				math.Min(dirty[previous].GetEndPoint().Y(), dirty[i].GetStartPoint().Y()),
				dirty[previous].GetEndPoint().Z(),
			))

		} else if lastLine2D.Intersects(curLine2D) {
			intersection, _ := lastLine2D.Intersection(curLine2D)

			clean[i] = clean[i].SetStartPoint(vector.NewVector3(
				intersection.X(),
				math.Min(
					scaledLastLine.YIntersection(intersection.X(), intersection.Y()),
					scaledCurLine.YIntersection(intersection.X(), intersection.Y()),
				),
				intersection.Y(),
			))
		} else {
			return nil, errors.New("Unable to achieve intersection before to current. Guess I'm the big stupid")
		}

		if sameNextDirection.X() == 0 && sameNextDirection.Y() == 0 {

			clean[i] = clean[i].SetEndPoint(vector.NewVector3(
				dirty[i].GetEndPoint().X(),
				math.Min(dirty[i].GetEndPoint().Y(), dirty[next].GetStartPoint().Y()),
				dirty[i].GetEndPoint().Z(),
			))

		} else if curLine2D.Intersects(nextLine2D) {
			intersection, _ := curLine2D.Intersection(nextLine2D)

			clean[i] = clean[i].SetEndPoint(vector.NewVector3(
				intersection.X(),
				math.Min(
					scaledNextLine.YIntersection(intersection.X(), intersection.Y()),
					scaledCurLine.YIntersection(intersection.X(), intersection.Y()),
				),
				intersection.Y(),
			))
		} else {
			return nil, errors.New("Unable to achieve intersection current to next. Guess I'm the big stupid")
		}
	}

	return clean, nil
}
