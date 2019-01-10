package internal

import (
	"image"
	"math"
	"sort"

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
			path := GeoPathFromRect(r)
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

func SmartRayCasting(x, y float64, objects []image.Rectangle) []Segment {
	center := geo.NewPoint(x, y)
	lines := []Segment{}

	rayLength := float64(1000) // Something "big" enough
	for _, rOuter := range objects {

		objPoints := geoPointsFromRect(rOuter)

		// Calculate the angle for all points
		for _, p := range objPoints {

			angle := geo.NewLine(center, p).Direction()

			// Offset angle in two directions, to go around corner
			for _, offset := range []float64{-0.005, 0.005} {
				points := []*geo.Point{}
				start := geo.NewPoint(x, y)
				end := geo.NewPoint(x, y).Add(geo.NewPoint(rayLength*math.Cos(angle+offset), rayLength*math.Sin(angle+offset)))
				ray := geo.NewLine(start, end)

				for _, rInner := range objects {
					path := GeoPathFromRect(rInner)
					tmp, _ := path.Intersection(ray)
					points = append(points, tmp...)
				}

				// // Find closest point
				min := math.Inf(1)
				minP := &geo.Point{}
				for i := range points {
					d := points[i].DistanceFrom(center)

					if d < min {
						min = d
						minP = points[i]
					}
				}
				lines = append(lines, Seg(x, y, minP.X(), minP.Y()))
			}
		}
	}

	sort.Slice(lines, func(i int, j int) bool {
		return lines[i].Direction() < lines[j].Direction()
	})
	return lines
}

func GeoPathFromRect(r image.Rectangle) *geo.Path {
	path := geo.NewPath()
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Min.Y)))
	path.Push(geo.NewPoint(float64(r.Max.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Max.Y)))
	path.Push(geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)))
	return path
}

func geoPointsFromRect(r image.Rectangle) []*geo.Point {
	return []*geo.Point{
		geo.NewPoint(float64(r.Min.X), float64(r.Min.Y)),
		geo.NewPoint(float64(r.Max.X), float64(r.Min.Y)),
		geo.NewPoint(float64(r.Max.X), float64(r.Max.Y)),
		geo.NewPoint(float64(r.Min.X), float64(r.Max.Y)),
	}
}
