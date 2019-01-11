package internal

import (
	"image"

	"github.com/SolarLune/resolv/resolv"
	geo "github.com/paulmach/go.geo"
)

type Segment struct {
	X1, Y1, X2, Y2 float64
}

func Seg(x1, y1, x2, y2 float64) Segment {
	return Segment{x1, y1, x2, y2}
}

func (s Segment) Direction() float64 {
	return geo.NewLine(geo.NewPoint(s.X1, s.Y1), geo.NewPoint(s.X2, s.Y2)).Direction()

	// return math.Atan((s.X2 - s.X1) / (s.Y2 - s.Y1))
}

func SegmentsFromRect(r image.Rectangle) []Segment {
	s := []Segment{}
	s = append(s, Seg(float64(r.Min.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Min.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Max.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Max.Y), float64(r.Max.X), float64(r.Min.Y)))
	s = append(s, Seg(float64(r.Max.X), float64(r.Min.Y), float64(r.Min.X), float64(r.Min.Y)))
	return s
}

func LinesFromRect(r image.Rectangle) []resolv.Shape {
	return []resolv.Shape{
		resolv.NewLine(int32(r.Min.X), int32(r.Min.Y), int32(r.Min.X), int32(r.Max.Y)),
		resolv.NewLine(int32(r.Min.X), int32(r.Max.Y), int32(r.Max.X), int32(r.Max.Y)),
		resolv.NewLine(int32(r.Max.X), int32(r.Max.Y), int32(r.Max.X), int32(r.Min.Y)),
		resolv.NewLine(int32(r.Max.X), int32(r.Min.Y), int32(r.Min.X), int32(r.Min.Y)),
	}
}

func ShapeFromRect(r image.Rectangle) resolv.Shape {
	r2 := resolv.NewRectangle(int32(r.Min.X), int32(r.Min.Y), int32(r.Dx()), int32(r.Dy()))
	return r2
}
