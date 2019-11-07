package main

import (
	"github.com/EliCDavis/vector"
)

// Spout is something a pipe connects too.
type Spout struct {
	outDirection vector.Vector3
	entrance     vector.Vector3
	opening      float64
}
