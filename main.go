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

	radius := 1.0

	model := mesh.Model{}

	// for x := 1.0; x < 5.0; x += 1.0 {
	// 	hopper := Hopper{
	// 		binHeight:       2.0,
	// 		taperHeight:     1.0,
	// 		radius:          radius,
	// 		heightOffGround: 1.0,
	// 		position:        vector.NewVector3(((radius*2.0)+1.0)*x, 0, 0),
	// 		rotation:        mesh.QuaternionZero(),
	// 	}
	// 	model = model.Merge(hopper.ToModel())
	// }

	hopper := Hopper{
		binHeight:       2.0,
		taperHeight:     1.0,
		radius:          radius,
		heightOffGround: 1.0,
		position:        vector.NewVector3(((radius * 2.0) + 1.0), 0, 0),
		rotation:        mesh.QuaternionZero(),
	}
	model = model.Merge(hopper.ToModel())

	hopper = Hopper{
		binHeight:       2.0,
		taperHeight:     2.0,
		radius:          radius,
		heightOffGround: 1.0,
		position:        vector.NewVector3(((radius*2.0)+1.0)*2.0, 0, 0),
		rotation:        mesh.QuaternionZero(),
	}
	model = model.Merge(hopper.ToModel())

	hopper = Hopper{
		binHeight:       3.0,
		taperHeight:     1.0,
		radius:          radius,
		heightOffGround: 1.0,
		position:        vector.NewVector3(((radius*2.0)+1.0)*3.0, 0, 0),
		rotation:        mesh.QuaternionZero(),
	}
	model = model.Merge(hopper.ToModel())

	hopper = Hopper{
		binHeight:       2.0,
		taperHeight:     1.0,
		radius:          radius,
		heightOffGround: 3.0,
		position:        vector.NewVector3(((radius*2.0)+1.0)*4.0, 0, 0),
		rotation:        mesh.QuaternionZero(),
	}
	model = model.Merge(hopper.ToModel())

	err := save(model, "out.obj")

	if err != nil {
		panic(err)
	}

}
