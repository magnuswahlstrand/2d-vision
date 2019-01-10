package internal

import (
	"image"
	"math"

	geo "github.com/paulmach/go.geo"
)

func BasicRayCasting(x, y float64, objects []image.Rectangle) []Segment {
	lines := []Segment{}

	sparser := 4.
	for i := 0.; i < (360 / sparser); i++ {
		length := float64(1000) // Something big enough

		angle := sparser * math.Pi * i / 180
		start := geo.NewPoint(x, y)
		end := geo.NewPoint(x, y).Add(geo.NewPoint(length*math.Cos(angle), length*math.Sin(angle)))

		line := geo.NewLine(start, end)

		// Check intersection with all walls

		points := []*geo.Point{}
		for _, r := range objects {
			path := geoPathFromRect(r)
			tmp, _ := path.Intersection(line)
			points = append(points, tmp...)
		}

		// Find closest point
		min := math.Inf(1)
		minP := &geo.Point{}
		for i := range points {

			d := points[i].DistanceFrom(geo.NewPoint(x, y))

			if d < min {
				min = d
				minP = points[i]
			}
		}

		lines = append(lines, Seg(x, y, minP.X(), minP.Y()))
	}
	return lines
}

func geoPathFromRect(r image.Rectangle) *geo.Path {
	path := geo.NewPath()
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	return path
}
