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
	mesh, err := DrawPlaneShape(
		vector.Vector3Zero(),
		vector.Vector3Up(),
		30.0,
		80,
	)

	if err != nil {
		panic(err)
	}

	err = save(mesh, "out.obj")

	if err != nil {
		panic(err)
	}

}
