package main

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/EliCDavis/mesh"
	"github.com/EliCDavis/meshedpotatoes/path"
	"github.com/EliCDavis/vector"
)

func save(mesh mesh.Model, name string) error {
	defer timeTrack(time.Now(), "Saving Model")
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {

	radius := 1.0

	model := mesh.Model{}

	for x := 1.0; x < 5.0; x += 1.0 {
		hopper := Hopper{
			binHeight:       2.0,
			taperHeight:     1.0,
			radius:          radius,
			heightOffGround: 1.0,
			position:        vector.NewVector3(((radius*2.0)+1.0)*x, 0, 0),
			rotation:        mesh.QuaternionZero(),
		}
		model = model.Merge(hopper.ToModel())
	}

	var bridgeDirection path.Path = []vector.Vector3{
		vector.NewVector3(1, 0, -2),
		vector.NewVector3(1, 0, 2),

		vector.NewVector3(5, 1, 2),
		vector.NewVector3(10, 1, 2),

		vector.NewVector3(14, 0, 2),
		vector.NewVector3(14, 0, -2),
	}

	bridge, err := CreateLoopingPlatform(bridgeDirection, 1.0)
	if err != nil {
		panic(err)
	}

	err = save(model.Merge(bridge), "out.obj")

	if err != nil {
		panic(err)
	}

}
