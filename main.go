package main

import (
	"bufio"
	"math"
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

	points := make([]vector.Vector3, 0)
	for i := 0.0; i < 100.0; i += 1.0 {
		points = append(points, vector.NewVector3(math.Sin(i/10.0)*10.0, i, 0))
	}

	var ls LineSegment3D = points

	mesh, err := ls.CreatePipe(1)

	if err != nil {
		panic(err)
	}

	err = save(mesh, "out.obj")

	if err != nil {
		panic(err)
	}

}
