package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"testing"

	"github.com/kyeett/2d-vision/internal"
	geo "github.com/paulmach/go.geo"
)

func Test_Intersect(t *testing.T) {
	path := geo.NewPath()
	path.Push(geo.NewPoint(0, 0))
	path.Push(geo.NewPoint(1, 1))

	line := geo.NewLine(geo.NewPoint(0, 1), geo.NewPoint(1, 0))

	// intersects does a simpler check for yes/no
	if path.Intersects(line) {
		// intersection will return the actual points and places on intersection
		points, segments := path.Intersection(line)

		for i, _ := range points {
			log.Printf("Intersection %d at %v with path segment %d", i, points[i], segments[i][0])
		}
	}
}

func Test_IntersectRect(t *testing.T) {
	// path := geo.NewPath()
	// path.Push(geo.NewPoint(0, 0))
	// path.Push(geo.NewPoint(1, 1))

	r := image.Rect(0, 0, 1, 1)
	path := internal.GeoPathFromRect(r)

	// geo.NewBoundFromPoints(r.Min, r.Max)

	line := geo.NewLine(geo.NewPoint(-0.1, 0.5), geo.NewPoint(1.1, 0.5))

	// intersects does a simpler check for yes/no
	if path.Intersects(line) {
		// intersection will return the actual points and places on intersection
		points, segments := path.Intersection(line)
		for i, _ := range points {
			log.Printf("Intersection %d at %v with path segment %d", i, points[i], segments[i][0])
		}
	}
}

func Test_Direction(t *testing.T) {
	// path := geo.NewPath()
	// path.Push(geo.NewPoint(0, 0))
	// path.Push(geo.NewPoint(1, 1))

	for _, s := range []internal.Segment{internal.Seg(0, 0, -1, -1), internal.Seg(0, 0, 1, 1), internal.Seg(0, 0, 0, 1)} {
		fmt.Println(180 * s.Direction() / math.Pi)
	}
	t.Fail()
}
