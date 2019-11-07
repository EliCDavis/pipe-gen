package main

import (
	"bufio"
	"os"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/vector"
)

func save(mesh mesh.Model, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	err = mesh.Save(w)
	if err != nil {
		return err
	}
	return w.Flush()
}

func main() {

	hopper := Hopper{
		binHeight:   2.0,
		taperHeight: 1.0,
		radius:      1.,
		position:    vector.Vector3Zero(),
		rotation:    mesh.QuaternionZero(),
	}

	mesh := hopper.ToModel()

	// points := make([]vector.Vector3, 0)
	// thickness := make([]float64, 0)
	// for i := 0.0; i < 1000.0; i += 1.0 {
	// 	points = append(points, vector.NewVector3(math.Sin(i/100.0)*10.0, i/10, 0))
	// 	thickness = append(thickness, math.Abs(math.Sin(i/100.0))+.1)
	// }

	// var ls LineSegment3D = points

	// mesh, err := ls.CreatePipeWithVarryingThickness(thickness)

	// if err != nil {
	// 	panic(err)
	// }

	err := save(mesh, "out.obj")

	if err != nil {
		panic(err)
	}

}
